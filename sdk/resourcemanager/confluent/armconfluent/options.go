//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armconfluent

// AccessClientCreateRoleBindingOptions contains the optional parameters for the AccessClient.CreateRoleBinding method.
type AccessClientCreateRoleBindingOptions struct {
	// placeholder for future optional parameters
}

// AccessClientDeleteRoleBindingOptions contains the optional parameters for the AccessClient.DeleteRoleBinding method.
type AccessClientDeleteRoleBindingOptions struct {
	// placeholder for future optional parameters
}

// AccessClientInviteUserOptions contains the optional parameters for the AccessClient.InviteUser method.
type AccessClientInviteUserOptions struct {
	// placeholder for future optional parameters
}

// AccessClientListClustersOptions contains the optional parameters for the AccessClient.ListClusters method.
type AccessClientListClustersOptions struct {
	// placeholder for future optional parameters
}

// AccessClientListEnvironmentsOptions contains the optional parameters for the AccessClient.ListEnvironments method.
type AccessClientListEnvironmentsOptions struct {
	// placeholder for future optional parameters
}

// AccessClientListInvitationsOptions contains the optional parameters for the AccessClient.ListInvitations method.
type AccessClientListInvitationsOptions struct {
	// placeholder for future optional parameters
}

// AccessClientListRoleBindingNameListOptions contains the optional parameters for the AccessClient.ListRoleBindingNameList
// method.
type AccessClientListRoleBindingNameListOptions struct {
	// placeholder for future optional parameters
}

// AccessClientListRoleBindingsOptions contains the optional parameters for the AccessClient.ListRoleBindings method.
type AccessClientListRoleBindingsOptions struct {
	// placeholder for future optional parameters
}

// AccessClientListServiceAccountsOptions contains the optional parameters for the AccessClient.ListServiceAccounts method.
type AccessClientListServiceAccountsOptions struct {
	// placeholder for future optional parameters
}

// AccessClientListUsersOptions contains the optional parameters for the AccessClient.ListUsers method.
type AccessClientListUsersOptions struct {
	// placeholder for future optional parameters
}

// MarketplaceAgreementsClientCreateOptions contains the optional parameters for the MarketplaceAgreementsClient.Create method.
type MarketplaceAgreementsClientCreateOptions struct {
	// Confluent Marketplace Agreement resource
	Body *AgreementResource
}

// MarketplaceAgreementsClientListOptions contains the optional parameters for the MarketplaceAgreementsClient.NewListPager
// method.
type MarketplaceAgreementsClientListOptions struct {
	// placeholder for future optional parameters
}

// OrganizationClientBeginCreateOptions contains the optional parameters for the OrganizationClient.BeginCreate method.
type OrganizationClientBeginCreateOptions struct {
	// Organization resource model
	Body *OrganizationResource

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// OrganizationClientBeginDeleteOptions contains the optional parameters for the OrganizationClient.BeginDelete method.
type OrganizationClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// OrganizationClientCreateAPIKeyOptions contains the optional parameters for the OrganizationClient.CreateAPIKey method.
type OrganizationClientCreateAPIKeyOptions struct {
	// placeholder for future optional parameters
}

// OrganizationClientDeleteClusterAPIKeyOptions contains the optional parameters for the OrganizationClient.DeleteClusterAPIKey
// method.
type OrganizationClientDeleteClusterAPIKeyOptions struct {
	// placeholder for future optional parameters
}

// OrganizationClientGetClusterAPIKeyOptions contains the optional parameters for the OrganizationClient.GetClusterAPIKey
// method.
type OrganizationClientGetClusterAPIKeyOptions struct {
	// placeholder for future optional parameters
}

// OrganizationClientGetClusterByIDOptions contains the optional parameters for the OrganizationClient.GetClusterByID method.
type OrganizationClientGetClusterByIDOptions struct {
	// placeholder for future optional parameters
}

// OrganizationClientGetEnvironmentByIDOptions contains the optional parameters for the OrganizationClient.GetEnvironmentByID
// method.
type OrganizationClientGetEnvironmentByIDOptions struct {
	// placeholder for future optional parameters
}

// OrganizationClientGetOptions contains the optional parameters for the OrganizationClient.Get method.
type OrganizationClientGetOptions struct {
	// placeholder for future optional parameters
}

// OrganizationClientGetSchemaRegistryClusterByIDOptions contains the optional parameters for the OrganizationClient.GetSchemaRegistryClusterByID
// method.
type OrganizationClientGetSchemaRegistryClusterByIDOptions struct {
	// placeholder for future optional parameters
}

// OrganizationClientListByResourceGroupOptions contains the optional parameters for the OrganizationClient.NewListByResourceGroupPager
// method.
type OrganizationClientListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// OrganizationClientListBySubscriptionOptions contains the optional parameters for the OrganizationClient.NewListBySubscriptionPager
// method.
type OrganizationClientListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// OrganizationClientListClustersOptions contains the optional parameters for the OrganizationClient.NewListClustersPager
// method.
type OrganizationClientListClustersOptions struct {
	// Pagination size
	PageSize *int32

	// An opaque pagination token to fetch the next set of records
	PageToken *string
}

// OrganizationClientListEnvironmentsOptions contains the optional parameters for the OrganizationClient.NewListEnvironmentsPager
// method.
type OrganizationClientListEnvironmentsOptions struct {
	// Pagination size
	PageSize *int32

	// An opaque pagination token to fetch the next set of records
	PageToken *string
}

// OrganizationClientListRegionsOptions contains the optional parameters for the OrganizationClient.ListRegions method.
type OrganizationClientListRegionsOptions struct {
	// placeholder for future optional parameters
}

// OrganizationClientListSchemaRegistryClustersOptions contains the optional parameters for the OrganizationClient.NewListSchemaRegistryClustersPager
// method.
type OrganizationClientListSchemaRegistryClustersOptions struct {
	// Pagination size
	PageSize *int32

	// An opaque pagination token to fetch the next set of records
	PageToken *string
}

// OrganizationClientUpdateOptions contains the optional parameters for the OrganizationClient.Update method.
type OrganizationClientUpdateOptions struct {
	// Updated Organization resource
	Body *OrganizationResourceUpdate
}

// OrganizationOperationsClientListOptions contains the optional parameters for the OrganizationOperationsClient.NewListPager
// method.
type OrganizationOperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// ValidationsClientValidateOrganizationOptions contains the optional parameters for the ValidationsClient.ValidateOrganization
// method.
type ValidationsClientValidateOrganizationOptions struct {
	// placeholder for future optional parameters
}

// ValidationsClientValidateOrganizationV2Options contains the optional parameters for the ValidationsClient.ValidateOrganizationV2
// method.
type ValidationsClientValidateOrganizationV2Options struct {
	// placeholder for future optional parameters
}
