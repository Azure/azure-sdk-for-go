//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"time"
)

// Type Declarations ---------------------------------------------------------------------

// AccessConditions identifies blob-specific access conditions which you optionally set.
type AccessConditions = exported.BlobAccessConditions

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = exported.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = exported.ModifiedAccessConditions

// AccessTier defines values for Blob Access Tier
type AccessTier = generated.AccessTier

// ClientOptions adds additional client options while constructing connection
type ClientOptions = exported.ClientOptions

// CpkInfo contains a group of parameters for client provided encryption key.
type CpkInfo = generated.CpkInfo

// CpkScopeInfo contains a group of parameters for client provided encryption scope.
type CpkScopeInfo = generated.CpkScopeInfo

// DeleteSnapshotsOptionType defines values for DeleteSnapshotsOptionType
type DeleteSnapshotsOptionType = generated.DeleteSnapshotsOptionType

// HTTPHeaders contains a group of parameters for the BlobClient.SetHTTPHeaders method.
type HTTPHeaders = generated.BlobHTTPHeaders

// ImmutabilityPolicyMode defines values for BlobImmutabilityPolicyMode
type ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyMode

// CopyStatusType defines values for CopyStatusType
type CopyStatusType = generated.CopyStatusType

// EncryptionAlgorithmType defines values for EncryptionAlgorithmType
type EncryptionAlgorithmType = generated.EncryptionAlgorithmType

// RehydratePriority - If an object is in rehydrate pending state then this header is returned with priority of rehydrate.
// Valid values are High and Standard.
type RehydratePriority = generated.RehydratePriority

// SASQueryParameters object represents the components that make up an Azure Storage SAS' query parameters.
// You parse a map of query parameters into its fields by calling NewSASQueryParameters(). You add the components
// to a query parameter map by calling AddToValues().
// NOTE: Changing any field requires computing a new SAS signature using a XxxSASSignatureValues type.
type SASQueryParameters = exported.SASQueryParameters

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// SourceModifiedAccessConditions contains a group of parameters for the BlobClient.StartCopyFromURL method.
type SourceModifiedAccessConditions = generated.SourceModifiedAccessConditions

// SASPermissions type simplifies creating the permissions string for an Azure Storage blob SAS.
// Initialize an instance of this type and then call its String method to set BlobSASSignatureValues's Permissions field.
type SASPermissions = exported.BlobSASPermissions

// Request Model Declaration -------------------------------------------------------------------------------------------

// DownloadOptions contains the optional parameters for the Client.Download method.
type DownloadOptions struct {
	// When set to true and specified together with the Range, the service returns the MD5 hash for the range, as long as the
	// range is less than or equal to 4 MB in size.
	RangeGetContentMD5 *bool

	// Optional, you can specify whether a particular range of the blob is read
	Offset *int64
	Count  *int64

	AccessConditions *AccessConditions
	CpkInfo          *CpkInfo
	CpkScopeInfo     *CpkScopeInfo
}

