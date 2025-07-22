# Release History

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