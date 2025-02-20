# Release History

## 2.0.0 (2025-02-28)
### Breaking Changes

- Function `NewClientFactory` parameter(s) have been changed from `(azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`

### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `RequestState` with values `RequestStateAccepted`, `RequestStateCanceled`, `RequestStateCreated`, `RequestStateFailed`, `RequestStateInProgress`, `RequestStateInvalid`, `RequestStateSucceeded`
- New function `*ClientFactory.NewGroupQuotaLimitsClient() *GroupQuotaLimitsClient`
- New function `*ClientFactory.NewGroupQuotaLimitsRequestClient() *GroupQuotaLimitsRequestClient`
- New function `*ClientFactory.NewGroupQuotaSubscriptionAllocationClient() *GroupQuotaSubscriptionAllocationClient`
- New function `*ClientFactory.NewGroupQuotaSubscriptionAllocationRequestClient() *GroupQuotaSubscriptionAllocationRequestClient`
- New function `*ClientFactory.NewGroupQuotaSubscriptionRequestsClient() *GroupQuotaSubscriptionRequestsClient`
- New function `*ClientFactory.NewGroupQuotaSubscriptionsClient() *GroupQuotaSubscriptionsClient`
- New function `*ClientFactory.NewGroupQuotasClient() *GroupQuotasClient`
- New function `NewGroupQuotaLimitsClient(azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaLimitsClient, error)`
- New function `*GroupQuotaLimitsClient.List(context.Context, string, string, string, string, *GroupQuotaLimitsClientListOptions) (GroupQuotaLimitsClientListResponse, error)`
- New function `NewGroupQuotaLimitsRequestClient(azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaLimitsRequestClient, error)`
- New function `*GroupQuotaLimitsRequestClient.Get(context.Context, string, string, string, *GroupQuotaLimitsRequestClientGetOptions) (GroupQuotaLimitsRequestClientGetResponse, error)`
- New function `*GroupQuotaLimitsRequestClient.NewListPager(string, string, string, string, *GroupQuotaLimitsRequestClientListOptions) *runtime.Pager[GroupQuotaLimitsRequestClientListResponse]`
- New function `*GroupQuotaLimitsRequestClient.BeginUpdate(context.Context, string, string, string, string, *GroupQuotaLimitsRequestClientBeginUpdateOptions) (*runtime.Poller[GroupQuotaLimitsRequestClientUpdateResponse], error)`
- New function `NewGroupQuotaSubscriptionAllocationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaSubscriptionAllocationClient, error)`
- New function `*GroupQuotaSubscriptionAllocationClient.List(context.Context, string, string, string, string, *GroupQuotaSubscriptionAllocationClientListOptions) (GroupQuotaSubscriptionAllocationClientListResponse, error)`
- New function `NewGroupQuotaSubscriptionAllocationRequestClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaSubscriptionAllocationRequestClient, error)`
- New function `*GroupQuotaSubscriptionAllocationRequestClient.Get(context.Context, string, string, string, string, *GroupQuotaSubscriptionAllocationRequestClientGetOptions) (GroupQuotaSubscriptionAllocationRequestClientGetResponse, error)`
- New function `*GroupQuotaSubscriptionAllocationRequestClient.NewListPager(string, string, string, string, *GroupQuotaSubscriptionAllocationRequestClientListOptions) *runtime.Pager[GroupQuotaSubscriptionAllocationRequestClientListResponse]`
- New function `*GroupQuotaSubscriptionAllocationRequestClient.BeginUpdate(context.Context, string, string, string, string, SubscriptionQuotaAllocationsList, *GroupQuotaSubscriptionAllocationRequestClientBeginUpdateOptions) (*runtime.Poller[GroupQuotaSubscriptionAllocationRequestClientUpdateResponse], error)`
- New function `NewGroupQuotaSubscriptionRequestsClient(azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaSubscriptionRequestsClient, error)`
- New function `*GroupQuotaSubscriptionRequestsClient.Get(context.Context, string, string, string, *GroupQuotaSubscriptionRequestsClientGetOptions) (GroupQuotaSubscriptionRequestsClientGetResponse, error)`
- New function `*GroupQuotaSubscriptionRequestsClient.NewListPager(string, string, *GroupQuotaSubscriptionRequestsClientListOptions) *runtime.Pager[GroupQuotaSubscriptionRequestsClientListResponse]`
- New function `NewGroupQuotaSubscriptionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaSubscriptionsClient, error)`
- New function `*GroupQuotaSubscriptionsClient.BeginCreateOrUpdate(context.Context, string, string, *GroupQuotaSubscriptionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[GroupQuotaSubscriptionsClientCreateOrUpdateResponse], error)`
- New function `*GroupQuotaSubscriptionsClient.BeginDelete(context.Context, string, string, *GroupQuotaSubscriptionsClientBeginDeleteOptions) (*runtime.Poller[GroupQuotaSubscriptionsClientDeleteResponse], error)`
- New function `*GroupQuotaSubscriptionsClient.Get(context.Context, string, string, *GroupQuotaSubscriptionsClientGetOptions) (GroupQuotaSubscriptionsClientGetResponse, error)`
- New function `*GroupQuotaSubscriptionsClient.NewListPager(string, string, *GroupQuotaSubscriptionsClientListOptions) *runtime.Pager[GroupQuotaSubscriptionsClientListResponse]`
- New function `*GroupQuotaSubscriptionsClient.BeginUpdate(context.Context, string, string, *GroupQuotaSubscriptionsClientBeginUpdateOptions) (*runtime.Poller[GroupQuotaSubscriptionsClientUpdateResponse], error)`
- New function `NewGroupQuotasClient(azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotasClient, error)`
- New function `*GroupQuotasClient.BeginCreateOrUpdate(context.Context, string, string, *GroupQuotasClientBeginCreateOrUpdateOptions) (*runtime.Poller[GroupQuotasClientCreateOrUpdateResponse], error)`
- New function `*GroupQuotasClient.BeginDelete(context.Context, string, string, *GroupQuotasClientBeginDeleteOptions) (*runtime.Poller[GroupQuotasClientDeleteResponse], error)`
- New function `*GroupQuotasClient.Get(context.Context, string, string, *GroupQuotasClientGetOptions) (GroupQuotasClientGetResponse, error)`
- New function `*GroupQuotasClient.NewListPager(string, *GroupQuotasClientListOptions) *runtime.Pager[GroupQuotasClientListResponse]`
- New function `*GroupQuotasClient.BeginUpdate(context.Context, string, string, *GroupQuotasClientBeginUpdateOptions) (*runtime.Poller[GroupQuotasClientUpdateResponse], error)`
- New struct `AllocatedQuotaToSubscriptionList`
- New struct `AllocatedToSubscription`
- New struct `AllocationRequestBase`
- New struct `AllocationRequestBaseProperties`
- New struct `AllocationRequestBasePropertiesName`
- New struct `AllocationRequestStatus`
- New struct `AllocationRequestStatusList`
- New struct `AllocationRequestStatusProperties`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New struct `GroupQuotaDetails`
- New struct `GroupQuotaDetailsName`
- New struct `GroupQuotaLimit`
- New struct `GroupQuotaLimitList`
- New struct `GroupQuotaLimitListProperties`
- New struct `GroupQuotaLimitProperties`
- New struct `GroupQuotaList`
- New struct `GroupQuotaRequestBase`
- New struct `GroupQuotaRequestBaseProperties`
- New struct `GroupQuotaRequestBasePropertiesName`
- New struct `GroupQuotaSubscriptionID`
- New struct `GroupQuotaSubscriptionIDList`
- New struct `GroupQuotaSubscriptionIDProperties`
- New struct `GroupQuotaSubscriptionRequestStatus`
- New struct `GroupQuotaSubscriptionRequestStatusList`
- New struct `GroupQuotaSubscriptionRequestStatusProperties`
- New struct `GroupQuotaUsagesBase`
- New struct `GroupQuotaUsagesBaseName`
- New struct `GroupQuotasEntity`
- New struct `GroupQuotasEntityBase`
- New struct `GroupQuotasEntityBasePatch`
- New struct `GroupQuotasEntityPatch`
- New struct `GroupQuotasEntityPatchProperties`
- New struct `GroupQuotasEntityProperties`
- New struct `LROResponse`
- New struct `LROResponseProperties`
- New struct `ProxyResource`
- New struct `Resource`
- New struct `ResourceBaseRequest`
- New struct `ResourceUsageList`
- New struct `ResourceUsages`
- New struct `SubmittedResourceRequestStatus`
- New struct `SubmittedResourceRequestStatusList`
- New struct `SubmittedResourceRequestStatusProperties`
- New struct `SubscriptionGroupQuotaAssignment`
- New struct `SubscriptionQuotaAllocationRequestList`
- New struct `SubscriptionQuotaAllocations`
- New struct `SubscriptionQuotaAllocationsList`
- New struct `SubscriptionQuotaAllocationsListProperties`
- New struct `SubscriptionQuotaAllocationsProperties`
- New struct `SubscriptionQuotaAllocationsStatusList`
- New struct `SubscriptionQuotaDetails`
- New struct `SubscriptionQuotaDetailsName`
- New struct `SystemData`


