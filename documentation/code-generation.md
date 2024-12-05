# Generate code

## Generate SDK packages

### Generate an Azure-SDK-for-Go service package

1. [Install AutoRest](https://github.com/Azure/autorest#installing-autorest).

1. Call autorest with the following arguments...

``` cmd
autorest path/to/readme/file --go --go-sdk-folder=<your/gopath/src/github.com/Azure/azure-sdk-for-go> --package-version=<version> --user-agent=<Azure-SDK-For-Go/version services> [--tag=choose/a/tag/in/the/readme/file]
```

For example...

``` cmd
autorest C:/azure-rest-api-specs/specification/advisor/resource-manager/readme.md --go --go-sdk-folder=C:/goWorkspace/src/github.com/Azure/azure-sdk-for-go --tag=package-2016-07-preview --package-version=v11.2.0-beta --user-agent='Azure-SDK-For-Go/v11.2.0-beta services'
```

- If you are looking to generate code based on a specific swagger file, you can replace `path/to/readme/file` with `--input-file=path/to/swagger/file`.
- If the readme file you want to use as input does not have golang tags yet, you can call autorest like this...

``` cmd
autorest path/to/readme/file --go --license-header=<MICROSOFT_APACHE_NO_VERSION> --namespace=<packageName> --output-folder=<your/gopath/src/github.com/Azure/azure-sdk-for-go/services/serviceName/mgmt/APIversion/packageName> --package-version=<version> --user-agent=<Azure-SDK-For-Go/version services> --clear-output-folder --can-clear-output-folder --tag=<choose/a/tag/in/the/readme/file>
```

For example...

``` cmd
autorest --input-file=https://raw.githubusercontent.com/Azure/azure-rest-api-specs/current/specification/network/resource-manager/Microsoft.Network/2017-10-01/loadBalancer.json --go --license-header=MICROSOFT_APACHE_NO_VERSION --namespace=lb --output-folder=C:/goWorkspace/src/github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network/lb --package-version=v11.2.0-beta --clear-output-folder --can-clear-output-folder
```

1. Run `go fmt` on the generated package folder.

1. To make sure the SDK has been generated correctly, also run `golint`, `go build` and `go vet`.

### Generate Azure SDK for Go service packages in bulk

All services, all API versions.

1. [Install AutoRest](https://github.com/Azure/autorest#installing-autorest).

This repo contains a tool to generate the SDK, which depends on the golang tags from the readme files in the Azure REST API specs repo. The tool assumes you have an [Azure REST API specs](https://github.com/Azure/azure-rest-api-specs) clone, and [golint](https://github.com/golang/lint) is installed.

1. `cd eng/tools/generator`

1. `go install`

1. Add `GOPATH/bin` to your `PATH`, in case it was not already there.

1. Call the generator tool like this...

``` cmd
generator –r [–v] [–l=logs/output/folder] –version=<version> path/to/your/swagger/repo/clone
```

For example...

``` cmd
generator –r –v –l=temp –version=v11.2.0-beta C:/azure-rest-api-specs
```

The generator tool already runs `go fmt`, `golint`, `go build` and `go vet`; so running them is not necessary.

#### Use the generator tool to generate a single package

1. Just call the generator tool specifying the service to be generated in the input folder.

``` cmd
generator –r [–v] [–l=logs/output/folder] –version=<version> path/to/your/swagger/repo/clone/specification/service
```

For example...

``` cmd
generator –r –v –l=temp –version=v11.2.0-beta C:/azure-rest-api-specs/specification/network
```

## Include a new package in the SDK

1. Submit a pull request to the Azure REST API specs repo adding the golang tags for the service and API versions in the service readme file, if the needed tags are not there yet.

1. Once the tags are available in the Azure REST API specs repo, generate the SDK.

1. In the changelog file, document the new generated SDK. Include the [autorest.go extension](https://github.com/Azure/autorest.go) version used, and the Azure REST API specs repo commit from where the SDK was generated.

1. Install [dep](https://github.com/golang/dep).

1. Run `dep ensure`.

1. Submit a pull request to this repo, and we will review it.

## Generate Azure SDK for Go profiles

Take a look into the [profile generator documentation](https://github.com/Azure/azure-sdk-for-go/tree/main/eng/tools/profileBuilder)

## Generate Go Mgmt SDK from Typespec
### Prerequisites
- [Node.js 18.x LTS](https://nodejs.org/en/download) or later
- [Go 1.23.x](https://go.dev/doc/install) or later
- [Git](https://git-scm.com/downloads)
- Identify the `tspconfig.yaml` file for your package in the [Rest API Spec Repo](https://github.com/Azure/azure-rest-api-specs) and ensure there is a configuration for the Go SDK similar to that shown below and in [this example](https://github.com/Azure/azure-rest-api-specs/blob/b09c9ec927456021dc549e111fa2cac3b4b00659/specification/contosowidgetmanager/Contoso.Management/tspconfig.yaml#L40)

     ```yaml
     options:
       "@azure-tools/typespec-go":
         service-dir: "SERVICE_DIRECTORY_NAME"
         package-dir: "PACKAGE_DIRECTORY_NAME"
         module: "github.com/Azure/azure-sdk-for-go/{service-dir}/{package-dir}"
         fix-const-stuttering: true
         flavor: "azure"
         generate-examples: true
         generate-fakes: true
         head-as-boolean: true
         inject-spans: true
         remove-unreferenced-types: true
    ```
- Local Clone of Rest API Spec Repo Fork
  - If you don't already have a fork, [Fork](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo#forking-a-repository) the [Rest API Spec Repo](https://github.com/Azure/azure-rest-api-specs).
  - Clone your fork of the repo.

    ```
      git clone https://github.com/{YOUR_GITHUB_USERNAME}/azure-rest-api-specs.git
    ```
- Local Clone of Go Language Repo Fork
  - If you don't already have a fork, [Fork](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo#forking-a-repository) the [Go Repo](https://github.com/Azure/azure-sdk-for-go).
  - Clone your fork of the repo.

    ```
      git clone https://github.com/{YOUR_GITHUB_USERNAME}/azure-sdk-for-go.git
    ```

### Steps
1. Complete the prerequisites listed above
2. Run [automation_init.sh](https://github.com/Azure/azure-sdk-for-go/blob/main/eng/scripts/automation_init.sh)
 
   ```sh
     cd "~/azure-sdk-for-go/eng/scripts" # navigate to the script directory
     ./automation_init.sh
   ```
   > On windows you can use Git Bash to run the script. This script also installs the generator as a global tool.
3. Create a local json file named `generatedInput.json` with content similar to that shown below

   ```json
      {
        "dryRun": false,
        "specFolder": "LOCAL_AZURE-REST-API-SPECS_REPO_ROOT", // e.g "C:\git\azure-sdk-for-go"
        "headSha": "SHA_OF_AZURE-REST-API-SPECS_REPO", // use ' git rev-parse HEAD '
        "repoHttpsUrl": "https://github.com/Azure/azure-rest-api-specs",
        "relatedTypeSpecProjectFolder": [
          "specification/SERVICE_DIRECTORY_NAME/PACKAGE_DIRECTORY_NAME/" // e.g specification/contosowidgetmanager/Contoso.Management
        ]
      }
   ```
4. Run the [Generator](https://github.com/chidozieononiwu/azure-sdk-for-go/tree/main/eng/tools/generator)
   ```sh
     generator automation-v2 "PATH_TO_generatedInput.json" generateOutput.json
   ```
   > generateOutput.json is a parameter for the name of the output file that will be created by the script.
   
5. View information about the generated SDK in `generateOutput.json`
6. Prepare your SDK for release. The necessary approvals, guidance for testing, documentation, and release pipelines is described in your release plan. More information about the Azure SDK Release Tool is [here](https://eng.ms/docs/products/azure-developer-experience/plan/release-plan)










