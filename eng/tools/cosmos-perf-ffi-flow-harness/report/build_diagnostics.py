"""Render the Go v2 FFI-flow page and the Go pprof/trace tops as explained HTML.

Inputs are the already-exported files in 04-ffi-and-scheduling/*-flow.json and
05-go-diagnostics/w*/*.txt; no raw profiles are needed. Called from build_html.main().
"""
import html
import json
import re

from build_html import ROOT, shell

WORKERS = ("32", "64")
COLORS = ["#1d4e89", "#2e8b57", "#d4801c", "#8e44ad"]
AZ = "github.com/Azure/azure-sdk-for-go/sdk/data/"


def esc(text):
    return html.escape(str(text))


def nice_max(value):
    if value <= 0:
        return 1
    for step in (1, 2, 2.5, 5, 10):
        mag = 10 ** (len(str(int(value))) - 1)
        if value <= step * mag:
            return step * mag
    return value


def line_chart(title, series, y_label, y_max=None, width=900, height=280):
    """series: list of (label, [values], dashed)."""
    left, right, top, bottom = 64, 16, 62, 40
    pw, ph = width - left - right, height - top - bottom
    n = max(len(v) for _, v, _ in series)
    ymax = y_max or nice_max(max(max(v) for _, v, _ in series) * 1.05)
    x = lambda i: left + pw * i / max(n - 1, 1)
    y = lambda v: top + ph - ph * v / ymax
    out = [f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {width} {height}" font-family="Segoe UI,sans-serif" font-size="12">',
           f'<rect width="{width}" height="{height}" fill="#fff"/>',
           f'<text x="{left}" y="18" font-size="14" font-weight="600">{esc(title)}</text>']
    for k in range(6):
        v = ymax * k / 5
        out.append(f'<line x1="{left}" x2="{left + pw}" y1="{y(v):.1f}" y2="{y(v):.1f}" stroke="#e5e7eb"/>')
        out.append(f'<text x="{left - 6}" y="{y(v) + 4:.1f}" text-anchor="end" fill="#57606a">{v:,.0f}</text>')
    for i in range(0, n, 10):
        out.append(f'<text x="{x(i):.1f}" y="{top + ph + 16}" text-anchor="middle" fill="#57606a">{i}</text>')
    out.append(f'<text x="{left + pw / 2}" y="{height - 6}" text-anchor="middle" fill="#57606a">seconds into the measurement window</text>')
    out.append(f'<text transform="translate(14 {top + ph / 2}) rotate(-90)" text-anchor="middle" fill="#57606a">{esc(y_label)}</text>')
    lx, ly = left, 40
    for idx, (label, values, dashed) in enumerate(series):
        c = COLORS[idx % len(COLORS)]
        pts = " ".join(f"{x(i):.1f},{y(v):.1f}" for i, v in enumerate(values))
        dash = ' stroke-dasharray="6 4"' if dashed else ""
        out.append(f'<polyline points="{pts}" fill="none" stroke="{c}" stroke-width="2"{dash}/>')
        w = 22 + 6.8 * len(label) + 20
        if lx + w > width - right and lx > left:
            lx, ly = left, ly + 16
        out.append(f'<line x1="{lx}" x2="{lx + 18}" y1="{ly - 4}" y2="{ly - 4}" stroke="{c}" stroke-width="3"{dash}/>'
                   f'<text x="{lx + 22}" y="{ly}" fill="#1f2328">{esc(label)}</text>')
        lx += w
    out.append("</svg>")
    return "".join(out)


def pct_bars(items, color="#1d4e89", width=360):
    top = max((v for _, v, _ in items), default=1) or 1
    return "".join(
        f'<div class=bar><span class=k>{esc(k)}</span><span class=t style="width:{max(v / top * width, 1):.0f}px;background:{color}"></span>'
        f'<span class=v style="width:auto">{esc(t)}</span></div>'
        for k, v, t in items
    )


