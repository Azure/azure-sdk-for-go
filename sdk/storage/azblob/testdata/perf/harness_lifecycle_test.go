// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/require"
)

type harnessTransport struct {
	mu       sync.Mutex
	requests []*http.Request
}

func (t *harnessTransport) Do(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	t.requests = append(t.requests, req.Clone(req.Context()))
	t.mu.Unlock()

	if req.Method == http.MethodGet && req.URL.Query().Get("comp") == "list" {
		body := `<?xml version="1.0" encoding="utf-8"?><EnumerationResults ServiceEndpoint="https://example.test/" ContainerName="container"><Blobs></Blobs><NextMarker /></EnumerationResults>`
		return &http.Response{Request: req, StatusCode: http.StatusOK, Status: "OK", Header: http.Header{"Content-Type": []string{"application/xml"}}, Body: io.NopCloser(strings.NewReader(body))}, nil
	}
	if req.Method == http.MethodGet {
		return (&transferTransport{}).Do(req)
	}
	if req.Method == http.MethodHead {
		return (&transferTransport{}).Do(req)
	}
	status := http.StatusCreated
	if req.Method == http.MethodDelete {
		status = http.StatusAccepted
	}
	return setupResponse(req, status), nil
}

func (t *harnessTransport) count(method string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	count := 0
	for _, req := range t.requests {
		if req.Method == method {
			count++
		}
	}
	return count
}

func (t *harnessTransport) lastListMaxResults() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	for index := len(t.requests) - 1; index >= 0; index-- {
		if t.requests[index].URL.Query().Get("comp") == "list" {
			return t.requests[index].URL.Query().Get("maxresults")
		}
	}
	return ""
}

func installHarnessFactory(t *testing.T, transport policy.Transporter) {
	t.Helper()
	savedFactory := containerClientFactory
	containerClientFactory = func(containerName string, _ *container.ClientOptions) (*container.Client, error) {
		return container.NewClientWithNoCredential("https://example.test/"+containerName, &container.ClientOptions{
			ClientOptions: azcore.ClientOptions{Transport: transport, Retry: policy.RetryOptions{MaxRetries: -1}},
		})
	}
	t.Cleanup(func() { containerClientFactory = savedFactory })
}

func TestUploadHarnessLifecycle(t *testing.T) {
	transport := &harnessTransport{}
	installHarnessFactory(t, transport)
	savedOptions, savedMethod := uploadTestOpts, uploadMethod
	uploadTestOpts.size, uploadMethod = 4, "single"
	t.Cleanup(func() { uploadTestOpts, uploadMethod = savedOptions, savedMethod })

	globalValue, err := NewUploadTest(context.Background(), perf.PerfTestOptions{})
	require.NoError(t, err)
	global := globalValue.(*uploadTestGlobal)
	options := perf.PerfTestOptions{Name: "worker"}
	worker, err := global.NewPerfTest(context.Background(), &options)
	require.NoError(t, err)
	require.NoError(t, worker.Run(context.Background()))
	require.NoError(t, worker.Cleanup(context.Background()))
	require.NoError(t, global.GlobalCleanup(context.Background()))
	require.GreaterOrEqual(t, transport.count(http.MethodPut), 2)
	require.Equal(t, 1, transport.count(http.MethodDelete))
}

func TestUploadBufferSetup(t *testing.T) {
	transport := &harnessTransport{}
	installHarnessFactory(t, transport)
	savedOptions, savedMethod := uploadTestOpts, uploadMethod
	uploadTestOpts.size, uploadMethod = 4, "buffer"
	t.Cleanup(func() { uploadTestOpts, uploadMethod = savedOptions, savedMethod })

	globalValue, err := NewUploadTest(context.Background(), perf.PerfTestOptions{})
	require.NoError(t, err)
	global := globalValue.(*uploadTestGlobal)
	require.Len(t, global.payload, 4)
	require.Empty(t, global.randomSeed)
	require.NoError(t, global.GlobalCleanup(context.Background()))
}

func TestDownloadHarnessLifecycle(t *testing.T) {
	transport := &harnessTransport{}
	installHarnessFactory(t, transport)
	savedOptions, savedMethod := downloadTestOpts, downloadMethod
	downloadTestOpts.size, downloadMethod = 4, "buffer"
	t.Cleanup(func() { downloadTestOpts, downloadMethod = savedOptions, savedMethod })

	globalValue, err := NewDownloadTest(context.Background(), perf.PerfTestOptions{})
	require.NoError(t, err)
	global := globalValue.(*downloadTestGlobal)
	options := perf.PerfTestOptions{Name: "worker"}
	worker, err := global.NewPerfTest(context.Background(), &options)
	require.NoError(t, err)
	require.NoError(t, worker.Run(context.Background()))
	require.NoError(t, worker.Cleanup(context.Background()))
	require.NoError(t, global.GlobalCleanup(context.Background()))
	require.Positive(t, transport.count(http.MethodGet))
	require.Equal(t, 1, transport.count(http.MethodDelete))
}

