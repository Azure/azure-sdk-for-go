import pathlib, json, re, collections
root = pathlib.Path("/mnt/cosmos-ffi-flow-scratch-20261105-a/go-v2-ffi-flow-latest-main-20261105-a")
def snap(p):
    s = json.loads(p.read_text())
    out = {}
    for tid, o in s["threadObservations"].items():
        st = o.get("status") or ""
        name = re.search(r"Name:\t(.*)", st)
        vol = re.search(r"voluntary_ctxt_switches:\t(\d+)", st)
        nv = re.search(r"nonvoluntary_ctxt_switches:\t(\d+)", st)
        ss = (o.get("schedstat") or "0 0 0").split()
        out[tid] = dict(name=name.group(1) if name else "?", run=int(ss[0]), wait=int(ss[1]), slices=int(ss[2]),
                        vol=int(vol.group(1)) if vol else 0, nv=int(nv.group(1)) if nv else 0)
    return out, s
res = []
for leg in sorted(root.glob("run/w*/*/repeat-*/go-v2-native")):
    a, b = leg/"raw/snapshot-measurement-start.json", leg/"raw/snapshot-measurement-end.json"
    m = leg/"metrics.json"
    if not (a.exists() and b.exists() and m.exists()):
        continue
    met = json.loads(m.read_text())
    ok = met["counts"]["success"]; wall = met["timing"]["wallSeconds"]
    s0, _ = snap(a); s1, _ = snap(b)
    g = collections.defaultdict(lambda: collections.Counter())
    for tid, e in s1.items():
        st = s0.get(tid, dict(run=0, wait=0, slices=0, vol=0, nv=0))
        d = {k: e[k]-st[k] for k in ("run","wait","slices","vol","nv")}
        name = e["name"]
        grp = g[name]; grp["threads"] += 1 if d["run"] > 1e6 else 0
        for k, v in d.items(): grp[k] += v
    row = {"leg": str(leg.relative_to(root/"run")).replace("/go-v2-native",""), "ok": ok, "wall": round(wall,1),
           "procCpuUs": round(met["timing"]["processCPUSeconds"]*1e6/ok,1), "tput": round(ok/wall)}
    row["thr"] = {n: dict(threads=c["threads"], runUs=round(c["run"]/1e3/ok,1), waitUs=round(c["wait"]/1e3/ok,1),
                  slices=round(c["slices"]/ok,2), vol=round(c["vol"]/ok,2), nv=round(c["nv"]/ok,3),
                  usPerSlice=round(c["run"]/1e3/max(c["slices"],1),1), cores=round(c["run"]/1e9/wall,2))
                  for n, c in g.items() if c["run"] > 1e7}
    res.append(row)
(root/"analysis/cpu-deep-v1").mkdir(parents=True, exist_ok=True)
(root/"analysis/cpu-deep-v1/thread-schedstat.json").write_text(json.dumps(res))
for r in res:
    print(r["leg"], "ok", r["ok"], "tput", r["tput"], "procCPUus/read", r["procCpuUs"])
    for n, t in r["thr"].items():
        print("   ", n, t)
print("VM_PY_COMPLETE")
