---
name: azsdk-common-apiview-feedback-resolution
license: MIT
metadata:
  version: "1.0.0"
description: "Retrieve and resolve feedback from APIView reviews on Azure SDK packages. **UTILITY SKILL**. USE FOR: \"APIView comments\", \"resolve API review feedback\", \"rename type per reviewer\", \"SDK API surface changes\", \"regenerate SDK after review\". INVOKES: azsdk_apiview_get_comments, azsdk_typespec_delegate_apiview_feedback, azsdk_run_typespec_validation, azsdk_package_generate_code."
compatibility:
  requires: "azure-sdk-mcp server, SDK pull request with APIView review link"
---

# APIView Feedback Resolution

## MCP Tools

| MCP Tool                                   | Purpose                     |
| ------------------------------------------ | --------------------------- |
| `azsdk_apiview_get_comments`               | Retrieve APIView comments   |
| `azsdk_typespec_delegate_apiview_feedback` | AI-resolve APIView feedback |
| `azsdk_run_typespec_validation`            | Validate TypeSpec changes   |
| `azsdk_package_generate_code`              | Regenerate SDK              |

## Steps

1. **Retrieve Comments** — Get APIView revision URL from SDK PR, run `azsdk_apiview_get_comments`.
2. **Categorize** — Group as Critical/Suggestions/Informational. See [feedback steps](references/feedback-resolution-steps.md).
3. **Resolve** — For TypeSpec changes, use `azsdk_typespec_delegate_apiview_feedback`. For code-only fixes, apply directly.
4. **Validate & Regenerate** — Run validation, regenerate SDK, build and test.
5. **Confirm** — Verify all items addressed, inform user to request re-review.

## MCP Prerequisites

Requires `azure-sdk-mcp` server connected and authenticated.

## CLI Fallback

Without MCP, review APIView comments in browser and apply fixes to TypeSpec or SDK code directly.
