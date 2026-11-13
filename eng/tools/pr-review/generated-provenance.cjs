// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

const { createHash } = require("node:crypto");
const { matchesGlob } = require("node:path").posix;

function requireSHA(value) {
  if (typeof value !== "string" || !/^[a-f0-9]{40}$/i.test(value)) {
    throw new Error("Missing or invalid commit/blob SHA.");
  }
  return value;
}

function requireCount(value, name) {
  if (!Number.isSafeInteger(value) || value < 0) throw new Error(`Missing or invalid ${name}.`);
  return value;
}

async function commitFiles(github, repo, sha) {
  const files = [];
  let commit;
  for (let page = 1; page <= 30; page++) {
    const { data } = await github.rest.repos.getCommit({ ...repo, ref: sha, per_page: 100, page });
    if (data.sha !== sha || !Array.isArray(data.files)) {
      throw new Error(`Incomplete commit response for ${sha}.`);
    }
    commit ??= data;
    if (!Array.isArray(commit.parents) || commit.parents.length > 1) {
      throw new Error(`Commit ${sha} is a merge or has incomplete parent metadata; manual review is required.`);
    }
    files.push(...data.files);
    if (data.files.length < 100) {
      if (new Set(files.map((file) => file.filename)).size !== files.length ||
          files.reduce((total, file) => total + requireCount(file.additions, "file additions"), 0) !== commit.stats?.additions ||
          files.reduce((total, file) => total + requireCount(file.deletions, "file deletions"), 0) !== commit.stats?.deletions) {
        throw new Error(`Incomplete or inconsistent file list for ${sha}.`);
      }
      return { commit, files };
    }
  }
  throw new Error(`Commit ${sha} reached GitHub's 3000-file response limit.`);
}

async function blob(github, repo, sha) {
  requireSHA(sha);
  const { data } = await github.rest.git.getBlob({ ...repo, file_sha: sha });
  if (data.sha !== sha || data.encoding !== "base64" || typeof data.content !== "string") {
    throw new Error(`Incomplete blob response for ${sha}.`);
  }
  const content = Buffer.from(data.content, "base64");
  const hash = createHash("sha1").update(`blob ${content.length}\0`).update(content).digest("hex");
  if (content.length !== data.size || hash !== sha) throw new Error(`Truncated or invalid blob ${sha}.`);
  return content;
}

async function previousBlob(github, repo, parent, filename) {
  const { data } = await github.rest.repos.getContent({ ...repo, path: filename, ref: parent });
  if (Array.isArray(data) || data.type !== "file") {
    throw new Error(`Cannot read ${filename} at parent ${parent} as a file.`);
  }
  return blob(github, repo, data.sha);
}

async function checkProvenance({ github, context, core }, { prNumber, headSHA, protectedFileMarker, filePatterns }) {
  const result = {
    status: "unverifiable",
    pr_number: prNumber,
    head_sha: headSHA,
    base_sha: null,
    base_ref: null,
    findings: [],
    errors: [],
  };
  try {
    requireSHA(headSHA);
    if (!Number.isSafeInteger(prNumber) || prNumber <= 0) throw new Error("Invalid PR number.");
    if (typeof protectedFileMarker !== "string" || protectedFileMarker.trim().length === 0) {
      throw new Error("protectedFileMarker must be non-empty literal text.");
    }
    if (!Array.isArray(filePatterns) || filePatterns.length === 0 ||
        filePatterns.some((pattern) => typeof pattern !== "string" || pattern.trim().length === 0)) {
      throw new Error("filePatterns must contain at least one non-empty filename glob.");
    }
    const patterns = filePatterns.map((pattern) => pattern.includes("/") ? pattern : `**/${pattern}`);
    const repo = context.repo;
    const { data: pr } = await github.rest.pulls.get({ ...repo, pull_number: prNumber });
    if (pr.state !== "open" || pr.head.sha !== headSHA) throw new Error("The PR head changed or the PR closed.");
    const baseSHA = requireSHA(pr.base.sha);
    if (typeof pr.base.ref !== "string" || pr.base.ref.length === 0) throw new Error("Missing PR base ref.");
    result.base_sha = baseSHA;
    result.base_ref = pr.base.ref;
    const commits = await github.paginate(github.rest.pulls.listCommits, { ...repo, pull_number: prNumber, per_page: 100 });
    if (commits.length === 0 || commits.length !== requireCount(pr.commits, "PR commit count") ||
        new Set(commits.map((commit) => commit.sha)).size !== commits.length || commits.at(-1).sha !== headSHA) {
      throw new Error("The PR commit list is incomplete or does not end at the validated head.");
    }
    for (const entry of commits) {
      const sha = requireSHA(entry.sha);
      // Use GitHub's linked login, not the user-controlled Git author display name.
      if (entry.author?.login === "azure-sdk") continue;
      const { commit, files: changed } = await commitFiles(github, repo, sha);
      const parent = commit.parents.length === 1 ? requireSHA(commit.parents[0].sha) : null;
      const author = entry.author?.login || "unlinked author";
      core.info(`Filtering ${changed.length} files in non-automation commit ${sha}.`);
      for (const file of changed) {
        if (typeof file.filename !== "string" || file.filename.length === 0 ||
            !["added", "removed", "modified", "renamed", "copied"].includes(file.status) ||
            (["renamed", "copied"].includes(file.status) &&
              (typeof file.previous_filename !== "string" || file.previous_filename.length === 0))) {
          throw new Error(`Incomplete file metadata in ${sha}.`);
        }
        if (![file.filename, file.previous_filename].some((name) =>
          typeof name === "string" && patterns.some((pattern) => matchesGlob(name, pattern)))) continue;
        if (!parent && file.status !== "added") throw new Error(`Missing parent for ${file.filename} in ${sha}.`);
        const before = file.status === "added" ? null :
          await previousBlob(github, repo, parent, file.previous_filename || file.filename);
        const after = file.status === "removed" ? null : await blob(github, repo, file.sha);
        const protectedFile = [before, after].some((content) => content?.includes(protectedFileMarker));
        if (protectedFile && (before === null || after === null || !before.equals(after) ||
            file.status === "renamed" || file.status === "copied")) {
          result.findings.push({ file: file.filename, sha, author });
        }
      }
    }
    const { data: current } = await github.rest.pulls.get({ ...repo, pull_number: prNumber });
    if (current.state !== "open" || current.head.sha !== headSHA ||
        current.base.sha !== result.base_sha || current.base.ref !== result.base_ref) {
      throw new Error("The PR head or base changed, or the PR closed during verification.");
    }
    result.status = result.findings.length === 0 ? "pass" : "manual-edits";
  } catch (error) {
    const message = error.message || String(error);
    core.error(`Generated-file provenance could not be verified: ${message}`);
    result.errors.push(message);
  }
  return result;
}

