// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

const fs = require("node:fs");
const path = require("node:path");
const { issueNumber, isEligible } = require("./issue-investigation.cjs");

function skip(core, reason) {
  core.info(`Skipping investigation handoff: ${reason}`);
  return null;
}

async function prepareRun({ github, context, core }) {
  const trigger = context.payload.workflow_run;
  const { data: repository } = await github.rest.repos.get(context.repo);
  const { data: run } = await github.rest.actions.getWorkflowRun({
    ...context.repo, run_id: issueNumber(trigger.id),
  });
  if (run.repository?.id !== repository.id || run.head_repository?.id !== repository.id ||
      run.path !== ".github/workflows/issue-triage.lock.yml" || run.name !== "Agentic Triage" ||
      !["issues", "workflow_dispatch"].includes(run.event) ||
      run.head_branch !== repository.default_branch || run.status !== "completed" ||
      run.conclusion !== "success" || run.run_attempt !== trigger.run_attempt) {
    return skip(core, "not a successful current-attempt triage run on this repository's default branch");
  }
  const jobs = await github.paginate(github.rest.actions.listJobsForWorkflowRunAttempt, {
    ...context.repo, run_id: run.id, attempt_number: run.run_attempt, per_page: 100,
  });
  for (const name of ["activation", "agent", "detection", "safe_outputs", "conclusion"]) {
    const matches = jobs.filter(job => job.name === name);
    if (matches.length !== 1 || matches[0].conclusion !== "success") {
      return skip(core, `triage job ${name} did not complete successfully`);
    }
  }
  const mentions = jobs.filter(job => job.name === "mention_owners");
  if (mentions.length !== 1 || !["success", "skipped"].includes(mentions[0].conclusion)) {
    return skip(core, "owner notification did not complete successfully");
  }
  return {
    id: run.id, attempt: run.run_attempt, sha: run.head_sha, event: run.event,
    repository: repository.full_name, repositoryId: repository.id,
    ref: `refs/heads/${repository.default_branch}`,
    mentions: mentions[0].conclusion,
  };
}

function targetFromActivation(info, prompt, run) {
  if (!info || typeof info.repository !== "string" || info.repository.toLowerCase() !== run.repository.toLowerCase() ||
      info.sha !== run.sha || info.ref !== run.ref ||
      issueNumber(info.run_id) !== run.id || issueNumber(info.run_attempt) !== run.attempt ||
      info.event_name !== run.event || info.workflow_name !== "Agentic Triage" || info.staged !== false) {
    throw new Error("Activation artifact does not match the completed triage run");
  }
  // This is the existing workflow's canonical target line, captured by activation
  // before the agent runs. Never infer the triggering issue from agent output.
  const targets = [...prompt.matchAll(/^Your task is to analyze issue #([1-9]\d*) and perform initial triage following the decision flow below\r?$/gm)];
  if (targets.length !== 1) {
    throw new Error("Activation prompt must identify exactly one canonical triage target");
  }
  return issueNumber(targets[0][1]);
}

async function checkTarget({ github, context, core, run, directory }) {
  const info = JSON.parse(fs.readFileSync(path.join(directory, "aw_info.json"), "utf8"));
  const prompt = fs.readFileSync(path.join(directory, "aw-prompts", "prompt.txt"), "utf8");
  const number = targetFromActivation(info, prompt, run);
  const { data: issue } = await github.rest.issues.get({ ...context.repo, issue_number: number });
  return isEligible(issue) ? number : skip(core, "issue does not have a clean investigation handoff");
}

function completedPlan(output, number, repository, mentionResult) {
  if (!output || !Array.isArray(output.items) ||
      (output.errors != null && (!Array.isArray(output.errors) || output.errors.length > 0))) {
    throw new Error("Cannot investigate a partially collected triage plan");
  }
  if (output.items.length === 0 || output.items.some(item => ["noop", "close_issue"].includes(item?.type))) {
    return null;
  }
  const allowed = new Set(["add_labels", "remove_labels", "add_comment", "assign_to_user", "mention_owners"]);
  for (const item of output.items) {
    if (!item || !allowed.has(item.type)) throw new Error("Unexpected completed triage output");
    if (item.repo != null &&
        (typeof item.repo !== "string" || item.repo.toLowerCase() !== repository.toLowerCase())) {
      throw new Error("Triage output targets another repository");
    }
    for (const key of ["item_number", "issue_number", "pull_number", "pull_request_number", "pr_number", "pr"]) {
      if (item[key] != null && issueNumber(item[key]) !== number) {
        throw new Error("Triage output does not match the activation target");
      }
    }
  }
  const mentions = output.items.filter(item => item.type === "mention_owners");
  const assignments = output.items.filter(item => item.type === "assign_to_user");
  const comments = output.items.filter(item => item.type === "add_comment");
  if (comments.some(item => typeof item.body !== "string" || !item.body.trim()) ||
      !comments.some(item => item.body.includes("Agentic Issue Triage")) || assignments.length > 1 ||
      !((mentions.length === 1 && comments.length === 1 && mentionResult === "success") ||
        (mentions.length === 0 && assignments.length === 1 && comments.length === 2 && mentionResult === "skipped"))) {
    return null;
  }
  return output;
}

function checkPlan({ core, run, number, directory }) {
  const output = JSON.parse(fs.readFileSync(path.join(directory, "agent_output.json"), "utf8"));
  return completedPlan(output, number, run.repository, run.mentions) ||
    skip(core, "triage did not produce a complete owner-routing and explanation plan");
}

async function dispatch({ github, context, core, run, number, detectionFile, parseDetectionLog, createDispatchHandler }) {
  const { verdict, error } = parseDetectionLog(fs.readFileSync(detectionFile, "utf8"));
  if (error || !verdict || ["prompt_injection", "secret_leak", "malicious_patch"].some(key => verdict[key] !== false)) {
    throw new Error("Investigation handoff requires a clean threat-detection verdict");
  }
  const { data: repository } = await github.rest.repos.get(context.repo);
  const { data: current } = await github.rest.actions.getWorkflowRun({ ...context.repo, run_id: run.id });
  if (repository.id !== run.repositoryId || `refs/heads/${repository.default_branch}` !== run.ref ||
      current.run_attempt !== run.attempt || current.head_sha !== run.sha ||
      current.head_branch !== repository.default_branch ||
      current.status !== "completed" || current.conclusion !== "success") {
    return skip(core, "triage provenance or default branch changed before dispatch");
  }
  const { data: issue } = await github.rest.issues.get({ ...context.repo, issue_number: issueNumber(number) });
  if (!isEligible(issue)) return skip(core, "issue eligibility changed before dispatch");
  const handler = await createDispatchHandler({
    workflows: ["issue-investigation"], workflow_files: { "issue-investigation": ".lock.yml" },
    aw_context_workflows: ["issue-investigation"], max: 1,
    allowed_refs: [run.ref], "target-ref": run.ref,
  });
  const result = await handler({
    workflow_name: "issue-investigation", ref: run.ref, inputs: { issue_number: String(number) },
  }, {});
  if (!result || result.success !== true || result.skipped || result.staged) {
    throw new Error("Native investigation dispatch failed or was not applied");
  }
  core.info(`Investigation dispatched for issue #${number} from completed triage run ${run.id}`);
  return result;
}

module.exports = { prepareRun, targetFromActivation, checkTarget, completedPlan, checkPlan, dispatch };