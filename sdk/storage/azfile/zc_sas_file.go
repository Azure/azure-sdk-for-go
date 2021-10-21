// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

// FileSASSignatureValues is used to generate a Shared Access Signature (SAS) for an Azure Storage share or file.
type FileSASSignatureValues struct {
	Version            string      `param:"sv"`  // If not specified, this defaults to SASVersion
	Protocol           SASProtocol `param:"spr"` // See the SASProtocol* constants
	StartTime          time.Time   `param:"st"`  // Not specified if IsZero
	ExpiryTime         time.Time   `param:"se"`  // Not specified if IsZero
	SnapshotTime       time.Time
	Permissions        string  `param:"sp"` // Create by initializing a ShareSASPermissions or FileSASPermissions and then call String()
	IPRange            IPRange `param:"sip"`
	Identifier         string  `param:"si"`
	ShareName          string
	FilePath           string // Ex: "directory/FileName" or "FileName". Use "" to create a Share SAS.
	CacheControl       string // rscc
	ContentDisposition string // rscd
	ContentEncoding    string // rsce
	ContentLanguage    string // rscl
	ContentType        string // rsct
}

// NewSASQueryParameters uses an account's shared key credential to sign this signature values to produce
// the proper SAS query parameters.
func (v FileSASSignatureValues) NewSASQueryParameters(sharedKeyCredential *SharedKeyCredential) (SASQueryParameters, error) {
	if sharedKeyCredential == nil {
		return SASQueryParameters{}, errors.New("sharedKeyCredential can't be nil")
	}

	resource := "s"
	if v.FilePath == "" {
		// Make sure the permission characters are in the correct order
		perms := &ShareSASPermissions{}
		if err := perms.Parse(v.Permissions); err != nil {
			return SASQueryParameters{}, err
		}
		v.Permissions = perms.String()
	} else {
		resource = "f"
		// Make sure the permission characters are in the correct order
		perms := &FileSASPermissions{}
		if err := perms.Parse(v.Permissions); err != nil {
			return SASQueryParameters{}, err
		}
		v.Permissions = perms.String()
	}
	if v.Version == "" {
		v.Version = SASVersion
	}
	startTime, expiryTime, snapshotTime := FormatTimesForSASSigning(v.StartTime, v.ExpiryTime, v.SnapshotTime)

	p := SASQueryParameters{
		// Common SAS parameters
		version:     v.Version,
		protocol:    v.Protocol,
		startTime:   v.StartTime,
		expiryTime:  v.ExpiryTime,
		permissions: v.Permissions,
		ipRange:     v.IPRange,

		// Share/File-specific SAS parameters
		resource:           resource,
		identifier:         v.Identifier,
		cacheControl:       v.CacheControl,
		contentDisposition: v.ContentDisposition,
		contentEncoding:    v.ContentEncoding,
		contentLanguage:    v.ContentLanguage,
		contentType:        v.ContentType,
		snapshotTime:       v.SnapshotTime,
	}

	// String to sign: http://msdn.microsoft.com/en-us/library/azure/dn140255.aspx
	stringToSign := strings.Join([]string{
		v.Permissions,
		startTime,
		expiryTime,
		getCanonicalName(sharedKeyCredential.AccountName(), v.ShareName, v.FilePath),
		v.Identifier,
		v.IPRange.String(),
		string(v.Protocol),
		v.Version,
		resource,
		snapshotTime,
		v.CacheControl,       // rscc
		v.ContentDisposition, // rscd
		v.ContentEncoding,    // rsce
		v.ContentLanguage,    // rscl
		v.ContentType},       // rsct
		"\n")

	signature, err := sharedKeyCredential.ComputeHMACSHA256(stringToSign)
	p.signature = signature
	return p, err
}

// getCanonicalName computes the canonical name for a share or file resource for SAS signing.
func getCanonicalName(account string, shareName string, filePath string) string {
	// Share: "/file/account/sharename"
	// File:  "/file/account/sharename/filename"
	// File:  "/file/account/sharename/directoryname/filename"
	elements := []string{"/file/", account, "/", shareName}
	if filePath != "" {
		dfp := strings.Replace(filePath, "\\", "/", -1)
		if dfp[0] == '/' {
			dfp = dfp[1:]
		}
		elements = append(elements, "/", dfp)
	}
	return strings.Join(elements, "")
}

// The ShareSASPermissions type simplifies creating the permissions string for an Azure Storage share SAS.
// Initialize an instance of this type and then call its String method to set FileSASSignatureValues's Permissions field.
type ShareSASPermissions struct {
	Read, Create, Write, Delete, List bool
}

// String produces the SAS permissions string for an Azure Storage share.
// Call this method to set FileSASSignatureValues's Permissions field.
func (p ShareSASPermissions) String() string {
	var b bytes.Buffer
	if p.Read {
		b.WriteRune('r')
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

// Parse initializes the ShareSASPermissions' fields from a string.
func (p *ShareSASPermissions) Parse(s string) error {
	*p = ShareSASPermissions{} // Clear the flags
	for _, r := range s {
		switch r {
		case 'r':
			p.Read = true
		case 'c':
			p.Create = true
		case 'w':
			p.Write = true
		case 'd':
			p.Delete = true
		case 'l':
			p.List = true
		default:
			return fmt.Errorf("invalid permission: '%v'", r)
		}
	}
	return nil
}

// The FileSASPermissions type simplifies creating the permissions string for an Azure Storage file SAS.
// Initialize an instance of this type and then call its String method to set FileSASSignatureValues's Permissions field.
type FileSASPermissions struct{ Read, Create, Write, Delete bool }

// String produces the SAS permissions string for an Azure Storage file.
// Call this method to set FileSASSignatureValues' Permissions field.
func (p FileSASPermissions) String() string {
	var b bytes.Buffer
	if p.Read {
		b.WriteRune('r')
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

// Parse initializes the FileSASPermissions' fields from a string.
func (p *FileSASPermissions) Parse(s string) error {
	*p = FileSASPermissions{} // Clear the flags
	for _, r := range s {
		switch r {
		case 'r':
			p.Read = true
		case 'c':
			p.Create = true
		case 'w':
			p.Write = true
		case 'd':
			p.Delete = true
		default:
			return fmt.Errorf("invalid permission: '%v'", r)
		}
	}
	return nil
}
