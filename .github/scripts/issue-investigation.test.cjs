// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

const assert = require("node:assert/strict");
const { test } = require("node:test");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const { issueNumber, isEligible, validateOutputs, checkEligibility, checkOutputs, prepareOutputs, validateTriageOutputs, nativeCommentPolicy } =
  require("./issue-investigation.cjs");

const repository = "Azure/azure-sdk-for-go";
const eligibleIssue = (category = "Client", service = "KeyVault") => ({
  state: "open",
  labels: [
    { name: category, color: "ffeb77" },
    { name: service, color: "e99695" },
    { name: "customer-reported", color: "3800e0" },
    { name: "question", color: "0052cc" },
    { name: "needs-team-attention", color: "ededed" },
  ],
});
const comment = () => ({ type: "add_comment", item_number: 42, body: "Investigation and next action" });
const close = () => ({ type: "close_issue", issue_number: 42, body: "Documented service behavior" });
const assign = () => ({ type: "assign_to_agent", issue_number: 42, agent: "copilot" });
const validate = items => validateOutputs({ items }, "42", repository);

test("accepts positive integer issue numbers and rejects ambiguous input", () => {
  for (const value of [1, "42", Number.MAX_SAFE_INTEGER]) {
    assert.equal(issueNumber(value), Number(value));
  }
  for (const value of [undefined, null, "", " ", " 42", "42 ", "01", 0, -1,
    "1.0", 1.5, "1e2", "42junk", "42\n", "1; echo bad", Infinity, 9007199254740992,
    ["42"], { toString: () => "42" }]) {
    assert.throws(() => issueNumber(value), /positive safe integer/);
  }
});

test("keeps every repository category and arbitrary service labels eligible", () => {
  for (const category of ["Central-EngSys", "Client", "Mgmt", "Mgmt-EngSys", "Provisioning", "Service"]) {
    for (const service of ["KeyVault", "Azure.Core", "Storage", "Compute", "Future Service"]) {
      assert.equal(isEligible(eligibleIssue(category, service)), true);
    }
  }
  const issue = eligibleIssue();
  issue.labels = issue.labels.map(label => ({
    name: label.name.toUpperCase(), color: label.color.toUpperCase(),
  }));
  assert.equal(isEligible(issue), true);
});

test("rejects closed issues, pull requests, missing and ambiguous handoffs", () => {
  assert.equal(isEligible({ ...eligibleIssue(), state: "closed" }), false);
  assert.equal(isEligible({ ...eligibleIssue(), pull_request: { url: "https://api.github.com/repos/Azure/azure-sdk-for-go/pulls/42" } }), false);
  for (const name of ["customer-reported", "KeyVault", "Client"]) {
    const issue = eligibleIssue();
    issue.labels = issue.labels.filter(label => label.name !== name);
    assert.equal(isEligible(issue), false);
  }
  for (const color of ["e99695", "ffeb77"]) {
    const issue = eligibleIssue();
    issue.labels.push({ name: "Another label", color });
    assert.equal(isEligible(issue), false);
  }
  for (const name of ["needs-triage", "needs-team-triage", "issue-addressed", "needs-author-feedback"]) {
    const issue = eligibleIssue();
    issue.labels.push({ name: name.toUpperCase(), color: "ededed" });
    assert.equal(isEligible(issue), false);
  }
  for (const labels of [null, ["Client"], [{ name: "Client" }]]) {
    assert.throws(() => isEligible({ state: "open", labels }), /names and colors/);
  }
});

test("accepts only the intended investigation outcome shapes", () => {
  validate([{ type: "noop", message: "Not enough evidence" }]);
  validate([comment()]);
  validate([close()]);
  validate([{ ...close(), state_reason: "not_planned" }]);
  validate([comment(), assign()]);
  validate([{ ...comment(), item_number: "42", repo: "azure/AZURE-SDK-FOR-GO" }]);
  for (const items of [[], [null], [{ type: "constructor" }], [{ type: "add_labels" }],
    [comment(), comment()], [assign()], [assign(), comment()], [close(), assign()],
    [comment(), close()], [comment(), { type: "noop" }]]) {
    assert.throws(() => validate(items));
  }
  for (const output of [null, {}, { items: {} }]) {
    assert.throws(() => validateOutputs(output, 42, repository));
  }
  for (const overrides of [{ state_reason: "completed" }, { state_reason: "duplicate" }, { duplicate_of: 43 }]) {
    assert.throws(() => validate([{ ...close(), ...overrides }]), /as not planned/);
  }
});

test("rejects cross-issue, alternate-target, cross-repository and empty-comment outputs", () => {
  for (const make of [comment, close, assign]) {
    const item = make();
    const field = item.type === "add_comment" ? "item_number" : "issue_number";
    const validateItem = value => validate(value.type === "assign_to_agent" ? [comment(), value] : [value]);
    for (const target of [undefined, 0, 41, "42junk"]) {
      assert.throws(() => validateItem({ ...item, [field]: target }));
    }
    for (const alias of ["item_number", "issue_number", "pull_number", "pull_request_number",
      "pr_number", "pr", "pr-number", "comment_id", "commentId", "comment-id", "reply_to_id", "target"]) {
      if (alias !== field) assert.throws(() => validateItem({ ...item, [alias]: 42 }));
    }
    for (const key of ["repo", "pull_request_repo"]) {
      for (const value of ["Azure/other", "", 42]) {
        assert.throws(() => validateItem({ ...item, [key]: value }));
      }
    }
  }
  for (const make of [comment, close]) {
    for (const body of [undefined, null, "", " \n", 42]) {
      assert.throws(() => validate([{ ...make(), body }]), /explanatory comment/);
    }
  }
});

