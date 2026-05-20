// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

type attemptCounter struct {
	attempts int
}

func (a *attemptCounter) Do(req *policy.Request) (*http.Response, error) {
	a.attempts++
	return req.Next()
}

// throttleTestPipeline builds an azcore client wired with a specific
// throttleRetryPolicy and an attempt counter for inspection in tests.
// azcore's built-in retry policy is disabled so the throttleRetryPolicy is the
// only thing retrying. The counter is placed *after* the throttle policy so it
// gets invoked on every retry the throttle policy issues.
func throttleTestPipeline(t *testing.T, srv *mock.Server, p *throttleRetryPolicy) (*azcore.Client, *attemptCounter) {
	t.Helper()
	counter := &attemptCounter{}
	internal, err := azcore.NewClient("azcosmosthrottletest", "v1.0.0",
		azruntime.PipelineOptions{
			PerRetry: []policy.Policy{p, counter},
		},
		&policy.ClientOptions{
			Transport: srv,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		})
	require.NoError(t, err)
	return internal, counter
}

func doThrottleRequest(t *testing.T, c *azcore.Client, url string) (*http.Response, error) {
	t.Helper()
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, url)
	require.NoError(t, err)
	return c.Pipeline().Do(req)
}

func TestThrottleRetry_SucceedsAfterRetries(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	srv.AppendResponse(mock.WithStatusCode(429), mock.WithHeader(cosmosHeaderRetryAfterMs, "1"))
	srv.AppendResponse(mock.WithStatusCode(429), mock.WithHeader(cosmosHeaderRetryAfterMs, "1"))
	srv.AppendResponse(mock.WithStatusCode(200))

	client, counter := throttleTestPipeline(t, srv, &throttleRetryPolicy{
		maxRetryAttempts: 5,
		maxRetryWaitTime: 5 * time.Second,
		defaultDelay:     time.Millisecond,
	})

	resp, err := doThrottleRequest(t, client, srv.URL())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 3, counter.attempts)
}

func TestThrottleRetry_ExhaustsAttempts(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	srv.SetResponse(mock.WithStatusCode(429), mock.WithHeader(cosmosHeaderRetryAfterMs, "1"))

	client, counter := throttleTestPipeline(t, srv, &throttleRetryPolicy{
		maxRetryAttempts: 3,
		maxRetryWaitTime: 10 * time.Second,
		defaultDelay:     time.Millisecond,
	})

	resp, err := doThrottleRequest(t, client, srv.URL())
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	// 1 initial attempt + 3 retries
	require.Equal(t, 4, counter.attempts)
}

func TestThrottleRetry_ExhaustsCumulativeWaitTime(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	// Each 429 asks for a 60ms delay; budget is 100ms so only one retry fits.
	srv.SetResponse(mock.WithStatusCode(429), mock.WithHeader(cosmosHeaderRetryAfterMs, "60"))

	client, counter := throttleTestPipeline(t, srv, &throttleRetryPolicy{
		maxRetryAttempts: 100,
		maxRetryWaitTime: 100 * time.Millisecond,
		defaultDelay:     time.Millisecond,
	})

	resp, err := doThrottleRequest(t, client, srv.URL())
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	// 1 initial + 1 retry (cumulative 60ms used; second retry would push past 100ms)
	require.Equal(t, 2, counter.attempts)
}

func TestThrottleRetry_MissingHeaderUsesDefault(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	// 429 with no retry-after header followed by success. The policy should fall
	// back to its defaultDelay value.
	srv.AppendResponse(mock.WithStatusCode(429))
	srv.AppendResponse(mock.WithStatusCode(200))

	client, counter := throttleTestPipeline(t, srv, &throttleRetryPolicy{
		maxRetryAttempts: 5,
		maxRetryWaitTime: time.Second,
		defaultDelay:     5 * time.Millisecond,
	})

	resp, err := doThrottleRequest(t, client, srv.URL())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 2, counter.attempts)
}

