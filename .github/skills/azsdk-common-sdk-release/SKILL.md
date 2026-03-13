---
name: azsdk-common-sdk-release
license: MIT
metadata:
  version: "1.0.0"
description: "**UTILITY SKILL** Check SDK package release readiness and trigger the release pipeline for Azure SDK packages. USE FOR: \"release SDK\", \"trigger release\", \"check release readiness\", \"release pipeline\", \"publish package\", \"ship SDK\". INVOKES: azsdk_release_sdk. FOR SINGLE OPERATIONS: Use azsdk_release_sdk with checkReady=true for readiness check only."
compatibility:
  requires: "azure-sdk-mcp server, SDK package merged on release branch"
  supports: ".NET, Java, JavaScript, Python, Go"
---

# SDK Release

## MCP Tools

| Tool                | Purpose                                                     |
| ------------------- | ----------------------------------------------------------- |
| `azsdk_release_sdk` | Check release readiness and/or trigger the release pipeline |

## Steps

1. **Collect Info** — Get `packageName` and `language` from the user. Optionally get `branch` (defaults to main).
2. **Check Readiness** — Run `azsdk_release_sdk` with `checkReady: true` to verify API review approval, changelog, package name approval, and release date.
3. **Review Results** — If not ready, display failing checks and guide user to resolve.
4. **Trigger Release** — Once ready, run `azsdk_release_sdk` with `checkReady: false`. Show pipeline link and inform user they must approve the release stage.

## MCP Prerequisites

Requires `azure-sdk-mcp` server. No CLI fallback — prompt user to configure MCP if unavailable.
