# Release History

## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*RecommendationsClient.List` has been removed
- Function `*ConfigurationsClient.ListBySubscription` has been removed
- Function `*RecommendationMetadataClient.List` has been removed
- Function `*ConfigurationsClient.ListByResourceGroup` has been removed
- Function `*SuppressionsClient.List` has been removed
- Function `*OperationsClient.List` has been removed

### Features Added

- New function `*SuppressionsClient.NewListPager(*SuppressionsClientListOptions) *runtime.Pager[SuppressionsClientListResponse]`
- New function `*RecommendationMetadataClient.NewListPager(*RecommendationMetadataClientListOptions) *runtime.Pager[RecommendationMetadataClientListResponse]`
- New function `*ConfigurationsClient.NewListByResourceGroupPager(string, *ConfigurationsClientListByResourceGroupOptions) *runtime.Pager[ConfigurationsClientListByResourceGroupResponse]`
- New function `*RecommendationsClient.NewListPager(*RecommendationsClientListOptions) *runtime.Pager[RecommendationsClientListResponse]`
- New function `*ConfigurationsClient.NewListBySubscriptionPager(*ConfigurationsClientListBySubscriptionOptions) *runtime.Pager[ConfigurationsClientListBySubscriptionResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `NewSuppressionsClient` return value(s) have been changed from `(*SuppressionsClient)` to `(*SuppressionsClient, error)`
- Function `*SuppressionsClient.List` return value(s) have been changed from `(*SuppressionsClientListPager)` to `(*runtime.Pager[SuppressionsClientListResponse])`
- Function `NewConfigurationsClient` return value(s) have been changed from `(*ConfigurationsClient)` to `(*ConfigurationsClient, error)`
- Function `*RecommendationsClient.List` return value(s) have been changed from `(*RecommendationsClientListPager)` to `(*runtime.Pager[RecommendationsClientListResponse])`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `*ConfigurationsClient.ListByResourceGroup` parameter(s) have been changed from `(context.Context, string, *ConfigurationsClientListByResourceGroupOptions)` to `(string, *ConfigurationsClientListByResourceGroupOptions)`
- Function `*ConfigurationsClient.ListByResourceGroup` return value(s) have been changed from `(ConfigurationsClientListByResourceGroupResponse, error)` to `(*runtime.Pager[ConfigurationsClientListByResourceGroupResponse])`
- Function `*ConfigurationsClient.ListBySubscription` return value(s) have been changed from `(*ConfigurationsClientListBySubscriptionPager)` to `(*runtime.Pager[ConfigurationsClientListBySubscriptionResponse])`
- Function `NewRecommendationsClient` return value(s) have been changed from `(*RecommendationsClient)` to `(*RecommendationsClient, error)`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*RecommendationMetadataClient.List` return value(s) have been changed from `(*RecommendationMetadataClientListPager)` to `(*runtime.Pager[RecommendationMetadataClientListResponse])`
- Function `NewRecommendationMetadataClient` return value(s) have been changed from `(*RecommendationMetadataClient)` to `(*RecommendationMetadataClient, error)`
- Type of `RecommendationProperties.Remediation` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Type of `RecommendationProperties.Actions` has been changed from `[]map[string]map[string]interface{}` to `[]map[string]interface{}`
- Type of `RecommendationProperties.ExposedMetadataProperties` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Type of `RecommendationProperties.Metadata` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Type of `ResourceMetadata.Action` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Function `*OperationsClientListPager.Err` has been removed
- Function `*SuppressionsClientListPager.Err` has been removed
- Function `*RecommendationMetadataClientListPager.NextPage` has been removed
- Function `*ConfigurationsClientListBySubscriptionPager.Err` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `Impact.ToPtr` has been removed
- Function `*RecommendationMetadataClientListPager.Err` has been removed
- Function `Risk.ToPtr` has been removed
- Function `*SuppressionsClientListPager.NextPage` has been removed
- Function `*RecommendationMetadataClientListPager.PageResponse` has been removed
- Function `CPUThreshold.ToPtr` has been removed
- Function `*RecommendationsClientListPager.Err` has been removed
- Function `*RecommendationsClientListPager.NextPage` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*ConfigurationsClientListBySubscriptionPager.PageResponse` has been removed
- Function `*SuppressionsClientListPager.PageResponse` has been removed
- Function `*RecommendationsClientListPager.PageResponse` has been removed
- Function `ConfigurationName.ToPtr` has been removed
- Function `Category.ToPtr` has been removed
- Function `*ConfigurationsClientListBySubscriptionPager.NextPage` has been removed
- Function `DigestConfigState.ToPtr` has been removed
- Function `Scenario.ToPtr` has been removed
- Struct `ConfigurationsClientCreateInResourceGroupResult` has been removed
- Struct `ConfigurationsClientCreateInSubscriptionResult` has been removed
- Struct `ConfigurationsClientListByResourceGroupResult` has been removed
- Struct `ConfigurationsClientListBySubscriptionPager` has been removed
- Struct `ConfigurationsClientListBySubscriptionResult` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `RecommendationMetadataClientGetResult` has been removed
- Struct `RecommendationMetadataClientListPager` has been removed
- Struct `RecommendationMetadataClientListResult` has been removed
- Struct `RecommendationsClientGenerateResult` has been removed
- Struct `RecommendationsClientGetResult` has been removed
- Struct `RecommendationsClientListPager` has been removed
- Struct `RecommendationsClientListResult` has been removed
- Struct `SuppressionsClientCreateResult` has been removed
- Struct `SuppressionsClientGetResult` has been removed
- Struct `SuppressionsClientListPager` has been removed
- Struct `SuppressionsClientListResult` has been removed
- Field `RecommendationsClientGetResult` of struct `RecommendationsClientGetResponse` has been removed
- Field `RawResponse` of struct `RecommendationsClientGetResponse` has been removed
- Field `RecommendationMetadataClientListResult` of struct `RecommendationMetadataClientListResponse` has been removed
- Field `RawResponse` of struct `RecommendationMetadataClientListResponse` has been removed
- Field `RecommendationsClientListResult` of struct `RecommendationsClientListResponse` has been removed
- Field `RawResponse` of struct `RecommendationsClientListResponse` has been removed
- Field `SuppressionsClientGetResult` of struct `SuppressionsClientGetResponse` has been removed
- Field `RawResponse` of struct `SuppressionsClientGetResponse` has been removed
- Field `ConfigurationsClientCreateInSubscriptionResult` of struct `ConfigurationsClientCreateInSubscriptionResponse` has been removed
- Field `RawResponse` of struct `ConfigurationsClientCreateInSubscriptionResponse` has been removed
- Field `RawResponse` of struct `RecommendationsClientGetGenerateStatusResponse` has been removed
- Field `SuppressionsClientListResult` of struct `SuppressionsClientListResponse` has been removed
- Field `RawResponse` of struct `SuppressionsClientListResponse` has been removed
- Field `SuppressionsClientCreateResult` of struct `SuppressionsClientCreateResponse` has been removed
- Field `RawResponse` of struct `SuppressionsClientCreateResponse` has been removed
- Field `RecommendationMetadataClientGetResult` of struct `RecommendationMetadataClientGetResponse` has been removed
- Field `RawResponse` of struct `RecommendationMetadataClientGetResponse` has been removed
- Field `ConfigurationsClientListBySubscriptionResult` of struct `ConfigurationsClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `ConfigurationsClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `SuppressionsClientDeleteResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `ConfigurationsClientListByResourceGroupResult` of struct `ConfigurationsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `ConfigurationsClientListByResourceGroupResponse` has been removed
- Field `RecommendationsClientGenerateResult` of struct `RecommendationsClientGenerateResponse` has been removed
- Field `RawResponse` of struct `RecommendationsClientGenerateResponse` has been removed
- Field `ConfigurationsClientCreateInResourceGroupResult` of struct `ConfigurationsClientCreateInResourceGroupResponse` has been removed
- Field `RawResponse` of struct `ConfigurationsClientCreateInResourceGroupResponse` has been removed

### Features Added

- New anonymous field `MetadataEntityListResult` in struct `RecommendationMetadataClientListResponse`
- New anonymous field `ConfigurationListResult` in struct `ConfigurationsClientListBySubscriptionResponse`
- New anonymous field `ResourceRecommendationBaseListResult` in struct `RecommendationsClientListResponse`
- New anonymous field `SuppressionContract` in struct `SuppressionsClientGetResponse`
- New anonymous field `ResourceRecommendationBase` in struct `RecommendationsClientGetResponse`
- New anonymous field `MetadataEntity` in struct `RecommendationMetadataClientGetResponse`
- New anonymous field `ConfigData` in struct `ConfigurationsClientCreateInResourceGroupResponse`
- New anonymous field `SuppressionContract` in struct `SuppressionsClientCreateResponse`
- New field `Location` in struct `RecommendationsClientGenerateResponse`
- New field `RetryAfter` in struct `RecommendationsClientGenerateResponse`
- New anonymous field `ConfigurationListResult` in struct `ConfigurationsClientListByResourceGroupResponse`
- New anonymous field `SuppressionContractListResult` in struct `SuppressionsClientListResponse`
- New anonymous field `OperationEntityListResult` in struct `OperationsClientListResponse`
- New anonymous field `ConfigData` in struct `ConfigurationsClientCreateInSubscriptionResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*SuppressionsClient.Create` parameter(s) have been changed from `(context.Context, string, string, string, SuppressionContract, *SuppressionsCreateOptions)` to `(context.Context, string, string, string, SuppressionContract, *SuppressionsClientCreateOptions)`
- Function `*SuppressionsClient.Create` return value(s) have been changed from `(SuppressionsCreateResponse, error)` to `(SuppressionsClientCreateResponse, error)`
- Function `*RecommendationMetadataClient.List` parameter(s) have been changed from `(*RecommendationMetadataListOptions)` to `(*RecommendationMetadataClientListOptions)`
- Function `*RecommendationMetadataClient.List` return value(s) have been changed from `(*RecommendationMetadataListPager)` to `(*RecommendationMetadataClientListPager)`
- Function `*SuppressionsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *SuppressionsGetOptions)` to `(context.Context, string, string, string, *SuppressionsClientGetOptions)`
- Function `*SuppressionsClient.Get` return value(s) have been changed from `(SuppressionsGetResponse, error)` to `(SuppressionsClientGetResponse, error)`
- Function `*RecommendationMetadataClient.Get` parameter(s) have been changed from `(context.Context, string, *RecommendationMetadataGetOptions)` to `(context.Context, string, *RecommendationMetadataClientGetOptions)`
- Function `*RecommendationMetadataClient.Get` return value(s) have been changed from `(RecommendationMetadataGetResponse, error)` to `(RecommendationMetadataClientGetResponse, error)`
- Function `*SuppressionsClient.List` parameter(s) have been changed from `(*SuppressionsListOptions)` to `(*SuppressionsClientListOptions)`
- Function `*SuppressionsClient.List` return value(s) have been changed from `(*SuppressionsListPager)` to `(*SuppressionsClientListPager)`
- Function `*ConfigurationsClient.CreateInResourceGroup` parameter(s) have been changed from `(context.Context, ConfigurationName, string, ConfigData, *ConfigurationsCreateInResourceGroupOptions)` to `(context.Context, ConfigurationName, string, ConfigData, *ConfigurationsClientCreateInResourceGroupOptions)`
- Function `*ConfigurationsClient.CreateInResourceGroup` return value(s) have been changed from `(ConfigurationsCreateInResourceGroupResponse, error)` to `(ConfigurationsClientCreateInResourceGroupResponse, error)`
- Function `*RecommendationsClient.Generate` parameter(s) have been changed from `(context.Context, *RecommendationsGenerateOptions)` to `(context.Context, *RecommendationsClientGenerateOptions)`
- Function `*RecommendationsClient.Generate` return value(s) have been changed from `(RecommendationsGenerateResponse, error)` to `(RecommendationsClientGenerateResponse, error)`
- Function `*ConfigurationsClient.ListBySubscription` parameter(s) have been changed from `(*ConfigurationsListBySubscriptionOptions)` to `(*ConfigurationsClientListBySubscriptionOptions)`
- Function `*ConfigurationsClient.ListBySubscription` return value(s) have been changed from `(*ConfigurationsListBySubscriptionPager)` to `(*ConfigurationsClientListBySubscriptionPager)`
- Function `*SuppressionsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, string, *SuppressionsDeleteOptions)` to `(context.Context, string, string, string, *SuppressionsClientDeleteOptions)`
- Function `*SuppressionsClient.Delete` return value(s) have been changed from `(SuppressionsDeleteResponse, error)` to `(SuppressionsClientDeleteResponse, error)`
- Function `*RecommendationsClient.GetGenerateStatus` parameter(s) have been changed from `(context.Context, string, *RecommendationsGetGenerateStatusOptions)` to `(context.Context, string, *RecommendationsClientGetGenerateStatusOptions)`
- Function `*RecommendationsClient.GetGenerateStatus` return value(s) have been changed from `(RecommendationsGetGenerateStatusResponse, error)` to `(RecommendationsClientGetGenerateStatusResponse, error)`
- Function `*ConfigurationsClient.CreateInSubscription` parameter(s) have been changed from `(context.Context, ConfigurationName, ConfigData, *ConfigurationsCreateInSubscriptionOptions)` to `(context.Context, ConfigurationName, ConfigData, *ConfigurationsClientCreateInSubscriptionOptions)`
- Function `*ConfigurationsClient.CreateInSubscription` return value(s) have been changed from `(ConfigurationsCreateInSubscriptionResponse, error)` to `(ConfigurationsClientCreateInSubscriptionResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*RecommendationsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *RecommendationsGetOptions)` to `(context.Context, string, string, *RecommendationsClientGetOptions)`
- Function `*RecommendationsClient.Get` return value(s) have been changed from `(RecommendationsGetResponse, error)` to `(RecommendationsClientGetResponse, error)`
- Function `*ConfigurationsClient.ListByResourceGroup` parameter(s) have been changed from `(context.Context, string, *ConfigurationsListByResourceGroupOptions)` to `(context.Context, string, *ConfigurationsClientListByResourceGroupOptions)`
- Function `*ConfigurationsClient.ListByResourceGroup` return value(s) have been changed from `(ConfigurationsListByResourceGroupResponse, error)` to `(ConfigurationsClientListByResourceGroupResponse, error)`
- Function `*RecommendationsClient.List` parameter(s) have been changed from `(*RecommendationsListOptions)` to `(*RecommendationsClientListOptions)`
- Function `*RecommendationsClient.List` return value(s) have been changed from `(*RecommendationsListPager)` to `(*RecommendationsClientListPager)`
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*ConfigurationsListBySubscriptionPager.PageResponse` has been removed
- Function `*ConfigurationsListBySubscriptionPager.Err` has been removed
- Function `*RecommendationMetadataListPager.PageResponse` has been removed
- Function `*SuppressionsListPager.Err` has been removed
- Function `ARMErrorResponseBody.Error` has been removed
- Function `*SuppressionsListPager.PageResponse` has been removed
- Function `*RecommendationsListPager.NextPage` has been removed
- Function `*RecommendationMetadataListPager.NextPage` has been removed
- Function `ArmErrorResponse.Error` has been removed
- Function `*RecommendationMetadataListPager.Err` has been removed
- Function `*ConfigurationsListBySubscriptionPager.NextPage` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*SuppressionsListPager.NextPage` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*RecommendationsListPager.Err` has been removed
- Function `*RecommendationsListPager.PageResponse` has been removed
- Struct `ConfigurationsCreateInResourceGroupOptions` has been removed
- Struct `ConfigurationsCreateInResourceGroupResponse` has been removed
- Struct `ConfigurationsCreateInResourceGroupResult` has been removed
- Struct `ConfigurationsCreateInSubscriptionOptions` has been removed
- Struct `ConfigurationsCreateInSubscriptionResponse` has been removed
- Struct `ConfigurationsCreateInSubscriptionResult` has been removed
- Struct `ConfigurationsListByResourceGroupOptions` has been removed
- Struct `ConfigurationsListByResourceGroupResponse` has been removed
- Struct `ConfigurationsListByResourceGroupResult` has been removed
- Struct `ConfigurationsListBySubscriptionOptions` has been removed
- Struct `ConfigurationsListBySubscriptionPager` has been removed
- Struct `ConfigurationsListBySubscriptionResponse` has been removed
- Struct `ConfigurationsListBySubscriptionResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `RecommendationMetadataGetOptions` has been removed
- Struct `RecommendationMetadataGetResponse` has been removed
- Struct `RecommendationMetadataGetResult` has been removed
- Struct `RecommendationMetadataListOptions` has been removed
- Struct `RecommendationMetadataListPager` has been removed
- Struct `RecommendationMetadataListResponse` has been removed
- Struct `RecommendationMetadataListResult` has been removed
- Struct `RecommendationsGenerateOptions` has been removed
- Struct `RecommendationsGenerateResponse` has been removed
- Struct `RecommendationsGenerateResult` has been removed
- Struct `RecommendationsGetGenerateStatusOptions` has been removed
- Struct `RecommendationsGetGenerateStatusResponse` has been removed
- Struct `RecommendationsGetOptions` has been removed
- Struct `RecommendationsGetResponse` has been removed
- Struct `RecommendationsGetResult` has been removed
- Struct `RecommendationsListOptions` has been removed
- Struct `RecommendationsListPager` has been removed
- Struct `RecommendationsListResponse` has been removed
- Struct `RecommendationsListResult` has been removed
- Struct `SuppressionsCreateOptions` has been removed
- Struct `SuppressionsCreateResponse` has been removed
- Struct `SuppressionsCreateResult` has been removed
- Struct `SuppressionsDeleteOptions` has been removed
- Struct `SuppressionsDeleteResponse` has been removed
- Struct `SuppressionsGetOptions` has been removed
- Struct `SuppressionsGetResponse` has been removed
- Struct `SuppressionsGetResult` has been removed
- Struct `SuppressionsListOptions` has been removed
- Struct `SuppressionsListPager` has been removed
- Struct `SuppressionsListResponse` has been removed
- Struct `SuppressionsListResult` has been removed
- Field `InnerError` of struct `ArmErrorResponse` has been removed
- Field `Resource` of struct `SuppressionContract` has been removed
- Field `Resource` of struct `ConfigData` has been removed
- Field `Resource` of struct `ResourceRecommendationBase` has been removed

### Features Added

- New function `*SuppressionsClientListPager.PageResponse() SuppressionsClientListResponse`
- New function `*OperationsClientListPager.Err() error`
- New function `*ConfigurationsClientListBySubscriptionPager.PageResponse() ConfigurationsClientListBySubscriptionResponse`
- New function `*RecommendationMetadataClientListPager.NextPage(context.Context) bool`
- New function `*RecommendationMetadataClientListPager.PageResponse() RecommendationMetadataClientListResponse`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*SuppressionsClientListPager.Err() error`
- New function `*ConfigurationsClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `*RecommendationsClientListPager.Err() error`
- New function `*RecommendationMetadataClientListPager.Err() error`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*ConfigurationsClientListBySubscriptionPager.Err() error`
- New function `*RecommendationsClientListPager.NextPage(context.Context) bool`
- New function `*RecommendationsClientListPager.PageResponse() RecommendationsClientListResponse`
- New function `*SuppressionsClientListPager.NextPage(context.Context) bool`
- New struct `ConfigurationsClientCreateInResourceGroupOptions`
- New struct `ConfigurationsClientCreateInResourceGroupResponse`
- New struct `ConfigurationsClientCreateInResourceGroupResult`
- New struct `ConfigurationsClientCreateInSubscriptionOptions`
- New struct `ConfigurationsClientCreateInSubscriptionResponse`
- New struct `ConfigurationsClientCreateInSubscriptionResult`
- New struct `ConfigurationsClientListByResourceGroupOptions`
- New struct `ConfigurationsClientListByResourceGroupResponse`
- New struct `ConfigurationsClientListByResourceGroupResult`
- New struct `ConfigurationsClientListBySubscriptionOptions`
- New struct `ConfigurationsClientListBySubscriptionPager`
- New struct `ConfigurationsClientListBySubscriptionResponse`
- New struct `ConfigurationsClientListBySubscriptionResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `RecommendationMetadataClientGetOptions`
- New struct `RecommendationMetadataClientGetResponse`
- New struct `RecommendationMetadataClientGetResult`
- New struct `RecommendationMetadataClientListOptions`
- New struct `RecommendationMetadataClientListPager`
- New struct `RecommendationMetadataClientListResponse`
- New struct `RecommendationMetadataClientListResult`
- New struct `RecommendationsClientGenerateOptions`
- New struct `RecommendationsClientGenerateResponse`
- New struct `RecommendationsClientGenerateResult`
- New struct `RecommendationsClientGetGenerateStatusOptions`
- New struct `RecommendationsClientGetGenerateStatusResponse`
- New struct `RecommendationsClientGetOptions`
- New struct `RecommendationsClientGetResponse`
- New struct `RecommendationsClientGetResult`
- New struct `RecommendationsClientListOptions`
- New struct `RecommendationsClientListPager`
- New struct `RecommendationsClientListResponse`
- New struct `RecommendationsClientListResult`
- New struct `SuppressionsClientCreateOptions`
- New struct `SuppressionsClientCreateResponse`
- New struct `SuppressionsClientCreateResult`
- New struct `SuppressionsClientDeleteOptions`
- New struct `SuppressionsClientDeleteResponse`
- New struct `SuppressionsClientGetOptions`
- New struct `SuppressionsClientGetResponse`
- New struct `SuppressionsClientGetResult`
- New struct `SuppressionsClientListOptions`
- New struct `SuppressionsClientListPager`
- New struct `SuppressionsClientListResponse`
- New struct `SuppressionsClientListResult`
- New field `ID` in struct `ResourceRecommendationBase`
- New field `Name` in struct `ResourceRecommendationBase`
- New field `Type` in struct `ResourceRecommendationBase`
- New field `Error` in struct `ArmErrorResponse`
- New field `ID` in struct `ConfigData`
- New field `Name` in struct `ConfigData`
- New field `Type` in struct `ConfigData`
- New field `Type` in struct `SuppressionContract`
- New field `ID` in struct `SuppressionContract`
- New field `Name` in struct `SuppressionContract`


## 0.1.0 (2021-11-16)

- Initial preview release.
