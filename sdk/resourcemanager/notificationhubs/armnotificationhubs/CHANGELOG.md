# Release History

## 0.3.0 (2021-12-15)
### Breaking Changes

- Field `SharedAccessAuthorizationRuleListResult` of struct `NamespacesListKeysResult` has been removed

### Features Added

- New anonymous field `ResourceListKeys` in struct `NamespacesListKeysResult`


## 0.2.0 (2021-10-29)

### Breaking Changes

- `arm.Connection` has been removed in `github.com/Azure/azure-sdk-for-go/sdk/azcore/v0.20.0`
- The parameters of `NewXXXClient` has been changed from `(con *arm.Connection, subscriptionID string)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`

## 0.1.0 (2021-10-26)

- Initial preview release.
