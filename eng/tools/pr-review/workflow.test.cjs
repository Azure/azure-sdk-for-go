// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

const assert = require("node:assert/strict");
const { readFileSync } = require("node:fs");
const path = require("node:path");
const { test } = require("node:test");
const { routeReviewRequest, prepareReview } = require("./review-request.cjs");

const sha = "a".repeat(40);
const nextSha = "b".repeat(40);
const created = "2026-10-19T19:00:00Z";
const started = "2026-10-19T19:00:02Z";
const needed = "mgmt-review-needed";
const progress = "mgmt-review-in-progress";
const bot = { id: 7, login: "github-actions[bot]", type: "Bot" };

function fixture() {
  const state = {
    context: {
      eventName: "workflow_run",
      actor: "maintainer",
      ref: "refs/heads/main",
      repo: { owner: "Azure", repo: "azure-sdk-for-go" },
      payload: {
        repository: { id: 1, full_name: "Azure/azure-sdk-for-go", default_branch: "main" },
        workflow_run: { id: 100 },
      },
    },
    run: {
      id: 100,
      run_attempt: 1,
      created_at: started,
      workflow_id: 10,
      repository: { id: 1 },
      head_repository: { id: 2 },
      head_sha: sha,
      actor: { id: 3, login: "maintainer", type: "User" },
      triggering_actor: bot,
      event: "pull_request",
      status: "completed",
      conclusion: "success",
      pull_requests: [],
    },
    pr: {
      number: 42,
      state: "open",
      base: { repo: { id: 1 } },
      head: { repo: { id: 2 }, sha },
      labels: [{ name: needed }],
    },
    events: [
      {
        id: 1000,
        event: "labeled",
        actor: { id: 3, login: "maintainer", type: "User" },
        created_at: created,
        label: { name: needed },
      },
    ],
    jobs: [{ name: "request", status: "completed", conclusion: "success" }],
    permission: "write",
    calls: { dispatches: [], additions: [], removals: [], actors: [], attempts: [], logs: [] },
  };
  state.associated = [structuredClone(state.pr)];
  state.core = {
    info: (message) => state.calls.logs.push(message),
    error: (message) => state.calls.logs.push(message),
  };
  state.changeLabel = (name, event, actor = bot) => {
    state.pr.labels = state.pr.labels.filter((label) => label.name !== name);
    if (event === "labeled") state.pr.labels.push({ name });
    state.events.push({
      id: Math.max(2000, ...state.events.map((event) => event.id)) + 1,
      event,
      label: { name },
      actor,
      created_at: started,
    });
  };
  state.github = {
    rest: {
      actions: {
        getWorkflow: async ({ workflow_id }) => {
          assert.equal(workflow_id, "pr-review-intake.yml");
          return { data: { id: 10 } };
        },
        getWorkflowRun: async ({ run_id }) => {
          assert.equal(run_id, 100);
          return { data: state.run };
        },
        listJobsForWorkflowRunAttempt: async ({ attempt_number }) => {
          state.calls.attempts.push(attempt_number);
          return { data: { jobs: state.jobs } };
        },
        createWorkflowDispatch: async (request) => state.calls.dispatches.push(request),
      },
      repos: {
        getCollaboratorPermissionLevel: async ({ username }) => {
          state.calls.actors.push(username);
          return { data: { permission: state.permission } };
        },
        listPullRequestsAssociatedWithCommit: async ({ commit_sha }) => {
          assert.equal(commit_sha, sha);
          return { data: structuredClone(state.associated) };
        },
      },
      pulls: {
        get: async ({ pull_number }) => {
          assert.equal(pull_number, 42);
          return { data: structuredClone(state.pr) };
        },
      },
      issues: {
        listEvents: async () => ({ data: structuredClone(state.events) }),
        addLabels: async (request) => {
          state.calls.additions.push(request);
          for (const label of request.labels) state.changeLabel(label, "labeled");
        },
        removeLabel: async (request) => {
          state.calls.removals.push(request);
          state.changeLabel(request.name, "unlabeled");
        },
      },
    },
    paginate: async (method, parameters) => {
      const { data } = await method(parameters);
      return Array.isArray(data) ? data : data.jobs;
    },
  };
  return state;
}

function worker(automatic = true) {
  const state = fixture();
  state.context.eventName = "workflow_dispatch";
  state.context.actor = automatic ? bot.login : "maintainer";
  state.context.payload.sender = automatic ? bot : state.run.actor;
  state.context.payload.inputs = {
    item_number: "42",
    head_sha: automatic ? sha : "",
    request_run_id: automatic ? "100" : "",
    request_event_id: automatic ? "1000" : "",
  };
  return state;
}

