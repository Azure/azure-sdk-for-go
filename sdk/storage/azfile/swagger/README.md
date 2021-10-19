# Code Generation - Azure File SDK for Golang

```bash
cd swagger
autorest README.md --use="https://github.com/Azure/autorest.go/releases/download/v4.0.0-preview.27/autorest-go-4.0.0-preview.27.tgz"
gofmt -w generated/*
```

### Settings

```yaml
clear-output-folder: true
license-header: MICROSOFT_MIT_NO_VERSION
input-file: "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/storage-dataplane-preview/specification/storage/data-plane/Microsoft.FileStorage/preview/2020-02-10/file.json"
module: "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile"
credential-scope: "https://storage.azure.com/.default"
output-folder: generated/
file-prefix: "zz_generated_"
modelerfour.lenient-model-deduplication: true
openapi-type: "data-plane"
verbose: true
security: AzureKey
module-version: "0.1.0"
```

``` yaml
directive:
- from: swagger-document
  where: $["x-ms-paths"]..responses..headers["x-ms-file-last-write-time"]
  transform: >
    $.format = "str";
- from: swagger-document
  where: $["x-ms-paths"]..responses..headers["x-ms-file-change-time"]
  transform: >
    $.format = "str";
- from: swagger-document
  where: $["x-ms-paths"]..responses..headers["x-ms-file-creation-time"]
  transform: >
    $.format = "str";
```

### Change new SMB file parameters to use default values

``` yaml
directive:
- from: swagger-document
  where: $.parameters.FileCreationTime
  transform: >
    $.format = "str";
    $.default = "now";
- from: swagger-document
  where: $.parameters.FileLastWriteTime
  transform: >
    $.format = "str";
    $.default = "now";
- from: swagger-document
  where: $.parameters.FileAttributes
  transform: >
    $.default = "none";
- from: swagger-document
  where: $.parameters.FilePermission
  transform: >
    $.default = "inherit";
```

### FileRangeWriteFromUrl Constant

This value is supposed to be the constant value update and these changes turn it from a parameter into a constant.

``` yaml
directive:
- from: swagger-document
  where: $.parameters.FileRangeWriteFromUrl
  transform: >
    delete $.default;
    delete $["x-ms-enum"];
    $["x-ms-parameter-location"] = "method";
```