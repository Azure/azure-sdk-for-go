#!/usr/bin/env bash
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
inputs=$(cd "$(dirname "$0")" && pwd)
[[ -f "$root/provenance/CAMPAIGN_COMPLETE" && ! -e "$root/comparison-v2" ]]
export HOME=/root GOCACHE="$root/build/go-cache"
python3 "$inputs/profile_attribution.py" "$root" derived-label-v2
python3 "$inputs/report.py" "$root" comparison-v2 derived-label-v2
python3 "$inputs/export_final.py"
python3 - "$root" <<'PY'
import base64, difflib, hashlib, json, lzma, pathlib, sys
root = pathlib.Path(sys.argv[1])
old, new = root / "export/share-v2", root / "export/share-v3"
updates = []
old_names = {p.relative_to(old).as_posix() for p in old.rglob("*") if p.is_file()}
new_names = {p.relative_to(new).as_posix() for p in new.rglob("*") if p.is_file()}
assert old_names == new_names
for name in sorted(new_names):
    a, b = (old / name).read_bytes(), (new / name).read_bytes()
    if a == b:
        continue
    before, after = a.decode().splitlines(keepends=True), b.decode().splitlines(keepends=True)
    operations = []
    for tag, i1, i2, j1, j2 in difflib.SequenceMatcher(a=before, b=after, autojunk=False).get_opcodes():
        if tag == "equal":
            operations.append({"copy":[i1,i2]})
        elif tag in ("replace", "insert"):
            operations.append({"text":after[j1:j2]})
    test = "".join("".join(before[op["copy"][0]:op["copy"][1]]) if "copy" in op else "".join(op["text"])
                   for op in operations).encode()
    assert test == b
    updates.append({"path":name,"oldSHA256":hashlib.sha256(a).hexdigest(),
                    "newSHA256":hashlib.sha256(b).hexdigest(),"operations":operations})
payload = {"baseFolder":"share-v2","newFolder":"share-v3","files":len(new_names),
           "correction":"C ABI export classifier must not match Rust benchmark namespaces",
           "rawMeasurementsChanged":False,"updates":updates}
data = lzma.compress(json.dumps(payload,separators=(",",":")).encode(), preset=9)
(root / "export/attribution-label-patch.json.xz").write_bytes(data)
encoded = base64.b64encode(data).decode()
(root / "export/attribution-label-patch.b64").write_text(encoded)
metadata = {"bytes":len(data),"sha256":hashlib.sha256(data).hexdigest(),
            "base64Bytes":len(encoded),"changedFiles":len(updates)}
(root / "export/attribution-label-patch-meta.json").write_text(json.dumps(metadata)+"\n")
print("LABEL_PATCH_META:"+json.dumps(metadata,separators=(",",":")))
PY
echo ATTRIBUTION_LABEL_CORRECTION_COMPLETE
