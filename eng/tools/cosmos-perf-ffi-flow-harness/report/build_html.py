"""Build the navigable HTML view of the team report (no third-party dependencies).

Run from anywhere: python build_html.py
Writes index.html at the report root and an .html sibling for every report .md file.
"""
import html
import json
import pathlib
import re
import sys

ROOT = pathlib.Path(__file__).resolve().parents[2]

PAGES = [
    ("Summary", "index.html"),
    ("Overview (README)", "README.html"),
    ("1 · Results dashboard", "01-results/dashboard.html"),
    ("1 · Results table", "01-results/report.html"),
    ("2 · Go v2 vs Rust CPU", "02-why-go-v2-costs-more-cpu/cpu-deep-dive.html"),
    ("3 · Rust hotspots", "03-rust-driver-inside-go-v2/HOTSPOTS.html"),
    ("3 · Rust-in-Go layers", "03-rust-driver-inside-go-v2/rust-in-go.html"),
    ("4 · Thread saturation", "04-ffi-and-scheduling/THREAD-SATURATION.html"),
    ("4 · Scheduler analysis", "04-ffi-and-scheduling/SCHEDULER-ANALYSIS.html"),
    ("4 · Completion-queue flow", "04-ffi-and-scheduling/dashboard.html"),
    ("5 · Go diagnostics", "05-go-diagnostics/index.html"),
    ("6 · Optimization proposals", "06-optimizations/PROPOSALS.html"),
    ("Methods", "methods/METHODS.html"),
    ("Provenance", "provenance/index.html"),
]
MD_FILES = [
    "README.md",
    "01-results/report.md",
    "03-rust-driver-inside-go-v2/HOTSPOTS.md",
    "04-ffi-and-scheduling/THREAD-SATURATION.md",
    "04-ffi-and-scheduling/SCHEDULER-ANALYSIS.md",
    "04-ffi-and-scheduling/report.md",
    "06-optimizations/PROPOSALS.md",
    "methods/METHODS.md",
]

CSS = """
*{box-sizing:border-box}body{margin:0;font:15px/1.55 -apple-system,Segoe UI,Roboto,sans-serif;color:#1f2328;background:#f6f8fa}
nav{position:fixed;top:0;left:0;bottom:0;width:250px;background:#0b2545;color:#fff;padding:18px 12px;overflow:auto}
nav h2{font-size:15px;margin:0 0 12px;color:#9ecbff}nav a{display:block;color:#dbe9ff;text-decoration:none;padding:5px 8px;border-radius:6px;font-size:13.5px}
nav a:hover{background:#13315c}nav a.cur{background:#1d4e89;color:#fff;font-weight:600}
main{margin-left:250px;padding:28px 40px;max-width:1250px}
h1{margin-top:0;font-size:26px}h2{border-bottom:1px solid #d0d7de;padding-bottom:4px;margin-top:32px}
table{border-collapse:collapse;margin:12px 0;background:#fff;font-size:13.5px}th,td{border:1px solid #d0d7de;padding:5px 9px;vertical-align:top}
th{background:#eef2f6}td.n,th.n{text-align:right;font-variant-numeric:tabular-nums}code{background:#eaeef2;padding:1px 4px;border-radius:4px;font-size:12.5px}
.cards{display:grid;grid-template-columns:repeat(auto-fill,minmax(230px,1fr));gap:14px;margin:16px 0}
.card{background:#fff;border:1px solid #d0d7de;border-radius:10px;padding:14px 16px}.card .v{font-size:26px;font-weight:700;color:#0b2545}.card .l{font-size:13px;color:#57606a}
.charts{display:grid;grid-template-columns:repeat(auto-fill,minmax(min(420px,100%),1fr));gap:14px}.charts figure{margin:0;background:#fff;border:1px solid #d0d7de;border-radius:10px;padding:10px}
@media (max-width:900px){nav{position:static;width:auto;max-height:none}main{margin-left:0;padding:16px}}
.charts img{width:100%}.charts figure.wide{grid-column:1/-1}.charts figure.wide img{max-width:1080px;display:block;margin:auto}
.measure{background:#fff;border:1px solid #d0d7de;border-radius:10px;padding:4px 18px;margin:12px 0}figcaption{font-size:13px;color:#57606a;margin-top:4px}
.bar{display:flex;align-items:center;gap:8px;margin:3px 0;font-size:13px}.bar .k{width:230px}.bar .t{height:16px;border-radius:3px}.bar .v{width:60px;text-align:right;font-variant-numeric:tabular-nums}
.note{background:#fff8c5;border:1px solid #d4a72c;border-radius:8px;padding:10px 14px;font-size:13.5px}
.links a{display:inline-block;margin:4px 8px 4px 0;padding:6px 12px;background:#fff;border:1px solid #d0d7de;border-radius:8px;text-decoration:none}
iframe{width:100%;height:85vh;border:1px solid #d0d7de;border-radius:8px;background:#fff}
"""


