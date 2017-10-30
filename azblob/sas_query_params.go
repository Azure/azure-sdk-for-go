package azblob

import (
	"net"
	"net/url"
	"strings"
	"time"
)

// SASTimeFormat represents the format of a SAS start or expiry time. Use it when formatting/parsing a time.Time.
const SASTimeFormat = "2006-01-02T15:04:05Z" //"2017-07-27T00:00:00Z" // ISO 8601

// https://docs.microsoft.com/en-us/rest/api/storageservices/constructing-a-service-sas

// A SASQueryParameters object represents the components that make up an Azure Storage SAS' query parameters.
// You parse a map of query parameters into its fields by calling NewSASQueryParameters(). You add the components
// to a query parameter map by calling AddToValues().
// NOTE: Changing any field requires computing a new SAS signature using a XxxSASSignatureValues type.
//
// This type defines the components used by all Azure Storage resources (Containers, Blobs, Files, & Queues).
type SASQueryParameters struct {
	// All members are immutable or values so copies of this struct are goroutine-safe.
	Version       string    `param:"sv"`
	Services      string    `param:"ss"`
	ResourceTypes string    `param:"srt"`
	Protocol      string    `param:"spr"`
	StartTime     time.Time `param:"st"`
	ExpiryTime    time.Time `param:"se"`
	IPRange       IPRange   `param:"sip"`
	Identifier    string    `param:"si"`
	Resource      string    `param:"sr"`
	Permissions   string    `param:"sp"`
	Signature     string    `param:"sig"`
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

// NewSASQueryParameters creates and initializes a SASQueryParameters object based on the
// query parameter map passed in values. Upon return, all SAS-related query parameters are
// removed from the map.
func NewSASQueryParameters(values url.Values) SASQueryParameters {
	p := SASQueryParameters{}
	for k, v := range values {
		val := v[0]
		isSASKey := true
		switch k {
		case "sv":
			p.Version = val
		case "ss":
			p.Services = val
		case "srt":
			p.ResourceTypes = val
		case "spr":
			p.Protocol = val
		case "st":
			p.StartTime, _ = time.Parse(SASTimeFormat, val)
		case "se":
			p.ExpiryTime, _ = time.Parse(SASTimeFormat, val)
		case "sip":
			dashIndex := strings.Index(val, "-")
			if dashIndex == -1 {
				p.IPRange.Start = net.ParseIP(val)
			} else {
				p.IPRange.Start = net.ParseIP(val[:dashIndex])
				p.IPRange.End = net.ParseIP(val[dashIndex+1:])
			}
		case "si":
			p.Identifier = val
		case "sr":
			p.Resource = val
		case "sp":
			p.Permissions = val
		case "sig":
			p.Signature = val
		default:
			isSASKey = false // We didn't recognize the query parameter
		}
		if isSASKey {
			delete(values, k)
		}
	}
	return p
}

// AddToValues adds the SAS components to the specified query parameters map.
func (p *SASQueryParameters) AddToValues(v url.Values) url.Values {
	if p.Version != "" {
		v.Add("sv", p.Version)
	}
	if p.Services != "" {
		v.Add("ss", p.Services)
	}
	if p.ResourceTypes != "" {
		v.Add("srt", p.ResourceTypes)
	}
	if p.Protocol != "" {
		v.Add("spr", p.Protocol)
	}
	if !p.StartTime.IsZero() {
		v.Add("st", p.StartTime.Format(SASTimeFormat))
	}
	if !p.ExpiryTime.IsZero() {
		v.Add("se", p.ExpiryTime.Format(SASTimeFormat))
	}
	if len(p.IPRange.Start) > 0 {
		v.Add("sip", p.IPRange.String())
	}
	if p.Identifier != "" {
		v.Add("si", p.Identifier)
	}
	if p.Resource != "" {
		v.Add("sr", p.Resource)
	}
	if p.Permissions != "" {
		v.Add("sp", p.Permissions)
	}
	if p.Signature != "" {
		v.Add("sig", p.Signature)
	}
	return v
}

// Encode encodes the SAS query parameters into URL encoded form sorted by key.
func (p *SASQueryParameters) Encode() string {
	v := url.Values{}
	p.AddToValues(v)
	return v.Encode()
}
