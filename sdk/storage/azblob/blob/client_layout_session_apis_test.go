// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
)

// Session authentication and Data Locality have different scopes, and this file pins both down.
//
// Session authentication covers five operations: Get Blob, Get Blob Properties, Put Blob,
// Put Block and Put Block List.
//
// Data Locality covers Get Blob and the managed downloads built on it, and nothing else. The
// four operations that gained session authentication did not gain a layout with it.
//
// Where the two do meet, on a locality-routed download chunk, they compose: Data Locality
// rewrites the request URL host to a layout endpoint while keeping the account Host header, and
// the session signature is built from the configured account name and the URL path, never the
// host, so rerouting a chunk cannot invalidate it.

const localityHost = "locality1.blob.core.windows.net"

// layoutSessionTransport serves CreateSession, GetLayout and blob reads, recording the
// authorization scheme and host of every read.
type layoutSessionTransport struct {
	mu sync.Mutex

	blobSize    int64
	layoutRange int64

	createSessionCalls int
	layoutCalls        int
	reads              []layoutSessionRead
	// ops records every request that is not CreateSession, so a test can assert what an
	// operation did and did not put on the wire.
	ops []layoutSessionOp
}

type layoutSessionRead struct {
	urlHost    string
	hostHeader string
	scheme     string
	start, end int64
}

// layoutSessionOp is one request as it reached the transport.
type layoutSessionOp struct {
	method     string
	comp       string
	scheme     string
	urlHost    string
	hostHeader string
}