def rel(from_page, to_page):
    depth = from_page.count("/")
    return "../" * depth + to_page


def shell(page, title, body):
    nav = "".join(
        f'<a href="{rel(page, href)}"{" class=cur" if href == page else ""}>{html.escape(label)}</a>'
        for label, href in PAGES
    )
    return (
        f"<!doctype html><html lang=en><head><meta charset=utf-8><title>{html.escape(title)}</title>"
        f"<meta name=viewport content='width=device-width,initial-scale=1'><style>{CSS}</style></head>"
        f"<body><nav><h2>Go v2 CPU analysis</h2>{nav}</nav><main>{body}</main></body></html>"
    )


def inline(text):
    parts = re.split(r"(`[^`]+`)", text)
    out = []
    for part in parts:
        if part.startswith("`") and part.endswith("`") and len(part) > 1:
            out.append(f"<code>{html.escape(part[1:-1])}</code>")
            continue
        s = html.escape(part, quote=False)
        s = re.sub(r"\*\*(.+?)\*\*", r"<strong>\1</strong>", s)

        def link(m):
            target = m.group(2)
            if not re.match(r"[a-z]+:", target):
                target = re.sub(r"\.md(#|$)", r".html\1", target)
                if target.endswith("/"):
                    folder = target.rstrip("/").split("/")[-1]
                    if folder in FOLDER_PAGE:
                        target += FOLDER_PAGE[folder]
            return f'<a href="{html.escape(target)}">{m.group(1)}</a>'

        s = re.sub(r"\[([^\]]+)\]\(([^)]+)\)", link, s)
        out.append(s)
    return "".join(out)


FOLDER_PAGE = {}
for _label, _href in PAGES:
    if "/" in _href:
        FOLDER_PAGE.setdefault(_href.split("/")[0], _href.split("/", 1)[1])
FOLDER_PAGE["scripts"] = "build_html.py"


def is_num(cell):
    return bool(re.fullmatch(r"[\*~≈–\-+]*[\d.,]+[\d.,%–\- ]*\**|–", cell.strip()))


def md_to_html(md):
    lines = md.splitlines()
    out, i = [], 0
    while i < len(lines):
        line = lines[i]
        if not line.strip():
            i += 1
            continue
        m = re.match(r"^(#{1,4}) (.*)", line)
        if m:
            n = len(m.group(1))
            out.append(f"<h{n}>{inline(m.group(2))}</h{n}>")
            i += 1
            continue
        if line.startswith("```"):
            j = i + 1
            while j < len(lines) and not lines[j].startswith("```"):
                j += 1
            out.append("<pre><code>" + html.escape("\n".join(lines[i + 1:j])) + "</code></pre>")
            i = j + 1
            continue
        if line.lstrip().startswith("|"):
            rows = []
            while i < len(lines) and lines[i].lstrip().startswith("|"):
                rows.append([c.strip() for c in re.split(r"(?<!\\)\|", lines[i].strip())[1:-1]])
                i += 1
            head, body = rows[0], [r for r in rows[1:] if not all(re.fullmatch(r":?-+:?", c) for c in r)]
            t = "<table><tr>" + "".join(f"<th>{inline(c)}</th>" for c in head) + "</tr>"
            for r in body:
                t += "<tr>" + "".join(
                    f'<td{" class=n" if is_num(c) else ""}>{inline(c.replace(chr(92) + "|", "|"))}</td>' for c in r
                ) + "</tr>"
            out.append(t + "</table>")
            continue
        m = re.match(r"^(\s*)([-*]|\d+\.) (.*)", line)
        if m:
            tag = "ol" if m.group(2)[0].isdigit() else "ul"
            items = []
            while i < len(lines):
                m = re.match(r"^(\s*)([-*]|\d+\.) (.*)", lines[i])
                if m and len(m.group(1)) == 0:
                    items.append(m.group(3))
                elif lines[i].startswith("  ") and lines[i].strip() and items:
                    items[-1] += " " + lines[i].strip()
                else:
                    break
                i += 1
            out.append(f"<{tag}>" + "".join(f"<li>{inline(x)}</li>" for x in items) + f"</{tag}>")
            continue
        para = [line]
        i += 1
        while i < len(lines) and lines[i].strip() and not re.match(r"^(#|\||\s*[-*] |\s*\d+\. |```)", lines[i]):
            para.append(lines[i])
            i += 1
        out.append("<p>" + inline(" ".join(p.strip() for p in para)) + "</p>")
    return "\n".join(out)


