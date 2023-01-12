//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blobbatch

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

// helper method to return the generated.BlobClient which is used for creating the sub-requests
func getGeneratedBlobClient(b *blob.Client) *generated.BlobClient {
	return base.InnerClient((*base.Client[generated.BlobClient])(b))
}

// NewBatchBuilder creates a new blob batch builder
func NewBatchBuilder() *BatchBuilder {
	return &BatchBuilder{}
}

// AddBlobBatchDelete adds delete operation to the batch operations list
//   - blobPath: Must be in the following format.
//   - blobName when using ContainerBatchClient, e.g. blob.txt
//   - /containerName/blobName when using ServiceBatchClient, e.g. /container/blob.txt
//   - DeleteOptions - delete options; pass nil to accept default values
func (b *BatchBuilder) AddBlobBatchDelete(blobPath string, o *DeleteOptions) {
	b.batchDeleteList = append(b.batchDeleteList, &BatchDeleteOptions{
		BlobPath:      blobPath,
		DeleteOptions: o,
	})
}

// AddBlobBatchSetTier adds set tier operation to the batch operations list
//   - blobPath: Must be in the following format.
//   - blobName when using ContainerBatchClient, e.g. blob.txt
//   - /containerName/blobName when using ServiceBatchClient, e.g. /container/blob.txt
//   - accessTier - defines values for Blob Access Tier
//   - SetTierOptions - set tier options; pass nil to accept default values
func (b *BatchBuilder) AddBlobBatchSetTier(blobPath string, accessTier blob.AccessTier, o *SetTierOptions) {
	b.batchSetTierList = append(b.batchSetTierList, &BatchSetTierOptions{
		BlobPath:       blobPath,
		AccessTier:     accessTier,
		SetTierOptions: o,
	})
}

// Reset is used for cleaning the batch builder list
func (b *BatchBuilder) Reset() {
	b.batchDeleteList = b.batchDeleteList[:0]
	b.batchSetTierList = b.batchSetTierList[:0]
}

// createBatchRequest creates a new batch request using the operations present in the batch builder list
//   - policy - auth policy for authorizing the sub requests. It is nil in case of SAS or User Delegation SAS
//   - url:
//   - container url in case of container batch client
//   - service url in case of service batch client
//   - batchID - used as batch boundary in the request body
//
// Example of a sub-request in the batch request body
//
//	--batch_357de4f7-6d0b-4e02-8cd2-6361411a9525
//	Content-Type: application/http
//	Content-Transfer-Encoding: binary
//	Content-ID: 0
//
//	DELETE /container0/blob0 HTTP/1.1
//	x-ms-date: Thu, 14 Jun 2018 16:46:54 GMT
//	Authorization: SharedKey account:G4jjBXA7LI/RnWKIOQ8i9xH4p76pAQ+4Fs4R1VxasaE=
//	Content-Length: 0
func (b *BatchBuilder) createBatchRequest(ctx context.Context, p policy.Policy, url *string, batchID *string) (string, error) {
	urlParts, err := blob.ParseURL(*url)
	if err != nil {
		return "", err
	}

	contentID := 1
	batchRequest := ""
	for _, d := range b.batchDeleteList {
		blobPath := d.getBlobPath(&urlParts)
		blobURL := fmt.Sprintf("%s://%s%s", urlParts.Scheme, urlParts.Host, blobPath)
		subReq, err := d.createDeleteRequest(ctx, p, batchID, &blobURL, &contentID)
		if err != nil {
			// TODO: handle error
			continue
		}
		batchRequest += subReq
		contentID++
	}

	for _, b := range b.batchSetTierList {
		blobPath := b.getBlobPath(&urlParts)
		blobURL := fmt.Sprintf("%s://%s%s", urlParts.Scheme, urlParts.Host, blobPath)
		subReq, err := b.createSetTierRequest(ctx, p, batchID, &blobURL, &contentID)
		if err != nil {
			// TODO: handle error
			continue
		}
		batchRequest += subReq
		contentID++
	}

	// add the last line of the request body. It looks like,
	// --batch_357de4f7-6d0b-4e02-8cd2-6361411a9525--
	batchRequest += shared.GetBatchRequestDelimiter(*batchID, true, true) + shared.HttpNewline
	return batchRequest, nil
}

