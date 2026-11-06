import hashlib
import html
import json
import pathlib

out = pathlib.Path(__file__).resolve().parent / "cpu-deep-v1"
data = json.loads((out / "cpu-deep.json").read_text())
legs = {(l["workers"], l["leg"]): l for l in data["legs"]}
D = lambda leg, key: {k: v for k, v in leg[key]}

GROUPS = [
    ("Native malloc/free/memcpy (Rust code, glibc)", ["serialization-body-allocation-memcopy"], "#d62728"),
    ("Go runtime + Go SDK", ["go-runtime-sdk"], "#1f77b4"),
    ("Kernel futex / context switch", ["kernel-scheduling-syscalls", "native-synchronization-syscalls"], "#9467bd"),
    ("Same Rust driver code, costlier in-process", ["cosmos-routing-session-transport", "native-standard-library", "http2-hpack",
                                                   "resolved-other", "tokio", "http-client", "tls", "shared-key-signing-hmac-headers",
                                                   "residual-unknown"], "#ff7f0e"),
    ("Kernel other / network / IRQ", ["kernel-other", "kernel-network-memory"], "#8c564b"),
    ("cgo / completion FFI", ["cgo-completion-ffi"], "#2ca02c"),
]


def grouped(leg):
    cats = D(leg, "selfCategoryUs")
    covered = {c for _, cs, _ in GROUPS for c in cs}
    missing = set(cats) - covered
    assert not missing, missing
    return [(name, sum(cats.get(c, 0) for c in cs), color) for name, cs, color in GROUPS]


def esc(s):
    return html.escape(str(s))


def waterfall(workers):
    go, drv, sdk = legs[(workers, "go-v2-native")], legs[(workers, "rust-sdk")], legs[(workers, "rust-sdk")]
    g, r = grouped(go), grouped(drv)
    deltas = [(n, gv - rv, c) for (n, gv, c), (_, rv, _) in zip(g, r)]
    total_delta = go["profileCPUUsPerRead"] - drv["profileCPUUsPerRead"]
    assert abs(sum(d for _, d, _ in deltas) - total_delta) < 0.01, (sum(d for _, d, _ in deltas), total_delta)
    scale, x0, w, rowh = 2.2, 300, 960, 30
    rows = [("Rust SDK (baseline)", 0, drv["profileCPUUsPerRead"], "#7f7f7f")]
    run = drv["profileCPUUsPerRead"]
    for n, d, c in deltas:
        rows.append(("+ " + n, run, d, c))
        run += d
    rows.append(("Go v2 over native driver", 0, go["profileCPUUsPerRead"], "#000"))
    h = rowh * len(rows) + 40
    svg = [f'<svg viewBox="0 0 {w} {h}" width="{w}" height="{h}" font-family="Segoe UI,Arial" font-size="13">']
    for i, (label, start, value, color) in enumerate(rows):
        y = 10 + i * rowh
        svg.append(f'<text x="{x0-8}" y="{y+18}" text-anchor="end">{esc(label)}</text>')
        svg.append(f'<rect x="{x0+start*scale:.1f}" y="{y+3}" width="{max(value,0)*scale:.1f}" height="{rowh-8}" fill="{color}"/>')
        tag = f"{value:.1f}" if i in (0, len(rows) - 1) else f"+{value:.1f}"
        svg.append(f'<text x="{x0+(start+value)*scale+6:.1f}" y="{y+18}">{tag} µs</text>')
    svg.append("</svg>")
    note = (f"Direct Rust driver in the same perf pass: {sdk['profileCPUUsPerRead']:.1f} µs/read. "
            f"Delta Go v2 − Rust SDK = {total_delta:+.1f} µs/read ({total_delta/drv['profileCPUUsPerRead']*100:+.0f}%).")
    return "".join(svg), note, deltas


def table(title, rows, cols, fmt=lambda v: f"{v:.2f}"):
    body = "".join("<tr>" + "".join(f"<td>{esc(c) if j == 0 else fmt(c)}</td>" for j, c in enumerate(r)) + "</tr>" for r in rows)
    return f"<h3>{esc(title)}</h3><table><tr>{''.join(f'<th>{esc(c)}</th>' for c in cols)}</tr>{body}</table>"


def compare(key, workers, keys=None, n=14, prefix=""):
    go, drv = D(legs[(workers, "go-v2-native")], key), D(legs[(workers, "rust-sdk")], key)
    names = keys or sorted(set(go) | set(drv), key=lambda k: -abs(go.get(k, 0) - drv.get(k, 0)))[:n]
    return [[prefix + k, go.get(k, 0), drv.get(k, 0), go.get(k, 0) - drv.get(k, 0)] for k in names]


