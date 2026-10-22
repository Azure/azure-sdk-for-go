// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package streaming_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSSEReader(t *testing.T, body io.ReadCloser, handler streaming.SSEHandler[frameView], options *streaming.SSEReaderOptions) *streaming.SSEReader[frameView] {
	t.Helper()
	return newSSEReaderWithContext(t, context.Background(), initialBodyConnector(body, nil), handler, options)
}

func newSSEReaderWithContext(t *testing.T, ctx context.Context, connect streaming.SSEConnector, handler streaming.SSEHandler[frameView], options *streaming.SSEReaderOptions) *streaming.SSEReader[frameView] {
	t.Helper()
	if handler.Decode == nil {
		handler.Decode = decodeView
	}
	reader, err := streaming.NewSSEReader(ctx, connect, handler, options)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, reader.Close()) })
	return reader
}

func initialBodyConnector(body io.ReadCloser, connect streaming.SSEConnector) streaming.SSEConnector {
	return func(ctx context.Context, lastEventID string) (io.ReadCloser, error) {
		if body != nil {
			initialBody := body
			body = nil
			return initialBody, nil
		}
		if connect == nil {
			return nil, errors.New("unexpected reconnect")
		}
		return connect(ctx, lastEventID)
	}
}

func eventBody(body string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(body))
}

type sseErrorReader struct{ err error }

func (r sseErrorReader) Read([]byte) (int, error) { return 0, r.err }

func TestMetadataOnlyBlocksDoNotDecode(t *testing.T) {
	for _, prefix := range []string{": heartbeat\n\n", "retry: 0\n\n", "id: first\n\n", "event: ignored\n\n", "unknown: field\n\n"} {
		t.Run(strings.TrimSpace(prefix), func(t *testing.T) {
			reader := newSSEReader(t, eventBody(prefix+"data: hello\n\n"), streaming.SSEHandler[frameView]{}, nil)
			event, err := reader.Next()
			require.NoError(t, err)
			assert.Equal(t, "hello", event.data)
			assert.Empty(t, event.eventType)
			_, err = reader.Next()
			require.ErrorIs(t, err, io.EOF)
		})
	}
}

func TestMetadataOnlyIDIsCommitted(t *testing.T) {
	reader := newSSEReader(t, eventBody("id: first\n\n: heartbeat\n\n"), streaming.SSEHandler[frameView]{}, nil)
	_, err := reader.Next()
	require.ErrorIs(t, err, io.EOF)
	assert.Equal(t, "first", reader.LastEventID())
}

func TestEmptyDataEventIsDelivered(t *testing.T) {
	reader := newSSEReader(t, eventBody("data:\n\n"), streaming.SSEHandler[frameView]{}, nil)
	event, err := reader.Next()
	require.NoError(t, err)
	assert.Empty(t, event.data)
	_, err = reader.Next()
	require.ErrorIs(t, err, io.EOF)
}

func TestIncompleteBlockDoesNotAdvanceID(t *testing.T) {
	for _, ending := range []string{"id: second\ndata: incomplete", "id: second\ndata: incomplete\n", "id: second\n"} {
		t.Run(ending, func(t *testing.T) {
			reader := newSSEReader(t, eventBody("id: first\ndata: complete\n\n"+ending), streaming.SSEHandler[frameView]{}, nil)
			event, err := reader.Next()
			require.NoError(t, err)
			assert.Equal(t, "complete", event.data)
			_, err = reader.Next()
			require.ErrorIs(t, err, io.EOF)
			assert.Equal(t, "first", reader.LastEventID())
		})
	}
}

func TestLeadingBOMIsIgnored(t *testing.T) {
	got, _ := collect(t, "\uFEFFdata: hello\n\n", decodeView)
	require.Len(t, got, 1)
	assert.Equal(t, "hello", got[0].data)
}