func TestListHarnessLifecycleAndPageSize(t *testing.T) {
	transport := &harnessTransport{}
	installHarnessFactory(t, transport)
	savedOptions, savedPageSize := listTestOpts, listPageSize
	listTestOpts, listPageSize = listTestOptions{count: 2, parallelism: 1}, 17
	t.Cleanup(func() { listTestOpts, listPageSize = savedOptions, savedPageSize })

	globalValue, err := NewListTest(context.Background(), perf.PerfTestOptions{})
	require.NoError(t, err)
	global := globalValue.(*listTestGlobal)
	options := perf.PerfTestOptions{Name: "worker"}
	worker, err := global.NewPerfTest(context.Background(), &options)
	require.NoError(t, err)
	require.NoError(t, worker.Run(context.Background()))
	require.NoError(t, worker.Cleanup(context.Background()))
	require.NoError(t, global.GlobalCleanup(context.Background()))
	require.Equal(t, "17", transport.lastListMaxResults())
	require.Equal(t, 1, transport.count(http.MethodDelete))
}

type pagingTransport struct {
	pages int
	fail  bool
}

func (t *pagingTransport) Do(req *http.Request) (*http.Response, error) {
	t.pages++
	if t.fail {
		return setupResponse(req, http.StatusInternalServerError), nil
	}
	next := ""
	if req.URL.Query().Get("marker") == "" {
		next = "next"
	}
	body := `<?xml version="1.0" encoding="utf-8"?><EnumerationResults ServiceEndpoint="https://example.test/" ContainerName="container"><Blobs></Blobs><NextMarker>` + next + `</NextMarker></EnumerationResults>`
	return &http.Response{Request: req, StatusCode: http.StatusOK, Status: "OK", Header: http.Header{"Content-Type": []string{"application/xml"}}, Body: io.NopCloser(strings.NewReader(body))}, nil
}

func TestListRunMultiplePagesAndError(t *testing.T) {
	savedPageSize := listPageSize
	listPageSize = 0
	t.Cleanup(func() { listPageSize = savedPageSize })

	transport := &pagingTransport{}
	test := &listPerfTest{containerClient: newTransferContainer(t, transport)}
	require.NoError(t, test.Run(context.Background()))
	require.Equal(t, 2, transport.pages)

	transport = &pagingTransport{fail: true}
	test.containerClient = newTransferContainer(t, transport)
	require.Error(t, test.Run(context.Background()))
	require.Equal(t, 1, transport.pages)
}

var errContainerFactory = errors.New("container factory failed")

func TestHarnessFactoryErrors(t *testing.T) {
	savedFactory := containerClientFactory
	containerClientFactory = func(string, *container.ClientOptions) (*container.Client, error) {
		return nil, errContainerFactory
	}
	t.Cleanup(func() { containerClientFactory = savedFactory })

	_, err := NewUploadTest(context.Background(), perf.PerfTestOptions{})
	require.ErrorIs(t, err, errContainerFactory)
	_, err = NewDownloadTest(context.Background(), perf.PerfTestOptions{})
	require.ErrorIs(t, err, errContainerFactory)
	_, err = NewListTest(context.Background(), perf.PerfTestOptions{})
	require.ErrorIs(t, err, errContainerFactory)
	options := perf.PerfTestOptions{}
	_, err = (&uploadTestGlobal{}).NewPerfTest(context.Background(), &options)
	require.ErrorIs(t, err, errContainerFactory)
	_, err = (&downloadTestGlobal{}).NewPerfTest(context.Background(), &options)
	require.ErrorIs(t, err, errContainerFactory)
	_, err = (&listTestGlobal{}).NewPerfTest(context.Background(), &options)
	require.ErrorIs(t, err, errContainerFactory)

	for _, global := range []interface {
		GlobalCleanup(context.Context) error
	}{
		&downloadTestGlobal{}, &listTestGlobal{},
	} {
		require.ErrorIs(t, global.GlobalCleanup(context.Background()), errContainerFactory)
	}
}
