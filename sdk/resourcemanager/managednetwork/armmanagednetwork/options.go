//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmanagednetwork

// GroupsClientBeginCreateOrUpdateOptions contains the optional parameters for the GroupsClient.BeginCreateOrUpdate method.
type GroupsClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// GroupsClientBeginDeleteOptions contains the optional parameters for the GroupsClient.BeginDelete method.
type GroupsClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// GroupsClientGetOptions contains the optional parameters for the GroupsClient.Get method.
type GroupsClientGetOptions struct {
	// placeholder for future optional parameters
}

// GroupsClientListByManagedNetworkOptions contains the optional parameters for the GroupsClient.NewListByManagedNetworkPager
// method.
type GroupsClientListByManagedNetworkOptions struct {
	// Skiptoken is only used if a previous operation returned a partial result. If a previous response contains a nextLink element,
	// the value of the nextLink element will include a skiptoken parameter that
	// specifies a starting point to use for subsequent calls.
	Skiptoken *string

	// May be used to limit the number of results in a page for list queries.
	Top *int32
}

// ManagedNetworksClientBeginDeleteOptions contains the optional parameters for the ManagedNetworksClient.BeginDelete method.
type ManagedNetworksClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ManagedNetworksClientBeginUpdateOptions contains the optional parameters for the ManagedNetworksClient.BeginUpdate method.
type ManagedNetworksClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ManagedNetworksClientCreateOrUpdateOptions contains the optional parameters for the ManagedNetworksClient.CreateOrUpdate
// method.
type ManagedNetworksClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// ManagedNetworksClientGetOptions contains the optional parameters for the ManagedNetworksClient.Get method.
type ManagedNetworksClientGetOptions struct {
	// placeholder for future optional parameters
}

// ManagedNetworksClientListByResourceGroupOptions contains the optional parameters for the ManagedNetworksClient.NewListByResourceGroupPager
// method.
type ManagedNetworksClientListByResourceGroupOptions struct {
	// Skiptoken is only used if a previous operation returned a partial result. If a previous response contains a nextLink element,
	// the value of the nextLink element will include a skiptoken parameter that
	// specifies a starting point to use for subsequent calls.
	Skiptoken *string

	// May be used to limit the number of results in a page for list queries.
	Top *int32
}

// ManagedNetworksClientListBySubscriptionOptions contains the optional parameters for the ManagedNetworksClient.NewListBySubscriptionPager
// method.
type ManagedNetworksClientListBySubscriptionOptions struct {
	// Skiptoken is only used if a previous operation returned a partial result. If a previous response contains a nextLink element,
	// the value of the nextLink element will include a skiptoken parameter that
	// specifies a starting point to use for subsequent calls.
	Skiptoken *string

	// May be used to limit the number of results in a page for list queries.
	Top *int32
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.NewListPager method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// PeeringPoliciesClientBeginCreateOrUpdateOptions contains the optional parameters for the PeeringPoliciesClient.BeginCreateOrUpdate
// method.
type PeeringPoliciesClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// PeeringPoliciesClientBeginDeleteOptions contains the optional parameters for the PeeringPoliciesClient.BeginDelete method.
type PeeringPoliciesClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// PeeringPoliciesClientGetOptions contains the optional parameters for the PeeringPoliciesClient.Get method.
type PeeringPoliciesClientGetOptions struct {
	// placeholder for future optional parameters
}

// PeeringPoliciesClientListByManagedNetworkOptions contains the optional parameters for the PeeringPoliciesClient.NewListByManagedNetworkPager
// method.
type PeeringPoliciesClientListByManagedNetworkOptions struct {
	// Skiptoken is only used if a previous operation returned a partial result. If a previous response contains a nextLink element,
	// the value of the nextLink element will include a skiptoken parameter that
	// specifies a starting point to use for subsequent calls.
	Skiptoken *string

	// May be used to limit the number of results in a page for list queries.
	Top *int32
}

// ScopeAssignmentsClientCreateOrUpdateOptions contains the optional parameters for the ScopeAssignmentsClient.CreateOrUpdate
// method.
type ScopeAssignmentsClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// ScopeAssignmentsClientDeleteOptions contains the optional parameters for the ScopeAssignmentsClient.Delete method.
type ScopeAssignmentsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// ScopeAssignmentsClientGetOptions contains the optional parameters for the ScopeAssignmentsClient.Get method.
type ScopeAssignmentsClientGetOptions struct {
	// placeholder for future optional parameters
}

// ScopeAssignmentsClientListOptions contains the optional parameters for the ScopeAssignmentsClient.NewListPager method.
type ScopeAssignmentsClientListOptions struct {
	// placeholder for future optional parameters
}
