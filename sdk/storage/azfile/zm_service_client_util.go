//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azfile

// ---------------------------------------------------------------------------------------------------------------------

type ListShareOptions struct {
}

// ---------------------------------------------------------------------------------------------------------------------

// ServiceGetPropertiesOptions provides set of options for ServiceClient.GetAccountInfo
type ServiceGetPropertiesOptions struct {
	// placeholder for future options
}

func (o *ServiceGetPropertiesOptions) format() *serviceClientGetPropertiesOptions {
	return nil
}

type ServiceGetPropertiesResponse struct {
	serviceClientGetPropertiesResponse
}

func toServiceGetPropertiesResponse(resp serviceClientGetPropertiesResponse) ServiceGetPropertiesResponse {
	return ServiceGetPropertiesResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ServiceSetPropertiesOptions struct {
	// The set of CORS rules.
	Cors []*CorsRule `xml:"Cors>CorsRule"`

	// A summary of request statistics grouped by API in hourly aggregates for files.
	HourMetrics *Metrics `xml:"HourMetrics"`

	// A summary of request statistics grouped by API in minute aggregates for files.
	MinuteMetrics *Metrics `xml:"MinuteMetrics"`

	// Protocol settings
	Protocol *ShareProtocolSettings `xml:"ProtocolSettings"`
}

func (o *ServiceSetPropertiesOptions) format() (StorageServiceProperties, *serviceClientSetPropertiesOptions) {
	if o == nil {
		return StorageServiceProperties{}, nil
	}

	return StorageServiceProperties{
		Cors:          o.Cors,
		HourMetrics:   o.HourMetrics,
		MinuteMetrics: o.MinuteMetrics,
		Protocol:      o.Protocol,
	}, nil
}

type ServiceSetPropertiesResponse struct {
	serviceClientSetPropertiesResponse
}

func toServiceSetPropertiesResponse(resp serviceClientSetPropertiesResponse) ServiceSetPropertiesResponse {
	return ServiceSetPropertiesResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ServiceListSharesOptions struct {
	// Include this parameter to specify one or more datasets to include in the response.
	Include []ListSharesIncludeType
	// A string value that identifies the portion of the list to be returned with the next list operation. The operation returns
	// a marker value within the response body if the list returned was not complete.
	// The marker value may then be used in a subsequent call to request the next set of list items. The marker value is opaque
	// to the client.
	Marker *string
	// Specifies the maximum number of entries to return. If the request does not specify maxresults, or specifies a value greater
	// than 5,000, the server will return up to 5,000 items.
	MaxResults *int32
	// Filters the results to return only entries whose name begins with the specified prefix.
	Prefix *string
}

func (o *ServiceListSharesOptions) format() *serviceClientListSharesSegmentOptions {
	if o == nil {
		return nil
	}
	return &serviceClientListSharesSegmentOptions{
		Include:    o.Include,
		Marker:     o.Marker,
		Maxresults: o.MaxResults,
		Prefix:     o.Prefix,
	}
}

type ServiceListSharesPager struct {
	*serviceClientListSharesSegmentPager
}

func toServiceListSharesPager(resp *serviceClientListSharesSegmentPager) *ServiceListSharesPager {
	return &ServiceListSharesPager{resp}
}
