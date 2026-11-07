// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

const assert = require("node:assert/strict");
const { test } = require("node:test");
const fs = require("node:fs");
const path = require("node:path");
const os = require("node:os");
const crypto = require("node:crypto");
const { prepareRun, targetFromActivation, checkTarget, completedPlan, checkPlan, dispatch } =
  require("./issue-investigation-handoff.cjs");

const repository = { id: 10, full_name: "Azure/azure-sdk-for-go", default_branch: "main" };
const context = { repo: { owner: "Azure", repo: "azure-sdk-for-go" }, payload: { workflow_run: { id: 100, run_attempt: 1 } } };
const originalRun = () => ({
  id: 100, run_attempt: 1, head_sha: "a".repeat(40), head_branch: "main",
  name: "Agentic Triage", path: ".github/workflows/issue-triage.lock.yml",
  event: "issues", status: "completed", conclusion: "success",
  repository: { id: 10 }, head_repository: { id: 10 },
});
const proof = () => ({
  id: 100, attempt: 1, sha: "a".repeat(40), event: "issues", repository: repository.full_name,
  repositoryId: 10, ref: "refs/heads/main", mentions: "success",
});
const info = () => ({
  run_id: 100, run_attempt: "1", repository: repository.full_name, sha: "a".repeat(40),
  ref: "refs/heads/main", event_name: "issues", workflow_name: "Agentic Triage", staged: false,
});
const prompt = number => `# Agentic Triage\nYour task is to analyze issue #${number} and perform initial triage following the decision flow below\n`;
const issue = () => ({ state: "open", labels: [
  { name: "customer-reported", color: "3800e0" }, { name: "KeyVault", color: "e99695" },
  { name: "Client", color: "ffeb77" },
] });
const comment = () => ({ type: "add_comment", item_number: 42, body: "## Agentic Issue Triage\nAnalysis" });
const mention = () => ({ type: "mention_owners", owners: "owner1, owner2", message: "Routing" });
const owner = () => ({ type: "assign_to_user", issue_number: 42, assignees: ["owner1"] });
const plan = () => ({ items: [owner(), mention(), comment()], errors: [] });
const clean = () => ({ verdict: { prompt_injection: false, secret_leak: false, malicious_patch: false } });
const jobs = () => ["activation", "agent", "detection", "safe_outputs", "conclusion", "mention_owners"]
  .map(name => ({ name, conclusion: "success" }));
const core = { info() {} };
function client() {
  return {
    rest: {
      repos: { async get() { return { data: repository }; } },
      actions: {
        async getWorkflowRun() { return { data: originalRun() }; },
        listJobsForWorkflowRunAttempt() {},
      },
      issues: { async get() { return { data: issue() }; } },
    },
    async paginate(_method, args) { assert.equal(args.attempt_number, 1); return jobs(); },
  };
}

// This golden baseline is deliberately the existing upstream workflow at fa71c24.
// An intentional triage change requires an explicit baseline update, not silent regeneration.
test("original triage source and generated runtime remain byte-equivalent to upstream", () => {
  const dir = path.join(__dirname, "..", "workflows");
  const hash = file => crypto.createHash("sha256")
    .update(fs.readFileSync(path.join(dir, file), "utf8").replaceAll("\r\n", "\n")).digest("hex");
  assert.equal(hash("issue-triage.md"), "f0ae5922221816f6301820bb8e0001fffc44e9a56a52ec85b28b838adf099ef5");
  assert.equal(hash("issue-triage.lock.yml"), "ea9d684beaa485a4524ee930e8e68b05de9d4ccd26bb9e029b116f5497265e33");
});

test("new handoff is opt-in, independent and does not grant issue writes", () => {
  const file = fs.readFileSync(path.join(__dirname, "..", "workflows", "issue-investigation-handoff.yml"), "utf8");
  assert.match(file, /workflow_run:\s+workflows: \["Agentic Triage"\]\s+types: \[completed\]/);
  assert.match(file, /if: vars\.ENABLE_ISSUE_INVESTIGATION == 'true'/);
  assert.doesNotMatch(file, /issues: write|contents: write|needs:/);
  assert.match(file, /ref: \$\{\{ github\.workflow_sha \}\}/);
  assert.match(file, /persist-credentials: false/);
});

