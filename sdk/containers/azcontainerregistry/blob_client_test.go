// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

// blobDigest returns the digest of the first blob i.e. layer in the image manifest
func blobDigest(t *testing.T, image, imageDigest string) string {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return fakeDigest
	}
	client, err := NewClient("https://"+testConfig.loginServer, testConfig.credential, nil)
	require.NoError(t, err)
	res, err := client.GetManifest(ctx, image, imageDigest, &ClientGetManifestOptions{Accept: to.Ptr("application/vnd.docker.distribution.manifest.v2+json")})
	require.NoError(t, err)
	reader, err := NewDigestValidationReader(*res.DockerContentDigest, res.ManifestData)
	require.NoError(t, err)
	manifest, err := io.ReadAll(reader)
	require.NoError(t, err)
	blobDigest := string(regexp.MustCompile("(sha256:[a-f0-9]{64})").Find(manifest))
	require.NotEmpty(t, blobDigest)
	_, hash, found := strings.Cut(blobDigest, ":")
	require.True(t, found)
	if recording.GetRecordMode() == recording.RecordingMode {
		require.NoError(t, recording.AddGeneralRegexSanitizer("00", hash, nil))
	}
	return blobDigest
}

func TestBlobClient(t *testing.T) {
	repository, digest := buildImage(t)
	blobDigest := blobDigest(t, repository, digest)

	t.Run("CheckBlobExists", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		res, err := client.CheckBlobExists(ctx, repository, blobDigest, nil)
		require.NoError(t, err)
		require.Equal(t, blobDigest, *res.DockerContentDigest)
	})

	t.Run("CheckChunkExists", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		res, err := client.CheckChunkExists(ctx, repository, blobDigest, "bytes=0-299", nil)
		require.NoError(t, err)
		require.NotEmpty(t, *res.ContentLength)
	})

	t.Run("CheckChunkExists_fail", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		_, err = client.CheckChunkExists(ctx, repository, "wrong digest", "bytes=0-299", nil)
		require.Error(t, err)
	})

	t.Run("CheckBlobExists_fail", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		_, err = client.CheckBlobExists(ctx, repository, "wrong digest", nil)
		require.Error(t, err)
	})

	t.Run("CompleteUpload_wrongDigest", func(t *testing.T) {
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
		uploadResp, err := client.uploadChunk(ctx, *startRes.Location, streaming.NopCloser(bytes.NewReader(blob)), nil)
		require.NoError(t, err)
		_, err = client.completeUpload(ctx, "sha256:00000000", *uploadResp.Location, nil)
		require.Error(t, err)
	})

	t.Run("DeleteBlob_fail", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		_, err = client.DeleteBlob(ctx, repository, "wrong digest", nil)
		require.Error(t, err)
	})

	t.Run("GetBlob", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		res, err := client.GetBlob(ctx, repository, blobDigest, nil)
		require.NoError(t, err)
		require.NotEmpty(t, *res.ContentLength)
		reader, err := NewDigestValidationReader(blobDigest, res.BlobData)
		require.NoError(t, err)
		if recording.GetRecordMode() == recording.PlaybackMode {
			reader.digestValidator = &sha256Validator{&fakeHash{}}
		}
		_, err = io.ReadAll(reader)
		require.NoError(t, err)
	})

	t.Run("GetBlob_fail", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		_, err = client.GetBlob(ctx, repository, "wrong digest", nil)
		require.Error(t, err)
	})

	t.Run("GetChunk", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		chunkSize := 1000
		current := 0
		blob := bytes.NewBuffer(nil)
		for {
			res, err := client.GetChunk(ctx, repository, blobDigest, fmt.Sprintf("bytes=%d-%d", current, current+chunkSize-1), nil)
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
		reader, err := NewDigestValidationReader(blobDigest, blob)
		require.NoError(t, err)
		if recording.GetRecordMode() == recording.PlaybackMode {
			reader.digestValidator = &sha256Validator{&fakeHash{}}
		}
		_, err = io.ReadAll(reader)
		require.NoError(t, err)
	})

	t.Run("GetChunk_fail", func(t *testing.T) {
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		_, err = client.GetChunk(ctx, repository, "wrong digest", "bytes=0-999", nil)
		require.Error(t, err)
	})

	t.Run("MountBlob", func(t *testing.T) {
		repository2, _ := buildImage(t)
		startRecording(t)
		endpoint, cred, options := getEndpointCredAndClientOptions(t)
		client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
		require.NoError(t, err)
		res, err := client.MountBlob(ctx, repository2, repository, blobDigest, nil)
		require.NoError(t, err)
		require.NotEmpty(t, res.Location)
	})
}

