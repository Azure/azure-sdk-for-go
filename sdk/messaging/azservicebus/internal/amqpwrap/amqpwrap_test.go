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
	for _, expectedErr := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(expectedErr.Error(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			r := NewMockgoamqpReceiver(ctrl)

			r.EXPECT().Close(gomock.Any()).Return(expectedErr)

			rw := &AMQPReceiverWrapper{Inner: r, ContextWithTimeoutFn: context.WithTimeout}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err := rw.Close(ctx)
			require.ErrorIs(t, err, expectedErr)
		})
	}
}

func TestAMQPSenderWrapper(t *testing.T) {
	for _, expectedErr := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(expectedErr.Error(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := NewMockAMQPSenderCloser(ctrl)

			s.EXPECT().Close(gomock.Any()).Return(expectedErr)

			rw := &AMQPSenderWrapper{Inner: s, ContextWithTimeoutFn: context.WithTimeout}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			require.ErrorIs(t, rw.Close(ctx), expectedErr)
		})
	}
}

func TestAMQPSessionWrapper(t *testing.T) {
	for _, expectedErr := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(fmt.Sprintf("Close() == %s", expectedErr.Error()), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := NewMockgoamqpSession(ctrl)

			s.EXPECT().Close(gomock.Any()).Return(expectedErr)

			rw := &AMQPSessionWrapper{Inner: s, ContextWithTimeoutFn: context.WithTimeout}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			require.ErrorIs(t, rw.Close(ctx), expectedErr)
		})
	}

	for _, expectedErr := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(fmt.Sprintf("NewReceiver() == %s", expectedErr.Error()), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			inner := NewMockgoamqpSession(ctrl)

			inner.EXPECT().NewReceiver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expectedErr)

			rw := &AMQPSessionWrapper{Inner: inner, ContextWithTimeoutFn: context.WithTimeout}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			receiver, err := rw.NewReceiver(ctx, "source", nil)
			require.ErrorIs(t, err, expectedErr)
			require.Nil(t, receiver)
		})
	}

	for _, expectedErr := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(fmt.Sprintf("NewSender() == %s", expectedErr.Error()), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			inner := NewMockgoamqpSession(ctrl)

			inner.EXPECT().NewSender(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expectedErr)

			rw := &AMQPSessionWrapper{Inner: inner}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			sender, err := rw.NewSender(ctx, "target", nil)
			require.ErrorIs(t, err, expectedErr)
			require.Nil(t, sender)
		})
	}
}

func TestAMQPConnWrapper(t *testing.T) {
	for _, expectedErr := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(expectedErr.Error(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			inner := NewMockgoamqpConn(ctrl)

			inner.EXPECT().NewSession(gomock.Any(), gomock.Any()).Return(nil, expectedErr)

			conn := &AMQPClientWrapper{Inner: inner}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			sess, err := conn.NewSession(ctx, nil)
			require.ErrorIs(t, err, expectedErr)
			require.Nil(t, sess)
		})
	}
}
