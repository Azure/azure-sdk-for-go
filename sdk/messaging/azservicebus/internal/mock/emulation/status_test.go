// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStatus(t *testing.T) {
	top := NewStatus(nil)

	top.CloseWithError(errors.New("Hello World"))
	top.CloseWithError(errors.New("Hello World 2"))

	require.EqualError(t, top.Err(), "Hello World")

	child := NewStatus(top)

	select {
	case <-child.Done():
		require.EqualError(t, child.Err(), "Hello World")
	default:
		require.Fail(t, "Should have been closed already")
	}
}

func TestMockReceiverDone(t *testing.T) {
	for _, source := range []string{"receiver", "session", "connection"} {
		t.Run(source, func(t *testing.T) {
			// Independent statuses ensure Done does not rely on asynchronous parent propagation.
			conn := &MockConnection{Status: NewStatus(nil)}
			sess := &MockSession{Conn: conn, Status: NewStatus(nil)}
			rcvr := &MockReceiver{Session: sess, Status: NewStatus(nil)}
			wantErr := errors.New(source + " closed")

			switch source {
			case "receiver":
				rcvr.Status.CloseWithError(wantErr)
			case "session":
				sess.Status.CloseWithError(wantErr)
			case "connection":
				conn.Status.CloseWithError(wantErr)
			}

			timeout := time.NewTimer(5 * time.Second)
			defer timeout.Stop()

			select {
			case err := <-rcvr.Done():
				require.ErrorIs(t, err, wantErr)
			case <-timeout.C:
				require.FailNow(t, "receiver did not report termination")
			}
		})
	}
}
