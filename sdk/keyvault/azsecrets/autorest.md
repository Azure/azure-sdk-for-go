## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file:
- https://github.com/Azure/azure-rest-api-specs/blob/e2ef44b87405b412403ccb005bfb3975411adf60/specification/keyvault/data-plane/Microsoft.KeyVault/stable/7.3/secrets.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: false
output-folder: internal/generated
openapi-type: "data-plane"
security: "AADToken"
security-scopes:  "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.41"
module-version: 0.8.0
export-clients: true

# export request creators and response handlers for use by pollers and pagers, and remove the secretVersion path param check
# (secretVersion == "" is legal for Key Vault but indescribable in OpenAPI)
directive:
  - from: keyvault_client.go
    where: $
    transform: >-
      return $.
        replace(/get(.*)CreateRequest/g, function(_, s) { return `Get${s}CreateRequest` }).
        replace(/get(.*)HandleResponse/g, function(_, s) { return `Get${s}HandleResponse` }).
        replace(/\sif secretVersion == "" \{\s+return nil, errors\.New\("parameter secretVersion cannot be empty"\)\s+\}\s/g, ``);
```
