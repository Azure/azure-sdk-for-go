//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpolicy

// AssignmentsClientCreateByIDResponse contains the response from method AssignmentsClient.CreateByID.
type AssignmentsClientCreateByIDResponse struct {
	// The policy assignment.
	Assignment
}

// AssignmentsClientCreateResponse contains the response from method AssignmentsClient.Create.
type AssignmentsClientCreateResponse struct {
	// The policy assignment.
	Assignment
}

// AssignmentsClientDeleteByIDResponse contains the response from method AssignmentsClient.DeleteByID.
type AssignmentsClientDeleteByIDResponse struct {
	// The policy assignment.
	Assignment
}

// AssignmentsClientDeleteResponse contains the response from method AssignmentsClient.Delete.
type AssignmentsClientDeleteResponse struct {
	// The policy assignment.
	Assignment
}

// AssignmentsClientGetByIDResponse contains the response from method AssignmentsClient.GetByID.
type AssignmentsClientGetByIDResponse struct {
	// The policy assignment.
	Assignment
}

// AssignmentsClientGetResponse contains the response from method AssignmentsClient.Get.
type AssignmentsClientGetResponse struct {
	// The policy assignment.
	Assignment
}

// AssignmentsClientListForManagementGroupResponse contains the response from method AssignmentsClient.NewListForManagementGroupPager.
type AssignmentsClientListForManagementGroupResponse struct {
	// List of policy assignments.
	AssignmentListResult
}

// AssignmentsClientListForResourceGroupResponse contains the response from method AssignmentsClient.NewListForResourceGroupPager.
type AssignmentsClientListForResourceGroupResponse struct {
	// List of policy assignments.
	AssignmentListResult
}

// AssignmentsClientListForResourceResponse contains the response from method AssignmentsClient.NewListForResourcePager.
type AssignmentsClientListForResourceResponse struct {
	// List of policy assignments.
	AssignmentListResult
}

// AssignmentsClientListResponse contains the response from method AssignmentsClient.NewListPager.
type AssignmentsClientListResponse struct {
	// List of policy assignments.
	AssignmentListResult
}

// AssignmentsClientUpdateByIDResponse contains the response from method AssignmentsClient.UpdateByID.
type AssignmentsClientUpdateByIDResponse struct {
	// The policy assignment.
	Assignment
}

// AssignmentsClientUpdateResponse contains the response from method AssignmentsClient.Update.
type AssignmentsClientUpdateResponse struct {
	// The policy assignment.
	Assignment
}

// DataPolicyManifestsClientGetByPolicyModeResponse contains the response from method DataPolicyManifestsClient.GetByPolicyMode.
type DataPolicyManifestsClientGetByPolicyModeResponse struct {
	// The data policy manifest.
	DataPolicyManifest
}

// DataPolicyManifestsClientListResponse contains the response from method DataPolicyManifestsClient.NewListPager.
type DataPolicyManifestsClientListResponse struct {
	// List of data policy manifests.
	DataPolicyManifestListResult
}

// DefinitionsClientCreateOrUpdateAtManagementGroupResponse contains the response from method DefinitionsClient.CreateOrUpdateAtManagementGroup.
type DefinitionsClientCreateOrUpdateAtManagementGroupResponse struct {
	// The policy definition.
	Definition
}

// DefinitionsClientCreateOrUpdateResponse contains the response from method DefinitionsClient.CreateOrUpdate.
type DefinitionsClientCreateOrUpdateResponse struct {
	// The policy definition.
	Definition
}

// DefinitionsClientDeleteAtManagementGroupResponse contains the response from method DefinitionsClient.DeleteAtManagementGroup.
type DefinitionsClientDeleteAtManagementGroupResponse struct {
	// placeholder for future response values
}

// DefinitionsClientDeleteResponse contains the response from method DefinitionsClient.Delete.
type DefinitionsClientDeleteResponse struct {
	// placeholder for future response values
}

// DefinitionsClientGetAtManagementGroupResponse contains the response from method DefinitionsClient.GetAtManagementGroup.
type DefinitionsClientGetAtManagementGroupResponse struct {
	// The policy definition.
	Definition
}

