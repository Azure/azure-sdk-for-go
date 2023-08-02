//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"time"
)

// SetAccessPolicyOptions provides set of configurations for Filesystem.SetAccessPolicy operation.
type SetAccessPolicyOptions struct {
	// Specifies whether data in the filesystem may be accessed publicly and the level of access.
	// If this header is not included in the request, filesystem data is private to the account owner.
	Access           *PublicAccessType
	AccessConditions *AccessConditions
	FilesystemACL    []*SignedIdentifier
}

func (o *SetAccessPolicyOptions) format() *container.SetAccessPolicyOptions {
	if o == nil {
		return nil
	}
	return &container.SetAccessPolicyOptions{
		Access:           o.Access,
		AccessConditions: exported.FormatContainerAccessConditions(o.AccessConditions),
		ContainerACL:     o.FilesystemACL,
	}
}

// CreateOptions contains the optional parameters for the Client.Create method.
type CreateOptions struct {
	// Specifies whether data in the filesystem may be accessed publicly and the level of access.
	Access *PublicAccessType

	// Optional. Specifies a user-defined name-value pair associated with the filesystem.
	Metadata map[string]*string

	// Optional. Specifies the encryption scope settings to set on the filesystem.
	CPKScopeInfo *CPKScopeInfo
}

func (o *CreateOptions) format() *container.CreateOptions {
	if o == nil {
		return nil
	}
	return &container.CreateOptions{
		Access:       o.Access,
		Metadata:     o.Metadata,
		CPKScopeInfo: o.CPKScopeInfo,
	}
}

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions struct {
	AccessConditions *AccessConditions
}

func (o *DeleteOptions) format() *container.DeleteOptions {
	if o == nil {
		return nil
	}
	return &container.DeleteOptions{
		AccessConditions: exported.FormatContainerAccessConditions(o.AccessConditions),
	}
}

// GetPropertiesOptions contains the optional parameters for the FilesystemClient.GetProperties method.
type GetPropertiesOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetPropertiesOptions) format() *container.GetPropertiesOptions {
	if o == nil {
		return nil
	}
	if o.LeaseAccessConditions == nil {
		o.LeaseAccessConditions = &LeaseAccessConditions{}
	}
	return &container.GetPropertiesOptions{
		LeaseAccessConditions: &container.LeaseAccessConditions{
			LeaseID: o.LeaseAccessConditions.LeaseID,
		},
	}
}

// SetMetadataOptions contains the optional parameters for the Client.SetMetadata method.
type SetMetadataOptions struct {
	Metadata         map[string]*string
	AccessConditions *AccessConditions
}

func (o *SetMetadataOptions) format() *container.SetMetadataOptions {
	if o == nil {
		return nil
	}
	accConditions := exported.FormatContainerAccessConditions(o.AccessConditions)
	return &container.SetMetadataOptions{
		Metadata:                 o.Metadata,
		LeaseAccessConditions:    accConditions.LeaseAccessConditions,
		ModifiedAccessConditions: accConditions.ModifiedAccessConditions,
	}
}

// GetAccessPolicyOptions contains the optional parameters for the Client.GetAccessPolicy method.
type GetAccessPolicyOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetAccessPolicyOptions) format() *container.GetAccessPolicyOptions {
	if o == nil {
		return nil
	}
	if o.LeaseAccessConditions == nil {
		o.LeaseAccessConditions = &LeaseAccessConditions{}
	}
	return &container.GetAccessPolicyOptions{
		LeaseAccessConditions: &container.LeaseAccessConditions{
			LeaseID: o.LeaseAccessConditions.LeaseID,
		},
	}
}

// CPKScopeInfo contains a group of parameters for the FilesystemClient.Create method.
type CPKScopeInfo = container.CPKScopeInfo

// AccessPolicy - An Access policy.
type AccessPolicy = container.AccessPolicy

// AccessPolicyPermission type simplifies creating the permissions string for a container's access policy.
// Initialize an instance of this type and then call its String method to set AccessPolicy's Permission field.
type AccessPolicyPermission = exported.AccessPolicyPermission

// SignedIdentifier - signed identifier.
type SignedIdentifier = container.SignedIdentifier

// ListPathsOptions contains the optional parameters for the Filesystem.ListPaths operation.
type ListPathsOptions struct {
	Marker     *string
	MaxResults *int32
	Prefix     *string
	Upn        *bool
}

func (o *ListPathsOptions) format() generated.FileSystemClientListPathsOptions {
	if o == nil {
		return generated.FileSystemClientListPathsOptions{}
	}

	return generated.FileSystemClientListPathsOptions{
		Continuation: o.Marker,
		MaxResults:   o.MaxResults,
		Path:         o.Prefix,
		Upn:          o.Upn,
	}
}

// ListDeletedPathsOptions contains the optional parameters for the Filesystem.ListDeletedPaths operation. PLACEHOLDER
type ListDeletedPathsOptions struct {
	Marker     *string
	MaxResults *int32
	Prefix     *string
}

func (o *ListDeletedPathsOptions) format() generated.FileSystemClientListBlobHierarchySegmentOptions {
	showOnly := "deleted"
	if o == nil {
		return generated.FileSystemClientListBlobHierarchySegmentOptions{Showonly: &showOnly}
	}
	return generated.FileSystemClientListBlobHierarchySegmentOptions{
		Marker:     o.Marker,
		MaxResults: o.MaxResults,
		Prefix:     o.Prefix,
		Showonly:   &showOnly,
	}
}

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

// UndeletePathOptions contains the optional parameters for the Filesystem.UndeletePath operation.
type UndeletePathOptions struct {
	// placeholder
}

func (o *UndeletePathOptions) format() *UndeletePathOptions {
	if o == nil {
		return nil
	}
	return &UndeletePathOptions{}
}

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// AccessConditions identifies blob-specific access conditions which you optionally set.
type AccessConditions = exported.AccessConditions

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = exported.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = exported.ModifiedAccessConditions
