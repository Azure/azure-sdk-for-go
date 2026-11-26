// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

const { createHash } = require("node:crypto");

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

function textReader({ github, context }, headSHA) {
  const files = new Map();
  return async (filename) => {
    if (!files.has(filename)) {
      files.set(filename, previousBlob(github, context.repo, headSHA, filename)
        .then((content) => new TextDecoder("utf-8", { fatal: true }).decode(content)));
    }
    return files.get(filename);
  };
}

module.exports = { requireSHA, requireCount, blob, previousBlob, textReader };
