//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pageblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

type CreateOptions struct {
	// Set for page blobs only. The sequence number is a user-controlled value that you can use to track requests. The value of
	// the sequence number must be between 0 and 2^63 - 1.
	SequenceNumber *int64

	// Optional. Used to set blob tags in various blob operations.
	Tags map[string]string

	// Optional. Specifies a user-defined name-value pair associated with the blob. If no name-value pairs are specified, the
	// operation will copy the metadata from the source blob or file to the destination blob. If one or more name-value pairs
	// are specified, the destination blob is created with the specified metadata, and metadata is not copied from the source
	// blob or file. Note that beginning with version 2009-09-19, metadata names must adhere to the naming rules for C# identifiers.
	// See Naming and Referencing Containers, Blobs, and Metadata for more information.
	Metadata map[string]string

	// Optional. Indicates the tier to be set on the page blob.
	Tier *PremiumPageBlobAccessTier

	HTTPHeaders *HTTPHeaders

	CpkInfo *CpkInfo

	CpkScopeInfo *CpkScopeInfo

	AccessConditions *AccessConditions
	// Specifies the date time when the blobs immutability policy is set to expire.
	ImmutabilityPolicyExpiry *time.Time
	// Specifies the immutability policy mode to set on the blob.
	ImmutabilityPolicyMode *ImmutabilityPolicyMode
	// Specified if a legal hold should be set on the blob.
	LegalHold *bool
}

type AccessTier = generated.AccessTier

type PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTier

const (
	PremiumPageBlobAccessTierP10 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP10
	PremiumPageBlobAccessTierP15 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP15
	PremiumPageBlobAccessTierP20 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP20
	PremiumPageBlobAccessTierP30 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP30
	PremiumPageBlobAccessTierP4  PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP4
	PremiumPageBlobAccessTierP40 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP40
	PremiumPageBlobAccessTierP50 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP50
	PremiumPageBlobAccessTierP6  PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP6
	PremiumPageBlobAccessTierP60 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP60
	PremiumPageBlobAccessTierP70 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP70
	PremiumPageBlobAccessTierP80 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP80
)

// ImmutabilityPolicyMode enum
type ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyMode

const (
	ImmutabilityPolicyModeMutable  ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyModeMutable
	ImmutabilityPolicyModeUnlocked ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyModeUnlocked
	ImmutabilityPolicyModeLocked   ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyModeLocked
)

// PossibleBlobImmutabilityPolicyModeValues returns the possible values for the BlobImmutabilityPolicyMode const type.
func PossibleBlobImmutabilityPolicyModeValues() []ImmutabilityPolicyMode {
	return []ImmutabilityPolicyMode{
		ImmutabilityPolicyModeMutable,
		ImmutabilityPolicyModeUnlocked,
		ImmutabilityPolicyModeLocked,
	}
}

type CpkInfo = generated.CpkInfo

type CpkScopeInfo = generated.CpkScopeInfo

type AccessConditions = exported.BlobAccessConditions

type HTTPHeaders = generated.BlobHTTPHeaders

// Redeclared Options

type DownloadOptions = blob.DownloadOptions

type DeleteOptions = blob.DeleteOptions

type UndeleteOptions = blob.UndeleteOptions

type SetTierOptions = blob.SetTierOptions

type GetPropertiesOptions = blob.GetPropertiesOptions

type SetHTTPHeadersOptions = blob.SetHTTPHeadersOptions

type SetMetadataOptions = blob.SetMetadataOptions

type CreateSnapshotOptions = blob.CreateSnapshotOptions

type StartCopyFromURLOptions = blob.StartCopyFromURLOptions

type AbortCopyFromURLOptions = blob.AbortCopyFromURLOptions

type SetTagsOptions = blob.SetTagsOptions

type GetTagsOptions = blob.GetTagsOptions

type CopyFromURLOptions = blob.CopyFromURLOptions
