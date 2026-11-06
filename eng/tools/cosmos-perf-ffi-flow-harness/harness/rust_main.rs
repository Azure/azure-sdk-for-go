// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

mod core;

use azure_core::{credentials::Secret, Bytes};
use azure_data_cosmos::{
    clients::ContainerClient, options::BinaryEncodingOptions, AccountEndpoint, AccountReference,
    CosmosClient, CosmosRuntime, RoutingStrategy,
};
use azure_data_cosmos_driver::{
    driver::{CosmosDriver, CosmosDriverRuntime},
    models::{ContainerReference, CosmosOperation, ItemReference},
    options::{
        ConnectionPoolOptions, DriverOptions, OperationOptions, OperationOptionsBuilder, Region,
    },
};
use serde_json::{json, Value};
use std::{
    env,
    error::Error,
    fs::{self, File, OpenOptions},
    io::{Read, Write},
    os::fd::FromRawFd,
    sync::{
        atomic::{AtomicU64, Ordering},
        Arc,
    },
    time::{Duration, Instant, SystemTime, UNIX_EPOCH},
};
use tokio::task::JoinSet;

type Result<T> = std::result::Result<T, Box<dyn Error + Send + Sync>>;

struct Options {
    endpoint: String,
    database: String,
    container: String,
    region: String,
    owner: String,
    leg: String,
    action: String,
    key_fd: i32,
    workers: usize,
    warmup: u64,
    measurement: u64,
}

#[derive(Clone)]
enum Client {
    Sdk(ContainerClient),
    Driver(Arc<CosmosDriver>, ContainerReference),
}

struct ReadResult {
    body: Bytes,
    retries: u64,
    gateway: bool,
    http2: bool,
}

fn binary_disabled() -> OperationOptions {
    OperationOptionsBuilder::new()
        .with_binary_encoding(BinaryEncodingOptions::new().with_enabled(false))
        .build()
}

impl Client {
    async fn read(&self, id: &str) -> Result<ReadResult> {
        match self {
            Self::Sdk(client) => {
                let response = client.read_item(id.to_owned(), id, None).await?;
                let diagnostics = response.diagnostics();
                let requests = diagnostics.requests();
                let last = requests
                    .last()
                    .ok_or("SDK response has no request diagnostics")?;
                let result = ReadResult {
                    gateway: last.transport_kind().is_gateway(),
                    http2: last.transport_http_version().is_http2(),
                    retries: requests.len().saturating_sub(1) as u64,
                    body: response.into_body().single()?,
                };
                Ok(result)
            }
            Self::Driver(driver, container) => {
                let reference =
                    ItemReference::from_name(container, id.to_owned().into(), id.to_owned());
                let response = driver
                    .execute_singleton_operation(
                        CosmosOperation::read_item(reference),
                        OperationOptions::default(),
                    )
                    .await?;
                let diagnostics = response.diagnostics();
                let requests = diagnostics.requests();
                let last = requests
                    .last()
                    .ok_or("driver response has no request diagnostics")?;
                let result = ReadResult {
                    gateway: last.transport_kind().is_gateway(),
                    http2: last.transport_http_version().is_http2(),
                    retries: requests.len().saturating_sub(1) as u64,
                    body: response.into_body().single()?,
                };
                Ok(result)
            }
        }
    }
}

fn options() -> Result<Options> {
    let args: Vec<String> = env::args().skip(1).collect();
    core::reject_secret_inputs(&args)?;
    let mut result = Options {
        endpoint: String::new(),
        database: String::new(),
        container: String::new(),
        region: String::new(),
        owner: String::new(),
        leg: String::new(),
        action: String::new(),
        key_fd: -1,
        workers: 0,
        warmup: 10,
        measurement: 60,
    };
    let mut index = 0;
    while index < args.len() {
        let argument = &args[index];
        if !argument.starts_with("--") {
            if !result.action.is_empty() {
                return Err("multiple actions".into());
            }
            result.action.clone_from(argument);
            index += 1;
            continue;
        }
        let value = args.get(index + 1).ok_or("missing flag value")?;
        match argument.as_str() {
            "--endpoint" => result.endpoint.clone_from(value),
            "--database" => result.database.clone_from(value),
            "--container" => result.container.clone_from(value),
            "--region" => result.region.clone_from(value),
            "--owner" => result.owner.clone_from(value),
            "--leg" => result.leg.clone_from(value),
            "--key-fd" => result.key_fd = value.parse()?,
            "--workers" => result.workers = value.parse()?,
            "--warmup-seconds" => result.warmup = value.parse()?,
            "--measurement-seconds" => result.measurement = value.parse()?,
            _ => return Err(format!("unknown flag {argument}").into()),
        }
        index += 2;
    }
    core::validate_owner(&result.owner)?;
    if result.key_fd < 3
        || result.endpoint.is_empty()
        || result.database.is_empty()
        || result.container.is_empty()
        || result.region.is_empty()
        || !matches!(result.leg.as_str(), "rust-sdk" | "rust-driver")
        || !matches!(result.action.as_str(), "preflight" | "run")
    {
        return Err("invalid or missing required arguments".into());
    }
    if result.action == "run" && (!matches!(result.workers, 32 | 64) || result.measurement == 0) {
        return Err("run requires workers=32|64 and positive measurement duration".into());
    }
    Ok(result)
}

