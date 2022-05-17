# Release History

## 1.0.0 (2022-05-17)
### Breaking Changes

- Function `*ZonesClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[ZonesClientDeleteResponse], error)` to `(*runtime.Poller[ZonesClientDeleteResponse], error)`
- Function `RecordSetListResult.MarshalJSON` has been removed
- Function `ResourceReference.MarshalJSON` has been removed
- Function `CloudErrorBody.MarshalJSON` has been removed
- Function `ResourceReferenceResultProperties.MarshalJSON` has been removed
- Function `ZoneListResult.MarshalJSON` has been removed


## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*ZonesClient.List` has been removed
- Function `*RecordSetsClient.ListByType` has been removed
- Function `*RecordSetsClient.ListAllByDNSZone` has been removed
- Function `*ZonesClient.ListByResourceGroup` has been removed
- Function `*RecordSetsClient.ListByDNSZone` has been removed

### Features Added

- New function `*RecordSetsClient.NewListAllByDNSZonePager(string, string, *RecordSetsClientListAllByDNSZoneOptions) *runtime.Pager[RecordSetsClientListAllByDNSZoneResponse]`
- New function `*ZonesClient.NewListByResourceGroupPager(string, *ZonesClientListByResourceGroupOptions) *runtime.Pager[ZonesClientListByResourceGroupResponse]`
- New function `*ZonesClient.NewListPager(*ZonesClientListOptions) *runtime.Pager[ZonesClientListResponse]`
- New function `*RecordSetsClient.NewListByDNSZonePager(string, string, *RecordSetsClientListByDNSZoneOptions) *runtime.Pager[RecordSetsClientListByDNSZoneResponse]`
- New function `*RecordSetsClient.NewListByTypePager(string, string, RecordType, *RecordSetsClientListByTypeOptions) *runtime.Pager[RecordSetsClientListByTypeResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `*ZonesClient.List` return value(s) have been changed from `(*ZonesClientListPager)` to `(*runtime.Pager[ZonesClientListResponse])`
- Function `*ZonesClient.ListByResourceGroup` return value(s) have been changed from `(*ZonesClientListByResourceGroupPager)` to `(*runtime.Pager[ZonesClientListByResourceGroupResponse])`
- Function `NewRecordSetsClient` return value(s) have been changed from `(*RecordSetsClient)` to `(*RecordSetsClient, error)`
- Function `*RecordSetsClient.ListByDNSZone` return value(s) have been changed from `(*RecordSetsClientListByDNSZonePager)` to `(*runtime.Pager[RecordSetsClientListByDNSZoneResponse])`
- Function `NewResourceReferenceClient` return value(s) have been changed from `(*ResourceReferenceClient)` to `(*ResourceReferenceClient, error)`
- Function `*RecordSetsClient.ListAllByDNSZone` return value(s) have been changed from `(*RecordSetsClientListAllByDNSZonePager)` to `(*runtime.Pager[RecordSetsClientListAllByDNSZoneResponse])`
- Function `NewZonesClient` return value(s) have been changed from `(*ZonesClient)` to `(*ZonesClient, error)`
- Function `*ZonesClient.BeginDelete` return value(s) have been changed from `(ZonesClientDeletePollerResponse, error)` to `(*armruntime.Poller[ZonesClientDeleteResponse], error)`
- Function `*RecordSetsClient.ListByType` return value(s) have been changed from `(*RecordSetsClientListByTypePager)` to `(*runtime.Pager[RecordSetsClientListByTypeResponse])`
- Function `*RecordSetsClientListByDNSZonePager.PageResponse` has been removed
- Function `*RecordSetsClientListByTypePager.PageResponse` has been removed
- Function `ZoneType.ToPtr` has been removed
- Function `*RecordSetsClientListAllByDNSZonePager.PageResponse` has been removed
- Function `*RecordSetsClientListAllByDNSZonePager.NextPage` has been removed
- Function `*RecordSetsClientListAllByDNSZonePager.Err` has been removed
- Function `*ZonesClientDeletePoller.FinalResponse` has been removed
- Function `*ZonesClientListPager.NextPage` has been removed
- Function `*RecordSetsClientListByDNSZonePager.Err` has been removed
- Function `*ZonesClientDeletePoller.ResumeToken` has been removed
- Function `*ZonesClientDeletePoller.Poll` has been removed
- Function `*ZonesClientDeletePoller.Done` has been removed
- Function `*ZonesClientListByResourceGroupPager.NextPage` has been removed
- Function `*ZonesClientListByResourceGroupPager.PageResponse` has been removed
- Function `*ZonesClientListPager.PageResponse` has been removed
- Function `RecordType.ToPtr` has been removed
- Function `*ZonesClientDeletePollerResponse.Resume` has been removed
- Function `*RecordSetsClientListByTypePager.Err` has been removed
- Function `*ZonesClientListPager.Err` has been removed
- Function `ZonesClientDeletePollerResponse.PollUntilDone` has been removed
- Function `*RecordSetsClientListByDNSZonePager.NextPage` has been removed
- Function `*ZonesClientListByResourceGroupPager.Err` has been removed
- Function `*RecordSetsClientListByTypePager.NextPage` has been removed
- Struct `RecordSetsClientCreateOrUpdateResult` has been removed
- Struct `RecordSetsClientGetResult` has been removed
- Struct `RecordSetsClientListAllByDNSZonePager` has been removed
- Struct `RecordSetsClientListAllByDNSZoneResult` has been removed
- Struct `RecordSetsClientListByDNSZonePager` has been removed
- Struct `RecordSetsClientListByDNSZoneResult` has been removed
- Struct `RecordSetsClientListByTypePager` has been removed
- Struct `RecordSetsClientListByTypeResult` has been removed
- Struct `RecordSetsClientUpdateResult` has been removed
- Struct `ResourceReferenceClientGetByTargetResourcesResult` has been removed
- Struct `ZonesClientCreateOrUpdateResult` has been removed
- Struct `ZonesClientDeletePoller` has been removed
- Struct `ZonesClientDeletePollerResponse` has been removed
- Struct `ZonesClientGetResult` has been removed
- Struct `ZonesClientListByResourceGroupPager` has been removed
- Struct `ZonesClientListByResourceGroupResult` has been removed
- Struct `ZonesClientListPager` has been removed
- Struct `ZonesClientListResult` has been removed
- Struct `ZonesClientUpdateResult` has been removed
- Field `RecordSetsClientUpdateResult` of struct `RecordSetsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `RecordSetsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `ZonesClientDeleteResponse` has been removed
- Field `RecordSetsClientListByDNSZoneResult` of struct `RecordSetsClientListByDNSZoneResponse` has been removed
- Field `RawResponse` of struct `RecordSetsClientListByDNSZoneResponse` has been removed
- Field `RecordSetsClientGetResult` of struct `RecordSetsClientGetResponse` has been removed
- Field `RawResponse` of struct `RecordSetsClientGetResponse` has been removed
- Field `ZonesClientGetResult` of struct `ZonesClientGetResponse` has been removed
- Field `RawResponse` of struct `ZonesClientGetResponse` has been removed
- Field `ZonesClientUpdateResult` of struct `ZonesClientUpdateResponse` has been removed
- Field `RawResponse` of struct `ZonesClientUpdateResponse` has been removed
- Field `RecordSetsClientListAllByDNSZoneResult` of struct `RecordSetsClientListAllByDNSZoneResponse` has been removed
- Field `RawResponse` of struct `RecordSetsClientListAllByDNSZoneResponse` has been removed
- Field `RecordSetsClientCreateOrUpdateResult` of struct `RecordSetsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `RecordSetsClientCreateOrUpdateResponse` has been removed
- Field `RecordSetsClientListByTypeResult` of struct `RecordSetsClientListByTypeResponse` has been removed
- Field `RawResponse` of struct `RecordSetsClientListByTypeResponse` has been removed
- Field `ZonesClientListResult` of struct `ZonesClientListResponse` has been removed
- Field `RawResponse` of struct `ZonesClientListResponse` has been removed
- Field `ResourceReferenceClientGetByTargetResourcesResult` of struct `ResourceReferenceClientGetByTargetResourcesResponse` has been removed
- Field `RawResponse` of struct `ResourceReferenceClientGetByTargetResourcesResponse` has been removed
- Field `ZonesClientCreateOrUpdateResult` of struct `ZonesClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `ZonesClientCreateOrUpdateResponse` has been removed
- Field `ZonesClientListByResourceGroupResult` of struct `ZonesClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `ZonesClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `RecordSetsClientDeleteResponse` has been removed

### Features Added

- New anonymous field `RecordSet` in struct `RecordSetsClientGetResponse`
- New anonymous field `RecordSetListResult` in struct `RecordSetsClientListAllByDNSZoneResponse`
- New anonymous field `RecordSetListResult` in struct `RecordSetsClientListByDNSZoneResponse`
- New anonymous field `Zone` in struct `ZonesClientGetResponse`
- New anonymous field `ResourceReferenceResult` in struct `ResourceReferenceClientGetByTargetResourcesResponse`
- New anonymous field `Zone` in struct `ZonesClientUpdateResponse`
- New anonymous field `RecordSetListResult` in struct `RecordSetsClientListByTypeResponse`
- New anonymous field `ZoneListResult` in struct `ZonesClientListByResourceGroupResponse`
- New anonymous field `RecordSet` in struct `RecordSetsClientUpdateResponse`
- New anonymous field `Zone` in struct `ZonesClientCreateOrUpdateResponse`
- New field `ResumeToken` in struct `ZonesClientBeginDeleteOptions`
- New anonymous field `ZoneListResult` in struct `ZonesClientListResponse`
- New anonymous field `RecordSet` in struct `RecordSetsClientCreateOrUpdateResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*RecordSetsClient.ListAllByDNSZone` parameter(s) have been changed from `(string, string, *RecordSetsListAllByDNSZoneOptions)` to `(string, string, *RecordSetsClientListAllByDNSZoneOptions)`
- Function `*RecordSetsClient.ListAllByDNSZone` return value(s) have been changed from `(*RecordSetsListAllByDNSZonePager)` to `(*RecordSetsClientListAllByDNSZonePager)`
- Function `*RecordSetsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, string, RecordType, *RecordSetsDeleteOptions)` to `(context.Context, string, string, string, RecordType, *RecordSetsClientDeleteOptions)`
- Function `*RecordSetsClient.Delete` return value(s) have been changed from `(RecordSetsDeleteResponse, error)` to `(RecordSetsClientDeleteResponse, error)`
- Function `*RecordSetsClient.Update` parameter(s) have been changed from `(context.Context, string, string, string, RecordType, RecordSet, *RecordSetsUpdateOptions)` to `(context.Context, string, string, string, RecordType, RecordSet, *RecordSetsClientUpdateOptions)`
- Function `*RecordSetsClient.Update` return value(s) have been changed from `(RecordSetsUpdateResponse, error)` to `(RecordSetsClientUpdateResponse, error)`
- Function `*ZonesClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, Zone, *ZonesCreateOrUpdateOptions)` to `(context.Context, string, string, Zone, *ZonesClientCreateOrUpdateOptions)`
- Function `*ZonesClient.CreateOrUpdate` return value(s) have been changed from `(ZonesCreateOrUpdateResponse, error)` to `(ZonesClientCreateOrUpdateResponse, error)`
- Function `*RecordSetsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, RecordType, RecordSet, *RecordSetsCreateOrUpdateOptions)` to `(context.Context, string, string, string, RecordType, RecordSet, *RecordSetsClientCreateOrUpdateOptions)`
- Function `*RecordSetsClient.CreateOrUpdate` return value(s) have been changed from `(RecordSetsCreateOrUpdateResponse, error)` to `(RecordSetsClientCreateOrUpdateResponse, error)`
- Function `*RecordSetsClient.ListByDNSZone` parameter(s) have been changed from `(string, string, *RecordSetsListByDNSZoneOptions)` to `(string, string, *RecordSetsClientListByDNSZoneOptions)`
- Function `*RecordSetsClient.ListByDNSZone` return value(s) have been changed from `(*RecordSetsListByDNSZonePager)` to `(*RecordSetsClientListByDNSZonePager)`
- Function `*RecordSetsClient.ListByType` parameter(s) have been changed from `(string, string, RecordType, *RecordSetsListByTypeOptions)` to `(string, string, RecordType, *RecordSetsClientListByTypeOptions)`
- Function `*RecordSetsClient.ListByType` return value(s) have been changed from `(*RecordSetsListByTypePager)` to `(*RecordSetsClientListByTypePager)`
- Function `*ZonesClient.ListByResourceGroup` parameter(s) have been changed from `(string, *ZonesListByResourceGroupOptions)` to `(string, *ZonesClientListByResourceGroupOptions)`
- Function `*ZonesClient.ListByResourceGroup` return value(s) have been changed from `(*ZonesListByResourceGroupPager)` to `(*ZonesClientListByResourceGroupPager)`
- Function `*RecordSetsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, RecordType, *RecordSetsGetOptions)` to `(context.Context, string, string, string, RecordType, *RecordSetsClientGetOptions)`
- Function `*RecordSetsClient.Get` return value(s) have been changed from `(RecordSetsGetResponse, error)` to `(RecordSetsClientGetResponse, error)`
- Function `*ZonesClient.Update` parameter(s) have been changed from `(context.Context, string, string, ZoneUpdate, *ZonesUpdateOptions)` to `(context.Context, string, string, ZoneUpdate, *ZonesClientUpdateOptions)`
- Function `*ZonesClient.Update` return value(s) have been changed from `(ZonesUpdateResponse, error)` to `(ZonesClientUpdateResponse, error)`
- Function `*ZonesClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ZonesGetOptions)` to `(context.Context, string, string, *ZonesClientGetOptions)`
- Function `*ZonesClient.Get` return value(s) have been changed from `(ZonesGetResponse, error)` to `(ZonesClientGetResponse, error)`
- Function `*ZonesClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *ZonesBeginDeleteOptions)` to `(context.Context, string, string, *ZonesClientBeginDeleteOptions)`
- Function `*ZonesClient.BeginDelete` return value(s) have been changed from `(ZonesDeletePollerResponse, error)` to `(ZonesClientDeletePollerResponse, error)`
- Function `*ZonesClient.List` parameter(s) have been changed from `(*ZonesListOptions)` to `(*ZonesClientListOptions)`
- Function `*ZonesClient.List` return value(s) have been changed from `(*ZonesListPager)` to `(*ZonesClientListPager)`
- Function `*ZonesDeletePollerResponse.Resume` has been removed
- Function `*ZonesDeletePoller.Done` has been removed
- Function `NewDNSResourceReferenceClient` has been removed
- Function `*ZonesListByResourceGroupPager.PageResponse` has been removed
- Function `*RecordSetsListByDNSZonePager.Err` has been removed
- Function `*ZonesListByResourceGroupPager.NextPage` has been removed
- Function `DNSResourceReferenceRequestProperties.MarshalJSON` has been removed
- Function `*RecordSetsListByDNSZonePager.PageResponse` has been removed
- Function `*ZonesListPager.PageResponse` has been removed
- Function `*RecordSetsListByDNSZonePager.NextPage` has been removed
- Function `*RecordSetsListAllByDNSZonePager.NextPage` has been removed
- Function `*RecordSetsListByTypePager.Err` has been removed
- Function `*RecordSetsListByTypePager.NextPage` has been removed
- Function `*RecordSetsListAllByDNSZonePager.Err` has been removed
- Function `*DNSResourceReferenceClient.GetByTargetResources` has been removed
- Function `*ZonesListPager.NextPage` has been removed
- Function `*RecordSetsListByTypePager.PageResponse` has been removed
- Function `DNSResourceReference.MarshalJSON` has been removed
- Function `*ZonesListPager.Err` has been removed
- Function `*ZonesDeletePoller.ResumeToken` has been removed
- Function `ZonesDeletePollerResponse.PollUntilDone` has been removed
- Function `CloudError.Error` has been removed
- Function `*ZonesDeletePoller.Poll` has been removed
- Function `*RecordSetsListAllByDNSZonePager.PageResponse` has been removed
- Function `DNSResourceReferenceResultProperties.MarshalJSON` has been removed
- Function `*ZonesDeletePoller.FinalResponse` has been removed
- Function `*ZonesListByResourceGroupPager.Err` has been removed
- Struct `DNSResourceReference` has been removed
- Struct `DNSResourceReferenceClient` has been removed
- Struct `DNSResourceReferenceGetByTargetResourcesOptions` has been removed
- Struct `DNSResourceReferenceGetByTargetResourcesResponse` has been removed
- Struct `DNSResourceReferenceGetByTargetResourcesResult` has been removed
- Struct `DNSResourceReferenceRequest` has been removed
- Struct `DNSResourceReferenceRequestProperties` has been removed
- Struct `DNSResourceReferenceResult` has been removed
- Struct `DNSResourceReferenceResultProperties` has been removed
- Struct `RecordSetsCreateOrUpdateOptions` has been removed
- Struct `RecordSetsCreateOrUpdateResponse` has been removed
- Struct `RecordSetsCreateOrUpdateResult` has been removed
- Struct `RecordSetsDeleteOptions` has been removed
- Struct `RecordSetsDeleteResponse` has been removed
- Struct `RecordSetsGetOptions` has been removed
- Struct `RecordSetsGetResponse` has been removed
- Struct `RecordSetsGetResult` has been removed
- Struct `RecordSetsListAllByDNSZoneOptions` has been removed
- Struct `RecordSetsListAllByDNSZonePager` has been removed
- Struct `RecordSetsListAllByDNSZoneResponse` has been removed
- Struct `RecordSetsListAllByDNSZoneResult` has been removed
- Struct `RecordSetsListByDNSZoneOptions` has been removed
- Struct `RecordSetsListByDNSZonePager` has been removed
- Struct `RecordSetsListByDNSZoneResponse` has been removed
- Struct `RecordSetsListByDNSZoneResult` has been removed
- Struct `RecordSetsListByTypeOptions` has been removed
- Struct `RecordSetsListByTypePager` has been removed
- Struct `RecordSetsListByTypeResponse` has been removed
- Struct `RecordSetsListByTypeResult` has been removed
- Struct `RecordSetsUpdateOptions` has been removed
- Struct `RecordSetsUpdateResponse` has been removed
- Struct `RecordSetsUpdateResult` has been removed
- Struct `ZonesBeginDeleteOptions` has been removed
- Struct `ZonesCreateOrUpdateOptions` has been removed
- Struct `ZonesCreateOrUpdateResponse` has been removed
- Struct `ZonesCreateOrUpdateResult` has been removed
- Struct `ZonesDeletePoller` has been removed
- Struct `ZonesDeletePollerResponse` has been removed
- Struct `ZonesDeleteResponse` has been removed
- Struct `ZonesGetOptions` has been removed
- Struct `ZonesGetResponse` has been removed
- Struct `ZonesGetResult` has been removed
- Struct `ZonesListByResourceGroupOptions` has been removed
- Struct `ZonesListByResourceGroupPager` has been removed
- Struct `ZonesListByResourceGroupResponse` has been removed
- Struct `ZonesListByResourceGroupResult` has been removed
- Struct `ZonesListOptions` has been removed
- Struct `ZonesListPager` has been removed
- Struct `ZonesListResponse` has been removed
- Struct `ZonesListResult` has been removed
- Struct `ZonesUpdateOptions` has been removed
- Struct `ZonesUpdateResponse` has been removed
- Struct `ZonesUpdateResult` has been removed
- Field `Resource` of struct `Zone` has been removed
- Field `InnerError` of struct `CloudError` has been removed

### Features Added

- New function `*ZonesClientListPager.PageResponse() ZonesClientListResponse`
- New function `*ZonesClientListPager.NextPage(context.Context) bool`
- New function `ResourceReferenceResultProperties.MarshalJSON() ([]byte, error)`
- New function `ZonesClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (ZonesClientDeleteResponse, error)`
- New function `*ZonesClientDeletePollerResponse.Resume(context.Context, *ZonesClient, string) error`
- New function `*RecordSetsClientListByDNSZonePager.PageResponse() RecordSetsClientListByDNSZoneResponse`
- New function `*ZonesClientDeletePoller.FinalResponse(context.Context) (ZonesClientDeleteResponse, error)`
- New function `*ZonesClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*RecordSetsClientListAllByDNSZonePager.Err() error`
- New function `*RecordSetsClientListByDNSZonePager.NextPage(context.Context) bool`
- New function `*RecordSetsClientListByTypePager.Err() error`
- New function `*ZonesClientListPager.Err() error`
- New function `*RecordSetsClientListByTypePager.PageResponse() RecordSetsClientListByTypeResponse`
- New function `*ZonesClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*ZonesClientListByResourceGroupPager.PageResponse() ZonesClientListByResourceGroupResponse`
- New function `*RecordSetsClientListByTypePager.NextPage(context.Context) bool`
- New function `*ZonesClientDeletePoller.ResumeToken() (string, error)`
- New function `ResourceReference.MarshalJSON() ([]byte, error)`
- New function `ResourceReferenceRequestProperties.MarshalJSON() ([]byte, error)`
- New function `*RecordSetsClientListAllByDNSZonePager.NextPage(context.Context) bool`
- New function `*ZonesClientListByResourceGroupPager.Err() error`
- New function `NewResourceReferenceClient(string, azcore.TokenCredential, *arm.ClientOptions) *ResourceReferenceClient`
- New function `*ResourceReferenceClient.GetByTargetResources(context.Context, ResourceReferenceRequest, *ResourceReferenceClientGetByTargetResourcesOptions) (ResourceReferenceClientGetByTargetResourcesResponse, error)`
- New function `*RecordSetsClientListAllByDNSZonePager.PageResponse() RecordSetsClientListAllByDNSZoneResponse`
- New function `*RecordSetsClientListByDNSZonePager.Err() error`
- New function `*ZonesClientDeletePoller.Done() bool`
- New struct `RecordSetsClientCreateOrUpdateOptions`
- New struct `RecordSetsClientCreateOrUpdateResponse`
- New struct `RecordSetsClientCreateOrUpdateResult`
- New struct `RecordSetsClientDeleteOptions`
- New struct `RecordSetsClientDeleteResponse`
- New struct `RecordSetsClientGetOptions`
- New struct `RecordSetsClientGetResponse`
- New struct `RecordSetsClientGetResult`
- New struct `RecordSetsClientListAllByDNSZoneOptions`
- New struct `RecordSetsClientListAllByDNSZonePager`
- New struct `RecordSetsClientListAllByDNSZoneResponse`
- New struct `RecordSetsClientListAllByDNSZoneResult`
- New struct `RecordSetsClientListByDNSZoneOptions`
- New struct `RecordSetsClientListByDNSZonePager`
- New struct `RecordSetsClientListByDNSZoneResponse`
- New struct `RecordSetsClientListByDNSZoneResult`
- New struct `RecordSetsClientListByTypeOptions`
- New struct `RecordSetsClientListByTypePager`
- New struct `RecordSetsClientListByTypeResponse`
- New struct `RecordSetsClientListByTypeResult`
- New struct `RecordSetsClientUpdateOptions`
- New struct `RecordSetsClientUpdateResponse`
- New struct `RecordSetsClientUpdateResult`
- New struct `ResourceReference`
- New struct `ResourceReferenceClient`
- New struct `ResourceReferenceClientGetByTargetResourcesOptions`
- New struct `ResourceReferenceClientGetByTargetResourcesResponse`
- New struct `ResourceReferenceClientGetByTargetResourcesResult`
- New struct `ResourceReferenceRequest`
- New struct `ResourceReferenceRequestProperties`
- New struct `ResourceReferenceResult`
- New struct `ResourceReferenceResultProperties`
- New struct `ZonesClientBeginDeleteOptions`
- New struct `ZonesClientCreateOrUpdateOptions`
- New struct `ZonesClientCreateOrUpdateResponse`
- New struct `ZonesClientCreateOrUpdateResult`
- New struct `ZonesClientDeletePoller`
- New struct `ZonesClientDeletePollerResponse`
- New struct `ZonesClientDeleteResponse`
- New struct `ZonesClientGetOptions`
- New struct `ZonesClientGetResponse`
- New struct `ZonesClientGetResult`
- New struct `ZonesClientListByResourceGroupOptions`
- New struct `ZonesClientListByResourceGroupPager`
- New struct `ZonesClientListByResourceGroupResponse`
- New struct `ZonesClientListByResourceGroupResult`
- New struct `ZonesClientListOptions`
- New struct `ZonesClientListPager`
- New struct `ZonesClientListResponse`
- New struct `ZonesClientListResult`
- New struct `ZonesClientUpdateOptions`
- New struct `ZonesClientUpdateResponse`
- New struct `ZonesClientUpdateResult`
- New field `Name` in struct `Zone`
- New field `Type` in struct `Zone`
- New field `Location` in struct `Zone`
- New field `Tags` in struct `Zone`
- New field `ID` in struct `Zone`
- New field `Error` in struct `CloudError`


## 0.1.0 (2021-12-07)

- Initial preview release.
