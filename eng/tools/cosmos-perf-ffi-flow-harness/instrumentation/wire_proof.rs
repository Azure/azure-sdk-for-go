// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

use std::sync::atomic::{AtomicU64, Ordering};

static ELIGIBLE: AtomicU64 = AtomicU64::new(0);
static BINARY: AtomicU64 = AtomicU64::new(0);
static TEXT: AtomicU64 = AtomicU64::new(0);
static INVALID: AtomicU64 = AtomicU64::new(0);
static TEXT_OUTPUT_REQUESTED: AtomicU64 = AtomicU64::new(0);
static GATEWAY_V1_HTTP2: AtomicU64 = AtomicU64::new(0);
static MIN_BYTES: AtomicU64 = AtomicU64::new(u64::MAX);
static MAX_BYTES: AtomicU64 = AtomicU64::new(0);

pub(crate) fn point_response(
    response: &crate::models::CosmosResponse,
    text_output_requested: bool,
) {
    let crate::models::ResponseBody::Bytes(body) = response.body() else {
        INVALID.fetch_add(1, Ordering::Relaxed);
        return;
    };
    ELIGIBLE.fetch_add(1, Ordering::Relaxed);
    if response
        .diagnostics_ref()
        .requests()
        .last()
        .is_some_and(|request| {
            request.transport_kind().is_gateway() && request.transport_http_version().is_http2()
        })
    {
        GATEWAY_V1_HTTP2.fetch_add(1, Ordering::Relaxed);
    }
    MIN_BYTES.fetch_min(body.len() as u64, Ordering::Relaxed);
    MAX_BYTES.fetch_max(body.len() as u64, Ordering::Relaxed);
    if text_output_requested {
        TEXT_OUTPUT_REQUESTED.fetch_add(1, Ordering::Relaxed);
    }
    if crate::binary_json::is_binary(body) {
        BINARY.fetch_add(1, Ordering::Relaxed);
    } else if serde_json::from_slice::<serde_json::Value>(body).is_ok() {
        TEXT.fetch_add(1, Ordering::Relaxed);
    } else {
        INVALID.fetch_add(1, Ordering::Relaxed);
    }
}

/// Writes representation counters only; no response bytes or headers are persisted.
pub fn write_requested() -> std::io::Result<()> {
    let Some(path) = std::env::var_os("COSMOS_PERF_WIRE_PROOF") else {
        return Ok(());
    };
    let eligible = ELIGIBLE.load(Ordering::Relaxed);
    let value = serde_json::json!({
        "schemaVersion": 1,
        "hook": "successful document read, before driver-side transcoding",
        "eligibleResponses": eligible,
        "binaryResponsesBeforeConversion": BINARY.load(Ordering::Relaxed),
        "textResponsesBeforeConversion": TEXT.load(Ordering::Relaxed),
        "invalidResponsesBeforeConversion": INVALID.load(Ordering::Relaxed),
        "textOutputRequests": TEXT_OUTPUT_REQUESTED.load(Ordering::Relaxed),
        "gatewayV1HTTP2Responses": GATEWAY_V1_HTTP2.load(Ordering::Relaxed),
        "minimumWireBodyBytes": if eligible == 0 { None } else { Some(MIN_BYTES.load(Ordering::Relaxed)) },
        "maximumWireBodyBytes": if eligible == 0 { None } else { Some(MAX_BYTES.load(Ordering::Relaxed)) },
    });
    let mut file = std::fs::OpenOptions::new()
        .create_new(true)
        .write(true)
        .open(path)?;
    serde_json::to_writer_pretty(&mut file, &value)?;
    std::io::Write::write_all(&mut file, b"\n")?;
    Ok(())
}