func (tr *layoutSessionTransport) Do(req *http.Request) (*http.Response, error) {
	query := req.URL.Query()
	scheme, _, _ := strings.Cut(req.Header.Get("Authorization"), " ")

	if req.Method == http.MethodPost && query.Get("comp") == "session" {
		tr.mu.Lock()
		tr.createSessionCalls++
		tr.mu.Unlock()
		body := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<CreateSessionResult>
	<AuthenticationType>HMAC</AuthenticationType>
	<Id>layout-session</Id>
	<Credentials>
		<SessionKey>%s</SessionKey>
		<SessionToken>layout-session-token</SessionToken>
	</Credentials>
	<Expiration>%s</Expiration>
</CreateSessionResult>`, sessionTestKey, time.Now().Add(time.Hour).Format(time.RFC1123))
		return newSessionTestResponse(req, http.StatusCreated, http.Header{}, []byte(body)), nil
	}

	tr.mu.Lock()
	tr.ops = append(tr.ops, layoutSessionOp{
		method: req.Method, comp: query.Get("comp"), scheme: scheme,
		urlHost: req.URL.Host, hostHeader: req.Host,
	})
	tr.mu.Unlock()

	if query.Get("comp") == "layout" {
		tr.mu.Lock()
		tr.layoutCalls++
		tr.mu.Unlock()
		return tr.layoutResponse(req), nil
	}

	if req.Method == http.MethodGet {
		start, end := parseXMSRange(req)
		if end >= tr.blobSize {
			end = tr.blobSize - 1
		}
		tr.mu.Lock()
		tr.reads = append(tr.reads, layoutSessionRead{
			urlHost: req.URL.Host, hostHeader: req.Host, scheme: scheme, start: start, end: end,
		})
		tr.mu.Unlock()

		header := http.Header{}
		header.Set("Content-Length", fmt.Sprintf("%d", end-start+1))
		header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, tr.blobSize))
		header.Set("ETag", `"layout-session-etag"`)
		return newSessionTestResponse(req, http.StatusPartialContent, header, make([]byte, end-start+1)), nil
	}

	header := http.Header{}
	header.Set("ETag", `"layout-session-etag"`)
	switch req.Method {
	case http.MethodHead:
		// Get Blob Properties
		header.Set("Content-Length", fmt.Sprintf("%d", tr.blobSize))
		return newSessionTestResponse(req, http.StatusOK, header, nil), nil
	case http.MethodPut:
		// Put Blob, Put Block, Put Block List
		return newSessionTestResponse(req, http.StatusCreated, header, nil), nil
	}
	return newSessionTestResponse(req, http.StatusOK, header, nil), nil
}

// layoutResponse describes the whole blob as ranges served by localityHost.
func (tr *layoutSessionTransport) layoutResponse(req *http.Request) *http.Response {
	type endpoint struct {
		Index int32  `xml:"Index,attr"`
		Value string `xml:"Value,attr"`
	}
	type rng struct {
		Start         int64 `xml:"Start,attr"`
		End           int64 `xml:"End,attr"`
		EndpointIndex int32 `xml:"EndpointIndex,attr"`
	}
	type endpoints struct {
		Endpoint []endpoint `xml:"Endpoint"`
	}
	type ranges struct {
		Range []rng `xml:"Range"`
	}
	type blobLayout struct {
		XMLName   xml.Name  `xml:"BlobLayout"`
		Endpoints endpoints `xml:"Endpoints"`
		Ranges    ranges    `xml:"Ranges"`
	}

	l := blobLayout{Endpoints: endpoints{Endpoint: []endpoint{{Index: 0, Value: "https://" + localityHost}}}}
	for start := int64(0); start < tr.blobSize; start += tr.layoutRange {
		end := start + tr.layoutRange - 1
		if end >= tr.blobSize {
			end = tr.blobSize - 1
		}
		l.Ranges.Range = append(l.Ranges.Range, rng{Start: start, End: end, EndpointIndex: 0})
	}

	data, _ := xml.Marshal(l)
	header := http.Header{}
	header.Set("x-ms-blob-content-length", fmt.Sprintf("%d", tr.blobSize))
	header.Set("ETag", `"layout-session-etag"`)
	return &http.Response{
		Request:    req,
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       io.NopCloser(bytes.NewReader(data)),
	}
}

func parseXMSRange(req *http.Request) (int64, int64) {
	// the generated clients write this header as a raw lowercase map key, which
	// http.Header.Get cannot see, so read it the way production code does
	raw := shared.HeaderValue(req.Header, "x-ms-range")
	var start, end int64
	if _, err := fmt.Sscanf(raw, "bytes=%d-%d", &start, &end); err != nil {
		return 0, -1
	}
	return start, end
}

func (tr *layoutSessionTransport) opsSnapshot() (ops []layoutSessionOp, layoutCalls int) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	return append([]layoutSessionOp(nil), tr.ops...), tr.layoutCalls
}

func (tr *layoutSessionTransport) snapshot() (reads []layoutSessionRead, layoutCalls, createSessions int) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	return append([]layoutSessionRead(nil), tr.reads...), tr.layoutCalls, tr.createSessionCalls
}

func newLayoutSessionClient(t *testing.T, tr *layoutSessionTransport) (*blob.Client, *blockblob.Client) {
	t.Helper()
	opts := &service.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: tr,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		},
		Session: azblob.SessionOptions{
			Mode:        azblob.SessionModeEnabled,
			AccountName: sessionTestAccount,
		},
	}
	svcClient, err := service.NewClient(sessionTestServiceURL, sessionTestTokenCredential{}, opts)
	require.NoError(t, err)
	contClient := svcClient.NewContainerClient(sessionTestContainer)
	return contClient.NewBlobClient(sessionTestBlob), contClient.NewBlockBlobClient(sessionTestBlob)
}

// A layout-aware download with sessions enabled must authenticate every read with the session,
// including the chunks the layout reroutes to a different endpoint.
func TestLayoutAwareDownloadUsesSessionOnEveryChunk(t *testing.T) {
	tr := &layoutSessionTransport{blobSize: 300, layoutRange: 100}
	blobClient, _ := newLayoutSessionClient(t, tr)

	n, err := blobClient.DownloadBuffer(context.Background(), make([]byte, 300), &blob.DownloadBufferOptions{
		LayoutAwareRouting: blob.LayoutAwareRoutingEnabled,
		BlockSize:          100,
	})
	require.NoError(t, err)
	require.Equal(t, int64(300), n)

	reads, layoutCalls, createSessions := tr.snapshot()
	require.Equal(t, 1, createSessions, "one container session covers the whole download")
	require.Equal(t, 1, layoutCalls)
	require.Len(t, reads, 3, "the initial read plus the two remaining chunks")

	for i, r := range reads {
		require.Equal(t, "Session", r.scheme,
			"read %d (bytes %d-%d) must be session-authenticated even after layout routing", i, r.start, r.end)
	}

	// the initial read precedes the layout, so it goes to the account endpoint
	require.Equal(t, sessionTestHost, reads[0].urlHost)
	for _, r := range reads[1:] {
		require.Equal(t, localityHost, r.urlHost, "remaining chunks are routed to the layout endpoint")
		require.Equal(t, sessionTestHost, r.hostHeader, "the account Host header must be preserved")
	}
}

// GetLayout itself is not a session-eligible operation, so it stays on the bearer token while the
// reads around it use the session.
func TestGetLayoutUsesBearerWhileReadsUseSession(t *testing.T) {
	tr := &layoutSessionTransport{blobSize: 300, layoutRange: 100}
	blobClient, _ := newLayoutSessionClient(t, tr)

	_, err := blobClient.DownloadBuffer(context.Background(), make([]byte, 300), &blob.DownloadBufferOptions{
		LayoutAwareRouting: blob.LayoutAwareRoutingEnabled,
		BlockSize:          100,
	})
	require.NoError(t, err)

	reads, layoutCalls, _ := tr.snapshot()
	require.Equal(t, 1, layoutCalls)
	require.NotEmpty(t, reads)
	for _, r := range reads {
		require.Equal(t, "Session", r.scheme)
	}
}

// With layout routing disabled the reads still use the session and stay on the account endpoint.
func TestLayoutDisabledDownloadStillUsesSession(t *testing.T) {
	tr := &layoutSessionTransport{blobSize: 300, layoutRange: 100}
	blobClient, _ := newLayoutSessionClient(t, tr)

	_, err := blobClient.DownloadBuffer(context.Background(), make([]byte, 300), &blob.DownloadBufferOptions{
		LayoutAwareRouting: blob.LayoutAwareRoutingDisabled,
		BlockSize:          100,
	})
	require.NoError(t, err)

	reads, layoutCalls, _ := tr.snapshot()
	require.Zero(t, layoutCalls)
	require.Len(t, reads, 3)
	for _, r := range reads {
		require.Equal(t, "Session", r.scheme)
		require.Equal(t, sessionTestHost, r.urlHost)
	}
}

// Data Locality is scoped to Get Blob and the managed downloads built on it. Gaining session
// authentication did not give the other four operations a layout: these tests fail if a layout
// lookup or a rerouted request ever appears on GetBlobProperties, PutBlob, PutBlock or
// PutBlockList.
//
// The client is configured exactly as the locality-routed download tests above are, so the
// layout machinery is available and would be exercised if the operation reached for it.
func TestNewSessionAPIsDoNotUseDataLocality(t *testing.T) {
	ctx := context.Background()
	payload := []byte("private drop payload")

	tests := []struct {
		name   string
		call   func(t *testing.T, blobClient *blob.Client, bbClient *blockblob.Client)
		method string
		comp   string
	}{
		{
			name: "GetBlobProperties",
			call: func(t *testing.T, blobClient *blob.Client, _ *blockblob.Client) {
				_, err := blobClient.GetProperties(ctx, nil)
				require.NoError(t, err)
			},
			method: http.MethodHead,
			comp:   "",
		},
		{
			name: "PutBlob",
			call: func(t *testing.T, _ *blob.Client, bbClient *blockblob.Client) {
				_, err := bbClient.Upload(ctx, streaming.NopCloser(bytes.NewReader(payload)), nil)
				require.NoError(t, err)
			},
			method: http.MethodPut,
			comp:   "",
		},
		{
			name: "PutBlock",
			call: func(t *testing.T, _ *blob.Client, bbClient *blockblob.Client) {
				_, err := bbClient.StageBlock(ctx, "YmxvY2sx", streaming.NopCloser(bytes.NewReader(payload)), nil)
				require.NoError(t, err)
			},
			method: http.MethodPut,
			comp:   "block",
		},
		{
			name: "PutBlockList",
			call: func(t *testing.T, _ *blob.Client, bbClient *blockblob.Client) {
				_, err := bbClient.CommitBlockList(ctx, []string{"YmxvY2sx"}, nil)
				require.NoError(t, err)
			},
			method: http.MethodPut,
			comp:   "blocklist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &layoutSessionTransport{blobSize: 300, layoutRange: 100}
			blobClient, bbClient := newLayoutSessionClient(t, tr)

			tt.call(t, blobClient, bbClient)

			ops, layoutCalls := tr.opsSnapshot()

			require.Zero(t, layoutCalls, "%s must not fetch a layout", tt.name)
			require.Len(t, ops, 1, "%s must put exactly one request on the wire", tt.name)

			op := ops[0]
			require.Equal(t, tt.method, op.method)
			require.Equal(t, tt.comp, op.comp)
			require.Equal(t, "Session", op.scheme, "%s must use session authentication", tt.name)
			require.Equal(t, sessionTestHost, op.urlHost,
				"%s must go to the account endpoint, not a layout endpoint", tt.name)
			require.Equal(t, sessionTestHost, op.hostHeader)

			// nothing was routed anywhere else
			for _, o := range ops {
				require.NotEqual(t, localityHost, o.urlHost,
					"%s must never be routed to a layout endpoint", tt.name)
			}
		})
	}
}

// The same four operations must stay off the layout when sessions are not in play either, so the
// absence of Data Locality is a property of the operation and not of the auth scheme.
func TestNewSessionAPIsDoNotUseDataLocalityWithoutSession(t *testing.T) {
	ctx := context.Background()

	tr := &layoutSessionTransport{blobSize: 300, layoutRange: 100}
	opts := &blockblob.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: tr,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		},
	}
	bbClient, err := blockblob.NewClientWithNoCredential(
		sessionTestServiceURL+sessionTestContainer+"/"+sessionTestBlob, opts)
	require.NoError(t, err)

	_, err = bbClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte("no session"))), nil)
	require.NoError(t, err)
	_, err = bbClient.StageBlock(ctx, "YmxvY2sx", streaming.NopCloser(bytes.NewReader([]byte("no session"))), nil)
	require.NoError(t, err)
	_, err = bbClient.CommitBlockList(ctx, []string{"YmxvY2sx"}, nil)
	require.NoError(t, err)
	_, err = bbClient.BlobClient().GetProperties(ctx, nil)
	require.NoError(t, err)

	ops, layoutCalls := tr.opsSnapshot()
	require.Zero(t, layoutCalls, "no layout lookup belongs on these operations")
	require.Len(t, ops, 4)
	for _, o := range ops {
		require.Equal(t, sessionTestHost, o.urlHost, "%s %s must use the account endpoint", o.method, o.comp)
		require.Empty(t, o.scheme, "an anonymous client sends no Authorization header")
	}
}

// flakySessionBody yields a couple of bytes and then reports an unexpected EOF, which is what a
// dropped connection looks like mid-body and what makes the retry reader re-issue the read.
type flakySessionBody struct{ remaining int }

func (b *flakySessionBody) Read(p []byte) (int, error) {
	if b.remaining > 0 && len(p) > 0 {
		b.remaining--
		p[0] = 'a'
		return 1, nil
	}
	return 0, io.ErrUnexpectedEOF
}

func (b *flakySessionBody) Close() error { return nil }

// retrySessionTransport serves CreateSession and ranged reads, failing the first read mid-body.
type retrySessionTransport struct {
	mu    sync.Mutex
	reads []layoutSessionOp
	calls int
}

func (tr *retrySessionTransport) Do(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodPost && req.URL.Query().Get("comp") == "session" {
		body := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<CreateSessionResult>
	<AuthenticationType>HMAC</AuthenticationType>
	<Id>retry-session</Id>
	<Credentials>
		<SessionKey>%s</SessionKey>
		<SessionToken>retry-session-token</SessionToken>
	</Credentials>
	<Expiration>%s</Expiration>
</CreateSessionResult>`, sessionTestKey, time.Now().Add(time.Hour).Format(time.RFC1123))
		return newSessionTestResponse(req, http.StatusCreated, http.Header{}, []byte(body)), nil
	}

	scheme, _, _ := strings.Cut(req.Header.Get("Authorization"), " ")

	tr.mu.Lock()
	tr.calls++
	call := tr.calls
	tr.reads = append(tr.reads, layoutSessionOp{
		method: req.Method, comp: req.URL.Query().Get("comp"), scheme: scheme,
		urlHost: req.URL.Host, hostHeader: req.Host,
	})
	tr.mu.Unlock()

	start, end := parseXMSRange(req)
	if end < start {
		start, end = 0, 9
	}
	header := http.Header{}
	header.Set("Content-Length", fmt.Sprintf("%d", end-start+1))
	header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, 10))
	header.Set("ETag", `"retry-session-etag"`)

	var body io.ReadCloser
	if call == 1 {
		body = &flakySessionBody{remaining: 2}
	} else {
		body = io.NopCloser(bytes.NewReader(make([]byte, end-start+1)))
	}
	return &http.Response{
		Request: req, StatusCode: http.StatusPartialContent, Header: header, Body: body,
	}, nil
}

