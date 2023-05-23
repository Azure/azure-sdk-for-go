//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"context"
	"crypto/sha256"
	"encoding"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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
		options.Cloud = defaultCloud
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

type wrappedReadSeeker struct {
	io.Reader
	io.Seeker
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

// newLimitTeeReader returns a Reader that writes to w what it reads from r with n bytes limit.
func newLimitTeeReader(r io.Reader, w io.Writer, n int64) io.Reader {
	return &limitTeeReader{r, w, n}
}

type limitTeeReader struct {
	r io.Reader
	w io.Writer
	n int64
}

func (lt *limitTeeReader) Read(p []byte) (int, error) {
	n, err := lt.r.Read(p)
	if n > 0 && lt.n > 0 {
		wn, werr := lt.w.Write(p[:n])
		if werr != nil {
			return wn, werr
		}
		lt.n -= int64(wn)
	}
	return n, err
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
	size, err := chunkData.Seek(0, io.SeekEnd) // Seek to the end to get the stream's size
	if err != nil {
		return BlobClientUploadChunkResponse{}, err
	}
	_, err = chunkData.Seek(0, io.SeekStart)
	if err != nil {
		return BlobClientUploadChunkResponse{}, err
	}
	wrappedChunkData := &wrappedReadSeeker{Reader: newLimitTeeReader(chunkData, blobDigestCalculator.h, size), Seeker: chunkData}
	var requestOptions *blobClientUploadChunkOptions
	if options != nil && options.RangeStart != nil && options.RangeEnd != nil {
		requestOptions = &blobClientUploadChunkOptions{ContentRange: to.Ptr(fmt.Sprintf("%d-%d", *options.RangeStart, *options.RangeEnd))}
	}
	resp, err := client.uploadChunk(ctx, location, streaming.NopCloser(wrappedChunkData), requestOptions)
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
