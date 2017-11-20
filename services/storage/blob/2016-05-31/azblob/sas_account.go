package azblob

import (
	"bytes"
	"strings"
	"time"
)

// SASVersion indicates the SAS version.
const SASVersion = "2015-04-05"

const (
	// SASProtocolHTTPS can be specified for a SAS protocol
	SASProtocolHTTPS = "https"

	// SASProtocolHTTPSandHTTP can be specified for a SAS protocol
	SASProtocolHTTPSandHTTP = "https,http"
)

// FormatTimesForSASSigning converts a time.Time to a snapshotTimeFormat string suitable for a
// SASField's StartTime or ExpiryTime fields. Returns "" if value.IsZero().
func FormatTimesForSASSigning(startTime, expiryTime time.Time) (string, string) {
	ss := ""
	if !startTime.IsZero() {
		ss = startTime.Format(SASTimeFormat) // "yyyy-MM-ddTHH:mm:ssZ"
	}
	se := ""
	if !expiryTime.IsZero() {
		se = expiryTime.Format(SASTimeFormat) // "yyyy-MM-ddTHH:mm:ssZ"
	}
	return ss, se
}

// AccountSASSignatureValues is used to generate a Shared Access Signature (SAS) for an Azure Storage account.
type AccountSASSignatureValues struct {
	Version       string    `param:"sv"`  // If not specified, this defaults to azstorage.SASVersion
	Protocol      string    `param:"spr"` // See the SASProtocol* constants
	StartTime     time.Time `param:"st"`  // Not specified if IsZero
	ExpiryTime    time.Time `param:"se"`  // Not specified if IsZero
	Permissions   string    `param:"sp"`
	IPRange       IPRange   `param:"sip"`
	Services      string    `param:"ss"`
	ResourceTypes string    `param:"srt"`
}

// NewSASQueryParameters uses an account's shared key credential to sign this signature values to produce
// the proper SAS query parameters.
func (v AccountSASSignatureValues) NewSASQueryParameters(sharedKeyCredential *SharedKeyCredential) SASQueryParameters {
	// https://docs.microsoft.com/en-us/rest/api/storageservices/Constructing-an-Account-SAS
	if v.ExpiryTime.IsZero() || v.Permissions == "" || v.ResourceTypes == "" || v.Services == "" {
		panic("Account SAS is missing at least one of these: ExpiryTime, Permissions, Service, or ResourceType")
	}
	if v.Version == "" {
		v.Version = SASVersion
	}
	startTime, expiryTime := FormatTimesForSASSigning(v.StartTime, v.ExpiryTime)

	stringToSign := strings.Join([]string{
		sharedKeyCredential.AccountName(),
		v.Permissions,
		v.Services,
		v.ResourceTypes,
		startTime,
		expiryTime,
		v.IPRange.String(),
		v.Protocol,
		v.Version,
		""}, // That right, the account SAS requires a terminating extra newline
		"\n")

	signature := sharedKeyCredential.ComputeHMACSHA256(stringToSign)
	p := SASQueryParameters{
		// Common SAS parameters
		Version:     v.Version,
		Protocol:    v.Protocol,
		StartTime:   v.StartTime,
		ExpiryTime:  v.ExpiryTime,
		Permissions: v.Permissions,
		IPRange:     v.IPRange,

		// Account-specific SAS parameters
		Services:      v.Services,
		ResourceTypes: v.ResourceTypes,

		// Calculated SAS signature
		Signature: signature,
	}
	return p
}

// The AccountSASPermissions type simplifies creating the permissions string for an Azure Storage Account SAS.
// Initialize an instance of this type and then call its String method to set AccountSASSignatureValues's Permissions field.
type AccountSASPermissions struct {
	Read, Write, Delete, List, Add, Create, Update, Process bool
}

// String produces the SAS permissions string for an Azure Storage account.
// Call this method to set AccountSASSignatureValues's Permissions field.
func (p AccountSASPermissions) String() string {
	var buffer bytes.Buffer
	if p.Read {
		buffer.WriteRune('r')
	}
	if p.Write {
		buffer.WriteRune('w')
	}
	if p.Delete {
		buffer.WriteRune('d')
	}
	if p.List {
		buffer.WriteRune('l')
	}
	if p.Add {
		buffer.WriteRune('a')
	}
	if p.Create {
		buffer.WriteRune('c')
	}
	if p.Update {
		buffer.WriteRune('u')
	}
	if p.Process {
		buffer.WriteRune('p')
	}
	return buffer.String()
}

// The AccountSASServices type simplifies creating the services string for an Azure Storage Account SAS.
// Initialize an instance of this type and then call its String method to set AccountSASSignatureValues's Services field.
type AccountSASServices struct {
	Blob, Queue, File bool
}

// String produces the SAS services string for an Azure Storage account.
// Call this method to set AccountSASSignatureValues's Services field.
func (s AccountSASServices) String() string {
	var buffer bytes.Buffer
	if s.Blob {
		buffer.WriteRune('b')
	}
	if s.Queue {
		buffer.WriteRune('q')
	}
	if s.File {
		buffer.WriteRune('f')
	}
	return buffer.String()
}

// The AccountSASResourceTypes type simplifies creating the resource types string for an Azure Storage Account SAS.
// Initialize an instance of this type and then call its String method to set AccountSASSignatureValues's ResourceTypes field.
type AccountSASResourceTypes struct {
	Service, Container, Object bool
}

// String produces the SAS resource types string for an Azure Storage account.
// Call this method to set AccountSASSignatureValues's ResourceTypes field.
func (rt AccountSASResourceTypes) String() string {
	var buffer bytes.Buffer
	if rt.Service {
		buffer.WriteRune('s')
	}
	if rt.Container {
		buffer.WriteRune('c')
	}
	if rt.Object {
		buffer.WriteRune('o')
	}
	return buffer.String()
}
