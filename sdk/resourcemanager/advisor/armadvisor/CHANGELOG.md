# Release History

## 2.0.0-beta.1 (2026-03-19)
### Breaking Changes

- Type of `RecommendationsClientGenerateResponse.RetryAfter` has been changed from `*string` to `*int32`
- Struct `Resource` has been removed

### Features Added

- New enum type `Aggregated` with values `AggregatedDay`, `AggregatedMonth`, `AggregatedWeek`
- New enum type `Control` with values `ControlBusinessContinuity`, `ControlDisasterRecovery`, `ControlHighAvailability`, `ControlMonitoringAndAlerting`, `ControlOther`, `ControlPersonalized`, `ControlPrioritizedRecommendations`, `ControlScalability`, `ControlServiceUpgradeAndRetirement`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `Duration` with values `Duration14`, `Duration21`, `Duration30`, `Duration60`, `Duration7`, `Duration90`
- New enum type `PredictionType` with values `PredictionTypePredictiveRightsizing`
- New enum type `Priority` with values `PriorityCritical`, `PriorityHigh`, `PriorityInformational`, `PriorityLow`, `PriorityMedium`
- New enum type `PriorityName` with values `PriorityNameHigh`, `PriorityNameLow`, `PriorityNameMedium`
- New enum type `Reason` with values `ReasonAlternativeSolution`, `ReasonExcessiveInvestment`, `ReasonIncompatible`, `ReasonRiskAccepted`, `ReasonTooComplex`, `ReasonUnclear`
- New enum type `ReasonForRejectionName` with values `ReasonForRejectionNameNotARisk`, `ReasonForRejectionNameRiskAccepted`
- New enum type `RecommendationStatusName` with values `RecommendationStatusNameApproved`, `RecommendationStatusNamePending`, `RecommendationStatusNameRejected`
- New enum type `ReviewStatus` with values `ReviewStatusCompleted`, `ReviewStatusInProgress`, `ReviewStatusNew`, `ReviewStatusTriaged`
- New enum type `State` with values `StateApproved`, `StateCompleted`, `StateDismissed`, `StateInProgress`, `StatePending`, `StatePostponed`, `StateRejected`
- New function `NewAssessmentTypesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AssessmentTypesClient, error)`
- New function `*AssessmentTypesClient.NewListPager(options *AssessmentTypesClientListOptions) *runtime.Pager[AssessmentTypesClientListResponse]`
- New function `NewAssessmentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AssessmentsClient, error)`
- New function `*AssessmentsClient.Delete(ctx context.Context, assessmentName string, options *AssessmentsClientDeleteOptions) (AssessmentsClientDeleteResponse, error)`
- New function `*AssessmentsClient.Get(ctx context.Context, assessmentName string, options *AssessmentsClientGetOptions) (AssessmentsClientGetResponse, error)`
- New function `*AssessmentsClient.NewListPager(options *AssessmentsClientListOptions) *runtime.Pager[AssessmentsClientListResponse]`
- New function `*AssessmentsClient.Put(ctx context.Context, assessmentName string, assessmentContract AssessmentResult, options *AssessmentsClientPutOptions) (AssessmentsClientPutResponse, error)`
- New function `NewClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*Client, error)`
- New function `*Client.NewAssessmentTypesClient() *AssessmentTypesClient`
- New function `*Client.NewAssessmentsClient() *AssessmentsClient`
- New function `*Client.NewConfigurationsClient() *ConfigurationsClient`
- New function `*Client.NewOperationsClient() *OperationsClient`
- New function `*Client.NewRecommendationMetadataClient() *RecommendationMetadataClient`
- New function `*Client.NewRecommendationsClient() *RecommendationsClient`
- New function `*Client.NewResiliencyReviewsClient() *ResiliencyReviewsClient`
- New function `*Client.NewScoresClient() *ScoresClient`
- New function `*Client.NewSuppressionsClient() *SuppressionsClient`
- New function `*Client.NewTriageRecommendationsClient() *TriageRecommendationsClient`
- New function `*Client.NewTriageResourcesClient() *TriageResourcesClient`
- New function `*Client.NewWorkloadsClient() *WorkloadsClient`
- New function `*Client.Predict(ctx context.Context, predictionRequest PredictionRequest, options *ClientPredictOptions) (ClientPredictResponse, error)`
- New function `*ClientFactory.NewAssessmentTypesClient() *AssessmentTypesClient`
- New function `*ClientFactory.NewAssessmentsClient() *AssessmentsClient`
- New function `*ClientFactory.NewClient() *Client`
- New function `*ClientFactory.NewResiliencyReviewsClient() *ResiliencyReviewsClient`
- New function `*ClientFactory.NewScoresClient() *ScoresClient`
- New function `*ClientFactory.NewTriageRecommendationsClient() *TriageRecommendationsClient`
- New function `*ClientFactory.NewTriageResourcesClient() *TriageResourcesClient`
- New function `*ClientFactory.NewWorkloadsClient() *WorkloadsClient`
- New function `PossiblePriorityValues() []Priority`
- New function `PossibleReasonValues() []Reason`
- New function `*RecommendationsClient.NewListByTenantPager(resourceURI string, options *RecommendationsClientListByTenantOptions) *runtime.Pager[RecommendationsClientListByTenantResponse]`
- New function `*RecommendationsClient.Patch(ctx context.Context, resourceURI string, recommendationID string, trackedProperties TrackedRecommendationPropertiesPayload, options *RecommendationsClientPatchOptions) (RecommendationsClientPatchResponse, error)`
- New function `NewResiliencyReviewsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ResiliencyReviewsClient, error)`
- New function `*ResiliencyReviewsClient.Get(ctx context.Context, reviewID string, options *ResiliencyReviewsClientGetOptions) (ResiliencyReviewsClientGetResponse, error)`
- New function `*ResiliencyReviewsClient.NewListPager(options *ResiliencyReviewsClientListOptions) *runtime.Pager[ResiliencyReviewsClientListResponse]`
- New function `NewScoresClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ScoresClient, error)`
- New function `*ScoresClient.Get(ctx context.Context, name string, options *ScoresClientGetOptions) (ScoresClientGetResponse, error)`
- New function `*ScoresClient.NewListPager(options *ScoresClientListOptions) *runtime.Pager[ScoresClientListResponse]`
- New function `NewTriageRecommendationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TriageRecommendationsClient, error)`
- New function `*TriageRecommendationsClient.ApproveTriageRecommendation(ctx context.Context, reviewID string, recommendationID string, options *TriageRecommendationsClientApproveTriageRecommendationOptions) (TriageRecommendationsClientApproveTriageRecommendationResponse, error)`
- New function `*TriageRecommendationsClient.Get(ctx context.Context, reviewID string, recommendationID string, options *TriageRecommendationsClientGetOptions) (TriageRecommendationsClientGetResponse, error)`
- New function `*TriageRecommendationsClient.NewListPager(reviewID string, options *TriageRecommendationsClientListOptions) *runtime.Pager[TriageRecommendationsClientListResponse]`
- New function `*TriageRecommendationsClient.RejectTriageRecommendation(ctx context.Context, reviewID string, recommendationID string, recommendationRejectBody RecommendationRejectBody, options *TriageRecommendationsClientRejectTriageRecommendationOptions) (TriageRecommendationsClientRejectTriageRecommendationResponse, error)`
- New function `*TriageRecommendationsClient.ResetTriageRecommendation(ctx context.Context, reviewID string, recommendationID string, options *TriageRecommendationsClientResetTriageRecommendationOptions) (TriageRecommendationsClientResetTriageRecommendationResponse, error)`
- New function `NewTriageResourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TriageResourcesClient, error)`
- New function `*TriageResourcesClient.Get(ctx context.Context, reviewID string, recommendationID string, recommendationResourceID string, options *TriageResourcesClientGetOptions) (TriageResourcesClientGetResponse, error)`
- New function `*TriageResourcesClient.NewListPager(reviewID string, recommendationID string, options *TriageResourcesClientListOptions) *runtime.Pager[TriageResourcesClientListResponse]`
- New function `NewWorkloadsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WorkloadsClient, error)`
- New function `*WorkloadsClient.NewListPager(options *WorkloadsClientListOptions) *runtime.Pager[WorkloadsClientListResponse]`
- New struct `AssessmentListResult`
- New struct `AssessmentResult`
- New struct `AssessmentResultProperties`
- New struct `AssessmentTypeListResult`
- New struct `AssessmentTypeResult`
- New struct `PredictionRequest`
- New struct `PredictionRequestProperties`
- New struct `PredictionResponse`
- New struct `PredictionResponseProperties`
- New struct `RecommendationPropertiesResourceWorkload`
- New struct `RecommendationPropertiesReview`
- New struct `RecommendationRejectBody`
- New struct `ResiliencyReview`
- New struct `ResiliencyReviewCollection`
- New struct `ResiliencyReviewProperties`
- New struct `ScoreEntity`
- New struct `ScoreEntityForAdvisor`
- New struct `ScoreEntityProperties`
- New struct `ScoreResponse`
- New struct `SystemData`
- New struct `TimeSeriesEntity`
- New struct `TrackedRecommendationProperties`
- New struct `TrackedRecommendationPropertiesPayload`
- New struct `TrackedRecommendationPropertiesPayloadProperties`
- New struct `TriageRecommendation`
- New struct `TriageRecommendationCollection`
- New struct `TriageRecommendationProperties`
- New struct `TriageResource`
- New struct `TriageResourceCollection`
- New struct `TriageResourceProperties`
- New struct `WorkloadListResult`
- New struct `WorkloadResult`
- New field `SystemData` in struct `ConfigData`
- New field `Duration` in struct `ConfigDataProperties`
- New field `SystemData` in struct `MetadataEntity`
- New field `Control`, `Notes`, `ResourceWorkload`, `Review`, `SourceSystem`, `Tracked`, `TrackedProperties` in struct `RecommendationProperties`
- New field `RetryAfter` in struct `RecommendationsClientGetGenerateStatusResponse`
- New field `SystemData` in struct `ResourceRecommendationBase`
- New field `SystemData` in struct `SuppressionContract`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/advisor/armadvisor` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).