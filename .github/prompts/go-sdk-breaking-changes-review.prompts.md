# Go SDK Breaking Changes Review

Follow these steps to help users review the Go SDK breaking changes. Show the high-level description to users to help them understand the flow. Use the provided tools to perform actions and gather information as needed.

## Step 1: Generate the SDK Locally

Follow [Generate Go SDK from API specification](./go-sdk-generation.instructions.md) to generate the SDK from the API specification if the breaking change info is not provided. Ensure the breaking changes, SDK package path and spec folder path are provided.

## Step 2: Review Breaking Changes

Follow the [Azure Go SDK Breaking Changes Review and Resolution Guide](../../documentation/sdk-breaking-changes-guide.md) with correct generation type section to review the breaking changes:

- Iterate through all patterns in the guide to find matching changelog entries.
  - If breaking changes match a pattern, extract the reason, spec pattern, impact, and resolution from the guide.
  - If the breaking changes can be resolved, update the spec file according to the guide. Then regenerate the SDK locally and verify if the breaking changes are resolved.
- For remaining breaking changes that do not match any patterns, provide the breaking change items to the user for manual review.

## Step 3: Gather Review Result

Provide comprehensive breaking changes analysis report in a new markdown file including: the breaking changes detected, their reason and impact, and any resolutions applied. Use structured format to make it easy to read and understand.
