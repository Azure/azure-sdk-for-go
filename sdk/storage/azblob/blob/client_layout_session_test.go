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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
)

// layoutAuthTracker records how each request in a data locality flow was authenticated. Get Blob
// Layout carries comp=layout, which makes it ineligible for session authentication, while the
// blob reads it routes are plain GETs and therefore eligible.
type layoutAuthTracker struct {
	mu sync.Mutex

	createSessionCount int
	layoutBearerCount  int
	layoutSessionCount int
	blobGetBearerCount int
	blobGetSessionCnt  int
}

func (p *layoutAuthTracker) Do(req *policy.Request) (*http.Response, error) {
	raw := req.Raw()
	comp := raw.URL.Query().Get("comp")
	isSessionAuth := strings.HasPrefix(raw.Header.Get("Authorization"), "Session ")

	p.mu.Lock()
	switch {
	case raw.Method == http.MethodPost && comp == "session":
		p.createSessionCount++
	case raw.Method == http.MethodGet && comp == "layout":
		if isSessionAuth {
			p.layoutSessionCount++
		} else {
			p.layoutBearerCount++
		}
	case raw.Method == http.MethodGet && comp == "":
		if isSessionAuth {
			p.blobGetSessionCnt++
		} else {
			p.blobGetBearerCount++
		}
	}
	p.mu.Unlock()

	return req.Next()
}

func (p *layoutAuthTracker) layoutCounts() (bearer, session int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.layoutBearerCount, p.layoutSessionCount
}

func (p *layoutAuthTracker) blobGetCounts() (bearer, session int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.blobGetBearerCount, p.blobGetSessionCnt
}

func (p *layoutAuthTracker) sessionsCreated() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.createSessionCount
}

// newSessionServiceClient returns a service client with session mode enabled and the given tracker
// installed, along with a shared key client used for test setup.
func newSessionServiceClient(s *BlobRecordedTestsSuite, tracker *layoutAuthTracker) (setupClient *service.Client, sessionClient *service.Client, err error) {
	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	setupOptions := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &setupOptions.ClientOptions)
	sharedKeyCred, err := service.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, nil, err
	}
	setupClient, err = service.NewClientWithSharedKeyCredential(serviceURL, sharedKeyCred, setupOptions)
	if err != nil {
		return nil, nil, err
	}

	cred, err := testcommon.GetGenericTokenCredential()
	if err != nil {
		return nil, nil, err
	}
	sessionOptions := &service.ClientOptions{
		Session: azblob.SessionOptions{
			Mode:        azblob.SessionModeEnabled,
			AccountName: accountName,
		},
	}
	testcommon.SetClientOptions(s.T(), &sessionOptions.ClientOptions)
	sessionOptions.PerRetryPolicies = append(sessionOptions.PerRetryPolicies, tracker)
	sessionClient, err = service.NewClient(serviceURL, cred, sessionOptions)
	if err != nil {
		return nil, nil, err
	}
	return setupClient, sessionClient, nil
}

