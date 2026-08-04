---
name: query-release-plan
description: Query an Azure SDK Release Plan from Azure DevOps (the `Release` project work item behind the auth-gated azsdk-releaseplan-dashboard URL) by release-plan id or package name, and print its state, release month, and per-language release status. Use for requests like "query release plan 35135", "what is the release plan for azure-resourcemanager-foo", "is the Java release approved/released for <module>", or "show release plan details".
---

# Query Release Plan

Read an Azure SDK Release Plan from Azure DevOps and print its key fields (state, plan type,
release month, spec approval, and per-language release status/version/PR).

A Release Plan is a work item in the `Release` project of `dev.azure.com/azure-sdk`. AutoPR bodies
link to it as `https://azsdk-releaseplan-dashboard-*.azurewebsites.net/?releaseplan=<id>`.

## Key insight: id vs. dashboard URL

- The dashboard URL is auth-gated — a plain web fetch returns HTTP 401, so you cannot read the
  plan's contents that way.
- The `?releaseplan=<n>` value is the plan's `Custom.ReleasePlanID` field, which is NOT always the
  ADO work item id. For newer plans the two coincide (e.g. `35135` = work item 35135), but for
  older plans they differ (e.g. `releaseplan=2101` resolves to work item 32864, whose title is
  "Release Plan - 2101 - ..."). Work item 2101 is a completely unrelated ancient plan — so never
  assume the dashboard number is the work item id.
- Resolve the dashboard number via WIQL on `Custom.ReleasePlanID`, then read the work item:
  `GET https://dev.azure.com/azure-sdk/Release/_apis/wit/workitems/<workitemid>?api-version=7.1`.
  The `--id` flag of the script does this for you; use `--work-item` to pass a raw work item id
  directly.

## Auth

Uses the Azure DevOps token for the corp account. The default `az` login may be a personal tenant
→ `TF400813`. Switch first:

```powershell
az account set --subscription "Azure SDK Developer Playground"
```

The script calls `az account get-access-token --resource 499b84ac-1321-427f-aa17-267ca6975798`
(`499b84ac-...` = Azure DevOps) itself; you only need to be signed in as corp.

## Workflow

Run the bundled script:

```powershell
$env:PYTHONIOENCODING="utf-8"

# by release-plan id (the number from the dashboard URL; resolved via ReleasePlanID)
python .github/skills/query-release-plan/scripts/query_release_plan.py --id 35135

# by raw ADO work item id (skip ReleasePlanID resolution)
python .github/skills/query-release-plan/scripts/query_release_plan.py --work-item 32864

# by Java package name (finds the newest matching plan)
python .github/skills/query-release-plan/scripts/query_release_plan.py --java-package azure-resourcemanager-containerservicepreparedimgspec

# by any-language package name
python .github/skills/query-release-plan/scripts/query_release_plan.py --package azure-mgmt-foo

# dump all raw fields
python .github/skills/query-release-plan/scripts/query_release_plan.py --id 35135 --json
```

### How it resolves inputs to a work item

- `--id <n>` — treated as the dashboard `Custom.ReleasePlanID`; WIQL finds the work item with that
  ReleasePlanID. If none has it, `<n>` is used as a raw work item id.
- `--work-item <n>` — read work item `<n>` directly (no ReleasePlanID lookup).
- `--java-package` / `--package` — WIQL on the package field(s) for
  `[System.WorkItemType]='Release Plan'`:
  - `--java-package` → `[Custom.JavaPackageName]='<pkg>'`
  - `--package` → any of `Custom.{Dotnet,Java,Go,Python,JavaScript}PackageName`
  - If no exact match, it retries with `CONTAINS` (the plan's package field may differ from the
    repo module name, e.g. module `...commvaultcontentstore` vs plan package
    `azure-resourcemanager-commvault`).

Newest matching id is shown; pass `--work-item` to pick a specific one when several match.

## Fields printed

Top-level: `System.Title`, `System.State`, `Custom.ReleasePlanType`, `Custom.SDKReleasemonth`,
`Custom.SDKtypetobereleased`, `Custom.SDKLanguages`, `Custom.APISpecApprovalStatus`,
`Custom.ReleasePlanSubmittedby`, `Custom.NamespaceApprovalIssue`.

Per language (`.NET, Java, Go, Python, JavaScript`) using the field suffix
(`Dotnet, Java, Go, Python, JavaScript`):

- `Custom.ReleaseStatusFor<Lang>` — e.g. `Released`, `Approval Pending`, `ready for review`
- `Custom.ReleasedVersionFor<Lang>`
- `Custom.SDKPullRequestFor<Lang>` (the SDK PR URL)
- `Custom.<Lang>PackageName`

Other useful raw fields (via `--json`): `Custom.GenerationStatusFor<Lang>`,
`Custom.SDKPullRequestStatusFor<Lang>`, `Custom.ReleasePipelineFor<Lang>`,
`Custom.SDKGenerationPipelineFor<Lang>`, `Custom.AttestationStatus`, `Custom.ApiSpecProjectPath`.

## Interpreting Java status

A `azure-resourcemanager-*` module whose CHANGELOG top entry is a dated (not `Unreleased`) beta but
which is not yet on Maven typically shows `Custom.ReleaseStatusForJava = Approval Pending`
(generation done, PR merged, publish not yet approved).

## Companion skills

- `query-released-azure-lib` — the Release Plan work item holds the authoritative per-language
  package names in `Custom.<Lang>PackageName`; feed those to that skill's `--dotnet/--java/...`
  flags when the collapsed service token differs by language.

## Success criteria

- Resolved the plan (by id or package) and printed its state + per-language status.
- Did not attempt to web-fetch the dashboard URL (it 401s); used the REST API.
