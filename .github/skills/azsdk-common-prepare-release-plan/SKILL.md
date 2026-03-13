---
name: azsdk-common-prepare-release-plan
license: MIT
metadata:
  version: "1.0.0"
description: "**UTILITY SKILL** Create and manage release plan work items for Azure SDK releases across languages. USE FOR: \"create release plan\", \"update release plan\", \"link SDK PR to plan\", \"namespace approval\", \"check release plan status\". INVOKES: azsdk_create_release_plan, azsdk_get_release_plan, azsdk_update_sdk_details_in_release_plan, azsdk_link_sdk_pull_request_to_release_plan, azsdk_link_namespace_approval_issue."
compatibility:
  requires: "azure-sdk-mcp server, API spec PR in Azure/azure-rest-api-specs"
---

# Prepare Release Plan

> Do not display Azure DevOps work item URLs. Only provide Release Plan Link and ID.

## MCP Tools

| Tool                                          | Purpose              |
| --------------------------------------------- | -------------------- |
| `azsdk_create_release_plan`                   | Create release plan  |
| `azsdk_get_release_plan`                      | Get plan details     |
| `azsdk_get_release_plan_for_spec_pr`          | Find plan by PR      |
| `azsdk_update_sdk_details_in_release_plan`    | Update SDK info      |
| `azsdk_link_sdk_pull_request_to_release_plan` | Link SDK PR          |
| `azsdk_link_namespace_approval_issue`         | Link namespace issue |

## Steps

1. **Prerequisites** — Check for API spec PR; prompt if unavailable.
2. **Check Existing** — Query by plan number or spec PR link.
3. **Gather Info** — Collect Service Tree IDs, timeline, API version. See [details](references/release-plan-details.md).
4. **Create** — Run `azsdk_create_release_plan`.
5. **SDK Details** — Map emitters to languages.
6. **Namespace** — For mgmt plane first releases, link approval issue.
7. **Link PRs** — Link SDK PRs to plan.

## MCP Prerequisites

Requires `azure-sdk-mcp` server. No CLI fallback available.
