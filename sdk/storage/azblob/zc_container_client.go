//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"net/http"
	"net/url"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ContainerClient represents a URL to the Azure Storage container allowing you to manipulate its blobs.
type ContainerClient struct {
	client    *containerClient
	sharedKey *SharedKeyCredential
}

// URL returns the URL endpoint used by the ContainerClient object.
func (c *ContainerClient) URL() string {
	return c.client.endpoint
}

// NewContainerClient creates a ContainerClient object using the specified URL, Azure AD credential, and options.
func NewContainerClient(containerURL string, cred azcore.TokenCredential, options *ClientOptions) *ContainerClient {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{internal.TokenScope}, nil)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := internal.NewConnection(containerURL, conOptions)

	return &ContainerClient{
		client: newContainerClient(conn.Endpoint(), conn.Pipeline()),
	}
}

// NewContainerClientWithNoCredential creates a ContainerClient object using the specified URL and options.
func NewContainerClientWithNoCredential(containerURL string, options *ClientOptions) *ContainerClient {
	conOptions := getConnectionOptions(options)
	conn := internal.NewConnection(containerURL, conOptions)

	return &ContainerClient{
		client: newContainerClient(conn.Endpoint(), conn.Pipeline()),
	}
}

// NewContainerClientWithSharedKey creates a ContainerClient object using the specified URL, shared key, and options.
func NewContainerClientWithSharedKey(containerURL string, cred *SharedKeyCredential, options *ClientOptions) *ContainerClient {
	authPolicy := newSharedKeyCredPolicy(cred)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := internal.NewConnection(containerURL, conOptions)

	return &ContainerClient{
		client:    newContainerClient(conn.Endpoint(), conn.Pipeline()),
		sharedKey: cred,
	}
}

// NewContainerClientFromConnectionString creates a ContainerClient object using connection string of an account
func NewContainerClientFromConnectionString(connectionString string, containerName string, options *ClientOptions) *ContainerClient {
	svcClient := NewServiceClientFromConnectionString(connectionString, options)
	return svcClient.NewContainerClient(containerName)
}

// NewBlobClient creates a new BlobClient object by concatenating blobName to the end of
// ContainerClient's URL. The new BlobClient uses the same request policy pipeline as the ContainerClient.
// To change the pipeline, create the BlobClient and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewBlobClient instead of calling this object's
// NewBlobClient method.
func (c *ContainerClient) NewBlobClient(blobName string) *BlobClient {
	blobURL := appendToURLPath(c.URL(), blobName)

	return &BlobClient{
		client:    newBlobClient(blobURL, c.client.pl),
		sharedKey: c.sharedKey,
	}
}

// NewAppendBlobClient creates a new AppendBlobURL object by concatenating blobName to the end of
// ContainerClient's URL. The new AppendBlobURL uses the same request policy pipeline as the ContainerClient.
// To change the pipeline, create the AppendBlobURL and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewAppendBlobClient instead of calling this object's
// NewAppendBlobClient method.
func (c *ContainerClient) NewAppendBlobClient(blobName string) *AppendBlobClient {
	blobURL := appendToURLPath(c.URL(), blobName)

	return &AppendBlobClient{
		BlobClient: BlobClient{
			client:    newBlobClient(blobURL, c.client.pl),
			sharedKey: c.sharedKey,
		},
		client: newAppendBlobClient(blobURL, c.client.pl),
	}
}

// NewBlockBlobClient creates a new BlockBlobClient object by concatenating blobName to the end of
// ContainerClient's URL. The new BlockBlobClient uses the same request policy pipeline as the ContainerClient.
// To change the pipeline, create the BlockBlobClient and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewBlockBlobClient instead of calling this object's
// NewBlockBlobClient method.
func (c *ContainerClient) NewBlockBlobClient(blobName string) *BlockBlobClient {
	blobURL := appendToURLPath(c.URL(), blobName)

	return &BlockBlobClient{
		BlobClient: BlobClient{
			client:    newBlobClient(blobURL, c.client.pl),
			sharedKey: c.sharedKey,
		},
		client: newBlockBlobClient(blobURL, c.client.pl),
	}
}

