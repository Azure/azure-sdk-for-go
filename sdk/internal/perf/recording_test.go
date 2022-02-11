// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"net/http"
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
	require.Contains(t, resp.Request.URL.String(), "https://www.bing.com")

	proxyTransportsSuite[t.Name()].SetMode("record")
	req, err = http.NewRequest("POST", "https://www.bing.com", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, "https://localhost:5001", resp.Request.URL.String())
	require.Contains(t, req.Header.Get(upstreamURIHeader), "https://www.bing.com")
	require.Equal(t, req.Header.Get(modeHeader), "record")
	require.Equal(t, req.Header.Get(idHeader), client.recID)

	proxyTransportsSuite[t.Name()].SetMode("playback")
	req, err = http.NewRequest("POST", "https://www.bing.com", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, "https://localhost:5001", resp.Request.URL.String())
	require.Contains(t, req.Header.Get(upstreamURIHeader), "https://www.bing.com")
	require.Equal(t, req.Header.Get(modeHeader), "playback")
	require.Equal(t, req.Header.Get(idHeader), client.recID)
}
