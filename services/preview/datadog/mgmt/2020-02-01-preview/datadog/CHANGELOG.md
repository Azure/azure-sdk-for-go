# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. MonitorsClient.Update
	- Returns
		- From: MonitorResource, error
		- To: MonitorsUpdateFuture, error
1. MonitorsClient.UpdateSender
	- Returns
		- From: *http.Response, error
		- To: MonitorsUpdateFuture, error

## Additive Changes

### New Funcs

1. *MonitorsUpdateFuture.UnmarshalJSON([]byte) error

### Struct Changes

#### New Structs

1. MonitorsUpdateFuture

#### New Struct Fields

1. MonitorResourceUpdateParameters.Sku
