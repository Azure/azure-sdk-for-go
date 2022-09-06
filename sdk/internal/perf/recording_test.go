// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecordingHTTPClient_Do(t *testing.T) {
	req, err := http.NewRequest("POST", "https://www.bing.com", nil)
	require.NoError(t, err)

	client := NewProxyTransport(&TransportOptions{
		TestName: t.Name(),
		proxyURL: "https://localhost:5001/",
	})
	require.NotNil(t, client)

	resp, err := client.Do(req)
	require.NoError(t, err)
	matched, err := regexp.MatchString(`https://www\d*\.bing\.com`, resp.Request.URL.String())
	require.NoError(t, err)
	require.True(t, matched)

	proxyTransportsSuite[t.Name()].SetMode("record")
	req, err = http.NewRequest("POST", "https://www.bing.com", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, "https://localhost:5001", resp.Request.URL.String())
	require.Contains(t, resp.Request.Header.Get(upstreamURIHeader), "https://www.bing.com")
	require.Equal(t, resp.Request.Header.Get(modeHeader), "record")
	require.Equal(t, resp.Request.Header.Get(idHeader), client.recID)

	proxyTransportsSuite[t.Name()].SetMode("playback")
	req, err = http.NewRequest("POST", "https://www.bing.com", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, "https://localhost:5001", resp.Request.URL.String())
	require.Contains(t, resp.Request.Header.Get(upstreamURIHeader), "https://www.bing.com")
	require.Equal(t, resp.Request.Header.Get(modeHeader), "playback")
	require.Equal(t, resp.Request.Header.Get(idHeader), client.recID)
}
