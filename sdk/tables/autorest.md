## Go

These settings apply only when `--go` is specified on the command line.

<!-- Original autorest command used by Chris Scott -->
<!-- autorest --use=@autorest/go@4.0.0-preview.20 https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/specification/cosmos-db/data-plane/readme.md --tag=package-2019-02 --file-prefix="zz_generated_" --modelerfour.lenient-model-deduplication --license-header=MICROSOFT_MIT_NO_VERSION --output-folder=aztables --module=aztables --openapi-type="data-plane" --credential-scope=none -->

``` yaml
go: true
version: "^3.0.0"
input-file: https://github.com/Azure/azure-rest-api-specs/blob/d744b6bcb95ab4034832ded556dbbe58f4287c5b/specification/cosmos-db/data-plane/Microsoft.Tables/preview/2019-02-02/table.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: false
output-folder: aztable
file-prefix: "zz_generated_"
# namespace: aztable
tag: package-2019-02
credential-scope: none
use: "@autorest/go@4.0.0-preview.23"
openapi-type: data-plane
```

<!-- This change will have to be made by hand for now, leaving this directive in case it is useful later -->
<!--
``` yaml
directive:
  # dynamically change TableEntityProperties from map[string]interface{} to []byte
  - from: swagger-document
    where: $.definitions.TableEntityProperties
    transform: >-
      $["type"] = "string";
      $["format"] = "byte";
      delete $.additionalProperties;
``` -->

### Go multi-api

``` yaml $(go) && $(multiapi)
batch:
  - tag: package-2019-02
```