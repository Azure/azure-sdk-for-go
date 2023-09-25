//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybridnetwork

// DevicesClientCreateOrUpdateResponse contains the response from method DevicesClient.BeginCreateOrUpdate.
type DevicesClientCreateOrUpdateResponse struct {
	// Device resource.
	Device
}

// DevicesClientDeleteResponse contains the response from method DevicesClient.BeginDelete.
type DevicesClientDeleteResponse struct {
	// placeholder for future response values
}

// DevicesClientGetResponse contains the response from method DevicesClient.Get.
type DevicesClientGetResponse struct {
	// Device resource.
	Device
}

// DevicesClientListByResourceGroupResponse contains the response from method DevicesClient.NewListByResourceGroupPager.
type DevicesClientListByResourceGroupResponse struct {
	// Response for devices API service call.
	DeviceListResult
}

// DevicesClientListBySubscriptionResponse contains the response from method DevicesClient.NewListBySubscriptionPager.
type DevicesClientListBySubscriptionResponse struct {
	// Response for devices API service call.
	DeviceListResult
}

// DevicesClientListRegistrationKeyResponse contains the response from method DevicesClient.ListRegistrationKey.
type DevicesClientListRegistrationKeyResponse struct {
	// The device registration key.
	DeviceRegistrationKey
}

// DevicesClientUpdateTagsResponse contains the response from method DevicesClient.UpdateTags.
type DevicesClientUpdateTagsResponse struct {
	// Device resource.
	Device
}

// NetworkFunctionVendorSKUsClientListBySKUResponse contains the response from method NetworkFunctionVendorSKUsClient.NewListBySKUPager.
type NetworkFunctionVendorSKUsClientListBySKUResponse struct {
	// The network function sku details.
	NetworkFunctionSKUDetails
}

// NetworkFunctionVendorSKUsClientListByVendorResponse contains the response from method NetworkFunctionVendorSKUsClient.NewListByVendorPager.
type NetworkFunctionVendorSKUsClientListByVendorResponse struct {
	// A list of available network function skus.
	NetworkFunctionSKUListResult
}

// NetworkFunctionVendorsClientListResponse contains the response from method NetworkFunctionVendorsClient.NewListPager.
type NetworkFunctionVendorsClientListResponse struct {
	// The network function vendor list result.
	NetworkFunctionVendorListResult
}

// NetworkFunctionsClientCreateOrUpdateResponse contains the response from method NetworkFunctionsClient.BeginCreateOrUpdate.
type NetworkFunctionsClientCreateOrUpdateResponse struct {
	// Network function resource response.
	NetworkFunction
}

// NetworkFunctionsClientDeleteResponse contains the response from method NetworkFunctionsClient.BeginDelete.
type NetworkFunctionsClientDeleteResponse struct {
	// placeholder for future response values
}

// NetworkFunctionsClientGetResponse contains the response from method NetworkFunctionsClient.Get.
type NetworkFunctionsClientGetResponse struct {
	// Network function resource response.
	NetworkFunction
}

// NetworkFunctionsClientListByResourceGroupResponse contains the response from method NetworkFunctionsClient.NewListByResourceGroupPager.
type NetworkFunctionsClientListByResourceGroupResponse struct {
	// Response for network function API service call.
	NetworkFunctionListResult
}

// NetworkFunctionsClientListBySubscriptionResponse contains the response from method NetworkFunctionsClient.NewListBySubscriptionPager.
type NetworkFunctionsClientListBySubscriptionResponse struct {
	// Response for network function API service call.
	NetworkFunctionListResult
}

// NetworkFunctionsClientUpdateTagsResponse contains the response from method NetworkFunctionsClient.UpdateTags.
type NetworkFunctionsClientUpdateTagsResponse struct {
	// Network function resource response.
	NetworkFunction
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// A list of the operations.
	OperationList
}

// RoleInstancesClientGetResponse contains the response from method RoleInstancesClient.Get.
type RoleInstancesClientGetResponse struct {
	// The role instance sub resource.
	RoleInstance
}

