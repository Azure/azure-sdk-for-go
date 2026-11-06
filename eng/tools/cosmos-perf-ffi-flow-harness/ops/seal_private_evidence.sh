#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
umask 077
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
[[ -f "$root/provenance/CAMPAIGN_COMPLETE" ]]
[[ ! -e "$root/export/private-evidence-v1.tar.xz" ]]
mkdir -p "$root/provenance/symbols"
libc=/usr/lib/x86_64-linux-gnu/libc.so.6
id=$(readelf -n "$libc" | awk '/Build ID:/ {print $3}')
debug=/mnt/cosmos-perf/cpu-isolation-synthetic-20261105-d/attribution-symbols/libc6-dbg/extracted/usr/lib/debug/.build-id/${id:0:2}/${id:2}.debug
cp "$libc" "$root/provenance/symbols/libc.so.6"
cp "$debug" "$root/provenance/symbols/libc.debug"
python3 - "$root" <<'PY'
import hashlib, json, pathlib, sys
root = pathlib.Path(sys.argv[1])
selected = ["run", "provenance", "build/bin", "src/go/clean",
            "src/rust/sdk/cosmos", "src/rust/Cargo.toml", "src/rust/Cargo.lock",
            "clean-comparison", "ffi-comparison", "comparison-v1", "run-spec.json"]
for category in ("clean", "proof", "flow"):
    selected += [f"build/go/{category}/sdk", f"build/go/{category}/adapter",
                 f"build/go/{category}/native", f"build/rust/{category}"]
selected += [path.relative_to(root).as_posix() for path in (root / "tools").glob("inputs-*")]
names = set()
for name in selected:
    path = root / name
    if path.is_file():
        names.add(name)
    elif path.is_dir():
        names.update(item.relative_to(root).as_posix() for item in path.rglob("*")
                     if item.is_file() and not item.is_symlink())
names.discard("provenance/RAW-EVIDENCE-SHA256SUMS")
names.discard("provenance/private-evidence-inputs.list0")
manifest = root / "provenance/RAW-EVIDENCE-SHA256SUMS"
with manifest.open("x") as handle:
    for name in sorted(names):
        digest = hashlib.sha256()
        with (root / name).open("rb") as source:
            for block in iter(lambda: source.read(1024 * 1024), b""):
                digest.update(block)
        handle.write(f"{digest.hexdigest()}  {name}\n")
names.add("provenance/RAW-EVIDENCE-SHA256SUMS")
(root / "provenance/private-evidence-inputs.list0").write_bytes(
    b"".join(name.encode() + b"\0" for name in sorted(names)))
print("PRIVATE_EVIDENCE_INPUT_FILES", len(names))
PY
tar --sort=name --mtime=@0 --owner=0 --group=0 --numeric-owner \
  -C "$root" --null -T "$root/provenance/private-evidence-inputs.list0" -c |
  xz -T2 -3 > "$root/export/private-evidence-v1.tar.xz"
xz --test "$root/export/private-evidence-v1.tar.xz"
chmod 0600 "$root/export/private-evidence-v1.tar.xz"
python3 - "$root" <<'PY'
import hashlib, json, pathlib, shutil, sys
root = pathlib.Path(sys.argv[1])
path = root / "export/private-evidence-v1.tar.xz"
def digest(path):
    value = hashlib.sha256()
    with path.open("rb") as handle:
        for block in iter(lambda: handle.read(1024 * 1024), b""):
            value.update(block)
    return value.hexdigest()
value = {"private": True, "path": str(path), "bytes":path.stat().st_size,
         "sha256":digest(path), "containsRawStackData":True, "notForUnreviewedSharing":True,
         "localDownloadComplete":False}
destination = pathlib.Path("/mnt/cosmos-perf/latest-main-20261105-private-evidence")
if shutil.disk_usage("/mnt/cosmos-perf").free >= path.stat().st_size + 256 * 2**20:
    destination.mkdir(mode=0o700, exist_ok=True)
    copy = destination / path.name
    with path.open("rb") as source, copy.open("xb") as target:
        shutil.copyfileobj(source, target, 8 * 2**20)
    copy.chmod(0o600)
    assert digest(copy) == value["sha256"]
    value["osDiskPreservedCopy"] = str(copy)
else:
    value["osDiskPreservedCopy"] = None
    value["retention"] = "scratch disk retained; insufficient OS-disk headroom for a second copy"
(root / "provenance/private-evidence.json").write_text(json.dumps(value, indent=2) + "\n")
print("PRIVATE_EVIDENCE_META:" + json.dumps(value, separators=(",", ":")))
print("PRIVATE_EVIDENCE_ARCHIVE_SEALED")
PY
