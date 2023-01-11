//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"hash"
	"io"
	"reflect"
)

// BlobClientOptions contains the optional parameters for the NewBlobClient method.
type BlobClientOptions struct {
	azcore.ClientOptions
}

// NewBlobClient creates a new instance of BlobClient with the specified values.
//   - endpoint - registry login URL
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - client options, pass nil to accept the default values.
func NewBlobClient(endpoint string, credential azcore.TokenCredential, options *BlobClientOptions) (*BlobClient, error) {
	if options == nil {
		options = &BlobClientOptions{}
	}

	if reflect.ValueOf(options.Cloud).IsZero() {
		options.Cloud = cloud.AzurePublic
	}
	c, ok := options.Cloud.Services[ServiceName]
	if !ok || c.Audience == "" {
		return nil, errors.New("provided Cloud field is missing Azure Container Registry configuration")
	}

	authClient := newAuthenticationClient(endpoint, &authenticationClientOptions{
		options.ClientOptions,
	})
	authPolicy := newAuthenticationPolicy(
		credential,
		[]string{c.Audience + "/.default"},
		authClient,
		nil,
	)

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &BlobClient{
		endpoint,
		pl,
	}, nil
}

// BlobDigestCalculator help to calculate blob digest when uploading blob.
// Don't use this type directly, use NewBlobDigestCalculator() instead.
type BlobDigestCalculator struct {
	h hash.Hash
}

type wrappedReadSeekCloser struct {
	io.Reader
	io.Seeker
	io.Closer
}

// NewBlobDigestCalculator creates a new calculator to help to calculate blob digest when uploading blob.
func NewBlobDigestCalculator() *BlobDigestCalculator {
	return &BlobDigestCalculator{
		h: sha256.New(),
	}
}

// UploadChunk - Upload a stream of data without completing the upload.
//
//   - location - Link acquired from upload start or previous chunk
//   - chunkData - Raw data of blob
//   - blobDigestCalculator - Calculator that help to calculate blob digest
//   - options - BlobClientUploadChunkOptions contains the optional parameters for the BlobClient.UploadChunk method.
func (client *BlobClient) UploadChunk(ctx context.Context, location string, chunkData io.ReadSeekCloser, blobDigestCalculator *BlobDigestCalculator, options *BlobClientUploadChunkOptions) (BlobClientUploadChunkResponse, error) {
	wrappedChunkData := &wrappedReadSeekCloser{Reader: io.TeeReader(chunkData, blobDigestCalculator.h), Seeker: chunkData, Closer: chunkData}
	return client.uploadChunk(ctx, location, wrappedChunkData, options)
}

// CompleteUpload - Complete the upload with previously uploaded content.
//
//   - digest - Digest of a BLOB
//   - location - Link acquired from upload start or previous chunk
//   - blobDigestCalculator - Calculator that help to calculate blob digest
//   - options - BlobClientCompleteUploadOptions contains the optional parameters for the BlobClient.CompleteUpload method.
func (client *BlobClient) CompleteUpload(ctx context.Context, location string, blobDigestCalculator *BlobDigestCalculator, options *BlobClientCompleteUploadOptions) (BlobClientCompleteUploadResponse, error) {
	return client.completeUpload(ctx, fmt.Sprintf("sha256:%x", blobDigestCalculator.h.Sum(nil)), location, options)
}
