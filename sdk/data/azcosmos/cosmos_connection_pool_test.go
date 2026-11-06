// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type cosmosPoolCounters struct {
	live  atomic.Int64
	peak  atomic.Int64
	dials atomic.Int64
}

func (c *cosmosPoolCounters) dial(dial func(context.Context, string, string) (net.Conn, error)) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		c.dials.Add(1)
		n := c.live.Add(1)
		for peak := c.peak.Load(); n > peak && !c.peak.CompareAndSwap(peak, n); peak = c.peak.Load() {
		}
		conn, err := dial(ctx, network, address)
		if err != nil {
			c.live.Add(-1)
			return nil, err
		}
		return &cosmosPoolCountedConn{Conn: conn, counters: c}, nil
	}
}

type cosmosPoolCountedConn struct {
	net.Conn
	counters *cosmosPoolCounters
	once     sync.Once
	err      error
}

func (c *cosmosPoolCountedConn) Close() error {
	c.once.Do(func() {
		c.err = c.Conn.Close()
		c.counters.live.Add(-1)
	})
	return c.err
}

func newPoolTestClient(t testing.TB, server *httptest.Server, limit int) (*http.Client, *cosmosPoolCounters) {
	t.Helper()
	transport := newCosmosHTTPTransport()
	transport.TLSClientConfig.RootCAs = server.Client().Transport.(*http.Transport).TLSClientConfig.RootCAs
	counters := &cosmosPoolCounters{}
	transport.DialContext = counters.dial(transport.DialContext)
	var rt http.RoundTripper = transport
	if limit > 0 {
		rt = newCosmosConnectionPool(transport, limit)
	}
	client := &http.Client{Transport: rt, Timeout: defaultHTTPRoundTripTimeout}
	t.Cleanup(client.CloseIdleConnections)
	return client, counters
}

func newPoolTestServer(t testing.TB, h2 bool, streams int, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewUnstartedServer(handler)
	server.EnableHTTP2 = h2
	server.Config.HTTP2 = &http.HTTP2Config{MaxConcurrentStreams: streams}
	server.StartTLS()
	t.Cleanup(server.Close)
	return server
}

func drainPoolResponse(t testing.TB, response *http.Response) {
	t.Helper()
	_, err := io.Copy(io.Discard, response.Body)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
}

type cosmosPoolResult struct {
	response *http.Response
	err      error
}

func receivePoolResult(t *testing.T, results <-chan cosmosPoolResult) cosmosPoolResult {
	t.Helper()
	select {
	case result := <-results:
		return result
	case <-time.After(10 * time.Second):
		t.Fatal("request did not complete")
		return cosmosPoolResult{}
	}
}

func TestCosmosConnectionPool_Saturation(t *testing.T) {
	for _, h2 := range []bool{false, true} {
		t.Run(strconv.FormatBool(h2), func(t *testing.T) {
			const limit = 32
			release := make(chan struct{})
			var releaseOnce sync.Once
			unblock := func() { releaseOnce.Do(func() { close(release) }) }
			server := newPoolTestServer(t, h2, 1, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/warm" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.(http.Flusher).Flush()
				select {
				case <-release:
				case <-r.Context().Done():
				}
			})
			client, counters := newPoolTestClient(t, server, limit)
			t.Cleanup(unblock)
			warm, err := client.Get(server.URL + "/warm")
			require.NoError(t, err)
			drainPoolResponse(t, warm)

			results := make(chan cosmosPoolResult, limit+8)
			for range limit {
				go func() {
					response, err := client.Get(server.URL)
					results <- cosmosPoolResult{response, err}
				}()
			}
			held := make([]*http.Response, 0, limit)
			for range limit {
				result := receivePoolResult(t, results)
				require.NoError(t, result.err)
				assert.Equal(t, map[bool]int{false: 1, true: 2}[h2], result.response.ProtoMajor)
				held = append(held, result.response)
			}
			assert.EqualValues(t, limit, counters.live.Load())

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			for i := range 8 {
				go func() {
					requestContext := t.Context()
					if i == 0 {
						requestContext = ctx
					}
					req, err := http.NewRequestWithContext(requestContext, http.MethodGet, server.URL, nil)
					if err != nil {
						results <- cosmosPoolResult{err: err}
						return
					}
					response, err := client.Do(req)
					results <- cosmosPoolResult{response, err}
				}()
			}
			pool := client.Transport.(*cosmosConnectionPool)
			require.Eventually(t, func() bool {
				pool.mu.Lock()
				defer pool.mu.Unlock()
				for _, endpoint := range pool.endpoints {
					return endpoint.users == 8
				}
				return false
			}, 5*time.Second, time.Millisecond)
			assert.EqualValues(t, limit, counters.peak.Load())
			cancel()
			canceled := receivePoolResult(t, results)
			require.ErrorIs(t, canceled.err, context.Canceled)
			unblock()
			for _, response := range held {
				drainPoolResponse(t, response)
			}
			for range 7 {
				result := receivePoolResult(t, results)
				require.NoError(t, result.err)
				drainPoolResponse(t, result.response)
			}
			assert.LessOrEqual(t, counters.peak.Load(), int64(limit))
			client.CloseIdleConnections()
			require.Eventually(t, func() bool { return counters.live.Load() == 0 }, 5*time.Second, time.Millisecond)
			pool.mu.Lock()
			assert.Empty(t, pool.endpoints)
			pool.mu.Unlock()
		})
	}
}

