# Code Generation - Azure Storage File SDK for Golang

<!-- autorest --use=@autorest/go@4.0.0-preview.35 https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/storage/data-plane/Microsoft.BlobStorage/preview/2020-10-02/blob.json --file-prefix="zz_generated_" --modelerfour.lenient-model-deduplication --license-header=MICROSOFT_MIT_NO_VERSION --output-folder=generated/ --module=azblob --openapi-type="data-plane" --credential-scope=none -->

```bash
cd swagger
autorest autorest.md
gofmt -w generated/*
```

### Settings

```yaml
go: true
clear-output-folder: false
version: "^3.0.0"
license-header: MICROSOFT_MIT_NO_VERSION
input-file: "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/storage/data-plane/Microsoft.FileStorage/preview/2020-10-02/file.json"
module: "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile"
credential-scope: "https://storage.azure.com/.default"
output-folder: internal/
file-prefix: "zz_generated_"
openapi-type: "data-plane"
verbose: true
security: AzureKey
module-version: "0.1.0"
modelerfour:
  group-parameters: false
  seal-single-value-enum-by-default: true
  lenient-model-deduplication: true
export-clients: false
use: "@autorest/go@4.0.0-preview.41"
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

### Metrics

``` yaml
directive:
- from: swagger-document
  where: $.definitions
  transform: >
    $.Metrics.type = "object";
```

### Use strings for dates in responses. (Easy to handle the conversion in clients)

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

### ShareFileRangeList

``` yaml
directive:
- from: swagger-document
  where: $.definitions
  transform: >
    $.ShareFileRangeList.xml = {
        "name": "Ranges"
    };
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

### Remove ShareName, Directory, and FileName - we have direct URIs

``` yaml
directive:
- from: swagger-document
  where: $["x-ms-paths"]
  transform: >
   Object.keys($).map(id => {
     if (id.includes('/{shareName}/{directory}/{fileName}'))
     {
       $[id.replace('/{shareName}/{directory}/{fileName}', '?shareName_dir_file')] = $[id];
       delete $[id];
     }
     else if (id.includes('/{shareName}/{directory}'))
     {
       $[id.replace('/{shareName}/{directory}', '?shareName_dir')] = $[id];
       delete $[id];
     }
     else if (id.includes('/{shareName}'))
     {
       $[id.replace('/{shareName}', '?shareName')] = $[id];
       delete $[id];
     }
   });
```

### ShareErrorCode

``` yaml
directive:
- from: swagger-document
  where: $.definitions.ErrorCode
  transform: >
    $["x-ms-enum"].name = "ShareErrorCode";
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

### Rename ShareFileHTTPHeaders to FileHttpHeader and remove file prefix from properties

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