// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/stretchr/testify/require"
)

// newPipelineForPolicy wraps the given policy in a pipeline targeting the supplied mock server.
// Retries are disabled so each pipeline invocation produces exactly one transport call, making
// expectation checks against the mock server's response sequence deterministic.
func newPipelineForPolicy(p policy.Policy, srv policy.Transporter) runtime.Pipeline {
	return runtime.NewPipeline("test", "v0.0.0",
		runtime.PipelineOptions{PerRetry: []policy.Policy{p}},
		&policy.ClientOptions{
			Transport: srv,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		},
	)
}

// newRequestWithBody constructs a *policy.Request with the supplied body bytes attached and
// the Content-Length populated. When body is nil the request has no body.
func newRequestWithBody(t *testing.T, body []byte) *policy.Request {
	t.Helper()
	req, err := runtime.NewRequest(context.Background(), http.MethodPut, "https://localhost")
	require.NoError(t, err)
	if body != nil {
		require.NoError(t, req.SetBody(streaming.NopCloser(bytes.NewReader(body)), "application/octet-stream"))
	}
	return req
}

// TestExpectContinuePolicyAddsHeaderOnContentBody verifies the Expect: 100-continue header is
// set only when a non-empty body is present.
func TestExpectContinuePolicyAddsHeaderOnContentBody(t *testing.T) {
	for _, hasBody := range []bool{true, false} {
		hasBody := hasBody
		t.Run("", func(t *testing.T) {
			srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeSrv()
			srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

			p := NewExpectContinuePolicy(&exported.ExpectContinueOptions{Mode: exported.ExpectContinueModeOn})
			require.NotNil(t, p)
			pl := newPipelineForPolicy(p, srv)

			var body []byte
			if hasBody {
				body = []byte("foo")
			}
			req := newRequestWithBody(t, body)
			_, err := pl.Do(req)
			require.NoError(t, err)

			if hasBody {
				require.Equal(t, "100-continue", req.Raw().Header.Get("Expect"))
			} else {
				require.Empty(t, req.Raw().Header.Get("Expect"))
			}
		})
	}
}