func TestCosmosConnectionPool_MultiplexingAndBodyClose(t *testing.T) {
	for _, h2 := range []bool{false, true} {
		t.Run(strconv.FormatBool(h2), func(t *testing.T) {
			server := newPoolTestServer(t, h2, 4, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/warm" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.(http.Flusher).Flush()
				<-r.Context().Done()
			})
			client, counters := newPoolTestClient(t, server, 1)
			warm, err := client.Get(server.URL + "/warm")
			require.NoError(t, err)
			drainPoolResponse(t, warm)
			count := 1
			if h2 {
				count = 4
			}
			var held []*http.Response
			for range count {
				response, err := client.Get(server.URL)
				require.NoError(t, err)
				held = append(held, response)
			}
			t.Cleanup(func() {
				for _, response := range held {
					_ = response.Body.Close()
				}
			})
			client.CloseIdleConnections()
			assert.EqualValues(t, 1, counters.live.Load(), "active responses must survive CloseIdleConnections")
			results := make(chan cosmosPoolResult, 1)
			go func() {
				response, err := client.Get(server.URL)
				results <- cosmosPoolResult{response, err}
			}()
			require.NoError(t, held[0].Body.Close())
			next := receivePoolResult(t, results)
			require.NoError(t, next.err)
			require.NoError(t, next.response.Body.Close())
			assert.EqualValues(t, 1, counters.peak.Load())
		})
	}
}

func TestCosmosConnectionPool_IdleCleanup(t *testing.T) {
	for _, h2 := range []bool{false, true} {
		t.Run(strconv.FormatBool(h2), func(t *testing.T) {
			server := newPoolTestServer(t, h2, 1, func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, "ok")
			})
			client, counters := newPoolTestClient(t, server, 2)
			pool := client.Transport.(*cosmosConnectionPool)
			pool.transport.IdleConnTimeout = 30 * time.Millisecond
			for range 2 {
				response, err := client.Get(server.URL)
				require.NoError(t, err)
				drainPoolResponse(t, response)
				require.Eventually(t, func() bool {
					pool.mu.Lock()
					defer pool.mu.Unlock()
					return counters.live.Load() == 0 && len(pool.endpoints) == 0
				}, 5*time.Second, time.Millisecond)
				pool.mu.Lock()
				assert.Empty(t, pool.endpoints)
				assert.Zero(t, pool.idle.Len())
				pool.mu.Unlock()
			}
		})
	}
}

