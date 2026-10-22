// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package streaming_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"testing/iotest"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
)

type observedSSEHTTPBody struct {
	io.ReadCloser
	reading   chan struct{}
	closed    chan struct{}
	readOnce  sync.Once
	closeOnce sync.Once
	closeErr  error
}

func (b *observedSSEHTTPBody) Read(p []byte) (int, error) {
	if b.reading != nil {
		b.readOnce.Do(func() { close(b.reading) })
	}
	return b.ReadCloser.Read(p)
}

func (b *observedSSEHTTPBody) Close() error {
	err := b.ReadCloser.Close()
	b.closeOnce.Do(func() { close(b.closed) })
	return errors.Join(err, b.closeErr)
}

func sseHTTPConnector(endpoint string, options *policy.ClientOptions) streaming.SSEConnector {
	pipeline := runtime.NewPipeline("sse-test", "v0.0.0", runtime.PipelineOptions{}, options)
	return func(ctx context.Context, lastEventID string) (io.ReadCloser, error) {
		return sendSSERequest(ctx, pipeline, endpoint, lastEventID)
	}
}

func waitForSSEHTTPSignal(t *testing.T, signal <-chan struct{}) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for SSE lifecycle signal")
	}
}

func TestSSEHTTPInitialResponseValidation(t *testing.T) {
	for _, tc := range []struct {
		name        string
		status      int
		contentType string
		wantErr     string
	}{
		{name: "no content", status: http.StatusNoContent},
		{name: "wrong content type", status: http.StatusOK, contentType: "application/json", wantErr: "unexpected SSE Content-Type"},
		{name: "malformed content type", status: http.StatusOK, contentType: "text/event-stream; charset", wantErr: "invalid SSE Content-Type"},
		{name: "HTTP error", status: http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				if tc.contentType != "" {
					w.Header().Set("Content-Type", tc.contentType)
				}
				w.WriteHeader(tc.status)
				if tc.status != http.StatusNoContent {
					_, _ = fmt.Fprint(w, "{}")
				}
			}))
			t.Cleanup(server.Close)
			bodies := make(chan *observedSSEHTTPBody, 1)
			transport := sseTransport(func(req *http.Request) (*http.Response, error) {
				resp, err := http.DefaultClient.Do(req)
				if err == nil {
					body := &observedSSEHTTPBody{ReadCloser: resp.Body, closed: make(chan struct{})}
					resp.Body = body
					bodies <- body
				}
				return resp, err
			})
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			t.Cleanup(cancel)
			connect := sseHTTPConnector(server.URL, &policy.ClientOptions{
				Transport: transport, Retry: policy.RetryOptions{MaxRetries: -1},
			})
			reader, err := streaming.NewSSEReader(ctx, connect, streaming.SSEHandler[frameView]{Decode: decodeView}, nil)
			if tc.status == http.StatusNoContent {
				require.NoError(t, err)
				require.NotNil(t, reader)
				t.Cleanup(func() { require.NoError(t, reader.Close()) })
				for range 2 {
					_, err = reader.Next()
					require.ErrorIs(t, err, io.EOF)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, reader)
				if tc.wantErr != "" {
					require.ErrorContains(t, err, tc.wantErr)
				} else {
					var responseErr *azcore.ResponseError
					require.ErrorAs(t, err, &responseErr)
					require.Equal(t, tc.status, responseErr.StatusCode)
				}
			}
			waitForSSEHTTPSignal(t, (<-bodies).closed)
			require.NoError(t, ctx.Err(), "closing the reader's child context must not cancel the operation context")
			require.EqualValues(t, 1, requests.Load())
		})
	}
}

func TestSSEHTTPIdleReadCanBeInterrupted(t *testing.T) {
	for _, mode := range []string{"close", "cancel", "deadline"} {
		t.Run(mode, func(t *testing.T) {
			var requests atomic.Int32
			disconnected := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				w.Header().Set("Content-Type", "text/event-stream")
				_, _ = fmt.Fprint(w, ": heartbeat\n\n")
				w.(http.Flusher).Flush()
				<-r.Context().Done()
				close(disconnected)
			}))
			t.Cleanup(server.Close)
			reading := make(chan struct{})
			transport := sseTransport(func(req *http.Request) (*http.Response, error) {
				resp, err := http.DefaultClient.Do(req)
				if err == nil {
					resp.Body = &observedSSEHTTPBody{ReadCloser: resp.Body, reading: reading, closed: make(chan struct{})}
				}
				return resp, err
			})
			timeout := 5 * time.Second
			if mode == "deadline" {
				timeout = time.Second
			}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			t.Cleanup(cancel)
			connect := sseHTTPConnector(server.URL, &policy.ClientOptions{
				Transport: transport, Retry: policy.RetryOptions{TryTimeout: time.Minute},
			})
			reader, err := streaming.NewSSEReader(ctx, connect, streaming.SSEHandler[frameView]{Decode: decodeView}, nil)
			require.NoError(t, err)
			require.NotNil(t, reader)
			t.Cleanup(func() { require.NoError(t, reader.Close()) })
			result := make(chan error, 1)
			go func() { _, err := reader.Next(); result <- err }()
			waitForSSEHTTPSignal(t, reading)
			wantErr := error(io.EOF)
			switch mode {
			case "close":
				require.NoError(t, reader.Close())
				require.NoError(t, ctx.Err())
			case "cancel":
				cancel()
				wantErr = context.Canceled
			case "deadline":
				wantErr = context.DeadlineExceeded
			}
			select {
			case err := <-result:
				require.ErrorIs(t, err, wantErr)
			case <-time.After(5 * time.Second):
				t.Fatal("SSE read did not stop")
			}
			waitForSSEHTTPSignal(t, disconnected)
			_, err = reader.Next()
			require.ErrorIs(t, err, io.EOF)
			require.EqualValues(t, 1, requests.Load())
		})
	}
}

