// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

const assert = require("node:assert/strict");
const { createHash } = require("node:crypto");
const { test } = require("node:test");
const { gatherModules } = require("./pr-module-review.cjs");
const { markReadyForReview } = require("./draft-review.cjs");
const { checkPipelines } = require("./pipeline-review.cjs");
const { checkBreakingChanges, moduleVersion, breakingItems } = require("./breaking-change-review.cjs");
const { checkOnboarding, namespaceEvidence, releasePipeline, serviceToken } = require("./onboarding-review.cjs");
const { devOpsReader, fetchSpecText, authorization } = require("./azure-devops-read.cjs");
const { runReviewSteps } = require("./review-steps.cjs");
const { assertReviewResults, assertAgentOutput, hasBlockers, formatReviewComment, prepareAgentOutput, publishReviewFailure } = require("./review-output.cjs");

const headSHA = "a".repeat(40);
const baseSHA = "b".repeat(40);
const specSHA = "c".repeat(40);
const path = "sdk/resourcemanager/compute/armcompute";
const dashboard = "https://azsdk-releaseplan-dashboard-example.westus-01.azurewebsites.net/?releaseplan=37";
const change = "Operation `*Client.Delete` has been removed";
const changelog = `# Release History\n## 2.0.0 (2026-11-13)\n### Breaking Changes\n\n- ${change}\n\n### Features Added\n- Not a breaking change\n## 1.1.0 (2026-01-01)\n### Breaking Changes\n- Older change\n`;
const suppression = `suppressions:\n  azure-sdk-for-go:\n    - package: ${path}\n      breaking-changes:\n        - '${change}'\n`;

function fixture(first = false) {
  const state = {
    pr: {
      number: 42, node_id: "PR42", state: "open", draft: false, head: { sha: headSHA },
      base: { sha: baseSHA, ref: "main" }, changed_files: 1, comments: 0, body: dashboard,
      mergeable: true, mergeable_state: "clean",
    },
    provenance: { status: "pass", pr_number: 42, head_sha: headSHA,
      base_sha: baseSHA, base_ref: "main", findings: [], errors: [] },
    files: [{ filename: `${path}/${first ? "ci.yml" : "version.go"}`, status: first ? "added" : "modified" }],
    comments: [],
    checks: [{ id: 1, name: "go - pullrequest", head_sha: headSHA, status: "completed", conclusion: "success",
      details_url: "https://dev.azure.com/azure-sdk/public/_build/results?buildId=1" }],
    texts: new Map([
      [`${path}/version.go`, 'package armcompute\nconst moduleVersion = "v2.0.0"\n'],
      [`${path}/CHANGELOG.md`, changelog],
      [`${path}/tsp-location.yaml`, `repo: Azure/azure-rest-api-specs\ncommit: ${specSHA}\ndirectory: specification/compute/resource-manager\n`],
      [`${path}/README.md`, "# Azure Compute management client\n\nManage Azure Compute.\n\n## Getting started\nTemplate text"],
    ]),
    specification: suppression,
    writes: [], reads: [], logs: [], adoCalls: [], fetched: [], pages: [],
    definitions: [{ id: 10, repository: { id: "Azure/azure-sdk-for-go", type: "GitHub" },
      process: { yamlFilename: `${path}/ci.yml` } }],
  };
  const blobs = new Map();
  const core = {
    info: (message) => state.logs.push(message),
    error: (message) => state.logs.push(message),
    warning: (message) => state.logs.push(message),
    setSecret: (value) => state.logs.push({ masked: value }),
  };
  const github = {
    rest: {
      pulls: {
        get: async () => ({ data: structuredClone(state.pr) }),
        listFiles: async () => ({ data: structuredClone(state.files) }),
      },
      issues: {
        listComments: async () => ({ data: structuredClone(state.comments) }),
        removeLabel: async (value) => state.writes.push(value),
        createComment: async (value) => state.writes.push({ kind: "comment", ...value }),
        addLabels: async (value) => state.writes.push(value),
      },
      checks: {
        listForRef: async ({ page, per_page, ref, filter }) => {
          assert.equal(ref, headSHA);
          assert.equal(filter, "latest");
          state.pages.push(page);
          return { data: { total_count: state.checks.length, check_runs: state.checks.slice((page - 1) * per_page, page * per_page) } };
        },
      },
      repos: {
        get: async () => ({ data: { allow_squash_merge: true } }),
        getContent: async ({ path: filename, ref }) => {
          assert.equal(ref, headSHA);
          state.reads.push(filename);
          if (!state.texts.has(filename)) throw Object.assign(new Error(`${filename} unavailable.`), { status: 404 });
          const content = Buffer.from(state.texts.get(filename));
          const sha = createHash("sha1").update(`blob ${content.length}\0`).update(content).digest("hex");
          blobs.set(sha, { sha, size: content.length, encoding: "base64", content: content.toString("base64") });
          return { data: { type: "file", sha } };
        },
      },
      git: { getBlob: async ({ file_sha }) => ({ data: blobs.get(file_sha) }) },
    },
    graphql: async (query, variables) => {
      assert.match(query, /markPullRequestReadyForReview/);
      assert.deepEqual(variables, { id: "PR42" });
      state.writes.push({ mutation: "ready" });
      return { markPullRequestReadyForReview: { pullRequest: {
        id: "PR42", isDraft: false, headRefOid: headSHA, baseRefOid: baseSHA, baseRefName: "main",
      } } };
    },
    paginate: async (method, options) => {
      assert.equal(options.per_page, 100);
      return (await method(options)).data;
    },
  };
  state.api = { github, core, context: { repo: { owner: "Azure", repo: "azure-sdk-for-go" } } };
  state.adoRead = async (project, route, parameters, body) => {
    state.adoCalls.push({ project, route, parameters, body });
    if (route === "wit/wiql") return { data: { workItems: [{ id: 99 }] } };
    if (route === "wit/workitems/99") return { data: { fields: {
      "Custom.ReleasePlanID": 37, "Custom.NamespaceApprovalIssue": "https://github.com/Azure/azure-sdk/issues/123",
    } } };
    if (route === "build/definitions") return { data: { count: state.definitions.length, value: state.definitions } };
    throw new Error(`Unexpected Azure DevOps route ${route}.`);
  };
  state.fetchText = async (url) => {
    state.fetched.push(url);
    return { status: 200, text: state.specification };
  };
  state.run = () => runReviewSteps(state.api, {
    prNumber: 42, headSHA, provenance: state.provenance, adoRead: state.adoRead, fetchText: state.fetchText,
  });
  return state;
}

