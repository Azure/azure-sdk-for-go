// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
	"github.com/stretchr/testify/require"
)

// sessionAuthTracker is a pipeline policy that tracks session-related and bearer token requests.
// It is installed as a per-retry policy, so it observes traffic on both the DFS pipeline and the
// inner blob pipeline.
type sessionAuthTracker struct {
	mu                 sync.Mutex
	createSessionCount int
	sessionAuthCount   int
	bearerAuthCount    int
}

func (p *sessionAuthTracker) Do(req *policy.Request) (*http.Response, error) {
	p.mu.Lock()

	// a CreateSession request is a POST with comp=session
	if req.Raw().Method == http.MethodPost && req.Raw().URL.Query().Get("comp") == "session" {
		p.createSessionCount++
	}

	authHeader := req.Raw().Header.Get("Authorization")
	switch {
	case strings.HasPrefix(authHeader, "Session "):
		p.sessionAuthCount++
	case strings.HasPrefix(authHeader, "Bearer "):
		p.bearerAuthCount++
	}
	p.mu.Unlock()

	return req.Next()
}

func (p *sessionAuthTracker) counts() (createSessionCount, sessionAuthCount, bearerAuthCount int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.createSessionCount, p.sessionAuthCount, p.bearerAuthCount
}

// newSessionServiceClient builds a datalake service client backed by a token credential with
// session-based authentication configured, optionally installing the supplied tracker.
func newSessionServiceClient(t *testing.T, accountName string, mode azdatalake.SessionMode, tracker *sessionAuthTracker) (*service.Client, error) {
	t.Helper()

	cred, err := testcommon.GetGenericTokenCredential()
	require.NoError(t, err)

	sessionOptions := &service.ClientOptions{
		Session: azdatalake.SessionOptions{
			Mode:        mode,
			AccountName: accountName,
		},
	}
	testcommon.SetClientOptions(t, &sessionOptions.ClientOptions)
	if tracker != nil {
		sessionOptions.PerRetryPolicies = append(sessionOptions.PerRetryPolicies, tracker)
	}

	return service.NewClient(fmt.Sprintf("https://%s.dfs.core.windows.net/", accountName), cred, sessionOptions)
}

func (s *RecordedTestSuite) TestFileDownloadWithSessionOptions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	// set up the filesystem and file with a shared key client
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, fsName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	fileName := testcommon.GenerateFileName(testName)
	uploadData := []byte("test data for session download")

	fClient := fsClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NoError(fClient.UploadBuffer(context.Background(), uploadData, nil))

	// read it back with session-based authentication
	sessionSvcClient, err := newSessionServiceClient(s.T(), accountName, azdatalake.SessionModeEnabled, nil)
	_require.NoError(err)

	sessionFileClient := sessionSvcClient.NewFileSystemClient(fsName).NewFileClient(fileName)
	resp, err := sessionFileClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	downloadedData, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.NoError(resp.Body.Close())
	_require.Equal(uploadData, downloadedData)
}

// AccountName is optional; when it is omitted it is derived from the client's URL.
func (s *RecordedTestSuite) TestFileDownloadWithSessionAccountNameDerivedFromURL() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, fsName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	fileName := testcommon.GenerateFileName(testName)
	uploadData := []byte("test data for a derived account name")

	fClient := fsClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NoError(fClient.UploadBuffer(context.Background(), uploadData, nil))

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	tracker := &sessionAuthTracker{}

	// AccountName is intentionally not set here
	sessionOptions := &service.ClientOptions{
		Session: azdatalake.SessionOptions{
			Mode: azdatalake.SessionModeEnabled,
		},
	}
	testcommon.SetClientOptions(s.T(), &sessionOptions.ClientOptions)
	sessionOptions.PerRetryPolicies = append(sessionOptions.PerRetryPolicies, tracker)

	sessionSvcClient, err := service.NewClient(fmt.Sprintf("https://%s.dfs.core.windows.net/", accountName), cred, sessionOptions)
	_require.NoError(err)

	sessionFileClient := sessionSvcClient.NewFileSystemClient(fsName).NewFileClient(fileName)
	resp, err := sessionFileClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	downloadedData, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.NoError(resp.Body.Close())
	_require.Equal(uploadData, downloadedData)

	createSessionCount, sessionAuthCount, _ := tracker.counts()
	_require.Equal(1, createSessionCount, "expected a session to be created")
	_require.Equal(1, sessionAuthCount, "expected the read to use session auth with the derived account name")
}

