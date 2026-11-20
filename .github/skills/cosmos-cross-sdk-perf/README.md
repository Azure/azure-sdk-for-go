# Cosmos cross-SDK performance campaign

This skill produces independently reproducible point-read performance evidence for Cosmos DB Go v1, Go v2
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
sudo "$SKILL/scripts/run-clean-leg.sh" \
  --leg go-v1 --workers 8 --repeat 1 --run-dir "$RUN" --config "$RUN/campaign.json" -- \
  /path/to/go-v1-adapter --config "$RUN/campaign.json"

python3 "$SKILL/scripts/campaign.py" report --run-dir "$RUN" --config "$RUN/campaign.json"
python3 "$SKILL/scripts/campaign.py" sanitize \
  --run-dir "$RUN" --output "$RUN/cosmos-perf-sanitized.tar.gz"
```

Repeat the ordered clean-score matrix twice, then run the separate perf and Go-diagnostic classes.
Adapter commands must write the schema files described in `SKILL.md`.

## Portable adapters

Checked-in adapter source lives under `adapters/go` and `adapters/rust`. Build and stage it only
from clean checkouts at the configured commits:

```bash
"$SKILL/scripts/build-go-adapters.sh" \
  --go-v2-checkout /src/azure-sdk-for-go \
  --go-v2-commit 7201f4b70fac1184e66379278566f1ac8d299109 \
  --native-module /src/azure-cosmos-driver \
  --native-commit 1e449be78bffef6ef99e6bc1810216a84a9dac16 \
  --output "$RUN/build/go"
"$SKILL/scripts/build-rust-adapter.sh" \
  --checkout /src/azure-sdk-for-rust \
  --commit f87618c80049eab5fd1ed8f16c7cf37488f43663 \
  --output "$RUN/build/rust"
"$SKILL/scripts/stage-adapters.sh" \
  --go-build "$RUN/build/go" --rust-build "$RUN/build/rust" \
  --output "$RUN/staged"
```

Each adapter supports `seed`, `verify`, `preflight`, `run`, and `cleanup`. Supply endpoint,
database, container, region (Rust), and an owner consisting only of lowercase letters, digits, and
hyphens. IDs and partition keys are deterministically scoped to that owner. Cleanup reads and
checks every expected item before deleting it; malformed or differently owned items are retained.

The account key is accepted only from an inherited descriptor numbered 3 or higher. Open that
descriptor immediately before launch, pass only its number through `--key-fd`, and close it after
the process starts. Do not place the key in arguments or environment variables. For example, if a
trusted local command writes only the key to standard output:

```bash
exec {cosmos_key_fd}< <(trusted-key-command)
"$SKILL/scripts/run-clean-leg.sh" \
  --leg go-v1 --workers 8 --repeat 1 --run-dir "$RUN" --config "$RUN/campaign.json" \
  --key-fd "$cosmos_key_fd" -- \
  "$RUN/staged/bin/cosmos-perf-go-v1" \
  --endpoint "$COSMOS_ENDPOINT" --database "$COSMOS_DATABASE" \
  --container "$COSMOS_CONTAINER" --owner "$CAMPAIGN_OWNER" run
exec {cosmos_key_fd}<&-
```

The launchers set only the descriptor number and non-secret run metadata in the child environment.
Adapters reject common key-bearing flags and environment variables. They write a raw adapter
result under the immutable evidence directory; `scripts/convert-adapter-output.py` deterministically
adds launcher-captured process/cgroup memory and emits the existing `metrics.json` schema.

The default pins correspond to Go v1.5.0, the verified Go v2 diagnostics-fix head from
Azure/azure-sdk-for-go#27873, the Rust SDK campaign commit, and the native-driver source and
distribution commits recorded by `Azure/azure-cosmos-driver` v0.2.0 provenance. Re-verify and
update all pins together before using newer sources. Flamegraphs use the pinned commit of Brendan
Gregg's canonical `brendangregg/FlameGraph` repository: `stackcollapse-perf.pl --all`,
`flamegraph.pl`, and `difffolded.pl`. The generator rejects a dirty or mismatched checkout and
records the commit in the campaign root.

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
  w8/
    clean-score/repeat-1/{go-v1,go-v2-native,rust}/
    clean-score/repeat-2/{go-v1,go-v2-native,rust}/
    perf/repeat-1/{go-v1,go-v2-native,rust}/
    perf/differential/
    go-diagnostic/repeat-1/{go-v1,go-v2-native}/
  w16/
  comparison/
    summary.json
    summary.csv
    setup-equivalence.csv
    report.md
    throughput.svg
    latency.svg
    cpu-flywheel.svg
  SHA256SUMS
  dashboard.html
```

Raw directories are made read-only after capture. Never modify them to repair a run; create a new
run. Only clean-score results are authoritative. Derived artifacts may be regenerated after
classifier changes.

## Contributing results

Run `sanitize`, inspect the archive contents, and submit the archive through the internal
collaboration channel chosen by the campaign owner. Include the generated manifest hash and the
three source pins. Do not commit result bundles, generated charts, Azure resource identifiers,
credentials, or customer data to this repository.
