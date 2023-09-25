//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmachinelearningservices

// ComputeClientBeginCreateOrUpdateOptions contains the optional parameters for the ComputeClient.BeginCreateOrUpdate method.
type ComputeClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ComputeClientBeginDeleteOptions contains the optional parameters for the ComputeClient.BeginDelete method.
type ComputeClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ComputeClientBeginRestartOptions contains the optional parameters for the ComputeClient.BeginRestart method.
type ComputeClientBeginRestartOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ComputeClientBeginStartOptions contains the optional parameters for the ComputeClient.BeginStart method.
type ComputeClientBeginStartOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ComputeClientBeginStopOptions contains the optional parameters for the ComputeClient.BeginStop method.
type ComputeClientBeginStopOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ComputeClientBeginUpdateOptions contains the optional parameters for the ComputeClient.BeginUpdate method.
type ComputeClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ComputeClientGetOptions contains the optional parameters for the ComputeClient.Get method.
type ComputeClientGetOptions struct {
	// placeholder for future optional parameters
}

// ComputeClientListKeysOptions contains the optional parameters for the ComputeClient.ListKeys method.
type ComputeClientListKeysOptions struct {
	// placeholder for future optional parameters
}

// ComputeClientListNodesOptions contains the optional parameters for the ComputeClient.NewListNodesPager method.
type ComputeClientListNodesOptions struct {
	// placeholder for future optional parameters
}

// ComputeClientListOptions contains the optional parameters for the ComputeClient.NewListPager method.
type ComputeClientListOptions struct {
	// Continuation token for pagination.
	Skip *string
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.NewListPager method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpointConnectionsClientCreateOrUpdateOptions contains the optional parameters for the PrivateEndpointConnectionsClient.CreateOrUpdate
// method.
type PrivateEndpointConnectionsClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpointConnectionsClientDeleteOptions contains the optional parameters for the PrivateEndpointConnectionsClient.Delete
// method.
type PrivateEndpointConnectionsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpointConnectionsClientGetOptions contains the optional parameters for the PrivateEndpointConnectionsClient.Get
// method.
type PrivateEndpointConnectionsClientGetOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpointConnectionsClientListOptions contains the optional parameters for the PrivateEndpointConnectionsClient.NewListPager
// method.
type PrivateEndpointConnectionsClientListOptions struct {
	// placeholder for future optional parameters
}

// PrivateLinkResourcesClientListOptions contains the optional parameters for the PrivateLinkResourcesClient.List method.
type PrivateLinkResourcesClientListOptions struct {
	// placeholder for future optional parameters
}

// QuotasClientListOptions contains the optional parameters for the QuotasClient.NewListPager method.
type QuotasClientListOptions struct {
	// placeholder for future optional parameters
}

// QuotasClientUpdateOptions contains the optional parameters for the QuotasClient.Update method.
type QuotasClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// UsagesClientListOptions contains the optional parameters for the UsagesClient.NewListPager method.
type UsagesClientListOptions struct {
	// placeholder for future optional parameters
}

// VirtualMachineSizesClientListOptions contains the optional parameters for the VirtualMachineSizesClient.List method.
type VirtualMachineSizesClientListOptions struct {
	// placeholder for future optional parameters
}

// WorkspaceConnectionsClientCreateOptions contains the optional parameters for the WorkspaceConnectionsClient.Create method.
type WorkspaceConnectionsClientCreateOptions struct {
	// placeholder for future optional parameters
}

// WorkspaceConnectionsClientDeleteOptions contains the optional parameters for the WorkspaceConnectionsClient.Delete method.
type WorkspaceConnectionsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// WorkspaceConnectionsClientGetOptions contains the optional parameters for the WorkspaceConnectionsClient.Get method.
type WorkspaceConnectionsClientGetOptions struct {
	// placeholder for future optional parameters
}

// WorkspaceConnectionsClientListOptions contains the optional parameters for the WorkspaceConnectionsClient.NewListPager
// method.
type WorkspaceConnectionsClientListOptions struct {
	// Category of the workspace connection.
	Category *string

	// Target of the workspace connection.
	Target *string
}

// WorkspaceFeaturesClientListOptions contains the optional parameters for the WorkspaceFeaturesClient.NewListPager method.
type WorkspaceFeaturesClientListOptions struct {
	// placeholder for future optional parameters
}

// WorkspaceSKUsClientListOptions contains the optional parameters for the WorkspaceSKUsClient.NewListPager method.
type WorkspaceSKUsClientListOptions struct {
	// placeholder for future optional parameters
}

// WorkspacesClientBeginCreateOrUpdateOptions contains the optional parameters for the WorkspacesClient.BeginCreateOrUpdate
// method.
type WorkspacesClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// WorkspacesClientBeginDeleteOptions contains the optional parameters for the WorkspacesClient.BeginDelete method.
type WorkspacesClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// WorkspacesClientBeginDiagnoseOptions contains the optional parameters for the WorkspacesClient.BeginDiagnose method.
type WorkspacesClientBeginDiagnoseOptions struct {
	// The parameter of diagnosing workspace health
	Parameters *DiagnoseWorkspaceParameters

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// WorkspacesClientBeginPrepareNotebookOptions contains the optional parameters for the WorkspacesClient.BeginPrepareNotebook
// method.
type WorkspacesClientBeginPrepareNotebookOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// WorkspacesClientBeginResyncKeysOptions contains the optional parameters for the WorkspacesClient.BeginResyncKeys method.
type WorkspacesClientBeginResyncKeysOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// WorkspacesClientGetOptions contains the optional parameters for the WorkspacesClient.Get method.
type WorkspacesClientGetOptions struct {
	// placeholder for future optional parameters
}

// WorkspacesClientListByResourceGroupOptions contains the optional parameters for the WorkspacesClient.NewListByResourceGroupPager
// method.
type WorkspacesClientListByResourceGroupOptions struct {
	// Continuation token for pagination.
	Skip *string
}

// WorkspacesClientListBySubscriptionOptions contains the optional parameters for the WorkspacesClient.NewListBySubscriptionPager
// method.
type WorkspacesClientListBySubscriptionOptions struct {
	// Continuation token for pagination.
	Skip *string
}

// WorkspacesClientListKeysOptions contains the optional parameters for the WorkspacesClient.ListKeys method.
type WorkspacesClientListKeysOptions struct {
	// placeholder for future optional parameters
}

// WorkspacesClientListNotebookAccessTokenOptions contains the optional parameters for the WorkspacesClient.ListNotebookAccessToken
// method.
type WorkspacesClientListNotebookAccessTokenOptions struct {
	// placeholder for future optional parameters
}

// WorkspacesClientListNotebookKeysOptions contains the optional parameters for the WorkspacesClient.ListNotebookKeys method.
type WorkspacesClientListNotebookKeysOptions struct {
	// placeholder for future optional parameters
}

// WorkspacesClientListOutboundNetworkDependenciesEndpointsOptions contains the optional parameters for the WorkspacesClient.ListOutboundNetworkDependenciesEndpoints
// method.
type WorkspacesClientListOutboundNetworkDependenciesEndpointsOptions struct {
	// placeholder for future optional parameters
}

// WorkspacesClientListStorageAccountKeysOptions contains the optional parameters for the WorkspacesClient.ListStorageAccountKeys
// method.
type WorkspacesClientListStorageAccountKeysOptions struct {
	// placeholder for future optional parameters
}

// WorkspacesClientUpdateOptions contains the optional parameters for the WorkspacesClient.Update method.
type WorkspacesClientUpdateOptions struct {
	// placeholder for future optional parameters
}