// DefinitionsClientGetBuiltInResponse contains the response from method DefinitionsClient.GetBuiltIn.
type DefinitionsClientGetBuiltInResponse struct {
	// The policy definition.
	Definition
}

// DefinitionsClientGetResponse contains the response from method DefinitionsClient.Get.
type DefinitionsClientGetResponse struct {
	// The policy definition.
	Definition
}

// DefinitionsClientListBuiltInResponse contains the response from method DefinitionsClient.NewListBuiltInPager.
type DefinitionsClientListBuiltInResponse struct {
	// List of policy definitions.
	DefinitionListResult
}

// DefinitionsClientListByManagementGroupResponse contains the response from method DefinitionsClient.NewListByManagementGroupPager.
type DefinitionsClientListByManagementGroupResponse struct {
	// List of policy definitions.
	DefinitionListResult
}

// DefinitionsClientListResponse contains the response from method DefinitionsClient.NewListPager.
type DefinitionsClientListResponse struct {
	// List of policy definitions.
	DefinitionListResult
}

// ExemptionsClientCreateOrUpdateResponse contains the response from method ExemptionsClient.CreateOrUpdate.
type ExemptionsClientCreateOrUpdateResponse struct {
	// The policy exemption.
	Exemption
}

// ExemptionsClientDeleteResponse contains the response from method ExemptionsClient.Delete.
type ExemptionsClientDeleteResponse struct {
	// placeholder for future response values
}

// ExemptionsClientGetResponse contains the response from method ExemptionsClient.Get.
type ExemptionsClientGetResponse struct {
	// The policy exemption.
	Exemption
}

// ExemptionsClientListForManagementGroupResponse contains the response from method ExemptionsClient.NewListForManagementGroupPager.
type ExemptionsClientListForManagementGroupResponse struct {
	// List of policy exemptions.
	ExemptionListResult
}

// ExemptionsClientListForResourceGroupResponse contains the response from method ExemptionsClient.NewListForResourceGroupPager.
type ExemptionsClientListForResourceGroupResponse struct {
	// List of policy exemptions.
	ExemptionListResult
}

// ExemptionsClientListForResourceResponse contains the response from method ExemptionsClient.NewListForResourcePager.
type ExemptionsClientListForResourceResponse struct {
	// List of policy exemptions.
	ExemptionListResult
}

// ExemptionsClientListResponse contains the response from method ExemptionsClient.NewListPager.
type ExemptionsClientListResponse struct {
	// List of policy exemptions.
	ExemptionListResult
}

// ExemptionsClientUpdateResponse contains the response from method ExemptionsClient.Update.
type ExemptionsClientUpdateResponse struct {
	// The policy exemption.
	Exemption
}

// SetDefinitionsClientCreateOrUpdateAtManagementGroupResponse contains the response from method SetDefinitionsClient.CreateOrUpdateAtManagementGroup.
type SetDefinitionsClientCreateOrUpdateAtManagementGroupResponse struct {
	// The policy set definition.
	SetDefinition
}

// SetDefinitionsClientCreateOrUpdateResponse contains the response from method SetDefinitionsClient.CreateOrUpdate.
type SetDefinitionsClientCreateOrUpdateResponse struct {
	// The policy set definition.
	SetDefinition
}

// SetDefinitionsClientDeleteAtManagementGroupResponse contains the response from method SetDefinitionsClient.DeleteAtManagementGroup.
type SetDefinitionsClientDeleteAtManagementGroupResponse struct {
	// placeholder for future response values
}

// SetDefinitionsClientDeleteResponse contains the response from method SetDefinitionsClient.Delete.
type SetDefinitionsClientDeleteResponse struct {
	// placeholder for future response values
}

// SetDefinitionsClientGetAtManagementGroupResponse contains the response from method SetDefinitionsClient.GetAtManagementGroup.
type SetDefinitionsClientGetAtManagementGroupResponse struct {
	// The policy set definition.
	SetDefinition
}

// SetDefinitionsClientGetBuiltInResponse contains the response from method SetDefinitionsClient.GetBuiltIn.
type SetDefinitionsClientGetBuiltInResponse struct {
	// The policy set definition.
	SetDefinition
}