test("routes fork reviews using GitHub metadata, not run.pull_requests or artifacts", async () => {
  const state = fixture();
  await routeReviewRequest(state);
  assert.deepEqual(state.calls.dispatches, [{
    ...state.context.repo,
    workflow_id: "mgmt-review.lock.yml",
    ref: "main",
    inputs: { item_number: "42", head_sha: sha, request_run_id: "100", request_event_id: "1000" },
  }]);
  assert.deepEqual(state.calls.actors, ["maintainer"]);
  assert.equal(state.calls.additions.length, 0);
});

for (const permission of ["admin", "maintain", "write", "triage", "read", "none"]) {
  test(`review authorization checks ${permission} access`, async () => {
    const state = fixture();
    state.permission = permission;
    if (["admin", "maintain", "write"].includes(permission)) {
      await routeReviewRequest(state);
      assert.equal(state.calls.dispatches.length, 1);
    } else {
      await assert.rejects(routeReviewRequest(state), /needs repository write/);
      assert.equal(state.calls.dispatches.length, 0);
    }
  });
}

for (const [name, mutate] of [
  ["another workflow", (s) => s.run.workflow_id++],
  ["another repository", (s) => s.run.repository.id++],
  ["a non-PR run", (s) => s.run.event = "push"],
  ["a failed run", (s) => s.run.conclusion = "failure"],
  ["an unfinished run", (s) => s.run.status = "in_progress"],
]) {
  test(`review router rejects ${name}`, async () => {
    const state = fixture();
    mutate(state);
    await assert.rejects(routeReviewRequest(state));
    assert.equal(state.calls.dispatches.length, 0);
  });
}

test("routers and workers only run trusted default-branch code", async () => {
  const review = fixture();
  review.context.ref = "refs/heads/untrusted";
  await assert.rejects(routeReviewRequest(review), /default branch/);
  const dispatch = worker();
  dispatch.context.ref = "refs/heads/untrusted";
  await assert.rejects(prepareReview(dispatch, "mgmt-review"), /default branch/);
});

for (const jobs of [
  [],
  [{ name: "other", status: "completed", conclusion: "success" }],
  [{ name: "request", status: "completed", conclusion: "failure" }],
  [{ name: "request", status: "in_progress", conclusion: null }],
  Array(2).fill({ name: "request", status: "completed", conclusion: "success" }),
]) {
  test(`review router rejects invalid current-attempt jobs: ${JSON.stringify(jobs)}`, async () => {
    const state = fixture();
    state.jobs = jobs;
    await assert.rejects(routeReviewRequest(state));
  });
}

test("skipped review intake cannot authorize a review on rerun", async () => {
  const state = fixture();
  state.run.run_attempt = 2;
  state.jobs[0].conclusion = "skipped";
  await routeReviewRequest(state);
  assert.equal(state.calls.dispatches.length, 0);
  assert.deepEqual(state.calls.attempts, [2]);
});

for (const [name, mutate] of [
  ["absent", (s) => s.events = []],
  ["removed", (s) => s.events[0].event = "unlabeled"],
  ["another actor", (s) => s.events[0].actor.id++],
  ["too old", (s) => s.events[0].created_at = "2026-10-19T18:55:01Z"],
  ["after run creation", (s) => s.events[0].created_at = "2026-10-19T19:00:03Z"],
]) {
  test(`a ${name} label event cannot authorize routing or execution`, async () => {
    const router = fixture();
    mutate(router);
    await routeReviewRequest(router);
    assert.equal(router.calls.dispatches.length, 0);
    const dispatch = worker();
    mutate(dispatch);
    assert.equal(await prepareReview(dispatch, "mgmt-review"), undefined);
    assert.equal(dispatch.calls.additions.length, 0);
  });
}

test("the original requester, not a maintainer rerunning the intake, is authorized", async () => {
  const state = fixture();
  state.run.actor.login = "outsider";
  state.events[0].actor.login = "outsider";
  state.permission = "read";
  await assert.rejects(routeReviewRequest(state), /outsider needs/);
});

test("review requests accept the exact five-minute boundary, not one millisecond older", async () => {
  const state = fixture();
  state.events[0].created_at = "2026-10-19T18:55:02.000Z";
  await routeReviewRequest(state);
  assert.equal(state.calls.dispatches.length, 1);
  state.calls.dispatches = [];
  state.events[0].created_at = "2026-10-19T18:55:01.999Z";
  await routeReviewRequest(state);
  assert.equal(state.calls.dispatches.length, 0);
});

test("Go management review retains its existing no-bot-request policy", async () => {
  const state = fixture();
  state.events[0].actor = { ...state.run.actor, type: "Bot", login: "azure-sdk-automation[bot]" };
  await assert.rejects(routeReviewRequest(state), /not authorized/);
});

