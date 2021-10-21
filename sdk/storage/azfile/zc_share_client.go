package azfile

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"net/url"
)

type ShareClient struct {
	client *shareClient
	u      url.URL
	cred   azcore.Credential
}

func NewShareClient(shareURL string, cred azcore.Credential, options *ClientOptions) (ShareClient, error) {
	u, err := url.Parse(shareURL)
	if err != nil {
		return ShareClient{}, err
	}
	con := newConnection(shareURL, cred, options.getConnectionOptions())
	return ShareClient{client: &shareClient{con: con}, u: *u, cred: cred}, nil
}

func (s ShareClient) URL() string {
	return s.u.String()
}

// WithSnapshot creates a new ShareURL object identical to the source but with the specified snapshot timestamp.
// Pass time.Time{} to remove the snapshot returning a URL to the base share.
func (s ShareClient) WithSnapshot(shareSnapshot string) ShareClient {
	shareURLParts := NewFileURLParts(s.URL())
	shareURLParts.ShareSnapshot = shareSnapshot
	u, _ := url.Parse(shareURLParts.URL())

	return ShareClient{
		client: &shareClient{
			&connection{u: shareURLParts.URL(), p: s.client.con.p},
		},
		u:    *u,
		cred: s.cred,
	}
}

// NewRootDirectoryClient creates a new DirectoryClient object using ShareClient's URL.
// The new DirectoryClient uses the same request policy pipeline as the ShareClient.
func (s ShareClient) NewRootDirectoryClient() DirectoryClient {
	directoryURL := s.URL()
	conn := &connection{directoryURL, s.client.con.p}
	return DirectoryClient{
		client: &directoryClient{
			con: conn,
		},
		cred: s.cred,
		u:    s.u,
	}
}

// NewDirectoryClient creates a new DirectoryURL object by concatenating directoryName to the end of
// ShareURL's URL. The new DirectoryURL uses the same request policy pipeline as the ShareURL.
// To change the pipeline, create the DirectoryURL and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewDirectoryURL instead of calling this object's
// NewDirectoryURL method.
func (s ShareClient) NewDirectoryClient(directoryName string) DirectoryClient {
	directoryURL := appendToURLPath(s.URL(), directoryName)
	u, _ := url.Parse(directoryURL)
	conn := &connection{directoryURL, s.client.con.p}
	return DirectoryClient{
		client: &directoryClient{
			con: conn,
		},
		u:    *u,
		cred: s.cred,
	}
}

// Create creates a new share within a storage account. If a share with the same name already exists, the operation fails.
// quotaInGB specifies the maximum size of the share in gigabytes, 0 means you accept service's default quota.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/create-share.
func (s ShareClient) Create(ctx context.Context, options *CreateShareOptions) (ShareCreateResponse, error) {
	formattedOptions := options.format()
	shareCreateResponse, err := s.client.Create(ctx, formattedOptions)
	return shareCreateResponse, handleError(err)
}

// CreateSnapshot creates a read-only snapshot of a share.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/snapshot-share.
func (s ShareClient) CreateSnapshot(ctx context.Context, options *CreateShareSnapshotOptions) (ShareCreateSnapshotResponse, error) {
	formattedOptions := options.format()
	shareCreateSnapshotResponse, err := s.client.CreateSnapshot(ctx, formattedOptions)
	return shareCreateSnapshotResponse, handleError(err)
}

// Delete marks the specified share or share snapshot for deletion.
// The share or share snapshot and any files contained within it are later deleted during garbage collection.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-share.
func (s ShareClient) Delete(ctx context.Context, options *DeleteShareOptions) (ShareDeleteResponse, error) {
	formattedOptions, leaseAccessConditions := options.format()
	shareDeleteResponse, err := s.client.Delete(ctx, formattedOptions, leaseAccessConditions)
	return shareDeleteResponse, handleError(err)
}

// GetProperties returns all user-defined metadata and system properties for the specified share or share snapshot.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/get-share-properties.
func (s ShareClient) GetProperties(ctx context.Context, options *GetSharePropertiesOptions) (ShareGetPropertiesResponse, error) {
	formattedOptions, leaseAccessCondition := options.format()
	shareGetPropertiesResponse, err := s.client.GetProperties(ctx, formattedOptions, leaseAccessCondition)
	return shareGetPropertiesResponse, handleError(err)
}

