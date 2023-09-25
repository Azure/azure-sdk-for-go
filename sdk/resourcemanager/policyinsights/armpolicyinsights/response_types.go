//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpolicyinsights

// AttestationsClientCreateOrUpdateAtResourceGroupResponse contains the response from method AttestationsClient.BeginCreateOrUpdateAtResourceGroup.
type AttestationsClientCreateOrUpdateAtResourceGroupResponse struct {
	// An attestation resource.
	Attestation
}

// AttestationsClientCreateOrUpdateAtResourceResponse contains the response from method AttestationsClient.BeginCreateOrUpdateAtResource.
type AttestationsClientCreateOrUpdateAtResourceResponse struct {
	// An attestation resource.
	Attestation
}

// AttestationsClientCreateOrUpdateAtSubscriptionResponse contains the response from method AttestationsClient.BeginCreateOrUpdateAtSubscription.
type AttestationsClientCreateOrUpdateAtSubscriptionResponse struct {
	// An attestation resource.
	Attestation
}

// AttestationsClientDeleteAtResourceGroupResponse contains the response from method AttestationsClient.DeleteAtResourceGroup.
type AttestationsClientDeleteAtResourceGroupResponse struct {
	// placeholder for future response values
}

// AttestationsClientDeleteAtResourceResponse contains the response from method AttestationsClient.DeleteAtResource.
type AttestationsClientDeleteAtResourceResponse struct {
	// placeholder for future response values
}

// AttestationsClientDeleteAtSubscriptionResponse contains the response from method AttestationsClient.DeleteAtSubscription.
type AttestationsClientDeleteAtSubscriptionResponse struct {
	// placeholder for future response values
}

// AttestationsClientGetAtResourceGroupResponse contains the response from method AttestationsClient.GetAtResourceGroup.
type AttestationsClientGetAtResourceGroupResponse struct {
	// An attestation resource.
	Attestation
}

// AttestationsClientGetAtResourceResponse contains the response from method AttestationsClient.GetAtResource.
type AttestationsClientGetAtResourceResponse struct {
	// An attestation resource.
	Attestation
}

// AttestationsClientGetAtSubscriptionResponse contains the response from method AttestationsClient.GetAtSubscription.
type AttestationsClientGetAtSubscriptionResponse struct {
	// An attestation resource.
	Attestation
}

// AttestationsClientListForResourceGroupResponse contains the response from method AttestationsClient.NewListForResourceGroupPager.
type AttestationsClientListForResourceGroupResponse struct {
	// List of attestations.
	AttestationListResult
}

// AttestationsClientListForResourceResponse contains the response from method AttestationsClient.NewListForResourcePager.
type AttestationsClientListForResourceResponse struct {
	// List of attestations.
	AttestationListResult
}

// AttestationsClientListForSubscriptionResponse contains the response from method AttestationsClient.NewListForSubscriptionPager.
type AttestationsClientListForSubscriptionResponse struct {
	// List of attestations.
	AttestationListResult
}

// OperationsClientListResponse contains the response from method OperationsClient.List.
type OperationsClientListResponse struct {
	// List of available operations.
	OperationsListResults
}

// PolicyEventsClientListQueryResultsForManagementGroupResponse contains the response from method PolicyEventsClient.NewListQueryResultsForManagementGroupPager.
type PolicyEventsClientListQueryResultsForManagementGroupResponse struct {
	// Query results.
	PolicyEventsQueryResults
}

// PolicyEventsClientListQueryResultsForPolicyDefinitionResponse contains the response from method PolicyEventsClient.NewListQueryResultsForPolicyDefinitionPager.
type PolicyEventsClientListQueryResultsForPolicyDefinitionResponse struct {
	// Query results.
	PolicyEventsQueryResults
}

// PolicyEventsClientListQueryResultsForPolicySetDefinitionResponse contains the response from method PolicyEventsClient.NewListQueryResultsForPolicySetDefinitionPager.
type PolicyEventsClientListQueryResultsForPolicySetDefinitionResponse struct {
	// Query results.
	PolicyEventsQueryResults
}

// PolicyEventsClientListQueryResultsForResourceGroupLevelPolicyAssignmentResponse contains the response from method PolicyEventsClient.NewListQueryResultsForResourceGroupLevelPolicyAssignmentPager.
type PolicyEventsClientListQueryResultsForResourceGroupLevelPolicyAssignmentResponse struct {
	// Query results.
	PolicyEventsQueryResults
}