test("ambiguous commit associations cannot select an arbitrary review target", async () => {
  const state = fixture();
  state.associated.push({ ...state.pr, number: 43 });
  await assert.rejects(routeReviewRequest(state), /Multiple PRs/);
  assert.equal(state.calls.dispatches.length, 0);
});

test("workers independently validate and claim exactly the routed event", async () => {
  const state = worker();
  assert.deepEqual(await prepareReview(state, "mgmt-review"), { number: 42, headSha: sha });
  assert.deepEqual(state.pr.labels, [{ name: progress }]);
  assert.deepEqual(state.calls.actors, ["maintainer"]);
  assert.deepEqual(state.calls.removals.map((request) => request.name), [needed]);
});

test("manual dispatch can select the current head without an intake or pending label", async () => {
  const state = worker(false);
  state.pr.labels = [];
  assert.deepEqual(await prepareReview(state, "mgmt-review"), { number: 42, headSha: sha });
  assert.deepEqual(state.pr.labels, [{ name: progress }]);
});

for (const [name, mutate, error] of [
  ["bot without intake", (s) => s.context.payload.inputs.request_run_id = "", /must include an intake/],
  ["mismatched SHA", (s) => s.context.payload.inputs.head_sha = nextSha, /does not match/],
  ["invalid PR number", (s) => s.context.payload.inputs.item_number = "42; echo bad", /positive integer/],
  ["invalid event ID", (s) => s.context.payload.inputs.request_event_id = "bad", /positive integer/],
]) {
  test(`workers reject ${name} before claiming a label`, async () => {
    const state = worker();
    mutate(state);
    await assert.rejects(prepareReview(state, "mgmt-review"), error);
    assert.equal(state.calls.additions.length, 0);
  });
}

test("replaced label event IDs and changed PR heads cannot be replayed", async () => {
  const replaced = worker();
  replaced.events[0].id++;
  assert.equal(await prepareReview(replaced, "mgmt-review"), undefined);
  const stale = worker();
  stale.pr.head.sha = nextSha;
  assert.equal(await prepareReview(stale, "mgmt-review"), undefined);
});

test("label replacement during deletion is restored and in-progress state is rolled back", async () => {
  const state = worker();
  const remove = state.github.rest.issues.removeLabel;
  state.github.rest.issues.removeLabel = async (request) => {
    if (request.name === needed) {
      state.changeLabel(needed, "unlabeled", state.run.actor);
      state.changeLabel(needed, "labeled", state.run.actor);
    }
    await remove(request);
  };
  await assert.rejects(prepareReview(state, "mgmt-review"), /changed while claiming/);
  assert.deepEqual(state.pr.labels, [{ name: needed }]);
});

test("failure to mark in-progress preserves the pending request and surfaces the error", async () => {
  const state = worker();
  state.github.rest.issues.addLabels = async () => { throw new Error("Labels API unavailable"); };
  await assert.rejects(prepareReview(state, "mgmt-review"), /Labels API unavailable/);
  assert.deepEqual(state.pr.labels, [{ name: needed }]);
  assert.equal(state.calls.removals.length, 0);
});

test("API errors propagate rather than starting a review", async () => {
  const state = worker();
  state.github.rest.pulls.get = async () => { throw new Error("PR API unavailable"); };
  await assert.rejects(prepareReview(state, "mgmt-review"), /PR API unavailable/);
});

const workflows = path.resolve(__dirname, "..", "..", "..", ".github", "workflows");
const text = (name) => readFileSync(path.join(workflows, name), "utf8").replace(/\r\n/g, "\n");

test("review intake has no permissions, secrets, checkout, or artifact exchange", () => {
  const intake = text("pr-review-intake.yml");
  assert.match(intake, /permissions: \{\}/);
  assert.doesNotMatch(intake, /uses:|secrets\.|concurrency:/);
  const router = text("pr-review-router.yml");
  assert.match(router, /ref: \$\{\{ github.sha \}\}/);
  assert.doesNotMatch(router, /download-artifact|workflow_run\.head_sha/);
});

test("management review workflows do not subscribe to the legacy PR target event", () => {
  for (const file of ["pr-review-intake.yml", "pr-review-router.yml", "mgmt-review.md", "mgmt-review.lock.yml"]) {
    assert.doesNotMatch(text(file), /^  pull_request_target:/m, file);
  }
});

test("reviewer preserves pending requests in the workflow-level concurrency queue", () => {
  for (const file of ["mgmt-review.md", "mgmt-review.lock.yml"]) {
    const concurrency = text(file).match(/^concurrency:\n((?:[ \t]+[^\n]*\n)+)/m)?.[1];
    assert.ok(concurrency, `${file} has workflow-level concurrency`);
    assert.match(concurrency, /^  group: "?gh-aw-\$\{\{ github.workflow \}\}-\$\{\{ github.event.inputs.item_number \}\}"?$/m);
    assert.match(concurrency, /^  cancel-in-progress: false$/m);
    assert.match(concurrency, /^  queue: max$/m);
  }
});

