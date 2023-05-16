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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestBlobClient_CompleteUpload(t *testing.T) {
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
	calculator := NewBlobDigestCalculator()
	uploadResp, err := client.UploadChunk(ctx, *startRes.Location, bytes.NewReader(blob), calculator, nil)
	require.NoError(t, err)
	completeResp, err := client.CompleteUpload(ctx, *uploadResp.Location, calculator, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *completeResp.DockerContentDigest)
}

func TestBlobClient_UploadChunk(t *testing.T) {
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
	calculator := NewBlobDigestCalculator()
	uploadResp, err := client.UploadChunk(ctx, *startRes.Location, bytes.NewReader(blob), calculator, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *uploadResp.Location)
	_, err = client.CancelUpload(ctx, *uploadResp.Location, nil)
	require.NoError(t, err)
}

func TestBlobClient_CompleteUpload_uploadByChunk(t *testing.T) {
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
	calculator := NewBlobDigestCalculator()
	oriReader := bytes.NewReader(blob)
	firstPart := io.NewSectionReader(oriReader, int64(0), int64(len(blob)/2))
	secondPart := io.NewSectionReader(oriReader, int64(len(blob)/2), int64(len(blob)-len(blob)/2))
	uploadResp, err := client.UploadChunk(ctx, *startRes.Location, firstPart, calculator, &BlobClientUploadChunkOptions{RangeStart: to.Ptr(int32(0)), RangeEnd: to.Ptr(int32(len(blob)/2 - 1))})
	require.NoError(t, err)
	require.NotEmpty(t, *uploadResp.Location)
	uploadResp, err = client.UploadChunk(ctx, *uploadResp.Location, secondPart, calculator, &BlobClientUploadChunkOptions{RangeStart: to.Ptr(int32(len(blob) / 2)), RangeEnd: to.Ptr(int32(len(blob) - 1))})
	require.NoError(t, err)
	require.NotEmpty(t, *uploadResp.Location)
	completeResp, err := client.CompleteUpload(ctx, *uploadResp.Location, calculator, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *completeResp.DockerContentDigest)
}

func TestNewBlobClient(t *testing.T) {
	client, err := NewBlobClient("test", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	wrongCloudConfig := cloud.Configuration{
		ActiveDirectoryAuthorityHost: "test", Services: map[cloud.ServiceName]cloud.ServiceConfiguration{},
	}
	_, err = NewBlobClient("test", nil, &BlobClientOptions{ClientOptions: azcore.ClientOptions{Cloud: wrongCloudConfig}})
	require.Errorf(t, err, "provided Cloud field is missing Azure Container Registry configuration")
}

func TestBlobDigestCalculator_saveAndRestoreState(t *testing.T) {
	calculator := NewBlobDigestCalculator()
	calculator.restoreState()
	calculator.saveState()
	calculator.restoreState()
	calculator.h.Write([]byte("test1"))
	sum := calculator.h.Sum(nil)
	calculator.saveState()
	calculator.h.Write([]byte("test2"))
	require.NotEqual(t, sum, calculator.h.Sum(nil))
	calculator.restoreState()
	require.Equal(t, sum, calculator.h.Sum(nil))
}

func TestBlobClient_CompleteUpload_uploadByChunkFailOver(t *testing.T) {
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
	calculator := NewBlobDigestCalculator()
	oriReader := bytes.NewReader(blob)
	firstPart := io.NewSectionReader(oriReader, int64(0), int64(len(blob)/2))
	secondPart := io.NewSectionReader(oriReader, int64(len(blob)/2), int64(len(blob)-len(blob)/2))
	uploadResp, err := client.UploadChunk(ctx, *startRes.Location, firstPart, calculator, &BlobClientUploadChunkOptions{RangeStart: to.Ptr(int32(0)), RangeEnd: to.Ptr(int32(len(blob)/2 - 1))})
	require.NoError(t, err)
	require.NotEmpty(t, *uploadResp.Location)
	sum := calculator.h.Sum(nil)
	// upload with a wrong range to test fail over
	_, err = client.UploadChunk(ctx, *uploadResp.Location, secondPart, calculator, &BlobClientUploadChunkOptions{RangeStart: to.Ptr(int32(-1)), RangeEnd: to.Ptr(int32(-1))})
	require.Error(t, err)
	require.Equal(t, sum, calculator.h.Sum(nil))
	uploadResp, err = client.UploadChunk(ctx, *uploadResp.Location, secondPart, calculator, &BlobClientUploadChunkOptions{RangeStart: to.Ptr(int32(len(blob) / 2)), RangeEnd: to.Ptr(int32(len(blob) - 1))})
	require.NoError(t, err)
	require.NotEmpty(t, *uploadResp.Location)
	completeResp, err := client.CompleteUpload(ctx, *uploadResp.Location, calculator, nil)
	require.NoError(t, err)
	require.NotEmpty(t, *completeResp.DockerContentDigest)
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

func TestBlobClient_GetChunk(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	ctx := context.Background()
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	digest := "sha256:042a816809aac8d0f7d7cacac7965782ee2ecac3f21bcf9f24b1de1a7387b769"
	calculator := NewBlobDigestCalculator()
	res, err := client.GetChunk(ctx, "alpine", digest, "bytes=0-999", calculator, nil)
	require.NoError(t, err)
	require.Equal(t, int64(1000), *res.ContentLength)
	res, err = client.GetChunk(ctx, "alpine", digest, "bytes=1000-1471", calculator, nil)
	require.NoError(t, err)
	require.Equal(t, int64(472), *res.ContentLength)
}

func TestBlobClient_GetChunk_wrongDigest(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusPartialContent), mock.WithHeader("Content-Range", "bytes 0-9/20"), mock.WithBody([]byte("0123456789")))
	srv.AppendResponse(mock.WithStatusCode(http.StatusPartialContent), mock.WithHeader("Content-Range", "bytes 10-19/20"), mock.WithBody([]byte("0123456789")))

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	blobClient := &BlobClient{
		srv.URL(),
		pl,
	}
	ctx := context.Background()
	calculator := NewBlobDigestCalculator()
	res, err := blobClient.GetChunk(ctx, "test", "sha256:test", "bytes=0-9", calculator, nil)
	require.NoError(t, err)
	require.Equal(t, int64(10), *res.ContentLength)
	_, err = blobClient.GetChunk(ctx, "test", "sha256:test", "bytes=10-19", calculator, nil)
	require.Error(t, err)
}