// PolicyEventsClientListQueryResultsForResourceGroupResponse contains the response from method PolicyEventsClient.NewListQueryResultsForResourceGroupPager.
type PolicyEventsClientListQueryResultsForResourceGroupResponse struct {
	// Query results.
	PolicyEventsQueryResults
}

// PolicyEventsClientListQueryResultsForResourceResponse contains the response from method PolicyEventsClient.NewListQueryResultsForResourcePager.
type PolicyEventsClientListQueryResultsForResourceResponse struct {
	// Query results.
	PolicyEventsQueryResults
}

// PolicyEventsClientListQueryResultsForSubscriptionLevelPolicyAssignmentResponse contains the response from method PolicyEventsClient.NewListQueryResultsForSubscriptionLevelPolicyAssignmentPager.
type PolicyEventsClientListQueryResultsForSubscriptionLevelPolicyAssignmentResponse struct {
	// Query results.
	PolicyEventsQueryResults
}

// PolicyEventsClientListQueryResultsForSubscriptionResponse contains the response from method PolicyEventsClient.NewListQueryResultsForSubscriptionPager.
type PolicyEventsClientListQueryResultsForSubscriptionResponse struct {
	// Query results.
	PolicyEventsQueryResults
}

// PolicyMetadataClientGetResourceResponse contains the response from method PolicyMetadataClient.GetResource.
type PolicyMetadataClientGetResourceResponse struct {
	// Policy metadata resource definition.
	PolicyMetadata
}

// PolicyMetadataClientListResponse contains the response from method PolicyMetadataClient.NewListPager.
type PolicyMetadataClientListResponse struct {
	// Collection of policy metadata resources.
	PolicyMetadataCollection
}

// PolicyRestrictionsClientCheckAtManagementGroupScopeResponse contains the response from method PolicyRestrictionsClient.CheckAtManagementGroupScope.
type PolicyRestrictionsClientCheckAtManagementGroupScopeResponse struct {
	// The result of a check policy restrictions evaluation on a resource.
	CheckRestrictionsResult
}

// PolicyRestrictionsClientCheckAtResourceGroupScopeResponse contains the response from method PolicyRestrictionsClient.CheckAtResourceGroupScope.
type PolicyRestrictionsClientCheckAtResourceGroupScopeResponse struct {
	// The result of a check policy restrictions evaluation on a resource.
	CheckRestrictionsResult
}

// PolicyRestrictionsClientCheckAtSubscriptionScopeResponse contains the response from method PolicyRestrictionsClient.CheckAtSubscriptionScope.
type PolicyRestrictionsClientCheckAtSubscriptionScopeResponse struct {
	// The result of a check policy restrictions evaluation on a resource.
	CheckRestrictionsResult
}

// PolicyStatesClientListQueryResultsForManagementGroupResponse contains the response from method PolicyStatesClient.NewListQueryResultsForManagementGroupPager.
type PolicyStatesClientListQueryResultsForManagementGroupResponse struct {
	// Query results.
	PolicyStatesQueryResults
}

// PolicyStatesClientListQueryResultsForPolicyDefinitionResponse contains the response from method PolicyStatesClient.NewListQueryResultsForPolicyDefinitionPager.
type PolicyStatesClientListQueryResultsForPolicyDefinitionResponse struct {
	// Query results.
	PolicyStatesQueryResults
}

// PolicyStatesClientListQueryResultsForPolicySetDefinitionResponse contains the response from method PolicyStatesClient.NewListQueryResultsForPolicySetDefinitionPager.
type PolicyStatesClientListQueryResultsForPolicySetDefinitionResponse struct {
	// Query results.
	PolicyStatesQueryResults
}

// PolicyStatesClientListQueryResultsForResourceGroupLevelPolicyAssignmentResponse contains the response from method PolicyStatesClient.NewListQueryResultsForResourceGroupLevelPolicyAssignmentPager.
type PolicyStatesClientListQueryResultsForResourceGroupLevelPolicyAssignmentResponse struct {
	// Query results.
	PolicyStatesQueryResults
}

// PolicyStatesClientListQueryResultsForResourceGroupResponse contains the response from method PolicyStatesClient.NewListQueryResultsForResourceGroupPager.
type PolicyStatesClientListQueryResultsForResourceGroupResponse struct {
	// Query results.
	PolicyStatesQueryResults
}

// PolicyStatesClientListQueryResultsForResourceResponse contains the response from method PolicyStatesClient.NewListQueryResultsForResourcePager.
type PolicyStatesClientListQueryResultsForResourceResponse struct {
	// Query results.
	PolicyStatesQueryResults
}

