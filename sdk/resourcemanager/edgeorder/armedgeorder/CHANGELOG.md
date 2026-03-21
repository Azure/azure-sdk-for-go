# Release History

## 2.0.0-beta.1 (2026-03-19)
### Breaking Changes

- Function `*ClientFactory.NewManagementClient` has been removed
- Function `NewManagementClient` has been removed
- Function `*ManagementClient.CancelOrderItem` has been removed
- Function `*ManagementClient.BeginCreateAddress` has been removed
- Function `*ManagementClient.BeginCreateOrderItem` has been removed
- Function `*ManagementClient.BeginDeleteAddressByName` has been removed
- Function `*ManagementClient.BeginDeleteOrderItemByName` has been removed
- Function `*ManagementClient.GetAddressByName` has been removed
- Function `*ManagementClient.GetOrderByName` has been removed
- Function `*ManagementClient.GetOrderItemByName` has been removed
- Function `*ManagementClient.NewListAddressesAtResourceGroupLevelPager` has been removed
- Function `*ManagementClient.NewListAddressesAtSubscriptionLevelPager` has been removed
- Function `*ManagementClient.NewListConfigurationsPager` has been removed
- Function `*ManagementClient.NewListOperationsPager` has been removed
- Function `*ManagementClient.NewListOrderAtResourceGroupLevelPager` has been removed
- Function `*ManagementClient.NewListOrderAtSubscriptionLevelPager` has been removed
- Function `*ManagementClient.NewListOrderItemsAtResourceGroupLevelPager` has been removed
- Function `*ManagementClient.NewListOrderItemsAtSubscriptionLevelPager` has been removed
- Function `*ManagementClient.NewListProductFamiliesMetadataPager` has been removed
- Function `*ManagementClient.NewListProductFamiliesPager` has been removed
- Function `*ManagementClient.BeginReturnOrderItem` has been removed
- Function `*ManagementClient.BeginUpdateAddress` has been removed
- Function `*ManagementClient.BeginUpdateOrderItem` has been removed
- Struct `BasicInformation` has been removed
- Struct `CommonProperties` has been removed
- Struct `ConfigurationFilters` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `ShippingDetails` has been removed
- Struct `TrackedResource` has been removed
- Field `ConfigurationFilters` of struct `ConfigurationsRequest` has been removed
- Field `ManagementRpDetails` of struct `OrderItemDetails` has been removed
- Field `Count`, `DeviceDetails` of struct `ProductDetails` has been removed

### Features Added

