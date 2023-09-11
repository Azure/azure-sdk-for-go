## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
go: true
version: "^3.0.0"
input-file:
- https://github.com/Azure/azure-rest-api-specs/blob/498ccf7ddf78ced8ef515f88b755b2eb3775de9e/specification/appconfiguration/data-plane/Microsoft.AppConfiguration/stable/1.0/appconfiguration.json
- appconfiguration_ext.json
openapi-type: "data-plane"
module: github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig
license-header: MICROSOFT_MIT_NO_VERSION
output-folder: ../azappconfig
clear-output-folder: false
override-client-name: Client
security: "AADToken"
use: "@autorest/go@4.0.0-preview.49"
export-clients: true
```

### Fix up enums

``` yaml
directive:

# PutKeyValue ==> AddSetting
- rename-operation:
    from: PutKeyValue
    to: AddSetting

# PutKeyValue ==> AddSetting
- rename-operation:
    from: GetKeyValue
    to: GetSetting

# PutKeyValue ==> AddSetting
- rename-operation:
    from: PutLock
    to: SetReadOnly

# etag -> ETag
#- from: swagger-document
#  where: $.paths..parameters..[?(@.name=='etag')]
#  rename-property:
#    from: etag
#    to: ETag

# rename KeyValueFieldsEtag to KeyValueFieldsETag
- from:
    - client.go
    - models.go
    - response_types.go
    - setting_selector.go
    - setting.go
    - constants.go
  where: $
  transform: return $.replace(/(KeyValueFieldsEtag)/g, "KeyValueFieldsETag");

# rename KeyValueFieldsLocked to KeyValueFieldsIsReadOnly
- from:
    - client.go
    - models.go
    - response_types.go
    - setting_selector.go
    - setting.go
    - constants.go
  where: $
  transform: return $.replace(/(KeyValueFieldsLocked)/g, "KeyValueFieldsIsReadOnly");

# rename KeyValueFieldsLocked to KeyValueFieldsIsReadOnly
- from:
    - models.go
    - models_serde.go    
    - setting.go
  where: $
  transform: return $.replace(/(Locked)/g, "IsReadOnly");

- from:
    - models.go    
    - models_serde.go
    - setting.go
  where: $
  transform: return $.replace(/(Etag)/g, "ETag");

# delete client name prefix from method options and response types
- from:
    - client.go
    - models.go
    - response_types.go
  where: $
  transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");

# change type of syncToken from *string to *syncTokenPolicy
- from: client.go
  where: $
  transform: return $.replace(/syncToken \*string/g, "syncTokenPolicy *syncTokenPolicy");


# Don't process the Sync-Token header, let the policy do that
- from: client.go
  where: $
  transform: return $.replace(/if client.syncToken != nil {\n\s+req.Raw\(\).Header\["Sync-Token"\] = \[\]string{\*client.syncToken}\n\s+}/g, "");
- from: client.go
  where: $
  transform: return $.replace(/if val := resp.Header.Get\("Sync-token"\); val != "" {\n\s+result.SyncToken = &val\n\s+}/g, "");

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
