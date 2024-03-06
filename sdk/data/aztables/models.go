// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	generated "github.com/Azure/azure-sdk-for-go/sdk/data/aztables/internal"
)

// CorsRule - CORS is an HTTP feature that enables a web application running under one domain to access resources in another domain. Web browsers implement
// a security restriction known as same-origin policy that
// prevents a web page from calling APIs in a different domain; CORS provides a secure way to allow one domain (the origin domain) to call APIs in another
// domain.
type CorsRule struct {
	// REQUIRED; The request headers that the origin domain may specify on the CORS request.
	AllowedHeaders *string `xml:"AllowedHeaders"`

	// REQUIRED; The methods (HTTP request verbs) that the origin domain may use for a CORS request. (comma separated)
	AllowedMethods *string `xml:"AllowedMethods"`

	// REQUIRED; The origin domains that are permitted to make a request against the service via CORS. The origin domain is the domain from which the request
	// originates. Note that the origin must be an exact
	// case-sensitive match with the origin that the user age sends to the service. You can also use the wildcard character '*' to allow all origin domains
	// to make requests via CORS.
	AllowedOrigins *string `xml:"AllowedOrigins"`

	// REQUIRED; The response headers that may be sent in the response to the CORS request and exposed by the browser to the request issuer.
	ExposedHeaders *string `xml:"ExposedHeaders"`

	// REQUIRED; The maximum amount time that a browser should cache the preflight OPTIONS request.
	MaxAgeInSeconds *int32 `xml:"MaxAgeInSeconds"`
}

func (c *CorsRule) toGenerated() *generated.CorsRule {
	if c == nil {
		return nil
	}

	return &generated.CorsRule{
		AllowedHeaders:  c.AllowedHeaders,
		AllowedMethods:  c.AllowedMethods,
		AllowedOrigins:  c.AllowedOrigins,
		ExposedHeaders:  c.ExposedHeaders,
		MaxAgeInSeconds: c.MaxAgeInSeconds,
	}
}

func fromGeneratedCors(c *generated.CorsRule) *CorsRule {
	if c == nil {
		return nil
	}

	return &CorsRule{
		AllowedHeaders:  c.AllowedHeaders,
		AllowedMethods:  c.AllowedMethods,
		AllowedOrigins:  c.AllowedOrigins,
		ExposedHeaders:  c.ExposedHeaders,
		MaxAgeInSeconds: c.MaxAgeInSeconds,
	}
}

func toGeneratedCorsRules(corsRules []*CorsRule) []*generated.CorsRule {
	if len(corsRules) == 0 {
		return nil
	}
	ret := make([]*generated.CorsRule, len(corsRules))
	for i := range corsRules {
		ret[i] = corsRules[i].toGenerated()
	}
	return ret
}

// ServiceProperties - Service Properties for a given table
type ServiceProperties struct {
	// The set of CORS rules.
	Cors []*CorsRule `xml:"Cors>CorsRule"`

	// A summary of request statistics grouped by API in hourly aggregates for tables.
	HourMetrics *Metrics `xml:"HourMetrics"`

	// Azure Analytics Logging settings.
	Logging *Logging `xml:"Logging"`

	// A summary of request statistics grouped by API in minute aggregates for tables.
	MinuteMetrics *Metrics `xml:"MinuteMetrics"`
}

func (t *ServiceProperties) toGenerated() *generated.TableServiceProperties {
	if t == nil {
		return &generated.TableServiceProperties{}
	}

	return &generated.TableServiceProperties{
		Cors:          toGeneratedCorsRules(t.Cors),
		HourMetrics:   toGeneratedMetrics(t.HourMetrics),
		Logging:       toGeneratedLogging(t.Logging),
		MinuteMetrics: toGeneratedMetrics(t.MinuteMetrics),
	}
}

// TableProperties contains the properties for a single Table
type TableProperties struct {
	// The name of the table.
	Name *string `json:"TableName,omitempty"`

	// The OData properties of the table in JSON format.
	Value []byte
}

// RetentionPolicy - The retention policy.
type RetentionPolicy struct {
	// REQUIRED; Indicates whether a retention policy is enabled for the service.
	Enabled *bool `xml:"Enabled"`

	// Indicates the number of days that metrics or logging or soft-deleted data should be retained. All data older than this value will be deleted.
	Days *int32 `xml:"Days"`
}

func toGeneratedRetentionPolicy(r *RetentionPolicy) *generated.RetentionPolicy {
	if r == nil {
		return &generated.RetentionPolicy{}
	}

	return &generated.RetentionPolicy{
		Enabled: r.Enabled,
		Days:    r.Days,
	}
}

func fromGeneratedRetentionPolicy(r *generated.RetentionPolicy) *RetentionPolicy {
	if r == nil {
		return &RetentionPolicy{}
	}

	return &RetentionPolicy{
		Enabled: r.Enabled,
		Days:    r.Days,
	}
}

