//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"context"
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
	startRes, err := client.StartUpload(ctx, "hello-world-test", nil)
	require.NoError(t, err)
	_, err = client.CancelUpload(ctx, *startRes.Location, nil)
	require.NoError(t, err)
}

func TestBlobClient_CheckBlobExists(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	res, err := client.CheckBlobExists(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	require.Equal(t, digest, *res.DockerContentDigest)
}

func TestBlobClient_CheckChunkExists(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	res, err := client.CheckChunkExists(ctx, "hello-world", digest, "bytes=0-299", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.ContentLength)
}

func TestBlobClient_completeUpload_wrongDigest(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	getRes, err := client.GetBlob(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	blob, err := io.ReadAll(getRes.BlobData)
	require.NoError(t, err)
	startRes, err := client.StartUpload(ctx, "hello-world-test", nil)
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
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	_, err = client.DeleteBlob(ctx, "hello-world-test", digest, nil)
	require.NoError(t, err)
}

func TestBlobClient_GetBlob(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	res, err := client.GetBlob(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.ContentLength)
}

func TestBlobClient_GetChunk(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:feac5306138255e28a9862d3f3d29025d0a4d0648855afe1acd6131af07138ac"
	res, err := client.GetChunk(ctx, "ubuntu", digest, "bytes=0-1000", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.ContentLength)
}

func TestBlobClient_GetUploadStatus(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	startRes, err := client.StartUpload(ctx, "hello-world-test", nil)
	require.NoError(t, err)
	checkResp, err := client.GetUploadStatus(ctx, *startRes.Location, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *checkResp.DockerUploadUUID)
	_, err = client.CancelUpload(ctx, *startRes.Location, nil)
	require.NoError(t, err)
}

func TestBlobClient_MountBlob(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	res, err := client.MountBlob(ctx, "hello-world", "hello-world-test", digest, nil)
	require.NoError(t, err)
	require.NotEmpty(t, res.Location)
}

func TestBlobClient_StartUpload(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	startRes, err := client.StartUpload(ctx, "hello-world-test", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *startRes.Location)
	_, err = client.CancelUpload(ctx, *startRes.Location, nil)
	require.NoError(t, err)
}
