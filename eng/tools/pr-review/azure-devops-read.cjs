// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

const { execFileSync } = require("node:child_process");

function authorization(core, pat) {
  if (pat != null && (typeof pat !== "string" || (pat.length > 0 && !pat.trim()))) {
    throw new Error("The Azure DevOps PAT must be non-empty text when configured.");
  }
  if (pat) {
    core.setSecret(pat);
    const header = `Basic ${Buffer.from(`:${pat}`).toString("base64")}`;
    core.setSecret(header);
    return header;
  }
  try {
    const token = execFileSync("az", ["account", "get-access-token", "--resource",
      "499b84ac-1321-427f-aa17-267ca6975798", "--query", "accessToken", "-o", "tsv"],
    { encoding: "utf8", stdio: ["ignore", "pipe", "pipe"], timeout: 30000 }).trim();
    if (!token) throw new Error("Empty access token.");
    core.setSecret(token);
    return `Bearer ${token}`;
  } catch (error) {
    core.warning(`Azure DevOps authentication unavailable (CLI status ${error.status ?? error.code ?? "unknown"}); use PR namespace comments and manually verify inaccessible release pipelines.`);
    return null;
  }
}

function devOpsReader(header, fetchImpl = fetch) {
  return async (project, route, query = {}, body) => {
    if (!header) throw new Error("Azure DevOps read authentication is unavailable.");
    const url = new URL(`https://dev.azure.com/azure-sdk/${project}/_apis/${route}`);
    url.search = new URLSearchParams({ "api-version": "7.1", ...query }).toString();
    const response = await fetchImpl(url, {
      method: body ? "POST" : "GET", redirect: "error", signal: AbortSignal.timeout(30000),
      headers: { Authorization: header, "Content-Type": "application/json" },
      ...(body ? { body: JSON.stringify(body) } : {}),
    });
    if (!response.ok) throw new Error(`Azure DevOps ${route} returned HTTP ${response.status}.`);
    let data;
    try {
      data = await response.json();
    } catch {
      throw new Error(`Azure DevOps ${route} did not return valid JSON.`);
    }
    return { data, continuation: response.headers.get("x-ms-continuationtoken") };
  };
}

async function fetchSpecText(url, fetchImpl = fetch) {
  if (new URL(url).origin !== "https://raw.githubusercontent.com") throw new Error("Invalid pinned-spec host.");
  const response = await fetchImpl(url, { redirect: "error", signal: AbortSignal.timeout(30000) });
  if (response.status === 404) return { status: 404, text: "" };
  if (!response.ok) throw new Error(`Pinned specification fetch returned HTTP ${response.status}.`);
  return { status: response.status, text: await response.text() };
}

module.exports = { authorization, devOpsReader, fetchSpecText };
