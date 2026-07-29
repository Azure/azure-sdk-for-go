// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
)

// authRequestTracker is a pipeline policy that tracks session-related and bearer token requests.
type authRequestTracker struct {
	mu                 sync.Mutex
	createSessionCount int
	sessionAuthCount   int
	bearerAuthCount    int
}

func (p *authRequestTracker) Do(req *policy.Request) (*http.Response, error) {
	p.mu.Lock()

	// Check if this is a CreateSession request (POST with comp=session)
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

func (p *authRequestTracker) counts() (createSessionCount, sessionAuthCount, bearerAuthCount int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.createSessionCount, p.sessionAuthCount, p.bearerAuthCount
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadWithSessionOptions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)

	// Create service client with SharedKeyCredential for setup
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	sharedKeyCred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)
	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, sharedKeyCred, options)
	_require.NoError(err)

	// Create container
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Create and upload blob
	blobName := testcommon.GenerateBlobName(testName)

	// Upload data
	uploadData := []byte("test data for session download")

	bbClient := containerClient.NewBlockBlobClient(blobName)

	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(uploadData)), nil)
	_require.NoError(err)

	// Create service client with TokenCredential and session-based authentication enabled
	sessionOptions := &service.ClientOptions{
		Session: azblob.SessionOptions{
			Mode:        azblob.SessionModeEnabled,
			AccountName: accountName,
		},
	}
	testcommon.SetClientOptions(s.T(), &sessionOptions.ClientOptions)
	sessionSvcClient, err := service.NewClient(serviceURL, cred, sessionOptions)
	_require.NoError(err)
	sessionContClient := sessionSvcClient.NewContainerClient(containerName)
	sessionBlobClient := sessionContClient.NewBlobClient(blobName)

	resp, err := sessionBlobClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	downloadedData, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	err = resp.Body.Close()
	_require.NoError(err)

	_require.Equal(uploadData, downloadedData)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadWithSessionModeOff() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	sharedKeyCred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)
	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, sharedKeyCred, options)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	uploadData := []byte("test data for session mode off")
	bbClient := containerClient.NewBlockBlobClient(blobName)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(uploadData)), nil)
	_require.NoError(err)

	sessionTracker := &authRequestTracker{}

	sessionOptions := &service.ClientOptions{
		Session: azblob.SessionOptions{
			Mode: azblob.SessionModeDisabled,
		},
	}
	testcommon.SetClientOptions(s.T(), &sessionOptions.ClientOptions)
	sessionOptions.PerRetryPolicies = append(sessionOptions.PerRetryPolicies, sessionTracker)
	sessionSvcClient, err := service.NewClient(serviceURL, cred, sessionOptions)
	_require.NoError(err)

	sessionBlobClient := sessionSvcClient.NewContainerClient(containerName).NewBlobClient(blobName)
	resp, err := sessionBlobClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	downloadedData, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_ = resp.Body.Close()
	_require.Equal(uploadData, downloadedData)

	createSessionCount, sessionAuthCount, _ := sessionTracker.counts()

	_require.Equal(0, createSessionCount, "Expected no CreateSession calls when SessionModeDisabled")
	_require.Equal(0, sessionAuthCount, "Expected no session-authenticated requests when SessionModeDisabled")
}

// The container scope of a session is resolved from the request URL, so a client created for the
// service endpoint acquires a session for whichever container the blob request targets.
func (s *BlobRecordedTestsSuite) TestBlobDownloadWithSessionContainerResolvedFromURL() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	sharedKeyCred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)
	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, sharedKeyCred, options)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	uploadData := []byte("test data for container resolved session")
	bbClient := containerClient.NewBlockBlobClient(blobName)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(uploadData)), nil)
	_require.NoError(err)

	sessionTracker := &authRequestTracker{}

	sessionOptions := &service.ClientOptions{
		Session: azblob.SessionOptions{
			Mode:        azblob.SessionModeEnabled,
			AccountName: accountName,
		},
	}
	testcommon.SetClientOptions(s.T(), &sessionOptions.ClientOptions)
	sessionOptions.PerRetryPolicies = append(sessionOptions.PerRetryPolicies, sessionTracker)
	sessionSvcClient, err := service.NewClient(serviceURL, cred, sessionOptions)
	_require.NoError(err)

	sessionBlobClient := sessionSvcClient.NewContainerClient(containerName).NewBlobClient(blobName)
	resp, err := sessionBlobClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	downloadedData, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_ = resp.Body.Close()
	_require.Equal(uploadData, downloadedData)

	createSessionCount, sessionAuthCount, _ := sessionTracker.counts()

	_require.Equal(1, createSessionCount, "Expected a session to be created for the container in the request URL")
	_require.Equal(1, sessionAuthCount, "Expected the download to use session authentication")
}