// Sessions are only used when explicitly enabled, so neither the disabled mode nor the zero-value
// default mode may create or use one.
func (s *RecordedTestSuite) TestFileDownloadWithSessionModeOff() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, fsName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	fileName := testcommon.GenerateFileName(testName)
	uploadData := []byte("test data for session mode off")

	fClient := fsClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NoError(fClient.UploadBuffer(context.Background(), uploadData, nil))

	tracker := &sessionAuthTracker{}

	// SessionModeDefault is the zero value, so it covers clients that never mention sessions
	for _, mode := range []azdatalake.SessionMode{azdatalake.SessionModeDisabled, azdatalake.SessionModeDefault} {
		sessionSvcClient, err := newSessionServiceClient(s.T(), accountName, mode, tracker)
		_require.NoError(err)

		sessionFileClient := sessionSvcClient.NewFileSystemClient(fsName).NewFileClient(fileName)
		resp, err := sessionFileClient.DownloadStream(context.Background(), nil)
		_require.NoError(err)

		downloadedData, err := io.ReadAll(resp.Body)
		_require.NoError(err)
		_ = resp.Body.Close()
		_require.Equal(uploadData, downloadedData)
	}

	createSessionCount, sessionAuthCount, _ := tracker.counts()
	_require.Equal(0, createSessionCount, "Expected no CreateSession calls when sessions are not enabled")
	_require.Equal(0, sessionAuthCount, "Expected no session-authenticated requests when sessions are not enabled")
}

// The session scope is the filesystem, resolved from the request URL, so a client created for the
// service endpoint acquires a session for whichever filesystem the request targets.
func (s *RecordedTestSuite) TestFileDownloadWithSessionFilesystemResolvedFromURL() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, fsName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	fileName := testcommon.GenerateFileName(testName)
	uploadData := []byte("test data for filesystem resolved session")

	fClient := fsClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NoError(fClient.UploadBuffer(context.Background(), uploadData, nil))

	tracker := &sessionAuthTracker{}
	sessionSvcClient, err := newSessionServiceClient(s.T(), accountName, azdatalake.SessionModeEnabled, tracker)
	_require.NoError(err)

	sessionFileClient := sessionSvcClient.NewFileSystemClient(fsName).NewFileClient(fileName)
	resp, err := sessionFileClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	downloadedData, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_ = resp.Body.Close()
	_require.Equal(uploadData, downloadedData)

	createSessionCount, sessionAuthCount, _ := tracker.counts()
	_require.Equal(1, createSessionCount, "Expected a session to be created for the filesystem in the request URL")
	_require.Equal(1, sessionAuthCount, "Expected the download to use session authentication")
}

func (s *UnrecordedTestSuite) TestFileDownloadWithSessionOptionsConcurrentDownloads() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, fsName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	const numFiles = 5
	uploadData := []byte("test data for concurrent session download")
	fileNames := make([]string, numFiles)

	for i := 0; i < numFiles; i++ {
		fileNames[i] = fmt.Sprintf("%s-file-%d", testcommon.GenerateFileName(testName), i)
		fClient := fsClient.NewFileClient(fileNames[i])
		_, err = fClient.Create(context.Background(), nil)
		_require.NoError(err)
		_require.NoError(fClient.UploadBuffer(context.Background(), uploadData, nil))
	}

	tracker := &sessionAuthTracker{}
	sessionSvcClient, err := newSessionServiceClient(s.T(), accountName, azdatalake.SessionModeEnabled, tracker)
	_require.NoError(err)

	sessionFSClient := sessionSvcClient.NewFileSystemClient(fsName)

	var wg sync.WaitGroup
	errChan := make(chan error, numFiles)

	for i := 0; i < numFiles; i++ {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()
			sessionFileClient := sessionFSClient.NewFileClient(fileName)
			resp, err := sessionFileClient.DownloadStream(context.Background(), nil)
			if err != nil {
				errChan <- err
				return
			}
			downloadedData, err := io.ReadAll(resp.Body)
			if err != nil {
				errChan <- err
				return
			}
			_ = resp.Body.Close()
			if !bytes.Equal(uploadData, downloadedData) {
				errChan <- fmt.Errorf("downloaded data mismatch for file %s", fileName)
			}
		}(fileNames[i])
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		_require.NoError(err)
	}

	createSessionCount, sessionAuthCount, _ := tracker.counts()
	_require.Equal(1, createSessionCount, "Expected exactly one CreateSession call due to caching")
	_require.Equal(numFiles, sessionAuthCount, "Expected all downloads to use session authentication")
}

