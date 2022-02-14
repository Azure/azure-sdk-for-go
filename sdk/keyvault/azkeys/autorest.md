## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file:
- https://github.com/Azure/azure-rest-api-specs/blob/8a061f1e9031450b9eb5546d242f2a28c93eaa69/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/common.json
- https://github.com/Azure/azure-rest-api-specs/blob/8a061f1e9031450b9eb5546d242f2a28c93eaa69/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/keys.json
- https://github.com/Azure/azure-rest-api-specs/blob/8a061f1e9031450b9eb5546d242f2a28c93eaa69/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/rbac.json
- https://github.com/Azure/azure-rest-api-specs/blob/8a061f1e9031450b9eb5546d242f2a28c93eaa69/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/securitydomain.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: true
output-folder: internal/generated
module: azkeys
openapi-type: "data-plane"
security: "AADToken"
security-scopes:  "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.35"
module-version: 0.3.0
export-clients: true
```