test("module gathering detects added ci.yml and both sides of cross-module renames", async () => {
  const state = fixture(true);
  state.files.push({ filename: "sdk/resourcemanager/network/armnetwork/file.go", status: "renamed",
    previous_filename: `${path}/file.go` });
  state.pr.changed_files++;
  const result = await gatherModules(state.api, 42, headSHA);
  assert.deepEqual(result.modules, [
    { path, first_onboard: true }, { path: "sdk/resourcemanager/network/armnetwork", first_onboard: false },
  ]);
  state.files[0].status = "modified";
  assert.equal((await gatherModules(state.api, 42, headSHA)).modules[0].first_onboard, false);
});

test("incomplete or duplicate changed-file lists do not imply no new modules", async () => {
  const state = fixture();
  state.pr.changed_files++;
  await assert.rejects(gatherModules(state.api, 42, headSHA), /incomplete/);
  state.files.push(state.files[0]);
  await assert.rejects(gatherModules(state.api, 42, headSHA), /inconsistent/);
});

test("draft conversion uses the trusted mutation once and checks its returned snapshot", async () => {
  const state = fixture();
  assert.deepEqual(await markReadyForReview(state.api, state.pr, headSHA), { status: "ready", converted: false });
  assert.equal(state.writes.length, 0);
  state.pr.draft = true;
  assert.equal((await markReadyForReview(state.api, state.pr, headSHA)).converted, true);
  assert.equal(state.writes.length, 1);
  state.api.github.graphql = async () => ({ markPullRequestReadyForReview: { pullRequest: {
    id: "PR42", isDraft: false, headRefOid: baseSHA, baseRefOid: baseSHA, baseRefName: "main",
  } } });
  await assert.rejects(markReadyForReview(state.api, state.pr, headSHA), /changed/);
});

