import base64
import collections
import csv
import hashlib
import json
import lzma
import pathlib
import re

root = pathlib.Path("/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a")
out_dir = root / "analysis/cpu-deep-v1"
out_dir.mkdir(parents=True, exist_ok=True)
_accepted_path = root / "provenance/accepted-legs.json"
accepted = json.loads(_accepted_path.read_text()) if _accepted_path.exists() else {}
_probe = json.loads((root / "run/w64/perf/repeat-1/go-v2-native/derived-label-v2/attribution.json").read_text())
print("ATTR_KEYS:" + ",".join(sorted(_probe)[:40]))

HEADER = re.compile(r"^(?P<comm>\S.*?)\s+(?P<pid>\d+)(?:/(?P<tid>\d+))?\s+(?:\[\d+\]\s+)?(?P<ts>\d+\.\d+):\s+(?P<period>\d+)?\s*(?P<event>[\w:-]+):")
FRAME = re.compile(r"^\s+[0-9a-f]+\s+(?P<sym>.+?)\s+\((?P<dso>[^()]*)\)\s*$")
HASH = re.compile(r"::h[0-9a-f]{16}$")
OFFSET = re.compile(r"\+0x[0-9a-f]+$")

INCLUSIVE = {
    "go: runtime.cgocall (Go->C calls)": r"^runtime\.cgocall$",
    "go: cgo callbacks (C->Go)": r"^runtime\.cgocallbackg?$",
    "go: runtime.mallocgc": r"^runtime\.mallocgc$",
    "go: GC background mark": r"^runtime\.gcBgMarkWorker$",
    "go: GC assist": r"^runtime\.gcAssistAlloc$",
    "go: scheduler schedule/findRunnable": r"^runtime\.(schedule|findRunnable)$",
    "go: park/sleep (futexsleep/stopm)": r"^runtime\.(futexsleep|stopm|notesleep)$",
    "go: wake (futexwakeup/startm/wakep)": r"^runtime\.(futexwakeup|startm|wakep|notewakeup)$",
    "go: netpoll": r"^runtime\.netpoll$",
    "go: sysmon": r"^runtime\.sysmon$",
    "go: chansend/chanrecv": r"^runtime\.(chansend|chanrecv)1?$",
    "go-sdk: reactor.drainOnce": r"azcosmos/v2\.\(\*reactor\)\.drainOnce$",
    "go-sdk: nativeDriver.execute": r"azcosmos/v2\.\(\*nativeDriver\)\.execute$",
    "go-sdk: translateCompletion": r"azcosmos/v2\.translateCompletion$",
    "go-sdk: diagnostics copy/JSON": r"azcosmos/v2\.(newDiagnosticsJSON|copyDiagnostics|copyDiagnosticAttempts)",
    "go-sdk: OperationOptions.toNative": r"azcosmos/v2\.OperationOptions\.toNative$",
    "go-sdk: cgo GoBytes/GoString": r"(^|\.)_?C?f?u?n?c?_?Go(Bytes|String)$|^runtime\.(gobytes|gostring)",
    "native: C ABI exports cosmos_*": r"^cosmos_[a-z0-9_]+$",
    "rust: tokio worker park/unpark": r"^tokio::runtime::(scheduler::multi_thread::worker::Context::park|park::)",
    "rust: tokio scheduler total": r"^tokio::runtime::scheduler::",
    "rust: h2": r"^<?h2::",
    "rust: hyper": r"^<?hyper(_util)?::",
    "rust: reqwest": r"^<?reqwest::",
    "rust: TLS (openssl/rustls)": r"^(SSL_|ssl3_|tls_|EVP_|<?rustls::|<?openssl::)",
    "rust: serde/serde_json": r"^<?serde(_json)?::",
    "rust: HMAC/SHA signing": r"(hmac|sha2::|HMAC|SHA256)",
    "rust: azure_data_cosmos_driver": r"^<?azure_data_cosmos_driver::",
    "rust: azure_data_cosmos SDK": r"^<?azure_data_cosmos::",
    "rust: azure_core": r"^<?azure_core::",
    "libc: malloc/free family": r"^(_int_malloc|_int_free|malloc|free|__libc_malloc|__libc_free|cfree|realloc|_int_realloc|malloc_consolidate|tcache_)",
    "libc: memcpy/memmove/memset": r"^__mem(cpy|move|set)",
    "kernel: futex": r"^(do_futex|__x64_sys_futex)$",
    "kernel: epoll": r"^(do_epoll_wait|ep_poll)$",
    "kernel: TCP send": r"^tcp_sendmsg$",
    "kernel: TCP recv": r"^tcp_recvmsg$",
    "kernel: __schedule (context switch)": r"^__schedule$",
    "kernel: all syscalls": r"^entry_SYSCALL_64",
    "kernel: interrupts/softirq": r"^(asm_common_interrupt|handle_softirqs|asm_sysvec)",
}
INCLUSIVE = {key: re.compile(value) for key, value in INCLUSIVE.items()}
SYSCALL_KIND = [
    ("futex", re.compile(r"futex")), ("epoll", re.compile(r"epoll|ep_poll")),
    ("tcp-send", re.compile(r"tcp_sendmsg|sock_sendmsg|__sys_sendto|sock_write")),
    ("tcp-recv", re.compile(r"tcp_recvmsg|sock_recvmsg|__sys_recvfrom|sock_read")),
    ("nanosleep", re.compile(r"nanosleep")), ("sched_yield", re.compile(r"sched_yield")),
    ("mmap/munmap/madvise", re.compile(r"mmap|munmap|madvise|brk")),
    ("read/write-other", re.compile(r"ksys_read|ksys_write|vfs_")),
    ("signal", re.compile(r"rt_sig|signal")),
]
ALLOC_LEAF = re.compile(r"^(_int_malloc|_int_free|malloc|free|__libc_|cfree|realloc|_int_realloc|malloc_consolidate|tcache|__mem|runtime\.(mallocgc|memmove|memclr|nextFreeFast|heapBits|scanobject|greyobject|findObject|spanOf|wbBuf|bulkBarrier|gcWriteBarrier|typedmemmove|growslice|makeslice|newobject|rawstring|gobytes|concatstring|slicebytetostring|mapassign|mallocgcSmall|mallocgcTiny|mallocgcLarge)|alloc::|core::ptr::drop_in_place|<alloc::|bytes::|__rust_(alloc|dealloc|realloc)|__rdl_)")


