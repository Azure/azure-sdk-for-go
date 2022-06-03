## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file:
- https://github.com/Azure/azure-rest-api-specs/blob/e2ef44b87405b412403ccb005bfb3975411adf60/specification/keyvault/data-plane/Microsoft.KeyVault/stable/7.3/keys.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: false
output-folder: internal/generated
openapi-type: "data-plane"
security: "AADToken"
security-scopes:  "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.42"
export-clients: true

# export request creators and response handlers for use by pollers and pagers, and remove the keyVersion path param check
# (keyVersion == "" is legal for Key Vault but indescribable by OpenAPI)
directive:
  - from: keyvault_client.go
    where: $
    transform: >-
      return $.
        replace(/get(.*)CreateRequest/g, function(_, s) { return `Get${s}CreateRequest` }).
        replace(/get(.*)HandleResponse/g, function(_, s) { return `Get${s}HandleResponse` }).
        replace(/if keyVersion == "" \{\s*.*\s*\}\s*/g, ``).
        replace(/(urlPath = strings\.ReplaceAll\(urlPath, "\{key-version\}", url\.PathEscape\(keyVersion\)\)\s+)/g, function(_, s) { return `${s}urlPath = strings.ReplaceAll(urlPath, "//", "/")\n\t` });
```
