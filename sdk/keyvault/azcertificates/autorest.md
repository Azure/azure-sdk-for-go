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
module: azcertificates
openapi-type: "data-plane"
security: "AADToken"
security-scopes:  "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.37"
module-version: 0.2.1
export-clients: true
```