func TestCosmosConnectionPool_FailedDialAndTLS(t *testing.T) {
	for _, failure := range []string{"dial", "TLS", "canceled"} {
		t.Run(failure, func(t *testing.T) {
			server := newPoolTestServer(t, true, 1, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})
			client, counters := newPoolTestClient(t, server, 1)
			pool := client.Transport.(*cosmosConnectionPool)
			dial := pool.transport.DialContext
			ctx := t.Context()
			injected := errors.New("dial failed")
			switch failure {
			case "dial":
				pool.transport.DialContext = func(context.Context, string, string) (net.Conn, error) { return nil, injected }
			case "TLS":
				pool.transport.TLSClientConfig.ServerName = "not-localhost.invalid"
			case "canceled":
				canceledContext, cancel := context.WithCancel(ctx)
				cancel()
				ctx = canceledContext
			}
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
			require.NoError(t, err)
			_, err = client.Do(req)
			require.Error(t, err)
			if failure == "dial" {
				assert.ErrorIs(t, err, injected)
			}
			if failure == "canceled" {
				assert.ErrorIs(t, err, context.Canceled)
			}
			assert.Zero(t, counters.live.Load())
			pool.mu.Lock()
			assert.Empty(t, pool.endpoints)
			pool.mu.Unlock()
			pool.transport.DialContext = dial
			pool.transport.TLSClientConfig.ServerName = ""
			response, err := client.Get(server.URL)
			require.NoError(t, err)
			drainPoolResponse(t, response)
		})
	}
}

func TestNewHTTPClient_Options(t *testing.T) {
	for _, limit := range []int{0, 1, 32, 64, -1} {
		client, err := NewHTTPClient(&HTTPClientOptions{MaxConnectionsPerEndpoint: limit})
		require.NoError(t, err)
		t.Cleanup(client.CloseIdleConnections)
		if limit == -1 {
			transport, ok := client.Transport.(*http.Transport)
			require.True(t, ok)
			assert.Zero(t, transport.MaxConnsPerHost)
		} else {
			pool, ok := client.Transport.(*cosmosConnectionPool)
			require.True(t, ok)
			if limit == 0 {
				limit = 32
			}
			assert.Equal(t, limit, pool.limit)
		}
	}
	client, err := NewHTTPClient(&HTTPClientOptions{MaxConnectionsPerEndpoint: -2})
	require.Error(t, err)
	assert.Nil(t, client)
}

func TestCosmosConnectionPool_ClosesRequestBodyOnQueueCancellation(t *testing.T) {
	server := newPoolTestServer(t, true, 1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	})
	client, _ := newPoolTestClient(t, server, 1)
	response, err := client.Get(server.URL)
	require.NoError(t, err)
	defer func() { assert.NoError(t, response.Body.Close()) }()
	body := &cosmosPoolTestBody{Reader: strings.NewReader("item")}
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL, body)
	require.NoError(t, err)
	_, err = client.Do(req)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.EqualValues(t, 1, body.closed.Load())
}

type cosmosPoolTestBody struct {
	*strings.Reader
	closed atomic.Int64
}

func (b *cosmosPoolTestBody) Close() error {
	b.closed.Add(1)
	return nil
}

func TestCosmosConnectionPool_EndpointScopeAndIdleLimits(t *testing.T) {
	for _, h2 := range []bool{false, true} {
		t.Run(strconv.FormatBool(h2), func(t *testing.T) {
			release := make(chan struct{})
			var once sync.Once
			unblock := func() { once.Do(func() { close(release) }) }
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.(http.Flusher).Flush()
				select {
				case <-release:
				case <-r.Context().Done():
				}
			}
			first := newPoolTestServer(t, h2, 1, handler)
			second := newPoolTestServer(t, h2, 1, handler)
			client, counters := newPoolTestClient(t, first, 4)
			pool := client.Transport.(*cosmosConnectionPool)
			pool.transport.MaxIdleConnsPerHost = 2
			pool.transport.MaxIdleConns = 3
			t.Cleanup(unblock)
			var held []*http.Response
			for _, endpoint := range []string{first.URL, second.URL} {
				for range 4 {
					response, err := client.Get(endpoint)
					require.NoError(t, err)
					held = append(held, response)
				}
			}
			assert.EqualValues(t, 8, counters.live.Load(), "the limit is per endpoint, not per client")
			unblock()
			for _, response := range held {
				drainPoolResponse(t, response)
			}
			require.Eventually(t, func() bool { return counters.live.Load() <= 3 }, 5*time.Second, time.Millisecond)
			pool.mu.Lock()
			for _, endpoint := range pool.endpoints {
				assert.LessOrEqual(t, endpoint.idle, 2)
			}
			assert.LessOrEqual(t, pool.idle.Len(), 3)
			pool.mu.Unlock()
		})
	}
}

