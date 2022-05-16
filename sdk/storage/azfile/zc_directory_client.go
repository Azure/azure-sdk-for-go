//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
)

// A DirectoryClient represents a URL to the Azure Storage directory allowing you to manipulate its directories and files.
type DirectoryClient struct {
	client    *directoryClient
	sharedKey *SharedKeyCredential
}

// URL returns the URL endpoint used by the ServiceClient object.
func (d DirectoryClient) URL() string {
	return d.client.endpoint
}

// NewDirectoryClient creates a DirectoryClient object using the specified URL, Azure AD credential, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewDirectoryClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*DirectoryClient, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{tokenScope}, nil)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := newConnection(serviceURL, conOptions)

	return &DirectoryClient{
		client: newDirectoryClient(conn.Endpoint(), conn.Pipeline()),
	}, nil
}

// NewDirectoryClientWithNoCredential creates a DirectoryClient object using the specified URL and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net?<SAS token>
func NewDirectoryClientWithNoCredential(serviceURL string, options *ClientOptions) (*DirectoryClient, error) {
	conOptions := getConnectionOptions(options)
	conn := newConnection(serviceURL, conOptions)

	return &DirectoryClient{
		client: newDirectoryClient(conn.Endpoint(), conn.Pipeline()),
	}, nil
}

// NewDirectoryClientWithSharedKey creates a DirectoryClient object using the specified URL, shared key, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewDirectoryClientWithSharedKey(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*DirectoryClient, error) {
	authPolicy := newSharedKeyCredPolicy(cred)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := newConnection(serviceURL, conOptions)

	return &DirectoryClient{
		client:    newDirectoryClient(conn.Endpoint(), conn.Pipeline()),
		sharedKey: cred,
	}, nil
}

// NewDirectoryClientFromConnectionString creates a DirectoryClient from the given connection string.
//nolint
func NewDirectoryClientFromConnectionString(connectionString string, options *ClientOptions) (*DirectoryClient, error) {
	endpoint, credential, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	return NewDirectoryClientWithSharedKey(endpoint, credential, options)
}

// NewFileClient creates a new FileURL object by concatenating fileName to the end of
// DirectoryClient's URL. The new FileURL uses the same request policy pipeline as the DirectoryClient.
// To change the pipeline, create the FileURL and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewFileURL instead of calling this object's
// NewFileURL method.
func (d *DirectoryClient) NewFileClient(fileName string) (*FileClient, error) {
	blobURL := appendToURLPath(d.URL(), fileName)

	return &FileClient{
		client:    newFileClient(blobURL, d.client.pl),
		sharedKey: d.sharedKey,
	}, nil
}

func (d *DirectoryClient) NewDirectoryClient(directoryName string) (*DirectoryClient, error) {
	directoryURL := appendToURLPath(d.URL(), directoryName)

	return &DirectoryClient{
		client:    newDirectoryClient(directoryURL, d.client.pl),
		sharedKey: d.sharedKey,
	}, nil
}

// Create creates a new directory within a storage account.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/create-directory.
// Pass default values for SMB properties (ex: "None" for file attributes).
// If permissions is empty, the default permission "inherit" is used.
// For SDDL strings over 9KB, upload using ShareURL.CreatePermission, and supply the permissionKey.
func (d DirectoryClient) Create(ctx context.Context, options *DirectoryCreateOptions) (DirectoryCreateResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, createOptions, err := options.format()
	if err != nil {
		return DirectoryCreateResponse{}, err
	}

	directoryCreateResponse, err := d.client.Create(ctx, fileAttributes, fileCreationTime, fileLastWriteTime, createOptions)
	return toDirectoryCreateResponse(directoryCreateResponse), handleError(err)
}

// Delete removes the specified empty directory. Note that the directory must be empty before it can be deleted..
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-directory.
func (d DirectoryClient) Delete(ctx context.Context, options *DirectoryDeleteOptions) (DirectoryDeleteResponse, error) {
	directoryDeleteOptions := options.format()
	directoryDeleteResponse, err := d.client.Delete(ctx, directoryDeleteOptions)
	return toDirectoryDeleteResponse(directoryDeleteResponse), handleError(err)
}

