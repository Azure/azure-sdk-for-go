import html
import json
import pathlib

out = pathlib.Path(__file__).resolve().parent / "cpu-deep-v1"
d = json.loads((out / "rust-layers.json").read_text())
deep = {(l["workers"], l["leg"]): l for l in json.loads((out / "cpu-deep.json").read_text())["legs"]}
L = {(l["workers"], l["leg"]): l for l in d["legs"]}
KINDS = ("compute", "alloc/copy", "kernel")
GROUP = {"Go runtime": "Go side", "Go SDK": "Go side", "native C ABI / completion": "Go↔Rust boundary (C ABI, completion queue)"}
COLORS = {"compute": "#ff7f0e", "alloc/copy": "#d62728", "kernel": "#9467bd"}
e = html.escape


def bars(items, total, scale=3.2):
    h = 26 * len(items) + 10
    svg = [f'<svg width="1100" height="{h}" font-family="Segoe UI,Arial" font-size="12">']
    for i, (name, kinds) in enumerate(items):
        y, x = 5 + i * 26, 300
        svg.append(f'<text x="292" y="{y+15}" text-anchor="end">{e(name)}</text>')
        for k in KINDS:
            w = kinds.get(k, 0) * scale
            svg.append(f'<rect x="{x:.1f}" y="{y+2}" width="{w:.1f}" height="18" fill="{COLORS[k]}"><title>{k}: {kinds.get(k,0):.1f} µs</title></rect>')
            x += w
        t = sum(kinds.values())
        svg.append(f'<text x="{x+6:.1f}" y="{y+15}">{t:.1f} µs ({t/total*100:.0f}%)</text>')
    svg.append("</svg>")
    return "".join(svg)


summary, parts = {}, []
for workers in (64, 32):
    g = L[(workers, "go-v2-native")]
    total = g["totalUs"]
    layers = sorted(g["layers"].items(), key=lambda kv: -sum(kv[1].values()))
    assert abs(sum(sum(v.values()) for _, v in layers) - total) < 0.01
    rust = [(n, v) for n, v in layers if n not in GROUP]
    rust_total = sum(sum(v.values()) for _, v in rust)
    kinds = {k: sum(v.get(k, 0) for _, v in rust) for k in KINDS}
    groups = {"Rust driver stack": rust_total}
    for n, v in layers:
        if n in GROUP:
            groups[GROUP[n]] = groups.get(GROUP[n], 0) + sum(v.values())
    summary[workers] = {"total": round(total, 2), "groups": {k: round(v, 2) for k, v in groups.items()},
                        "rustKinds": {k: round(v, 2) for k, v in kinds.items()},
                        "rustLayers": {n: round(sum(v.values()), 2) for n, v in rust}}
    dp = deep[(workers, "go-v2-native")]
    threads = "".join(f"<tr><td>{e(k)}</td><td>{v:.1f}</td></tr>" for k, v in dp["threadGroupsUs"][:4])
    leaves = "".join(f"<tr><td>{e(n)}</td><td>{sum(v.values()):.1f}</td><td>"
                     + ", ".join(f"{e(a)} {b:.1f}" for a, b in g["layerTopLeaves"].get(n, [])[:5]) + "</td></tr>" for n, v in layers)
    fns = "".join(f"<tr><td>{e(k)}</td><td>{v:.2f}</td></tr>" for k, v in g["driverInnermostFn"][:20])
    futex = "".join(f"<tr><td>{e(k)}</td><td>{v:.2f}</td></tr>" for k, v in dp["syscallOriginUs"][:12])
    alloc = "".join(f"<tr><td>{e(k[:150])}</td><td>{v:.2f}</td></tr>" for k, v in dp["allocOwnerUs"][:15])
    parts.append(f"""<h2>W{workers}: Go v2 process, {total:.1f} µs CPU per read</h2>
<p><b>Rust driver stack {rust_total:.1f} µs ({rust_total/total*100:.0f}%)</b> — allocation/copy {kinds['alloc/copy']:.1f}, compute {kinds['compute']:.1f}, kernel {kinds['kernel']:.1f} µs.
Go↔Rust boundary {groups.get('Go↔Rust boundary (C ABI, completion queue)',0):.1f} µs; Go side {groups.get('Go side',0):.1f} µs.</p>
<p><span style="color:{COLORS['compute']}">■ compute</span> <span style="color:{COLORS['alloc/copy']}">■ allocation/copy (glibc malloc/free, memcpy)</span> <span style="color:{COLORS['kernel']}">■ kernel (futex, context switch, network)</span></p>
{bars([(n, v) for n, v in layers], total)}
<div class="grid">
<div><h3>CPU by OS thread</h3><table><tr><th>Thread</th><th>µs/read</th></tr>{threads}</table>
<h3>Innermost driver function (everything beneath it, incl. malloc it triggers)</h3><table><tr><th>Function</th><th>µs/read</th></tr>{fns}</table></div>
<div><h3>Top leaf functions inside each layer</h3><table><tr><th>Layer</th><th>µs</th><th>Top leaves (µs)</th></tr>{leaves}</table>
<h3>System calls by user-space caller</h3><table><tr><th>syscall ← caller</th><th>µs/read</th></tr>{futex}</table>
<h3>Allocation/copy CPU by first non-allocator frame</h3><table><tr><th>Owner</th><th>µs/read</th></tr>{alloc}</table></div></div>""")

page = f"""<!doctype html><meta charset="utf-8"><title>Rust driver inside Go v2: what consumes what</title>
<style>body{{font-family:Segoe UI,Arial;margin:24px;max-width:1600px}}table{{border-collapse:collapse;font-size:12px;margin-bottom:12px}}
td,th{{border:1px solid #ccc;padding:3px 6px;text-align:right}}td:first-child,th:first-child,td:nth-child(3){{text-align:left;max-width:720px}}
.grid{{display:grid;grid-template-columns:1fr 1fr;gap:18px}}</style>
<h1>Rust driver inside the Go v2 process: what consumes what</h1>
<p>Go v2 process only — no other SDK used as a baseline. Each perf sample is assigned to exactly one layer: the innermost non-kernel frame matching a layer
(driver stage, reqwest, hyper, h2, TLS, Tokio, C ABI/completion, Go SDK, Go runtime). Generic std/alloc/http/bytes/hashbrown/libc frames roll up to the nearest layer,
so <code>malloc</code> called from h2 counts as h2. Layers sum exactly to the profile's CPU/read. Kind comes from the leaf frame. Source: existing perf captures (49 Hz, DWARF stacks).</p>
{''.join(parts)}
<h2>Caveats</h2><ul><li>One perf run per worker count; rows under ~2 µs are noisy.</li>
<li>"driver: other" is driver code outside transport/routing/dataflow/pipeline modules (cosmos_driver execution, diagnostics, models).</li>
<li>Innermost-driver-function and allocation-owner tables are alternative views; do not add them to the layer table.</li></ul>"""
(out / "rust-in-go.html").write_text(page, encoding="utf-8")
(out / "rust-in-go-summary.json").write_text(json.dumps(summary, indent=1))
g = L[(64, "go-v2-native")]
print(json.dumps(summary[64], indent=1))
for k, v in g["driverInnermostFn"][:14]:
    print(f"{v:6.2f} {k[:130]}")
for n in ("driver: other", "native C ABI / completion", "tokio runtime", "Go runtime"):
    print(n, g["layerTopLeaves"].get(n))
