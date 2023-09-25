//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdnsresolver

// DNSForwardingRulesetsClientCreateOrUpdateResponse contains the response from method DNSForwardingRulesetsClient.BeginCreateOrUpdate.
type DNSForwardingRulesetsClientCreateOrUpdateResponse struct {
	// Describes a DNS forwarding ruleset.
	DNSForwardingRuleset
}

// DNSForwardingRulesetsClientDeleteResponse contains the response from method DNSForwardingRulesetsClient.BeginDelete.
type DNSForwardingRulesetsClientDeleteResponse struct {
	// placeholder for future response values
}

// DNSForwardingRulesetsClientGetResponse contains the response from method DNSForwardingRulesetsClient.Get.
type DNSForwardingRulesetsClientGetResponse struct {
	// Describes a DNS forwarding ruleset.
	DNSForwardingRuleset
}

// DNSForwardingRulesetsClientListByResourceGroupResponse contains the response from method DNSForwardingRulesetsClient.NewListByResourceGroupPager.
type DNSForwardingRulesetsClientListByResourceGroupResponse struct {
	// The response to an enumeration operation on DNS forwarding rulesets.
	DNSForwardingRulesetListResult
}

// DNSForwardingRulesetsClientListByVirtualNetworkResponse contains the response from method DNSForwardingRulesetsClient.NewListByVirtualNetworkPager.
type DNSForwardingRulesetsClientListByVirtualNetworkResponse struct {
	// The response to an enumeration operation on Virtual Network DNS Forwarding Ruleset.
	VirtualNetworkDNSForwardingRulesetListResult
}

// DNSForwardingRulesetsClientListResponse contains the response from method DNSForwardingRulesetsClient.NewListPager.
type DNSForwardingRulesetsClientListResponse struct {
	// The response to an enumeration operation on DNS forwarding rulesets.
	DNSForwardingRulesetListResult
}

// DNSForwardingRulesetsClientUpdateResponse contains the response from method DNSForwardingRulesetsClient.BeginUpdate.
type DNSForwardingRulesetsClientUpdateResponse struct {
	// Describes a DNS forwarding ruleset.
	DNSForwardingRuleset
}

// DNSResolversClientCreateOrUpdateResponse contains the response from method DNSResolversClient.BeginCreateOrUpdate.
type DNSResolversClientCreateOrUpdateResponse struct {
	// Describes a DNS resolver.
	DNSResolver
}

// DNSResolversClientDeleteResponse contains the response from method DNSResolversClient.BeginDelete.
type DNSResolversClientDeleteResponse struct {
	// placeholder for future response values
}

// DNSResolversClientGetResponse contains the response from method DNSResolversClient.Get.
type DNSResolversClientGetResponse struct {
	// Describes a DNS resolver.
	DNSResolver
}

// DNSResolversClientListByResourceGroupResponse contains the response from method DNSResolversClient.NewListByResourceGroupPager.
type DNSResolversClientListByResourceGroupResponse struct {
	// The response to an enumeration operation on DNS resolvers.
	ListResult
}

// DNSResolversClientListByVirtualNetworkResponse contains the response from method DNSResolversClient.NewListByVirtualNetworkPager.
type DNSResolversClientListByVirtualNetworkResponse struct {
	// The response to an enumeration operation on sub-resources.
	SubResourceListResult
}

// DNSResolversClientListResponse contains the response from method DNSResolversClient.NewListPager.
type DNSResolversClientListResponse struct {
	// The response to an enumeration operation on DNS resolvers.
	ListResult
}

// DNSResolversClientUpdateResponse contains the response from method DNSResolversClient.BeginUpdate.
type DNSResolversClientUpdateResponse struct {
	// Describes a DNS resolver.
	DNSResolver
}

// ForwardingRulesClientCreateOrUpdateResponse contains the response from method ForwardingRulesClient.CreateOrUpdate.
type ForwardingRulesClientCreateOrUpdateResponse struct {
	// Describes a forwarding rule within a DNS forwarding ruleset.
	ForwardingRule
}

// ForwardingRulesClientDeleteResponse contains the response from method ForwardingRulesClient.Delete.
type ForwardingRulesClientDeleteResponse struct {
	// placeholder for future response values
}

// ForwardingRulesClientGetResponse contains the response from method ForwardingRulesClient.Get.
type ForwardingRulesClientGetResponse struct {
	// Describes a forwarding rule within a DNS forwarding ruleset.
	ForwardingRule
}