test("compiled reviewer retains validated targets in both agent and publisher configs", () => {
  const source = text("mgmt-review.md");
  const compiled = text("mgmt-review.lock.yml");
  assert.match(source, /bots: \[github-actions\]/);
  assert.match(source, /cancel-in-progress: false/);
  assert.match(source, /do not check out or execute PR code/);
  assert.match(compiled, /^  agent:\n    needs:\n(?:      - [\w-]+\n)*      - validate_request\n/m);
  assert.match(compiled, /^  safe_outputs:\n    needs:\n(?:      - [\w-]+\n)*      - validate_request\n/m);
  const config = (key) => JSON.parse(JSON.parse(
    compiled.match(new RegExp(`^\\s+${key}: (.+)$`, "m"))[1],
  ));
  const agent = config("GH_AW_SAFE_OUTPUTS_CONFIG");
  const publisher = config("GH_AW_SAFE_OUTPUTS_HANDLER_CONFIG");
  for (const handler of ["add_comment", "add_labels", "remove_labels"]) {
    assert.equal(agent[handler].target, "${{ needs.validate_request.outputs.pr_number }}");
    assert.deepEqual(agent[handler], publisher[handler]);
  }
  assert.deepEqual(agent.add_labels.allowed, ["mgmt-review-analyzed"]);
  assert.deepEqual(agent.remove_labels.allowed, ["mgmt-review-in-progress"]);
  const safeOutputs = compiled.slice(compiled.indexOf("\n  safe_outputs:\n"));
  assert.ok(safeOutputs.indexOf("Reject stale review outputs") < safeOutputs.indexOf("Process Safe Outputs"));
});

test("the actual safe-output guard cleans up rejected reviews without consuming new requests", async (t) => {
  const source = text("mgmt-review.md");
  const script = source.split("    - name: Reject stale review outputs\n")[1]
    .split("  threat-detection:")[0].split("        script: |\n")[1]
    .split("\n").map((line) => line.slice(10)).join("\n");
  const AsyncFunction = Object.getPrototypeOf(async function () {}).constructor;
  const guard = new AsyncFunction("github", "context", "process", "core", script);
  const process = { env: { REVIEW_PR_NUMBER: "42", REVIEW_HEAD_SHA: sha } };
  const run = (state) => guard(state.github, state.context, process, state.core);

  await t.test("current head proceeds without removing labels", async () => {
    const state = fixture();
    state.pr.labels = [{ name: needed }, { name: progress }];
    await run(state);
    assert.deepEqual(state.pr.labels, [{ name: needed }, { name: progress }]);
    assert.equal(state.calls.removals.length, 0);
  });

  for (const [name, mutate] of [
    ["changed head", (state) => state.pr.head.sha = nextSha],
    ["closed PR", (state) => state.pr.state = "closed"],
  ]) {
    await t.test(name, async () => {
      const state = fixture();
      state.pr.labels = [{ name: needed }, { name: progress }];
      mutate(state);
      await assert.rejects(run(state), /No review outputs/);
      assert.deepEqual(state.pr.labels, [{ name: needed }]);
      assert.deepEqual(state.calls.removals, [{
        ...state.context.repo,
        issue_number: 42,
        name: progress,
      }]);
    });
  }

  await t.test("already-absent in-progress label still rejects stale output", async () => {
    const state = fixture();
    state.pr.head.sha = nextSha;
    state.github.rest.issues.removeLabel = async (request) => {
      state.calls.removals.push(request);
      throw Object.assign(new Error("Label not found"), { status: 404 });
    };
    await assert.rejects(run(state), /No review outputs/);
    assert.deepEqual(state.pr.labels, [{ name: needed }]);
    assert.deepEqual(state.calls.removals.map((request) => request.name), [progress]);
    assert.ok(state.calls.logs.includes("The in-progress label was already removed."));
  });

  await t.test("cleanup failures are surfaced and no other labels are touched", async () => {
    const state = fixture();
    state.pr.labels = [{ name: needed }, { name: progress }];
    state.pr.head.sha = nextSha;
    state.github.rest.issues.removeLabel = async (request) => {
      state.calls.removals.push(request);
      throw Object.assign(new Error("Label cleanup denied"), { status: 403 });
    };
    await assert.rejects(run(state), /Label cleanup denied/);
    assert.deepEqual(state.pr.labels, [{ name: needed }, { name: progress }]);
    assert.deepEqual(state.calls.removals.map((request) => request.name), [progress]);
  });
});
