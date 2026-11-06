#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then exec /bin/bash "$0" "$@"; fi
set -Eeuo pipefail
root=/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a
find "$root/tools/perf-home/.debug" -maxdepth 7 \( -type f -o -type l \) -printf '%P -> %l\n' | head -n 12
libc=/usr/lib/x86_64-linux-gnu/libc.so.6
echo REAL_LIBC_SECTIONS
readelf -S "$libc" | grep -E 'eh_frame|debug_frame|symtab|debuglink'
debug=/mnt/cosmos-perf/cpu-isolation-synthetic-20261105-d/attribution-symbols/libc6-dbg/extracted/usr/lib/debug/.build-id/6a/6cb5a80fd5e5bb5bf3aeaa111add98e88843a5.debug
echo DEBUG_FILE_SECTIONS
readelf -S "$debug" 2>/dev/null | grep -E 'eh_frame|debug_frame|symtab|debuglink'
echo PERF_INLINE_OPTIONS
perf script -h 2>&1 | grep -E 'inline|symfs|buildid' || true
echo PERF_SYMBOL_CACHE_INSPECTED