// ForwardingRulesClientListResponse contains the response from method ForwardingRulesClient.NewListPager.
type ForwardingRulesClientListResponse struct {
	// The response to an enumeration operation on forwarding rules within a DNS forwarding ruleset.
	ForwardingRuleListResult
}

// ForwardingRulesClientUpdateResponse contains the response from method ForwardingRulesClient.Update.
type ForwardingRulesClientUpdateResponse struct {
	// Describes a forwarding rule within a DNS forwarding ruleset.
	ForwardingRule
}

// InboundEndpointsClientCreateOrUpdateResponse contains the response from method InboundEndpointsClient.BeginCreateOrUpdate.
type InboundEndpointsClientCreateOrUpdateResponse struct {
	// Describes an inbound endpoint for a DNS resolver.
	InboundEndpoint
}

// InboundEndpointsClientDeleteResponse contains the response from method InboundEndpointsClient.BeginDelete.
type InboundEndpointsClientDeleteResponse struct {
	// placeholder for future response values
}

// InboundEndpointsClientGetResponse contains the response from method InboundEndpointsClient.Get.
type InboundEndpointsClientGetResponse struct {
	// Describes an inbound endpoint for a DNS resolver.
	InboundEndpoint
}

// InboundEndpointsClientListResponse contains the response from method InboundEndpointsClient.NewListPager.
type InboundEndpointsClientListResponse struct {
	// The response to an enumeration operation on inbound endpoints for a DNS resolver.
	InboundEndpointListResult
}

// InboundEndpointsClientUpdateResponse contains the response from method InboundEndpointsClient.BeginUpdate.
type InboundEndpointsClientUpdateResponse struct {
	// Describes an inbound endpoint for a DNS resolver.
	InboundEndpoint
}

// OutboundEndpointsClientCreateOrUpdateResponse contains the response from method OutboundEndpointsClient.BeginCreateOrUpdate.
type OutboundEndpointsClientCreateOrUpdateResponse struct {
	// Describes an outbound endpoint for a DNS resolver.
	OutboundEndpoint
}

// OutboundEndpointsClientDeleteResponse contains the response from method OutboundEndpointsClient.BeginDelete.
type OutboundEndpointsClientDeleteResponse struct {
	// placeholder for future response values
}

// OutboundEndpointsClientGetResponse contains the response from method OutboundEndpointsClient.Get.
type OutboundEndpointsClientGetResponse struct {
	// Describes an outbound endpoint for a DNS resolver.
	OutboundEndpoint
}

// OutboundEndpointsClientListResponse contains the response from method OutboundEndpointsClient.NewListPager.
type OutboundEndpointsClientListResponse struct {
	// The response to an enumeration operation on outbound endpoints for a DNS resolver.
	OutboundEndpointListResult
}

// OutboundEndpointsClientUpdateResponse contains the response from method OutboundEndpointsClient.BeginUpdate.
type OutboundEndpointsClientUpdateResponse struct {
	// Describes an outbound endpoint for a DNS resolver.
	OutboundEndpoint
}

// VirtualNetworkLinksClientCreateOrUpdateResponse contains the response from method VirtualNetworkLinksClient.BeginCreateOrUpdate.
type VirtualNetworkLinksClientCreateOrUpdateResponse struct {
	// Describes a virtual network link.
	VirtualNetworkLink
}

// VirtualNetworkLinksClientDeleteResponse contains the response from method VirtualNetworkLinksClient.BeginDelete.
type VirtualNetworkLinksClientDeleteResponse struct {
	// placeholder for future response values
}

// VirtualNetworkLinksClientGetResponse contains the response from method VirtualNetworkLinksClient.Get.
type VirtualNetworkLinksClientGetResponse struct {
	// Describes a virtual network link.
	VirtualNetworkLink
}

// VirtualNetworkLinksClientListResponse contains the response from method VirtualNetworkLinksClient.NewListPager.
type VirtualNetworkLinksClientListResponse struct {
	// The response to an enumeration operation on virtual network links.
	VirtualNetworkLinkListResult
}

// VirtualNetworkLinksClientUpdateResponse contains the response from method VirtualNetworkLinksClient.BeginUpdate.
type VirtualNetworkLinksClientUpdateResponse struct {
	// Describes a virtual network link.
	VirtualNetworkLink
}