// PolicyStatesClientListQueryResultsForSubscriptionLevelPolicyAssignmentResponse contains the response from method PolicyStatesClient.NewListQueryResultsForSubscriptionLevelPolicyAssignmentPager.
type PolicyStatesClientListQueryResultsForSubscriptionLevelPolicyAssignmentResponse struct {
	// Query results.
	PolicyStatesQueryResults
}

// PolicyStatesClientListQueryResultsForSubscriptionResponse contains the response from method PolicyStatesClient.NewListQueryResultsForSubscriptionPager.
type PolicyStatesClientListQueryResultsForSubscriptionResponse struct {
	// Query results.
	PolicyStatesQueryResults
}

// PolicyStatesClientSummarizeForManagementGroupResponse contains the response from method PolicyStatesClient.SummarizeForManagementGroup.
type PolicyStatesClientSummarizeForManagementGroupResponse struct {
	// Summarize action results.
	SummarizeResults
}

// PolicyStatesClientSummarizeForPolicyDefinitionResponse contains the response from method PolicyStatesClient.SummarizeForPolicyDefinition.
type PolicyStatesClientSummarizeForPolicyDefinitionResponse struct {
	// Summarize action results.
	SummarizeResults
}

// PolicyStatesClientSummarizeForPolicySetDefinitionResponse contains the response from method PolicyStatesClient.SummarizeForPolicySetDefinition.
type PolicyStatesClientSummarizeForPolicySetDefinitionResponse struct {
	// Summarize action results.
	SummarizeResults
}

// PolicyStatesClientSummarizeForResourceGroupLevelPolicyAssignmentResponse contains the response from method PolicyStatesClient.SummarizeForResourceGroupLevelPolicyAssignment.
type PolicyStatesClientSummarizeForResourceGroupLevelPolicyAssignmentResponse struct {
	// Summarize action results.
	SummarizeResults
}

// PolicyStatesClientSummarizeForResourceGroupResponse contains the response from method PolicyStatesClient.SummarizeForResourceGroup.
type PolicyStatesClientSummarizeForResourceGroupResponse struct {
	// Summarize action results.
	SummarizeResults
}

// PolicyStatesClientSummarizeForResourceResponse contains the response from method PolicyStatesClient.SummarizeForResource.
type PolicyStatesClientSummarizeForResourceResponse struct {
	// Summarize action results.
	SummarizeResults
}

// PolicyStatesClientSummarizeForSubscriptionLevelPolicyAssignmentResponse contains the response from method PolicyStatesClient.SummarizeForSubscriptionLevelPolicyAssignment.
type PolicyStatesClientSummarizeForSubscriptionLevelPolicyAssignmentResponse struct {
	// Summarize action results.
	SummarizeResults
}

// PolicyStatesClientSummarizeForSubscriptionResponse contains the response from method PolicyStatesClient.SummarizeForSubscription.
type PolicyStatesClientSummarizeForSubscriptionResponse struct {
	// Summarize action results.
	SummarizeResults
}

// PolicyStatesClientTriggerResourceGroupEvaluationResponse contains the response from method PolicyStatesClient.BeginTriggerResourceGroupEvaluation.
type PolicyStatesClientTriggerResourceGroupEvaluationResponse struct {
	// placeholder for future response values
}

// PolicyStatesClientTriggerSubscriptionEvaluationResponse contains the response from method PolicyStatesClient.BeginTriggerSubscriptionEvaluation.
type PolicyStatesClientTriggerSubscriptionEvaluationResponse struct {
	// placeholder for future response values
}

// PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse contains the response from method PolicyTrackedResourcesClient.NewListQueryResultsForManagementGroupPager.
type PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse struct {
	// Query results.
	PolicyTrackedResourcesQueryResults
}

// PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse contains the response from method PolicyTrackedResourcesClient.NewListQueryResultsForResourceGroupPager.
type PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse struct {
	// Query results.
	PolicyTrackedResourcesQueryResults
}

// PolicyTrackedResourcesClientListQueryResultsForResourceResponse contains the response from method PolicyTrackedResourcesClient.NewListQueryResultsForResourcePager.
type PolicyTrackedResourcesClientListQueryResultsForResourceResponse struct {
	// Query results.
	PolicyTrackedResourcesQueryResults
}

// PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse contains the response from method PolicyTrackedResourcesClient.NewListQueryResultsForSubscriptionPager.
type PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse struct {
	// Query results.
	PolicyTrackedResourcesQueryResults
}

