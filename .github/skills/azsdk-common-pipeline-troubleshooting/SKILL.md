---
name: azsdk-common-pipeline-troubleshooting
license: MIT
metadata:
  version: "1.0.0"
description: "Diagnose and resolve failures in Azure SDK CI and generation pipelines. **UTILITY SKILL**. USE FOR: \"pipeline failed\", \"build failure\", \"CI check failing\", \"SDK generation error\", \"reproduce pipeline locally\", \"debug SDK pipeline\". INVOKES: azsdk_analyze_pipeline, azsdk_verify_setup, azsdk_package_build_code, azsdk_package_run_check, azsdk_package_pack."
compatibility:
  requires: "azure-sdk-mcp server, Azure DevOps pipeline build ID"
---

# Pipeline Troubleshooting

## MCP Prerequisites

Requires `azure-sdk-mcp` server for pipeline analysis and local reproduction.

## MCP Tools

| Tool                       | Purpose                      |
| -------------------------- | ---------------------------- |
| `azsdk_analyze_pipeline`   | Analyze pipeline failures    |
| `azsdk_verify_setup`       | Verify local environment     |
| `azsdk_package_build_code` | Reproduce build locally      |
| `azsdk_package_run_check`  | Run validation checks        |
| `azsdk_package_pack`       | Create SDK artifact packages |

## Steps

1. **Identify** — Get build ID, run `azsdk_analyze_pipeline`. Categorize failure type.
2. **Analyze** — See [failure patterns](references/failure-patterns.md) for common causes.
3. **Reproduce** — Run `azsdk_verify_setup`, then `azsdk_package_build_code` or `azsdk_package_run_check`.
4. **Fix** — Apply direct edits for code or TypeSpec changes.
5. **Verify** — Confirm fix locally, push changes, monitor pipeline re-run.

## Examples

- "My pipeline build 12345 failed, help me debug it"
- "Reproduce CI failure locally for azure-sdk-for-python"

## Troubleshooting

If `azsdk_analyze_pipeline` returns no data, verify the build ID and MCP connection.

## CLI Fallback

Without MCP: view pipeline logs in Azure DevOps browser UI, download and inspect failure stages manually.
