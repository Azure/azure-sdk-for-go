// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

const { parseDocument } = require("yaml");
const { SemVer } = require("semver");
const { requireSHA } = require("./github-data.cjs");

function parseYAML(text, name) {
  const document = parseDocument(text, { uniqueKeys: true });
  if (document.errors.length) throw new Error(`Invalid ${name}: ${document.errors[0].code}.`);
  if (document.warnings.length) throw new Error(`Unsupported ${name}: ${document.warnings[0].code}.`);
  return document.toJS({ maxAliasCount: 50 });
}

function moduleVersion(source) {
  // Remove comments and non-string literals before extracting the SDK's literal const.
  const clean = source.replace(/\/\/[^\n]*|\/\*[\s\S]*?\*\/|`[^`]*`|'(?:\\.|[^'\\])*'|"(?:\\.|[^"\\])*"/g,
    (token) => token.startsWith('"') ? token : token.replace(/[^\n]/g, " "));
  const declarations = [...clean.matchAll(/^\s*const\s+moduleVersion(?:\s+string)?\s*=\s*"([^"\n]*)"\s*;?\s*$/gm)];
  for (const block of clean.matchAll(/\bconst\s*\(([\s\S]*?)\)/g)) {
    declarations.push(...block[1].matchAll(/^\s*moduleVersion(?:\s+string)?\s*=\s*"([^"\n]*)"\s*;?\s*$/gm));
  }
  if (declarations.length !== 1) throw new Error("Cannot determine a unique literal moduleVersion constant.");
  const version = declarations[0][1].replace(/^v/, "");
  const parsed = new SemVer(version);
  return { version, skip: parsed.major === 0 || parsed.prerelease.length > 0 ||
    (parsed.major === 1 && parsed.minor === 0 && parsed.patch === 0) };
}

function breakingItems(changelog, version) {
  const lines = changelog.split(/\r?\n/);
  const releases = lines.flatMap((line, index) => {
    const match = line.match(/^##\s+(\S+)\s+\([^)]*\)\s*$/);
    return match?.[1].replace(/^v/, "") === version ? [index] : [];
  });
  if (releases.length !== 1) throw new Error(`Cannot find a unique CHANGELOG release entry for ${version}.`);
  const start = releases[0] + 1;
  let end = lines.findIndex((line, index) => index >= start && /^##\s/.test(line));
  if (end === -1) end = lines.length;
  const release = lines.slice(start, end);
  const sections = release.flatMap((line, index) => /^### Breaking Changes\s*$/.test(line) ? [index] : []);
  if (sections.length === 0) return [];
  if (sections.length !== 1) throw new Error("The release has multiple Breaking Changes sections.");
  const items = [];
  for (const line of release.slice(sections[0] + 1)) {
    if (/^#{1,6}\s/.test(line)) break;
    if (!line.trim()) continue;
    const bullet = line.match(/^[-*+]\s+(.+)$/);
    if (bullet) items.push(bullet[1]);
    else if (/^\s+/.test(line) && items.length) items[items.length - 1] += `\n${line.trimStart()}`;
    else throw new Error("The Breaking Changes section contains non-bullet text.");
  }
  return items;
}

function approvals(data, path) {
  if (data == null) return [];
  if (typeof data !== "object" || Array.isArray(data)) throw new Error("Invalid suppression document structure.");
  const section = data.suppressions;
  if (section == null) return [];
  if (typeof section !== "object" || Array.isArray(section)) throw new Error("Invalid suppressions structure.");
  const entries = section["azure-sdk-for-go"];
  if (entries == null) return [];
  if (!Array.isArray(entries)) throw new Error("Invalid Go suppression list.");
  const approved = [];
  for (const entry of entries) {
    if (!entry || typeof entry !== "object" || Array.isArray(entry) || typeof entry.package !== "string") {
      throw new Error("Invalid Go suppression entry.");
    }
    if (entry.package !== path) continue;
    const changes = entry["breaking-changes"];
    if (changes == null) continue;
    if (!Array.isArray(changes) || changes.some((item) => typeof item !== "string")) {
      throw new Error("Invalid matching-package breaking-changes list.");
    }
    approved.push(...changes);
  }
  return approved;
}

async function checkBreakingChanges(path, readText, fetchText) {
  const result = { path, status: "unverifiable", version: null, items: [], source_url: null, errors: [] };
  try {
    const version = moduleVersion(await readText(`${path}/version.go`));
    result.version = version.version;
    if (version.skip) return { ...result, status: "skipped" };
    const items = breakingItems(await readText(`${path}/CHANGELOG.md`), version.version);
    if (items.length === 0) return { ...result, status: "skipped" };
    result.items = items.map((text) => ({ text, status: "Unverifiable", reason: "" }));
    const location = parseYAML(await readText(`${path}/tsp-location.yaml`), "tsp-location.yaml");
    const repo = typeof location?.repo === "string" &&
      location.repo.replace(/^https:\/\/github\.com\//, "").replace(/\.git$/, "");
    if (!repo || !/^[\w.-]+\/[\w.-]+$/.test(repo) ||
        repo.split("/").some((part) => part === "." || part === "..")) throw new Error("Missing or invalid TypeSpec repo.");
    const commit = requireSHA(location.commit);
    if (typeof location.directory !== "string" || !location.directory ||
        location.directory.split("/").some((part) => !part || part === "." || part === ".." || part.includes("\\"))) {
      throw new Error("Missing or invalid TypeSpec directory.");
    }
    const directory = location.directory.split("/").map(encodeURIComponent).join("/");
    result.source_url = `https://github.com/${repo}/blob/${commit}/${directory}/sdk-suppressions.yaml`;
    const response = await fetchText(`https://raw.githubusercontent.com/${repo}/${commit}/${directory}/sdk-suppressions.yaml`);
    result.source_missing = response.status === 404;
    if (response.status !== 404 && !(response.status >= 200 && response.status < 300)) {
      throw new Error(`Pinned specification fetch returned HTTP ${response.status}.`);
    }
    const allowed = response.status === 404 ? [] : approvals(parseYAML(response.text, "sdk-suppressions.yaml"), path);
    result.items = items.map((text) => ({
      text, status: allowed.includes(text) ? "Approved" : "Unapproved",
      reason: allowed.includes(text) ? "" : "No exact approval in the matching Go package's breaking-changes list.",
    }));
    result.status = result.items.every((item) => item.status === "Approved") ? "approved" : "unapproved";
  } catch (error) {
    const message = error?.message || String(error);
    result.errors.push(message);
    result.items = result.items.map((item) => ({ ...item, status: "Unverifiable", reason: message }));
  }
  return result;
}

module.exports = { checkBreakingChanges, moduleVersion, breakingItems };
