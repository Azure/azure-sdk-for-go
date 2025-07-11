# Go SDK Breaking Changes Review

Follow these steps to help users review the Go SDK breaking changes. Show the high-level description to users to help them understand the flow. Use the provided tools to perform actions and gather information as needed.

## Overview

1. **Generate Go SDK locally** - Reproduce the breaking changes locally by generating the Go SDK from the API specification.
2. **Gather review info** - Gather enough info before starting the review.
3. **Review breaking changes** - Categorize each breaking change and find correlated change of API specification.
4. **Resolve breaking changes** - Use client customization to address the breaking changes.
5. **Gather results** - Gather results and provide report to user.

## Step 1: Generate the SDK Locally

Follow [Generate Go SDK from API specification](./go-sdk-generation.instructions.md) to generate the SDK from the API specification.

**Prerequisites**: Ensure the SDK generation completed successfully with breaking changes detected.

## Step 2: Gather Review Information

### Extract Generation Result

**Goal**: Identify the type of SDK generation and extract breaking change information.

**Actions**:

1. Read `generateOutput.json` file from the workspace root folder.
2. Extract the following information:
   - Package path (`<package-path>`)
   - Changelog content (`<changelog-content>`)
   - Breaking change items (`<breaking-change-items>`)

**Error Handling**: If `generateOutput.json` is missing or corrupted, regenerate the SDK first.

**Success Criteria**: Successfully extract `<package-path>`, `<changelog-content>`, and `<breaking-change-items>`.

### Determine SDK Generation Type

**Goal**: Determine the type of SDK generation based on the changed files.

**Actions**:

1. Run `git diff` command to analyze changes under `<package-path>` in current workspace.
2. Determine the generation type based on the following criteria:
   - **"TypeSpec Migration"**: If diff contains removal of `autorest.md` file
   - **"TypeSpec Update"**: If diff contains any changes to `*.tsp` files
   - **"Update Swagger"**: All other cases

**Success Criteria**: Determine `<sdk-generation-type>` from the three possible values.

## Step 3: Review Breaking Changes

### Categorize Breaking Changes

**Goal**: Group breaking change items, find the reasons and impacts, and identify potential resolutions.

**Actions**:

1. For items in `<breaking-change-items>`:
   - Follow the [Azure Go SDK Breaking Changes Review and Resolution Guide](../../documentation/sdk-breaking-changes-guide.md) with `<sdk-generation-type>` section
   - Group the breaking changes into the corresponding categories
     - If no category matches, keep the items as "Unknown Breaking Change" and inform user to have manual review
   - Determine the root cause and impact
   - Identify potential resolution approaches

**Success Criteria**: All breaking change items are categorized with type, reason, and potential resolution.

## Step 4: Resolve Breaking Changes

**Goal**: Resolve the breaking changes using client customization if applicable.

**Actions**:

1. Ask the user: "Would you like to resolve the breaking changes using client customization?"

2. If the user responds **"No"**:

   - Skip to Step 5 (Gather Review Result)

3. If the user responds **"Yes"**:

   - If customizations are required, check if `client.tsp` exists in <related-typeSpec-project-folder>
     - **If `client.tsp` exists**:
       - Review existing customizations
       - Propose modifications to resolve breaking changes
     - **If `client.tsp` does not exist**:
       - Create a new `client.tsp` file
       - Import `main.tsp` file in `client.tsp`
       - Using the root namespace defined in `main.tsp`
       - Add customizations to resolve breaking changes
   - If config changes are needed, propose modifications to `tsconfig.yaml` in <related-typeSpec-project-folder>

4. Revert all changes under <package-path>, and generate the SDK locally again to see if the breaking changes are resolved.

**Error Handling**: If customization cannot resolve specific breaking changes, document the limitation and suggest alternative approaches.

**Success Criteria**: Breaking changes are resolved through customization where possible, or limitations are documented.

## Step 5: Gather Review Result

**Goal**: Provide comprehensive breaking changes analysis report.

**Actions**:

Generate the result in the following format:

#### Breaking Changes Summary

```
⚠️ **Breaking Changes Detected**: [count of <breaking-change-items>] changes found
📦 **Affected Package**: [<package-path>]
🔧 **SDK Generation Type**: [<sdk-generation-type>]
```

#### Individual Analysis

For a group of breaking change items:

```
🚨 **Breaking Changes**: [changes-changelog]
📋 **Reason**: [breaking-changes-reason]
🔧 **Resolution**: [resolution-method]
💡 **Impact**: [impact-assessment]
```

#### Overall Assessment

```
✅ **Resolvable**: [count] breaking changes can be resolved through customization
⚠️ **Requires Attention**: [count] breaking changes need manual intervention
📋 **Next Steps**: [recommended-actions]
```

**Success Criteria**: Complete analysis report is generated with all available information and clear next steps.
