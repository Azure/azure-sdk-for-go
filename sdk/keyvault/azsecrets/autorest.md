## Go

These settings apply only when `--go` is specified on the command line.

<!-- Original autorest command used by Chris Scott -->
<!-- autorest --use=@autorest/go@4.0.0-preview.20 https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/specification/cosmos-db/data-plane/readme.md --tag=package-2019-02 --file-prefix="zz_generated_" --modelerfour.lenient-model-deduplication --license-header=MICROSOFT_MIT_NO_VERSION --output-folder=aztables --module=aztables --openapi-type="data-plane" --credential-scope=none -->

``` yaml
go: true
version: "^3.0.0"
input-file: https://github.com/Azure/azure-rest-api-specs/blob/1e2c9f3ec93078da8078389941531359e274f32a/specification/keyvault/data-plane/Microsoft.KeyVault/stable/7.2/secrets.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: false
output-folder: internal
tag: package-7.2
credential-scope: none
use: "@autorest/go@4.0.0-preview.27"
module-version: 0.1.0
# security: "AADToken"
# security-scopes: "https://storage.azure.com/.default"
```

<!-- ### Go multi-api

``` yaml $(go) && $(multiapi)
batch:
  - tag: package-2019-02
``` -->