## 1.1.0-beta.2 (2025-01-22)
### Breaking Changes

- Function `*GroupQuotaSubscriptionAllocationRequestClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, string, AllocationRequestStatus, *GroupQuotaSubscriptionAllocationRequestClientBeginUpdateOptions)` to `(context.Context, string, string, string, string, SubscriptionQuotaAllocationsList, *GroupQuotaSubscriptionAllocationRequestClientBeginUpdateOptions)`
- Function `*GroupQuotaSubscriptionAllocationRequestClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *GroupQuotaSubscriptionAllocationRequestClientGetOptions)` to `(context.Context, string, string, string, string, *GroupQuotaSubscriptionAllocationRequestClientGetOptions)`
- Type of `GroupQuotaLimit.Properties` has been changed from `*GroupQuotaDetails` to `*GroupQuotaLimitProperties`
- Type of `GroupQuotaLimitsRequestClientBeginUpdateOptions.GroupQuotaRequest` has been changed from `*SubmittedResourceRequestStatus` to `*GroupQuotaLimitList`
- Type of `GroupQuotasEntity.Properties` has been changed from `*GroupQuotasEntityBase` to `*GroupQuotasEntityProperties`
- Type of `GroupQuotasEntityPatch.Properties` has been changed from `*GroupQuotasEntityBasePatch` to `*GroupQuotasEntityPatchProperties`
- Type of `SubscriptionQuotaAllocations.Properties` has been changed from `*SubscriptionQuotaDetails` to `*SubscriptionQuotaAllocationsProperties`
- Enum `EnforcementState` has been removed
- Enum `EnvironmentType` has been removed
- Enum `GroupingIDType` has been removed
- Function `*ClientFactory.NewGroupQuotaLocationSettingsClient` has been removed
- Function `*ClientFactory.NewGroupQuotaUsagesClient` has been removed
- Function `*GroupQuotaLimitsClient.Get` has been removed
- Function `*GroupQuotaLimitsRequestClient.BeginCreateOrUpdate` has been removed
- Function `NewGroupQuotaLocationSettingsClient` has been removed
- Function `*GroupQuotaLocationSettingsClient.BeginCreateOrUpdate` has been removed
- Function `*GroupQuotaLocationSettingsClient.Get` has been removed
- Function `*GroupQuotaLocationSettingsClient.NewListPager` has been removed
- Function `*GroupQuotaLocationSettingsClient.BeginUpdate` has been removed
- Function `*GroupQuotaSubscriptionAllocationClient.Get` has been removed
- Function `*GroupQuotaSubscriptionAllocationRequestClient.BeginCreateOrUpdate` has been removed
- Function `NewGroupQuotaUsagesClient` has been removed
- Function `*GroupQuotaUsagesClient.NewListPager` has been removed
- Operation `*GroupQuotaLimitsClient.NewListPager` does not support pagination anymore, use `*GroupQuotaLimitsClient.List` instead.
- Operation `*GroupQuotaSubscriptionAllocationClient.NewListPager` does not support pagination anymore, use `*GroupQuotaSubscriptionAllocationClient.List` instead.
- Struct `AdditionalAttributes` has been removed
- Struct `AdditionalAttributesPatch` has been removed
- Struct `BillingAccountID` has been removed
- Struct `GroupQuotasEnforcementListResponse` has been removed
- Struct `GroupQuotasEnforcementResponse` has been removed
- Struct `GroupQuotasEnforcementResponseProperties` has been removed
- Struct `GroupingID` has been removed
- Field `Region` of struct `GroupQuotaDetails` has been removed
- Field `ID`, `Name`, `SystemData`, `Type` of struct `GroupQuotaLimit` has been removed
- Field `NextLink`, `Value` of struct `GroupQuotaLimitList` has been removed
- Field `SubmittedResourceRequestStatus` of struct `GroupQuotaLimitsRequestClientUpdateResponse` has been removed
- Field `AllocationRequestStatus` of struct `GroupQuotaSubscriptionAllocationRequestClientUpdateResponse` has been removed
- Field `AdditionalAttributes` of struct `GroupQuotasEntityBase` has been removed
- Field `AdditionalAttributes` of struct `GroupQuotasEntityBasePatch` has been removed
- Field `ID`, `Name`, `SystemData`, `Type` of struct `SubscriptionQuotaAllocations` has been removed
- Field `NextLink`, `Value` of struct `SubscriptionQuotaAllocationsList` has been removed
- Field `Region` of struct `SubscriptionQuotaDetails` has been removed

