//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container

import (
	"fmt"
	"reflect"

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

// CpkScopeInfo contains a group of parameters for the ContainerClient.Create method.
type CpkScopeInfo = generated.ContainerCpkScopeInfo

// BlobProperties - Properties of a blob
type BlobProperties = generated.BlobPropertiesInternal

// BlobItem - An Azure Storage blob
type BlobItem = generated.BlobItemInternal

// AccessConditions identifies container-specific access conditions which you optionally set.
type AccessConditions = exported.ContainerAccessConditions

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = exported.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = exported.ModifiedAccessConditions

// AccessPolicy - An Access policy
type AccessPolicy = generated.AccessPolicy

// AccessPolicyPermission type simplifies creating the permissions string for a container's access policy.
// Initialize an instance of this type and then call its String method to set AccessPolicy's Permission field.
type AccessPolicyPermission = exported.AccessPolicyPermission

// SignedIdentifier - signed identifier
type SignedIdentifier = generated.SignedIdentifier

// Request Model Declaration -------------------------------------------------------------------------------------------

// CreateOptions contains the optional parameters for the Client.Create method.
type CreateOptions struct {
	// Specifies whether data in the container may be accessed publicly and the level of access
	Access *PublicAccessType

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata map[string]string

	// Optional. Specifies the encryption scope settings to set on the container.
	CpkScopeInfo *CpkScopeInfo
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions struct {
	AccessConditions *AccessConditions
}

func (o *DeleteOptions) format() (*generated.ContainerClientDeleteOptions, *generated.LeaseAccessConditions, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatContainerAccessConditions(o.AccessConditions)
	return nil, leaseAccessConditions, modifiedAccessConditions
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
func (o *GetPropertiesOptions) format() (*generated.ContainerClientGetPropertiesOptions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// ListBlobsInclude indicates what additional information the service should return with each blob.
type ListBlobsInclude struct {
	Copy, Metadata, Snapshots, UncommittedBlobs, Deleted, Tags, Versions, LegalHold, ImmutabilityPolicy, DeletedWithVersions bool
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
		include = append(include, generated.ListBlobsIncludeItemDeletedwithversions)
	}
	if l.ImmutabilityPolicy {
		include = append(include, generated.ListBlobsIncludeItemImmutabilitypolicy)
	}
	if l.LegalHold {
		include = append(include, generated.ListBlobsIncludeItemLegalhold)
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
		include = append(include, generated.ListBlobsIncludeItemUncommittedblobs)
	}
	if l.Versions {
		include = append(include, generated.ListBlobsIncludeItemVersions)
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
	// Specifies the maximum number of containers to return. If the request does not specify maxresults, or specifies a value
	// greater than 5000, the server will return up to 5000 items. Note that if the
	// listing operation crosses a partition boundary, then the service will return a continuation token for retrieving the remainder
	// of the results. For this reason, it is possible that the service will
	// return fewer results than specified by maxresults, or than the default of 5000.
	MaxResults *int32
	// Filters the results to return only containers whose name begins with the specified prefix.
	Prefix *string
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
	// Specifies the maximum number of containers to return. If the request does not specify maxresults, or specifies a value
	// greater than 5000, the server will return up to 5000 items. Note that if the
	// listing operation crosses a partition boundary, then the service will return a continuation token for retrieving the remainder
	// of the results. For this reason, it is possible that the service will
	// return fewer results than specified by maxresults, or than the default of 5000.
	MaxResults *int32
	// Filters the results to return only containers whose name begins with the specified prefix.
	Prefix *string
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

// SetMetadataOptions contains the optional parameters for the Client.SetMetadata method.
type SetMetadataOptions struct {
	Metadata                 map[string]string
	LeaseAccessConditions    *LeaseAccessConditions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *SetMetadataOptions) format() (*generated.ContainerClientSetMetadataOptions, *generated.LeaseAccessConditions, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	return &generated.ContainerClientSetMetadataOptions{Metadata: o.Metadata}, o.LeaseAccessConditions, o.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// GetAccessPolicyOptions contains the optional parameters for the Client.GetAccessPolicy method.
type GetAccessPolicyOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetAccessPolicyOptions) format() (*generated.ContainerClientGetAccessPolicyOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// SetAccessPolicyOptions provides set of configurations for ContainerClient.SetAccessPolicy operation
type SetAccessPolicyOptions struct {
	// Specifies whether data in the container may be accessed publicly and the level of access
	Access *PublicAccessType
	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	AccessConditions *AccessConditions
}

func (o *SetAccessPolicyOptions) format() (*generated.ContainerClientSetAccessPolicyOptions, *LeaseAccessConditions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}
	lac, mac := exported.FormatContainerAccessConditions(o.AccessConditions)
	return &generated.ContainerClientSetAccessPolicyOptions{
		Access: o.Access,
	}, lac, mac
}

// ---------------------------------------------------------------------------------------------------------------------

type BatchDeleteOptions struct {
	BlobName          *string
	BlobDeleteOptions *blob.DeleteOptions
	VersionID         *string
	Snapshot          *string
}

func (o *BatchDeleteOptions) createDeleteSubRequest(ctx context.Context) (*policy.Request, error)  {
				//options *BlobClientDeleteOptions, 
				//leaseAccessConditions *LeaseAccessConditions, 
				//modifiedAccessConditions *ModifiedAccessConditions) 
	if o.BlobName == nil {
		return nil, fmt.Errorf("blob name not provided")
	}

	req, err := runtime.NewRequest(ctx, http.MethodDelete, o.BlobName)
	if err != nil {
		return nil, err
	}

	reqQP := req.Raw().URL.Query()
	if o !=  nil && o.Snapshot != nil {
		reqQP.Set("snapshot", *o.Snapshot)
	}
	if o !=  nil && o.VersionID != nil {
		reqQP.Set("versionid", *o.VersionID)
	}
	if o !=  nil && o.BlobDeleteOptions != nil && o.BlobDeleteOptions.DeleteType != nil {
		reqQP.Set("deletetype", string(*options.DeleteType))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	if leaseAccessConditions != nil && leaseAccessConditions.LeaseID != nil {
		req.Raw().Header["x-ms-lease-id"] = []string{*leaseAccessConditions.LeaseID}
	}
	if o !=  nil && options.DeleteSnapshots != nil {
		req.Raw().Header["x-ms-delete-snapshots"] = []string{string(*options.DeleteSnapshots)}
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfModifiedSince != nil {
		req.Raw().Header["If-Modified-Since"] = []string{modifiedAccessConditions.IfModifiedSince.Format(time.RFC1123)}
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfUnmodifiedSince != nil {
		req.Raw().Header["If-Unmodified-Since"] = []string{modifiedAccessConditions.IfUnmodifiedSince.Format(time.RFC1123)}
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfMatch != nil {
		req.Raw().Header["If-Match"] = []string{string(*modifiedAccessConditions.IfMatch)}
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfNoneMatch != nil {
		req.Raw().Header["If-None-Match"] = []string{string(*modifiedAccessConditions.IfNoneMatch)}
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfTags != nil {
		req.Raw().Header["x-ms-if-tags"] = []string{*modifiedAccessConditions.IfTags}
	}
	req.Raw().Header["x-ms-version"] = []string{"2020-10-02"}
	if o !=  nil && options.RequestID != nil {
		req.Raw().Header["x-ms-client-request-id"] = []string{*options.RequestID}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

func (o *BatchDeleteOptions) format() (string, error) {
	if o == nil {
		return "", nil
	}

	&basics, o.AccessConditions.LeaseAccessConditions, o.AccessConditions.ModifiedAccessConditions

}