test("pending or absent parent checks defer other checks without waiting", async () => {
  for (const checks of [[], [{ id: 1, name: "go - pullrequest", head_sha: headSHA, status: "in_progress" }]]) {
    const state = fixture(true);
    state.checks = checks;
    const result = await state.run();
    assert.equal(result.status, "pending");
    assert.deepEqual(state.reads, []);
    assert.deepEqual(state.adoCalls, []);
    assert.deepEqual(result.breaking_changes, []);
    assert.equal(hasBlockers(result), true);
  }
});

test("the review's own queued checks do not keep the ADO pipeline pending", async () => {
  const state = fixture();
  state.checks.push({ id: 2, name: "Management Release Assistant", head_sha: headSHA, status: "queued" });
  assert.equal((await checkPipelines(state.api, headSHA)).status, "completed");
});

test("pipeline categories and log links come from current-head API data", async () => {
  const state = fixture();
  const names = ["Build/Test on windows_go_1273", "Analyze", "generate_job_matrix"];
  state.checks.push(...names.map((name, index) => ({
    id: index + 2, head_sha: headSHA, name: `go - pullrequest (Build ${name})`,
    status: "completed", conclusion: "failure", details_url: `https://dev.azure.com/azure-sdk/public/_build/results?buildId=${index + 2}`,
    output: { title: "failed", summary: "Actual API summary." },
  })));
  const result = await checkPipelines(state.api, headSHA);
  assert.deepEqual(result.failures.map((failure) => failure.category), ["build-test", "analyze", "matrix"]);
  assert.equal(result.failures[0].target_url, state.checks[1].details_url);
  assert.equal(result.failures[0].summary, "Actual API summary.");
});

test("check-run pagination verifies counts, duplicates, and head identity", async () => {
  const state = fixture();
  state.checks.push(...Array.from({ length: 100 }, (_, index) => ({
    id: index + 2, head_sha: headSHA, name: `Other ${index}`, status: "completed",
  })));
  assert.equal((await checkPipelines(state.api, headSHA)).status, "completed");
  assert.deepEqual(state.pages, [1, 2]);
  state.checks[100].id = 1;
  await assert.rejects(checkPipelines(state.api, headSHA), /duplicate/);
  state.checks[100].id = 101;
  state.checks[100].head_sha = baseSHA;
  await assert.rejects(checkPipelines(state.api, headSHA), /stale/);
  state.api.github.rest.checks.listForRef = async () => ({ data: { total_count: 3, check_runs: state.checks.slice(0, 1) } });
  await assert.rejects(checkPipelines(state.api, headSHA), /incomplete/);
});

for (const version of ["v0.9.0", "v2.0.0-beta.1", "v1.0.0+build.1"]) {
  test(`${version} skips breaking-change lookup`, async () => {
    const state = fixture();
    state.texts.set(`${path}/version.go`, `package p\nconst (\nmoduleVersion = "${version}"\n)\n`);
    const result = await state.run();
    assert.equal(result.breaking_changes[0].status, "skipped");
    assert.deepEqual(state.reads, [`${path}/version.go`]);
    assert.deepEqual(state.fetched, []);
  });
}

test("module version extraction ignores commented and raw-string declarations", () => {
  const source = '/* const moduleVersion = "v0.1.0" */\nconst text = `const moduleVersion = "v0.2.0"`\nconst (\nmoduleVersion = "v2.0.0+build.1"\n)\n';
  assert.deepEqual(moduleVersion(source), { version: "2.0.0+build.1", skip: false });
  assert.throws(() => moduleVersion('const moduleVersion = "v2.0.0"\nconst moduleVersion = "v0.1.0"'), /unique/);
});

test("only bullets from the new release's Breaking Changes section are checked", () => {
  assert.deepEqual(breakingItems(changelog, "2.0.0"), [change]);
  assert.throws(() => breakingItems(changelog, "3.0.0"), /release entry/);
  assert.deepEqual(breakingItems("## 2.0.0 (date)\n### Features Added\n- A\n", "2.0.0"), []);
  assert.deepEqual(breakingItems("## 2.0.0 (date)\n### Breaking Changes\n\n### Other Changes\n- A\n", "2.0.0"), []);
});

test("suppression lookup stays at the PR's exact pinned specification commit", async () => {
  const state = fixture();
  const result = await state.run();
  assert.equal(result.breaking_changes[0].status, "approved");
  assert.deepEqual(result.breaking_changes[0].items, [{ text: change, status: "Approved", reason: "" }]);
  assert.deepEqual(state.fetched, [`https://raw.githubusercontent.com/Azure/azure-rest-api-specs/${specSHA}/specification/compute/resource-manager/sdk-suppressions.yaml`]);
});

