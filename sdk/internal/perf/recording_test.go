// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const proxyManualStartEnv = "PROXY_MANUAL_START"

func TestRecordingHTTPClient_Do(t *testing.T) {
	// Ignore manual start in pipeline tests, we always want to exercise install
	os.Setenv(proxyManualStartEnv, "false")

	proxy, err := recording.StartTestProxy("", nil)
	require.NoError(t, err)
	defer func() {
		err := recording.StopTestProxy(proxy)
		if err != nil {
			panic(err)
		}
	}()

	req, err := http.NewRequest("POST", "https://www.bing.com", nil)
	require.NoError(t, err)

	proxyURL := fmt.Sprintf("https://localhost:%d", proxy.Options.ProxyPort)
	client := NewProxyTransport(&TransportOptions{
		TestName: t.Name(),
		proxyURL: proxyURL,
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
	require.Equal(t, proxyURL, resp.Request.URL.String())
	require.Contains(t, resp.Request.Header.Get(upstreamURIHeader), "https://www.bing.com")
	require.Equal(t, resp.Request.Header.Get(modeHeader), "record")
	require.Equal(t, resp.Request.Header.Get(idHeader), client.recID)

	proxyTransportsSuite[t.Name()].SetMode("playback")
	req, err = http.NewRequest("POST", "https://www.bing.com", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, proxyURL, resp.Request.URL.String())
	require.Contains(t, resp.Request.Header.Get(upstreamURIHeader), "https://www.bing.com")
	require.Equal(t, resp.Request.Header.Get(modeHeader), "playback")
	require.Equal(t, resp.Request.Header.Get(idHeader), client.recID)
}
