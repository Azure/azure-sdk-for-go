//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetBodilessMatcher(t *testing.T) {
	temp := recordMode
	recordMode = RecordingMode
	defer func() { recordMode = temp }()

	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.NoError(t, err)

	err = Stop(t, nil)
	require.NoError(t, err)

	// Run a second request to with different body to verify it works
	recordMode = PlaybackMode

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	err = SetBodilessMatcher(t, nil)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", "https://localhost:5001", bytes.NewReader([]byte("abcdef")))
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	_, err = client.Do(req)
	require.NoError(t, err)

	err = Stop(t, nil)
	require.NoError(t, err)

	err = ResetProxy(nil)
	require.NoError(t, err)
}

func TestSetBodilessMatcherNilTest(t *testing.T) {
	temp := recordMode
	recordMode = RecordingMode
	defer func() { recordMode = temp }()

	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.NoError(t, err)

	err = Stop(t, nil)
	require.NoError(t, err)

	// Run a second request to with different body to verify it works
	recordMode = PlaybackMode

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	err = SetBodilessMatcher(nil, nil)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", "https://localhost:5001", bytes.NewReader([]byte("abcdef")))
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	_, err = client.Do(req)
	require.NoError(t, err)

	err = Stop(t, nil)
	require.NoError(t, err)

	err = ResetProxy(nil)
	require.NoError(t, err)
}

func TestSetDefaultMatcher(t *testing.T) {
	temp := recordMode
	recordMode = RecordingMode
	defer func() { recordMode = temp }()

	err := Start(t, packagePath, nil)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))

	client, err := GetHTTPClient(t)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.NoError(t, err)

	err = Stop(t, nil)
	require.NoError(t, err)

	// Run a second request to with different body to verify it works
	recordMode = PlaybackMode

	err = Start(t, packagePath, nil)
	require.NoError(t, err)

	err = SetDefaultMatcher(nil, &SetDefaultMatcherOptions{ExcludedHeaders: []string{"ExampleHeader"}})
	require.NoError(t, err)

	req, err = http.NewRequest("POST", "https://localhost:5001", nil)
	require.NoError(t, err)

	req.Header.Set(UpstreamURIHeader, "https://bing.com")
	req.Header.Set(ModeHeader, GetRecordMode())
	req.Header.Set(IDHeader, GetRecordingId(t))
	req.Header.Set("ExampleHeader", "blah-blah-blah")

	err = handleProxyResponse(client.Do(req))
	require.NoError(t, err)

	err = Stop(t, nil)
	require.NoError(t, err)

	err = ResetProxy(nil)
	require.NoError(t, err)
}

func TestAddDefaults(t *testing.T) {
	require.Equal(t, 4, len(addDefaults([]string{})))
	require.Equal(t, 4, len(addDefaults([]string{":path"})))
	require.Equal(t, 4, len(addDefaults([]string{":path", ":authority"})))
	require.Equal(t, 4, len(addDefaults([]string{":path", ":authority", ":method"})))
	require.Equal(t, 4, len(addDefaults([]string{":path", ":authority", ":method", ":scheme"})))
	require.Equal(t, 5, len(addDefaults([]string{":path", ":authority", ":method", ":scheme", "extra"})))
}
