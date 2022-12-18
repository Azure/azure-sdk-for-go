//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"bytes"
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestContainerRegistryBlobClient_CancelUpload(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	startRes, err := client.StartUpload(ctx, "hello-world-test", nil)
	require.NoError(t, err)
	_, err = client.CancelUpload(ctx, *startRes.Location, nil)
	require.NoError(t, err)
}

func TestContainerRegistryBlobClient_CheckBlobExists(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	res, err := client.CheckBlobExists(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	require.Equal(t, digest, *res.DockerContentDigest)
}

func TestContainerRegistryBlobClient_CheckChunkExists(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	res, err := client.CheckChunkExists(ctx, "hello-world", digest, "bytes=0-299", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.ContentLength)
}

func TestContainerRegistryBlobClient_CompleteUpload(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	getRes, err := client.GetBlob(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	blob, err := io.ReadAll(getRes.Body)
	startRes, err := client.StartUpload(ctx, "hello-world-test", nil)
	require.NoError(t, err)
	uploadResp, err := client.UploadChunk(ctx, *startRes.Location, streaming.NopCloser(bytes.NewReader(blob)), nil)
	require.NoError(t, err)
	completeResp, err := client.CompleteUpload(ctx, digest, *uploadResp.Location, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *completeResp.DockerContentDigest)
}

func TestContainerRegistryBlobClient_DeleteBlob(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	_, err = client.DeleteBlob(ctx, "hello-world-test", digest, nil)
	require.NoError(t, err)
}

func TestContainerRegistryBlobClient_GetBlob(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	res, err := client.GetBlob(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.ContentLength)
}

func TestContainerRegistryBlobClient_GetChunk(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:feac5306138255e28a9862d3f3d29025d0a4d0648855afe1acd6131af07138ac"
	res, err := client.GetChunk(ctx, "ubuntu", digest, "bytes=0-1000", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *res.ContentLength)
}

func TestContainerRegistryBlobClient_GetUploadStatus(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	startRes, err := client.StartUpload(ctx, "hello-world-test", nil)
	require.NoError(t, err)
	checkResp, err := client.GetUploadStatus(ctx, *startRes.Location, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *checkResp.DockerUploadUUID)
	_, err = client.CancelUpload(ctx, *startRes.Location, nil)
	require.NoError(t, err)
}

func TestContainerRegistryBlobClient_MountBlob(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	res, err := client.MountBlob(ctx, "hello-world", "hello-world-test", digest, nil)
	require.NoError(t, err)
	require.NotEmpty(t, res.Location)
}

func TestContainerRegistryBlobClient_StartUpload(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	startRes, err := client.StartUpload(ctx, "hello-world-test", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *startRes.Location)
	_, err = client.CancelUpload(ctx, *startRes.Location, nil)
	require.NoError(t, err)
}

func TestContainerRegistryBlobClient_UploadChunk(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewContainerRegistryBlobClient("https://azacrlivetest.azurecr.io", cred, &ContainerRegistryBlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	getRes, err := client.GetBlob(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	blob, err := io.ReadAll(getRes.Body)
	startRes, err := client.StartUpload(ctx, "hello-world-test", nil)
	require.NoError(t, err)
	uploadResp, err := client.UploadChunk(ctx, *startRes.Location, streaming.NopCloser(bytes.NewReader(blob)), nil)
	require.NoError(t, err)
	require.NotEmpty(t, *uploadResp.Location)
	_, err = client.CancelUpload(ctx, *uploadResp.Location, nil)
	require.NoError(t, err)
}