def file_list_page(page, title, folder, intro):
    rows = "".join(
        f'<tr><td><a href="{p.relative_to(folder).as_posix()}">{html.escape(p.relative_to(folder).as_posix())}</a></td>'
        f"<td class=n>{p.stat().st_size:,}</td></tr>"
        for p in sorted(folder.rglob("*"))
        if p.is_file() and p.name != "index.html"
    )
    body = f"<h1>{html.escape(title)}</h1><p>{intro}</p><table><tr><th>File</th><th class=n>Bytes</th></tr>{rows}</table>"
    (ROOT / page).write_text(shell(page, title, body), encoding="utf-8")


def bars(items, color, unit="µs"):
    top = max(v for _, v in items)
    return "".join(
        f'<div class=bar><span class=k>{html.escape(k)}</span><span class=t style="width:{v / top * 420:.0f}px;background:{color}"></span>'
        f"<span class=v>{v:.1f} {unit}</span></div>"
        for k, v in items
    )


def build_index():
    rig = json.loads((ROOT / "03-rust-driver-inside-go-v2/rust-in-go-summary.json").read_text(encoding="utf-8"))["64"]
    groups = sorted(rig["groups"].items(), key=lambda x: -x[1])
    layers = sorted(rig["rustLayers"].items(), key=lambda x: -x[1])
    themes = [
        ("glibc malloc/free/memcpy", 87), ("HTTP stack (h2+hyper+reqwest+TLS)", 69), ("Kernel wake/sleep", 26),
        ("Response/request headers", 17), ("Operation orchestration", 16), ("Futex inside allocator", 15),
        ("Session token", 14), ("Diagnostics", 8.6), ("Request signing", 6),
    ]
    charts = "".join(
        f'<figure{" class=wide" if "flywheel" in f else ""}><img src="{f if "/" in f else "01-results/" + f}" alt="{c}"><figcaption>{c}</figcaption></figure>'
        for f, c in [
            ("throughput.svg", "Throughput (reads/s)"), ("cpu-per-read.svg", "CPU µs per read"),
            ("latency-p99.svg", "p99 latency"), ("memory-rss.svg", "RSS"),
            ("04-ffi-and-scheduling/w64-flow-throughput.svg",
             "Go v2 FFI flow, W64: submissions, Rust completions, Go deliveries and resumes are equal every second"),
            ("w64-go-v2-native-cpu-flywheel.svg",
             "Go v2 CPU flywheel, W64 (perf run). The ring is all 4 vCPUs; coloured slices are this process, "
             "grey is idle or other processes. Only the black slice is truly unattributed."),
        ]
    )
    measured = """
<h2>How CPU and thread utilization were measured</h2>
<div class=measure>
<p>Three independent measurements, each from its own run class. Their denominators are never added together.</p>
<ol>
<li><strong>Process CPU (clean runs, authoritative).</strong> Each adapter reads its own user + system CPU at
the start and end of the 60 s measurement window (10 s warmup excluded): Go via <code>getrusage(RUSAGE_SELF)</code>,
Rust via <code>/proc/self/stat</code>. For Go v2 this covers every thread in the process, including the Rust
driver's Tokio threads. <em>CPU µs/read</em> = CPU seconds ÷ successful reads; <em>busy cores</em> = CPU seconds ÷
wall seconds. No profiler runs in this class.</li>
<li><strong>Per-thread utilization (same clean runs).</strong> The launcher snapshots
<code>/proc/&lt;pid&gt;/task/&lt;tid&gt;/schedstat</code> (run ns, run-queue wait ns, timeslices) and <code>status</code>
(thread name, voluntary/involuntary switches) at measurement start and end. Threads named
<code>azurecosmosdriv</code> are the driver's Tokio workers; all others are Go. <em>Thread busy</em> = run ns ÷ wall ns
(≈ 0.93 per driver thread); run-queue wait ≈ 1 µs/read shows they were not starved of CPU.</li>
<li><strong>Where the CPU goes (separate perf run).</strong> Linux <code>perf</code> sampled the Go v2 process at
49 Hz with call chains. Each sample's share is scaled to that run's measured CPU µs/read, then split by thread name
and by the innermost matching stack frame (layer). perf adds overhead, so its absolute throughput is lower; use it for
proportions only.</li>
</ol>
<p>Caveats: one perf repeat; schedstat gives window averages, not peaks. Details: <a href="methods/METHODS.html">Methods</a>.</p>
</div>"""
    body = f"""
<h1>Go v2 + native Rust driver: why it costs more CPU per read</h1>
<p>Campaign <code>go-v2-ffi-flow-latest-main-20261105-a</code> · 4-vCPU Linux VM · Gateway V1 HTTP/2, binary encoding off ·
1,024 × 1,024-byte documents · point reads at 32 and 64 closed-loop workers · 2 clean repeats, 1 perf repeat.</p>
<h2 id=quick-summary>Quick summary</h2>
<p>Numbers are from W64, Go v2 process only. The CPU breakdown comes from one perf run, so read it as directional.</p>
<h3>Top bottlenecks</h3>
<ol>
<li><strong>The two Rust driver threads set the ceiling.</strong> They run about 93% busy and cap Go v2 at roughly 8.3k reads/s. Going from 32 to 64 workers adds about 3% throughput but roughly doubles latency (about 4 ms to 7.8 ms).</li>
<li><strong>Most of the CPU is inside the Rust driver.</strong> Of 294 µs per read, about 228 µs runs on the driver threads and about 66 µs on Go threads.</li>
<li><strong>Memory allocation and copying is the biggest single cost.</strong> glibc malloc, free and memcpy take about 87 µs per read. Another 15 µs or so is futex waiting under malloc and free, which suggests allocator lock contention (inferred, not measured).</li>
<li><strong>The HTTP stack costs about 69 µs per read.</strong> This covers HTTP/2 header encoding (HPACK), header maps, URL parsing, connection-pool checkout and TLS.</li>
<li><strong>Per-read Cosmos bookkeeping adds up to about 55 µs.</strong> Response-header parsing about 17, orchestration and hedging about 16, session tokens about 14, diagnostics about 9. Much of this is repeated work that could be cached or built only when needed.</li>
<li><strong>On the Go side,</strong> copying each response body out of C is about 31% of Go allocations, and a diagnostics object is built for every read.</li>
</ol>
<h3>FFI pressure (Go ↔ Rust boundary)</h3>
<p>Measured with flow counters and gauges in Go v2 and the driver, in separate 60-second instrumented runs (not clean scores).</p>
<ul>
<li><strong>Nothing is lost or backed up.</strong> Submitted, finished in Rust, drained by Go and resumed counts are equal every second (494,550 each at W64). The ready queue stays near empty (max 43).</li>
<li><strong>Workers wait on the driver, not on Go.</strong> In-flight reads always equal the worker count, which is expected for a closed-loop test. What matters is where they sit: almost all are inside Rust, not waiting for Go to pick them up. Go's block profile agrees: about 95% of wait time is workers parked in <code>awaitCompletion</code>. That wait includes the network round trip and queueing behind the busy driver threads, not only Rust CPU.</li>
<li><strong>The handoff is fast.</strong> From result delivery to the worker running again takes about 22 µs, roughly 0.3–0.5% of read latency. Go drains about one completion per wakeup (mean 1.08), so it keeps up.</li>
<li><strong>The boundary still costs some CPU.</strong> The C interface, cgo calls and completion handling take about 22 µs per read (about 8% of Go v2 CPU), plus the body copy above. That's worth trimming for efficiency, but it isn't what limits throughput.</li>
</ul>
<p><strong>Candidate fixes to A/B test:</strong> a different allocator for the driver (mimalloc or jemalloc), 3 Tokio threads instead of 2, fewer header and URL allocations, caching session tokens and options, building diagnostics lazily, and avoiding the Go body copy. <a href="06-optimizations/PROPOSALS.html">Proposals</a></p>
<h2>Key numbers</h2>
<div class=cards>
<div class=card><div class=v>290 µs</div><div class=l>Go v2 CPU per read at W64 (Rust SDK 172 µs, direct Rust driver 172 µs)</div></div>
<div class=card><div class=v>8.3k/s</div><div class=l>Go v2 reads/s at W64, vs 10.6k for Rust SDK; W32 → W64 adds only 2.3%</div></div>
<div class=card><div class=v>2 × 93%</div><div class=l>Both driver Tokio threads saturated → ~8k reads/s ceiling while the VM uses ~2.4 of 4 cores</div></div>
<div class=card><div class=v>228 / 294 µs</div><div class=l>Share of Go v2 CPU (perf, W64) spent on the Rust driver threads</div></div>
<div class=card><div class=v>~87 µs</div><div class=l>glibc malloc/free/memcpy inside the driver per read, plus ~15 µs futex waits under the allocator</div></div>
</div>
<h2>More detail</h2>
<ol>
<li><strong>The ceiling is two busy Tokio threads.</strong> They almost never wait for a CPU and are rarely preempted, so adding workers only adds queueing (p50 doubles from W32 to W64). <a href="04-ffi-and-scheduling/THREAD-SATURATION.html">Thread saturation</a></li>
<li><strong>Most CPU is in the Rust driver, and much of it is allocation and copying.</strong> After that come the HTTP/2 stack, header handling, session-token work, orchestration and diagnostics. <a href="03-rust-driver-inside-go-v2/HOTSPOTS.html">Rust hotspots</a></li>
<li><strong>The Go scheduler is not the main bottleneck.</strong> Runnable delay is flat at ~22 µs; completions are handed off one at a time (mean batch 1.08). <a href="04-ffi-and-scheduling/SCHEDULER-ANALYSIS.html">Scheduler analysis</a></li>
<li><strong>Best first experiments:</strong> switch the driver allocator to mimalloc/jemalloc, and use 3 Tokio workers on 4 vCPUs. Then cut per-read allocations in headers, session token, diagnostics and completions. <a href="06-optimizations/PROPOSALS.html">Proposals</a></li>
</ol>
<div class=note>Caveats: function attribution is from one perf repeat (differences of a few µs are noise). Process CPU, perf and pprof are separate denominators — never add them. Allocator lock contention is inferred, not measured. The Tokio 2/3/4 evidence comes from an unmatched Gateway V2 campaign.</div>
{measured}
<h2>Headline charts</h2><div class=charts>{charts}</div>
<h2>Where Go v2's 294 µs per read goes (W64, perf)</h2>{bars(groups, "#1d4e89")}
<h2>Rust driver layers inside Go v2 (W64, exclusive)</h2>{bars(layers, "#2e8b57")}
<h2>Rust driver themes (W64, approximate, overlapping)</h2>{bars(themes, "#c0392b")}
<h2>Explore</h2><div class=links>{"".join(f'<a href="{h}">{html.escape(l)}</a>' for l, h in PAGES[1:])}</div>
"""
    (ROOT / "index.html").write_text(shell("index.html", "Go v2 CPU analysis — summary", body), encoding="utf-8")


