Generated from https://github.com/Azure/azure-rest-api-specs/tree/92ab22b49bd085116af0c61fada2c6c360702e9e/specification/managementgroups/resource-manager/readme.md tag: `package-2020-02`

Code generator @microsoft.azure/autorest.go@2.1.175


## Breaking Changes

### Removed Constants

1. Type.MicrosoftManagementmanagementGroup
1. Type1.Subscriptions

## Struct Changes

### Removed Struct Fields

1. BaseClient.OperationResultID
1. BaseClient.Skip
1. BaseClient.Skiptoken
1. BaseClient.Top

## Signature Changes

### Const Types

1. MicrosoftManagementmanagementGroups changed type from Type1 to Type

### Funcs

1. Client.GetDescendants
	- Params
		- From: context.Context, string
		- To: context.Context, string, string, *int32
1. Client.GetDescendantsComplete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string, *int32
1. Client.GetDescendantsPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string, *int32
1. Client.List
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. Client.ListComplete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. Client.ListPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. EntitiesClient.List
	- Params
		- From: context.Context, string, string, string, string, string, string
		- To: context.Context, string, *int32, *int32, string, string, string, string, string, string
1. EntitiesClient.ListComplete
	- Params
		- From: context.Context, string, string, string, string, string, string
		- To: context.Context, string, *int32, *int32, string, string, string, string, string, string
1. EntitiesClient.ListPreparer
	- Params
		- From: context.Context, string, string, string, string, string, string
		- To: context.Context, string, *int32, *int32, string, string, string, string, string, string
1. New
	- Params
		- From: string, *int32, *int32, string
		- To: <none>
1. NewClient
	- Params
		- From: string, *int32, *int32, string
		- To: <none>
1. NewClientWithBaseURI
	- Params
		- From: string, string, *int32, *int32, string
		- To: string
1. NewEntitiesClient
	- Params
		- From: string, *int32, *int32, string
		- To: <none>
1. NewEntitiesClientWithBaseURI
	- Params
		- From: string, string, *int32, *int32, string
		- To: string
1. NewHierarchySettingsClient
	- Params
		- From: string, *int32, *int32, string
		- To: <none>
1. NewHierarchySettingsClientWithBaseURI
	- Params
		- From: string, string, *int32, *int32, string
		- To: string
1. NewOperationsClient
	- Params
		- From: string, *int32, *int32, string
		- To: <none>
1. NewOperationsClientWithBaseURI
	- Params
		- From: string, string, *int32, *int32, string
		- To: string
1. NewSubscriptionsClient
	- Params
		- From: string, *int32, *int32, string
		- To: <none>
1. NewSubscriptionsClientWithBaseURI
	- Params
		- From: string, string, *int32, *int32, string
		- To: string
1. NewWithBaseURI
	- Params
		- From: string, string, *int32, *int32, string
		- To: string

### New Constants

1. Type1.Type1MicrosoftManagementmanagementGroups
1. Type1.Type1Subscriptions
