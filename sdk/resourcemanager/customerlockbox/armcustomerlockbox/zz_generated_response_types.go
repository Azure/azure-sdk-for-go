//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcustomerlockbox

// GetClientTenantOptedInResponse contains the response from method GetClient.TenantOptedIn.
type GetClientTenantOptedInResponse struct {
	TenantOptInResponse
}

// OperationsClientListResponse contains the response from method OperationsClient.List.
type OperationsClientListResponse struct {
	OperationListResult
}

// PostClientDisableLockboxResponse contains the response from method PostClient.DisableLockbox.
type PostClientDisableLockboxResponse struct {
	// placeholder for future response values
}

// PostClientEnableLockboxResponse contains the response from method PostClient.EnableLockbox.
type PostClientEnableLockboxResponse struct {
	// placeholder for future response values
}

// RequestsClientGetResponse contains the response from method RequestsClient.Get.
type RequestsClientGetResponse struct {
	LockboxRequestResponse
}

// RequestsClientListResponse contains the response from method RequestsClient.List.
type RequestsClientListResponse struct {
	RequestListResult
}

// RequestsClientUpdateStatusResponse contains the response from method RequestsClient.UpdateStatus.
type RequestsClientUpdateStatusResponse struct {
	Approval
}
