# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. ProviderOperationsMetadataClient.Get
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. ProviderOperationsMetadataClient.GetPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. ProviderOperationsMetadataClient.List
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. ProviderOperationsMetadataClient.ListComplete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. ProviderOperationsMetadataClient.ListPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string

## Additive Changes

### New Funcs

1. ErrorAdditionalInfo.MarshalJSON() ([]byte, error)
1. ErrorDetail.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ErrorAdditionalInfo
1. ErrorDetail
1. ErrorResponse