// GetProperties returns the directory's metadata and system properties.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/get-directory-properties.
func (d DirectoryClient) GetProperties(ctx context.Context, options *DirectoryGetPropertiesOptions) (DirectoryGetPropertiesResponse, error) {
	directoryGetPropertiesOptions := options.format()
	directoryGetPropertiesResponse, err := d.client.GetProperties(ctx, directoryGetPropertiesOptions)
	return toDirectoryGetPropertiesResponse(directoryGetPropertiesResponse), handleError(err)
}

// SetProperties sets the directory's metadata and system properties.
// Preserve values for SMB properties.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/set-directory-properties.
func (d DirectoryClient) SetProperties(ctx context.Context, options *DirectorySetPropertiesOptions) (DirectorySetPropertiesResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, directorySetPropertiesOptions := options.format()

	directorySetPropertiesResponse, err := d.client.SetProperties(ctx, fileAttributes, fileCreationTime, fileLastWriteTime, directorySetPropertiesOptions)
	return toDirectorySetPropertiesResponse(directorySetPropertiesResponse), handleError(err)
}

// SetMetadata sets the directory's metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-directory-metadata.
func (d DirectoryClient) SetMetadata(ctx context.Context, metadata map[string]string, options *DirectorySetMetadataOptions) (DirectorySetMetadataResponse, error) {
	formattedOptions, err := options.format(metadata)
	if err != nil {
		return DirectorySetMetadataResponse{}, err
	}

	directorySetMetadataResponse, err := d.client.SetMetadata(ctx, formattedOptions)
	return toDirectorySetMetadataResponse(directorySetMetadataResponse), handleError(err)
}

// ListFilesAndDirectories returns a single segment of files and directories starting from the specified Marker.
// Use an empty Marker to start enumeration from the beginning. File and directory names are returned in lexicographic order.
// After getting a segment, process it, and then call ListFilesAndDirectoriesSegment again (passing the the previously-returned
// Marker) to get the next segment. This method lists the contents only for a single level of the directory hierarchy.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/list-directories-and-files.
func (d DirectoryClient) ListFilesAndDirectories(options *DirectoryListFilesAndDirectoriesOptions) *runtime.Pager[DirectoryListFilesAndDirectoriesResponse] {
	formattedOptions := options.format()
	return runtime.NewPager(runtime.PagingHandler[DirectoryListFilesAndDirectoriesResponse]{
		More: func(page DirectoryListFilesAndDirectoriesResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *DirectoryListFilesAndDirectoriesResponse) (DirectoryListFilesAndDirectoriesResponse, error) {
			if page.ListFilesAndDirectoriesSegmentResponse.NextMarker == nil || len(*page.ListFilesAndDirectoriesSegmentResponse.NextMarker) == 0 {
				return DirectoryListFilesAndDirectoriesResponse{}, handleError(errors.New("unexpected missing NextMarker"))
			}

			req, err := d.client.listFilesAndDirectoriesSegmentCreateRequest(ctx, formattedOptions)
			if err != nil {
				return DirectoryListFilesAndDirectoriesResponse{}, handleError(err)
			}
			queryValues, err := url.ParseQuery(req.Raw().URL.RawQuery)
			if err != nil {
				return DirectoryListFilesAndDirectoriesResponse{}, handleError(err)
			}
			queryValues.Set("marker", *page.ListFilesAndDirectoriesSegmentResponse.NextMarker)
			req.Raw().URL.RawQuery = queryValues.Encode()
			if err != nil {
				return DirectoryListFilesAndDirectoriesResponse{}, handleError(err)
			}

			resp, err := d.client.pl.Do(req)
			if err != nil {
				return DirectoryListFilesAndDirectoriesResponse{}, handleError(err)
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return DirectoryListFilesAndDirectoriesResponse{}, handleError(runtime.NewResponseError(resp))
			}
			generatedResp, err := d.client.listFilesAndDirectoriesSegmentHandleResponse(resp)
			return toDirectoryListFilesAndDirectoriesResponse(generatedResp), handleError(err)
		},
	})
}
