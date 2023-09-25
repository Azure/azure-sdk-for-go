//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armprivatedns

// PrivateZonesClientCreateOrUpdateResponse contains the response from method PrivateZonesClient.BeginCreateOrUpdate.
type PrivateZonesClientCreateOrUpdateResponse struct {
	// Describes a Private DNS zone.
	PrivateZone
}

// PrivateZonesClientDeleteResponse contains the response from method PrivateZonesClient.BeginDelete.
type PrivateZonesClientDeleteResponse struct {
	// placeholder for future response values
}

// PrivateZonesClientGetResponse contains the response from method PrivateZonesClient.Get.
type PrivateZonesClientGetResponse struct {
	// Describes a Private DNS zone.
	PrivateZone
}

// PrivateZonesClientListByResourceGroupResponse contains the response from method PrivateZonesClient.NewListByResourceGroupPager.
type PrivateZonesClientListByResourceGroupResponse struct {
	// The response to a Private DNS zone list operation.
	PrivateZoneListResult
}

// PrivateZonesClientListResponse contains the response from method PrivateZonesClient.NewListPager.
type PrivateZonesClientListResponse struct {
	// The response to a Private DNS zone list operation.
	PrivateZoneListResult
}

// PrivateZonesClientUpdateResponse contains the response from method PrivateZonesClient.BeginUpdate.
type PrivateZonesClientUpdateResponse struct {
	// Describes a Private DNS zone.
	PrivateZone
}

// RecordSetsClientCreateOrUpdateResponse contains the response from method RecordSetsClient.CreateOrUpdate.
type RecordSetsClientCreateOrUpdateResponse struct {
	// Describes a DNS record set (a collection of DNS records with the same name and type) in a Private DNS zone.
	RecordSet
}

// RecordSetsClientDeleteResponse contains the response from method RecordSetsClient.Delete.
type RecordSetsClientDeleteResponse struct {
	// placeholder for future response values
}

// RecordSetsClientGetResponse contains the response from method RecordSetsClient.Get.
type RecordSetsClientGetResponse struct {
	// Describes a DNS record set (a collection of DNS records with the same name and type) in a Private DNS zone.
	RecordSet
}

// RecordSetsClientListByTypeResponse contains the response from method RecordSetsClient.NewListByTypePager.
type RecordSetsClientListByTypeResponse struct {
	// The response to a record set list operation.
	RecordSetListResult
}

// RecordSetsClientListResponse contains the response from method RecordSetsClient.NewListPager.
type RecordSetsClientListResponse struct {
	// The response to a record set list operation.
	RecordSetListResult
}

// RecordSetsClientUpdateResponse contains the response from method RecordSetsClient.Update.
type RecordSetsClientUpdateResponse struct {
	// Describes a DNS record set (a collection of DNS records with the same name and type) in a Private DNS zone.
	RecordSet
}

// VirtualNetworkLinksClientCreateOrUpdateResponse contains the response from method VirtualNetworkLinksClient.BeginCreateOrUpdate.
type VirtualNetworkLinksClientCreateOrUpdateResponse struct {
	// Describes a link to virtual network for a Private DNS zone.
	VirtualNetworkLink
}

// VirtualNetworkLinksClientDeleteResponse contains the response from method VirtualNetworkLinksClient.BeginDelete.
type VirtualNetworkLinksClientDeleteResponse struct {
	// placeholder for future response values
}

// VirtualNetworkLinksClientGetResponse contains the response from method VirtualNetworkLinksClient.Get.
type VirtualNetworkLinksClientGetResponse struct {
	// Describes a link to virtual network for a Private DNS zone.
	VirtualNetworkLink
}

// VirtualNetworkLinksClientListResponse contains the response from method VirtualNetworkLinksClient.NewListPager.
type VirtualNetworkLinksClientListResponse struct {
	// The response to a list virtual network link to Private DNS zone operation.
	VirtualNetworkLinkListResult
}

// VirtualNetworkLinksClientUpdateResponse contains the response from method VirtualNetworkLinksClient.BeginUpdate.
type VirtualNetworkLinksClientUpdateResponse struct {
	// Describes a link to virtual network for a Private DNS zone.
	VirtualNetworkLink
}

