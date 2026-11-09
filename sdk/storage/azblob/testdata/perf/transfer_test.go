// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/require"
)

type transferTransport struct {
	calls atomic.Int32
}

func (t *transferTransport) Do(req *http.Request) (*http.Response, error) {
	t.calls.Add(1)
	if req.Method == http.MethodHead {
		header := http.Header{"Content-Length": []string{"4"}, "X-Ms-Blob-Type": []string{"BlockBlob"}}
		return &http.Response{Request: req, StatusCode: http.StatusOK, Status: "OK", Header: header, Body: http.NoBody}, nil
	}
	if req.Method == http.MethodGet {
		body := "data"
		header := http.Header{
			"Content-Length":           []string{"4"},
			"Content-Range":            []string{"bytes 0-3/4"},
			"X-Ms-Blob-Content-Length": []string{"4"},
			"X-Ms-Blob-Type":           []string{"BlockBlob"},
		}
		return &http.Response{Request: req, StatusCode: http.StatusPartialContent, Status: "Partial Content", Header: header, Body: io.NopCloser(strings.NewReader(body))}, nil
	}
	return setupResponse(req, http.StatusCreated), nil
}

func newTransferContainer(t *testing.T, transport policy.Transporter) *container.Client {
	t.Helper()
	client, err := container.NewClientWithNoCredential("https://example.test/container", &container.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: transport, Retry: policy.RetryOptions{MaxRetries: -1}},
	})
	require.NoError(t, err)
	return client
}

func TestUploadMethods(t *testing.T) {
	savedMethod := uploadMethod
	t.Cleanup(func() { uploadMethod = savedMethod })
	for _, method := range []string{"single", "buffer", "stream"} {
		t.Run(method, func(t *testing.T) {
			transport := &transferTransport{}
			client := newTransferContainer(t, transport)
			global := &uploadTestGlobal{size: 4, payload: []byte("data"), randomSeed: []byte("data")}
			test := &uploadPerfTest{uploadTestGlobal: global, blobClient: client.NewBlockBlobClient("blob")}
			uploadMethod = method

			require.NoError(t, test.Run(context.Background()))
			require.Positive(t, transport.calls.Load())
		})
	}
}

func TestUploadRunRejectsUnknownMethod(t *testing.T) {
	savedMethod := uploadMethod
	t.Cleanup(func() { uploadMethod = savedMethod })
	uploadMethod = "invalid"
	err := (&uploadPerfTest{}).Run(context.Background())
	require.ErrorContains(t, err, "unknown --upload-method")
}

func TestDownloadMethods(t *testing.T) {
	savedMethod := downloadMethod
	t.Cleanup(func() { downloadMethod = savedMethod })
	for _, method := range []string{"stream", "buffer"} {
		t.Run(method, func(t *testing.T) {
			transport := &transferTransport{}
			client := newTransferContainer(t, transport)
			global := &downloadTestGlobal{size: 4}
			test := &downloadPerfTest{
				downloadTestGlobal: global,
				PerfTestOptions:    perf.PerfTestOptions{},
				blobClient:         client.NewBlockBlobClient("blob").BlobClient(),
				buffer:             make([]byte, 4),
			}
			downloadMethod = method

			require.NoError(t, test.Run(context.Background()))
			require.Positive(t, transport.calls.Load())
			if method == "buffer" {
				require.Equal(t, []byte("data"), test.buffer)
			}
		})
	}
}

func TestDownloadRunRejectsUnknownMethod(t *testing.T) {
	savedMethod := downloadMethod
	t.Cleanup(func() { downloadMethod = savedMethod })
	downloadMethod = "invalid"
	err := (&downloadPerfTest{}).Run(context.Background())
	require.ErrorContains(t, err, "unknown --download-method")
}

func TestRandomStreamSeekAndBounds(t *testing.T) {
	stream := newRandomStream([]byte("ab"), 5)
	data, err := io.ReadAll(stream)
	require.NoError(t, err)
	require.Equal(t, []byte("ababa"), data)

	position, err := stream.Seek(1, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(1), position)
	data, err = io.ReadAll(stream)
	require.NoError(t, err)
	require.Equal(t, []byte("baba"), data)

	_, err = stream.Seek(-1, io.SeekStart)
	require.Error(t, err)
	_, err = stream.Seek(-10, io.SeekCurrent)
	require.Error(t, err)
	_, err = stream.Seek(-10, io.SeekEnd)
	require.Error(t, err)
	_, err = stream.Seek(0, 99)
	require.Error(t, err)
	position, err = stream.Seek(-1, io.SeekEnd)
	require.NoError(t, err)
	require.Equal(t, int64(4), position)
	require.NoError(t, stream.Close())
	require.NoError(t, NopCloser(strings.NewReader("data")).Close())
	require.Error(t, validateTransferOptions("upload", -1))
}
