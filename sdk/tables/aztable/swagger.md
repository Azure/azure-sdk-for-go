## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
input-file: https://raw.githubusercontent.com/Azure/azure-rest-api-specs/master/specification/cosmos-db/data-plane/Microsoft.Tables/preview/2019-02-02/table.json
license-header: MICROSOFT_MIT_NO_VERSION
namespace: aztables
clear-output-folder: true
output-folder: ./
```

``` yaml
directive:
  # dynamically change TableEntityProperties from map[string]interface{} to []byte
  - from: source-file-go
    where: $.definitions.TableEntityProperties
    transform: >-
      $["type"] = "array"
      $["items"] = {
          "type": "string",
          "format": "byte",
      }
      $lib.log($);
```

### Go multi-api

``` yaml $(go) && $(multiapi)
batch:
  - tag: package-2019-02
```