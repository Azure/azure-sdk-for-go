//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestSetMatcherRecordingOptions(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	parsed, err := url.Parse(srv.URL())
	require.NoError(t, err)
	port, err := strconv.ParseInt(parsed.Port(), 10, 0)
	require.NoError(t, err)
	ro := RecordingOptions{ProxyPort: int(port)}
	t.Run("SetBodilessMatcher", func(t *testing.T) {
		err := SetBodilessMatcher(t, &MatcherOptions{RecordingOptions: ro})
		require.NoError(t, err)
	})
	t.Run("SetDefaultMatcher", func(t *testing.T) {
		err = SetDefaultMatcher(nil, &SetDefaultMatcherOptions{RecordingOptions: ro})
		require.NoError(t, err)
	})
}

type matchersTests struct {
	suite.Suite
	proxy *TestProxyInstance
}

func TestMatchers(t *testing.T) {
	suite.Run(t, new(matchersTests))
}

func (s *matchersTests) SetupSuite() {
	// Ignore manual start in pipeline tests, we always want to exercise install
	os.Setenv(proxyManualStartEnv, "false")
	proxy, err := StartTestProxy("", nil)
	s.proxy = proxy
	require.NoError(s.T(), err)
}

func (s *matchersTests) TearDownSuite() {
	err1 := StopTestProxy(s.proxy)
	err2 := os.RemoveAll("./testdata/recordings/TestMatchers/")
	require.NoError(s.T(), err1)
	require.NoError(s.T(), err2)
}

func (s *matchersTests) TestSetBodilessMatcher() {
	require := require.New(s.T())
	temp := recordMode
	recordMode = RecordingMode
	defer func() { recordMode = temp }()

	err := Start(s.T(), packagePath, nil)
	require.NoError(err)

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	_, err = client.Do(req)
	require.NoError(err)

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Run a second request to with different body to verify it works
	recordMode = PlaybackMode

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	err = SetBodilessMatcher(s.T(), nil)
	require.NoError(err)

	req, err = http.NewRequest("POST", defaultOptions().baseURL(), bytes.NewReader([]byte("abcdef")))
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))

	_, err = client.Do(req)
	require.NoError(err)

	err = Stop(s.T(), nil)
	require.NoError(err)

	err = ResetProxy(nil)
	require.NoError(err)
}

func (s *matchersTests) TestSetBodilessMatcherNilTest() {
	require := require.New(s.T())
	temp := recordMode
	recordMode = RecordingMode
	defer func() { recordMode = temp }()

	err := Start(s.T(), packagePath, nil)
	require.NoError(err)

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	_, err = client.Do(req)
	require.NoError(err)

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Run a second request to with different body to verify it works
	recordMode = PlaybackMode

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	err = SetBodilessMatcher(nil, nil)
	require.NoError(err)

	req, err = http.NewRequest("POST", defaultOptions().baseURL(), bytes.NewReader([]byte("abcdef")))
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))

	_, err = client.Do(req)
	require.NoError(err)

	err = Stop(s.T(), nil)
	require.NoError(err)

	err = ResetProxy(nil)
	require.NoError(err)
}

func (s *matchersTests) TestSetDefaultMatcher() {
	require := require.New(s.T())
	temp := recordMode
	recordMode = RecordingMode
	defer func() { recordMode = temp }()

	err := Start(s.T(), packagePath, nil)
	require.NoError(err)

	req, err := http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))

	client, err := GetHTTPClient(s.T())
	require.NoError(err)

	_, err = client.Do(req)
	require.NoError(err)

	err = Stop(s.T(), nil)
	require.NoError(err)

	// Run a second request to with different body to verify it works
	recordMode = PlaybackMode

	err = Start(s.T(), packagePath, nil)
	require.NoError(err)

	err = SetDefaultMatcher(nil, &SetDefaultMatcherOptions{ExcludedHeaders: []string{"ExampleHeader"}})
	require.NoError(err)

	req, err = http.NewRequest("POST", defaultOptions().baseURL(), nil)
	require.NoError(err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(s.T()))
	req.Header.Set("ExampleHeader", "blah-blah-blah")

	err = handleProxyResponse(client.Do(req))
	require.NoError(err)

	err = Stop(s.T(), nil)
	require.NoError(err)

	err = ResetProxy(nil)
	require.NoError(err)
}

func (s *matchersTests) TestAddDefaults() {
	require := require.New(s.T())
	require.Equal(4, len(addDefaults([]string{})))
	require.Equal(4, len(addDefaults([]string{":path"})))
	require.Equal(4, len(addDefaults([]string{":path", ":authority"})))
	require.Equal(4, len(addDefaults([]string{":path", ":authority", ":method"})))
	require.Equal(4, len(addDefaults([]string{":path", ":authority", ":method", ":scheme"})))
	require.Equal(5, len(addDefaults([]string{":path", ":authority", ":method", ":scheme", "extra"})))
}
