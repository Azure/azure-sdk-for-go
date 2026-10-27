// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type proxyDiagnostics struct {
	mu     sync.Mutex
	events []string
}

func newProxyDiagnostics(t *testing.T, proxy *httputil.ReverseProxy) *proxyDiagnostics {
	t.Helper()
	d := &proxyDiagnostics{}
	t.Cleanup(func() {
		if t.Failed() {
			d.mu.Lock()
			defer d.mu.Unlock()
			for _, event := range d.events {
				t.Log(event)
			}
		}
	})
	transport := proxy.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	proxy.Transport = &diagnosticProxyTransport{base: transport, diagnostics: d}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		d.record("proxy failure method=%s path=%s error=%s", r.Method, r.URL.Path, proxyTransportError(err))
		w.WriteHeader(http.StatusBadGateway)
	}
	return d
}

func (d *proxyDiagnostics) record(format string, args ...any) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.events) >= 512 {
		if len(d.events) == 512 {
			d.events = append(d.events, "proxy diagnostics truncated after 512 events")
		}
		return
	}
	d.events = append(d.events, time.Now().UTC().Format(time.RFC3339Nano)+" "+fmt.Sprintf(format, args...))
}

func TestProxyDiagnostics(t *testing.T) {
	for _, unavailable := range []bool{false, true} {
		t.Run(fmt.Sprintf("upstream unavailable=%t", unavailable), func(t *testing.T) {
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
			t.Cleanup(upstream.Close)
			target, err := url.Parse(upstream.URL)
			require.NoError(t, err)
			proxy := httputil.NewSingleHostReverseProxy(target)
			d := newProxyDiagnostics(t, proxy)
			server := httptest.NewUnstartedServer(proxy)
			d.observeServer(server)
			server.Start()
			t.Cleanup(server.Close)
			if unavailable {
				upstream.Close()
			}
			request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/?secret=hidden", nil)
			require.NoError(t, err)
			request.Header.Set("Authorization", "synthetic-private-header")
			response, err := server.Client().Do(request)
			require.NoError(t, err)
			_, err = io.Copy(io.Discard, response.Body)
			require.NoError(t, err)
			require.NoError(t, response.Body.Close())
			server.Close()
			d.mu.Lock()
			events := strings.Join(d.events, "\n")
			d.mu.Unlock()
			require.Contains(t, events, "request arrived")
			require.Contains(t, events, "state=new")
			if unavailable {
				require.Equal(t, http.StatusBadGateway, response.StatusCode)
				require.Contains(t, events, "upstream round trip failed")
				require.Contains(t, events, "proxy failure")
			} else {
				require.Equal(t, http.StatusNoContent, response.StatusCode)
				require.Contains(t, events, "upstream response status=204")
			}
			require.NotContains(t, events, "synthetic-private-header")
			require.NotContains(t, events, "secret=hidden")
		})
	}
}

func (d *proxyDiagnostics) observeServer(server *httptest.Server) {
	handler := server.Config.Handler
	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d.record("request arrived peer=%s method=%s path=%s protocol=%s", r.RemoteAddr, r.Method, r.URL.Path, r.Proto)
		handler.ServeHTTP(w, r)
		d.record("handler returned peer=%s method=%s path=%s", r.RemoteAddr, r.Method, r.URL.Path)
	})
	server.Config.ConnState = func(conn net.Conn, state http.ConnState) {
		d.record("connection peer=%s local=%s state=%s", conn.RemoteAddr(), conn.LocalAddr(), state)
	}
	d.record("proxy listener=%s", server.Listener.Addr())
}

func (d *proxyDiagnostics) operation(_ context.Context, diagnostic OperationDiagnostic) {
	if diagnostic.Diagnostics != nil {
		d.record("native diagnostics=%s", diagnostic.Diagnostics.JSON)
	}
}

type diagnosticProxyTransport struct {
	base        http.RoundTripper
	diagnostics *proxyDiagnostics
}

func (t *diagnosticProxyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	d := t.diagnostics
	trace := &httptrace.ClientTrace{
		ConnectStart: func(network, addr string) { d.record("upstream connect network=%s address=%s", network, addr) },
		ConnectDone: func(network, addr string, err error) {
			d.record("upstream connected network=%s address=%s error=%v", network, addr, err)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			d.record("upstream connection peer=%s reused=%t idle=%t", info.Conn.RemoteAddr(), info.Reused, info.WasIdle)
		},
		WroteRequest:         func(info httptrace.WroteRequestInfo) { d.record("upstream request written error=%v", info.Err) },
		GotFirstResponseByte: func() { d.record("upstream first response byte") },
	}
	response, err := t.base.RoundTrip(r.WithContext(httptrace.WithClientTrace(r.Context(), trace)))
	if err != nil {
		d.record("upstream round trip failed: %s", proxyTransportError(err))
	} else {
		d.record("upstream response status=%d protocol=%s", response.StatusCode, response.Proto)
	}

	return response, err
}

func proxyTransportError(err error) string {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return fmt.Sprintf("%s: %v", urlErr.Op, urlErr.Err)
	}
	return fmt.Sprintf("%T: %v", err, err)
}
