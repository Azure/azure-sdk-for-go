## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file:
- https://github.com/Azure/azure-rest-api-specs/blob/ecdce42924ed0f7e60a32c74bc0eb674ca6d4aae/specification/keyvault/data-plane/Microsoft.KeyVault/stable/7.2/secrets.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: true
output-folder: internal
tag: package-7.2
credential-scope: none
use: "@autorest/go@4.0.0-preview.36"
module-version: 0.6.0
export-clients: true
```