async fn client(options: &Options) -> Result<Client> {
    // SAFETY: the controller transfers ownership of a private inherited descriptor >= 3.
    let key_file = unsafe { File::from_raw_fd(options.key_fd) };
    let mut key = Vec::new();
    key_file.take(16 * 1024).read_to_end(&mut key)?;
    while key.last().is_some_and(|byte| byte.is_ascii_whitespace()) {
        key.pop();
    }
    if key.is_empty() {
        return Err("empty authentication descriptor".into());
    }
    let pool = ConnectionPoolOptions::builder()
        .with_gateway_v2_disabled(true)
        .build()?;
    let result = match options.leg.as_str() {
        "rust-sdk" => {
            let runtime = CosmosRuntime::builder()
                .with_connection_pool(pool)
                .build()
                .await?;
            let endpoint: AccountEndpoint = options.endpoint.parse()?;
            let account = AccountReference::with_authentication_key(
                endpoint,
                Secret::from(String::from_utf8(key.clone())?),
            );
            let client = CosmosClient::builder()
                .with_runtime(runtime)
                .with_binary_encoding_options(BinaryEncodingOptions::new().with_enabled(false))
                .build(
                    account,
                    RoutingStrategy::PreferredRegions(vec![Region::new(options.region.clone())]),
                )
                .await?;
            Client::Sdk(
                client
                    .database_client(&options.database)
                    .container_client(&options.container, None)
                    .await?,
            )
        }
        "rust-driver" => {
            let runtime = CosmosDriverRuntime::builder()
                .with_connection_pool(pool)
                .build()
                .await?;
            let account = azure_data_cosmos_driver::models::AccountReference::with_account_key(
                options.endpoint.parse()?,
                Secret::from(String::from_utf8(key.clone())?),
            );
            let driver = runtime
                .create_driver(
                    DriverOptions::builder(account)
                        .with_operation_options(binary_disabled())
                        .with_preferred_regions(vec![Region::new(options.region.clone())])
                        .build(),
                )
                .await?;
            let container = driver
                .resolve_container(&options.database, &options.container, binary_disabled())
                .await?;
            Client::Driver(driver, container)
        }
        _ => unreachable!(),
    };
    key.fill(0);
    Ok(result)
}

async fn preflight(client: &Client, options: &Options) -> Result<Value> {
    let response = client.read(&core::item_id(&options.owner, 0)).await?;
    core::validate_document(&response.body, &options.owner, 0)?;
    if !response.gateway || !response.http2 || response.retries != 0 {
        return Err("preflight requires Gateway V1, HTTP/2, and zero retries".into());
    }
    Ok(json!({
        "transportKind":"gateway-v1", "negotiatedHTTPVersion":"HTTP/2", "binaryEncoding":false,
        "binaryDisablement":"explicit SDK default and direct-driver operation options",
        "wireProof":"separate pre-conversion proof required; caller JSON is not wire proof",
        "callerBodyBytes":response.body.len()
    }))
}

struct PhaseResult {
    success: u64,
    errors: u64,
    retries: u64,
    latencies: Vec<u64>,
    #[cfg(feature = "flow-diagnostic")]
    operation_buckets: std::collections::BTreeMap<u64, u64>,
}

