# Release History

## 2.0.0-beta.4 (2024-04-26)
### Breaking Changes

- Function `*DiscoverySolutionClient.NewListPager` parameter(s) have been changed from `(string, *DiscoverySolutionClientListOptions)` to `(*DiscoverySolutionClientListOptions)`

### Features Added

- New value `QuestionTypeDateTimePicker`, `QuestionTypeMultiSelect` added to enum type `QuestionType`
- New value `SolutionTypeSelfHelp`, `SolutionTypeTroubleshooters` added to enum type `SolutionType`
- New value `TypeInput` added to enum type `Type`
- New enum type `ValidationScope` with values `ValidationScopeGUIDFormat`, `ValidationScopeIPAddressFormat`, `ValidationScopeNone`, `ValidationScopeNumberOnlyFormat`, `ValidationScopeURLFormat`
- New function `*ClientFactory.NewDiscoverySolutionNLPSubscriptionScopeClient() *DiscoverySolutionNLPSubscriptionScopeClient`
- New function `*ClientFactory.NewDiscoverySolutionNLPTenantScopeClient() *DiscoverySolutionNLPTenantScopeClient`
- New function `*ClientFactory.NewSimplifiedSolutionsClient() *SimplifiedSolutionsClient`
- New function `*ClientFactory.NewSolutionSelfHelpClient() *SolutionSelfHelpClient`
- New function `NewDiscoverySolutionNLPSubscriptionScopeClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DiscoverySolutionNLPSubscriptionScopeClient, error)`
- New function `*DiscoverySolutionNLPSubscriptionScopeClient.Post(context.Context, *DiscoverySolutionNLPSubscriptionScopeClientPostOptions) (DiscoverySolutionNLPSubscriptionScopeClientPostResponse, error)`
- New function `NewDiscoverySolutionNLPTenantScopeClient(azcore.TokenCredential, *arm.ClientOptions) (*DiscoverySolutionNLPTenantScopeClient, error)`
- New function `*DiscoverySolutionNLPTenantScopeClient.Post(context.Context, *DiscoverySolutionNLPTenantScopeClientPostOptions) (DiscoverySolutionNLPTenantScopeClientPostResponse, error)`
- New function `NewSimplifiedSolutionsClient(azcore.TokenCredential, *arm.ClientOptions) (*SimplifiedSolutionsClient, error)`
- New function `*SimplifiedSolutionsClient.BeginCreate(context.Context, string, string, SimplifiedSolutionsResource, *SimplifiedSolutionsClientBeginCreateOptions) (*runtime.Poller[SimplifiedSolutionsClientCreateResponse], error)`
- New function `*SimplifiedSolutionsClient.Get(context.Context, string, string, *SimplifiedSolutionsClientGetOptions) (SimplifiedSolutionsClientGetResponse, error)`
- New function `*SolutionClient.WarmUp(context.Context, string, string, *SolutionClientWarmUpOptions) (SolutionClientWarmUpResponse, error)`
- New function `NewSolutionSelfHelpClient(azcore.TokenCredential, *arm.ClientOptions) (*SolutionSelfHelpClient, error)`
- New function `*SolutionSelfHelpClient.Get(context.Context, string, *SolutionSelfHelpClientGetOptions) (SolutionSelfHelpClientGetResponse, error)`
- New struct `ClassificationService`
- New struct `DiscoveryNlpRequest`
- New struct `DiscoveryNlpResponse`
- New struct `NlpSolutions`
- New struct `ReplacementMapsSelfHelp`
- New struct `SectionSelfHelp`
- New struct `SimplifiedSolutionsResource`
- New struct `SimplifiedSolutionsResourceProperties`
- New struct `SolutionNlpMetadataResource`
- New struct `SolutionResourceSelfHelp`
- New struct `SolutionWarmUpRequestBody`
- New struct `SolutionsResourcePropertiesSelfHelp`
- New field `Status`, `Version` in struct `AutomatedCheckResult`
- New field `ValidationScope` in struct `ResponseValidationProperties`
- New field `EstimatedCompletionTime` in struct `SolutionsDiagnostic`
- New field `QuestionTitle` in struct `StepInput`


## 2.0.0-beta.3 (2023-12-22)
### Breaking Changes

- Type of `StepInput.QuestionType` has been changed from `*string` to `*QuestionType`

### Features Added

- New value `DiagnosticProvisioningStateRunning` added to enum type `DiagnosticProvisioningState`
- New value `SolutionProvisioningStatePartialComplete`, `SolutionProvisioningStateRunning` added to enum type `SolutionProvisioningState`
- New field `SystemData` in struct `SolutionResource`


## 2.0.0-beta.2 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0-beta.1 (2023-10-27)
### Breaking Changes

- Type of `DiagnosticResourceProperties.ProvisioningState` has been changed from `*ProvisioningState` to `*DiagnosticProvisioningState`
- Type of `SolutionMetadataProperties.SolutionType` has been changed from `*string` to `*SolutionType`
- Type of `SolutionMetadataResource.Properties` has been changed from `*SolutionMetadataProperties` to `*Solutions`
- Enum `ProvisioningState` has been removed
- Function `*DiagnosticsClient.CheckNameAvailability` has been removed
- Field `RequiredParameterSets` of struct `SolutionMetadataProperties` has been removed

### Features Added

