# Release History

## 1.0.0 (2024-06-21)
### Breaking Changes

- Type of `ReportResourcePatch.Properties` has been changed from `*ReportProperties` to `*ReportPatchProperties`
- `CategoryStatusHealthy`, `CategoryStatusUnhealthy` from enum `CategoryStatus` has been removed
- `ControlFamilyStatusHealthy`, `ControlFamilyStatusUnhealthy` from enum `ControlFamilyStatus` has been removed
- `ResourceStatusNotApplicable` from enum `ResourceStatus` has been removed
- Enum `AssessmentSeverity` has been removed
- Enum `CategoryType` has been removed
- Enum `ComplianceState` has been removed
- Enum `ControlFamilyType` has been removed
- Enum `ControlType` has been removed
- Enum `IsPass` has been removed
- Function `*ClientFactory.NewReportsClient` has been removed
- Function `*ClientFactory.NewSnapshotsClient` has been removed
- Function `NewReportsClient` has been removed
- Function `*ReportsClient.NewListPager` has been removed
- Function `NewSnapshotsClient` has been removed
- Function `*SnapshotsClient.NewListPager` has been removed
- Struct `Assessment` has been removed
- Struct `AssessmentResource` has been removed
- Struct `ReportResourceList` has been removed
- Struct `SnapshotResourceList` has been removed
- Field `CategoryType` of struct `Category` has been removed
- Field `ComplianceState`, `ControlType`, `PolicyDescription`, `PolicyDisplayName`, `PolicyID`, `ResourceGroup`, `StatusChangeDate`, `SubscriptionID` of struct `ComplianceReportItem` has been removed
- Field `Assessments`, `ControlShortName`, `ControlType` of struct `Control` has been removed
- Field `FamilyName`, `FamilyStatus`, `FamilyType` of struct `ControlFamily` has been removed
- Field `ID`, `ReportName` of struct `ReportProperties` has been removed
- Field `ResourceName`, `Tags` of struct `ResourceMetadata` has been removed
- Field `ID` of struct `SnapshotProperties` has been removed

### Features Added

