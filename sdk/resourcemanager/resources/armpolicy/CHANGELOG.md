# Release History

## 2.0.0-beta.1 (2026-06-03)
### Breaking Changes

- Field `Ancestors`, `ComplianceReasonCode`, `PolicyAssignmentDisplayName`, `PolicyDefinitionDisplayName`, `PolicyDefinitionGroupNames`, `PolicyExemptionIDs`, `PolicySetDefinitionCategory`, `PolicySetDefinitionDisplayName`, `ResourceLocation` of struct `LogInfo` has been removed

### Features Added

- New value `SelectorKindGroupPrincipalID`, `SelectorKindUserPrincipalID` added to enum type `SelectorKind`
- New enum type `AliasPathAttributes` with values `AliasPathAttributesModifiable`, `AliasPathAttributesNone`
- New enum type `AliasPathTokenType` with values `AliasPathTokenTypeAny`, `AliasPathTokenTypeArray`, `AliasPathTokenTypeBoolean`, `AliasPathTokenTypeInteger`, `AliasPathTokenTypeNotSpecified`, `AliasPathTokenTypeNumber`, `AliasPathTokenTypeObject`, `AliasPathTokenTypeString`
- New enum type `AliasPatternType` with values `AliasPatternTypeExtract`, `AliasPatternTypeNotSpecified`
- New enum type `AliasType` with values `AliasTypeMask`, `AliasTypeNotSpecified`, `AliasTypePlainText`
- New enum type `AssignmentScopeValidation` with values `AssignmentScopeValidationDefault`, `AssignmentScopeValidationDoNotValidate`
- New enum type `ExemptionCategory` with values `ExemptionCategoryMitigated`, `ExemptionCategoryWaiver`
- New enum type `PolicyAction` with values `PolicyActionAllow`, `PolicyActionAudit`, `PolicyActionDeny`, `PolicyActionError`, `PolicyActionUnknown`
- New function `*ClientFactory.NewDataPolicyManifestsClient() *DataPolicyManifestsClient`
- New function `*ClientFactory.NewEnrollmentsClient() *EnrollmentsClient`
- New function `*ClientFactory.NewExemptionsClient() *ExemptionsClient`
- New function `*ClientFactory.NewVariableValuesClient() *VariableValuesClient`
- New function `*ClientFactory.NewVariablesClient() *VariablesClient`
- New function `NewDataPolicyManifestsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*DataPolicyManifestsClient, error)`
- New function `*DataPolicyManifestsClient.GetByPolicyMode(ctx context.Context, policyMode string, options *DataPolicyManifestsClientGetByPolicyModeOptions) (DataPolicyManifestsClientGetByPolicyModeResponse, error)`
- New function `*DataPolicyManifestsClient.NewListPager(options *DataPolicyManifestsClientListOptions) *runtime.Pager[DataPolicyManifestsClientListResponse]`
- New function `NewEnrollmentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*EnrollmentsClient, error)`
- New function `*EnrollmentsClient.CreateOrUpdate(ctx context.Context, scope string, policyEnrollmentName string, parameters Enrollment, options *EnrollmentsClientCreateOrUpdateOptions) (EnrollmentsClientCreateOrUpdateResponse, error)`
- New function `*EnrollmentsClient.Delete(ctx context.Context, scope string, policyEnrollmentName string, options *EnrollmentsClientDeleteOptions) (EnrollmentsClientDeleteResponse, error)`
- New function `*EnrollmentsClient.Get(ctx context.Context, scope string, policyEnrollmentName string, options *EnrollmentsClientGetOptions) (EnrollmentsClientGetResponse, error)`
- New function `*EnrollmentsClient.NewListForManagementGroupPager(managementGroupID string, options *EnrollmentsClientListForManagementGroupOptions) *runtime.Pager[EnrollmentsClientListForManagementGroupResponse]`
- New function `*EnrollmentsClient.NewListForResourceGroupPager(resourceGroupName string, options *EnrollmentsClientListForResourceGroupOptions) *runtime.Pager[EnrollmentsClientListForResourceGroupResponse]`
- New function `*EnrollmentsClient.NewListForResourcePager(resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, options *EnrollmentsClientListForResourceOptions) *runtime.Pager[EnrollmentsClientListForResourceResponse]`
- New function `*EnrollmentsClient.NewListPager(options *EnrollmentsClientListOptions) *runtime.Pager[EnrollmentsClientListResponse]`
- New function `*EnrollmentsClient.Update(ctx context.Context, scope string, policyEnrollmentName string, parameters EnrollmentUpdate, options *EnrollmentsClientUpdateOptions) (EnrollmentsClientUpdateResponse, error)`
- New function `NewExemptionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ExemptionsClient, error)`
- New function `*ExemptionsClient.CreateOrUpdate(ctx context.Context, scope string, policyExemptionName string, parameters Exemption, options *ExemptionsClientCreateOrUpdateOptions) (ExemptionsClientCreateOrUpdateResponse, error)`
- New function `*ExemptionsClient.Delete(ctx context.Context, scope string, policyExemptionName string, options *ExemptionsClientDeleteOptions) (ExemptionsClientDeleteResponse, error)`
- New function `*ExemptionsClient.Get(ctx context.Context, scope string, policyExemptionName string, options *ExemptionsClientGetOptions) (ExemptionsClientGetResponse, error)`
- New function `*ExemptionsClient.NewListForManagementGroupPager(managementGroupID string, options *ExemptionsClientListForManagementGroupOptions) *runtime.Pager[ExemptionsClientListForManagementGroupResponse]`
- New function `*ExemptionsClient.NewListForResourceGroupPager(resourceGroupName string, options *ExemptionsClientListForResourceGroupOptions) *runtime.Pager[ExemptionsClientListForResourceGroupResponse]`
- New function `*ExemptionsClient.NewListForResourcePager(resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, options *ExemptionsClientListForResourceOptions) *runtime.Pager[ExemptionsClientListForResourceResponse]`
- New function `*ExemptionsClient.NewListPager(options *ExemptionsClientListOptions) *runtime.Pager[ExemptionsClientListResponse]`
- New function `*ExemptionsClient.Update(ctx context.Context, scope string, policyExemptionName string, parameters ExemptionUpdate, options *ExemptionsClientUpdateOptions) (ExemptionsClientUpdateResponse, error)`
- New function `NewVariableValuesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VariableValuesClient, error)`
- New function `*VariableValuesClient.CreateOrUpdate(ctx context.Context, variableName string, variableValueName string, parameters VariableValue, options *VariableValuesClientCreateOrUpdateOptions) (VariableValuesClientCreateOrUpdateResponse, error)`
- New function `*VariableValuesClient.CreateOrUpdateAtManagementGroup(ctx context.Context, managementGroupID string, variableName string, variableValueName string, parameters VariableValue, options *VariableValuesClientCreateOrUpdateAtManagementGroupOptions) (VariableValuesClientCreateOrUpdateAtManagementGroupResponse, error)`
- New function `*VariableValuesClient.Delete(ctx context.Context, variableName string, variableValueName string, options *VariableValuesClientDeleteOptions) (VariableValuesClientDeleteResponse, error)`
- New function `*VariableValuesClient.DeleteAtManagementGroup(ctx context.Context, managementGroupID string, variableName string, variableValueName string, options *VariableValuesClientDeleteAtManagementGroupOptions) (VariableValuesClientDeleteAtManagementGroupResponse, error)`
- New function `*VariableValuesClient.Get(ctx context.Context, variableName string, variableValueName string, options *VariableValuesClientGetOptions) (VariableValuesClientGetResponse, error)`
- New function `*VariableValuesClient.GetAtManagementGroup(ctx context.Context, managementGroupID string, variableName string, variableValueName string, options *VariableValuesClientGetAtManagementGroupOptions) (VariableValuesClientGetAtManagementGroupResponse, error)`
- New function `*VariableValuesClient.NewListForManagementGroupPager(managementGroupID string, variableName string, options *VariableValuesClientListForManagementGroupOptions) *runtime.Pager[VariableValuesClientListForManagementGroupResponse]`
- New function `*VariableValuesClient.NewListPager(variableName string, options *VariableValuesClientListOptions) *runtime.Pager[VariableValuesClientListResponse]`
- New function `NewVariablesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VariablesClient, error)`
- New function `*VariablesClient.CreateOrUpdate(ctx context.Context, variableName string, parameters Variable, options *VariablesClientCreateOrUpdateOptions) (VariablesClientCreateOrUpdateResponse, error)`
- New function `*VariablesClient.CreateOrUpdateAtManagementGroup(ctx context.Context, managementGroupID string, variableName string, parameters Variable, options *VariablesClientCreateOrUpdateAtManagementGroupOptions) (VariablesClientCreateOrUpdateAtManagementGroupResponse, error)`
- New function `*VariablesClient.Delete(ctx context.Context, variableName string, options *VariablesClientDeleteOptions) (VariablesClientDeleteResponse, error)`
- New function `*VariablesClient.DeleteAtManagementGroup(ctx context.Context, managementGroupID string, variableName string, options *VariablesClientDeleteAtManagementGroupOptions) (VariablesClientDeleteAtManagementGroupResponse, error)`
- New function `*VariablesClient.Get(ctx context.Context, variableName string, options *VariablesClientGetOptions) (VariablesClientGetResponse, error)`
- New function `*VariablesClient.GetAtManagementGroup(ctx context.Context, managementGroupID string, variableName string, options *VariablesClientGetAtManagementGroupOptions) (VariablesClientGetAtManagementGroupResponse, error)`
- New function `*VariablesClient.NewListForManagementGroupPager(managementGroupID string, options *VariablesClientListForManagementGroupOptions) *runtime.Pager[VariablesClientListForManagementGroupResponse]`
- New function `*VariablesClient.NewListPager(options *VariablesClientListOptions) *runtime.Pager[VariablesClientListResponse]`
- New struct `Alias`
- New struct `AliasPath`
- New struct `AliasPathMetadata`
- New struct `AliasPattern`
- New struct `DataEffect`
- New struct `DataManifestCustomResourceFunctionDefinition`
- New struct `DataManifestResourceFunctionsDefinition`
- New struct `DataPolicyManifest`
- New struct `DataPolicyManifestListResult`
- New struct `DataPolicyManifestProperties`
- New struct `Enrollment`
- New struct `EnrollmentListResult`
- New struct `EnrollmentProperties`
- New struct `EnrollmentUpdate`
- New struct `EnrollmentUpdateProperties`
- New struct `Exemption`
- New struct `ExemptionListResult`
- New struct `ExemptionProperties`
- New struct `ExemptionUpdate`
- New struct `ExemptionUpdateProperties`
- New struct `ResourceTypeAliases`
- New struct `SelfServeExemptionSettings`
- New struct `TokenEvaluatedRequestDetails`
- New struct `Variable`
- New struct `VariableColumn`
- New struct `VariableListResult`
- New struct `VariableProperties`
- New struct `VariableValue`
- New struct `VariableValueColumnValue`
- New struct `VariableValueListResult`
- New struct `VariableValueProperties`
- New field `SelfServeExemptionSettings` in struct `AssignmentProperties`
- New field `SelfServeExemptionSettings` in struct `AssignmentUpdateProperties`
- New field `AdditionalInfo`, `EndpointKind`, `PolicyAction`, `PolicyEvaluationDetails` in struct `ExternalEvaluationEndpointInvocationResult`
- New field `RequestDetails` in struct `TokenResponse`


