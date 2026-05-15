// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated_blob

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// KeyInfo - Key information
type KeyInfo struct {
	// REQUIRED; The date-time the key expires in ISO 8601 UTC time
	Expiry *string `xml:"Expiry"`

	// REQUIRED; The date-time the key is active in ISO 8601 UTC time
	Start *string `xml:"Start"`

	// The delegated user tenant id in Azure AD
	DelegatedUserTenantID *string `xml:"DelegatedUserTid"`
}

// ListFileSystemsSegmentResponse - An enumeration of containers
type ListFileSystemsSegmentResponse struct {
	// REQUIRED
	FileSystemItems []*FileSystemItem `xml:"Containers>Container"`

	// REQUIRED
	ServiceEndpoint *string `xml:"ServiceEndpoint,attr"`
	Marker          *string `xml:"Marker"`
	MaxResults      *int32  `xml:"MaxResults"`
	NextMarker      *string `xml:"NextMarker"`
	Prefix          *string `xml:"Prefix"`
}

// FileSystemItem - An Azure Storage container
type FileSystemItem struct {
	// REQUIRED
	Name *string `xml:"Name"`

	// REQUIRED; Properties of a container
	Properties *FileSystemProperties `xml:"Properties"`
	Deleted    *bool                 `xml:"Deleted"`

	// Dictionary of
	Metadata map[string]*string `xml:"Metadata"`
	Version  *string            `xml:"Version"`
}

// FileSystemProperties - Properties of a container
type FileSystemProperties struct {
	// REQUIRED
	ETag *azcore.ETag `xml:"Etag"`

	// REQUIRED
	LastModified           *time.Time `xml:"Last-Modified"`
	DefaultEncryptionScope *string    `xml:"DefaultEncryptionScope"`
	DeletedTime            *time.Time `xml:"DeletedTime"`
	HasImmutabilityPolicy  *bool      `xml:"HasImmutabilityPolicy"`
	HasLegalHold           *bool      `xml:"HasLegalHold"`

	// Indicates if version level worm is enabled on this container.
	IsImmutableStorageWithVersioningEnabled *bool              `xml:"ImmutableStorageWithVersioningEnabled"`
	LeaseDuration                           *LeaseDurationType `xml:"LeaseDuration"`
	LeaseState                              *LeaseStateType    `xml:"LeaseState"`
	LeaseStatus                             *LeaseStatusType   `xml:"LeaseStatus"`
	PreventEncryptionScopeOverride          *bool              `xml:"DenyEncryptionScopeOverride"`
	PublicAccess                            *PublicAccessType  `xml:"PublicAccess"`
	RemainingRetentionDays                  *int32             `xml:"RemainingRetentionDays"`
}

// UserDelegationKey - A user delegation key
type UserDelegationKey struct {
	// REQUIRED; The date-time the key expires
	SignedExpiry *time.Time `xml:"SignedExpiry"`

	// REQUIRED; The Azure Active Directory object ID in GUID format.
	SignedOID *string `xml:"SignedOid"`

	// REQUIRED; Abbreviation of the Azure Storage service that accepts the key
	SignedService *string `xml:"SignedService"`

	// REQUIRED; The date-time the key is active
	SignedStart *time.Time `xml:"SignedStart"`

	// REQUIRED; The Azure Active Directory tenant ID in GUID format
	SignedTID *string `xml:"SignedTid"`

	// REQUIRED; The service version that created the key
	SignedVersion *string `xml:"SignedVersion"`

	// REQUIRED; The key as a base64 string
	Value *string `xml:"Value"`

	// The delegated user tenant id in Azure AD. Return if DelegatedUserTid is specified.
	SignedDelegatedUserTenantID *string `xml:"SignedDelegatedUserTid"`
}
