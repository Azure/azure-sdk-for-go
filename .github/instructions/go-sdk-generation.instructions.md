# Generate Go SDK from API specification

Follow these steps to help users generate Go SDK from specific API specifications. Show the high-level description to users to help them understand the flow. Use the provided tools to perform actions and gather information as needed.

## Overview

1. **Prerequisites**: Ensure user has the correct environment to generate the Go SDK.
2. **Prepare generation config**: Get the information of API specification from user's input and prepare the generation config file.
3. **Run generation command**: Execute command to perform the Go SDK generation.
4. **Gather results and handle errors**: Retrieve the generation log and handle any problems.

## Prerequisites

Execute powershell [script](../../eng/scripts/Check-SDKGenerationPrerequisites.ps1) to verify all prerequisites are met.

## Step 1: Prepare Generation Config

**Goal**: Generate the `generatedInput.json` file.

**Actions**:

Execute the PowerShell [script](../../eng/scripts/Prepare-GenerationConfig.ps1) to prepare the generation config. The script supports multiple input methods:

- **From open editor files**: If `tspconfig.yaml` or `main.tsp` files are currently open in the editor, use their parent directory as the input:

  ```powershell
  .\eng\scripts\Prepare-GenerationConfig.ps1 -InputPath "<path from open editor files>"
  ```

- **From local file path**: If you have a local `tspconfig.yaml` file path:

  ```powershell
  .\eng\scripts\Prepare-GenerationConfig.ps1 -InputPath <local file path>"
  ```

- **From GitHub PR link**: If you have a GitHub PR URL:

  1. Prompt user: "Do you have a local copy of the spec repository? If yes, please provide the local repository path."
  2. Use the appropriate command based on the response:

  ```powershell
  # With local repository
  .\eng\scripts\Prepare-GenerationConfig.ps1 -PrUrl "<pr link>" -LocalRepoPath "<local repository path>"

  # Without local repository (will clone automatically)
  .\eng\scripts\Prepare-GenerationConfig.ps1 -PrUrl "<pr link>"
  ```

**Success Criteria**: `generatedInput.json` config is successfully generated in the workspace root folder.

## Step 2: Run Generation Command

**Goal**: Generate Go SDK according to the config.

**Actions**:

1. Ensure current workspace's git status is clean and up-to-date with the remote main branch.
2. Run generation command in the root folder of current workspace:

```bash
generator automation-v2 "<absolute path of generatedInput.json from last step>" generateOutput.json
```

3. Wait for the result. This step may take several minutes to complete.
4. Display the content of the `generateOutput.json` file from the workspace root folder.

**Success Criteria**: Generation command completes execution.

## Step 3: Gather Results and Handle Errors

**Goal**: Gather the results of Go SDK generation and provide suggestions for handling errors if generation fails.

**Actions**:

1. **Success Detection**: If the log contains `Finish processing typespec project` keyword and no `[ERROR]` keywords are found, then generation succeeded. Check the output and summarize the results.

2. **Error Handling**: If any `[ERROR]` logs or stack traces are found, follow the [Azure Go SDK Automation Troubleshooting Guide](../../documentation/sdk-automation-tsg.md) to suggest how to resolve the errors.

**Success Criteria**: Successfully detect the success or failure of Go SDK generation and provide a summary to the user.
