# Release History

## 0.4.0 (2022-04-18)
### Breaking Changes

- Function `*ProfilesClient.ListBySubscription` has been removed
- Function `*ProfilesClient.ListByResourceGroup` has been removed

### Features Added

- New function `*ProfilesClient.NewListBySubscriptionPager(*ProfilesClientListBySubscriptionOptions) *runtime.Pager[ProfilesClientListBySubscriptionResponse]`
- New function `*ProfilesClient.NewListByResourceGroupPager(string, *ProfilesClientListByResourceGroupOptions) *runtime.Pager[ProfilesClientListByResourceGroupResponse]`


## 0.3.0 (2022-04-13)
### Breaking Changes

- Function `*ProfilesClient.ListByResourceGroup` parameter(s) have been changed from `(context.Context, string, *ProfilesClientListByResourceGroupOptions)` to `(string, *ProfilesClientListByResourceGroupOptions)`
- Function `*ProfilesClient.ListByResourceGroup` return value(s) have been changed from `(ProfilesClientListByResourceGroupResponse, error)` to `(*runtime.Pager[ProfilesClientListByResourceGroupResponse])`
- Function `*ProfilesClient.ListBySubscription` parameter(s) have been changed from `(context.Context, *ProfilesClientListBySubscriptionOptions)` to `(*ProfilesClientListBySubscriptionOptions)`
- Function `*ProfilesClient.ListBySubscription` return value(s) have been changed from `(ProfilesClientListBySubscriptionResponse, error)` to `(*runtime.Pager[ProfilesClientListBySubscriptionResponse])`
- Function `NewUserMetricsKeysClient` return value(s) have been changed from `(*UserMetricsKeysClient)` to `(*UserMetricsKeysClient, error)`
- Function `NewGeographicHierarchiesClient` return value(s) have been changed from `(*GeographicHierarchiesClient)` to `(*GeographicHierarchiesClient, error)`
- Function `*HeatMapClient.Get` parameter(s) have been changed from `(context.Context, string, string, Enum8, *HeatMapClientGetOptions)` to `(context.Context, string, string, *HeatMapClientGetOptions)`
- Function `*EndpointsClient.Update` parameter(s) have been changed from `(context.Context, string, string, string, string, Endpoint, *EndpointsClientUpdateOptions)` to `(context.Context, string, string, EndpointType, string, Endpoint, *EndpointsClientUpdateOptions)`
- Function `*EndpointsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, string, Endpoint, *EndpointsClientCreateOrUpdateOptions)` to `(context.Context, string, string, EndpointType, string, Endpoint, *EndpointsClientCreateOrUpdateOptions)`
- Function `*EndpointsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, string, string, *EndpointsClientDeleteOptions)` to `(context.Context, string, string, EndpointType, string, *EndpointsClientDeleteOptions)`
- Function `NewEndpointsClient` return value(s) have been changed from `(*EndpointsClient)` to `(*EndpointsClient, error)`
- Function `NewProfilesClient` return value(s) have been changed from `(*ProfilesClient)` to `(*ProfilesClient, error)`
- Function `*EndpointsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, string, *EndpointsClientGetOptions)` to `(context.Context, string, string, EndpointType, string, *EndpointsClientGetOptions)`
- Function `NewHeatMapClient` return value(s) have been changed from `(*HeatMapClient)` to `(*HeatMapClient, error)`
- Const `Enum8Default` has been removed
- Function `EndpointMonitorStatus.ToPtr` has been removed
- Function `MonitorProtocol.ToPtr` has been removed
- Function `EndpointStatus.ToPtr` has been removed
- Function `TrafficRoutingMethod.ToPtr` has been removed
- Function `ProfileStatus.ToPtr` has been removed
- Function `ProfileMonitorStatus.ToPtr` has been removed
- Function `TrafficViewEnrollmentStatus.ToPtr` has been removed
- Function `Enum8.ToPtr` has been removed
- Function `AllowedEndpointRecordType.ToPtr` has been removed
- Function `PossibleEnum8Values` has been removed
- Struct `EndpointsClientCreateOrUpdateResult` has been removed
- Struct `EndpointsClientDeleteResult` has been removed
- Struct `EndpointsClientGetResult` has been removed
- Struct `EndpointsClientUpdateResult` has been removed
- Struct `GeographicHierarchiesClientGetDefaultResult` has been removed
- Struct `HeatMapClientGetResult` has been removed
- Struct `ProfilesClientCheckTrafficManagerRelativeDNSNameAvailabilityResult` has been removed
- Struct `ProfilesClientCreateOrUpdateResult` has been removed
- Struct `ProfilesClientDeleteResult` has been removed
- Struct `ProfilesClientGetResult` has been removed
- Struct `ProfilesClientListByResourceGroupResult` has been removed
- Struct `ProfilesClientListBySubscriptionResult` has been removed
- Struct `ProfilesClientUpdateResult` has been removed
- Struct `UserMetricsKeysClientCreateOrUpdateResult` has been removed
- Struct `UserMetricsKeysClientDeleteResult` has been removed
- Struct `UserMetricsKeysClientGetResult` has been removed
- Field `UserMetricsKeysClientGetResult` of struct `UserMetricsKeysClientGetResponse` has been removed
- Field `RawResponse` of struct `UserMetricsKeysClientGetResponse` has been removed
- Field `EndpointsClientUpdateResult` of struct `EndpointsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `EndpointsClientUpdateResponse` has been removed
- Field `EndpointsClientDeleteResult` of struct `EndpointsClientDeleteResponse` has been removed
- Field `RawResponse` of struct `EndpointsClientDeleteResponse` has been removed
- Field `EndpointsClientCreateOrUpdateResult` of struct `EndpointsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `EndpointsClientCreateOrUpdateResponse` has been removed
- Field `ProfilesClientGetResult` of struct `ProfilesClientGetResponse` has been removed
- Field `RawResponse` of struct `ProfilesClientGetResponse` has been removed
- Field `ProfilesClientListByResourceGroupResult` of struct `ProfilesClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `ProfilesClientListByResourceGroupResponse` has been removed
- Field `HeatMapClientGetResult` of struct `HeatMapClientGetResponse` has been removed
- Field `RawResponse` of struct `HeatMapClientGetResponse` has been removed
- Field `ProfilesClientUpdateResult` of struct `ProfilesClientUpdateResponse` has been removed
- Field `RawResponse` of struct `ProfilesClientUpdateResponse` has been removed
- Field `UserMetricsKeysClientCreateOrUpdateResult` of struct `UserMetricsKeysClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `UserMetricsKeysClientCreateOrUpdateResponse` has been removed
- Field `EndpointsClientGetResult` of struct `EndpointsClientGetResponse` has been removed
- Field `RawResponse` of struct `EndpointsClientGetResponse` has been removed
- Field `ProfilesClientDeleteResult` of struct `ProfilesClientDeleteResponse` has been removed
- Field `RawResponse` of struct `ProfilesClientDeleteResponse` has been removed
- Field `ProfilesClientCreateOrUpdateResult` of struct `ProfilesClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `ProfilesClientCreateOrUpdateResponse` has been removed
- Field `GeographicHierarchiesClientGetDefaultResult` of struct `GeographicHierarchiesClientGetDefaultResponse` has been removed
- Field `RawResponse` of struct `GeographicHierarchiesClientGetDefaultResponse` has been removed
- Field `ProfilesClientCheckTrafficManagerRelativeDNSNameAvailabilityResult` of struct `ProfilesClientCheckTrafficManagerRelativeDNSNameAvailabilityResponse` has been removed
- Field `RawResponse` of struct `ProfilesClientCheckTrafficManagerRelativeDNSNameAvailabilityResponse` has been removed
- Field `UserMetricsKeysClientDeleteResult` of struct `UserMetricsKeysClientDeleteResponse` has been removed
- Field `RawResponse` of struct `UserMetricsKeysClientDeleteResponse` has been removed
- Field `ProfilesClientListBySubscriptionResult` of struct `ProfilesClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `ProfilesClientListBySubscriptionResponse` has been removed

