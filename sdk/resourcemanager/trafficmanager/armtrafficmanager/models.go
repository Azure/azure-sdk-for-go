//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armtrafficmanager

import "time"

// CheckTrafficManagerRelativeDNSNameAvailabilityParameters - Parameters supplied to check Traffic Manager name operation.
type CheckTrafficManagerRelativeDNSNameAvailabilityParameters struct {
	// The name of the resource.
	Name *string

	// The type of the resource.
	Type *string
}

// DNSConfig - Class containing DNS settings in a Traffic Manager profile.
type DNSConfig struct {
	// The relative DNS name provided by this Traffic Manager profile. This value is combined with the DNS domain name used by
// Azure Traffic Manager to form the fully-qualified domain name (FQDN) of the
// profile.
	RelativeName *string

	// The DNS Time-To-Live (TTL), in seconds. This informs the local DNS resolvers and DNS clients how long to cache DNS responses
// provided by this Traffic Manager profile.
	TTL *int64

	// READ-ONLY; The fully-qualified domain name (FQDN) of the Traffic Manager profile. This is formed from the concatenation
// of the RelativeName with the DNS domain used by Azure Traffic Manager.
	Fqdn *string
}

// DeleteOperationResult - The result of the request or operation.
type DeleteOperationResult struct {
	// READ-ONLY; The result of the operation or request.
	OperationResult *bool
}

// Endpoint - Class representing a Traffic Manager endpoint.
type Endpoint struct {
	// Fully qualified resource Id for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficManagerProfiles/{resourceName}
	ID *string

	// The name of the resource
	Name *string

	// The properties of the Traffic Manager endpoint.
	Properties *EndpointProperties

	// The type of the resource. Ex- Microsoft.Network/trafficManagerProfiles.
	Type *string
}

// EndpointProperties - Class representing a Traffic Manager endpoint properties.
type EndpointProperties struct {
	// If Always Serve is enabled, probing for endpoint health will be disabled and endpoints will be included in the traffic
// routing method.
	AlwaysServe *AlwaysServe

	// List of custom headers.
	CustomHeaders []*EndpointPropertiesCustomHeadersItem

	// Specifies the location of the external or nested endpoints when using the 'Performance' traffic routing method.
	EndpointLocation *string

	// The monitoring status of the endpoint.
	EndpointMonitorStatus *EndpointMonitorStatus

	// The status of the endpoint. If the endpoint is Enabled, it is probed for endpoint health and is included in the traffic
// routing method.
	EndpointStatus *EndpointStatus

	// The list of countries/regions mapped to this endpoint when using the 'Geographic' traffic routing method. Please consult
// Traffic Manager Geographic documentation for a full list of accepted values.
	GeoMapping []*string

	// The minimum number of endpoints that must be available in the child profile in order for the parent profile to be considered
// available. Only applicable to endpoint of type 'NestedEndpoints'.
	MinChildEndpoints *int64

	// The minimum number of IPv4 (DNS record type A) endpoints that must be available in the child profile in order for the parent
// profile to be considered available. Only applicable to endpoint of type
// 'NestedEndpoints'.
	MinChildEndpointsIPv4 *int64

	// The minimum number of IPv6 (DNS record type AAAA) endpoints that must be available in the child profile in order for the
// parent profile to be considered available. Only applicable to endpoint of type
// 'NestedEndpoints'.
	MinChildEndpointsIPv6 *int64

	// The priority of this endpoint when using the 'Priority' traffic routing method. Possible values are from 1 to 1000, lower
// values represent higher priority. This is an optional parameter. If specified,
// it must be specified on all endpoints, and no two endpoints can share the same priority value.
	Priority *int64

	// The list of subnets, IP addresses, and/or address ranges mapped to this endpoint when using the 'Subnet' traffic routing
// method. An empty list will match all ranges not covered by other endpoints.
	Subnets []*EndpointPropertiesSubnetsItem

	// The fully-qualified DNS name or IP address of the endpoint. Traffic Manager returns this value in DNS responses to direct
// traffic to this endpoint.
	Target *string

	// The Azure Resource URI of the of the endpoint. Not applicable to endpoints of type 'ExternalEndpoints'.
	TargetResourceID *string

	// The weight of this endpoint when using the 'Weighted' traffic routing method. Possible values are from 1 to 1000.
	Weight *int64
}

// EndpointPropertiesCustomHeadersItem - Custom header name and value.
type EndpointPropertiesCustomHeadersItem struct {
	// Header name.
	Name *string

	// Header value.
	Value *string
}

