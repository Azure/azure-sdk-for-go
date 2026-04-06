// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestBlobCustomClient(t *testing.T) {
	repository, digest := buildImage(t)
	blobDigest := blobDigest(t, repository, digest)

	t.Run("CompleteUpload", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		getRes, err := client.GetBlob(ctx, repository, blobDigest, nil)
		require.NoError(t, err)
		blob, err := io.ReadAll(getRes.BlobData)
		require.NoError(t, err)
		startRes, err := client.StartUpload(ctx, "hello-world", nil)
		require.NoError(t, err)
		calculator := NewBlobDigestCalculator()
		if recording.GetRecordMode() == recording.PlaybackMode {
			calculator.h = &fakeHash{}
		}
		uploadResp, err := client.UploadChunk(ctx, *startRes.Location, bytes.NewReader(blob), calculator, nil)
		require.NoError(t, err)
		completeResp, err := client.CompleteUpload(ctx, *uploadResp.Location, calculator, nil)
		require.NoError(t, err)
		require.NotEmpty(t, *completeResp.DockerContentDigest)
	})

	t.Run("CompleteUpload_uploadByChunk", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		getRes, err := client.GetBlob(ctx, repository, blobDigest, nil)
		require.NoError(t, err)
		blob, err := io.ReadAll(getRes.BlobData)
		require.NoError(t, err)
		startRes, err := client.StartUpload(ctx, "hello-world", nil)
		require.NoError(t, err)
		calculator := NewBlobDigestCalculator()
		if recording.GetRecordMode() == recording.PlaybackMode {
			calculator.h = &fakeHash{}
		}
		oriReader := bytes.NewReader(blob)
		size := int64(len(blob))
		chunkSize := int64(736)
		current := int64(0)
		location := *startRes.Location
		for {
			end := current + chunkSize
			if end > size {
				end = size
			}
			chunkReader := io.NewSectionReader(oriReader, current, end-current)
			uploadResp, err := client.UploadChunk(ctx, location, chunkReader, calculator, &BlobClientUploadChunkOptions{RangeStart: to.Ptr(int32(current)), RangeEnd: to.Ptr(int32(end - 1))})
			require.NoError(t, err)
			require.NotEmpty(t, *uploadResp.Location)
			location = *uploadResp.Location
			current = end
			if current >= size {
				break
			}
		}
		completeResp, err := client.CompleteUpload(ctx, location, calculator, nil)
		require.NoError(t, err)
		require.NotEmpty(t, *completeResp.DockerContentDigest)
	})

	t.Run("CompleteUpload_uploadByChunkFailOver", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		getRes, err := client.GetBlob(ctx, repository, blobDigest, nil)
		require.NoError(t, err)
		blob, err := io.ReadAll(getRes.BlobData)
		require.NoError(t, err)
		startRes, err := client.StartUpload(ctx, "hello-world", nil)
		require.NoError(t, err)
		calculator := NewBlobDigestCalculator()
		if recording.GetRecordMode() == recording.PlaybackMode {
			calculator.h = &fakeHash{}
		}
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
	})

	t.Run("UploadChunk", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		getRes, err := client.GetBlob(ctx, repository, blobDigest, nil)
		require.NoError(t, err)
		blob, err := io.ReadAll(getRes.BlobData)
		require.NoError(t, err)
		startRes, err := client.StartUpload(ctx, "hello-world", nil)
		require.NoError(t, err)
		calculator := NewBlobDigestCalculator()
		if recording.GetRecordMode() == recording.PlaybackMode {
			calculator.h = &fakeHash{}
		}
		uploadResp, err := client.UploadChunk(ctx, *startRes.Location, bytes.NewReader(blob), calculator, nil)
		require.NoError(t, err)
		require.NotEmpty(t, *uploadResp.Location)
		_, err = client.CancelUpload(ctx, *uploadResp.Location, nil)
		require.NoError(t, err)
	})
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

func TestBlobCustomClient_UploadChunk_retry(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusGatewayTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusGatewayTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))

	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	require.NoError(t, err)
	client := &BlobClient{
		azcoreClient,
		srv.URL(),
	}
	chunkData := bytes.NewReader([]byte("test"))
	calculator := NewBlobDigestCalculator()
	_, err = client.UploadChunk(ctx, "location", chunkData, calculator, nil)
	require.NoError(t, err)
	require.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", fmt.Sprintf("%x", calculator.h.Sum(nil)))
}

// fakeHash is a fake hash.Hash for playback mode
type fakeHash struct{}

func (f *fakeHash) Write(p []byte) (int, error) { return 0, nil }

func (f *fakeHash) Sum(b []byte) []byte { return []byte{0} }

func (f *fakeHash) Reset() {}

func (f *fakeHash) Size() int { return 0 }

func (f *fakeHash) BlockSize() int { return 1 }

func (f *fakeHash) MarshalBinary() ([]byte, error) { return nil, nil }