func TestBlobClient_CancelUpload(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
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
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.CancelUpload(ctx, "wrong location", nil)
	require.Error(t, err)
}

func TestBlobClient_CheckBlobExists_empty(t *testing.T) {
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.CheckBlobExists(ctx, "", "digest", nil)
	require.Error(t, err)
	_, err = client.CheckBlobExists(ctx, "repository", "", nil)
	require.Error(t, err)
}

func TestBlobClient_CheckChunkExists_empty(t *testing.T) {
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.CheckChunkExists(ctx, "", "digest", "range", nil)
	require.Error(t, err)
	_, err = client.CheckChunkExists(ctx, "name", "", "range", nil)
	require.Error(t, err)
}

func TestBlobClient_DeleteBlob(t *testing.T) {
	repository, imgDigest := buildImage(t)
	blobDigest := blobDigest(t, repository, imgDigest)
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.DeleteBlob(ctx, repository, blobDigest, nil)
	require.NoError(t, err)
}

func TestBlobClient_DeleteBlob_empty(t *testing.T) {
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.DeleteBlob(ctx, "", "digest", nil)
	require.Error(t, err)
	_, err = client.DeleteBlob(ctx, "name", "", nil)
	require.Error(t, err)
}

func TestBlobClient_GetBlob_empty(t *testing.T) {
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.GetBlob(ctx, "", "digest", nil)
	require.Error(t, err)
	_, err = client.GetBlob(ctx, "repository", "", nil)
	require.Error(t, err)
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
	resp, err := client.GetBlob(ctx, "name", "sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", nil)
	require.NoError(t, err)
	reader, err := NewDigestValidationReader("sha256:wrong", resp.BlobData)
	require.NoError(t, err)
	_, err = io.ReadAll(reader)
	require.Error(t, err, ErrMismatchedHash)
}

func TestBlobClient_GetChunk_empty(t *testing.T) {
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.GetChunk(ctx, "", "digest", "bytes=0-999", nil)
	require.Error(t, err)
	_, err = client.GetChunk(ctx, "repository", "", "bytes=0-999", nil)
	require.Error(t, err)
}

func TestBlobClient_GetUploadStatus(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
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
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.GetUploadStatus(ctx, "wrong location", nil)
	require.Error(t, err)
}

func TestBlobClient_MountBlob_fail(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	_, err = client.MountBlob(ctx, "wrong name", "wrong from", "wrong mount", nil)
	require.Error(t, err)
}

func TestBlobClient_MountBlob_empty(t *testing.T) {
	client, err := NewBlobClient("endpoint", nil, nil)
	require.NoError(t, err)
	_, err = client.MountBlob(ctx, "", "from", "mount", nil)
	require.Error(t, err)
}

func TestBlobClient_StartUpload(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client, err := NewBlobClient(endpoint, cred, &BlobClientOptions{ClientOptions: options})
	require.NoError(t, err)
	startRes, err := client.StartUpload(ctx, "hello-world", nil)
	require.NoError(t, err)
	require.NotEmpty(t, *startRes.Location)
	_, err = client.CancelUpload(ctx, *startRes.Location, nil)
	require.NoError(t, err)
}

func TestBlobClient_StartUpload_empty(t *testing.T) {
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
