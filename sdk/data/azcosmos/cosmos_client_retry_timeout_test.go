// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type timeoutTestTransport func(*http.Request) (*http.Response, error)

func (f timeoutTestTransport) Do(r *http.Request) (*http.Response, error) { return f(r) }

type timeoutTestAttemptCounter struct{ attempts int }

func (c *timeoutTestAttemptCounter) Do(req *policy.Request) (*http.Response, error) {
	c.attempts++
	return req.Next()
}

type timeoutTestVeto struct{ error }

func (*timeoutTestVeto) NonRetriable()   {}
func (v *timeoutTestVeto) Unwrap() error { return v.error }

func timeoutTestResponse(r *http.Request, status int) *http.Response {
	return &http.Response{StatusCode: status, Status: http.StatusText(status), Header: http.Header{},
		Request: r, Body: io.NopCloser(strings.NewReader(`{"id":"item"}`))}
}

func timeoutTestContainer(t *testing.T, transport policy.Transporter, opts policy.RetryOptions, enabled bool) (*ContainerClient, *timeoutTestAttemptCounter) {
	t.Helper()
	endpoint, err := url.Parse("https://timeout.invalid/")
	require.NoError(t, err)
	lc := CreateMockLC(*endpoint, false)
	lc.enableCrossRegionRetries = enabled
	gem := &globalEndpointManager{clientEndpoint: endpoint.String(), locationCache: lc}
	counter := &timeoutTestAttemptCounter{}
	internal, err := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{
		PerCall:  []policy.Policy{&retryCallerContextPolicy{}},
		PerRetry: []policy.Policy{counter, &clientRetryPolicy{gem: gem}},
	}, &policy.ClientOptions{Transport: transport, Retry: opts})
	require.NoError(t, err)
	client := &Client{endpoint: endpoint.String(), endpointUrl: endpoint, internal: internal, gem: gem}
	container, err := client.NewContainer("db", "coll")
	require.NoError(t, err)
	return container, counter
}

func timeoutTestCall(ctx context.Context, c *ContainerClient, op string) error {
	pk := NewPartitionKeyString("pk")
	switch op {
	case "create":
		_, err := c.CreateItem(ctx, pk, []byte(`{"id":"item"}`), nil)
		return err
	case "upsert":
		_, err := c.UpsertItem(ctx, pk, []byte(`{"id":"item"}`), nil)
		return err
	case "patch":
		ops := PatchOperations{}
		ops.AppendIncrement("/counter", 1)
		_, err := c.PatchItem(ctx, pk, "item", ops, nil)
		return err
	case "batch":
		batch := c.NewTransactionalBatch(pk)
		batch.ReadItem("item", nil)
		_, err := c.ExecuteTransactionalBatch(ctx, batch, nil)
		return err
	case "query":
		_, err := c.NewQueryItemsPager("SELECT * FROM c", pk, nil).NextPage(ctx)
		return err
	default:
		_, err := c.ReadItem(ctx, pk, "item", nil)
		return err
	}
}

func TestReadAttemptTimeoutTLS(t *testing.T) {
	for _, tc := range []struct {
		op, phase string
		retry     bool
	}{
		{"read", "headers", true}, {"read", "body", true}, {"query", "headers", true},
		{"query", "body", false}, {"create", "headers", false}, {"create", "body", false},
		{"upsert", "headers", false}, {"upsert", "body", false}, {"patch", "headers", false}, {"patch", "body", false},
	} {
		t.Run(tc.op+"/"+tc.phase, func(t *testing.T) {
			var sends atomic.Int32
			stop := make(chan struct{})
			var server *httptest.Server
			server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path == "/" {
					_, err := fmt.Fprintf(w, `{"writableLocations":[{"name":"East US","databaseAccountEndpoint":%q}],"readableLocations":[{"name":"East US","databaseAccountEndpoint":%q}]}`, server.URL+"/", server.URL+"/")
					assert.NoError(t, err)
					return
				}
				if sends.Add(1) == 1 {
					if tc.phase == "body" {
						if tc.op == "create" {
							w.WriteHeader(http.StatusCreated)
						} else {
							w.WriteHeader(http.StatusOK)
						}
						w.(http.Flusher).Flush()
					}
					select {
					case <-r.Context().Done():
					case <-stop:
					}
					return
				}
				if tc.op == "query" {
					_, err := fmt.Fprint(w, `{"Documents":[{"id":"item"}]}`)
					assert.NoError(t, err)
				} else {
					_, err := fmt.Fprint(w, `{"id":"item"}`)
					assert.NoError(t, err)
				}
			}))
			defer func() { close(stop); server.Close() }()
			key, err := NewKeyCredential("dGVzdA==")
			require.NoError(t, err)
			client, err := NewClientWithKey(server.URL, key, &ClientOptions{ClientOptions: azcore.ClientOptions{
				Transport: server.Client(),
				Retry: policy.RetryOptions{MaxRetries: 3, TryTimeout: 100 * time.Millisecond,
					RetryDelay: 1100 * time.Millisecond, MaxRetryDelay: 1100 * time.Millisecond},
			}})
			require.NoError(t, err)
			defer client.Close()
			container, err := client.NewContainer("db", "coll")
			require.NoError(t, err)
			start := time.Now()
			err = timeoutTestCall(context.Background(), container, tc.op)
			t.Logf("sends=%d elapsed=%s error=%v", sends.Load(), time.Since(start), err)
			if tc.retry {
				require.NoError(t, err)
				assert.Equal(t, int32(2), sends.Load())
				assert.GreaterOrEqual(t, time.Since(start), time.Second)
			} else {
				require.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Equal(t, int32(1), sends.Load())
			}
		})
	}
}