## 1.0.0 (2026-03-19)
### Breaking Changes

- Enum `AliasPathAttributes` has been removed
- Enum `AliasPathTokenType` has been removed
- Enum `AliasPatternType` has been removed
- Enum `AliasType` has been removed
- Enum `AssignmentScopeValidation` has been removed
- Enum `ExemptionCategory` has been removed
- Function `*ClientFactory.NewDataPolicyManifestsClient` has been removed
- Function `*ClientFactory.NewExemptionsClient` has been removed
- Function `*ClientFactory.NewVariableValuesClient` has been removed
- Function `*ClientFactory.NewVariablesClient` has been removed
- Function `NewDataPolicyManifestsClient` has been removed
- Function `*DataPolicyManifestsClient.GetByPolicyMode` has been removed
- Function `*DataPolicyManifestsClient.NewListPager` has been removed
- Function `NewExemptionsClient` has been removed
- Function `*ExemptionsClient.CreateOrUpdate` has been removed
- Function `*ExemptionsClient.Delete` has been removed
- Function `*ExemptionsClient.Get` has been removed
- Function `*ExemptionsClient.NewListForManagementGroupPager` has been removed
- Function `*ExemptionsClient.NewListForResourceGroupPager` has been removed
- Function `*ExemptionsClient.NewListForResourcePager` has been removed
- Function `*ExemptionsClient.NewListPager` has been removed
- Function `*ExemptionsClient.Update` has been removed
- Function `NewVariableValuesClient` has been removed
- Function `*VariableValuesClient.CreateOrUpdate` has been removed
- Function `*VariableValuesClient.CreateOrUpdateAtManagementGroup` has been removed
- Function `*VariableValuesClient.Delete` has been removed
- Function `*VariableValuesClient.DeleteAtManagementGroup` has been removed
- Function `*VariableValuesClient.Get` has been removed
- Function `*VariableValuesClient.GetAtManagementGroup` has been removed
- Function `*VariableValuesClient.NewListForManagementGroupPager` has been removed
- Function `*VariableValuesClient.NewListPager` has been removed
- Function `NewVariablesClient` has been removed
- Function `*VariablesClient.CreateOrUpdate` has been removed
- Function `*VariablesClient.CreateOrUpdateAtManagementGroup` has been removed
- Function `*VariablesClient.Delete` has been removed
- Function `*VariablesClient.DeleteAtManagementGroup` has been removed
- Function `*VariablesClient.Get` has been removed
- Function `*VariablesClient.GetAtManagementGroup` has been removed
- Function `*VariablesClient.NewListForManagementGroupPager` has been removed
- Function `*VariablesClient.NewListPager` has been removed
- Struct `Alias` has been removed
- Struct `AliasPath` has been removed
- Struct `AliasPathMetadata` has been removed
- Struct `AliasPattern` has been removed
- Struct `DataEffect` has been removed
- Struct `DataManifestCustomResourceFunctionDefinition` has been removed
- Struct `DataManifestResourceFunctionsDefinition` has been removed
- Struct `DataPolicyManifest` has been removed
- Struct `DataPolicyManifestListResult` has been removed
- Struct `DataPolicyManifestProperties` has been removed
- Struct `Exemption` has been removed
- Struct `ExemptionListResult` has been removed
- Struct `ExemptionProperties` has been removed
- Struct `ExemptionUpdate` has been removed
- Struct `ExemptionUpdateProperties` has been removed
- Struct `ResourceTypeAliases` has been removed
- Struct `Variable` has been removed
- Struct `VariableColumn` has been removed
- Struct `VariableListResult` has been removed
- Struct `VariableProperties` has been removed
- Struct `VariableValue` has been removed
- Struct `VariableValueColumnValue` has been removed
- Struct `VariableValueListResult` has been removed
- Struct `VariableValueProperties` has been removed
- Field `Expand` of struct `AssignmentsClientGetByIDOptions` has been removed

