//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// Client represents a URL to the Azure Storage share allowing you to manipulate its directories and files.
type Client base.Client[generated.ShareClient]

// NewClient creates an instance of Client with the specified values.
//   - shareURL - the URL of the share e.g. https://<account>.file.core.windows.net/share
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(shareURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a share or with a shared access signature (SAS) token.
//   - shareURL - the URL of the share e.g. https://<account>.file.core.windows.net/share?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(shareURL string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - shareURL - the URL of the share e.g. https://<account>.file.core.windows.net/share
//   - cred - a SharedKeyCredential created with the matching share's storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(shareURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - shareName - the name of the share within the storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, shareName string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

func (s *Client) generated() *generated.ShareClient {
	return base.InnerClient((*base.Client[generated.ShareClient])(s))
}

func (s *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.ShareClient])(s))
}

// URL returns the URL endpoint used by the Client object.
func (s *Client) URL() string {
	return "s.generated().Endpoint()"
}

// NewDirectoryClient creates a new directory.Client object by concatenating directoryName to the end of this Client's URL.
// The new directory.Client uses the same request policy pipeline as the Client.
func (s *Client) NewDirectoryClient(directoryName string) *directory.Client {
	return nil
}

// NewRootDirectoryClient creates a new directory.Client object using the Client's URL.
// The new directory.Client uses the same request policy pipeline as the Client.
func (s *Client) NewRootDirectoryClient() *directory.Client {
	return nil
}

// WithSnapshot creates a new Client object identical to the source but with the specified share snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base share.
func (s *Client) WithSnapshot(shareSnapshot string) (*Client, error) {
	return nil, nil
}

// Create operation creates a new share within a storage account. If a share with the same name already exists, the operation fails.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/create-share.
func (s *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	return CreateResponse{}, nil
}

// Delete operation marks the specified share for deletion. The share and any files contained within it are later deleted during garbage collection.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/delete-share.
func (s *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	return DeleteResponse{}, nil
}

// Restore operation restores a share that had previously been soft-deleted.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/restore-share.
func (s *Client) Restore(ctx context.Context, deletedShareVersion string, options *RestoreOptions) (RestoreResponse, error) {
	return RestoreResponse{}, nil
}

// GetProperties operation returns all user-defined metadata and system properties for the specified share or share snapshot.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-share-properties.
func (s *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	return GetPropertiesResponse{}, nil
}

// SetProperties operation sets properties for the specified share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-share-properties.
func (s *Client) SetProperties(ctx context.Context, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	return SetPropertiesResponse{}, nil
}

// CreateSnapshot operation creates a read-only snapshot of a share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/snapshot-share.
func (s *Client) CreateSnapshot(ctx context.Context, options *CreateSnapshotOptions) (CreateSnapshotResponse, error) {
	return CreateSnapshotResponse{}, nil
}

// GetAccessPolicy operation returns information about stored access policies specified on the share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-share-acl.
func (s *Client) GetAccessPolicy(ctx context.Context, o *GetAccessPolicyOptions) (GetAccessPolicyResponse, error) {
	return GetAccessPolicyResponse{}, nil
}

// SetAccessPolicy operation sets a stored access policy for use with shared access signatures.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-share-acl.
func (s *Client) SetAccessPolicy(ctx context.Context, o *SetAccessPolicyOptions) (SetAccessPolicyResponse, error) {
	return SetAccessPolicyResponse{}, nil
}

// CreatePermission operation creates a permission (a security descriptor) at the share level.
// The created security descriptor can be used for the files and directories in the share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/create-permission.
func (s *Client) CreatePermission(ctx context.Context, sharePermission string, o *CreatePermissionOptions) (CreatePermissionResponse, error) {
	return CreatePermissionResponse{}, nil
}

// GetPermission operation gets the SDDL permission string from the service using a known permission key.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-permission.
func (s *Client) GetPermission(ctx context.Context, filePermissionKey string, o *GetPermissionOptions) (GetPermissionResponse, error) {
	return GetPermissionResponse{}, nil
}

// SetMetadata operation sets one or more user-defined name-value pairs for the specified share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-share-metadata.
func (s *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	return SetMetadataResponse{}, nil

}

// GetStatistics operation retrieves statistics related to the share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-share-stats.
func (s *Client) GetStatistics(ctx context.Context, options *GetStatisticsOptions) (GetStatisticsResponse, error) {
	return GetStatisticsResponse{}, nil
}

// TODO: should the lease methods be part of a new lease client like azblob?

// AcquireLease operation can be used to request a new lease.
// The lease duration must be between 15 and 60 seconds, or infinite (-1).
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *Client) AcquireLease(ctx context.Context, duration int32, options *AcquireLeaseOptions) (AcquireLeaseResponse, error) {
	// TODO: update generated code to make duration as required parameter
	return AcquireLeaseResponse{}, nil
}

// BreakLease operation can be used to break the lease, if the file share has an active lease. Once a lease is broken, it cannot be renewed.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *Client) BreakLease(ctx context.Context, options *BreakLeaseOptions) (BreakLeaseResponse, error) {
	return BreakLeaseResponse{}, nil
}

// ChangeLease operation can be used to change the lease ID of an active lease.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *Client) ChangeLease(ctx context.Context, leaseID string, proposedLeaseID string, options *ChangeLeaseOptions) (ChangeLeaseResponse, error) {
	// TODO: update generated code to make proposedLeaseID as required parameter
	return ChangeLeaseResponse{}, nil
}

// ReleaseLease operation can be used to free the lease if it is no longer needed so that another client may immediately acquire a lease against the file share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *Client) ReleaseLease(ctx context.Context, leaseID string, options *ReleaseLeaseOptions) (ReleaseLeaseResponse, error) {
	return ReleaseLeaseResponse{}, nil
}

// RenewLease operation can be used to renew an existing lease.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *Client) RenewLease(ctx context.Context, leaseID string, options *RenewLeaseOptions) (RenewLeaseResponse, error) {
	return RenewLeaseResponse{}, nil
}