func TestSSEHTTPCancellationWhileOpening(t *testing.T) {
	started := make(chan struct{})
	disconnected := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
		close(disconnected)
	}))
	t.Cleanup(server.Close)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	connect := sseHTTPConnector(server.URL, nil)
	type openResult struct {
		reader *streaming.SSEReader[frameView]
		err    error
	}
	result := make(chan openResult, 1)
	go func() {
		reader, err := streaming.NewSSEReader(ctx, connect, streaming.SSEHandler[frameView]{Decode: decodeView}, nil)
		result <- openResult{reader: reader, err: err}
	}()
	waitForSSEHTTPSignal(t, started)
	cancel()
	select {
	case got := <-result:
		require.ErrorIs(t, got.err, context.Canceled)
		require.Nil(t, got.reader)
	case <-time.After(5 * time.Second):
		t.Fatal("opening the SSE stream did not stop")
	}
	waitForSSEHTTPSignal(t, disconnected)
}

func TestSSEHTTPInitialConnectionPipelineRetries(t *testing.T) {
	for _, fail := range []bool{false, true} {
		name := "retry succeeds"
		if fail {
			name = "retries exhausted"
		}
		t.Run(name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if requests.Add(1) == 1 || fail {
					http.Error(w, "unavailable", http.StatusServiceUnavailable)
					return
				}
				w.Header().Set("Content-Type", "text/event-stream")
				_, _ = fmt.Fprint(w, "event: message\ndata: {\"message\":\"hello\"}\n\n")
			}))
			t.Cleanup(server.Close)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			t.Cleanup(cancel)
			connect := sseHTTPConnector(server.URL, &policy.ClientOptions{
				Retry: policy.RetryOptions{MaxRetries: 1, RetryDelay: time.Nanosecond, MaxRetryDelay: time.Microsecond},
			})
			reader, err := streaming.NewSSEReader(ctx, connect, streaming.SSEHandler[frameView]{Decode: decodeView}, nil)
			if fail {
				var responseErr *azcore.ResponseError
				require.ErrorAs(t, err, &responseErr)
				require.Equal(t, http.StatusServiceUnavailable, responseErr.StatusCode)
				require.Nil(t, reader)
			} else {
				require.NoError(t, err)
				require.NotNil(t, reader)
				t.Cleanup(func() { require.NoError(t, reader.Close()) })
				event, err := reader.Next()
				require.NoError(t, err)
				require.Equal(t, "hello", event.fields["message"])
				_, err = reader.Next()
				require.ErrorIs(t, err, io.EOF)
			}
			require.EqualValues(t, 2, requests.Load())
		})
	}
}

func TestSSEHTTPInitialTransportEOFIsNotCompletion(t *testing.T) {
	var requests atomic.Int32
	connect := sseHTTPConnector("https://fake.local", &policy.ClientOptions{
		Retry: policy.RetryOptions{MaxRetries: -1},
		Transport: sseTransport(func(*http.Request) (*http.Response, error) {
			requests.Add(1)
			return nil, io.EOF
		}),
	})
	reader, err := streaming.NewSSEReader(context.Background(), connect, streaming.SSEHandler[frameView]{Decode: decodeView}, nil)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.Nil(t, reader)
	require.EqualValues(t, 1, requests.Load())
}

func TestSSEHTTPInitialResponseIgnoresCleanupEOF(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNoContent} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			body := &observedSSEHTTPBody{ReadCloser: io.NopCloser(http.NoBody), closed: make(chan struct{}), closeErr: io.EOF}
			connect := sseHTTPConnector("https://fake.local", &policy.ClientOptions{
				Transport: sseTransport(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: status,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       body,
						Request:    req,
					}, nil
				}),
			})
			reader, err := streaming.NewSSEReader(context.Background(), connect, streaming.SSEHandler[frameView]{Decode: decodeView}, nil)
			if status == http.StatusOK {
				require.Error(t, err)
				require.NotErrorIs(t, err, io.EOF)
				require.Nil(t, reader)
				require.ErrorContains(t, err, "unexpected SSE Content-Type")
			} else {
				require.NoError(t, err)
				require.NotNil(t, reader)
				t.Cleanup(func() { require.NoError(t, reader.Close()) })
				_, err = reader.Next()
				require.ErrorIs(t, err, io.EOF)
			}
			waitForSSEHTTPSignal(t, body.closed)
		})
	}
}

func TestSSEHTTPRejectedResponseEOFIsNotCompletion(t *testing.T) {
	readErr := fmt.Errorf("reading response: %w", io.EOF)
	body := &observedSSEHTTPBody{
		ReadCloser: io.NopCloser(iotest.ErrReader(readErr)),
		closed:     make(chan struct{}),
	}
	var requests atomic.Int32
	connect := sseHTTPConnector("https://fake.local", &policy.ClientOptions{
		Retry: policy.RetryOptions{MaxRetries: -1},
		Transport: sseTransport(func(req *http.Request) (*http.Response, error) {
			requests.Add(1)
			return &http.Response{
				StatusCode: http.StatusCreated,
				Header:     make(http.Header),
				Body:       body,
				Request:    req,
			}, nil
		}),
	})
	reader, err := streaming.NewSSEReader(context.Background(), connect, streaming.SSEHandler[frameView]{Decode: decodeView}, nil)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.NotErrorIs(t, err, io.EOF)
	require.ErrorContains(t, err, readErr.Error())
	require.Nil(t, reader)
	require.EqualValues(t, 1, requests.Load())
	waitForSSEHTTPSignal(t, body.closed)
}
