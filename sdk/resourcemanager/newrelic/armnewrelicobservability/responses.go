//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnewrelicobservability

// AccountsClientListResponse contains the response from method AccountsClient.NewListPager.
type AccountsClientListResponse struct {
	// Response of get all accounts Operation.
	AccountsListResponse
}

// BillingInfoClientGetResponse contains the response from method BillingInfoClient.Get.
type BillingInfoClientGetResponse struct {
	// Marketplace Subscription and Organization details to which resource gets billed into.
	BillingInfoResponse
}

// ConnectedPartnerResourcesClientListResponse contains the response from method ConnectedPartnerResourcesClient.NewListPager.
type ConnectedPartnerResourcesClientListResponse struct {
	// List of all active newrelic deployments.
	ConnectedPartnerResourcesListResponse
}

// MonitoredSubscriptionsClientCreateorUpdateResponse contains the response from method MonitoredSubscriptionsClient.BeginCreateorUpdate.
type MonitoredSubscriptionsClientCreateorUpdateResponse struct {
	// The request to update subscriptions needed to be monitored by the NewRelic monitor resource.
	MonitoredSubscriptionProperties
}

// MonitoredSubscriptionsClientDeleteResponse contains the response from method MonitoredSubscriptionsClient.BeginDelete.
type MonitoredSubscriptionsClientDeleteResponse struct {
	// placeholder for future response values
}

// MonitoredSubscriptionsClientGetResponse contains the response from method MonitoredSubscriptionsClient.Get.
type MonitoredSubscriptionsClientGetResponse struct {
	// The request to update subscriptions needed to be monitored by the NewRelic monitor resource.
	MonitoredSubscriptionProperties
}

// MonitoredSubscriptionsClientListResponse contains the response from method MonitoredSubscriptionsClient.NewListPager.
type MonitoredSubscriptionsClientListResponse struct {
	MonitoredSubscriptionPropertiesList
}

// MonitoredSubscriptionsClientUpdateResponse contains the response from method MonitoredSubscriptionsClient.BeginUpdate.
type MonitoredSubscriptionsClientUpdateResponse struct {
	// The request to update subscriptions needed to be monitored by the NewRelic monitor resource.
	MonitoredSubscriptionProperties
}

// MonitorsClientCreateOrUpdateResponse contains the response from method MonitorsClient.BeginCreateOrUpdate.
type MonitorsClientCreateOrUpdateResponse struct {
	// A Monitor Resource by NewRelic
	NewRelicMonitorResource
}

// MonitorsClientDeleteResponse contains the response from method MonitorsClient.BeginDelete.
type MonitorsClientDeleteResponse struct {
	// placeholder for future response values
}

// MonitorsClientGetMetricRulesResponse contains the response from method MonitorsClient.GetMetricRules.
type MonitorsClientGetMetricRulesResponse struct {
	// Set of rules for sending metrics for the Monitor resource.
	MetricRules
}

// MonitorsClientGetMetricStatusResponse contains the response from method MonitorsClient.GetMetricStatus.
type MonitorsClientGetMetricStatusResponse struct {
	// Response of get metrics status Operation.
	MetricsStatusResponse
}

// MonitorsClientGetResponse contains the response from method MonitorsClient.Get.
type MonitorsClientGetResponse struct {
	// A Monitor Resource by NewRelic
	NewRelicMonitorResource
}

// MonitorsClientListAppServicesResponse contains the response from method MonitorsClient.NewListAppServicesPager.
type MonitorsClientListAppServicesResponse struct {
	// Response of a list app services Operation.
	AppServicesListResponse
}

// MonitorsClientListByResourceGroupResponse contains the response from method MonitorsClient.NewListByResourceGroupPager.
type MonitorsClientListByResourceGroupResponse struct {
	// The response of a NewRelicMonitorResource list operation.
	NewRelicMonitorResourceListResult
}

// MonitorsClientListBySubscriptionResponse contains the response from method MonitorsClient.NewListBySubscriptionPager.
type MonitorsClientListBySubscriptionResponse struct {
	// The response of a NewRelicMonitorResource list operation.
	NewRelicMonitorResourceListResult
}

// MonitorsClientListHostsResponse contains the response from method MonitorsClient.NewListHostsPager.
type MonitorsClientListHostsResponse struct {
	// Response of a list VM Host Operation.
	VMHostsListResponse
}

// MonitorsClientListLinkedResourcesResponse contains the response from method MonitorsClient.NewListLinkedResourcesPager.
type MonitorsClientListLinkedResourcesResponse struct {
	// Response of a list operation.
	LinkedResourceListResponse
}

// MonitorsClientListMonitoredResourcesResponse contains the response from method MonitorsClient.NewListMonitoredResourcesPager.
type MonitorsClientListMonitoredResourcesResponse struct {
	// List of all the resources being monitored by NewRelic monitor resource
	MonitoredResourceListResponse
}

// MonitorsClientSwitchBillingResponse contains the response from method MonitorsClient.SwitchBilling.
type MonitorsClientSwitchBillingResponse struct {
	// A Monitor Resource by NewRelic
	NewRelicMonitorResource

	// RetryAfter contains the information returned from the Retry-After header response.
	RetryAfter *int32
}

// MonitorsClientUpdateResponse contains the response from method MonitorsClient.Update.
type MonitorsClientUpdateResponse struct {
	// A Monitor Resource by NewRelic
	NewRelicMonitorResource
}

// MonitorsClientVMHostPayloadResponse contains the response from method MonitorsClient.VMHostPayload.
type MonitorsClientVMHostPayloadResponse struct {
	// Response of payload to be passed while installing VM agent.
	VMExtensionPayload
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to get the next set of results.
	OperationListResult
}

// OrganizationsClientListResponse contains the response from method OrganizationsClient.NewListPager.
type OrganizationsClientListResponse struct {
	// Response of get all organizations Operation.
	OrganizationsListResponse
}

// PlansClientListResponse contains the response from method PlansClient.NewListPager.
type PlansClientListResponse struct {
	// Response of get all plan data Operation.
	PlanDataListResponse
}

// TagRulesClientCreateOrUpdateResponse contains the response from method TagRulesClient.BeginCreateOrUpdate.
type TagRulesClientCreateOrUpdateResponse struct {
	// A tag rule belonging to NewRelic account
	TagRule
}

// TagRulesClientDeleteResponse contains the response from method TagRulesClient.BeginDelete.
type TagRulesClientDeleteResponse struct {
	// placeholder for future response values
}

// TagRulesClientGetResponse contains the response from method TagRulesClient.Get.
type TagRulesClientGetResponse struct {
	// A tag rule belonging to NewRelic account
	TagRule
}

// TagRulesClientListByNewRelicMonitorResourceResponse contains the response from method TagRulesClient.NewListByNewRelicMonitorResourcePager.
type TagRulesClientListByNewRelicMonitorResourceResponse struct {
	// The response of a TagRule list operation.
	TagRuleListResult
}

// TagRulesClientUpdateResponse contains the response from method TagRulesClient.Update.
type TagRulesClientUpdateResponse struct {
	// A tag rule belonging to NewRelic account
	TagRule
}
