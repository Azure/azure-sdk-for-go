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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
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
type Client base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client]

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
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

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
	//base.SharedKeyComposite((*base.CompositeClient[generated.BlobClient, generated.BlockBlobClient])(bb))
	dirClientWithDFS, _, _ := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(d))
	return dirClientWithDFS
}

func (d *Client) generatedDirClientWithBlob() *generated.PathClient {
	_, dirClientWithBlob, _ := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(d))
	return dirClientWithBlob
}

func (d *Client) blobClient() *blockblob.Client {
	_, _, blobClient := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(d))
	return blobClient
}

func (d *Client) getClientOptions() *base.ClientOptions {
	return base.GetCompositeClientOptions((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(d))
}

func (d *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKeyComposite((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(d))
}

func (d *Client) identityCredential() *azcore.TokenCredential {
	return base.IdentityCredentialComposite((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(d))
}

// DFSURL returns the URL endpoint used by the Client object.
func (d *Client) DFSURL() string {
	return d.generatedDirClientWithDFS().Endpoint()
}

// BlobURL returns the URL endpoint used by the Client object.
func (d *Client) BlobURL() string {
	return d.generatedDirClientWithBlob().Endpoint()
}

//TODO: create method to get file client - this will require block blob to have a method to get another block blob

// Create creates a new file (dfs1).
func (d *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	lac, mac, httpHeaders, createOpts, cpkOpts := options.format()
	resp, err := d.generatedDirClientWithDFS().Create(ctx, createOpts, httpHeaders, lac, mac, nil, cpkOpts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// Delete deletes directory and any path under it (dfs1).
func (d *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	lac, mac, deleteOpts := path.FormatDeleteOptions(options, true)
	resp, err := d.generatedDirClientWithDFS().Delete(ctx, deleteOpts, lac, mac)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// GetProperties gets the properties of a file (blob3)
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
	newPathURL, newBlobURL := shared.GetURLs(runtime.JoinPaths(firstPart, newName))
	parsedNewURL, _ := url.Parse(d.DFSURL())
	return parsedNewURL.Path, newPathURL, newBlobURL
}

// Rename renames a file (dfs1)
func (d *Client) Rename(ctx context.Context, newName string, options *RenameOptions) (RenameResponse, error) {
	newPathWithoutURL, newBlobURL, newPathURL := d.renamePathInURL(newName)
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
		return RenameResponse{}, err
	}
	newDirClient := (*Client)(base.NewPathClient(newPathURL, newBlobURL, newBlobClient, d.generatedDirClientWithDFS().InternalClient().WithClientName(shared.DirectoryClient), d.sharedKey(), d.identityCredential(), d.getClientOptions()))
	resp, err := newDirClient.generatedDirClientWithDFS().Create(ctx, createOpts, nil, lac, mac, smac, nil)
	return RenameResponse{
		Response:           resp,
		NewDirectoryClient: newDirClient,
	}, exported.ConvertToDFSError(err)
}

// SetAccessControl sets the owner, owning group, and permissions for a file or directory (dfs1).
func (d *Client) SetAccessControl(ctx context.Context, options *SetAccessControlOptions) (SetAccessControlResponse, error) {
	opts, lac, mac, err := path.FormatSetAccessControlOptions(options)
	if err != nil {
		return SetAccessControlResponse{}, err
	}
	resp, err := d.generatedDirClientWithDFS().SetAccessControl(ctx, opts, lac, mac)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

func (d *Client) setAccessControlHelper(mode generated.PathSetAccessControlRecursiveMode, listOptions *generated.PathClientSetAccessControlRecursiveOptions) *runtime.Pager[SetAccessControlRecursiveResponse] {
	return runtime.NewPager(runtime.PagingHandler[SetAccessControlRecursiveResponse]{
		More: func(page SetAccessControlRecursiveResponse) bool {
			return page.Continuation != nil && len(*page.Continuation) > 0
		},
		Fetcher: func(ctx context.Context, page *SetAccessControlRecursiveResponse) (SetAccessControlRecursiveResponse, error) {
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
				return SetAccessControlRecursiveResponse{}, err
			}
			resp, err := d.generatedDirClientWithDFS().InternalClient().Pipeline().Do(req)
			err = exported.ConvertToDFSError(err)
			if err != nil {
				return SetAccessControlRecursiveResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SetAccessControlRecursiveResponse{}, runtime.NewResponseError(resp)
			}
			newResp, err := d.generatedDirClientWithDFS().SetAccessControlRecursiveHandleResponse(resp)
			return newResp, exported.ConvertToDFSError(err)
		},
	})

}

// NewSetAccessControlRecursivePager sets the owner, owning group, and permissions for a file or directory (dfs1).
func (d *Client) NewSetAccessControlRecursivePager(ACL string, options *SetAccessControlRecursiveOptions) *runtime.Pager[SetAccessControlRecursiveResponse] {
	mode, listOptions := options.format(ACL, "set")
	return d.setAccessControlHelper(mode, listOptions)
}

// NewUpdateAccessControlRecursivePager updates the owner, owning group, and permissions for a file or directory (dfs1).
func (d *Client) NewUpdateAccessControlRecursivePager(ACL string, options *UpdateAccessControlRecursiveOptions) *runtime.Pager[UpdateAccessControlRecursiveResponse] {
	mode, listOptions := options.format(ACL, "modify")
	return d.setAccessControlHelper(mode, listOptions)
}

// NewRemoveAccessControlRecursivePager removes the owner, owning group, and permissions for a file or directory (dfs1).
func (d *Client) NewRemoveAccessControlRecursivePager(ACL string, options *RemoveAccessControlRecursiveOptions) *runtime.Pager[RemoveAccessControlRecursiveResponse] {
	mode, listOptions := options.format(ACL, "remove")
	return d.setAccessControlHelper(mode, listOptions)
}

// GetAccessControl gets the owner, owning group, and permissions for a file or directory (dfs1).
func (d *Client) GetAccessControl(ctx context.Context, options *GetAccessControlOptions) (GetAccessControlResponse, error) {
	opts, lac, mac := path.FormatGetAccessControlOptions(options)
	resp, err := d.generatedDirClientWithDFS().GetProperties(ctx, opts, lac, mac)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// SetMetadata sets the metadata for a file or directory (blob3).
func (d *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	opts, metadata := path.FormatSetMetadataOptions(options)
	resp, err := d.blobClient().SetMetadata(ctx, metadata, opts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// SetHTTPHeaders sets the HTTP headers for a file or directory (blob3).
func (d *Client) SetHTTPHeaders(ctx context.Context, httpHeaders HTTPHeaders, options *SetHTTPHeadersOptions) (SetHTTPHeadersResponse, error) {
	opts, blobHTTPHeaders := path.FormatSetHTTPHeadersOptions(options, httpHeaders)
	resp, err := d.blobClient().SetHTTPHeaders(ctx, blobHTTPHeaders, opts)
	newResp := SetHTTPHeadersResponse{}
	path.FormatSetHTTPHeadersResponse(&newResp, &resp)
	err = exported.ConvertToDFSError(err)
	return newResp, err
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at blob.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (f *Client) GetSASURL(permissions sas.DirectoryPermissions, expiry time.Time, o *GetSASURLOptions) (string, error) {
	if f.sharedKey() == nil {
		return "", datalakeerror.MissingSharedKeyCredential
	}

	urlParts, err := sas.ParseURL(f.BlobURL())
	err = exported.ConvertToDFSError(err)
	if err != nil {
		return "", err
	}

	st := path.FormatGetSASURLOptions(o)

	qps, err := sas.DatalakeSignatureValues{
		DirectoryPath:  urlParts.PathName,
		FilesystemName: urlParts.FilesystemName,
		Version:        sas.Version,
		Permissions:    permissions.String(),
		StartTime:      st,
		ExpiryTime:     expiry.UTC(),
	}.SignWithSharedKey(f.sharedKey())

	err = exported.ConvertToDFSError(err)
	if err != nil {
		return "", err
	}

	endpoint := f.BlobURL() + "?" + qps.Encode()

	return endpoint, nil
}
