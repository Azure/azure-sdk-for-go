# Unreleased

## Breaking Changes

### Removed Funcs

1. ErrorResponse.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. ClusterErrorResponse
1. DataExportErrorResponse
1. ErrorContract

#### Removed Struct Fields

1. ErrorResponse.AdditionalInfo
1. ErrorResponse.Code
1. ErrorResponse.Details
1. ErrorResponse.Message
1. ErrorResponse.Target

## Additive Changes

### New Constants

1. WorkspaceSkuNameEnum.WorkspaceSkuNameEnumLACluster

### New Funcs

1. ErrorDetail.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ErrorDetail

#### New Struct Fields

1. ErrorResponse.Error
1. WorkspaceProperties.CreatedDate
1. WorkspaceProperties.Features
1. WorkspaceProperties.ForceCmkForQuery
1. WorkspaceProperties.ModifiedDate
