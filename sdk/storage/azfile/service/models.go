// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// UserDelegationCredential contains an account's name and its user delegation key.
type UserDelegationCredential = exported.UserDelegationCredential

// KeyInfo contains KeyInfo struct.
type KeyInfo = generated.KeyInfo

// GetUserDelegationCredentialOptions contains optional parameters for GetUserDelegationKey method.
type GetUserDelegationCredentialOptions struct{}

func (o *GetUserDelegationCredentialOptions) format() *generated.ServiceClientGetUserDelegationKeyOptions {
	return nil
}

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	return exported.NewSharedKeyCredential(accountName, accountKey)
}

// CreateShareOptions contains the optional parameters for the share.Client.Create method.
type CreateShareOptions = share.CreateOptions

// DeleteShareOptions contains the optional parameters for the share.Client.Delete method.
type DeleteShareOptions = share.DeleteOptions

// RestoreShareOptions contains the optional parameters for the share.Client.Restore method.
type RestoreShareOptions = share.RestoreOptions

// ---------------------------------------------------------------------------------------------------------------------

// GetPropertiesOptions provides set of options for Client.GetProperties
type GetPropertiesOptions struct {
	// placeholder for future options
}

func (o *GetPropertiesOptions) format() *generated.ServiceClientGetPropertiesOptions {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// SetPropertiesOptions provides set of options for Client.SetProperties
type SetPropertiesOptions struct {
	// The set of CORS rules.
	CORS []*CORSRule

	// A summary of request statistics grouped by API in hourly aggregates for files.
	HourMetrics *Metrics

	// A summary of request statistics grouped by API in minute aggregates for files.
	MinuteMetrics *Metrics

	// Protocol settings
	Protocol *ProtocolSettings
}

func (o *SetPropertiesOptions) format() (generated.StorageServiceProperties, *generated.ServiceClientSetPropertiesOptions) {
	if o == nil {
		return generated.StorageServiceProperties{}, nil
	}

	formatMetrics(o.HourMetrics)
	formatMetrics(o.MinuteMetrics)

	return generated.StorageServiceProperties{
		CORS:          o.CORS,
		HourMetrics:   o.HourMetrics,
		MinuteMetrics: o.MinuteMetrics,
		Protocol:      o.Protocol,
	}, nil
}

// update version of Storage Analytics to configure. Use 1.0 for this value.
func formatMetrics(m *Metrics) {
	if m == nil {
		return
	}

	m.Version = to.Ptr(shared.StorageAnalyticsVersion)
}

// StorageServiceProperties - Storage service properties.
type StorageServiceProperties = generated.StorageServiceProperties

// CORSRule - CORS is an HTTP feature that enables a web application running under one domain to access resources in
// another domain. Web browsers implement a security restriction known as same-origin policy that
// prevents a web page from calling APIs in a different domain; CORS provides a secure way to allow one domain (the origin
// domain) to call APIs in another domain.
type CORSRule = generated.CORSRule

// Metrics - Storage Analytics metrics for file service.
type Metrics = generated.Metrics

// RetentionPolicy - The retention policy.
type RetentionPolicy = generated.RetentionPolicy

// ProtocolSettings - Protocol settings
type ProtocolSettings = generated.ProtocolSettings

// SMBSettings - Settings for SMB protocol.
type SMBSettings = generated.SMBSettings

// SMBMultichannel - Settings for SMB multichannel
type SMBMultichannel = generated.SMBMultichannel

// ---------------------------------------------------------------------------------------------------------------------

// ListSharesOptions contains the optional parameters for the Client.NewListSharesPager method.
type ListSharesOptions struct {
	// Include this parameter to specify one or more datasets to include in the responseBody.
	Include ListSharesInclude

	// A string value that identifies the portion of the list to be returned with the next list operation. The operation returns
	// a marker value within the responseBody body if the list returned was not complete.
	// The marker value may then be used in a subsequent call to request the next set of list items. The marker value is opaque
	// to the client.
	Marker *string

	// Specifies the maximum number of entries to return. If the request does not specify maxresults, or specifies a value greater
	// than 5,000, the server will return up to 5,000 items.
	MaxResults *int32

	// Filters the results to return only entries whose name begins with the specified prefix.
	Prefix *string
}

// ListSharesInclude indicates what additional information the service should return with each share.
type ListSharesInclude struct {
	// Tells the service whether to return metadata for each share.
	Metadata bool

	// Tells the service whether to return soft-deleted shares.
	Deleted bool

	// Tells the service whether to return share snapshots.
	Snapshots bool
}

// Share - A listed Azure Storage share item.
type Share = generated.Share

// ShareProperties - Properties of a share.
type ShareProperties = generated.ShareProperties

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
