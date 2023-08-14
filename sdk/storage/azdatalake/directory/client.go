//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated_blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/path"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Datalake Storage service.
type Client base.CompositeClient[generated.PathClient, generated_blob.BlobClient, blockblob.Client]

// NewClient creates an instance of Client with the specified values.
//   - directoryURL - the URL of the directory e.g. https://<account>.dfs.core.windows.net/fs/dir
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(directoryURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	blobURL, directoryURL := shared.GetURLs(directoryURL)

	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.DirectoryClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
	}
	blobClientOpts := blockblob.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobClient, _ := blockblob.NewClient(blobURL, cred, &blobClientOpts)
	dirClient := base.NewPathClient(directoryURL, blobURL, blobClient, azClient, nil, &cred, (*base.ClientOptions)(conOptions))

	return (*Client)(dirClient), nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - directoryURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/fs/dir?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(directoryURL string, options *ClientOptions) (*Client, error) {
	blobURL, directoryURL := shared.GetURLs(directoryURL)

	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.DirectoryClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
	}
	blobClientOpts := blockblob.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobClient, _ := blockblob.NewClientWithNoCredential(blobURL, &blobClientOpts)
	dirClient := base.NewPathClient(directoryURL, blobURL, blobClient, azClient, nil, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(dirClient), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - directoryURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/fs/dir
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(directoryURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	blobURL, directoryURL := shared.GetURLs(directoryURL)

	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.DirectoryClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
	}
	blobClientOpts := blockblob.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSharedKey, err := cred.ConvertToBlobSharedKey()
	if err != nil {
		return nil, err
	}
	blobClient, _ := blockblob.NewClientWithSharedKeyCredential(blobURL, blobSharedKey, &blobClientOpts)
	dirClient := base.NewPathClient(directoryURL, blobURL, blobClient, azClient, cred, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(dirClient), nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, dirPath, fsName string, options *ClientOptions) (*Client, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	dirPath = strings.ReplaceAll(dirPath, "\\", "/")
	parsed.ServiceURL = runtime.JoinPaths(parsed.ServiceURL, fsName, dirPath)

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewClientWithNoCredential(parsed.ServiceURL, options)
}

func (d *Client) generatedDirClientWithDFS() *generated.PathClient {
	dirClientWithDFS, _, _ := base.InnerClients((*base.CompositeClient[generated.PathClient, generated_blob.BlobClient, blockblob.Client])(d))
	return dirClientWithDFS
}

func (d *Client) generatedDirClientWithBlob() *generated_blob.BlobClient {
	_, dirClientWithBlob, _ := base.InnerClients((*base.CompositeClient[generated.PathClient, generated_blob.BlobClient, blockblob.Client])(d))
	return dirClientWithBlob
}

func (d *Client) blobClient() *blockblob.Client {
	_, _, blobClient := base.InnerClients((*base.CompositeClient[generated.PathClient, generated_blob.BlobClient, blockblob.Client])(d))
	return blobClient
}

func (d *Client) getClientOptions() *base.ClientOptions {
	return base.GetCompositeClientOptions((*base.CompositeClient[generated.PathClient, generated_blob.BlobClient, blockblob.Client])(d))
}

func (d *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKeyComposite((*base.CompositeClient[generated.PathClient, generated_blob.BlobClient, blockblob.Client])(d))
}

func (d *Client) identityCredential() *azcore.TokenCredential {
	return base.IdentityCredentialComposite((*base.CompositeClient[generated.PathClient, generated_blob.BlobClient, blockblob.Client])(d))
}

// DFSURL returns the URL endpoint used by the Client object.
func (d *Client) DFSURL() string {
	return d.generatedDirClientWithDFS().Endpoint()
}

// BlobURL returns the URL endpoint used by the Client object.
func (d *Client) BlobURL() string {
	return d.generatedDirClientWithBlob().Endpoint()
}

// NewFileClient creates a new directory.Client object by concatenating directoryName to the end of this Client's URL.
// The new directory.Client uses the same request policy pipeline as the Client.
func (d *Client) NewFileClient(fileName string) (*file.Client, error) {
	fileName = url.PathEscape(fileName)
	fileURL := runtime.JoinPaths(d.DFSURL(), fileName)
	newBlobURL, fileURL := shared.GetURLs(fileURL)
	var newBlobClient *blockblob.Client
	var err error
	if d.identityCredential() != nil {
		newBlobClient, err = blockblob.NewClient(newBlobURL, *d.identityCredential(), nil)
	} else if d.sharedKey() != nil {
		blobSharedKey, _ := d.sharedKey().ConvertToBlobSharedKey()
		newBlobClient, err = blockblob.NewClientWithSharedKeyCredential(newBlobURL, blobSharedKey, nil)
	} else {
		newBlobClient, err = blockblob.NewClientWithNoCredential(newBlobURL, nil)
	}
	if err != nil {
		return nil, exported.ConvertToDFSError(err)
	}
	return (*file.Client)(base.NewPathClient(fileURL, newBlobURL, newBlobClient, d.generatedDirClientWithDFS().InternalClient().WithClientName(shared.FileClient), d.sharedKey(), d.identityCredential(), d.getClientOptions())), nil
}

// Create creates a new directory.
func (d *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	lac, mac, httpHeaders, createOpts, cpkOpts := options.format()
	resp, err := d.generatedDirClientWithDFS().Create(ctx, createOpts, httpHeaders, lac, mac, nil, cpkOpts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// Delete deletes directory and any path under it.
func (d *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	lac, mac, deleteOpts := path.FormatDeleteOptions(options, true)
	resp, err := d.generatedDirClientWithDFS().Delete(ctx, deleteOpts, lac, mac)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// GetProperties gets the properties of a directory.
func (d *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts := path.FormatGetPropertiesOptions(options)
	var respFromCtx *http.Response
	ctxWithResp := runtime.WithCaptureResponse(ctx, &respFromCtx)
	resp, err := d.blobClient().GetProperties(ctxWithResp, opts)
	newResp := path.FormatGetPropertiesResponse(&resp, respFromCtx)
	err = exported.ConvertToDFSError(err)
	return newResp, err
}

func (d *Client) renamePathInURL(newName string) (string, string, string) {
	endpoint := d.DFSURL()
	separator := "/"
	// Find the index of the last occurrence of the separator
	lastIndex := strings.LastIndex(endpoint, separator)
	// Split the string based on the last occurrence of the separator
	firstPart := endpoint[:lastIndex] // From the beginning of the string to the last occurrence of the separator
	newBlobURL, newPathURL := shared.GetURLs(runtime.JoinPaths(firstPart, newName))
	parsedNewURL, _ := url.Parse(d.DFSURL())
	return parsedNewURL.Path, newPathURL, newBlobURL
}

// Rename renames a directory.
func (d *Client) Rename(ctx context.Context, newName string, options *RenameOptions) error {
	newPathWithoutURL, newPathURL, newBlobURL := d.renamePathInURL(newName)
	lac, mac, smac, createOpts := path.FormatRenameOptions(options, newPathWithoutURL)
	var newBlobClient *blockblob.Client
	var err error
	if d.identityCredential() != nil {
		newBlobClient, err = blockblob.NewClient(newBlobURL, *d.identityCredential(), nil)
	} else if d.sharedKey() != nil {
		blobSharedKey, _ := d.sharedKey().ConvertToBlobSharedKey()
		newBlobClient, err = blockblob.NewClientWithSharedKeyCredential(newBlobURL, blobSharedKey, nil)
	} else {
		newBlobClient, err = blockblob.NewClientWithNoCredential(newBlobURL, nil)
	}
	if err != nil {
		return exported.ConvertToDFSError(err)
	}
	newDirClient := (*Client)(base.NewPathClient(newPathURL, newBlobURL, newBlobClient, d.generatedDirClientWithDFS().InternalClient().WithClientName(shared.DirectoryClient), d.sharedKey(), d.identityCredential(), d.getClientOptions()))
	_, err = newDirClient.generatedDirClientWithDFS().Create(ctx, createOpts, nil, lac, mac, smac, nil)
	//return RenameResponse{
	//	Response:           resp,
	//	NewDirectoryClient: newDirClient,
	//}, exported.ConvertToDFSError(err)
	return exported.ConvertToDFSError(err)
}

// SetAccessControl sets the owner, owning group, and permissions for a directory.
func (d *Client) SetAccessControl(ctx context.Context, options *SetAccessControlOptions) (SetAccessControlResponse, error) {
	opts, lac, mac, err := path.FormatSetAccessControlOptions(options)
	if err != nil {
		return SetAccessControlResponse{}, err
	}
	resp, err := d.generatedDirClientWithDFS().SetAccessControl(ctx, opts, lac, mac)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

func (d *Client) setAccessControlPager(mode generated.PathSetAccessControlRecursiveMode, listOptions *generated.PathClientSetAccessControlRecursiveOptions) *runtime.Pager[generated.PathClientSetAccessControlRecursiveResponse] {
	return runtime.NewPager(runtime.PagingHandler[generated.PathClientSetAccessControlRecursiveResponse]{
		More: func(page generated.PathClientSetAccessControlRecursiveResponse) bool {
			return page.Continuation != nil && len(*page.Continuation) > 0
		},
		Fetcher: func(ctx context.Context, page *generated.PathClientSetAccessControlRecursiveResponse) (generated.PathClientSetAccessControlRecursiveResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = d.generatedDirClientWithDFS().SetAccessControlRecursiveCreateRequest(ctx, mode, listOptions)
				err = exported.ConvertToDFSError(err)
			} else {
				listOptions.Continuation = page.Continuation
				req, err = d.generatedDirClientWithDFS().SetAccessControlRecursiveCreateRequest(ctx, mode, listOptions)
				err = exported.ConvertToDFSError(err)
			}
			if err != nil {
				return generated.PathClientSetAccessControlRecursiveResponse{}, err
			}
			resp, err := d.generatedDirClientWithDFS().InternalClient().Pipeline().Do(req)
			err = exported.ConvertToDFSError(err)
			if err != nil {
				return generated.PathClientSetAccessControlRecursiveResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return generated.PathClientSetAccessControlRecursiveResponse{}, runtime.NewResponseError(resp)
			}
			newResp, err := d.generatedDirClientWithDFS().SetAccessControlRecursiveHandleResponse(resp)
			return newResp, exported.ConvertToDFSError(err)
		},
	})

}

func (d *Client) setAccessControlRecursiveHelper(mode generated.PathSetAccessControlRecursiveMode, listOptions *generated.PathClientSetAccessControlRecursiveOptions, options *SetAccessControlRecursiveOptions) (SetAccessControlRecursiveResponse, error) {
	pager := d.setAccessControlPager(mode, listOptions)
	counter := *options.MaxBatches
	continueOnFailure := listOptions.ForceFlag
	totalSuccessfulDirs := int32(0)
	totalSuccessfulFiles := int32(0)
	totalFailureCount := int32(0)
	finalResponse := SetAccessControlRecursiveResponse{
		DirectoriesSuccessful: &totalSuccessfulDirs,
		FilesSuccessful:       &totalSuccessfulFiles,
		FailureCount:          &totalFailureCount,
		FailedEntries:         []*ACLFailedEntry{},
	}
	for pager.More() && counter != 0 {
		resp, err := pager.NextPage(context.Background())
		if err != nil {
			return finalResponse, exported.ConvertToDFSError(err)
		}
		finalResponse.DirectoriesSuccessful = to.Ptr(*finalResponse.DirectoriesSuccessful + *resp.DirectoriesSuccessful)
		finalResponse.FilesSuccessful = to.Ptr(*finalResponse.FilesSuccessful + *resp.FilesSuccessful)
		finalResponse.FailureCount = to.Ptr(*finalResponse.FailureCount + *resp.FailureCount)
		finalResponse.FailedEntries = append(finalResponse.FailedEntries, resp.FailedEntries...)
		counter = counter - 1
		if !*continueOnFailure && *resp.FailureCount > 0 {
			return finalResponse, exported.ConvertToDFSError(err)
		}
	}
	return finalResponse, nil
}

// SetAccessControlRecursive sets the owner, owning group, and permissions for a directory.
func (d *Client) SetAccessControlRecursive(ACL string, options *SetAccessControlRecursiveOptions) (SetAccessControlRecursiveResponse, error) {
	if options == nil {
		options = &SetAccessControlRecursiveOptions{}
	}
	mode, listOptions := options.format(ACL, "set")
	return d.setAccessControlRecursiveHelper(mode, listOptions, options)
}

// UpdateAccessControlRecursive updates the owner, owning group, and permissions for a directory.
func (d *Client) UpdateAccessControlRecursive(ACL string, options *UpdateAccessControlRecursiveOptions) (SetAccessControlRecursiveResponse, error) {
	if options == nil {
		options = &UpdateAccessControlRecursiveOptions{}
	}
	mode, listOptions := options.format(ACL, "modify")
	return d.setAccessControlRecursiveHelper(mode, listOptions, options)
}

// RemoveAccessControlRecursive removes the owner, owning group, and permissions for a directory.
func (d *Client) RemoveAccessControlRecursive(ACL string, options *RemoveAccessControlRecursiveOptions) (SetAccessControlRecursiveResponse, error) {
	if options == nil {
		options = &RemoveAccessControlRecursiveOptions{}
	}
	mode, listOptions := options.format(ACL, "remove")
	return d.setAccessControlRecursiveHelper(mode, listOptions, options)
}

// GetAccessControl gets the owner, owning group, and permissions for a directory.
func (d *Client) GetAccessControl(ctx context.Context, options *GetAccessControlOptions) (GetAccessControlResponse, error) {
	opts, lac, mac := path.FormatGetAccessControlOptions(options)
	resp, err := d.generatedDirClientWithDFS().GetProperties(ctx, opts, lac, mac)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// SetMetadata sets the metadata for a directory.
func (d *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	opts, metadata := path.FormatSetMetadataOptions(options)
	resp, err := d.blobClient().SetMetadata(ctx, metadata, opts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// SetHTTPHeaders sets the HTTP headers for a directory.
func (d *Client) SetHTTPHeaders(ctx context.Context, httpHeaders HTTPHeaders, options *SetHTTPHeadersOptions) (SetHTTPHeadersResponse, error) {
	opts, blobHTTPHeaders := path.FormatSetHTTPHeadersOptions(options, httpHeaders)
	resp, err := d.blobClient().SetHTTPHeaders(ctx, blobHTTPHeaders, opts)
	newResp := SetHTTPHeadersResponse{}
	path.FormatSetHTTPHeadersResponse(&newResp, &resp)
	err = exported.ConvertToDFSError(err)
	return newResp, err
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at directory.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (d *Client) GetSASURL(permissions sas.DirectoryPermissions, expiry time.Time, o *GetSASURLOptions) (string, error) {
	if d.sharedKey() == nil {
		return "", datalakeerror.MissingSharedKeyCredential
	}

	urlParts, err := sas.ParseURL(d.BlobURL())
	err = exported.ConvertToDFSError(err)
	if err != nil {
		return "", err
	}

	st := path.FormatGetSASURLOptions(o)

	qps, err := sas.DatalakeSignatureValues{
		DirectoryPath:  urlParts.PathName,
		FileSystemName: urlParts.FileSystemName,
		Version:        sas.Version,
		Permissions:    permissions.String(),
		StartTime:      st,
		ExpiryTime:     expiry.UTC(),
	}.SignWithSharedKey(d.sharedKey())

	err = exported.ConvertToDFSError(err)
	if err != nil {
		return "", err
	}

	endpoint := d.BlobURL() + "?" + qps.Encode()

	return endpoint, nil
}