### Features Added

- New struct `GroupQuotaLimitListProperties`
- New struct `GroupQuotaLimitProperties`
- New struct `GroupQuotasEntityPatchProperties`
- New struct `GroupQuotasEntityProperties`
- New struct `SubscriptionQuotaAllocationsListProperties`
- New struct `SubscriptionQuotaAllocationsProperties`
- New field `ResourceName` in struct `GroupQuotaDetails`
- New field `ID`, `Name`, `Properties`, `SystemData`, `Type` in struct `GroupQuotaLimitList`
- New anonymous field `GroupQuotaLimitList` in struct `GroupQuotaLimitsRequestClientUpdateResponse`
- New anonymous field `SubscriptionQuotaAllocationsList` in struct `GroupQuotaSubscriptionAllocationRequestClientUpdateResponse`
- New field `ID`, `Name`, `Properties`, `SystemData`, `Type` in struct `SubscriptionQuotaAllocationsList`
- New field `ResourceName` in struct `SubscriptionQuotaDetails`


## 1.1.0-beta.1 (2024-04-26)
### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `EnforcementState` with values `EnforcementStateDisabled`, `EnforcementStateEnabled`, `EnforcementStateNotAvailable`
- New enum type `EnvironmentType` with values `EnvironmentTypeNonProduction`, `EnvironmentTypeProduction`
- New enum type `GroupingIDType` with values `GroupingIDTypeBillingID`, `GroupingIDTypeServiceTreeID`
- New enum type `RequestState` with values `RequestStateAccepted`, `RequestStateCanceled`, `RequestStateCreated`, `RequestStateFailed`, `RequestStateInProgress`, `RequestStateInvalid`, `RequestStateSucceeded`
- New function `*ClientFactory.NewGroupQuotaLimitsClient() *GroupQuotaLimitsClient`
- New function `*ClientFactory.NewGroupQuotaLimitsRequestClient() *GroupQuotaLimitsRequestClient`
- New function `*ClientFactory.NewGroupQuotaLocationSettingsClient() *GroupQuotaLocationSettingsClient`
- New function `*ClientFactory.NewGroupQuotaSubscriptionAllocationClient() *GroupQuotaSubscriptionAllocationClient`
- New function `*ClientFactory.NewGroupQuotaSubscriptionAllocationRequestClient() *GroupQuotaSubscriptionAllocationRequestClient`
- New function `*ClientFactory.NewGroupQuotaSubscriptionRequestsClient() *GroupQuotaSubscriptionRequestsClient`
- New function `*ClientFactory.NewGroupQuotaSubscriptionsClient() *GroupQuotaSubscriptionsClient`
- New function `*ClientFactory.NewGroupQuotaUsagesClient() *GroupQuotaUsagesClient`
- New function `*ClientFactory.NewGroupQuotasClient() *GroupQuotasClient`
- New function `NewGroupQuotaLimitsClient(azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaLimitsClient, error)`
- New function `*GroupQuotaLimitsClient.Get(context.Context, string, string, string, string, string, *GroupQuotaLimitsClientGetOptions) (GroupQuotaLimitsClientGetResponse, error)`
- New function `*GroupQuotaLimitsClient.NewListPager(string, string, string, string, *GroupQuotaLimitsClientListOptions) *runtime.Pager[GroupQuotaLimitsClientListResponse]`
- New function `NewGroupQuotaLimitsRequestClient(azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaLimitsRequestClient, error)`
- New function `*GroupQuotaLimitsRequestClient.BeginCreateOrUpdate(context.Context, string, string, string, string, *GroupQuotaLimitsRequestClientBeginCreateOrUpdateOptions) (*runtime.Poller[GroupQuotaLimitsRequestClientCreateOrUpdateResponse], error)`
- New function `*GroupQuotaLimitsRequestClient.Get(context.Context, string, string, string, *GroupQuotaLimitsRequestClientGetOptions) (GroupQuotaLimitsRequestClientGetResponse, error)`
- New function `*GroupQuotaLimitsRequestClient.NewListPager(string, string, string, string, *GroupQuotaLimitsRequestClientListOptions) *runtime.Pager[GroupQuotaLimitsRequestClientListResponse]`
- New function `*GroupQuotaLimitsRequestClient.BeginUpdate(context.Context, string, string, string, string, *GroupQuotaLimitsRequestClientBeginUpdateOptions) (*runtime.Poller[GroupQuotaLimitsRequestClientUpdateResponse], error)`
- New function `NewGroupQuotaLocationSettingsClient(azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaLocationSettingsClient, error)`
- New function `*GroupQuotaLocationSettingsClient.BeginCreateOrUpdate(context.Context, string, string, string, string, *GroupQuotaLocationSettingsClientBeginCreateOrUpdateOptions) (*runtime.Poller[GroupQuotaLocationSettingsClientCreateOrUpdateResponse], error)`
- New function `*GroupQuotaLocationSettingsClient.Get(context.Context, string, string, string, string, *GroupQuotaLocationSettingsClientGetOptions) (GroupQuotaLocationSettingsClientGetResponse, error)`
- New function `*GroupQuotaLocationSettingsClient.NewListPager(string, string, string, *GroupQuotaLocationSettingsClientListOptions) *runtime.Pager[GroupQuotaLocationSettingsClientListResponse]`
- New function `*GroupQuotaLocationSettingsClient.BeginUpdate(context.Context, string, string, string, string, *GroupQuotaLocationSettingsClientBeginUpdateOptions) (*runtime.Poller[GroupQuotaLocationSettingsClientUpdateResponse], error)`
- New function `NewGroupQuotaSubscriptionAllocationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaSubscriptionAllocationClient, error)`
- New function `*GroupQuotaSubscriptionAllocationClient.Get(context.Context, string, string, string, string, *GroupQuotaSubscriptionAllocationClientGetOptions) (GroupQuotaSubscriptionAllocationClientGetResponse, error)`
- New function `*GroupQuotaSubscriptionAllocationClient.NewListPager(string, string, string, *GroupQuotaSubscriptionAllocationClientListOptions) *runtime.Pager[GroupQuotaSubscriptionAllocationClientListResponse]`
- New function `NewGroupQuotaSubscriptionAllocationRequestClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaSubscriptionAllocationRequestClient, error)`
- New function `*GroupQuotaSubscriptionAllocationRequestClient.BeginCreateOrUpdate(context.Context, string, string, string, string, AllocationRequestStatus, *GroupQuotaSubscriptionAllocationRequestClientBeginCreateOrUpdateOptions) (*runtime.Poller[GroupQuotaSubscriptionAllocationRequestClientCreateOrUpdateResponse], error)`
- New function `*GroupQuotaSubscriptionAllocationRequestClient.Get(context.Context, string, string, string, *GroupQuotaSubscriptionAllocationRequestClientGetOptions) (GroupQuotaSubscriptionAllocationRequestClientGetResponse, error)`
- New function `*GroupQuotaSubscriptionAllocationRequestClient.NewListPager(string, string, string, string, *GroupQuotaSubscriptionAllocationRequestClientListOptions) *runtime.Pager[GroupQuotaSubscriptionAllocationRequestClientListResponse]`
- New function `*GroupQuotaSubscriptionAllocationRequestClient.BeginUpdate(context.Context, string, string, string, string, AllocationRequestStatus, *GroupQuotaSubscriptionAllocationRequestClientBeginUpdateOptions) (*runtime.Poller[GroupQuotaSubscriptionAllocationRequestClientUpdateResponse], error)`
- New function `NewGroupQuotaSubscriptionRequestsClient(azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaSubscriptionRequestsClient, error)`
- New function `*GroupQuotaSubscriptionRequestsClient.Get(context.Context, string, string, string, *GroupQuotaSubscriptionRequestsClientGetOptions) (GroupQuotaSubscriptionRequestsClientGetResponse, error)`
- New function `*GroupQuotaSubscriptionRequestsClient.NewListPager(string, string, *GroupQuotaSubscriptionRequestsClientListOptions) *runtime.Pager[GroupQuotaSubscriptionRequestsClientListResponse]`
- New function `NewGroupQuotaSubscriptionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaSubscriptionsClient, error)`
- New function `*GroupQuotaSubscriptionsClient.BeginCreateOrUpdate(context.Context, string, string, *GroupQuotaSubscriptionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[GroupQuotaSubscriptionsClientCreateOrUpdateResponse], error)`
- New function `*GroupQuotaSubscriptionsClient.BeginDelete(context.Context, string, string, *GroupQuotaSubscriptionsClientBeginDeleteOptions) (*runtime.Poller[GroupQuotaSubscriptionsClientDeleteResponse], error)`
- New function `*GroupQuotaSubscriptionsClient.Get(context.Context, string, string, *GroupQuotaSubscriptionsClientGetOptions) (GroupQuotaSubscriptionsClientGetResponse, error)`
- New function `*GroupQuotaSubscriptionsClient.NewListPager(string, string, *GroupQuotaSubscriptionsClientListOptions) *runtime.Pager[GroupQuotaSubscriptionsClientListResponse]`
- New function `*GroupQuotaSubscriptionsClient.BeginUpdate(context.Context, string, string, *GroupQuotaSubscriptionsClientBeginUpdateOptions) (*runtime.Poller[GroupQuotaSubscriptionsClientUpdateResponse], error)`
- New function `NewGroupQuotaUsagesClient(azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotaUsagesClient, error)`
- New function `*GroupQuotaUsagesClient.NewListPager(string, string, string, string, *GroupQuotaUsagesClientListOptions) *runtime.Pager[GroupQuotaUsagesClientListResponse]`
- New function `NewGroupQuotasClient(azcore.TokenCredential, *arm.ClientOptions) (*GroupQuotasClient, error)`
- New function `*GroupQuotasClient.BeginCreateOrUpdate(context.Context, string, string, *GroupQuotasClientBeginCreateOrUpdateOptions) (*runtime.Poller[GroupQuotasClientCreateOrUpdateResponse], error)`
- New function `*GroupQuotasClient.BeginDelete(context.Context, string, string, *GroupQuotasClientBeginDeleteOptions) (*runtime.Poller[GroupQuotasClientDeleteResponse], error)`
- New function `*GroupQuotasClient.Get(context.Context, string, string, *GroupQuotasClientGetOptions) (GroupQuotasClientGetResponse, error)`
- New function `*GroupQuotasClient.NewListPager(string, *GroupQuotasClientListOptions) *runtime.Pager[GroupQuotasClientListResponse]`
- New function `*GroupQuotasClient.BeginUpdate(context.Context, string, string, *GroupQuotasClientBeginUpdateOptions) (*runtime.Poller[GroupQuotasClientUpdateResponse], error)`
- New struct `AdditionalAttributes`
- New struct `AdditionalAttributesPatch`
- New struct `AllocatedQuotaToSubscriptionList`
- New struct `AllocatedToSubscription`
- New struct `AllocationRequestBase`
- New struct `AllocationRequestBaseProperties`
- New struct `AllocationRequestBasePropertiesName`
- New struct `AllocationRequestStatus`
- New struct `AllocationRequestStatusList`
- New struct `AllocationRequestStatusProperties`
- New struct `BillingAccountID`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New struct `GroupQuotaDetails`
- New struct `GroupQuotaDetailsName`
- New struct `GroupQuotaLimit`
- New struct `GroupQuotaLimitList`
- New struct `GroupQuotaList`
- New struct `GroupQuotaRequestBase`
- New struct `GroupQuotaRequestBaseProperties`
- New struct `GroupQuotaRequestBasePropertiesName`
- New struct `GroupQuotaSubscriptionID`
- New struct `GroupQuotaSubscriptionIDList`
- New struct `GroupQuotaSubscriptionIDProperties`
- New struct `GroupQuotaSubscriptionRequestStatus`
- New struct `GroupQuotaSubscriptionRequestStatusList`
- New struct `GroupQuotaSubscriptionRequestStatusProperties`
- New struct `GroupQuotaUsagesBase`
- New struct `GroupQuotaUsagesBaseName`
- New struct `GroupQuotasEnforcementListResponse`
- New struct `GroupQuotasEnforcementResponse`
- New struct `GroupQuotasEnforcementResponseProperties`
- New struct `GroupQuotasEntity`
- New struct `GroupQuotasEntityBase`
- New struct `GroupQuotasEntityBasePatch`
- New struct `GroupQuotasEntityPatch`
- New struct `GroupingID`
- New struct `LROResponse`
- New struct `LROResponseProperties`
- New struct `ProxyResource`
- New struct `Resource`
- New struct `ResourceBaseRequest`
- New struct `ResourceUsageList`
- New struct `ResourceUsages`
- New struct `SubmittedResourceRequestStatus`
- New struct `SubmittedResourceRequestStatusList`
- New struct `SubmittedResourceRequestStatusProperties`
- New struct `SubscriptionGroupQuotaAssignment`
- New struct `SubscriptionQuotaAllocationRequestList`
- New struct `SubscriptionQuotaAllocations`
- New struct `SubscriptionQuotaAllocationsList`
- New struct `SubscriptionQuotaAllocationsStatusList`
- New struct `SubscriptionQuotaDetails`
- New struct `SubscriptionQuotaDetailsName`
- New struct `SystemData`


## 1.0.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.

### Other Changes

- Release stable version.


## 0.6.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.6.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/quota/armquota` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
