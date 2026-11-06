"""Re-lay out the CPU flywheel SVGs for the team report.

Reads the original charts from share-v4/comparison-v1 (never modified) and writes uncramped copies
to team-report/01-results. Data values are copied verbatim; only positions and labels change.
"""
import html
import math
import pathlib
import re

ROOT = pathlib.Path(__file__).resolve().parents[2]
SRC = ROOT.parent / "share-v4" / "comparison-v1"
DST = ROOT / "01-results"

CX, CY, R, STROKE = 200, 300, 130, 52
CIRC = 2 * math.pi * R
LEG_X, CAP_X, US_X = 410, 790, 930
ROW0, ROW_H = 104, 27


def esc(s):
    return html.escape(s, quote=True)


def relayout(src):
    s = src.read_text(encoding="utf-8")
    title = html.unescape(re.search(r"<title>(.*?)</title>", s).group(1))
    arcs = re.findall(r'<circle [^>]*stroke="(#[0-9a-fA-F]+)"[^>]*stroke-dasharray="([\d.]+) [\d.]+"'
                      r'[^>]*stroke-dashoffset="(-?[\d.]+)"', s)
    old_r = float(re.search(r'<circle [^>]* r="([\d.]+)"', s).group(1))
    old_circ = 2 * math.pi * old_r
    rows = re.findall(r'<rect x="\d+" y="\d+" width="10" height="10" fill="(#[0-9a-fA-F]+)"/>\s*'
                      r'<text [^>]*>(.*?)</text>\s*<text [^>]*>(.*?)</text>\s*<text [^>]*>(.*?)</text>', s)
    centre = re.findall(r'<text x="\d+" y="\d+" text-anchor="middle"[^>]*>(.*?)</text>', s)
    notes = [html.unescape(n) for n in re.findall(r'<text x="24" y="\d+"[^>]*font-size="12">(.*?)</text>', s)]
    if len(arcs) != len(rows) or len(centre) != 2 or not notes:
        raise SystemExit(f"{src.name}: unexpected structure")
    used = float(re.match(r"([\d.]+)%", html.unescape(centre[0])).group(1))
    idle = round(100 - used, 2)

    height = ROW0 + (len(rows) + 1) * ROW_H + 40 + len(notes) * 22 + 20
    height = max(height, CY + R + STROKE)
    out = [f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1080 {height}" role="img" font-family="sans-serif">',
           f"<title>{esc(title)}</title>",
           '<rect width="100%" height="100%" fill="white"/>',
           f'<text x="24" y="36" font-size="21" font-weight="600">{esc(title)}</text>',
           f'<circle cx="{CX}" cy="{CY}" r="{R}" fill="none" stroke="#e2e8f0" stroke-width="{STROKE}"/>']
    for color, dash, off in arcs:
        length = float(dash) / old_circ * CIRC
        offset = float(off) / old_circ * CIRC
        out.append(f'<circle cx="{CX}" cy="{CY}" r="{R}" fill="none" stroke="{color}" stroke-width="{STROKE}" '
                   f'stroke-dasharray="{length:.6f} {CIRC - length:.6f}" stroke-dashoffset="{offset:.6f}" '
                   f'transform="rotate(-90 {CX} {CY})"/>')
    out.append(f'<text x="{CX}" y="{CY - 2}" text-anchor="middle" font-size="30" font-weight="600">{centre[0]}</text>')
    out.append(f'<text x="{CX}" y="{CY + 24}" text-anchor="middle" font-size="12" fill="#475569">of 4 vCPUs used</text>')
    out.append(f'<text x="{CX}" y="{CY + 40}" text-anchor="middle" font-size="12" fill="#475569">by the benchmark process</text>')

    hy = ROW0 - 24
    out.append(f'<text x="{LEG_X}" y="{hy}" font-size="12" font-weight="600" fill="#334155">Category</text>')
    out.append(f'<text x="{CAP_X + 90}" y="{hy}" text-anchor="end" font-size="12" font-weight="600" fill="#334155">% of 4 vCPUs</text>')
    out.append(f'<text x="{US_X + 100}" y="{hy}" text-anchor="end" font-size="12" font-weight="600" fill="#334155">CPU µs/read</text>')
    out.append(f'<line x1="{LEG_X - 22}" y1="{hy + 8}" x2="1040" y2="{hy + 8}" stroke="#cbd5e1"/>')
    for i, (color, label, cap, us) in enumerate(rows):
        y = ROW0 + i * ROW_H
        if i % 2:
            out.append(f'<rect x="{LEG_X - 22}" y="{y - 17}" width="{1040 - LEG_X + 22}" height="{ROW_H}" fill="#f8fafc"/>')
        cap_v = html.unescape(cap).replace(" capacity", "")
        us_v = html.unescape(us).replace(" us/read", "")
        out.append(f'<rect x="{LEG_X - 16}" y="{y - 10}" width="11" height="11" rx="2" fill="{color}"/>')
        out.append(f'<text x="{LEG_X}" y="{y}" font-size="13">{label}</text>')
        out.append(f'<text x="{CAP_X + 90}" y="{y}" text-anchor="end" font-family="monospace" font-size="13">{esc(cap_v)}</text>')
        out.append(f'<text x="{US_X + 100}" y="{y}" text-anchor="end" font-family="monospace" font-size="13">{esc(us_v)}</text>')
    y = ROW0 + len(rows) * ROW_H
    out.append(f'<line x1="{LEG_X - 22}" y1="{y - 19}" x2="1040" y2="{y - 19}" stroke="#cbd5e1"/>')
    out.append(f'<rect x="{LEG_X - 16}" y="{y - 10}" width="11" height="11" rx="2" fill="#e2e8f0" stroke="#94a3b8"/>')
    out.append(f'<text x="{LEG_X}" y="{y}" font-size="13" fill="#475569">Grey: idle or other processes (not this process)</text>')
    out.append(f'<text x="{CAP_X + 90}" y="{y}" text-anchor="end" font-family="monospace" font-size="13" fill="#475569">{idle:.2f}%</text>')
    out.append(f'<text x="{US_X + 100}" y="{y}" text-anchor="end" font-family="monospace" font-size="13" fill="#475569">–</text>')

    ny = max(y, CY + R + STROKE // 2) + 44
    for n in notes:
        if n.startswith("Gray remainder"):
            continue
        out.append(f'<text x="24" y="{ny}" font-size="12" fill="#475569">{esc(n)}</text>')
        ny += 22
    out.append("</svg>")
    return "\n".join(out) + "\n"


def main():
    files = sorted(SRC.glob("*-cpu-flywheel.svg"))
    if not files:
        raise SystemExit(f"no flywheels under {SRC}")
    for f in files:
        (DST / f.name).write_text(relayout(f), encoding="utf-8")
        print("relaid", f.name)


if __name__ == "__main__":
    main()
