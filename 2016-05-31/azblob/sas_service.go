package azblob

import (
	"bytes"
	"strings"
	"time"
)

// BlobSASSignatureValues is used to generate a Shared Access Signature (SAS) for an Azure Storage container or blob.
type BlobSASSignatureValues struct {
	Version            string    `param:"sv"`  // If not specified, this defaults to SASVersion
	Protocol           string    `param:"spr"` // See the SASProtocol* constants
	StartTime          time.Time `param:"st"`  // Not specified if IsZero
	ExpiryTime         time.Time `param:"se"`  // Not specified if IsZero
	Permissions        string    `param:"sp"`
	IPRange            IPRange   `param:"sip"`
	ContainerName      string
	BlobName           string // Use "" to create a Container SAS
	Identifier         string `param:"si"`
	CacheControl       string // rscc
	ContentDisposition string // rscd
	ContentEncoding    string // rsce
	ContentLanguage    string // rscl
	ContentType        string // rsct
}

// NewSASQueryParameters uses an account's shared key credential to sign this signature values to produce
// the proper SAS query parameters.
func (v BlobSASSignatureValues) NewSASQueryParameters(sharedKeyCredential *SharedKeyCredential) SASQueryParameters {
	if sharedKeyCredential == nil {
		panic("sharedKeyCredential can't be nil")
	}

	resource := "c"
	if v.BlobName != "" {
		resource = "b"
	}
	if v.Version == "" {
		v.Version = SASVersion
	}
	startTime, expiryTime := FormatTimesForSASSigning(v.StartTime, v.ExpiryTime)

	// String to sign: http://msdn.microsoft.com/en-us/library/azure/dn140255.aspx
	stringToSign := strings.Join([]string{
		v.Permissions,
		startTime,
		expiryTime,
		getCanonicalName(sharedKeyCredential.AccountName(), v.ContainerName, v.BlobName),
		v.Identifier,
		v.IPRange.String(),
		v.Protocol,
		v.Version,
		v.CacheControl,       // rscc
		v.ContentDisposition, // rscd
		v.ContentEncoding,    // rsce
		v.ContentLanguage,    // rscl
		v.ContentType},       // rsct
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

		// Container/Blob-specific SAS parameters
		Resource:   resource,
		Identifier: v.Identifier,

		// Calculated SAS signature
		Signature: signature,
	}
	return p
}

// getCanonicalName computes the canonical name for a container or blob resource for SAS signing.
func getCanonicalName(account string, containerName string, blobName string) string {
	// Container: "/blob/account/containername"
	// Blob:      "/blob/account/containername/blobname"
	elems := []string{"/blob/", account, "/", containerName}
	if blobName != "" {
		elems = append(elems, "/", strings.Replace(blobName, "\\", "/", -1))
	}
	return strings.Join(elems, "")
}

// The ContainerSASPermissions type simplifies creating the permissions string for an Azure Storage container SAS.
// Initialize an instance of this type and then call its String method to set BlobSASSignatureValues's Permissions field.
type ContainerSASPermissions struct {
	Read, Add, Create, Write, Delete, List bool
}

// String produces the SAS permissions string for an Azure Storage container.
// Call this method to set BlobSASSignatureValues's Permissions field.
func (p ContainerSASPermissions) String() string {
	var b bytes.Buffer
	if p.Read {
		b.WriteRune('r')
	}
	if p.Add {
		b.WriteRune('a')
	}
	if p.Create {
		b.WriteRune('c')
	}
	if p.Write {
		b.WriteRune('w')
	}
	if p.Delete {
		b.WriteRune('d')
	}
	if p.List {
		b.WriteRune('l')
	}
	return b.String()
}

// Parse initializes the ContainerSASPermissions's fields from a string.
func (p *ContainerSASPermissions) Parse(s string) {
	p.Read = strings.ContainsRune(s, 'r')
	p.Add = strings.ContainsRune(s, 'a')
	p.Create = strings.ContainsRune(s, 'c')
	p.Write = strings.ContainsRune(s, 'w')
	p.Delete = strings.ContainsRune(s, 'd')
	p.List = strings.ContainsRune(s, 'l')
}

// The BlobSASPermissions type simplifies creating the permissions string for an Azure Storage blob SAS.
// Initialize an instance of this type and then call its String method to set BlobSASSignatureValues's Permissions field.
type BlobSASPermissions struct{ Read, Add, Create, Write, Delete bool }

// String produces the SAS permissions string for an Azure Storage blob.
// Call this method to set BlobSASSignatureValues's Permissions field.
func (p BlobSASPermissions) String() string {
	var b bytes.Buffer
	if p.Read {
		b.WriteRune('r')
	}
	if p.Add {
		b.WriteRune('a')
	}
	if p.Create {
		b.WriteRune('c')
	}
	if p.Write {
		b.WriteRune('w')
	}
	if p.Delete {
		b.WriteRune('d')
	}
	return b.String()
}

// Parse initializes the BlobSASPermissions's fields from a string.
func (p *BlobSASPermissions) Parse(s string) {
	p.Read = strings.ContainsRune(s, 'r')
	p.Add = strings.ContainsRune(s, 'a')
	p.Create = strings.ContainsRune(s, 'c')
	p.Write = strings.ContainsRune(s, 'w')
	p.Delete = strings.ContainsRune(s, 'd')
}