func TestCosmosConnectionPool_DialCancellation(t *testing.T) {
	server := newPoolTestServer(t, true, 1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	for _, handshake := range []bool{false, true} {
		t.Run(strconv.FormatBool(handshake), func(t *testing.T) {
			transport := newCosmosHTTPTransport()
			transport.TLSClientConfig.RootCAs = server.Client().Transport.(*http.Transport).TLSClientConfig.RootCAs
			dial := transport.DialContext
			var first atomic.Bool
			var peer net.Conn
			counters := &cosmosPoolCounters{}
			transport.DialContext = counters.dial(func(ctx context.Context, network, address string) (net.Conn, error) {
				if first.CompareAndSwap(false, true) {
					if handshake {
						conn, other := net.Pipe()
						peer = other
						return conn, nil
					}
					<-ctx.Done()
					return nil, ctx.Err()
				}
				return dial(ctx, network, address)
			})
			pool := newCosmosConnectionPool(transport, 1)
			client := &http.Client{Transport: pool}
			t.Cleanup(client.CloseIdleConnections)
			ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
			require.NoError(t, err)
			_, err = client.Do(req)
			require.ErrorIs(t, err, context.DeadlineExceeded)
			if peer != nil {
				require.NoError(t, peer.Close())
			}
			assert.Zero(t, counters.live.Load())
			pool.mu.Lock()
			assert.Empty(t, pool.endpoints)
			pool.mu.Unlock()
			response, err := client.Get(server.URL)
			require.NoError(t, err)
			drainPoolResponse(t, response)
			assert.EqualValues(t, 1, counters.peak.Load())
		})
	}
}

func TestCosmosConnectionPool_GoAwayRetainsBudget(t *testing.T) {
	release := make(chan struct{})
	var once sync.Once
	unblock := func() { once.Do(func() { close(release) }) }
	server := newPoolTestServer(t, true, 2, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		select {
		case <-release:
		case <-r.Context().Done():
		}
	})
	client, counters := newPoolTestClient(t, server, 1)
	t.Cleanup(unblock)
	response, err := client.Get(server.URL)
	require.NoError(t, err)
	defer func() { assert.NoError(t, response.Body.Close()) }()
	pool := client.Transport.(*cosmosConnectionPool)
	pool.mu.Lock()
	var conn *http.ClientConn
	for _, endpoint := range pool.endpoints {
		conn = endpoint.conns[0].conn
	}
	pool.mu.Unlock()
	require.NotNil(t, conn)
	require.Equal(t, 1, conn.Available())
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	shutdown := make(chan error, 1)
	go func() { shutdown <- server.Config.Shutdown(ctx) }()
	require.Eventually(t, func() bool {
		return conn.Available() == 0 && conn.InFlight() == 1 && conn.Err() == nil
	}, time.Second, time.Millisecond)
	results := make(chan cosmosPoolResult, 1)
	go func() {
		next, err := client.Get(server.URL)
		results <- cosmosPoolResult{next, err}
	}()
	require.Eventually(t, func() bool {
		pool.mu.Lock()
		defer pool.mu.Unlock()
		for _, endpoint := range pool.endpoints {
			return endpoint.users == 1
		}
		return false
	}, time.Second, time.Millisecond)
	assert.EqualValues(t, 1, counters.dials.Load(), "a draining connection must keep its slot")
	unblock()
	drainPoolResponse(t, response)
	result := receivePoolResult(t, results)
	require.Error(t, result.err, "the server listener is closed")
	require.NoError(t, <-shutdown)
	assert.EqualValues(t, 1, counters.peak.Load())
	require.Eventually(t, func() bool { return counters.live.Load() == 0 }, time.Second, time.Millisecond)
}

