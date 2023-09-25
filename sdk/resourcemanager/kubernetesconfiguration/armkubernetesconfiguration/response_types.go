//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armkubernetesconfiguration

// ExtensionsClientCreateResponse contains the response from method ExtensionsClient.BeginCreate.
type ExtensionsClientCreateResponse struct {
	// The Extension object.
	Extension
}

// ExtensionsClientDeleteResponse contains the response from method ExtensionsClient.BeginDelete.
type ExtensionsClientDeleteResponse struct {
	// placeholder for future response values
}

// ExtensionsClientGetResponse contains the response from method ExtensionsClient.Get.
type ExtensionsClientGetResponse struct {
	// The Extension object.
	Extension
}

// ExtensionsClientListResponse contains the response from method ExtensionsClient.NewListPager.
type ExtensionsClientListResponse struct {
	// Result of the request to list Extensions. It contains a list of Extension objects and a URL link to get the next set of
// results.
	ExtensionsList
}

// ExtensionsClientUpdateResponse contains the response from method ExtensionsClient.BeginUpdate.
type ExtensionsClientUpdateResponse struct {
	// The Extension object.
	Extension
}

// FluxConfigOperationStatusClientGetResponse contains the response from method FluxConfigOperationStatusClient.Get.
type FluxConfigOperationStatusClientGetResponse struct {
	// The current status of an async operation.
	OperationStatusResult
}

// FluxConfigurationsClientCreateOrUpdateResponse contains the response from method FluxConfigurationsClient.BeginCreateOrUpdate.
type FluxConfigurationsClientCreateOrUpdateResponse struct {
	// The Flux Configuration object returned in Get & Put response.
	FluxConfiguration
}

// FluxConfigurationsClientDeleteResponse contains the response from method FluxConfigurationsClient.BeginDelete.
type FluxConfigurationsClientDeleteResponse struct {
	// placeholder for future response values
}

// FluxConfigurationsClientGetResponse contains the response from method FluxConfigurationsClient.Get.
type FluxConfigurationsClientGetResponse struct {
	// The Flux Configuration object returned in Get & Put response.
	FluxConfiguration
}

// FluxConfigurationsClientListResponse contains the response from method FluxConfigurationsClient.NewListPager.
type FluxConfigurationsClientListResponse struct {
	// Result of the request to list Flux Configurations. It contains a list of FluxConfiguration objects and a URL link to get
// the next set of results.
	FluxConfigurationsList
}

// FluxConfigurationsClientUpdateResponse contains the response from method FluxConfigurationsClient.BeginUpdate.
type FluxConfigurationsClientUpdateResponse struct {
	// The Flux Configuration object returned in Get & Put response.
	FluxConfiguration
}

// OperationStatusClientGetResponse contains the response from method OperationStatusClient.Get.
type OperationStatusClientGetResponse struct {
	// The current status of an async operation.
	OperationStatusResult
}

// OperationStatusClientListResponse contains the response from method OperationStatusClient.NewListPager.
type OperationStatusClientListResponse struct {
	// The async operations in progress, in the cluster.
	OperationStatusList
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// Result of the request to list operations.
	ResourceProviderOperationList
}

// SourceControlConfigurationsClientCreateOrUpdateResponse contains the response from method SourceControlConfigurationsClient.CreateOrUpdate.
type SourceControlConfigurationsClientCreateOrUpdateResponse struct {
	// The SourceControl Configuration object returned in Get & Put response.
	SourceControlConfiguration
}

// SourceControlConfigurationsClientDeleteResponse contains the response from method SourceControlConfigurationsClient.BeginDelete.
type SourceControlConfigurationsClientDeleteResponse struct {
	// placeholder for future response values
}

// SourceControlConfigurationsClientGetResponse contains the response from method SourceControlConfigurationsClient.Get.
type SourceControlConfigurationsClientGetResponse struct {
	// The SourceControl Configuration object returned in Get & Put response.
	SourceControlConfiguration
}

// SourceControlConfigurationsClientListResponse contains the response from method SourceControlConfigurationsClient.NewListPager.
type SourceControlConfigurationsClientListResponse struct {
	// Result of the request to list Source Control Configurations. It contains a list of SourceControlConfiguration objects and
// a URL link to get the next set of results.
	SourceControlConfigurationList
}

