# Refresh SDK Task — Update Go SDKs with Latest TypeSpec

## Overview

Update the following Azure SDK for Go modules to their latest TypeSpec-generated code, ensuring each compiles, is properly formatted, and is submitted as an individual pull request.

## Target Services

| # | Service | Spec Directory |
|---|---------|----------------|
| 1 | alertsmanagement | `C:\w\azure-rest-api-specs\specification\alertsmanagement\resource-manager\Microsoft.AlertsManagement\AlertsManagement\tspconfig.yaml` |
| 2 | azurestackhci | `C:\w\azure-rest-api-specs\specification\azurestackhci\resource-manager\Microsoft.AzureStackHCI\StackHCI\tspconfig.yaml` |
| 3 | commerce | `C:\w\azure-rest-api-specs\specification\commerce\resource-manager\Microsoft.Commerce\Commerce\tspconfig.yaml` |
| 4 | communication | `C:\w\azure-rest-api-specs\specification\communication\Communication.Management\tspconfig.yaml` |
| 5 | containerservice | `C:\w\azure-rest-api-specs\specification\containerservice\resource-manager\Microsoft.ContainerService\aks\tspconfig.yaml` |
| 6 | eventgrid | `C:\w\azure-rest-api-specs\specification\eventgrid\resource-manager\Microsoft.EventGrid\EventGrid\tspconfig.yaml` |
| 7 | kubernetesconfigurationextensions | `C:\w\azure-rest-api-specs\specification\kubernetesconfiguration\resource-manager\Microsoft.KubernetesConfiguration\extensions\tspconfig.yaml` |
| 8 | kubernetesconfigurationextensiontypes | `C:\w\azure-rest-api-specs\specification\kubernetesconfiguration\resource-manager\Microsoft.KubernetesConfiguration\extensionTypes\tspconfig.yaml` |
| 9 | kubernetesconfigurationfluxconfigurations | `C:\w\azure-rest-api-specs\specification\kubernetesconfiguration\resource-manager\Microsoft.KubernetesConfiguration\fluxConfigurations\tspconfig.yaml` |
| 10 | kubernetesconfigurationprivatelinkscopes | `C:\w\azure-rest-api-specs\specification\kubernetesconfiguration\resource-manager\Microsoft.KubernetesConfiguration\privateLinkScopes\tspconfig.yaml` |
| 11 | maps | `C:\w\azure-rest-api-specs\specification\maps\resource-manager\Microsoft.Maps\Maps\tspconfig.yaml` |
| 12 | marketplace | `C:\w\azure-rest-api-specs\specification\marketplace\resource-manager\Microsoft.Marketplace\Marketplace\tspconfig.yaml` |
| 13 | purview | `C:\w\azure-rest-api-specs\specification\purview\resource-manager\Microsoft.Purview\Purview\tspconfig.yaml` |
| 14 | recoveryservicessiterecovery | `C:\w\azure-rest-api-specs\specification\recoveryservicessiterecovery\resource-manager\Microsoft.RecoveryServices\SiteRecovery\tspconfig.yaml` |
| 15 | search | `C:\w\azure-rest-api-specs\specification\search\resource-manager\Microsoft.Search\Search\tspconfig.yaml` |
| 16 | frontdoor | `C:\w\azure-rest-api-specs\specification\frontdoor\FrontDoor.Management\tspconfig.yaml` |
| 17 | networkfunction | `C:\w\azure-rest-api-specs\specification\networkfunction\NetworkFunction.Management\tspconfig.yaml` |

## Repositories

- **SDK repo (Go):** `C:/w/azure-sdk-for-go`
- **API specs repo:** `C:/w/azure-rest-api-specs`

## Workflow — Repeat for Each Service

For **each** service listed above, perform the following steps in order:

### Step 1: Generate SDK with Latest TypeSpec

- Locate to the root of the Azure SDK for Go repository:
  ```
  cd C:/w/azure-sdk-for-go
  ```
- Use the `generator generate` cmd ([doc](C:\w\azure-sdk-for-go\eng\tools\generator\README.md)) to regenerate the Go SDK from the latest TypeSpec spec.

### Step 2: Build & Verify

- Compile the generated code to make sure there are no build errors:
  ```
  cd sdk/resourcemanager/<service>/<package>
  go build ./...
  ```
- Run `gofmt` / `goimports` to ensure the code is properly formatted:
  ```
  gofmt -w .
  goimports -w .
  ```
- Fix any compilation or formatting issues before proceeding.

### Step 3: Create a New Branch

- From the `main` branch, create a dedicated branch for this service:
  ```
  git checkout -b refresh/<service>-latest-typespec
  ```

### Step 4: Commit the Changes

- Stage and commit all generated/modified files:
  ```
  git add .
  git commit -m "Regenerate <service> Go SDK with latest TypeSpec"
  ```

### Step 5: Push & Create a Pull Request

- Push the branch to the remote:
  ```
  git push origin refresh/<service>-latest-typespec
  ```
- Create a pull request targeting the `main` branch with:
  - **Title:** `[Refresh] Regenerate <service> Go SDK with latest TypeSpec`
  - **Description:** Automated regeneration of the `<service>` Go SDK using the latest TypeSpec definitions.

### Step 6: Return to Main

- Switch back to the `main` branch before starting the next service:
  ```
  git checkout main
  ```
- Repeat from Step 1 for the next service in the list.

## Important Notes

- Process each service **one at a time**, completing all steps before moving on.
- Each service gets its **own branch and its own PR** — do not bundle multiple services into a single PR.
- If a service has multiple TypeSpec projects (e.g., multiple `tspconfig.yaml` files), handle each one and include all changes in the same service branch/PR.
- If generation or compilation fails for a service, note the error, skip that service, and continue with the next one. Come back to failed services afterward.
- Always ensure you are on the `main` branch (with latest changes pulled) before starting a new service.
