// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"net"
	"net/url"
	"time"
)

// FormatTimesForSASSigning converts a time.Time to a snapshotTimeFormat string suitable for a
// SASField's StartTime or ExpiryTime fields. Returns "" if value.IsZero().
func FormatTimesForSASSigning(startTime, expiryTime time.Time) (string, string) {
	ss := ""
	if !startTime.IsZero() {
		ss = formatSASTimeWithDefaultFormat(&startTime)
	}
	se := ""
	if !expiryTime.IsZero() {
		se = formatSASTimeWithDefaultFormat(&expiryTime)
	}
	return ss, se
}

// sasTimeFormat represents the format of a SAS start or expiry time. Use it when formatting/parsing a time.Time.
const sasTimeFormat = "2006-01-02T15:04:05Z" // "2017-07-27T00:00:00Z" // ISO 8601

// formatSASTimeWithDefaultFormat format time with ISO 8601 in "yyyy-MM-ddTHH:mm:ssZ".
func formatSASTimeWithDefaultFormat(t *time.Time) string {
	return formatSASTime(t, sasTimeFormat) // By default, "yyyy-MM-ddTHH:mm:ssZ" is used
}

// formatSASTime format time with given format, use ISO 8601 in "yyyy-MM-ddTHH:mm:ssZ" by default.
func formatSASTime(t *time.Time, format string) string {
	if format != "" {
		return t.Format(format)
	}
	return t.Format(sasTimeFormat) // By default, "yyyy-MM-ddTHH:mm:ssZ" is used
}

// https://learn.microsoft.com/rest/api/storageservices/constructing-a-service-sas

// SASQueryParameters represents the components that make up an Azure Storage SAS' query parameters.
// You parse a map of query parameters into its fields by calling Sign(). You add the components
// to a query parameter map by calling AddToValues().
// NOTE: Changing any field requires computing a new SAS signature using a XxxSASSignatureValues type.
type SASQueryParameters struct {
	// All members are immutable or values so copies of this struct are goroutine-safe.
	version       string      `param:"sv"`
	services      string      `param:"ss"`
	resourceTypes string      `param:"srt"`
	protocol      SASProtocol `param:"spr"`
	startTime     time.Time   `param:"st"`
	expiryTime    time.Time   `param:"se"`
	ipRange       IPRange     `param:"sip"`
	identifier    string      `param:"si"`
	resource      string      `param:"sr"`
	permissions   string      `param:"sp"`
	signature     string      `param:"sig"`
	signedVersion string      `param:"skv"`
	tableName     string      `param:"tn"`
	startPk       string      `param:"spk"`
	startRk       string      `param:"srk"`
	endPk         string      `param:"epk"`
	endRk         string      `param:"erk"`

	// private member used for startTime and expiryTime formatting.
	stTimeFormat string
	seTimeFormat string
}

// SignedVersion gets the signed version of a SASQueryParameter
func (p *SASQueryParameters) SignedVersion() string {
	return p.signedVersion
}

// Version gets the version of a SASQueryParameter
func (p *SASQueryParameters) Version() string {
	return p.version
}

// Services gets the services of a SASQueryParameter
func (p *SASQueryParameters) Services() string {
	return p.services
}

// ResourceTypes gets the resource types of a SASQueryParameter
func (p *SASQueryParameters) ResourceTypes() string {
	return p.resourceTypes
}

// Protocol gets the protocol of a SASQueryParameter
func (p *SASQueryParameters) Protocol() SASProtocol {
	return p.protocol
}

// StartTime gets the start time of a SASQueryParameter
func (p *SASQueryParameters) StartTime() time.Time {
	return p.startTime
}

// ExpiryTime gets the expiry time of a SASQueryParameter
func (p *SASQueryParameters) ExpiryTime() time.Time {
	return p.expiryTime
}

// IPRange gets the IP Range of a SASQueryParameter
func (p *SASQueryParameters) IPRange() IPRange {
	return p.ipRange
}

// Identifier gets the identifier of a SASQueryParameter
func (p *SASQueryParameters) Identifier() string {
	return p.identifier
}

// Resource gets the resource of a SASQueryParameter
func (p *SASQueryParameters) Resource() string {
	return p.resource
}

// Permissions gets the permissions of a SASQueryParameter
func (p *SASQueryParameters) Permissions() string {
	return p.permissions
}

// Signature gets the signature of a SASQueryParameter
func (p *SASQueryParameters) Signature() string {
	return p.signature
}

// StartPartitionKey gets the start PartitionKey of a SASQueryParameter
func (p *SASQueryParameters) StartPartitionKey() string {
	return p.startPk
}

// StartRowKey gets the start RowKey of a SASQueryParameter
func (p *SASQueryParameters) StartRowKey() string {
	return p.startRk
}

// EndPartitionKey gets the end PartitionKey of a SASQueryParameter
func (p *SASQueryParameters) EndPartitionKey() string {
	return p.endPk
}

// EndRowKey gets the end RowKey of a SASQueryParameter
func (p *SASQueryParameters) EndRowKey() string {
	return p.endRk
}

// IPRange represents a SAS IP range's start IP and (optionally) end IP.
type IPRange struct {
	Start net.IP // Not specified if length = 0
	End   net.IP // Not specified if length = 0
}

// String returns a string representation of an IPRange.
func (ipr *IPRange) String() string {
	if len(ipr.Start) == 0 {
		return ""
	}
	start := ipr.Start.String()
	if len(ipr.End) == 0 {
		return start
	}
	return start + "-" + ipr.End.String()
}

// addToValues adds the SAS components to the specified query parameters map.
func (p *SASQueryParameters) addToValues(v url.Values) url.Values {
	if p.version != "" {
		v.Add("sv", p.version)
	}
	if p.services != "" {
		v.Add("ss", p.services)
	}
	if p.resourceTypes != "" {
		v.Add("srt", p.resourceTypes)
	}
	if p.protocol != "" {
		v.Add("spr", string(p.protocol))
	}
	if !p.startTime.IsZero() {
		v.Add("st", formatSASTime(&(p.startTime), p.stTimeFormat))
	}
	if !p.expiryTime.IsZero() {
		v.Add("se", formatSASTime(&(p.expiryTime), p.seTimeFormat))
	}
	if len(p.ipRange.Start) > 0 {
		v.Add("sip", p.ipRange.String())
	}
	if p.identifier != "" {
		v.Add("si", p.identifier)
	}
	if p.resource != "" {
		v.Add("sr", p.resource)
	}
	if p.permissions != "" {
		v.Add("sp", p.permissions)
	}
	if p.signature != "" {
		v.Add("sig", p.signature)
	}
	if p.tableName != "" {
		v.Add("tn", p.tableName)
	}
	if p.startPk != "" {
		v.Add("spk", p.startPk)
	}
	if p.endPk != "" {
		v.Add("epk", p.endPk)
	}
	if p.startRk != "" {
		v.Add("srk", p.startRk)
	}
	if p.endRk != "" {
		v.Add("erk", p.endRk)
	}
	return v
}

// Encode encodes the SAS query parameters into URL encoded form sorted by key.
func (p *SASQueryParameters) Encode() string {
	v := url.Values{}
	p.addToValues(v)
	return v.Encode()
}
