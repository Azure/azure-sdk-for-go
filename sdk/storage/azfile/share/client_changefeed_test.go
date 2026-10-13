// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"github.com/stretchr/testify/require"
)

const (
	changeFeedEnabledHeader   = "x-ms-file-enable-change-feed"
	changeFeedRetentionHeader = "x-ms-file-change-feed-retention-in-days"
	changeFeedContainerHeader = "x-ms-file-blob-container-for-xfiles-change-feed"
	changeFeedContainerPrefix = "$fileschangefeed-"
)

func (s *ShareRecordedTestsSuite) TestShareCreateWithChangeFeed() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.Create(context.Background(), &share.CreateOptions{
		EnableChangeFeed:          to.Ptr(true),
		ChangeFeedRetentionInDays: to.Ptr(int32(3)),
	})
	_require.NoError(err)

	getResp, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(getResp.EnableChangeFeed)
	_require.True(*getResp.EnableChangeFeed)
	_require.NotNil(getResp.ChangeFeedRetentionInDays)
	_require.Equal(int32(3), *getResp.ChangeFeedRetentionInDays)
	_require.NotNil(getResp.ChangeFeedBlobContainerName)
	_require.True(strings.HasPrefix(*getResp.ChangeFeedBlobContainerName, changeFeedContainerPrefix))
}

func (s *ShareRecordedTestsSuite) TestShareSetPropertiesChangeFeed() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{
		EnableChangeFeed:          to.Ptr(true),
		ChangeFeedRetentionInDays: to.Ptr(int32(5)),
	})
	_require.NoError(err)

	getResp, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(getResp.EnableChangeFeed)
	_require.True(*getResp.EnableChangeFeed)
	_require.NotNil(getResp.ChangeFeedRetentionInDays)
	_require.Equal(int32(5), *getResp.ChangeFeedRetentionInDays)
	_require.NotNil(getResp.ChangeFeedBlobContainerName)
	_require.True(strings.HasPrefix(*getResp.ChangeFeedBlobContainerName, changeFeedContainerPrefix))
}

func (s *ShareRecordedTestsSuite) TestShareSetPropertiesChangeFeedRetentionInDaysOnly() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{
		EnableChangeFeed:          to.Ptr(true),
		ChangeFeedRetentionInDays: to.Ptr(int32(5)),
	})
	_require.NoError(err)

	// update only the retention period, leaving change feed enabled
	_, err = shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{
		ChangeFeedRetentionInDays: to.Ptr(int32(30)),
	})
	_require.NoError(err)

	getResp, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(getResp.EnableChangeFeed)
	_require.True(*getResp.EnableChangeFeed)
	_require.NotNil(getResp.ChangeFeedRetentionInDays)
	_require.Equal(int32(30), *getResp.ChangeFeedRetentionInDays)
}

func (s *ShareRecordedTestsSuite) TestShareCreateWithoutChangeFeed() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	getResp, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	if getResp.EnableChangeFeed != nil {
		_require.False(*getResp.EnableChangeFeed)
	}
	_require.Nil(getResp.ChangeFeedBlobContainerName)
}

// changeFeedTransport records the requests it handles and answers each with the configured response headers.
type changeFeedTransport struct {
	requests []*http.Request
	header   http.Header
}