func TestReconnectAfterEmptyAndIDlessSegments(t *testing.T) {
	segments := []string{"", "data: one\n\n", ""}
	calls := 0
	reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(eventBody(""),
		func(_ context.Context, lastEventID string) (io.ReadCloser, error) {
			assert.Empty(t, lastEventID)
			if calls == len(segments)-1 {
				calls++
				return nil, io.EOF
			}
			body := segments[calls]
			calls++
			return eventBody(body), nil
		}),
		streaming.SSEHandler[frameView]{Reconnect: true}, &streaming.SSEReaderOptions{ReconnectDelay: time.Nanosecond})
	event, err := reader.Next()
	require.NoError(t, err)
	assert.Equal(t, "one", event.data)
	_, err = reader.Next()
	require.ErrorIs(t, err, io.EOF)
	assert.Equal(t, 3, calls)
}

func TestReconnectPreservesIDUntilReset(t *testing.T) {
	for _, reset := range []bool{false, true} {
		t.Run(strconv.FormatBool(reset), func(t *testing.T) {
			var ids []string
			reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(eventBody("id: first\ndata: one\n\n"),
				func(_ context.Context, lastEventID string) (io.ReadCloser, error) {
					ids = append(ids, lastEventID)
					if len(ids) == 2 {
						return nil, io.EOF
					}
					body := "data: two\n\n"
					if reset {
						body = "id:\n" + body
					}
					return eventBody(body), nil
				}),
				streaming.SSEHandler[frameView]{Reconnect: true}, &streaming.SSEReaderOptions{ReconnectDelay: time.Nanosecond})
			_, err := reader.Next()
			require.NoError(t, err)
			_, err = reader.Next()
			require.NoError(t, err)
			wantID := "first"
			if reset {
				wantID = ""
			}
			assert.Equal(t, wantID, reader.LastEventID())
			_, err = reader.Next()
			require.ErrorIs(t, err, io.EOF)
			assert.Equal(t, []string{"first", wantID}, ids)
		})
	}
}

func TestResumeSeedsIDFromOptions(t *testing.T) {
	for _, tc := range []struct {
		name    string
		options *streaming.SSEReaderOptions
		wantID  string
	}{
		{name: "nil options"},
		{name: "zero options", options: &streaming.SSEReaderOptions{}},
		{name: "checkpoint", options: &streaming.SSEReaderOptions{LastEventID: "prior"}, wantID: "prior"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader := newSSEReader(t, eventBody("data: one\n\n"), streaming.SSEHandler[frameView]{}, tc.options)
			assert.Equal(t, tc.wantID, reader.LastEventID())
			_, err := reader.Next()
			require.NoError(t, err)
			assert.Equal(t, tc.wantID, reader.LastEventID())
			resetReader := newSSEReader(t, eventBody("id:\ndata: reset\n\n"), streaming.SSEHandler[frameView]{}, tc.options)
			_, err = resetReader.Next()
			require.NoError(t, err)
			assert.Empty(t, resetReader.LastEventID())
		})
	}
}

func TestReconnectReadErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		err     error
		recover bool
	}{
		{name: "unexpected EOF", err: io.ErrUnexpectedEOF, recover: true},
		{name: "network failure", err: &net.OpError{Op: "read", Net: "tcp", Err: errors.New("connection reset")}, recover: true},
		{name: "attempt timeout", err: context.DeadlineExceeded, recover: true},
		{name: "non-retriable", err: errorinfo.NonRetriableError(io.ErrUnexpectedEOF)},
		{name: "unknown read error", err: errors.New("fatal read error")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(io.NopCloser(io.MultiReader(strings.NewReader("id: first\ndata: one\n\n"), sseErrorReader{err: tc.err})),
				func(_ context.Context, lastEventID string) (io.ReadCloser, error) {
					calls++
					assert.Equal(t, "first", lastEventID)
					return eventBody("data: two\n\n"), nil
				}),
				streaming.SSEHandler[frameView]{Reconnect: true}, &streaming.SSEReaderOptions{ReconnectDelay: time.Nanosecond})
			_, err := reader.Next()
			require.NoError(t, err)
			event, err := reader.Next()
			if tc.recover {
				require.NoError(t, err)
				assert.Equal(t, "two", event.data)
				assert.Equal(t, 1, calls)
			} else {
				require.ErrorIs(t, err, tc.err)
				assert.Zero(t, calls)
				_, err = reader.Next()
				require.ErrorIs(t, err, io.EOF)
			}
		})
	}
}

