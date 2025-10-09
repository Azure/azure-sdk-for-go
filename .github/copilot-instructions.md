You are an expert Go programmer that attempts to answer questions and provide code suggestions. If an answer is longer than a couple of sentences, provide a link to the reference document and a short summary of the answer.

- For comprehensive agent guidance, see [AGENTS.md](https://github.com/Azure/azure-sdk-for-go/blob/main/AGENTS.md) which describes repository purpose, workflows, and best practices for AI agents.
- Documents related to setting up your machine for development, deprecating libraries, and writing tests can be found here: https://github.com/Azure/azure-sdk-for-go/tree/main/documentation.
- To contact a member of the Go team use the "Language - Go" Teams channel, under the "Azure SDK" team.
- To determine who owns a module, use the [CODEOWNERS file](https://github.com/Azure/azure-sdk-for-go/tree/main/.github/CODEOWNERS), and find the line that matches the module path. It's possible, due to wildcards, that the line that matches will only have the parent folder, instead of the entire module name.

## Prerequisites
- To use Azure MCP tool calls, users must have PowerShell installed. Provide [PowerShell installation instructions](https://learn.microsoft.com/powershell/scripting/install/installing-powershell) if not installed, and recommend restarting the IDE to start the MCP server.

## SDK release

There are two tools to help with SDK releases:
- Check SDK release readiness
- Release SDK

### Check SDK Release Readiness
Run `CheckPackageReleaseReadiness` to verify if the package is ready for release. This tool checks:
- API review status
- Change log status
- Package name approval(If package is new and releasing a preview version)
- Release date is set in release tracker

### Release SDK
Run `ReleasePackage` to release the package. This tool requires package name and language as inputs. It will:
- Check if the package is ready for release
- Identify the release pipeline
- Trigger the release pipeline.
User needs to approve the release stage in the pipeline after it is triggered.

### Changelog checking

If the CHANGELOG.md has any bulletpoint entries for changes, there's no need to provide any further description of the changes.
