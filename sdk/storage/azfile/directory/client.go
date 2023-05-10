//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"net/http"
	"net/url"
	"strings"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Storage directory allowing you to manipulate its directories and files.
type Client base.Client[generated.DirectoryClient]

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a directory or with a shared access signature (SAS) token.
//   - directoryURL - the URL of the directory e.g. https://<account>.file.core.windows.net/share/directory?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(directoryURL string, options *ClientOptions) (*Client, error) {
	conOptions := shared.GetClientOptions(options)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*Client)(base.NewDirectoryClient(directoryURL, pl, nil)), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - directoryURL - the URL of the directory e.g. https://<account>.file.core.windows.net/share/directory
//   - cred - a SharedKeyCredential created with the matching directory's storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(directoryURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*Client)(base.NewDirectoryClient(directoryURL, pl, cred)), nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - shareName - the name of the share within the storage account
//   - directoryPath - the path of the directory within the share
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, shareName string, directoryPath string, options *ClientOptions) (*Client, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	directoryPath = strings.ReplaceAll(directoryPath, "\\", "/")
	parsed.ServiceURL = runtime.JoinPaths(parsed.ServiceURL, shareName, directoryPath)

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewClientWithNoCredential(parsed.ServiceURL, options)
}

func (d *Client) generated() *generated.DirectoryClient {
	return base.InnerClient((*base.Client[generated.DirectoryClient])(d))
}

func (d *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.DirectoryClient])(d))
}

// URL returns the URL endpoint used by the Client object.
func (d *Client) URL() string {
	return d.generated().Endpoint()
}

// NewSubdirectoryClient creates a new Client object by concatenating subDirectoryName to the end of this Client's URL.
// The new subdirectory Client uses the same request policy pipeline as the parent directory Client.
func (d *Client) NewSubdirectoryClient(subDirectoryName string) *Client {
	subDirectoryName = url.PathEscape(strings.TrimRight(subDirectoryName, "/"))
	subDirectoryURL := runtime.JoinPaths(d.URL(), subDirectoryName)
	return (*Client)(base.NewDirectoryClient(subDirectoryURL, d.generated().Pipeline(), d.sharedKey()))
}

// NewFileClient creates a new file.Client object by concatenating fileName to the end of this Client's URL.
// The new file.Client uses the same request policy pipeline as the Client.
func (d *Client) NewFileClient(fileName string) *file.Client {
	fileName = url.PathEscape(fileName)
	fileURL := runtime.JoinPaths(d.URL(), fileName)
	return (*file.Client)(base.NewFileClient(fileURL, d.generated().Pipeline(), d.sharedKey()))
}

// Create operation creates a new directory under the specified share or parent directory.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/create-directory.
func (d *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, opts := options.format()
	resp, err := d.generated().Create(ctx, fileAttributes, fileCreationTime, fileLastWriteTime, opts)
	return resp, err
}

// Delete operation removes the specified empty directory. Note that the directory must be empty before it can be deleted.
// Deleting directories that aren't empty returns error 409 (Directory Not Empty).
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/delete-directory.
func (d *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	opts := options.format()
	resp, err := d.generated().Delete(ctx, opts)
	return resp, err
}

// GetProperties operation returns all system properties for the specified directory, and it can also be used to check the existence of a directory.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-directory-properties.
func (d *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts := options.format()
	resp, err := d.generated().GetProperties(ctx, opts)
	return resp, err
}

// SetProperties operation sets system properties for the specified directory.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-directory-properties.
func (d *Client) SetProperties(ctx context.Context, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, opts := options.format()
	resp, err := d.generated().SetProperties(ctx, fileAttributes, fileCreationTime, fileLastWriteTime, opts)
	return resp, err
}

// SetMetadata operation sets user-defined metadata for the specified directory.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-directory-metadata.
func (d *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	opts := options.format()
	resp, err := d.generated().SetMetadata(ctx, opts)
	return resp, err
}

// ForceCloseHandles operation closes a handle or handles opened on a directory.
//   - handleID - Specifies the handle ID to be closed. Use an asterisk (*) as a wildcard string to specify all handles.
//
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/force-close-handles.
func (d *Client) ForceCloseHandles(ctx context.Context, handleID string, options *ForceCloseHandlesOptions) (ForceCloseHandlesResponse, error) {
	opts := options.format()
	resp, err := d.generated().ForceCloseHandles(ctx, handleID, opts)
	return resp, err
}

// ListHandles operation returns a list of open handles on a directory.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-handles.
func (d *Client) ListHandles(ctx context.Context, options *ListHandlesOptions) (ListHandlesResponse, error) {
	opts := options.format()
	resp, err := d.generated().ListHandles(ctx, opts)
	return resp, err
}

// NewListFilesAndDirectoriesPager operation returns a pager for the files and directories starting from the specified Marker.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-directories-and-files.
func (d *Client) NewListFilesAndDirectoriesPager(options *ListFilesAndDirectoriesOptions) *runtime.Pager[ListFilesAndDirectoriesResponse] {
	listOptions := generated.DirectoryClientListFilesAndDirectoriesSegmentOptions{}
	if options != nil {
		listOptions.Include = options.Include.format()
		listOptions.IncludeExtendedInfo = options.IncludeExtendedInfo
		listOptions.Marker = options.Marker
		listOptions.Maxresults = options.MaxResults
		listOptions.Prefix = options.Prefix
		listOptions.Sharesnapshot = options.ShareSnapshot
	}

	return runtime.NewPager(runtime.PagingHandler[ListFilesAndDirectoriesResponse]{
		More: func(page ListFilesAndDirectoriesResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *ListFilesAndDirectoriesResponse) (ListFilesAndDirectoriesResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = d.generated().ListFilesAndDirectoriesSegmentCreateRequest(ctx, &listOptions)
			} else {
				listOptions.Marker = page.NextMarker
				req, err = d.generated().ListFilesAndDirectoriesSegmentCreateRequest(ctx, &listOptions)
			}
			if err != nil {
				return ListFilesAndDirectoriesResponse{}, err
			}
			resp, err := d.generated().Pipeline().Do(req)
			if err != nil {
				return ListFilesAndDirectoriesResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListFilesAndDirectoriesResponse{}, runtime.NewResponseError(resp)
			}
			return d.generated().ListFilesAndDirectoriesSegmentHandleResponse(resp)
		},
	})
}
