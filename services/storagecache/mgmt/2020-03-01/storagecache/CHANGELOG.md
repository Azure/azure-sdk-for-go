# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Structs

1. SetObject

### Signature Changes

#### Funcs

1. CachesClient.DeleteResponder
	- Returns
		- From: SetObject, error
		- To: autorest.Response, error
1. CachesClient.FlushResponder
	- Returns
		- From: SetObject, error
		- To: autorest.Response, error
1. CachesClient.StartResponder
	- Returns
		- From: SetObject, error
		- To: autorest.Response, error
1. CachesClient.StopResponder
	- Returns
		- From: SetObject, error
		- To: autorest.Response, error
1. CachesClient.UpgradeFirmwareResponder
	- Returns
		- From: SetObject, error
		- To: autorest.Response, error
1. StorageTargetsClient.DeleteResponder
	- Returns
		- From: SetObject, error
		- To: autorest.Response, error

#### Struct Fields

1. CachesDeleteFuture.Result changed type from func(CachesClient) (SetObject, error) to func(CachesClient) (autorest.Response, error)
1. CachesFlushFuture.Result changed type from func(CachesClient) (SetObject, error) to func(CachesClient) (autorest.Response, error)
1. CachesStartFuture.Result changed type from func(CachesClient) (SetObject, error) to func(CachesClient) (autorest.Response, error)
1. CachesStopFuture.Result changed type from func(CachesClient) (SetObject, error) to func(CachesClient) (autorest.Response, error)
1. CachesUpgradeFirmwareFuture.Result changed type from func(CachesClient) (SetObject, error) to func(CachesClient) (autorest.Response, error)
1. StorageTargetsDeleteFuture.Result changed type from func(StorageTargetsClient) (SetObject, error) to func(StorageTargetsClient) (autorest.Response, error)

## Additive Changes

### New Funcs

1. *AscOperation.UnmarshalJSON([]byte) error
1. AscOperation.MarshalJSON() ([]byte, error)
1. AscOperationProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. AscOperationProperties

#### New Struct Fields

1. AscOperation.*AscOperationProperties
