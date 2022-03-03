# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Structs

1. RestorePointProvisioningDetails

#### Removed Struct Fields

1. PublicIPAddressSku.PublicIPAddressSkuName
1. PublicIPAddressSku.PublicIPAddressSkuTier
1. RestorePoint.ConsistencyMode
1. RestorePoint.ExcludeDisks
1. RestorePoint.ProvisioningDetails
1. RestorePoint.ProvisioningState
1. RestorePoint.SourceMetadata

### Signature Changes

#### Struct Fields

1. OrchestrationServiceStateInput.ServiceName changed type from *string to OrchestrationServiceNames

## Additive Changes

### New Funcs

1. *DiskRestorePointGrantAccessFuture.UnmarshalJSON([]byte) error
1. *DiskRestorePointRevokeAccessFuture.UnmarshalJSON([]byte) error
1. *RestorePoint.UnmarshalJSON([]byte) error
1. DiskRestorePointClient.GrantAccess(context.Context, string, string, string, string, GrantAccessData) (DiskRestorePointGrantAccessFuture, error)
1. DiskRestorePointClient.GrantAccessPreparer(context.Context, string, string, string, string, GrantAccessData) (*http.Request, error)
1. DiskRestorePointClient.GrantAccessResponder(*http.Response) (AccessURI, error)
1. DiskRestorePointClient.GrantAccessSender(*http.Request) (DiskRestorePointGrantAccessFuture, error)
1. DiskRestorePointClient.RevokeAccess(context.Context, string, string, string, string) (DiskRestorePointRevokeAccessFuture, error)
1. DiskRestorePointClient.RevokeAccessPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. DiskRestorePointClient.RevokeAccessResponder(*http.Response) (autorest.Response, error)
1. DiskRestorePointClient.RevokeAccessSender(*http.Request) (DiskRestorePointRevokeAccessFuture, error)
1. RestorePointProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. DiskRestorePointGrantAccessFuture
1. DiskRestorePointRevokeAccessFuture
1. RestorePointProperties

#### New Struct Fields

1. PublicIPAddressSku.Name
1. PublicIPAddressSku.Tier
1. RestorePoint.*RestorePointProperties
1. RestorePointSourceMetadata.Location
