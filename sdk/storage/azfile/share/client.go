//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Storage share allowing you to manipulate its directories and files.
type Client base.Client[generated.ShareClient]

// NewClient creates an instance of Client with the specified values.
//   - shareURL - the URL of the share e.g. https://<account>.file.core.windows.net/share
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
//
// Note that the only share-level operations that support token credential authentication are CreatePermission and GetPermission.
// Also note that ClientOptions.FileRequestIntent is currently required for token authentication.
func NewClient(shareURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	audience := base.GetAudience((*base.ClientOptions)(options))
	conOptions := shared.GetClientOptions(options)
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{audience}, &policy.BearerTokenOptions{
		InsecureAllowCredentialWithHTTP: conOptions.InsecureAllowCredentialWithHTTP,
	})
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	return (*Client)(base.NewShareClient(shareURL, azClient, nil, (*base.ClientOptions)(conOptions))), nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a share or with a shared access signature (SAS) token.
//   - shareURL - the URL of the share e.g. https://<account>.file.core.windows.net/share?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(shareURL string, options *ClientOptions) (*Client, error) {
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	return (*Client)(base.NewShareClient(shareURL, azClient, nil, (*base.ClientOptions)(conOptions))), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - shareURL - the URL of the share e.g. https://<account>.file.core.windows.net/share
//   - cred - a SharedKeyCredential created with the matching share's storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(shareURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	return (*Client)(base.NewShareClient(shareURL, azClient, cred, (*base.ClientOptions)(conOptions))), nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - shareName - the name of the share within the storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, shareName string, options *ClientOptions) (*Client, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	parsed.ServiceURL = runtime.JoinPaths(parsed.ServiceURL, shareName)

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewClientWithNoCredential(parsed.ServiceURL, options)
}

func (s *Client) generated() *generated.ShareClient {
	return base.InnerClient((*base.Client[generated.ShareClient])(s))
}

func (s *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.ShareClient])(s))
}

func (s *Client) getClientOptions() *base.ClientOptions {
	return base.GetClientOptions((*base.Client[generated.ShareClient])(s))
}

// URL returns the URL endpoint used by the Client object.
func (s *Client) URL() string {
	return s.generated().Endpoint()
}

// NewDirectoryClient creates a new directory.Client object by concatenating directoryName to the end of this Client's URL.
// The new directory.Client uses the same request policy pipeline as the Client.
func (s *Client) NewDirectoryClient(directoryName string) *directory.Client {
	directoryName = url.PathEscape(strings.TrimRight(directoryName, "/"))
	directoryURL := runtime.JoinPaths(s.URL(), directoryName)
	return (*directory.Client)(base.NewDirectoryClient(directoryURL, s.generated().InternalClient().WithClientName(exported.ModuleName), s.sharedKey(), s.getClientOptions()))
}

// NewRootDirectoryClient creates a new directory.Client object for the root of the share using the Client's URL.
// The new directory.Client uses the same request policy pipeline as the Client.
func (s *Client) NewRootDirectoryClient() *directory.Client {
	rootDirURL := s.URL()
	return (*directory.Client)(base.NewDirectoryClient(rootDirURL, s.generated().InternalClient().WithClientName(exported.ModuleName), s.sharedKey(), s.getClientOptions()))
}

// WithSnapshot creates a new Client object identical to the source but with the specified share snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base share.
func (s *Client) WithSnapshot(shareSnapshot string) (*Client, error) {
	p, err := sas.ParseURL(s.URL())
	if err != nil {
		return nil, err
	}
	p.ShareSnapshot = shareSnapshot
	clientOptions := base.GetClientOptions((*base.Client[generated.ShareClient])(s))

	return (*Client)(base.NewShareClient(p.String(), s.generated().InternalClient(), s.sharedKey(), clientOptions)), nil
}

// Create operation creates a new share within a storage account. If a share with the same name already exists, the operation fails.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/create-share.
func (s *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	opts := options.format()
	resp, err := s.generated().Create(ctx, opts)
	return resp, err
}

// Delete operation marks the specified share for deletion. The share and any files contained within it are later deleted during garbage collection.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/delete-share.
func (s *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := s.generated().Delete(ctx, opts, leaseAccessConditions)
	return resp, err
}

