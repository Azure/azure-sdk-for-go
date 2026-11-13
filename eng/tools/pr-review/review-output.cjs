// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

const { assertReviewCurrent, publishReviewComment } = require("./generated-provenance.cjs");

function validErrors(errors) {
  return Array.isArray(errors) && errors.every((error) => typeof error === "string" && error.length > 0);
}

function ciFailuresBlock(review) {
  return review.pipeline.failures.length > 0 && !(review.mergeable === true && review.squash_merge_allowed === true &&
    ["clean", "unstable"].includes(review.mergeable_state));
}

function assertReviewResults(review, provenance) {
  if (!review || review.pr_number !== provenance.pr_number || review.head_sha !== provenance.head_sha ||
      review.base_sha !== provenance.base_sha || review.base_ref !== provenance.base_ref ||
      !["completed", "pending", "unverifiable"].includes(review.status) ||
      !validErrors(review.errors) || !Array.isArray(review.modules) ||
      !validErrors(review.mergeability_errors) ||
      !Array.isArray(review.onboarding) || !Array.isArray(review.breaking_changes) ||
      (review.status === "unverifiable" && review.errors.length === 0) ||
      (review.status !== "unverifiable" && (review.errors.length !== 0 || review.draft?.status !== "ready")) ||
      (review.status === "pending" && review.pipeline?.status !== "pending") ||
      (review.status === "completed" && review.pipeline?.status !== "completed")) {
    throw new Error("Missing, inconsistent, or stale deterministic review results; no outputs may be published.");
  }
  if (review.status === "completed") {
    const paths = review.modules.map((module) => module.path);
    if (new Set(paths).size !== paths.length ||
        review.modules.some((module) => !/^sdk\/resourcemanager\/[^/]+\/arm[^/]+$/.test(module.path) ||
          typeof module.first_onboard !== "boolean")) {
      throw new Error("The deterministic affected-module list is invalid.");
    }
    const newPaths = review.modules.filter((module) => module.first_onboard).map((module) => module.path);
    const samePaths = (items, expected) => items.length === expected.length &&
      new Set(items.map((item) => item.path)).size === items.length &&
      items.every((item) => expected.includes(item.path));
    const unavailable = (item) => item.status === "unverifiable" &&
      (!validErrors(item.errors) || item.errors.length === 0);
    if (!Array.isArray(review.pipeline.failures) ||
        review.pipeline.classification_required !== ciFailuresBlock(review) ||
        review.pipeline.failures.some((check) => typeof check.name !== "string" || !check.name ||
          !["failure", "cancelled", "timed_out", "action_required", "stale"].includes(check.conclusion) ||
          !["build-test", "analyze", "matrix", "other"].includes(check.category) ||
          typeof check.title !== "string" || typeof check.summary !== "string" ||
          (check.target_url !== null && typeof check.target_url !== "string")) ||
        ((typeof review.mergeable !== "boolean" || typeof review.mergeable_state !== "string" ||
          review.mergeable_state === "unknown") && review.mergeability_errors.length === 0) ||
        !samePaths(review.breaking_changes, paths) || !samePaths(review.onboarding, newPaths) ||
        review.breaking_changes.some((item) => !["skipped", "approved", "unapproved", "unverifiable"].includes(item.status) ||
          !Array.isArray(item.items) || !validErrors(item.errors) ||
          item.items.some((change) => typeof change.text !== "string" || !change.text ||
            typeof change.reason !== "string" || !["Approved", "Unapproved", "Unverifiable"].includes(change.status) ||
            (change.status !== "Approved" && !change.reason)) ||
          (item.status === "skipped" && (item.items.length > 0 || item.errors.length > 0)) ||
          (item.status === "approved" && (item.items.length === 0 || item.errors.length > 0 ||
            item.items.some((change) => change.status !== "Approved"))) ||
          (item.status === "unapproved" && (item.errors.length > 0 ||
            !item.items.some((change) => change.status === "Unapproved"))) ||
          (item.status === "unverifiable" && item.errors.length === 0)) ||
        review.onboarding.some((item) => !["satisfied", "missing", "unverifiable"].includes(item.pipeline?.status) ||
          !["satisfied", "missing", "unverifiable"].includes(item.namespace?.status) ||
          !["semantic-review-required", "unverifiable"].includes(item.readme?.status) ||
          unavailable(item.pipeline) || unavailable(item.namespace) || unavailable(item.readme) ||
          (item.pipeline.status !== "unverifiable" && (!Array.isArray(item.pipeline.definition_ids) ||
            item.pipeline.definition_ids.some((id) => !Number.isSafeInteger(id) || id <= 0) ||
            (item.pipeline.status === "satisfied") !== (item.pipeline.definition_ids.length > 0))) ||
          (item.namespace.status === "satisfied" && (typeof item.namespace.issue !== "string" || !item.namespace.issue.trim())) ||
          (item.readme.status === "semantic-review-required" && typeof item.readme.text !== "string"))) {
      throw new Error("The deterministic per-module review results are incomplete or invalid.");
    }
  }
}

