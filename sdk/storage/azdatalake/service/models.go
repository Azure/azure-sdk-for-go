//go:build go1.18

//

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated_blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"time"
)
import blobSAS "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"

// CreateFileSystemOptions contains the optional parameters for the FileSystem Create method.
type CreateFileSystemOptions = filesystem.CreateOptions

// DeleteFileSystemOptions contains the optional parameters for the FileSystem Delete method.
type DeleteFileSystemOptions = filesystem.DeleteOptions

// GetUserDelegationCredentialOptions contains optional parameters for GetUserDelegationKey method.
type GetUserDelegationCredentialOptions struct {
	// placeholder for future options
}

func (o *GetUserDelegationCredentialOptions) format() *generated_blob.ServiceClientGetUserDelegationKeyOptions {
	return nil
}

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties method.
type GetPropertiesOptions struct {
	// placeholder for future options
}

func (o *GetPropertiesOptions) format() *service.GetPropertiesOptions {
	if o == nil {
		return nil
	}
	return &service.GetPropertiesOptions{}
}

// SetPropertiesOptions provides set of options for Client.SetProperties
type SetPropertiesOptions struct {
	// CORS The set of CORS rules.
	CORS []*CORSRule
	// DefaultServiceVersion The default version to use for requests to the Datalake service if an incoming request's version is not specified. Possible
	// values include version 2008-10-27 and all more recent versions.
	DefaultServiceVersion *string
	// DeleteRetentionPolicy the retention policy which determines how long the associated data should persist.
	DeleteRetentionPolicy *RetentionPolicy
	// HourMetrics a summary of request statistics grouped by API in hour or minute aggregates
	// If version is not set - we default to "1.0"
	HourMetrics *Metrics
	// Logging Azure Analytics Logging settings.
	// If version is not set - we default to "1.0"
	Logging *Logging
	// MinuteMetrics a summary of request statistics grouped by API in hour or minute aggregates
	// If version is not set - we default to "1.0"
	MinuteMetrics *Metrics
	// StaticWebsite The properties that enable an account to host a static website.
	StaticWebsite *StaticWebsite
}

func (o *SetPropertiesOptions) format() *service.SetPropertiesOptions {
	if o == nil {
		return nil
	}
	return &service.SetPropertiesOptions{
		CORS:                  o.CORS,
		DefaultServiceVersion: o.DefaultServiceVersion,
		DeleteRetentionPolicy: o.DeleteRetentionPolicy,
		HourMetrics:           o.HourMetrics,
		Logging:               o.Logging,
		MinuteMetrics:         o.MinuteMetrics,
		StaticWebsite:         o.StaticWebsite,
	}
}

// ListFileSystemsInclude indicates what additional information the service should return with each filesystem.
type ListFileSystemsInclude struct {
	// Metadata tells the service whether to return metadata for each filesystem.
	Metadata *bool
	// Deleted tells the service whether to return soft-deleted filesystems.
	Deleted *bool
	// System tells the service whether to return system filesystems.
	System *bool
}

// ListFileSystemsOptions contains the optional parameters for the ListFileSystems method.
type ListFileSystemsOptions struct {
	// Include tells the service whether to return filesystem metadata.
	Include ListFileSystemsInclude
	// Marker is the continuation token to use when continuing the operation.
	Marker *string
	// MaxResults sets the maximum number of paths that will be returned per page.
	MaxResults *int32
	// Prefix filters the results to return only filesystems whose names begin with the specified prefix.
	Prefix *string
}

// GetSASURLOptions contains the optional parameters for the Client.GetSASURL method.
type GetSASURLOptions struct {
	// StartTime is the time after which the SAS will become valid.
	StartTime *time.Time
}

func (o *GetSASURLOptions) format(resources sas.AccountResourceTypes, permissions sas.AccountPermissions) (blobSAS.AccountResourceTypes, blobSAS.AccountPermissions, *service.GetSASURLOptions) {
	res := blobSAS.AccountResourceTypes{
		Service:   resources.Service,
		Container: resources.Container,
		Object:    resources.Object,
	}
	perms := blobSAS.AccountPermissions{
		Read:    permissions.Read,
		Write:   permissions.Write,
		Delete:  permissions.Delete,
		List:    permissions.List,
		Add:     permissions.Add,
		Create:  permissions.Create,
		Update:  permissions.Update,
		Process: permissions.Process,
	}
	if o == nil {
		return res, perms, nil
	}

	return res, perms, &service.GetSASURLOptions{
		StartTime: o.StartTime,
	}
}

// KeyInfo contains KeyInfo struct.
type KeyInfo = generated_blob.KeyInfo

// CORSRule - CORS is an HTTP feature that enables a web application running under one domain to access resources in another
// domain. Web browsers implement a security restriction known as same-origin policy that
// prevents a web page from calling APIs in a different domain; CORS provides a secure way to allow one domain (the origin
// domain) to call APIs in another domain.
type CORSRule = service.CORSRule

// RetentionPolicy - the retention policy which determines how long the associated data should persist.
type RetentionPolicy = service.RetentionPolicy

// Metrics - a summary of request statistics grouped by API in hour or minute aggregates
type Metrics = service.Metrics

// Logging - Azure Analytics Logging settings.
type Logging = service.Logging

// StaticWebsite - The properties that enable an account to host a static website.
type StaticWebsite = service.StaticWebsite

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// UserDelegationCredential contains an account's name and its user delegation key.
type UserDelegationCredential = exported.UserDelegationCredential

// UserDelegationKey contains UserDelegationKey.
type UserDelegationKey = exported.UserDelegationKey

// AccessConditions identifies blob-specific access conditions which you optionally set.
type AccessConditions = exported.AccessConditions

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = exported.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = exported.ModifiedAccessConditions

// CPKScopeInfo contains a group of parameters for the FileSystemClient.Create method.
type CPKScopeInfo = filesystem.CPKScopeInfo

// StorageServiceProperties - Storage Service Properties. Returned in GetServiceProperties call.
type StorageServiceProperties = service.StorageServiceProperties

// ListFileSystemsSegmentResponse contains fields from the ListFileSystems operation
type ListFileSystemsSegmentResponse = generated_blob.ListFileSystemsSegmentResponse

// FileSystemItem contains fields from the ListFileSystems operation
type FileSystemItem = generated_blob.FileSystemItem

// FileSystemProperties contains fields from the ListFileSystems operation
type FileSystemProperties = generated_blob.FileSystemProperties
