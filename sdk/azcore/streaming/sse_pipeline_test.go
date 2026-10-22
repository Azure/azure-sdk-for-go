// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package streaming_test

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sseTransport func(*http.Request) (*http.Response, error)

func (t sseTransport) Do(request *http.Request) (*http.Response, error) { return t(request) }

func sendSSERequest(ctx context.Context, pipeline runtime.Pipeline, endpoint, lastEventID string) (io.ReadCloser, error) {
	request, err := runtime.NewRequest(ctx, http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}
	if lastEventID != "" {
		request.Raw().Header.Set("Last-Event-ID", lastEventID)
	}
	runtime.SkipBodyDownload(request)
	response, err := pipeline.Do(request)
	return runtime.SSEResponse(response, err, http.StatusOK)
}

func TestSSEReconnectRespectsPipelineRetryPolicy(t *testing.T) {
	for _, tc := range []struct {
		name        string
		override    bool
		veto        bool
		disabled    bool
		wantRetries int
	}{
		{name: "client retry limit", wantRetries: 2},
		{name: "context retry override", override: true, wantRetries: 1},
		{name: "ShouldRetry veto", veto: true},
		{name: "pipeline retries disabled", disabled: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requests := 0
			transport := sseTransport(func(request *http.Request) (*http.Response, error) {
				requests++
				status := http.StatusInternalServerError
				body := "{}"
				if requests == 1 {
					assert.NotContains(t, request.Header, "Last-Event-Id")
					status = http.StatusOK
					body = "retry: 0\nid: first\ndata: one\n\n"
				} else {
					assert.Equal(t, "first", request.Header.Get("Last-Event-ID"))
				}
				return &http.Response{
					StatusCode: status,
					Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
					Body:       eventBody(body),
					Request:    request,
				}, nil
			})
			retry := policy.RetryOptions{MaxRetries: 2, RetryDelay: time.Nanosecond, MaxRetryDelay: time.Microsecond}
			if tc.disabled {
				retry.MaxRetries = -1
			}
			if tc.veto {
				retry.ShouldRetry = func(*http.Response, error) bool { return false }
			}
			pipeline := runtime.NewPipeline("sse-test", "v0.0.0", runtime.PipelineOptions{}, &policy.ClientOptions{
				Retry:     retry,
				Transport: transport,
			})
			ctx := context.Background()
			if tc.override {
				ctx = policy.WithRetryOptions(ctx, policy.RetryOptions{
					MaxRetries:    1,
					RetryDelay:    time.Nanosecond,
					MaxRetryDelay: time.Microsecond,
				})
			}
			connects := 0
			reader, err := streaming.NewSSEReader(ctx, func(ctx context.Context, lastEventID string) (io.ReadCloser, error) {
				connects++
				return sendSSERequest(ctx, pipeline, "https://example.com/events", lastEventID)
			}, streaming.SSEHandler[frameView]{
				Decode:    decodeView,
				Reconnect: true,
			}, nil)
			require.NoError(t, err)
			defer func() {
				_ = reader.Close()
			}()
			assert.Equal(t, 1, connects)
			_, err = reader.Next()
			require.NoError(t, err)
			_, err = reader.Next()
			require.Error(t, err)
			assert.Equal(t, 2, connects)
			assert.Equal(t, 2+tc.wantRetries, requests)
			_, err = reader.Next()
			require.ErrorIs(t, err, io.EOF)
			assert.Equal(t, 2, connects)
		})
	}
}