### Features Added

- New value `EnforcementModeEnroll` added to enum type `EnforcementMode`
- New value `OverrideKindDefinitionVersion` added to enum type `OverrideKind`
- New enum type `AssignmentType` with values `AssignmentTypeCustom`, `AssignmentTypeNotSpecified`, `AssignmentTypeSystem`, `AssignmentTypeSystemHidden`
- New enum type `ExternalEndpointResult` with values `ExternalEndpointResultFailed`, `ExternalEndpointResultSucceeded`
- New enum type `PolicyTokenResult` with values `PolicyTokenResultFailed`, `PolicyTokenResultSucceeded`
- New function `*ClientFactory.NewTokensClient() *TokensClient`
- New function `NewTokensClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TokensClient, error)`
- New function `*TokensClient.Acquire(ctx context.Context, parameters TokenRequest, options *TokensClientAcquireOptions) (TokensClientAcquireResponse, error)`
- New function `*TokensClient.AcquireAtManagementGroup(ctx context.Context, managementGroupName string, parameters TokenRequest, options *TokensClientAcquireAtManagementGroupOptions) (TokensClientAcquireAtManagementGroupResponse, error)`
- New struct `ExternalEvaluationEndpointInvocationResult`
- New struct `ExternalEvaluationEndpointSettings`
- New struct `ExternalEvaluationEnforcementSettings`
- New struct `LogInfo`
- New struct `TokenOperation`
- New struct `TokenRequest`
- New struct `TokenResponse`
- New field `AssignmentType`, `InstanceID` in struct `AssignmentProperties`
- New field `ExternalEvaluationEnforcementSettings` in struct `DefinitionProperties`
- New field `ExternalEvaluationEnforcementSettings` in struct `DefinitionVersionProperties`


