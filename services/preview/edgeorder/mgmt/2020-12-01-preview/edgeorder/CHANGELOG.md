# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. BaseClient.ListOrderItemsAtSubscriptionLevel
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. BaseClient.ListOrderItemsAtSubscriptionLevelComplete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. BaseClient.ListOrderItemsAtSubscriptionLevelPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string

#### Struct Fields

1. OrderItemDetails.ForwardShippingDetails changed type from *ShippingDetails to *ForwardShippingDetails
1. OrderItemDetails.ManagementRpDetails changed type from interface{} to *ResourceProviderDetails
1. OrderItemDetails.ReverseShippingDetails changed type from *ShippingDetails to *ReverseShippingDetails

## Additive Changes

### New Funcs

1. ForwardShippingDetails.MarshalJSON() ([]byte, error)
1. ResourceProviderDetails.MarshalJSON() ([]byte, error)
1. ReverseShippingDetails.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ForwardShippingDetails
1. ResourceProviderDetails
1. ReverseShippingDetails

#### New Struct Fields

1. DeviceDetails.ManagementResourceTenantID
1. OrderItemDetails.ManagementRpDetailsList
1. ProductDetails.ProductDoubleEncryptionStatus
1. ProductFamilyProperties.ResourceProviderDetails
1. ReturnOrderItemDetails.ServiceTag
1. ReturnOrderItemDetails.ShippingBoxRequired
