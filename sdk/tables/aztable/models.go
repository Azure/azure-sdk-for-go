// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
)

// ListOptions contains a group of parameters for the Table.Query method.
type ListOptions struct {
	// OData filter expression.
	Filter *string
	// Specifies the media type for the response.
	Format *ODataMetadataFormat
	// Select expression using OData notation. Limits the columns on each record to just those requested, e.g. "$select=PolicyAssignmentId, ResourceId".
	Select *string
	// Maximum number of records to return.
	Top *int32
}

func (l *ListOptions) toQueryOptions() *generated.QueryOptions {
	if l == nil {
		return &generated.QueryOptions{}
	}

	return &generated.QueryOptions{
		Filter: l.Filter,
		// Format: l.Format,  // TODO: Fix
		Select: l.Select,
		Top:    l.Top,
	}
}

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

// TableServiceProperties - Table Service Properties.
type TableServiceProperties struct {
	// The set of CORS rules.
	Cors []*CorsRule `xml:"Cors>CorsRule"`

	// A summary of request statistics grouped by API in hourly aggregates for tables.
	HourMetrics *generated.Metrics `xml:"HourMetrics"`

	// Azure Analytics Logging settings.
	Logging *Logging `xml:"Logging"`

	// A summary of request statistics grouped by API in minute aggregates for tables.
	MinuteMetrics *generated.Metrics `xml:"MinuteMetrics"`
}

func (t *TableServiceProperties) toGenerated() *generated.TableServiceProperties {
	if t == nil {
		return &generated.TableServiceProperties{}
	}

	return &generated.TableServiceProperties{
		Cors:          toGeneratedCorsRules(t.Cors),
		HourMetrics:   t.HourMetrics,
		Logging:       toGeneratedLogging(t.Logging),
		MinuteMetrics: t.MinuteMetrics,
	}
}

func toGeneratedCorsRules(corsRules []*CorsRule) []*generated.CorsRule {
	var ret []*generated.CorsRule
	for _, c := range corsRules {
		ret = append(ret, c.toGenerated())
	}
	return ret
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