func TestCosmosConnectionPool_Proxy(t *testing.T) {
	for _, h2 := range []bool{false, true} {
		t.Run(strconv.FormatBool(h2), func(t *testing.T) {
			server := newPoolTestServer(t, h2, 1, func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.WriteString(w, "proxied")
			})
			var tunnels atomic.Int64
			proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !assert.Equal(t, http.MethodConnect, r.Method) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				target, err := net.Dial("tcp", r.Host)
				if !assert.NoError(t, err) {
					w.WriteHeader(http.StatusBadGateway)
					return
				}
				defer func() { _ = target.Close() }()
				connection, _, err := w.(http.Hijacker).Hijack()
				if !assert.NoError(t, err) {
					return
				}
				defer func() { _ = connection.Close() }()
				tunnels.Add(1)
				_, err = io.WriteString(connection, "HTTP/1.1 200 Connection Established\r\n\r\n")
				if !assert.NoError(t, err) {
					return
				}
				done := make(chan struct{})
				go func() {
					_, _ = io.Copy(target, connection)
					_ = target.Close()
					close(done)
				}()
				_, _ = io.Copy(connection, target)
				_ = connection.Close()
				<-done
			}))
			t.Cleanup(proxy.Close)
			proxyURL, err := url.Parse(proxy.URL)
			require.NoError(t, err)
			transport := newCosmosHTTPTransport()
			transport.TLSClientConfig.RootCAs = server.Client().Transport.(*http.Transport).TLSClientConfig.RootCAs
			transport.Proxy = func(req *http.Request) (*url.URL, error) {
				assert.Equal(t, strings.TrimPrefix(server.URL, "https://"), req.URL.Host)
				return proxyURL, nil
			}
			client := &http.Client{Transport: newCosmosConnectionPool(transport, 1), Timeout: 5 * time.Second}
			t.Cleanup(client.CloseIdleConnections)
			for range 2 {
				response, err := client.Get(server.URL)
				require.NoError(t, err)
				assert.Equal(t, map[bool]int{false: 1, true: 2}[h2], response.ProtoMajor)
				drainPoolResponse(t, response)
			}
			assert.EqualValues(t, 1, tunnels.Load())
		})
	}
}

func TestCosmosConnectionPool_DoesNotReplay(t *testing.T) {
	for _, h2 := range []bool{false, true} {
		t.Run(strconv.FormatBool(h2), func(t *testing.T) {
			var requests atomic.Int64
			var server *httptest.Server
			server = newPoolTestServer(t, h2, 1, func(w http.ResponseWriter, r *http.Request) {
				if requests.Add(1) == 1 {
					server.CloseClientConnections()
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})
			client, counters := newPoolTestClient(t, server, 1)
			_, err := client.Post(server.URL, "application/json", strings.NewReader("{}"))
			require.Error(t, err)
			assert.EqualValues(t, 1, requests.Load())
			response, err := client.Get(server.URL)
			require.NoError(t, err)
			drainPoolResponse(t, response)
			assert.EqualValues(t, 1, counters.peak.Load())
		})
	}
}

func TestCosmosConnectionPool_ClientTimeoutIncludesBody(t *testing.T) {
	server := newPoolTestServer(t, true, 1, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	})
	client, _ := newPoolTestClient(t, server, 1)
	client.Timeout = 50 * time.Millisecond
	response, err := client.Get(server.URL)
	require.NoError(t, err)
	_, err = io.ReadAll(response.Body)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	if err := response.Body.Close(); err != nil {
		assert.ErrorContains(t, err, "canceled")
	}

}

func TestCosmosConnectionPool_Trace(t *testing.T) {
	for _, h2 := range []bool{false, true} {
		for _, limit := range []int{0, 1} {
			t.Run(strconv.FormatBool(h2)+"/"+strconv.Itoa(limit), func(t *testing.T) {
				server := newPoolTestServer(t, h2, 1, func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				})
				client, _ := newPoolTestClient(t, server, limit)
				var get, got, wrote atomic.Int64
				trace := &httptrace.ClientTrace{
					GetConn:      func(string) { get.Add(1) },
					GotConn:      func(httptrace.GotConnInfo) { got.Add(1) },
					WroteRequest: func(httptrace.WroteRequestInfo) { wrote.Add(1) },
				}
				for range 3 {
					ctx := httptrace.WithClientTrace(t.Context(), trace)
					request, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
					require.NoError(t, err)
					response, err := client.Do(request)
					require.NoError(t, err)
					drainPoolResponse(t, response)
				}
				t.Logf("GetConn=%d GotConn=%d WroteRequest=%d", get.Load(), got.Load(), wrote.Load())
				if limit > 0 {
					assert.EqualValues(t, 3, get.Load())
					assert.EqualValues(t, 3, got.Load())
					assert.EqualValues(t, 3, wrote.Load())
				}
			})
		}
	}
}
