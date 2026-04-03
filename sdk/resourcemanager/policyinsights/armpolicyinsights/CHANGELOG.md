# Release History

## 0.10.0 (2026-03-26)
### Breaking Changes

- Function `*AttestationsClient.NewListForResourceGroupPager` parameter(s) have been changed from `(resourceGroupName string, queryOptions *QueryOptions, options *AttestationsClientListForResourceGroupOptions)` to `(resourceGroupName string, options *AttestationsClientListForResourceGroupOptions)`
- Function `*AttestationsClient.NewListForResourcePager` parameter(s) have been changed from `(resourceID string, queryOptions *QueryOptions, options *AttestationsClientListForResourceOptions)` to `(resourceID string, options *AttestationsClientListForResourceOptions)`
- Function `*AttestationsClient.NewListForSubscriptionPager` parameter(s) have been changed from `(queryOptions *QueryOptions, options *AttestationsClientListForSubscriptionOptions)` to `(options *AttestationsClientListForSubscriptionOptions)`
- Function `NewComponentPolicyStatesClient` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Function `*ComponentPolicyStatesClient.ListQueryResultsForPolicyDefinition` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, policyDefinitionName string, componentPolicyStatesResource ComponentPolicyStatesResource, options *ComponentPolicyStatesClientListQueryResultsForPolicyDefinitionOptions)` to `(ctx context.Context, policyDefinitionName string, componentPolicyStatesResource ComponentPolicyStatesResource, options *ComponentPolicyStatesClientListQueryResultsForPolicyDefinitionOptions)`
- Function `*ComponentPolicyStatesClient.ListQueryResultsForResourceGroup` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, resourceGroupName string, componentPolicyStatesResource ComponentPolicyStatesResource, options *ComponentPolicyStatesClientListQueryResultsForResourceGroupOptions)` to `(ctx context.Context, resourceGroupName string, componentPolicyStatesResource ComponentPolicyStatesResource, options *ComponentPolicyStatesClientListQueryResultsForResourceGroupOptions)`
- Function `*ComponentPolicyStatesClient.ListQueryResultsForResourceGroupLevelPolicyAssignment` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, resourceGroupName string, policyAssignmentName string, componentPolicyStatesResource ComponentPolicyStatesResource, options *ComponentPolicyStatesClientListQueryResultsForResourceGroupLevelPolicyAssignmentOptions)` to `(ctx context.Context, resourceGroupName string, policyAssignmentName string, componentPolicyStatesResource ComponentPolicyStatesResource, options *ComponentPolicyStatesClientListQueryResultsForResourceGroupLevelPolicyAssignmentOptions)`
- Function `*ComponentPolicyStatesClient.ListQueryResultsForSubscription` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, componentPolicyStatesResource ComponentPolicyStatesResource, options *ComponentPolicyStatesClientListQueryResultsForSubscriptionOptions)` to `(ctx context.Context, componentPolicyStatesResource ComponentPolicyStatesResource, options *ComponentPolicyStatesClientListQueryResultsForSubscriptionOptions)`
- Function `*ComponentPolicyStatesClient.ListQueryResultsForSubscriptionLevelPolicyAssignment` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, policyAssignmentName string, componentPolicyStatesResource ComponentPolicyStatesResource, options *ComponentPolicyStatesClientListQueryResultsForSubscriptionLevelPolicyAssignmentOptions)` to `(ctx context.Context, policyAssignmentName string, componentPolicyStatesResource ComponentPolicyStatesResource, options *ComponentPolicyStatesClientListQueryResultsForSubscriptionLevelPolicyAssignmentOptions)`
- Function `NewPolicyEventsClient` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Function `*PolicyEventsClient.NewListQueryResultsForManagementGroupPager` parameter(s) have been changed from `(policyEventsResource PolicyEventsResourceType, managementGroupName string, queryOptions *QueryOptions, options *PolicyEventsClientListQueryResultsForManagementGroupOptions)` to `(policyEventsResource PolicyEventsResourceType, managementGroupName string, options *PolicyEventsClientListQueryResultsForManagementGroupOptions)`
- Function `*PolicyEventsClient.NewListQueryResultsForPolicyDefinitionPager` parameter(s) have been changed from `(policyEventsResource PolicyEventsResourceType, subscriptionID string, policyDefinitionName string, queryOptions *QueryOptions, options *PolicyEventsClientListQueryResultsForPolicyDefinitionOptions)` to `(policyEventsResource PolicyEventsResourceType, policyDefinitionName string, options *PolicyEventsClientListQueryResultsForPolicyDefinitionOptions)`
- Function `*PolicyEventsClient.NewListQueryResultsForPolicySetDefinitionPager` parameter(s) have been changed from `(policyEventsResource PolicyEventsResourceType, subscriptionID string, policySetDefinitionName string, queryOptions *QueryOptions, options *PolicyEventsClientListQueryResultsForPolicySetDefinitionOptions)` to `(policyEventsResource PolicyEventsResourceType, policySetDefinitionName string, options *PolicyEventsClientListQueryResultsForPolicySetDefinitionOptions)`
- Function `*PolicyEventsClient.NewListQueryResultsForResourceGroupLevelPolicyAssignmentPager` parameter(s) have been changed from `(policyEventsResource PolicyEventsResourceType, subscriptionID string, resourceGroupName string, policyAssignmentName string, queryOptions *QueryOptions, options *PolicyEventsClientListQueryResultsForResourceGroupLevelPolicyAssignmentOptions)` to `(resourceGroupName string, policyEventsResource PolicyEventsResourceType, policyAssignmentName string, options *PolicyEventsClientListQueryResultsForResourceGroupLevelPolicyAssignmentOptions)`
- Function `*PolicyEventsClient.NewListQueryResultsForResourceGroupPager` parameter(s) have been changed from `(policyEventsResource PolicyEventsResourceType, subscriptionID string, resourceGroupName string, queryOptions *QueryOptions, options *PolicyEventsClientListQueryResultsForResourceGroupOptions)` to `(resourceGroupName string, policyEventsResource PolicyEventsResourceType, options *PolicyEventsClientListQueryResultsForResourceGroupOptions)`
- Function `*PolicyEventsClient.NewListQueryResultsForResourcePager` parameter(s) have been changed from `(policyEventsResource PolicyEventsResourceType, resourceID string, queryOptions *QueryOptions, options *PolicyEventsClientListQueryResultsForResourceOptions)` to `(policyEventsResource PolicyEventsResourceType, resourceID string, options *PolicyEventsClientListQueryResultsForResourceOptions)`
- Function `*PolicyEventsClient.NewListQueryResultsForSubscriptionLevelPolicyAssignmentPager` parameter(s) have been changed from `(policyEventsResource PolicyEventsResourceType, subscriptionID string, policyAssignmentName string, queryOptions *QueryOptions, options *PolicyEventsClientListQueryResultsForSubscriptionLevelPolicyAssignmentOptions)` to `(policyEventsResource PolicyEventsResourceType, policyAssignmentName string, options *PolicyEventsClientListQueryResultsForSubscriptionLevelPolicyAssignmentOptions)`
- Function `*PolicyEventsClient.NewListQueryResultsForSubscriptionPager` parameter(s) have been changed from `(policyEventsResource PolicyEventsResourceType, subscriptionID string, queryOptions *QueryOptions, options *PolicyEventsClientListQueryResultsForSubscriptionOptions)` to `(policyEventsResource PolicyEventsResourceType, options *PolicyEventsClientListQueryResultsForSubscriptionOptions)`
- Function `*PolicyMetadataClient.NewListPager` parameter(s) have been changed from `(queryOptions *QueryOptions, options *PolicyMetadataClientListOptions)` to `(options *PolicyMetadataClientListOptions)`
- Function `NewPolicyStatesClient` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Function `*PolicyStatesClient.BeginTriggerResourceGroupEvaluation` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, resourceGroupName string, options *PolicyStatesClientBeginTriggerResourceGroupEvaluationOptions)` to `(ctx context.Context, resourceGroupName string, options *PolicyStatesClientBeginTriggerResourceGroupEvaluationOptions)`
- Function `*PolicyStatesClient.BeginTriggerSubscriptionEvaluation` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, options *PolicyStatesClientBeginTriggerSubscriptionEvaluationOptions)` to `(ctx context.Context, options *PolicyStatesClientBeginTriggerSubscriptionEvaluationOptions)`
- Function `*PolicyStatesClient.NewListQueryResultsForManagementGroupPager` parameter(s) have been changed from `(policyStatesResource PolicyStatesResource, managementGroupName string, queryOptions *QueryOptions, options *PolicyStatesClientListQueryResultsForManagementGroupOptions)` to `(policyStatesResource PolicyStatesResource, managementGroupName string, options *PolicyStatesClientListQueryResultsForManagementGroupOptions)`
- Function `*PolicyStatesClient.NewListQueryResultsForPolicyDefinitionPager` parameter(s) have been changed from `(policyStatesResource PolicyStatesResource, subscriptionID string, policyDefinitionName string, queryOptions *QueryOptions, options *PolicyStatesClientListQueryResultsForPolicyDefinitionOptions)` to `(policyStatesResource PolicyStatesResource, policyDefinitionName string, options *PolicyStatesClientListQueryResultsForPolicyDefinitionOptions)`
- Function `*PolicyStatesClient.NewListQueryResultsForPolicySetDefinitionPager` parameter(s) have been changed from `(policyStatesResource PolicyStatesResource, subscriptionID string, policySetDefinitionName string, queryOptions *QueryOptions, options *PolicyStatesClientListQueryResultsForPolicySetDefinitionOptions)` to `(policyStatesResource PolicyStatesResource, policySetDefinitionName string, options *PolicyStatesClientListQueryResultsForPolicySetDefinitionOptions)`
- Function `*PolicyStatesClient.NewListQueryResultsForResourceGroupLevelPolicyAssignmentPager` parameter(s) have been changed from `(policyStatesResource PolicyStatesResource, subscriptionID string, resourceGroupName string, policyAssignmentName string, queryOptions *QueryOptions, options *PolicyStatesClientListQueryResultsForResourceGroupLevelPolicyAssignmentOptions)` to `(resourceGroupName string, policyStatesResource PolicyStatesResource, policyAssignmentName string, options *PolicyStatesClientListQueryResultsForResourceGroupLevelPolicyAssignmentOptions)`
- Function `*PolicyStatesClient.NewListQueryResultsForResourceGroupPager` parameter(s) have been changed from `(policyStatesResource PolicyStatesResource, subscriptionID string, resourceGroupName string, queryOptions *QueryOptions, options *PolicyStatesClientListQueryResultsForResourceGroupOptions)` to `(resourceGroupName string, policyStatesResource PolicyStatesResource, options *PolicyStatesClientListQueryResultsForResourceGroupOptions)`
- Function `*PolicyStatesClient.NewListQueryResultsForResourcePager` parameter(s) have been changed from `(policyStatesResource PolicyStatesResource, resourceID string, queryOptions *QueryOptions, options *PolicyStatesClientListQueryResultsForResourceOptions)` to `(policyStatesResource PolicyStatesResource, resourceID string, options *PolicyStatesClientListQueryResultsForResourceOptions)`
- Function `*PolicyStatesClient.NewListQueryResultsForSubscriptionLevelPolicyAssignmentPager` parameter(s) have been changed from `(policyStatesResource PolicyStatesResource, subscriptionID string, policyAssignmentName string, queryOptions *QueryOptions, options *PolicyStatesClientListQueryResultsForSubscriptionLevelPolicyAssignmentOptions)` to `(policyStatesResource PolicyStatesResource, policyAssignmentName string, options *PolicyStatesClientListQueryResultsForSubscriptionLevelPolicyAssignmentOptions)`
- Function `*PolicyStatesClient.NewListQueryResultsForSubscriptionPager` parameter(s) have been changed from `(policyStatesResource PolicyStatesResource, subscriptionID string, queryOptions *QueryOptions, options *PolicyStatesClientListQueryResultsForSubscriptionOptions)` to `(policyStatesResource PolicyStatesResource, options *PolicyStatesClientListQueryResultsForSubscriptionOptions)`
- Function `*PolicyStatesClient.SummarizeForManagementGroup` parameter(s) have been changed from `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, managementGroupName string, queryOptions *QueryOptions, options *PolicyStatesClientSummarizeForManagementGroupOptions)` to `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, managementGroupName string, options *PolicyStatesClientSummarizeForManagementGroupOptions)`
- Function `*PolicyStatesClient.SummarizeForPolicyDefinition` parameter(s) have been changed from `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, subscriptionID string, policyDefinitionName string, queryOptions *QueryOptions, options *PolicyStatesClientSummarizeForPolicyDefinitionOptions)` to `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, policyDefinitionName string, options *PolicyStatesClientSummarizeForPolicyDefinitionOptions)`
- Function `*PolicyStatesClient.SummarizeForPolicySetDefinition` parameter(s) have been changed from `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, subscriptionID string, policySetDefinitionName string, queryOptions *QueryOptions, options *PolicyStatesClientSummarizeForPolicySetDefinitionOptions)` to `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, policySetDefinitionName string, options *PolicyStatesClientSummarizeForPolicySetDefinitionOptions)`
- Function `*PolicyStatesClient.SummarizeForResource` parameter(s) have been changed from `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, resourceID string, queryOptions *QueryOptions, options *PolicyStatesClientSummarizeForResourceOptions)` to `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, resourceID string, options *PolicyStatesClientSummarizeForResourceOptions)`
- Function `*PolicyStatesClient.SummarizeForResourceGroup` parameter(s) have been changed from `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, subscriptionID string, resourceGroupName string, queryOptions *QueryOptions, options *PolicyStatesClientSummarizeForResourceGroupOptions)` to `(ctx context.Context, resourceGroupName string, policyStatesSummaryResource PolicyStatesSummaryResourceType, options *PolicyStatesClientSummarizeForResourceGroupOptions)`
- Function `*PolicyStatesClient.SummarizeForResourceGroupLevelPolicyAssignment` parameter(s) have been changed from `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, subscriptionID string, resourceGroupName string, policyAssignmentName string, queryOptions *QueryOptions, options *PolicyStatesClientSummarizeForResourceGroupLevelPolicyAssignmentOptions)` to `(ctx context.Context, resourceGroupName string, policyStatesSummaryResource PolicyStatesSummaryResourceType, policyAssignmentName string, options *PolicyStatesClientSummarizeForResourceGroupLevelPolicyAssignmentOptions)`
- Function `*PolicyStatesClient.SummarizeForSubscription` parameter(s) have been changed from `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, subscriptionID string, queryOptions *QueryOptions, options *PolicyStatesClientSummarizeForSubscriptionOptions)` to `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, options *PolicyStatesClientSummarizeForSubscriptionOptions)`
- Function `*PolicyStatesClient.SummarizeForSubscriptionLevelPolicyAssignment` parameter(s) have been changed from `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, subscriptionID string, policyAssignmentName string, queryOptions *QueryOptions, options *PolicyStatesClientSummarizeForSubscriptionLevelPolicyAssignmentOptions)` to `(ctx context.Context, policyStatesSummaryResource PolicyStatesSummaryResourceType, policyAssignmentName string, options *PolicyStatesClientSummarizeForSubscriptionLevelPolicyAssignmentOptions)`
- Function `*PolicyTrackedResourcesClient.NewListQueryResultsForManagementGroupPager` parameter(s) have been changed from `(managementGroupName string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForManagementGroupOptions)` to `(managementGroupName string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, options *PolicyTrackedResourcesClientListQueryResultsForManagementGroupOptions)`
- Function `*PolicyTrackedResourcesClient.NewListQueryResultsForResourceGroupPager` parameter(s) have been changed from `(resourceGroupName string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForResourceGroupOptions)` to `(resourceGroupName string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, options *PolicyTrackedResourcesClientListQueryResultsForResourceGroupOptions)`
- Function `*PolicyTrackedResourcesClient.NewListQueryResultsForResourcePager` parameter(s) have been changed from `(resourceID string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForResourceOptions)` to `(resourceID string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, options *PolicyTrackedResourcesClientListQueryResultsForResourceOptions)`
- Function `*PolicyTrackedResourcesClient.NewListQueryResultsForSubscriptionPager` parameter(s) have been changed from `(policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForSubscriptionOptions)` to `(policyTrackedResourcesResource PolicyTrackedResourcesResourceType, options *PolicyTrackedResourcesClientListQueryResultsForSubscriptionOptions)`
- Function `*RemediationsClient.NewListDeploymentsAtManagementGroupPager` parameter(s) have been changed from `(managementGroupID string, remediationName string, queryOptions *QueryOptions, options *RemediationsClientListDeploymentsAtManagementGroupOptions)` to `(managementGroupID string, remediationName string, options *RemediationsClientListDeploymentsAtManagementGroupOptions)`
- Function `*RemediationsClient.NewListDeploymentsAtResourceGroupPager` parameter(s) have been changed from `(resourceGroupName string, remediationName string, queryOptions *QueryOptions, options *RemediationsClientListDeploymentsAtResourceGroupOptions)` to `(resourceGroupName string, remediationName string, options *RemediationsClientListDeploymentsAtResourceGroupOptions)`
- Function `*RemediationsClient.NewListDeploymentsAtResourcePager` parameter(s) have been changed from `(resourceID string, remediationName string, queryOptions *QueryOptions, options *RemediationsClientListDeploymentsAtResourceOptions)` to `(resourceID string, remediationName string, options *RemediationsClientListDeploymentsAtResourceOptions)`
- Function `*RemediationsClient.NewListDeploymentsAtSubscriptionPager` parameter(s) have been changed from `(remediationName string, queryOptions *QueryOptions, options *RemediationsClientListDeploymentsAtSubscriptionOptions)` to `(remediationName string, options *RemediationsClientListDeploymentsAtSubscriptionOptions)`
- Function `*RemediationsClient.NewListForManagementGroupPager` parameter(s) have been changed from `(managementGroupID string, queryOptions *QueryOptions, options *RemediationsClientListForManagementGroupOptions)` to `(managementGroupID string, options *RemediationsClientListForManagementGroupOptions)`
- Function `*RemediationsClient.NewListForResourceGroupPager` parameter(s) have been changed from `(resourceGroupName string, queryOptions *QueryOptions, options *RemediationsClientListForResourceGroupOptions)` to `(resourceGroupName string, options *RemediationsClientListForResourceGroupOptions)`
- Function `*RemediationsClient.NewListForResourcePager` parameter(s) have been changed from `(resourceID string, queryOptions *QueryOptions, options *RemediationsClientListForResourceOptions)` to `(resourceID string, options *RemediationsClientListForResourceOptions)`
- Function `*RemediationsClient.NewListForSubscriptionPager` parameter(s) have been changed from `(queryOptions *QueryOptions, options *RemediationsClientListForSubscriptionOptions)` to `(options *RemediationsClientListForSubscriptionOptions)`
- Struct `ErrorDefinitionAutoGenerated` has been removed
- Struct `ErrorDefinitionAutoGenerated2` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ErrorResponseAutoGenerated` has been removed
- Struct `ErrorResponseAutoGenerated2` has been removed
- Struct `QueryFailure` has been removed
- Struct `QueryFailureError` has been removed
- Struct `QueryOptions` has been removed
- Struct `Resource` has been removed

