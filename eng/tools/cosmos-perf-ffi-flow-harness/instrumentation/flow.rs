// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

use std::sync::atomic::{AtomicU64, AtomicUsize, Ordering};
use std::sync::OnceLock;
use std::time::{SystemTime, UNIX_EPOCH};

const SECONDS: usize = 1024;
const EVENTS: usize = 8;
const NAMES: [&str; EVENTS] = [
    "admitted",
    "admission_aborted",
    "first_poll",
    "terminal",
    "enqueued",
    "enqueue_rejected",
    "drained",
    "task_submitted",
];

struct Bucket {
    counts: [AtomicU64; EVENTS],
    ready_high: AtomicU64,
    outstanding_high: AtomicU64,
    active_high: AtomicU64,
    ready_last: AtomicU64,
    outstanding_last: AtomicU64,
}

impl Bucket {
    const fn new() -> Self {
        Self {
            counts: [const { AtomicU64::new(0) }; EVENTS],
            ready_high: AtomicU64::new(0),
            outstanding_high: AtomicU64::new(0),
            active_high: AtomicU64::new(0),
            ready_last: AtomicU64::new(0),
            outstanding_last: AtomicU64::new(0),
        }
    }
}

static BUCKETS: [[Bucket; SECONDS]; 3] = [const { [const { Bucket::new() }; SECONDS] }; 3];
static EPOCH: OnceLock<u64> = OnceLock::new();
static PHASE: AtomicUsize = AtomicUsize::new(0);
static ACTIVE: AtomicU64 = AtomicU64::new(0);
static READY: AtomicU64 = AtomicU64::new(0);
static OUTSTANDING: AtomicU64 = AtomicU64::new(0);
static OVERFLOW: AtomicU64 = AtomicU64::new(0);
static START_NS: [AtomicU64; 3] = [const { AtomicU64::new(0) }; 3];
static END_NS: [AtomicU64; 3] = [const { AtomicU64::new(0) }; 3];

fn now_ns() -> u64 {
    SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .expect("system clock before Unix epoch")
        .as_nanos()
        .try_into()
        .expect("system clock exceeds u64 nanoseconds")
}

fn bucket() -> Option<&'static Bucket> {
    let second = now_ns() / 1_000_000_000;
    let epoch = *EPOCH.get_or_init(|| second);
    let Some(offset) = second.checked_sub(epoch) else {
        OVERFLOW.fetch_add(1, Ordering::Relaxed);
        return None;
    };
    if offset >= SECONDS as u64 {
        OVERFLOW.fetch_add(1, Ordering::Relaxed);
        return None;
    }
    Some(&BUCKETS[PHASE.load(Ordering::Acquire)][offset as usize])
}

pub(crate) fn event(index: usize, amount: u64) {
    if let Some(bucket) = bucket() {
        bucket.counts[index].fetch_add(amount, Ordering::Relaxed);
    }
}

pub(crate) fn queue(ready: usize, outstanding: u32) {
    READY.store(ready as u64, Ordering::Relaxed);
    OUTSTANDING.store(outstanding as u64, Ordering::Relaxed);
    if let Some(bucket) = bucket() {
        bucket.ready_high.fetch_max(ready as u64, Ordering::Relaxed);
        bucket
            .outstanding_high
            .fetch_max(outstanding as u64, Ordering::Relaxed);
        bucket.ready_last.store(ready as u64, Ordering::Relaxed);
        bucket
            .outstanding_last
            .store(outstanding as u64, Ordering::Relaxed);
    }
}

pub(crate) fn first_poll() {
    let current = ACTIVE.fetch_add(1, Ordering::Relaxed) + 1;
    if let Some(bucket) = bucket() {
        bucket.counts[2].fetch_add(1, Ordering::Relaxed);
        bucket.active_high.fetch_max(current, Ordering::Relaxed);
    }
}

pub(crate) fn terminal() {
    let previous = ACTIVE.fetch_sub(1, Ordering::Relaxed);
    assert!(previous > 0, "terminal without first poll");
    event(3, 1);
}

#[no_mangle]
pub extern "C" fn cosmos_perf_flow_phase(phase: u32) -> u64 {
    assert!(phase < 3, "invalid flow phase");
    let now = now_ns();
    EPOCH.get_or_init(|| now / 1_000_000_000);
    let previous = PHASE.swap(phase as usize, Ordering::AcqRel);
    END_NS[previous].store(now, Ordering::Release);
    START_NS[phase as usize].store(now, Ordering::Release);
    now
}

