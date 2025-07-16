# Go SDK Breaking Changes Review

Follow these steps to help users review the Go SDK breaking changes. Show the high-level description to users to help them understand the flow. Use the provided tools to perform actions and gather information as needed.

## Step 1: Generate the SDK Locally

Follow [Generate Go SDK from API specification](./go-sdk-generation.instructions.md) to generate the SDK from the API specification.

**Prerequisites**: Ensure the SDK generation completed successfully with breaking changes detected.

## Step 2: Review Breaking Changes

Follow the [Azure Go SDK Breaking Changes Review and Resolution Guide](../../documentation/sdk-breaking-changes-guide.md) with correct generation type section to review the breaking changes.

If a breaking change could be resolved, try to update the spec file to resolve the breaking change. If the spec file is not editable, document the breaking change and its impact.

If any changes are made to the spec file, regenerate the SDK locally and check if the breaking changes are resolved.

## Step 3: Gather Review Result

Provide comprehensive breaking changes analysis report.

### Breaking Changes Summary

```
⚠️ **Breaking Changes Detected**: [count of <breaking-change-items>] changes found
📦 **Affected Package**: [<package-path>]
🔧 **SDK Generation Type**: [<sdk-generation-type>]
```

### Individual Analysis

For a group of breaking change items:

```
🚨 **Breaking Changes**: [changes-changelog]
📋 **Reason**: [breaking-changes-reason]
🔧 **Resolution**: [resolution-method]
💡 **Impact**: [impact-assessment]
```

### Overall Assessment

```
✅ **Resolvable**: [count] breaking changes can be resolved through customization
⚠️ **Requires Attention**: [count] breaking changes need manual intervention
📋 **Next Steps**: [recommended-actions]
```