// SetProperties sets service-defined properties for the specified share.
// quotaInGB specifies the maximum size of the share in gigabytes, 0 means no quote and uses service's default value.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/set-share-properties.
func (s ShareClient) SetProperties(ctx context.Context, options *SetSharePropertiesOptions) (ShareSetPropertiesResponse, error) {
	formattedOptions, leaseAccessConditions, err := options.format()
	if err != nil {
		return ShareSetPropertiesResponse{}, handleError(err)
	}
	shareSetPropertiesResponse, err := s.client.SetProperties(ctx, formattedOptions, leaseAccessConditions)
	return shareSetPropertiesResponse, handleError(err)
}

// SetMetadata sets the share's metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-share-metadata.
func (s ShareClient) SetMetadata(ctx context.Context, metadata map[string]string, options *SetShareMetadataOptions) (ShareSetMetadataResponse, error) {
	formattedOptions, leaseAccessConditions, err := options.format(metadata)
	if err != nil {
		return ShareSetMetadataResponse{}, err
	}
	shareSetMetadataResponse, err := s.client.SetMetadata(ctx, formattedOptions, leaseAccessConditions)
	return shareSetMetadataResponse, handleError(err)
}

// GetPermissions returns information about stored access policies specified on the share.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-share-acl.
func (s ShareClient) GetPermissions(ctx context.Context, options *GetShareAccessPolicyOptions) (ShareGetAccessPolicyResponse, error) {
	leaseAccessConditions := options.format()
	shareGetAccessPolicyResponse, err := s.client.GetAccessPolicy(ctx, nil, leaseAccessConditions)
	return shareGetAccessPolicyResponse, handleError(err)
}

// CreatePermission uploads a SDDL permission string, and returns a permission key to use in conjunction with a file or folder.
// Note that this is only required for 9KB or larger permission strings.
// Furthermore, note that SDDL strings should be converted to a portable format before being uploaded.
// In order to make a SDDL portable, please replace well-known SIDs with their domain specific counterpart.
// Well-known SIDs are listed here: https://docs.microsoft.com/en-us/windows/win32/secauthz/sid-strings
// More info about SDDL strings can be located at: https://docs.microsoft.com/en-us/windows/win32/secauthz/security-descriptor-string-format
func (s ShareClient) CreatePermission(ctx context.Context, sharePermission string, options *CreateSharePermissionOptions) (ShareCreatePermissionResponse, error) {
	permission := SharePermission{Permission: to.StringPtr(sharePermission)}
	shareCreatePermissionResponse, err := s.client.CreatePermission(ctx, permission, nil)
	return shareCreatePermissionResponse, handleError(err)
}

// GetPermission obtains a SDDL permission string from the service using a known permission key.
func (s ShareClient) GetPermission(ctx context.Context, filePermissionKey string, options *GetSharePermissionOptions) (ShareGetPermissionResponse, error) {
	shareGetPermissionResponse, err := s.client.GetPermission(ctx, filePermissionKey, nil)
	return shareGetPermissionResponse, handleError(err)
}

// SetPermissions sets a stored access policy for use with shared access signatures.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-share-acl.
func (s ShareClient) SetPermissions(ctx context.Context, shareACLs []*SignedIdentifier, options *SetShareAccessPolicyOptions) (ShareSetAccessPolicyResponse, error) {
	shareSetAccessPolicyOptions, leaseAccessConditions := options.format(shareACLs)
	shareSetAccessPolicyResponse, err := s.client.SetAccessPolicy(ctx, shareSetAccessPolicyOptions, leaseAccessConditions)
	return shareSetAccessPolicyResponse, handleError(err)
}

// GetStatistics retrieves statistics related to the share.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/get-share-stats.
func (s ShareClient) GetStatistics(ctx context.Context, options *GetShareStatisticsOptions) (ShareGetStatisticsResponse, error) {
	formattedOptions, leaseAccessConditions := options.format()
	shareGetStatisticsResponse, err := s.client.GetStatistics(ctx, formattedOptions, leaseAccessConditions)
	return shareGetStatisticsResponse, handleError(err)
}
