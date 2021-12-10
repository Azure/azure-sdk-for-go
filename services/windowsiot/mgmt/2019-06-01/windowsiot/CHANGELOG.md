# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. ErrorDetails.Code
1. ErrorDetails.Details
1. ErrorDetails.Message
1. ErrorDetails.Target

### Signature Changes

#### Funcs

1. ServicesClient.CreateOrUpdate
	- Params
		- From: context.Context, string, string, DeviceServiceProperties, string
		- To: context.Context, string, string, DeviceService, string
1. ServicesClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, string, DeviceServiceProperties, string
		- To: context.Context, string, string, DeviceService, string
1. ServicesClient.Update
	- Params
		- From: context.Context, string, string, DeviceServiceProperties, string
		- To: context.Context, string, string, DeviceService, string
1. ServicesClient.UpdatePreparer
	- Params
		- From: context.Context, string, string, DeviceServiceProperties, string
		- To: context.Context, string, string, DeviceService, string

## Additive Changes

### Struct Changes

#### New Structs

1. ErrorDetailsError

#### New Struct Fields

1. ErrorDetails.Error
1. OperationEntity.IsDataAction
1. OperationEntity.Origin