test("reads only the dispatched issue and propagates API errors", async () => {
  const calls = [];
  const args = {
    context: { repo: { owner: "Azure", repo: "azure-sdk-for-go" } },
    core: { info() {} },
    number: "42",
    github: { rest: { issues: { async get(parameters) {
      calls.push(parameters);
      return { data: eligibleIssue() };
    } } } },
  };
  assert.equal(await checkEligibility(args), true);
  assert.deepEqual(calls, [{ owner: "Azure", repo: "azure-sdk-for-go", issue_number: 42 }]);
  args.github.rest.issues.get = async () => { throw new Error("API unavailable"); };
  await assert.rejects(checkEligibility(args), /API unavailable/);
  await assert.rejects(checkEligibility({ ...args, number: "0" }), /positive safe integer/);
});

test("rechecks eligibility before applying outputs and rejects malformed artifacts", async () => {
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), "issue-investigation-"));
  try {
    const outputFile = path.join(directory, "output.json");
    let issue = eligibleIssue();
    let calls = 0;
    const args = {
      outputFile, number: "42",
      context: { repo: { owner: "Azure", repo: "azure-sdk-for-go" } },
      core: { info() {} },
      github: { rest: { issues: { async get() { calls++; return { data: issue }; } } } },
    };
    fs.writeFileSync(outputFile, JSON.stringify({ items: [comment(), assign()] }));
    await checkOutputs(args);
    assert.equal(calls, 1);
    issue = { ...issue, state: "closed" };
    await assert.rejects(checkOutputs(args), /eligibility changed/);
    fs.writeFileSync(outputFile, JSON.stringify({ items: [{ ...comment(), item_number: 43 }] }));
    await assert.rejects(checkOutputs(args), /dispatched issue/);
    assert.equal(calls, 2);
    fs.writeFileSync(outputFile, JSON.stringify({ items: [{ type: "noop", message: "No action needed" }] }));
    await checkOutputs(args);
    assert.equal(calls, 2);
    fs.writeFileSync(outputFile, "{");
    await assert.rejects(checkOutputs(args), SyntaxError);
    await assert.rejects(checkOutputs({ ...args, outputFile: path.join(directory, "missing") }), /ENOENT/);
  } finally {
    fs.rmSync(directory, { recursive: true, force: true });
  }
});

test("rejects incomplete collection and pending-suggestion investigation plans", () => {
  for (const errors of ["", {}, ["Invalid preceding output"]]) {
    assert.throws(() => validateOutputs({ items: [comment()], errors }, 42, repository), /partially collected/);
  }
  validateOutputs({ items: [comment()], errors: [] }, 42, repository);
  for (const item of [close(), assign()]) {
    const items = item.type === "assign_to_agent" ? [comment(), { ...item, suggest: true }] : [{ ...item, suggest: true }];
    assert.throws(() => validate(items), /pending suggestions/);
  }
  for (const message of [undefined, "", " \n", 42]) {
    assert.throws(() => validate([{ type: "noop", message }]), /requires an explanation/);
  }
});

