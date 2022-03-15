# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. BaseClient.ProviderName
1. BaseClient.ResourceName
1. BaseClient.ResourceType

### Signature Changes

#### Funcs

1. ManagementAssociationsClient.CreateOrUpdate
	- Params
		- From: context.Context, string, string, ManagementAssociation
		- To: context.Context, string, string, string, string, string, ManagementAssociation
1. ManagementAssociationsClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, string, ManagementAssociation
		- To: context.Context, string, string, string, string, string, ManagementAssociation
1. ManagementAssociationsClient.Delete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string, string, string
1. ManagementAssociationsClient.DeletePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string, string, string
1. ManagementAssociationsClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string, string, string
1. ManagementAssociationsClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string, string, string
1. New
	- Params
		- From: string, string, string, string
		- To: string
1. NewManagementAssociationsClient
	- Params
		- From: string, string, string, string
		- To: string
1. NewManagementAssociationsClientWithBaseURI
	- Params
		- From: string, string, string, string, string
		- To: string, string
1. NewManagementConfigurationsClient
	- Params
		- From: string, string, string, string
		- To: string
1. NewManagementConfigurationsClientWithBaseURI
	- Params
		- From: string, string, string, string, string
		- To: string, string
1. NewOperationsClient
	- Params
		- From: string, string, string, string
		- To: string
1. NewOperationsClientWithBaseURI
	- Params
		- From: string, string, string, string, string
		- To: string, string
1. NewSolutionsClient
	- Params
		- From: string, string, string, string
		- To: string
1. NewSolutionsClientWithBaseURI
	- Params
		- From: string, string, string, string, string
		- To: string, string
1. NewWithBaseURI
	- Params
		- From: string, string, string, string, string
		- To: string, string