func TestThrottleRetry_ExplicitZeroRetryAfterIsHonored(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	// Several 429s explicitly asking for "retry immediately" (header value "0"),
	// followed by success. If the policy treated an explicit 0 as "missing" and
	// fell back to defaultDelay (5s here) it would either exceed the 100ms
	// cumulative budget (so the first retry would be skipped and the test would
	// receive a 429) or take much longer than the assertion below allows.
	for i := 0; i < 4; i++ {
		srv.AppendResponse(mock.WithStatusCode(429), mock.WithHeader(cosmosHeaderRetryAfterMs, "0"))
	}
	srv.AppendResponse(mock.WithStatusCode(200))

	client, counter := throttleTestPipeline(t, srv, &throttleRetryPolicy{
		maxRetryAttempts: 10,
		maxRetryWaitTime: 100 * time.Millisecond,
		defaultDelay:     5 * time.Second,
	})

	start := time.Now()
	resp, err := doThrottleRequest(t, client, srv.URL())
	elapsed := time.Since(start)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 5, counter.attempts, "expected 4 retries + 1 success")
	require.Less(t, elapsed, time.Second, "explicit zero retry-after should not have waited the default delay")
}

func TestThrottleRetry_Non429PassesThrough(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	srv.SetResponse(mock.WithStatusCode(503))

	client, counter := throttleTestPipeline(t, srv, newThrottleRetryPolicy(&ThrottlingRetryOptions{
		MaxRetryAttempts: 5,
		MaxRetryWaitTime: 5 * time.Second,
	}))

	resp, err := doThrottleRequest(t, client, srv.URL())
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	require.Equal(t, 1, counter.attempts)
}

func TestThrottleRetry_ContextCancellationAbortsRetry(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	// Ask for a long retry-after so the policy is asleep when the context is cancelled.
	srv.SetResponse(mock.WithStatusCode(429), mock.WithHeader(cosmosHeaderRetryAfterMs, "5000"))

	client, counter := throttleTestPipeline(t, srv, &throttleRetryPolicy{
		maxRetryAttempts: 10,
		maxRetryWaitTime: time.Minute,
		defaultDelay:     time.Second,
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	req, err := azruntime.NewRequest(ctx, http.MethodGet, srv.URL())
	require.NoError(t, err)
	_, err = client.Pipeline().Do(req)
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
	// Exactly one transport attempt: the retry was aborted while sleeping.
	require.Equal(t, 1, counter.attempts)
}

func TestThrottleRetry_DefaultsAppliedWhenOptionsNil(t *testing.T) {
	p := newThrottleRetryPolicy(nil)
	require.Equal(t, defaultMaxThrottleRetryAttempts, p.maxRetryAttempts)
	require.Equal(t, defaultMaxThrottleRetryWaitTime, p.maxRetryWaitTime)
	require.Equal(t, defaultThrottleRetryDelay, p.defaultDelay)
}

func TestThrottleRetry_DefaultsAppliedWhenOptionsZero(t *testing.T) {
	p := newThrottleRetryPolicy(&ThrottlingRetryOptions{})
	require.Equal(t, defaultMaxThrottleRetryAttempts, p.maxRetryAttempts)
	require.Equal(t, defaultMaxThrottleRetryWaitTime, p.maxRetryWaitTime)
}

func TestThrottleRetry_NegativeAttemptsDisablesRetry(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	srv.SetResponse(mock.WithStatusCode(429), mock.WithHeader(cosmosHeaderRetryAfterMs, "1"))

	client, counter := throttleTestPipeline(t, srv,
		newThrottleRetryPolicy(&ThrottlingRetryOptions{MaxRetryAttempts: -1}))

	resp, err := doThrottleRequest(t, client, srv.URL())
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	require.Equal(t, 1, counter.attempts)
}

func TestReadRetryAfterMs(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   time.Duration
		wantOK bool
	}{
		{"missing", "", 0, false},
		{"integer", "1500", 1500 * time.Millisecond, true},
		{"float", "12.5", 12500 * time.Microsecond, true},
		{"explicit-zero", "0", 0, true},
		{"invalid", "not-a-number", 0, false},
		{"negative", "-10", 0, false},
		{"nan", "NaN", 0, false},
		{"positive-inf", "Inf", 0, false},
		{"negative-inf", "-Inf", 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{Header: http.Header{}}
			if tc.header != "" {
				resp.Header.Set(cosmosHeaderRetryAfterMs, tc.header)
			}
			got, ok := readRetryAfterMs(resp)
			require.Equal(t, tc.want, got)
			require.Equal(t, tc.wantOK, ok)
		})
	}
	got, ok := readRetryAfterMs(nil)
	require.Equal(t, time.Duration(0), got)
	require.False(t, ok)
}

