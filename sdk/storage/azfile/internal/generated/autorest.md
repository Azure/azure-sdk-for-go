# Code Generation - Azure Blob SDK for Golang

### Settings

```yaml
go: true
clear-output-folder: false
version: "^3.0.0"
license-header: MICROSOFT_MIT_NO_VERSION
input-file: "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/storage/data-plane/Microsoft.FileStorage/preview/2020-10-02/file.json"
credential-scope: "https://storage.azure.com/.default"
output-folder: ../generated
file-prefix: "zz_"
openapi-type: "data-plane"
verbose: true
security: AzureKey
modelerfour:
  group-parameters: false
  seal-single-value-enum-by-default: true
  lenient-model-deduplication: true
export-clients: true
use: "@autorest/go@4.0.0-preview.43"
```

### Don't include share name, directory, or file name in path - we have direct URIs.

``` yaml
directive:
- from: swagger-document
  where: $["x-ms-paths"]
  transform: >
    for (const property in $)
    {
        if (property.includes('/{shareName}/{directory}/{fileName}'))
        {
            $[property]["parameters"] = $[property]["parameters"].filter(function(param) { return (typeof param['$ref'] === "undefined") || (false == param['$ref'].endsWith("#/parameters/ShareName") && false == param['$ref'].endsWith("#/parameters/DirectoryPath") && false == param['$ref'].endsWith("#/parameters/FilePath"))});
        } 
        else if (property.includes('/{shareName}/{directory}'))
        {
            $[property]["parameters"] = $[property]["parameters"].filter(function(param) { return (typeof param['$ref'] === "undefined") || (false == param['$ref'].endsWith("#/parameters/ShareName") && false == param['$ref'].endsWith("#/parameters/DirectoryPath"))});
        }
        else if (property.includes('/{shareName}'))
        {
            $[property]["parameters"] = $[property]["parameters"].filter(function(param) { return (typeof param['$ref'] === "undefined") || (false == param['$ref'].endsWith("#/parameters/ShareName"))});
        }
    }
```

### Use strings for dates in responses

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
- from: swagger-document
  where: $.parameters.FileChangeTime
  transform: >
    $.format = "str";
```
