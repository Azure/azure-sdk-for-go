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

// This file exercises the operations that session authentication covers end to end: the requests
// the generated clients actually build are run through the real pipeline, so the assertions are
// about what reached the transport rather than about IsRequestEligible in isolation.

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
	body   []byte
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
	if req.Body != nil {
		op.body, _ = io.ReadAll(req.Body)
	}

	tr.mu.Lock()
	tr.ops = append(tr.ops, op)
	reject := tr.rejectSession
	tr.mu.Unlock()

	if reject != nil && scheme == "Session" && reject(op) {
		return newSessionTestResponse(req, http.StatusUnauthorized, http.Header{}, nil), nil
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

// ---- the five supported operations use session authentication ------------------------------

func TestSessionAuthUsedForGetBlob(t *testing.T) {
	tr := &sessionAPITransport{}
	blobClient, _ := newSessionTestClients(t, tr)

	resp, err := blobClient.DownloadStream(context.Background(), nil)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	ops := tr.recorded()
	require.Len(t, ops, 1)
	require.Equal(t, http.MethodGet, ops[0].method)
	require.Equal(t, "Session", ops[0].scheme)
}

func TestSessionAuthUsedForGetBlobProperties(t *testing.T) {
	tr := &sessionAPITransport{}
	blobClient, _ := newSessionTestClients(t, tr)

	_, err := blobClient.GetProperties(context.Background(), nil)
	require.NoError(t, err)

	ops := tr.recorded()
	require.Len(t, ops, 1)
	require.Equal(t, http.MethodHead, ops[0].method)
	require.Equal(t, "Session", ops[0].scheme, "GetBlobProperties must use session authentication")
}

func TestSessionAuthUsedForPutBlob(t *testing.T) {
	tr := &sessionAPITransport{}
	_, bbClient := newSessionTestClients(t, tr)

	payload := []byte("put blob payload")
	_, err := bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(payload)), nil)
	require.NoError(t, err)

	ops := tr.recorded()
	require.Len(t, ops, 1)
	require.Equal(t, http.MethodPut, ops[0].method)
	require.Empty(t, ops[0].comp)
	require.Equal(t, "Session", ops[0].scheme, "PutBlob must use session authentication")
	require.Equal(t, payload, ops[0].body)
}

func TestSessionAuthUsedForPutBlock(t *testing.T) {
	tr := &sessionAPITransport{}
	_, bbClient := newSessionTestClients(t, tr)

	payload := []byte("put block payload")
	_, err := bbClient.StageBlock(context.Background(), "YmxvY2sx", streaming.NopCloser(bytes.NewReader(payload)), nil)
	require.NoError(t, err)

	ops := tr.recorded()
	require.Len(t, ops, 1)
	require.Equal(t, http.MethodPut, ops[0].method)
	require.Equal(t, "block", ops[0].comp)
	require.Equal(t, "Session", ops[0].scheme, "PutBlock must use session authentication")
	require.Equal(t, payload, ops[0].body)
}

func TestSessionAuthUsedForPutBlockList(t *testing.T) {
	tr := &sessionAPITransport{}
	_, bbClient := newSessionTestClients(t, tr)

	_, err := bbClient.CommitBlockList(context.Background(), []string{"YmxvY2sx"}, nil)
	require.NoError(t, err)

	ops := tr.recorded()
	require.Len(t, ops, 1)
	require.Equal(t, http.MethodPut, ops[0].method)
	require.Equal(t, "blocklist", ops[0].comp)
	require.Equal(t, "Session", ops[0].scheme, "PutBlockList must use session authentication")
	require.NotEmpty(t, ops[0].body, "the block list body must reach the service")
}

// A full staged upload exercises all three write operations against one cached session.
func TestSessionAuthStagedUploadUsesOneSession(t *testing.T) {
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
	require.Equal(t, 3, sessionOps)
	require.Equal(t, 0, bearerOps)
	require.Equal(t, 1, createSessions, "the container session is created once and reused")
}

// ---- unsupported operations stay on bearer authentication ----------------------------------

