// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:generate mockgen -source amqpwrap.go -package amqpwrap -copyright_file ../mock/testdata/copyright.txt -destination mock_amqp_test.go

package amqpwrap

import (
	"context"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAMQPReceiverWrapper(t *testing.T) {
	for _, e := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(e.Error(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			r := NewMockgoamqpReceiver(ctrl)

			r.EXPECT().Close(gomock.Any()).Return(e)

			rw := &AMQPReceiverWrapper{inner: r}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			require.ErrorIs(t, ErrConnResetNeeded, rw.Close(ctx))
		})
	}
}

func TestAMQPSenderWrapper(t *testing.T) {
	for _, e := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(e.Error(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := NewMockAMQPSenderCloser(ctrl)

			s.EXPECT().Close(gomock.Any()).Return(e)

			rw := &AMQPSenderWrapper{inner: s}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			require.ErrorIs(t, ErrConnResetNeeded, rw.Close(ctx))
		})
	}
}

func TestAMQPSessionWrapper(t *testing.T) {
	for _, e := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(fmt.Sprintf("Close() == %s", e.Error()), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := NewMockgoamqpSession(ctrl)

			s.EXPECT().Close(gomock.Any()).Return(e)

			rw := &AMQPSessionWrapper{Inner: s}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			require.ErrorIs(t, ErrConnResetNeeded, rw.Close(ctx))
		})
	}

	for _, e := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(fmt.Sprintf("NewReceiver() == %s", e.Error()), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			inner := NewMockgoamqpSession(ctrl)

			inner.EXPECT().NewReceiver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, e)

			rw := &AMQPSessionWrapper{Inner: inner}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			receiver, err := rw.NewReceiver(ctx, "source", nil)
			require.ErrorIs(t, err, ErrConnResetNeeded)
			require.Nil(t, receiver)
		})
	}

	for _, e := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(fmt.Sprintf("NewSender() == %s", e.Error()), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			inner := NewMockgoamqpSession(ctrl)

			inner.EXPECT().NewSender(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, e)

			rw := &AMQPSessionWrapper{Inner: inner}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			sender, err := rw.NewSender(ctx, "target", nil)
			require.ErrorIs(t, err, ErrConnResetNeeded)
			require.Nil(t, sender)
		})
	}
}

func TestAMQPConnWrapper(t *testing.T) {
	for _, e := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(e.Error(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			inner := NewMockgoamqpConn(ctrl)

			inner.EXPECT().NewSession(gomock.Any(), gomock.Any()).Return(nil, e)

			conn := &AMQPClientWrapper{Inner: inner}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			sess, err := conn.NewSession(ctx, nil)
			require.ErrorIs(t, err, ErrConnResetNeeded)
			require.Nil(t, sess)
		})
	}
}
