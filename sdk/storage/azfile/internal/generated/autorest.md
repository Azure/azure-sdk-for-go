# Code Generation - Azure File SDK for Golang

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

### Don't include share name, directory, or file name in path - we have direct URIs

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

### Add Last-Modified to SetMetadata

``` yaml
directive:
- from: swagger-document
  where: $["x-ms-paths"]["/{shareName}/{directory}/{fileName}?comp=metadata"]
  transform: >
    $.put.responses["200"].headers["Last-Modified"] = {
        "type": "string",
        "format": "date-time-rfc1123",
        "description": "Returns the date and time the file was last modified. Any operation that modifies the file, including an update of the file's metadata or properties, changes the last-modified time of the file."
    }
```

### Add Content-MD5 to Put Range from URL

``` yaml
directive:
- from: swagger-document
  where: $["x-ms-paths"]["/{shareName}/{directory}/{fileName}?comp=range&fromURL"]
  transform: >
    $.put.responses["201"].headers["Content-MD5"] = {
        "type": "string",
        "format": "byte",
        "description": "This header is returned so that the client can check for message content integrity. The value of this header is computed by the File service; it is not necessarily the same value as may have been specified in the request headers."
    }
```

### ShareServiceProperties, ShareMetrics, ShareCorsRule, and ShareRetentionPolicy

``` yaml
directive:
- rename-model:
    from: Metrics
    to: ShareMetrics
- rename-model:
    from: CorsRule
    to: ShareCorsRule
- rename-model:
    from: RetentionPolicy
    to: ShareRetentionPolicy
- rename-model:
    from: StorageServiceProperties
    to: ShareServiceProperties
    
- from: swagger-document
  where: $.definitions
  transform: >
    $.ShareMetrics.properties.IncludeAPIs["x-ms-client-name"] = "IncludeApis";
    $.ShareServiceProperties.xml = {"name": "StorageServiceProperties"};
    $.ShareCorsRule.xml = {"name": "CorsRule"};
- from: swagger-document
  where: $.parameters
  transform: >
    $.StorageServiceProperties.name = "ShareServiceProperties";
```

### Rename FileHttpHeaders to ShareFileHTTPHeaders and remove file prefix from properties

``` yaml
directive:
- from: swagger-document
  where: $.parameters
  transform: >
    $.FileCacheControl["x-ms-parameter-grouping"].name = "share-file-http-headers";
    $.FileCacheControl["x-ms-client-name"] = "cacheControl";
    $.FileContentDisposition["x-ms-parameter-grouping"].name = "share-file-http-headers";
    $.FileContentDisposition["x-ms-client-name"] = "contentDisposition";
    $.FileContentEncoding["x-ms-parameter-grouping"].name = "share-file-http-headers";
    $.FileContentEncoding["x-ms-client-name"] = "contentEncoding";
    $.FileContentLanguage["x-ms-parameter-grouping"].name = "share-file-http-headers";
    $.FileContentLanguage["x-ms-client-name"] = "contentLanguage";
    $.FileContentMD5["x-ms-parameter-grouping"].name = "share-file-http-headers";
    $.FileContentMD5["x-ms-client-name"] = "contentMd5";
    $.FileContentType["x-ms-parameter-grouping"].name = "share-file-http-headers";
    $.FileContentType["x-ms-client-name"] = "contentType";
```