def explain(measures, shows, means):
    return (f'<div class=measure><p><strong>What it measures:</strong> {measures}</p>'
            f'<p><strong>What it shows:</strong> {shows}</p><p><strong>What it means:</strong> {means}</p></div>')


# ---------------------------------------------------------------- FFI flow page

def build_flow_page():
    folder = ROOT / "04-ffi-and-scheduling"
    data = {w: json.loads((folder / f"w{w}-flow.json").read_text(encoding="utf-8")) for w in WORKERS}
    cards, sections = [], []
    for w in WORKERS:
        d = data[w]
        b = d["buckets"]
        rate = d["success"] / d["flowWallSeconds"]
        latency_ms = int(w) / rate * 1000
        sub = [x["goRates"]["submitAccepted"] for x in b]
        term = [x["nativeRates"]["terminal"] for x in b]
        comp = [x["goRates"]["delivered"] for x in b]
        res = [x["goRates"]["resumed"] for x in b]
        gap = [abs(s - c) for s, c in zip(sub, comp)]
        inner = b[1:-1]
        gap_inner = [abs(x["goRates"]["submitAccepted"] - x["goRates"]["delivered"]) for x in inner]
        hist = {int(k): v for k, v in d["batchHistogram"].items()}
        batches = sum(hist.values())
        ones, twos = hist.get(1, 0), hist.get(2, 0)
        ready = [x["nativeReadyHigh"] for x in b]
        admitted = [x["nativeAdmittedUndrainedHigh"] for x in b]
        polled = [x["nativePolledUnfinishedHigh"] for x in b]
        ready_sorted = sorted(ready)
        ready_median = ready_sorted[len(ready_sorted) // 2]
        tot = d["goTotals"]
        cards.append(
            f'<div class=card><div class=v>W{w}: {rate:,.0f}/s</div><div class=l>{d["success"]:,} reads in {d["flowWallSeconds"]:.0f} s; '
            f'submitted = completed = resumed; 0 errors, 0 rejected</div></div>')
        files = folder / f"w{w}-flow-throughput.svg", folder / f"w{w}-flow-gauges.svg"
        files[0].write_text(line_chart(f"W{w}: per-second FFI flow (all four counters)", [
            ("Go submissions/s", sub, False), ("Rust terminal/s", term, True),
            ("Go completions/s", comp, False), ("Go waiter resumes/s", res, True)], "operations per second",
            y_max=nice_max(max(sub) * 1.15)), encoding="utf-8")
        files[1].write_text(line_chart(f"W{w}: native queue and in-flight gauges (per-second high-water)", [
            ("In flight inside Rust (admitted, not drained)", admitted, False),
            ("Polled by Tokio, not finished", polled, True),
            ("Ready queue waiting for Go", ready, False)], "operations",
            y_max=nice_max(int(w) * 1.15)), encoding="utf-8")
        hist_items = [(f"batch of {k}", v, f"{v:,} ({v / batches:.1%})") for k, v in sorted(hist.items())[:6]]
        rest = batches - sum(v for _, v, _ in hist_items)
        if rest:
            hist_items.append(("larger batches", rest, f"{rest:,} ({rest / batches:.1%})"))
        sections.append(f"""
<h2 id=w{w}>W{w} ({w} closed-loop workers)</h2>
<div class=charts>
<figure class=wide><img src="{files[0].name}" alt="W{w} flow rates"><figcaption>Lines overlap because the four counters are equal each second.
Largest per-second submit/complete gap inside the window: {max(gap_inner):,.0f} ops (first/last partial seconds: {max(gap):,.0f}).</figcaption></figure>
</div>
{explain(
    "Four counters sampled per second: Go handed an operation to the driver (<em>submission</em>), Rust finished it "
    "(<em>terminal</em>), Go's reactor drained it from the C completion queue (<em>completion</em>), and the waiting worker "
    "goroutine woke up (<em>resume</em>).",
    f"All four lines sit on top of each other at about {rate:,.0f}/s for the whole minute. Totals match exactly: "
    f"{tot['submitAccepted']:,} submitted, {tot['delivered']:,} delivered, {tot['resumed']:,} resumed, {tot['submitErrors']} submit errors.",
    "The FFI handoff loses nothing, retries nothing and doesn't back up. Throughput is flat, so there is no warm-up drift, "
    "GC stall or periodic hiccup. Whatever limits throughput, it isn't requests piling up at the Go/Rust boundary.")}
<div class=charts><figure class=wide><img src="{files[1].name}" alt="W{w} queue gauges"><figcaption>High-water values per second.
Ready queue: median {ready_median}, peak {max(ready)}.</figcaption></figure></div>
{explain(
    "Inside the native driver, per second: the most operations admitted but not yet handed back to Go (in flight), "
    "the most already polled by a Tokio worker but not finished, and the most finished results waiting in the ready "
    "queue for Go to pick up.",
    f"In-flight and polled-unfinished are pinned at {max(admitted)} = the worker count every second. The ready queue is "
    f"usually {ready_median} or fewer, with occasional spikes to {max(ready)}.",
    f"Nearly every worker's read is inside Rust at any moment; Go has almost nothing outstanding on its side. The ready "
    f"queue stays near empty, so Go collects results as fast as Rust produces them. The spikes are brief bursts that "
    f"clear within the same second. By Little's law the average read takes about {latency_ms:.1f} ms "
    f"({w} ÷ {rate:,.0f}/s), almost all of it in the driver, the network and the service.")}
<h3>Completion batch sizes (W{w})</h3>
<div class=measure>{pct_bars(hist_items, "#2e8b57")}</div>
{explain(
    "Each time Go's reactor wakes, it drains everything in the completion queue as one batch. The histogram counts batch sizes.",
    f"{ones / batches:.0%} of wakeups find exactly one completion and {twos / batches:.0%} find two; mean "
    f"{d['meanCompletionBatchSize']:.2f}. That's {d['completionBatchesPerSecond']:,.0f} wakeups/s for {d['completionsPerSecond']:,.0f} completions/s.",
    "Batch size is not queue depth: small batches mean Go drains quickly, not that work is scarce. The cost is one "
    "reactor wakeup and cgo crossing per read. Batching more would save a little Go CPU but add latency; this isn't the bottleneck.")}
<h3>Handoff latency (W{w})</h3>
{explain(
    "Time from Rust publishing a completion to the waiting goroutine running again (publication plus Go scheduling).",
    f"Mean {d['deliveryToResumeMeanUs']:.1f} µs, max {d['deliveryToResumeMaxUs'] / 1000:.1f} ms.",
    f"About {d['deliveryToResumeMeanUs'] / (latency_ms * 1000):.1%} of the ~{latency_ms:.1f} ms average read. It doesn't grow from "
    "W32 to W64, so the Go scheduler isn't what limits throughput.")}
""")
    body = f"""
<h1>Go v2 FFI-flow diagnostics</h1>
<p>How work crosses between Go and the Rust driver: submitted from Go, executed by Rust's Tokio threads, published to
a C completion queue, drained by one Go reactor goroutine, then handed to the waiting worker. These are separate
instrumented runs (60 s measurement), not authoritative clean scores.</p>
<div class=cards>{''.join(cards)}
<div class=card><div class=v>~22 µs</div><div class=l>Mean Rust-to-Go handoff at both W32 and W64</div></div></div>
<div class=note><strong>Bottom line.</strong> The boundary is healthy: nothing is lost or queued, and Go picks up results
in ~22 µs. Doubling workers from 32 to 64 adds only ~3% throughput because every extra request just waits inside the
driver, whose two Tokio threads are saturated (<a href="THREAD-SATURATION.html">Thread saturation</a>). Work on making
the driver cheaper per read, not on the handoff.</div>
<p><a href="#w32">W32</a> · <a href="#w64">W64</a> · Data: <a href="w32-flow.json">w32-flow.json</a>, <a href="w64-flow.json">w64-flow.json</a> · Table: <a href="report.html">report</a></p>
{''.join(sections)}
"""
    page = "04-ffi-and-scheduling/dashboard.html"
    (ROOT / page).write_text(shell(page, "Go v2 FFI-flow diagnostics", body), encoding="utf-8")


# ---------------------------------------------------------------- Go pprof page

ROW = re.compile(r"^\s*(\S+)\s+([\d.]+)%\s+[\d.]+%\s+(\S+)\s+([\d.]+)%\s+(.+?)\s*$")


def parse_top(path):
    text = path.read_text(encoding="utf-8")
    header = {}
    for line in text.splitlines():
        m = re.match(r"^(Type|Duration|Time): (.*)$", line)
        if m:
            header[m.group(1)] = m.group(2)
    total = re.search(r"of ([\d.]+\w*) total", text)
    rows = [ROW.match(l) for l in text.splitlines()]
    rows = [(m.group(5), m.group(1), float(m.group(2)), m.group(3), float(m.group(4))) for m in rows if m]
    return header, total.group(1) if total else "0", rows


def short(name):
    return name.replace(AZ, "").replace("azcosmos/v2.", "v2.")


def top_table(path, mode):
    header, total, rows = parse_top(path)
    if not rows:
        return f"<p><em>No samples (total {esc(total)}).</em></p>"
    key = 4 if mode == "cum" else 2
    rows = sorted(rows, key=lambda r: -r[key])[:12]
    trs = "".join(
        f'<tr><td title="{esc(n)}"><code>{esc(short(n)[:78])}</code></td><td class=n>{esc(f)}</td>'
        f'<td><span class=t style="display:inline-block;height:12px;width:{fp * 1.4:.0f}px;background:#1d4e89"></span> {fp:.1f}%</td>'
        f'<td class=n>{esc(c)}</td><td><span class=t style="display:inline-block;height:12px;width:{cp * 1.4:.0f}px;background:#9ecbff"></span> {cp:.1f}%</td></tr>'
        for n, f, fp, c, cp in rows)
    meta = " · ".join(esc(f"{k}: {v}") for k, v in header.items())
    return (f'<p class=meta style="font-size:12.5px;color:#57606a">{meta} · total {esc(total)} · sorted by {"cumulative" if mode == "cum" else "self (flat)"} · '
            f'<a href="{path.parent.name}/{path.name}">raw text</a></p>'
            f'<table><tr><th>Function</th><th class=n>flat</th><th>flat %</th><th class=n>cum</th><th>cum %</th></tr>{trs}</table>')


def seconds(value):
    m = re.match(r"([\d.]+)(us|ms|s)$", value)
    if not m:
        return 0.0
    return float(m.group(1)) * {"us": 1e-6, "ms": 1e-3, "s": 1.0}[m.group(2)]


def sched_rows(path):
    keep = []
    for line in path.read_text(encoding="utf-8").splitlines():
        m = re.match(r"\s*(go-v2-native|azurecosmosdriv):\((\d+)\)\s*\|\s*([\d.]+) ms\s*\|\s*(\d+)\s*\|\s*avg:\s*([\d.]+) ms\s*\|\s*max:\s*([\d.]+) ms", line)
        if m:
            keep.append(m.groups())
    bug = re.search(r"([\d.]+)% context switch bugs", path.read_text(encoding="utf-8"))
    return keep, bug.group(1) if bug else "?"


PROFILES = [
    ("cpu-top.txt", "flat", "CPU (Go pprof, 15 s)",
     "Where on-CPU time was sampled by Go's profiler during a 15 s window. <code>[libc.so.6]</code> and <code>[binary]</code> are "
     "C/Rust frames Go cannot symbolize.",
     "About 1.95 cores' worth of samples. Roughly half lands in unsymbolized libc and Rust code, ~12–15% in "
     "<code>runtime.cgocall</code> (crossing into C), and only a few percent in Go SDK functions such as <code>executeItem</code> "
     "and <code>translateCompletion</code>.",
     "Go-side SDK code is a small slice of CPU; most cost is inside the native driver, which pprof can't break down. "
     "Use the perf-based pages (<a href=\"../03-rust-driver-inside-go-v2/HOTSPOTS.html\">Rust hotspots</a>) for that. pprof percentages are a "
     "different denominator from perf and process CPU; don't add them."),
    ("heap-alloc-bytes-top.txt", "flat", "Go heap: bytes allocated",
     "Total bytes Go allocated since process start (includes warmup). Go heap only, not Rust/C allocations.",
     "<code>_Cfunc_GoBytes</code> (copying each response body from C into Go) is ~31%. Completion handling "
     "(<code>awaitCompletion</code>, <code>moveCompletion</code>) and diagnostics copies (<code>newDiagnosticsJSON</code>, "
     "<code>copyDiagnostics*</code>) are most of the rest; request building (<code>buildNativeItemRequest</code>, "
     "<code>toNative</code>) is smaller.",
     "Every read copies its body across the boundary and builds a Go diagnostics object. Avoiding the body copy and "
     "building diagnostics lazily are the main Go-side allocation savings (proposal C-items)."),
    ("heap-alloc-objects-top.txt", "flat", "Go heap: objects allocated",
     "Count of Go allocations (same window). High counts drive GC work even when bytes are small.",
     "Spread across <code>GoString</code>, <code>OperationOptions.toNative</code>, <code>awaitCompletion</code>, "
     "<code>buildNativeItemRequest</code>, the benchmark's own <code>itemID</code>, partition key conversion and diagnostics.",
     "Many small per-read objects: option and partition key conversion each read could be cached. Note <code>main.itemID</code> is the benchmark, not the SDK."),
    ("heap-live-top.txt", "flat", "Go heap: live memory at end",
     "Memory still reachable when the profile was taken.",
     "About 6 MB, mostly benchmark latency buffers, the profiler itself and TLS/HTTP setup.",
     "No Go-side leak or growth. Process RSS is dominated by the native side and runtime, not the Go heap."),
    ("block-top.txt", "cum", "Blocking (goroutine wait time)",
     "Total time goroutines spent blocked on channels, selects and wait groups, summed over all goroutines. Wait time, not CPU.",
     "~95% is <code>selectgo</code> inside <code>awaitCompletion</code>: each worker waiting for its read to come back from Rust. "
     "Total ÷ 15 s ≈ the number of workers, i.e. every worker is waiting almost all the time.",
     "Expected for a closed-loop benchmark and consistent with the flow page: requests sit in the driver, not in Go. Not a hotspot."),
    ("mutex-top.txt", "cum", "Mutex contention",
     "Time goroutines waited to acquire Go mutexes.",
     "About 1 ms total over 15 s, mostly runtime-internal locks.",
     "No meaningful Go lock contention, including the container-handle cache and pending-operation map."),
    ("trace-sync-top.txt", "cum", "Runtime trace: sync blocking (5 s)",
     "Same idea as the block profile, from the 5 s runtime trace.",
     "Again ~97% is workers waiting in <code>awaitCompletion</code>; total ÷ 5 s ≈ the worker count.",
     "Confirms the block profile with an independent capture."),
    ("trace-syscall-top.txt", "flat", "Runtime trace: time in syscalls/cgo (5 s)",
     "Time goroutines spent inside syscalls or cgo calls.",
     "Most is <code>cosmos_completion_queue_wait</code>, the reactor goroutine sleeping in C until Rust publishes a result. "
     "Submit and free calls are each well under a second total.",
     "The reactor is idle most of the time, so draining completions isn't a bottleneck; each cgo submit/free is cheap."),
    ("trace-sched-top.txt", "cum", "Runtime trace: scheduler latency (5 s)",
     "Time goroutines were runnable but waiting for a Go processor (P), attributed to the code that woke them.",
     "Under 1 s total over 5 s; ~96% is attributed to <code>chansend1</code> in <code>reactor.deliver</code>, the point where "
     "the reactor wakes a worker. Spread over ~40k reads, this matches the ~22 µs handoff on the flow page.",
     "The wakeup site, not a blocking send. The Go scheduler adds ~22 µs per read and doesn't worsen from W32 to W64."),
    ("trace-net-top.txt", "flat", "Runtime trace: network blocking",
     "Time blocked in Go's network poller.",
     "Zero: all network I/O happens in Rust, not in Go.",
     "Expected for Go v2; Go never touches the socket."),
]


def build_go_diag_page():
    folder = ROOT / "05-go-diagnostics"
    parts, toc = [], []
    for name, mode, title, measures, shows, means in PROFILES:
        anchor = name.replace(".txt", "")
        toc.append(f'<a href="#{anchor}">{esc(title)}</a>')
        cols = "".join(f"<figure><h3 style='margin:4px 0'>W{w}</h3>{top_table(folder / f'w{w}' / name, mode)}</figure>"
                       for w in WORKERS)
        parts.append(f'<h2 id={anchor}>{esc(title)}</h2>{explain(measures, shows, means)}<div class=charts style="grid-template-columns:1fr">{cols}</div>')
    sched = []
    for w in WORKERS:
        rows, bug = sched_rows(folder / f"w{w}" / "os-scheduler-latency.txt")
        trs = "".join(f"<tr><td><code>{esc(t)}</code></td><td class=n>{n}</td><td class=n>{float(r):,.0f}</td><td class=n>{int(s):,}</td>"
                      f"<td class=n>{float(a) * 1000:.0f}</td><td class=n>{float(mx):.1f}</td></tr>" for t, n, r, s, a, mx in rows)
        sched.append(f"<figure><h3 style='margin:4px 0'>W{w}</h3><table><tr><th>Thread group</th><th class=n>Threads</th><th class=n>Run ms</th>"
                     f"<th class=n>Switches</th><th class=n>Avg delay µs</th><th class=n>Max delay ms</th></tr>{trs}</table>"
                     f"<p class=meta style='font-size:12.5px;color:#57606a'>perf reported {bug}% context-switch inconsistencies · "
                     f"<a href='w{w}/os-scheduler-latency.txt'>raw text</a></p></figure>")
    parts.append('<h2 id=os-sched>OS scheduler latency (perf sched, PID-attached)</h2>' + explain(
        "Linux <code>perf sched latency</code> for the benchmark process: run time, context switches and wait-to-run delay per thread group.",
        "The two driver threads (<code>azurecosmosdriv</code>) do most of the running with far fewer switches than Go's threads. "
        "But perf flagged ~90% of context switches as inconsistent.",
        "Treat this as qualitative only. The trustworthy per-thread numbers come from schedstat on the "
        "<a href=\"../04-ffi-and-scheduling/THREAD-SATURATION.html\">Thread saturation</a> page.")
        + f'<div class=charts style="grid-template-columns:1fr">{"".join(sched)}</div>')
    other = sorted(p.relative_to(folder).as_posix() for p in folder.rglob("*.json"))
    body = f"""
<h1>Go diagnostics (pprof and runtime trace)</h1>
<p>Go-side profiles from a separate diagnostic run (15 s pprof window, 5 s runtime trace) at W32 and W64. Rendered from
the exported <code>go tool pprof -top</code> text; raw <code>.pprof</code> files stay on the VM. Not authoritative scores,
and Go pprof can't see inside the Rust driver.</p>
<div class=note><strong>Bottom line.</strong> Go-side code is a small share of CPU and has no lock contention, leaks or scheduler
problem. Its measurable costs are copying each response body out of C and building diagnostics objects per read. Goroutines
spend their time waiting for the Rust driver, which is where the CPU goes.</div>
<p>{' · '.join(toc)} · <a href="#os-sched">OS scheduler latency</a></p>
{''.join(parts)}
<h2>Other files</h2><p>{' · '.join(f'<a href="{esc(p)}">{esc(p)}</a>' for p in other)}</p>
"""
    page = "05-go-diagnostics/index.html"
    (ROOT / page).write_text(shell(page, "Go diagnostics (pprof and runtime trace)", body), encoding="utf-8")


def main():
    build_flow_page()
    build_go_diag_page()
