# Go SDK Breaking Changes Review

Follow these steps to help users review the Go SDK breaking changes. Show the high-level description to users to help them understand the flow.

## Step 1: Prerequisites

- Run `go run . environment` under `../../eng/tools/generator` to check the environment.
- Ensure the SDK repository and spec repository are accessible.

## Step 2: Generate Go SDK

- If user provide a local spec path, then run `go run . generate <sdk-repo-path> <spec-repo-path> --tsp-config=<tsp-config-path>` under `../../eng/tools/generator` to generate the Go SDK.
- If user provide a GitHub PR link, then run `go run . generate <sdk-repo-path> <spec-repo-path> --github-pr=<pr-link>` under `../../eng/tools/generator` to generate the Go SDK.
- Refer the [README](../../eng/tools/generator/README.md) for detailed usage.

## Step 3: Review Breaking Changes

### 3.1 Determine Review Approach

- Check the generation type from the generator output
- If generation type is `MigrateToTypeSpec`, follow the [Azure Go SDK Breaking Changes Review and Resolution Guide for TypeSpec Migration](../../documentation/sdk-breaking-changes-guide-migration.md) to review the breaking changes.
- If generation type is anything else, follow the [Azure Go SDK Breaking Changes Review and Resolution Guide](../../documentation/sdk-breaking-changes-guide.md) to review the breaking changes.

### 3.2 Analyze Breaking Changes Systematically

Iterate through all items in the guide:

- For each item, check if the changelog entries could match the `Changelog Pattern`.
  - If match, get the `Reason`, `Spec Pattern`, `Impact`, and `Resolution` from the guide.
  - If the breaking changes can be resolved, follow the resolution guide to apply the necessary changes.
- If there is any changes made in the spec folder:
  - Revert the code changes under the package path.
  - Regenerate the SDK locally to see if the changelogs are updated.
- For breaking changes that don't match any patterns, document them for manual review

## Step 4: Gather Review Result

Provide comprehensive breaking changes analysis report in a new markdown file including: the breaking changes detected, their reason and impact, and any resolutions applied. Use structured format to make it easy to read and understand.
