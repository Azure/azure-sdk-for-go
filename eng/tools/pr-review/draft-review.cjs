// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

async function markReadyForReview({ github }, pr, headSHA) {
  if (typeof pr.draft !== "boolean") throw new Error("The PR draft state is unavailable.");
  if (!pr.draft) return { status: "ready", converted: false };
  if (typeof pr.node_id !== "string" || !pr.node_id) throw new Error("The PR node ID is unavailable.");
  const data = await github.graphql(`
    mutation($id: ID!) {
      markPullRequestReadyForReview(input: {pullRequestId: $id}) {
        pullRequest { id isDraft headRefOid baseRefOid baseRefName }
      }
    }`, { id: pr.node_id });
  const current = data.markPullRequestReadyForReview?.pullRequest;
  if (current?.id !== pr.node_id || current.isDraft !== false || current.headRefOid !== headSHA ||
      current.baseRefOid !== pr.base.sha || current.baseRefName !== pr.base.ref) {
    throw new Error("Draft conversion failed or the PR changed during conversion.");
  }
  return { status: "ready", converted: true };
}

module.exports = { markReadyForReview };