test("preflight accepts only successful same-repository default-branch triage attempts", async () => {
  const github = client();
  assert.deepEqual(await prepareRun({ github, context, core }), proof());
  for (const change of [{ event: "pull_request" }, { path: ".github/workflows/other.yml" },
    { name: "Other" }, { head_branch: "feature" }, { status: "in_progress" },
    { conclusion: "failure" }, { run_attempt: 2 }, { repository: { id: 11 } }, { head_repository: { id: 11 } }]) {
    github.rest.actions.getWorkflowRun = async () => ({ data: { ...originalRun(), ...change } });
    assert.equal(await prepareRun({ github, context, core }), null);
  }
  github.rest.actions.getWorkflowRun = async () => ({ data: originalRun() });
  for (const name of ["activation", "agent", "detection", "safe_outputs", "conclusion", "mention_owners"]) {
    github.paginate = async () => jobs().map(job => job.name === name ? { ...job, conclusion: "failure" } : job);
    assert.equal(await prepareRun({ github, context, core }), null);
  }
  github.paginate = async () => jobs().map(job => job.name === "mention_owners" ? { ...job, conclusion: "skipped" } : job);
  assert.equal((await prepareRun({ github, context, core })).mentions, "skipped");
  github.rest.actions.getWorkflowRun = async () => ({ data: { ...originalRun(), event: "workflow_dispatch" } });
  assert.equal((await prepareRun({ github, context, core })).event, "workflow_dispatch");
});

test("target comes only from canonical activation text and matching provenance", () => {
  assert.equal(targetFromActivation(info(), prompt(42), proof()), 42);
  assert.equal(targetFromActivation({ ...info(), repository: "azure/AZURE-SDK-FOR-GO" }, prompt(42), proof()), 42);
  for (const change of [{ repository: "Azure/other" }, { sha: "b".repeat(40) }, { ref: "refs/heads/feature" },
    { run_id: 101 }, { run_attempt: 2 }, { event_name: "pull_request" }, { workflow_name: "Other" }, { staged: true }]) {
    assert.throws(() => targetFromActivation({ ...info(), ...change }, prompt(42), proof()), /does not match/);
  }
  for (const text of ["", prompt(0), prompt("42junk"), prompt(42) + prompt(43), prompt("9007199254740992")]) {
    assert.throws(() => targetFromActivation(info(), text, proof()));
  }
});

test("existing fallback, closed and noop plans do not trigger investigation", () => {
  assert.deepEqual(completedPlan(plan(), 42, repository.full_name, "success"), plan());
  assert.ok(completedPlan({ items: [owner(), comment(), comment()] }, 42, repository.full_name, "skipped"));
  for (const items of [[], [{ type: "noop", message: "Already labeled" }],
    [{ type: "close_issue", issue_number: 42, body: "Closed" }], [comment()],
    [mention(), comment(), comment()], [owner(), comment()], [mention(), { ...comment(), body: "Generic acknowledgment" }]]) {
    assert.equal(completedPlan({ items }, 42, repository.full_name, "success"), null);
  }
  assert.equal(completedPlan(plan(), 42, repository.full_name, "skipped"), null);
  for (const output of [null, {}, { items: [] , errors: {} }, { ...plan(), errors: ["Invalid item"] },
    { items: [mention(), { ...comment(), item_number: 43 }] },
    { items: [mention(), { ...comment(), repo: "Azure/other" }] }, { items: [{ type: "dispatch_workflow" }] }]) {
    assert.throws(() => completedPlan(output, 42, repository.full_name, "success"));
  }
});