func (s *BlobUnrecordedTestsSuite) TestBlobDownloadWithSessionOptionsConcurrentDownloads() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)

	// Create service client with SharedKeyCredential for setup
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	sharedKeyCred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)
	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, sharedKeyCred, options)
	_require.NoError(err)

	// Create container
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Create multiple blobs for concurrent download
	const numBlobs = 5
	uploadData := []byte("test data for concurrent session download")
	blobNames := make([]string, numBlobs)

	for i := 0; i < numBlobs; i++ {
		blobNames[i] = fmt.Sprintf("%s-blob-%d", testcommon.GenerateBlobName(testName), i)
		bbClient := containerClient.NewBlockBlobClient(blobNames[i])
		_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(uploadData)), nil)
		_require.NoError(err)
	}

	// Create a policy to track session requests
	sessionTracker := &authRequestTracker{}

	// Create service client with TokenCredential and session-based authentication enabled
	sessionOptions := &service.ClientOptions{
		Session: azblob.SessionOptions{
			Mode:        azblob.SessionModeEnabled,
			AccountName: accountName,
		},
	}
	testcommon.SetClientOptions(s.T(), &sessionOptions.ClientOptions)
	sessionOptions.PerRetryPolicies = append(sessionOptions.PerRetryPolicies, sessionTracker)
	sessionSvcClient, err := service.NewClient(serviceURL, cred, sessionOptions)
	_require.NoError(err)

	sessionContClient := sessionSvcClient.NewContainerClient(containerName)

	// Perform concurrent downloads
	var wg sync.WaitGroup
	errChan := make(chan error, numBlobs)

	for i := 0; i < numBlobs; i++ {
		wg.Add(1)
		go func(blobName string) {
			defer wg.Done()
			sessionBlobClient := sessionContClient.NewBlobClient(blobName)
			resp, err := sessionBlobClient.DownloadStream(context.Background(), nil)
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
				errChan <- fmt.Errorf("downloaded data mismatch for blob %s", blobName)
			}
		}(blobNames[i])
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		_require.NoError(err)
	}

	// Verify that only one CreateSession call was made (session should be cached)
	createSessionCount, sessionAuthCount, _ := sessionTracker.counts()

	_require.Equal(1, createSessionCount, "Expected exactly one CreateSession call due to caching")
	_require.Equal(numBlobs, sessionAuthCount, "Expected all downloads to use session authentication")
}

func (s *BlobUnrecordedTestsSuite) TestBlobDownloadWithSessionOptionsLargeFileDownloadBuffer() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)

	// Create service client with SharedKeyCredential for setup
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	sharedKeyCred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)
	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, sharedKeyCred, options)
	_require.NoError(err)

	// Create container
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Create a large blob (10 MB to trigger chunked download)
	blobName := testcommon.GenerateBlobName(testName)
	const fileSize = 10 * 1024 * 1024 // 10 MB
	uploadData := make([]byte, fileSize)
	for i := range uploadData {
		uploadData[i] = byte(i % 256)
	}

	bbClient := containerClient.NewBlockBlobClient(blobName)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(uploadData)), nil)
	_require.NoError(err)

	// Create a policy to track session requests
	sessionTracker := &authRequestTracker{}

	// Create service client with TokenCredential and session-based authentication enabled
	sessionOptions := &service.ClientOptions{
		Session: azblob.SessionOptions{
			Mode:        azblob.SessionModeEnabled,
			AccountName: accountName,
		},
	}
	testcommon.SetClientOptions(s.T(), &sessionOptions.ClientOptions)
	sessionOptions.PerRetryPolicies = append(sessionOptions.PerRetryPolicies, sessionTracker)
	sessionSvcClient, err := service.NewClient(serviceURL, cred, sessionOptions)
	_require.NoError(err)

	sessionContClient := sessionSvcClient.NewContainerClient(containerName)
	sessionBlobClient := sessionContClient.NewBlobClient(blobName)

	// Download to buffer with small block size to trigger multiple chunks
	buffer := make([]byte, fileSize)
	downloaded, err := sessionBlobClient.DownloadBuffer(context.Background(), buffer, &blob.DownloadBufferOptions{
		BlockSize:   4 * 1024 * 1024, // 4 MB blocks
		Concurrency: 2,
	})
	_require.NoError(err)
	_require.Equal(int64(fileSize), downloaded)
	_require.Equal(uploadData, buffer)

	// Verify session was used
	createSessionCount, sessionAuthCount, _ := sessionTracker.counts()

	_require.Equal(1, createSessionCount, "Expected exactly one CreateSession call")
	_require.Greater(sessionAuthCount, 0, "Expected at least one request with session authentication")
}

