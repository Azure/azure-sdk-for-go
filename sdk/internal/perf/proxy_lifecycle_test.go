// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type proxyNoOpTest struct{}

func (proxyNoOpTest) Run(context.Context) error     { return nil }
func (proxyNoOpTest) Cleanup(context.Context) error { return nil }

func TestProxyPlaybackIsStopped(t *testing.T) {
	var mu sync.Mutex
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests[r.URL.Path]++
		mu.Unlock()
		if r.URL.Path == "/record/start" || r.URL.Path == "/playback/start" {
			w.Header().Set(idHeader, "recording-id")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := &RecordingHTTPClient{
		defaultClient: server.Client(),
		options:       TransportOptions{TestName: "ProxyTest-0", proxyURL: server.URL},
		mode:          liveMode,
	}
	runner := newPerfRunner(PerfMethods{}, "ProxyTest")
	runner.tests = []PerfTest{proxyNoOpTest{}}
	runner.proxyTransports["ProxyTest-0"] = transport

	require.NoError(t, runner.bootstrapProxies())
	require.NoError(t, runner.stopProxyPlaybacks())

	mu.Lock()
	defer mu.Unlock()
	require.Equal(t, 1, requests["/record/start"])
	require.Equal(t, 1, requests["/record/stop"])
	require.Equal(t, 1, requests["/playback/start"])
	require.Equal(t, 1, requests["/playback/stop"])
}

func TestProxyStartRejectsErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(idHeader, "misleading-id")
		http.Error(w, "proxy unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()
	client := &RecordingHTTPClient{
		defaultClient: server.Client(),
		options:       TransportOptions{proxyURL: server.URL},
		mode:          playbackMode,
		recID:         "recording-id",
	}

	err := client.start()

	require.Error(t, err)
	require.Contains(t, err.Error(), "status 503")
}
