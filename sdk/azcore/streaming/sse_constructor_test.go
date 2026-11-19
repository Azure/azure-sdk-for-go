// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package streaming_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
	"testing/synctest"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialConnectionIsEagerAndSeedsCursor(t *testing.T) {
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
			type contextKey struct{}
			key := contextKey{}
			ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), key, "operation"), time.Hour)
			defer cancel()
			var readerCtx context.Context
			calls, decodes := 0, 0
			body := &observedSSEBody{ReadCloser: eventBody("data: one\n\n")}
			reader, err := streaming.NewSSEReader(ctx, func(connectCtx context.Context, lastEventID string) (io.ReadCloser, error) {
				calls++
				readerCtx = connectCtx
				assert.Equal(t, tc.wantID, lastEventID)
				return body, nil
			}, streaming.SSEHandler[frameView]{
				Decode: func(frame streaming.SSEFrame) (frameView, bool, error) {
					decodes++
					return decodeView(frame)
				},
			}, tc.options)
			require.NoError(t, err)
			t.Cleanup(func() { assert.NoError(t, reader.Close()) })
			assert.Equal(t, 1, calls)
			assert.Zero(t, decodes)
			require.NotNil(t, readerCtx)
			assert.NotSame(t, ctx, readerCtx)
			assert.Equal(t, "operation", readerCtx.Value(key))
			wantDeadline, _ := ctx.Deadline()
			deadline, ok := readerCtx.Deadline()
			assert.True(t, ok)
			assert.Equal(t, wantDeadline, deadline)
			assert.Equal(t, tc.wantID, reader.LastEventID())
			event, err := reader.Next()
			require.NoError(t, err)
			assert.Equal(t, "one", event.data)
			assert.Equal(t, 1, decodes)
			assert.Equal(t, tc.wantID, reader.LastEventID())
			require.NoError(t, reader.Close())
			assert.ErrorIs(t, readerCtx.Err(), context.Canceled)
			assert.NoError(t, ctx.Err())
			assert.EqualValues(t, 1, body.closed.Load())
			_, err = reader.Next()
			require.ErrorIs(t, err, io.EOF)
			assert.Equal(t, 1, calls)
		})
	}
}

func TestInitialConnectionCompletionAndFailures(t *testing.T) {
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
		{name: "completion with EOF close failure", withBody: true, connectErr: io.EOF, closeErr: io.EOF},
		{name: "failure", connectErr: connectErr},
		{name: "failure with body", withBody: true, connectErr: connectErr},
		{name: "failure with close failure", withBody: true, connectErr: connectErr, closeErr: closeErr},
		{name: "failure with EOF close failure", withBody: true, connectErr: connectErr, closeErr: io.EOF},
		{name: "failure with wrapped EOF close failure", withBody: true, connectErr: connectErr, closeErr: fmt.Errorf("closing body: %w", io.EOF)},
		{name: "nil body"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			var readerCtx context.Context
			var body *observedSSEBody
			if tc.withBody {
				body = &observedSSEBody{ReadCloser: eventBody(""), closeErr: tc.closeErr}
			}
			reader, err := streaming.NewSSEReader(context.Background(), func(ctx context.Context, _ string) (io.ReadCloser, error) {
				calls++
				readerCtx = ctx
				if tc.withBody {
					return body, tc.connectErr
				}
				return nil, tc.connectErr
			}, streaming.SSEHandler[frameView]{Decode: decodeView, Reconnect: true}, &streaming.SSEReaderOptions{LastEventID: "prior"})
			if tc.connectErr == io.EOF && tc.closeErr == nil {
				require.NoError(t, err)
				require.NotNil(t, reader)
				t.Cleanup(func() { assert.NoError(t, reader.Close()) })
				assert.Equal(t, "prior", reader.LastEventID())
				for range 2 {
					_, err = reader.Next()
					require.ErrorIs(t, err, io.EOF)
				}
			} else {
				require.Error(t, err)
				assert.Nil(t, reader)
				if tc.closeErr != nil && (tc.connectErr == io.EOF || !errors.Is(tc.closeErr, io.EOF)) {
					assert.ErrorIs(t, err, tc.closeErr)
					if !errors.Is(tc.closeErr, io.EOF) {
						assert.NotErrorIs(t, err, io.EOF)
					}
				}
				if tc.connectErr != nil && tc.connectErr != io.EOF {
					assert.ErrorIs(t, err, tc.connectErr)
					assert.NotErrorIs(t, err, io.EOF)
				}
				if tc.connectErr == nil {
					assert.ErrorContains(t, err, "nil body")
				}
			}
			assert.Equal(t, 1, calls)
			require.NotNil(t, readerCtx)
			assert.ErrorIs(t, readerCtx.Err(), context.Canceled)
			if body != nil {
				assert.EqualValues(t, 1, body.closed.Load())
			}
		})
	}
}

func TestInitialConnectionWithCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	reader, err := streaming.NewSSEReader(ctx, func(context.Context, string) (io.ReadCloser, error) {
		calls++
		return nil, errors.New("unexpected connection")
	}, streaming.SSEHandler[frameView]{Decode: decodeView}, nil)
	require.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, reader)
	assert.Zero(t, calls)
}

func TestInitialConnectionDeadline(t *testing.T) {
	for _, tc := range []struct {
		name     string
		lateBody bool
		closeErr error
	}{
		{name: "connector honors deadline"},
		{name: "connector returns late body", lateBody: true},
		{name: "late body cleanup EOF", lateBody: true, closeErr: io.EOF},
		{name: "late body cleanup wrapped EOF", lateBody: true, closeErr: fmt.Errorf("closing body: %w", io.EOF)},
		{name: "late body cleanup failure", lateBody: true, closeErr: errors.New("close failed")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				calls := 0
				body := &observedSSEBody{ReadCloser: eventBody("data: late\n\n"), closeErr: tc.closeErr}
				reader, err := streaming.NewSSEReader(ctx, func(connectCtx context.Context, _ string) (io.ReadCloser, error) {
					calls++
					<-connectCtx.Done()
					if tc.lateBody {
						return body, nil
					}
					return nil, connectCtx.Err()
				}, streaming.SSEHandler[frameView]{Decode: decodeView}, nil)
				if tc.closeErr != nil && !errors.Is(tc.closeErr, io.EOF) {
					require.ErrorIs(t, err, tc.closeErr)
				} else {
					require.ErrorIs(t, err, context.DeadlineExceeded)
				}
				require.NotErrorIs(t, err, io.EOF)
				assert.Nil(t, reader)
				assert.Equal(t, 1, calls)
				if tc.lateBody {
					assert.EqualValues(t, 1, body.closed.Load())
				} else {
					require.NoError(t, body.Close())
				}
			})
		})
	}
}
