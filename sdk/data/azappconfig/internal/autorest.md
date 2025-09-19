## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file:
- https://github.com/Azure/azure-rest-api-specs/blob/c1af3ab8e803da2f40fc90217a6d023bc13b677f/specification/appconfiguration/data-plane/Microsoft.AppConfiguration/stable/2023-11-01/appconfiguration.json
- appconfiguration_ext.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: false
file-prefix: "zz_"
output-folder: generated
openapi-type: "data-plane"
security: "AADToken"
use: "@autorest/go@4.0.0-preview.63"
slice-elements-byval: true
modelerfour:
  lenient-model-deduplication: true
```

### Fix up parameter names
```yaml
directive:
# Directive renaming "KeyValueFields" value to "SettingFields".
- from: swagger-document
  where: '$.parameters.KeyValueFields.items.x-ms-enum'
  transform: >
    $["name"] = "SettingFields";
```

### Fix up pagers
```yaml
directive:
- from: swagger-document
  where: $.paths.*.get.x-ms-pageable
  transform: >
    $.operationName = "GetNextPage";
- from: zz_azureappconfiguration_client.go
  where: $
  transform: >
    return $.
      replace(/urlPath\s+:=\s+"\/\{nextLink\}"/, "urlPath := nextLink").
      replace(/\s+urlPath\s+=\s+strings\.ReplaceAll\(urlPath, "\{nextLink\}", nextLink\)/, "");
```

```yaml
directive:
- from: zz_azureappconfiguration_client.go
  where: $
  transform: >
    return $.replace(/createSnapshot\(/g, "CreateSnapshot(");
```

```yaml
directive:
- from: zz_constants.go
  where: $
  transform: >
    return $.replace(/SnapshotFieldsEtag/g, "SnapshotFieldsETag");
```
