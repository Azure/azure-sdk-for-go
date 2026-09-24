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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
)

// Data Locality rewrites the request URL host to a layout endpoint while keeping the account Host
// header, and session authentication signs the request. This file proves the two coexist: the
// signature is built from the configured account name and the URL path, never the host, so
// rerouting a chunk cannot invalidate it.

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
}

type layoutSessionRead struct {
	urlHost    string
	hostHeader string
	scheme     string
	start, end int64
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

	return newSessionTestResponse(req, http.StatusOK, http.Header{}, nil), nil
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
	raw := req.Header.Get("x-ms-range")
	if raw == "" {
		if v := req.Header["x-ms-range"]; len(v) > 0 {
			raw = v[0]
		}
	}
	var start, end int64
	if _, err := fmt.Sscanf(raw, "bytes=%d-%d", &start, &end); err != nil {
		return 0, -1
	}
	return start, end
}

func (tr *layoutSessionTransport) snapshot() (reads []layoutSessionRead, layoutCalls, createSessions int) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	return append([]layoutSessionRead(nil), tr.reads...), tr.layoutCalls, tr.createSessionCalls
}

func newLayoutSessionClient(t *testing.T, tr *layoutSessionTransport) *blob.Client {
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
	return svcClient.NewContainerClient(sessionTestContainer).NewBlobClient(sessionTestBlob)
}

// A layout-aware download with sessions enabled must authenticate every read with the session,
// including the chunks the layout reroutes to a different endpoint.
func TestLayoutAwareDownloadUsesSessionOnEveryChunk(t *testing.T) {
	tr := &layoutSessionTransport{blobSize: 300, layoutRange: 100}
	blobClient := newLayoutSessionClient(t, tr)

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
	blobClient := newLayoutSessionClient(t, tr)

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
	blobClient := newLayoutSessionClient(t, tr)

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
