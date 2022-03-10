## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file: https://github.com/Azure/azure-rest-api-specs/blob/d744b6bcb95ab4034832ded556dbbe58f4287c5b/specification/cosmos-db/data-plane/Microsoft.Tables/preview/2019-02-02/table.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: false
output-folder: internal
tag: package-2019-02
credential-scope: none
use: "@autorest/go@4.0.0-preview.36"
module-version: 0.5.1
security: "AADToken"
security-scopes: "https://storage.azure.com/.default"
modelerfour:
  group-parameters: false
```

### Go multi-api

``` yaml $(go) && $(multiapi)
batch:
  - tag: package-2019-02
```