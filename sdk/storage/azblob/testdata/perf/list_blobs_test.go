// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/require"
)

type blockingUploadTransport struct {
	started chan struct{}
	release chan struct{}
	active  atomic.Int32
	maximum atomic.Int32
	calls   atomic.Int32
}

func (t *blockingUploadTransport) Do(req *http.Request) (*http.Response, error) {
	t.calls.Add(1)
	active := t.active.Add(1)
	defer t.active.Add(-1)
	for {
		maximum := t.maximum.Load()
		if active <= maximum || t.maximum.CompareAndSwap(maximum, active) {
			break
		}
	}
	t.started <- struct{}{}

	select {
	case <-t.release:
		return uploadResponse(req), nil
	case <-req.Context().Done():
		return nil, req.Context().Err()
	}
}

type cancelingUploadTransport struct {
	started  chan struct{}
	release  chan struct{}
	failOnce sync.Once
	canceled atomic.Int32
	calls    atomic.Int32
}

func (t *cancelingUploadTransport) Do(req *http.Request) (*http.Response, error) {
	t.calls.Add(1)
	t.started <- struct{}{}
	<-t.release

	var uploadErr error
	t.failOnce.Do(func() {
		uploadErr = errSeedUpload
	})
	if uploadErr != nil {
		return nil, uploadErr
	}

	<-req.Context().Done()
	t.canceled.Add(1)
	return nil, req.Context().Err()
}

type recordingUploadTransport struct {
	calls atomic.Int32
}

func (t *recordingUploadTransport) Do(req *http.Request) (*http.Response, error) {
	t.calls.Add(1)
	return uploadResponse(req), nil
}

var errSeedUpload = errors.New("seed upload failed")

func uploadResponse(req *http.Request) *http.Response {
	return &http.Response{
		Request:    req,
		Status:     "Created",
		StatusCode: http.StatusCreated,
		Header:     http.Header{},
		Body:       http.NoBody,
	}
}

func newListTestContainerClient(t *testing.T, transport policy.Transporter) *container.Client {
	t.Helper()
	client, err := container.NewClientWithNoCredential("https://example.test/container", &container.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: transport,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		},
	})
	require.NoError(t, err)
	return client
}

func waitForUploads(t *testing.T, started <-chan struct{}, count int) {
	t.Helper()
	for range count {
		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for upload worker to start")
		}
	}
}

func TestSeedBlobsLimitsWorkers(t *testing.T) {
	const (
		count       = 8
		parallelism = 3
	)
	transport := &blockingUploadTransport{
		started: make(chan struct{}, count),
		release: make(chan struct{}),
	}
	client := newListTestContainerClient(t, transport)
	test := &listTestGlobal{blobPrefix: "blob"}
	done := make(chan error, 1)
	go func() {
		done <- test.seedBlobs(context.Background(), client, count, parallelism)
	}()

	waitForUploads(t, transport.started, parallelism)
	callsAtLimit := transport.calls.Load()
	maximumAtLimit := transport.maximum.Load()
	close(transport.release)

	require.NoError(t, <-done)
	require.Equal(t, int32(parallelism), callsAtLimit)
	require.Equal(t, int32(parallelism), maximumAtLimit)
	require.Equal(t, int32(count), transport.calls.Load())
	require.LessOrEqual(t, transport.maximum.Load(), int32(parallelism))
}

func TestSeedBlobsCancelsWorkersAfterError(t *testing.T) {
	const parallelism = 3
	transport := &cancelingUploadTransport{
		started: make(chan struct{}, parallelism),
		release: make(chan struct{}),
	}
	client := newListTestContainerClient(t, transport)
	test := &listTestGlobal{blobPrefix: "blob"}
	done := make(chan error, 1)
	go func() {
		done <- test.seedBlobs(context.Background(), client, 10, parallelism)
	}()

	waitForUploads(t, transport.started, parallelism)
	close(transport.release)

	require.ErrorIs(t, <-done, errSeedUpload)
	require.Equal(t, int32(parallelism), transport.calls.Load())
	require.Equal(t, int32(parallelism-1), transport.canceled.Load())
}

func TestSeedBlobsWithCanceledContext(t *testing.T) {
	transport := &recordingUploadTransport{}
	client := newListTestContainerClient(t, transport)
	test := &listTestGlobal{blobPrefix: "blob"}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := test.seedBlobs(ctx, client, 5, 2)

	require.ErrorIs(t, err, context.Canceled)
	require.Zero(t, transport.calls.Load())
}
