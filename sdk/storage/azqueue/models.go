//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azqueue

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/sas"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	return exported.NewSharedKeyCredential(accountName, accountKey)
}

// URLParts object represents the components that make up an Azure Storage Queue URL.
// NOTE: Changing any SAS-related field requires computing a new SAS signature.
type URLParts = sas.URLParts

// ParseURL parses a URL initializing URLParts' fields including any SAS-related & snapshot query parameters. Any other
// query parameters remain in the UnparsedParams field. This method overwrites all fields in the URLParts object.
func ParseURL(u string) (URLParts, error) {
	return sas.ParseURL(u)
}

// ================================================================

// CorsRule - CORS is an HTTP feature that enables a web application running under one domain to access resources in another
// domain. Web browsers implement a security restriction known as same-origin policy that
// prevents a web page from calling APIs in a different domain; CORS provides a secure way to allow one domain (the origin
// domain) to call APIs in another domain
type CorsRule = generated.CorsRule

// GeoReplication - Geo-Replication information for the Secondary Storage Service
type GeoReplication = generated.GeoReplication

// RetentionPolicy - the retention policy which determines how long the associated data should persist
type RetentionPolicy = generated.RetentionPolicy

// Metrics - a summary of request statistics grouped by API in hour or minute aggregates for queues
type Metrics = generated.Metrics

// Logging - Azure Analytics Logging settings.
type Logging = generated.Logging

// TODO: CreateQueueOptions = queue.CreateOptions
// TODO: DeleteQueueOptions = queue.DeleteOptions

// StorageServiceProperties - Storage Service Properties.
type StorageServiceProperties = generated.StorageServiceProperties

// StorageServiceStats - Stats for the storage service.
type StorageServiceStats = generated.StorageServiceStats

// ---------------------------------------------------------------------------------------------------------------------

// ListQueuesOptions provides set of configurations for ListQueues operation
type ListQueuesOptions struct {
	Include ListQueuesInclude

	// A string value that identifies the portion of the list of queues to be returned with the next listing operation. The
	// operation returns the NextMarker value within the response body if the listing operation did not return all queues
	// remaining to be listed with the current page. The NextMarker value can be used as the value for the marker parameter in
	// a subsequent call to request the next page of list items. The marker value is opaque to the client.
	Marker *string

	// Specifies the maximum number of queues to return. If the request does not specify max results, or specifies a value
	// greater than 5000, the server will return up to 5000 items. Note that if the listing operation crosses a partition boundary,
	// then the service will return a continuation token for retrieving the remainder of the results. For this reason, it is possible
	// that the service will return fewer results than specified by max results, or than the default of 5000.
	MaxResults *int32

	// Filters the results to return only queues whose name begins with the specified prefix.
	Prefix *string
}

// ListQueuesInclude indicates what additional information the service should return with each queue.
type ListQueuesInclude struct {
	// Tells the service whether to return metadata for each queue.
	Metadata bool
}

// ---------------------------------------------------------------------------------------------------------------------

// SetPropertiesOptions provides set of options for ServiceClient.SetProperties
type SetPropertiesOptions struct {
	// The set of CORS rules.
	Cors []*CorsRule

	// a summary of request statistics grouped by API in hour or minute aggregates for queues
	HourMetrics *Metrics

	// Azure Analytics Logging settings.
	Logging *Logging

	// a summary of request statistics grouped by API in hour or minute aggregates for queues
	MinuteMetrics *Metrics
}

func (o *SetPropertiesOptions) format() (generated.StorageServiceProperties, *generated.ServiceClientSetPropertiesOptions) {
	if o == nil {
		return generated.StorageServiceProperties{}, nil
	}

	return generated.StorageServiceProperties{
		Cors:          o.Cors,
		HourMetrics:   o.HourMetrics,
		Logging:       o.Logging,
		MinuteMetrics: o.MinuteMetrics,
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------

// GetPropertiesOptions contains the optional parameters for the ServiceClient.GetProperties method.
type GetPropertiesOptions struct {
	// placeholder for future options
}

func (o *GetPropertiesOptions) format() *generated.ServiceClientGetPropertiesOptions {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// GetStatisticsOptions provides set of options for ServiceClient.GetStatistics
type GetStatisticsOptions struct {
	// placeholder for future options
}

func (o *GetStatisticsOptions) format() *generated.ServiceClientGetStatisticsOptions {
	return nil
}
