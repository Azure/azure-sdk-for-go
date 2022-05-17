# Release History

## 1.0.0 (2022-05-17)
### Breaking Changes

- Function `OperationList.MarshalJSON` has been removed
- Function `AssignmentReportDetails.MarshalJSON` has been removed
- Function `AssignmentReportList.MarshalJSON` has been removed
- Function `AssignmentReportProperties.MarshalJSON` has been removed
- Function `AssignmentList.MarshalJSON` has been removed

### Features Added

- New const `CreatedByTypeKey`
- New const `CreatedByTypeUser`
- New const `CreatedByTypeApplication`
- New const `CreatedByTypeManagedIdentity`
- New function `SystemData.MarshalJSON() ([]byte, error)`
- New function `*SystemData.UnmarshalJSON([]byte) error`
- New function `PossibleCreatedByTypeValues() []CreatedByType`
- New struct `AssignmentReportsVMSSClientGetOptions`
- New struct `AssignmentReportsVMSSClientGetResponse`
- New struct `AssignmentReportsVMSSClientListOptions`
- New struct `AssignmentReportsVMSSClientListResponse`
- New struct `AssignmentsVMSSClientDeleteOptions`
- New struct `AssignmentsVMSSClientDeleteResponse`
- New struct `AssignmentsVMSSClientGetOptions`
- New struct `AssignmentsVMSSClientGetResponse`
- New struct `AssignmentsVMSSClientListOptions`
- New struct `AssignmentsVMSSClientListResponse`
- New struct `SystemData`
- New field `SystemData` in struct `Assignment`
- New field `AssignmentSource` in struct `Navigation`


## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*AssignmentsClient.SubscriptionList` has been removed
- Function `*AssignmentsClient.List` has been removed
- Function `*AssignmentsClient.RGList` has been removed
- Function `*OperationsClient.List` has been removed
- Function `*HCRPAssignmentsClient.List` has been removed

### Features Added