// Get Blob Layout is a comp=layout request, so it must keep using bearer authentication even when
// session mode is enabled, while a plain blob read on the same client uses a session.
func (s *BlobRecordedTestsSuite) TestGetLayoutWithSessionEnabled() {
	_require := require.New(s.T())
	testName := s.T().Name()

	tracker := &layoutAuthTracker{}
	svcClient, sessionSvcClient, err := newSessionServiceClient(s, tracker)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	uploadData := []byte(testcommon.BlockBlobDefaultData)
	_, err = containerClient.NewBlockBlobClient(blobName).Upload(context.Background(), streaming.NopCloser(bytes.NewReader(uploadData)), nil)
	_require.NoError(err)

	sessionBlobClient := sessionSvcClient.NewContainerClient(containerName).NewBlobClient(blobName)

	// Get Blob Layout over the session-enabled client.
	pager := sessionBlobClient.GetLayoutPager(nil)
	_require.NotNil(pager)
	_require.True(pager.More())

	resp, err := pager.NextPage(context.Background())
	_require.NoError(err)
	_require.Equal(len(uploadData), int(*resp.BlobContentLength))
	_require.Equal(1, len(resp.Ranges.Range))
	_require.Equal(1, len(resp.Endpoints.Endpoint))
	_require.Equal(int64(0), *resp.Ranges.Range[0].Start)
	_require.Equal(int64(len(uploadData)-1), *resp.Ranges.Range[0].End)
	_require.Equal(*resp.Endpoints.Endpoint[0].Index, *resp.Ranges.Range[0].EndpointIndex)
	_require.NotNil(*resp.Endpoints.Endpoint[0].Value)
	_require.False(pager.More())

	layoutBearer, layoutSession := tracker.layoutCounts()
	_require.Equal(1, layoutBearer, "expected Get Blob Layout to use bearer auth")
	_require.Equal(0, layoutSession, "comp=layout requests are not eligible for session auth")
	_require.Equal(0, tracker.sessionsCreated(), "expected no session to be created for Get Blob Layout alone")

	// A plain blob read on the same client is session eligible, so it mints and uses a session.
	downloadResp, err := sessionBlobClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)
	downloadedData, err := io.ReadAll(downloadResp.Body)
	_require.NoError(err)
	_require.NoError(downloadResp.Body.Close())
	_require.Equal(uploadData, downloadedData)

	blobGetBearer, blobGetSession := tracker.blobGetCounts()
	_require.Equal(1, tracker.sessionsCreated(), "expected the blob read to create a session")
	_require.Equal(1, blobGetSession, "expected the blob read to use session auth")
	_require.Equal(0, blobGetBearer, "expected no blob read to fall back to bearer auth")
}

// Layout-aware routing and session authentication combined: the layout lookup uses bearer auth and
// every chunk it routes is read with session auth.
func (s *BlobUnrecordedTestsSuite) TestDownloadBufferLayoutAwareRoutingWithSessionEnabled() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	setupOptions := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &setupOptions.ClientOptions)
	sharedKeyCred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)
	svcClient, err := service.NewClientWithSharedKeyCredential(serviceURL, sharedKeyCred, setupOptions)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// A blob large enough to be downloaded in several chunks.
	const blobSize = 16 * 1024 * 1024 // 16 MiB
	const blockSize = 4 * 1024 * 1024 // 4 MiB
	blobName := testcommon.GenerateBlobName(testName)
	blobContentReader, expectedData := testcommon.GenerateData(blobSize)
	_, err = containerClient.NewBlockBlobClient(blobName).UploadStream(context.Background(), blobContentReader,
		&blockblob.UploadStreamOptions{BlockSize: blockSize})
	_require.NoError(err)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	tracker := &layoutAuthTracker{}
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

	sessionBlobClient := sessionSvcClient.NewContainerClient(containerName).NewBlobClient(blobName)

	buff := make([]byte, blobSize)
	n, err := sessionBlobClient.DownloadBuffer(context.Background(), buff, &blob.DownloadBufferOptions{
		LayoutAwareRouting: blob.LayoutAwareRoutingEnabled,
		BlockSize:          blockSize,
		Concurrency:        2,
	})
	_require.NoError(err)
	_require.Equal(int64(blobSize), n)
	_require.Equal(expectedData, buff)

	layoutBearer, layoutSession := tracker.layoutCounts()
	_require.Greater(layoutBearer, 0, "expected the layout lookup to use bearer auth")
	_require.Equal(0, layoutSession, "comp=layout requests are not eligible for session auth")

	blobGetBearer, blobGetSession := tracker.blobGetCounts()
	_require.Equal(1, tracker.sessionsCreated(), "expected exactly one session for the container")
	_require.Equal(blobSize/blockSize, blobGetSession, "expected every routed chunk read to use session auth")
	_require.Equal(0, blobGetBearer, "expected no chunk read to fall back to bearer auth")
}
