#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail

mount=/mnt/cosmos-ffi-flow-scratch-20261105-a
campaign=go-v2-ffi-flow-latest-main-20261105-a
ownership=__OWNERSHIP_PROOF__
lun=__SCRATCH_LUN__
expected_bytes=34359738368

udevadm settle --timeout=30
device=
for candidate in "/dev/disk/azure/scsi1/lun$lun" "/dev/disk/azure/data/by-lun/$lun"; do
  if [[ -b "$candidate" ]]; then
    resolved=$(readlink -f "$candidate")
    [[ -z "$device" || "$device" == "$resolved" ]] || {
      echo "Conflicting Azure LUN device links; refusing to format" >&2
      exit 1
    }
    device=$resolved
  fi
done
[[ -n "$device" && -b "$device" ]] || {
  echo "Exact approved Azure LUN is not discoverable; no device was formatted" >&2
  exit 1
}
[[ "$(blockdev --getsize64 "$device")" == "$expected_bytes" ]]
root_device=$(findmnt -rn -o SOURCE /)
if lsblk -snrpo NAME "$root_device" | grep -Fxq "$device"; then
  echo "Refusing to format a root-filesystem ancestor" >&2
  exit 1
fi

if mountpoint -q "$mount"; then
  [[ "$(readlink -f "$(findmnt -rn -M "$mount" -o SOURCE)")" == "$device" ]]
  python3 - "$mount/.campaign-ownership.json" "$ownership" <<'PY'
import json, pathlib, sys
value = json.loads(pathlib.Path(sys.argv[1]).read_text())
assert value["ownershipProof"] == sys.argv[2], "existing mount is not this campaign's disk"
PY
else
  if [[ -e "$mount" && -n "$(find "$mount" -mindepth 1 -maxdepth 1 -print -quit)" ]]; then
    echo "Refusing to hide existing files beneath a new mount" >&2
    exit 1
  fi
  [[ "$(lsblk -nrpo TYPE "$device")" == disk ]] || {
    echo "Unexpected child partitions; no device was formatted" >&2
    exit 1
  }
  [[ -z "$(lsblk -nrpo MOUNTPOINT "$device" | tr -d '[:space:]')" ]]
  filesystem=$(lsblk -dnro FSTYPE "$device")
  fresh=false
  if [[ -z "$filesystem" ]]; then
    [[ -z "$(wipefs --no-act "$device")" ]] || {
      echo "Unexpected disk signature; refusing to format" >&2
      exit 1
    }
    mkfs.ext4 -q -m 0 -L cosmosflow32 "$device"
    fresh=true
  elif [[ "$filesystem" != ext4 ]]; then
    echo "Unexpected existing filesystem; refusing to format" >&2
    exit 1
  fi
  mkdir -p "$mount"
  mount -o noatime "$device" "$mount"
  if [[ "$fresh" == true ]]; then
    python3 - "$mount/.campaign-ownership.json" "$ownership" "$campaign" <<'PY'
import json, os, pathlib, sys
path = pathlib.Path(sys.argv[1])
with path.open("x") as handle:
    json.dump({"ownershipProof": sys.argv[2], "campaign": sys.argv[3],
               "newDiskOnly": True, "sizeGiB": 32}, handle)
    handle.write("\n")
os.chmod(path, 0o600)
PY
  else
    python3 - "$mount/.campaign-ownership.json" "$ownership" <<'PY'
import json, pathlib, sys
value = json.loads(pathlib.Path(sys.argv[1]).read_text())
assert value["ownershipProof"] == sys.argv[2], "filesystem ownership not verified"
PY
  fi
fi

root=$mount/$campaign
mkdir -p "$root"/{src/go,build,private,provenance,tools}
chmod 700 "$root/private"
python3 - "$mount" "$root" <<'PY'
import json, os, pathlib, shutil, sys
mount, root = map(pathlib.Path, sys.argv[1:])
free = shutil.disk_usage(mount).free
assert free >= 30 * 2**30, "new scratch disk has less than 30 GiB free"
assert os.cpu_count() == 4, "four-vCPU VM is required"
print(json.dumps({"scratchMounted": True, "freeGiB": round(free / 2**30, 2),
                  "cpus": os.cpu_count(), "root": str(root), "preExistingDisksChanged": False}))
PY
echo SCRATCH_MOUNT_COMPLETE
