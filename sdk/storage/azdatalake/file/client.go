//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

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
	"net/url"
	"strings"
	"time"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Datalake Storage service.
type Client base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client]

// NewClient creates an instance of Client with the specified values.
//   - fileURL - the URL of the blob e.g. https://<account>.dfs.core.windows.net/fs/file.txt
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(fileURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	blobURL, fileURL := shared.GetURLs(fileURL)
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.FileClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
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
	fileClient := base.NewPathClient(fileURL, blobURL, blobClient, azClient, nil, &cred, (*base.ClientOptions)(conOptions))

	return (*Client)(fileClient), nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - fileURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/fs/file.txt?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(fileURL string, options *ClientOptions) (*Client, error) {
	blobURL, fileURL := shared.GetURLs(fileURL)

	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.FileClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
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
	fileClient := base.NewPathClient(fileURL, blobURL, blobClient, azClient, nil, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(fileClient), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - fileURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/fs/file.txt
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(fileURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	blobURL, fileURL := shared.GetURLs(fileURL)

	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.FileClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
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
	fileClient := base.NewPathClient(fileURL, blobURL, blobClient, azClient, cred, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(fileClient), nil
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

func (f *Client) generatedFileClientWithDFS() *generated.PathClient {
	//base.SharedKeyComposite((*base.CompositeClient[generated.BlobClient, generated.BlockBlobClient])(bb))
	dirClientWithDFS, _, _ := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(f))
	return dirClientWithDFS
}

func (f *Client) generatedFileClientWithBlob() *generated.PathClient {
	_, dirClientWithBlob, _ := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(f))
	return dirClientWithBlob
}

func (f *Client) blobClient() *blockblob.Client {
	_, _, blobClient := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(f))
	return blobClient
}

func (f *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKeyComposite((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(f))
}

func (f *Client) identityCredential() *azcore.TokenCredential {
	return base.IdentityCredentialComposite((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(f))
}

func (f *Client) getClientOptions() *base.ClientOptions {
	return base.GetCompositeClientOptions((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(f))
}

// DFSURL returns the URL endpoint used by the Client object.
func (f *Client) DFSURL() string {
	return f.generatedFileClientWithDFS().Endpoint()
}

// BlobURL returns the URL endpoint used by the Client object.
func (f *Client) BlobURL() string {
	return f.generatedFileClientWithBlob().Endpoint()
}

// Create creates a new file (dfs1).
func (f *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	lac, mac, httpHeaders, createOpts, cpkOpts := options.format()
	resp, err := f.generatedFileClientWithDFS().Create(ctx, createOpts, httpHeaders, lac, mac, nil, cpkOpts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// Delete deletes a file (dfs1).
func (f *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	lac, mac, deleteOpts := options.format()
	resp, err := f.generatedFileClientWithDFS().Delete(ctx, deleteOpts, lac, mac)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// GetProperties gets the properties of a file (blob3)
func (f *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts := path.FormatGetPropertiesOptions(options)
	// TODO: format response + add acls, owner, group, permissions to it
	resp, err := f.blobClient().GetProperties(ctx, opts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

func (f *Client) renamePathInURL(newName string) (string, string, string) {
	endpoint := f.DFSURL()
	separator := "/"
	// Find the index of the last occurrence of the separator
	lastIndex := strings.LastIndex(endpoint, separator)
	// Split the string based on the last occurrence of the separator
	firstPart := endpoint[:lastIndex] // From the beginning of the string to the last occurrence of the separator
	newPathURL, newBlobURL := shared.GetURLs(runtime.JoinPaths(firstPart, newName))
	parsedNewURL, _ := url.Parse(f.DFSURL())
	return parsedNewURL.Path, newPathURL, newBlobURL
}

// Rename renames a file (dfs1)
func (f *Client) Rename(ctx context.Context, newName string, options *RenameOptions) (RenameResponse, error) {
	newPathWithoutURL, newBlobURL, newPathURL := f.renamePathInURL(newName)
	lac, mac, smac, createOpts := options.format(newPathWithoutURL)
	var newBlobClient *blockblob.Client
	if f.identityCredential() != nil {
		newBlobClient, _ = blockblob.NewClient(newBlobURL, *f.identityCredential(), nil)
	} else if f.sharedKey() != nil {
		blobSharedKey, _ := f.sharedKey().ConvertToBlobSharedKey()
		newBlobClient, _ = blockblob.NewClientWithSharedKeyCredential(newBlobURL, blobSharedKey, nil)
	} else {
		newBlobClient, _ = blockblob.NewClientWithNoCredential(newBlobURL, nil)
	}
	newFileClient := (*Client)(base.NewPathClient(newPathURL, newBlobURL, newBlobClient, f.generatedFileClientWithDFS().InternalClient().WithClientName(shared.FileClient), f.sharedKey(), f.identityCredential(), f.getClientOptions()))
	resp, err := newFileClient.generatedFileClientWithDFS().Create(ctx, createOpts, nil, lac, mac, smac, nil)
	return RenameResponse{
		Response:      resp,
		NewFileClient: newFileClient,
	}, exported.ConvertToDFSError(err)
}

// SetExpiry operation sets an expiry time on an existing file (blob2).
func (f *Client) SetExpiry(ctx context.Context, expiryType SetExpiryType, o *SetExpiryOptions) (SetExpiryResponse, error) {
	expMode, opts := expiryType.Format(o)
	resp, err := f.generatedFileClientWithBlob().SetExpiry(ctx, expMode, opts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

//// Upload uploads data to a file.
//func (f *Client) Upload(ctx context.Context) {
//
//}
//
//// Append appends data to a file.
//func (f *Client) Append(ctx context.Context) {
//
//}
//
//// Flush flushes previous uploaded data to a file.
//func (f *Client) Flush(ctx context.Context) {
//
//}
//
//// Download downloads data from a file.
//func (f *Client) Download(ctx context.Context) {
//
//}

// SetAccessControl sets the owner, owning group, and permissions for a file or directory (dfs1).
func (f *Client) SetAccessControl(ctx context.Context, options *SetAccessControlOptions) (SetAccessControlResponse, error) {
	opts, lac, mac, err := path.FormatSetAccessControlOptions(options)
	if err != nil {
		return SetAccessControlResponse{}, err
	}
	resp, err := f.generatedFileClientWithDFS().SetAccessControl(ctx, opts, lac, mac)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// UpdateAccessControl updates the owner, owning group, and permissions for a file or directory (dfs1).
func (f *Client) UpdateAccessControl(ctx context.Context, ACL string, options *UpdateAccessControlOptions) (UpdateAccessControlResponse, error) {
	opts, mode := options.format(ACL)
	resp, err := f.generatedFileClientWithDFS().SetAccessControlRecursive(ctx, mode, opts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// GetAccessControl gets the owner, owning group, and permissions for a file or directory (dfs1).
func (f *Client) GetAccessControl(ctx context.Context, options *GetAccessControlOptions) (GetAccessControlResponse, error) {
	opts, lac, mac := path.FormatGetAccessControlOptions(options)
	resp, err := f.generatedFileClientWithDFS().GetProperties(ctx, opts, lac, mac)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// RemoveAccessControl removes the owner, owning group, and permissions for a file or directory (dfs1).
func (f *Client) RemoveAccessControl(ctx context.Context, ACL string, options *RemoveAccessControlOptions) (RemoveAccessControlResponse, error) {
	opts, mode := options.format(ACL)
	resp, err := f.generatedFileClientWithDFS().SetAccessControlRecursive(ctx, mode, opts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// SetMetadata sets the metadata for a file or directory (blob3).
func (f *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	opts, metadata := path.FormatSetMetadataOptions(options)
	resp, err := f.blobClient().SetMetadata(ctx, metadata, opts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// SetHTTPHeaders sets the HTTP headers for a file or directory (blob3).
func (f *Client) SetHTTPHeaders(ctx context.Context, httpHeaders HTTPHeaders, options *SetHTTPHeadersOptions) (SetHTTPHeadersResponse, error) {
	opts, blobHTTPHeaders := path.FormatSetHTTPHeadersOptions(options, httpHeaders)
	resp, err := f.blobClient().SetHTTPHeaders(ctx, blobHTTPHeaders, opts)
	newResp := SetHTTPHeadersResponse{}
	path.FormatSetHTTPHeadersResponse(&newResp, &resp)
	err = exported.ConvertToDFSError(err)
	return newResp, err
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at blob.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (f *Client) GetSASURL(permissions sas.FilePermissions, expiry time.Time, o *GetSASURLOptions) (string, error) {
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
		FilePath:       urlParts.PathName,
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
