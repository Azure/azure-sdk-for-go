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

func getGeneratedBlobClient(b *blob.Client) *generated.BlobClient {
	return base.InnerClient((*base.Client[generated.BlobClient])(b))
}

func NewBatchBuilder() *BatchBuilder {
	return &BatchBuilder{}
}

func (b *BatchBuilder) AddBlobBatchDelete(blobPath string, o *DeleteOptions) {
	b.batchDeleteList = append(b.batchDeleteList, &BatchDeleteOptions{
		BlobPath:      &blobPath,
		DeleteOptions: o,
	})
}

func (b *BatchBuilder) AddBlobBatchSetTier(blobPath string, accessTier blob.AccessTier, o *SetTierOptions) {
	b.batchSetTierList = append(b.batchSetTierList, &BatchSetTierOptions{
		BlobPath:       &blobPath,
		AccessTier:     accessTier,
		SetTierOptions: o,
	})
}

func (b *BatchBuilder) Reset() {
	b.batchDeleteList = b.batchDeleteList[:0]
	b.batchSetTierList = b.batchSetTierList[:0]
}

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
			// handle error
			continue
		}
		batchRequest += subReq
		contentID++
	}

	// TODO: add for batch set tier
	for _, b := range b.batchSetTierList {
		blobPath := b.getBlobPath(&urlParts)
		blobURL := fmt.Sprintf("%s://%s%s", urlParts.Scheme, urlParts.Host, blobPath)
		subReq, err := b.createSetTierRequest(ctx, p, batchID, &blobURL, &contentID)
		if err != nil {
			// handle error
			continue
		}
		batchRequest += subReq
		contentID++
	}
	batchRequest += shared.GetBatchRequestDelimiter(*batchID, true, true) + shared.HttpNewline
	return batchRequest, nil
}

// TODO: add query parameters
func (d *BatchDeleteOptions) getBlobPath(urlParts *blob.URLParts) string {
	urlPath := *d.BlobPath
	if len(urlParts.ContainerName) != 0 {
		urlPath = "/" + urlParts.ContainerName + "/" + urlPath
	}

	return urlPath
}

func (d *BatchDeleteOptions) createDeleteRequest(ctx context.Context, p policy.Policy, batchID *string, blobURL *string, contentID *int) (string, error) {
	blobClient, err := blob.NewClientWithNoCredential(*blobURL, nil)
	if err != nil {
		return "", err
	}
	req, err := getGeneratedBlobClient(blobClient).DeleteCreateRequest(ctx, nil, nil, nil)
	if err != nil {
		return "", err
	}

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

// TODO: add query parameters
func (b *BatchSetTierOptions) getBlobPath(urlParts *blob.URLParts) string {
	urlPath := *b.BlobPath
	if len(urlParts.ContainerName) != 0 {
		urlPath = "/" + urlParts.ContainerName + "/" + urlPath
	}

	return urlPath
}

func (b *BatchSetTierOptions) createSetTierRequest(ctx context.Context, p policy.Policy, batchID *string, blobURL *string, contentID *int) (string, error) {
	blobClient, err := blob.NewClientWithNoCredential(*blobURL, nil)
	if err != nil {
		return "", err
	}
	req, err := getGeneratedBlobClient(blobClient).SetTierCreateRequest(ctx, b.AccessTier, nil, nil, nil)
	if err != nil {
		return "", err
	}

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
