// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

const { requireCount } = require("./github-data.cjs");

async function gatherModules({ github, context }, pr) {
  const files = await github.paginate(github.rest.pulls.listFiles, {
    ...context.repo, pull_number: pr.number, per_page: 100,
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
  return [...modules.values()].sort((a, b) => a.path.localeCompare(b.path));
}

module.exports = { gatherModules };