// EndpointPropertiesSubnetsItem - Subnet first address, scope, and/or last address.
type EndpointPropertiesSubnetsItem struct {
	// First address in the subnet.
	First *string

	// Last address in the subnet.
	Last *string

	// Block size (number of leading bits in the subnet mask).
	Scope *int32
}

// GeographicHierarchy - Class representing the Geographic hierarchy used with the Geographic traffic routing method.
type GeographicHierarchy struct {
	// Fully qualified resource Id for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficManagerProfiles/{resourceName}
	ID *string

	// The name of the resource
	Name *string

	// The properties of the Geographic Hierarchy resource.
	Properties *GeographicHierarchyProperties

	// The type of the resource. Ex- Microsoft.Network/trafficManagerProfiles.
	Type *string
}

// GeographicHierarchyProperties - Class representing the properties of the Geographic hierarchy used with the Geographic
// traffic routing method.
type GeographicHierarchyProperties struct {
	// The region at the root of the hierarchy from all the regions in the hierarchy can be retrieved.
	GeographicHierarchy *Region
}

// HeatMapEndpoint - Class which is a sparse representation of a Traffic Manager endpoint.
type HeatMapEndpoint struct {
	// A number uniquely identifying this endpoint in query experiences.
	EndpointID *int32

	// The ARM Resource ID of this Traffic Manager endpoint.
	ResourceID *string
}

// HeatMapModel - Class representing a Traffic Manager HeatMap.
type HeatMapModel struct {
	// Fully qualified resource Id for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficManagerProfiles/{resourceName}
	ID *string

	// The name of the resource
	Name *string

	// The properties of the Traffic Manager HeatMap.
	Properties *HeatMapProperties

	// The type of the resource. Ex- Microsoft.Network/trafficManagerProfiles.
	Type *string
}

// HeatMapProperties - Class representing a Traffic Manager HeatMap properties.
type HeatMapProperties struct {
	// The ending of the time window for this HeatMap, exclusive.
	EndTime *time.Time

	// The endpoints used in this HeatMap calculation.
	Endpoints []*HeatMapEndpoint

	// The beginning of the time window for this HeatMap, inclusive.
	StartTime *time.Time

	// The traffic flows produced in this HeatMap calculation.
	TrafficFlows []*TrafficFlow
}

// MonitorConfig - Class containing endpoint monitoring settings in a Traffic Manager profile.
type MonitorConfig struct {
	// List of custom headers.
	CustomHeaders []*MonitorConfigCustomHeadersItem

	// List of expected status code ranges.
	ExpectedStatusCodeRanges []*MonitorConfigExpectedStatusCodeRangesItem

	// The monitor interval for endpoints in this profile. This is the interval at which Traffic Manager will check the health
// of each endpoint in this profile.
	IntervalInSeconds *int64

	// The path relative to the endpoint domain name used to probe for endpoint health.
	Path *string

	// The TCP port used to probe for endpoint health.
	Port *int64

	// The profile-level monitoring status of the Traffic Manager profile.
	ProfileMonitorStatus *ProfileMonitorStatus

	// The protocol (HTTP, HTTPS or TCP) used to probe for endpoint health.
	Protocol *MonitorProtocol

	// The monitor timeout for endpoints in this profile. This is the time that Traffic Manager allows endpoints in this profile
// to response to the health check.
	TimeoutInSeconds *int64

	// The number of consecutive failed health check that Traffic Manager tolerates before declaring an endpoint in this profile
// Degraded after the next failed health check.
	ToleratedNumberOfFailures *int64
}

// MonitorConfigCustomHeadersItem - Custom header name and value.
type MonitorConfigCustomHeadersItem struct {
	// Header name.
	Name *string

	// Header value.
	Value *string
}

// MonitorConfigExpectedStatusCodeRangesItem - Min and max value of a status code range.
type MonitorConfigExpectedStatusCodeRangesItem struct {
	// Max status code.
	Max *int32

	// Min status code.
	Min *int32
}

// NameAvailability - Class representing a Traffic Manager Name Availability response.
type NameAvailability struct {
	// Descriptive message that explains why the name is not available, when applicable.
	Message *string

	// The relative name.
	Name *string

	// Describes whether the relative name is available or not.
	NameAvailable *bool

	// The reason why the name is not available, when applicable.
	Reason *string

	// Traffic Manager profile resource type.
	Type *string
}