func (s *BlobUnrecordedTestsSuite) TestBlobDownloadWithSessionOptionsLargeFileDownloadFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)

	// Create service client with SharedKeyCredential for setup
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	sharedKeyCred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)
	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, sharedKeyCred, options)
	_require.NoError(err)

	// Create container
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Create a large blob (10 MB to trigger chunked download)
	blobName := testcommon.GenerateBlobName(testName)
	const fileSize = 10 * 1024 * 1024 // 10 MB
	uploadData := make([]byte, fileSize)
	for i := range uploadData {
		uploadData[i] = byte(i % 256)
	}

	bbClient := containerClient.NewBlockBlobClient(blobName)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(uploadData)), nil)
	_require.NoError(err)

	// Create a policy to track session requests
	sessionTracker := &authRequestTracker{}

	// Create service client with TokenCredential and session-based authentication enabled
	sessionOptions := &service.ClientOptions{
		Session: azblob.SessionOptions{
			Mode:        azblob.SessionModeEnabled,
			AccountName: accountName,
		},
	}
	testcommon.SetClientOptions(s.T(), &sessionOptions.ClientOptions)
	sessionOptions.PerRetryPolicies = append(sessionOptions.PerRetryPolicies, sessionTracker)
	sessionSvcClient, err := service.NewClient(serviceURL, cred, sessionOptions)
	_require.NoError(err)

	sessionContClient := sessionSvcClient.NewContainerClient(containerName)
	sessionBlobClient := sessionContClient.NewBlobClient(blobName)

	// Create temp file for download
	tmpFile, err := os.CreateTemp("", "session-download-test-*")
	_require.NoError(err)
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	// Download to file with small block size to trigger multiple chunks
	downloaded, err := sessionBlobClient.DownloadFile(context.Background(), tmpFile, &blob.DownloadFileOptions{
		BlockSize:   4 * 1024 * 1024, // 4 MB blocks
		Concurrency: 2,
	})
	_require.NoError(err)
	_require.Equal(int64(fileSize), downloaded)

	// Verify file content
	_, err = tmpFile.Seek(0, io.SeekStart)
	_require.NoError(err)
	downloadedData, err := io.ReadAll(tmpFile)
	_require.NoError(err)
	_require.Equal(uploadData, downloadedData)

	// Verify session was used
	createSessionCount, sessionAuthCount, _ := sessionTracker.counts()

	_require.Equal(1, createSessionCount, "Expected exactly one CreateSession call")
	_require.Greater(sessionAuthCount, 0, "Expected at least one request with session authentication")
}

func (s *BlobUnrecordedTestsSuite) TestBlobDownloadWithSessionOptionsSessionExpiration() {
	s.T().Skip("Skipping: this test sleeps for ~6 minutes to validate session expiration, which is too slow for standard test runs")
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)

	// Create service client with SharedKeyCredential for setup
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	sharedKeyCred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)
	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, sharedKeyCred, options)
	_require.NoError(err)

	// Create container
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Create blob
	blobName := testcommon.GenerateBlobName(testName)
	uploadData := []byte("test data for session expiration test")

	bbClient := containerClient.NewBlockBlobClient(blobName)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(uploadData)), nil)
	_require.NoError(err)

	// Create a policy to track session requests
	sessionTracker := &authRequestTracker{}

	// Create service client with TokenCredential and session-based authentication enabled
	sessionOptions := &service.ClientOptions{
		Session: azblob.SessionOptions{
			Mode:        azblob.SessionModeEnabled,
			AccountName: accountName,
		},
	}
	testcommon.SetClientOptions(s.T(), &sessionOptions.ClientOptions)
	sessionOptions.PerRetryPolicies = append(sessionOptions.PerRetryPolicies, sessionTracker)
	sessionSvcClient, err := service.NewClient(serviceURL, cred, sessionOptions)
	_require.NoError(err)

	sessionContClient := sessionSvcClient.NewContainerClient(containerName)
	sessionBlobClient := sessionContClient.NewBlobClient(blobName)

	// Issue reads every 30 seconds for 6 minutes (13 reads total: 0s, 30s, 60s, ... 360s)
	const readInterval = 30 * time.Second
	const totalDuration = 6 * time.Minute
	numReads := int(totalDuration/readInterval) + 1

	s.T().Logf("Starting session expiration test: %d reads over %v", numReads, totalDuration)

	for i := 0; i < numReads; i++ {
		if i > 0 {
			s.T().Logf("Waiting %v before read %d...", readInterval, i+1)
			time.Sleep(readInterval)
		}

		s.T().Logf("Issuing read %d/%d at %v", i+1, numReads, time.Now().Format(time.RFC3339))

		resp, err := sessionBlobClient.DownloadStream(context.Background(), nil)
		_require.NoError(err)

		downloadedData, err := io.ReadAll(resp.Body)
		_require.NoError(err)
		err = resp.Body.Close()
		_require.NoError(err)

		_require.Equal(uploadData, downloadedData)

		// Log current session counts
		currentCreateCount, currentAuthCount, _ := sessionTracker.counts()
		s.T().Logf("After read %d: CreateSession calls=%d, Session auth requests=%d", i+1, currentCreateCount, currentAuthCount)
	}

	// Verify session was created at least twice (initial + after expiration)
	createSessionCount, sessionAuthCount, _ := sessionTracker.counts()

	s.T().Logf("Final counts: CreateSession calls=%d, Session auth requests=%d", createSessionCount, sessionAuthCount)

	_require.GreaterOrEqual(createSessionCount, 2, "Expected at least 2 CreateSession calls due to session expiration (sessions expire after ~5 minutes)")
	_require.Equal(numReads, sessionAuthCount, "Expected all reads to use session authentication")
}

