//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhanaonazure

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// List of HANA operations
	OperationList
}

// ProviderInstancesClientCreateResponse contains the response from method ProviderInstancesClient.BeginCreate.
type ProviderInstancesClientCreateResponse struct {
	// A provider instance associated with a SAP monitor.
	ProviderInstance
}

// ProviderInstancesClientDeleteResponse contains the response from method ProviderInstancesClient.BeginDelete.
type ProviderInstancesClientDeleteResponse struct {
	// placeholder for future response values
}

// ProviderInstancesClientGetResponse contains the response from method ProviderInstancesClient.Get.
type ProviderInstancesClientGetResponse struct {
	// A provider instance associated with a SAP monitor.
	ProviderInstance
}

// ProviderInstancesClientListResponse contains the response from method ProviderInstancesClient.NewListPager.
type ProviderInstancesClientListResponse struct {
	// The response from the List provider instances operation.
	ProviderInstanceListResult
}

// SapMonitorsClientCreateResponse contains the response from method SapMonitorsClient.BeginCreate.
type SapMonitorsClientCreateResponse struct {
	// SAP monitor info on Azure (ARM properties and SAP monitor properties)
	SapMonitor
}

// SapMonitorsClientDeleteResponse contains the response from method SapMonitorsClient.BeginDelete.
type SapMonitorsClientDeleteResponse struct {
	// placeholder for future response values
}

// SapMonitorsClientGetResponse contains the response from method SapMonitorsClient.Get.
type SapMonitorsClientGetResponse struct {
	// SAP monitor info on Azure (ARM properties and SAP monitor properties)
	SapMonitor
}

// SapMonitorsClientListResponse contains the response from method SapMonitorsClient.NewListPager.
type SapMonitorsClientListResponse struct {
	// The response from the List SAP monitors operation.
	SapMonitorListResult
}

// SapMonitorsClientUpdateResponse contains the response from method SapMonitorsClient.Update.
type SapMonitorsClientUpdateResponse struct {
	// SAP monitor info on Azure (ARM properties and SAP monitor properties)
	SapMonitor
}