// Profile - Class representing a Traffic Manager profile.
type Profile struct {
	// Fully qualified resource Id for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficManagerProfiles/{resourceName}
	ID *string

	// The Azure Region where the resource lives
	Location *string

	// The name of the resource
	Name *string

	// The properties of the Traffic Manager profile.
	Properties *ProfileProperties

	// Resource tags.
	Tags map[string]*string

	// The type of the resource. Ex- Microsoft.Network/trafficManagerProfiles.
	Type *string
}

// ProfileListResult - The list Traffic Manager profiles operation response.
type ProfileListResult struct {
	// Gets the list of Traffic manager profiles.
	Value []*Profile
}

// ProfileProperties - Class representing the Traffic Manager profile properties.
type ProfileProperties struct {
	// The list of allowed endpoint record types.
	AllowedEndpointRecordTypes []*AllowedEndpointRecordType

	// The DNS settings of the Traffic Manager profile.
	DNSConfig *DNSConfig

	// The list of endpoints in the Traffic Manager profile.
	Endpoints []*Endpoint

	// Maximum number of endpoints to be returned for MultiValue routing type.
	MaxReturn *int64

	// The endpoint monitoring settings of the Traffic Manager profile.
	MonitorConfig *MonitorConfig

	// The status of the Traffic Manager profile.
	ProfileStatus *ProfileStatus

	// The traffic routing method of the Traffic Manager profile.
	TrafficRoutingMethod *TrafficRoutingMethod

	// Indicates whether Traffic View is 'Enabled' or 'Disabled' for the Traffic Manager profile. Null, indicates 'Disabled'.
// Enabling this feature will increase the cost of the Traffic Manage profile.
	TrafficViewEnrollmentStatus *TrafficViewEnrollmentStatus
}

// ProxyResource - The resource model definition for a ARM proxy resource. It will have everything other than required location
// and tags
type ProxyResource struct {
	// Fully qualified resource Id for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficManagerProfiles/{resourceName}
	ID *string

	// The name of the resource
	Name *string

	// The type of the resource. Ex- Microsoft.Network/trafficManagerProfiles.
	Type *string
}

// QueryExperience - Class representing a Traffic Manager HeatMap query experience properties.
type QueryExperience struct {
	// REQUIRED; The id of the endpoint from the 'endpoints' array which these queries were routed to.
	EndpointID *int32

	// REQUIRED; The number of queries originating from this location.
	QueryCount *int32

	// The latency experienced by queries originating from this location.
	Latency *float64
}

// Region - Class representing a region in the Geographic hierarchy used with the Geographic traffic routing method.
type Region struct {
	// The code of the region
	Code *string

	// The name of the region
	Name *string

	// The list of Regions grouped under this Region in the Geographic Hierarchy.
	Regions []*Region
}

// Resource - The core properties of ARM resources
type Resource struct {
	// Fully qualified resource Id for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficManagerProfiles/{resourceName}
	ID *string

	// The name of the resource
	Name *string

	// The type of the resource. Ex- Microsoft.Network/trafficManagerProfiles.
	Type *string
}

// TrackedResource - The resource model definition for a ARM tracked top level resource
type TrackedResource struct {
	// Fully qualified resource Id for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficManagerProfiles/{resourceName}
	ID *string

	// The Azure Region where the resource lives
	Location *string

	// The name of the resource
	Name *string

	// Resource tags.
	Tags map[string]*string

	// The type of the resource. Ex- Microsoft.Network/trafficManagerProfiles.
	Type *string
}

// TrafficFlow - Class representing a Traffic Manager HeatMap traffic flow properties.
type TrafficFlow struct {
	// The approximate latitude that these queries originated from.
	Latitude *float64

	// The approximate longitude that these queries originated from.
	Longitude *float64

	// The query experiences produced in this HeatMap calculation.
	QueryExperiences []*QueryExperience

	// The IP address that this query experience originated from.
	SourceIP *string
}

// UserMetricsModel - Class representing Traffic Manager User Metrics.
type UserMetricsModel struct {
	// Fully qualified resource Id for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficManagerProfiles/{resourceName}
	ID *string

	// The name of the resource
	Name *string

	// The properties of the Traffic Manager User Metrics.
	Properties *UserMetricsProperties

	// The type of the resource. Ex- Microsoft.Network/trafficManagerProfiles.
	Type *string
}

// UserMetricsProperties - Class representing a Traffic Manager Real User Metrics key response.
type UserMetricsProperties struct {
	// The key returned by the User Metrics operation.
	Key *string
}