func (s *UnrecordedTestSuite) TestFileDownloadWithSessionOptionsLargeFileDownloadBuffer() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, fsName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	// a 10 MB file downloads in multiple chunks
	fileName := testcommon.GenerateFileName(testName)
	const fileSize = 10 * 1024 * 1024
	uploadData := make([]byte, fileSize)
	for i := range uploadData {
		uploadData[i] = byte(i % 256)
	}

	fClient := fsClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NoError(fClient.UploadBuffer(context.Background(), uploadData, nil))

	tracker := &sessionAuthTracker{}
	sessionSvcClient, err := newSessionServiceClient(s.T(), accountName, azdatalake.SessionModeEnabled, tracker)
	_require.NoError(err)

	sessionFileClient := sessionSvcClient.NewFileSystemClient(fsName).NewFileClient(fileName)

	buffer := make([]byte, fileSize)
	downloaded, err := sessionFileClient.DownloadBuffer(context.Background(), buffer, &file.DownloadBufferOptions{
		ChunkSize:   4 * 1024 * 1024,
		Concurrency: 2,
	})
	_require.NoError(err)
	_require.Equal(int64(fileSize), downloaded)
	_require.Equal(uploadData, buffer)

	createSessionCount, sessionAuthCount, _ := tracker.counts()
	_require.Equal(1, createSessionCount, "Expected exactly one CreateSession call")
	_require.Greater(sessionAuthCount, 0, "Expected at least one request with session authentication")
}

func (s *UnrecordedTestSuite) TestFileDownloadWithSessionOptionsLargeFileDownloadFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, fsName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	fileName := testcommon.GenerateFileName(testName)
	const fileSize = 10 * 1024 * 1024
	uploadData := make([]byte, fileSize)
	for i := range uploadData {
		uploadData[i] = byte(i % 256)
	}

	fClient := fsClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NoError(fClient.UploadBuffer(context.Background(), uploadData, nil))

	tracker := &sessionAuthTracker{}
	sessionSvcClient, err := newSessionServiceClient(s.T(), accountName, azdatalake.SessionModeEnabled, tracker)
	_require.NoError(err)

	sessionFileClient := sessionSvcClient.NewFileSystemClient(fsName).NewFileClient(fileName)

	tmpFile, err := os.CreateTemp("", "session-download-test-*")
	_require.NoError(err)
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	downloaded, err := sessionFileClient.DownloadFile(context.Background(), tmpFile, &file.DownloadFileOptions{
		ChunkSize:   4 * 1024 * 1024,
		Concurrency: 2,
	})
	_require.NoError(err)
	_require.Equal(int64(fileSize), downloaded)

	_, err = tmpFile.Seek(0, io.SeekStart)
	_require.NoError(err)
	downloadedData, err := io.ReadAll(tmpFile)
	_require.NoError(err)
	_require.Equal(uploadData, downloadedData)

	createSessionCount, sessionAuthCount, _ := tracker.counts()
	_require.Equal(1, createSessionCount, "Expected exactly one CreateSession call")
	_require.Greater(sessionAuthCount, 0, "Expected at least one request with session authentication")
}