func TestReadAttemptTimeoutBoundaries(t *testing.T) {
	for _, mode := range []string{"recover", "persistent", "persistent-deadline", "408", "wait-only", "disabled", "no-outer", "outer-veto", "veto", "veto-backoff", "query-body", "caller-transport", "caller-backoff", "caller-cancel-transport", "caller-cancel-backoff", "caller-outer", "caller-before", "write-dial", "write-timeout", "write-no-outer", "batch"} {
		t.Run(mode, func(t *testing.T) {
			sends := 0
			marker := errors.New("synthetic transport failure")
			veto := &timeoutTestVeto{marker}
			var firstSend, secondSend time.Time
			transport := timeoutTestTransport(func(r *http.Request) (*http.Response, error) {
				sends++
				if sends == 1 {
					firstSend = time.Now()
				}
				if sends == 2 {
					secondSend = time.Now()
				}
				if mode == "408" {
					return timeoutTestResponse(r, 408), nil
				}
				if mode == "query-body" {
					resp := timeoutTestResponse(r, 200)
					resp.Body = &timeoutTestBody{ctx: r.Context(), err: marker}
					return resp, nil
				}
				if sends > 1 && mode != "persistent" && mode != "persistent-deadline" {
					return timeoutTestResponse(r, 200), nil
				}
				if mode == "caller-outer" {
					return nil, marker
				}
				if mode == "write-dial" {
					<-r.Context().Done()
					return nil, &net.OpError{Op: "dial", Net: "tcp", Err: r.Context().Err()}
				}
				if mode == "veto" || mode == "persistent-deadline" || mode == "write-timeout" || mode == "write-no-outer" || mode == "batch" || mode == "caller-transport" || mode == "caller-cancel-transport" {
					<-r.Context().Done()
					if mode == "veto" {
						return nil, &url.Error{Op: "GET", URL: r.URL.String(), Err: fmt.Errorf("nested: %w", errors.Join(r.Context().Err(), veto))}
					}
					return nil, r.Context().Err()
				}
				if mode == "veto-backoff" {
					return nil, &url.Error{Op: "GET", URL: r.URL.String(), Err: fmt.Errorf("nested: %w", errors.Join(io.EOF, veto))}
				}
				return nil, &url.Error{Op: "GET", URL: r.URL.String(), Err: errors.Join(io.EOF, marker)}
			})
			opts := policy.RetryOptions{MaxRetries: 10, TryTimeout: 100 * time.Millisecond,
				RetryDelay: 1100 * time.Millisecond, MaxRetryDelay: 1100 * time.Millisecond}
			if mode == "wait-only" {
				opts.MaxRetries = 3
				opts.RetryDelay = time.Millisecond
				opts.MaxRetryDelay = time.Millisecond
			}
			if mode == "no-outer" || mode == "write-no-outer" {
				opts.MaxRetries = -1
			}
			if mode == "outer-veto" {
				opts.ShouldRetry = func(*http.Response, error) bool { return false }
			}
			ctx := context.Background()
			cancel := func() {}
			if strings.HasPrefix(mode, "caller-") {
				if mode == "caller-before" {
					ctx, cancel = context.WithCancel(ctx)
					cancel()
				} else if strings.HasPrefix(mode, "caller-cancel-") {
					ctx, cancel = context.WithCancel(ctx)
					timer := time.AfterFunc(50*time.Millisecond, cancel)
					defer timer.Stop()
				} else {
					duration := 50 * time.Millisecond
					if mode == "caller-outer" {
						duration = 150 * time.Millisecond
					}
					ctx, cancel = context.WithTimeout(ctx, duration)
				}
			}

			defer cancel()
			c, counter := timeoutTestContainer(t, transport, opts, mode != "disabled")
			op := "read"
			if strings.HasPrefix(mode, "write-") {
				op = "create"
			} else if mode == "batch" {
				op = "batch"
			} else if mode == "query-body" {
				op = "query"
			}
			err := timeoutTestCall(ctx, c, op)
			t.Logf("outer=%d sends=%d error=%v", counter.attempts, sends, err)
			switch mode {
			case "recover":
				require.NoError(t, err)
				assert.Equal(t, 2, sends)
				assert.Equal(t, 2, counter.attempts)
				assert.GreaterOrEqual(t, secondSend.Sub(firstSend), time.Second)
			case "persistent", "persistent-deadline":
				require.Error(t, err)
				assert.Equal(t, 5, sends)
				assert.Equal(t, 5, counter.attempts)
			case "408":
				require.Error(t, err)
				assert.Equal(t, 2, sends)
				assert.Equal(t, 2, counter.attempts)
			case "wait-only":
				require.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Equal(t, 1, sends)
				assert.Equal(t, 4, counter.attempts)
			case "veto", "veto-backoff":
				require.ErrorIs(t, err, marker)
				var got *timeoutTestVeto
				require.ErrorAs(t, err, &got)
				assert.Same(t, veto, got)
				assert.Equal(t, 1, sends)
				assert.Equal(t, 1, counter.attempts)
				var original *url.Error
				assert.ErrorAs(t, err, &original)
			case "query-body":
				require.ErrorIs(t, err, marker)
				require.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Equal(t, 1, sends)
				assert.Equal(t, 1, counter.attempts)
			default:
				require.Error(t, err)
				if strings.HasPrefix(mode, "caller-") {
					assert.ErrorIs(t, err, ctx.Err())
				}
				assert.Equal(t, 1, sends)
				assert.Equal(t, 1, counter.attempts)
			}
		})
	}
}

