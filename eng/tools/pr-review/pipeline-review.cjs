// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

const { requireCount } = require("./github-data.cjs");

function category(name) {
  if (/^go - pullrequest \(Build Build\/Test on .+_go_\d+\)$/.test(name)) return "build-test";
  if (name === "go - pullrequest (Build Analyze)") return "analyze";
  if (name === "go - pullrequest (Build generate_job_matrix)") return "matrix";
  return "other";
}

async function checkPipelines({ github, context }, headSHA) {
  const checks = [];
  let expected;
  for (let page = 1; page <= 100; page++) {
    const { data } = await github.rest.checks.listForRef({
      ...context.repo, ref: headSHA, filter: "latest", per_page: 100, page,
    });
    const total = requireCount(data.total_count, "check-run count");
    expected ??= total;
    if (total !== expected || !Array.isArray(data.check_runs)) throw new Error("Inconsistent check-run pagination.");
    checks.push(...data.check_runs);
    if (checks.length === expected) break;
    if (data.check_runs.length < 100 || checks.length > expected || page === 100) {
      throw new Error("The check-run list is incomplete.");
    }
  }
  if (new Set(checks.map((check) => check.id)).size !== checks.length ||
      checks.some((check) => !Number.isSafeInteger(check.id) || check.head_sha !== headSHA ||
        typeof check.name !== "string")) {
    throw new Error("The check-run list has invalid, duplicate, or stale metadata.");
  }
  const pipeline = checks.filter((check) => check.name === "go - pullrequest" ||
    check.name.startsWith("go - pullrequest (Build "));
  if (pipeline.some((check) => !["queued", "in_progress", "completed"].includes(check.status))) {
    throw new Error("A pipeline check has an invalid status.");
  }
  const parent = pipeline.filter((check) => check.name === "go - pullrequest");
  if (parent.length > 1) throw new Error("Multiple current go - pullrequest parent checks were returned.");
  if (parent.length === 0 || pipeline.some((check) => check.status !== "completed")) {
    return { status: "pending", failures: [] };
  }
  const conclusions = ["success", "neutral", "skipped", "failure", "cancelled", "timed_out", "action_required", "stale"];
  if (pipeline.some((check) => !conclusions.includes(check.conclusion))) {
    throw new Error("A completed pipeline check has no valid conclusion.");
  }
  return {
    status: "completed",
    failures: pipeline.filter((check) => !["success", "neutral", "skipped"].includes(check.conclusion)).map((check) => ({
      name: check.name,
      conclusion: check.conclusion,
      category: category(check.name),
      title: check.output?.title || "",
      summary: check.output?.summary || "",
      target_url: check.details_url || check.target_url || null,
    })),
  };
}

module.exports = { checkPipelines };
