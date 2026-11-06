#!/usr/bin/env python3
import argparse
import pathlib


def once(text, old, new):
    if text.count(old) != 1:
        raise ValueError(f"Expected unique source anchor: {old[:100]!r}")
    return text.replace(old, new, 1)


def patch_native(root):
    crate = root / "sdk/cosmos/azure_data_cosmos_driver_native"
    path = crate / "Cargo.toml"
    if "__internal_perf_flow" in path.read_text() or (crate / "src/perf_flow.rs").exists():
        raise ValueError("native flow source is already patched")
    path.write_text(once(path.read_text(), "[features]\n", "[features]\n__internal_perf_flow = []\n"))
    path = crate / "src/lib.rs"
    path.write_text(once(
        path.read_text(), "#[macro_use]\npub mod string;",
        """macro_rules! flow {
    ($name:ident($($argument:expr),*)) => {
        #[cfg(feature = "__internal_perf_flow")]
        crate::perf_flow::$name($($argument),*);
    };
}
#[cfg(feature = "__internal_perf_flow")]
mod perf_flow;

#[macro_use]
pub mod string;"""))
    (crate / "src/perf_flow.rs").write_text(
        pathlib.Path(__file__).with_name("flow.rs").read_text())
    path = crate / "src/runtime.rs"
    path.write_text(once(
        path.read_text(),
        "let tokio = tokio::runtime::Builder::new_multi_thread()\n            .enable_all()",
        "let tokio = tokio::runtime::Builder::new_multi_thread()\n            .worker_threads(2)\n            .enable_all()"))
    path = crate / "src/submit.rs"
    text = once(
        path.read_text(),
        "    runtime.tokio.spawn(async move {\n",
        "    flow!(event(7, 1));\n    runtime.tokio.spawn(async move {\n        flow!(first_poll());\n")
    text = once(
        text,
        "        let rc = CompletionQueue::enqueue_into_inner(&ctx.queue, completion);",
        "        flow!(terminal());\n        let rc = CompletionQueue::enqueue_into_inner(&ctx.queue, completion);")
    path.write_text(text)
    path = crate / "src/completion.rs"
    text = once(
        path.read_text(),
        "        guard.in_flight += 1;\n        Ok(())",
        "        guard.in_flight += 1;\n        flow!(event(0, 1));\n        flow!(queue(guard.deque.len(), guard.in_flight));\n        Ok(())")
    start = text.index("    pub(crate) fn release_in_flight(&self)")
    end = text.index("\n    ///", start)
    part = once(text[start:end], "        maybe_mark_drained(&mut guard);",
                "        flow!(event(1, 1));\n        flow!(queue(guard.deque.len(), guard.in_flight));\n        maybe_mark_drained(&mut guard);")
    text = text[:start] + part + text[end:]
    start = text.index("    pub(crate) fn enqueue_into_inner(")
    end = text.index("\n}\n", start)
    part = text[start:end]
    for reason in ("CosmosErrorCodeQueueShutdown", "CosmosErrorCodeQueueFull"):
        part = once(part,
                    f"            return CosmosErrorCode::{reason};",
                    f"            flow!(event(5, 1));\n            flow!(queue(guard.deque.len(), guard.in_flight));\n            return CosmosErrorCode::{reason};")
    part = once(part, "        guard.deque.push_back(c);",
                "        guard.deque.push_back(c);\n        flow!(event(4, 1));\n        flow!(queue(guard.deque.len(), guard.in_flight));")
    text = text[:start] + part + text[end:]
    text = once(
        text,
        "    let c = guard.deque.pop_front()?;\n    guard.in_flight = guard.in_flight.saturating_sub(1);",
        "    let c = guard.deque.pop_front()?;\n    guard.in_flight = guard.in_flight.saturating_sub(1);\n    flow!(event(6, 1));\n    flow!(queue(guard.deque.len(), guard.in_flight));")
    path.write_text(text)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("root", type=pathlib.Path)
    args = parser.parse_args()
    patch_native(args.root)
    print("FLOW_NATIVE_PATCH_COMPLETE")
