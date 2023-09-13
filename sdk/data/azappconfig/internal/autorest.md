## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file:
- https://github.com/Azure/azure-rest-api-specs/blob/e01d8afe9be7633ed36db014af16d47fec01f737/specification/appconfiguration/data-plane/Microsoft.AppConfiguration/stable/1.0/appconfiguration.json
- appconfiguration_ext.json
license-header: MICROSOFT_MIT_NO_VERSION
clear-output-folder: false
file-prefix: "zz_"
output-folder: generated
openapi-type: "data-plane"
security: "AADToken"
use: "@autorest/go@4.0.0-preview.51"
```

### Fix up enums

``` yaml
directive:
- from: swagger-document
  where: $.paths./kv
  transform: >
    $.get.parameters[6].items["x-ms-enum"] = {
        "name": "SettingFields",
        "modelAsString": true
    };
    $.head.parameters[6].items["x-ms-enum"] = {
        "name": "SettingFields",
        "modelAsString": true
    };
- from: swagger-document
  where: $.paths./kv/{key}
  transform: >
    $.get.parameters[7].items["x-ms-enum"] = {
        "name": "SettingFields",
        "modelAsString": true
    };
    $.head.parameters[7].items["x-ms-enum"] = {
        "name": "SettingFields",
        "modelAsString": true
    };
- from: swagger-document
  where: $.paths./labels
  transform: >
    $.get.parameters[5].items["x-ms-enum"] = {
        "name": "LabelFields",
        "modelAsString": true
    };
    $.head.parameters[5].items["x-ms-enum"] = {
        "name": "LabelFields",
        "modelAsString": true
    };
- from: swagger-document
  where: $.paths./revisions
  transform: >
    $.get.parameters[6].items["x-ms-enum"] = {
        "name": "SettingFields",
        "modelAsString": true
    };
    $.head.parameters[6].items["x-ms-enum"] = {
        "name": "SettingFields",
        "modelAsString": true
    };
```

### Fix up pagers
```yaml
directive:
- from: swagger-document
  where: $.paths.*.get.x-ms-pageable
  transform: >
    $.operationName = "GetNextPage";
- from: azureappconfiguration_client.go
  where: $
  transform: >
    return $.
      replace(/urlPath\s+:=\s+"\/\{nextLink\}"/, "urlPath := nextLink").
      replace(/\s+urlPath\s+=\s+strings\.ReplaceAll\(urlPath, "\{nextLink\}", nextLink\)/, "");
```
