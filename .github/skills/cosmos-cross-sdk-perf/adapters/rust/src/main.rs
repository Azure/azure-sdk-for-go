// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// cSpell:words rfind serde usize

mod core;

use azure_core::credentials::Secret;
use azure_data_cosmos::clients::ContainerClient;
use azure_data_cosmos::error::CosmosError;
use azure_data_cosmos::{
    options::{BinaryEncodingOptions, ConnectionPoolOptions, Region},
    AccountEndpoint, AccountReference, CosmosClient, CosmosRuntime, RoutingStrategy,
};
use core::{
    item_id, key_index, make_document, owned_by, validate_document, validate_owner, SEED_COUNT,
};
use serde::Serialize;
use serde_json::{json, Value};
use std::{
    env,
    error::Error,
    fs::{self, File},
    io::{Read, Write},
    os::fd::FromRawFd,
    sync::{
        atomic::{AtomicU64, Ordering},
        Arc,
    },
    time::{Duration, Instant},
};
use tokio::task::JoinSet;

type BoxError = Box<dyn Error + Send + Sync>;

#[derive(Clone)]
struct Options {
    endpoint: String,
    database: String,
    container: String,
    region: String,
    owner: String,
    key_fd: i32,
    workers: usize,
    warmup: Duration,
    measurement: Duration,
}

#[derive(Serialize)]
#[serde(rename_all = "camelCase")]
struct AdapterResult {
    workload: Value,
    transport_proof: Value,
    counts: Value,
    timing: Value,
    latency_milliseconds: Value,
    queue_delay_milliseconds: Option<Value>,
    gauges: Value,
}

#[tokio::main(flavor = "multi_thread")]
async fn main() -> Result<(), BoxError> {
    let arguments: Vec<String> = env::args().skip(1).collect();
    core::reject_secret_inputs(&arguments)?;
    let (options, action) = parse_options(&arguments)?;
    let mut key = read_key(options.key_fd)?;
    let client = create_client(&options, &key).await?;
    key.fill(0);
    let container = client
        .database_client(&options.database)
        .container_client(&options.container, None)
        .await?;

    match action.as_str() {
        "seed" => seed(&container, &options).await?,
        "verify" => verify(&container, &options).await?,
        "preflight" => println!("{}", preflight(&container, &options).await?),
        "run" => run_benchmark(container, &options).await?,
        "cleanup" => cleanup(&container, &options).await?,
        _ => return Err(format!("unknown action {action:?}").into()),
    }
    Ok(())
}

fn parse_options(arguments: &[String]) -> Result<(Options, String), BoxError> {
    let mut endpoint = String::new();
    let mut database = String::new();
    let mut container = String::new();
    let mut region = String::new();
    let mut owner = String::new();
    let mut key_fd = env::var("COSMOS_PERF_KEY_FD")
        .ok()
        .and_then(|v| v.parse().ok());
    let mut workers = env::var("COSMOS_PERF_WORKERS")
        .ok()
        .and_then(|v| v.parse().ok())
        .unwrap_or(0);
    let mut warmup = 30;
    let mut measurement = 300;
    let mut action = None;
    let mut index = 0;
    while index < arguments.len() {
        let argument = &arguments[index];
        if !argument.starts_with("--") {
            if action.replace(argument.clone()).is_some() {
                return Err("exactly one action is required".into());
            }
            index += 1;
            continue;
        }
        let value = arguments.get(index + 1).ok_or("flag value is missing")?;
        match argument.as_str() {
            "--endpoint" => endpoint = value.clone(),
            "--database" => database = value.clone(),
            "--container" => container = value.clone(),
            "--region" => region = value.clone(),
            "--owner" => owner = value.clone(),
            "--key-fd" => key_fd = Some(value.parse()?),
            "--workers" => workers = value.parse()?,
            "--warmup-seconds" => warmup = value.parse()?,
            "--measurement-seconds" => measurement = value.parse()?,
            _ => return Err(format!("unknown flag {argument:?}").into()),
        }
        index += 2;
    }
    validate_owner(&owner)?;
    if endpoint.is_empty() || database.is_empty() || container.is_empty() || region.is_empty() {
        return Err("--endpoint, --database, --container, and --region are required".into());
    }
    let key_fd = key_fd.ok_or("--key-fd is required")?;
    if key_fd < 3 {
        return Err("--key-fd must be at least 3".into());
    }
    Ok((
        Options {
            endpoint,
            database,
            container,
            region,
            owner,
            key_fd,
            workers,
            warmup: Duration::from_secs(warmup),
            measurement: Duration::from_secs(measurement),
        },
        action.ok_or("one action is required")?,
    ))
}