function hasBlockers(review) {
  return review.status !== "completed" || review.errors.length > 0 ||
    review.mergeability_errors.length > 0 ||
    review.mergeable === false || review.mergeable_state === "blocked" ||
    ciFailuresBlock(review) ||
    review.breaking_changes.some((item) => !["skipped", "approved"].includes(item.status)) ||
    review.onboarding.some((item) => item.pipeline.status !== "satisfied" || item.namespace.status !== "satisfied" ||
      item.readme.status === "unverifiable");
}

function assertAgentOutput(review, output) {
  const blocked = hasBlockers(review);
  const visit = (value) => {
    if (!value || typeof value !== "object") return;
    for (const [key, child] of Object.entries(value)) {
      if (key === "body" && typeof child === "string" && /(^|\n)#{1,6}\s+[^\n]*\bPR is ready to merge\b/i.test(child) &&
          (blocked || review.onboarding.some((item) => !child.includes(`${item.path}/README.md`)))) {
        throw new Error("The agent claimed merge readiness despite deterministic review blockers.");
      }
      visit(child);
    }
  };
  visit(output);
}

function markdown(value) {
  return String(value).replace(/[\\`*_{}[\]()<>!|#]/g, "\\$&").replace(/[\r\n]/g, " ");
}

function formatReviewComment(review, body = "") {
  const marker = "<!-- gh-aw-workflow-id: mgmt-review -->";
  if (review.status === "pending") {
    return `${marker}\n## ⏳ Pipeline Checks Still Running\n\nThe \`go - pullrequest\` pipeline checks have not completed yet. Analysis cannot proceed until all checks finish.\n\n**Action required:** Re-trigger this workflow by re-applying the \`mgmt-review-needed\` label after the pipeline checks have completed.`;
  }
  if (review.status === "unverifiable") {
    return `${marker}\n## Deterministic Review - Manual Review Required\n\n${review.errors.map((error) => `- ${markdown(error)}`).join("\n")}\n\n**Action required:** Resolve the verification failure or request maintainer review, then re-apply \`mgmt-review-needed\`. This PR is not ready to merge.`;
  }
  const lines = [marker];
  for (const item of review.breaking_changes.filter((item) => item.status !== "skipped")) {
    lines.push("", "## Breaking-Change Review", "",
      `**Module:** \`${markdown(item.path)}\`${item.version ? ` (\`${markdown(item.version)}\`)` : ""}`, "");
    if (item.items.length) {
      lines.push("| Breaking change | Status | Reason |", "| --- | --- | --- |");
      for (const change of item.items) {
        lines.push(`| ${markdown(change.text)} | ${change.status} | ${markdown(change.reason)} |`);
      }
    }
    for (const error of item.errors) lines.push(`- Verification failure: ${markdown(error)}`);
    const url = item.source_missing ? item.source_url.replace("/blob/", "/tree/").replace(/\/sdk-suppressions\.yaml$/, "") : item.source_url;
    if (url) lines.push("", `Approval source: [${item.source_missing ? "pinned specification directory" : "sdk-suppressions.yaml"}](${url.replace(/[()]/g, (char) => `%${char.charCodeAt(0).toString(16)}`)}).`);
    if (item.status === "approved") lines.push("", "All breaking changes are approved.");
    else lines.push("", "**Action required:** A reviewer must manually review each unapproved or unverifiable item and confirm the required Go SDK breaking-change approval.");
  }
  const outstanding = [];
  for (const item of review.onboarding) {
    const module = `\`${markdown(item.path)}\``;
    if (item.pipeline.status === "missing") outstanding.push(`- **Pipeline setup** (${module}): comment \`/azp run prepare-pipelines\` on this PR.`);
    if (item.pipeline.status === "unverifiable") outstanding.push(`- **Pipeline setup** (${module}): ${item.pipeline.errors.map(markdown).join("; ")}. A maintainer must verify the release pipeline.`);
    if (item.namespace.status === "missing") outstanding.push(`- **Namespace approval** (${module}): no approval was established by the available Release Plan or namespace-related PR comment signals. Get namespace approval and comment its GitHub issue or PR link.`);
    if (item.namespace.status === "unverifiable") outstanding.push(`- **Namespace approval** (${module}): ${item.namespace.errors.map(markdown).join("; ")}. A maintainer must verify approval.`);
    if (item.readme.status === "unverifiable") outstanding.push(`- **README service name** (\`${markdown(item.path)}/README.md\`): ${item.readme.errors.map(markdown).join("; ")}. Resolve the access failure before naming review.`);
  }
  if (review.mergeability_errors.length) lines.push("", "## Mergeability Verification", "",
    ...review.mergeability_errors.map((error) => `- ${markdown(error)}`),
    "", "**Action required:** A maintainer must verify the mergeability information before declaring readiness.");
  if (review.mergeable === false) lines.push("", "## Merge Conflicts", "",
    "**Action required:** Resolve the conflicts with the target branch before merging.");
  if (review.mergeable_state === "blocked") lines.push("", "## Merging Is Blocked", "",
    "**Action required:** Resolve GitHub's outstanding merge requirements before merging.");
  const semantic = body.replaceAll(marker, "").trim();
  if (semantic) lines.push("", semantic);
  if (outstanding.length) {
    const heading = "## ⚠️ First On-Board Service — Action Required";
    if (!semantic.split("\n").includes(heading)) lines.push("", heading, "",
      "This PR cannot be merged until the following onboarding requirements are completed:", "");
    lines.push(...outstanding);
  }
  return lines.join("\n");
}

function prepareAgentOutput(review, output) {
  assertAgentOutput(review, output);
  const items = Array.isArray(output) ? output : output?.items;
  const labelNames = { remove_labels: "mgmt-review-in-progress", add_labels: "mgmt-review-analyzed" };
  if (!Array.isArray(items) || items.filter((item) => item?.type === "add_comment").length !== 1 ||
      items.some((item) => item?.type === "add_comment" ?
        typeof item.body !== "string" || !item.body.trim() :
        !Object.hasOwn(labelNames, item?.type) || !Array.isArray(item.labels) ||
        item.labels.length !== 1 || item.labels[0] !== labelNames[item.type])) {
    throw new Error("The agent must supply exactly one comment and only the permitted completion-label actions.");
  }
  const comment = items.find((item) => item.type === "add_comment");
  const prepared = [
    { ...comment, body: formatReviewComment(review, comment.body) },
    { type: "remove_labels", labels: [labelNames.remove_labels] },
    { type: "add_labels", labels: [labelNames.add_labels] },
  ];
  return Array.isArray(output) ? prepared : { ...output, items: prepared };
}

async function publishReviewFailure(api, { review, provenance, prNumber, headSHA, jobResult }) {
  if (provenance.status !== "pass") throw new Error("Deterministic review reporting requires passing provenance.");
  if (!review) {
    const message = `The trusted review_checks job ended ${jobResult} without producing review results. Inspect its workflow logs and request a fresh review.`;
    api.core.error(message);
    review = {
      pr_number: prNumber, head_sha: headSHA, base_sha: provenance.base_sha, base_ref: provenance.base_ref,
      status: "unverifiable", errors: [message], modules: [], onboarding: [], breaking_changes: [], mergeability_errors: [],
    };
  }
  assertReviewResults(review, provenance);
  await assertReviewCurrent(api, provenance, prNumber, headSHA);
  if (review.status === "completed") throw new Error("Completed reviews must use the agent publisher.");
  await publishReviewComment(api, prNumber, formatReviewComment(review));
}

module.exports = {
  assertReviewResults, assertAgentOutput, hasBlockers, ciFailuresBlock, formatReviewComment, prepareAgentOutput, publishReviewFailure,
};
