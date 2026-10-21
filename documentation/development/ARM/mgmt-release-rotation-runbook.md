# Go mgmt. SDK release rotation runbook

Follow these steps for generated management SDK PRs under `sdk/resourcemanager`, ending with PR approval. **Review the Management Release Assistant result; do not repeat its CI, generated-file, breaking-change, or onboarding checks.** Human review also covers Copilot comments. Route handwritten or generator changes to the Go owner.

**Contacts:** the **service team** is the contact in the PR/release plan; the **Go owner** is the Go management SDK maintainer listed under `/sdk/resourcemanager/` in [CODEOWNERS](https://github.com/Azure/azure-sdk-for-go/blob/main/.github/CODEOWNERS). Ask the Go owner to coordinate vendor support when needed.

## Step 1. Start from the daily mail

Open **Azure SDK for Go management release PR check - YYYY-MM-DD**. Release owners can receive this mail directly or through forwarding by the Go owner. Keep using it until a shared cross-language dashboard provides equivalent coverage and the Go owner confirms the transition. Use the open-PR sections below.

| Mail section | Next action |
| --- | --- |
| `ON HOLD - HoldOn label` | Leave the PR on hold until the on hold problem is resolved. |
| `CI INCOMPLETE - go - pullrequest not completed` | Wait for the top-level PR pipeline to finish, then manually add `mgmt-review-needed` to trigger review analysis. Wait for the final workflow result, then continue with Step 2. If the pipeline is stuck or absent, ask the Go owner. |
| `READY - needs SDK team review` | Open the PR and continue with Step 2. Human review is still required. |
| `ACTION REQUIRED - breaking-change review` | Open the linked workflow comment and send its affected-module/status summary to the Go owner for a decision. Do not approve around it. |
| `NO ACTION - already approved` | No action is needed; do not submit a duplicate review. |
| `NO ACTION - needs service team follow-up` | Relay the existing workflow's required action to the service team; wait for its fix instead of restarting review. |
| `Pending / errors` | Send the reported reason to the Go owner for report/workflow troubleshooting. Missing analysis is not a clean review. |

The mail is a snapshot. Open the PR for current status; **Since last action** reflects PR activity, including bots, not how long a partner has waited. A missing/failed mail is not an empty queue: contact the Go owner and use the [live management AutoPR query](https://github.com/Azure/azure-sdk-for-go/pulls?q=is%3Apr+is%3Aopen+%22AutoPR+sdk-resourcemanager%2F%22+in%3Atitle) as a fallback.

## Step 2. Open the PR and find the current review result

Read only the information needed for the next steps:

- **Description / release plan:** module, target version, and service-team contact.
- **Conversation:** latest final Management Release Assistant comment and existing human decisions; ignore workflow-started/completed notifications.
- **Copilot review:** outstanding findings and their dispositions.
- **Labels:** `HoldOn` and workflow status labels.

The daily report already triggers or refreshes eligible analyses. If the result is missing, apply `mgmt-review-needed` after CI completes if authorized; otherwise ask the Go owner. Wait for a final result. Do not duplicate an active run, analyze a held/CI-incomplete PR, or treat `mgmt-review-analyzed` alone as approval.

## Step 3. Act on the workflow result

Use the [Management Release Assistant](https://github.com/Azure/azure-sdk-for-go/blob/main/.github/workflows/mgmt-review.md) findings and links rather than reconstructing its checks.

| Workflow result | Reviewer action |
| --- | --- |
| Pipeline Checks Still Running | Wait for CI, then request fresh analysis. |
| Next Steps to Merge | Ask the service team to apply the listed fixes; send pipeline infrastructure or generator failures to the Go owner. Refresh analysis after fixes; do not redo CI diagnosis. |
| Manual Edits to Auto-Generated Files Detected | Ask the manual edit author to revert the listed manual edits. The service team fixes source-spec issues and regenerates; the Go owner handles generator issues. Do not patch generated output. |
| Breaking-Change Review | Use the reported statuses. Escalate **Unapproved**, **Unverifiable**, or verification failures to the Go owner; do not repeat suppression matching or grant exceptions. All-approved items alone are not a ready verdict. |
| First On-Board Service - Action Required | Ask the service team to complete all listed namespace, pipeline, and README service-name requirements and resolve any verification failures. Refresh analysis after fixes. |
| PR is ready to merge | Continue with Step 4; this is advisory, not human approval. |

Do not approve while any reported blocker remains. Ask the Go owner to resolve incomplete or contradictory results. Keep the next action on the existing PR thread.

## Step 4. Review Copilot comments

Ask the service team to fix actionable spec/client-customization findings and regenerate; send generator issues or unclear findings to the Go owner. Do not make output-only edits. For intentional generated patterns, explain the disposition using the [generated-code guidance](https://github.com/Azure/azure-sdk-for-go/blob/main/.github/instructions/go-code.instructions.md#generated-code--do-not-suggest-changes-to-any-of-the-patterns-below). Resolve threads after a fix or justified decision.

### Common partner questions during these steps

| Question | Response |
| --- | --- |
| "CI is green / the bot says ready. Can you approve now?" | Explain any unresolved workflow findings or Copilot comments. Approve once the Step 5 requirements are met. |
| "Why is this waiting, and can the team unblock it?" | Point to the mail category and the service-team action or Go-owner escalation in Step 3. Do not bypass holds/checks or grant exceptions. |
| "Must we make every Copilot change?" | Fix real issues upstream; explain intentional patterns using the guidance in Step 4. Do not edit generated output just to resolve a thread. |

## Step 5. Approve

Approve when the workflow reports readiness, Copilot findings have justified dispositions, and no blocking human review or hold remains.