### Features Added

- New const `EndpointTypeExternalEndpoints`
- New const `EndpointTypeNestedEndpoints`
- New const `EndpointTypeAzureEndpoints`
- New function `PossibleEndpointTypeValues() []EndpointType`
- New anonymous field `ProfileListResult` in struct `ProfilesClientListByResourceGroupResponse`
- New anonymous field `ProfileListResult` in struct `ProfilesClientListBySubscriptionResponse`
- New anonymous field `Profile` in struct `ProfilesClientCreateOrUpdateResponse`
- New anonymous field `HeatMapModel` in struct `HeatMapClientGetResponse`
- New anonymous field `DeleteOperationResult` in struct `EndpointsClientDeleteResponse`
- New anonymous field `Endpoint` in struct `EndpointsClientGetResponse`
- New anonymous field `GeographicHierarchy` in struct `GeographicHierarchiesClientGetDefaultResponse`
- New anonymous field `NameAvailability` in struct `ProfilesClientCheckTrafficManagerRelativeDNSNameAvailabilityResponse`
- New anonymous field `Endpoint` in struct `EndpointsClientUpdateResponse`
- New anonymous field `Endpoint` in struct `EndpointsClientCreateOrUpdateResponse`
- New anonymous field `DeleteOperationResult` in struct `ProfilesClientDeleteResponse`
- New anonymous field `UserMetricsModel` in struct `UserMetricsKeysClientGetResponse`
- New anonymous field `DeleteOperationResult` in struct `UserMetricsKeysClientDeleteResponse`
- New anonymous field `Profile` in struct `ProfilesClientUpdateResponse`
- New anonymous field `UserMetricsModel` in struct `UserMetricsKeysClientCreateOrUpdateResponse`
- New anonymous field `Profile` in struct `ProfilesClientGetResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*GeographicHierarchiesClient.GetDefault` parameter(s) have been changed from `(context.Context, *GeographicHierarchiesGetDefaultOptions)` to `(context.Context, *GeographicHierarchiesClientGetDefaultOptions)`
- Function `*GeographicHierarchiesClient.GetDefault` return value(s) have been changed from `(GeographicHierarchiesGetDefaultResponse, error)` to `(GeographicHierarchiesClientGetDefaultResponse, error)`
- Function `*ProfilesClient.CheckTrafficManagerRelativeDNSNameAvailability` parameter(s) have been changed from `(context.Context, CheckTrafficManagerRelativeDNSNameAvailabilityParameters, *ProfilesCheckTrafficManagerRelativeDNSNameAvailabilityOptions)` to `(context.Context, CheckTrafficManagerRelativeDNSNameAvailabilityParameters, *ProfilesClientCheckTrafficManagerRelativeDNSNameAvailabilityOptions)`
- Function `*ProfilesClient.CheckTrafficManagerRelativeDNSNameAvailability` return value(s) have been changed from `(ProfilesCheckTrafficManagerRelativeDNSNameAvailabilityResponse, error)` to `(ProfilesClientCheckTrafficManagerRelativeDNSNameAvailabilityResponse, error)`
- Function `*ProfilesClient.ListByResourceGroup` parameter(s) have been changed from `(context.Context, string, *ProfilesListByResourceGroupOptions)` to `(context.Context, string, *ProfilesClientListByResourceGroupOptions)`
- Function `*ProfilesClient.ListByResourceGroup` return value(s) have been changed from `(ProfilesListByResourceGroupResponse, error)` to `(ProfilesClientListByResourceGroupResponse, error)`
- Function `*EndpointsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, string, *EndpointsGetOptions)` to `(context.Context, string, string, string, string, *EndpointsClientGetOptions)`
- Function `*EndpointsClient.Get` return value(s) have been changed from `(EndpointsGetResponse, error)` to `(EndpointsClientGetResponse, error)`
- Function `*ProfilesClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, Profile, *ProfilesCreateOrUpdateOptions)` to `(context.Context, string, string, Profile, *ProfilesClientCreateOrUpdateOptions)`
- Function `*ProfilesClient.CreateOrUpdate` return value(s) have been changed from `(ProfilesCreateOrUpdateResponse, error)` to `(ProfilesClientCreateOrUpdateResponse, error)`
- Function `*ProfilesClient.ListBySubscription` parameter(s) have been changed from `(context.Context, *ProfilesListBySubscriptionOptions)` to `(context.Context, *ProfilesClientListBySubscriptionOptions)`
- Function `*ProfilesClient.ListBySubscription` return value(s) have been changed from `(ProfilesListBySubscriptionResponse, error)` to `(ProfilesClientListBySubscriptionResponse, error)`
- Function `*EndpointsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, string, Endpoint, *EndpointsCreateOrUpdateOptions)` to `(context.Context, string, string, string, string, Endpoint, *EndpointsClientCreateOrUpdateOptions)`
- Function `*EndpointsClient.CreateOrUpdate` return value(s) have been changed from `(EndpointsCreateOrUpdateResponse, error)` to `(EndpointsClientCreateOrUpdateResponse, error)`
- Function `*ProfilesClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *ProfilesDeleteOptions)` to `(context.Context, string, string, *ProfilesClientDeleteOptions)`
- Function `*ProfilesClient.Delete` return value(s) have been changed from `(ProfilesDeleteResponse, error)` to `(ProfilesClientDeleteResponse, error)`
- Function `*ProfilesClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ProfilesGetOptions)` to `(context.Context, string, string, *ProfilesClientGetOptions)`
- Function `*ProfilesClient.Get` return value(s) have been changed from `(ProfilesGetResponse, error)` to `(ProfilesClientGetResponse, error)`
- Function `*ProfilesClient.Update` parameter(s) have been changed from `(context.Context, string, string, Profile, *ProfilesUpdateOptions)` to `(context.Context, string, string, Profile, *ProfilesClientUpdateOptions)`
- Function `*ProfilesClient.Update` return value(s) have been changed from `(ProfilesUpdateResponse, error)` to `(ProfilesClientUpdateResponse, error)`
- Function `*EndpointsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, string, string, *EndpointsDeleteOptions)` to `(context.Context, string, string, string, string, *EndpointsClientDeleteOptions)`
- Function `*EndpointsClient.Delete` return value(s) have been changed from `(EndpointsDeleteResponse, error)` to `(EndpointsClientDeleteResponse, error)`
- Function `*EndpointsClient.Update` parameter(s) have been changed from `(context.Context, string, string, string, string, Endpoint, *EndpointsUpdateOptions)` to `(context.Context, string, string, string, string, Endpoint, *EndpointsClientUpdateOptions)`
- Function `*EndpointsClient.Update` return value(s) have been changed from `(EndpointsUpdateResponse, error)` to `(EndpointsClientUpdateResponse, error)`
- Function `*HeatMapClient.Get` parameter(s) have been changed from `(context.Context, string, string, Enum8, *HeatMapGetOptions)` to `(context.Context, string, string, Enum8, *HeatMapClientGetOptions)`
- Function `*HeatMapClient.Get` return value(s) have been changed from `(HeatMapGetResponse, error)` to `(HeatMapClientGetResponse, error)`
- Function `CloudError.Error` has been removed
- Function `*TrafficManagerUserMetricsKeysClient.CreateOrUpdate` has been removed
- Function `*TrafficManagerUserMetricsKeysClient.Get` has been removed
- Function `HeatMapModel.MarshalJSON` has been removed
- Function `*TrafficManagerUserMetricsKeysClient.Delete` has been removed
- Function `NewTrafficManagerUserMetricsKeysClient` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `TrafficManagerGeographicHierarchy.MarshalJSON` has been removed
- Function `UserMetricsModel.MarshalJSON` has been removed
- Struct `EndpointsCreateOrUpdateOptions` has been removed
- Struct `EndpointsCreateOrUpdateResponse` has been removed
- Struct `EndpointsCreateOrUpdateResult` has been removed
- Struct `EndpointsDeleteOptions` has been removed
- Struct `EndpointsDeleteResponse` has been removed
- Struct `EndpointsDeleteResult` has been removed
- Struct `EndpointsGetOptions` has been removed
- Struct `EndpointsGetResponse` has been removed
- Struct `EndpointsGetResult` has been removed
- Struct `EndpointsUpdateOptions` has been removed
- Struct `EndpointsUpdateResponse` has been removed
- Struct `EndpointsUpdateResult` has been removed
- Struct `GeographicHierarchiesGetDefaultOptions` has been removed
- Struct `GeographicHierarchiesGetDefaultResponse` has been removed
- Struct `GeographicHierarchiesGetDefaultResult` has been removed
- Struct `HeatMapGetOptions` has been removed
- Struct `HeatMapGetResponse` has been removed
- Struct `HeatMapGetResult` has been removed
- Struct `ProfilesCheckTrafficManagerRelativeDNSNameAvailabilityOptions` has been removed
- Struct `ProfilesCheckTrafficManagerRelativeDNSNameAvailabilityResponse` has been removed
- Struct `ProfilesCheckTrafficManagerRelativeDNSNameAvailabilityResult` has been removed
- Struct `ProfilesCreateOrUpdateOptions` has been removed
- Struct `ProfilesCreateOrUpdateResponse` has been removed
- Struct `ProfilesCreateOrUpdateResult` has been removed
- Struct `ProfilesDeleteOptions` has been removed
- Struct `ProfilesDeleteResponse` has been removed
- Struct `ProfilesDeleteResult` has been removed
- Struct `ProfilesGetOptions` has been removed
- Struct `ProfilesGetResponse` has been removed
- Struct `ProfilesGetResult` has been removed
- Struct `ProfilesListByResourceGroupOptions` has been removed
- Struct `ProfilesListByResourceGroupResponse` has been removed
- Struct `ProfilesListByResourceGroupResult` has been removed
- Struct `ProfilesListBySubscriptionOptions` has been removed
- Struct `ProfilesListBySubscriptionResponse` has been removed
- Struct `ProfilesListBySubscriptionResult` has been removed
- Struct `ProfilesUpdateOptions` has been removed
- Struct `ProfilesUpdateResponse` has been removed
- Struct `ProfilesUpdateResult` has been removed
- Struct `TrafficManagerGeographicHierarchy` has been removed
- Struct `TrafficManagerNameAvailability` has been removed
- Struct `TrafficManagerUserMetricsKeysClient` has been removed
- Struct `TrafficManagerUserMetricsKeysCreateOrUpdateOptions` has been removed
- Struct `TrafficManagerUserMetricsKeysCreateOrUpdateResponse` has been removed
- Struct `TrafficManagerUserMetricsKeysCreateOrUpdateResult` has been removed
- Struct `TrafficManagerUserMetricsKeysDeleteOptions` has been removed
- Struct `TrafficManagerUserMetricsKeysDeleteResponse` has been removed
- Struct `TrafficManagerUserMetricsKeysDeleteResult` has been removed
- Struct `TrafficManagerUserMetricsKeysGetOptions` has been removed
- Struct `TrafficManagerUserMetricsKeysGetResponse` has been removed
- Struct `TrafficManagerUserMetricsKeysGetResult` has been removed
- Field `ProxyResource` of struct `Endpoint` has been removed
- Field `ProxyResource` of struct `UserMetricsModel` has been removed
- Field `ProxyResource` of struct `HeatMapModel` has been removed
- Field `Resource` of struct `ProxyResource` has been removed
- Field `TrackedResource` of struct `Profile` has been removed
- Field `Resource` of struct `TrackedResource` has been removed
- Field `InnerError` of struct `CloudError` has been removed

