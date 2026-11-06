#!/usr/bin/env python3
import pathlib
import shutil
import sys

from patch_native import once, patch_native
from patch_wire import patch_wire

root = pathlib.Path(sys.argv[1])
inputs = pathlib.Path(__file__).parent
patch_native(root)
patch_wire(root)
adapter = root / "sdk/cosmos/cosmos_perf_flow"
(adapter / "src").mkdir(parents=True, exist_ok=False)
shutil.copyfile(inputs / "rust_main.rs", adapter / "src/main.rs")
shutil.copyfile(inputs / "core.rs", adapter / "src/core.rs")
(adapter / "Cargo.toml").write_text("""[package]
name = "cosmos_perf_flow"
version = "0.0.0"
edition.workspace = true
publish = false

[features]
wire-proof = ["azure_data_cosmos_driver/__internal_perf_wire"]
flow-diagnostic = []

[dependencies]
azure_core = { workspace = true, default-features = false }
azure_data_cosmos = { path = "../azure_data_cosmos", features = ["key_auth", "unstable_patch", "fault_injection"] }
azure_data_cosmos_driver = { path = "../azure_data_cosmos_driver", default-features = false, features = ["tokio", "rustls", "unstable_patch", "fault_injection"] }
serde_json.workspace = true
tokio = { workspace = true, features = ["rt-multi-thread", "macros", "time"] }
""")
path = root / "Cargo.toml"
text = path.read_text()
if '"sdk/cosmos/cosmos_perf_flow"' in text:
    raise SystemExit("runner is already a workspace member")
path.write_text(once(text, "members = [", 'members = [\n  "sdk/cosmos/cosmos_perf_flow",'))
print("FLOW_RUST_PREPARE_COMPLETE")
