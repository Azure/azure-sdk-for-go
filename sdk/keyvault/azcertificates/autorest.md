## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file:
- https://github.com/Azure/azure-rest-api-specs/blob/ecdce42924ed0f7e60a32c74bc0eb674ca6d4aae/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/common.json
- https://github.com/Azure/azure-rest-api-specs/blob/ecdce42924ed0f7e60a32c74bc0eb674ca6d4aae/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/certificates.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: true
output-folder: internal/generated
module: azcertificates
openapi-type: "data-plane"
security: "AADToken"
security-scopes:  "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.35"
module-version: 0.1.0
export-clients: true
```
