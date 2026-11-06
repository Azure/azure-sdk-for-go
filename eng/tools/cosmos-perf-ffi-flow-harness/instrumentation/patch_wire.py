#!/usr/bin/env python3
import pathlib
import re
from patch_native import once


def normalize_wire_call(text):
    old = r"crate::perf_wire::point_response\(\s*response\.body\(\),\s*plan\.transcodes_response_to_text\(\),?\s*\);"
    new = r"crate::perf_wire::point_response\(\s*response,\s*plan\.transcodes_response_to_text\(\),?\s*\);"
    if len(re.findall(old, text)) == 1:
        return re.sub(old, "crate::perf_wire::point_response(response, plan.transcodes_response_to_text());", text)
    if len(re.findall(new, text)) == 1:
        return text
    raise ValueError("wire proof call site mismatch")


def patch_wire(root):
    inputs = pathlib.Path(__file__).parent
    driver = root / "sdk/cosmos/azure_data_cosmos_driver"
    native = root / "sdk/cosmos/azure_data_cosmos_driver_native"
    path = driver / "Cargo.toml"
    if "__internal_perf_wire" in path.read_text() or (driver / "src/perf_wire.rs").exists():
        raise ValueError("wire proof source is already patched")
    path.write_text(once(path.read_text(), "[features]\n", "[features]\n__internal_perf_wire = []\n"))
    path = native / "Cargo.toml"
    path.write_text(once(path.read_text(), "[features]\n",
                        '[features]\n__internal_perf_wire = ["azure_data_cosmos_driver/__internal_perf_wire"]\n'))
    path = driver / "src/lib.rs"
    path.write_text(once(path.read_text(), "pub mod options;\n",
                        'pub mod options;\n#[cfg(feature = "__internal_perf_wire")]\n#[doc(hidden)]\npub mod perf_wire;\n'))
    (driver / "src/perf_wire.rs").write_text((inputs / "wire_proof.rs").read_text())
    path = driver / "src/driver/cosmos_driver.rs"
    text = path.read_text()
    start = text.index("    fn finalize_plan_response(\n")
    end = text.index("    async fn uses_cache_only_topology(\n", start)
    part = once(text[start:end], "        let response = response?;\n", """        let response = response?;
        #[cfg(feature = "__internal_perf_wire")]
        if plan.operation.operation_type() == crate::models::OperationType::Read
            && plan.operation.resource_type() == crate::models::ResourceType::Document
        {
            if let Some(response) = response.as_ref() {
                if u16::from(response.status().status_code()) == 200 {
                    crate::perf_wire::point_response(response, plan.transcodes_response_to_text());
                }
            }
        }
""")
    path.write_text(text[:start] + part + text[end:])
    path = native / "src/runtime.rs"
    path.write_text(once(
        path.read_text(),
        "        unsafe { ManuallyDrop::take(&mut self.tokio) }.shutdown_background();\n",
        """        unsafe { ManuallyDrop::take(&mut self.tokio) }.shutdown_background();
        #[cfg(feature = "__internal_perf_wire")]
        if let Err(error) = azure_data_cosmos_driver::perf_wire::write_requested() {
            eprintln!("wire proof output failed: {error}");
        }
"""))