### Features Added

- New function `*UserMetricsKeysClient.Delete(context.Context, *UserMetricsKeysClientDeleteOptions) (UserMetricsKeysClientDeleteResponse, error)`
- New function `NewUserMetricsKeysClient(string, azcore.TokenCredential, *arm.ClientOptions) *UserMetricsKeysClient`
- New function `*UserMetricsKeysClient.CreateOrUpdate(context.Context, *UserMetricsKeysClientCreateOrUpdateOptions) (UserMetricsKeysClientCreateOrUpdateResponse, error)`
- New function `*UserMetricsKeysClient.Get(context.Context, *UserMetricsKeysClientGetOptions) (UserMetricsKeysClientGetResponse, error)`
- New struct `EndpointsClientCreateOrUpdateOptions`
- New struct `EndpointsClientCreateOrUpdateResponse`
- New struct `EndpointsClientCreateOrUpdateResult`
- New struct `EndpointsClientDeleteOptions`
- New struct `EndpointsClientDeleteResponse`
- New struct `EndpointsClientDeleteResult`
- New struct `EndpointsClientGetOptions`
- New struct `EndpointsClientGetResponse`
- New struct `EndpointsClientGetResult`
- New struct `EndpointsClientUpdateOptions`
- New struct `EndpointsClientUpdateResponse`
- New struct `EndpointsClientUpdateResult`
- New struct `GeographicHierarchiesClientGetDefaultOptions`
- New struct `GeographicHierarchiesClientGetDefaultResponse`
- New struct `GeographicHierarchiesClientGetDefaultResult`
- New struct `GeographicHierarchy`
- New struct `HeatMapClientGetOptions`
- New struct `HeatMapClientGetResponse`
- New struct `HeatMapClientGetResult`
- New struct `NameAvailability`
- New struct `ProfilesClientCheckTrafficManagerRelativeDNSNameAvailabilityOptions`
- New struct `ProfilesClientCheckTrafficManagerRelativeDNSNameAvailabilityResponse`
- New struct `ProfilesClientCheckTrafficManagerRelativeDNSNameAvailabilityResult`
- New struct `ProfilesClientCreateOrUpdateOptions`
- New struct `ProfilesClientCreateOrUpdateResponse`
- New struct `ProfilesClientCreateOrUpdateResult`
- New struct `ProfilesClientDeleteOptions`
- New struct `ProfilesClientDeleteResponse`
- New struct `ProfilesClientDeleteResult`
- New struct `ProfilesClientGetOptions`
- New struct `ProfilesClientGetResponse`
- New struct `ProfilesClientGetResult`
- New struct `ProfilesClientListByResourceGroupOptions`
- New struct `ProfilesClientListByResourceGroupResponse`
- New struct `ProfilesClientListByResourceGroupResult`
- New struct `ProfilesClientListBySubscriptionOptions`
- New struct `ProfilesClientListBySubscriptionResponse`
- New struct `ProfilesClientListBySubscriptionResult`
- New struct `ProfilesClientUpdateOptions`
- New struct `ProfilesClientUpdateResponse`
- New struct `ProfilesClientUpdateResult`
- New struct `UserMetricsKeysClient`
- New struct `UserMetricsKeysClientCreateOrUpdateOptions`
- New struct `UserMetricsKeysClientCreateOrUpdateResponse`
- New struct `UserMetricsKeysClientCreateOrUpdateResult`
- New struct `UserMetricsKeysClientDeleteOptions`
- New struct `UserMetricsKeysClientDeleteResponse`
- New struct `UserMetricsKeysClientDeleteResult`
- New struct `UserMetricsKeysClientGetOptions`
- New struct `UserMetricsKeysClientGetResponse`
- New struct `UserMetricsKeysClientGetResult`
- New field `Type` in struct `HeatMapModel`
- New field `ID` in struct `HeatMapModel`
- New field `Name` in struct `HeatMapModel`
- New field `ID` in struct `Endpoint`
- New field `Name` in struct `Endpoint`
- New field `Type` in struct `Endpoint`
- New field `ID` in struct `TrackedResource`
- New field `Name` in struct `TrackedResource`
- New field `Type` in struct `TrackedResource`
- New field `Type` in struct `Profile`
- New field `ID` in struct `Profile`
- New field `Location` in struct `Profile`
- New field `Name` in struct `Profile`
- New field `Tags` in struct `Profile`
- New field `Type` in struct `ProxyResource`
- New field `ID` in struct `ProxyResource`
- New field `Name` in struct `ProxyResource`
- New field `Type` in struct `UserMetricsModel`
- New field `ID` in struct `UserMetricsModel`
- New field `Name` in struct `UserMetricsModel`
- New field `Error` in struct `CloudError`


## 0.1.0 (2021-12-22)

- Init release.