for (const [label, yaml, status] of [
  ["wrong language", suppression.replace("azure-sdk-for-go", "azure-sdk-for-python"), "unapproved"],
  ["wrong package", suppression.replace(path, "sdk/resourcemanager/network/armnetwork"), "unapproved"],
  ["import suffix", suppression.replace(path, `${path}/v2`), "unapproved"],
  ["nonexact text", suppression.replace("has been removed", "has been removed."), "unapproved"],
  ["empty document", "", "unapproved"],
  ["missing list", `suppressions:\n  azure-sdk-for-go:\n    - package: ${path}\n`, "unapproved"],
  ["invalid YAML", "suppressions: [", "unverifiable"],
  ["duplicate keys", "suppressions: {}\nsuppressions: {}\n", "unverifiable"],
  ["invalid Go structure", "suppressions:\n  azure-sdk-for-go: {}\n", "unverifiable"],
  ["invalid approval type", `suppressions:\n  azure-sdk-for-go:\n    - package: ${path}\n      breaking-changes: [42]\n`, "unverifiable"],
]) {
  test(`${label} never approves a breaking change`, async () => {
    const state = fixture();
    state.specification = yaml;
    const result = await state.run();
    assert.equal(result.breaking_changes[0].status, status);
    assert.equal(hasBlockers(result), true);
    assertReviewResults(result, state.provenance);
  });
}

for (const status of [404, 401, 403, 500]) {
  test(`suppression HTTP ${status} has the correct missing-versus-unverifiable result`, async () => {
    const state = fixture();
    const review = await checkBreakingChanges(path, async (filename) => state.texts.get(filename),
      async () => ({ status, text: suppression }));
    assert.equal(review.status, status === 404 ? "unapproved" : "unverifiable");
  });
}

test("invalid versions and missing release/configuration are reported, never skipped", async () => {
  for (const filename of ["version.go", "CHANGELOG.md", "tsp-location.yaml"]) {
    const state = fixture();
    state.texts.delete(`${path}/${filename}`);
    const result = await state.run();
    assert.equal(result.breaking_changes[0].status, "unverifiable");
    assert.equal(result.breaking_changes[0].errors.length, 1);
    assert.ok(state.logs.some((message) => typeof message === "string" && message.includes(filename)));
  }
  assert.throws(() => moduleVersion('const moduleVersion = "garbage"'), /Invalid Version/);
});

test("namespace resolution uses ReleasePlanID, not the dashboard value as work-item ID", async () => {
  const state = fixture(true);
  const result = await state.run();
  assert.equal(result.onboarding[0].namespace.source, "release-plan");
  assert.match(state.adoCalls[0].body.query, /\[Custom\.ReleasePlanID\] = '37'/);
  assert.equal(state.adoCalls[1].route, "wit/workitems/99");
  assert.ok(!state.adoCalls.some((call) => call.route === "wit/workitems/37"));
});

for (const kind of ["issues", "pull"]) {
  test(`namespace ${kind} links retain the manual fallback after ADO failure`, async () => {
    const state = fixture();
    state.comments = [
      { id: 1, body: `<!-- gh-aw-workflow-id: mgmt-review --> namespace https://github.com/Azure/test/${kind}/1` },
      { id: 2, body: `Namespace review: https://github.com/Azure/test/${kind}/123` },
    ];
    const result = await namespaceEvidence({ ...state.api, prNumber: 42, commentCount: 2 }, dashboard,
      async () => { throw new Error("ADO read denied."); });
    assert.equal(result.status, "satisfied");
    assert.equal(result.issue, `https://github.com/Azure/test/${kind}/123`);
    assert.ok(state.logs.some((message) => message.includes("ADO read denied")));
  });
}

test("workflow comments alone or incomplete comment lists do not establish namespace approval", async () => {
  const state = fixture();
  state.comments = [{ id: 1, body: "<!-- gh-aw-workflow-id: mgmt-review --> namespace https://github.com/Azure/test/issues/1" }];
  const result = await namespaceEvidence({ ...state.api, prNumber: 42, commentCount: 1 }, "", state.adoRead);
  assert.equal(result.status, "missing");
  await assert.rejects(namespaceEvidence({ ...state.api, prNumber: 42, commentCount: 2 }, "", state.adoRead), /incomplete/);
});

