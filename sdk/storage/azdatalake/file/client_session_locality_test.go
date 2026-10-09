// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
	"github.com/stretchr/testify/require"
)

// These tests run DataLake clients against a fake account that serves both the DFS and the blob
// endpoint, and assert on what reached the wire: which endpoint each request went to and how it
// was authenticated. They port the behavior covered by the .NET STG105 DataLake
// SessionAuthenticationTests and DataLakeFileClientDataLocalityTests.

const (
	dlAccount    = "dlaccount"
	dlDFSHost    = dlAccount + ".dfs.core.windows.net"
	dlBlobHost   = dlAccount + ".blob.core.windows.net"
	dlDFSService = "https://" + dlDFSHost + "/"
	dlLayoutHost = "locality0.blob.core.windows.net"
	dlFileSize   = 300
	dlSessionKey = "ZmFrZWtleQ==" // base64 "fakekey"
	dlETag       = `"dl-etag"`
)

type dlTokenCredential struct{}

func (dlTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "fake-bearer", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

// dlRequest is one request as it reached the wire, other than Create Session.
type dlRequest struct {
	method     string
	urlHost    string
	hostHeader string
	path       string
	comp       string
	scheme     string
	token      string
	ifMatch    string
	rangeHdr   string
}

// dlAccountFake serves the DFS and blob endpoints of one account.
type dlAccountFake struct {
	mu sync.Mutex
	// createSessions counts Create Session requests per file system (container), and records the
	// host they were sent to.
	createSessions     map[string]int
	createSessionHosts []string
	requests           []dlRequest

	// downloadHint, when true, adds x-ms-download-hint: layout to Get Blob responses.
	downloadHint bool
	// layoutPages are the Get Layout pages, served in order by marker.
	layoutPages []string
	// rejectSession, when set, returns 401 for a session-authenticated request it selects.
	rejectSession func(r dlRequest) bool
	// failFirstRead makes the first Get Blob body fail part way through.
	failFirstRead bool
	reads         int
}

func newDLAccount() *dlAccountFake {
	return &dlAccountFake{createSessions: map[string]int{}}
}

func (a *dlAccountFake) Do(req *http.Request) (*http.Response, error) {
	q := req.URL.Query()
	container, _, _ := strings.Cut(strings.TrimPrefix(req.URL.Path, "/"), "/")

	if req.Method == http.MethodPost && q.Get("comp") == "session" {
		a.mu.Lock()
		a.createSessions[container]++
		n := a.createSessions[container]
		a.createSessionHosts = append(a.createSessionHosts, req.URL.Host)
		a.mu.Unlock()
		body := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?><CreateSessionResult><Credentials><SessionKey>%s</SessionKey>`+
			`<SessionToken>%s-%d</SessionToken></Credentials><Expiration>%s</Expiration></CreateSessionResult>`,
			dlSessionKey, container, n, time.Now().Add(time.Hour).UTC().Format(time.RFC1123))
		return dlResponse(req, http.StatusCreated, http.Header{}, []byte(body)), nil
	}

	r := dlRequest{
		method:     req.Method,
		urlHost:    req.URL.Host,
		hostHeader: req.Host,
		path:       req.URL.Path,
		comp:       q.Get("comp"),
		ifMatch:    req.Header.Get("If-Match"),
		rangeHdr:   firstNonEmpty(req.Header.Get("x-ms-range"), headerRaw(req.Header, "x-ms-range")),
	}
	auth := req.Header.Get("Authorization")
	r.scheme, _, _ = strings.Cut(auth, " ")
	if r.scheme == "Session" {
		_, cred, _ := strings.Cut(auth, " ")
		r.token, _, _ = strings.Cut(cred, ":")
	}
	a.mu.Lock()
	a.requests = append(a.requests, r)
	reject := a.rejectSession
	a.mu.Unlock()

	if reject != nil && r.scheme == "Session" && reject(r) {
		return dlResponse(req, http.StatusUnauthorized, http.Header{}, []byte("<Error><Code>InvalidAuthenticationInfo</Code></Error>")), nil
	}

	h := http.Header{}
	h.Set("ETag", dlETag)
	h.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	switch {
	case req.Method == http.MethodGet && r.comp == "layout":
		return a.layoutResponse(req, q.Get("marker")), nil
	case req.Method == http.MethodGet && r.comp == "":
		return a.readResponse(req, r), nil
	case req.Method == http.MethodHead:
		h.Set("Content-Length", strconv.Itoa(dlFileSize))
		h.Set("x-ms-resource-type", "file")
		return dlResponse(req, http.StatusOK, h, nil), nil
	case req.Method == http.MethodPut:
		return dlResponse(req, http.StatusCreated, h, nil), nil
	case req.Method == http.MethodPatch && q.Get("action") == "flush":
		return dlResponse(req, http.StatusOK, h, nil), nil
	case req.Method == http.MethodPatch:
		return dlResponse(req, http.StatusAccepted, h, nil), nil
	case req.Method == http.MethodDelete:
		return dlResponse(req, http.StatusOK, h, nil), nil
	}
	return dlResponse(req, http.StatusOK, h, nil), nil
}

func (a *dlAccountFake) readResponse(req *http.Request, r dlRequest) *http.Response {
	start, end := int64(0), int64(dlFileSize-1)
	if r.rangeHdr != "" {
		_, _ = fmt.Sscanf(r.rangeHdr, "bytes=%d-%d", &start, &end)
		if end >= dlFileSize {
			end = dlFileSize - 1
		}
	}
	data := bytes.Repeat([]byte{'d'}, int(end-start+1))
	h := http.Header{}
	h.Set("ETag", dlETag)
	h.Set("Content-Length", strconv.Itoa(len(data)))
	h.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, dlFileSize))
	if a.downloadHint {
		h.Set("x-ms-download-hint", "layout")
	}
	a.mu.Lock()
	a.reads++
	failThis := a.failFirstRead && a.reads == 1
	a.mu.Unlock()
	resp := dlResponse(req, http.StatusPartialContent, h, data)
	if failThis {
		resp.Body = &failingBody{data: data[:len(data)/2]}
	}
	return resp
}

func (a *dlAccountFake) layoutResponse(req *http.Request, marker string) *http.Response {
	i := 0
	if marker != "" {
		i, _ = strconv.Atoi(strings.TrimPrefix(marker, "m"))
	}
	h := http.Header{}
	h.Set("ETag", dlETag)
	h.Set("x-ms-blob-content-length", strconv.Itoa(dlFileSize))
	h.Set("x-ms-blob-content-type", "application/octet-stream")
	if i >= len(a.layoutPages) {
		return dlResponse(req, http.StatusOK, h, []byte(`<BlobLayout></BlobLayout>`))
	}
	return dlResponse(req, http.StatusOK, h, []byte(a.layoutPages[i]))
}

// layoutPage is a Get Layout page with one range per endpoint host.
func layoutPage(nextMarker string, ranges ...[3]int64) string {
	var b strings.Builder
	b.WriteString(`<BlobLayout><Endpoints><Endpoint Index="0" Value="https://` + dlLayoutHost + `"/></Endpoints><Ranges>`)
	for _, r := range ranges {
		fmt.Fprintf(&b, `<Range Start="%d" End="%d" EndpointIndex="%d"/>`, r[0], r[1], r[2])
	}
	b.WriteString(`</Ranges>`)
	if nextMarker != "" {
		b.WriteString(`<NextMarker>` + nextMarker + `</NextMarker>`)
	}
	b.WriteString(`</BlobLayout>`)
	return b.String()
}

func (a *dlAccountFake) snapshot() ([]dlRequest, map[string]int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	counts := map[string]int{}
	for k, v := range a.createSessions {
		counts[k] = v
	}
	return append([]dlRequest(nil), a.requests...), counts
}

func (a *dlAccountFake) last() dlRequest {
	reqs, _ := a.snapshot()
	return reqs[len(reqs)-1]
}

type failingBody struct {
	data []byte
	read bool
}

func (b *failingBody) Read(p []byte) (int, error) {
	if b.read {
		return 0, errors.New("connection reset")
	}
	b.read = true
	return copy(p, b.data), nil
}

func (b *failingBody) Close() error { return nil }

func dlResponse(req *http.Request, status int, h http.Header, body []byte) *http.Response {
	if body == nil {
		body = []byte{}
	}
	return &http.Response{Request: req, StatusCode: status, Status: http.StatusText(status), Header: h, Body: io.NopCloser(bytes.NewReader(body))}
}

func headerRaw(h http.Header, key string) string {
	if v := h[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func dlClientOptions(a *dlAccountFake) azcore.ClientOptions {
	return azcore.ClientOptions{Transport: a, Retry: policy.RetryOptions{MaxRetries: -1}}
}

func newDLFileClient(t *testing.T, a *dlAccountFake, path string, session azdatalake.SessionOptions) *file.Client {
	t.Helper()
	c, err := file.NewClient(dlDFSService+path, dlTokenCredential{}, &file.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: a, Retry: policy.RetryOptions{MaxRetries: -1}},
		Session:       session,
	})
	require.NoError(t, err)
	return c
}

func readAll(t *testing.T, c *file.Client, opts *file.DownloadStreamOptions) {
	t.Helper()
	resp, err := c.DownloadStream(context.Background(), opts)
	require.NoError(t, err)
	_, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
}

var enabled = azdatalake.SessionOptions{Mode: azdatalake.SessionModeEnabled}

// ---- Session ---------------------------------------------------------------------------------

func TestDataLakeFileReadUsesSession(t *testing.T) {
	a := newDLAccount()
	c := newDLFileClient(t, a, "fs/dir/file.txt", enabled)

	readAll(t, c, nil)

	r := a.last()
	require.Equal(t, http.MethodGet, r.method)
	require.Equal(t, dlBlobHost, r.urlHost, "file reads go to the blob endpoint")
	require.Equal(t, "Session", r.scheme)
	_, sessions := a.snapshot()
	require.Equal(t, map[string]int{"fs": 1}, sessions, "one session for the file system")
	require.Equal(t, []string{dlBlobHost}, a.createSessionHosts, "Create Session is sent to the blob endpoint")
}

func TestDataLakeSessionDefaultAndDisabledUseBearer(t *testing.T) {
	for _, mode := range []azdatalake.SessionMode{azdatalake.SessionModeAuto, azdatalake.SessionModeDisabled} {
		t.Run(string(mode), func(t *testing.T) {
			a := newDLAccount()
			c := newDLFileClient(t, a, "fs/file.txt", azdatalake.SessionOptions{Mode: mode})
			readAll(t, c, nil)
			require.Equal(t, "Bearer", a.last().scheme)
			_, sessions := a.snapshot()
			require.Empty(t, sessions)
		})
	}
}

// Requests to the DFS endpoint (create, append, flush, delete) and every blob request other than
// a file read stay on the bearer token, even with sessions enabled.
func TestDataLakeOnlyFileReadsUseSession(t *testing.T) {
	a := newDLAccount()
	c := newDLFileClient(t, a, "fs/file.txt", enabled)
	ctx := context.Background()

	_, err := c.Create(ctx, nil)
	require.NoError(t, err)
	_, err = c.AppendData(ctx, 0, streamingBody("hello"), nil)
	require.NoError(t, err)
	_, err = c.FlushData(ctx, 5, nil)
	require.NoError(t, err)
	_, _ = c.GetProperties(ctx, nil)
	readAll(t, c, nil)
	_, err = c.Delete(ctx, nil)
	require.NoError(t, err)

	reqs, sessions := a.snapshot()
	for _, r := range reqs {
		isRead := r.method == http.MethodGet && r.comp == "" && r.urlHost == dlBlobHost
		if isRead {
			require.Equal(t, "Session", r.scheme, "file read")
			continue
		}
		require.Equal(t, "Bearer", r.scheme, "%s %s on %s must use the bearer token", r.method, r.path, r.urlHost)
		if r.urlHost == dlDFSHost {
			require.NotEqual(t, "Session", r.scheme, "the DFS endpoint never uses sessions")
		}
	}
	require.Equal(t, map[string]int{"fs": 1}, sessions)
}

func TestDataLakeSessionIgnoredForSharedKeyAndSAS(t *testing.T) {
	a := newDLAccount()
	cred, err := azdatalake.NewSharedKeyCredential(dlAccount, dlSessionKey)
	require.NoError(t, err)
	opts := &file.ClientOptions{ClientOptions: dlClientOptions(a), Session: enabled}

	sk, err := file.NewClientWithSharedKeyCredential(dlDFSService+"fs/file.txt", cred, opts)
	require.NoError(t, err, "shared key clients ignore session options")
	readAll(t, sk, nil)
	require.Equal(t, "SharedKey", a.last().scheme)

	sas, err := file.NewClientWithNoCredential(dlDFSService+"fs/file.txt?sv=2026-12-06&sig=x", &file.ClientOptions{ClientOptions: dlClientOptions(a), Session: enabled})
	require.NoError(t, err, "SAS clients ignore session options")
	readAll(t, sas, nil)
	require.Equal(t, "", a.last().scheme)

	_, sessions := a.snapshot()
	require.Empty(t, sessions)
}

func TestDataLakeImplicitProviderSharedAcrossDerivedClients(t *testing.T) {
	a := newDLAccount()
	svc, err := service.NewClient(dlDFSService, dlTokenCredential{}, &service.ClientOptions{ClientOptions: dlClientOptions(a), Session: enabled})
	require.NoError(t, err)

	fs := svc.NewFileSystemClient("fs")
	readAll(t, fs.NewFileClient("a.txt"), nil)
	readAll(t, fs.NewFileClient("b.txt"), nil)
	dir := fs.NewDirectoryClient("dir")
	nested, err := dir.NewFileClient("c.txt")
	require.NoError(t, err)
	readAll(t, nested, nil)
	sub, err := dir.NewSubdirectoryClient("sub")
	require.NoError(t, err)
	deeper, err := sub.NewFileClient("d.txt")
	require.NoError(t, err)
	readAll(t, deeper, nil)

	readAll(t, svc.NewFileSystemClient("other").NewFileClient("e.txt"), nil)

	_, sessions := a.snapshot()
	require.Equal(t, map[string]int{"fs": 1, "other": 1}, sessions,
		"multiple files in one file system share its session; file systems have their own")
}

func TestDataLakeImplicitProvidersNotSharedAcrossIndependentClients(t *testing.T) {
	a := newDLAccount()
	readAll(t, newDLFileClient(t, a, "fs/a.txt", enabled), nil)
	readAll(t, newDLFileClient(t, a, "fs/b.txt", enabled), nil)
	_, sessions := a.snapshot()
	require.Equal(t, map[string]int{"fs": 2}, sessions)
}

// A provider shared by independently created clients keeps its sessions across them, including
// clients created after the earlier ones are gone.
func TestDataLakeSharedProviderSurvivesClientRecreation(t *testing.T) {
	a := newDLAccount()
	provider, err := azdatalake.NewContainerSessionProvider(dlDFSService, dlTokenCredential{}, &azcore.ClientOptions{Transport: a})
	require.NoError(t, err)
	shared := azdatalake.SessionOptions{Mode: azdatalake.SessionModeEnabled, Provider: provider}

	for i := range 3 {
		readAll(t, newDLFileClient(t, a, fmt.Sprintf("fs/file%d.txt", i), shared), nil)
	}
	fsClient, err := filesystem.NewClient(dlDFSService+"fs", dlTokenCredential{}, &filesystem.ClientOptions{ClientOptions: dlClientOptions(a), Session: shared})
	require.NoError(t, err)
	readAll(t, fsClient.NewFileClient("x.txt"), nil)

	_, sessions := a.snapshot()
	require.Equal(t, map[string]int{"fs": 1}, sessions)
	require.Equal(t, []string{dlBlobHost}, a.createSessionHosts, "a DFS service URL is converted to the blob endpoint")
}

func TestDataLakeConcurrentColdReadsCreateOneSession(t *testing.T) {
	a := newDLAccount()
	c := newDLFileClient(t, a, "fs/file.txt", enabled)
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			readAll(t, c, nil)
		}()
	}
	wg.Wait()
	_, sessions := a.snapshot()
	require.Equal(t, map[string]int{"fs": 1}, sessions)
}

func TestDataLakeUnauthorizedSessionFallsBackThenReacquires(t *testing.T) {
	a := newDLAccount()
	a.rejectSession = func(r dlRequest) bool { return r.token == "fs-1" }
	c := newDLFileClient(t, a, "fs/file.txt", enabled)

	readAll(t, c, nil)
	reqs, sessions := a.snapshot()
	require.Equal(t, "Session", reqs[len(reqs)-2].scheme)
	require.Equal(t, "Bearer", reqs[len(reqs)-1].scheme, "the rejected read is sent again with the bearer token")
	require.Equal(t, map[string]int{"fs": 1}, sessions)

	readAll(t, c, nil)
	r := a.last()
	require.Equal(t, "Session", r.scheme)
	require.Equal(t, "fs-2", r.token, "the next read uses a new session")
}

func TestDataLakeSessionAccountName(t *testing.T) {
	a := newDLAccount()
	opts := func(s azdatalake.SessionOptions) *file.ClientOptions {
		return &file.ClientOptions{ClientOptions: dlClientOptions(a), Session: s}
	}

	_, err := file.NewClient("https://files.contoso.com/fs/file.txt?sig=secret", dlTokenCredential{}, opts(enabled))
	require.Error(t, err, "explicitly enabled sessions need an account name")
	require.NotContains(t, err.Error(), "secret")

	_, err = file.NewClient("https://files.contoso.com/fs/file.txt", dlTokenCredential{}, opts(azdatalake.SessionOptions{}))
	require.NoError(t, err, "the default mode never fails construction")

	c, err := file.NewClient("https://files.contoso.com/fs/file.txt", dlTokenCredential{},
		opts(azdatalake.SessionOptions{Mode: azdatalake.SessionModeEnabled, AccountName: dlAccount}))
	require.NoError(t, err)
	readAll(t, c, nil)
	require.Equal(t, "Session", a.last().scheme)
}

// ---- Data locality ---------------------------------------------------------------------------

func TestDataLakeGetLayoutPager(t *testing.T) {
	a := newDLAccount()
	a.layoutPages = []string{
		layoutPage("m1", [3]int64{0, 149, 0}),
		layoutPage("", [3]int64{150, 299, 0}),
	}
	c := newDLFileClient(t, a, "fs/file.txt", azdatalake.SessionOptions{})

	pager := c.GetLayoutPager(&file.GetLayoutOptions{Range: &azdatalake.HTTPRange{Offset: 0, Count: dlFileSize}})
	var ranges []*file.LayoutRange
	pages := 0
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		pages++
		require.NotNil(t, page.Ranges)
		ranges = append(ranges, page.Ranges.Range...)
		require.Equal(t, int64(dlFileSize), *page.FileContentLength, "blob properties map to file properties")
		require.Equal(t, "application/octet-stream", *page.FileContentType)
		require.Equal(t, azcore.ETag(dlETag), *page.ETag)
		require.Equal(t, "https://"+dlLayoutHost, *page.Endpoints.Endpoint[0].Value)
	}
	require.Equal(t, 2, pages)
	require.Len(t, ranges, 2)

	reqs, _ := a.snapshot()
	var layoutReqs []dlRequest
	for _, r := range reqs {
		if r.comp == "layout" {
			layoutReqs = append(layoutReqs, r)
		}
	}
	require.Len(t, layoutReqs, 2)
	for _, r := range layoutReqs {
		require.Equal(t, dlBlobHost, r.urlHost, "Get Layout is forwarded to the blob endpoint")
		require.Equal(t, "bytes=0-299", r.rangeHdr)
	}
	require.Empty(t, layoutReqs[0].ifMatch)
	require.Equal(t, dlETag, layoutReqs[1].ifMatch, "later pages are locked to the first page's ETag")
}

func TestDataLakeDownloadStreamLayoutEndpoint(t *testing.T) {
	a := newDLAccount()
	c := newDLFileClient(t, a, "fs/file.txt", azdatalake.SessionOptions{})

	readAll(t, c, &file.DownloadStreamOptions{LayoutEndpoint: "https://" + dlLayoutHost})
	r := a.last()
	require.Equal(t, dlLayoutHost, r.urlHost, "the read is routed to the layout endpoint")
	require.Equal(t, dlBlobHost, r.hostHeader, "and still authenticates against the account")
}

func TestDataLakeDownloadStreamRetryStaysOnLayoutEndpoint(t *testing.T) {
	a := newDLAccount()
	a.failFirstRead = true
	c := newDLFileClient(t, a, "fs/file.txt", azdatalake.SessionOptions{})

	// a Range is required: DataLake's NewRetryReader dereferences it
	resp, err := c.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range:          &azdatalake.HTTPRange{Offset: 0, Count: dlFileSize},
		LayoutEndpoint: "https://" + dlLayoutHost,
	})
	require.NoError(t, err)
	body := resp.NewRetryReader(context.Background(), &file.RetryReaderOptions{MaxRetries: 2})
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	require.NoError(t, body.Close())
	require.Len(t, data, dlFileSize)

	reqs, _ := a.snapshot()
	require.GreaterOrEqual(t, len(reqs), 2, "the failed read was retried")
	for _, r := range reqs {
		require.Equal(t, dlLayoutHost, r.urlHost, "every attempt stays on the layout endpoint")
	}
}

func TestDataLakeManagedDownloadRoutingDefaultDisabled(t *testing.T) {
	a := newDLAccount()
	a.downloadHint = true
	a.layoutPages = []string{layoutPage("", [3]int64{0, 299, 0})}
	c := newDLFileClient(t, a, "fs/file.txt", azdatalake.SessionOptions{})

	_, err := c.DownloadBuffer(context.Background(), make([]byte, dlFileSize), &file.DownloadBufferOptions{ChunkSize: 100})
	require.NoError(t, err)
	tmp, err := os.CreateTemp(t.TempDir(), "dl")
	require.NoError(t, err)
	defer func() { require.NoError(t, tmp.Close()) }()
	_, err = c.DownloadFile(context.Background(), tmp, &file.DownloadFileOptions{ChunkSize: 100})
	require.NoError(t, err)

	reqs, _ := a.snapshot()
	for _, r := range reqs {
		require.NotEqual(t, "layout", r.comp, "the default doesn't fetch a layout")
		require.NotEqual(t, dlLayoutHost, r.urlHost, "the default doesn't route by layout")
	}
}

func TestDataLakeManagedDownloadRoutingEnabled(t *testing.T) {
	for _, tt := range []struct {
		name     string
		download func(t *testing.T, c *file.Client) error
	}{
		{"DownloadBuffer", func(t *testing.T, c *file.Client) error {
			_, err := c.DownloadBuffer(context.Background(), make([]byte, dlFileSize), &file.DownloadBufferOptions{
				ChunkSize: 100, Concurrency: 1, LayoutAwareRouting: file.LayoutAwareRoutingEnabled,
			})
			return err
		}},
		{"DownloadFile", func(t *testing.T, c *file.Client) error {
			tmp, err := os.CreateTemp(t.TempDir(), "dl")
			require.NoError(t, err)
			defer func() { require.NoError(t, tmp.Close()) }()
			_, err = c.DownloadFile(context.Background(), tmp, &file.DownloadFileOptions{
				ChunkSize: 100, Concurrency: 1, LayoutAwareRouting: file.LayoutAwareRoutingEnabled,
			})
			return err
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			a := newDLAccount()
			a.downloadHint = true
			a.layoutPages = []string{layoutPage("", [3]int64{100, 299, 0})}
			c := newDLFileClient(t, a, "fs/file.txt", enabled)

			require.NoError(t, tt.download(t, c))

			reqs, _ := a.snapshot()
			var layoutCalls, routed, initial int
			for _, r := range reqs {
				switch {
				case r.comp == "layout":
					layoutCalls++
					require.Equal(t, "bytes=100-299", r.rangeHdr, "only the remaining range is enumerated")
					require.Equal(t, dlETag, r.ifMatch, "the layout is pinned to the initial read's ETag")
					require.Equal(t, "Bearer", r.scheme, "Get Layout is not session-eligible")
				case r.method == http.MethodGet && r.urlHost == dlLayoutHost:
					routed++
					require.Equal(t, dlBlobHost, r.hostHeader)
					require.Equal(t, "Session", r.scheme, "routed chunks still use the session")
					require.Equal(t, dlETag, r.ifMatch)
				case r.method == http.MethodGet:
					initial++
				}
			}
			require.Equal(t, 1, layoutCalls)
			require.Equal(t, 2, routed, "the chunks after the initial read are routed")
			require.Equal(t, 1, initial)
			_, sessions := a.snapshot()
			require.Equal(t, map[string]int{"fs": 1}, sessions)
			require.Equal(t, []string{dlBlobHost}, a.createSessionHosts, "the session for a routed chunk is created at the account endpoint")
		})
	}
}

func streamingBody(s string) *readSeekNopCloser {
	return &readSeekNopCloser{Reader: bytes.NewReader([]byte(s))}
}

type readSeekNopCloser struct{ *bytes.Reader }

func (readSeekNopCloser) Close() error { return nil }