func (c *changeFeedTransport) Do(req *http.Request) (*http.Response, error) {
	c.requests = append(c.requests, req)
	header := c.header
	if header == nil {
		header = http.Header{}
	}
	status := http.StatusOK
	if req.Method == http.MethodPut && req.URL.Query().Get("comp") == "" {
		status = http.StatusCreated
	}
	return &http.Response{
		Request:    req,
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     header.Clone(),
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
}

// rawHeader returns the request header value as set by the generated client, which uses lowercase header keys.
func rawHeader(req *http.Request, key string) string {
	if v := req.Header[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}

func newChangeFeedShareClient(t *testing.T, transport policy.Transporter) *share.Client {
	shareClient, err := share.NewClientWithNoCredential("https://account.file.core.windows.net/share", &share.ClientOptions{ClientOptions: policy.ClientOptions{
		Transport: transport,
		Retry:     policy.RetryOptions{MaxRetries: -1},
	}})
	require.NoError(t, err)
	return shareClient
}

func TestShareChangeFeedRequestHeaders(t *testing.T) {
	transport := &changeFeedTransport{}
	shareClient := newChangeFeedShareClient(t, transport)

	_, err := shareClient.Create(context.Background(), &share.CreateOptions{
		EnableChangeFeed:          to.Ptr(true),
		ChangeFeedRetentionInDays: to.Ptr(int32(7)),
	})
	require.NoError(t, err)
	_, err = shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{
		EnableChangeFeed:          to.Ptr(false),
		ChangeFeedRetentionInDays: to.Ptr(int32(365)),
	})
	require.NoError(t, err)
	_, err = shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{
		ChangeFeedRetentionInDays: to.Ptr(int32(1)),
	})
	require.NoError(t, err)

	require.Len(t, transport.requests, 3)
	require.Equal(t, "true", rawHeader(transport.requests[0], changeFeedEnabledHeader))
	require.Equal(t, "7", rawHeader(transport.requests[0], changeFeedRetentionHeader))
	require.Equal(t, "false", rawHeader(transport.requests[1], changeFeedEnabledHeader))
	require.Equal(t, "365", rawHeader(transport.requests[1], changeFeedRetentionHeader))
	require.NotContains(t, transport.requests[2].Header, changeFeedEnabledHeader)
	require.Equal(t, "1", rawHeader(transport.requests[2], changeFeedRetentionHeader))
}

func TestShareChangeFeedRequestHeadersOmittedByDefault(t *testing.T) {
	transport := &changeFeedTransport{}
	shareClient := newChangeFeedShareClient(t, transport)

	_, err := shareClient.Create(context.Background(), nil)
	require.NoError(t, err)
	_, err = shareClient.Create(context.Background(), &share.CreateOptions{Quota: to.Ptr(int32(10))})
	require.NoError(t, err)
	_, err = shareClient.SetProperties(context.Background(), nil)
	require.NoError(t, err)
	_, err = shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{Quota: to.Ptr(int32(10))})
	require.NoError(t, err)

	require.Len(t, transport.requests, 4)
	for _, req := range transport.requests {
		require.NotContains(t, req.Header, changeFeedEnabledHeader)
		require.NotContains(t, req.Header, changeFeedRetentionHeader)
	}
}

func TestShareGetPropertiesChangeFeedResponseHeaders(t *testing.T) {
	containerName := changeFeedContainerPrefix + "8f2a6c1e-3b4d-4e5f-9a0b-1c2d3e4f5a6b"
	shareClient := newChangeFeedShareClient(t, &changeFeedTransport{header: http.Header{
		http.CanonicalHeaderKey(changeFeedEnabledHeader):   []string{"true"},
		http.CanonicalHeaderKey(changeFeedRetentionHeader): []string{"30"},
		http.CanonicalHeaderKey(changeFeedContainerHeader): []string{containerName},
	}})

	getResp, err := shareClient.GetProperties(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, getResp.EnableChangeFeed)
	require.True(t, *getResp.EnableChangeFeed)
	require.NotNil(t, getResp.ChangeFeedRetentionInDays)
	require.Equal(t, int32(30), *getResp.ChangeFeedRetentionInDays)
	require.NotNil(t, getResp.ChangeFeedBlobContainerName)
	require.Equal(t, containerName, *getResp.ChangeFeedBlobContainerName)
}

func TestShareGetPropertiesChangeFeedResponseHeadersAbsent(t *testing.T) {
	shareClient := newChangeFeedShareClient(t, &changeFeedTransport{})

	getResp, err := shareClient.GetProperties(context.Background(), nil)
	require.NoError(t, err)
	require.Nil(t, getResp.EnableChangeFeed)
	require.Nil(t, getResp.ChangeFeedRetentionInDays)
	require.Nil(t, getResp.ChangeFeedBlobContainerName)
}