## 0.10.0 (2025-03-18)
### Features Added

- New function `*ClientFactory.NewDefinitionVersionsClient() *DefinitionVersionsClient`
- New function `*ClientFactory.NewSetDefinitionVersionsClient() *SetDefinitionVersionsClient`
- New function `NewDefinitionVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DefinitionVersionsClient, error)`
- New function `*DefinitionVersionsClient.CreateOrUpdate(context.Context, string, string, DefinitionVersion, *DefinitionVersionsClientCreateOrUpdateOptions) (DefinitionVersionsClientCreateOrUpdateResponse, error)`
- New function `*DefinitionVersionsClient.CreateOrUpdateAtManagementGroup(context.Context, string, string, string, DefinitionVersion, *DefinitionVersionsClientCreateOrUpdateAtManagementGroupOptions) (DefinitionVersionsClientCreateOrUpdateAtManagementGroupResponse, error)`
- New function `*DefinitionVersionsClient.Delete(context.Context, string, string, *DefinitionVersionsClientDeleteOptions) (DefinitionVersionsClientDeleteResponse, error)`
- New function `*DefinitionVersionsClient.DeleteAtManagementGroup(context.Context, string, string, string, *DefinitionVersionsClientDeleteAtManagementGroupOptions) (DefinitionVersionsClientDeleteAtManagementGroupResponse, error)`
- New function `*DefinitionVersionsClient.Get(context.Context, string, string, *DefinitionVersionsClientGetOptions) (DefinitionVersionsClientGetResponse, error)`
- New function `*DefinitionVersionsClient.GetAtManagementGroup(context.Context, string, string, string, *DefinitionVersionsClientGetAtManagementGroupOptions) (DefinitionVersionsClientGetAtManagementGroupResponse, error)`
- New function `*DefinitionVersionsClient.GetBuiltIn(context.Context, string, string, *DefinitionVersionsClientGetBuiltInOptions) (DefinitionVersionsClientGetBuiltInResponse, error)`
- New function `*DefinitionVersionsClient.ListAll(context.Context, *DefinitionVersionsClientListAllOptions) (DefinitionVersionsClientListAllResponse, error)`
- New function `*DefinitionVersionsClient.ListAllAtManagementGroup(context.Context, string, *DefinitionVersionsClientListAllAtManagementGroupOptions) (DefinitionVersionsClientListAllAtManagementGroupResponse, error)`
- New function `*DefinitionVersionsClient.ListAllBuiltins(context.Context, *DefinitionVersionsClientListAllBuiltinsOptions) (DefinitionVersionsClientListAllBuiltinsResponse, error)`
- New function `*DefinitionVersionsClient.NewListBuiltInPager(string, *DefinitionVersionsClientListBuiltInOptions) *runtime.Pager[DefinitionVersionsClientListBuiltInResponse]`
- New function `*DefinitionVersionsClient.NewListByManagementGroupPager(string, string, *DefinitionVersionsClientListByManagementGroupOptions) *runtime.Pager[DefinitionVersionsClientListByManagementGroupResponse]`
- New function `*DefinitionVersionsClient.NewListPager(string, *DefinitionVersionsClientListOptions) *runtime.Pager[DefinitionVersionsClientListResponse]`
- New function `NewSetDefinitionVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SetDefinitionVersionsClient, error)`
- New function `*SetDefinitionVersionsClient.CreateOrUpdate(context.Context, string, string, SetDefinitionVersion, *SetDefinitionVersionsClientCreateOrUpdateOptions) (SetDefinitionVersionsClientCreateOrUpdateResponse, error)`
- New function `*SetDefinitionVersionsClient.CreateOrUpdateAtManagementGroup(context.Context, string, string, string, SetDefinitionVersion, *SetDefinitionVersionsClientCreateOrUpdateAtManagementGroupOptions) (SetDefinitionVersionsClientCreateOrUpdateAtManagementGroupResponse, error)`
- New function `*SetDefinitionVersionsClient.Delete(context.Context, string, string, *SetDefinitionVersionsClientDeleteOptions) (SetDefinitionVersionsClientDeleteResponse, error)`
- New function `*SetDefinitionVersionsClient.DeleteAtManagementGroup(context.Context, string, string, string, *SetDefinitionVersionsClientDeleteAtManagementGroupOptions) (SetDefinitionVersionsClientDeleteAtManagementGroupResponse, error)`
- New function `*SetDefinitionVersionsClient.Get(context.Context, string, string, *SetDefinitionVersionsClientGetOptions) (SetDefinitionVersionsClientGetResponse, error)`
- New function `*SetDefinitionVersionsClient.GetAtManagementGroup(context.Context, string, string, string, *SetDefinitionVersionsClientGetAtManagementGroupOptions) (SetDefinitionVersionsClientGetAtManagementGroupResponse, error)`
- New function `*SetDefinitionVersionsClient.GetBuiltIn(context.Context, string, string, *SetDefinitionVersionsClientGetBuiltInOptions) (SetDefinitionVersionsClientGetBuiltInResponse, error)`
- New function `*SetDefinitionVersionsClient.ListAll(context.Context, *SetDefinitionVersionsClientListAllOptions) (SetDefinitionVersionsClientListAllResponse, error)`
- New function `*SetDefinitionVersionsClient.ListAllAtManagementGroup(context.Context, string, *SetDefinitionVersionsClientListAllAtManagementGroupOptions) (SetDefinitionVersionsClientListAllAtManagementGroupResponse, error)`
- New function `*SetDefinitionVersionsClient.ListAllBuiltins(context.Context, *SetDefinitionVersionsClientListAllBuiltinsOptions) (SetDefinitionVersionsClientListAllBuiltinsResponse, error)`
- New function `*SetDefinitionVersionsClient.NewListBuiltInPager(string, *SetDefinitionVersionsClientListBuiltInOptions) *runtime.Pager[SetDefinitionVersionsClientListBuiltInResponse]`
- New function `*SetDefinitionVersionsClient.NewListByManagementGroupPager(string, string, *SetDefinitionVersionsClientListByManagementGroupOptions) *runtime.Pager[SetDefinitionVersionsClientListByManagementGroupResponse]`
- New function `*SetDefinitionVersionsClient.NewListPager(string, *SetDefinitionVersionsClientListOptions) *runtime.Pager[SetDefinitionVersionsClientListResponse]`
- New struct `DefinitionVersion`
- New struct `DefinitionVersionListResult`
- New struct `DefinitionVersionProperties`
- New struct `SetDefinitionVersion`
- New struct `SetDefinitionVersionListResult`
- New struct `SetDefinitionVersionProperties`
- New field `DefinitionVersion`, `EffectiveDefinitionVersion`, `LatestDefinitionVersion` in struct `AssignmentProperties`
- New field `Expand` in struct `AssignmentsClientGetByIDOptions`
- New field `Expand` in struct `AssignmentsClientGetOptions`
- New field `Expand` in struct `AssignmentsClientListForManagementGroupOptions`
- New field `Expand` in struct `AssignmentsClientListForResourceGroupOptions`
- New field `Expand` in struct `AssignmentsClientListForResourceOptions`
- New field `Expand` in struct `AssignmentsClientListOptions`
- New field `Version`, `Versions` in struct `DefinitionProperties`
- New field `DefinitionVersion`, `EffectiveDefinitionVersion`, `LatestDefinitionVersion` in struct `DefinitionReference`
- New field `Schema` in struct `ParameterDefinitionsValue`
- New field `Version`, `Versions` in struct `SetDefinitionProperties`
- New field `Expand` in struct `SetDefinitionsClientGetAtManagementGroupOptions`
- New field `Expand` in struct `SetDefinitionsClientGetBuiltInOptions`
- New field `Expand` in struct `SetDefinitionsClientGetOptions`
- New field `Expand` in struct `SetDefinitionsClientListBuiltInOptions`
- New field `Expand` in struct `SetDefinitionsClientListByManagementGroupOptions`
- New field `Expand` in struct `SetDefinitionsClientListOptions`


