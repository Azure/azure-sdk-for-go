//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"hash"
	"io"
	"reflect"
	"strconv"
	"strings"
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
	h         hash.Hash
	hashState []byte
}

// NewBlobDigestCalculator creates a new calculator to help to calculate blob digest when uploading blob.
func NewBlobDigestCalculator() *BlobDigestCalculator {
	return &BlobDigestCalculator{
		h: sha256.New(),
	}
}

func (b *BlobDigestCalculator) saveState() {
	b.hashState, _ = b.h.(encoding.BinaryMarshaler).MarshalBinary()
}

func (b *BlobDigestCalculator) restoreState() {
	if b.hashState == nil {
		return
	}
	_ = b.h.(encoding.BinaryUnmarshaler).UnmarshalBinary(b.hashState)
}

// BlobClientUploadChunkOptions contains the optional parameters for the BlobClient.UploadChunk method.
type BlobClientUploadChunkOptions struct {
	// Start of range for the blob to be uploaded.
	RangeStart *int32
	// End of range for the blob to be uploaded, inclusive.
	RangeEnd *int32
}

// UploadChunk - Upload a stream of data without completing the upload.
//
//   - location - Link acquired from upload start or previous chunk
//   - chunkData - Raw data of blob
//   - blobDigestCalculator - Calculator that help to calculate blob digest
//   - options - BlobClientUploadChunkOptions contains the optional parameters for the BlobClient.UploadChunk method.
func (client *BlobClient) UploadChunk(ctx context.Context, location string, chunkData io.ReadSeeker, blobDigestCalculator *BlobDigestCalculator, options *BlobClientUploadChunkOptions) (BlobClientUploadChunkResponse, error) {
	blobDigestCalculator.saveState()
	_, err := chunkData.Seek(0, io.SeekStart)
	if err != nil {
		return BlobClientUploadChunkResponse{}, err
	}
	payload, err := io.ReadAll(chunkData)
	if err != nil {
		return BlobClientUploadChunkResponse{}, err
	}
	_, err = blobDigestCalculator.h.Write(payload)
	if err != nil {
		return BlobClientUploadChunkResponse{}, err
	}
	_, err = chunkData.Seek(0, io.SeekStart)
	if err != nil {
		return BlobClientUploadChunkResponse{}, err
	}
	var requestOptions *blobClientUploadChunkOptions
	if options != nil && options.RangeStart != nil && options.RangeEnd != nil {
		requestOptions = &blobClientUploadChunkOptions{ContentRange: to.Ptr(fmt.Sprintf("%d-%d", *options.RangeStart, *options.RangeEnd))}
	}
	resp, err := client.uploadChunk(ctx, location, streaming.NopCloser(chunkData), requestOptions)
	if err != nil {
		blobDigestCalculator.restoreState()
	}
	return resp, err
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

// GetBlob - Retrieve the blob from the registry identified by digest.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-07-01
//   - name - Name of the image (including the namespace)
//   - digest - Digest of a BLOB
//   - options - BlobClientGetBlobOptions contains the optional parameters for the BlobClient.GetBlob method.
func (client *BlobClient) GetBlob(ctx context.Context, name string, digest string, options *BlobClientGetBlobOptions) (BlobClientGetBlobResponse, error) {
	resp, err := client.getBlob(ctx, name, digest, options)
	if err != nil {
		return resp, err
	}
	payload, err := io.ReadAll(resp.BlobData)
	_ = resp.BlobData.Close()
	if err != nil {
		return BlobClientGetBlobResponse{}, err
	}
	payloadDigest := fmt.Sprintf("sha256:%x", sha256.Sum256(payload))
	if digest != payloadDigest {
		return BlobClientGetBlobResponse{}, fmt.Errorf("retrieved blob digest %s does not match required digest %s", payloadDigest, digest)
	}
	if resp.DockerContentDigest != nil && *resp.DockerContentDigest != payloadDigest {
		return BlobClientGetBlobResponse{}, fmt.Errorf("retrieved blob digest %s does not match server-computed digest %s", payloadDigest, *resp.DockerContentDigest)
	}
	resp.BlobData = io.NopCloser(bytes.NewReader(payload))
	return resp, err
}

// GetChunk - Retrieve the blob from the registry identified by digest. This endpoint may also support RFC7233 compliant range
// requests. Support can be detected by issuing a HEAD request. If the header
// Accept-Range: bytes is returned, range requests can be used to fetch partial content.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-07-01
//   - name - Name of the image (including the namespace)
//   - digest - Digest of a BLOB
//   - rangeParam - Format : bytes=-, HTTP Range header specifying blob chunk.
//   - options - BlobClientGetChunkOptions contains the optional parameters for the BlobClient.GetChunk method.
func (client *BlobClient) GetChunk(ctx context.Context, name string, digest string, rangeParam string, blobDigestCalculator *BlobDigestCalculator, options *BlobClientGetChunkOptions) (BlobClientGetChunkResponse, error) {
	resp, err := client.getChunk(ctx, name, digest, rangeParam, options)
	if err != nil {
		return resp, err
	}
	payload, err := io.ReadAll(resp.ChunkData)
	_ = resp.ChunkData.Close()
	if err != nil {
		return BlobClientGetChunkResponse{}, err
	}
	_, err = blobDigestCalculator.h.Write(payload)
	if err != nil {
		return BlobClientGetChunkResponse{}, err
	}
	totalSize, err := strconv.Atoi(strings.Split(*resp.ContentRange, "/")[1])
	if err != nil {
		return BlobClientGetChunkResponse{}, err
	}
	currentRangeEnd, err := strconv.Atoi(strings.Split(strings.Split(*resp.ContentRange, "/")[0], "-")[1])
	if err != nil {
		return BlobClientGetChunkResponse{}, err
	}
	if totalSize == currentRangeEnd+1 {
		clientDigest := fmt.Sprintf("sha256:%x", blobDigestCalculator.h.Sum(nil))
		if digest != clientDigest {
			return BlobClientGetChunkResponse{}, fmt.Errorf("client-computed manifest digest %s does not match required digest %s", clientDigest, digest)
		}
	}
	resp.ChunkData = io.NopCloser(bytes.NewReader(payload))
	return resp, err
}