func TestSSEReconnectUsesOperationContextInsteadOfAttemptContext(t *testing.T) {
	requests := 0
	var attemptCtx context.Context
	transport := sseTransport(func(request *http.Request) (*http.Response, error) {
		requests++
		require.NoError(t, request.Context().Err())
		body := "retry: 0\nid: first\ndata: one\n\n"
		if requests == 1 {
			attemptCtx = request.Context()
		}
		if requests > 1 {
			body = "data: two\n\ndata: [DONE]\n\n"
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
			Body:       eventBody(body),
			Request:    request,
		}, nil
	})
	pipeline := runtime.NewPipeline("sse-test", "v0.0.0", runtime.PipelineOptions{}, &policy.ClientOptions{
		Retry:     policy.RetryOptions{TryTimeout: time.Minute},
		Transport: transport,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var readerCtx context.Context
	reader, err := streaming.NewSSEReader(ctx, func(connectCtx context.Context, lastEventID string) (io.ReadCloser, error) {
		if attemptCtx != nil {
			assert.ErrorIs(t, attemptCtx.Err(), context.Canceled)
			assert.Same(t, readerCtx, connectCtx)
		} else {
			readerCtx = connectCtx
		}
		require.NoError(t, connectCtx.Err())
		require.NoError(t, ctx.Err())
		return sendSSERequest(connectCtx, pipeline, "https://example.com/events", lastEventID)
	}, streaming.SSEHandler[frameView]{
		Decode:    decodeView,
		Reconnect: true,
	}, nil)
	require.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()
	first, err := reader.Next()
	require.NoError(t, err)
	second, err := reader.Next()
	require.NoError(t, err)
	assert.Equal(t, "one", first.data)
	assert.Equal(t, "two", second.data)
	assert.Equal(t, "first", reader.LastEventID())
	_, err = reader.Next()
	require.ErrorIs(t, err, io.EOF)
	assert.Equal(t, 2, requests)
}

func TestSSEReconnectCompletionAndInitialCheckpoint(t *testing.T) {
	for _, lastID := range []string{"", "prior"} {
		t.Run("checkpoint="+lastID, func(t *testing.T) {
			requests := 0
			completionBody := &observedSSEBody{ReadCloser: eventBody("")}
			transport := sseTransport(func(request *http.Request) (*http.Response, error) {
				requests++
				assert.Equal(t, lastID, request.Header.Get("Last-Event-ID"))
				if lastID == "" {
					assert.NotContains(t, request.Header, "Last-Event-Id")
				}
				response := &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"text/event-stream; charset=utf-8"}},
					Body:       eventBody("retry: 0\ndata: one\n\n"),
					Request:    request,
				}
				if requests > 1 {
					response.StatusCode = http.StatusNoContent
					response.Body = completionBody
				}
				return response, nil
			})
			pipeline := runtime.NewPipeline("sse-test", "v0.0.0", runtime.PipelineOptions{}, &policy.ClientOptions{Transport: transport})
			ctx := context.Background()
			reader := newSSEReaderWithContext(t, ctx, func(ctx context.Context, lastEventID string) (io.ReadCloser, error) {
				return sendSSERequest(ctx, pipeline, "https://example.com/events", lastEventID)
			}, streaming.SSEHandler[frameView]{Reconnect: true}, &streaming.SSEReaderOptions{LastEventID: lastID})
			assert.Equal(t, 1, requests)
			event, err := reader.Next()
			require.NoError(t, err)
			assert.Equal(t, "one", event.data)
			assert.Equal(t, lastID, reader.LastEventID())
			_, err = reader.Next()
			require.ErrorIs(t, err, io.EOF)
			_, err = reader.Next()
			require.ErrorIs(t, err, io.EOF)
			assert.Equal(t, 2, requests)
			assert.EqualValues(t, 1, completionBody.closed.Load())
		})
	}
}

func TestSSEInitialConnectionRespectsPipelineRetries(t *testing.T) {
	for _, fail := range []bool{false, true} {
		name := "retry succeeds"
		if fail {
			name = "retries exhausted"
		}
		t.Run(name, func(t *testing.T) {
			requests, connects := 0, 0
			transport := sseTransport(func(request *http.Request) (*http.Response, error) {
				requests++
				assert.Equal(t, "prior", request.Header.Get("Last-Event-ID"))
				status := http.StatusServiceUnavailable
				body := "{}"
				if requests == 2 && !fail {
					status = http.StatusOK
					body = "data: one\n\n"
				}
				return &http.Response{
					StatusCode: status,
					Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
					Body:       eventBody(body),
					Request:    request,
				}, nil
			})
			pipeline := runtime.NewPipeline("sse-test", "v0.0.0", runtime.PipelineOptions{}, &policy.ClientOptions{
				Retry:     policy.RetryOptions{MaxRetries: 1, RetryDelay: time.Nanosecond, MaxRetryDelay: time.Microsecond},
				Transport: transport,
			})
			reader, err := streaming.NewSSEReader(context.Background(), func(ctx context.Context, lastEventID string) (io.ReadCloser, error) {
				connects++
				return sendSSERequest(ctx, pipeline, "https://example.com/events", lastEventID)
			}, streaming.SSEHandler[frameView]{Decode: decodeView, Reconnect: true}, &streaming.SSEReaderOptions{LastEventID: "prior"})
			if fail {
				require.Error(t, err)
				assert.Nil(t, reader)
			} else {
				require.NoError(t, err)
				t.Cleanup(func() { assert.NoError(t, reader.Close()) })
				event, err := reader.Next()
				require.NoError(t, err)
				assert.Equal(t, "one", event.data)
			}
			assert.Equal(t, 1, connects)
			assert.Equal(t, 2, requests)
		})
	}
}