func TestFiniteStreamReadErrorDoesNotReconnect(t *testing.T) {
	calls := 0
	reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(io.NopCloser(sseErrorReader{err: io.ErrUnexpectedEOF}),
		func(context.Context, string) (io.ReadCloser, error) {
			calls++
			return nil, errors.New("unexpected reconnect")
		}),
		streaming.SSEHandler[frameView]{}, nil)
	_, err := reader.Next()
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	assert.Zero(t, calls)
}

func TestDecodeErrorIsTerminal(t *testing.T) {
	calls := 0
	reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(eventBody("data: {invalid}\n\ndata: valid\n\n"),
		func(context.Context, string) (io.ReadCloser, error) {
			calls++
			return nil, errors.New("unexpected reconnect")
		}),
		streaming.SSEHandler[frameView]{Reconnect: true}, nil)
	_, err := reader.Next()
	require.Error(t, err)
	_, err = reader.Next()
	require.ErrorIs(t, err, io.EOF)
	assert.Zero(t, calls)
}

func TestTerminalEventDoesNotReconnect(t *testing.T) {
	calls := 0
	reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(eventBody("id: first\ndata: [DONE]\n\n"),
		func(context.Context, string) (io.ReadCloser, error) {
			calls++
			return nil, errors.New("unexpected reconnect")
		}),
		streaming.SSEHandler[frameView]{Reconnect: true}, nil)
	_, err := reader.Next()
	require.ErrorIs(t, err, io.EOF)
	assert.Zero(t, calls)
}

func TestCloseStopsBufferedEventsAndReconnect(t *testing.T) {
	calls := 0
	reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(eventBody("id: first\ndata: one\n\ndata: two\n\n"),
		func(context.Context, string) (io.ReadCloser, error) {
			calls++
			return nil, errors.New("unexpected reconnect")
		}),
		streaming.SSEHandler[frameView]{Reconnect: true}, nil)
	_, err := reader.Next()
	require.NoError(t, err)
	require.NoError(t, reader.Close())
	require.NoError(t, reader.Close())
	_, err = reader.Next()
	require.ErrorIs(t, err, io.EOF)
	assert.Zero(t, calls)
}

func TestReconnectDelays(t *testing.T) {
	for _, tc := range []struct {
		name    string
		body    string
		options *streaming.SSEReaderOptions
		want    time.Duration
	}{
		{name: "default", body: "data: one\n\n", want: 3 * time.Second},
		{name: "configured", body: "data: one\n\n", options: &streaming.SSEReaderOptions{ReconnectDelay: 2 * time.Second}, want: 2 * time.Second},
		{name: "server metadata", body: "retry: 1500\n\ndata: one\n\n", want: 1500 * time.Millisecond},
		{name: "server zero", body: "retry: 0\n\ndata: one\n\n"},
		{name: "incomplete retry block", body: "data: one\n\nretry: 1500\n", want: 1500 * time.Millisecond},
		{name: "invalid retry", body: "retry: invalid\ndata: one\n\n", want: 3 * time.Second},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var elapsed time.Duration
				start := time.Now()
				reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(eventBody(tc.body),
					func(context.Context, string) (io.ReadCloser, error) {
						elapsed = time.Since(start)
						return nil, io.EOF
					}),
					streaming.SSEHandler[frameView]{Reconnect: true}, tc.options)
				_, err := reader.Next()
				require.NoError(t, err)
				_, err = reader.Next()
				require.ErrorIs(t, err, io.EOF)
				assert.Equal(t, tc.want, elapsed)
			})
		})
	}
}

func TestCancellationAndCloseInterruptReconnectDelay(t *testing.T) {
	for _, closeReader := range []bool{false, true} {
		t.Run(strconv.FormatBool(closeReader), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				var calls atomic.Int32
				reader := newSSEReaderWithContext(t, ctx, initialBodyConnector(eventBody(""),
					func(context.Context, string) (io.ReadCloser, error) {
						calls.Add(1)
						return nil, errors.New("unexpected reconnect")
					}),
					streaming.SSEHandler[frameView]{Reconnect: true}, nil)
				result := make(chan error, 1)
				go func() { _, err := reader.Next(); result <- err }()
				synctest.Wait()
				if closeReader {
					require.NoError(t, reader.Close())
					require.ErrorIs(t, <-result, io.EOF)
				} else {
					cancel()
					require.ErrorIs(t, <-result, context.Canceled)
				}
				assert.Zero(t, calls.Load())
			})
		})
	}
}