test("assignment is released only after successful native comment delivery", async () => {
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), "issue-handoff-"));
  try {
    const outputFile = path.join(directory, "output.json");
    const original = { items: [comment(), assign()], errors: [] };
    let delivered = 0;
    const outputs = [];
    let issue = eligibleIssue();
    const args = {
      outputFile, number: "42",
      context: { repo: { owner: "Azure", repo: "azure-sdk-for-go" } },
      core: { info() {}, setOutput(name, value) { outputs.push([name, value]); } },
      github: { rest: { issues: { async get() { return { data: issue }; } } } },
      async postComment(item) {
        assert.deepEqual(item, comment());
        assert.deepEqual(JSON.parse(fs.readFileSync(outputFile, "utf8")), original);
        delivered++;
        return { success: true, commentId: 100 };
      },
    };
    const reset = () => fs.writeFileSync(outputFile, JSON.stringify(original));
    const assertUnchanged = () => assert.deepEqual(JSON.parse(fs.readFileSync(outputFile, "utf8")), original);
    reset();
    await prepareOutputs(args);
    assert.equal(delivered, 1);
    assert.deepEqual(JSON.parse(fs.readFileSync(outputFile, "utf8")), { items: [assign()], errors: [] });
    assert.deepEqual(outputs, [["delivered_comment_id", "100"]]);

    for (const result of [undefined, { success: false }, { success: true, skipped: true, commentId: 100 },
      { success: true, staged: true }, { success: true }]) {
      reset();
      await assert.rejects(prepareOutputs({ ...args, postComment: async () => result }), /not delivered/);
      assertUnchanged();
    }
    reset();
    await assert.rejects(prepareOutputs({ ...args, postComment: async () => { throw new Error("HTTP 403"); } }), /HTTP 403/);
    assertUnchanged();
    await assert.rejects(prepareOutputs({ ...args, postComment: undefined }), /Native comment handler is required/);
    assertUnchanged();

    await prepareOutputs({ ...args, staged: true, postComment: async () => { throw new Error("Must not post in staged mode"); } });
    assertUnchanged();
    assert.equal(delivered, 1);

    await assert.rejects(prepareOutputs({ ...args, postComment: async () => {
      issue = { ...issue, state: "closed" };
      return { success: true, commentId: 100 };
    } }), /eligibility changed after comment/);
    assertUnchanged();
    issue = eligibleIssue();
    fs.writeFileSync(outputFile, JSON.stringify({ items: [comment()] }));
    await prepareOutputs({ ...args, postComment: async () => { throw new Error("Standalone comments stay native"); } });
    assert.deepEqual(JSON.parse(fs.readFileSync(outputFile, "utf8")), { items: [comment()] });
    fs.writeFileSync(outputFile, JSON.stringify({ items: [close()] }));
    await prepareOutputs({ ...args, postComment: async () => { throw new Error("Closure stays native"); } });
    assert.deepEqual(JSON.parse(fs.readFileSync(outputFile, "utf8")), { items: [close()] });
  } finally {
    fs.rmSync(directory, { recursive: true, force: true });
  }
});

test("triage dispatch requires a complete same-issue explanation and routing plan", () => {
  const dispatch = () => ({ type: "dispatch_workflow", workflow_name: "issue-investigation", inputs: { issue_number: "42" } });
  const mention = () => ({ type: "mention_owners", owners: "owner1, owner2", message: "Routing" });
  const owner = () => ({ type: "assign_to_user", issue_number: 42, assignees: ["owner1"] });
  const check = (items, errors = []) => validateTriageOutputs({ items, errors }, 42, repository);
  check([mention(), comment(), dispatch()]);
  check([owner(), comment(), comment(), dispatch()]);
  check([owner(), mention(), comment(), dispatch()]);
  check([{ type: "noop", message: "Manual triage" }]);
  check([comment()]);
  for (const items of [
    [dispatch()], [comment(), dispatch()], [owner(), comment(), dispatch()],
    [mention(), dispatch()], [mention(), comment(), comment(), dispatch()],
    [owner(), owner(), comment(), comment(), dispatch()],
    [mention(), comment(), { ...dispatch(), workflow_name: "other" }],
    [mention(), comment(), { ...dispatch(), inputs: { issue_number: "43" } }],
    [dispatch(), mention(), comment()], [mention(), comment(), dispatch(), dispatch()],
    [mention(), { ...comment(), item_number: 43 }, dispatch()],
    [mention(), { ...comment(), repo: "Azure/other" }, dispatch()],
    [mention(), { ...comment(), body: " " }, dispatch()],
    [mention(), comment(), close(), dispatch()],
    [mention(), comment(), { type: "noop" }, dispatch()],
  ]) {
    assert.throws(() => check(items));
  }
  for (const errors of ["", {}, ["Missing owner output"]]) {
    assert.throws(() => check([mention(), comment(), dispatch()], errors), /partially collected/);
  }
  assert.throws(() => validateTriageOutputs(null, 42, repository), /triage safe-output/);
});

test("native comment prerequisite reuses the entire compiled policy", () => {
  const workflowFile = path.join(__dirname, "..", "workflows", "issue-investigation.lock.yml");
  const workflow = fs.readFileSync(workflowFile, "utf8");
  const domains = [...workflow.matchAll(/^ +GH_AW_ALLOWED_DOMAINS: ("[^\r\n]*")\r?$/gm)]
    .map(match => JSON.parse(match[1]));
  const policy = nativeCommentPolicy(workflowFile, 42);
  assert.ok(domains.length >= 2);
  assert.ok(domains.every(value => value === policy.domains));
  assert.equal(policy.config.target, "42");
  assert.equal(policy.config.max, 1);
  for (const host of ["github.com", "api.github.com", "learn.microsoft.com", "feedback.azure.com", "pkg.go.dev", "proxy.golang.org"]) {
    assert.ok(policy.domains.split(",").includes(host));
  }
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), "native-comment-policy-"));
  try {
    const file = path.join(directory, "workflow.yml");
    for (const changed of ["", workflow.replaceAll("GH_AW_ALLOWED_DOMAINS:", "REMOVED_DOMAINS:"),
      workflow.replaceAll('\\"add_comment\\":{\\"max\\":1', '\\"add_comment\\":{\\"max\\":2'),
      workflow + '\n          GH_AW_ALLOWED_DOMAINS: "different.example"\n']) {
      fs.writeFileSync(file, changed);
      assert.throws(() => nativeCommentPolicy(file, 42));
    }
  } finally {
    fs.rmSync(directory, { recursive: true, force: true });
  }
});