- New enum type `AggregationType` with values `AggregationTypeAvg`, `AggregationTypeCount`, `AggregationTypeMax`, `AggregationTypeMin`, `AggregationTypeSum`
- New enum type `AutomatedCheckResultType` with values `AutomatedCheckResultTypeError`, `AutomatedCheckResultTypeInformation`, `AutomatedCheckResultTypeSuccess`, `AutomatedCheckResultTypeWarning`
- New enum type `Confidence` with values `ConfidenceHigh`, `ConfidenceLow`, `ConfidenceMedium`
- New enum type `DiagnosticProvisioningState` with values `DiagnosticProvisioningStateCanceled`, `DiagnosticProvisioningStateFailed`, `DiagnosticProvisioningStatePartialComplete`, `DiagnosticProvisioningStateSucceeded`
- New enum type `ExecutionStatus` with values `ExecutionStatusFailed`, `ExecutionStatusRunning`, `ExecutionStatusSuccess`, `ExecutionStatusWarning`
- New enum type `Name` with values `NameProblemClassificationID`, `NameReplacementKey`, `NameSolutionID`
- New enum type `QuestionContentType` with values `QuestionContentTypeHTML`, `QuestionContentTypeMarkdown`, `QuestionContentTypeText`
- New enum type `QuestionType` with values `QuestionTypeDropdown`, `QuestionTypeMultiLineInfoBox`, `QuestionTypeRadioButton`, `QuestionTypeTextInput`
- New enum type `ResultType` with values `ResultTypeCommunity`, `ResultTypeDocumentation`
- New enum type `SolutionProvisioningState` with values `SolutionProvisioningStateCanceled`, `SolutionProvisioningStateFailed`, `SolutionProvisioningStateSucceeded`
- New enum type `SolutionType` with values `SolutionTypeDiagnostics`, `SolutionTypeSolutions`
- New enum type `TroubleshooterProvisioningState` with values `TroubleshooterProvisioningStateAutoContinue`, `TroubleshooterProvisioningStateCanceled`, `TroubleshooterProvisioningStateFailed`, `TroubleshooterProvisioningStateRunning`, `TroubleshooterProvisioningStateSucceeded`
- New enum type `Type` with values `TypeAutomatedCheck`, `TypeDecision`, `TypeInsight`, `TypeSolution`
- New function `NewCheckNameAvailabilityClient(azcore.TokenCredential, *arm.ClientOptions) (*CheckNameAvailabilityClient, error)`
- New function `*CheckNameAvailabilityClient.Post(context.Context, string, *CheckNameAvailabilityClientPostOptions) (CheckNameAvailabilityClientPostResponse, error)`
- New function `*ClientFactory.NewCheckNameAvailabilityClient() *CheckNameAvailabilityClient`
- New function `*ClientFactory.NewSolutionClient() *SolutionClient`
- New function `*ClientFactory.NewTroubleshootersClient() *TroubleshootersClient`
- New function `NewSolutionClient(azcore.TokenCredential, *arm.ClientOptions) (*SolutionClient, error)`
- New function `*SolutionClient.BeginCreate(context.Context, string, string, SolutionResource, *SolutionClientBeginCreateOptions) (*runtime.Poller[SolutionClientCreateResponse], error)`
- New function `*SolutionClient.Get(context.Context, string, string, *SolutionClientGetOptions) (SolutionClientGetResponse, error)`
- New function `*SolutionClient.BeginUpdate(context.Context, string, string, SolutionPatchRequestBody, *SolutionClientBeginUpdateOptions) (*runtime.Poller[SolutionClientUpdateResponse], error)`
- New function `NewTroubleshootersClient(azcore.TokenCredential, *arm.ClientOptions) (*TroubleshootersClient, error)`
- New function `*TroubleshootersClient.Continue(context.Context, string, string, *TroubleshootersClientContinueOptions) (TroubleshootersClientContinueResponse, error)`
- New function `*TroubleshootersClient.Create(context.Context, string, string, TroubleshooterResource, *TroubleshootersClientCreateOptions) (TroubleshootersClientCreateResponse, error)`
- New function `*TroubleshootersClient.End(context.Context, string, string, *TroubleshootersClientEndOptions) (TroubleshootersClientEndResponse, error)`
- New function `*TroubleshootersClient.Get(context.Context, string, string, *TroubleshootersClientGetOptions) (TroubleshootersClientGetResponse, error)`
- New function `*TroubleshootersClient.Restart(context.Context, string, string, *TroubleshootersClientRestartOptions) (TroubleshootersClientRestartResponse, error)`
- New struct `AutomatedCheckResult`
- New struct `ContinueRequestBody`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `Filter`
- New struct `FilterGroup`
- New struct `MetricsBasedChart`
- New struct `ReplacementMaps`
- New struct `ResponseOption`
- New struct `ResponseValidationProperties`
- New struct `RestartTroubleshooterResponse`
- New struct `SearchResult`
- New struct `Section`
- New struct `SolutionPatchRequestBody`
- New struct `SolutionResource`
- New struct `SolutionResourceProperties`
- New struct `Solutions`
- New struct `SolutionsDiagnostic`
- New struct `SolutionsTroubleshooters`
- New struct `Step`
- New struct `StepInput`
- New struct `TriggerCriterion`
- New struct `TroubleshooterInstanceProperties`
- New struct `TroubleshooterResource`
- New struct `TroubleshooterResponse`
- New struct `Video`
- New struct `VideoGroup`
- New struct `VideoGroupVideo`
- New struct `WebResult`
- New field `RequiredInputs` in struct `SolutionMetadataProperties`


## 1.0.0 (2023-06-23)
### Other Changes

- Release stable version.

## 0.1.0 (2023-04-28)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/selfhelp/armselfhelp` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).