async fn phase(
    client: &Client,
    options: &Options,
    seconds: u64,
    measured: bool,
) -> Result<PhaseResult> {
    let sequence = Arc::new(AtomicU64::new(0));
    let deadline = Instant::now() + Duration::from_secs(seconds);
    let mut tasks = JoinSet::new();
    for _ in 0..options.workers {
        let client = client.clone();
        let sequence = sequence.clone();
        let owner = options.owner.clone();
        tasks.spawn(async move {
            let mut result = PhaseResult {
                success: 0,
                errors: 0,
                retries: 0,
                latencies: Vec::new(),
                #[cfg(feature = "flow-diagnostic")]
                operation_buckets: std::collections::BTreeMap::new(),
            };
            while Instant::now() < deadline {
                let index = core::key_index(sequence.fetch_add(1, Ordering::Relaxed));
                let id = core::item_id(&owner, index);
                let start = Instant::now();
                match client.read(&id).await {
                    Ok(response) => {
                        if measured {
                            result.latencies.push(start.elapsed().as_nanos() as u64);
                        }
                        result.retries += response.retries;
                        if response.body.is_empty() || !response.gateway || !response.http2 {
                            result.errors += 1;
                            eprintln!(
                                "invalid read response: bytes={} gateway={} http2={}",
                                response.body.len(),
                                response.gateway,
                                response.http2
                            );
                            break;
                        }
                        result.success += 1;
                        #[cfg(feature = "flow-diagnostic")]
                        if measured {
                            let second = SystemTime::now()
                                .duration_since(UNIX_EPOCH)
                                .expect("clock before Unix epoch")
                                .as_secs();
                            *result.operation_buckets.entry(second).or_default() += 1;
                        }
                    }
                    Err(error) => {
                        result.errors += 1;
                        if let Some(error) =
                            error.downcast_ref::<azure_data_cosmos_driver::error::CosmosError>()
                        {
                            eprintln!("point read failed: status={}", error.status());
                        } else {
                            eprintln!("point read failed outside Cosmos status handling");
                        }
                        break;
                    }
                }
            }
            result
        });
    }
    let mut result = PhaseResult {
        success: 0,
        errors: 0,
        retries: 0,
        latencies: Vec::new(),
        #[cfg(feature = "flow-diagnostic")]
        operation_buckets: std::collections::BTreeMap::new(),
    };
    while let Some(joined) = tasks.join_next().await {
        let mut worker = joined?;
        result.success += worker.success;
        result.errors += worker.errors;
        result.retries += worker.retries;
        result.latencies.append(&mut worker.latencies);
        #[cfg(feature = "flow-diagnostic")]
        for (second, count) in worker.operation_buckets {
            *result.operation_buckets.entry(second).or_default() += count;
        }
    }
    Ok(result)
}

fn cpu_seconds() -> Result<f64> {
    let stat = fs::read_to_string("/proc/self/stat")?;
    let end = stat.rfind(')').ok_or("invalid process stat")?;
    let fields: Vec<&str> = stat[end + 2..].split_whitespace().collect();
    let ticks: f64 = env::var("COSMOS_PERF_CLK_TCK")?.parse()?;
    if ticks <= 0.0 {
        return Err("invalid CLK_TCK".into());
    }
    Ok((fields[11].parse::<f64>()? + fields[12].parse::<f64>()?) / ticks)
}

fn percentiles(mut values: Vec<u64>) -> Value {
    values.sort_unstable();
    let get = |q: f64| {
        if values.is_empty() {
            0.0
        } else {
            values[((values.len() - 1) as f64 * q) as usize] as f64 / 1e6
        }
    };
    json!({"p50":get(0.5),"p90":get(0.9),"p95":get(0.95),"p99":get(0.99),"p999":get(0.999)})
}

fn barrier(name: &str) -> Result<()> {
    let Some(directory) = env::var_os("COSMOS_PERF_BARRIER_DIR") else {
        return Ok(());
    };
    let directory = std::path::Path::new(&directory);
    OpenOptions::new()
        .create_new(true)
        .write(true)
        .open(directory.join(format!("request-{name}")))?;
    let deadline = Instant::now() + Duration::from_secs(30);
    loop {
        match fs::metadata(directory.join(format!("ready-{name}"))) {
            Ok(_) => return Ok(()),
            Err(error) if error.kind() == std::io::ErrorKind::NotFound => {}
            Err(error) => return Err(error.into()),
        }
        if Instant::now() > deadline {
            return Err("measurement barrier timed out".into());
        }
        std::thread::sleep(Duration::from_millis(1));
    }
}

fn monotonic_ns() -> Result<u64> {
    #[repr(C)]
    struct Timespec {
        seconds: std::os::raw::c_long,
        nanoseconds: std::os::raw::c_long,
    }
    unsafe extern "C" {
        fn clock_gettime(clock: i32, time: *mut Timespec) -> i32;
    }
    let mut time = Timespec {
        seconds: 0,
        nanoseconds: 0,
    };
    // SAFETY: Linux CLOCK_MONOTONIC writes one correctly aligned timespec.
    if unsafe { clock_gettime(1, &mut time) } != 0 {
        return Err(std::io::Error::last_os_error().into());
    }
    Ok(time.seconds as u64 * 1_000_000_000 + time.nanoseconds as u64)
}

