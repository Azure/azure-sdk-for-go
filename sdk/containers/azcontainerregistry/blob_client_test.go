//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

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
	digest := "sha256:042a816809aac8d0f7d7cacac7965782ee2ecac3f21bcf9f24b1de1a7387b769"
	res, err := client.CheckBlobExists(ctx, "alpine", digest, nil)
	require.NoError(t, err)
	require.Equal(t, digest, *res.DockerContentDigest)
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
	digest := "sha256:042a816809aac8d0f7d7cacac7965782ee2ecac3f21bcf9f24b1de1a7387b769"
	res, err := client.CheckChunkExists(ctx, "alpine", digest, "bytes=0-299", nil)
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
	digest := "sha256:042a816809aac8d0f7d7cacac7965782ee2ecac3f21bcf9f24b1de1a7387b769"
	getRes, err := client.GetBlob(ctx, "alpine", digest, nil)
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
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:042a816809aac8d0f7d7cacac7965782ee2ecac3f21bcf9f24b1de1a7387b769"
	_, err = client.DeleteBlob(ctx, "alpine", digest, nil)
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
	digest := "sha256:042a816809aac8d0f7d7cacac7965782ee2ecac3f21bcf9f24b1de1a7387b769"
	res, err := client.GetBlob(ctx, "alpine", digest, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.ContentLength)
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
	digest := "sha256:042a816809aac8d0f7d7cacac7965782ee2ecac3f21bcf9f24b1de1a7387b769"
	res, err := client.GetChunk(ctx, "alpine", digest, "bytes=0-999", nil)
	require.NoError(t, err)
	require.Equal(t, int64(1000), *res.ContentLength)
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
	digest := "sha256:042a816809aac8d0f7d7cacac7965782ee2ecac3f21bcf9f24b1de1a7387b769"
	res, err := client.MountBlob(ctx, "hello-world", "alpine", digest, nil)
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
	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, nil)
	client := &BlobClient{
		"wrong-endpoint",
		pl,
	}
	ctx := context.Background()
	_, err := client.CancelUpload(ctx, "location", nil)
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