fn read_key(fd: i32) -> Result<Vec<u8>, BoxError> {
    let mut file = unsafe { File::from_raw_fd(fd) };
    let mut value = Vec::new();
    file.take(16 * 1024).read_to_end(&mut value)?;
    while value.last().is_some_and(|byte| byte.is_ascii_whitespace()) {
        value.pop();
    }
    if value.is_empty() {
        return Err("key descriptor was empty".into());
    }
    Ok(value)
}

async fn create_client(options: &Options, key: &[u8]) -> Result<CosmosClient, BoxError> {
    let pool = ConnectionPoolOptions::builder()
        .with_gateway_v2_disabled(true)
        .build()?;
    let runtime = CosmosRuntime::builder()
        .with_connection_pool(pool)
        .build()
        .await?;
    let endpoint: AccountEndpoint = options.endpoint.parse()?;
    let secret = Secret::from(String::from_utf8(key.to_vec())?);
    let account = AccountReference::with_authentication_key(endpoint, secret);
    Ok(CosmosClient::builder()
        .with_runtime(runtime)
        .with_binary_encoding_options(BinaryEncodingOptions::new().with_enabled(false))
        .build(
            account,
            RoutingStrategy::ProximityTo(Region::new(options.region.clone())),
        )
        .await?)
}

async fn read_raw(
    container: &ContainerClient,
    owner: &str,
    index: usize,
) -> Result<(Vec<u8>, u64, Value), BoxError> {
    let id = item_id(owner, index);
    let response = container.read_item(id.as_str(), id.as_str(), None).await?;
    let diagnostics = response.diagnostics();
    let requests = diagnostics.requests();
    let retries = requests.len().saturating_sub(1) as u64;
    let proof = requests.last().map(|request| {
        json!({
            "transportKind": if request.transport_kind().is_gateway() { "gateway-v1" } else { "gateway-v2" },
            "negotiatedHTTPVersion": if request.transport_http_version().is_http2() { "HTTP/2" } else { "HTTP/1.1" },
            "binaryEncoding": false
        })
    }).unwrap_or(Value::Null);
    let body = response.into_body().single()?;
    Ok((body.to_vec(), retries, proof))
}

async fn seed(container: &ContainerClient, options: &Options) -> Result<(), BoxError> {
    for index in 0..SEED_COUNT {
        let id = item_id(&options.owner, index);
        let body = make_document(&options.owner, index)?;
        let value: Value = serde_json::from_slice(&body)?;
        match read_raw(container, &options.owner, index).await {
            Ok((existing, _, _)) if owned_by(&existing, &options.owner) => {
                container
                    .upsert_item(id.as_str(), id.as_str(), value, None)
                    .await?;
            }
            Ok(_) => return Err(format!("refusing to overwrite non-owned item {id:?}").into()),
            Err(error) if is_not_found(error.as_ref()) => {
                container
                    .create_item(id.as_str(), id.as_str(), value, None)
                    .await?;
            }
            Err(error) => return Err(error),
        }
    }
    println!(
        "{}",
        json!({"action":"seed","owner":options.owner,"count":SEED_COUNT})
    );
    Ok(())
}

async fn verify(container: &ContainerClient, options: &Options) -> Result<(), BoxError> {
    for index in 0..SEED_COUNT {
        let (body, _, _) = read_raw(container, &options.owner, index).await?;
        validate_document(&body, &options.owner, index)?;
    }
    println!(
        "{}",
        json!({"action":"verify","owner":options.owner,"count":SEED_COUNT})
    );
    Ok(())
}

async fn cleanup(container: &ContainerClient, options: &Options) -> Result<(), BoxError> {
    let mut deleted = 0;
    let mut skipped = 0;
    for index in 0..SEED_COUNT {
        let id = item_id(&options.owner, index);
        match read_raw(container, &options.owner, index).await {
            Ok((body, _, _)) if owned_by(&body, &options.owner) => {
                container
                    .delete_item(id.as_str(), id.as_str(), None)
                    .await?;
                deleted += 1;
            }
            Ok(_) => skipped += 1,
            Err(error) if is_not_found(error.as_ref()) => {}
            Err(error) => return Err(error),
        }
    }
    println!(
        "{}",
        json!({"action":"cleanup","owner":options.owner,"deleted":deleted,"nonOwnedSkipped":skipped})
    );
    Ok(())
}

async fn preflight(container: &ContainerClient, options: &Options) -> Result<Value, BoxError> {
    let (body, _, proof) = read_raw(container, &options.owner, 0).await?;
    validate_document(&body, &options.owner, 0)?;
    if proof["transportKind"] != "gateway-v1" || proof["negotiatedHTTPVersion"] != "HTTP/2" {
        return Err(format!("transport preflight failed: {proof}").into());
    }
    Ok(proof)
}