func TestWriteAttemptTimeoutDialProvenance(t *testing.T) {
	var dials atomic.Int32
	httpTransport := &http.Transport{DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
		dials.Add(1)
		<-ctx.Done()
		return nil, &net.OpError{Op: "dial", Net: "tcp", Err: ctx.Err()}
	}}
	defer httpTransport.CloseIdleConnections()
	httpClient := &http.Client{Transport: httpTransport}
	var original error
	transport := timeoutTestTransport(func(r *http.Request) (*http.Response, error) {
		resp, err := httpClient.Do(r)
		original = err
		return resp, err
	})
	c, attempts := timeoutTestContainer(t, transport, policy.RetryOptions{
		MaxRetries: 3, TryTimeout: 100 * time.Millisecond, RetryDelay: time.Millisecond,
	}, true)
	err := timeoutTestCall(context.Background(), c, "create")
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.ErrorIs(t, original, context.DeadlineExceeded)
	var opErr *net.OpError
	assert.False(t, errors.As(original, &opErr), "net/http discards dial provenance on native cancellation")
	assert.Equal(t, connectionErrorAmbiguous, classifyNetworkError(original))
	assert.Equal(t, 1, attempts.attempts)
	assert.Equal(t, int32(1), dials.Load())
}

func TestReadAttemptTimeoutResumesRoutingAndReservation(t *testing.T) {
	endpoint, err := url.Parse("https://timeout.invalid/")
	require.NoError(t, err)
	lc := CreateMockLC(*endpoint, false)
	gem := &globalEndpointManager{clientEndpoint: endpoint.String(), locationCache: lc}
	var observed retryContext
	transport := timeoutTestTransport(func(r *http.Request) (*http.Response, error) {
		return timeoutTestResponse(r, 200), nil
	})
	capture := timeoutTestStateCapture{state: &observed}
	pipeline := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{
		PerCall:  []policy.Policy{&retryCallerContextPolicy{}},
		PerRetry: []policy.Policy{&clientRetryPolicy{gem: gem}, capture},
	}, &policy.ClientOptions{Transport: transport, Retry: policy.RetryOptions{MaxRetries: -1}})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, endpoint.String())
	require.NoError(t, err)
	req.SetOperationValue(pipelineRequestOptions{})
	state := retryContext{retryCount: 2, sessionRetryCount: 1, preferredLocationIndex: 1,
		serverErrorRetryCount: 1, sameRegionRetryCount: 3, crossRegionFailoverDone: true,
		requestTimeoutRetryDone: true, resolveFromHead: true, retryNotBefore: time.Now().Add(30 * time.Millisecond)}
	req.SetOperationValue(interruptedReadRetry{state: state, err: context.DeadlineExceeded})
	start := time.Now()
	_, err = pipeline.Do(req)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, time.Since(start), 30*time.Millisecond)
	state.retryNotBefore = time.Time{}
	assert.Equal(t, state, observed)
	var consumed interruptedReadRetry
	require.True(t, req.OperationValue(&consumed))
	assert.Nil(t, consumed.err)
}

