//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azfile

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

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

// ShareMetricProperties defines convenience struct for Metrics,
type ShareMetricProperties struct {
	// Enabled - Indicates whether metrics are enabled for the File service.
	Enabled *bool

	// Version - The version of Storage Analytics to configure.
	// Version string, comment out version, as it's mandatory and should be 1.0
	// IncludeAPIs - Indicates whether metrics should generate summary statistics for called API operations.
	IncludeAPIs *bool

	// RetentionPolicyEnabled - Indicates whether a retention policy is enabled for the File service.
	RetentionPolicyEnabled *bool
	// RetentionDays - Indicates the number of days that metrics data should be retained.
	RetentionDays *int32
}

type ServiceSetPropertiesOptions struct {
	// The set of CORS rules.
	Cors []*ShareCorsRule

	// A summary of request statistics grouped by API in hourly aggregates for files.
	HourMetrics *ShareMetricProperties

	// A summary of request statistics grouped by API in minute aggregates for files.
	MinuteMetrics *ShareMetricProperties

	// Protocol settings
	Protocol *ShareProtocolSettings `xml:"ProtocolSettings"`
}

func (mp *ShareMetricProperties) toMetrics() *ShareMetrics {
	if mp == nil {
		return nil
	}

	metrics := ShareMetrics{
		Version: to.Ptr(StorageAnalyticsVersion),
	}

	if mp.Enabled != nil && *mp.Enabled {
		metrics.Enabled = mp.Enabled
		metrics.IncludeApis = mp.IncludeAPIs
		metrics.RetentionPolicy = &ShareRetentionPolicy{
			Enabled: mp.RetentionPolicyEnabled,
			Days:    mp.RetentionDays,
		}
	}

	return &metrics
}

func (o *ServiceSetPropertiesOptions) format() (ShareServiceProperties, *serviceClientSetPropertiesOptions) {
	if o == nil {
		return ShareServiceProperties{}, nil
	}

	return ShareServiceProperties{
		Cors:          o.Cors,
		HourMetrics:   o.HourMetrics.toMetrics(),
		MinuteMetrics: o.MinuteMetrics.toMetrics(),
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
	// Include this parameter to specify one or more datasets to include in the responseBody.
	Include []ListSharesIncludeType
	// A string value that identifies the portion of the list to be returned with the next list operation. The operation returns
	// a marker value within the responseBody body if the list returned was not complete.
	// The marker value may then be used in a subsequent call to request the next set of list items. The marker value is opaque
	// to the client.
	Marker *string
	// Specifies the maximum number of entries to return. If the request does not specify maxresults, or specifies a value greater
	// than 5,000, the server will return up to 5,000 items.
	MaxResults *int32
	// Filters the results to return only entries whose name begins with the specified prefix.
	Prefix *string
}

func (o *ServiceListSharesOptions) format() serviceClientListSharesSegmentOptions {
	if o == nil {
		return serviceClientListSharesSegmentOptions{}
	}
	return serviceClientListSharesSegmentOptions{
		Include:    o.Include,
		Marker:     o.Marker,
		Maxresults: o.MaxResults,
		Prefix:     o.Prefix,
	}
}

type ServiceListSharesResponse struct {
	serviceClientListSharesSegmentResponse
}

func toServiceListSharesResponse(resp serviceClientListSharesSegmentResponse) ServiceListSharesResponse {
	return ServiceListSharesResponse{resp}
}
