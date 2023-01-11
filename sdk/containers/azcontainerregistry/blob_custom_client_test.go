//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestBlobClient_CompleteUpload(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient("https://azacrlivetest.azurecr.io", cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	getRes, err := client.GetBlob(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	blob, err := io.ReadAll(getRes.BlobData)
	startRes, err := client.StartUpload(ctx, "hello-world-test", nil)
	require.NoError(t, err)
	calculator := NewBlobDigestCalculator()
	uploadResp, err := client.UploadChunk(ctx, *startRes.Location, streaming.NopCloser(bytes.NewReader(blob)), calculator, nil)
	require.NoError(t, err)
	completeResp, err := client.CompleteUpload(ctx, *uploadResp.Location, calculator, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *completeResp.DockerContentDigest)
}

func TestBlobClient_UploadChunk(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient("https://azacrlivetest.azurecr.io", cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	getRes, err := client.GetBlob(ctx, "hello-world", digest, nil)
	require.NoError(t, err)
	blob, err := io.ReadAll(getRes.BlobData)
	startRes, err := client.StartUpload(ctx, "hello-world-test", nil)
	require.NoError(t, err)
	calculator := NewBlobDigestCalculator()
	uploadResp, err := client.UploadChunk(ctx, *startRes.Location, streaming.NopCloser(bytes.NewReader(blob)), calculator, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *uploadResp.Location)
	_, err = client.CancelUpload(ctx, *uploadResp.Location, nil)
	require.NoError(t, err)
}

func TestNewBlobClient(t *testing.T) {
	client, err := NewBlobClient("test", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	wrongCloudConfig := cloud.Configuration{
		ActiveDirectoryAuthorityHost: "test", Services: map[cloud.ServiceName]cloud.ServiceConfiguration{},
	}
	client, err = NewBlobClient("test", nil, &BlobClientOptions{ClientOptions: azcore.ClientOptions{Cloud: wrongCloudConfig}})
	require.Errorf(t, err, "provided Cloud field is missing Azure Container Registry configuration")
}
