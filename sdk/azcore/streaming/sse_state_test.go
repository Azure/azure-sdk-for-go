// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package streaming

import (
	"context"
	"io"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

// Pause the context check after a reconnect installs its response body.
type reconnectStateContext struct {
	context.Context
	armed   atomic.Bool
	check   chan struct{}
	release chan struct{}
}

func (c *reconnectStateContext) Err() error {
	if c.armed.CompareAndSwap(true, false) {
		close(c.check)
		<-c.release
	}
	return c.Context.Err()
}

func TestReconnectStateCheckInterruption(t *testing.T) {
	for _, closeReader := range []bool{true, false} {
		name := "operation cancellation"
		if closeReader {
			name = "explicit close"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			body, writer := io.Pipe()
			defer func() { _ = writer.Close() }()
			var stateContext *reconnectStateContext
			calls := 0
			reader, err := NewSSEReader(ctx, func(context.Context, string) (io.ReadCloser, error) {
				calls++
				if calls == 1 {
					return io.NopCloser(strings.NewReader("retry: 0\n\n")), nil
				}
				stateContext.armed.Store(true)
				return body, nil
			}, SSEHandler[string]{
				Decode: func(frame SSEFrame) (string, bool, error) {
					return string(frame.Data), false, nil
				},
				Reconnect: true,
			}, nil)
			require.NoError(t, err)
			defer func() { _ = reader.Close() }()
			stateContext = &reconnectStateContext{
				Context: reader.ctx,
				check:   make(chan struct{}),
				release: make(chan struct{}),
			}
			reader.ctx = stateContext
			result := make(chan error, 1)
			go func() {
				_, err := reader.Next()
				result <- err
			}()
			<-stateContext.check
			closeResult := make(chan error, 1)
			if closeReader {
				// Complete Close at the paused check if it isn't protected by
				// the state lock; otherwise let Close wait for that lock.
				if reader.mu.TryLock() {
					reader.mu.Unlock()
					closeResult <- reader.Close()
				} else {
					go func() { closeResult <- reader.Close() }()
				}
			} else {
				cancel()
			}
			close(stateContext.release)
			wantErr := error(context.Canceled)
			if closeReader {
				wantErr = io.EOF
			}
			require.ErrorIs(t, <-result, wantErr)
			if closeReader {
				require.NoError(t, <-closeResult)
			}
			require.Equal(t, 2, calls)
			_, err = reader.Next()
			require.ErrorIs(t, err, io.EOF)
			_, err = body.Read(make([]byte, 1))
			require.ErrorIs(t, err, io.ErrClosedPipe)
		})
	}
}
