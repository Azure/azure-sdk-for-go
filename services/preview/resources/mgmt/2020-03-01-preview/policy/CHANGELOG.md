# Unreleased

## Breaking Changes

### Removed Funcs

1. AzureEntityResource.MarshalJSON() ([]byte, error)
1. ErrorDetail.MarshalJSON() ([]byte, error)
1. ProxyResource.MarshalJSON() ([]byte, error)
1. Resource.MarshalJSON() ([]byte, error)
1. TrackedResource.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. AzureEntityResource
1. ErrorDetail
1. ProxyResource
1. Resource
1. TrackedResource

### Signature Changes

#### Funcs

1. AssignmentsClient.List
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, *int32
1. AssignmentsClient.ListComplete
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, *int32
1. AssignmentsClient.ListForResource
	- Params
		- From: context.Context, string, string, string, string, string, string, string, *int32
		- To: context.Context, string, string, string, string, string, string, *int32
1. AssignmentsClient.ListForResourceComplete
	- Params
		- From: context.Context, string, string, string, string, string, string, string, *int32
		- To: context.Context, string, string, string, string, string, string, *int32
1. AssignmentsClient.ListForResourceGroup
	- Params
		- From: context.Context, string, string, string, *int32
		- To: context.Context, string, string, *int32
1. AssignmentsClient.ListForResourceGroupComplete
	- Params
		- From: context.Context, string, string, string, *int32
		- To: context.Context, string, string, *int32
1. AssignmentsClient.ListForResourceGroupPreparer
	- Params
		- From: context.Context, string, string, string, *int32
		- To: context.Context, string, string, *int32
1. AssignmentsClient.ListForResourcePreparer
	- Params
		- From: context.Context, string, string, string, string, string, string, string, *int32
		- To: context.Context, string, string, string, string, string, string, *int32
1. AssignmentsClient.ListPreparer
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, *int32
1. DefinitionsClient.CreateOrUpdate
	- Params
		- From: context.Context, string, Definition, string
		- To: context.Context, string, Definition
1. DefinitionsClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, Definition, string
		- To: context.Context, string, Definition
1. DefinitionsClient.Delete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. DefinitionsClient.DeletePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. DefinitionsClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. DefinitionsClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. DefinitionsClient.List
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, *int32
1. DefinitionsClient.ListComplete
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, *int32
1. DefinitionsClient.ListPreparer
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, *int32
1. ExemptionsClient.CreateOrUpdate
	- Params
		- From: context.Context, Exemption, string, string
		- To: context.Context, string, string, Exemption
1. ExemptionsClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, Exemption, string, string
		- To: context.Context, string, string, Exemption
1. ExemptionsClient.List
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. ExemptionsClient.ListComplete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. ExemptionsClient.ListForResource
	- Params
		- From: context.Context, string, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string
1. ExemptionsClient.ListForResourceComplete
	- Params
		- From: context.Context, string, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string
1. ExemptionsClient.ListForResourceGroup
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. ExemptionsClient.ListForResourceGroupComplete
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. ExemptionsClient.ListForResourceGroupPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. ExemptionsClient.ListForResourcePreparer
	- Params
		- From: context.Context, string, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string
1. ExemptionsClient.ListPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. New
	- Params
		- From: <none>
		- To: string
1. NewAssignmentsClient
	- Params
		- From: <none>
		- To: string
1. NewAssignmentsClientWithBaseURI
	- Params
		- From: string
		- To: string, string
1. NewDefinitionsClient
	- Params
		- From: <none>
		- To: string
1. NewDefinitionsClientWithBaseURI
	- Params
		- From: string
		- To: string, string
1. NewExemptionsClient
	- Params
		- From: <none>
		- To: string
1. NewExemptionsClientWithBaseURI
	- Params
		- From: string
		- To: string, string
1. NewSetDefinitionsClient
	- Params
		- From: <none>
		- To: string
1. NewSetDefinitionsClientWithBaseURI
	- Params
		- From: string
		- To: string, string
1. NewWithBaseURI
	- Params
		- From: string
		- To: string, string
1. SetDefinitionsClient.CreateOrUpdate
	- Params
		- From: context.Context, string, SetDefinition, string
		- To: context.Context, string, SetDefinition
1. SetDefinitionsClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, SetDefinition, string
		- To: context.Context, string, SetDefinition
1. SetDefinitionsClient.Delete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. SetDefinitionsClient.DeletePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. SetDefinitionsClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. SetDefinitionsClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. SetDefinitionsClient.List
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, *int32
1. SetDefinitionsClient.ListComplete
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, *int32
1. SetDefinitionsClient.ListPreparer
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, *int32

#### Struct Fields

1. CloudError.Error changed type from *ErrorDetail to *ErrorResponse

## Additive Changes

### New Funcs

1. ErrorResponse.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ErrorResponse

#### New Struct Fields

1. BaseClient.SubscriptionID