func TestCloseCancelsPendingConnect(t *testing.T) {
	started := make(chan struct{})
	reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(eventBody("retry: 0\n\n"),
		func(ctx context.Context, _ string) (io.ReadCloser, error) {
			close(started)
			<-ctx.Done()
			return nil, ctx.Err()
		}),
		streaming.SSEHandler[frameView]{Reconnect: true}, nil)
	result := make(chan error, 1)
	go func() { _, err := reader.Next(); result <- err }()
	<-started
	require.NoError(t, reader.Close())
	require.ErrorIs(t, <-result, io.EOF)
}

type observedSSEBody struct {
	io.ReadCloser
	closed   atomic.Int32
	closeErr error
}

func (b *observedSSEBody) Close() error {
	b.closed.Add(1)
	return errors.Join(b.ReadCloser.Close(), b.closeErr)
}

func TestCloseDiscardsLateConnectResponse(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	body := &observedSSEBody{ReadCloser: eventBody("data: late\n\n")}
	reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(eventBody("retry: 0\n\n"),
		func(context.Context, string) (io.ReadCloser, error) {
			close(started)
			<-release
			return body, nil
		}),
		streaming.SSEHandler[frameView]{Reconnect: true}, nil)
	result := make(chan error, 1)
	go func() { _, err := reader.Next(); result <- err }()
	<-started
	require.NoError(t, reader.Close())
	close(release)
	require.ErrorIs(t, <-result, io.EOF)
	assert.EqualValues(t, 1, body.closed.Load())
}

type notifiedSSEBody struct {
	io.ReadCloser
	started chan struct{}
}

func (b *notifiedSSEBody) Read(p []byte) (int, error) {
	close(b.started)
	return b.ReadCloser.Read(p)
}

func TestCloseInterruptsBlockedRead(t *testing.T) {
	read, write := io.Pipe()
	defer func() {
		_ = write.Close()
	}()
	body := &notifiedSSEBody{ReadCloser: read, started: make(chan struct{})}
	reader := newSSEReader(t, body, streaming.SSEHandler[frameView]{}, nil)
	result := make(chan error, 1)
	go func() { _, err := reader.Next(); result <- err }()
	<-body.started
	require.NoError(t, reader.Close())
	require.ErrorIs(t, <-result, io.EOF)
}

func TestZeroReaderIsComplete(t *testing.T) {
	var reader streaming.SSEReader[frameView]
	assert.Empty(t, reader.LastEventID())
	_, err := reader.Next()
	require.ErrorIs(t, err, io.EOF)
	require.NoError(t, reader.Close())
}

func TestRetryDurationDoesNotOverflow(t *testing.T) {
	if strconv.IntSize < 64 {
		t.Skip("requires a retry value larger than a time.Duration")
	}
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var calls atomic.Int32
		body := "retry: " + strconv.FormatInt(math.MaxInt64, 10) + "\n\n"
		reader := newSSEReaderWithContext(t, ctx, initialBodyConnector(eventBody(body),
			func(context.Context, string) (io.ReadCloser, error) {
				calls.Add(1)
				return nil, errors.New("unexpected reconnect")
			}),
			streaming.SSEHandler[frameView]{Reconnect: true}, nil)
		result := make(chan error, 1)
		go func() { _, err := reader.Next(); result <- err }()
		synctest.Wait()
		assert.Zero(t, calls.Load())
		cancel()
		require.ErrorIs(t, <-result, context.Canceled)
	})
}

func TestInvalidReaderConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name      string
		ctx       context.Context
		noConnect bool
		handler   streaming.SSEHandler[frameView]
		options   *streaming.SSEReaderOptions
		wantErr   string
	}{
		{name: "nil context", handler: streaming.SSEHandler[frameView]{Decode: decodeView}, wantErr: "context must not be nil"},
		{name: "nil decoder", ctx: context.Background(), wantErr: "Decode must not be nil"},
		{name: "nil connector", ctx: context.Background(), noConnect: true, handler: streaming.SSEHandler[frameView]{Decode: decodeView}, wantErr: "Connect must not be nil"},
		{name: "negative delay", ctx: context.Background(), handler: streaming.SSEHandler[frameView]{Decode: decodeView}, options: &streaming.SSEReaderOptions{ReconnectDelay: -time.Second}, wantErr: "ReconnectDelay must not be negative"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &observedSSEBody{ReadCloser: eventBody("")}
			t.Cleanup(func() { assert.NoError(t, body.Close()) })
			calls := 0
			connect := streaming.SSEConnector(func(context.Context, string) (io.ReadCloser, error) {
				calls++
				return body, nil
			})
			if tc.noConnect {
				connect = nil
			}
			reader, err := streaming.NewSSEReader(tc.ctx, connect, tc.handler, tc.options)
			require.ErrorContains(t, err, tc.wantErr)
			assert.Nil(t, reader)
			assert.Zero(t, calls)
			assert.Zero(t, body.closed.Load())
		})
	}
}

func TestConnectFailureIsNotRetried(t *testing.T) {
	for _, invalidResponse := range []bool{false, true} {
		t.Run(strconv.FormatBool(invalidResponse), func(t *testing.T) {
			calls := 0
			connectErr := errors.New("pipeline exhausted retries")
			reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(eventBody("retry: 0\n\n"),
				func(context.Context, string) (io.ReadCloser, error) {
					calls++
					if invalidResponse {
						return nil, nil
					}
					return nil, connectErr
				}),
				streaming.SSEHandler[frameView]{Reconnect: true}, nil)
			_, err := reader.Next()
			if invalidResponse {
				require.ErrorContains(t, err, "nil body")
			} else {
				require.ErrorIs(t, err, connectErr)
			}
			_, err = reader.Next()
			require.ErrorIs(t, err, io.EOF)
			assert.Equal(t, 1, calls)
		})
	}
}

func TestConnectCompletionAndErrorBodyCleanup(t *testing.T) {
	connectErr := errors.New("connection failed")
	closeErr := errors.New("body close failed")
	for _, tc := range []struct {
		name       string
		withBody   bool
		connectErr error
		closeErr   error
	}{
		{name: "completion", connectErr: io.EOF},
		{name: "completion with body", withBody: true, connectErr: io.EOF},
		{name: "completion with close failure", withBody: true, connectErr: io.EOF, closeErr: closeErr},
		{name: "failure with body", withBody: true, connectErr: connectErr},
		{name: "failure with close failure", withBody: true, connectErr: connectErr, closeErr: closeErr},
		{name: "failure with EOF close failure", withBody: true, connectErr: connectErr, closeErr: io.EOF},
		{name: "failure with wrapped EOF close failure", withBody: true, connectErr: connectErr, closeErr: fmt.Errorf("closing body: %w", io.EOF)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			body := &observedSSEBody{ReadCloser: eventBody(""), closeErr: tc.closeErr}
			reader := newSSEReaderWithContext(t, context.Background(), initialBodyConnector(eventBody("retry: 0\n\n"),
				func(context.Context, string) (io.ReadCloser, error) {
					calls++
					if tc.withBody {
						return body, tc.connectErr
					}
					return nil, tc.connectErr
				}),
				streaming.SSEHandler[frameView]{Reconnect: true}, nil)
			_, err := reader.Next()
			if tc.closeErr != nil && !errors.Is(tc.closeErr, io.EOF) {
				require.ErrorIs(t, err, tc.closeErr)
				assert.NotErrorIs(t, err, io.EOF)
			}
			if tc.connectErr != io.EOF || tc.closeErr == nil {
				require.ErrorIs(t, err, tc.connectErr)
			}
			if tc.connectErr != io.EOF {
				assert.NotErrorIs(t, err, io.EOF)
			}
			if tc.withBody {
				assert.EqualValues(t, 1, body.closed.Load())
			}
			_, err = reader.Next()
			require.ErrorIs(t, err, io.EOF)
			assert.Equal(t, 1, calls)
		})
	}
}