async fn run_benchmark(container: ContainerClient, options: &Options) -> Result<(), BoxError> {
    if !matches!(options.workers, 8 | 16) {
        return Err("--workers must be 8 or 16".into());
    }
    let proof = preflight(&container, options).await?;
    run_phase(container.clone(), options, options.warmup, false).await?;
    let cpu_start = process_cpu_seconds()?;
    let started = Instant::now();
    let measured = run_phase(container, options, options.measurement, true).await?;
    let wall = started.elapsed().as_secs_f64();
    let cpu = process_cpu_seconds()? - cpu_start;
    if measured.errors != 0 {
        return Err(format!("measurement completed with {} errors", measured.errors).into());
    }
    let result = AdapterResult {
        workload: json!({
            "partitionKeyPath":"/pk","seedCount":SEED_COUNT,"documentBytes":core::DOCUMENT_BYTES,
            "operation":"point-read","workers":options.workers,
            "warmupSeconds":options.warmup.as_secs(),"measurementSeconds":options.measurement.as_secs(),
            "closedLoop":true,"keySequence":"round-robin-0-through-1023",
            "rawResponseMode":true,"validationOutsideTimedPath":true,
            "gatewayProtocol":"GatewayV1HTTP2","binaryEncoding":false,"authentication":"account-key"
        }),
        transport_proof: proof,
        counts: json!({"success":measured.success,"error":measured.errors,"retry":measured.retries,"throttle":0}),
        timing: json!({"wallSeconds":wall,"processCPUSeconds":cpu}),
        latency_milliseconds: percentiles(measured.latencies),
        queue_delay_milliseconds: None,
        gauges: json!({"queueHighWater":null,"inFlightHighWater":options.workers,"completionHighWater":null}),
    };
    let output = env::var("COSMOS_PERF_ADAPTER_RESULT")
        .map_err(|_| "COSMOS_PERF_ADAPTER_RESULT is required for run")?;
    let mut file = fs::OpenOptions::new()
        .create_new(true)
        .write(true)
        .open(output)?;
    serde_json::to_writer_pretty(&mut file, &result)?;
    file.write_all(b"\n")?;
    Ok(())
}

struct PhaseResult {
    success: u64,
    errors: u64,
    retries: u64,
    latencies: Vec<u64>,
}

async fn run_phase(
    container: ContainerClient,
    options: &Options,
    duration: Duration,
    record: bool,
) -> Result<PhaseResult, BoxError> {
    let sequence = Arc::new(AtomicU64::new(0));
    let deadline = Instant::now() + duration;
    let mut tasks = JoinSet::new();
    for _ in 0..options.workers {
        let container = container.clone();
        let owner = options.owner.clone();
        let sequence = sequence.clone();
        tasks.spawn(async move {
            let mut latencies = Vec::with_capacity(4096);
            let mut success = 0;
            let mut retries = 0;
            while Instant::now() < deadline {
                let index = key_index(sequence.fetch_add(1, Ordering::Relaxed));
                let started = Instant::now();
                let (body, operation_retries, _) = read_raw(&container, &owner, index).await?;
                let elapsed = started.elapsed();
                if record {
                    latencies.push(elapsed.as_nanos() as u64);
                }
                validate_document(&body, &owner, index)?;
                retries += operation_retries;
                success += 1;
            }
            Ok::<_, BoxError>((success, retries, latencies))
        });
    }
    let mut result = PhaseResult {
        success: 0,
        errors: 0,
        retries: 0,
        latencies: Vec::new(),
    };
    while let Some(joined) = tasks.join_next().await {
        let (success, retries, mut latencies) = joined??;
        result.success += success;
        result.retries += retries;
        result.latencies.append(&mut latencies);
    }
    Ok(result)
}

fn percentiles(mut values: Vec<u64>) -> Value {
    values.sort_unstable();
    let get = |quantile: f64| -> f64 {
        if values.is_empty() {
            return 0.0;
        }
        let index = ((values.len() - 1) as f64 * quantile) as usize;
        values[index] as f64 / 1_000_000.0
    };
    json!({"p50":get(0.5),"p90":get(0.9),"p95":get(0.95),"p99":get(0.99),"p999":get(0.999)})
}

fn is_not_found(error: &(dyn Error + 'static)) -> bool {
    error
        .downcast_ref::<CosmosError>()
        .is_some_and(|error| error.status().is_not_found())
}

fn process_cpu_seconds() -> Result<f64, BoxError> {
    let stat = fs::read_to_string("/proc/self/stat")?;
    let close = stat.rfind(')').ok_or("invalid /proc/self/stat")?;
    let fields: Vec<&str> = stat[close + 2..].split_whitespace().collect();
    let user: f64 = fields.get(11).ok_or("missing user CPU")?.parse()?;
    let system: f64 = fields.get(12).ok_or("missing system CPU")?.parse()?;
    Ok((user + system) / 100.0)
}
