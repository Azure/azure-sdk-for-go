## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file:
- https://github.com/Azure/azure-rest-api-specs/blob/e2ef44b87405b412403ccb005bfb3975411adf60/specification/keyvault/data-plane/Microsoft.KeyVault/stable/7.3/certificates.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: true
output-folder: internal/generated
openapi-type: "data-plane"
security: "AADToken"
security-scopes:  "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.40"
module-version: 0.4.0
export-clients: true

# remove the empty certificateVersion path param check.  it's legal for KV but can't be described in OpenAPI
directive:
  - from: constants.go
    where: $
    transform: >-
      return $.
        replace(/moduleName\s+=\s+"generated"/, `ModuleName = "azcertificates"`).
        replace(/moduleVersion\s+=/, `ModuleVersion =`);
  - from: keyvault_client.go
    where: $
    transform: >-
      return $.
        replaceAll(/\sif certificateVersion == "" \{\s+return nil, errors\.New\("parameter certificateVersion cannot be empty"\)\s+\}\s/g, ``);
```