// NewPageBlobClient creates a new PageBlobURL object by concatenating blobName to the end of ContainerClient's URL. The new PageBlobURL uses the same request policy pipeline as the ContainerClient.
// To change the pipeline, create the PageBlobURL and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewPageBlobClient instead of calling this object's
// NewPageBlobClient method.
func (c *ContainerClient) NewPageBlobClient(blobName string) *PageBlobClient {
	blobURL := appendToURLPath(c.URL(), blobName)

	return &PageBlobClient{
		BlobClient: BlobClient{
			client:    newBlobClient(blobURL, c.client.pl),
			sharedKey: c.sharedKey,
		},
		client: newPageBlobClient(blobURL, c.client.pl),
	}
}

// Create creates a new container within a storage account. If a container with the same name already exists, the operation fails.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/create-container.
func (c *ContainerClient) Create(ctx context.Context, options *ContainerCreateOptions) (ContainerCreateResponse, error) {
	basics, cpkInfo := options.format()
	resp, err := c.client.Create(ctx, basics, cpkInfo)

	return toContainerCreateResponse(resp), err
}

// Delete marks the specified container for deletion. The container and any blobs contained within it are later deleted during garbage collection.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-container.
func (c *ContainerClient) Delete(ctx context.Context, o *ContainerDeleteOptions) (ContainerDeleteResponse, error) {
	basics, leaseInfo, accessConditions := o.format()
	resp, err := c.client.Delete(ctx, basics, leaseInfo, accessConditions)

	return toContainerDeleteResponse(resp), err
}

// GetProperties returns the container's properties.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-container-metadata.
func (c *ContainerClient) GetProperties(ctx context.Context, o *ContainerGetPropertiesOptions) (ContainerGetPropertiesResponse, error) {
	// NOTE: GetMetadata actually calls GetProperties internally because GetProperties returns the metadata AND the properties.
	// This allows us to not expose a GetProperties method at all simplifying the API.
	// The optionals are nil, like they were in track 1.5
	options, leaseAccess := o.format()
	resp, err := c.client.GetProperties(ctx, options, leaseAccess)

	return toContainerGetPropertiesResponse(resp), err
}

// SetMetadata sets the container's metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-container-metadata.
func (c *ContainerClient) SetMetadata(ctx context.Context, o *ContainerSetMetadataOptions) (ContainerSetMetadataResponse, error) {
	metadataOptions, lac, mac := o.format()
	resp, err := c.client.SetMetadata(ctx, metadataOptions, lac, mac)

	return toContainerSetMetadataResponse(resp), err
}

// GetAccessPolicy returns the container's access policy. The access policy indicates whether container's blobs may be accessed publicly.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-container-acl.
func (c *ContainerClient) GetAccessPolicy(ctx context.Context, o *ContainerGetAccessPolicyOptions) (ContainerGetAccessPolicyResponse, error) {
	options, ac := o.format()
	resp, err := c.client.GetAccessPolicy(ctx, options, ac)

	return toContainerGetAccessPolicyResponse(resp), err
}

// SetAccessPolicy sets the container's permissions. The access policy indicates whether blobs in a container may be accessed publicly.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-container-acl.
func (c *ContainerClient) SetAccessPolicy(ctx context.Context, o *ContainerSetAccessPolicyOptions) (ContainerSetAccessPolicyResponse, error) {
	accessPolicy, mac, lac := o.format()
	resp, err := c.client.SetAccessPolicy(ctx, accessPolicy, mac, lac)

	return toContainerSetAccessPolicyResponse(resp), err
}