func (s *UnrecordedTestSuite) TestFileDownloadWithSessionOptionsMultipleFilesSingleSession() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, fsName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	fileData := []byte("test data for multiple file session download")
	fileNames := make([]string, 4)
	for i := range fileNames {
		fileNames[i] = fmt.Sprintf("%s-%d", testcommon.GenerateFileName(testName), i)
		fClient := fsClient.NewFileClient(fileNames[i])
		_, err = fClient.Create(context.Background(), nil)
		_require.NoError(err)
		_require.NoError(fClient.UploadBuffer(context.Background(), fileData, nil))
	}

	tracker := &sessionAuthTracker{}
	sessionSvcClient, err := newSessionServiceClient(s.T(), accountName, azdatalake.SessionModeEnabled, tracker)
	_require.NoError(err)

	sessionFSClient := sessionSvcClient.NewFileSystemClient(fsName)

	for _, fileName := range fileNames {
		sessionFileClient := sessionFSClient.NewFileClient(fileName)
		resp, err := sessionFileClient.DownloadStream(context.Background(), nil)
		_require.NoError(err)

		downloadedData, err := io.ReadAll(resp.Body)
		_require.NoError(err)
		_require.NoError(resp.Body.Close())
		_require.Equal(fileData, downloadedData)
	}

	createSessionCount, sessionAuthCount, _ := tracker.counts()
	_require.Equal(1, createSessionCount, "expected only one session to be created")
	_require.Equal(len(fileNames), sessionAuthCount, "expected each read to use session auth")
}

// A SessionProvider supplied through SessionOptions is shared by every client it is injected into.
// Sessions are filesystem scoped, so independently created clients that target the same filesystem
// reuse a single session, while a client targeting a different filesystem mints a new one.
func (s *RecordedTestSuite) TestFileSharedSessionProviderReusedAcrossClients() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	// two files in the first filesystem and one in the second
	fileData := []byte("test data for a shared session provider")
	firstFSName := testcommon.GenerateFileSystemName(testName) + "1"
	secondFSName := testcommon.GenerateFileSystemName(testName) + "2"
	fileNames := []string{
		testcommon.GenerateFileName(testName) + "-1",
		testcommon.GenerateFileName(testName) + "-2",
	}

	firstFSClient := testcommon.CreateNewFileSystem(context.Background(), _require, firstFSName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, firstFSClient)
	for _, fileName := range fileNames {
		fClient := firstFSClient.NewFileClient(fileName)
		_, err = fClient.Create(context.Background(), nil)
		_require.NoError(err)
		_require.NoError(fClient.UploadBuffer(context.Background(), fileData, nil))
	}

	secondFSClient := testcommon.CreateNewFileSystem(context.Background(), _require, secondFSName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, secondFSClient)
	fClient := secondFSClient.NewFileClient(fileNames[0])
	_, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NoError(fClient.UploadBuffer(context.Background(), fileData, nil))

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	tracker := &sessionAuthTracker{}

	// the tracker is installed on the provider's pipeline too so CreateSession calls are counted
	providerOptions := &azdatalake.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &providerOptions.ClientOptions)
	providerOptions.PerRetryPolicies = append(providerOptions.PerRetryPolicies, tracker)

	serviceURL := fmt.Sprintf("https://%s.dfs.core.windows.net/", accountName)
	provider, err := azdatalake.NewFilesystemSessionProvider(cred, serviceURL, providerOptions)
	_require.NoError(err)

	// every file gets its own client, but they all share the injected provider
	downloadWithSharedProvider := func(fsName, fileName string) {
		fileOptions := &file.ClientOptions{
			Session: azdatalake.SessionOptions{
				Mode:        azdatalake.SessionModeEnabled,
				AccountName: accountName,
				Provider:    provider,
			},
		}
		testcommon.SetClientOptions(s.T(), &fileOptions.ClientOptions)
		fileOptions.PerRetryPolicies = append(fileOptions.PerRetryPolicies, tracker)

		fileURL := fmt.Sprintf("%s%s/%s", serviceURL, fsName, fileName)
		sessionFileClient, err := file.NewClient(fileURL, cred, fileOptions)
		_require.NoError(err)

		resp, err := sessionFileClient.DownloadStream(context.Background(), nil)
		_require.NoError(err)

		downloadedData, err := io.ReadAll(resp.Body)
		_require.NoError(err)
		_require.NoError(resp.Body.Close())
		_require.Equal(fileData, downloadedData)
	}

	// the two clients in the first filesystem share a session
	for _, fileName := range fileNames {
		downloadWithSharedProvider(firstFSName, fileName)
	}

	createSessionCount, sessionAuthCount, _ := tracker.counts()
	_require.Equal(1, createSessionCount, "expected clients in the same filesystem to share one session")
	_require.Equal(len(fileNames), sessionAuthCount, "expected every client to use session auth")

	// a client in a different filesystem needs a session of its own
	downloadWithSharedProvider(secondFSName, fileNames[0])

	createSessionCount, sessionAuthCount, _ = tracker.counts()
	_require.Equal(2, createSessionCount, "expected a second session for the second filesystem")
	_require.Equal(len(fileNames)+1, sessionAuthCount, "expected every client to use session auth")
}

