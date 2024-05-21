//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

const alpineBlobDigest = "sha256:042a816809aac8d0f7d7cacac7965782ee2ecac3f21bcf9f24b1de1a7387b769"

func TestBlobClient_CancelUpload(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	startRes, err := client.StartUpload(ctx, "hello-world", nil)
	require.NoError(t, err)
	_, err = client.CancelUpload(ctx, *startRes.Location, nil)
	require.NoError(t, err)
}

func TestBlobClient_CancelUpload_fail(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.CancelUpload(ctx, "wrong location", nil)
	require.Error(t, err)
}

func TestBlobClient_CheckBlobExists(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.CheckBlobExists(ctx, "alpine", alpineBlobDigest, nil)
	require.NoError(t, err)
	require.Equal(t, alpineBlobDigest, *res.DockerContentDigest)
}

func TestBlobClient_CheckBlobExists_fail(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.CheckBlobExists(ctx, "alpine", "wrong digest", nil)
	require.Error(t, err)
}

func TestBlobClient_CheckBlobExists_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.CheckBlobExists(ctx, "", "digest", nil)
	require.Error(t, err)
	_, err = client.CheckBlobExists(ctx, "alpine", "", nil)
	require.Error(t, err)
}

func TestBlobClient_CheckChunkExists(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.CheckChunkExists(ctx, "alpine", alpineBlobDigest, "bytes=0-299", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.ContentLength)
}

func TestBlobClient_CheckChunkExists_fail(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.CheckChunkExists(ctx, "alpine", "wrong digest", "bytes=0-299", nil)
	require.Error(t, err)
}

func TestBlobClient_CheckChunkExists_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.CheckChunkExists(ctx, "", "digest", "range", nil)
	require.Error(t, err)
	_, err = client.CheckChunkExists(ctx, "name", "", "range", nil)
	require.Error(t, err)
}

func TestBlobClient_completeUpload_wrongDigest(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	getRes, err := client.GetBlob(ctx, "alpine", alpineBlobDigest, nil)
	require.NoError(t, err)
	blob, err := io.ReadAll(getRes.BlobData)
	require.NoError(t, err)
	startRes, err := client.StartUpload(ctx, "hello-world", nil)
	require.NoError(t, err)
	uploadResp, err := client.uploadChunk(ctx, *startRes.Location, streaming.NopCloser(bytes.NewReader(blob)), nil)
	require.NoError(t, err)
	_, err = client.completeUpload(ctx, "sha256:00000000", *uploadResp.Location, nil)
	require.Error(t, err)
}

func TestBlobClient_DeleteBlob(t *testing.T) {
	d := "sha256:bfe296a525011f7eb76075d688c681ca4feaad5afe3b142b36e30f1a171dc99a"
	require.NotEqual(t, d, alpineBlobDigest, "test bug: deleting a blob used in other tests")
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteBlob(ctx, "alpine", d, nil)
	require.NoError(t, err)
}

func TestBlobClient_DeleteBlob_fail(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteBlob(ctx, "alpine", "wrong digest", nil)
	require.Error(t, err)
}

func TestBlobClient_DeleteBlob_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.DeleteBlob(ctx, "", "digest", nil)
	require.Error(t, err)
	_, err = client.DeleteBlob(ctx, "name", "", nil)
	require.Error(t, err)
}

func TestBlobClient_GetBlob(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.GetBlob(ctx, "alpine", alpineBlobDigest, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.ContentLength)
	reader, err := NewDigestValidationReader(alpineBlobDigest, res.BlobData)
	require.NoError(t, err)
	_, err = io.ReadAll(reader)
	require.NoError(t, err)
}

func TestBlobClient_GetBlob_wrongDigest(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("test")))

	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	require.NoError(t, err)
	client := &BlobClient{
		azcoreClient,
		srv.URL(),
	}
	ctx := context.Background()
	resp, err := client.GetBlob(ctx, "name", "sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", nil)
	require.NoError(t, err)
	reader, err := NewDigestValidationReader("sha256:wrong", resp.BlobData)
	require.NoError(t, err)
	_, err = io.ReadAll(reader)
	require.Error(t, err, ErrMismatchedHash)
}

func TestBlobClient_GetBlob_fail(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetBlob(ctx, "alpine", "wrong digest", nil)
	require.Error(t, err)
}

