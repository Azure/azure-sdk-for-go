## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file: https://github.com/Azure/azure-rest-api-specs/blob/d744b6bcb95ab4034832ded556dbbe58f4287c5b/specification/cosmos-db/data-plane/Microsoft.Tables/preview/2019-02-02/table.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: false
output-folder: internal
file-prefix: "zz_"
tag: package-2019-02
credential-scope: none
use: "@autorest/go@4.0.0-preview.59"
security: "AADToken"
security-scopes: "https://storage.azure.com/.default"
honor-body-placement: true
modelerfour:
  group-parameters: false
  seal-single-value-enum-by-default: true

directive:
  - from: zz_table_client.go
    where: $
    transform: >-
      return $.
        replace(/\(client \*TableClient\) deleteEntityCreateRequest\(/, `(client *TableClient) DeleteEntityCreateRequest(`).
        replace(/\(client \*TableClient\) insertEntityCreateRequest\(/, `(client *TableClient) InsertEntityCreateRequest(`).
        replace(/\(client \*TableClient\) mergeEntityCreateRequest\(/, `(client *TableClient) MergeEntityCreateRequest(`).
        replace(/\(client \*TableClient\) updateEntityCreateRequest\(/, `(client *TableClient) UpdateEntityCreateRequest(`).
        replace(/= client\.deleteEntityCreateRequest\(/, `= client.DeleteEntityCreateRequest(`).
        replace(/= client\.insertEntityCreateRequest\(/, `= client.InsertEntityCreateRequest(`).
        replace(/= client\.mergeEntityCreateRequest\(/, `= client.MergeEntityCreateRequest(`).
        replace(/= client\.updateEntityCreateRequest\(/, `= client.UpdateEntityCreateRequest(`).
        replace(/if rowKey == "" \{\s*.*\s*\}\s*/g, ``);
  - from:
      - zz_time_rfc1123.go
      - zz_time_rfc3339.go
    where: $
    transform: return $.replace(/UnmarshalText\(data\s+\[\]byte\)\s+(?:error|\(error\))\s+\{\s/g, `UnmarshalText(data []byte) error {\n\tif len(data) == 0 {\n\t\treturn nil\n\t}\n`);
```

### Go multi-api

``` yaml $(go) && $(multiapi)
batch:
  - tag: package-2019-02
```
