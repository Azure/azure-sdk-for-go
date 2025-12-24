// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package sas

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
)

// SignatureValues is used to generate a Shared Access Signature (SAS) for an Azure Storage file or share.
// For more information on creating service sas, see https://docs.microsoft.com/rest/api/storageservices/constructing-a-service-sas
// User Delegation SAS not supported for files service
type SignatureValues struct {
	Version                     string    `param:"sv"`  // If not specified, this defaults to Version
	Protocol                    Protocol  `param:"spr"` // See the Protocol* constants
	StartTime                   time.Time `param:"st"`  // Not specified if IsZero
	ExpiryTime                  time.Time `param:"se"`  // Not specified if IsZero
	SnapshotTime                time.Time
	Permissions                 string  `param:"sp"` // Create by initializing SharePermissions or FilePermissions and then call String()
	IPRange                     IPRange `param:"sip"`
	Identifier                  string  `param:"si"`
	ShareName                   string
	FilePath                    string // Ex: "directory/FileName". Use "" to create a Share SAS and file path for File SAS.
	CacheControl                string // rscc
	ContentDisposition          string // rscd
	ContentEncoding             string // rsce
	ContentLanguage             string // rscl
	ContentType                 string // rsct
	AuthorizedObjectID          string // saoid
	UnauthorizedObjectID        string // suoid
	CorrelationID               string // scid
	SignedDelegatedUserObjectID string // sduoid
}

// SignWithSharedKey uses an account's SharedKeyCredential to sign this signature values to produce the proper SAS query parameters.
func (v SignatureValues) SignWithSharedKey(sharedKeyCredential *SharedKeyCredential) (QueryParameters, error) {
	if v.Identifier == "" && (v.ExpiryTime.IsZero() || v.Permissions == "") {
		return QueryParameters{}, errors.New("service SAS is missing at least one of these: ExpiryTime or Permissions")
	}

	resource := "s"
	if v.FilePath == "" {
		// Make sure the permission characters are in the correct order
		perms, err := parseSharePermissions(v.Permissions)
		if err != nil {
			return QueryParameters{}, err
		}
		v.Permissions = perms.String()
	} else {
		resource = "f"
		// Make sure the permission characters are in the correct order
		perms, err := parseFilePermissions(v.Permissions)
		if err != nil {
			return QueryParameters{}, err
		}
		v.Permissions = perms.String()
	}

	if v.Version == "" {
		v.Version = Version
	}
	startTime, expiryTime, _ := formatTimesForSigning(v.StartTime, v.ExpiryTime, v.SnapshotTime)

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
		v.CacheControl,       // rscc
		v.ContentDisposition, // rscd
		v.ContentEncoding,    // rsce
		v.ContentLanguage,    // rscl
		v.ContentType},       // rsct
		"\n")

	signature, err := exported.ComputeHMACSHA256(sharedKeyCredential, stringToSign)
	if err != nil {
		return QueryParameters{}, err
	}

	p := QueryParameters{
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
		shareSnapshotTime:  v.SnapshotTime,
		// Calculated SAS signature
		signature: signature,
	}

	return p, nil
}

