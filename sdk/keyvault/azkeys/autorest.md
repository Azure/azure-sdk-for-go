## Go

These settings apply only when `--go` is specified on the command line.

<!-- autorest --use=@autorest/go@4.0.0-preview.29 https://github.com/Azure/azure-rest-api-specs/tree/1e2c9f3ec93078da8078389941531359e274f32a/specification/keyvault/data-plane --tag=package-7.2 --output-folder=internal --module-azkeys --module-version=0.1.0 --openapi-type="data-plane" --security="AADToken" --security-scopes="https://vault.azure.net/.default" -->

``` yaml
go: true
version: "^3.0.0"
input-file:
# - https://github.com/Azure/azure-rest-api-specs/blob/1e2c9f3ec93078da8078389941531359e274f32a/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/backuprestore.json
# - https://github.com/Azure/azure-rest-api-specs/blob/1e2c9f3ec93078da8078389941531359e274f32a/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/certificates.json
- https://github.com/Azure/azure-rest-api-specs/blob/1e2c9f3ec93078da8078389941531359e274f32a/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/common.json
- https://github.com/Azure/azure-rest-api-specs/blob/1e2c9f3ec93078da8078389941531359e274f32a/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/keys.json
- https://github.com/Azure/azure-rest-api-specs/blob/1e2c9f3ec93078da8078389941531359e274f32a/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/rbac.json
# - https://github.com/Azure/azure-rest-api-specs/blob/1e2c9f3ec93078da8078389941531359e274f32a/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/secrets.json
- https://github.com/Azure/azure-rest-api-specs/blob/1e2c9f3ec93078da8078389941531359e274f32a/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/securitydomain.json
# - https://github.com/Azure/azure-rest-api-specs/blob/1e2c9f3ec93078da8078389941531359e274f32a/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.3-preview/storage.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: false
output-folder: internal
module: azkeys
openapi-type: "data-plane"
security: "AADToken"
security-scopes:  "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.29"
module-version: 0.1.0
```
