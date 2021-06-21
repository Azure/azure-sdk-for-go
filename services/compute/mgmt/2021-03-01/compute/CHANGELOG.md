# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. PublicIPAddressSku.PublicIPAddressSkuName
1. PublicIPAddressSku.PublicIPAddressSkuTier

## Additive Changes

### New Funcs

1. *DiskRestorePointGrantAccessFuture.UnmarshalJSON([]byte) error
1. *DiskRestorePointRevokeAccessFuture.UnmarshalJSON([]byte) error
1. DiskRestorePointClient.GrantAccess(context.Context, string, string, string, string, GrantAccessData) (DiskRestorePointGrantAccessFuture, error)
1. DiskRestorePointClient.GrantAccessPreparer(context.Context, string, string, string, string, GrantAccessData) (*http.Request, error)
1. DiskRestorePointClient.GrantAccessResponder(*http.Response) (AccessURI, error)
1. DiskRestorePointClient.GrantAccessSender(*http.Request) (DiskRestorePointGrantAccessFuture, error)
1. DiskRestorePointClient.RevokeAccess(context.Context, string, string, string, string) (DiskRestorePointRevokeAccessFuture, error)
1. DiskRestorePointClient.RevokeAccessPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. DiskRestorePointClient.RevokeAccessResponder(*http.Response) (autorest.Response, error)
1. DiskRestorePointClient.RevokeAccessSender(*http.Request) (DiskRestorePointRevokeAccessFuture, error)

### Struct Changes

#### New Structs

1. DiskRestorePointGrantAccessFuture
1. DiskRestorePointRevokeAccessFuture

#### New Struct Fields

1. PublicIPAddressSku.Name
1. PublicIPAddressSku.Tier