def short(symbol):
    symbol = HASH.sub("", OFFSET.sub("", symbol))
    symbol = symbol.replace("github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2.", "azcosmos/v2.")
    symbol = re.sub(r"\[[0-9a-f]{16}\]", "", symbol)
    return symbol[:160]


def thread_group(comm):
    comm = re.sub(r"[-_ ]?\d+$", "", comm.strip())
    return comm or "?"


def language(symbol, dso):
    if dso.startswith("[kernel") or dso == "[kernel.kallsyms]":
        return "kernel"
    if symbol.startswith(("runtime.", "sync.", "internal/", "main.", "github.com/", "azcosmos/", "syscall.", "time.", "context.")):
        return "go"
    if "libc.so" in dso or "ld-linux" in dso:
        return "libc"
    return "native"


def parse_script(path):
    samples = []
    current = None
    with path.open(errors="replace") as handle:
        for line in handle:
            if not line.strip():
                if current:
                    samples.append(current)
                    current = None
                continue
            if not line[0].isspace():
                match = HEADER.match(line)
                if not match:
                    raise ValueError("unparsed perf header: " + line[:120])
                if current:
                    samples.append(current)
                current = {"comm": match["comm"], "tid": match["tid"] or match["pid"], "ts": match["ts"], "frames": []}
            else:
                match = FRAME.match(line)
                if match and current is not None:
                    current["frames"].append((OFFSET.sub("", match["sym"]), match["dso"]))
    if current:
        samples.append(current)
    return samples


def ts_key(value):
    return value.rstrip("0").rstrip(".")


RUST_ESC = {"SP": "@", "BP": "*", "RF": "&", "LT": "<", "GT": ">", "LP": "(", "RP": ")", "C": ","}


def rust_decode(text):
    if text.startswith("_$"):
        text = text[1:]
    text = re.sub(r"\$(SP|BP|RF|LT|GT|LP|RP|C)\$", lambda m: RUST_ESC[m[1]], text)
    text = re.sub(r"\$u([0-9a-f]{2,6})\$", lambda m: chr(int(m[1], 16)), text)
    return text.replace("..", "::")


