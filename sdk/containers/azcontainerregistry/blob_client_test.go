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
	startRes, err := client.StartUpload(ctx, "hello-world", nil)
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
	digest := "sha256:042a816809aac8d0f7d7cacac7965782ee2ecac3f21bcf9f24b1de1a7387b769"
	res, err := client.CheckBlobExists(ctx, "alpine", digest, nil)
	require.NoError(t, err)
	require.Equal(t, digest, *res.DockerContentDigest)
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