// Restore operation restores a share that had previously been soft-deleted.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/restore-share.
func (s *Client) Restore(ctx context.Context, deletedShareVersion string, options *RestoreOptions) (RestoreResponse, error) {
	urlParts, err := sas.ParseURL(s.URL())
	if err != nil {
		return RestoreResponse{}, err
	}

	opts := &generated.ShareClientRestoreOptions{
		DeletedShareName:    &urlParts.ShareName,
		DeletedShareVersion: &deletedShareVersion,
	}
	resp, err := s.generated().Restore(ctx, opts)
	return resp, err
}

// GetProperties operation returns all user-defined metadata and system properties for the specified share or share snapshot.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-share-properties.
func (s *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := s.generated().GetProperties(ctx, opts, leaseAccessConditions)
	return resp, err
}

// SetProperties operation sets properties for the specified share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-share-properties.
func (s *Client) SetProperties(ctx context.Context, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := s.generated().SetProperties(ctx, opts, leaseAccessConditions)
	return resp, err
}

// CreateSnapshot operation creates a read-only snapshot of a share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/snapshot-share.
func (s *Client) CreateSnapshot(ctx context.Context, options *CreateSnapshotOptions) (CreateSnapshotResponse, error) {
	opts := options.format()
	resp, err := s.generated().CreateSnapshot(ctx, opts)
	return resp, err
}

// GetAccessPolicy operation returns information about stored access policies specified on the share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-share-acl.
func (s *Client) GetAccessPolicy(ctx context.Context, options *GetAccessPolicyOptions) (GetAccessPolicyResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := s.generated().GetAccessPolicy(ctx, opts, leaseAccessConditions)
	return resp, err
}

// SetAccessPolicy operation sets a stored access policy for use with shared access signatures.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-share-acl.
func (s *Client) SetAccessPolicy(ctx context.Context, options *SetAccessPolicyOptions) (SetAccessPolicyResponse, error) {
	opts, acl, leaseAccessConditions, err := options.format()
	if err != nil {
		return SetAccessPolicyResponse{}, err
	}

	resp, err := s.generated().SetAccessPolicy(ctx, acl, opts, leaseAccessConditions)
	return resp, err
}

// CreatePermission operation creates a permission (a security descriptor) at the share level.
// The created security descriptor can be used for the files and directories in the share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/create-permission.
func (s *Client) CreatePermission(ctx context.Context, sharePermission string, options *CreatePermissionOptions) (CreatePermissionResponse, error) {
	permission, opts := options.format(sharePermission)
	resp, err := s.generated().CreatePermission(ctx, permission, opts)
	return resp, err
}

// GetPermission operation gets the SDDL permission string from the service using a known permission key.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-permission.
func (s *Client) GetPermission(ctx context.Context, filePermissionKey string, options *GetPermissionOptions) (GetPermissionResponse, error) {
	opts := options.format()
	resp, err := s.generated().GetPermission(ctx, filePermissionKey, opts)
	return resp, err
}

// SetMetadata operation sets one or more user-defined name-value pairs for the specified share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-share-metadata.
func (s *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := s.generated().SetMetadata(ctx, opts, leaseAccessConditions)
	return resp, err
}

// GetStatistics operation retrieves statistics related to the share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-share-stats.
func (s *Client) GetStatistics(ctx context.Context, options *GetStatisticsOptions) (GetStatisticsResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := s.generated().GetStatistics(ctx, opts, leaseAccessConditions)
	return resp, err
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at share.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (s *Client) GetSASURL(permissions sas.SharePermissions, expiry time.Time, o *GetSASURLOptions) (string, error) {
	if s.sharedKey() == nil {
		return "", fileerror.MissingSharedKeyCredential
	}
	st := o.format()

	urlParts, err := sas.ParseURL(s.URL())
	if err != nil {
		return "", err
	}

	t, err := time.Parse(sas.SnapshotTimeFormat, urlParts.ShareSnapshot)
	if err != nil {
		t = time.Time{}
	}

	qps, err := sas.SignatureValues{
		Version:      sas.Version,
		ShareName:    urlParts.ShareName,
		SnapshotTime: t,
		Permissions:  permissions.String(),
		StartTime:    st,
		ExpiryTime:   expiry.UTC(),
	}.SignWithSharedKey(s.sharedKey())
	if err != nil {
		return "", err
	}

	endpoint := s.URL() + "?" + qps.Encode()

	return endpoint, nil
}
