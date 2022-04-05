## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file:
- https://github.com/Azure/azure-rest-api-specs/blob/e2ef44b87405b412403ccb005bfb3975411adf60/specification/keyvault/data-plane/Microsoft.KeyVault/stable/7.3/secrets.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: true
output-folder: internal
tag: package-7.2
credential-scope: none
use: "@autorest/go@4.0.0-preview.37"
module-version: 0.7.0
export-clients: true
```