### Features Added

- New field `Filter`, `Top` in struct `AttestationsClientListForResourceGroupOptions`
- New field `Filter`, `Top` in struct `AttestationsClientListForResourceOptions`
- New field `Filter`, `Top` in struct `AttestationsClientListForSubscriptionOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyEventsClientListQueryResultsForManagementGroupOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyEventsClientListQueryResultsForPolicyDefinitionOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyEventsClientListQueryResultsForPolicySetDefinitionOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyEventsClientListQueryResultsForResourceGroupLevelPolicyAssignmentOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyEventsClientListQueryResultsForResourceGroupOptions`
- New field `Apply`, `Expand`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyEventsClientListQueryResultsForResourceOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyEventsClientListQueryResultsForSubscriptionLevelPolicyAssignmentOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyEventsClientListQueryResultsForSubscriptionOptions`
- New field `SystemData` in struct `PolicyMetadata`
- New field `Top` in struct `PolicyMetadataClientListOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyStatesClientListQueryResultsForManagementGroupOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyStatesClientListQueryResultsForPolicyDefinitionOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyStatesClientListQueryResultsForPolicySetDefinitionOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyStatesClientListQueryResultsForResourceGroupLevelPolicyAssignmentOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyStatesClientListQueryResultsForResourceGroupOptions`
- New field `Apply`, `Expand`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyStatesClientListQueryResultsForResourceOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyStatesClientListQueryResultsForSubscriptionLevelPolicyAssignmentOptions`
- New field `Apply`, `Filter`, `From`, `OrderBy`, `Select`, `SkipToken`, `To`, `Top` in struct `PolicyStatesClientListQueryResultsForSubscriptionOptions`
- New field `Filter`, `From`, `To`, `Top` in struct `PolicyStatesClientSummarizeForManagementGroupOptions`
- New field `Filter`, `From`, `To`, `Top` in struct `PolicyStatesClientSummarizeForPolicyDefinitionOptions`
- New field `Filter`, `From`, `To`, `Top` in struct `PolicyStatesClientSummarizeForPolicySetDefinitionOptions`
- New field `Filter`, `From`, `To`, `Top` in struct `PolicyStatesClientSummarizeForResourceGroupLevelPolicyAssignmentOptions`
- New field `Filter`, `From`, `To`, `Top` in struct `PolicyStatesClientSummarizeForResourceGroupOptions`
- New field `Filter`, `From`, `To`, `Top` in struct `PolicyStatesClientSummarizeForResourceOptions`
- New field `Filter`, `From`, `To`, `Top` in struct `PolicyStatesClientSummarizeForSubscriptionLevelPolicyAssignmentOptions`
- New field `Filter`, `From`, `To`, `Top` in struct `PolicyStatesClientSummarizeForSubscriptionOptions`
- New field `Filter`, `Top` in struct `PolicyTrackedResourcesClientListQueryResultsForManagementGroupOptions`
- New field `Filter`, `Top` in struct `PolicyTrackedResourcesClientListQueryResultsForResourceGroupOptions`
- New field `Filter`, `Top` in struct `PolicyTrackedResourcesClientListQueryResultsForResourceOptions`
- New field `Filter`, `Top` in struct `PolicyTrackedResourcesClientListQueryResultsForSubscriptionOptions`
- New field `RetryAfter` in struct `RemediationsClientCreateOrUpdateAtManagementGroupResponse`
- New field `RetryAfter` in struct `RemediationsClientCreateOrUpdateAtResourceResponse`
- New field `Top` in struct `RemediationsClientListDeploymentsAtManagementGroupOptions`
- New field `Top` in struct `RemediationsClientListDeploymentsAtResourceGroupOptions`
- New field `Top` in struct `RemediationsClientListDeploymentsAtResourceOptions`
- New field `Top` in struct `RemediationsClientListDeploymentsAtSubscriptionOptions`
- New field `Filter`, `Top` in struct `RemediationsClientListForManagementGroupOptions`
- New field `Filter`, `Top` in struct `RemediationsClientListForResourceGroupOptions`
- New field `Filter`, `Top` in struct `RemediationsClientListForResourceOptions`
- New field `Filter`, `Top` in struct `RemediationsClientListForSubscriptionOptions`


