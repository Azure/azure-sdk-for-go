// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/stretchr/testify/require"
)

// wireTransport answers Create Session with a session and every other request with 200, and
// records the authorization scheme of each request that isn't Create Session.
type wireTransport struct {
	mu             sync.Mutex
	schemes        []string
	createSessions int
}

func (w *wireTransport) Do(req *http.Request) (*http.Response, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if req.URL.Query().Get("comp") == "session" {
		w.createSessions++
		body := `<?xml version="1.0" encoding="utf-8"?><CreateSessionResult><Credentials><SessionKey>dGVzdC1rZXk=</SessionKey>` +
			`<SessionToken>tok</SessionToken></Credentials><Expiration>` + time.Now().Add(time.Hour).UTC().Format(time.RFC1123) +
			`</Expiration></CreateSessionResult>`
		return &http.Response{StatusCode: http.StatusCreated, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
	}
	scheme, _, _ := strings.Cut(req.Header.Get("Authorization"), " ")
	w.schemes = append(w.schemes, scheme)
	return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody, Request: req}, nil
}

func (w *wireTransport) lastScheme() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.schemes[len(w.schemes)-1]
}

func newWireClient(t *testing.T, tr *wireTransport, cred azcore.TokenCredential, sharedKey *exported.SharedKeyCredential, mode exported.SessionMode, rawURL string) *azcore.Client {
	t.Helper()
	opts := &ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: tr, Retry: policy.RetryOptions{MaxRetries: -1}},
		Session:       exported.SessionOptions{Mode: mode},
	}
	c, err := GetAzClient(rawURL, cred, sharedKey, opts)
	require.NoError(t, err)
	return c
}

// sendThrough sends a request with the given method, query and headers through the client's
// pipeline and returns the authorization scheme that reached the wire.
func sendThrough(t *testing.T, c *azcore.Client, tr *wireTransport, method, rawURL string, header map[string]string) string {
	t.Helper()
	req, err := runtime.NewRequest(context.Background(), method, rawURL)
	require.NoError(t, err)
	for k, v := range header {
		req.Raw().Header[k] = []string{v}
	}
	resp, err := c.Pipeline().Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	return tr.lastScheme()
}

const wireBlobURL = "https://account.blob.core.windows.net/container/blob"

func TestSessionAuthOnTheWire(t *testing.T) {
	sharedKey, err := exported.NewSharedKeyCredential("account", "dGVzdC1rZXk=")
	require.NoError(t, err)

	for _, tt := range []struct {
		name       string
		cred       azcore.TokenCredential
		sharedKey  *exported.SharedKeyCredential
		url        string
		mode       exported.SessionMode
		wantScheme string
	}{
		{"TokenEnabled", fakeTokenCredential{}, nil, wireBlobURL, exported.SessionModeEnabled, "Session"},
		{"TokenAuto", fakeTokenCredential{}, nil, wireBlobURL, exported.SessionModeAuto, "Bearer"},
		{"TokenDisabled", fakeTokenCredential{}, nil, wireBlobURL, exported.SessionModeDisabled, "Bearer"},
		{"SharedKeyEnabled", nil, sharedKey, wireBlobURL, exported.SessionModeEnabled, "SharedKey"},
		{"SASEnabled", nil, nil, wireBlobURL + "?sv=2026-12-06&sig=x", exported.SessionModeEnabled, ""},
		{"AnonymousEnabled", nil, nil, wireBlobURL, exported.SessionModeEnabled, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tr := &wireTransport{}
			c := newWireClient(t, tr, tt.cred, tt.sharedKey, tt.mode, tt.url)
			require.Equal(t, tt.wantScheme, sendThrough(t, c, tr, http.MethodGet, tt.url, nil))
			if tt.wantScheme != "Session" {
				require.Zero(t, tr.createSessions, "only token credential clients with sessions enabled create sessions")
			}
		})
	}
}

// The private drop also signed Get Blob Properties, Put Blob, Put Block and Put Block List with a
// session. The public preview signs only Get Blob, so these use the bearer token.
func TestPrivateDropSessionOperationsUseBearer(t *testing.T) {
	tr := &wireTransport{}
	c := newWireClient(t, tr, fakeTokenCredential{}, nil, exported.SessionModeEnabled, wireBlobURL)

	require.Equal(t, "Session", sendThrough(t, c, tr, http.MethodGet, wireBlobURL, nil), "Get Blob uses a session")
	require.Equal(t, "Bearer", sendThrough(t, c, tr, http.MethodHead, wireBlobURL, nil), "Get Blob Properties")
	require.Equal(t, "Bearer", sendThrough(t, c, tr, http.MethodPut, wireBlobURL, map[string]string{"x-ms-blob-type": "BlockBlob"}), "Put Blob")
	require.Equal(t, "Bearer", sendThrough(t, c, tr, http.MethodPut, wireBlobURL+"?comp=block&blockid=YQ%3D%3D", nil), "Put Block")
	require.Equal(t, "Bearer", sendThrough(t, c, tr, http.MethodPut, wireBlobURL+"?comp=blocklist", nil), "Put Block List")
	require.Equal(t, "Bearer", sendThrough(t, c, tr, http.MethodGet, wireBlobURL+"?comp=layout", nil), "Get Layout")
	require.Equal(t, 1, tr.createSessions)
}

func TestNoCredentialHasNoAuthPolicy(t *testing.T) {
	tr := &wireTransport{}
	c := newWireClient(t, tr, nil, nil, exported.SessionModeAuto, wireBlobURL)
	require.Equal(t, "", sendThrough(t, c, tr, http.MethodGet, wireBlobURL, nil))
}
