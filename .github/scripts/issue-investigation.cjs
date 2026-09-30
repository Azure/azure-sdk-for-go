// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

const fs = require("node:fs");

const blockedLabels = [
  "needs-triage",
  "needs-team-triage",
  "issue-addressed",
  "needs-author-feedback",
];

function issueNumber(value) {
  const text = String(value);
  if ((typeof value !== "string" && typeof value !== "number") ||
      !/^[1-9]\d*$/.test(text) || !Number.isSafeInteger(Number(text))) {
    throw new Error("issue_number must be a positive safe integer");
  }
  return Number(text);
}

function isEligible(issue) {
  const labels = issue.labels;
  if (!Array.isArray(labels) || labels.some(label =>
    !label || typeof label.name !== "string" || typeof label.color !== "string")) {
    throw new Error("Issue labels must include names and colors");
  }
  const names = new Set(labels.map(label => label.name.toLowerCase()));
  const countColor = color => labels.filter(label =>
    label.color.toLowerCase() === color).length;
  return !issue.pull_request && issue.state === "open" &&
    countColor("e99695") === 1 && countColor("ffeb77") === 1 &&
    names.has("customer-reported") && !blockedLabels.some(label => names.has(label));
}

function validateOutputs(output, expectedIssue, repository) {
  const expected = issueNumber(expectedIssue);
  if (!output || !Array.isArray(output.items) || output.items.length === 0) {
    throw new Error("Expected a nonempty safe-output items array");
  }
  if (output.errors != null && (!Array.isArray(output.errors) || output.errors.length > 0)) {
    throw new Error("Cannot apply a partially collected investigation plan");
  }
  const targetFields = {
    add_comment: "item_number",
    close_issue: "issue_number",
    assign_to_agent: "issue_number",
  };
  const alternateTargets = [
    "item_number", "issue_number", "pull_number", "pull_request_number",
    "pr_number", "pr", "pr-number", "comment_id", "commentId", "comment-id",
    "reply_to_id", "target",
  ];
  const counts = new Map();
  for (const item of output.items) {
    if (!item || typeof item !== "object") {
      throw new Error("Invalid safe-output item");
    }
    const field = targetFields[item.type];
    if (!Object.hasOwn(targetFields, item.type) && item.type !== "noop") {
      throw new Error("Unexpected investigation output type");
    }
    counts.set(item.type, (counts.get(item.type) || 0) + 1);
    if (counts.get(item.type) > 1) {
      throw new Error("Only one output of each type is allowed");
    }
    if (item.type === "noop") {
      if (typeof item.message !== "string" || !item.message.trim()) {
        throw new Error("A noop requires an explanation");
      }
      continue;
    }
    if (item.suggest === true) {
      throw new Error("Investigation outputs must not request pending suggestions");
    }
    // Native handlers accept explicit targets in preference to configured targets.
    if (issueNumber(item[field]) !== expected ||
        alternateTargets.some(key => key !== field && item[key] != null)) {
      throw new Error("Investigation output must target only the dispatched issue");
    }
    for (const key of ["repo", "pull_request_repo"]) {
      if (item[key] != null &&
          (typeof item[key] !== "string" ||
           item[key].toLowerCase() !== repository.toLowerCase())) {
        throw new Error("Investigation output must target only the current repository");
      }
    }
    if ((item.type === "add_comment" || item.type === "close_issue") &&
        (typeof item.body !== "string" || !item.body.trim())) {
      throw new Error("Investigation decisions require an explanatory comment");
    }
    if (item.type === "close_issue" &&
        ((item.state_reason != null && item.state_reason !== "not_planned") ||
         item.duplicate_of != null)) {
      throw new Error("Investigation may only close service-controlled issues as not planned");
    }
  }
  if (counts.has("noop") && output.items.length !== 1) {
    throw new Error("A noop cannot be combined with issue changes");
  }
  if (counts.has("close_issue") && output.items.length !== 1) {
    throw new Error("Use close_issue with its explanation, without other outputs");
  }
  if (counts.has("assign_to_agent") &&
      (output.items.length !== 2 || output.items[0].type !== "add_comment")) {
    throw new Error("Copilot assignment requires a preceding investigation comment");
  }
}

async function checkEligibility({ github, context, core, number }) {
  const { data: issue } = await github.rest.issues.get({
    ...context.repo, issue_number: issueNumber(number),
  });
  const eligible = isEligible(issue);
  core.info(eligible ? "Issue is eligible for investigation" :
    "Skipping investigation: issue is closed, is a pull request, or has an incomplete triage handoff");
  return eligible;
}

