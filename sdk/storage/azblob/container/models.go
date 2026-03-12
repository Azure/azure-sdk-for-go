// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container

import (
	"reflect"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	return exported.NewSharedKeyCredential(accountName, accountKey)
}

// Request Model Declaration -------------------------------------------------------------------------------------------

// CPKScopeInfo contains a group of parameters for the ContainerClient.Create method.
type CPKScopeInfo = generated.ContainerCPKScopeInfo

// BlobFlatListSegment - List of BlobItem.
type BlobFlatListSegment = generated.BlobFlatListSegment

// BlobHierarchyListSegment - List of BlobItem and BlobPrefix.
type BlobHierarchyListSegment = generated.BlobHierarchyListSegment

// BlobProperties - Properties of a blob.
type BlobProperties = generated.BlobProperties

// BlobItem - An Azure Storage blob.
type BlobItem = generated.BlobItem

// BlobTags - Blob tags.
type BlobTags = generated.BlobTags

// BlobPrefix is a blob's prefix when hierarchically listing blobs.
type BlobPrefix = generated.BlobPrefix

// BlobTag - a key/value pair on a blob.
type BlobTag = generated.BlobTag

// AccessConditions identifies container-specific access conditions which you optionally set.
type AccessConditions = exported.ContainerAccessConditions

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = exported.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = exported.ModifiedAccessConditions

// AccessPolicy - An Access policy.
type AccessPolicy = generated.AccessPolicy

// AccessPolicyPermission type simplifies creating the permissions string for a container's access policy.
// Initialize an instance of this type and then call its String method to set AccessPolicy's Permission field.
type AccessPolicyPermission = exported.AccessPolicyPermission

// SignedIdentifier - signed identifier.
type SignedIdentifier = generated.SignedIdentifier

// Request Model Declaration -------------------------------------------------------------------------------------------

// CreateOptions contains the optional parameters for the Client.Create method.
type CreateOptions struct {
	// Specifies whether data in the container may be accessed publicly and the level of access.
	Access *PublicAccessType

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata map[string]*string

	// Optional. Specifies the encryption scope settings to set on the container.
	CPKScopeInfo *CPKScopeInfo
}

