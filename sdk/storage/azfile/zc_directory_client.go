package azfile

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/url"
)

// A DirectoryClient represents a URL to the Azure Storage directory allowing you to manipulate its directories and files.
type DirectoryClient struct {
	client *directoryClient
	u      url.URL
	cred   azcore.Credential
}

// NewDirectoryClient creates a DirectoryClient object using the specified URL and request policy pipeline.
// Note: p can't be nil.
func NewDirectoryClient(directoryURL string, cred azcore.Credential, options *ClientOptions) (DirectoryClient, error) {
	u, err := url.Parse(directoryURL)
	if err != nil {
		return DirectoryClient{}, err
	}
	return DirectoryClient{client: &directoryClient{
		con: newConnection(directoryURL, cred, options.getConnectionOptions()),
	}, u: *u, cred: cred}, nil
}

// URL returns the URL endpoint used by the DirectoryClient object.
func (d DirectoryClient) URL() string {
	return d.u.String()
}

// NewFileClient creates a new FileURL object by concatenating fileName to the end of
// DirectoryClient's URL. The new FileURL uses the same request policy pipeline as the DirectoryClient.
// To change the pipeline, create the FileURL and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewFileURL instead of calling this object's
// NewFileURL method.
func (d DirectoryClient) NewFileClient(fileName string) (FileClient, error) {
	blobURL := appendToURLPath(d.URL(), fileName)
	u, err := url.Parse(blobURL)
	if err != nil {
		return FileClient{}, err
	}
	return FileClient{
		client: &fileClient{con: &connection{u: blobURL, p: d.client.con.p}},
		u:      *u,
		cred:   d.cred,
	}, nil
}

// Create creates a new directory within a storage account.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/create-directory.
// Pass default values for SMB properties (ex: "None" for file attributes).
// If permissions is empty, the default permission "inherit" is used.
// For SDDL strings over 9KB, upload using ShareURL.CreatePermission, and supply the permissionKey.
func (d DirectoryClient) Create(ctx context.Context, options *CreateDirectoryOptions) (DirectoryCreateResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, directoryCreateOptions, err := options.format()
	if err != nil {
		return DirectoryCreateResponse{}, err
	}

	return d.client.Create(ctx, fileAttributes, fileCreationTime, fileLastWriteTime, directoryCreateOptions)
}

// Delete removes the specified empty directory. Note that the directory must be empty before it can be deleted..
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-directory.
func (d DirectoryClient) Delete(ctx context.Context) (DirectoryDeleteResponse, error) {
	return d.client.Delete(ctx, nil)
}

// GetProperties returns the directory's metadata and system properties.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/get-directory-properties.
func (d DirectoryClient) GetProperties(ctx context.Context, options *GetDirectoryPropertiesOptions) (DirectoryGetPropertiesResponse, error) {
	directoryGetPropertiesOptions := options.format()
	return d.client.GetProperties(ctx, directoryGetPropertiesOptions)
}

//// SetProperties sets the directory's metadata and system properties.
//// Preserves values for SMB properties.
//// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/set-directory-properties.
//func (d DirectoryClient) SetProperties(ctx context.Context, options *SetDirectoryPropertiesOptions) (*DirectorySetPropertiesResponse, error) {
//	permStr, permKey, fileAttr, fileCreateTime, FileLastWriteTime, err := properties.selectSMBPropertyValues(true, defaultPreserveString, defaultPreserveString, defaultPreserveString)
//
//	fileAttributes string, fileCreationTime time.Time, fileLastWriteTime time.Time, formattedOptions :=
//
//	return d.client.SetProperties(ctx, fileAttr, fileCreateTime, FileLastWriteTime, nil, permStr, permKey)
//}

// SetMetadata sets the directory's metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-directory-metadata.
func (d DirectoryClient) SetMetadata(ctx context.Context, options *SetDirectoryMetadataOptions) (DirectorySetMetadataResponse, error) {
	formattedOptions := options.format()
	return d.client.SetMetadata(ctx, formattedOptions)
}

//// ListFilesAndDirectoriesSegment returns a single segment of files and directories starting from the specified Marker.
//// Use an empty Marker to start enumeration from the beginning. File and directory names are returned in lexicographic order.
//// After getting a segment, process it, and then call ListFilesAndDirectoriesSegment again (passing the the previously-returned
//// Marker) to get the next segment. This method lists the contents only for a single level of the directory hierarchy.
//// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/list-directories-and-files.
//func (d DirectoryClient) ListFilesAndDirectoriesSegment(ctx context.Context, marker Marker, o ListFilesAndDirectoriesOptions) (*ListFilesAndDirectoriesSegmentResponse, error) {
//	prefix, maxResults := o.format()
//	return d.client.ListFilesAndDirectoriesSegment(ctx, prefix, nil, marker.Val, maxResults, nil)
//}
