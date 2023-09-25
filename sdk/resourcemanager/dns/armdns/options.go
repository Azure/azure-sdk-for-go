//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdns

// RecordSetsClientCreateOrUpdateOptions contains the optional parameters for the RecordSetsClient.CreateOrUpdate method.
type RecordSetsClientCreateOrUpdateOptions struct {
	// The etag of the record set. Omit this value to always overwrite the current record set. Specify the last-seen etag value
// to prevent accidentally overwriting any concurrent changes.
	IfMatch *string

	// Set to '*' to allow a new record set to be created, but to prevent updating an existing record set. Other values will be
// ignored.
	IfNoneMatch *string
}

// RecordSetsClientDeleteOptions contains the optional parameters for the RecordSetsClient.Delete method.
type RecordSetsClientDeleteOptions struct {
	// The etag of the record set. Omit this value to always delete the current record set. Specify the last-seen etag value to
// prevent accidentally deleting any concurrent changes.
	IfMatch *string
}

// RecordSetsClientGetOptions contains the optional parameters for the RecordSetsClient.Get method.
type RecordSetsClientGetOptions struct {
	// placeholder for future optional parameters
}

// RecordSetsClientListAllByDNSZoneOptions contains the optional parameters for the RecordSetsClient.NewListAllByDNSZonePager
// method.
type RecordSetsClientListAllByDNSZoneOptions struct {
	// The suffix label of the record set name that has to be used to filter the record set enumerations. If this parameter is
// specified, Enumeration will return only records that end with .
	RecordSetNameSuffix *string

	// The maximum number of record sets to return. If not specified, returns up to 100 record sets.
	Top *int32
}

// RecordSetsClientListByDNSZoneOptions contains the optional parameters for the RecordSetsClient.NewListByDNSZonePager method.
type RecordSetsClientListByDNSZoneOptions struct {
	// The suffix label of the record set name that has to be used to filter the record set enumerations. If this parameter is
// specified, Enumeration will return only records that end with .
	Recordsetnamesuffix *string

	// The maximum number of record sets to return. If not specified, returns up to 100 record sets.
	Top *int32
}

// RecordSetsClientListByTypeOptions contains the optional parameters for the RecordSetsClient.NewListByTypePager method.
type RecordSetsClientListByTypeOptions struct {
	// The suffix label of the record set name that has to be used to filter the record set enumerations. If this parameter is
// specified, Enumeration will return only records that end with .
	Recordsetnamesuffix *string

	// The maximum number of record sets to return. If not specified, returns up to 100 record sets.
	Top *int32
}

// RecordSetsClientUpdateOptions contains the optional parameters for the RecordSetsClient.Update method.
type RecordSetsClientUpdateOptions struct {
	// The etag of the record set. Omit this value to always overwrite the current record set. Specify the last-seen etag value
// to prevent accidentally overwriting concurrent changes.
	IfMatch *string
}

// ResourceReferenceClientGetByTargetResourcesOptions contains the optional parameters for the ResourceReferenceClient.GetByTargetResources
// method.
type ResourceReferenceClientGetByTargetResourcesOptions struct {
	// placeholder for future optional parameters
}

// ZonesClientBeginDeleteOptions contains the optional parameters for the ZonesClient.BeginDelete method.
type ZonesClientBeginDeleteOptions struct {
	// The etag of the DNS zone. Omit this value to always delete the current zone. Specify the last-seen etag value to prevent
// accidentally deleting any concurrent changes.
	IfMatch *string

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ZonesClientCreateOrUpdateOptions contains the optional parameters for the ZonesClient.CreateOrUpdate method.
type ZonesClientCreateOrUpdateOptions struct {
	// The etag of the DNS zone. Omit this value to always overwrite the current zone. Specify the last-seen etag value to prevent
// accidentally overwriting any concurrent changes.
	IfMatch *string

	// Set to '*' to allow a new DNS zone to be created, but to prevent updating an existing zone. Other values will be ignored.
	IfNoneMatch *string
}

// ZonesClientGetOptions contains the optional parameters for the ZonesClient.Get method.
type ZonesClientGetOptions struct {
	// placeholder for future optional parameters
}

// ZonesClientListByResourceGroupOptions contains the optional parameters for the ZonesClient.NewListByResourceGroupPager
// method.
type ZonesClientListByResourceGroupOptions struct {
	// The maximum number of record sets to return. If not specified, returns up to 100 record sets.
	Top *int32
}

// ZonesClientListOptions contains the optional parameters for the ZonesClient.NewListPager method.
type ZonesClientListOptions struct {
	// The maximum number of DNS zones to return. If not specified, returns up to 100 zones.
	Top *int32
}

// ZonesClientUpdateOptions contains the optional parameters for the ZonesClient.Update method.
type ZonesClientUpdateOptions struct {
	// The etag of the DNS zone. Omit this value to always overwrite the current zone. Specify the last-seen etag value to prevent
// accidentally overwriting any concurrent changes.
	IfMatch *string
}