EXISTING_HTML = [
    "01-results/dashboard.html",
    "02-why-go-v2-costs-more-cpu/cpu-deep-dive.html",
    "03-rust-driver-inside-go-v2/rust-in-go.html",
]
HOME_BAR = ('<div id=report-home style="font:13px sans-serif;background:#0b2545;padding:6px 12px">'
            '<a style="color:#dbe9ff" href="../index.html">&larr; Report summary</a></div>')


def add_home_bar(page):
    path = ROOT / page
    text = path.read_text(encoding="utf-8")
    if "id=report-home" in text:
        return
    text, n = re.subn(r"(<body[^>]*>)", lambda m: m.group(1) + HOME_BAR, text, count=1)
    if n == 0:
        text = HOME_BAR + text
    path.write_text(text, encoding="utf-8")


def main():
    sys.dont_write_bytecode = True
    import relayout_flywheels
    relayout_flywheels.main()
    for page in EXISTING_HTML:
        add_home_bar(page)
    for md in MD_FILES:
        src = ROOT / md
        page = md[:-3] + ".html"
        text = src.read_text(encoding="utf-8")
        title = next((l[2:] for l in text.splitlines() if l.startswith("# ")), md)
        (ROOT / page).write_text(shell(page, title, md_to_html(text)), encoding="utf-8")
    import build_diagnostics
    build_diagnostics.main()
    file_list_page("provenance/index.html", "Provenance", ROOT / "provenance",
                   "Source commits, toolchain versions, binary SHA-256 and build IDs for every leg, plus the run specification.")
    build_index()


if __name__ == "__main__":
    main()