#[no_mangle]
pub extern "C" fn cosmos_perf_flow_epoch_seconds() -> u64 {
    *EPOCH.get_or_init(|| now_ns() / 1_000_000_000)
}

pub(crate) fn snapshot() -> serde_json::Value {
    let mut rows = Vec::new();
    for (phase, buckets) in BUCKETS.iter().enumerate() {
        for (second, bucket) in buckets.iter().enumerate() {
            let counts: Vec<u64> = bucket
                .counts
                .iter()
                .map(|value| value.load(Ordering::Relaxed))
                .collect();
            if counts.iter().all(|value| *value == 0) {
                continue;
            }
            rows.push(serde_json::json!({
                "phase":phase, "unixSecond":cosmos_perf_flow_epoch_seconds() + second as u64,
                "counts": NAMES.iter().zip(counts).collect::<std::collections::BTreeMap<_, _>>(),
                "readyHigh":bucket.ready_high.load(Ordering::Relaxed),
                "admittedUndrainedHigh":bucket.outstanding_high.load(Ordering::Relaxed),
                "polledUnfinishedHigh":bucket.active_high.load(Ordering::Relaxed),
                "readyLast":bucket.ready_last.load(Ordering::Relaxed),
                "admittedUndrainedLast":bucket.outstanding_last.load(Ordering::Relaxed),
            }));
        }
    }
    serde_json::json!({
        "schemaVersion":1, "rows":rows, "overflow":OVERFLOW.load(Ordering::Relaxed),
        "activeAtSnapshot":ACTIVE.load(Ordering::Relaxed),
        "readyAtSnapshot":READY.load(Ordering::Relaxed),
        "admittedUndrainedAtSnapshot":OUTSTANDING.load(Ordering::Relaxed),
        "phaseStartUnixNs":START_NS.iter().map(|v|v.load(Ordering::Acquire)).collect::<Vec<_>>(),
        "phaseEndUnixNs":END_NS.iter().map(|v|v.load(Ordering::Acquire)).collect::<Vec<_>>(),
        "definitions":{
            "readyHigh":"actual ready deque length sampled under the existing queue mutex after enqueue and drain",
            "admittedUndrainedHigh":"reserved native operations, including ready completions not yet drained",
            "polledUnfinishedHigh":"spawned oneshot tasks first-polled but not terminal; includes wrapper processing",
            "lastGauges":"last queue-lock observation in the bucket, not a time average; idle buckets may have no observation",
            "phases":"0=setup/teardown, 1=excluded warmup, 2=measurement including final closed-loop drain"
        }
    })
}

#[no_mangle]
/// Writes one new native diagnostic ledger.
///
/// # Safety
/// `path` must point to a readable, NUL-terminated string for the duration of this call.
pub unsafe extern "C" fn cosmos_perf_flow_write(path: *const std::ffi::c_char) -> i32 {
    if path.is_null() {
        return 1;
    }
    // SAFETY: the diagnostic caller supplies a live, NUL-terminated CString.
    let path = unsafe { std::ffi::CStr::from_ptr(path) };
    let Ok(path) = path.to_str() else { return 2 };
    let mut file = match std::fs::OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(path)
    {
        Ok(file) => file,
        Err(error) => {
            eprintln!("flow diagnostic output open failed: {error}");
            return 3;
        }
    };
    if let Err(error) = serde_json::to_writer_pretty(&mut file, &snapshot()) {
        eprintln!("flow diagnostic output write failed: {error}");
        return 4;
    }
    0
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn bucket_event_and_queue_accounting() {
        cosmos_perf_flow_phase(2);
        event(0, 3);
        first_poll();
        terminal();
        queue(2, 3);
        queue(1, 2);
        let value = snapshot();
        assert_eq!(value["overflow"], 0);
        assert_eq!(value["activeAtSnapshot"], 0);
        let rows = value["rows"].as_array().unwrap();
        assert_eq!(
            rows.iter()
                .map(|r| r["counts"]["admitted"].as_u64().unwrap())
                .sum::<u64>(),
            3
        );
        assert_eq!(
            rows.iter()
                .map(|r| r["counts"]["first_poll"].as_u64().unwrap())
                .sum::<u64>(),
            1
        );
        assert_eq!(
            rows.iter()
                .map(|r| r["counts"]["terminal"].as_u64().unwrap())
                .sum::<u64>(),
            1
        );
        assert_eq!(
            rows.iter().map(|r| r["readyHigh"].as_u64().unwrap()).max(),
            Some(2)
        );
        cosmos_perf_flow_phase(0);
        assert!(END_NS[2].load(Ordering::Acquire) > 0);
    }
}