def rust_demangle(symbol):
    symbol = re.sub(r"\.llvm\.\d+$", "", symbol)
    if symbol.startswith("_ZN") and symbol.endswith("E"):
        body, parts, index = symbol[3:-1], [], 0
        while index < len(body) and body[index].isdigit():
            end = index
            while body[end].isdigit():
                end += 1
            length = int(body[index:end])
            parts.append(body[end:end + length])
            index = end + length
        if index == len(body) and parts:
            if re.fullmatch(r"h[0-9a-f]{16}", parts[-1]):
                parts.pop()
            return "::".join(rust_decode(p) for p in parts)
    if "$LT$" in symbol or "$u20$" in symbol or ".." in symbol and "$" in symbol:
        return rust_decode(symbol)
    return symbol


def demangle_all(samples):
    import subprocess
    raw = sorted({s for sample in samples for s, _ in sample["frames"]})
    out = subprocess.run(["c++filt"], input="\n".join(raw), capture_output=True, text=True, check=True).stdout.splitlines()
    assert len(out) == len(raw)
    table = {a: short(rust_demangle(b)) for a, b in zip(raw, out)}
    for sample in samples:
        sample["frames"] = [(table[s], d) for s, d in sample["frames"]]


def leg_analysis(workers, leg):
    relative = accepted.get(f"w{workers}/perf/repeat-1/{leg}", f"run/w{workers}/perf/repeat-1/{leg}")
    directory = root / relative
    derived = directory / "derived-label-v2"
    attribution = json.loads((derived / "attribution.json").read_text())
    rows = list(csv.DictReader((derived / "sample-mapping.csv").open()))
    script = parse_script(directory / "raw/perf-script.txt")
    demangle_all(script)
    by_ts = collections.defaultdict(list)
    for sample in script:
        by_ts[ts_key(sample["ts"])].append(sample)
    total_period = sum(int(row["period"]) for row in rows)
    assert abs(total_period - attribution["totalPeriodNs"]) <= 1, (total_period, attribution["totalPeriodNs"])
    us_per_period = attribution["profileCPUUsPerRead"] / total_period
    groups = collections.Counter()
    matrix = collections.Counter()
    self_cat = collections.Counter()
    alloc_owner = collections.Counter()
    alloc_owner_lang = collections.Counter()
    alloc_leaf = collections.Counter()
    syscall_kind = collections.Counter()
    syscall_origin = collections.Counter()
    sched_leaf = collections.Counter()
    go_runtime_leaf = collections.Counter()
    inclusive = collections.Counter()
    leaf_lang = collections.Counter()
    unmatched = 0
    missing_ts = []
    leaf_mismatch = []
    used = collections.Counter()
    for row in rows:
        key = ts_key(row["timestamp"])
        candidates = by_ts.get(key, [])
        index = used[key]
        if index >= len(candidates):
            unmatched += 1
            missing_ts.append(row["timestamp"])
            continue
        used[key] += 1
        sample = candidates[index]
        weight = int(row["period"]) * us_per_period
        frames = sample["frames"]
        if frames and short(row["leaf"]) != frames[0][0] and row["originalLeaf"] == "":
            if not frames[0][0].startswith(short(row["leaf"])[:40]):
                leaf_mismatch.append((short(row["leaf"]), frames[0][0]))
        group = thread_group(sample["comm"])
        groups[group] += weight
        category = row["selfCategory"]
        self_cat[category] += weight
        matrix[(category, row["callerCategory"])] += weight
        leaf_dso = frames[0][1] if frames else row["dso"]
        leaf_lang[language(short(row["leaf"]), leaf_dso)] += weight
        seen = set()
        for symbol, _ in frames:
            for name, pattern in INCLUSIVE.items():
                if name not in seen and pattern.search(symbol):
                    seen.add(name)
                    inclusive[name] += weight
        if category == "serialization-body-allocation-memcopy":
            alloc_leaf[short(row["leaf"])] += weight
            owner = next(((s, d) for s, d in frames if not ALLOC_LEAF.search(s) and not d.startswith("[kernel")), None)
            if owner:
                alloc_owner[owner[0]] += weight
                alloc_owner_lang[language(*owner)] += weight
            else:
                alloc_owner["<no owning frame>"] += weight
        if category == "go-runtime-sdk":
            go_runtime_leaf[short(row["leaf"])] += weight
        if category == "kernel-scheduling-syscalls":
            sched_leaf[short(row["leaf"])] += weight
        syscall_index = next((i for i, (s, _) in enumerate(frames) if s.startswith("entry_SYSCALL_64")), None)
        if syscall_index is not None:
            kernel_part = " ".join(s for s, _ in frames[:syscall_index])
            kind = next((name for name, pattern in SYSCALL_KIND if pattern.search(kernel_part)), "other")
            syscall_kind[kind] += weight
            user = next((s for s, d in frames[syscall_index + 1:]
                         if not d.startswith("[kernel") and not re.match(r"^(syscall|__syscall|__GI_|__libc_|__futex|futex_|epoll_wait|__internal_syscall|internal/runtime/syscall|runtime\.(futex|epollwait|usleep|nanotime|write1|read|osyield|raise|tgkill|madvise|sysMmap|sysMunmap|mmap|munmap)|__lll_|pthread_cond|___pthread|__pthread|parking_lot|std::sys::|std::thread::park|mio::|tokio::runtime::(io|park))", s)), "<syscall wrapper only>")
            syscall_origin[f"{kind} <- {user}"] += weight
    print("JOIN", workers, leg, "missingTs", len(missing_ts), "leafMismatch", len(leaf_mismatch),
          missing_ts[:2], leaf_mismatch[:3], list(by_ts)[:2], flush=True)
    assert unmatched == 0, f"{workers}/{leg}: {unmatched} unmatched samples"
    assert len(leaf_mismatch) <= 0.01 * len(rows), f"{workers}/{leg}: leaf mismatch {len(leaf_mismatch)}"
    total = sum(self_cat.values())
    assert abs(total - attribution["profileCPUUsPerRead"]) < 1e-6 * total + 1e-6

    def top(counter, count=20):
        return [[key if isinstance(key, str) else " | ".join(key), round(value, 4)]
                for key, value in counter.most_common(count)]

    return {
        "workers": workers, "leg": leg, "samples": len(rows), "scriptSamples": len(script),
        "profileCPUUsPerRead": attribution["profileCPUUsPerRead"],
        "profileThroughput": attribution["profileThroughput"],
        "processCPUSeconds": attribution["processCPUSeconds"],
        "sampledToProcessCPUPercent": attribution["sampledToProcessCPUPercent"],
        "threadGroupsUs": top(groups, 30),
        "leafLanguageUs": top(leaf_lang),
        "selfCategoryUs": top(self_cat, 30),
        "selfByCallerUs": top(matrix, 60),
        "allocLeafUs": top(alloc_leaf, 20),
        "allocOwnerUs": top(alloc_owner, 40),
        "allocOwnerLanguageUs": top(alloc_owner_lang),
        "goRuntimeLeafUs": top(go_runtime_leaf, 25),
        "kernelSchedLeafUs": top(sched_leaf, 20),
        "syscallKindUs": top(syscall_kind),
        "syscallOriginUs": top(syscall_origin, 30),
        "inclusiveUs": sorted([[key, round(value, 4)] for key, value in inclusive.items()], key=lambda r: -r[1]),
    }


