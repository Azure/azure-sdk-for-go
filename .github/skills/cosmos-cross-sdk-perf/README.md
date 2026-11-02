# Cosmos cross-SDK performance campaign

This skill produces independently reproducible performance evidence for Cosmos DB Go v1, Go v2
over the Rust native driver, and the standalone Rust SDK. It provisions or reuses Azure resources,
enforces a materially matched workload, captures process and native evidence, validates CPU
accounting, and renders deterministic JSON, CSV, Markdown, and SVG reports.

## End-to-end example

```bash
SKILL=.github/skills/cosmos-cross-sdk-perf
RUN=$HOME/cosmos-perf-runs/$(date -u +%Y%m%dT%H%M%SZ)
mkdir -p "$RUN"
cp "$SKILL/config/default.json" "$RUN/campaign.json"

python3 "$SKILL/scripts/campaign.py" validate --config "$RUN/campaign.json"
python3 "$SKILL/scripts/campaign.py" prerequisites --phase control
python3 "$SKILL/scripts/campaign.py" provision \
  --config "$RUN/campaign.json" --state "$RUN/azure-state.json" --dry-run
# Review the plan, then repeat without --dry-run.

# Copy the run directory and skill to the provisioned VM. Build each pinned adapter, then:
sudo "$SKILL/scripts/profile-leg.sh" \
  --leg go-v1 --run-dir "$RUN" --config "$RUN/campaign.json" -- \
  /path/to/go-v1-adapter --config "$RUN/campaign.json"

python3 "$SKILL/scripts/campaign.py" report --run-dir "$RUN" --config "$RUN/campaign.json"
python3 "$SKILL/scripts/campaign.py" sanitize \
  --run-dir "$RUN" --output "$RUN/cosmos-perf-sanitized.tar.gz"
```

Repeat `profile-leg.sh` for `go-v2-native` and `rust`. Adapter commands must write the schema files
described in `SKILL.md`.

The default pins correspond to Go v1.5.0, the verified Go v2 diagnostics-fix head from
Azure/azure-sdk-for-go#27873, the Rust SDK campaign commit, and the native-driver source and
distribution commits recorded by `Azure/azure-cosmos-driver` v0.2.0 provenance. Re-verify and
update all pins together before using newer sources.

## Runtime and cost

Allow 15-30 minutes for first-time VM/toolchain setup, 30 seconds of excluded warmup and 300 seconds
of measurement per leg, plus profile processing. A complete sequential run usually takes 25-60
minutes. Azure charges accrue for the Linux VM, managed disk, networking, Cosmos DB throughput and
requests, and any newly created account. The default is a four-vCPU VM in the verified Cosmos
region. Review current Azure pricing before provisioning and tear down campaign-owned resources
when finished.

## Output layout

```text
run/
  azure-state.json
  config/effective.json
  go-v1/
    raw/                 # immutable evidence
    derived/             # classified tables and charts
    metrics.json
    perf-self.csv
    provenance.json
  go-v2-native/
  rust/
  comparison/
    comparison.json
    comparison.csv
    report.md
    throughput.svg
    latency.svg
    cpu-flywheel.svg
  SHA256SUMS
```

Raw directories are made read-only after capture. Never modify them to repair a run; create a new
run. Derived artifacts may be regenerated after classifier changes.

## Contributing results

Run `sanitize`, inspect the archive contents, and submit the archive through the internal
collaboration channel chosen by the campaign owner. Include the generated manifest hash and the
three source pins. Do not commit result bundles, generated charts, Azure resource identifiers,
credentials, or customer data to this repository.
