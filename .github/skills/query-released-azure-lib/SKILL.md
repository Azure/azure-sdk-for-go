---
name: query-released-azure-lib
description: Check which Azure SDK languages (.NET, Java, Python, JS, Go) have released a given management or client library, by reading the per-language release inventory CSVs at azure-sdk/_data/releases/latest. Handles the different package-name conventions per language. Use for requests like "did .NET release <lib>", "which languages released azure-resourcemanager-foo", "is <service> released in Python/Go", or "release status of <lib> across languages".
---

# Query Released Azure Lib

Report which languages have released an Azure SDK library (and at what version / date), for a
given service, by reading the published per-language release inventory CSVs.

## Data source

`https://raw.githubusercontent.com/Azure/azure-sdk/main/_data/releases/latest/<lang>-packages.csv`
for `lang ∈ {dotnet, java, python, js, go}`. These are the same CSVs that back
`https://azure.github.io/azure-sdk/releases/latest` — whose `?search=` box is client-side JS, so a
static web fetch of that page is useless; fetch the CSVs instead.

Each row has `VersionGA` and `VersionPreview`. Interpretation:

- non-empty `VersionGA` → GA released (that version)
- else non-empty `VersionPreview` → preview/beta released
- row present but both empty → listed but no version yet
- no row → not released in that language

> Note: Java's CSV has an extra `GroupId` column, so parse by header name (the script uses
> `csv.DictReader`), never by fixed column index.

## The naming problem (why matching is non-trivial)

The same service is named differently per language, and a single field is not enough:

- `.NET` `Azure.ResourceManager.<Svc>[.<Sub>]` (PascalCase, dots; 3–4 segments)
- `Java` `azure-resourcemanager-<svc>`
- `Python` `azure-mgmt-<svc>` (lowercase, often no separators)
- `JS` `@azure/arm-<svc>`
- `Go` `sdk/resourcemanager/<mid>/arm<suffix>`

`ServiceName` and `RepoPath` are not reliable keys: for a sub-service like Compute BulkActions,
`.NET`/`Java`/`JS` `RepoPath` collapses to the parent `compute`, and a service can exist in both a
dot/dash form (`Azure.ResourceManager.Compute.BulkActions`) and a flattened form
(`Azure.ResourceManager.ComputeBulkActions`).

### Matching approach: collapsed token

Reduce both the query and each row to a collapsed token — strip the language prefix, drop
`.`/`-`/`/`, lowercase — then compare. So `compute-bulkactions`, `Compute.BulkActions`, and
`computebulkactions` all collapse to `computebulkactions`. For Go, the `arm` infix is handled by
emitting `<mid>`, `<suffix>`, and the combined `<mid><suffix>` token, so
`sdk/resourcemanager/compute/armbulkactions` also collapses to `computebulkactions`. Matching also
considers the collapsed `ServiceName`. All matching rows are printed (a service may have several
variants), so ambiguity is visible rather than hidden.

## Workflow

Run the bundled script:

```powershell
$env:PYTHONIOENCODING="utf-8"

# by shared service token (recommended) — collapsed match across all languages
python .github/skills/query-released-azure-lib/scripts/query_released_azure_lib.py --service commvaultcontentstore

# emphasise a single language (e.g. is .NET released?)
python .github/skills/query-released-azure-lib/scripts/query_released_azure_lib.py --service networkcloud --lang dotnet

# give exact per-language package names when the service token differs by language
# (e.g. from the query-release-plan skill's Custom.<Lang>PackageName fields)
python .github/skills/query-released-azure-lib/scripts/query_released_azure_lib.py \
    --dotnet Azure.ResourceManager.Commvault --java azure-resourcemanager-commvault \
    --python azure-mgmt-commvault --js @azure/arm-commvault --go armcommvault
```

Output ends with a summary: `released:`, `not released:`, and a dedicated `.NET:` line
(released / NOT released / not checked).

## Companion skills

- `query-release-plan` — the Release Plan work item holds the authoritative per-language package
  names in `Custom.<Lang>PackageName`; feed those to the `--dotnet/--java/...` flags here when the
  collapsed token differs by language.

## Success criteria

- Every requested language is classified released / not-released, with version(s) and date(s), and
  the `.NET` status is called out explicitly.
- Package-name convention differences are handled via the collapsed token; all matching rows are
  shown, not just the first.
