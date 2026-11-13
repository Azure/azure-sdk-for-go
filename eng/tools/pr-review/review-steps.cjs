// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

const { gatherModules } = require("./pr-module-review.cjs");
const { markReadyForReview } = require("./draft-review.cjs");
const { checkPipelines } = require("./pipeline-review.cjs");
const { checkOnboarding } = require("./onboarding-review.cjs");
const { checkBreakingChanges } = require("./breaking-change-review.cjs");
const { textReader } = require("./github-data.cjs");
const { authorization, devOpsReader, fetchSpecText } = require("./azure-devops-read.cjs");
const { assertReviewCurrent } = require("./generated-provenance.cjs");
const { ciFailuresBlock } = require("./review-output.cjs");

async function runReviewSteps(api, { prNumber, headSHA, provenance, azureDevOpsPAT, adoRead, fetchText = fetchSpecText }) {
  const result = {
    pr_number: prNumber, head_sha: headSHA, base_sha: provenance.base_sha, base_ref: provenance.base_ref,
    status: "unverifiable", modules: [], changed_file_count: null, draft: null, pipeline: null,
    onboarding: [], breaking_changes: [], mergeable: null, mergeable_state: null,
    squash_merge_allowed: null, mergeability_errors: [], errors: [],
  };
  try {
    await assertReviewCurrent(api, provenance, prNumber, headSHA);
    if (provenance.status !== "pass") throw new Error("Other review steps require passing provenance.");
    const gathered = await gatherModules(api, prNumber, headSHA);
    const { pr } = gathered;
    if (pr.base.sha !== result.base_sha || pr.base.ref !== result.base_ref) throw new Error("The PR base changed.");
    result.modules = gathered.modules;
    result.changed_file_count = pr.changed_files;
    result.mergeable = pr.mergeable;
    result.mergeable_state = pr.mergeable_state;
    if (typeof pr.mergeable !== "boolean" || typeof pr.mergeable_state !== "string" || pr.mergeable_state === "unknown") {
      const message = "GitHub mergeability is not available yet; request a fresh review when it is computed.";
      result.mergeability_errors.push(message);
      api.core.warning(message);
    }
    result.draft = await markReadyForReview(api, pr, headSHA);
    result.pipeline = await checkPipelines(api, headSHA);
    if (result.pipeline.status === "pending") {
      result.status = "pending";
    } else {
      if (result.pipeline.failures.length > 0) {
        try {
          const { data: repo } = await api.github.rest.repos.get(api.context.repo);
          if (typeof repo.allow_squash_merge !== "boolean") throw new Error("The repository squash-merge setting is unavailable.");
          result.squash_merge_allowed = repo.allow_squash_merge;
        } catch (error) {
          const message = error?.message || String(error);
          result.mergeability_errors.push(message);
          api.core.error(`Mergeability verification failed: ${message}`);
        }
      }
      result.pipeline.classification_required = ciFailuresBlock(result);
      const readText = textReader(api, headSHA);
      if (result.modules.some((module) => module.first_onboard)) {
        adoRead ||= devOpsReader(authorization(api.core, azureDevOpsPAT));
        result.onboarding = await checkOnboarding({ ...api, prNumber, commentCount: pr.comments },
          result.modules, pr.body, readText, adoRead);
      }
      for (const module of result.modules) {
        const review = await checkBreakingChanges(module.path, readText, fetchText);
        for (const error of review.errors) api.core.error(`Breaking-change review for ${module.path}: ${error}`);
        result.breaking_changes.push(review);
      }
      result.status = "completed";
    }
    await assertReviewCurrent(api, provenance, prNumber, headSHA);
  } catch (error) {
    result.status = "unverifiable";
    const message = error?.message || String(error);
    result.errors.push(message);
    api.core.error(`Deterministic review could not be completed: ${message}`);
  }
  return result;
}

module.exports = { runReviewSteps };
