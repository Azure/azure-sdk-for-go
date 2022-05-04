//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armportal

// DashboardsClientCreateOrUpdateResponse contains the response from method DashboardsClient.CreateOrUpdate.
type DashboardsClientCreateOrUpdateResponse struct {
	Dashboard
}

// DashboardsClientDeleteResponse contains the response from method DashboardsClient.Delete.
type DashboardsClientDeleteResponse struct {
	// placeholder for future response values
}

// DashboardsClientGetResponse contains the response from method DashboardsClient.Get.
type DashboardsClientGetResponse struct {
	Dashboard
}

// DashboardsClientListByResourceGroupResponse contains the response from method DashboardsClient.ListByResourceGroup.
type DashboardsClientListByResourceGroupResponse struct {
	DashboardListResult
}

// DashboardsClientListBySubscriptionResponse contains the response from method DashboardsClient.ListBySubscription.
type DashboardsClientListBySubscriptionResponse struct {
	DashboardListResult
}

// DashboardsClientUpdateResponse contains the response from method DashboardsClient.Update.
type DashboardsClientUpdateResponse struct {
	Dashboard
}

// ListTenantConfigurationViolationsClientListResponse contains the response from method ListTenantConfigurationViolationsClient.List.
type ListTenantConfigurationViolationsClientListResponse struct {
	ViolationsList
}

// OperationsClientListResponse contains the response from method OperationsClient.List.
type OperationsClientListResponse struct {
	ResourceProviderOperationList
}

// TenantConfigurationsClientCreateResponse contains the response from method TenantConfigurationsClient.Create.
type TenantConfigurationsClientCreateResponse struct {
	Configuration
}

// TenantConfigurationsClientDeleteResponse contains the response from method TenantConfigurationsClient.Delete.
type TenantConfigurationsClientDeleteResponse struct {
	// placeholder for future response values
}

// TenantConfigurationsClientGetResponse contains the response from method TenantConfigurationsClient.Get.
type TenantConfigurationsClientGetResponse struct {
	Configuration
}

// TenantConfigurationsClientListResponse contains the response from method TenantConfigurationsClient.List.
type TenantConfigurationsClientListResponse struct {
	ConfigurationList
}
