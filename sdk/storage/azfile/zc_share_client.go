//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

type ShareClient struct {
	client    *shareClient
	sharedKey *SharedKeyCredential
}

func (s *ShareClient) URL() string {
	return s.client.endpoint
}

func NewShareClient(shareURL string, cred azcore.TokenCredential, options *ClientOptions) (*ShareClient, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{tokenScope}, nil)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := newConnection(shareURL, conOptions)

	return &ShareClient{
		client: newShareClient(conn.Endpoint(), conn.Pipeline()),
	}, nil
}

// NewShareClientWithNoCredential creates a SHareClient object using the specified URL and options.
func NewShareClientWithNoCredential(shareURL string, options *ClientOptions) (*ShareClient, error) {
	conOptions := getConnectionOptions(options)
	conn := newConnection(shareURL, conOptions)

	return &ShareClient{
		client: newShareClient(conn.Endpoint(), conn.Pipeline()),
	}, nil
}

// NewShareClientWithSharedKey creates a ShareClient object using the specified URL, shared key, and options.
func NewShareClientWithSharedKey(containerURL string, cred *SharedKeyCredential, options *ClientOptions) (*ShareClient, error) {
	authPolicy := newSharedKeyCredPolicy(cred)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := newConnection(containerURL, conOptions)

	return &ShareClient{
		client:    newShareClient(conn.Endpoint(), conn.Pipeline()),
		sharedKey: cred,
	}, nil
}

// NewShareClientFromConnectionString creates a service client from the given connection string.
//nolint
func NewShareClientFromConnectionString(connectionString string, shareURL string, options *ClientOptions) (*ShareClient, error) {
	_, credential, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	return NewShareClientWithSharedKey(shareURL, credential, options)
}

// WithSnapshot creates a new ShareURL object identical to the source but with the specified snapshot timestamp.
// Pass time.Time{} to remove the snapshot returning a URL to the base share.
func (s *ShareClient) WithSnapshot(shareSnapshot string) (*ShareClient, error) {
	shareURLParts, err := NewFileURLParts(s.URL())
	if err != nil {
		return nil, err
	}
	shareURLParts.ShareSnapshot = shareSnapshot
	shareURL := shareURLParts.URL()

	return &ShareClient{
		client:    newShareClient(shareURL, s.client.pl),
		sharedKey: s.sharedKey,
	}, nil
}

// NewRootDirectoryClient creates a new DirectoryClient object using ShareClient's URL.
// The new DirectoryClient uses the same request policy pipeline as the ShareClient.
func (s *ShareClient) NewRootDirectoryClient() (*DirectoryClient, error) {
	directoryURL := s.URL()
	return &DirectoryClient{
		client:    newDirectoryClient(directoryURL, s.client.pl),
		sharedKey: s.sharedKey,
	}, nil
}

// NewDirectoryClient creates a new DirectoryURL object by concatenating directoryName to the end of
// ShareURL's URL. The new DirectoryURL uses the same request policy pipeline as the ShareURL.
// To change the pipeline, create the DirectoryURL and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewDirectoryURL instead of calling this object's
// NewDirectoryURL method.
func (s *ShareClient) NewDirectoryClient(directoryName string) (*DirectoryClient, error) {
	directoryURL := appendToURLPath(s.URL(), directoryName)

	return &DirectoryClient{
		client:    newDirectoryClient(directoryURL, s.client.pl),
		sharedKey: s.sharedKey,
	}, nil
}

// Create creates a new share within a storage account. If a share with the same name already exists, the operation fails.
// quotaInGB specifies the maximum size of the share in gigabytes, 0 means you accept service's default quota.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/create-share.
func (s *ShareClient) Create(ctx context.Context, options *ShareCreateOptions) (ShareCreateResponse, error) {
	formattedOptions := options.format()
	createResponse, err := s.client.Create(ctx, formattedOptions)

	return toShareCreateResponse(createResponse), handleError(err)
}

// CreateSnapshot creates a read-only snapshot of a share.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/snapshot-share.
func (s *ShareClient) CreateSnapshot(ctx context.Context, options *ShareCreateSnapshotOptions) (ShareCreateSnapshotResponse, error) {
	formattedOptions := options.format()
	createSnapshotResponse, err := s.client.CreateSnapshot(ctx, formattedOptions)

	return toShareCreateSnapshotResponse(createSnapshotResponse), handleError(err)
}

// Delete marks the specified share or share snapshot for deletion.
// The share or share snapshot and any files contained within it are later deleted during garbage collection.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-share.
func (s *ShareClient) Delete(ctx context.Context, options *ShareDeleteOptions) (ShareDeleteResponse, error) {
	formattedOptions, leaseAccessConditions := options.format()
	deleteResponse, err := s.client.Delete(ctx, formattedOptions, leaseAccessConditions)

	return toShareDeleteResponse(deleteResponse), handleError(err)
}

