Generated from https://github.com/Azure/azure-rest-api-specs/tree/0b17e6a5e811fd7b122d383b4942441d95e5e8cf/specification/storage/resource-manager/readme.md tag: `package-2020-08-preview`

Code generator @microsoft.azure/autorest.go@2.1.169


## Breaking Changes

## Signature Changes

### Funcs

1. AccountsClient.ListByResourceGroup
	- Returns
		- From: AccountListResult, error
		- To: AccountListResultPage, error

### New Funcs

1. AccountsClient.ListByResourceGroupComplete(context.Context, string) (AccountListResultIterator, error)

## Struct Changes

### New Structs

1. ManagementPolicyVersion

### New Struct Fields

1. AccountProperties.AllowSharedKeyAccess
1. AccountPropertiesCreateParameters.AllowSharedKeyAccess
1. AccountPropertiesUpdateParameters.AllowSharedKeyAccess
1. ChangeFeed.RetentionInDays
1. ManagementPolicyAction.Version
1. ManagementPolicySnapShot.TierToArchive
1. ManagementPolicySnapShot.TierToCool