func (o *CreateOptions) format() *generated.ContainerClientCreateOptions {
	if o == nil {
		return nil
	}

	opts := &generated.ContainerClientCreateOptions{
		Access:   o.Access,
		Metadata: o.Metadata,
	}
	if o.CPKScopeInfo != nil {
		opts.DefaultEncryptionScope = o.CPKScopeInfo.DefaultEncryptionScope
		opts.PreventEncryptionScopeOverride = o.CPKScopeInfo.PreventEncryptionScopeOverride
	}

	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions struct {
	AccessConditions *AccessConditions
}

func (o *DeleteOptions) format() *generated.ContainerClientDeleteOptions {
	if o == nil {
		return nil
	}

	opts := &generated.ContainerClientDeleteOptions{}
	if o.AccessConditions != nil {
		if o.AccessConditions.LeaseAccessConditions != nil {
			opts.LeaseID = o.AccessConditions.LeaseAccessConditions.LeaseID
		}
		if o.AccessConditions.ModifiedAccessConditions != nil {
			opts.IfModifiedSince = o.AccessConditions.ModifiedAccessConditions.IfModifiedSince
			opts.IfUnmodifiedSince = o.AccessConditions.ModifiedAccessConditions.IfUnmodifiedSince
		}
	}

	// Note: missing o.AccessConditions.ModifiedAccessConditions.IfMatch and o.AccessConditions.ModifiedAccessConditions.IfNoneMatch
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// RestoreOptions contains the optional parameters for the Client.Restore method.
type RestoreOptions struct {
	// placeholder for future options
}

// ---------------------------------------------------------------------------------------------------------------------

// GetPropertiesOptions contains the optional parameters for the ContainerClient.GetProperties method.
type GetPropertiesOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

// ContainerClientGetPropertiesOptions contains the optional parameters for the ContainerClient.GetProperties method.
func (o *GetPropertiesOptions) format() *generated.ContainerClientGetPropertiesOptions {
	if o == nil {
		return nil
	}

	opts := &generated.ContainerClientGetPropertiesOptions{}
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}

	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// ListBlobsInclude indicates what additional information the service should return with each blob.
type ListBlobsInclude struct {
	Copy, Metadata, Snapshots, UncommittedBlobs, Deleted, Tags, Versions, LegalHold, ImmutabilityPolicy, DeletedWithVersions, Permissions bool
}

func (l ListBlobsInclude) format() []generated.ListBlobsIncludeItem {
	if reflect.ValueOf(l).IsZero() {
		return nil
	}

	include := []generated.ListBlobsIncludeItem{}

	if l.Copy {
		include = append(include, generated.ListBlobsIncludeItemCopy)
	}
	if l.Deleted {
		include = append(include, generated.ListBlobsIncludeItemDeleted)
	}
	if l.DeletedWithVersions {
		include = append(include, generated.ListBlobsIncludeItemDeletedWithVersions)
	}
	if l.ImmutabilityPolicy {
		include = append(include, generated.ListBlobsIncludeItemImmutabilityPolicy)
	}
	if l.LegalHold {
		include = append(include, generated.ListBlobsIncludeItemLegalHold)
	}
	if l.Metadata {
		include = append(include, generated.ListBlobsIncludeItemMetadata)
	}
	if l.Snapshots {
		include = append(include, generated.ListBlobsIncludeItemSnapshots)
	}
	if l.Tags {
		include = append(include, generated.ListBlobsIncludeItemTags)
	}
	if l.UncommittedBlobs {
		include = append(include, generated.ListBlobsIncludeItemUncommittedBlobs)
	}
	if l.Versions {
		include = append(include, generated.ListBlobsIncludeItemVersions)
	}
	if l.Permissions {
		include = append(include, generated.ListBlobsIncludeItemPermissions)
	}
	return include
}

// ListBlobsFlatOptions contains the optional parameters for the ContainerClient.ListBlobFlatSegment method.
type ListBlobsFlatOptions struct {
	// Include this parameter to specify one or more datasets to include in the response.
	Include ListBlobsInclude
	// A string value that identifies the portion of the list of containers to be returned with the next listing operation. The
	// operation returns the NextMarker value within the response body if the listing
	// operation did not return all containers remaining to be listed with the current page. The NextMarker value can be used
	// as the value for the marker parameter in a subsequent call to request the next
	// page of list items. The marker value is opaque to the client.
	Marker *string
	// Specifies the maximum number of containers to return. If the request does not specify MaxResults, or specifies a value
	// greater than 5000, the server will return up to 5000 items. Note that if the
	// listing operation crosses a partition boundary, then the service will return a continuation token for retrieving the remainder
	// of the results. For this reason, it is possible that the service will
	// return fewer results than specified by MaxResults, or than the default of 5000.
	MaxResults *int32
	// Filters the results to return only containers whose name begins with the specified prefix.
	Prefix *string
	// Specifies the relative path to list paths from. For non-recursive list, only one entity level is supported; For recursive
	// list, multiple entity levels are supported. (Inclusive)
	StartFrom *string
}

// ---------------------------------------------------------------------------------------------------------------------

// ListBlobsHierarchyOptions provides set of configurations for Client.NewListBlobsHierarchyPager
type ListBlobsHierarchyOptions struct {
	// Include this parameter to specify one or more datasets to include in the response.
	Include ListBlobsInclude
	// A string value that identifies the portion of the list of containers to be returned with the next listing operation. The
	// operation returns the NextMarker value within the response body if the listing
	// operation did not return all containers remaining to be listed with the current page. The NextMarker value can be used
	// as the value for the marker parameter in a subsequent call to request the next
	// page of list items. The marker value is opaque to the client.
	Marker *string
	// Specifies the maximum number of containers to return. If the request does not specify MaxResults, or specifies a value
	// greater than 5000, the server will return up to 5000 items. Note that if the
	// listing operation crosses a partition boundary, then the service will return a continuation token for retrieving the remainder
	// of the results. For this reason, it is possible that the service will
	// return fewer results than specified by MaxResults, or than the default of 5000.
	MaxResults *int32
	// Filters the results to return only containers whose name begins with the specified prefix.
	Prefix *string
	// Specifies the relative path to list paths from. For non-recursive list, only one entity level is supported; For recursive
	// list, multiple entity levels are supported. (Inclusive)
	StartFrom *string
}

// ContainerClientListBlobHierarchySegmentOptions contains the optional parameters for the ContainerClient.ListBlobHierarchySegment method.
func (o *ListBlobsHierarchyOptions) format() generated.ContainerClientListBlobHierarchySegmentOptions {
	if o == nil {
		return generated.ContainerClientListBlobHierarchySegmentOptions{}
	}

	return generated.ContainerClientListBlobHierarchySegmentOptions{
		Include:    o.Include.format(),
		Marker:     o.Marker,
		Maxresults: o.MaxResults,
		Prefix:     o.Prefix,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// GetSASURLOptions contains the optional parameters for the Client.GetSASURL method.
type GetSASURLOptions struct {
	StartTime *time.Time
}

func (o *GetSASURLOptions) format() time.Time {
	if o == nil {
		return time.Time{}
	}

	var st time.Time
	if o.StartTime != nil {
		st = o.StartTime.UTC()
	} else {
		st = time.Time{}
	}
	return st
}

// ---------------------------------------------------------------------------------------------------------------------

// SetMetadataOptions contains the optional parameters for the Client.SetMetadata method.
type SetMetadataOptions struct {
	Metadata                 map[string]*string
	LeaseAccessConditions    *LeaseAccessConditions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *SetMetadataOptions) format() *generated.ContainerClientSetMetadataOptions {
	if o == nil {
		return nil
	}

	opts := &generated.ContainerClientSetMetadataOptions{
		Metadata: o.Metadata,
	}
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	if o.ModifiedAccessConditions != nil {
		opts.IfModifiedSince = o.ModifiedAccessConditions.IfModifiedSince
	}

	// Note: missing mapping for most of o.ModifiedAccessConditions
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// GetAccessPolicyOptions contains the optional parameters for the Client.GetAccessPolicy method.
type GetAccessPolicyOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetAccessPolicyOptions) format() *generated.ContainerClientGetAccessPolicyOptions {
	if o == nil {
		return nil
	}

	opts := &generated.ContainerClientGetAccessPolicyOptions{}
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}

	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// SetAccessPolicyOptions provides set of configurations for ContainerClient.SetAccessPolicy operation.
type SetAccessPolicyOptions struct {
	// Specifies whether data in the container may be accessed publicly and the level of access.
	// If this header is not included in the request, container data is private to the account owner.
	Access           *PublicAccessType
	AccessConditions *AccessConditions
	ContainerACL     []*SignedIdentifier
}

func (o *SetAccessPolicyOptions) format() *generated.ContainerClientSetAccessPolicyOptions {
	if o == nil {
		return nil
	}

	opts := &generated.ContainerClientSetAccessPolicyOptions{
		Access: o.Access,
	}
	if o.AccessConditions != nil {
		if o.AccessConditions.LeaseAccessConditions != nil {
			opts.LeaseID = o.AccessConditions.LeaseAccessConditions.LeaseID
		}
		if o.AccessConditions.ModifiedAccessConditions != nil {
			opts.IfModifiedSince = o.AccessConditions.ModifiedAccessConditions.IfModifiedSince
			opts.IfUnmodifiedSince = o.AccessConditions.ModifiedAccessConditions.IfUnmodifiedSince
		}
	}

	// Note: missing mapping for o.AccessConditions.ModifiedAccessConditions.IfMatch,  o.AccessConditions.ModifiedAccessConditions.IfNoneMatch,
	return opts
}

func formatTime(c *SignedIdentifier) error {
	if c == nil {
		return nil
	}

	if c.AccessPolicy == nil {
		return nil
	}

	if c.AccessPolicy.Start != nil {
		st, err := time.Parse(time.RFC3339, c.AccessPolicy.Start.UTC().Format(time.RFC3339))
		if err != nil {
			return err
		}
		c.AccessPolicy.Start = &st
	}
	if c.AccessPolicy.Expiry != nil {
		et, err := time.Parse(time.RFC3339, c.AccessPolicy.Expiry.UTC().Format(time.RFC3339))
		if err != nil {
			return err
		}
		c.AccessPolicy.Expiry = &et
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// GetAccountInfoOptions provides set of options for Client.GetAccountInfo
type GetAccountInfoOptions struct {
	// placeholder for future options
}

func (o *GetAccountInfoOptions) format() *generated.ContainerClientGetAccountInfoOptions {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// BatchDeleteOptions contains the optional parameters for the BatchBuilder.Delete method.
type BatchDeleteOptions struct {
	blob.DeleteOptions
	VersionID *string
	Snapshot  *string
}

func (o *BatchDeleteOptions) format() *generated.BlobClientDeleteOptions {
	if o == nil {
		return nil
	}

	opts := &generated.BlobClientDeleteOptions{
		DeleteSnapshots: o.DeleteSnapshots,
		BlobDeleteType:  o.BlobDeleteType, // None by default
		Snapshot:        o.Snapshot,
		VersionID:       o.VersionID,
	}
	if o.AccessConditions != nil {
		if o.AccessConditions.LeaseAccessConditions != nil {
			opts.LeaseID = o.AccessConditions.LeaseAccessConditions.LeaseID
		}
		if o.AccessConditions.ModifiedAccessConditions != nil {
			opts.IfMatch = o.AccessConditions.ModifiedAccessConditions.IfMatch
			opts.IfModifiedSince = o.AccessConditions.ModifiedAccessConditions.IfModifiedSince
			opts.IfNoneMatch = o.AccessConditions.ModifiedAccessConditions.IfNoneMatch
			opts.IfUnmodifiedSince = o.AccessConditions.ModifiedAccessConditions.IfUnmodifiedSince
		}
	}

	return opts
}

// BatchSetTierOptions contains the optional parameters for the BatchBuilder.SetTier method.
type BatchSetTierOptions struct {
	blob.SetTierOptions
	VersionID *string
	Snapshot  *string
}

func (o *BatchSetTierOptions) format() *generated.BlobClientSetTierOptions {
	if o == nil {
		return nil
	}

	opts := &generated.BlobClientSetTierOptions{
		RehydratePriority: o.RehydratePriority,
		Snapshot:          o.Snapshot,
		VersionID:         o.VersionID,
	}
	if o.AccessConditions != nil && o.AccessConditions.LeaseAccessConditions != nil {
		opts.LeaseID = o.AccessConditions.LeaseAccessConditions.LeaseID
	}

	// Note: missing mapping for o.AccessConditions.ModifiedAccessConditions
	return opts
}

// SubmitBatchOptions contains the optional parameters for the Client.SubmitBatch method.
type SubmitBatchOptions struct {
	// placeholder for future options
}

func (o *SubmitBatchOptions) format() *generated.ContainerClientSubmitBatchOptions {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// FilterBlobsOptions provides set of options for Client.FilterBlobs.
type FilterBlobsOptions struct {
	// A string value that identifies the portion of the list of containers to be returned with the next listing operation. The
	// operation returns the NextMarker value within the response body if the listing
	// operation did not return all containers remaining to be listed with the current page. The NextMarker value can be used
	// as the value for the marker parameter in a subsequent call to request the next
	// page of list items. The marker value is opaque to the client.
	Marker *string
	// Specifies the maximum number of containers to return. If the request does not specify maxresults, or specifies a value
	// greater than 5000, the server will return up to 5000 items. Note that if the
	// listing operation crosses a partition boundary, then the service will return a continuation token for retrieving the remainder
	// of the results. For this reason, it is possible that the service will
	// return fewer results than specified by maxresults, or than the default of 5000.
	MaxResults *int32
}

func (o *FilterBlobsOptions) format() *generated.ContainerClientFilterBlobsOptions {
	if o == nil {
		return nil
	}
	return &generated.ContainerClientFilterBlobsOptions{
		Marker:     o.Marker,
		Maxresults: o.MaxResults,
	}
}
