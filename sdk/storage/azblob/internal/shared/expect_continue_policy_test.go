// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
)

// mockTransport captures the Expect header as seen by the transport layer.
type mockTransport struct {
	statusCode    int
	lastExpectHdr string
}

func (m *mockTransport) Do(req *http.Request) (*http.Response, error) {
	m.lastExpectHdr = req.Header.Get("Expect")
	return &http.Response{
		StatusCode: m.statusCode,
		Header:     http.Header{},
		Body:       http.NoBody,
		Request:    req,
	}, nil
}

func newTestPipeline(ecPolicy *ExpectContinuePolicy, transport *mockTransport) runtime.Pipeline {
	return runtime.NewPipeline("test", "v1.0.0",
		runtime.PipelineOptions{
			PerRetry: []policy.Policy{ecPolicy},
		},
		&policy.ClientOptions{
			Transport: transport,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		},
	)
}

func newPutRequest(contentLength int64) (*policy.Request, error) {
	req, err := runtime.NewRequest(context.Background(), http.MethodPut, "https://example.blob.core.windows.net/container/blob")
	if err != nil {
		return nil, err
	}
	if contentLength > 0 {
		body := strings.NewReader(strings.Repeat("x", int(contentLength)))
		err = req.SetBody(readSeekCloser{body}, "application/octet-stream")
		if err != nil {
			return nil, err
		}
		req.Raw().ContentLength = contentLength
	}
	return req, nil
}

type readSeekCloser struct {
	*strings.Reader
}

func (r readSeekCloser) Close() error { return nil }

func TestExpectContinue_AppliedAboveThreshold(t *testing.T) {
	t.Setenv(EnvExpectContinueDisabled, "")
	transport := &mockTransport{statusCode: http.StatusCreated}
	ecPolicy := NewExpectContinuePolicy()
	pl := newTestPipeline(ecPolicy, transport)

	// 8 MiB + 1 byte = above threshold
	req, err := newPutRequest(expectContinueThreshold + 1)
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, "100-continue", transport.lastExpectHdr)
}

func TestExpectContinue_AppliedAtThreshold(t *testing.T) {
	t.Setenv(EnvExpectContinueDisabled, "")
	transport := &mockTransport{statusCode: http.StatusCreated}
	ecPolicy := NewExpectContinuePolicy()
	pl := newTestPipeline(ecPolicy, transport)

	// Exactly 8 MiB = at threshold, should get header
	req, err := newPutRequest(expectContinueThreshold)
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, "100-continue", transport.lastExpectHdr)
}

func TestExpectContinue_NotAppliedBelowThreshold(t *testing.T) {
	t.Setenv(EnvExpectContinueDisabled, "")
	transport := &mockTransport{statusCode: http.StatusCreated}
	ecPolicy := NewExpectContinuePolicy()
	pl := newTestPipeline(ecPolicy, transport)

	req, err := newPutRequest(1024)
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Empty(t, transport.lastExpectHdr)
}

func TestExpectContinue_NotAppliedToGET(t *testing.T) {
	t.Setenv(EnvExpectContinueDisabled, "")
	transport := &mockTransport{statusCode: http.StatusOK}
	ecPolicy := NewExpectContinuePolicy()
	pl := newTestPipeline(ecPolicy, transport)

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://example.blob.core.windows.net/container/blob")
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Empty(t, transport.lastExpectHdr)
}

func TestExpectContinue_NotAppliedZeroContentLength(t *testing.T) {
	t.Setenv(EnvExpectContinueDisabled, "")
	transport := &mockTransport{statusCode: http.StatusCreated}
	ecPolicy := NewExpectContinuePolicy()
	pl := newTestPipeline(ecPolicy, transport)

	req, err := runtime.NewRequest(context.Background(), http.MethodPut, "https://example.blob.core.windows.net/container/blob")
	require.NoError(t, err)
	req.Raw().ContentLength = 0
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Empty(t, transport.lastExpectHdr)
}

func TestExpectContinue_DisabledByEnvVar1(t *testing.T) {
	t.Setenv(EnvExpectContinueDisabled, "1")
	transport := &mockTransport{statusCode: http.StatusCreated}
	ecPolicy := NewExpectContinuePolicy()
	pl := newTestPipeline(ecPolicy, transport)

	req, err := newPutRequest(expectContinueThreshold + 1)
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Empty(t, transport.lastExpectHdr)
}

func TestExpectContinue_DisabledByEnvVarTrue(t *testing.T) {
	t.Setenv(EnvExpectContinueDisabled, "true")
	transport := &mockTransport{statusCode: http.StatusCreated}
	ecPolicy := NewExpectContinuePolicy()
	pl := newTestPipeline(ecPolicy, transport)

	req, err := newPutRequest(expectContinueThreshold + 1)
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Empty(t, transport.lastExpectHdr)
}

func TestExpectContinue_NotDisabledByOtherEnvValues(t *testing.T) {
	t.Setenv(EnvExpectContinueDisabled, "false")
	transport := &mockTransport{statusCode: http.StatusCreated}
	ecPolicy := NewExpectContinuePolicy()
	pl := newTestPipeline(ecPolicy, transport)

	req, err := newPutRequest(expectContinueThreshold + 1)
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, "100-continue", transport.lastExpectHdr)
}