- New function `*AssignmentsClient.NewSubscriptionListPager(*AssignmentsClientSubscriptionListOptions) *runtime.Pager[AssignmentsClientSubscriptionListResponse]`
- New function `*HCRPAssignmentsClient.NewListPager(string, string, *HCRPAssignmentsClientListOptions) *runtime.Pager[HCRPAssignmentsClientListResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `*AssignmentsClient.NewRGListPager(string, *AssignmentsClientRGListOptions) *runtime.Pager[AssignmentsClientRGListResponse]`
- New function `*AssignmentsClient.NewListPager(string, string, *AssignmentsClientListOptions) *runtime.Pager[AssignmentsClientListResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `NewAssignmentReportsClient` return value(s) have been changed from `(*AssignmentReportsClient)` to `(*AssignmentReportsClient, error)`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `NewAssignmentsClient` return value(s) have been changed from `(*AssignmentsClient)` to `(*AssignmentsClient, error)`
- Function `*AssignmentsClient.List` parameter(s) have been changed from `(context.Context, string, string, *AssignmentsClientListOptions)` to `(string, string, *AssignmentsClientListOptions)`
- Function `*AssignmentsClient.List` return value(s) have been changed from `(AssignmentsClientListResponse, error)` to `(*runtime.Pager[AssignmentsClientListResponse])`
- Function `NewHCRPAssignmentReportsClient` return value(s) have been changed from `(*HCRPAssignmentReportsClient)` to `(*HCRPAssignmentReportsClient, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsClientListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsClientListResponse, error)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `*HCRPAssignmentsClient.List` parameter(s) have been changed from `(context.Context, string, string, *HCRPAssignmentsClientListOptions)` to `(string, string, *HCRPAssignmentsClientListOptions)`
- Function `*HCRPAssignmentsClient.List` return value(s) have been changed from `(HCRPAssignmentsClientListResponse, error)` to `(*runtime.Pager[HCRPAssignmentsClientListResponse])`
- Function `*AssignmentsClient.RGList` parameter(s) have been changed from `(context.Context, string, *AssignmentsClientRGListOptions)` to `(string, *AssignmentsClientRGListOptions)`
- Function `*AssignmentsClient.RGList` return value(s) have been changed from `(AssignmentsClientRGListResponse, error)` to `(*runtime.Pager[AssignmentsClientRGListResponse])`
- Function `*AssignmentsClient.SubscriptionList` parameter(s) have been changed from `(context.Context, *AssignmentsClientSubscriptionListOptions)` to `(*AssignmentsClientSubscriptionListOptions)`
- Function `*AssignmentsClient.SubscriptionList` return value(s) have been changed from `(AssignmentsClientSubscriptionListResponse, error)` to `(*runtime.Pager[AssignmentsClientSubscriptionListResponse])`
- Function `NewHCRPAssignmentsClient` return value(s) have been changed from `(*HCRPAssignmentsClient)` to `(*HCRPAssignmentsClient, error)`
- Function `ComplianceStatus.ToPtr` has been removed
- Function `ActionAfterReboot.ToPtr` has been removed
- Function `Type.ToPtr` has been removed
- Function `Kind.ToPtr` has been removed
- Function `ProvisioningState.ToPtr` has been removed
- Function `ConfigurationMode.ToPtr` has been removed
- Function `AssignmentType.ToPtr` has been removed
- Struct `AssignmentReportsClientGetResult` has been removed
- Struct `AssignmentReportsClientListResult` has been removed
- Struct `AssignmentsClientCreateOrUpdateResult` has been removed
- Struct `AssignmentsClientGetResult` has been removed
- Struct `AssignmentsClientListResult` has been removed
- Struct `AssignmentsClientRGListResult` has been removed
- Struct `AssignmentsClientSubscriptionListResult` has been removed
- Struct `HCRPAssignmentReportsClientGetResult` has been removed
- Struct `HCRPAssignmentReportsClientListResult` has been removed
- Struct `HCRPAssignmentsClientCreateOrUpdateResult` has been removed
- Struct `HCRPAssignmentsClientGetResult` has been removed
- Struct `HCRPAssignmentsClientListResult` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `AssignmentsClientGetResult` of struct `AssignmentsClientGetResponse` has been removed
- Field `RawResponse` of struct `AssignmentsClientGetResponse` has been removed
- Field `AssignmentReportsClientListResult` of struct `AssignmentReportsClientListResponse` has been removed
- Field `RawResponse` of struct `AssignmentReportsClientListResponse` has been removed
- Field `AssignmentReportsClientGetResult` of struct `AssignmentReportsClientGetResponse` has been removed
- Field `RawResponse` of struct `AssignmentReportsClientGetResponse` has been removed
- Field `HCRPAssignmentsClientCreateOrUpdateResult` of struct `HCRPAssignmentsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `HCRPAssignmentsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `AssignmentsClientDeleteResponse` has been removed
- Field `AssignmentsClientCreateOrUpdateResult` of struct `AssignmentsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `AssignmentsClientCreateOrUpdateResponse` has been removed
- Field `AssignmentsClientListResult` of struct `AssignmentsClientListResponse` has been removed
- Field `RawResponse` of struct `AssignmentsClientListResponse` has been removed
- Field `HCRPAssignmentReportsClientListResult` of struct `HCRPAssignmentReportsClientListResponse` has been removed
- Field `RawResponse` of struct `HCRPAssignmentReportsClientListResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `HCRPAssignmentsClientDeleteResponse` has been removed
- Field `HCRPAssignmentsClientListResult` of struct `HCRPAssignmentsClientListResponse` has been removed
- Field `RawResponse` of struct `HCRPAssignmentsClientListResponse` has been removed
- Field `AssignmentsClientSubscriptionListResult` of struct `AssignmentsClientSubscriptionListResponse` has been removed
- Field `RawResponse` of struct `AssignmentsClientSubscriptionListResponse` has been removed
- Field `HCRPAssignmentsClientGetResult` of struct `HCRPAssignmentsClientGetResponse` has been removed
- Field `RawResponse` of struct `HCRPAssignmentsClientGetResponse` has been removed
- Field `AssignmentsClientRGListResult` of struct `AssignmentsClientRGListResponse` has been removed
- Field `RawResponse` of struct `AssignmentsClientRGListResponse` has been removed
- Field `HCRPAssignmentReportsClientGetResult` of struct `HCRPAssignmentReportsClientGetResponse` has been removed
- Field `RawResponse` of struct `HCRPAssignmentReportsClientGetResponse` has been removed

### Features Added

- New struct `ErrorResponse`
- New struct `ErrorResponseError`
- New anonymous field `AssignmentList` in struct `AssignmentsClientListResponse`
- New anonymous field `Assignment` in struct `AssignmentsClientGetResponse`
- New anonymous field `AssignmentReportList` in struct `HCRPAssignmentReportsClientListResponse`
- New anonymous field `AssignmentList` in struct `AssignmentsClientRGListResponse`
- New anonymous field `AssignmentReport` in struct `AssignmentReportsClientGetResponse`
- New anonymous field `AssignmentList` in struct `AssignmentsClientSubscriptionListResponse`
- New anonymous field `Assignment` in struct `HCRPAssignmentsClientGetResponse`
- New anonymous field `AssignmentReport` in struct `HCRPAssignmentReportsClientGetResponse`
- New anonymous field `Assignment` in struct `AssignmentsClientCreateOrUpdateResponse`
- New anonymous field `AssignmentList` in struct `HCRPAssignmentsClientListResponse`
- New anonymous field `AssignmentReportList` in struct `AssignmentReportsClientListResponse`
- New anonymous field `Assignment` in struct `HCRPAssignmentsClientCreateOrUpdateResponse`
- New anonymous field `OperationList` in struct `OperationsClientListResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-21)
### Breaking Changes

- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsListOptions)` to `(context.Context, *OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsListResponse, error)` to `(OperationsClientListResponse, error)`
- Type of `AssignmentReportResource.Properties` has been changed from `map[string]interface{}` to `interface{}`
- Function `*GuestConfigurationAssignmentProperties.UnmarshalJSON` has been removed
- Function `*GuestConfigurationHCRPAssignmentsClient.Delete` has been removed
- Function `*GuestConfigurationHCRPAssignmentReportsClient.Get` has been removed
- Function `GuestConfigurationAssignmentList.MarshalJSON` has been removed
- Function `*GuestConfigurationHCRPAssignmentsClient.CreateOrUpdate` has been removed
- Function `NewGuestConfigurationHCRPAssignmentReportsClient` has been removed
- Function `*GuestConfigurationAssignmentReportsClient.List` has been removed
- Function `GuestConfigurationAssignmentReportProperties.MarshalJSON` has been removed
- Function `*GuestConfigurationHCRPAssignmentsClient.Get` has been removed
- Function `NewGuestConfigurationAssignmentsClient` has been removed
- Function `*GuestConfigurationAssignmentReportsClient.Get` has been removed
- Function `*GuestConfigurationAssignmentsClient.Delete` has been removed
- Function `*GuestConfigurationAssignmentsClient.List` has been removed
- Function `*GuestConfigurationAssignmentReportProperties.UnmarshalJSON` has been removed
- Function `*GuestConfigurationHCRPAssignmentReportsClient.List` has been removed
- Function `AssignmentReport.MarshalJSON` has been removed
- Function `*GuestConfigurationAssignmentsClient.SubscriptionList` has been removed
- Function `NewGuestConfigurationHCRPAssignmentsClient` has been removed
- Function `*GuestConfigurationAssignmentsClient.RGList` has been removed
- Function `GuestConfigurationAssignmentReportList.MarshalJSON` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*GuestConfigurationAssignmentsClient.Get` has been removed
- Function `*GuestConfigurationAssignmentsClient.CreateOrUpdate` has been removed
- Function `GuestConfigurationAssignmentProperties.MarshalJSON` has been removed
- Function `GuestConfigurationNavigation.MarshalJSON` has been removed
- Function `*AssignmentReport.UnmarshalJSON` has been removed
- Function `*GuestConfigurationHCRPAssignmentsClient.List` has been removed
- Function `NewGuestConfigurationAssignmentReportsClient` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ErrorResponseError` has been removed
- Struct `GuestConfigurationAssignment` has been removed
- Struct `GuestConfigurationAssignmentList` has been removed
- Struct `GuestConfigurationAssignmentProperties` has been removed
- Struct `GuestConfigurationAssignmentReport` has been removed
- Struct `GuestConfigurationAssignmentReportList` has been removed
- Struct `GuestConfigurationAssignmentReportProperties` has been removed
- Struct `GuestConfigurationAssignmentReportsClient` has been removed
- Struct `GuestConfigurationAssignmentReportsGetOptions` has been removed
- Struct `GuestConfigurationAssignmentReportsGetResponse` has been removed
- Struct `GuestConfigurationAssignmentReportsGetResult` has been removed
- Struct `GuestConfigurationAssignmentReportsListOptions` has been removed
- Struct `GuestConfigurationAssignmentReportsListResponse` has been removed
- Struct `GuestConfigurationAssignmentReportsListResult` has been removed
- Struct `GuestConfigurationAssignmentsClient` has been removed
- Struct `GuestConfigurationAssignmentsCreateOrUpdateOptions` has been removed
- Struct `GuestConfigurationAssignmentsCreateOrUpdateResponse` has been removed
- Struct `GuestConfigurationAssignmentsCreateOrUpdateResult` has been removed
- Struct `GuestConfigurationAssignmentsDeleteOptions` has been removed
- Struct `GuestConfigurationAssignmentsDeleteResponse` has been removed
- Struct `GuestConfigurationAssignmentsGetOptions` has been removed
- Struct `GuestConfigurationAssignmentsGetResponse` has been removed
- Struct `GuestConfigurationAssignmentsGetResult` has been removed
- Struct `GuestConfigurationAssignmentsListOptions` has been removed
- Struct `GuestConfigurationAssignmentsListResponse` has been removed
- Struct `GuestConfigurationAssignmentsListResult` has been removed
- Struct `GuestConfigurationAssignmentsRGListOptions` has been removed
- Struct `GuestConfigurationAssignmentsRGListResponse` has been removed
- Struct `GuestConfigurationAssignmentsRGListResult` has been removed
- Struct `GuestConfigurationAssignmentsSubscriptionListOptions` has been removed
- Struct `GuestConfigurationAssignmentsSubscriptionListResponse` has been removed
- Struct `GuestConfigurationAssignmentsSubscriptionListResult` has been removed
- Struct `GuestConfigurationHCRPAssignmentReportsClient` has been removed
- Struct `GuestConfigurationHCRPAssignmentReportsGetOptions` has been removed
- Struct `GuestConfigurationHCRPAssignmentReportsGetResponse` has been removed
- Struct `GuestConfigurationHCRPAssignmentReportsGetResult` has been removed
- Struct `GuestConfigurationHCRPAssignmentReportsListOptions` has been removed
- Struct `GuestConfigurationHCRPAssignmentReportsListResponse` has been removed
- Struct `GuestConfigurationHCRPAssignmentReportsListResult` has been removed
- Struct `GuestConfigurationHCRPAssignmentsClient` has been removed
- Struct `GuestConfigurationHCRPAssignmentsCreateOrUpdateOptions` has been removed
- Struct `GuestConfigurationHCRPAssignmentsCreateOrUpdateResponse` has been removed
- Struct `GuestConfigurationHCRPAssignmentsCreateOrUpdateResult` has been removed
- Struct `GuestConfigurationHCRPAssignmentsDeleteOptions` has been removed
- Struct `GuestConfigurationHCRPAssignmentsDeleteResponse` has been removed
- Struct `GuestConfigurationHCRPAssignmentsGetOptions` has been removed
- Struct `GuestConfigurationHCRPAssignmentsGetResponse` has been removed
- Struct `GuestConfigurationHCRPAssignmentsGetResult` has been removed
- Struct `GuestConfigurationHCRPAssignmentsListOptions` has been removed
- Struct `GuestConfigurationHCRPAssignmentsListResponse` has been removed
- Struct `GuestConfigurationHCRPAssignmentsListResult` has been removed
- Struct `GuestConfigurationNavigation` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `Resource` of struct `ProxyResource` has been removed
- Field `ComplianceStatus` of struct `AssignmentReport` has been removed
- Field `OperationType` of struct `AssignmentReport` has been removed
- Field `Assignment` of struct `AssignmentReport` has been removed
- Field `EndTime` of struct `AssignmentReport` has been removed
- Field `StartTime` of struct `AssignmentReport` has been removed
- Field `ReportID` of struct `AssignmentReport` has been removed
- Field `Resources` of struct `AssignmentReport` has been removed
- Field `VM` of struct `AssignmentReport` has been removed

### Features Added

- New function `*CommonAssignmentReport.UnmarshalJSON([]byte) error`
- New function `NewAssignmentsClient(string, azcore.TokenCredential, *arm.ClientOptions) *AssignmentsClient`
- New function `*AssignmentsClient.List(context.Context, string, string, *AssignmentsClientListOptions) (AssignmentsClientListResponse, error)`
- New function `AssignmentReportProperties.MarshalJSON() ([]byte, error)`
- New function `*HCRPAssignmentReportsClient.Get(context.Context, string, string, string, string, *HCRPAssignmentReportsClientGetOptions) (HCRPAssignmentReportsClientGetResponse, error)`
- New function `*HCRPAssignmentsClient.List(context.Context, string, string, *HCRPAssignmentsClientListOptions) (HCRPAssignmentsClientListResponse, error)`
- New function `*HCRPAssignmentReportsClient.List(context.Context, string, string, string, *HCRPAssignmentReportsClientListOptions) (HCRPAssignmentReportsClientListResponse, error)`
- New function `NewAssignmentReportsClient(string, azcore.TokenCredential, *arm.ClientOptions) *AssignmentReportsClient`
- New function `*AssignmentReportsClient.List(context.Context, string, string, string, *AssignmentReportsClientListOptions) (AssignmentReportsClientListResponse, error)`
- New function `*AssignmentsClient.RGList(context.Context, string, *AssignmentsClientRGListOptions) (AssignmentsClientRGListResponse, error)`
- New function `*AssignmentProperties.UnmarshalJSON([]byte) error`
- New function `*HCRPAssignmentsClient.Delete(context.Context, string, string, string, *HCRPAssignmentsClientDeleteOptions) (HCRPAssignmentsClientDeleteResponse, error)`
- New function `*AssignmentsClient.Get(context.Context, string, string, string, *AssignmentsClientGetOptions) (AssignmentsClientGetResponse, error)`
- New function `AssignmentList.MarshalJSON() ([]byte, error)`
- New function `*AssignmentsClient.Delete(context.Context, string, string, string, *AssignmentsClientDeleteOptions) (AssignmentsClientDeleteResponse, error)`
- New function `AssignmentReportList.MarshalJSON() ([]byte, error)`
- New function `AssignmentProperties.MarshalJSON() ([]byte, error)`
- New function `*HCRPAssignmentsClient.Get(context.Context, string, string, string, *HCRPAssignmentsClientGetOptions) (HCRPAssignmentsClientGetResponse, error)`
- New function `*AssignmentReportProperties.UnmarshalJSON([]byte) error`
- New function `NewHCRPAssignmentReportsClient(string, azcore.TokenCredential, *arm.ClientOptions) *HCRPAssignmentReportsClient`
- New function `*AssignmentReportsClient.Get(context.Context, string, string, string, string, *AssignmentReportsClientGetOptions) (AssignmentReportsClientGetResponse, error)`
- New function `*HCRPAssignmentsClient.CreateOrUpdate(context.Context, string, string, string, Assignment, *HCRPAssignmentsClientCreateOrUpdateOptions) (HCRPAssignmentsClientCreateOrUpdateResponse, error)`
- New function `*AssignmentsClient.CreateOrUpdate(context.Context, string, string, string, Assignment, *AssignmentsClientCreateOrUpdateOptions) (AssignmentsClientCreateOrUpdateResponse, error)`
- New function `NewHCRPAssignmentsClient(string, azcore.TokenCredential, *arm.ClientOptions) *HCRPAssignmentsClient`
- New function `Navigation.MarshalJSON() ([]byte, error)`
- New function `CommonAssignmentReport.MarshalJSON() ([]byte, error)`
- New function `*AssignmentsClient.SubscriptionList(context.Context, *AssignmentsClientSubscriptionListOptions) (AssignmentsClientSubscriptionListResponse, error)`
- New struct `Assignment`
- New struct `AssignmentList`
- New struct `AssignmentProperties`
- New struct `AssignmentReportList`
- New struct `AssignmentReportProperties`
- New struct `AssignmentReportsClient`
- New struct `AssignmentReportsClientGetOptions`
- New struct `AssignmentReportsClientGetResponse`
- New struct `AssignmentReportsClientGetResult`
- New struct `AssignmentReportsClientListOptions`
- New struct `AssignmentReportsClientListResponse`
- New struct `AssignmentReportsClientListResult`
- New struct `AssignmentsClient`
- New struct `AssignmentsClientCreateOrUpdateOptions`
- New struct `AssignmentsClientCreateOrUpdateResponse`
- New struct `AssignmentsClientCreateOrUpdateResult`
- New struct `AssignmentsClientDeleteOptions`
- New struct `AssignmentsClientDeleteResponse`
- New struct `AssignmentsClientGetOptions`
- New struct `AssignmentsClientGetResponse`
- New struct `AssignmentsClientGetResult`
- New struct `AssignmentsClientListOptions`
- New struct `AssignmentsClientListResponse`
- New struct `AssignmentsClientListResult`
- New struct `AssignmentsClientRGListOptions`
- New struct `AssignmentsClientRGListResponse`
- New struct `AssignmentsClientRGListResult`
- New struct `AssignmentsClientSubscriptionListOptions`
- New struct `AssignmentsClientSubscriptionListResponse`
- New struct `AssignmentsClientSubscriptionListResult`
- New struct `CommonAssignmentReport`
- New struct `HCRPAssignmentReportsClient`
- New struct `HCRPAssignmentReportsClientGetOptions`
- New struct `HCRPAssignmentReportsClientGetResponse`
- New struct `HCRPAssignmentReportsClientGetResult`
- New struct `HCRPAssignmentReportsClientListOptions`
- New struct `HCRPAssignmentReportsClientListResponse`
- New struct `HCRPAssignmentReportsClientListResult`
- New struct `HCRPAssignmentsClient`
- New struct `HCRPAssignmentsClientCreateOrUpdateOptions`
- New struct `HCRPAssignmentsClientCreateOrUpdateResponse`
- New struct `HCRPAssignmentsClientCreateOrUpdateResult`
- New struct `HCRPAssignmentsClientDeleteOptions`
- New struct `HCRPAssignmentsClientDeleteResponse`
- New struct `HCRPAssignmentsClientGetOptions`
- New struct `HCRPAssignmentsClientGetResponse`
- New struct `HCRPAssignmentsClientGetResult`
- New struct `HCRPAssignmentsClientListOptions`
- New struct `HCRPAssignmentsClientListResponse`
- New struct `HCRPAssignmentsClientListResult`
- New struct `Navigation`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `Name` in struct `AssignmentReport`
- New field `Properties` in struct `AssignmentReport`
- New field `Location` in struct `ProxyResource`
- New field `Name` in struct `ProxyResource`
- New field `ID` in struct `ProxyResource`
- New field `Type` in struct `ProxyResource`


## 0.1.0 (2021-12-07)

- Initial preview release.
