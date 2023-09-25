//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armoperationsmanagement

// ManagementAssociationsClientCreateOrUpdateResponse contains the response from method ManagementAssociationsClient.CreateOrUpdate.
type ManagementAssociationsClientCreateOrUpdateResponse struct {
	// The container for solution.
	ManagementAssociation
}

// ManagementAssociationsClientDeleteResponse contains the response from method ManagementAssociationsClient.Delete.
type ManagementAssociationsClientDeleteResponse struct {
	// placeholder for future response values
}

// ManagementAssociationsClientGetResponse contains the response from method ManagementAssociationsClient.Get.
type ManagementAssociationsClientGetResponse struct {
	// The container for solution.
	ManagementAssociation
}

// ManagementAssociationsClientListBySubscriptionResponse contains the response from method ManagementAssociationsClient.ListBySubscription.
type ManagementAssociationsClientListBySubscriptionResponse struct {
	// the list of ManagementAssociation response
	ManagementAssociationPropertiesList
}

// ManagementConfigurationsClientCreateOrUpdateResponse contains the response from method ManagementConfigurationsClient.CreateOrUpdate.
type ManagementConfigurationsClientCreateOrUpdateResponse struct {
	// The container for solution.
	ManagementConfiguration
}

// ManagementConfigurationsClientDeleteResponse contains the response from method ManagementConfigurationsClient.Delete.
type ManagementConfigurationsClientDeleteResponse struct {
	// placeholder for future response values
}

// ManagementConfigurationsClientGetResponse contains the response from method ManagementConfigurationsClient.Get.
type ManagementConfigurationsClientGetResponse struct {
	// The container for solution.
	ManagementConfiguration
}

// ManagementConfigurationsClientListBySubscriptionResponse contains the response from method ManagementConfigurationsClient.ListBySubscription.
type ManagementConfigurationsClientListBySubscriptionResponse struct {
	// the list of ManagementConfiguration response
	ManagementConfigurationPropertiesList
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// Result of the request to list solution operations.
	OperationListResult
}

// SolutionsClientCreateOrUpdateResponse contains the response from method SolutionsClient.BeginCreateOrUpdate.
type SolutionsClientCreateOrUpdateResponse struct {
	// The container for solution.
	Solution
}

// SolutionsClientDeleteResponse contains the response from method SolutionsClient.BeginDelete.
type SolutionsClientDeleteResponse struct {
	// placeholder for future response values
}

// SolutionsClientGetResponse contains the response from method SolutionsClient.Get.
type SolutionsClientGetResponse struct {
	// The container for solution.
	Solution
}

// SolutionsClientListByResourceGroupResponse contains the response from method SolutionsClient.ListByResourceGroup.
type SolutionsClientListByResourceGroupResponse struct {
	// the list of solution response
	SolutionPropertiesList
}

// SolutionsClientListBySubscriptionResponse contains the response from method SolutionsClient.ListBySubscription.
type SolutionsClientListBySubscriptionResponse struct {
	// the list of solution response
	SolutionPropertiesList
}

// SolutionsClientUpdateResponse contains the response from method SolutionsClient.BeginUpdate.
type SolutionsClientUpdateResponse struct {
	// The container for solution.
	Solution
}