func (s *BlobUnrecordedTestsSuite) TestBlobDownloadWithSessionOptionsMultipleBlobsSingleSession() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	sharedKeyCred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, sharedKeyCred, options)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobData := []byte("test data for multiple blob session download")
	blobNames := make([]string, 4)
	for i := range blobNames {
		blobNames[i] = fmt.Sprintf("%s-%d", testcommon.GenerateBlobName(testName), i)
		bbClient := containerClient.NewBlockBlobClient(blobNames[i])
		_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(blobData)), nil)
		_require.NoError(err)
	}

	sessionTracker := &authRequestTracker{}

	sessionOptions := &service.ClientOptions{
		Session: azblob.SessionOptions{
			Mode:        azblob.SessionModeEnabled,
			AccountName: accountName,
		},
	}
	testcommon.SetClientOptions(s.T(), &sessionOptions.ClientOptions)
	sessionOptions.PerRetryPolicies = append(sessionOptions.PerRetryPolicies, sessionTracker)
	sessionSvcClient, err := service.NewClient(serviceURL, cred, sessionOptions)
	_require.NoError(err)

	sessionContClient := sessionSvcClient.NewContainerClient(containerName)

	for _, blobName := range blobNames {
		sessionBlobClient := sessionContClient.NewBlobClient(blobName)
		resp, err := sessionBlobClient.DownloadStream(context.Background(), nil)
		_require.NoError(err)

		downloadedData, err := io.ReadAll(resp.Body)
		_require.NoError(err)
		_require.NoError(resp.Body.Close())
		_require.Equal(blobData, downloadedData)
	}

	createSessionCount, sessionAuthCount, _ := sessionTracker.counts()

	_require.Equal(1, createSessionCount, "expected only one session to be created")
	_require.Equal(len(blobNames), sessionAuthCount, "expected each GET to use session auth")
}

func (s *BlobRecordedTestsSuite) TestBlobRandomRestCallsUseBearerExceptGetUsesSession() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	sharedKeyCred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, sharedKeyCred, options)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blobName)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("hello world"))), nil)
	_require.NoError(err)

	tracker := &authRequestTracker{}

	sessionOptions := &service.ClientOptions{
		Session: azblob.SessionOptions{
			Mode:        azblob.SessionModeEnabled,
			AccountName: accountName,
		},
	}
	testcommon.SetClientOptions(s.T(), &sessionOptions.ClientOptions)
	sessionOptions.PerRetryPolicies = append(sessionOptions.PerRetryPolicies, tracker)
	sessionSvcClient, err := service.NewClient(serviceURL, cred, sessionOptions)
	_require.NoError(err)

	sessionContClient := sessionSvcClient.NewContainerClient(containerName)
	sessionBlobClient := sessionContClient.NewBlobClient(blobName)

	// GET should use session auth.
	resp, err := sessionBlobClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)
	_, err = io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.NoError(resp.Body.Close())

	// Non-GET REST calls should use bearer auth.
	_, _ = sessionBlobClient.SetMetadata(context.Background(), map[string]*string{"a": to.Ptr("b")}, nil)

	_, _ = sessionBlobClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{
		BlobContentType: to.Ptr("text/plain"),
	}, nil)

	// Get Tags
	_, _ = sessionBlobClient.GetTags(context.Background(), nil)

	createSessionCount, sessionAuthCount, bearerAuthCount := tracker.counts()

	_require.Equal(1, createSessionCount, "expected a session to be created")
	_require.Equal(1, sessionAuthCount, "expected GET call to use session auth")
	_require.Equal(4, bearerAuthCount, "expected non-GET and create session REST calls to use bearer auth")
}