## 0.9.0 (2025-07-25)
### Breaking Changes

- Type of `PolicyEvaluationResult.EvaluationDetails` has been changed from `*PolicyEvaluationDetails` to `*CheckRestrictionEvaluationDetails`

### Features Added

- New value `FieldRestrictionResultAudit` added to enum type `FieldRestrictionResult`
- New enum type `ComponentPolicyStatesResource` with values `ComponentPolicyStatesResourceLatest`
- New function `*ClientFactory.NewComponentPolicyStatesClient() *ComponentPolicyStatesClient`
- New function `NewComponentPolicyStatesClient(azcore.TokenCredential, *arm.ClientOptions) (*ComponentPolicyStatesClient, error)`
- New function `*ComponentPolicyStatesClient.ListQueryResultsForPolicyDefinition(context.Context, string, string, ComponentPolicyStatesResource, *ComponentPolicyStatesClientListQueryResultsForPolicyDefinitionOptions) (ComponentPolicyStatesClientListQueryResultsForPolicyDefinitionResponse, error)`
- New function `*ComponentPolicyStatesClient.ListQueryResultsForResource(context.Context, string, ComponentPolicyStatesResource, *ComponentPolicyStatesClientListQueryResultsForResourceOptions) (ComponentPolicyStatesClientListQueryResultsForResourceResponse, error)`
- New function `*ComponentPolicyStatesClient.ListQueryResultsForResourceGroup(context.Context, string, string, ComponentPolicyStatesResource, *ComponentPolicyStatesClientListQueryResultsForResourceGroupOptions) (ComponentPolicyStatesClientListQueryResultsForResourceGroupResponse, error)`
- New function `*ComponentPolicyStatesClient.ListQueryResultsForResourceGroupLevelPolicyAssignment(context.Context, string, string, string, ComponentPolicyStatesResource, *ComponentPolicyStatesClientListQueryResultsForResourceGroupLevelPolicyAssignmentOptions) (ComponentPolicyStatesClientListQueryResultsForResourceGroupLevelPolicyAssignmentResponse, error)`
- New function `*ComponentPolicyStatesClient.ListQueryResultsForSubscription(context.Context, string, ComponentPolicyStatesResource, *ComponentPolicyStatesClientListQueryResultsForSubscriptionOptions) (ComponentPolicyStatesClientListQueryResultsForSubscriptionResponse, error)`
- New function `*ComponentPolicyStatesClient.ListQueryResultsForSubscriptionLevelPolicyAssignment(context.Context, string, string, ComponentPolicyStatesResource, *ComponentPolicyStatesClientListQueryResultsForSubscriptionLevelPolicyAssignmentOptions) (ComponentPolicyStatesClientListQueryResultsForSubscriptionLevelPolicyAssignmentResponse, error)`
- New struct `CheckRestrictionEvaluationDetails`
- New struct `ComponentExpressionEvaluationDetails`
- New struct `ComponentPolicyEvaluationDetails`
- New struct `ComponentPolicyState`
- New struct `ComponentPolicyStatesQueryResults`
- New struct `PolicyEffectDetails`
- New field `IncludeAuditEffect` in struct `CheckRestrictionsRequest`
- New field `PolicyEffect`, `Reason` in struct `FieldRestriction`
- New field `IsDataAction` in struct `Operation`
- New field `EffectDetails` in struct `PolicyEvaluationResult`
- New field `ResourceIDs` in struct `RemediationFilters`


## 0.8.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.7.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.7.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.6.0 (2022-10-07)
### Features Added

- New field `Metadata` in struct `AttestationProperties`
- New field `AssessmentDate` in struct `AttestationProperties`


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/policyinsights/armpolicyinsights` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).