- New value `AvailabilityStageDiscoverable` added to enum type `AvailabilityStage`
- New value `LinkTypeDiscoverable` added to enum type `LinkType`
- New value `OrderItemTypeExternal` added to enum type `OrderItemType`
- New value `StageNameReadyToSetup` added to enum type `StageName`
- New enum type `AddressClassification` with values `AddressClassificationShipping`, `AddressClassificationSite`
- New enum type `AutoProvisioningStatus` with values `AutoProvisioningStatusDisabled`, `AutoProvisioningStatusEnabled`
- New enum type `ChildConfigurationType` with values `ChildConfigurationTypeAdditionalConfiguration`, `ChildConfigurationTypeDeviceConfiguration`
- New enum type `DevicePresenceVerificationStatus` with values `DevicePresenceVerificationStatusCompleted`, `DevicePresenceVerificationStatusNotInitiated`
- New enum type `FulfillmentType` with values `FulfillmentTypeExternal`, `FulfillmentTypeMicrosoft`
- New enum type `IdentificationType` with values `IdentificationTypeNotSupported`, `IdentificationTypeSerialNumber`
- New enum type `OrderMode` with values `OrderModeDefault`, `OrderModeDoNotFulfill`
- New enum type `ProvisioningState` with values `ProvisioningStateCanceled`, `ProvisioningStateCreating`, `ProvisioningStateFailed`, `ProvisioningStateSucceeded`
- New enum type `ProvisioningSupport` with values `ProvisioningSupportCloudBased`, `ProvisioningSupportManual`
- New enum type `TermCommitmentType` with values `TermCommitmentTypeNone`, `TermCommitmentTypeTimed`, `TermCommitmentTypeTrial`
- New function `NewAddressesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AddressesClient, error)`
- New function `*AddressesClient.BeginCreate(ctx context.Context, resourceGroupName string, addressName string, addressResource AddressResource, options *AddressesClientBeginCreateOptions) (*runtime.Poller[AddressesClientCreateResponse], error)`
- New function `*AddressesClient.BeginDelete(ctx context.Context, resourceGroupName string, addressName string, options *AddressesClientBeginDeleteOptions) (*runtime.Poller[AddressesClientDeleteResponse], error)`
- New function `*AddressesClient.Get(ctx context.Context, resourceGroupName string, addressName string, options *AddressesClientGetOptions) (AddressesClientGetResponse, error)`
- New function `*AddressesClient.NewListByResourceGroupPager(resourceGroupName string, options *AddressesClientListByResourceGroupOptions) *runtime.Pager[AddressesClientListByResourceGroupResponse]`
- New function `*AddressesClient.NewListBySubscriptionPager(options *AddressesClientListBySubscriptionOptions) *runtime.Pager[AddressesClientListBySubscriptionResponse]`
- New function `*AddressesClient.BeginUpdate(ctx context.Context, resourceGroupName string, addressName string, addressUpdateParameter AddressUpdateParameter, options *AddressesClientBeginUpdateOptions) (*runtime.Poller[AddressesClientUpdateResponse], error)`
- New function `*ClientFactory.NewAddressesClient() *AddressesClient`
- New function `*ClientFactory.NewOperationsClient() *OperationsClient`
- New function `*ClientFactory.NewOrderItemsClient() *OrderItemsClient`
- New function `*ClientFactory.NewOrdersClient() *OrdersClient`
- New function `*ClientFactory.NewProductsAndConfigurationsClient() *ProductsAndConfigurationsClient`
- New function `NewOperationsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationsClient, error)`
- New function `*OperationsClient.NewListPager(options *OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `NewOrderItemsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*OrderItemsClient, error)`
- New function `*OrderItemsClient.Cancel(ctx context.Context, resourceGroupName string, orderItemName string, cancellationReason CancellationReason, options *OrderItemsClientCancelOptions) (OrderItemsClientCancelResponse, error)`
- New function `*OrderItemsClient.BeginCreate(ctx context.Context, resourceGroupName string, orderItemName string, orderItemResource OrderItemResource, options *OrderItemsClientBeginCreateOptions) (*runtime.Poller[OrderItemsClientCreateResponse], error)`
- New function `*OrderItemsClient.BeginDelete(ctx context.Context, resourceGroupName string, orderItemName string, options *OrderItemsClientBeginDeleteOptions) (*runtime.Poller[OrderItemsClientDeleteResponse], error)`
- New function `*OrderItemsClient.Get(ctx context.Context, resourceGroupName string, orderItemName string, options *OrderItemsClientGetOptions) (OrderItemsClientGetResponse, error)`
- New function `*OrderItemsClient.NewListByResourceGroupPager(resourceGroupName string, options *OrderItemsClientListByResourceGroupOptions) *runtime.Pager[OrderItemsClientListByResourceGroupResponse]`
- New function `*OrderItemsClient.NewListBySubscriptionPager(options *OrderItemsClientListBySubscriptionOptions) *runtime.Pager[OrderItemsClientListBySubscriptionResponse]`
- New function `*OrderItemsClient.BeginReturn(ctx context.Context, resourceGroupName string, orderItemName string, returnOrderItemDetails ReturnOrderItemDetails, options *OrderItemsClientBeginReturnOptions) (*runtime.Poller[OrderItemsClientReturnResponse], error)`
- New function `*OrderItemsClient.BeginUpdate(ctx context.Context, resourceGroupName string, orderItemName string, orderItemUpdateParameter OrderItemUpdateParameter, options *OrderItemsClientBeginUpdateOptions) (*runtime.Poller[OrderItemsClientUpdateResponse], error)`
- New function `NewOrdersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*OrdersClient, error)`
- New function `*OrdersClient.Get(ctx context.Context, resourceGroupName string, location string, orderName string, options *OrdersClientGetOptions) (OrdersClientGetResponse, error)`
- New function `*OrdersClient.NewListByResourceGroupPager(resourceGroupName string, options *OrdersClientListByResourceGroupOptions) *runtime.Pager[OrdersClientListByResourceGroupResponse]`
- New function `*OrdersClient.NewListBySubscriptionPager(options *OrdersClientListBySubscriptionOptions) *runtime.Pager[OrdersClientListBySubscriptionResponse]`
- New function `NewProductsAndConfigurationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ProductsAndConfigurationsClient, error)`
- New function `*ProductsAndConfigurationsClient.NewListConfigurationsPager(configurationsRequest ConfigurationsRequest, options *ProductsAndConfigurationsClientListConfigurationsOptions) *runtime.Pager[ProductsAndConfigurationsClientListConfigurationsResponse]`
- New function `*ProductsAndConfigurationsClient.NewListProductFamiliesMetadataPager(options *ProductsAndConfigurationsClientListProductFamiliesMetadataOptions) *runtime.Pager[ProductsAndConfigurationsClientListProductFamiliesMetadataResponse]`
- New function `*ProductsAndConfigurationsClient.NewListProductFamiliesPager(productFamiliesRequest ProductFamiliesRequest, options *ProductsAndConfigurationsClientListProductFamiliesOptions) *runtime.Pager[ProductsAndConfigurationsClientListProductFamiliesResponse]`
- New struct `AdditionalConfiguration`
- New struct `CategoryInformation`
- New struct `ChildConfiguration`
- New struct `ChildConfigurationFilter`
- New struct `ChildConfigurationProperties`
- New struct `ConfigurationDeviceDetails`
- New struct `ConfigurationFilter`
- New struct `DevicePresenceVerificationDetails`
- New struct `GroupedChildConfigurations`
- New struct `OrderItemDetailsUpdateParameter`
- New struct `ProductDetailsUpdateParameter`
- New struct `ProvisioningDetails`
- New struct `SiteDetails`
- New struct `TermCommitmentInformation`
- New struct `TermCommitmentPreferences`
- New struct `TermTypeDetails`
- New struct `UserAssignedIdentity`
- New field `AddressClassification`, `ProvisioningState` in struct `AddressProperties`
- New field `TermTypeDetails` in struct `BillingMeterDetails`
- New field `ChildConfigurationTypes`, `FulfilledBy`, `GroupedChildConfigurations`, `ProvisioningSupport`, `SupportedTermCommitmentDurations` in struct `ConfigurationProperties`
- New field `ConfigurationFilter` in struct `ConfigurationsRequest`
- New field `DisplaySerialNumber`, `ProvisioningDetails`, `ProvisioningSupport` in struct `DeviceDetails`
- New field `ConfigurationIDDisplayName` in struct `HierarchyInformation`
- New field `OrderItemMode`, `SiteDetails` in struct `OrderItemDetails`
- New field `ProvisioningState` in struct `OrderItemProperties`
- New field `Identity` in struct `OrderItemResource`
- New field `Identity` in struct `OrderItemUpdateParameter`
- New field `OrderItemDetails` in struct `OrderItemUpdateProperties`
- New field `OrderMode` in struct `OrderProperties`
- New field `TermCommitmentPreferences` in struct `Preferences`
- New field `ChildConfigurationDeviceDetails`, `IdentificationType`, `OptInAdditionalConfigurations`, `ParentDeviceDetails`, `ParentProvisioningDetails`, `TermCommitmentInformation` in struct `ProductDetails`
- New field `FulfilledBy` in struct `ProductFamilyProperties`
- New field `FulfilledBy` in struct `ProductLineProperties`
- New field `FulfilledBy` in struct `ProductProperties`
- New field `UserAssignedIdentities` in struct `ResourceIdentity`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/edgeorder/armedgeorder` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).