- New value `CategoryStatusFailed`, `CategoryStatusNotApplicable`, `CategoryStatusPassed`, `CategoryStatusPendingApproval` added to enum type `CategoryStatus`
- New value `ControlFamilyStatusFailed`, `ControlFamilyStatusNotApplicable`, `ControlFamilyStatusPassed`, `ControlFamilyStatusPendingApproval` added to enum type `ControlFamilyStatus`
- New value `ControlStatusPendingApproval` added to enum type `ControlStatus`
- New value `ProvisioningStateFixing`, `ProvisioningStateVerifying` added to enum type `ProvisioningState`
- New value `ReportStatusReviewing` added to enum type `ReportStatus`
- New enum type `CheckNameAvailabilityReason` with values `CheckNameAvailabilityReasonAlreadyExists`, `CheckNameAvailabilityReasonInvalid`
- New enum type `ContentType` with values `ContentTypeApplicationJSON`
- New enum type `DeliveryStatus` with values `DeliveryStatusFailed`, `DeliveryStatusNotStarted`, `DeliveryStatusSucceeded`
- New enum type `EnableSSLVerification` with values `EnableSSLVerificationFalse`, `EnableSSLVerificationTrue`
- New enum type `EvidenceType` with values `EvidenceTypeAutoCollectedEvidence`, `EvidenceTypeData`, `EvidenceTypeFile`
- New enum type `InputType` with values `InputTypeBoolean`, `InputTypeDate`, `InputTypeEmail`, `InputTypeGroup`, `InputTypeMultiSelectCheckbox`, `InputTypeMultiSelectDropdown`, `InputTypeMultiSelectDropdownCustom`, `InputTypeMultilineText`, `InputTypeNone`, `InputTypeNumber`, `InputTypeSingleSelectDropdown`, `InputTypeSingleSelection`, `InputTypeTelephone`, `InputTypeText`, `InputTypeURL`, `InputTypeUpload`, `InputTypeYearPicker`, `InputTypeYesNoNa`
- New enum type `IsRecommendSolution` with values `IsRecommendSolutionFalse`, `IsRecommendSolutionTrue`
- New enum type `NotificationEvent` with values `NotificationEventAssessmentFailure`, `NotificationEventGenerateSnapshotFailed`, `NotificationEventGenerateSnapshotSuccess`, `NotificationEventReportConfigurationChanges`, `NotificationEventReportDeletion`
- New enum type `ResourceOrigin` with values `ResourceOriginAWS`, `ResourceOriginAzure`, `ResourceOriginGCP`
- New enum type `ResponsibilityEnvironment` with values `ResponsibilityEnvironmentAWS`, `ResponsibilityEnvironmentAzure`, `ResponsibilityEnvironmentGCP`, `ResponsibilityEnvironmentGeneral`
- New enum type `ResponsibilitySeverity` with values `ResponsibilitySeverityHigh`, `ResponsibilitySeverityLow`, `ResponsibilitySeverityMedium`
- New enum type `ResponsibilityStatus` with values `ResponsibilityStatusFailed`, `ResponsibilityStatusNotApplicable`, `ResponsibilityStatusPassed`, `ResponsibilityStatusPendingApproval`
- New enum type `ResponsibilityType` with values `ResponsibilityTypeAutomated`, `ResponsibilityTypeManual`, `ResponsibilityTypeScopedManual`
- New enum type `Result` with values `ResultFailed`, `ResultSucceeded`
- New enum type `Rule` with values `RuleAzureApplication`, `RuleCharLength`, `RuleCreditCardPCI`, `RuleDomains`, `RuleDynamicDropdown`, `RulePreventNonEnglishChar`, `RulePublicSOX`, `RulePublisherVerification`, `RuleRequired`, `RuleURL`, `RuleUSPrivacyShield`, `RuleUrls`, `RuleValidEmail`, `RuleValidGUID`
- New enum type `SendAllEvents` with values `SendAllEventsFalse`, `SendAllEventsTrue`
- New enum type `UpdateWebhookKey` with values `UpdateWebhookKeyFalse`, `UpdateWebhookKeyTrue`
- New enum type `WebhookKeyEnabled` with values `WebhookKeyEnabledFalse`, `WebhookKeyEnabledTrue`
- New enum type `WebhookStatus` with values `WebhookStatusDisabled`, `WebhookStatusEnabled`
- New function `*ClientFactory.NewEvidenceClient() *EvidenceClient`
- New function `*ClientFactory.NewProviderActionsClient() *ProviderActionsClient`
- New function `*ClientFactory.NewScopingConfigurationClient() *ScopingConfigurationClient`
- New function `*ClientFactory.NewWebhookClient() *WebhookClient`
- New function `NewEvidenceClient(azcore.TokenCredential, *arm.ClientOptions) (*EvidenceClient, error)`
- New function `*EvidenceClient.CreateOrUpdate(context.Context, string, string, EvidenceResource, *EvidenceClientCreateOrUpdateOptions) (EvidenceClientCreateOrUpdateResponse, error)`
- New function `*EvidenceClient.Delete(context.Context, string, string, *EvidenceClientDeleteOptions) (EvidenceClientDeleteResponse, error)`
- New function `*EvidenceClient.Download(context.Context, string, string, EvidenceFileDownloadRequest, *EvidenceClientDownloadOptions) (EvidenceClientDownloadResponse, error)`
- New function `*EvidenceClient.Get(context.Context, string, string, *EvidenceClientGetOptions) (EvidenceClientGetResponse, error)`
- New function `*EvidenceClient.NewListByReportPager(string, *EvidenceClientListByReportOptions) *runtime.Pager[EvidenceClientListByReportResponse]`
- New function `NewProviderActionsClient(azcore.TokenCredential, *arm.ClientOptions) (*ProviderActionsClient, error)`
- New function `*ProviderActionsClient.CheckNameAvailability(context.Context, CheckNameAvailabilityRequest, *ProviderActionsClientCheckNameAvailabilityOptions) (ProviderActionsClientCheckNameAvailabilityResponse, error)`
- New function `*ProviderActionsClient.GetCollectionCount(context.Context, GetCollectionCountRequest, *ProviderActionsClientGetCollectionCountOptions) (ProviderActionsClientGetCollectionCountResponse, error)`
- New function `*ProviderActionsClient.GetOverviewStatus(context.Context, GetOverviewStatusRequest, *ProviderActionsClientGetOverviewStatusOptions) (ProviderActionsClientGetOverviewStatusResponse, error)`
- New function `*ProviderActionsClient.ListInUseStorageAccounts(context.Context, ListInUseStorageAccountsRequest, *ProviderActionsClientListInUseStorageAccountsOptions) (ProviderActionsClientListInUseStorageAccountsResponse, error)`
- New function `*ProviderActionsClient.BeginOnboard(context.Context, OnboardRequest, *ProviderActionsClientBeginOnboardOptions) (*runtime.Poller[ProviderActionsClientOnboardResponse], error)`
- New function `*ProviderActionsClient.BeginTriggerEvaluation(context.Context, TriggerEvaluationRequest, *ProviderActionsClientBeginTriggerEvaluationOptions) (*runtime.Poller[ProviderActionsClientTriggerEvaluationResponse], error)`
- New function `*ReportClient.BeginFix(context.Context, string, *ReportClientBeginFixOptions) (*runtime.Poller[ReportClientFixResponse], error)`
- New function `*ReportClient.GetScopingQuestions(context.Context, string, *ReportClientGetScopingQuestionsOptions) (ReportClientGetScopingQuestionsResponse, error)`
- New function `*ReportClient.NewListPager(*ReportClientListOptions) *runtime.Pager[ReportClientListResponse]`
- New function `*ReportClient.NestedResourceCheckNameAvailability(context.Context, string, CheckNameAvailabilityRequest, *ReportClientNestedResourceCheckNameAvailabilityOptions) (ReportClientNestedResourceCheckNameAvailabilityResponse, error)`
- New function `*ReportClient.BeginSyncCertRecord(context.Context, string, SyncCertRecordRequest, *ReportClientBeginSyncCertRecordOptions) (*runtime.Poller[ReportClientSyncCertRecordResponse], error)`
- New function `*ReportClient.BeginVerify(context.Context, string, *ReportClientBeginVerifyOptions) (*runtime.Poller[ReportClientVerifyResponse], error)`
- New function `NewScopingConfigurationClient(azcore.TokenCredential, *arm.ClientOptions) (*ScopingConfigurationClient, error)`
- New function `*ScopingConfigurationClient.CreateOrUpdate(context.Context, string, string, ScopingConfigurationResource, *ScopingConfigurationClientCreateOrUpdateOptions) (ScopingConfigurationClientCreateOrUpdateResponse, error)`
- New function `*ScopingConfigurationClient.Delete(context.Context, string, string, *ScopingConfigurationClientDeleteOptions) (ScopingConfigurationClientDeleteResponse, error)`
- New function `*ScopingConfigurationClient.Get(context.Context, string, string, *ScopingConfigurationClientGetOptions) (ScopingConfigurationClientGetResponse, error)`
- New function `*ScopingConfigurationClient.NewListPager(string, *ScopingConfigurationClientListOptions) *runtime.Pager[ScopingConfigurationClientListResponse]`
- New function `*SnapshotClient.NewListPager(string, *SnapshotClientListOptions) *runtime.Pager[SnapshotClientListResponse]`
- New function `NewWebhookClient(azcore.TokenCredential, *arm.ClientOptions) (*WebhookClient, error)`
- New function `*WebhookClient.CreateOrUpdate(context.Context, string, string, WebhookResource, *WebhookClientCreateOrUpdateOptions) (WebhookClientCreateOrUpdateResponse, error)`
- New function `*WebhookClient.Delete(context.Context, string, string, *WebhookClientDeleteOptions) (WebhookClientDeleteResponse, error)`
- New function `*WebhookClient.Get(context.Context, string, string, *WebhookClientGetOptions) (WebhookClientGetResponse, error)`
- New function `*WebhookClient.NewListPager(string, *WebhookClientListOptions) *runtime.Pager[WebhookClientListResponse]`
- New function `*WebhookClient.Update(context.Context, string, string, WebhookResourcePatch, *WebhookClientUpdateOptions) (WebhookClientUpdateResponse, error)`
- New struct `CertSyncRecord`
- New struct `CheckNameAvailabilityRequest`
- New struct `CheckNameAvailabilityResponse`
- New struct `ControlSyncRecord`
- New struct `EvidenceFileDownloadRequest`
- New struct `EvidenceFileDownloadResponse`
- New struct `EvidenceFileDownloadResponseEvidenceFile`
- New struct `EvidenceProperties`
- New struct `EvidenceResource`
- New struct `EvidenceResourceListResult`
- New struct `GetCollectionCountRequest`
- New struct `GetCollectionCountResponse`
- New struct `GetOverviewStatusRequest`
- New struct `GetOverviewStatusResponse`
- New struct `ListInUseStorageAccountsRequest`
- New struct `ListInUseStorageAccountsResponse`
- New struct `OnboardRequest`
- New struct `OnboardResponse`
- New struct `QuickAssessment`
- New struct `Recommendation`
- New struct `RecommendationSolution`
- New struct `ReportFixResult`
- New struct `ReportPatchProperties`
- New struct `ReportResourceListResult`
- New struct `ReportVerificationResult`
- New struct `Responsibility`
- New struct `ResponsibilityResource`
- New struct `ScopingAnswer`
- New struct `ScopingConfigurationProperties`
- New struct `ScopingConfigurationResource`
- New struct `ScopingConfigurationResourceListResult`
- New struct `ScopingQuestion`
- New struct `ScopingQuestions`
- New struct `SnapshotResourceListResult`
- New struct `StatusItem`
- New struct `StorageInfo`
- New struct `SyncCertRecordRequest`
- New struct `SyncCertRecordResponse`
- New struct `TriggerEvaluationProperty`
- New struct `TriggerEvaluationRequest`
- New struct `TriggerEvaluationResponse`
- New struct `WebhookProperties`
- New struct `WebhookResource`
- New struct `WebhookResourceListResult`
- New struct `WebhookResourcePatch`
- New field `ControlFamilyName`, `ControlStatus`, `ResourceOrigin`, `ResourceStatus`, `ResourceStatusChangeDate`, `ResponsibilityDescription`, `ResponsibilityTitle` in struct `ComplianceReportItem`
- New field `ControlName`, `Responsibilities` in struct `Control`
- New field `ControlFamilyName`, `ControlFamilyStatus` in struct `ControlFamily`
- New field `NotApplicableCount`, `PendingCount` in struct `OverviewStatus`
- New field `CertRecords`, `Errors`, `StorageInfo` in struct `ReportProperties`
- New field `AccountID`, `ResourceOrigin` in struct `ResourceMetadata`


## 0.3.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.2.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.1.0 (2022-10-26)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcomplianceautomation/armappcomplianceautomation` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).