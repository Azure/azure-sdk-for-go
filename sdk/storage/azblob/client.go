//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

// ClientOptions adds additional client options while constructing connection
type ClientOptions = exported.ClientOptions

// Client represents a URL to an Azure Storage blob; the blob may be a block blob, append blob, or page blob.
type Client struct {
	svc *service.Client
}

// NewClient creates a BlobClient object using the specified URL, Azure AD credential, and options.
func NewClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	svcClient, err := service.NewClient(serviceURL, cred, options)
	if err != nil {
		return nil, err
	}

	return &Client{
		svc: svcClient,
	}, nil
}

// NewClientWithNoCredential creates a BlobClient object using the specified URL and options.
func NewClientWithNoCredential(serviceURL string, options *ClientOptions) (*Client, error) {
	svcClient, err := service.NewClientWithNoCredential(serviceURL, options)
	if err != nil {
		return nil, err
	}

	return &Client{
		svc: svcClient,
	}, nil
}

// NewClientWithSharedKey creates a BlobClient object using the specified URL, shared key, and options.
func NewClientWithSharedKey(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	svcClient, err := service.NewClientWithSharedKey(serviceURL, cred, (*service.ClientOptions)(options))
	if err != nil {
		return nil, err
	}

	return &Client{
		svc: svcClient,
	}, nil
}

// NewClientFromConnectionString creates BlobClient from a connection String
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	containerClient, err := service.NewClientFromConnectionString(connectionString, (*service.ClientOptions)(options))
	if err != nil {
		return nil, err
	}
	return &Client{
		svc: containerClient,
	}, nil
}

// URL returns the URL endpoint used by the BlobClient object.
func (c *Client) URL() string {
	return c.svc.URL()
}

// CreateContainer is a lifecycle method to creates a new container under the specified account.
// If the container with the same name already exists, a ResourceExistsError will be raised.
// This method returns a client with which to interact with the newly created container.
func (c *Client) CreateContainer(ctx context.Context, containerName string, o *CreateContainerOptions) (CreateContainerResponse, error) {
	return c.svc.CreateContainer(ctx, containerName, o)
}

// DeleteContainer is a lifecycle method that marks the specified container for deletion.
// The container and any blobs contained within it are later deleted during garbage collection.
// If the container is not found, a ResourceNotFoundError will be raised.
func (c *Client) DeleteContainer(ctx context.Context, containerName string, o *DeleteContainerOptions) (DeleteContainerResponse, error) {
	return c.svc.DeleteContainer(ctx, containerName, o)
}

// DeleteBlob marks the specified blob or snapshot for deletion. The blob is later deleted during garbage collection.
// Note that deleting a blob also deletes all its snapshots.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-blob.
func (c *Client) DeleteBlob(ctx context.Context, containerName string, blobName string, o *DeleteBlobOptions) (DeleteBlobResponse, error) {
	return c.svc.NewContainerClient(containerName).NewBlobClient(blobName).Delete(ctx, o)
}

// NewListBlobsPager returns a pager for blobs starting from the specified Marker. Use an empty
// Marker to start enumeration from the beginning. Blob names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-blobs.
func (c *Client) NewListBlobsPager(containerName string, o *ListBlobsOptions) *runtime.Pager[ListBlobsResponse] {
	return c.svc.NewContainerClient(containerName).NewListBlobsFlatPager(o)
}

// NewListContainersPager operation returns a pager of the containers under the specified account.
// Use an empty Marker to start enumeration from the beginning. Container names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-containers2.
func (c *Client) NewListContainersPager(o *ListContainersOptions) *runtime.Pager[ListContainersResponse] {
	return c.svc.NewListContainersPager(o)
}

// TODO: Come here and expose highlevel functions here
//func (c *Client) Upload(ctx context.Context, containerName string, blobName string, data io.Reader, o *UploadOptions) (UploadResponse, error) {
//	o = shared.CopyOptions(o)
//	if o.TransferManager == nil {
//		// create a default transfer manager
//		if o.MaxBuffers == 0 {
//			o.MaxBuffers = 1
//		}
//		if o.BufferSize < blobrt.OneMB {
//			o.BufferSize = blobrt.OneMB
//		}
//		var err error
//		o.TransferManager, err = blobrt.NewStaticBuffer(o.BufferSize, o.MaxBuffers)
//		if err != nil {
//			return UploadResponse{}, fmt.Errorf("failed to create default transfer manager: %s", err)
//		}
//	} else {
//		// wrap in nop closer so we don't close caller's TM (caller is responsible for closing it)
//		o.TransferManager = &nopClosingTransferManager{o.TransferManager}
//	}
//
//	bb := c.svc.NewContainerClient(containerName).NewBlockBlobClient(blobName)
//	result, err := blockblob.Upload .ConcurrentUpload(ctx, data, bb, o)
//	if err != nil {
//		return UploadResponse{}, err
//	}
//	return result, nil
//}

// Download reads a range of bytes from a blob. The response also includes the blob's properties and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob.
func (c *Client) Download(ctx context.Context, containerName string, blobName string, o *DownloadOptions) (DownloadResponse, error) {
	o = shared.CopyOptions(o)
	return c.svc.NewContainerClient(containerName).NewBlobClient(blobName).Download(ctx, o.BlobOptions)
}

// ServiceClient returns the underlying *service.Client for this client.
func (c *Client) ServiceClient() *service.Client {
	return c.svc
}