// trackingBody is a strings.Reader-backed body that records how many times Seek(0,0)
// was called so tests can assert that the throttle policy rewinds the body across retries.
type trackingBody struct {
	*strings.Reader
	rewinds int
	closes  int
}

func newTrackingBody(s string) *trackingBody {
	return &trackingBody{Reader: strings.NewReader(s)}
}

func (b *trackingBody) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		b.rewinds++
	}
	return b.Reader.Seek(offset, whence)
}

func (b *trackingBody) Close() error {
	b.closes++
	return nil
}

// bodyEchoCounter is a per-retry policy that drains the request body and records what it
// saw on each transport attempt. It's used to prove that retries see the full body again.
type bodyEchoCounter struct {
	bodies [][]byte
}

func (b *bodyEchoCounter) Do(req *policy.Request) (*http.Response, error) {
	if raw := req.Raw(); raw != nil && raw.Body != nil {
		data, err := io.ReadAll(raw.Body)
		if err != nil {
			return nil, err
		}
		b.bodies = append(b.bodies, data)
		raw.Body = io.NopCloser(bytes.NewReader(data))
	} else {
		b.bodies = append(b.bodies, nil)
	}
	return req.Next()
}

func TestThrottleRetry_RewindsRequestBodyAcrossRetries(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	srv.AppendResponse(mock.WithStatusCode(429), mock.WithHeader(cosmosHeaderRetryAfterMs, "1"))
	srv.AppendResponse(mock.WithStatusCode(429), mock.WithHeader(cosmosHeaderRetryAfterMs, "1"))
	srv.AppendResponse(mock.WithStatusCode(200))

	echo := &bodyEchoCounter{}
	internal, err := azcore.NewClient("azcosmosthrottletest", "v1.0.0",
		azruntime.PipelineOptions{
			PerRetry: []policy.Policy{
				&throttleRetryPolicy{maxRetryAttempts: 5, maxRetryWaitTime: 5 * time.Second, defaultDelay: time.Millisecond},
				echo,
			},
		},
		&policy.ClientOptions{
			Transport: srv,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		})
	require.NoError(t, err)

	body := newTrackingBody(`{"id":"42"}`)
	req, err := azruntime.NewRequest(context.Background(), http.MethodPost, srv.URL())
	require.NoError(t, err)
	require.NoError(t, req.SetBody(body, "application/json"))

	resp, err := internal.Pipeline().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// All three transport attempts should have observed the full request body.
	require.Len(t, echo.bodies, 3)
	for i, b := range echo.bodies {
		require.Equal(t, `{"id":"42"}`, string(b), "transport attempt %d saw a truncated body", i)
	}
	// Two retries means the throttle policy rewound the body at least twice.
	require.GreaterOrEqual(t, body.rewinds, 2)
}

// fullPipelineClient wires up the same pipeline newClient builds (minus the
// globalEndpointManager bits) so we can verify the cosmos-level retry config
// in concert with azcore's retry policy.
func fullPipelineClient(t *testing.T, srv *mock.Server, opts *ClientOptions) (*azcore.Client, *attemptCounter) {
	t.Helper()
	if opts == nil {
		opts = &ClientOptions{}
	}
	clientOpts := opts.ClientOptions
	clientOpts.Transport = srv
	if clientOpts.Retry.RetryDelay == 0 {
		// Keep azcore's exponential backoff from making tests slow.
		clientOpts.Retry.RetryDelay = time.Millisecond
	}
	if clientOpts.Retry.StatusCodes == nil && clientOpts.Retry.ShouldRetry == nil {
		clientOpts.Retry.StatusCodes = defaultAzcoreRetryStatusCodesWithout429()
	}
	counter := &attemptCounter{}
	internal, err := azcore.NewClient("azcosmosthrottletest", "v1.0.0",
		azruntime.PipelineOptions{
			PerRetry: []policy.Policy{
				newThrottleRetryPolicy(&opts.ThrottlingRetryOptions),
				counter,
			},
		},
		&clientOpts)
	require.NoError(t, err)
	return internal, counter
}