// SignWithUserDelegation uses an account's UserDelegationCredential to sign this signature values to produce the proper SAS query parameters.
func (v SignatureValues) SignWithUserDelegation(userDelegationCredential *UserDelegationCredential) (QueryParameters, error) {
	if userDelegationCredential == nil {
		return QueryParameters{}, fmt.Errorf("cannot sign SAS query without User Delegation Key")
	}

	if v.ExpiryTime.IsZero() || v.Permissions == "" {
		return QueryParameters{}, errors.New("user delegation SAS is missing at least one of these: ExpiryTime or Permissions")
	}

	resource := "s"
	if v.FilePath == "" {
		perms, err := parseSharePermissions(v.Permissions)
		if err != nil {
			return QueryParameters{}, err
		}
		v.Permissions = perms.String()
	} else {
		resource = "f"
		perms, err := parseFilePermissions(v.Permissions)
		if err != nil {
			return QueryParameters{}, err
		}
		v.Permissions = perms.String()
	}

	if v.Version == "" {
		v.Version = Version
	}
	startTime, expiryTime, _ := formatTimesForSigning(v.StartTime, v.ExpiryTime, v.SnapshotTime)

	udk := exported.GetUDKParams(userDelegationCredential)
	udkStart, udkExpiry, _ := formatTimesForSigning(*udk.SignedStart, *udk.SignedExpiry, time.Time{})

	stringToSign := strings.Join([]string{
		v.Permissions,
		startTime,
		expiryTime,
		getCanonicalName(exported.GetAccountName(userDelegationCredential), v.ShareName, v.FilePath),
		*udk.SignedOid,
		*udk.SignedTid,
		udkStart,
		udkExpiry,
		*udk.SignedService,
		*udk.SignedVersion,
		v.AuthorizedObjectID,
		v.UnauthorizedObjectID,
		v.CorrelationID,
		"",
		v.SignedDelegatedUserObjectID,
		v.IPRange.String(),
		string(v.Protocol),
		v.Version,
		resource,
		"",
		"",
		v.CacheControl,
		v.ContentDisposition,
		v.ContentEncoding,
		v.ContentLanguage,
		v.ContentType},
		"\n")

	signature, err := exported.ComputeUDCHMACSHA256(userDelegationCredential, stringToSign)
	if err != nil {
		return QueryParameters{}, err
	}

	p := QueryParameters{
		version:     v.Version,
		protocol:    v.Protocol,
		startTime:   v.StartTime,
		expiryTime:  v.ExpiryTime,
		permissions: v.Permissions,
		ipRange:     v.IPRange,

		resource:                    resource,
		identifier:                  v.Identifier,
		cacheControl:                v.CacheControl,
		contentDisposition:          v.ContentDisposition,
		contentEncoding:             v.ContentEncoding,
		contentLanguage:             v.ContentLanguage,
		contentType:                 v.ContentType,
		shareSnapshotTime:           v.SnapshotTime,
		authorizedObjectID:          v.AuthorizedObjectID,
		unauthorizedObjectID:        v.UnauthorizedObjectID,
		correlationID:               v.CorrelationID,
		signedDelegatedUserObjectID: v.SignedDelegatedUserObjectID,
		signature:                   signature,
	}

	p.signedOID = *udk.SignedOid
	p.signedTID = *udk.SignedTid
	p.signedStart = *udk.SignedStart
	p.signedExpiry = *udk.SignedExpiry
	p.signedService = *udk.SignedService
	p.signedVersion = *udk.SignedVersion

	return p, nil
}

// getCanonicalName computes the canonical name for a share or file resource for SAS signing.
func getCanonicalName(account string, shareName string, filePath string) string {
	// Share: "/file/account/sharename"
	// File:  "/file/account/sharename/filename"
	// File:  "/file/account/sharename/directoryname/filename"
	elements := []string{"/file/", account, "/", shareName}
	if filePath != "" {
		dfp := strings.ReplaceAll(filePath, "\\", "/")
		if dfp[0] == '/' {
			dfp = dfp[1:]
		}
		elements = append(elements, "/", dfp)
	}
	return strings.Join(elements, "")
}

// SharePermissions type simplifies creating the permissions string for an Azure Storage share SAS.
// Initialize an instance of this type and then call its String method to set SignatureValues' Permissions field.
// All permissions descriptions can be found here: https://docs.microsoft.com/en-us/rest/api/storageservices/create-service-sas#permissions-for-a-share
type SharePermissions struct {
	Read, Create, Write, Delete, List bool
}

// String produces the SAS permissions string for an Azure Storage share.
// Call this method to set SignatureValues' Permissions field.
func (p *SharePermissions) String() string {
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

// parseSharePermissions initializes SharePermissions' fields from a string.
func parseSharePermissions(s string) (SharePermissions, error) {
	p := SharePermissions{} // Clear the flags
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
			return SharePermissions{}, fmt.Errorf("invalid permission: '%v'", r)
		}
	}
	return p, nil
}

// FilePermissions type simplifies creating the permissions string for an Azure Storage file SAS.
// Initialize an instance of this type and then call its String method to set SignatureValues' Permissions field.
// All permissions descriptions can be found here: https://docs.microsoft.com/en-us/rest/api/storageservices/create-service-sas#permissions-for-a-file
type FilePermissions struct {
	Read, Create, Write, Delete bool
}

// String produces the SAS permissions string for an Azure Storage file.
// Call this method to set SignatureValues' Permissions field.
func (p *FilePermissions) String() string {
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

// parseFilePermissions initializes the FilePermissions' fields from a string.
func parseFilePermissions(s string) (FilePermissions, error) {
	p := FilePermissions{} // Clear the flags
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
			return FilePermissions{}, fmt.Errorf("invalid permission: '%v'", r)
		}
	}
	return p, nil
}
