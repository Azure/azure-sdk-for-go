//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdynatrace

// MonitorsClientCreateOrUpdateResponse contains the response from method MonitorsClient.BeginCreateOrUpdate.
type MonitorsClientCreateOrUpdateResponse struct {
	// Dynatrace Monitor Resource
	MonitorResource
}

// MonitorsClientDeleteResponse contains the response from method MonitorsClient.BeginDelete.
type MonitorsClientDeleteResponse struct {
	// placeholder for future response values
}

// MonitorsClientGetMarketplaceSaaSResourceDetailsResponse contains the response from method MonitorsClient.GetMarketplaceSaaSResourceDetails.
type MonitorsClientGetMarketplaceSaaSResourceDetailsResponse struct {
	// Marketplace SaaS resource details linked to the given tenant Id
	MarketplaceSaaSResourceDetailsResponse
}

// MonitorsClientGetMetricStatusResponse contains the response from method MonitorsClient.GetMetricStatus.
type MonitorsClientGetMetricStatusResponse struct {
	// Response of get metrics status operation
	MetricsStatusResponse
}

// MonitorsClientGetResponse contains the response from method MonitorsClient.Get.
type MonitorsClientGetResponse struct {
	// Dynatrace Monitor Resource
	MonitorResource
}

// MonitorsClientGetSSODetailsResponse contains the response from method MonitorsClient.GetSSODetails.
type MonitorsClientGetSSODetailsResponse struct {
	// SSO details from the Dynatrace partner
	SSODetailsResponse
}

// MonitorsClientGetVMHostPayloadResponse contains the response from method MonitorsClient.GetVMHostPayload.
type MonitorsClientGetVMHostPayloadResponse struct {
	// Response of payload to be passed while installing VM agent.
	VMExtensionPayload
}

// MonitorsClientListAppServicesResponse contains the response from method MonitorsClient.NewListAppServicesPager.
type MonitorsClientListAppServicesResponse struct {
	// Response of a list App Services Operation.
	AppServiceListResponse
}

// MonitorsClientListByResourceGroupResponse contains the response from method MonitorsClient.NewListByResourceGroupPager.
type MonitorsClientListByResourceGroupResponse struct {
	// The response of a MonitorResource list operation.
	MonitorResourceListResult
}

// MonitorsClientListBySubscriptionIDResponse contains the response from method MonitorsClient.NewListBySubscriptionIDPager.
type MonitorsClientListBySubscriptionIDResponse struct {
	// The response of a MonitorResource list operation.
	MonitorResourceListResult
}

// MonitorsClientListHostsResponse contains the response from method MonitorsClient.NewListHostsPager.
type MonitorsClientListHostsResponse struct {
	// Response of a list VM Host Operation.
	VMHostsListResponse
}

// MonitorsClientListLinkableEnvironmentsResponse contains the response from method MonitorsClient.NewListLinkableEnvironmentsPager.
type MonitorsClientListLinkableEnvironmentsResponse struct {
	// Response for getting all the linkable environments
	LinkableEnvironmentListResponse
}

// MonitorsClientListMonitoredResourcesResponse contains the response from method MonitorsClient.NewListMonitoredResourcesPager.
type MonitorsClientListMonitoredResourcesResponse struct {
	// List of all the resources being monitored by Dynatrace monitor resource
	MonitoredResourceListResponse
}

// MonitorsClientUpdateResponse contains the response from method MonitorsClient.Update.
type MonitorsClientUpdateResponse struct {
	// Dynatrace Monitor Resource
	MonitorResource
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to get the next set of results.
	OperationListResult
}

// SingleSignOnClientCreateOrUpdateResponse contains the response from method SingleSignOnClient.BeginCreateOrUpdate.
type SingleSignOnClientCreateOrUpdateResponse struct {
	// Single sign-on configurations for a given monitor resource.
	SingleSignOnResource
}

// SingleSignOnClientGetResponse contains the response from method SingleSignOnClient.Get.
type SingleSignOnClientGetResponse struct {
	// Single sign-on configurations for a given monitor resource.
	SingleSignOnResource
}

// SingleSignOnClientListResponse contains the response from method SingleSignOnClient.NewListPager.
type SingleSignOnClientListResponse struct {
	// The response of a DynatraceSingleSignOnResource list operation.
	SingleSignOnResourceListResult
}

// TagRulesClientCreateOrUpdateResponse contains the response from method TagRulesClient.BeginCreateOrUpdate.
type TagRulesClientCreateOrUpdateResponse struct {
	// Tag rules for a monitor resource
	TagRule
}

// TagRulesClientDeleteResponse contains the response from method TagRulesClient.BeginDelete.
type TagRulesClientDeleteResponse struct {
	// placeholder for future response values
}

// TagRulesClientGetResponse contains the response from method TagRulesClient.Get.
type TagRulesClientGetResponse struct {
	// Tag rules for a monitor resource
	TagRule
}

// TagRulesClientListResponse contains the response from method TagRulesClient.NewListPager.
type TagRulesClientListResponse struct {
	// The response of a TagRule list operation.
	TagRuleListResult
}