func (s *RecordedTestSuite) TestFileRandomRestCallsUseBearerExceptGetUsesSession() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, fsName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	fileName := testcommon.GenerateFileName(testName)
	fClient := fsClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NoError(fClient.UploadBuffer(context.Background(), []byte("hello world"), nil))

	tracker := &sessionAuthTracker{}
	sessionSvcClient, err := newSessionServiceClient(s.T(), accountName, azdatalake.SessionModeEnabled, tracker)
	_require.NoError(err)

	sessionFileClient := sessionSvcClient.NewFileSystemClient(fsName).NewFileClient(fileName)

	// a read is eligible for session auth
	resp, err := sessionFileClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)
	_, err = io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.NoError(resp.Body.Close())

	// everything else falls back to bearer auth
	_, _ = sessionFileClient.SetMetadata(context.Background(), map[string]*string{"a": to.Ptr("b")}, nil)
	_, _ = sessionFileClient.SetHTTPHeaders(context.Background(), file.HTTPHeaders{
		ContentType: to.Ptr("text/plain"),
	}, nil)
	_, _ = sessionFileClient.GetProperties(context.Background(), nil)

	createSessionCount, sessionAuthCount, bearerAuthCount := tracker.counts()

	_require.Equal(1, createSessionCount, "expected a session to be created")
	_require.Equal(1, sessionAuthCount, "expected the read to use session auth")
	// the CreateSession call plus the three non-read operations all use bearer auth
	_require.GreaterOrEqual(bearerAuthCount, 4, "expected non-read REST calls to use bearer auth")
}

// TestSessionProviderUsesBlobEndpoint asserts that sessions are minted against the blob endpoint
// even when the caller supplies a DFS URL, since CreateSession is a blob endpoint operation.
func TestSessionProviderUsesBlobEndpoint(t *testing.T) {
	_require := require.New(t)

	transport := &endpointRecordingTransport{}
	opts := &azdatalake.ClientOptions{}
	opts.Transport = transport
	opts.Retry = policy.RetryOptions{MaxRetries: -1}

	cred := &staticTokenCredential{}
	provider, err := azdatalake.NewFilesystemSessionProvider(cred, "https://fakeaccount.dfs.core.windows.net/myfs/myfile", opts)
	_require.NoError(err)
	_require.NotNil(provider)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		"https://fakeaccount.dfs.core.windows.net/myfs/myfile", nil)
	_require.NoError(err)

	_, _ = provider.GetSession(req)

	_require.NotEmpty(transport.hosts, "the provider must attempt a CreateSession call")
	_require.Contains(transport.hosts[0], ".blob.", "CreateSession must target the blob endpoint")
	_require.NotContains(transport.hosts[0], ".dfs.")
}

// endpointRecordingTransport records the host of every request it sees.
type endpointRecordingTransport struct {
	mu    sync.Mutex
	hosts []string
}

func (e *endpointRecordingTransport) Do(req *http.Request) (*http.Response, error) {
	e.mu.Lock()
	e.hosts = append(e.hosts, req.URL.Host)
	e.mu.Unlock()
	return &http.Response{
		StatusCode: http.StatusServiceUnavailable,
		Status:     http.StatusText(http.StatusServiceUnavailable),
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader("")),
		Request:    req,
	}, nil
}

// staticTokenCredential hands out a fixed token so no identity provider is contacted.
type staticTokenCredential struct{}

func (staticTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "fake-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}
