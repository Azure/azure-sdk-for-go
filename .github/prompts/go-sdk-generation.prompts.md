# Generate Go SDK from API specification

Follow these steps to help users generate Go SDK from specific API specifications. Show the high-level description to users to help them understand the flow. Use the provided tools to perform actions and gather information as needed.

## Step 1: Prerequisites

Use tool to check environment and ensure SDK repo and spec repo are accessible.

## Step 2: Generate Go SDK

Use tool to generate Go SDK from the provided API specification.

## Step 3: Gather Results and Handle Errors

Provide a summary of the results and handle any errors according to the [Azure Go SDK Automation Troubleshooting Guide](../../documentation/sdk-automation-tsg.md) if generation fails.

## Output Format

### Success Analysis

```
✅ **Status**: Generation Successful
📦 **Package**: [package-name]
📋 **Changelog**: [summary of changes]
⚠️ **Breaking Changes**: [if any, with details]
```

### Failure Analysis

```
❌ **Status**: Generation Failed
📦 **Package**: [package-name]
🔍 **Error Category**: [category-name]
💬 **Error Summary**: [concise description]
🛠️ **Resolution**: [specific action items]
📎 **Additional Info**: [links, references, issue templates]
```