function markdown(value) {
  return String(value).replace(/[\\`*_{}[\]()<>!|#]/g, "\\$&").replace(/[\r\n]/g, " ");
}

function renderProvenance(result) {
  const lines = ["<!-- gh-aw-workflow-id: mgmt-review -->"];
  if (result.findings.length > 0) {
    lines.push("## Manual Edits to Auto-Generated Files Detected", "",
      "Non-automation commits changed protected files:", "");
    for (const finding of result.findings) {
      lines.push(`- ${markdown(finding.file)} - commit ${markdown(finding.sha.slice(0, 12))}, author ${markdown(finding.author)}`);
    }
  }
  if (result.status === "unverifiable") {
    lines.push("", "## Generated-File Review - Manual Review Required", "",
      "Generated-file provenance could not be verified:", "");
    for (const error of result.errors) lines.push(`- ${markdown(error)}`);
  }
  lines.push("", "**Action required:** A maintainer must review the listed edits or verification failures. Revert unauthorized generated-file changes or resolve the verification failure, then re-apply `mgmt-review-needed`.",
    "", "This PR is not ready to merge based on the generated-file provenance check.");
  return lines.join("\n");
}

function validateResult(result, prNumber, headSHA) {
  requireSHA(headSHA);
  const unknownBase = result?.base_sha === null && result.base_ref === null;
  if (!Number.isSafeInteger(prNumber) || prNumber <= 0 ||
      result?.pr_number !== prNumber || result.head_sha !== headSHA ||
      !["pass", "manual-edits", "unverifiable"].includes(result.status) ||
      !Array.isArray(result.findings) || !Array.isArray(result.errors) ||
      (unknownBase && result.status !== "unverifiable") ||
      (!unknownBase && (typeof result.base_sha !== "string" || !/^[a-f0-9]{40}$/i.test(result.base_sha) ||
        typeof result.base_ref !== "string" || result.base_ref.length === 0)) ||
      (result.status === "pass" && (result.findings.length > 0 || result.errors.length > 0)) ||
      (result.status === "manual-edits" && (result.findings.length === 0 || result.errors.length > 0)) ||
      (result.status === "unverifiable" && result.errors.length === 0)) {
    throw new Error("Missing, inconsistent, or stale deterministic provenance result; no review outputs may be published.");
  }
}

async function removeInProgress({ github, context, core }, prNumber) {
  try {
    await github.rest.issues.removeLabel({
      ...context.repo, issue_number: prNumber, name: "mgmt-review-in-progress",
    });
  } catch (error) {
    if (error.status !== 404) throw error;
    core.info("The in-progress label was already removed.");
  }
}

async function assertReviewCurrent(api, result, prNumber, headSHA) {
  validateResult(result, prNumber, headSHA);
  const { data: pr } = await api.github.rest.pulls.get({ ...api.context.repo, pull_number: prNumber });
  if (pr.state !== "open" || pr.head.sha !== headSHA ||
      (result.base_sha !== null && (pr.base?.sha !== result.base_sha || pr.base?.ref !== result.base_ref))) {
    await removeInProgress(api, prNumber);
    throw new Error("The PR head or base changed or the PR closed during review. No review outputs were published; request a new review.");
  }
}

async function publishProvenance(api, { result, prNumber, headSHA }) {
  await assertReviewCurrent(api, result, prNumber, headSHA);
  if (result.status === "pass") throw new Error("Passing provenance must use the normal review publisher.");
  await api.github.rest.issues.createComment({
    ...api.context.repo, issue_number: prNumber, body: renderProvenance(result),
  });
  await removeInProgress(api, prNumber);
  await api.github.rest.issues.addLabels({
    ...api.context.repo, issue_number: prNumber, labels: ["mgmt-review-analyzed"],
  });
}

module.exports = { checkProvenance, assertReviewCurrent, publishProvenance, renderProvenance };