// GetProperties returns all user-defined metadata and system properties for the specified share or share snapshot.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/get-share-properties.
func (s *ShareClient) GetProperties(ctx context.Context, options *ShareGetPropertiesOptions) (ShareGetPropertiesResponse, error) {
	formattedOptions, leaseAccessCondition := options.format()
	getPropertiesResponse, err := s.client.GetProperties(ctx, formattedOptions, leaseAccessCondition)

	return toShareGetPropertiesResponse(getPropertiesResponse), handleError(err)
}

// SetProperties sets service-defined properties for the specified share.
// quotaInGB specifies the maximum size of the share in gigabytes, 0 means no quote and uses service's default value.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/set-share-properties.
func (s *ShareClient) SetProperties(ctx context.Context, options *ShareSetPropertiesOptions) (ShareSetPropertiesResponse, error) {
	formattedOptions, leaseAccessConditions, err := options.format()
	if err != nil {
		return ShareSetPropertiesResponse{}, handleError(err)
	}
	setPropertiesResponse, err := s.client.SetProperties(ctx, formattedOptions, leaseAccessConditions)

	return toShareSetPropertiesResponse(setPropertiesResponse), handleError(err)
}

// SetMetadata sets the share's metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-share-metadata.
func (s *ShareClient) SetMetadata(ctx context.Context, metadata map[string]string, options *ShareSetMetadataOptions) (ShareSetMetadataResponse, error) {
	formattedOptions, leaseAccessConditions, err := options.format(metadata)
	if err != nil {
		return ShareSetMetadataResponse{}, err
	}

	setMetadataResponse, err := s.client.SetMetadata(ctx, formattedOptions, leaseAccessConditions)

	return toShareSetMetadataResponse(setMetadataResponse), handleError(err)
}

// GetPermissions returns information about stored access policies specified on the share.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-share-acl.
func (s *ShareClient) GetPermissions(ctx context.Context, o *ShareGetAccessPolicyOptions) (ShareGetAccessPolicyResponse, error) {
	formattedOptions, leaseAccessConditions := o.format()
	getAccessPolicyResponse, err := s.client.GetAccessPolicy(ctx, formattedOptions, leaseAccessConditions)

	return toShareGetAccessPolicyResponse(getAccessPolicyResponse), handleError(err)
}

// CreatePermission uploads a SDDL permission string, and returns a permission key to use in conjunction with a file or folder.
// Note that this is only required for 9KB or larger permission strings.
// Furthermore, note that SDDL strings should be converted to a portable format before being uploaded.
// In order to make a SDDL portable, please replace well-known SIDs with their domain specific counterpart.
// Well-known SIDs are listed here: https://docs.microsoft.com/en-us/windows/win32/secauthz/sid-strings
// More info about SDDL strings can be located at: https://docs.microsoft.com/en-us/windows/win32/secauthz/security-descriptor-string-format
func (s *ShareClient) CreatePermission(ctx context.Context, sharePermission string, o *ShareCreatePermissionOptions) (ShareCreatePermissionResponse, error) {
	permission, formattedOptions := o.format(sharePermission)
	createPermissionResponse, err := s.client.CreatePermission(ctx, permission, formattedOptions)

	return toShareCreatePermissionResponse(createPermissionResponse), handleError(err)
}

// GetPermission obtains a SDDL permission string from the service using a known permission key.
func (s *ShareClient) GetPermission(ctx context.Context, filePermissionKey string, o *ShareGetPermissionOptions) (ShareGetPermissionResponse, error) {
	formattedOptions := o.format()
	getPermissionResponse, err := s.client.GetPermission(ctx, filePermissionKey, formattedOptions)
	return toShareGetPermissionResponse(getPermissionResponse), handleError(err)
}

// SetPermissions sets a stored access policy for use with shared access signatures.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-share-acl.
func (s *ShareClient) SetPermissions(ctx context.Context, shareACLs []*SignedIdentifier, options *SetShareAccessPolicyOptions) (ShareSetAccessPolicyResponse, error) {
	shareSetAccessPolicyOptions, leaseAccessConditions := options.format(shareACLs)
	setAccessPolicyResponse, err := s.client.SetAccessPolicy(ctx, shareSetAccessPolicyOptions, leaseAccessConditions)

	return toShareSetAccessPolicyResponse(setAccessPolicyResponse), handleError(err)
}

// GetStatistics retrieves statistics related to the share.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/get-share-stats.
func (s *ShareClient) GetStatistics(ctx context.Context, options *ShareGetStatisticsOptions) (ShareGetStatisticsResponse, error) {
	formattedOptions, leaseAccessConditions := options.format()
	shareGetStatisticsResponse, err := s.client.GetStatistics(ctx, formattedOptions, leaseAccessConditions)

	return toShareGetStatisticsResponse(shareGetStatisticsResponse), handleError(err)
}
