// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

type failingReadCloser struct{ err error }

func (r failingReadCloser) Read([]byte) (int, error) { return 0, r.err }
func (failingReadCloser) Close() error               { return nil }

func proxyClient(proxyURL, mode, recordingID string, transport http.RoundTripper) *RecordingHTTPClient {
	return &RecordingHTTPClient{
		defaultClient: &http.Client{Transport: transport},
		options:       TransportOptions{proxyURL: proxyURL},
		mode:          mode,
		recID:         recordingID,
	}
}

func TestProxyTransportConstructionAndRewrite(t *testing.T) {
	defer snapshotFlags(t)()
	client := NewProxyTransport(nil)
	require.NotNil(t, client)
	require.Equal(t, liveMode, client.mode)
	client.SetMode(playbackMode)
	require.Equal(t, playbackMode, client.mode)

	request, err := http.NewRequest(http.MethodGet, "https://upstream.test/path", nil)
	require.NoError(t, err)
	request.Header.Set("original", "value")
	client.options.proxyURL = "https://proxy.test:8443"
	client.recID = "recording"
	rewritten, err := client.replaceAuthority(request)
	require.NoError(t, err)
	require.Equal(t, "upstream.test", request.URL.Host)
	require.Equal(t, "proxy.test:8443", rewritten.URL.Host)
	require.Equal(t, "https://upstream.test", rewritten.Header.Get(upstreamURIHeader))
	require.Equal(t, playbackMode, rewritten.Header.Get(modeHeader))
	require.Equal(t, "recording", rewritten.Header.Get(idHeader))
	require.Empty(t, request.Header.Get(modeHeader))

	client.options.proxyURL = "%"
	_, err = client.replaceAuthority(request)
	require.ErrorContains(t, err, "error parsing url")
	_, err = client.Do(request)
	require.Error(t, err)
}

func TestProxyStartErrors(t *testing.T) {
	networkErr := errors.New("network failed")
	client := proxyClient("%", "record", "", http.DefaultTransport)
	require.ErrorContains(t, client.start(), "creating a START request")

	client = proxyClient("https://proxy.test", "record", "", roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, networkErr
	}))
	require.ErrorContains(t, client.start(), "communicating with the test proxy")

	for _, status := range []int{http.StatusContinue + 99, http.StatusMultipleChoices} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			client := proxyClient("https://proxy.test", "record", "", roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{Request: req, StatusCode: status, Body: io.NopCloser(strings.NewReader("failed")), Header: http.Header{idHeader: []string{"misleading"}}}, nil
			}))
			require.ErrorContains(t, client.start(), "status")
		})
	}

	client = proxyClient("https://proxy.test", "record", "", roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{Request: req, StatusCode: http.StatusOK, Body: failingReadCloser{err: io.ErrUnexpectedEOF}, Header: http.Header{}}, nil
	}))
	require.ErrorContains(t, client.start(), io.ErrUnexpectedEOF.Error())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("no id"))
	}))
	defer server.Close()
	client = &RecordingHTTPClient{defaultClient: server.Client(), options: TransportOptions{proxyURL: server.URL}, mode: "record"}
	require.ErrorContains(t, client.start(), "recording ID was not returned")
}

func TestProxyStopErrors(t *testing.T) {
	client := proxyClient("%", "record", "recording", http.DefaultTransport)
	require.ErrorContains(t, client.stop(), "creating a STOP request")

	client = proxyClient("https://proxy.test", "record", "", http.DefaultTransport)
	require.ErrorContains(t, client.stop(), "recording ID was never set")

	networkErr := errors.New("network failed")
	client = proxyClient("https://proxy.test", "record", "recording", roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, networkErr
	}))
	require.ErrorContains(t, client.stop(), "communicating with the test proxy")

	client = proxyClient("https://proxy.test", "record", "recording", roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{Request: req, StatusCode: http.StatusBadRequest, Body: io.NopCloser(strings.NewReader("stop failed")), Header: http.Header{}}, nil
	}))
	require.ErrorContains(t, client.stop(), "stop failed")

	client = proxyClient("https://proxy.test", "record", "recording", roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{Request: req, StatusCode: http.StatusBadRequest, Body: failingReadCloser{err: io.ErrUnexpectedEOF}, Header: http.Header{}}, nil
	}))
	require.ErrorContains(t, client.stop(), io.ErrUnexpectedEOF.Error())
}