LAYERS = [
    ("TLS", r"^(SSL_|ssl3_|tls_|EVP_|ossl_|OPENSSL|CRYPTO_|BIO_|aesni|sha256_block|<?rustls::|<?openssl::)"),
    ("h2 (HTTP/2 framing+HPACK)", r"^<?h2::"),
    ("hyper", r"^<?hyper(_util)?::"),
    ("reqwest", r"^<?reqwest::"),
    ("driver: transport", r"azure_data_cosmos_driver::driver::transport"),
    ("driver: routing/session", r"azure_data_cosmos_driver::driver::routing"),
    ("driver: dataflow/topology", r"azure_data_cosmos_driver::driver::dataflow"),
    ("driver: pipeline", r"azure_data_cosmos_driver::driver::pipeline"),
    ("driver: other", r"azure_data_cosmos_driver::"),
    ("azure_core/typespec", r"^<?(azure_core|typespec)(_client_core)?::"),
    ("Rust SDK (azure_data_cosmos)", r"^<?azure_data_cosmos::"),
    ("native C ABI / completion", r"^(cosmos_[a-z0-9_]+$|<?azurecosmosdriver::)"),
    ("Go SDK", r"azcosmos/v2\."),
    ("Go runtime", r"^runtime\."),
    ("benchmark adapter", r"^<?(cosmos_perf|azure_data_cosmos_perf|main::|perf_adapter)"),
    ("tokio runtime", r"^<?tokio::"),
]
LAYERS = [(n, re.compile(p)) for n, p in LAYERS]


