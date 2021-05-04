# Unreleased Content

## Breaking Changes

### Struct Changes

#### Removed Structs

1. ChangeableAttributes

### Signature Changes

#### Funcs

1. ManifestsClient.UpdateAttributes
	- Params
		- From: context.Context, string, string, *ChangeableAttributes
		- To: context.Context, string, string, *ManifestChangeableAttributes
1. ManifestsClient.UpdateAttributesPreparer
	- Params
		- From: context.Context, string, string, *ChangeableAttributes
		- To: context.Context, string, string, *ManifestChangeableAttributes
1. RepositoryClient.UpdateAttributes
	- Params
		- From: context.Context, string, *ChangeableAttributes
		- To: context.Context, string, *RepositoryChangeableAttributes
1. RepositoryClient.UpdateAttributesPreparer
	- Params
		- From: context.Context, string, *ChangeableAttributes
		- To: context.Context, string, *RepositoryChangeableAttributes
1. TagClient.UpdateAttributes
	- Params
		- From: context.Context, string, string, *ChangeableAttributes
		- To: context.Context, string, string, *TagChangeableAttributes
1. TagClient.UpdateAttributesPreparer
	- Params
		- From: context.Context, string, string, *ChangeableAttributes
		- To: context.Context, string, string, *TagChangeableAttributes

#### Struct Fields

1. ManifestAttributesBase.ChangeableAttributes changed type from *ChangeableAttributes to *ManifestChangeableAttributes
1. RepositoryAttributes.ChangeableAttributes changed type from *ChangeableAttributes to *RepositoryChangeableAttributes
1. TagAttributesBase.ChangeableAttributes changed type from *ChangeableAttributes to *TagChangeableAttributes

## Additive Changes

### Struct Changes

#### New Structs

1. RepositoryChangeableAttributes
1. TagChangeableAttributes
