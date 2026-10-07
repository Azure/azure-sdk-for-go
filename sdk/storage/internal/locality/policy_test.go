// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package locality

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
)

// sent records where each attempt went: the URL host it was sent to and its Host header.
type sent struct {
	urlHost string
	host    string
}

// recordingTransport records every attempt and replies with the queued status codes.
type recordingTransport struct {
	mu       sync.Mutex
	attempts []sent
	statuses []int
}

func (tr *recordingTransport) Do(req *http.Request) (*http.Response, error) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.attempts = append(tr.attempts, sent{urlHost: req.URL.Host, host: req.Host})
	status := http.StatusOK
	if len(tr.statuses) > 0 {
		status, tr.statuses = tr.statuses[0], tr.statuses[1:]
	}
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader("")), Header: http.Header{}, Request: req}, nil
}

func newPipeline(tr policy.Transporter) runtime.Pipeline {
	return runtime.NewPipeline("test", "v0.0.0",
		runtime.PipelineOptions{PerCall: []policy.Policy{NewPolicy()}},
		&policy.ClientOptions{
			Transport: tr,
			Retry:     policy.RetryOptions{MaxRetries: 2, RetryDelay: 1, MaxRetryDelay: 1},
		})
}

func send(t *testing.T, pl runtime.Pipeline, ctx context.Context, rawURL string) *http.Request {
	t.Helper()
	req, err := runtime.NewRequest(ctx, http.MethodGet, rawURL)
	require.NoError(t, err)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	return req.Raw()
}

func TestWithEndpoint(t *testing.T) {
	require.Equal(t, "", Endpoint(context.Background()))
	require.Equal(t, "https://layout.example", Endpoint(WithEndpoint(context.Background(), "https://layout.example")))

	ctx := context.Background()
	require.Equal(t, ctx, WithEndpoint(ctx, ""), "an empty endpoint leaves the context unchanged")
}

func TestRewritesHostAndKeepsAccountHostHeader(t *testing.T) {
	for _, endpoint := range []string{
		"layout.blob.core.windows.net",
		"https://layout.blob.core.windows.net",
		"https://layout.blob.core.windows.net/",
		"https://layout.blob.core.windows.net/ignored/path?q=1#frag",
	} {
		t.Run(endpoint, func(t *testing.T) {
			tr := &recordingTransport{}
			raw := send(t, newPipeline(tr), WithEndpoint(context.Background(), endpoint),
				"https://account.blob.core.windows.net/container/blob?snapshot=x")

			require.Equal(t, []sent{{urlHost: "layout.blob.core.windows.net", host: "account.blob.core.windows.net"}}, tr.attempts)
			require.Equal(t, "/container/blob", raw.URL.Path, "only the host is rewritten")
			require.Equal(t, "snapshot=x", raw.URL.RawQuery, "only the host is rewritten")
		})
	}
}

func TestRewritesPort(t *testing.T) {
	tr := &recordingTransport{}
	send(t, newPipeline(tr), WithEndpoint(context.Background(), "https://10.0.0.5:8443"),
		"https://127.0.0.1:10000/devstoreaccount1/container/blob")
	require.Equal(t, []sent{{urlHost: "10.0.0.5:8443", host: "127.0.0.1:10000"}}, tr.attempts)
}

func TestNoEndpointIsNoOp(t *testing.T) {
	tr := &recordingTransport{}
	raw := send(t, newPipeline(tr), context.Background(), "https://account.blob.core.windows.net/container/blob")
	require.Equal(t, "account.blob.core.windows.net", raw.URL.Host)
	// http.NewRequest already sets Host to the URL host; the policy leaves both alone
	require.Equal(t, []sent{{urlHost: "account.blob.core.windows.net", host: "account.blob.core.windows.net"}}, tr.attempts)
}

func TestRetriesStayOnLayoutEndpoint(t *testing.T) {
	tr := &recordingTransport{statuses: []int{http.StatusServiceUnavailable, http.StatusServiceUnavailable, http.StatusOK}}
	send(t, newPipeline(tr), WithEndpoint(context.Background(), "layout.blob.core.windows.net"),
		"https://account.blob.core.windows.net/container/blob")

	want := sent{urlHost: "layout.blob.core.windows.net", host: "account.blob.core.windows.net"}
	require.Equal(t, []sent{want, want, want}, tr.attempts, "every attempt reaches the layout endpoint with the account Host")
}

func TestConcurrentRequestsAreIndependent(t *testing.T) {
	tr := &recordingTransport{}
	pl := newPipeline(tr)
	var wg sync.WaitGroup
	errs := make(chan error, 100)
	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()
			if i%2 == 0 {
				ctx = WithEndpoint(ctx, "layout.blob.core.windows.net")
			}
			req, err := runtime.NewRequest(ctx, http.MethodGet, "https://account.blob.core.windows.net/c/b")
			if err != nil {
				errs <- err
				return
			}
			resp, err := pl.Do(req)
			if err != nil {
				errs <- err
				return
			}
			_ = resp.Body.Close()
			routed := req.Raw().URL.Host == "layout.blob.core.windows.net"
			if routed != (i%2 == 0) {
				errs <- errors.New("a request was routed by another request's endpoint")
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
}
