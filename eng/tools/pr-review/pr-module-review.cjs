// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

const { requireSHA, requireCount } = require("./github-data.cjs");

async function gatherModules({ github, context }, prNumber, headSHA) {
  requireSHA(headSHA);
  if (!Number.isSafeInteger(prNumber) || prNumber <= 0) throw new Error("Invalid PR number.");
  const { data: pr } = await github.rest.pulls.get({ ...context.repo, pull_number: prNumber });
  if (pr.number !== prNumber || pr.state !== "open" || pr.head?.sha !== headSHA) {
    throw new Error("The PR identity or head changed, or the PR closed.");
  }
  const files = await github.paginate(github.rest.pulls.listFiles, {
    ...context.repo, pull_number: prNumber, per_page: 100,
  });
  if (!Array.isArray(files) || files.length !== requireCount(pr.changed_files, "PR changed-file count") ||
      new Set(files.map((file) => file.filename)).size !== files.length) {
    throw new Error("The PR changed-file list is incomplete or inconsistent.");
  }
  const modules = new Map();
  for (const file of files) {
    if (typeof file.filename !== "string" || !file.filename ||
        !["added", "removed", "modified", "renamed", "copied", "changed", "unchanged"].includes(file.status) ||
        (["renamed", "copied"].includes(file.status) && typeof file.previous_filename !== "string")) {
      throw new Error("The PR changed-file list has incomplete file metadata.");
    }
    for (const filename of [file.filename, file.previous_filename]) {
      const path = filename?.match(/^(sdk\/resourcemanager\/[^/]+\/arm[^/]+)\//)?.[1];
      if (!path) continue;
      const module = modules.get(path) || { path, first_onboard: false };
      module.first_onboard ||= filename === file.filename && file.status === "added" && filename === `${path}/ci.yml`;
      modules.set(path, module);
    }
  }
  return {
    pr,
    modules: [...modules.values()].sort((a, b) => a.path.localeCompare(b.path)),
    files: files.map(({ filename, status, previous_filename }) => ({
      filename, status, ...(previous_filename ? { previous_filename } : {}),
    })),
  };
}

module.exports = { gatherModules };
