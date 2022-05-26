# Unreleased

## Breaking Changes

### Removed Constants

1. ItemType.ItemTypeFolder

### Signature Changes

#### Funcs

1. WorkbooksClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, *bool
1. WorkbooksClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, *bool

## Additive Changes

### New Constants

1. ItemType.ItemTypeNone