// getBlobBatch returns the blob path including the query parameters in the following format.
// e.g. /containerName/blobName?queryParams
// TODO: add query parameters in the blob path, make this a common method for BatchDeleteOptions and BatchSetTierOptions
func (d *BatchDeleteOptions) getBlobPath(urlParts *blob.URLParts) string {
	urlPath := d.BlobPath
	if len(urlParts.ContainerName) != 0 {
		urlPath = "/" + urlParts.ContainerName + "/" + urlPath
	}

	return urlPath
}

// createDeleteRequest is used for creating the sub-request for delete operation using the generated.BlobClient
//   - policy - auth policy for authorizing the sub-request
//   - batchID - used as batch boundary in the request body
//   - blobURL - e.g. https://<account>.blob.core.windows.net/container/blob.txt?queryParams
//   - contentID - used in the Content-ID header for each sub-request response which is used for tracking
func (d *BatchDeleteOptions) createDeleteRequest(ctx context.Context, p policy.Policy, batchID *string, blobURL *string, contentID *int) (string, error) {
	blobClient, err := blob.NewClientWithNoCredential(*blobURL, nil)
	if err != nil {
		return "", err
	}
	req, err := getGeneratedBlobClient(blobClient).DeleteCreateRequest(ctx, nil, nil, nil)
	if err != nil {
		return "", err
	}

	// update the sub-request headers. Add x-ms-date and remove x-ms-version header
	shared.UpdateSubRequestHeaders(req)

	// adding authorization header to the request by executing the policy
	if p != nil {
		resp, err := p.Do(req)
		if resp != nil && err != nil {
			return fmt.Sprintf("%v: %v", resp.StatusCode, resp.Status), err
		}
	}

	batchSubRequest := shared.CreateSubReqHeader(*batchID, *contentID)
	batchSubRequest += shared.BuildSubRequest(req)

	return batchSubRequest, nil
}

// getBlobBatch returns the blob path including the query parameters in the following format.
// e.g. /containerName/blobName?queryParams
// TODO: add query parameters in the blob path, make this a common method for BatchDeleteOptions and BatchSetTierOptions
func (b *BatchSetTierOptions) getBlobPath(urlParts *blob.URLParts) string {
	urlPath := b.BlobPath
	if len(urlParts.ContainerName) != 0 {
		urlPath = "/" + urlParts.ContainerName + "/" + urlPath
	}

	return urlPath
}

// createSetTierRequest is used for creating the sub-request for set tier operation using the generated.BlobClient
//   - policy - auth policy for authorizing the sub-request
//   - batchID - used as batch boundary in the request body
//   - blobURL - e.g. https://<account>.blob.core.windows.net/container/blob.txt?queryParams
//   - contentID - used in the Content-ID header for each sub-request response which is used for tracking
func (b *BatchSetTierOptions) createSetTierRequest(ctx context.Context, p policy.Policy, batchID *string, blobURL *string, contentID *int) (string, error) {
	blobClient, err := blob.NewClientWithNoCredential(*blobURL, nil)
	if err != nil {
		return "", err
	}
	req, err := getGeneratedBlobClient(blobClient).SetTierCreateRequest(ctx, b.AccessTier, nil, nil, nil)
	if err != nil {
		return "", err
	}

	// update the sub-request headers. Add x-ms-date and remove x-ms-version header
	shared.UpdateSubRequestHeaders(req)

	// adding authorization header to the request by executing the policy
	if p != nil {
		resp, err := p.Do(req)
		if resp != nil && err != nil {
			return fmt.Sprintf("%v: %v", resp.StatusCode, resp.Status), err
		}
	}

	batchSubRequest := shared.CreateSubReqHeader(*batchID, *contentID)
	batchSubRequest += shared.BuildSubRequest(req)

	return batchSubRequest, nil
}