func TestSessionAuthNotUsedForUnsupportedOperations(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		call func(t *testing.T, blobClient *blob.Client, bbClient *blockblob.Client)
	}{
		{
			name: "DeleteBlob",
			call: func(t *testing.T, blobClient *blob.Client, _ *blockblob.Client) {
				_, err := blobClient.Delete(ctx, nil)
				require.NoError(t, err)
			},
		},
		{
			name: "SetBlobMetadata",
			call: func(t *testing.T, blobClient *blob.Client, _ *blockblob.Client) {
				_, err := blobClient.SetMetadata(ctx, nil, nil)
				require.NoError(t, err)
			},
		},
		{
			name: "SetBlobTier",
			call: func(t *testing.T, blobClient *blob.Client, _ *blockblob.Client) {
				_, _ = blobClient.SetTier(ctx, blob.AccessTierCool, nil)
			},
		},
		{
			name: "GetBlockList",
			call: func(t *testing.T, _ *blob.Client, bbClient *blockblob.Client) {
				_, _ = bbClient.GetBlockList(ctx, blockblob.BlockListTypeAll, nil)
			},
		},
		{
			name: "CopyFromURL",
			call: func(t *testing.T, blobClient *blob.Client, _ *blockblob.Client) {
				_, _ = blobClient.StartCopyFromURL(ctx, "https://other.blob.core.windows.net/c/b", nil)
			},
		},
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

// A rejected session must not cost the request its body: the policy rewinds before handing the
// request to the bearer token policy, so the service still receives the complete payload.
func TestSessionAuthPutBlobFallsBackToBearerWithIntactBody(t *testing.T) {
	payload := bytes.Repeat([]byte("put-blob-body-"), 512)

	// reject only the first session-authenticated attempt
	var once sync.Once
	tr := &sessionAPITransport{
		rejectSession: func(recordedOp) bool {
			reject := false
			once.Do(func() { reject = true })
			return reject
		},
	}

	_, bbClient := newSessionTestClients(t, tr)

	_, err := bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(payload)), nil)
	require.NoError(t, err, "the bearer fallback must complete the upload")

	ops := tr.recorded()
	require.Len(t, ops, 2, "one rejected session attempt and one bearer fallback")
	require.Equal(t, "Session", ops[0].scheme)
	require.Equal(t, payload, ops[0].body, "the session attempt sends the full body")
	require.Equal(t, "Bearer", ops[1].scheme, "the rejected request falls back to bearer authentication")
	require.Equal(t, payload, ops[1].body, "the body must be rewound intact for the bearer fallback")
}

func TestSessionAuthPutBlockFallsBackToBearerWithIntactBody(t *testing.T) {
	payload := bytes.Repeat([]byte("put-block-body-"), 512)

	var once sync.Once
	tr := &sessionAPITransport{
		rejectSession: func(op recordedOp) bool {
			reject := false
			once.Do(func() { reject = true })
			return reject
		},
	}

	_, bbClient := newSessionTestClients(t, tr)

	_, err := bbClient.StageBlock(context.Background(), "YmxvY2sx", streaming.NopCloser(bytes.NewReader(payload)), nil)
	require.NoError(t, err)

	ops := tr.recorded()
	require.Len(t, ops, 2)
	require.Equal(t, "Session", ops[0].scheme)
	require.Equal(t, "block", ops[0].comp)
	require.Equal(t, "Bearer", ops[1].scheme)
	require.Equal(t, payload, ops[1].body, "the staged block body must survive the fallback")
}

// After a 401 the cached session is discarded, so the next eligible request acquires a new one.
func TestSessionAuthUnauthorizedInvalidatesAndReacquires(t *testing.T) {
	var once sync.Once
	tr := &sessionAPITransport{
		rejectSession: func(op recordedOp) bool {
			reject := false
			once.Do(func() { reject = true })
			return reject
		},
	}

	_, bbClient := newSessionTestClients(t, tr)
	ctx := context.Background()

	// first upload: session rejected, falls back to bearer, session invalidated
	_, err := bbClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte("first"))), nil)
	require.NoError(t, err)

	// second upload: a new session is acquired and used
	_, err = bbClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte("second"))), nil)
	require.NoError(t, err)

	ops := tr.recorded()
	require.Len(t, ops, 3)
	require.Equal(t, "Session", ops[0].scheme)
	require.Equal(t, "Bearer", ops[1].scheme)
	require.Equal(t, "Session", ops[2].scheme, "the next eligible request uses the replacement session")

	_, _, createSessions := tr.counts()
	require.Equal(t, 2, createSessions, "the rejected session is discarded and a new one created")
}

// GetBlobProperties carries no body, but it must still fall back cleanly.
func TestSessionAuthGetBlobPropertiesFallsBackToBearer(t *testing.T) {
	var once sync.Once
	tr := &sessionAPITransport{
		rejectSession: func(op recordedOp) bool {
			reject := false
			once.Do(func() { reject = true })
			return reject
		},
	}

	blobClient, _ := newSessionTestClients(t, tr)

	_, err := blobClient.GetProperties(context.Background(), nil)
	require.NoError(t, err)

	ops := tr.recorded()
	require.Len(t, ops, 2)
	require.Equal(t, "Session", ops[0].scheme)
	require.Equal(t, "Bearer", ops[1].scheme)
}
