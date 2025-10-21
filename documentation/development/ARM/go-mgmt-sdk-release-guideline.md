# Go Management SDK Release Guideline

## Goal:
Generate and release Azure SDK for Go package from Swagger or TypeSpec.

---

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Release Process](#release-process)
4. [Special Case in PR](#special-case-in-pr)
5. [Add Live Test and Record](#add-live-test-and-record)
6. [Breaking Change Review](#breaking-change-review)

---

## Prerequisites

### 1. GitHub Account
- Link to Microsoft/Azure organization: [Link GitHub Account](https://dev.azure.com/azure-sdk/internal/_wiki/wikis/internal.wiki/111/Linking-Your-GitHub-Account)
- Optionally, update your GitHub profile with your Microsoft email and real name.

### 2. Permissions
- **SDK Repo**: Request SDK permissions [here](https://coreidentity.microsoft.com/manage/Entitlement/entitlement/azuresdkpart-heqj).
- **Asset Repo**: For vendors, join the Azure SDK team via [this link](https://github.com/orgs/Azure/teams/azure-sdk-write-vendors).
- **Azure Subscription**: Vendors must request access to the TME subscription or ask MSFT employees to add them.

### 3. Repositories to Clone
- [azure-rest-api-specs](https://github.com/Azure/azure-rest-api-specs)
- [azure-sdk-for-go](https://github.com/Azure/azure-sdk-for-go)

### 4. Install Necessary Tools
- **Generator Tool**: Run `go get github.com/Azure/azure-sdk-for-go/tools/generator`. If this fails, run `cd eng/tools/generator && go install` (Reference [link](https://github.com/Azure/azure-sdk-for-go/blob/72fe46870cff900262d54be73aa9a1eccfde12f2/documentation/code-generation.md#L46)).
- **Test Proxy**: Install by following [this guide](https://github.com/Azure/azure-sdk-tools/blob/main/tools/test-proxy/Azure.Sdk.Tools.TestProxy/README.md#installation).
- **TSP Client**: Run `npm install -g @azure-tools/typespec-client-generator-cli` (Ensure Node v20.18.0 is installed).

---

## Release Process

### 1. Release Request
- Release requests are initiated by the service team and should be submitted monthly.
- A request is created as an issue at: Azure's private `sdk-release-request` repo
- To process the request, use the following query: 
`is:issue state:open label:Go -label:HoldOn created:2024-09-27..2024-10-15`

### 2. General Release Guidelines
- All requests should be completed before the 4th Friday of the month.
- The issue may be labeled "ReadyForApiTest," or the live test about the service is ok indicating readiness.It appears in specs repo PR,[like](Azure/azure-rest-api-specs#33035)
- Ensure the "API Availability Check" passes for all services before releasing.

### 3. Update Existing RP

#### Steps:
1. If the issue lacks the `PRready` label, manually create a PR.
 - If the issue has the `Inconsistent` tag, update the version tag in `autorest.md`.
 - If the issue has the `TypeSpec` label, generate the SDK code using `generator release-v2`:
   ```
   generator release-v2 <path-to-sdk> <path-to-specs> <RP-name> <package-name> <tspconfig.yaml>
   ```
 - If no `TypeSpec` label, update the commit ID and tag in `autorest.md` and then:
   ```
   generator release-v2 <path-to-sdk> <path-to-specs> --update-spec-version=false
   ```
2. After generating the code, fix any issues, and if live tests are included:
 - Run `test-proxy restore` to restore assets.
 - Update API version in test records and push changes.

3. Commit changes and create a PR.
 - Handle PR errors (e.g., broken links) by commenting `/check-enforcer override`.
 - Wait for approval and automatic merging in Pipelines.

4. After pipeline completion, approve the merge and close the corresponding issue.

### 4. Releases New Service
#### Steps for New Service:
1. Pull the latest changes from:
 - [azure-rest-api-specs](https://github.com/Azure/azure-rest-api-specs)
 - [azure-sdk-for-go](https://github.com/Azure/azure-sdk-for-go)
 
2. Create a folder structure under `azure-sdk-for-go/sdks/resourcemanager` that matches the spec repo.
 - Copy the template from [aztemplate](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/template/aztemplate).

3. Run the generator and ensure the service name is correct in `README.md`.

4. Update the pipeline configuration in `ci.yml` and add `UsePipelineProxy: false` if you add live test.

5. Handle PR errors and ensure the pipeline runs as expected.
 - For new projects with minimal CI checks, use `/azp run prepare-pipelines`.

---


## Special Case in PR

### Handling Multiple Versions:
- Release services in ascending order of versions (e.g., FirstGA before FirstBeta).
- For services with multiple versions, the generated code will not change; only `changelog.md` should be updated.

### Different Swagger and Go Package Names:
- If the Swagger name differs from the Go package name, use the following command for code generation:
`generator release-v2 <sdk-path> <specs-path> resources <package-name>`

---

## Add Live Test and Record
It is not a required step in release process, we can ensure whether the api version has been deployed by live test.
### Install Scripts:
1. Run the following PowerShell commands to download necessary scripts:
 ```ps
 Invoke-WebRequest -OutFile "generate-assets-json.ps1" https://raw.githubusercontent.com/Azure/azure-sdk-tools/main/eng/common/testproxy/onboarding/generate-assets-json.ps1
 Invoke-WebRequest -OutFile "common-asset-functions.ps1" https://raw.githubusercontent.com/Azure/azure-sdk-tools/main/eng/common/testproxy/onboarding/common-asset-functions.ps1
```
2. Run the script in the service path:
`.\generate-assets-json.ps1 -InitialPush`
This will create a config file and push recordings to the Azure SDK Assets repo.
3. Before testing, create a utils_test.go file as the entry point for live tests. Modify "package" and pathToPackage to match your service.
4. Set the test mode to "live" using:
`$ENV:AZURE_RECORD_MODE="live"`
5. Once tests pass, switch to "playback" mode and ensure all tests pass in both modes.
6. Push the final assets with test-proxy push --assets-json-path assets.json.

---

## Breaking Change Review

### Steps:
1. Check for breaking changes by reviewing the Issues under sdk-release-request with both "Breaking Change" and "Stable" labels.
1. If there’s an associated PR, check for approval in the comments.
1. If no PR exists, locate the readme.md for the corresponding service and check the commit history for breaking changes.
1. [example to see details](https://github.com/Azure/azure-sdk-for-go/pull/23343)
