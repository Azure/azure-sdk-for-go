// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

function serviceToken(path) {
  const [mid, suffix] = path.split("/").slice(-2);
  const normalize = (value) => value.toLowerCase().replace(/^arm/, "").replace(/[^a-z0-9]/g, "");
  const left = normalize(mid);
  const right = normalize(suffix);
  return left === right ? left : left + right;
}

function releasePlanID(body) {
  const ids = new Set();
  for (const match of (body || "").matchAll(/https:\/\/[^\s<>()[\]"']+/g)) {
    const url = new URL(match[0]);
    if (!url.hostname.startsWith("azsdk-releaseplan-dashboard-") || !url.hostname.endsWith(".azurewebsites.net")) continue;
    const value = url.searchParams.get("releaseplan");
    if (!value || !/^[1-9]\d*$/.test(value) || !Number.isSafeInteger(Number(value))) {
      throw new Error("The Release Plan link has an invalid releaseplan value.");
    }
    ids.add(value);
  }
  if (ids.size > 1) throw new Error("Multiple Release Plan IDs need manual confirmation.");
  return [...ids][0];
}

async function namespaceEvidence(api, body, adoRead) {
  const notes = [];
  let issue = null;
  try {
    const id = releasePlanID(body);
    if (id) {
      const { data } = await adoRead("Release", "wit/wiql", {}, {
        query: `SELECT [System.Id] FROM WorkItems WHERE [System.WorkItemType] = 'Release Plan' AND [Custom.ReleasePlanID] = '${id}'`,
      });
      if (!Array.isArray(data.workItems) || data.workItems.length !== 1 ||
          !Number.isSafeInteger(data.workItems[0].id) || data.workItems[0].id <= 0) {
        throw new Error("No unique Release Plan matches Custom.ReleasePlanID.");
      }
      const { data: item } = await adoRead("Release", `wit/workitems/${data.workItems[0].id}`,
        { fields: "Custom.ReleasePlanID,Custom.NamespaceApprovalIssue" });
      if (String(item.fields?.["Custom.ReleasePlanID"]) !== id) throw new Error("The resolved Release Plan ID does not match.");
      const value = item.fields?.["Custom.NamespaceApprovalIssue"];
      if (value != null && typeof value !== "string") throw new Error("Invalid NamespaceApprovalIssue field.");
      if (typeof value === "string" && value.trim()) {
        return { status: "satisfied", source: "release-plan", issue: value, notes };
      }
    }
  } catch (error) {
    const message = error?.message || String(error);
    notes.push(message);
    api.core.warning(`Release Plan namespace signal unavailable: ${message}`);
  }
  const comments = await api.github.paginate(api.github.rest.issues.listComments, {
    ...api.context.repo, issue_number: api.prNumber, per_page: 100,
  });
  if (!Array.isArray(comments) || comments.some((comment) => typeof comment.body !== "string")) {
    throw new Error("Namespace-review comments are unavailable or invalid.");
  }
  if (!Number.isSafeInteger(api.commentCount) || comments.length !== api.commentCount ||
      new Set(comments.map((comment) => comment.id)).size !== comments.length) {
    throw new Error("The namespace-review comment list is incomplete or inconsistent.");
  }
  for (const comment of comments) {
    if (comment.body.includes("<!-- gh-aw-workflow-id: mgmt-review -->")) continue;
    const link = comment.body.match(/https:\/\/github\.com\/[\w.-]+\/[\w.-]+\/(?:issues|pull)\/[1-9]\d*(?=$|[\s/#?)>\]".,])/);
    if (/namespace/i.test(comment.body) && link) {
      issue = link[0];
      break;
    }
  }
  return { status: issue ? "satisfied" : "missing", source: issue ? "pr-comment" : null, issue, notes };
}

async function releasePipeline(path, repo, adoRead) {
  const definitions = [];
  const tokens = new Set();
  let continuation;
  do {
    const response = await adoRead("internal", "build/definitions", {
      repositoryId: `${repo.owner}/${repo.repo}`, repositoryType: "GitHub",
      yamlFilename: `${path}/ci.yml`, includeAllProperties: "true", processType: "2",
      ...(continuation ? { continuationToken: continuation } : {}),
    });
    if (!Array.isArray(response.data.value) || response.data.count !== response.data.value.length) {
      throw new Error("Invalid release-pipeline definition response.");
    }
    definitions.push(...response.data.value);
    continuation = response.continuation;
    if (continuation && (tokens.has(continuation) || tokens.size >= 100)) {
      throw new Error("Incomplete release-pipeline definition pagination.");
    }
    tokens.add(continuation);
  } while (continuation);
  if (new Set(definitions.map((definition) => definition.id)).size !== definitions.length ||
      definitions.some((definition) => !Number.isSafeInteger(definition.id) || definition.id <= 0 ||
      definition.repository?.id !== `${repo.owner}/${repo.repo}` ||
      definition.repository?.type !== "GitHub" ||
      definition.process?.yamlFilename?.replace(/^\//, "") !== `${path}/ci.yml`)) {
    throw new Error("Release-pipeline definitions do not match the module and repository.");
  }
  return { status: definitions.length ? "satisfied" : "missing", definition_ids: definitions.map((item) => item.id) };
}

async function checkOnboarding(api, modules, body, readText, adoRead) {
  const newModules = modules.filter((module) => module.first_onboard);
  if (!newModules.length) return [];
  let namespace;
  try {
    namespace = await namespaceEvidence(api, body, adoRead);
  } catch (error) {
    const message = error?.message || String(error);
    api.core.error(`Namespace approval could not be verified: ${message}`);
    namespace = { status: "unverifiable", errors: [message] };
  }
  const results = [];
  for (const { path } of newModules) {
    const item = { path, service_token: serviceToken(path), namespace, pipeline: null, readme: null };
    try {
      item.pipeline = await releasePipeline(path, api.context.repo, adoRead);
    } catch (error) {
      const message = error?.message || String(error);
      api.core.error(`Release pipelines for ${path} could not be verified: ${message}`);
      item.pipeline = { status: "unverifiable", errors: [message] };
    }
    try {
      const text = await readText(`${path}/README.md`);
      item.readme = { status: "semantic-review-required", text: text.split(/^##\s/m)[0] };
    } catch (error) {
      const message = error?.message || String(error);
      api.core.error(`README for ${path} could not be verified: ${message}`);
      item.readme = { status: "unverifiable", errors: [message] };
    }
    results.push(item);
  }
  return results;
}

module.exports = { checkOnboarding, namespaceEvidence, releasePipeline, serviceToken };