sections = []
summary = {}
for workers in (64, 32):
    svg, note, deltas = waterfall(workers)
    summary[workers] = {n: round(d, 2) for n, d, _ in deltas}
    go = legs[(workers, "go-v2-native")]
    threads = [[k, v] for k, v in go["threadGroupsUs"]]
    futex = [[k, v] for k, v in go["syscallOriginUs"] if k.startswith("futex")][:10]
    drv_futex = [[k, v] for k, v in legs[(workers, "rust-sdk")]["syscallOriginUs"] if k.startswith("futex")][:6]
    alloc = [[k[:120], v] for k, v in go["allocOwnerUs"][:12]]
    go_incl = [[k, v] for k, v in go["inclusiveUs"] if k.startswith(("go", "native: C ABI"))]
    sections.append(f"""
<h2>W{workers}: where the extra CPU per read goes</h2>
<p>{esc(note)}</p>{svg}
<div class="grid">
{table(f"W{workers} Go process: CPU by thread (µs/read)", threads, ["Thread group", "µs/read"])}
{table("Self CPU by language of leaf frame (µs/read)", compare("leafLanguageUs", workers), ["Leaf", "Go v2", "Rust SDK", "Δ"])}
{table("Futex syscall CPU by user-space caller — Go leg", futex, ["futex ← caller", "µs/read"])}
{table("Futex syscall CPU by caller — Rust SDK", drv_futex, ["futex ← caller", "µs/read"])}
{table("Allocation/copy Self CPU by owning frame — Go leg", alloc, ["First non-allocator frame", "µs/read"])}
{table("Inclusive Go-side and C-ABI cost (overlapping; do not sum)", go_incl, ["Frame", "µs/read"])}
{table("Inclusive library cost, Go v2 vs Rust SDK (overlapping)", compare("inclusiveUs", workers, n=12), ["Frame", "Go v2", "Rust SDK", "Δ"])}
</div>""")

page = f"""<!doctype html><meta charset="utf-8"><title>Why Go v2 uses more CPU per read</title>
<style>body{{font-family:Segoe UI,Arial;margin:24px;max-width:1500px}}table{{border-collapse:collapse;font-size:12px;margin-bottom:8px}}
td,th{{border:1px solid #ccc;padding:3px 6px;text-align:right}}td:first-child,th:first-child{{text-align:left;max-width:640px}}
.grid{{display:grid;grid-template-columns:1fr 1fr;gap:16px}}.box{{background:#f6f8fa;border-left:4px solid #0969da;padding:10px 14px}}</style>
<h1>Why Go v2 (over the native driver) uses more CPU per read</h1>
<div class="box"><b>Headline (W64):</b> the Go process spends ~{legs[(64,'go-v2-native')]['profileCPUUsPerRead']:.0f} µs/read versus
~{legs[(64,'rust-sdk')]['profileCPUUsPerRead']:.0f} µs/read for the standalone Rust SDK. Only about a fifth of the gap is Go code itself.
<ol>
<li><b>The Rust driver gets more expensive when hosted in the Go process.</b> Its own threads (<code>azurecosmosdriv</code>) use
{dict(map(tuple,legs[(64,'go-v2-native')]['threadGroupsUs']))['azurecosmosdriv']:.0f} µs/read — more than the entire standalone Rust SDK process.</li>
<li><b>glibc malloc/free is the largest single bucket</b> (+{summary[64]['Native malloc/free/memcpy (Rust code, glibc)']:.0f} µs). The owners are Rust frames (transport send, header maps, session tokens,
completion free), not Go allocation. Futex calls from <code>malloc</code>/<code>_int_free</code>/<code>realloc</code> rise, consistent with arena-lock contention and cross-thread frees
when Go threads free buffers allocated on Tokio threads.</li>
<li><b>Thread hand-off cost</b> (+{summary[64]['Kernel futex / context switch']:.0f} µs kernel): Go runtime <code>notesleep/notewakeup</code>, completion-queue
<code>enqueue</code> wakeups and condvar waits drive futex and context-switch CPU.</li>
<li><b>Go runtime + SDK</b> (+{summary[64]['Go runtime + Go SDK']:.0f} µs): scheduler <code>findRunnable</code>/steal, cgo entry/exit, completion translation; GC and Go allocation are small.</li>
</ol>
Hypothesis to test next (not proven here): co-locating more runnable threads than 4 vCPUs plus glibc arena contention inflates native costs. Candidate experiments: jemalloc/mimalloc for the native library,
fewer driver worker threads, or freeing native buffers on the thread that allocated them.</div>
{''.join(sections)}
<h2>Method and caveats</h2><ul>
<li>Source: existing <code>perf</code> cpu-clock captures (49 Hz, DWARF call chains) of the perf-class runs; no new measurements. Each sample weight = period × profile CPU µs/read ÷ total period.</li>
<li>Every classified Self row ({sum(l['samples'] for l in data['legs'])} samples) was joined to its full raw stack by timestamp; 100% matched. Waterfall buckets are Self categories and sum exactly to the measured delta.</li>
<li>Inclusive tables overlap and must not be summed. Thread, language and category views are alternative partitions of the same total.</li>
<li>Perf runs are slower than clean-score runs (perf overhead); use these for composition, the clean runs for absolute scores (W64: Go v2 289.9 vs Rust SDK 172.0 µs/read).</li>
<li>"Contention" is inferred from futex-from-malloc stacks and thread mix, not from lock-hold measurements. OS <code>perf sched</code> delay data was unusable and is not used.</li>
<li>Input SHA-256 (cpu-deep.json.xz): {hashlib.sha256((out/'cpu-deep.json.xz').read_bytes()).hexdigest()}</li></ul>"""
(out / "cpu-deep-dive.html").write_text(page, encoding="utf-8")
(out / "cpu-deep-summary.json").write_text(json.dumps(summary, indent=2))
print(json.dumps(summary, indent=1))
