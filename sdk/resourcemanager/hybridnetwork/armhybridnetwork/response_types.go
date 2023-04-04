//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armhybridnetwork

// DevicesClientCreateOrUpdateResponse contains the response from method DevicesClient.BeginCreateOrUpdate.
type DevicesClientCreateOrUpdateResponse struct {
	Device
}

// DevicesClientDeleteResponse contains the response from method DevicesClient.BeginDelete.
type DevicesClientDeleteResponse struct {
	// placeholder for future response values
}

// DevicesClientGetResponse contains the response from method DevicesClient.Get.
type DevicesClientGetResponse struct {
	Device
}

// DevicesClientListByResourceGroupResponse contains the response from method DevicesClient.NewListByResourceGroupPager.
type DevicesClientListByResourceGroupResponse struct {
	DeviceListResult
}

// DevicesClientListBySubscriptionResponse contains the response from method DevicesClient.NewListBySubscriptionPager.
type DevicesClientListBySubscriptionResponse struct {
	DeviceListResult
}

// DevicesClientListRegistrationKeyResponse contains the response from method DevicesClient.ListRegistrationKey.
type DevicesClientListRegistrationKeyResponse struct {
	DeviceRegistrationKey
}

// DevicesClientUpdateTagsResponse contains the response from method DevicesClient.UpdateTags.
type DevicesClientUpdateTagsResponse struct {
	Device
}

// NetworkFunctionVendorSKUsClientListBySKUResponse contains the response from method NetworkFunctionVendorSKUsClient.NewListBySKUPager.
type NetworkFunctionVendorSKUsClientListBySKUResponse struct {
	NetworkFunctionSKUDetails
}

// NetworkFunctionVendorSKUsClientListByVendorResponse contains the response from method NetworkFunctionVendorSKUsClient.NewListByVendorPager.
type NetworkFunctionVendorSKUsClientListByVendorResponse struct {
	NetworkFunctionSKUListResult
}

// NetworkFunctionVendorsClientListResponse contains the response from method NetworkFunctionVendorsClient.NewListPager.
type NetworkFunctionVendorsClientListResponse struct {
	NetworkFunctionVendorListResult
}

// NetworkFunctionsClientCreateOrUpdateResponse contains the response from method NetworkFunctionsClient.BeginCreateOrUpdate.
type NetworkFunctionsClientCreateOrUpdateResponse struct {
	NetworkFunction
}

// NetworkFunctionsClientDeleteResponse contains the response from method NetworkFunctionsClient.BeginDelete.
type NetworkFunctionsClientDeleteResponse struct {
	// placeholder for future response values
}

// NetworkFunctionsClientExecuteRequestResponse contains the response from method NetworkFunctionsClient.BeginExecuteRequest.
type NetworkFunctionsClientExecuteRequestResponse struct {
	// placeholder for future response values
}

// NetworkFunctionsClientGetResponse contains the response from method NetworkFunctionsClient.Get.
type NetworkFunctionsClientGetResponse struct {
	NetworkFunction
}

// NetworkFunctionsClientListByResourceGroupResponse contains the response from method NetworkFunctionsClient.NewListByResourceGroupPager.
type NetworkFunctionsClientListByResourceGroupResponse struct {
	NetworkFunctionListResult
}

// NetworkFunctionsClientListBySubscriptionResponse contains the response from method NetworkFunctionsClient.NewListBySubscriptionPager.
type NetworkFunctionsClientListBySubscriptionResponse struct {
	NetworkFunctionListResult
}

// NetworkFunctionsClientUpdateTagsResponse contains the response from method NetworkFunctionsClient.UpdateTags.
type NetworkFunctionsClientUpdateTagsResponse struct {
	NetworkFunction
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	OperationListResult
}

// RoleInstancesClientGetResponse contains the response from method RoleInstancesClient.Get.
type RoleInstancesClientGetResponse struct {
	RoleInstance
}

// RoleInstancesClientListResponse contains the response from method RoleInstancesClient.NewListPager.
type RoleInstancesClientListResponse struct {
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
	VendorNetworkFunction
}

// VendorNetworkFunctionsClientGetResponse contains the response from method VendorNetworkFunctionsClient.Get.
type VendorNetworkFunctionsClientGetResponse struct {
	VendorNetworkFunction
}

// VendorNetworkFunctionsClientListResponse contains the response from method VendorNetworkFunctionsClient.NewListPager.
type VendorNetworkFunctionsClientListResponse struct {
	VendorNetworkFunctionListResult
}

// VendorSKUPreviewClientCreateOrUpdateResponse contains the response from method VendorSKUPreviewClient.BeginCreateOrUpdate.
type VendorSKUPreviewClientCreateOrUpdateResponse struct {
	PreviewSubscription
}

// VendorSKUPreviewClientDeleteResponse contains the response from method VendorSKUPreviewClient.BeginDelete.
type VendorSKUPreviewClientDeleteResponse struct {
	// placeholder for future response values
}

// VendorSKUPreviewClientGetResponse contains the response from method VendorSKUPreviewClient.Get.
type VendorSKUPreviewClientGetResponse struct {
	PreviewSubscription
}

// VendorSKUPreviewClientListResponse contains the response from method VendorSKUPreviewClient.NewListPager.
type VendorSKUPreviewClientListResponse struct {
	PreviewSubscriptionsList
}

// VendorSKUsClientCreateOrUpdateResponse contains the response from method VendorSKUsClient.BeginCreateOrUpdate.
type VendorSKUsClientCreateOrUpdateResponse struct {
	VendorSKU
}

// VendorSKUsClientDeleteResponse contains the response from method VendorSKUsClient.BeginDelete.
type VendorSKUsClientDeleteResponse struct {
	// placeholder for future response values
}

// VendorSKUsClientGetResponse contains the response from method VendorSKUsClient.Get.
type VendorSKUsClientGetResponse struct {
	VendorSKU
}

// VendorSKUsClientListCredentialResponse contains the response from method VendorSKUsClient.ListCredential.
type VendorSKUsClientListCredentialResponse struct {
	SKUCredential
}

// VendorSKUsClientListResponse contains the response from method VendorSKUsClient.NewListPager.
type VendorSKUsClientListResponse struct {
	VendorSKUListResult
}

// VendorsClientCreateOrUpdateResponse contains the response from method VendorsClient.BeginCreateOrUpdate.
type VendorsClientCreateOrUpdateResponse struct {
	Vendor
}

// VendorsClientDeleteResponse contains the response from method VendorsClient.BeginDelete.
type VendorsClientDeleteResponse struct {
	// placeholder for future response values
}

// VendorsClientGetResponse contains the response from method VendorsClient.Get.
type VendorsClientGetResponse struct {
	Vendor
}

// VendorsClientListBySubscriptionResponse contains the response from method VendorsClient.NewListBySubscriptionPager.
type VendorsClientListBySubscriptionResponse struct {
	VendorListResult
}