func (tr *retrySessionTransport) snapshotReads() []layoutSessionOp {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	return append([]layoutSessionOp(nil), tr.reads...)
}

// When a locality-routed read is retried, it must stay on the layout endpoint and it must be
// signed again: the session policy sits inside the retry loop, so every attempt carries its own
// signature rather than reusing the first one.
func TestLocalityRoutedRetryReappliesSession(t *testing.T) {
	tr := &retrySessionTransport{}

	opts := &service.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: tr,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		},
		Session: azblob.SessionOptions{
			Mode:        azblob.SessionModeEnabled,
			AccountName: sessionTestAccount,
		},
	}
	svcClient, err := service.NewClient(sessionTestServiceURL, sessionTestTokenCredential{}, opts)
	require.NoError(t, err)
	blobClient := svcClient.NewContainerClient(sessionTestContainer).NewBlobClient(sessionTestBlob)

	ctx := context.Background()
	dr, err := blobClient.DownloadStream(ctx, &blob.DownloadStreamOptions{
		Range:          blob.HTTPRange{Offset: 0, Count: 10},
		LayoutEndpoint: "https://" + localityHost,
	})
	require.NoError(t, err)

	body := dr.NewRetryReader(ctx, &blob.RetryReaderOptions{MaxRetries: 2})
	_, err = io.ReadAll(body)
	require.NoError(t, err)
	require.NoError(t, body.Close())

	reads := tr.snapshotReads()
	require.GreaterOrEqual(t, len(reads), 2, "the failed read must have been retried")
	for i, r := range reads {
		require.Equal(t, localityHost, r.urlHost, "read %d must stay on the layout endpoint", i)
		require.Equal(t, sessionTestHost, r.hostHeader, "read %d must keep the account Host header", i)
		require.Equal(t, "Session", r.scheme, "read %d must be signed again on retry", i)
	}
}
