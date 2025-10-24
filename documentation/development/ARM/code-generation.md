# Generate Go Mgmt SDK from Typespec

### Prerequisites
- [Node.js 18.x LTS](https://nodejs.org/en/download) or later
- [Go 1.23.x](https://go.dev/doc/install) or later
- [Git](https://git-scm.com/downloads)
- Identify the `tspconfig.yaml` file for your package in the [Rest API Spec Repo](https://github.com/Azure/azure-rest-api-specs) and ensure there is a configuration for the Go SDK similar to that shown below and in [this example](https://github.com/Azure/azure-rest-api-specs/blob/main/specification/contosowidgetmanager/Contoso.Management/tspconfig.yaml#L40)

     ```yaml
     options:
       "@azure-tools/typespec-go":
         service-dir: "SERVICE_DIRECTORY_NAME"
         package-dir: "PACKAGE_DIRECTORY_NAME" # Should start with 'arm' following namespace review result. e.g. https://github.com/Azure/azure-sdk/issues/8290
         module: "github.com/Azure/azure-sdk-for-go/{service-dir}/{package-dir}"
         fix-const-stuttering: true
         flavor: "azure"
         generate-examples: true
         generate-fakes: true
         head-as-boolean: true
         inject-spans: true
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
        "specFolder": "LOCAL_AZURE-REST-API-SPECS_REPO_ROOT", // e.g. "C:/git/azure-rest-api-specs"
        "headSha": "SHA_OF_AZURE-REST-API-SPECS_REPO", // use ' git rev-parse HEAD ' on the local azure-rest-api-specs repo root 
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