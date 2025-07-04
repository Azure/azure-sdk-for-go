# Generate Go SDK from API specification

Follow these steps to help users generate Go SDK from specific API specifications. Show the high-level description to users to help them understand the flow. Use the provided tools to perform actions and gather information as needed.

## Go SDK Generation Process

1. **Prerequisites**: Ensure user has the correct environment to generate the Go SDK.
2. **Prepare generation config**: Get the information of API specification from user's input and prepare the generation config file.
3. **Run generation command**: Execute command to perform the Go SDK generation.
4. **Gather results and handle errors**: Retrieve the generation log and handle any problems.

## Prerequisites

1. Ensure Go version is 1.23 or later.
2. Install the `generator` tool (v0.1.0) with the command:

```bash
go install github.com/Azure/azure-sdk-for-go/eng/tools/generator@v0.1.0
```

3. Ensure Node.js version is 20 or later.
4. Install `tsp-client` tool (v0.21.0) with the command:

```bash
npm install -g @azure-tools/typespec-client-generator-cli@v0.21.0
```

5. Ensure GitHub CLI tool is installed and successfully logged in.
6. Ensure Git is installed.

## Step 1: Prepare Generation Config

**Goal**: Generate the `generatedInput.json` file.

**Actions**:

1. Get API specification root path (`specPath`):

- **From open editor files**: If `tspconfig.yaml` or `main.tsp` files are currently open in the editor, use their parent directory as the `specPath`.
- **From local file path**: If the user provides a local `tspconfig.yaml` file path, use its parent directory as the `specPath`.
- **From GitHub PR link**: If user provides a GitHub PR URL:
  - Prompt user: "Do you have a local copy of the spec repository?"
  - **If yes**:
    - Use the local repository path
  - **If no**:
    - Clone the repository to system temp directory
  - Use GitHub CLI to checkout the PR branch: `gh pr checkout <PR_NUMBER>`
  - Analyze the PR's changed files to locate `tspconfig.yaml`
  - Use the parent directory of the found `tspconfig.yaml` as the `specPath`

2. Split the `specPath` into two parts:

   - The part before `/specification/` is the `specFolder`.
   - The part after the first part is the `projectFolder`.

3. Use `git rev-parse HEAD` under `specFolder` to get the `headSha` value.

4. Use `git remote -v` under `specFolder` to determine the `repoHttpsUrl` value to be either `https://github.com/Azure/azure-rest-api-specs` or `https://github.com/Azure/azure-rest-api-specs-pr`.

5. Generate or replace the `generatedInput.json` file in the root folder of current workspace with the following format:

```json
{
  "specFolder": "<specFolder>",
  "headSha": "<headSha>",
  "repoHttpsUrl": "<repoHttpsUrl>",
  "relatedTypeSpecProjectFolder": ["<projectFolder>"]
}
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