// RoleInstancesClientListResponse contains the response from method RoleInstancesClient.NewListPager.
type RoleInstancesClientListResponse struct {
	// List of role instances of vendor network function.
	NetworkFunctionRoleInstanceListResult
}

// RoleInstancesClientRestartResponse contains the response from method RoleInstancesClient.BeginRestart.
type RoleInstancesClientRestartResponse struct {
	// placeholder for future response values
}

// RoleInstancesClientStartResponse contains the response from method RoleInstancesClient.BeginStart.
type RoleInstancesClientStartResponse struct {
	// placeholder for future response values
}

// RoleInstancesClientStopResponse contains the response from method RoleInstancesClient.BeginStop.
type RoleInstancesClientStopResponse struct {
	// placeholder for future response values
}

// VendorNetworkFunctionsClientCreateOrUpdateResponse contains the response from method VendorNetworkFunctionsClient.BeginCreateOrUpdate.
type VendorNetworkFunctionsClientCreateOrUpdateResponse struct {
	// Vendor network function sub resource.
	VendorNetworkFunction
}

// VendorNetworkFunctionsClientGetResponse contains the response from method VendorNetworkFunctionsClient.Get.
type VendorNetworkFunctionsClientGetResponse struct {
	// Vendor network function sub resource.
	VendorNetworkFunction
}

// VendorNetworkFunctionsClientListResponse contains the response from method VendorNetworkFunctionsClient.NewListPager.
type VendorNetworkFunctionsClientListResponse struct {
	// Response for vendors API service call.
	VendorNetworkFunctionListResult
}

// VendorSKUPreviewClientCreateOrUpdateResponse contains the response from method VendorSKUPreviewClient.BeginCreateOrUpdate.
type VendorSKUPreviewClientCreateOrUpdateResponse struct {
	// Customer subscription which can use a sku.
	PreviewSubscription
}

// VendorSKUPreviewClientDeleteResponse contains the response from method VendorSKUPreviewClient.BeginDelete.
type VendorSKUPreviewClientDeleteResponse struct {
	// placeholder for future response values
}

// VendorSKUPreviewClientGetResponse contains the response from method VendorSKUPreviewClient.Get.
type VendorSKUPreviewClientGetResponse struct {
	// Customer subscription which can use a sku.
	PreviewSubscription
}

// VendorSKUPreviewClientListResponse contains the response from method VendorSKUPreviewClient.NewListPager.
type VendorSKUPreviewClientListResponse struct {
	// A list of customer subscriptions which can use a sku.
	PreviewSubscriptionsList
}

// VendorSKUsClientCreateOrUpdateResponse contains the response from method VendorSKUsClient.BeginCreateOrUpdate.
type VendorSKUsClientCreateOrUpdateResponse struct {
	// Sku sub resource.
	VendorSKU
}

// VendorSKUsClientDeleteResponse contains the response from method VendorSKUsClient.BeginDelete.
type VendorSKUsClientDeleteResponse struct {
	// placeholder for future response values
}

// VendorSKUsClientGetResponse contains the response from method VendorSKUsClient.Get.
type VendorSKUsClientGetResponse struct {
	// Sku sub resource.
	VendorSKU
}

// VendorSKUsClientListResponse contains the response from method VendorSKUsClient.NewListPager.
type VendorSKUsClientListResponse struct {
	// Response for list vendor sku API service call.
	VendorSKUListResult
}

// VendorsClientCreateOrUpdateResponse contains the response from method VendorsClient.BeginCreateOrUpdate.
type VendorsClientCreateOrUpdateResponse struct {
	// Vendor resource.
	Vendor
}

// VendorsClientDeleteResponse contains the response from method VendorsClient.BeginDelete.
type VendorsClientDeleteResponse struct {
	// placeholder for future response values
}

// VendorsClientGetResponse contains the response from method VendorsClient.Get.
type VendorsClientGetResponse struct {
	// Vendor resource.
	Vendor
}

// VendorsClientListBySubscriptionResponse contains the response from method VendorsClient.NewListBySubscriptionPager.
type VendorsClientListBySubscriptionResponse struct {
	// Response for vendors API service call.
	VendorListResult
}

