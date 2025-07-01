# Azure Go SDK Copilot Instructions

This document provides instructions for GitHub Copilot to assist with Go SDK generation and automation analysis for the Azure SDK for Go repository.

## Go SDK Generation

Generate Go SDK from API specifications using the following process:

- **Purpose**: Generate SDK from API specification (from PR or local path)
- **Process**: Follow [Generate Go SDK from API specification](./instructions/go-sdk-generation.instructions.md)
- **Requirements**: Complete all steps sequentially without skipping or repeating steps
- **Input Sources**:
  - GitHub PR containing API specification changes
  - Local file path to API specification

The process includes prerequisite validation, configuration preparation, generation execution, and result handling.

## Go SDK Automation Analysis

Analyze Go SDK automation results from GitHub PRs:

- **Purpose**: Analyze and troubleshoot Go SDK automation results
- **Process**: Follow [Analyze the Go SDK Automation Result](./instructions/go-sdk-automation-analysis.instructions.md)
- **Requirements**: Complete all steps sequentially without skipping or repeating steps
- **Input**: GitHub PR with automation results to analyze

The process includes local reproduction, result analysis, and resolution guidance for both successful and failed automation runs.

## General Guidelines

- Always follow the step-by-step instructions in the referenced documents
- Do not skip any steps in the process
- Do not repeat completed steps
- Provide clear status updates and error handling guidance