// TestExpectContinuePolicyRespectsThreshold verifies the header is applied only when the
// request's Content-Length meets or exceeds ContentLengthThreshold.
func TestExpectContinuePolicyRespectsThreshold(t *testing.T) {
	cases := []struct {
		threshold    int64
		bodyLength   int
		expectHeader bool
	}{
		{threshold: 1024, bodyLength: 2048, expectHeader: true},
		{threshold: 2048, bodyLength: 1024, expectHeader: false},
		{threshold: 1024, bodyLength: 1024, expectHeader: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run("", func(t *testing.T) {
			srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeSrv()
			srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

			p := NewExpectContinuePolicy(&exported.ExpectContinueOptions{
				Mode:                   exported.ExpectContinueModeOn,
				ContentLengthThreshold: to.Ptr(tc.threshold),
			})
			require.NotNil(t, p)
			pl := newPipelineForPolicy(p, srv)

			req := newRequestWithBody(t, bytes.Repeat([]byte("a"), tc.bodyLength))
			_, err := pl.Do(req)
			require.NoError(t, err)

			if tc.expectHeader {
				require.Equal(t, "100-continue", req.Raw().Header.Get("Expect"))
			} else {
				require.Empty(t, req.Raw().Header.Get("Expect"))
			}
		})
	}
}

// TestExpectContinueOnThrottlePolicyAddsHeaderOnlyAfterError verifies that, in
// ExpectContinueModeApplyOnThrottle mode, a 429/500/503 response opens a window during which
// subsequent requests carry the header.
func TestExpectContinueOnThrottlePolicyAddsHeaderOnlyAfterError(t *testing.T) {
	for _, status := range []int{http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusServiceUnavailable} {
		status := status
		t.Run("", func(t *testing.T) {
			srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeSrv()
			srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
			srv.AppendResponse(mock.WithStatusCode(status))
			srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

			p := NewExpectContinuePolicy(&exported.ExpectContinueOptions{Mode: exported.ExpectContinueModeApplyOnThrottle})
			require.NotNil(t, p)
			pl := newPipelineForPolicy(p, srv)

			// message 1 doesn't get expect header (no prior throttle)
			req1 := newRequestWithBody(t, []byte("foo"))
			_, err := pl.Do(req1)
			require.NoError(t, err)
			require.Empty(t, req1.Raw().Header.Get("Expect"))

			// message 2 doesn't get expect header but triggers future messages
			req2 := newRequestWithBody(t, []byte("foo"))
			_, err = pl.Do(req2)
			require.NoError(t, err)
			require.Empty(t, req2.Raw().Header.Get("Expect"))

			// message 3 gets expect header
			req3 := newRequestWithBody(t, []byte("foo"))
			_, err = pl.Do(req3)
			require.NoError(t, err)
			require.Equal(t, "100-continue", req3.Raw().Header.Get("Expect"))
		})
	}
}

// TestExpectContinueOnThrottlePolicyRespectsThreshold verifies that, in ApplyOnThrottle mode,
// the header is applied only when the request's Content-Length meets or exceeds
// ContentLengthThreshold.
func TestExpectContinueOnThrottlePolicyRespectsThreshold(t *testing.T) {
	cases := []struct {
		threshold    int64
		bodyLength   int
		expectHeader bool
	}{
		{threshold: 1024, bodyLength: 2048, expectHeader: true},
		{threshold: 2048, bodyLength: 1024, expectHeader: false},
		{threshold: 1024, bodyLength: 1024, expectHeader: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run("", func(t *testing.T) {
			srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeSrv()
			srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
			srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

			p := NewExpectContinuePolicy(&exported.ExpectContinueOptions{
				Mode:                   exported.ExpectContinueModeApplyOnThrottle,
				ContentLengthThreshold: to.Ptr(tc.threshold),
			})
			require.NotNil(t, p)
			pl := newPipelineForPolicy(p, srv)

			// message 1 doesn't get expect header but triggers future messages
			req1 := newRequestWithBody(t, bytes.Repeat([]byte("a"), tc.bodyLength))
			_, err := pl.Do(req1)
			require.NoError(t, err)
			require.Empty(t, req1.Raw().Header.Get("Expect"))

			// message 2 may or may not get expect header depending on threshold
			req2 := newRequestWithBody(t, bytes.Repeat([]byte("a"), tc.bodyLength))
			_, err = pl.Do(req2)
			require.NoError(t, err)

			if tc.expectHeader {
				require.Equal(t, "100-continue", req2.Raw().Header.Get("Expect"))
			} else {
				require.Empty(t, req2.Raw().Header.Get("Expect"))
			}
		})
	}
}

// TestExpectContinueOnThrottlePolicyRevertsAfterBackoff verifies that, after the throttle window
// elapses, the header is no longer applied. The throttle interval is overridden on the policy
// struct so the test runs fast.
func TestExpectContinueOnThrottlePolicyRevertsAfterBackoff(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

	backoff := 10 * time.Millisecond
	throttlePolicy := &expectContinueOnThrottlePolicy{throttleInterval: backoff}
	pl := newPipelineForPolicy(throttlePolicy, srv)

	// message 1 doesn't get expect header but triggers future messages
	req1 := newRequestWithBody(t, []byte("foo"))
	_, err := pl.Do(req1)
	require.NoError(t, err)
	require.Empty(t, req1.Raw().Header.Get("Expect"))

	// message 2 gets expect header
	req2 := newRequestWithBody(t, []byte("foo"))
	_, err = pl.Do(req2)
	require.NoError(t, err)
	require.Equal(t, "100-continue", req2.Raw().Header.Get("Expect"))

	// wait out the throttle window
	time.Sleep(2 * backoff)

	// message 3 no longer gets expect header
	req3 := newRequestWithBody(t, []byte("foo"))
	_, err = pl.Do(req3)
	require.NoError(t, err)
	require.Empty(t, req3.Raw().Header.Get("Expect"))
}

// TestNewExpectContinuePolicyOffReturnsNil verifies that mode Off causes no policy to be added.
func TestNewExpectContinuePolicyOffReturnsNil(t *testing.T) {
	p := NewExpectContinuePolicy(&exported.ExpectContinueOptions{Mode: exported.ExpectContinueModeOff})
	require.Nil(t, p)
}

// TestNewExpectContinuePolicyDefaultsToApplyOnThrottle verifies that nil options produces the
// ApplyOnThrottle policy.
func TestNewExpectContinuePolicyDefaultsToApplyOnThrottle(t *testing.T) {
	p := NewExpectContinuePolicy(nil)
	require.NotNil(t, p)
	_, ok := p.(*expectContinueOnThrottlePolicy)
	require.True(t, ok, "expected *expectContinueOnThrottlePolicy, got %T", p)
}

// TestNewExpectContinuePolicyEnvVarDisables verifies that setting the disable env var causes
// no policy to be returned, regardless of supplied options.
func TestNewExpectContinuePolicyEnvVarDisables(t *testing.T) {
	for _, v := range []string{"1", "true", "True", "TRUE"} {
		v := v
		t.Run(v, func(t *testing.T) {
			t.Setenv(DisableExpectContinueHeaderEnvVar, v)
			resetExpectContinueEnvCacheForTest(t)
			require.Nil(t, NewExpectContinuePolicy(nil))
			require.Nil(t, NewExpectContinuePolicy(&exported.ExpectContinueOptions{Mode: exported.ExpectContinueModeOn}))
			require.Nil(t, NewExpectContinuePolicy(&exported.ExpectContinueOptions{Mode: exported.ExpectContinueModeApplyOnThrottle}))
		})
	}
}

// TestNewExpectContinuePolicyEnvVarFalsyValuesEnable verifies that the policy is added when the
// env var is unset or set to a falsy value.
func TestNewExpectContinuePolicyEnvVarFalsyValuesEnable(t *testing.T) {
	for _, v := range []string{"0", "false", "False", "not-a-bool"} {
		v := v
		t.Run(v, func(t *testing.T) {
			t.Setenv(DisableExpectContinueHeaderEnvVar, v)
			resetExpectContinueEnvCacheForTest(t)
			require.NotNil(t, NewExpectContinuePolicy(nil))
		})
	}
}

// resetExpectContinueEnvCacheForTest discards the memoized env-var lookup so a subsequent call
// to NewExpectContinuePolicy re-reads the current environment. The previous value is restored
// when the test (or subtest) completes.
func resetExpectContinueEnvCacheForTest(t *testing.T) {
	t.Helper()
	prev := envCheckExpectContinueDisabled
	envCheckExpectContinueDisabled = newExpectContinueEnvCheck()
	t.Cleanup(func() { envCheckExpectContinueDisabled = prev })
}

// TestExpectContinueEnvCacheIsMemoized verifies the env-var lookup is performed exactly once
// per cache lifetime: a change to the env variable after the first read is not observed by
// subsequent calls until the cache is reset.
func TestExpectContinueEnvCacheIsMemoized(t *testing.T) {
	resetExpectContinueEnvCacheForTest(t)

	t.Setenv(DisableExpectContinueHeaderEnvVar, "true")
	// Prime the cache while the variable is "true" - the policy must be disabled.
	require.Nil(t, NewExpectContinuePolicy(nil))

	// Change the env variable. The cache should not pick up the new value.
	t.Setenv(DisableExpectContinueHeaderEnvVar, "false")
	require.Nil(t, NewExpectContinuePolicy(nil), "cache should still report disabled until reset")

	// Resetting the cache forces a re-read.
	resetExpectContinueEnvCacheForTest(t)
	require.NotNil(t, NewExpectContinuePolicy(nil))
}

// TestExpectContinuePolicyIgnoresUnknownContentLength verifies the header is not added when
// content length is unknown (e.g. -1 from chunked encoding).
func TestExpectContinuePolicyIgnoresUnknownContentLength(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

	p := NewExpectContinuePolicy(&exported.ExpectContinueOptions{Mode: exported.ExpectContinueModeOn})
	require.NotNil(t, p)
	pl := newPipelineForPolicy(p, srv)

	req := newRequestWithBody(t, []byte("foo"))
	// Simulate an unknown content length (chunked encoding).
	req.Raw().ContentLength = -1
	_, err := pl.Do(req)
	require.NoError(t, err)
	require.Empty(t, req.Raw().Header.Get("Expect"))
}

// TestNewExpectContinuePolicyDefaultThrottleIntervalIsOneMinute verifies that the default
// throttle interval is one minute.
func TestNewExpectContinuePolicyDefaultThrottleIntervalIsOneMinute(t *testing.T) {
	p := NewExpectContinuePolicy(nil)
	require.NotNil(t, p)
	tp, ok := p.(*expectContinueOnThrottlePolicy)
	require.True(t, ok, "expected *expectContinueOnThrottlePolicy, got %T", p)
	require.Equal(t, time.Minute, tp.throttleInterval)
}

// TestNewExpectContinuePolicyHonorsAutoInterval verifies that a caller-supplied AutoInterval
// overrides the default throttle interval.
func TestNewExpectContinuePolicyHonorsAutoInterval(t *testing.T) {
	custom := 250 * time.Millisecond
	p := NewExpectContinuePolicy(&exported.ExpectContinueOptions{ThrottleInterval: &custom})
	require.NotNil(t, p)
	tp, ok := p.(*expectContinueOnThrottlePolicy)
	require.True(t, ok, "expected *expectContinueOnThrottlePolicy, got %T", p)
	require.Equal(t, custom, tp.throttleInterval)
}

// TestExpectContinueOnThrottlePolicyAutoIntervalEndToEnd verifies that the user-supplied
// AutoInterval is observed end-to-end through the factory and the pipeline: the header is set
// while the configured window is open and removed after it elapses.
func TestExpectContinueOnThrottlePolicyAutoIntervalEndToEnd(t *testing.T) {
	srv, closeSrv := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer closeSrv()
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

	backoff := 10 * time.Millisecond
	p := NewExpectContinuePolicy(&exported.ExpectContinueOptions{ThrottleInterval: &backoff})
	require.NotNil(t, p)
	pl := newPipelineForPolicy(p, srv)

	// message 1 doesn't get expect header but triggers future messages
	req1 := newRequestWithBody(t, []byte("foo"))
	_, err := pl.Do(req1)
	require.NoError(t, err)
	require.Empty(t, req1.Raw().Header.Get("Expect"))

	// message 2 gets expect header
	req2 := newRequestWithBody(t, []byte("foo"))
	_, err = pl.Do(req2)
	require.NoError(t, err)
	require.Equal(t, "100-continue", req2.Raw().Header.Get("Expect"))

	// wait out the user-supplied window
	time.Sleep(2 * backoff)

	// message 3 no longer gets expect header
	req3 := newRequestWithBody(t, []byte("foo"))
	_, err = pl.Do(req3)
	require.NoError(t, err)
	require.Empty(t, req3.Raw().Header.Get("Expect"))
}