test("release-pipeline creation requires a definition for the exact repository and module YAML", async () => {
  const state = fixture();
  assert.equal((await releasePipeline(path, state.api.context.repo, state.adoRead)).status, "satisfied");
  const request = state.adoCalls[0];
  assert.equal(request.project, "internal");
  assert.equal(request.parameters.yamlFilename, `${path}/ci.yml`);
  state.definitions = [];
  assert.equal((await releasePipeline(path, state.api.context.repo, state.adoRead)).status, "missing");
  state.definitions = [{ id: 1, repository: { id: "Another/repo", type: "GitHub" }, process: { yamlFilename: `${path}/ci.yml` } }];
  await assert.rejects(releasePipeline(path, state.api.context.repo, state.adoRead), /do not match/);
});

test("inaccessible pipelines and README files remain verification errors, not naming issues", async () => {
  const state = fixture(true);
  state.texts.delete(`${path}/README.md`);
  const results = await checkOnboarding({ ...state.api, prNumber: 42, commentCount: 0 },
    [{ path, first_onboard: true }], "", async () => { throw new Error("README access denied."); },
    async () => { throw new Error("Pipeline access denied."); });
  assert.equal(results[0].pipeline.status, "unverifiable");
  assert.deepEqual(results[0].readme, { status: "unverifiable", errors: ["README access denied."] });
});

test("service-token normalization retains the existing single/concatenated rules", () => {
  assert.equal(serviceToken("sdk/resourcemanager/network/armnetwork"), "network");
  assert.equal(serviceToken("sdk/resourcemanager/compute/armbulkactions"), "computebulkactions");
});

test("new-module naming stays with the agent and excludes installation boilerplate", async () => {
  const state = fixture(true);
  const result = await state.run();
  assert.equal(result.onboarding[0].readme.status, "semantic-review-required");
  assert.match(result.onboarding[0].readme.text, /Azure Compute/);
  assert.doesNotMatch(result.onboarding[0].readme.text, /Getting started|Template text/);
  assertReviewResults(result, state.provenance);
  assert.throws(() => assertAgentOutput(result, { items: [{ body: "## PR is ready to merge" }] }), /readiness/);
  assert.doesNotThrow(() => assertAgentOutput(result, { items: [{ body:
    `## PR is ready to merge\nService-name review confirmed for ${path}/README.md.` }] }));
});

test("provenance failures keep the approved fail-fast behavior", async () => {
  const state = fixture();
  state.provenance.status = "manual-edits";
  state.provenance.findings = [{ file: "generated.go", sha: headSHA, author: "developer" }];
  const result = await state.run();
  assert.equal(result.status, "unverifiable");
  assert.deepEqual(state.reads, []);
  assert.deepEqual(state.pages, []);
});

test("head changes during mechanical review reject the result and clear the old progress label", async () => {
  const state = fixture();
  const list = state.api.github.rest.checks.listForRef;
  state.api.github.rest.checks.listForRef = async (request) => {
    state.pr.head.sha = baseSHA;
    return list(request);
  };
  const result = await state.run();
  assert.equal(result.status, "unverifiable");
  assert.match(result.errors[0], /changed/);
  assert.equal(state.writes[0].name, "mgmt-review-in-progress");
});

test("publisher validation rejects missing module results and canonical or emoji ready claims with blockers", async () => {
  const state = fixture();
  const result = await state.run();
  assertReviewResults(result, state.provenance);
  assert.doesNotThrow(() => assertAgentOutput(result, { items: [{ body: "## PR is ready to merge" }] }));
  const missing = { ...result, breaking_changes: [] };
  assert.throws(() => assertReviewResults(missing, state.provenance), /incomplete/);
  const blocked = structuredClone(result);
  blocked.breaking_changes[0].status = "unapproved";
  blocked.breaking_changes[0].items[0].status = "Unapproved";
  blocked.breaking_changes[0].items[0].reason = "No exact approval.";
  assertReviewResults(blocked, state.provenance);
  for (const heading of ["## PR is ready to merge", "## ✅ PR is ready to merge"]) {
    assert.throws(() => assertAgentOutput(blocked, { items: [{ arguments: { body: heading } }] }), /readiness/);
  }
  assert.doesNotThrow(() => assertAgentOutput(blocked, { items: [{ body: "## PR is not ready to merge" }] }));
  assert.throws(() => assertReviewResults({ ...result, base_ref: "release" }, state.provenance), /stale/);
});

