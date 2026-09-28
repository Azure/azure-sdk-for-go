// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

const assert = require("node:assert/strict");
const { test } = require("node:test");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const { issueNumber, isEligible, validateOutputs, checkEligibility, checkOutputs } =
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
    fs.writeFileSync(outputFile, JSON.stringify({ items: [{ type: "noop" }] }));
    await checkOutputs(args);
    assert.equal(calls, 2);
    fs.writeFileSync(outputFile, "{");
    await assert.rejects(checkOutputs(args), SyntaxError);
    await assert.rejects(checkOutputs({ ...args, outputFile: path.join(directory, "missing") }), /ENOENT/);
  } finally {
    fs.rmSync(directory, { recursive: true, force: true });
  }
});