async function checkOutputs({ github, context, core, number, outputFile }) {
  const output = JSON.parse(fs.readFileSync(outputFile, "utf8"));
  validateOutputs(output, number, `${context.repo.owner}/${context.repo.repo}`);
  if (output.items[0].type === "noop") return output;
  if (!await checkEligibility({ github, context, core, number })) {
    throw new Error("Issue eligibility changed during investigation; no outputs will be applied");
  }
  return output;
}

function nativeCommentPolicy(workflowFile, expectedIssue) {
  const workflow = fs.readFileSync(workflowFile, "utf8");
  // The pinned compiler emits these environment values as JSON-quoted YAML scalars.
  // Reuse its policy instead of maintaining a narrower, second domain allowlist.
  function literal(name) {
    const pattern = new RegExp(`^ +${name}: ("[^\\r\\n]*")\\r?$`, "gm");
    const values = new Set([...workflow.matchAll(pattern)].map(match => JSON.parse(match[1])));
    if (values.size !== 1) throw new Error(`Missing or inconsistent native ${name} policy`);
    return [...values][0];
  }
  const config = JSON.parse(literal("GH_AW_SAFE_OUTPUTS_HANDLER_CONFIG")).add_comment;
  if (!config || config.target !== "${{ github.event.inputs.issue_number }}" || config.max !== 1) {
    throw new Error("Unexpected native investigation comment policy");
  }
  return {
    domains: literal("GH_AW_ALLOWED_DOMAINS"),
    config: { ...config, target: String(issueNumber(expectedIssue)) },
  };
}

async function prepareOutputs({ postComment, staged = false, ...args }) {
  const output = await checkOutputs(args);
  if (!output.items.some(item => item.type === "assign_to_agent") || staged) return;
  if (typeof postComment !== "function") {
    throw new Error("Native comment handler is required before assignment");
  }
  // Native safe outputs continue after individual failures. Deliver the prerequisite
  // with the native handler here, then leave only assignment for normal processing.
  const result = await postComment(output.items[0]);
  if (!result || result.success !== true || result.skipped || result.staged || !result.commentId) {
    throw new Error("Investigation comment was not delivered; assignment is blocked");
  }
  if (!await checkEligibility(args)) {
    throw new Error("Issue eligibility changed after comment delivery; assignment is blocked");
  }
  args.core.setOutput("delivered_comment_id", String(result.commentId));
  output.items = output.items.slice(1);
  fs.writeFileSync(args.outputFile, JSON.stringify(output));
}

function validateTriageOutputs(output, expectedIssue, repository) {
  if (!output || !Array.isArray(output.items)) {
    throw new Error("Expected a triage safe-output items array");
  }
  const dispatches = output.items.filter(item => item?.type === "dispatch_workflow");
  if (dispatches.length === 0) return;
  if (output.errors != null && (!Array.isArray(output.errors) || output.errors.length > 0)) {
    throw new Error("Cannot dispatch a partially collected triage plan");
  }
  const expected = issueNumber(expectedIssue);
  const dispatch = dispatches[0];
  if (dispatches.length !== 1 || output.items.at(-1) !== dispatch ||
      dispatch.workflow_name !== "issue-investigation" ||
      issueNumber(dispatch.inputs?.issue_number) !== expected) {
    throw new Error("Investigation dispatch must be last and target the triaged issue");
  }
  const mentions = output.items.filter(item => item?.type === "mention_owners");
  const assignments = output.items.filter(item => item?.type === "assign_to_user");
  const comments = output.items.filter(item => item?.type === "add_comment");
  if (assignments.length > 1 ||
      !((mentions.length === 1 && comments.length === 1) ||
        (mentions.length === 0 && assignments.length === 1 && comments.length === 2))) {
    throw new Error("Investigation dispatch requires an owner-routing plan");
  }
  if (comments.length === 0 || comments.some(item => typeof item.body !== "string" || !item.body.trim())) {
    throw new Error("Investigation dispatch requires a triage explanation");
  }
  for (const item of output.items) {
    if (!item || typeof item !== "object") throw new Error("Invalid triage output");
    if (item.repo != null &&
        (typeof item.repo !== "string" || item.repo.toLowerCase() !== repository.toLowerCase())) {
      throw new Error("Triage handoff must stay in the current repository");
    }
    for (const key of ["item_number", "issue_number", "pull_number", "pull_request_number", "pr_number", "pr"]) {
      if (item[key] != null && issueNumber(item[key]) !== expected) {
        throw new Error("Triage handoff must target only the triaged issue");
      }
    }
    if (["close_issue", "noop"].includes(item.type)) {
      throw new Error("Closed or no-action triage plans cannot dispatch investigation");
    }
  }
}

module.exports = {
  issueNumber, isEligible, validateOutputs, checkEligibility, checkOutputs,
  prepareOutputs, validateTriageOutputs, nativeCommentPolicy,
};
