//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armprivatedns

// PrivateZonesClientBeginCreateOrUpdateOptions contains the optional parameters for the PrivateZonesClient.BeginCreateOrUpdate
// method.
type PrivateZonesClientBeginCreateOrUpdateOptions struct {
	// The ETag of the Private DNS zone. Omit this value to always overwrite the current zone. Specify the last-seen ETag value
// to prevent accidentally overwriting any concurrent changes.
	IfMatch *string

	// Set to '*' to allow a new Private DNS zone to be created, but to prevent updating an existing zone. Other values will be
// ignored.
	IfNoneMatch *string

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// PrivateZonesClientBeginDeleteOptions contains the optional parameters for the PrivateZonesClient.BeginDelete method.
type PrivateZonesClientBeginDeleteOptions struct {
	// The ETag of the Private DNS zone. Omit this value to always delete the current zone. Specify the last-seen ETag value to
// prevent accidentally deleting any concurrent changes.
	IfMatch *string

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// PrivateZonesClientBeginUpdateOptions contains the optional parameters for the PrivateZonesClient.BeginUpdate method.
type PrivateZonesClientBeginUpdateOptions struct {
	// The ETag of the Private DNS zone. Omit this value to always overwrite the current zone. Specify the last-seen ETag value
// to prevent accidentally overwriting any concurrent changes.
	IfMatch *string

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// PrivateZonesClientGetOptions contains the optional parameters for the PrivateZonesClient.Get method.
type PrivateZonesClientGetOptions struct {
	// placeholder for future optional parameters
}

// PrivateZonesClientListByResourceGroupOptions contains the optional parameters for the PrivateZonesClient.NewListByResourceGroupPager
// method.
type PrivateZonesClientListByResourceGroupOptions struct {
	// The maximum number of record sets to return. If not specified, returns up to 100 record sets.
	Top *int32
}

// PrivateZonesClientListOptions contains the optional parameters for the PrivateZonesClient.NewListPager method.
type PrivateZonesClientListOptions struct {
	// The maximum number of Private DNS zones to return. If not specified, returns up to 100 zones.
	Top *int32
}

// RecordSetsClientCreateOrUpdateOptions contains the optional parameters for the RecordSetsClient.CreateOrUpdate method.
type RecordSetsClientCreateOrUpdateOptions struct {
	// The ETag of the record set. Omit this value to always overwrite the current record set. Specify the last-seen ETag value
// to prevent accidentally overwriting any concurrent changes.
	IfMatch *string

	// Set to '*' to allow a new record set to be created, but to prevent updating an existing record set. Other values will be
// ignored.
	IfNoneMatch *string
}

// RecordSetsClientDeleteOptions contains the optional parameters for the RecordSetsClient.Delete method.
type RecordSetsClientDeleteOptions struct {
	// The ETag of the record set. Omit this value to always delete the current record set. Specify the last-seen ETag value to
// prevent accidentally deleting any concurrent changes.
	IfMatch *string
}

// RecordSetsClientGetOptions contains the optional parameters for the RecordSetsClient.Get method.
type RecordSetsClientGetOptions struct {
	// placeholder for future optional parameters
}

// RecordSetsClientListByTypeOptions contains the optional parameters for the RecordSetsClient.NewListByTypePager method.
type RecordSetsClientListByTypeOptions struct {
	// The suffix label of the record set name to be used to filter the record set enumeration. If this parameter is specified,
// the returned enumeration will only contain records that end with ".".
	Recordsetnamesuffix *string

	// The maximum number of record sets to return. If not specified, returns up to 100 record sets.
	Top *int32
}

// RecordSetsClientListOptions contains the optional parameters for the RecordSetsClient.NewListPager method.
type RecordSetsClientListOptions struct {
	// The suffix label of the record set name to be used to filter the record set enumeration. If this parameter is specified,
// the returned enumeration will only contain records that end with ".".
	Recordsetnamesuffix *string

	// The maximum number of record sets to return. If not specified, returns up to 100 record sets.
	Top *int32
}

// RecordSetsClientUpdateOptions contains the optional parameters for the RecordSetsClient.Update method.
type RecordSetsClientUpdateOptions struct {
	// The ETag of the record set. Omit this value to always overwrite the current record set. Specify the last-seen ETag value
// to prevent accidentally overwriting concurrent changes.
	IfMatch *string
}

// VirtualNetworkLinksClientBeginCreateOrUpdateOptions contains the optional parameters for the VirtualNetworkLinksClient.BeginCreateOrUpdate
// method.
type VirtualNetworkLinksClientBeginCreateOrUpdateOptions struct {
	// The ETag of the virtual network link to the Private DNS zone. Omit this value to always overwrite the current virtual network
// link. Specify the last-seen ETag value to prevent accidentally overwriting
// any concurrent changes.
	IfMatch *string

	// Set to '*' to allow a new virtual network link to the Private DNS zone to be created, but to prevent updating an existing
// link. Other values will be ignored.
	IfNoneMatch *string

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VirtualNetworkLinksClientBeginDeleteOptions contains the optional parameters for the VirtualNetworkLinksClient.BeginDelete
// method.
type VirtualNetworkLinksClientBeginDeleteOptions struct {
	// The ETag of the virtual network link to the Private DNS zone. Omit this value to always delete the current zone. Specify
// the last-seen ETag value to prevent accidentally deleting any concurrent
// changes.
	IfMatch *string

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VirtualNetworkLinksClientBeginUpdateOptions contains the optional parameters for the VirtualNetworkLinksClient.BeginUpdate
// method.
type VirtualNetworkLinksClientBeginUpdateOptions struct {
	// The ETag of the virtual network link to the Private DNS zone. Omit this value to always overwrite the current virtual network
// link. Specify the last-seen ETag value to prevent accidentally overwriting
// any concurrent changes.
	IfMatch *string

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VirtualNetworkLinksClientGetOptions contains the optional parameters for the VirtualNetworkLinksClient.Get method.
type VirtualNetworkLinksClientGetOptions struct {
	// placeholder for future optional parameters
}

// VirtualNetworkLinksClientListOptions contains the optional parameters for the VirtualNetworkLinksClient.NewListPager method.
type VirtualNetworkLinksClientListOptions struct {
	// The maximum number of virtual network links to return. If not specified, returns up to 100 virtual network links.
	Top *int32
}

