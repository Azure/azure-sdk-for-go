# Go Mgmt SDK Release Guideline

Goal: Generate and release Azure SDK for Go package from Swagger or TypeSpec.

## Table of contents

* [Prerequisites](#Prerequisites)
* [Release Process](#Release-Process)
* [Update existing RP](#Update-existing-RP)
* [Releases new service](#Releases-new-service)
* [Breaking change review](#Breaking-change-review)

## Prerequisites

* Have a GitHub Account link to Microsoft/Azure organization, follow the link: https://dev.azure.com/azure-sdk/internal/_wiki/wikis/internal.wiki/111/Linking-Your-GitHub-Account
* (Optional) Update the account profile with the MS email and real name.
* Get the necessary permissions
  - Get the SDK repo permissions from: https://coreidentity.microsoft.com/manage/Entitlement/entitlement/azuresdkpart-heqj
  - Get the Asset repo permission by request join: 
    - (For vendors)https://github.com/orgs/Azure/teams/azure-sdk-write-vendors
  - Get the Azure subscription permission:
    - (For vendor) Request for the TME subscription or Ask the MSFTEs to add you in.
* Clone the following GitHub repositories:
  - https://github.com/Azure/azure-rest-api-specs
  - https://github.com/Azure/azure-sdk-for-go 
* Install the necessary tools by:
  - Generator tool, run: `go get github.com/Azure/azure-sdk-for-go/tools/generator`
  - Test-proxy, check: https://github.com/Azure/azure-sdk-tools/blob/main/tools/test-proxy/Azure.Sdk.Tools.TestProxy/README.md#installation
  - Tsp-client, run: `npm install -g @azure-tools/typespec-client-generator-cli`

## Release Process

* Release request:
    - Release requests are sent by the service team, telling us what packages we should release each month.
    - The release request is created as an issue at: https://github.com/Azure/sdk-release-request/issues
    - Usually, we should handle the release request meets with the following rules:
        - Issues have label: `Go`. And don’t have label: `HoldOn`
        - Issues were created from the second Friday of last month to the second Friday of the current month.
        - We could get the target issues by query like: `is:issue state:open label:Go -label:HoldOn created:2024-09-27..2024-10-15`
 - `ETA:All the release requests must be finished before the fourth Fridays of the current month.`

## Update existing RP 
Normally we deal with the RP which already have released packages before, for this kind of release request, please follow the steps:
1.	If the issue doesn’t have the `PRready` label. We’ll need to manually create a PR for it. Usually, this kind of issues will also have `Inconsistent` tag, which means we’ll need to update the version tag in autorest.md. Go to Step 2 
If the issue has the `PRready` label, go to Step 4
2.	We’ll generate the SDK in our local machine and create a PR for it.
    - a)	If the issue has the `TypeSpec`label, run following generate command:
generator release-v2 <azure-sdk-for-go path> <azure-rest-api-specs path> <RP name> <package name> --tsp-config “<path to tspconfig.yaml>”
      - azure-sdk-for-go path: Local path to your sdk repo folder, e.g. D:\Azure\azure-sdk-for-go
      - azure-rest-api-specs path: Local path to your sdk repo folder, e.g. D:\Azure\azure-rest-api-specs
      - RP name: service name, should be the same as one in the spec repo, e.g. apimanagement
      - package name: release package name, e.g. armapimanagement 
      - path to tspconfig.yaml: Path to the TSP config file of target RP in at local. e.g. D:\Azure\azure -rest-api-specs\specification\azurefleet\AzureFleet.Management\tspconfig.yaml
    - b)	If the issue doesn’t have the `TypeSpec`label, update the commit id and the tag in the autorest.md with the one in the release request, like:
        ```go
            azure-arm: true
            require:
                - https://github.com/Azure/azure-rest-api-specs/blob/e60df62e9e0d88462e6abba81a76d94eab000f0d/specification/containerinstance/resource-manager/readme.md
                - https://github.com/Azure/azure-rest-api-specs/blob/e60df62e9e0d88462e6abba81a76d94eab000f0d/specification/containerinstance/resource-manager/readme.go.md
            license-header: MICROSOFT_MIT_NO_VERSION
            module-version: 2.4.0
            tag: package-2023-05
            https://github.com/Azure/sdk-release-request/issues/5569
        ```
 
        Check the link’s commit id and the new tag, the autorest.md should be changed to:
        ```go
            azure-arm: true
            require:
                - https://github.com/Azure/azure-rest-api-specs/blob/ 655f4c80528b2aa2d5e52767e9a1bf7dd2a0655a/specification/containerinstance/resource-manager/readme.md
                - https://github.com/Azure/azure-rest-api-specs/blob/655f4c80528b2aa2d5e52767e9a1bf7dd2a0655a /specification/containerinstance/resource-manager/readme.go.md
            license-header: MICROSOFT_MIT_NO_VERSION
            module-version: 2.4.0
            tag: package-preview-2024-05
        ```

        Then, run following generate command:
            - generator release-v2 <azure-sdk-for-go path> <azure-rest-api-specs path> <RP name> <package name> --update-spec-version=false
3.	Fix the generated code if there are any issues.
4.	If the project contains files which end with '_live_test', we need to update the test records.
    - a.	If there are only API version changes, we could do the following steps:
        - i.	RUN: `test-proxy restore --assets-json-path assets.json` in the target project directory.
        - ii.	Get the new API version by finding any file in the project that ends with '_client', and in the comments, locate 'Generated from API version {api version}'.
        - iii.	Go to the '/.assets' folder and find the issued recording JSON file. Replace all the old version with the new API version found in step 2.
    - b.	If we need to re-record the tests, we could do the following steps:
        - i.	(GO Test  Chenjie fill)
5.	Run `test-proxy push --assets-json-path assets.json` in the project directory.
6.	Commit changes and create the PR.
7.	Check the PR pipeline and handle the errors:
    - a.	When a broken link error appears in the CI of a new major version or a new project's PR, leave a comment in PR: '/check-enforcer override'.
    - b.	When the CI gets stuck, try closing the PR and reopening it.
8.	Wait for the PR to get approved and merged.
9.	Search for the corresponding service’s pipeline in [Pipelines](https://dev.azure.com/azure-sdk/internal/_build?definitionScope=%5Cgo), wait for the automatic merge to complete, there will be a run in the pipeline. If not, manually click "Run Pipeline" to start a new run.
10.	When the Pipelines reports a broken link error, click and index the broken link and then start a new pipeline.
11.	After the pipeline completes, go to "Stages", select the rightmost node, and click "Approve".
12.	Leave a comment on the corresponding issue in azure-sdk-request repo and close the issue (e.g.<https://github.com/Azure/sdk-release-request/issues/5369#issuecomment-2301398185>)

## Releases new service 

If the service hasn’t released a package before, the service team would like to release a first beta version or first GA, please follow the below steps:
1.	Pull the latest main form:
    - a.	https://github.com/Azure/azure-rest-api-specs
    - b.	https://github.com/Azure/azure-sdk-for-go 
2.	Create a folder under azure-sdk-for-go\sdk\resourcemanager with the same name in the spec repo. Then create another folder under it, with the name armxxxxx, for instance armmonitor.
3.	Copy the template from the aztemplate directory[https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/template/aztemplate]. Be sure to update the contents as required, replacing all occurrences of template/aztemplate with the correct values.
4.	Run the generator 
5.	(Chenjie Fill)
6.	
7.	Check the PR pipeline and handle the errors:
    - a.	When a broken link error appears in the CI of a new major version or a new project's PR, leave a comment in PR: '/check-enforcer override'.
    - b.	When PR is a new project and the CI check is minimal, leave the comment '/azp run prepare-pipelines' to create new pipelines.
    - c.	When the CI gets stuck, try closing the PR and reopening it.
8.	Wait for the PR to get approved and merged.
9.	Search for the corresponding service’s pipeline in [Pipelines](https://dev.azure.com/azure-sdk/internal/_build?definitionScope=%5Cgo), wait for the automatic merge to complete, there will be a run in the pipeline. If not, manually click "Run Pipeline" to start a new run.
10.	When the Pipelines reports a broken link error, click and index the broken link and then start a new pipeline.
11.	After the pipeline completes, go to "Stages", select the rightmost node, and click "Approve".
12.	Leave a comment on the corresponding issue in azure-sdk-request repo and close the issue (e.g.<https://github.com/Azure/sdk-release-request/issues/5369#issuecomment-2301398185>)

## Breaking change review

We’ll need to check if the breaking changes in the PR are expected.
•	Where to find the target
1.	Check the Issues under sdk-release-request that have both "Breaking Change" and "Stable" labels. (or new major version)
2.	If there is a PR under the issue, go to the PR. Try finding the approved PR in the PR comments section.
Otherwise, go to the corresponding path in the azure-rest-api-specs Code and find the readme.md (e.g., azure-rest-api-specs/specification/{corresponding folder}/resource-manager/readme.md). Under #basic information, find the tag (e.g., package-2024-07), and go to the corresponding folder (e.g., azure-rest-api-specs/specification/{corresponding folder}/resource-manager/Microsoft.Batch/stable/2018-07-01). Click on History on the right, find the Commit where the Breaking Change occurred, and click on the number next to the Commit to enter the corresponding PR.Refine the the CHANGELOG.md of the corresponding azure-sdk-for-go if needed. There are two ways to determine if the breaking change was introduced by this PR.(changes related to "resource" could be ignored)