type timeoutTestStateCapture struct{ state *retryContext }

func (c timeoutTestStateCapture) Do(req *policy.Request) (*http.Response, error) {
	req.OperationValue(c.state)
	return req.Next()
}

type timeoutTestBody struct {
	ctx context.Context
	err error
}

func (b *timeoutTestBody) Read([]byte) (int, error) {
	<-b.ctx.Done()
	return 0, errors.Join(b.ctx.Err(), b.err)
}

func (*timeoutTestBody) Close() error { return nil }

func TestReadAttemptTimeoutRoutesReservedRetry(t *testing.T) {
	for _, mode := range []string{"connection", "not-sent", "408"} {
		t.Run(mode, func(t *testing.T) {
			a, err := url.Parse("https://region-a.invalid/")
			require.NoError(t, err)
			b, err := url.Parse("https://region-b.invalid/")
			require.NoError(t, err)
			lc := CreateMockLC(*a, false)
			lc.locationInfo.availReadLocations = []regionId{newRegionId("East US"), newRegionId("Central US")}
			lc.locationInfo.prefLocations = []regionId{newRegionId("East US"), newRegionId("Central US")}
			lc.locationInfo.availReadEndpointsByLocation = map[regionId]url.URL{newRegionId("East US"): *a, newRegionId("Central US"): *b}
			lc.locationInfo.readEndpoints = []url.URL{*a, *b}
			gem := &globalEndpointManager{clientEndpoint: a.String(), locationCache: lc}
			hosts := []string{}
			transport := timeoutTestTransport(func(r *http.Request) (*http.Response, error) {
				hosts = append(hosts, r.URL.Host)
				if mode == "408" && len(hosts) == 1 {
					return timeoutTestResponse(r, 408), nil
				}
				if mode == "connection" && len(hosts) <= 4 {
					return nil, io.EOF
				}
				if mode == "not-sent" && len(hosts) <= 4 {
					return nil, &net.DNSError{Err: "synthetic DNS failure", Name: r.URL.Host}
				}
				return timeoutTestResponse(r, 200), nil
			})
			counter := &timeoutTestAttemptCounter{}
			internal, err := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{
				PerCall:  []policy.Policy{&retryCallerContextPolicy{}},
				PerRetry: []policy.Policy{counter, &clientRetryPolicy{gem: gem}},
			}, &policy.ClientOptions{Transport: transport, Retry: policy.RetryOptions{
				MaxRetries: 10, TryTimeout: 100 * time.Millisecond,
				RetryDelay: 1100 * time.Millisecond, MaxRetryDelay: 1100 * time.Millisecond,
			}})
			require.NoError(t, err)
			client := &Client{endpoint: a.String(), endpointUrl: a, internal: internal, gem: gem}
			container, err := client.NewContainer("db", "coll")
			require.NoError(t, err)
			require.NoError(t, timeoutTestCall(context.Background(), container, "read"))
			expected := []string{a.Host, a.Host, a.Host, a.Host, b.Host}
			if mode == "408" {
				expected = []string{a.Host, b.Host}
			}
			assert.Equal(t, expected, hosts)
			assert.Equal(t, len(expected), counter.attempts)
			if mode == "not-sent" {
				assert.Len(t, lc.locationUnavailabilityInfoMap, 1)
			} else {
				assert.Empty(t, lc.locationUnavailabilityInfoMap, "ambiguous reads must not mark an endpoint")
			}
		})
	}
}