func TestBlobClient_GetBlob_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.GetBlob(ctx, "", "digest", nil)
	require.Error(t, err)
	_, err = client.GetBlob(ctx, "alpine", "", nil)
	require.Error(t, err)
}

func TestBlobClient_GetChunk(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	chunkSize := 1000
	current := 0
	blob := bytes.NewBuffer(nil)
	for {
		res, err := client.GetChunk(ctx, "alpine", alpineBlobDigest, fmt.Sprintf("bytes=%d-%d", current, current+chunkSize-1), nil)
		require.NoError(t, err)
		chunk, err := io.ReadAll(res.ChunkData)
		require.NoError(t, err)
		_, err = blob.Write(chunk)
		require.NoError(t, err)
		totalSize, _ := strconv.Atoi(strings.Split(*res.ContentRange, "/")[1])
		currentRangeEnd, _ := strconv.Atoi(strings.Split(strings.Split(*res.ContentRange, "/")[0], "-")[1])
		if totalSize == currentRangeEnd+1 {
			break
		}
		current += chunkSize
	}
	reader, err := NewDigestValidationReader(alpineBlobDigest, blob)
	require.NoError(t, err)
	_, err = io.ReadAll(reader)
	require.NoError(t, err)
}

func TestBlobClient_GetChunk_fail(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetChunk(ctx, "alpine", "wrong digest", "bytes=0-999", nil)
	require.Error(t, err)
}

func TestBlobClient_GetChunk_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.GetChunk(ctx, "", "digest", "bytes=0-999", nil)
	require.Error(t, err)
	_, err = client.GetChunk(ctx, "alpine", "", "bytes=0-999", nil)
	require.Error(t, err)
}

func TestBlobClient_GetUploadStatus(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	startRes, err := client.StartUpload(ctx, "hello-world", nil)
	require.NoError(t, err)
	checkResp, err := client.GetUploadStatus(ctx, *startRes.Location, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *checkResp.DockerUploadUUID)
	_, err = client.CancelUpload(ctx, *startRes.Location, nil)
	require.NoError(t, err)
}

func TestBlobClient_GetUploadStatus_fail(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetUploadStatus(ctx, "wrong location", nil)
	require.Error(t, err)
}

func TestBlobClient_MountBlob(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	res, err := client.MountBlob(ctx, "hello-world", "alpine", alpineBlobDigest, nil)
	require.NoError(t, err)
	require.NotEmpty(t, res.Location)
}

func TestBlobClient_MountBlob_fail(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.MountBlob(ctx, "wrong name", "wrong from", "wrong mount", nil)
	require.Error(t, err)
}

func TestBlobClient_MountBlob_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.MountBlob(ctx, "", "from", "mount", nil)
	require.Error(t, err)
}

func TestBlobClient_StartUpload(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	startRes, err := client.StartUpload(ctx, "hello-world", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *startRes.Location)
	_, err = client.CancelUpload(ctx, *startRes.Location, nil)
	require.NoError(t, err)
}

func TestBlobClient_StartUpload_empty(t *testing.T) {
	ctx := context.Background()
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.StartUpload(ctx, "", nil)
	require.Error(t, err)
}

func TestBlobClient_wrongEndpoint(t *testing.T) {
	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{}, nil)
	require.NoError(t, err)
	client := &BlobClient{
		azcoreClient,
		"wrong-endpoint",
	}
	ctx := context.Background()
	_, err = client.CancelUpload(ctx, "location", nil)
	require.Error(t, err)
	_, err = client.CheckBlobExists(ctx, "name", "digest", nil)
	require.Error(t, err)
	_, err = client.CheckChunkExists(ctx, "name", "digest", "range", nil)
	require.Error(t, err)
	_, err = client.completeUpload(ctx, "digest", "location", nil)
	require.Error(t, err)
	_, err = client.DeleteBlob(ctx, "name", "digest", nil)
	require.Error(t, err)
	_, err = client.GetBlob(ctx, "name", "digest", nil)
	require.Error(t, err)
	_, err = client.GetChunk(ctx, "name", "digest", "range", nil)
	require.Error(t, err)
	_, err = client.GetUploadStatus(ctx, "location", nil)
	require.Error(t, err)
	_, err = client.MountBlob(ctx, "name", "from", "mount", nil)
	require.Error(t, err)
	_, err = client.StartUpload(ctx, "name", nil)
	require.Error(t, err)
	_, err = client.uploadChunk(ctx, "digest", nil, nil)
	require.Error(t, err)
}
