# Code Generation - Azure Datalake SDK for Golang

### Settings

```yaml
go: true
clear-output-folder: false
version: "^3.0.0"
license-header: MICROSOFT_MIT_NO_VERSION
input-file: "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/storage/data-plane/Microsoft.StorageDataLake/preview/2020-10-02/DataLakeStorage.json"
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
use: "@autorest/go@4.0.0-preview.45"
```

### Remove Filesystem and PathName from parameter list since they are not needed
``` yaml
directive:
- from: swagger-document
  where: $["x-ms-paths"]
  transform: >
    for (const property in $)
    {
        if (property.includes('/{filesystem}/{path}'))
        {
            $[property]["parameters"] = $[property]["parameters"].filter(function(param) { return (typeof param['$ref'] === "undefined") || (false == param['$ref'].endsWith("#/parameters/FileSystem") && false == param['$ref'].endsWith("#/parameters/Path"))});
        }
        else if (property.includes('/{filesystem}'))
        {
            $[property]["parameters"] = $[property]["parameters"].filter(function(param) { return (typeof param['$ref'] === "undefined") || (false == param['$ref'].endsWith("#/parameters/FileSystem"))});
        }
    }
```

### Remove pager methods and export various generated methods in filesystem client

``` yaml
directive:
  - from: zz_filesystem_client.go
    where: $
    transform: >-
      return $.
        replace(/func \(client \*FileSystemClient\) NewListBlobHierarchySegmentPager\(.+\/\/ listBlobHierarchySegmentCreateRequest creates the ListBlobHierarchySegment request/s, `//\n// ListBlobHierarchySegmentCreateRequest creates the ListBlobHierarchySegment request`).
        replace(/\(client \*FileSystemClient\) listBlobHierarchySegmentCreateRequest\(/, `(client *FileSystemClient) ListBlobHierarchySegmentCreateRequest(`).
        replace(/\(client \*FileSystemClient\) listBlobHierarchySegmentHandleResponse\(/, `(client *FileSystemClient) ListBlobHierarchySegmentHandleResponse(`);
```

### Remove pager methods and export various generated methods in filesystem client

``` yaml
directive:
  - from: zz_filesystem_client.go
    where: $
    transform: >-
      return $.
        replace(/func \(client \*FileSystemClient\) NewListPathsPager\(.+\/\/ listPathsCreateRequest creates the ListPaths request/s, `//\n// ListPathsCreateRequest creates the ListPaths request`).
        replace(/\(client \*FileSystemClient\) listPathsCreateRequest\(/, `(client *FileSystemClient) ListPathsCreateRequest(`).
        replace(/\(client \*FileSystemClient\) listPathsHandleResponse\(/, `(client *FileSystemClient) ListPathsHandleResponse(`);
```

### Remove pager methods and export various generated methods in service client

``` yaml
directive:
  - from: zz_service_client.go
    where: $
    transform: >-
      return $.
        replace(/func \(client \*ServiceClient\) NewListFileSystemsPager\(.+\/\/ listFileSystemsCreateRequest creates the ListFileSystems request/s, `//\n// ListFileSystemsCreateRequest creates the ListFileSystems request`).
        replace(/\(client \*ServiceClient\) listFileSystemsCreateRequest\(/, `(client *FileSystemClient) ListFileSystemsCreateRequest(`).
        replace(/\(client \*ServiceClient\) listFileSystemsHandleResponse\(/, `(client *FileSystemClient) ListFileSystemsHandleResponse(`);
```