test("read clients restrict credential forwarding and preserve HTTP errors", async () => {
  const state = fixture();
  const header = authorization(state.api.core, "read-only-pat");
  const read = devOpsReader(header, async (url, options) => {
    assert.equal(url.origin, "https://dev.azure.com");
    assert.equal(options.headers.Authorization, header);
    assert.equal(options.redirect, "error");
    return { ok: false, status: 403 };
  });
  await assert.rejects(read("Release", "wit/wiql"), /HTTP 403/);
  await assert.rejects(devOpsReader(null)("internal", "build/definitions"), /authentication/);
  await assert.rejects(fetchSpecText("https://example.com/spec"), /host/);
  const notFound = await fetchSpecText("https://raw.githubusercontent.com/Azure/spec/commit/file", async (_, options) => {
    assert.equal(options.redirect, "error");
    assert.equal(options.headers, undefined);
    return { status: 404, ok: false };
  });
  assert.deepEqual(notFound, { status: 404, text: "" });
});

test("an unreadable squash setting does not suppress independent module checks", async () => {
  const state = fixture();
  state.checks[0].conclusion = "failure";
  state.api.github.rest.repos.get = async () => { throw new Error("Repository settings denied."); };
  const result = await state.run();
  assert.equal(result.status, "completed");
  assert.equal(result.breaking_changes[0].status, "approved");
  assert.deepEqual(result.mergeability_errors, ["Repository settings denied."]);
  assert.equal(hasBlockers(result), true);
  assertReviewResults(result, state.provenance);
});

test("unknown or blocked mergeability cannot produce merge-ready guidance", async () => {
  for (const value of [null, true]) {
    const state = fixture();
    state.pr.mergeable = value;
    state.pr.mergeable_state = value ? "blocked" : "unknown";
    const result = await state.run();
    assertReviewResults(result, state.provenance);
    assert.equal(hasBlockers(result), true);
    assert.throws(() => prepareAgentOutput(result, { items: [{ body: "## PR is ready to merge" }] }), /readiness/);
  }
});

test("unresolved YAML tags and thrown non-Error failures remain unverifiable", async () => {
  const state = fixture();
  state.specification = suppression.replace(`'${change}'`, `!unsupported '${change}'`);
  assert.equal((await state.run()).breaking_changes[0].status, "unverifiable");
  const result = await checkBreakingChanges(path, async () => { throw "Read denied."; }, state.fetchText);
  assert.deepEqual(result.errors, ["Read denied."]);
});

test("trusted rendering retains every approval item and omits skipped releases", async () => {
  const state = fixture();
  state.specification = "";
  const result = await state.run();
  assertReviewResults(result, state.provenance);
  const output = prepareAgentOutput(result, { items: [{ type: "add_comment", body: "<!-- gh-aw-workflow-id: mgmt-review -->\nCI diagnosis." }] });
  const body = output.items[0].body;
  assert.match(body, /\\\*Client\.Delete/);
  assert.match(body, /Unapproved/);
  assert.match(body, /CI diagnosis/);
  assert.equal(body.split("<!-- gh-aw-workflow-id: mgmt-review -->").length - 1, 1);
  state.texts.set(`${path}/version.go`, 'const moduleVersion = "v2.0.0-beta.1"');
  const skipped = formatReviewComment(await state.run(), "## PR is ready to merge");
  assert.doesNotMatch(skipped, /Breaking-Change Review/);
});

