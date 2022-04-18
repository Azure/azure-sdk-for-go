# Release History

## 0.3.0 (2022-04-18)
### Breaking Changes

- Function `*EnergyServicesClient.ListBySubscription` has been removed
- Function `*EnergyServicesClient.ListByResourceGroup` has been removed

### Features Added

- New function `*EnergyServicesClient.NewListByResourceGroupPager(string, *EnergyServicesClientListByResourceGroupOptions) *runtime.Pager[EnergyServicesClientListByResourceGroupResponse]`
- New function `*EnergyServicesClient.NewListBySubscriptionPager(*EnergyServicesClientListBySubscriptionOptions) *runtime.Pager[EnergyServicesClientListBySubscriptionResponse]`


## 0.2.0 (2022-04-12)
### Breaking Changes

- Function `NewEnergyServicesClient` return value(s) have been changed from `(*EnergyServicesClient)` to `(*EnergyServicesClient, error)`
- Function `*EnergyServicesClient.ListBySubscription` return value(s) have been changed from `(*EnergyServicesClientListBySubscriptionPager)` to `(*runtime.Pager[EnergyServicesClientListBySubscriptionResponse])`
- Function `*EnergyServicesClient.BeginCreate` return value(s) have been changed from `(EnergyServicesClientCreatePollerResponse, error)` to `(*armruntime.Poller[EnergyServicesClientCreateResponse], error)`
- Function `*EnergyServicesClient.BeginDelete` return value(s) have been changed from `(EnergyServicesClientDeletePollerResponse, error)` to `(*armruntime.Poller[EnergyServicesClientDeleteResponse], error)`
- Function `NewLocationsClient` return value(s) have been changed from `(*LocationsClient)` to `(*LocationsClient, error)`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*EnergyServicesClient.ListByResourceGroup` return value(s) have been changed from `(*EnergyServicesClientListByResourceGroupPager)` to `(*runtime.Pager[EnergyServicesClientListByResourceGroupResponse])`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Function `*EnergyServicesClientListByResourceGroupPager.NextPage` has been removed
- Function `*EnergyServicesClientCreatePoller.Done` has been removed
- Function `*EnergyServicesClientListByResourceGroupPager.Err` has been removed
- Function `CheckNameAvailabilityReason.ToPtr` has been removed
- Function `*EnergyServicesClientCreatePollerResponse.Resume` has been removed
- Function `*EnergyServicesClientDeletePoller.ResumeToken` has been removed
- Function `*EnergyServicesClientDeletePoller.Done` has been removed
- Function `*EnergyServicesClientCreatePoller.FinalResponse` has been removed
- Function `*EnergyServicesClientListBySubscriptionPager.PageResponse` has been removed
- Function `*EnergyServicesClientCreatePoller.Poll` has been removed
- Function `*EnergyServicesClientDeletePoller.Poll` has been removed
- Function `ProvisioningState.ToPtr` has been removed
- Function `EnergyServicesClientDeletePollerResponse.PollUntilDone` has been removed
- Function `ActionType.ToPtr` has been removed
- Function `*EnergyServicesClientDeletePoller.FinalResponse` has been removed
- Function `*EnergyServicesClientListBySubscriptionPager.Err` has been removed
- Function `EnergyServicesClientCreatePollerResponse.PollUntilDone` has been removed
- Function `*EnergyServicesClientListByResourceGroupPager.PageResponse` has been removed
- Function `*EnergyServicesClientListBySubscriptionPager.NextPage` has been removed
- Function `*EnergyServicesClientDeletePollerResponse.Resume` has been removed
- Function `CreatedByType.ToPtr` has been removed
- Function `Origin.ToPtr` has been removed
- Function `*EnergyServicesClientCreatePoller.ResumeToken` has been removed
- Struct `EnergyServicesClientCreatePoller` has been removed
- Struct `EnergyServicesClientCreatePollerResponse` has been removed
- Struct `EnergyServicesClientCreateResult` has been removed
- Struct `EnergyServicesClientDeletePoller` has been removed
- Struct `EnergyServicesClientDeletePollerResponse` has been removed
- Struct `EnergyServicesClientGetResult` has been removed
- Struct `EnergyServicesClientListByResourceGroupPager` has been removed
- Struct `EnergyServicesClientListByResourceGroupResult` has been removed
- Struct `EnergyServicesClientListBySubscriptionPager` has been removed
- Struct `EnergyServicesClientListBySubscriptionResult` has been removed
- Struct `EnergyServicesClientUpdateResult` has been removed
- Struct `LocationsClientCheckNameAvailabilityResult` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `LocationsClientCheckNameAvailabilityResult` of struct `LocationsClientCheckNameAvailabilityResponse` has been removed
- Field `RawResponse` of struct `LocationsClientCheckNameAvailabilityResponse` has been removed
- Field `RawResponse` of struct `EnergyServicesClientDeleteResponse` has been removed
- Field `EnergyServicesClientUpdateResult` of struct `EnergyServicesClientUpdateResponse` has been removed
- Field `RawResponse` of struct `EnergyServicesClientUpdateResponse` has been removed
- Field `EnergyServicesClientGetResult` of struct `EnergyServicesClientGetResponse` has been removed
- Field `RawResponse` of struct `EnergyServicesClientGetResponse` has been removed
- Field `EnergyServicesClientListByResourceGroupResult` of struct `EnergyServicesClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `EnergyServicesClientListByResourceGroupResponse` has been removed
- Field `EnergyServicesClientListBySubscriptionResult` of struct `EnergyServicesClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `EnergyServicesClientListBySubscriptionResponse` has been removed
- Field `EnergyServicesClientCreateResult` of struct `EnergyServicesClientCreateResponse` has been removed
- Field `RawResponse` of struct `EnergyServicesClientCreateResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed

### Features Added

- New anonymous field `EnergyService` in struct `EnergyServicesClientUpdateResponse`
- New field `ResumeToken` in struct `EnergyServicesClientBeginDeleteOptions`
- New field `ResumeToken` in struct `EnergyServicesClientBeginCreateOptions`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `EnergyService` in struct `EnergyServicesClientGetResponse`
- New anonymous field `EnergyServiceList` in struct `EnergyServicesClientListByResourceGroupResponse`
- New anonymous field `EnergyServiceList` in struct `EnergyServicesClientListBySubscriptionResponse`
- New anonymous field `EnergyService` in struct `EnergyServicesClientCreateResponse`
- New anonymous field `CheckNameAvailabilityResponse` in struct `LocationsClientCheckNameAvailabilityResponse`


## 0.1.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.1.0 (2022-01-14)

- Init release.