test("activation and plan artifacts are data and missing artifacts fail explicitly", async () => {
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), "handoff-artifacts-"));
  try {
    fs.mkdirSync(path.join(directory, "aw-prompts"));
    fs.writeFileSync(path.join(directory, "aw_info.json"), JSON.stringify(info()));
    fs.writeFileSync(path.join(directory, "aw-prompts", "prompt.txt"), prompt(42));
    fs.writeFileSync(path.join(directory, "agent_output.json"), JSON.stringify(plan()));
    const github = client();
    assert.equal(await checkTarget({ github, context, core, run: proof(), directory }), 42);
    assert.deepEqual(checkPlan({ core, run: proof(), number: 42, directory }), plan());
    github.rest.issues.get = async () => ({ data: { ...issue(), state: "closed" } });
    assert.equal(await checkTarget({ github, context, core, run: proof(), directory }), null);
    fs.writeFileSync(path.join(directory, "agent_output.json"), JSON.stringify({ items: [{ type: "noop" }] }));
    assert.equal(checkPlan({ core, run: proof(), number: 42, directory }), null);
    await assert.rejects(checkTarget({ github, context, core, run: proof(), directory: path.join(directory, "missing") }), /ENOENT/);
  } finally {
    fs.rmSync(directory, { recursive: true, force: true });
  }
});

test("native dispatch always receives an explicit verified default-branch ref", async () => {
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), "handoff-dispatch-"));
  try {
    const detectionFile = path.join(directory, "detection.log");
    fs.writeFileSync(detectionFile, "Native detection evidence");
    const github = client();
    const calls = [];
    const args = { github, context, core, run: proof(), number: 42, detectionFile,
      parseDetectionLog: clean,
      async createDispatchHandler(config) {
        calls.push(config);
        return async item => { calls.push(item); return { success: true }; };
      } };
    assert.equal((await dispatch(args)).success, true);
    assert.deepEqual(calls[0].allowed_refs, ["refs/heads/main"]);
    assert.equal(calls[0]["target-ref"], "refs/heads/main");
    assert.deepEqual(calls[1], { workflow_name: "issue-investigation", ref: "refs/heads/main", inputs: { issue_number: "42" } });
    for (const result of [{ error: "Parse failed" }, { verdict: {} },
      { verdict: { ...clean().verdict, prompt_injection: true } },
      { verdict: { ...clean().verdict, secret_leak: true } },
      { verdict: { ...clean().verdict, malicious_patch: true } }]) {
      await assert.rejects(dispatch({ ...args, parseDetectionLog: () => result }), /clean threat/);
    }
    assert.equal(calls.length, 2);
    for (const ref of [undefined, "refs/heads/feature", "main"]) {
      assert.equal(await dispatch({ ...args, run: { ...proof(), ref } }), null);
    }
    assert.equal(calls.length, 2);
    github.rest.issues.get = async () => ({ data: { ...issue(), state: "closed" } });
    assert.equal(await dispatch(args), null);
    github.rest.issues.get = async () => ({ data: issue() });
    for (const change of [{ run_attempt: 2 }, { head_sha: "b".repeat(40) }, { status: "in_progress" }, { conclusion: "failure" }]) {
      github.rest.actions.getWorkflowRun = async () => ({ data: { ...originalRun(), ...change } });
      assert.equal(await dispatch(args), null);
    }
    github.rest.actions.getWorkflowRun = async () => ({ data: originalRun() });
    github.rest.repos.get = async () => ({ data: { ...repository, default_branch: "new-default" } });
    assert.equal(await dispatch(args), null);
    github.rest.repos.get = async () => ({ data: repository });
    for (const result of [{ success: false }, { success: true, staged: true }, { success: true, skipped: true }]) {
      await assert.rejects(dispatch({ ...args, createDispatchHandler: async () => async () => result }), /dispatch failed/);
    }
    const branch = "stable/triage";
    const newRun = { ...proof(), ref: `refs/heads/${branch}` };
    github.rest.repos.get = async () => ({ data: { ...repository, default_branch: branch } });
    github.rest.actions.getWorkflowRun = async () => ({ data: { ...originalRun(), head_branch: branch } });
    await dispatch({ ...args, run: newRun });
    assert.equal(calls.at(-1).ref, "refs/heads/stable/triage");
  } finally {
    fs.rmSync(directory, { recursive: true, force: true });
  }
});

test("API failures propagate without modifying the completed upstream workflow", async () => {
  const github = client();
  github.rest.repos.get = async () => { throw new Error("API unavailable"); };
  await assert.rejects(prepareRun({ github, context, core }), /API unavailable/);
  const yaml = fs.readFileSync(path.join(__dirname, "..", "workflows", "issue-triage.lock.yml"), "utf8");
  assert.doesNotMatch(yaml, /issue-investigation|ENABLE_ISSUE_INVESTIGATION/);
});