func (o *DownloadOptions) format() (*generated.BlobClientDownloadOptions, *generated.LeaseAccessConditions, *generated.CpkInfo, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil
	}

	offset := int64(0)
	count := int64(CountToEnd)

	if o.Offset != nil {
		offset = *o.Offset
	}

	if o.Count != nil {
		count = *o.Count
	}

	basics := generated.BlobClientDownloadOptions{
		RangeGetContentMD5: o.RangeGetContentMD5,
		Range:              shared.HTTPRange{Offset: offset, Count: count}.Format(),
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return &basics, leaseAccessConditions, o.CpkInfo, modifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// DownloadToWriterAtOptions identifies options used by the DownloadToBuffer and DownloadToFile functions.
type DownloadToWriterAtOptions struct {
	// BlockSize specifies the block size to use for each parallel download; the default size is DefaultDownloadBlockSize.
	BlockSize int64

	// Progress is a function that is invoked periodically as bytes are received.
	Progress func(bytesTransferred int64)

	// BlobAccessConditions indicates the access conditions used when making HTTP GET requests against the blob.
	AccessConditions *AccessConditions

	// ClientProvidedKeyOptions indicates the client provided key by name and/or by value to encrypt/decrypt data.
	CpkInfo      *CpkInfo
	CpkScopeInfo *CpkScopeInfo

	// Parallelism indicates the maximum number of blocks to download in parallel (0=default)
	Parallelism uint16

	// RetryReaderOptionsPerBlock is used when downloading each block.
	RetryReaderOptionsPerBlock RetryReaderOptions
}

func (o *DownloadToWriterAtOptions) getBlobPropertiesOptions() *GetPropertiesOptions {
	return &GetPropertiesOptions{
		AccessConditions: o.AccessConditions,
		CpkInfo:          o.CpkInfo,
	}
}

func (o *DownloadToWriterAtOptions) getDownloadBlobOptions(offSet, count int64, rangeGetContentMD5 *bool) *DownloadOptions {
	return &DownloadOptions{
		AccessConditions:   o.AccessConditions,
		CpkInfo:            o.CpkInfo,
		CpkScopeInfo:       o.CpkScopeInfo,
		Offset:             &offSet,
		Count:              &count,
		RangeGetContentMD5: rangeGetContentMD5,
	}
}

// DownloadToBufferOptions identifies options used by the DownloadToBuffer and DownloadToFile functions.
type DownloadToBufferOptions = DownloadToWriterAtOptions

// DownloadToFileOptions identifies options used by the DownloadToBuffer and DownloadToFile functions.
type DownloadToFileOptions = DownloadToWriterAtOptions

// ---------------------------------------------------------------------------------------------------------------------

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions struct {
	// Required if the blob has associated snapshots. Specify one of the following two options: include: Delete the base blob
	// and all of its snapshots. only: Delete only the blob's snapshots and not the blob itself
	DeleteSnapshots  *DeleteSnapshotsOptionType
	AccessConditions *AccessConditions
}

func (o *DeleteOptions) format() (*generated.BlobClientDeleteOptions, *generated.LeaseAccessConditions, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	basics := generated.BlobClientDeleteOptions{
		DeleteSnapshots: o.DeleteSnapshots,
	}

	if o.AccessConditions == nil {
		return &basics, nil, nil
	}

	return &basics, o.AccessConditions.LeaseAccessConditions, o.AccessConditions.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// UndeleteOptions contains the optional parameters for the Client.Undelete method.
type UndeleteOptions struct {
	// placeholder for future options
}

func (o *UndeleteOptions) format() *generated.BlobClientUndeleteOptions {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// SetTierOptions contains the optional parameters for the Client.SetTier method.
type SetTierOptions struct {
	// Optional: Indicates the priority with which to rehydrate an archived blob.
	RehydratePriority *RehydratePriority

	AccessConditions *AccessConditions
}

func (o *SetTierOptions) format() (*generated.BlobClientSetTierOptions, *generated.LeaseAccessConditions, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return &generated.BlobClientSetTierOptions{RehydratePriority: o.RehydratePriority}, leaseAccessConditions, modifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties method
type GetPropertiesOptions struct {
	AccessConditions *AccessConditions
	CpkInfo          *CpkInfo
}

func (o *GetPropertiesOptions) format() (*generated.BlobClientGetPropertiesOptions,
	*generated.LeaseAccessConditions, *generated.CpkInfo, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return nil, leaseAccessConditions, o.CpkInfo, modifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// SetHTTPHeadersOptions contains the optional parameters for the Client.SetHTTPHeaders method.
type SetHTTPHeadersOptions struct {
	AccessConditions *AccessConditions
}

func (o *SetHTTPHeadersOptions) format() (*generated.BlobClientSetHTTPHeadersOptions, *generated.LeaseAccessConditions, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return nil, leaseAccessConditions, modifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// SetMetadataOptions provides set of configurations for Set Metadata on blob operation
type SetMetadataOptions struct {
	AccessConditions *AccessConditions
	CpkInfo          *CpkInfo
	CpkScopeInfo     *CpkScopeInfo
}

func (o *SetMetadataOptions) format() (*generated.LeaseAccessConditions, *CpkInfo,
	*CpkScopeInfo, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return leaseAccessConditions, o.CpkInfo, o.CpkScopeInfo, modifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// CreateSnapshotOptions contains the optional parameters for the Client.CreateSnapshot method.
type CreateSnapshotOptions struct {
	Metadata         map[string]string
	AccessConditions *AccessConditions
	CpkInfo          *CpkInfo
	CpkScopeInfo     *CpkScopeInfo
}

func (o *CreateSnapshotOptions) format() (*generated.BlobClientCreateSnapshotOptions, *generated.CpkInfo,
	*generated.CpkScopeInfo, *generated.ModifiedAccessConditions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)

	return &generated.BlobClientCreateSnapshotOptions{
		Metadata: o.Metadata,
	}, o.CpkInfo, o.CpkScopeInfo, modifiedAccessConditions, leaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// StartCopyFromURLOptions contains the optional parameters for the Client.StartCopyFromURL method.
type StartCopyFromURLOptions struct {
	// Specifies the date time when the blobs immutability policy is set to expire.
	ImmutabilityPolicyExpiry *time.Time
	// Specifies the immutability policy mode to set on the blob.
	ImmutabilityPolicyMode *ImmutabilityPolicyMode
	// Specified if a legal hold should be set on the blob.
	LegalHold *bool
	// Optional. Used to set blob tags in various blob operations.
	BlobTags map[string]string
	// Optional. Specifies a user-defined name-value pair associated with the blob. If no name-value pairs are specified, the
	// operation will copy the metadata from the source blob or file to the destination blob. If one or more name-value pairs
	// are specified, the destination blob is created with the specified metadata, and metadata is not copied from the source
	// blob or file. Note that beginning with version 2009-09-19, metadata names must adhere to the naming rules for C# identifiers.
	// See Naming and Referencing Containers, Blobs, and Metadata for more information.
	Metadata map[string]string
	// Optional: Indicates the priority with which to rehydrate an archived blob.
	RehydratePriority *RehydratePriority
	// Overrides the sealed state of the destination blob. Service version 2019-12-12 and newer.
	SealBlob *bool
	// Optional. Indicates the tier to be set on the blob.
	Tier *AccessTier

	SourceModifiedAccessConditions *SourceModifiedAccessConditions

	AccessConditions *AccessConditions
}

func (o *StartCopyFromURLOptions) format() (*generated.BlobClientStartCopyFromURLOptions,
	*generated.SourceModifiedAccessConditions, *generated.ModifiedAccessConditions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil
	}

	basics := generated.BlobClientStartCopyFromURLOptions{
		BlobTagsString:           shared.SerializeBlobTagsToStrPtr(o.BlobTags),
		Metadata:                 o.Metadata,
		RehydratePriority:        o.RehydratePriority,
		SealBlob:                 o.SealBlob,
		Tier:                     o.Tier,
		ImmutabilityPolicyExpiry: o.ImmutabilityPolicyExpiry,
		ImmutabilityPolicyMode:   o.ImmutabilityPolicyMode,
		LegalHold:                o.LegalHold,
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return &basics, o.SourceModifiedAccessConditions, modifiedAccessConditions, leaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// AbortCopyFromURLOptions contains the optional parameters for the Client.AbortCopyFromURL method.
type AbortCopyFromURLOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *AbortCopyFromURLOptions) format() (*generated.BlobClientAbortCopyFromURLOptions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}
	return nil, o.LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// SetTagsOptions contains the optional parameters for the Client.SetTags method.
type SetTagsOptions struct {
	// The version id parameter is an opaque DateTime value that, when present,
	// specifies the version of the blob to operate on. It's for service version 2019-10-10 and newer.
	VersionID *string
	// Optional header, Specifies the transactional crc64 for the body, to be validated by the service.
	TransactionalContentCRC64 []byte
	// Optional header, Specifies the transactional md5 for the body, to be validated by the service.
	TransactionalContentMD5 []byte

	Tags map[string]string

	AccessConditions *AccessConditions
}

func (o *SetTagsOptions) format() (*generated.BlobClientSetTagsOptions, *ModifiedAccessConditions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	options := &generated.BlobClientSetTagsOptions{
		Tags:                      shared.SerializeBlobTags(o.Tags),
		TransactionalContentMD5:   o.TransactionalContentMD5,
		TransactionalContentCRC64: o.TransactionalContentCRC64,
		VersionID:                 o.VersionID,
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return options, modifiedAccessConditions, leaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// GetTagsOptions contains the optional parameters for the Client.GetTags method.
type GetTagsOptions struct {
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the blob snapshot to retrieve.
	Snapshot *string
	// The version id parameter is an opaque DateTime value that, when present, specifies the version of the blob to operate on.
	// It's for service version 2019-10-10 and newer.
	VersionID *string

	BlobAccessConditions *AccessConditions
}

func (o *GetTagsOptions) format() (*generated.BlobClientGetTagsOptions, *generated.ModifiedAccessConditions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	options := &generated.BlobClientGetTagsOptions{
		Snapshot:  o.Snapshot,
		VersionID: o.VersionID,
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.BlobAccessConditions)
	return options, modifiedAccessConditions, leaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// CopyFromURLOptions contains the optional parameters for the Client.CopyFromURL method.
type CopyFromURLOptions struct {
	// Optional. Used to set blob tags in various blob operations.
	BlobTags map[string]string
	// Only Bearer type is supported. Credentials should be a valid OAuth access token to copy source.
	CopySourceAuthorization *string
	// Specifies the date time when the blobs immutability policy is set to expire.
	ImmutabilityPolicyExpiry *time.Time
	// Specifies the immutability policy mode to set on the blob.
	ImmutabilityPolicyMode *ImmutabilityPolicyMode
	// Specified if a legal hold should be set on the blob.
	LegalHold *bool
	// Optional. Specifies a user-defined name-value pair associated with the blob. If no name-value pairs are specified, the
	// operation will copy the metadata from the source blob or file to the destination
	// blob. If one or more name-value pairs are specified, the destination blob is created with the specified metadata, and metadata
	// is not copied from the source blob or file. Note that beginning with
	// version 2009-09-19, metadata names must adhere to the naming rules for C# identifiers. See Naming and Referencing Containers,
	// Blobs, and Metadata for more information.
	Metadata map[string]string
	// Specify the md5 calculated for the range of bytes that must be read from the copy source.
	SourceContentMD5 []byte
	// Optional. Indicates the tier to be set on the blob.
	Tier *AccessTier

	SourceModifiedAccessConditions *SourceModifiedAccessConditions

	BlobAccessConditions *AccessConditions
}

func (o *CopyFromURLOptions) format() (*generated.BlobClientCopyFromURLOptions, *generated.SourceModifiedAccessConditions, *generated.ModifiedAccessConditions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil
	}

	options := &generated.BlobClientCopyFromURLOptions{
		BlobTagsString:           shared.SerializeBlobTagsToStrPtr(o.BlobTags),
		CopySourceAuthorization:  o.CopySourceAuthorization,
		ImmutabilityPolicyExpiry: o.ImmutabilityPolicyExpiry,
		ImmutabilityPolicyMode:   o.ImmutabilityPolicyMode,
		LegalHold:                o.LegalHold,
		Metadata:                 o.Metadata,
		SourceContentMD5:         o.SourceContentMD5,
		Tier:                     o.Tier,
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.BlobAccessConditions)
	return options, o.SourceModifiedAccessConditions, modifiedAccessConditions, leaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// LeaseBreakNaturally tells ContainerClient's or BlobClient's BreakLease method to break the lease using service semantics.
const LeaseBreakNaturally = -1

func leasePeriodPointer(period int32) *int32 {
	if period != LeaseBreakNaturally {
		return &period
	} else {
		return nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// AcquireLeaseOptions contains the optional parameters for the LeaseClient.AcquireLease method.
type AcquireLeaseOptions struct {
	// Specifies the Duration of the lease, in seconds, or negative one (-1) for a lease that never expires. A non-infinite lease
	// can be between 15 and 60 seconds. A lease Duration cannot be changed using renew or change.
	Duration *int32

	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *AcquireLeaseOptions) format() (generated.BlobClientAcquireLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return generated.BlobClientAcquireLeaseOptions{}, nil
	}
	return generated.BlobClientAcquireLeaseOptions{
		Duration: o.Duration,
	}, o.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// BreakLeaseOptions contains the optional parameters for the LeaseClient.BreakLease method.
type BreakLeaseOptions struct {
	// For a break operation, proposed Duration the lease should continue before it is broken, in seconds, between 0 and 60. This
	// break period is only used if it is shorter than the time remaining on the lease. If longer, the time remaining on the lease
	// is used. A new lease will not be available before the break period has expired, but the lease may be held for longer than
	// the break period. If this header does not appear with a break operation, a fixed-Duration lease breaks after the remaining
	// lease period elapses, and an infinite lease breaks immediately.
	BreakPeriod              *int32
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BreakLeaseOptions) format() (*generated.BlobClientBreakLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	if o.BreakPeriod != nil {
		period := leasePeriodPointer(*o.BreakPeriod)
		return &generated.BlobClientBreakLeaseOptions{
			BreakPeriod: period,
		}, o.ModifiedAccessConditions
	}

	return nil, o.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// ChangeLeaseOptions contains the optional parameters for the LeaseClient.ChangeLease method.
type ChangeLeaseOptions struct {
	ProposedLeaseID          *string
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ChangeLeaseOptions) format() (*string, *generated.BlobClientChangeLeaseOptions, *ModifiedAccessConditions, error) {
	generatedUuid, err := uuid.New()
	if err != nil {
		return nil, nil, nil, err
	}
	leaseID := to.Ptr(generatedUuid.String())
	if o == nil {
		return leaseID, nil, nil, nil
	}

	if o.ProposedLeaseID == nil {
		o.ProposedLeaseID = leaseID
	}

	return o.ProposedLeaseID, nil, o.ModifiedAccessConditions, nil
}

// ---------------------------------------------------------------------------------------------------------------------

// RenewLeaseOptions contains the optional parameters for the LeaseClient.RenewLease method.
type RenewLeaseOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *RenewLeaseOptions) format() (*generated.BlobClientRenewLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// ReleaseLeaseOptions contains the optional parameters for the LeaseClient.ReleaseLease method.
type ReleaseLeaseOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ReleaseLeaseOptions) format() (*generated.BlobClientReleaseLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------
