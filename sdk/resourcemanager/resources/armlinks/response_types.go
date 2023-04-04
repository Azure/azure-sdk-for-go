//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armlinks

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	OperationListResult
}

// ResourceLinksClientCreateOrUpdateResponse contains the response from method ResourceLinksClient.CreateOrUpdate.
type ResourceLinksClientCreateOrUpdateResponse struct {
	ResourceLink
}

// ResourceLinksClientDeleteResponse contains the response from method ResourceLinksClient.Delete.
type ResourceLinksClientDeleteResponse struct {
	// placeholder for future response values
}

// ResourceLinksClientGetResponse contains the response from method ResourceLinksClient.Get.
type ResourceLinksClientGetResponse struct {
	ResourceLink
}

// ResourceLinksClientListAtSourceScopeResponse contains the response from method ResourceLinksClient.NewListAtSourceScopePager.
type ResourceLinksClientListAtSourceScopeResponse struct {
	ResourceLinkResult
}

// ResourceLinksClientListAtSubscriptionResponse contains the response from method ResourceLinksClient.NewListAtSubscriptionPager.
type ResourceLinksClientListAtSubscriptionResponse struct {
	ResourceLinkResult
}