// ListBlobsFlat returns a pager for blobs starting from the specified Marker. Use an empty
// Marker to start enumeration from the beginning. Blob names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-blobs.
func (c *ContainerClient) ListBlobsFlat(o *ContainerListBlobsFlatOptions) *runtime.Pager[ContainerListBlobsFlatResponse] {
	listOptions := o.format()
	return runtime.NewPager(runtime.PagingHandler[ContainerListBlobsFlatResponse]{
		More: func(page ContainerListBlobsFlatResponse) bool {
			if page.Marker == nil || len(*page.Marker) == 0 {
				return false
			}
			return true
		},
		Fetcher: func(ctx context.Context, page *ContainerListBlobsFlatResponse) (ContainerListBlobsFlatResponse, error) {
			var marker *string
			if page != nil {
				if page.NextMarker != nil {
					marker = page.NextMarker
				}
			} else {
				// If provided by the user, then use the one from options bag
				marker = listOptions.Marker
			}

			req, err := c.client.listBlobFlatSegmentCreateRequest(ctx, &listOptions)
			if err != nil {
				return ContainerListBlobsFlatResponse{}, err
			}
			if marker != nil {
				queryValues, err := url.ParseQuery(req.Raw().URL.RawQuery)
				if err != nil {
					return ContainerListBlobsFlatResponse{}, err
				}
				queryValues.Set("marker", *marker)
				req.Raw().URL.RawQuery = queryValues.Encode()
				if err != nil {
					return ContainerListBlobsFlatResponse{}, err
				}
			}

			resp, err := c.client.pl.Do(req)
			if err != nil {
				return ContainerListBlobsFlatResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ContainerListBlobsFlatResponse{}, runtime.NewResponseError(resp)
			}
			generatedResp, err := c.client.listBlobFlatSegmentHandleResponse(resp)
			return toContainerListBlobsFlatResponse(generatedResp), err
		},
	})
}

// ListBlobsHierarchy returns a channel of blobs starting from the specified Marker. Use an empty
// Marker to start enumeration from the beginning. Blob names are returned in lexicographic order.
// After getting a segment, process it, and then call ListBlobsHierarchicalSegment again (passing the the
// previously-returned Marker) to get the next segment.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-blobs.
// AutoPagerTimeout specifies the amount of time with no read operations before the channel times out and closes. Specify no time and it will be ignored.
// AutoPagerBufferSize specifies the channel's buffer size.
// Both the blob item channel and error channel should be watched. Only one error will be released via this channel (or a nil error, to register a clean exit.)
func (c *ContainerClient) ListBlobsHierarchy(delimiter string, o *ContainerListBlobsHierarchyOptions) *runtime.Pager[ContainerListBlobHierarchyResponse] {
	listOptions := o.format()
	return runtime.NewPager(runtime.PagingHandler[ContainerListBlobHierarchyResponse]{
		More: func(page ContainerListBlobHierarchyResponse) bool {
			if page.Marker == nil || len(*page.Marker) == 0 {
				return false
			}
			return true
		},
		Fetcher: func(ctx context.Context, page *ContainerListBlobHierarchyResponse) (ContainerListBlobHierarchyResponse, error) {
			var marker *string
			if page != nil {
				if page.NextMarker != nil {
					marker = page.NextMarker
				}
			} else {
				// If provided by the user, then use the one from options bag
				marker = listOptions.Marker
			}

			req, err := c.client.listBlobHierarchySegmentCreateRequest(ctx, delimiter, &listOptions)
			if err != nil {
				return ContainerListBlobHierarchyResponse{}, err
			}
			if marker != nil {
				queryValues, err := url.ParseQuery(req.Raw().URL.RawQuery)
				if err != nil {
					return ContainerListBlobHierarchyResponse{}, err
				}
				queryValues.Set("marker", *marker)
				req.Raw().URL.RawQuery = queryValues.Encode()
				if err != nil {
					return ContainerListBlobHierarchyResponse{}, err
				}
			}

			resp, err := c.client.pl.Do(req)
			if err != nil {
				return ContainerListBlobHierarchyResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ContainerListBlobHierarchyResponse{}, runtime.NewResponseError(resp)
			}
			generatedResp, err := c.client.listBlobHierarchySegmentHandleResponse(resp)
			return toContainerListBlobHierarchyResponse(generatedResp), err
		},
	})
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at container.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (c *ContainerClient) GetSASURL(permissions ContainerSASPermissions, start time.Time, expiry time.Time) (string, error) {
	if c.sharedKey == nil {
		return "", errors.New("SAS can only be signed with a SharedKeyCredential")
	}

	urlParts, err := NewBlobURLParts(c.URL())
	if err != nil {
		return "", err
	}

	// Containers do not have snapshots, nor versions.
	urlParts.SAS, err = BlobSASSignatureValues{
		ContainerName: urlParts.ContainerName,
		Permissions:   permissions.String(),
		StartTime:     start.UTC(),
		ExpiryTime:    expiry.UTC(),
	}.NewSASQueryParameters(c.sharedKey)

	return urlParts.URL(), err
}
