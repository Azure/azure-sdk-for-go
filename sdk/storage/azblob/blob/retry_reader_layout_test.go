// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

// flakyBody yields a few bytes and then fails, which is what drives the retry reader to re-issue
// the read.
type flakyBody struct {
	remaining int
}

func (b *flakyBody) Read(p []byte) (int, error) {
	if b.remaining > 0 && len(p) > 0 {
		b.remaining--
		p[0] = 'a'
		return 1, nil
	}
	// the retry reader re-issues on an unexpected EOF, which is what a dropped connection
	// surfaces as mid-body
	return 0, io.ErrUnexpectedEOF
}

func (b *flakyBody) Close() error { return nil }

// layoutRetryTransport records the host of every read and fails the first one mid-body.
type layoutRetryTransport struct {
	mu sync.Mutex
	// urlHosts is the host each request was actually sent to.
	urlHosts []string
	// hostHeaders is the Host header of each request, which must stay the account host.
	hostHeaders []string
	calls       int
}

func (t *layoutRetryTransport) Do(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	t.calls++
	call := t.calls
	t.urlHosts = append(t.urlHosts, req.URL.Host)
	t.hostHeaders = append(t.hostHeaders, req.Host)
	t.mu.Unlock()

	start, end := requestedRange(rawHeaderValue(req.Header, "x-ms-range"))
	if end < start {
		start, end = 0, 9
	}
	length := end - start + 1

	header := http.Header{}
	header.Set("Content-Length", strconv.FormatInt(length, 10))
	header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, 10))
	header.Set("ETag", `"retry-etag"`)

	var body io.ReadCloser
	if call == 1 {
		// fail partway through so the retry reader has to issue a second request
		body = &flakyBody{remaining: 2}
	} else {
		body = io.NopCloser(bytes.NewReader(make([]byte, length)))
	}

	return &http.Response{
		Request:    req,
		StatusCode: http.StatusPartialContent,
		Header:     header,
		Body:       body,
	}, nil
}

func (t *layoutRetryTransport) snapshot() ([]string, []string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]string(nil), t.urlHosts...), append([]string(nil), t.hostHeaders...)
}

// A read that the layout routed to a particular endpoint must stay on that endpoint when the
// retry reader re-issues it. Before the layout endpoint was carried on the response, the retry
// silently fell back to the account endpoint, which is the locality the routing existed to avoid.
func TestRetryReaderPreservesLayoutEndpoint(t *testing.T) {
	const (
		accountHost = "fakeaccount.blob.core.windows.net"
		layoutHost  = "locality1.blob.core.windows.net"
	)

	tr := &layoutRetryTransport{}
	client, err := NewClientWithNoCredential("https://"+accountHost+"/container/blob", &ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: tr,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		},
	})
	require.NoError(t, err)

	ctx := context.Background()
	dr, err := client.DownloadStream(ctx, &DownloadStreamOptions{
		Range:          HTTPRange{Offset: 0, Count: 10},
		LayoutEndpoint: "https://" + layoutHost,
	})
	require.NoError(t, err)

	body := dr.NewRetryReader(ctx, &RetryReaderOptions{MaxRetries: 2})
	_, err = io.ReadAll(body)
	require.NoError(t, err)
	require.NoError(t, body.Close())

	urlHosts, hostHeaders := tr.snapshot()
	require.GreaterOrEqual(t, len(urlHosts), 2, "the failed read must have been retried")
	for i, h := range urlHosts {
		require.Equal(t, layoutHost, h, "request %d must stay on the layout endpoint", i)
	}
	for i, h := range hostHeaders {
		require.Equal(t, accountHost, h, "request %d must keep the account Host header", i)
	}
}

// Without a layout endpoint the retry must go back to the account endpoint, unchanged.
func TestRetryReaderWithoutLayoutEndpointUsesAccountHost(t *testing.T) {
	const accountHost = "fakeaccount.blob.core.windows.net"

	tr := &layoutRetryTransport{}
	client, err := NewClientWithNoCredential("https://"+accountHost+"/container/blob", &ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: tr,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		},
	})
	require.NoError(t, err)

	ctx := context.Background()
	dr, err := client.DownloadStream(ctx, &DownloadStreamOptions{Range: HTTPRange{Offset: 0, Count: 10}})
	require.NoError(t, err)

	body := dr.NewRetryReader(ctx, &RetryReaderOptions{MaxRetries: 2})
	_, err = io.ReadAll(body)
	require.NoError(t, err)
	require.NoError(t, body.Close())

	urlHosts, _ := tr.snapshot()
	require.GreaterOrEqual(t, len(urlHosts), 2)
	for i, h := range urlHosts {
		require.Equal(t, accountHost, h, "request %d must use the account endpoint", i)
	}
}
