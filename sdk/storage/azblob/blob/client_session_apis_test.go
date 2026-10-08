// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
)

// This file exercises session authentication end to end, against the operation it covers and the
// ones it doesn't: the requests the generated clients actually build are run through the real
// pipeline, so the assertions are about what reached the transport rather than about request
// eligibility in isolation.

const (
	sessionTestAccount    = "fakeaccount"
	sessionTestServiceURL = "https://fakeaccount.blob.core.windows.net/"
	sessionTestHost       = "fakeaccount.blob.core.windows.net"
	sessionTestContainer  = "testcontainer"
	sessionTestBlob       = "testblob"
	// a session key must be valid base64: the policy hands it to NewSharedKeyCredential
	sessionTestKey = "dGVzdC1zZXNzaW9uLWtleQ=="
)

// sessionTestTokenCredential returns a static token so the bearer token policy can run without
// contacting an identity provider.
type sessionTestTokenCredential struct{}

func (sessionTestTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "fake-bearer-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

// recordedOp is one non-CreateSession request as it reached the transport.
type recordedOp struct {
	method string
	comp   string
	// scheme is the first token of the Authorization header, i.e. "Session" or "Bearer".
	scheme string
	// token is the session token a Session-authenticated request was signed with.
	token string
	// container is the first path segment of the request URL.
	container string
	body      []byte
}

// sessionAPITransport serves CreateSession and the blob operations under test, recording the
// authorization scheme and body of every operation request.
type sessionAPITransport struct {
	mu sync.Mutex

	createSessionCalls int
	ops                []recordedOp

	// rejectSession, when set, returns 401 for a session-authenticated request it selects. It is
	// how the tests drive the invalidate-and-fall-back-to-bearer path.
	rejectSession func(op recordedOp) bool

	// rejected tracks the bodies of the 401 responses, so tests can check they were closed.
	rejected []*closeTrackingBody
}

func (tr *sessionAPITransport) Do(req *http.Request) (*http.Response, error) {
	query := req.URL.Query()

	if req.Method == http.MethodPost && query.Get("comp") == "session" {
		tr.mu.Lock()
		tr.createSessionCalls++
		n := tr.createSessionCalls
		tr.mu.Unlock()
		body := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<CreateSessionResult>
	<AuthenticationType>HMAC</AuthenticationType>
	<Id>test-session-id</Id>
	<Credentials>
		<SessionKey>%s</SessionKey>
		<SessionToken>session-token-%d</SessionToken>
	</Credentials>
	<Expiration>%s</Expiration>
</CreateSessionResult>`, sessionTestKey, n, time.Now().Add(time.Hour).Format(time.RFC1123))
		return newSessionTestResponse(req, http.StatusCreated, http.Header{}, []byte(body)), nil
	}

	scheme, _, _ := strings.Cut(req.Header.Get("Authorization"), " ")
	op := recordedOp{method: req.Method, comp: query.Get("comp"), scheme: scheme}
	if scheme == "Session" {
		_, credential, _ := strings.Cut(req.Header.Get("Authorization"), " ")
		op.token, _, _ = strings.Cut(credential, ":")
	}
	op.container, _, _ = strings.Cut(strings.TrimPrefix(req.URL.Path, "/"), "/")
	if req.Body != nil {
		op.body, _ = io.ReadAll(req.Body)
	}

	tr.mu.Lock()
	tr.ops = append(tr.ops, op)
	reject := tr.rejectSession
	tr.mu.Unlock()

	if reject != nil && scheme == "Session" && reject(op) {
		body := &closeTrackingBody{Reader: bytes.NewReader([]byte("<Error><Code>InvalidAuthenticationInfo</Code></Error>"))}
		tr.mu.Lock()
		tr.rejected = append(tr.rejected, body)
		tr.mu.Unlock()
		resp := newSessionTestResponse(req, http.StatusUnauthorized, http.Header{}, nil)
		resp.Body = body
		return resp, nil
	}

	return sessionTestSuccessResponse(req, op)
}

// sessionTestSuccessResponse returns the response shape the generated client expects for op.
func sessionTestSuccessResponse(req *http.Request, op recordedOp) (*http.Response, error) {
	header := http.Header{}
	header.Set("ETag", `"session-test-etag"`)

	switch {
	case op.method == http.MethodGet && op.comp == "":
		// Get Blob
		data := []byte("blob contents")
		header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
		return newSessionTestResponse(req, http.StatusOK, header, data), nil
	case op.method == http.MethodHead && op.comp == "":
		// Get Blob Properties
		header.Set("Content-Length", "13")
		return newSessionTestResponse(req, http.StatusOK, header, nil), nil
	case op.method == http.MethodPut && (op.comp == "" || op.comp == "block" || op.comp == "blocklist"):
		// Put Blob, Put Block, Put Block List
		return newSessionTestResponse(req, http.StatusCreated, header, nil), nil
	case op.method == http.MethodDelete:
		return newSessionTestResponse(req, http.StatusAccepted, header, nil), nil
	default:
		return newSessionTestResponse(req, http.StatusOK, header, nil), nil
	}
}

func newSessionTestResponse(req *http.Request, status int, header http.Header, body []byte) *http.Response {
	if body == nil {
		body = []byte{}
	}
	return &http.Response{
		Request:    req,
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     header,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
}

// counts returns the number of operation requests authenticated with each scheme.
func (tr *sessionAPITransport) counts() (sessionOps, bearerOps, createSessions int) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	for _, op := range tr.ops {
		switch op.scheme {
		case "Session":
			sessionOps++
		case "Bearer":
			bearerOps++
		}
	}
	return sessionOps, bearerOps, tr.createSessionCalls
}

func (tr *sessionAPITransport) recorded() []recordedOp {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	return append([]recordedOp(nil), tr.ops...)
}

// newSessionTestClients builds a session-enabled container client backed by tr.
func newSessionTestClients(t *testing.T, tr *sessionAPITransport) (*blob.Client, *blockblob.Client) {
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

// ---- Get Blob is the only operation a session authenticates ---------------------------------

func TestSessionAuthUsedForGetBlob(t *testing.T) {
	tr := &sessionAPITransport{}
	blobClient, _ := newSessionTestClients(t, tr)

	resp, err := blobClient.DownloadStream(context.Background(), nil)
	require.NoError(t, err)
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	require.Equal(t, "blob contents", string(data))

	ops := tr.recorded()
	require.Len(t, ops, 1)
	require.Equal(t, http.MethodGet, ops[0].method)
	require.Equal(t, "Session", ops[0].scheme)
}

// The private drop also signed Get Blob Properties, Put Blob, Put Block and Put Block List with a
// session. The public preview matches the other Azure Storage SDKs and signs only Get Blob, so
// these four use the bearer token and never create a session.
func TestPrivateDropSessionOperationsUseBearer(t *testing.T) {
	ctx := context.Background()
	putBlobBody := bytes.Repeat([]byte("put-blob-body-"), 64)
	putBlockBody := bytes.Repeat([]byte("put-block-body-"), 64)

	tests := []struct {
		name     string
		method   string
		comp     string
		wantBody []byte
		call     func(t *testing.T, blobClient *blob.Client, bbClient *blockblob.Client)
	}{
		{"GetBlobProperties", http.MethodHead, "", nil, func(t *testing.T, blobClient *blob.Client, _ *blockblob.Client) {
			_, err := blobClient.GetProperties(ctx, nil)
			require.NoError(t, err)
		}},
		{"PutBlob", http.MethodPut, "", putBlobBody, func(t *testing.T, _ *blob.Client, bbClient *blockblob.Client) {
			_, err := bbClient.Upload(ctx, streaming.NopCloser(bytes.NewReader(putBlobBody)), nil)
			require.NoError(t, err)
		}},
		{"PutBlock", http.MethodPut, "block", putBlockBody, func(t *testing.T, _ *blob.Client, bbClient *blockblob.Client) {
			_, err := bbClient.StageBlock(ctx, "YmxvY2sx", streaming.NopCloser(bytes.NewReader(putBlockBody)), nil)
			require.NoError(t, err)
		}},
		{"PutBlockList", http.MethodPut, "blocklist", nil, func(t *testing.T, _ *blob.Client, bbClient *blockblob.Client) {
			_, err := bbClient.CommitBlockList(ctx, []string{"YmxvY2sx"}, nil)
			require.NoError(t, err)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &sessionAPITransport{}
			blobClient, bbClient := newSessionTestClients(t, tr)

			tt.call(t, blobClient, bbClient)

			ops := tr.recorded()
			require.Len(t, ops, 1)
			require.Equal(t, tt.method, ops[0].method)
			require.Equal(t, tt.comp, ops[0].comp)
			require.Equal(t, "Bearer", ops[0].scheme, "%s must use the bearer token", tt.name)
			if tt.wantBody != nil {
				require.Equal(t, tt.wantBody, ops[0].body)
			}
			_, _, createSessions := tr.counts()
			require.Zero(t, createSessions, "%s must not create a session", tt.name)
		})
	}
}

func TestSessionAuthStagedUploadCreatesNoSession(t *testing.T) {
	tr := &sessionAPITransport{}
	_, bbClient := newSessionTestClients(t, tr)

	ctx := context.Background()
	_, err := bbClient.StageBlock(ctx, "YmxvY2sx", streaming.NopCloser(bytes.NewReader([]byte("one"))), nil)
	require.NoError(t, err)
	_, err = bbClient.StageBlock(ctx, "YmxvY2sy", streaming.NopCloser(bytes.NewReader([]byte("two"))), nil)
	require.NoError(t, err)
	_, err = bbClient.CommitBlockList(ctx, []string{"YmxvY2sx", "YmxvY2sy"}, nil)
	require.NoError(t, err)

	sessionOps, bearerOps, createSessions := tr.counts()
	require.Equal(t, 0, sessionOps)
	require.Equal(t, 3, bearerOps)
	require.Equal(t, 0, createSessions)
}

func TestSessionAuthNotUsedForOtherOperations(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		call func(t *testing.T, blobClient *blob.Client, bbClient *blockblob.Client)
	}{
		{"DeleteBlob", func(t *testing.T, blobClient *blob.Client, _ *blockblob.Client) {
			_, err := blobClient.Delete(ctx, nil)
			require.NoError(t, err)
		}},
		{"SetBlobMetadata", func(t *testing.T, blobClient *blob.Client, _ *blockblob.Client) {
			_, err := blobClient.SetMetadata(ctx, nil, nil)
			require.NoError(t, err)
		}},
		{"SetBlobTier", func(t *testing.T, blobClient *blob.Client, _ *blockblob.Client) {
			_, _ = blobClient.SetTier(ctx, blob.AccessTierCool, nil)
		}},
		{"GetTags", func(t *testing.T, blobClient *blob.Client, _ *blockblob.Client) {
			_, _ = blobClient.GetTags(ctx, nil)
		}},
		{"GetBlockList", func(t *testing.T, _ *blob.Client, bbClient *blockblob.Client) {
			_, _ = bbClient.GetBlockList(ctx, blockblob.BlockListTypeAll, nil)
		}},
		{"CopyFromURL", func(t *testing.T, blobClient *blob.Client, _ *blockblob.Client) {
			_, _ = blobClient.StartCopyFromURL(ctx, "https://other.blob.core.windows.net/c/b", nil)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &sessionAPITransport{}
			blobClient, bbClient := newSessionTestClients(t, tr)

			tt.call(t, blobClient, bbClient)

			sessionOps, bearerOps, createSessions := tr.counts()
			require.Equal(t, 0, sessionOps, "%s must not use session authentication", tt.name)
			require.Equal(t, 1, bearerOps)
			require.Equal(t, 0, createSessions, "an ineligible request must not create a session")
		})
	}
}

// ---- 401 invalidates the session and falls back to bearer ----------------------------------

func rejectFirstSessionRequest() func(recordedOp) bool {
	var once sync.Once
	return func(recordedOp) bool {
		reject := false
		once.Do(func() { reject = true })
		return reject
	}
}

// A rejected session sends this request once more with a bearer token, drains the rejected
// response, and leaves the next request to create a new session.
func TestSessionAuthUnauthorizedFallsBackThenReacquires(t *testing.T) {
	tr := &sessionAPITransport{rejectSession: rejectFirstSessionRequest()}
	blobClient, _ := newSessionTestClients(t, tr)
	ctx := context.Background()

	resp, err := blobClient.DownloadStream(ctx, nil)
	require.NoError(t, err, "the bearer fallback completes the request")
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	require.Equal(t, "blob contents", string(data))

	ops := tr.recorded()
	require.Len(t, ops, 2, "one rejected session attempt and one bearer attempt")
	require.Equal(t, "Session", ops[0].scheme)
	require.Equal(t, "Bearer", ops[1].scheme)
	_, _, createSessions := tr.counts()
	require.Equal(t, 1, createSessions, "the rejected request doesn't wait for a new session")
	require.True(t, tr.rejectedBodiesClosed(), "the rejected response is drained and closed")

	resp, err = blobClient.DownloadStream(ctx, nil)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	ops = tr.recorded()
	require.Len(t, ops, 3)
	require.Equal(t, "Session", ops[2].scheme, "the next request uses a replacement session")
	_, _, createSessions = tr.counts()
	require.Equal(t, 2, createSessions)
}

// Many downloads in flight with the same session are all rejected; together they replace the
// session once rather than once each.
func TestSessionAuthConcurrentUnauthorizedReplacesSessionOnce(t *testing.T) {
	tr := &sessionAPITransport{rejectSession: func(op recordedOp) bool { return op.token == "session-token-1" }}
	blobClient, _ := newSessionTestClients(t, tr)
	ctx := context.Background()

	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := blobClient.DownloadStream(ctx, nil)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())
		}()
	}
	wg.Wait()

	resp, err := blobClient.DownloadStream(ctx, nil)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	_, _, createSessions := tr.counts()
	require.Equal(t, 2, createSessions, "one replacement session")
	ops := tr.recorded()
	require.Equal(t, "Session", ops[len(ops)-1].scheme)
	require.Equal(t, "session-token-2", ops[len(ops)-1].token)
}

// A rejected session in one container leaves other containers' sessions alone.
func TestSessionAuthInvalidationIsPerContainer(t *testing.T) {
	tr := &sessionAPITransport{rejectSession: func(op recordedOp) bool { return op.container == "a" && op.token == "session-token-1" }}
	svcClient := newSessionTestServiceClient(t, tr, nil)
	download := func(container string) {
		downloadOnce(t, svcClient.NewContainerClient(container).NewBlobClient("blob"))
	}

	download("a") // session-token-1, rejected
	download("b") // session-token-2
	download("a") // session-token-3
	download("b") // still session-token-2

	_, _, createSessions := tr.counts()
	require.Equal(t, 3, createSessions)
	ops := tr.recorded()
	require.Equal(t, "session-token-2", ops[len(ops)-1].token, "container b keeps its session")
}

// ---- session providers ---------------------------------------------------------------------

func newSessionTestServiceClient(t *testing.T, tr *sessionAPITransport, provider azblob.SessionProvider) *service.Client {
	t.Helper()
	opts := &service.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: tr, Retry: policy.RetryOptions{MaxRetries: -1}},
		Session:       azblob.SessionOptions{Mode: azblob.SessionModeEnabled, AccountName: sessionTestAccount, Provider: provider},
	}
	svcClient, err := service.NewClient(sessionTestServiceURL, sessionTestTokenCredential{}, opts)
	require.NoError(t, err)
	return svcClient
}

func downloadOnce(t *testing.T, client *blob.Client) {
	t.Helper()
	resp, err := client.DownloadStream(context.Background(), nil)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
}

func TestSessionImplicitProviderIsSharedWithDerivedClients(t *testing.T) {
	tr := &sessionAPITransport{}
	svcClient := newSessionTestServiceClient(t, tr, nil)

	downloadOnce(t, svcClient.NewContainerClient(sessionTestContainer).NewBlobClient("one"))
	downloadOnce(t, svcClient.NewContainerClient(sessionTestContainer).NewBlobClient("two"))

	_, _, createSessions := tr.counts()
	require.Equal(t, 1, createSessions)
}

func TestSessionImplicitProvidersAreNotSharedAcrossIndependentClients(t *testing.T) {
	tr := &sessionAPITransport{}
	downloadOnce(t, newSessionTestServiceClient(t, tr, nil).NewContainerClient(sessionTestContainer).NewBlobClient("blob"))
	downloadOnce(t, newSessionTestServiceClient(t, tr, nil).NewContainerClient(sessionTestContainer).NewBlobClient("blob"))

	_, _, createSessions := tr.counts()
	require.Equal(t, 2, createSessions)
}

// A provider passed to independently created clients keeps one session per container for all of
// them, including clients created after earlier ones were discarded.
func TestSessionSharedProviderOutlivesClients(t *testing.T) {
	tr := &sessionAPITransport{}
	provider, err := azblob.NewContainerSessionProvider(sessionTestServiceURL, sessionTestTokenCredential{}, &azblob.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: tr, Retry: policy.RetryOptions{MaxRetries: -1}},
	})
	require.NoError(t, err)

	for range 3 {
		// a new client each time, as an application recreating its clients would
		svcClient := newSessionTestServiceClient(t, tr, provider)
		downloadOnce(t, svcClient.NewContainerClient(sessionTestContainer).NewBlobClient("blob"))
	}
	blobClient, err := blob.NewClient(sessionTestServiceURL+sessionTestContainer+"/other", sessionTestTokenCredential{}, &blob.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: tr, Retry: policy.RetryOptions{MaxRetries: -1}},
		Session:       azblob.SessionOptions{Mode: azblob.SessionModeEnabled, Provider: provider},
	})
	require.NoError(t, err)
	downloadOnce(t, blobClient)

	_, _, createSessions := tr.counts()
	require.Equal(t, 1, createSessions, "one session for the container across every client")

	// a different container gets its own session from the same provider
	downloadOnce(t, newSessionTestServiceClient(t, tr, provider).NewContainerClient("othercontainer").NewBlobClient("blob"))
	_, _, createSessions = tr.counts()
	require.Equal(t, 2, createSessions)
}

// closeTrackingBody records whether a response body was closed.
type closeTrackingBody struct {
	*bytes.Reader
	closed atomic.Bool
}

func (b *closeTrackingBody) Close() error {
	b.closed.Store(true)
	return nil
}

// rejectedBodiesClosed reports whether every 401 response handed to the pipeline was closed.
func (tr *sessionAPITransport) rejectedBodiesClosed() bool {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	for _, b := range tr.rejected {
		if !b.closed.Load() {
			return false
		}
	}
	return len(tr.rejected) > 0
}