// SetDefinitionsClientGetResponse contains the response from method SetDefinitionsClient.Get.
type SetDefinitionsClientGetResponse struct {
	// The policy set definition.
	SetDefinition
}

// SetDefinitionsClientListBuiltInResponse contains the response from method SetDefinitionsClient.NewListBuiltInPager.
type SetDefinitionsClientListBuiltInResponse struct {
	// List of policy set definitions.
	SetDefinitionListResult
}

// SetDefinitionsClientListByManagementGroupResponse contains the response from method SetDefinitionsClient.NewListByManagementGroupPager.
type SetDefinitionsClientListByManagementGroupResponse struct {
	// List of policy set definitions.
	SetDefinitionListResult
}

// SetDefinitionsClientListResponse contains the response from method SetDefinitionsClient.NewListPager.
type SetDefinitionsClientListResponse struct {
	// List of policy set definitions.
	SetDefinitionListResult
}

// VariableValuesClientCreateOrUpdateAtManagementGroupResponse contains the response from method VariableValuesClient.CreateOrUpdateAtManagementGroup.
type VariableValuesClientCreateOrUpdateAtManagementGroupResponse struct {
	// The variable value.
	VariableValue
}

// VariableValuesClientCreateOrUpdateResponse contains the response from method VariableValuesClient.CreateOrUpdate.
type VariableValuesClientCreateOrUpdateResponse struct {
	// The variable value.
	VariableValue
}

// VariableValuesClientDeleteAtManagementGroupResponse contains the response from method VariableValuesClient.DeleteAtManagementGroup.
type VariableValuesClientDeleteAtManagementGroupResponse struct {
	// placeholder for future response values
}

// VariableValuesClientDeleteResponse contains the response from method VariableValuesClient.Delete.
type VariableValuesClientDeleteResponse struct {
	// placeholder for future response values
}

// VariableValuesClientGetAtManagementGroupResponse contains the response from method VariableValuesClient.GetAtManagementGroup.
type VariableValuesClientGetAtManagementGroupResponse struct {
	// The variable value.
	VariableValue
}

// VariableValuesClientGetResponse contains the response from method VariableValuesClient.Get.
type VariableValuesClientGetResponse struct {
	// The variable value.
	VariableValue
}

// VariableValuesClientListForManagementGroupResponse contains the response from method VariableValuesClient.NewListForManagementGroupPager.
type VariableValuesClientListForManagementGroupResponse struct {
	// List of variable values.
	VariableValueListResult
}

// VariableValuesClientListResponse contains the response from method VariableValuesClient.NewListPager.
type VariableValuesClientListResponse struct {
	// List of variable values.
	VariableValueListResult
}

// VariablesClientCreateOrUpdateAtManagementGroupResponse contains the response from method VariablesClient.CreateOrUpdateAtManagementGroup.
type VariablesClientCreateOrUpdateAtManagementGroupResponse struct {
	// The variable.
	Variable
}

// VariablesClientCreateOrUpdateResponse contains the response from method VariablesClient.CreateOrUpdate.
type VariablesClientCreateOrUpdateResponse struct {
	// The variable.
	Variable
}

// VariablesClientDeleteAtManagementGroupResponse contains the response from method VariablesClient.DeleteAtManagementGroup.
type VariablesClientDeleteAtManagementGroupResponse struct {
	// placeholder for future response values
}

// VariablesClientDeleteResponse contains the response from method VariablesClient.Delete.
type VariablesClientDeleteResponse struct {
	// placeholder for future response values
}

// VariablesClientGetAtManagementGroupResponse contains the response from method VariablesClient.GetAtManagementGroup.
type VariablesClientGetAtManagementGroupResponse struct {
	// The variable.
	Variable
}

// VariablesClientGetResponse contains the response from method VariablesClient.Get.
type VariablesClientGetResponse struct {
	// The variable.
	Variable
}

// VariablesClientListForManagementGroupResponse contains the response from method VariablesClient.NewListForManagementGroupPager.
type VariablesClientListForManagementGroupResponse struct {
	// List of variables.
	VariableListResult
}

// VariablesClientListResponse contains the response from method VariablesClient.NewListPager.
type VariablesClientListResponse struct {
	// List of variables.
	VariableListResult
}
