# Release History

## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*ServiceClient.ListByResourceGroup` has been removed
- Function `*OperationsClient.List` has been removed
- Function `*ServiceClient.ListBySubscription` has been removed

### Features Added

- New function `*ServiceClient.NewListBySubscriptionPager(*ServiceClientListBySubscriptionOptions) *runtime.Pager[ServiceClientListBySubscriptionResponse]`
- New function `*ServiceClient.NewListByResourceGroupPager(string, *ServiceClientListByResourceGroupOptions) *runtime.Pager[ServiceClientListByResourceGroupResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `NewServiceClient` return value(s) have been changed from `(*ServiceClient)` to `(*ServiceClient, error)`
- Function `*ServiceClient.BeginDelete` return value(s) have been changed from `(ServiceClientDeletePollerResponse, error)` to `(*armruntime.Poller[ServiceClientDeleteResponse], error)`
- Function `*ServiceClient.ListByResourceGroup` return value(s) have been changed from `(*ServiceClientListByResourceGroupPager)` to `(*runtime.Pager[ServiceClientListByResourceGroupResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*ServiceClient.BeginCreateOrUpdate` return value(s) have been changed from `(ServiceClientCreateOrUpdatePollerResponse, error)` to `(*armruntime.Poller[ServiceClientCreateOrUpdateResponse], error)`
- Function `*ServiceClient.ListBySubscription` return value(s) have been changed from `(*ServiceClientListBySubscriptionPager)` to `(*runtime.Pager[ServiceClientListBySubscriptionResponse])`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Function `Origin.ToPtr` has been removed
- Function `*ServiceClientDeletePoller.FinalResponse` has been removed
- Function `*ServiceClientDeletePoller.ResumeToken` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*ServiceClientCreateOrUpdatePoller.ResumeToken` has been removed
- Function `*ServiceClientCreateOrUpdatePoller.Done` has been removed
- Function `ServiceClientCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*ServiceClientCreateOrUpdatePollerResponse.Resume` has been removed
- Function `ActionType.ToPtr` has been removed
- Function `*ServiceClientCreateOrUpdatePoller.Poll` has been removed
- Function `*ServiceClientListByResourceGroupPager.PageResponse` has been removed
- Function `*ServiceClientListBySubscriptionPager.NextPage` has been removed
- Function `*ServiceClientDeletePoller.Done` has been removed
- Function `*ServiceClientListBySubscriptionPager.Err` has been removed
- Function `CreatedByType.ToPtr` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `KeyType.ToPtr` has been removed
- Function `ProvisioningState.ToPtr` has been removed
- Function `*ServiceClientCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*ServiceClientListByResourceGroupPager.NextPage` has been removed
- Function `*ServiceClientListBySubscriptionPager.PageResponse` has been removed
- Function `*ServiceClientDeletePoller.Poll` has been removed
- Function `*ServiceClientDeletePollerResponse.Resume` has been removed
- Function `ServiceClientDeletePollerResponse.PollUntilDone` has been removed
- Function `*ServiceClientListByResourceGroupPager.Err` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `ServiceClientCheckNameAvailabilityResult` has been removed
- Struct `ServiceClientCreateOrUpdatePoller` has been removed
- Struct `ServiceClientCreateOrUpdatePollerResponse` has been removed
- Struct `ServiceClientCreateOrUpdateResult` has been removed
- Struct `ServiceClientDeletePoller` has been removed
- Struct `ServiceClientDeletePollerResponse` has been removed
- Struct `ServiceClientGetResult` has been removed
- Struct `ServiceClientLinkNotificationHubResult` has been removed
- Struct `ServiceClientListByResourceGroupPager` has been removed
- Struct `ServiceClientListByResourceGroupResult` has been removed
- Struct `ServiceClientListBySubscriptionPager` has been removed
- Struct `ServiceClientListBySubscriptionResult` has been removed
- Struct `ServiceClientListKeysResult` has been removed
- Struct `ServiceClientRegenerateKeyResult` has been removed
- Struct `ServiceClientUpdateResult` has been removed
- Field `ServiceClientCheckNameAvailabilityResult` of struct `ServiceClientCheckNameAvailabilityResponse` has been removed
- Field `RawResponse` of struct `ServiceClientCheckNameAvailabilityResponse` has been removed
- Field `ServiceClientListByResourceGroupResult` of struct `ServiceClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `ServiceClientListByResourceGroupResponse` has been removed
- Field `ServiceClientRegenerateKeyResult` of struct `ServiceClientRegenerateKeyResponse` has been removed
- Field `RawResponse` of struct `ServiceClientRegenerateKeyResponse` has been removed
- Field `ServiceClientListBySubscriptionResult` of struct `ServiceClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `ServiceClientListBySubscriptionResponse` has been removed
- Field `ServiceClientCreateOrUpdateResult` of struct `ServiceClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `ServiceClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `ServiceClientDeleteResponse` has been removed
- Field `ServiceClientLinkNotificationHubResult` of struct `ServiceClientLinkNotificationHubResponse` has been removed
- Field `RawResponse` of struct `ServiceClientLinkNotificationHubResponse` has been removed
- Field `ServiceClientGetResult` of struct `ServiceClientGetResponse` has been removed
- Field `RawResponse` of struct `ServiceClientGetResponse` has been removed
- Field `ServiceClientListKeysResult` of struct `ServiceClientListKeysResponse` has been removed
- Field `RawResponse` of struct `ServiceClientListKeysResponse` has been removed
- Field `ServiceClientUpdateResult` of struct `ServiceClientUpdateResponse` has been removed
- Field `RawResponse` of struct `ServiceClientUpdateResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed

### Features Added

- New anonymous field `ServiceResourceList` in struct `ServiceClientListByResourceGroupResponse`
- New anonymous field `ServiceResource` in struct `ServiceClientGetResponse`
- New anonymous field `LinkedNotificationHub` in struct `ServiceClientLinkNotificationHubResponse`
- New anonymous field `ServiceResource` in struct `ServiceClientCreateOrUpdateResponse`
- New field `ResumeToken` in struct `ServiceClientBeginCreateOrUpdateOptions`
- New anonymous field `ServiceKeys` in struct `ServiceClientListKeysResponse`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `ServiceResourceList` in struct `ServiceClientListBySubscriptionResponse`
- New anonymous field `ServiceKeys` in struct `ServiceClientRegenerateKeyResponse`
- New anonymous field `ServiceResource` in struct `ServiceClientUpdateResponse`
- New field `ResumeToken` in struct `ServiceClientBeginDeleteOptions`
- New anonymous field `NameAvailability` in struct `ServiceClientCheckNameAvailabilityResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `NewCommunicationServiceClient` has been removed
- Function `CommunicationServiceResourceList.MarshalJSON` has been removed
- Function `*CommunicationServiceClient.CheckNameAvailability` has been removed
- Function `*CommunicationServiceClient.LinkNotificationHub` has been removed
- Function `CommunicationServiceResource.MarshalJSON` has been removed
- Function `*CommunicationServiceClient.ListBySubscription` has been removed
- Function `*CommunicationServiceDeletePollerResponse.Resume` has been removed
- Function `*CommunicationServiceClient.ListByResourceGroup` has been removed
- Function `*CommunicationServiceClient.BeginDelete` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*CommunicationServiceCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*CommunicationServiceCreateOrUpdatePoller.Poll` has been removed
- Function `*CommunicationServiceListByResourceGroupPager.Err` has been removed
- Function `*CommunicationServiceClient.Update` has been removed
- Function `*CommunicationServiceDeletePoller.Done` has been removed
- Function `*CommunicationServiceClient.BeginCreateOrUpdate` has been removed
- Function `*CommunicationServiceCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*CommunicationServiceListByResourceGroupPager.NextPage` has been removed
- Function `*CommunicationServiceCreateOrUpdatePoller.Done` has been removed
- Function `*CommunicationServiceListBySubscriptionPager.Err` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `*CommunicationServiceListBySubscriptionPager.NextPage` has been removed
- Function `*CommunicationServiceDeletePoller.Poll` has been removed
- Function `*CommunicationServiceDeletePoller.ResumeToken` has been removed
- Function `*CommunicationServiceListBySubscriptionPager.PageResponse` has been removed
- Function `LocationResource.MarshalJSON` has been removed
- Function `CommunicationServiceCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `CommunicationServiceDeletePollerResponse.PollUntilDone` has been removed
- Function `*CommunicationServiceCreateOrUpdatePoller.ResumeToken` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*CommunicationServiceClient.ListKeys` has been removed
- Function `*CommunicationServiceListByResourceGroupPager.PageResponse` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*CommunicationServiceDeletePoller.FinalResponse` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*CommunicationServiceClient.Get` has been removed
- Function `*CommunicationServiceClient.RegenerateKey` has been removed
- Struct `CommunicationServiceBeginCreateOrUpdateOptions` has been removed
- Struct `CommunicationServiceBeginDeleteOptions` has been removed
- Struct `CommunicationServiceCheckNameAvailabilityOptions` has been removed
- Struct `CommunicationServiceCheckNameAvailabilityResponse` has been removed
- Struct `CommunicationServiceCheckNameAvailabilityResult` has been removed
- Struct `CommunicationServiceClient` has been removed
- Struct `CommunicationServiceCreateOrUpdatePoller` has been removed
- Struct `CommunicationServiceCreateOrUpdatePollerResponse` has been removed
- Struct `CommunicationServiceCreateOrUpdateResponse` has been removed
- Struct `CommunicationServiceCreateOrUpdateResult` has been removed
- Struct `CommunicationServiceDeletePoller` has been removed
- Struct `CommunicationServiceDeletePollerResponse` has been removed
- Struct `CommunicationServiceDeleteResponse` has been removed
- Struct `CommunicationServiceGetOptions` has been removed
- Struct `CommunicationServiceGetResponse` has been removed
- Struct `CommunicationServiceGetResult` has been removed
- Struct `CommunicationServiceKeys` has been removed
- Struct `CommunicationServiceLinkNotificationHubOptions` has been removed
- Struct `CommunicationServiceLinkNotificationHubResponse` has been removed
- Struct `CommunicationServiceLinkNotificationHubResult` has been removed
- Struct `CommunicationServiceListByResourceGroupOptions` has been removed
- Struct `CommunicationServiceListByResourceGroupPager` has been removed
- Struct `CommunicationServiceListByResourceGroupResponse` has been removed
- Struct `CommunicationServiceListByResourceGroupResult` has been removed
- Struct `CommunicationServiceListBySubscriptionOptions` has been removed
- Struct `CommunicationServiceListBySubscriptionPager` has been removed
- Struct `CommunicationServiceListBySubscriptionResponse` has been removed
- Struct `CommunicationServiceListBySubscriptionResult` has been removed
- Struct `CommunicationServiceListKeysOptions` has been removed
- Struct `CommunicationServiceListKeysResponse` has been removed
- Struct `CommunicationServiceListKeysResult` has been removed
- Struct `CommunicationServiceProperties` has been removed
- Struct `CommunicationServiceRegenerateKeyOptions` has been removed
- Struct `CommunicationServiceRegenerateKeyResponse` has been removed
- Struct `CommunicationServiceRegenerateKeyResult` has been removed
- Struct `CommunicationServiceResource` has been removed
- Struct `CommunicationServiceResourceList` has been removed
- Struct `CommunicationServiceUpdateOptions` has been removed
- Struct `CommunicationServiceUpdateResponse` has been removed
- Struct `CommunicationServiceUpdateResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `*ServiceClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*ServiceClientListBySubscriptionPager.PageResponse() ServiceClientListBySubscriptionResponse`
- New function `*ServiceClient.Get(context.Context, string, string, *ServiceClientGetOptions) (ServiceClientGetResponse, error)`
- New function `*ServiceClientDeletePoller.FinalResponse(context.Context) (ServiceClientDeleteResponse, error)`
- New function `NewServiceClient(string, azcore.TokenCredential, *arm.ClientOptions) *ServiceClient`
- New function `*ServiceClientListBySubscriptionPager.Err() error`
- New function `*ServiceClient.RegenerateKey(context.Context, string, string, RegenerateKeyParameters, *ServiceClientRegenerateKeyOptions) (ServiceClientRegenerateKeyResponse, error)`
- New function `*ServiceClient.Update(context.Context, string, string, *ServiceClientUpdateOptions) (ServiceClientUpdateResponse, error)`
- New function `*ServiceClient.CheckNameAvailability(context.Context, *ServiceClientCheckNameAvailabilityOptions) (ServiceClientCheckNameAvailabilityResponse, error)`
- New function `*ServiceClient.ListByResourceGroup(string, *ServiceClientListByResourceGroupOptions) *ServiceClientListByResourceGroupPager`
- New function `*ServiceClient.ListKeys(context.Context, string, string, *ServiceClientListKeysOptions) (ServiceClientListKeysResponse, error)`
- New function `*ServiceClient.BeginDelete(context.Context, string, string, *ServiceClientBeginDeleteOptions) (ServiceClientDeletePollerResponse, error)`
- New function `ServiceClientCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (ServiceClientCreateOrUpdateResponse, error)`
- New function `*ServiceClientListByResourceGroupPager.Err() error`
- New function `*ServiceClientDeletePoller.ResumeToken() (string, error)`
- New function `*ServiceClient.BeginCreateOrUpdate(context.Context, string, string, *ServiceClientBeginCreateOrUpdateOptions) (ServiceClientCreateOrUpdatePollerResponse, error)`
- New function `ServiceClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (ServiceClientDeleteResponse, error)`
- New function `*ServiceClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `ServiceResourceList.MarshalJSON() ([]byte, error)`
- New function `*OperationsClientListPager.Err() error`
- New function `*ServiceClientCreateOrUpdatePoller.Done() bool`
- New function `*ServiceClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*ServiceClient.LinkNotificationHub(context.Context, string, string, *ServiceClientLinkNotificationHubOptions) (ServiceClientLinkNotificationHubResponse, error)`
- New function `*ServiceClientCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*ServiceClientCreateOrUpdatePollerResponse.Resume(context.Context, *ServiceClient, string) error`
- New function `*ServiceClient.ListBySubscription(*ServiceClientListBySubscriptionOptions) *ServiceClientListBySubscriptionPager`
- New function `*ServiceClientDeletePollerResponse.Resume(context.Context, *ServiceClient, string) error`
- New function `*ServiceClientCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*ServiceClientCreateOrUpdatePoller.FinalResponse(context.Context) (ServiceClientCreateOrUpdateResponse, error)`
- New function `ServiceResource.MarshalJSON() ([]byte, error)`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*ServiceClientListByResourceGroupPager.PageResponse() ServiceClientListByResourceGroupResponse`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*ServiceClientDeletePoller.Done() bool`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `ServiceClient`
- New struct `ServiceClientBeginCreateOrUpdateOptions`
- New struct `ServiceClientBeginDeleteOptions`
- New struct `ServiceClientCheckNameAvailabilityOptions`
- New struct `ServiceClientCheckNameAvailabilityResponse`
- New struct `ServiceClientCheckNameAvailabilityResult`
- New struct `ServiceClientCreateOrUpdatePoller`
- New struct `ServiceClientCreateOrUpdatePollerResponse`
- New struct `ServiceClientCreateOrUpdateResponse`
- New struct `ServiceClientCreateOrUpdateResult`
- New struct `ServiceClientDeletePoller`
- New struct `ServiceClientDeletePollerResponse`
- New struct `ServiceClientDeleteResponse`
- New struct `ServiceClientGetOptions`
- New struct `ServiceClientGetResponse`
- New struct `ServiceClientGetResult`
- New struct `ServiceClientLinkNotificationHubOptions`
- New struct `ServiceClientLinkNotificationHubResponse`
- New struct `ServiceClientLinkNotificationHubResult`
- New struct `ServiceClientListByResourceGroupOptions`
- New struct `ServiceClientListByResourceGroupPager`
- New struct `ServiceClientListByResourceGroupResponse`
- New struct `ServiceClientListByResourceGroupResult`
- New struct `ServiceClientListBySubscriptionOptions`
- New struct `ServiceClientListBySubscriptionPager`
- New struct `ServiceClientListBySubscriptionResponse`
- New struct `ServiceClientListBySubscriptionResult`
- New struct `ServiceClientListKeysOptions`
- New struct `ServiceClientListKeysResponse`
- New struct `ServiceClientListKeysResult`
- New struct `ServiceClientRegenerateKeyOptions`
- New struct `ServiceClientRegenerateKeyResponse`
- New struct `ServiceClientRegenerateKeyResult`
- New struct `ServiceClientUpdateOptions`
- New struct `ServiceClientUpdateResponse`
- New struct `ServiceClientUpdateResult`
- New struct `ServiceKeys`
- New struct `ServiceProperties`
- New struct `ServiceResource`
- New struct `ServiceResourceList`
- New field `Error` in struct `ErrorResponse`


## 0.1.0 (2021-12-01)

- Initial preview release.