// TestThrottleRetry_NoDoubleRetryWith429 verifies that, given the cosmos pipeline's
// default retry configuration, azcore's retry policy does not also retry 429s after
// the throttleRetryPolicy has exhausted its attempts. Without the StatusCodes override
// in newClient, azcore's default StatusCodes (which include 429) would retry the
// whole pipeline three additional times, multiplying the attempt count by 4.
func TestThrottleRetry_NoDoubleRetryWith429(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	srv.SetResponse(mock.WithStatusCode(429), mock.WithHeader(cosmosHeaderRetryAfterMs, "1"))

	client, counter := fullPipelineClient(t, srv, &ClientOptions{
		ThrottlingRetryOptions: ThrottlingRetryOptions{
			MaxRetryAttempts: 2,
			MaxRetryWaitTime: 10 * time.Second,
		},
	})

	resp, err := doThrottleRequest(t, client, srv.URL())
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	// Strictly 1 initial + 2 throttle retries. If azcore also retried 429s
	// the count would be 12 (3 * 4).
	require.Equal(t, 3, counter.attempts)
}

// TestThrottleRetry_AzcoreStillRetriesOther5xx ensures that excluding 429 from
// azcore's default retry StatusCodes leaves the other transient codes (e.g. 503)
// intact, so non-429 retries still happen.
func TestThrottleRetry_AzcoreStillRetriesOther5xx(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	srv.AppendResponse(mock.WithStatusCode(503))
	srv.AppendResponse(mock.WithStatusCode(503))
	srv.AppendResponse(mock.WithStatusCode(200))

	client, counter := fullPipelineClient(t, srv, nil)

	resp, err := doThrottleRequest(t, client, srv.URL())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	// azcore retried twice and then the third attempt succeeded.
	require.Equal(t, 3, counter.attempts)
}

// TestThrottleRetry_CallerStatusCodesPreserved verifies that we only override
// azcore's StatusCodes when the caller hasn't supplied their own. If they
// explicitly opt in to retrying 429 at the azcore layer, we respect that
// (even though it stacks with the throttle policy).
func TestThrottleRetry_CallerStatusCodesPreserved(t *testing.T) {
	srv, closeFn := mock.NewTLSServer()
	defer closeFn()

	srv.SetResponse(mock.WithStatusCode(429), mock.WithHeader(cosmosHeaderRetryAfterMs, "1"))

	caller429 := []int{http.StatusTooManyRequests}
	client, counter := fullPipelineClient(t, srv, &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Retry: policy.RetryOptions{
				MaxRetries:  1,
				StatusCodes: caller429,
				RetryDelay:  time.Millisecond,
			},
		},
		ThrottlingRetryOptions: ThrottlingRetryOptions{
			MaxRetryAttempts: 1,
			MaxRetryWaitTime: 10 * time.Second,
		},
	})

	resp, err := doThrottleRequest(t, client, srv.URL())
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	// Throttle: 1 initial + 1 retry = 2 per azcore try. azcore: 1 initial + 1 retry = 2 tries.
	// Total transport attempts = 2 * 2 = 4.
	require.Equal(t, 4, counter.attempts)
}

func TestDefaultAzcoreRetryStatusCodesWithout429(t *testing.T) {
	codes := defaultAzcoreRetryStatusCodesWithout429()
	for _, c := range codes {
		require.NotEqual(t, http.StatusTooManyRequests, c)
	}
	// Sanity check: the other transient HTTP failures azcore retries by default
	// are still in the list.
	for _, want := range []int{
		http.StatusRequestTimeout,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	} {
		require.Contains(t, codes, want)
	}
}