// RemediationsClientCancelAtManagementGroupResponse contains the response from method RemediationsClient.CancelAtManagementGroup.
type RemediationsClientCancelAtManagementGroupResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientCancelAtResourceGroupResponse contains the response from method RemediationsClient.CancelAtResourceGroup.
type RemediationsClientCancelAtResourceGroupResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientCancelAtResourceResponse contains the response from method RemediationsClient.CancelAtResource.
type RemediationsClientCancelAtResourceResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientCancelAtSubscriptionResponse contains the response from method RemediationsClient.CancelAtSubscription.
type RemediationsClientCancelAtSubscriptionResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientCreateOrUpdateAtManagementGroupResponse contains the response from method RemediationsClient.CreateOrUpdateAtManagementGroup.
type RemediationsClientCreateOrUpdateAtManagementGroupResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientCreateOrUpdateAtResourceGroupResponse contains the response from method RemediationsClient.CreateOrUpdateAtResourceGroup.
type RemediationsClientCreateOrUpdateAtResourceGroupResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientCreateOrUpdateAtResourceResponse contains the response from method RemediationsClient.CreateOrUpdateAtResource.
type RemediationsClientCreateOrUpdateAtResourceResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientCreateOrUpdateAtSubscriptionResponse contains the response from method RemediationsClient.CreateOrUpdateAtSubscription.
type RemediationsClientCreateOrUpdateAtSubscriptionResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientDeleteAtManagementGroupResponse contains the response from method RemediationsClient.DeleteAtManagementGroup.
type RemediationsClientDeleteAtManagementGroupResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientDeleteAtResourceGroupResponse contains the response from method RemediationsClient.DeleteAtResourceGroup.
type RemediationsClientDeleteAtResourceGroupResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientDeleteAtResourceResponse contains the response from method RemediationsClient.DeleteAtResource.
type RemediationsClientDeleteAtResourceResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientDeleteAtSubscriptionResponse contains the response from method RemediationsClient.DeleteAtSubscription.
type RemediationsClientDeleteAtSubscriptionResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientGetAtManagementGroupResponse contains the response from method RemediationsClient.GetAtManagementGroup.
type RemediationsClientGetAtManagementGroupResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientGetAtResourceGroupResponse contains the response from method RemediationsClient.GetAtResourceGroup.
type RemediationsClientGetAtResourceGroupResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientGetAtResourceResponse contains the response from method RemediationsClient.GetAtResource.
type RemediationsClientGetAtResourceResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientGetAtSubscriptionResponse contains the response from method RemediationsClient.GetAtSubscription.
type RemediationsClientGetAtSubscriptionResponse struct {
	// The remediation definition.
	Remediation
}

// RemediationsClientListDeploymentsAtManagementGroupResponse contains the response from method RemediationsClient.NewListDeploymentsAtManagementGroupPager.
type RemediationsClientListDeploymentsAtManagementGroupResponse struct {
	// List of deployments for a remediation.
	RemediationDeploymentsListResult
}

// RemediationsClientListDeploymentsAtResourceGroupResponse contains the response from method RemediationsClient.NewListDeploymentsAtResourceGroupPager.
type RemediationsClientListDeploymentsAtResourceGroupResponse struct {
	// List of deployments for a remediation.
	RemediationDeploymentsListResult
}

// RemediationsClientListDeploymentsAtResourceResponse contains the response from method RemediationsClient.NewListDeploymentsAtResourcePager.
type RemediationsClientListDeploymentsAtResourceResponse struct {
	// List of deployments for a remediation.
	RemediationDeploymentsListResult
}

// RemediationsClientListDeploymentsAtSubscriptionResponse contains the response from method RemediationsClient.NewListDeploymentsAtSubscriptionPager.
type RemediationsClientListDeploymentsAtSubscriptionResponse struct {
	// List of deployments for a remediation.
	RemediationDeploymentsListResult
}

// RemediationsClientListForManagementGroupResponse contains the response from method RemediationsClient.NewListForManagementGroupPager.
type RemediationsClientListForManagementGroupResponse struct {
	// List of remediations.
	RemediationListResult
}

// RemediationsClientListForResourceGroupResponse contains the response from method RemediationsClient.NewListForResourceGroupPager.
type RemediationsClientListForResourceGroupResponse struct {
	// List of remediations.
	RemediationListResult
}

// RemediationsClientListForResourceResponse contains the response from method RemediationsClient.NewListForResourcePager.
type RemediationsClientListForResourceResponse struct {
	// List of remediations.
	RemediationListResult
}

// RemediationsClientListForSubscriptionResponse contains the response from method RemediationsClient.NewListForSubscriptionPager.
type RemediationsClientListForSubscriptionResponse struct {
	// List of remediations.
	RemediationListResult
}

