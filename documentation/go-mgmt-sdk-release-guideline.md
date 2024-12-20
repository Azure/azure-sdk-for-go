# Go Mgmt SDK Release Guideline

Goal: Generate and release Azure SDK for Go package from Swagger or TypeSpec.

## Table of contents

* [Prerequisites](#Prerequisites)
* [Release Process](#Release-Process)
* [Update existing RP](#Update-existing-RP)
* [Special case in PR](#Special-case-in-PR)
* [Releases new service](#Releases-new-service)
* [Add live test and record](#Add-live-test-and-record)
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
  - Generator tool, run: `go get github.com/Azure/azure-sdk-for-go/tools/generator`，if this command does not work,  run：`cd eng/tools/generator && go install`,referenced from  (https://github.com/Azure/azure-sdk-for-go/blob/72fe46870cff900262d54be73aa9a1eccfde12f2/documentation/code-generation.md#L46)
  - Test-proxy, check: https://github.com/Azure/azure-sdk-tools/blob/main/tools/test-proxy/Azure.Sdk.Tools.TestProxy/README.md#installation
  - Tsp-client, run: `npm install -g @azure-tools/typespec-client-generator-cli`
  NOTE：node version: v20.18.0，download url: https://nodejs.org/en/download/prebuilt-installer

## Release Process

* Release request:
    - Release requests are sent by the service team, telling us what packages we should release each month.
    - The release request is created as an issue at: https://github.com/Azure/sdk-release-request/issues
    - Usually, we should handle the release request meets with the following rules:
        - Issues have label: `Go`. And don’t have label: `HoldOn`
        - Issues were created from the second Friday of last month to the second Friday of the current month.
        - We could get the target issues by query like: `is:issue state:open label:Go -label:HoldOn created:2024-09-27..2024-10-15`
 - `ETA:All the release requests must be finished before the fourth Fridays of the current month.`
 - `ETA:All the release requests must be labelized with "ReadyForApiTest"`, that means that the service to release is ready

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
      - before running the command to generate code, you should make sure that there is `readme.go.md` under path `specification/xxxx/resource-manager` or there is node `"@azure-tools/typespec-go"` config in `tspconfig.yaml` file,[referenced from](https://github.com/Azure/azure-rest-api-specs/blob/main/specification/liftrneon/Neon.Postgres.Management/tspconfig.yaml) 
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
 
        update the link’s commit id and the  tag in the file  autorest.md , tag should be changed to the tag mentioned in `sdk-release-request/issues/xxx` which showed as `"Readme Tag"` [to see detais](https://github.com/Azure/sdk-release-request/issues/5759) , commit id should be change to the lastes commit id in the `"link"`
    
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
6.	Check that whether generated  api version  is consistent with the version commented in client,[api version check see details](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/resourcemanager/neonpostgres/armneonpostgres/operations_client.go#L38),[client version see details](https://github.com/Azure/azure-rest-api-specs/blob/main/specification/liftrneon/Neon.Postgres.Management/main.tsp#L28),if is inconsistent, reported that
7.	Commit changes and create the PR.
8.	Check the PR pipeline and handle the errors:
    - a.	When a broken link error appears in the CI of a new major version or a new project's PR, leave a comment in PR: '/check-enforcer override'.
    - b.	When the CI gets stuck, try closing the PR and reopening it.
9.	Wait for the PR to get approved and merged.
10.	Search for the corresponding service’s pipeline in [Pipelines](https://dev.azure.com/azure-sdk/internal/_build?definitionScope=%5Cgo), wait for the automatic merge to complete, there will be a run in the pipeline. If not, manually click "Run Pipeline" to start a new run.
11.	When the Pipelines reports a broken link error, click and index the broken link and then start a new pipeline.
12.	After the pipeline completes, go to "Review", select the rightmost node, and click "Approve".
13.	Leave a comment on the corresponding issue in azure-sdk-request repo and close the issue (e.g.<https://github.com/Azure/sdk-release-request/issues/5369#issuecomment-2301398185>)

## Special case in PR
1.	If a service has multiple version to release at same window process, you need to release the service order by the version in ascending order.that means the small version needs to be released first, and then the next version should be generated auomatically again depending the released code.
2.	If the current version is `FirstGA` and the last version is `FirstBeta`,after generating,no codes will be changed, only the file `changelog.md` will be modified


## Releases new service 

If the service hasn’t released a package before, the service team would like to release a first beta version or first GA, please follow the below steps:
1.	Pull the latest main form:
    - a.	https://github.com/Azure/azure-rest-api-specs
    - b.	https://github.com/Azure/azure-sdk-for-go 
2.	Create a folder under azure-sdk-for-go\sdk\resourcemanager with the same name in the spec repo. Then create another folder under it, with the name armxxxxx, for instance armmonitor.
3.	Copy the template from the aztemplate directory[https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/template/aztemplate]. Be sure to update the contents as required, replacing all occurrences of template/aztemplate with the correct values.
4.	Run the generator 
5.	service name check for readme,file path:sdk/resourcemanager/xxxxrservice/armxxxservice/README.md. if service name is made up multi words `test service`, the service name auto generated in readme.md file may be changed to `TestService`,we need to replace `TestService` to `Test Service`
6.	update file:sdk/resourcemanager/xxxxrservice/armxxxservice/ci.yml, add `UsePipelineProxy: false` under `extends->parameters` if not exist
7.	Check the PR pipeline and handle the errors:
    - a.	When a broken link error appears in the CI of a new major version or a new project's PR, leave a comment in PR: '/check-enforcer override'.
    - b.	When PR is a new project and the CI check is minimal, leave the comment '/azp run prepare-pipelines' to create new pipelines.
    - c.	When the CI gets stuck, try closing the PR and reopening it.
8.	Wait for the PR to get approved and merged.
9.	Search for the corresponding service’s pipeline in [Pipelines](https://dev.azure.com/azure-sdk/internal/_build?definitionScope=%5Cgo), wait for the automatic merge to complete, there will be a run in the pipeline. If not, manually click "Run Pipeline" to start a new run.
10.	When the Pipelines reports a broken link error, click and index the broken link and then start a new pipeline.
11.	After the pipeline completes, go to "Stages", select the rightmost node, and click "Approve".
12.	Leave a comment on the corresponding issue in azure-sdk-request repo and close the issue (e.g.<https://github.com/Azure/sdk-release-request/issues/5369#issuecomment-2301398185>)

## Add live test and record
1. Using following commands to install the script(Do not directly download from GitHub, may have execution issues):
    ```cmd
     1)`Invoke-WebRequest -OutFile "generate-assets-json.ps1" https://raw.githubusercontent.com/Azure/azure-sdk-tools/main/eng/common/testproxy/onboarding/generate-assets-json.ps1`
    Or
    `wget https://raw.githubusercontent.com/Azure/azure-sdk-tools/main/eng/common/testproxy/onboarding/generate-assets-json.ps1 -o generate-assets-json.ps1`
    
   2) `Invoke-WebRequest -OutFile "generate-assets-json.ps1" https://raw.githubusercontent.com/Azure/azure-sdk-tools/main/eng/common/testproxy/onboarding/common-asset-functions.ps1`
    Or
    `wget https://raw.githubusercontent.com/Azure/azure-sdk-tools/main/eng/common/testproxy/onboarding/generate-assets-json.ps1 -o common-asset-functions.ps1`
    ```
    then you downloaded two files: `common-asset-functions.ps1` and `generate-assets-json.ps1`
2. Run `Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass` to disable unsigned check
3. Run the script under the service path:`sdk/resourcemanager/xxx/armxxx`,you can copy these two files into this path,and run `.\generate-assets-json.ps1  -InitialPush`, and then the file `assets.json` will be auto generated f under the path .
 Explain(With - InitialPush parameter, the script will do):
    - Create a config file under service path.(Details see below)
    - Remove all the recording files inside the service.
    - Generate a unique ID and added to the config file, push the recording files to Azure SDK Assets repo, stored as a Tag with that ID
    - Without - InitialPush, the script will only create config file. We still need to manually call test-proxy push -a <path-to-assets.json>  to finish rest of the steps.
4. Before creating a live test file, you need to create file named `utils_test.go`,which is the entrance of live test,referenced from [`the example`](https://github.com/jliusan/azure-sdk-for-go/blob/sdk-release-guideline/sdk/resourcemanager/compute/armcompute/utils_test.go),change the `"package"` name and const `"pathToPackage"` to the current service
5. The you can create a live test file named like `_live_test.go`,referenced from [`the example`](https://github.com/jliusan/azure-sdk-for-go/blob/sdk-release-guideline/sdk/resourcemanager/compute/armcompute/virtualmachineextensionimage_live_test.go)
6. Before you run the test cases under the mode `"live"`, you need to set `$ENV:AZURE_RECORD_MODE="live"`,[to see details](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/developer_setup.md#write-tests)
7. If all the test cases pass, you still need to change the test mode to `"playback"`, and ensure that test cases pass too under the test mode `"playback"`,otherwise the service ci pipeline will get failure in a new PR
8. At last, you need to run `test-proxy push --assets-json-path assets.json` under the service path

## Breaking change review

We’ll need to check if the breaking changes in the PR are expected.
•	Where to find the target
1.	Check the Issues under sdk-release-request that have both "Breaking Change" and "Stable" labels. (or new major version)
2.	If there is a PR under the issue, go to the PR. Try finding the approved PR in the PR comments section.
Otherwise, go to the corresponding path in the azure-rest-api-specs Code and find the readme.md (e.g., azure-rest-api-specs/specification/{corresponding folder}/resource-manager/readme.md). Under #basic information, find the tag (e.g., package-2024-07), and go to the corresponding folder (e.g., azure-rest-api-specs/specification/{corresponding folder}/resource-manager/Microsoft.Batch/stable/2018-07-01). Click on History on the right, find the Commit where the Breaking Change occurred, and click on the number next to the Commit to enter the corresponding PR.Refine the the CHANGELOG.md of the corresponding azure-sdk-for-go if needed. There are two ways to determine if the breaking change was introduced by this PR.(changes related to "resource" could be ignored).
[example to see details](https://github.com/Azure/azure-sdk-for-go/pull/23343)