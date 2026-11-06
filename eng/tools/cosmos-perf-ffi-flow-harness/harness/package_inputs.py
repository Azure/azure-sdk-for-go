#!/usr/bin/env python3
import base64
import hashlib
import io
import pathlib
import sys
import zipfile

root = pathlib.Path(__file__).parent
names = [
    "flow.rs", "flow.go", "flow_cgo.go", "flow_test.go",
    "barrier_go.go", "barrier_go_test.go", "rust_main.rs", "core.rs",
    "patch_native.py", "prepare_rust.py", "validate_flow.py", "run_leg.py",
    "calibrate.py", "report.py", "run-spec.json", "stage_sources.sh",
    "wire_proof.rs", "patch_wire.py", "build_rust.sh", "resume_rust.sh",
    "patch_go.py", "diagnostics_go.go", "build_go.sh", "resume_go.sh", "write_build_provenance.py",
    "run_campaign.sh", "resume_campaign.sh", "run_clean.sh", "run_diagnostics.sh", "run_ffi.sh", "flow_summary.py", "clean_summary.py", "analyze_go_diagnostics.py", "finalize_campaign.sh", "export_final.sh", "export_final.py", "correct_attribution_labels.sh", "refresh_rust_runner.sh", "validate_rust_flow.py", "profile_attribution.py", "attribute.py",
]
build_rust = sys.argv[1:] == ["build-rust"]
resume_rust = sys.argv[1:] == ["resume-rust"]
build_go = sys.argv[1:] == ["build-go"]
resume_go = sys.argv[1:] == ["resume-go"]
run_campaign = sys.argv[1:] == ["run"]
resume_campaign = sys.argv[1:] == ["resume-run"]
refresh_rust = sys.argv[1:] == ["refresh-rust"]
clean_only = sys.argv[1:] == ["clean"]
diagnostics_only = sys.argv[1:] == ["diagnostics"]
ffi_only = sys.argv[1:] == ["ffi"]
finalize = sys.argv[1:] == ["finalize"]
export_final = sys.argv[1:] == ["export"]
correct_labels = sys.argv[1:] == ["correct-labels"]
if sys.argv[1:] and not (build_rust or resume_rust or build_go or resume_go or run_campaign or resume_campaign or refresh_rust or clean_only or diagnostics_only or ffi_only or finalize or export_final or correct_labels):
    raise SystemExit("usage: package_inputs.py [build-rust|resume-rust|build-go|resume-go|run|resume-run|refresh-rust|clean|diagnostics|ffi|finalize|export]")
label = "inputs-v3" if resume_rust else ("inputs-v2" if build_rust else "inputs-v1")
marker = "INPUTS_V3_COMPLETE" if resume_rust else ("INPUTS_V2_COMPLETE" if build_rust else "INPUTS_COMPLETE")
entry = "resume_rust.sh" if resume_rust else ("build_rust.sh" if build_rust else "stage_sources.sh")
if build_go:
    entry = "build_go.sh"
if resume_go:
    entry = "resume_go.sh"
if run_campaign:
    entry = "run_campaign.sh"
if resume_campaign:
    entry = "resume_campaign.sh"
if refresh_rust:
    entry = "refresh_rust_runner.sh"
if clean_only:
    entry = "run_clean.sh"
if diagnostics_only:
    entry = "run_diagnostics.sh"
if ffi_only:
    entry = "run_ffi.sh"
if finalize:
    entry = "finalize_campaign.sh"
if export_final:
    entry = "export_final.sh"
if correct_labels:
    entry = "correct_attribution_labels.sh"
buffer = io.BytesIO()
with zipfile.ZipFile(buffer, "w", compression=zipfile.ZIP_DEFLATED, compresslevel=9) as archive:
    for name in names:
        info = zipfile.ZipInfo(name)
        info.compress_type = zipfile.ZIP_DEFLATED
        archive.writestr(info, (root / name).read_text().replace("\r\n", "\n").encode())
data = buffer.getvalue()
digest = hashlib.sha256(data).hexdigest()
if resume_rust:
    label = "inputs-resume-" + digest[:12]
    marker = "INPUTS_RESUME_" + digest[:12]
if build_go or resume_go:
    label = "inputs-go-" + digest[:12]
    marker = "INPUTS_GO_" + digest[:12]
if run_campaign or resume_campaign or refresh_rust or clean_only or diagnostics_only or ffi_only or finalize or export_final or correct_labels:
    action = ("correct-labels" if correct_labels else ("export" if export_final else ("finalize" if finalize else ("ffi" if ffi_only else ("clean" if clean_only else ("diagnostics" if diagnostics_only else
              ("refresh-rust" if refresh_rust else ("resume-run" if resume_campaign else "run"))))))))
    label = "inputs-" + action + "-" + digest[:12]
    marker = "INPUTS_" + action.replace("-", "_").upper() + "_" + digest[:12]
(root / (label + ".zip")).write_bytes(data)
script = """#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
[[ -f /mnt/cosmos-ffi-flow-scratch-20261105-a/.campaign-ownership.json ]]
[[ ! -e "$root/tools/INPUTS_COMPLETE" ]]
mkdir -p "$root/tools/inputs-v1"
base64 --decode > "$root/tools/inputs-v1.zip" <<'INPUT_PAYLOAD'
""" + base64.b64encode(data).decode() + "\nINPUT_PAYLOAD\n"
script += f"""[[ "$(sha256sum "$root/tools/inputs-v1.zip" | awk '{{print $1}}')" == "{digest}" ]]
python3 - "$root/tools/inputs-v1.zip" "$root/tools/inputs-v1" <<'PY'
import pathlib, sys, zipfile
with zipfile.ZipFile(sys.argv[1]) as archive:
    assert set(archive.namelist()) == {set(names)!r}, "unexpected input archive paths"
    for name in archive.namelist():
        with (pathlib.Path(sys.argv[2]) / name).open("xb") as handle:
            handle.write(archive.read(name))
PY
printf '%s\\n' "{digest}" > "$root/tools/INPUTS_COMPLETE"
exec /bin/bash "$root/tools/inputs-v1/stage_sources.sh"
"""
script = script.replace("inputs-v1", label).replace("INPUTS_COMPLETE", marker)
script = script.replace("exec /bin/bash \"$root/tools/" + label + "/stage_sources.sh\"",
                        "exec /bin/bash \"$root/tools/" + label + "/" + entry + "\"")
(root / ("run-campaign-inputs.sh" if (run_campaign or resume_campaign or refresh_rust or clean_only or diagnostics_only or ffi_only or finalize or export_final or correct_labels) else ("build-go-inputs.sh" if (build_go or resume_go) else ("resume-rust-inputs.sh" if resume_rust else ("build-rust-inputs.sh" if build_rust else "stage-inputs.sh"))))).write_text(
    script, encoding="utf-8", newline="\n")
print("INPUT_PACKAGE_SHA256=" + digest)