fn bounds(start: u64, end: u64) -> Result<()> {
    if let Some(directory) = env::var_os("COSMOS_PERF_BARRIER_DIR") {
        let mut file = OpenOptions::new()
            .create_new(true)
            .write(true)
            .open(std::path::Path::new(&directory).join("measurement-bounds.json"))?;
        serde_json::to_writer(
            &mut file,
            &json!({"monotonicStartNs":start,"monotonicEndNs":end}),
        )?;
        file.write_all(b"\n")?;
    }
    Ok(())
}

#[tokio::main(flavor = "multi_thread", worker_threads = 2)]
async fn main() -> Result<()> {
    let options = options()?;
    let client = client(&options).await?;
    let proof = preflight(&client, &options).await?;
    if options.action == "preflight" {
        println!("{proof}");
        #[cfg(feature = "wire-proof")]
        azure_data_cosmos_driver::perf_wire::write_requested()?;
        return Ok(());
    }
    let warmup = phase(&client, &options, options.warmup, false).await?;
    let warmup_counts = json!({"success":warmup.success,"error":warmup.errors,
                              "retry":warmup.retries,"throttle":null,"excludedFromScores":true});
    eprintln!("WARMUP_COUNTS {warmup_counts}");
    if let Ok(path) = env::var("COSMOS_PERF_ADAPTER_RESULT") {
        let path = std::path::Path::new(&path)
            .parent()
            .ok_or("invalid result directory")?
            .join("warmup-counters.json");
        let mut file = OpenOptions::new().write(true).create_new(true).open(path)?;
        serde_json::to_writer_pretty(&mut file, &warmup_counts)?;
        file.write_all(b"\n")?;
    }
    if warmup.errors != 0 {
        return Err("warmup completed with failed operations".into());
    }
    barrier("start")?;
    let cpu = cpu_seconds()?;
    let monotonic_start = monotonic_ns()?;
    let start_unix_ns = SystemTime::now().duration_since(UNIX_EPOCH)?.as_nanos();
    let start = Instant::now();
    let result = phase(&client, &options, options.measurement, true).await?;
    let wall = start.elapsed().as_secs_f64();
    let end_unix_ns = SystemTime::now().duration_since(UNIX_EPOCH)?.as_nanos();
    let monotonic_end = monotonic_ns()?;
    let cpu = cpu_seconds()? - cpu;
    bounds(monotonic_start, monotonic_end)?;
    barrier("end")?;
    #[cfg(feature = "flow-diagnostic")]
    let operation_buckets = Some(&result.operation_buckets);
    #[cfg(not(feature = "flow-diagnostic"))]
    let operation_buckets: Option<&std::collections::BTreeMap<u64, u64>> = None;
    let value = json!({
        "workload":{
            "partitionKeyPath":"/pk","seedCount":1024,"documentBytes":1024,
            "operation":"point-read","workers":options.workers,"warmupSeconds":options.warmup,
            "measurementSeconds":options.measurement,"closedLoop":true,
            "keySequence":"round-robin-0-through-1023","rawResponseMode":true,
            "validationOutsideTimedPath":true,"gatewayProtocol":"GatewayV1HTTP2",
            "binaryEncoding":false,"authentication":"account-key","tokioWorkerThreads":2
        },
        "transportProof":proof,
        "counts":{"success":result.success,"error":result.errors,"retry":result.retries,"throttle":0},
        "timing":{"wallSeconds":wall,"processCPUSeconds":cpu,"measurementStartUnixNs":start_unix_ns,
                  "measurementEndUnixNs":end_unix_ns},
        "latencyMilliseconds":percentiles(result.latencies),
        "queueDelayMilliseconds":null,
        "gauges":{"queueHighWater":null,"inFlightHighWater":null,"completionHighWater":null},
        "leg":options.leg,"operationBuckets":operation_buckets,"warmupCounts":warmup_counts
    });
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(env::var("COSMOS_PERF_ADAPTER_RESULT")?)?;
    serde_json::to_writer_pretty(&mut file, &value)?;
    file.write_all(b"\n")?;
    if result.errors != 0 || result.retries != 0 {
        return Err("measurement invalid: errors/retries".into());
    }
    Ok(())
}