test("absent suppressions link to the pinned directory and never infer approval", async () => {
  const state = fixture();
  state.fetchText = async () => ({ status: 404, text: "" });
  const body = formatReviewComment(await state.run());
  assert.match(body, new RegExp(`/tree/${specSHA}/specification/compute/resource-manager`));
  assert.doesNotMatch(body, /\/blob\//);
  assert.match(body, /Unapproved/);
});

test("the publisher combines only outstanding mechanical and semantic onboarding items", async () => {
  const state = fixture(true);
  state.definitions = [];
  const result = await state.run();
  const heading = "## ⚠️ First On-Board Service — Action Required";
  const body = formatReviewComment(result, `CI diagnosis.\n\n${heading}\n\n- **README service name**: confirm ${path}/README.md.`);
  assert.equal(body.split(heading).length - 1, 1);
  assert.match(body, /Pipeline setup/);
  assert.match(body, /README service name/);
  assert.doesNotMatch(body, /\*\*Namespace approval\*\*/);
  assert.ok(body.indexOf("CI diagnosis") < body.indexOf("Pipeline setup"));
});

test("pending CI is published without the agent and updates only completion labels", async () => {
  const state = fixture();
  state.checks[0].status = "in_progress";
  const review = await state.run();
  await publishReviewFailure(state.api, { review, provenance: state.provenance, prNumber: 42, headSHA });
  assert.equal(state.writes.length, 3);
  assert.match(state.writes[0].body, /Pipeline Checks Still Running/);
  assert.doesNotMatch(state.writes[0].body, /Breaking-Change Review|First On-Board|ready to merge/);
  assert.equal(state.writes[1].name, "mgmt-review-in-progress");
  assert.deepEqual(state.writes[2].labels, ["mgmt-review-analyzed"]);
});

test("a setup failure without results produces explicit manual-review guidance", async () => {
  const state = fixture();
  await publishReviewFailure(state.api, { review: null, provenance: state.provenance, prNumber: 42, headSHA, jobResult: "failure" });
  assert.match(state.writes[0].body, /review\\_checks job ended failure/);
  assert.match(state.writes[0].body, /not ready to merge/);
  assert.ok(state.logs.some((message) => message.includes("without producing review results")));
});

test("pending reporters reject stale snapshots and non-passing provenance", async () => {
  const state = fixture();
  state.checks[0].status = "queued";
  const review = await state.run();
  state.pr.base.sha = specSHA;
  await assert.rejects(publishReviewFailure(state.api, {
    review, provenance: state.provenance, prNumber: 42, headSHA,
  }), /No review outputs/);
  assert.equal(state.writes.filter((item) => item.kind === "comment").length, 0);
  state.provenance.status = "manual-edits";
  await assert.rejects(publishReviewFailure(state.api, {
    review, provenance: state.provenance, prNumber: 42, headSHA,
  }), /passing provenance/);
});

test("malformed approval rows and missing onboarding evidence cannot be published", async () => {
  const state = fixture(true);
  const review = await state.run();
  const invalid = structuredClone(review);
  invalid.breaking_changes[0].items[0].status = "Unknown";
  assert.throws(() => assertReviewResults(invalid, state.provenance), /incomplete or invalid/);
  const missing = structuredClone(review);
  missing.onboarding[0].pipeline.definition_ids = [];
  assert.throws(() => assertReviewResults(missing, state.provenance), /incomplete or invalid/);
});

test("completion labels are scripted after the comment even if the agent omits or misorders them", async () => {
  const review = await fixture().run();
  for (const items of [
    [{ type: "add_comment", body: "## PR is ready to merge" }],
    [{ type: "add_labels", labels: ["mgmt-review-analyzed"] },
      { type: "add_comment", body: "## PR is ready to merge" }],
  ]) {
    const result = prepareAgentOutput(review, { items });
    assert.deepEqual(result.items.map((item) => item.type), ["add_comment", "remove_labels", "add_labels"]);
    assert.deepEqual(result.items.slice(1).map((item) => item.labels), [["mgmt-review-in-progress"], ["mgmt-review-analyzed"]]);
  }
});

test("missing comments and unrelated label actions cannot complete a review", async () => {
  const review = await fixture().run();
  for (const items of [
    [],
    [{ type: "add_comment", body: "one" }, { type: "add_comment", body: "two" }],
    [{ type: "add_comment", body: "review" }, { type: "remove_labels", labels: ["mgmt-review-needed"] }],
  ]) {
    assert.throws(() => prepareAgentOutput(review, { items }), /exactly one comment/);
  }
});

test("CI classification follows the existing squash-merge eligibility exception", async () => {
  const state = fixture();
  state.checks[0].conclusion = "failure";
  const eligible = await state.run();
  assert.equal(eligible.pipeline.classification_required, false);
  state.pr.mergeable_state = "blocked";
  const blocked = await state.run();
  assert.equal(blocked.pipeline.classification_required, true);
  assertReviewResults(blocked, state.provenance);
});