## 0.9.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.8.0 (2023-10-27)
### Features Added

- New enum type `AssignmentScopeValidation` with values `AssignmentScopeValidationDefault`, `AssignmentScopeValidationDoNotValidate`
- New enum type `OverrideKind` with values `OverrideKindPolicyEffect`
- New enum type `SelectorKind` with values `SelectorKindPolicyDefinitionReferenceID`, `SelectorKindResourceLocation`, `SelectorKindResourceType`, `SelectorKindResourceWithoutLocation`
- New function `*ClientFactory.NewVariableValuesClient() *VariableValuesClient`
- New function `*ClientFactory.NewVariablesClient() *VariablesClient`
- New function `*ExemptionsClient.Update(context.Context, string, string, ExemptionUpdate, *ExemptionsClientUpdateOptions) (ExemptionsClientUpdateResponse, error)`
- New function `NewVariableValuesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VariableValuesClient, error)`
- New function `*VariableValuesClient.CreateOrUpdate(context.Context, string, string, VariableValue, *VariableValuesClientCreateOrUpdateOptions) (VariableValuesClientCreateOrUpdateResponse, error)`
- New function `*VariableValuesClient.CreateOrUpdateAtManagementGroup(context.Context, string, string, string, VariableValue, *VariableValuesClientCreateOrUpdateAtManagementGroupOptions) (VariableValuesClientCreateOrUpdateAtManagementGroupResponse, error)`
- New function `*VariableValuesClient.Delete(context.Context, string, string, *VariableValuesClientDeleteOptions) (VariableValuesClientDeleteResponse, error)`
- New function `*VariableValuesClient.DeleteAtManagementGroup(context.Context, string, string, string, *VariableValuesClientDeleteAtManagementGroupOptions) (VariableValuesClientDeleteAtManagementGroupResponse, error)`
- New function `*VariableValuesClient.Get(context.Context, string, string, *VariableValuesClientGetOptions) (VariableValuesClientGetResponse, error)`
- New function `*VariableValuesClient.GetAtManagementGroup(context.Context, string, string, string, *VariableValuesClientGetAtManagementGroupOptions) (VariableValuesClientGetAtManagementGroupResponse, error)`
- New function `*VariableValuesClient.NewListForManagementGroupPager(string, string, *VariableValuesClientListForManagementGroupOptions) *runtime.Pager[VariableValuesClientListForManagementGroupResponse]`
- New function `*VariableValuesClient.NewListPager(string, *VariableValuesClientListOptions) *runtime.Pager[VariableValuesClientListResponse]`
- New function `NewVariablesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VariablesClient, error)`
- New function `*VariablesClient.CreateOrUpdate(context.Context, string, Variable, *VariablesClientCreateOrUpdateOptions) (VariablesClientCreateOrUpdateResponse, error)`
- New function `*VariablesClient.CreateOrUpdateAtManagementGroup(context.Context, string, string, Variable, *VariablesClientCreateOrUpdateAtManagementGroupOptions) (VariablesClientCreateOrUpdateAtManagementGroupResponse, error)`
- New function `*VariablesClient.Delete(context.Context, string, *VariablesClientDeleteOptions) (VariablesClientDeleteResponse, error)`
- New function `*VariablesClient.DeleteAtManagementGroup(context.Context, string, string, *VariablesClientDeleteAtManagementGroupOptions) (VariablesClientDeleteAtManagementGroupResponse, error)`
- New function `*VariablesClient.Get(context.Context, string, *VariablesClientGetOptions) (VariablesClientGetResponse, error)`
- New function `*VariablesClient.GetAtManagementGroup(context.Context, string, string, *VariablesClientGetAtManagementGroupOptions) (VariablesClientGetAtManagementGroupResponse, error)`
- New function `*VariablesClient.NewListForManagementGroupPager(string, *VariablesClientListForManagementGroupOptions) *runtime.Pager[VariablesClientListForManagementGroupResponse]`
- New function `*VariablesClient.NewListPager(*VariablesClientListOptions) *runtime.Pager[VariablesClientListResponse]`
- New struct `AssignmentUpdateProperties`
- New struct `ExemptionUpdate`
- New struct `ExemptionUpdateProperties`
- New struct `Override`
- New struct `ResourceSelector`
- New struct `Selector`
- New struct `Variable`
- New struct `VariableColumn`
- New struct `VariableListResult`
- New struct `VariableProperties`
- New struct `VariableValue`
- New struct `VariableValueColumnValue`
- New struct `VariableValueListResult`
- New struct `VariableValueProperties`
- New field `Overrides`, `ResourceSelectors` in struct `AssignmentProperties`
- New field `Properties` in struct `AssignmentUpdate`
- New field `AssignmentScopeValidation`, `ResourceSelectors` in struct `ExemptionProperties`


## 0.7.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.7.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.6.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
