# Release History

## 0.2.0 (2022-02-22)
### Breaking Changes

- Function `*PolicyStatesClient.SummarizeForManagementGroup` parameter(s) have been changed from `(context.Context, Enum6, Enum0, string, *QueryOptions)` to `(context.Context, PolicyStatesSummaryResourceType, string, *QueryOptions)`
- Function `*PolicyStatesClient.SummarizeForSubscription` parameter(s) have been changed from `(context.Context, Enum6, string, *QueryOptions)` to `(context.Context, PolicyStatesSummaryResourceType, string, *QueryOptions)`
- Function `*RemediationsClient.CreateOrUpdateAtManagementGroup` parameter(s) have been changed from `(context.Context, Enum0, string, string, Remediation, *RemediationsClientCreateOrUpdateAtManagementGroupOptions)` to `(context.Context, string, string, Remediation, *RemediationsClientCreateOrUpdateAtManagementGroupOptions)`
- Function `*RemediationsClient.GetAtManagementGroup` parameter(s) have been changed from `(context.Context, Enum0, string, string, *RemediationsClientGetAtManagementGroupOptions)` to `(context.Context, string, string, *RemediationsClientGetAtManagementGroupOptions)`
- Function `*PolicyEventsClient.ListQueryResultsForPolicyDefinition` parameter(s) have been changed from `(Enum1, string, Enum4, string, *QueryOptions)` to `(PolicyEventsResourceType, string, string, *QueryOptions)`
- Function `*PolicyStatesClient.ListQueryResultsForSubscriptionLevelPolicyAssignment` parameter(s) have been changed from `(PolicyStatesResource, string, Enum4, string, *QueryOptions)` to `(PolicyStatesResource, string, string, *QueryOptions)`
- Function `*PolicyEventsClient.ListQueryResultsForSubscription` parameter(s) have been changed from `(Enum1, string, *QueryOptions)` to `(PolicyEventsResourceType, string, *QueryOptions)`
- Function `*PolicyStatesClient.ListQueryResultsForPolicyDefinition` parameter(s) have been changed from `(PolicyStatesResource, string, Enum4, string, *QueryOptions)` to `(PolicyStatesResource, string, string, *QueryOptions)`
- Function `*RemediationsClient.CancelAtManagementGroup` parameter(s) have been changed from `(context.Context, Enum0, string, string, *RemediationsClientCancelAtManagementGroupOptions)` to `(context.Context, string, string, *RemediationsClientCancelAtManagementGroupOptions)`
- Function `*RemediationsClient.ListDeploymentsAtManagementGroup` parameter(s) have been changed from `(Enum0, string, string, *QueryOptions)` to `(string, string, *QueryOptions)`
- Function `*PolicyEventsClient.ListQueryResultsForResourceGroup` parameter(s) have been changed from `(Enum1, string, string, *QueryOptions)` to `(PolicyEventsResourceType, string, string, *QueryOptions)`
- Function `*PolicyStatesClient.SummarizeForResourceGroup` parameter(s) have been changed from `(context.Context, Enum6, string, string, *QueryOptions)` to `(context.Context, PolicyStatesSummaryResourceType, string, string, *QueryOptions)`
- Function `*PolicyStatesClient.SummarizeForSubscriptionLevelPolicyAssignment` parameter(s) have been changed from `(context.Context, Enum6, string, Enum4, string, *QueryOptions)` to `(context.Context, PolicyStatesSummaryResourceType, string, string, *QueryOptions)`
- Function `*PolicyStatesClient.SummarizeForResource` parameter(s) have been changed from `(context.Context, Enum6, string, *QueryOptions)` to `(context.Context, PolicyStatesSummaryResourceType, string, *QueryOptions)`
- Function `*PolicyStatesClient.ListQueryResultsForManagementGroup` parameter(s) have been changed from `(PolicyStatesResource, Enum0, string, *QueryOptions)` to `(PolicyStatesResource, string, *QueryOptions)`
- Function `*PolicyEventsClient.ListQueryResultsForResourceGroupLevelPolicyAssignment` parameter(s) have been changed from `(Enum1, string, string, Enum4, string, *QueryOptions)` to `(PolicyEventsResourceType, string, string, string, *QueryOptions)`
- Function `*PolicyStatesClient.ListQueryResultsForPolicySetDefinition` parameter(s) have been changed from `(PolicyStatesResource, string, Enum4, string, *QueryOptions)` to `(PolicyStatesResource, string, string, *QueryOptions)`
- Function `*PolicyTrackedResourcesClient.ListQueryResultsForManagementGroup` parameter(s) have been changed from `(Enum0, string, Enum1, *QueryOptions)` to `(string, PolicyTrackedResourcesResourceType, *QueryOptions)`
- Function `*PolicyTrackedResourcesClient.ListQueryResultsForResourceGroup` parameter(s) have been changed from `(string, Enum1, *QueryOptions)` to `(string, PolicyTrackedResourcesResourceType, *QueryOptions)`
- Function `*RemediationsClient.DeleteAtManagementGroup` parameter(s) have been changed from `(context.Context, Enum0, string, string, *RemediationsClientDeleteAtManagementGroupOptions)` to `(context.Context, string, string, *RemediationsClientDeleteAtManagementGroupOptions)`
- Function `*PolicyEventsClient.ListQueryResultsForPolicySetDefinition` parameter(s) have been changed from `(Enum1, string, Enum4, string, *QueryOptions)` to `(PolicyEventsResourceType, string, string, *QueryOptions)`
- Function `*PolicyEventsClient.ListQueryResultsForManagementGroup` parameter(s) have been changed from `(Enum1, Enum0, string, *QueryOptions)` to `(PolicyEventsResourceType, string, *QueryOptions)`
- Function `*PolicyStatesClient.SummarizeForPolicySetDefinition` parameter(s) have been changed from `(context.Context, Enum6, string, Enum4, string, *QueryOptions)` to `(context.Context, PolicyStatesSummaryResourceType, string, string, *QueryOptions)`
- Function `*PolicyTrackedResourcesClient.ListQueryResultsForResource` parameter(s) have been changed from `(string, Enum1, *QueryOptions)` to `(string, PolicyTrackedResourcesResourceType, *QueryOptions)`
- Function `*PolicyEventsClient.ListQueryResultsForSubscriptionLevelPolicyAssignment` parameter(s) have been changed from `(Enum1, string, Enum4, string, *QueryOptions)` to `(PolicyEventsResourceType, string, string, *QueryOptions)`
- Function `*RemediationsClient.ListForManagementGroup` parameter(s) have been changed from `(Enum0, string, *QueryOptions)` to `(string, *QueryOptions)`
- Function `*PolicyStatesClient.SummarizeForPolicyDefinition` parameter(s) have been changed from `(context.Context, Enum6, string, Enum4, string, *QueryOptions)` to `(context.Context, PolicyStatesSummaryResourceType, string, string, *QueryOptions)`
- Function `*PolicyTrackedResourcesClient.ListQueryResultsForSubscription` parameter(s) have been changed from `(Enum1, *QueryOptions)` to `(PolicyTrackedResourcesResourceType, *QueryOptions)`
- Function `*PolicyEventsClient.ListQueryResultsForResource` parameter(s) have been changed from `(Enum1, string, *QueryOptions)` to `(PolicyEventsResourceType, string, *QueryOptions)`
- Function `*PolicyStatesClient.SummarizeForResourceGroupLevelPolicyAssignment` parameter(s) have been changed from `(context.Context, Enum6, string, string, Enum4, string, *QueryOptions)` to `(context.Context, PolicyStatesSummaryResourceType, string, string, string, *QueryOptions)`
- Function `*PolicyStatesClient.ListQueryResultsForResourceGroupLevelPolicyAssignment` parameter(s) have been changed from `(PolicyStatesResource, string, string, Enum4, string, *QueryOptions)` to `(PolicyStatesResource, string, string, string, *QueryOptions)`
- Type of `ExpressionEvaluationDetails.ExpressionValue` has been changed from `map[string]interface{}` to `interface{}`
- Type of `ExpressionEvaluationDetails.TargetValue` has been changed from `map[string]interface{}` to `interface{}`
- Type of `CheckRestrictionsResourceDetails.ResourceContent` has been changed from `map[string]interface{}` to `interface{}`
- Type of `PolicyMetadataProperties.Metadata` has been changed from `map[string]interface{}` to `interface{}`
- Type of `PolicyMetadataSlimProperties.Metadata` has been changed from `map[string]interface{}` to `interface{}`
- Const `Enum6Latest` has been removed
- Const `Enum1Default` has been removed
- Const `Enum0MicrosoftManagement` has been removed
- Const `Enum4MicrosoftAuthorization` has been removed
- Function `Enum6.ToPtr` has been removed
- Function `PossibleEnum4Values` has been removed
- Function `ErrorDefinitionAutoGenerated.MarshalJSON` has been removed
- Function `Enum1.ToPtr` has been removed
- Function `PossibleEnum6Values` has been removed
- Function `Enum4.ToPtr` has been removed
- Function `PossibleEnum0Values` has been removed
- Function `PossibleEnum1Values` has been removed
- Function `ErrorDefinitionAutoGenerated2.MarshalJSON` has been removed
- Function `Enum0.ToPtr` has been removed
- Struct `ErrorDefinitionAutoGenerated` has been removed
- Struct `ErrorDefinitionAutoGenerated2` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ErrorResponseAutoGenerated` has been removed
- Struct `ErrorResponseAutoGenerated2` has been removed
- Struct `QueryFailure` has been removed
- Struct `QueryFailureError` has been removed

### Features Added

- New const `PolicyStatesSummaryResourceTypeLatest`
- New const `PolicyEventsResourceTypeDefault`
- New const `PolicyTrackedResourcesResourceTypeDefault`
- New function `PolicyTrackedResourcesResourceType.ToPtr() *PolicyTrackedResourcesResourceType`
- New function `PossiblePolicyStatesSummaryResourceTypeValues() []PolicyStatesSummaryResourceType`
- New function `PolicyStatesSummaryResourceType.ToPtr() *PolicyStatesSummaryResourceType`
- New function `PossiblePolicyTrackedResourcesResourceTypeValues() []PolicyTrackedResourcesResourceType`
- New function `PossiblePolicyEventsResourceTypeValues() []PolicyEventsResourceType`
- New function `PolicyEventsResourceType.ToPtr() *PolicyEventsResourceType`


## 0.1.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.1.0 (2022-01-14)

- Init release.