def layer_of(frames):
    for symbol, dso in frames:
        if dso.startswith("[kernel"):
            continue
        for name, pattern in LAYERS:
            if pattern.search(symbol):
                return name
    return "unattributed (std/libc only stack)"


def modname(sym):
    s = sym.lstrip("<")
    s = s.split(" as ")[0]
    parts = re.split(r"::", s)
    out = []
    for p in parts:
        if not p or p[0].isupper() or "{" in p or "<" in p or "$" in p:
            break
        out.append(p)
    return "::".join(out[1:4]) or s[:60]


def tidy(sym):
    s = re.sub(r"::_\$u7b\$\$u7b\$closure\$u7d\$\$u7d\$|::\{\{closure\}\}", "{closure}", sym)
    s = s.replace("azure_data_cosmos_driver::", "")
    return s[:110]


def other(workers, leg):
    relative = accepted.get(f"w{workers}/perf/repeat-1/{leg}", f"run/w{workers}/perf/repeat-1/{leg}")
    directory = root / relative
    attribution = json.loads((directory / "derived-label-v2/attribution.json").read_text())
    rows = list(csv.DictReader((directory / "derived-label-v2/sample-mapping.csv").open()))
    script = parse_script(directory / "raw/perf-script.txt")
    demangle_all(script)
    by_ts = collections.defaultdict(list)
    for sample in script:
        by_ts[ts_key(sample["ts"])].append(sample)
    scale = attribution["profileCPUUsPerRead"] / sum(int(r["period"]) for r in rows)
    used = collections.Counter()
    mods, fns, outer = collections.Counter(), collections.Counter(), collections.Counter()
    kinds = collections.defaultdict(collections.Counter)
    for row in rows:
        key = ts_key(row["timestamp"])
        sample = by_ts[key][used[key]]
        used[key] += 1
        frames = sample["frames"]
        if layer_of(frames) != "driver: other":
            continue
        w = int(row["period"]) * scale
        fn = next(s for s, d in frames if not d.startswith("[kernel") and "azure_data_cosmos_driver::" in s)
        if frames[0][1].startswith("[kernel"):
            kind = "K"
        elif row["selfCategory"] == "serialization-body-allocation-memcopy":
            kind = "A"
        else:
            kind = "C"
        mods[modname(fn)] += w
        fns[tidy(fn)] += w
        kinds[tidy(fn)][kind] += w
        drv = [s for s, d in frames if "azure_data_cosmos_driver::" in s]
        outer[tidy(drv[-1])] += w
    return {"w": workers, "total": round(sum(mods.values()), 2), "mods": mods.most_common(12),
            "fns": [[k, round(v, 2), {a: round(b, 1) for a, b in kinds[k].items()}] for k, v in fns.most_common(22)],
            "outer": outer.most_common(6)}


res = [other(w, "go-v2-native") for w in (64, 32)]
(out_dir / "driver-other.json").write_text(json.dumps(res))
for r in res[:1]:
    print("TOTAL", r["total"])
    for k, v in r["mods"]:
        print(f"M {v:5.1f} {k}")
    for k, v, kd in r["fns"]:
        print(f"F {v:5.1f} {kd} {k}")
    for k, v in r["outer"]:
        print(f"O {v:5.1f} {k}")
print("W32TOTAL", res[1]["total"], [[k, round(v, 1)] for k, v in res[1]["mods"][:6]])
print("VM_PY_COMPLETE")
