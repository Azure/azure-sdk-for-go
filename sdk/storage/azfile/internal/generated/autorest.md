# Code Generation - Azure File SDK for Golang

### Settings

```yaml
go: true
clear-output-folder: false
version: "^3.0.0"
license-header: MICROSOFT_MIT_NO_VERSION
input-file: "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/b6472ffd34d5d4a155101b41b4eb1f356abff600/specification/storage/data-plane/Microsoft.FileStorage/stable/2026-02-06/file.json"
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
honor-body-placement: true
use: "@autorest/go@4.0.0-preview.61"
```

### Updating service version to 2026-02-06

```yaml
directive:
- from: 
  - zz_directory_client.go
  - zz_file_client.go
  - zz_share_client.go
  - zz_service_client.go
  where: $
  transform: >-
    return $.
      replaceAll(`[]string{"2025-11-05"}`, `[]string{ServiceVersion}`);
```
### Changing casing of NfsFileType, Nfs, ShareNfsSettingsEncryptionInTransit and ShareNfsSettings

```yaml
directive:
- from: 
  - zz_constants.go
  - zz_options.go
  - zz_response_types.go
  - zz_file_client.go
  - zz_directory_client.go
  - zz_models.go
  where: $
  transform: >-
    return $.
      replaceAll(`NfsFileType`, `NFSFileType`).
      replaceAll(`ShareNfsSettings`, `ShareNFSSettings`).
      replaceAll(`ShareNfsSettingsEncryptionInTransit`, `ShareNFSSettingsEncryptionInTransit`).
      replaceAll(`Nfs *`, `NFS *`);
```

### Updating Header Names XMSFileShareSnapshotUsageBytes and XMSFileShareUsageBytes

```yaml
directive:
- from: 
  - zz_response_types.go
  - zz_share_client.go
  where: $
  transform: >-
    return $.
      replaceAll(`XMSFileShareSnapshotUsageBytes`, `FileShareSnapshotUsageBytes`).
      replaceAll(`XMSFileShareUsageBytes`, `FileShareUsageBytes`);
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

### use azcore.ETag

``` yaml
directive:
- from: zz_models.go
  where: $
  transform: >-
    return $.
      replace(/import "time"/, `import (\n\t"time"\n\t"github.com/Azure/azure-sdk-for-go/sdk/azcore"\n)`).
      replace(/Etag\s+\*string/g, `ETag *azcore.ETag`);

- from: zz_response_types.go
  where: $
  transform: >-
    return $.
      replace(/"time"/, `"time"\n\t"github.com/Azure/azure-sdk-for-go/sdk/azcore"`).
      replace(/ETag\s+\*string/g, `ETag *azcore.ETag`);

- from:
  - zz_directory_client.go
  - zz_file_client.go
  - zz_share_client.go
  where: $
  transform: >-
    return $.
      replace(/"github\.com\/Azure\/azure\-sdk\-for\-go\/sdk\/azcore\/policy"/, `"github.com/Azure/azure-sdk-for-go/sdk/azcore"\n\t"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"`).
      replace(/result\.ETag\s+=\s+&val/g, `result.ETag = (*azcore.ETag)(&val)`);
```

### Rename models - remove `Share` prefix

``` yaml
directive:
- rename-model:
    from: ShareProtocolSettings
    to: ProtocolSettings
- rename-model:
    from: ShareSmbSettings
    to: SMBSettings
```

### Capitalise SMB field

``` yaml
directive:
- from:
  - zz_directory_client.go
  - zz_file_client.go
  - zz_models.go
  - zz_options.go
  - zz_share_client.go
  - zz_response_types.go
  where: $
  transform: >-
    return $.
      replace(/SmbMultichannel/g, `SMBMultichannel`).
      replace(/copyFileSmbInfo/g, `copyFileSMBInfo`).
      replace(/CopyFileSmbInfo/g, `CopyFileSMBInfo`).
      replace(/Smb\s+\*ShareSMBSettings/g, `SMB *ShareSMBSettings`).
      replace(/EnableSmbDirectoryLease/g, `EnableSMBDirectoryLease`); 
```

### Fixing casing of SignedTid and SignedOid

``` yaml
directive:
- from: zz_models.go
  where: $
  transform: >-
    return $.
      replace(/SignedOid\s+\*string/g, `SignedOID *string`).
      replace(/SignedTid\s+\*string/g, `SignedTID *string`);
```

### Rename models - remove `Item` and `Internal` suffix

``` yaml
directive:
- rename-model:
    from: DirectoryItem
    to: Directory
- rename-model:
    from: FileItem
    to: File
- rename-model:
    from: HandleItem
    to: Handle
- rename-model:
    from: ShareItemInternal
    to: Share
- rename-model:
    from: SharePropertiesInternal
    to: ShareProperties
```

### Remove `Items` and `List` suffix

``` yaml
directive:
  - from: source-file-go
    where: $
    transform: >-
      return $.
        replace(/DirectoryItems/g, "Directories").
        replace(/FileItems/g, "Files").
        replace(/ShareItems/g, "Shares").
        replace(/HandleList/g, "Handles");
