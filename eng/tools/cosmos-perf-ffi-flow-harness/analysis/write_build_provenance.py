#!/usr/bin/env python3
import hashlib
import json
import os
import pathlib
import subprocess
import sys

root = pathlib.Path(sys.argv[1])


def digest(path):
    value = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            value.update(chunk)
    return value.hexdigest()


rust = root / "src/rust"
rust_names = subprocess.check_output(
    ["git", "-C", str(rust), "ls-files", "-c", "-o", "--exclude-standard"], text=True).splitlines()
rust_sources = {name: digest(rust / name) for name in sorted(set(rust_names)) if (rust / name).is_file()}
source_manifest = root / "provenance/rust-source-sha256.json"
if source_manifest.exists():
    backup = source_manifest.with_name("rust-source-sha256-before-" + digest(source_manifest)[:12] + ".json")
    if not backup.exists():
        backup.write_bytes(source_manifest.read_bytes())
source_manifest.write_text(json.dumps(rust_sources, indent=2) + "\n")
wire_paths = {
    "goConverter": root / "build/go/clean/sdk/convert_native.go",
    "goAdapter": root / "build/go/clean/adapter/sdk_v2.go",
    "rustAdapter": rust / "sdk/cosmos/cosmos_perf_flow/src/main.rs",
    "rustOperationOptions": rust / "sdk/cosmos/azure_data_cosmos_driver/src/options/operation_options.rs",
    "rustOperationPipeline": rust / "sdk/cosmos/azure_data_cosmos_driver/src/driver/pipeline/operation_pipeline.rs",
    "rustResponseFinalization": rust / "sdk/cosmos/azure_data_cosmos_driver/src/driver/cosmos_driver.rs",
}
wire_sources = {name: digest(path) for name, path in wire_paths.items()}
wire_fingerprint = digest(root / "provenance/rust-source.patch") + json.dumps(wire_sources, sort_keys=True)
wire_hash = hashlib.sha256(wire_fingerprint.encode()).hexdigest()
(root / "provenance/wire-source-sha256.json").write_text(json.dumps({
    "sha256": wire_hash, "sourceFiles": wire_sources,
    "definition": "common source/configuration subset for wire options and transport; not an executable hash",
}, indent=2) + "\n")
go = "/usr/local/go/bin/go"
for category in ("clean", "proof", "flow"):
    sdk = root / "build/go" / category / "sdk"
    hashes = {p.relative_to(sdk).as_posix(): digest(p) for p in sorted(sdk.rglob("*")) if p.is_file()}
    (root / "provenance" / f"go-{category}-source-sha256.json").write_text(
        json.dumps(hashes, indent=2) + "\n")
    for kind in ("go-v2-native", "rust-bench"):
        binary = root / "build/bin" / category / kind
        environment = os.environ.copy()
        environment["HOME"] = "/root"
        environment["GOCACHE"] = str(root / "build/go-cache")
        notes = subprocess.check_output(["readelf", "-n", str(binary)], text=True)
        import re
        elf = re.findall(r"Build ID: ([0-9a-f]+)", notes)
        if not elf:
            raise ValueError(f"missing ELF build ID: {category}/{kind}")
        provenance = {
            "schemaVersion": 1, "campaign": root.name, "buildClass": category,
            "rustSourceCommit": "77ebe786ad1012bc9ac56105e08c05fb4c65b321",
            "goOptimizedBaseCommit": "b20b0fd00ac20c78bdb2d138fb6960b35e428af0",
            "nativeVersion": "0.3.0",
            "binding": "experimental Linux-only native 0.3.0 header/version/local-distribution port",
            "publishedNativeDistributionUsedUnmodified": False,
            "distributionPackagingBaseCommit": "1e449be78bffef6ef99e6bc1810216a84a9dac16",
            "wireSourceHash": wire_hash,
            "binary": {"sha256": digest(binary), "bytes": binary.stat().st_size,
                       "elfBuildId": elf},
            "nativeArchiveSHA256": digest(root / "build/rust" / category / "libazurecosmosdriver.a"),
            "localABIHeaderSHA256": digest(root / "build/go" / category / "sdk/azurecosmosdriver.h"),
            "rustSourcePatchSHA256": digest(root / "provenance/rust-source.patch"),
            "rustSourceManifestSHA256": digest(root / "provenance/rust-source-sha256.json"),
            "inputPackages": {p.name: digest(p) for p in sorted((root / "tools").glob("*.zip"))},
            "goCompiler": subprocess.check_output([go, "version"], text=True).strip(),
            "rustCompiler": (root / "provenance/rustc-version.txt").read_text(),
            "rustFlags": (root / "provenance/rustflags.txt").read_text().strip(),
            "cargoLockSHA256": digest(root / "provenance/Cargo.lock"),
            "features": {
                "wireProof": category == "proof", "ffiFlow": category == "flow" and kind == "go-v2-native",
                "rustOperationBuckets": category == "flow" and kind == "rust-bench",
                "tokioConsole": False,
            },
        }
        if kind == "go-v2-native":
            provenance["binary"]["goBuildId"] = subprocess.check_output(
                [go, "tool", "buildid", str(binary)], text=True, env=environment).strip()
        target = pathlib.Path(str(binary) + ".provenance.json")
        if target.exists():
            backup = target.with_name(target.name + ".superseded-" + digest(target)[:12])
            if not backup.exists():
                backup.write_bytes(target.read_bytes())
        with target.open("w") as handle:
            json.dump(provenance, handle, indent=2)
            handle.write("\n")
print("BUILD_PROVENANCE_SEALED")