// Logging - Azure Analytics Logging settings.
type Logging struct {
	// REQUIRED; Indicates whether all delete requests should be logged.
	Delete *bool `xml:"Delete"`

	// REQUIRED; Indicates whether all read requests should be logged.
	Read *bool `xml:"Read"`

	// REQUIRED; The retention policy.
	RetentionPolicy *RetentionPolicy `xml:"RetentionPolicy"`

	// REQUIRED; The version of Analytics to configure.
	Version *string `xml:"Version"`

	// REQUIRED; Indicates whether all write requests should be logged.
	Write *bool `xml:"Write"`
}

func toGeneratedLogging(l *Logging) *generated.Logging {
	if l == nil {
		return nil
	}

	return &generated.Logging{
		Delete:          l.Delete,
		Read:            l.Read,
		RetentionPolicy: toGeneratedRetentionPolicy(l.RetentionPolicy),
		Version:         l.Version,
		Write:           l.Write,
	}
}

func fromGeneratedLogging(g *generated.Logging) *Logging {
	if g == nil {
		return nil
	}

	return &Logging{
		Delete:          g.Delete,
		Read:            g.Read,
		Write:           g.Write,
		Version:         g.Version,
		RetentionPolicy: (*RetentionPolicy)(g.RetentionPolicy),
	}
}

// Metrics are the metrics for a Table
type Metrics struct {
	// REQUIRED; Indicates whether metrics are enabled for the Table service.
	Enabled *bool `xml:"Enabled"`

	// Indicates whether metrics should generate summary statistics for called API operations.
	IncludeAPIs *bool `xml:"IncludeAPIs"`

	// The retention policy.
	RetentionPolicy *RetentionPolicy `xml:"RetentionPolicy"`

	// The version of Analytics to configure.
	Version *string `xml:"Version"`
}

func toGeneratedMetrics(m *Metrics) *generated.Metrics {
	if m == nil {
		return nil
	}

	return &generated.Metrics{
		Enabled:         m.Enabled,
		IncludeAPIs:     m.IncludeAPIs,
		Version:         m.Version,
		RetentionPolicy: toGeneratedRetentionPolicy(m.RetentionPolicy),
	}
}

func fromGeneratedMetrics(m *generated.Metrics) *Metrics {
	if m == nil {
		return &Metrics{}
	}

	return &Metrics{
		Enabled:         m.Enabled,
		IncludeAPIs:     m.IncludeAPIs,
		Version:         m.Version,
		RetentionPolicy: fromGeneratedRetentionPolicy(m.RetentionPolicy),
	}
}

// SignedIdentifier - A signed identifier.
type SignedIdentifier struct {
	// REQUIRED; The access policy.
	AccessPolicy *AccessPolicy `xml:"AccessPolicy"`

	// REQUIRED; A unique id.
	ID *string `xml:"Id"`
}

func toGeneratedSignedIdentifier(s *SignedIdentifier) *generated.SignedIdentifier {
	if s == nil {
		return nil
	}

	return &generated.SignedIdentifier{
		ID:           s.ID,
		AccessPolicy: toGeneratedAccessPolicy(s.AccessPolicy),
	}
}

func fromGeneratedSignedIdentifier(s *generated.SignedIdentifier) *SignedIdentifier {
	if s == nil {
		return nil
	}

	return &SignedIdentifier{
		ID:           s.ID,
		AccessPolicy: fromGeneratedAccessPolicy(s.AccessPolicy),
	}

}

// AccessPolicy - An Access policy.
type AccessPolicy struct {
	// REQUIRED; The datetime that the policy expires.
	Expiry *time.Time `xml:"Expiry"`

	// REQUIRED; The permissions for the acl policy.
	Permission *string `xml:"Permission"`

	// REQUIRED; The datetime from which the policy is active.
	Start *time.Time `xml:"Start"`
}

func toGeneratedAccessPolicy(a *AccessPolicy) *generated.AccessPolicy {
	if a == nil {
		return nil
	}

	expiry := a.Expiry
	if expiry != nil {
		expiry = to.Ptr(expiry.UTC())
	}

	start := a.Start
	if start != nil {
		start = to.Ptr(start.UTC())
	}

	return &generated.AccessPolicy{
		Expiry:     expiry,
		Permission: a.Permission,
		Start:      start,
	}
}

func fromGeneratedAccessPolicy(g *generated.AccessPolicy) *AccessPolicy {
	if g == nil {
		return nil
	}

	return &AccessPolicy{
		Expiry:     g.Expiry,
		Permission: g.Permission,
		Start:      g.Start,
	}
}

// GeoReplication represents the GeoReplication status of an account
type GeoReplication struct {
	// REQUIRED; A GMT date/time value, to the second. All primary writes preceding this value are guaranteed to be available for read operations at the secondary.
	// Primary writes after this point in time may or may
	// not be available for reads.
	LastSyncTime *time.Time `xml:"LastSyncTime"`

	// REQUIRED; The status of the secondary location.
	Status *GeoReplicationStatus `xml:"Status"`
}

func fromGeneratedGeoReplication(g *generated.GeoReplication) *GeoReplication {
	if g == nil {
		return nil
	}

	return &GeoReplication{
		LastSyncTime: g.LastSyncTime,
		Status:       toGeneratedStatusType(g.Status),
	}
}
