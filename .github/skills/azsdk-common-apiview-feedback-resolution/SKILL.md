---
name: azsdk-common-apiview-feedback-resolution
license: MIT
metadata:
  version: "1.0.0"
description: "Resolve feedback from APIView reviews on Azure SDK packages. **UTILITY SKILL**. USE FOR: \"APIView comments\", \"resolve API review feedback\", \"rename type per reviewer\", \"SDK API surface changes\", \"regenerate SDK after review\". DO NOT USE FOR: general code review, non-APIView feedback, manual SDK editing. INVOKES: azsdk_apiview_get_comments, azsdk_typespec_delegate_apiview_feedback."
compatibility:
  requires: "azure-sdk-mcp server, SDK pull request with APIView review link"
---

# APIView Feedback Resolution

## MCP Tools

| Tool | Purpose |
|------|---------|
| `azsdk_apiview_get_comments` | Retrieve APIView comments |
| `azsdk_typespec_delegate_apiview_feedback` | AI-resolve feedback |
| `azsdk_run_typespec_validation` | Validate TypeSpec changes |
| `azsdk_package_generate_code` | Regenerate SDK |

## Steps

1. **Retrieve** — Get APIView URL from SDK PR, run `azsdk_apiview_get_comments`.
2. **Categorize** — Group as Critical/Suggestions/Informational. See [feedback steps](references/feedback-resolution-steps.md).
3. **Resolve** — Use `azsdk_typespec_delegate_apiview_feedback` for TypeSpec changes; apply code-only fixes directly.
4. **Validate** — Run validation, regenerate SDK, build and test.
5. **Confirm** — Verify all items addressed, inform user to request re-review.

## Examples

- "Resolve the APIView comments on my SDK pull request"
- "What feedback did the API reviewer leave on my package?"

## Troubleshooting

- **No comments returned**: Verify the PR has an APIView revision link and MCP server is connected.
- **Validation fails**: Re-run `azsdk_run_typespec_validation` after fixing TypeSpec errors.
- **MCP unavailable**: Review APIView comments in browser and apply fixes directly.
