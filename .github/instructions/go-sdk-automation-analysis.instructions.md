# Analyze the Go SDK Automation Result

Follow these steps to help users analyze Go SDK automation results from GitHub PRs. This guide focuses on interpreting results and providing resolution guidance.

## Process Overview

1. **Generate Go SDK locally** - Reproduce the automation result locally by generating the Go SDK from the API specification.
2. **Provide guidance** - Offer specific resolution steps and resources.

## Step 1: Generate Go SDK

1. Follow [Generate Go SDK from API specification](./go-sdk-generation.instructions.md) to reproduce the automation result locally
2. Collect all generation logs, output files, and error messages for analysis

## Step 2: Analyze Results

1. **Success Detection**: If the log contains `Finish processing typespec project` keyword and no `[ERROR]` keywords are found, then generation succeeded. Check the output and summarize the results.

2. **Error Handling**: If any `[ERROR]` logs or stack traces are found, follow the [Azure Go SDK Automation Troubleshooting Guide](../../documentation/sdk-automation-tsg.md) to suggest how to resolve the errors.

### Step 3: Resolution Guidance

1. Provide specific actionable steps based on error category
2. Include relevant documentation links
3. Format issue templates for external reporting when needed

## Output Format

### Success Analysis

```
✅ **Status**: Automation Successful
📦 **Package**: [package-name]
📋 **Changelog**: [summary of changes]
⚠️ **Breaking Changes**: [if any, with details]
```

### Failure Analysis

```
❌ **Status**: Automation Failed
📦 **Package**: [package-name]
🔍 **Error Category**: [category-name]
💬 **Error Summary**: [concise description]
🛠️ **Resolution**: [specific action items]
📎 **Additional Info**: [links, references, issue templates]
```