```

### Rename `FileID` to `ID` (except for Handle object)

``` yaml
directive:
- from: swagger-document
  where: $.definitions
  transform: >
    $.Directory.properties.FileId["x-ms-client-name"] = "ID";
    $.File.properties.FileId["x-ms-client-name"] = "ID";
    $.Handle.properties.HandleId["x-ms-client-name"] = "ID";

- from:
  - zz_directory_client.go
  - zz_file_client.go
  - zz_response_types.go
  where: $
  transform: >-
    return $.
      replace(/FileID/g, `ID`);
```


### Change CORS acronym to be all caps and rename `FileParentID` to `ParentID`

``` yaml
directive:
  - from: source-file-go
    where: $
    transform: >-
      return $.
        replace(/Cors/g, "CORS").
        replace(/FileParentID/g, "ParentID");
```

### Change cors xml to be correct

``` yaml
directive:
  - from: source-file-go
    where: $
    transform: >-
      return $.
        replace(/xml:"CORS>CORSRule"/g, "xml:\"Cors>CorsRule\"");
```

### Remove pager methods and export various generated methods in service client

``` yaml
directive:
  - from: zz_service_client.go
    where: $
    transform: >-
      return $.
        replace(/func \(client \*ServiceClient\) NewListSharesSegmentPager\(.+\/\/ listSharesSegmentCreateRequest creates the ListSharesSegment request/s, `//\n// listSharesSegmentCreateRequest creates the ListSharesSegment request`).
        replace(/\(client \*ServiceClient\) listSharesSegmentCreateRequest\(/, `(client *ServiceClient) ListSharesSegmentCreateRequest(`).
        replace(/\(client \*ServiceClient\) listSharesSegmentHandleResponse\(/, `(client *ServiceClient) ListSharesSegmentHandleResponse(`);
```

### Use string type for FileCreationTime and FileLastWriteTime

``` yaml
directive:
- from: swagger-document
  where: $.parameters.FileCreationTime
  transform: >
    $.format = "str";
- from: swagger-document
  where: $.parameters.FileLastWriteTime
  transform: >
    $.format = "str";
- from: swagger-document
  where: $.parameters.FileChangeTime
  transform: >
    $.format = "str";
```

### Remove pager methods and export various generated methods in directory client

``` yaml
directive:
  - from: zz_directory_client.go
    where: $
    transform: >-
      return $.
        replace(/func \(client \*DirectoryClient\) NewListFilesAndDirectoriesSegmentPager\(.+\/\/ listFilesAndDirectoriesSegmentCreateRequest creates the ListFilesAndDirectoriesSegment request/s, `//\n// listFilesAndDirectoriesSegmentCreateRequest creates the ListFilesAndDirectoriesSegment request`).
        replace(/\(client \*DirectoryClient\) listFilesAndDirectoriesSegmentCreateRequest\(/, `(client *DirectoryClient) ListFilesAndDirectoriesSegmentCreateRequest(`).
        replace(/\(client \*DirectoryClient\) listFilesAndDirectoriesSegmentHandleResponse\(/, `(client *DirectoryClient) ListFilesAndDirectoriesSegmentHandleResponse(`);
```

### Fix time format for parsing the response headers: x-ms-file-creation-time, x-ms-file-last-write-time, x-ms-file-change-time

``` yaml
directive:
  - from:
    - zz_directory_client.go
    - zz_file_client.go
    where: $
    transform: >-
      return $.
        replace(/fileCreationTime,\s+err\s+\:=\s+time\.Parse\(time\.RFC1123,\s+val\)/g, `fileCreationTime, err := time.Parse(ISO8601, val)`).
        replace(/fileLastWriteTime,\s+err\s+\:=\s+time\.Parse\(time\.RFC1123,\s+val\)/g, `fileLastWriteTime, err := time.Parse(ISO8601, val)`).
        replace(/fileChangeTime,\s+err\s+\:=\s+time\.Parse\(time\.RFC1123,\s+val\)/g, `fileChangeTime, err := time.Parse(ISO8601, val)`);
```

### Change `Duration` parameter in leases to be required

``` yaml
directive:
- from: swagger-document
  where: $.parameters.LeaseDuration
  transform: >
    $.required = true;
```

### Convert ShareUsageBytes to int64

``` yaml
directive:
  - from: zz_models.go
    where: $
    transform: >-
      return $.
        replace(/ShareUsageBytes\s+\*int32/g, `ShareUsageBytes *int64`);
```

### Convert StringEncoded to string type

``` yaml
directive:
  - from: zz_models.go
    where: $
    transform: >-
      return $.
        replace(/\*StringEncoded/g, `*string`);
```

### Removing UnmarshalXML for Handle to create custom UnmarshalXML function

``` yaml
directive:
- from: swagger-document
  where: $.definitions
  transform: >
    $.Handle["x-ms-go-omit-serde-methods"] = true;
```

### Convert FileAttributes to an optional parameter

``` yaml
directive:
- from: swagger-document
  where: $.parameters.FileAttributes
  transform: >
    $.required = false;
```

### Rename ProvisionedBandwidthMiBps response field

``` yaml
directive:
- from: swagger-document
  where: $["x-ms-paths"]["/{shareName}?restype=share"]
  transform: >
    $.get.responses["200"].headers["x-ms-share-provisioned-bandwidth-mibps"]["x-ms-client-name"] = "ProvisionedBandwidthMiBps"
```