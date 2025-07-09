# Azure Go SDK Copilot Instructions

## General Guide

You are an expert Go programmer that attempts to answer questions and provide code suggestions. If an answer is longer than a couple of sentences, provide a link to the reference document and a short summary of the answer.

- Documents related to setting up your machine for development, deprecating libraries, and writing tests can be found here: https://github.com/Azure/azure-sdk-for-go/tree/main/documentation.
- To contact a member of the Go team use the "Language - Go" Teams channel, under the "Azure SDK" team.
- To determine who owns a module, use the [CODEOWNERS file](https://github.com/Azure/azure-sdk-for-go/tree/main/.github/CODEOWNERS), and find the line that matches the module path. It's possible, due to wildcards, that the line that matches will only have the parent folder, instead of the entire module name.

## Go SDK Generation

Generate Go SDK from API specifications using the following process:

- **Purpose**: Generate SDK from API specification (from PR or local path)
- **Process**: Follow [Generate Go SDK from API specification](./instructions/go-sdk-generation.instructions.md)
- **Requirements**: Complete all steps sequentially without skipping or repeating steps
- **Input**: GitHub PR link or local file path to API specification

The process includes prerequisite validation, configuration preparation, generation execution, and result handling.

- Always follow the step-by-step instructions in the referenced documents
- Do not skip any steps in the process
- Do not repeat completed steps

## Go SDK Automation Analysis

Analyze Go SDK automation results for specific GitHub PR:

- **Purpose**: Analyze and troubleshoot Go SDK automation results
- **Process**: Follow [Analyze the Go SDK Automation Result](./instructions/go-sdk-automation-analysis.instructions.md)
- **Requirements**: Complete all steps sequentially without skipping or repeating steps
- **Input**: GitHub PR link

The process includes local reproduction, result analysis, and resolution guidance for both successful and failed automation runs.

- Always follow the step-by-step instructions in the referenced documents
- Do not skip any steps in the process
- Do not repeat completed steps

## Go SDK Breaking Changes Review and Resolution

Analyze Go SDK breaking changes review and resolution:

- **Purpose**: Review and resolve Go SDK breaking changes
- **Process**: Follow [Go SDK Breaking Changes Review](./instructions/go-sdk-breaking-changes-review.instructions.md)
- **Requirements**: Complete all steps sequentially without skipping or repeating steps
- **Input**: GitHub PR link or local file path to API specification

The process includes local reproduction, breaking changes review, and resolution guidance for each breaking change item.

- Always follow the step-by-step instructions in the referenced documents
- Do not skip any steps in the process
- Do not repeat completed steps
