// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
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
	}{
		{"missing", "", 0},
		{"integer", "1500", 1500 * time.Millisecond},
		{"float", "12.5", 12500 * time.Microsecond},
		{"invalid", "not-a-number", 0},
		{"negative", "-10", 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{Header: http.Header{}}
			if tc.header != "" {
				resp.Header.Set(cosmosHeaderRetryAfterMs, tc.header)
			}
			require.Equal(t, tc.want, readRetryAfterMs(resp))
		})
	}
	require.Equal(t, time.Duration(0), readRetryAfterMs(nil))
}
