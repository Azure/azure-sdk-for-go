// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:generate mockgen -source amqpwrap.go -package amqpwrap -copyright_file ../mock/testdata/copyright.txt -destination mock_amqp_test.go

package amqpwrap

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/Azure/go-amqp"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAMQPReceiverWrapper(t *testing.T) {
	t.Run("errors are wrapped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		inner := NewMockgoamqpReceiver(ctrl)

		inner.EXPECT().LinkName().Return("receiver").AnyTimes()

		inner.EXPECT().Receive(gomock.Any(), gomock.Any()).Return(nil, errors.New("receive failed"))
		inner.EXPECT().AcceptMessage(gomock.Any(), gomock.Any()).Return(errors.New("accept failed"))
		inner.EXPECT().ModifyMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("modify failed"))
		inner.EXPECT().RejectMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("reject failed"))
		inner.EXPECT().ReleaseMessage(gomock.Any(), gomock.Any()).Return(errors.New("release failed"))
		inner.EXPECT().IssueCredit(gomock.Any()).Return(errors.New("issue credit failed"))

		inner.EXPECT().Close(test.NotCancelledAndHasTimeout).Return(errors.New("close failed"))
		inner.EXPECT().Close(test.CancelledAndHasTimeout).Return(context.Canceled)

		rw := &AMQPReceiverWrapper{Inner: inner, ContextWithTimeoutFn: test.NewContextWithTimeoutForTests, connID: uint64(101)}

		assertErr := func(err error, msg string) {
			t.Helper()
			var wrapErr Error
			require.ErrorAs(t, err, &wrapErr)
			require.EqualError(t, wrapErr, msg)
			require.Equal(t, uint64(101), wrapErr.ConnID)
			require.Equal(t, "receiver", wrapErr.LinkName)
		}

		_, err := rw.Receive(context.Background(), nil)
		assertErr(err, "receive failed")

		err = rw.AcceptMessage(context.Background(), nil)
		assertErr(err, "accept failed")

		err = rw.ModifyMessage(context.Background(), nil, nil)
		assertErr(err, "modify failed")

		err = rw.ReleaseMessage(context.Background(), nil)
		assertErr(err, "release failed")

		err = rw.RejectMessage(context.Background(), nil, nil)
		assertErr(err, "reject failed")

		err = rw.IssueCredit(uint32(100))
		assertErr(err, "issue credit failed")

		err = rw.Close(context.Background())
		assertErr(err, "close failed")

		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		err = rw.Close(cancelledCtx)
		require.ErrorIs(t, err, context.Canceled)
		assertErr(err, "context canceled")
	})

	t.Run("normal usage", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		inner := NewMockgoamqpReceiver(ctrl)

		inner.EXPECT().LinkName().Return("receiver").AnyTimes()
		inner.EXPECT().IssueCredit(gomock.Any()).Return(nil)
		inner.EXPECT().Receive(test.NotCancelled, gomock.Any()).Return(&amqp.Message{}, nil)
		inner.EXPECT().Receive(test.Cancelled, gomock.Any()).Return(nil, context.Canceled)
		inner.EXPECT().Prefetched().Return(&amqp.Message{})
		inner.EXPECT().Prefetched().Return(nil)
		inner.EXPECT().LinkSourceFilterValue("hello").Return("world")

		rw := &AMQPReceiverWrapper{Inner: inner, ContextWithTimeoutFn: test.NewContextWithTimeoutForTests, connID: uint64(101)}

		require.Equal(t, uint64(101), rw.ConnID())
		require.Equal(t, "world", rw.LinkSourceFilterValue("hello"))

		require.Equal(t, uint32(0), rw.Credits())

		err := rw.IssueCredit(10)
		require.NoError(t, err)

		require.Equal(t, uint32(10), rw.Credits())

		msg, err := rw.Receive(context.Background(), nil)
		require.NotNil(t, msg)
		require.NoError(t, err)

		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		msg, err = rw.Receive(cancelledCtx, nil)
		require.Nil(t, msg)
		require.ErrorIs(t, err, context.Canceled)

		require.Equal(t, uint32(9), rw.Credits())

		msg = rw.Prefetched()
		require.NotNil(t, msg)

		require.Equal(t, uint32(8), rw.Credits())

		msg = rw.Prefetched()
		require.Nil(t, msg)

		require.Equal(t, uint32(8), rw.Credits(), "no message returned, no credits used")
	})
}

func TestAMQPSenderWrapper(t *testing.T) {
	t.Run("errors are wrapped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		s := NewMockAMQPSenderCloser(ctrl)

		assertErr := func(err error, msg string) {
			t.Helper()
			var wrapErr Error

			require.ErrorAs(t, err, &wrapErr)
			require.EqualError(t, wrapErr, msg)
			require.Equal(t, uint64(101), wrapErr.ConnID)
			require.Equal(t, "sender", wrapErr.LinkName)
		}

		s.EXPECT().LinkName().Return("sender").AnyTimes()
		s.EXPECT().Send(test.NotCancelled, gomock.Any(), gomock.Any()).Return(errors.New("send failed"))

		s.EXPECT().Close(test.CancelledAndHasTimeout).Return(context.Canceled)
		s.EXPECT().Close(test.NotCancelledAndHasTimeout).Return(errors.New("close failed"))

		sw := &AMQPSenderWrapper{Inner: s, ContextWithTimeoutFn: test.NewContextWithTimeoutForTests, connID: 101}

		err := sw.Send(context.Background(), nil, nil)
		assertErr(err, "send failed")

		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		err = sw.Close(cancelledCtx)
		require.ErrorIs(t, err, context.Canceled)
		assertErr(err, "context canceled")

		err = sw.Close(context.Background())
		assertErr(err, "close failed")
	})

	t.Run("", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		s := NewMockAMQPSenderCloser(ctrl)

		s.EXPECT().MaxMessageSize().Return(uint64(99))

		sw := &AMQPSenderWrapper{Inner: s, ContextWithTimeoutFn: test.NewContextWithTimeoutForTests, connID: 101}
		require.Equal(t, uint64(99), sw.MaxMessageSize())
		require.Equal(t, uint64(101), sw.ConnID())
	})
}

func TestAMQPSessionWrapper(t *testing.T) {
	t.Run("ConnID is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sess := NewMockgoamqpSession(ctrl)

		sess.EXPECT().NewReceiver(gomock.Any(), gomock.Any(), gomock.Any()).Return(&amqp.Receiver{}, nil)
		sess.EXPECT().NewSender(gomock.Any(), gomock.Any(), gomock.Any()).Return(&amqp.Sender{}, nil)

		sessWrapper := &AMQPSessionWrapper{connID: uint64(101), Inner: sess, ContextWithTimeoutFn: context.WithTimeout}

		require.Equal(t, uint64(101), sessWrapper.ConnID())

		rc, err := sessWrapper.NewReceiver(context.Background(), "source", "1", nil)
		require.NoError(t, err)
		require.Equal(t, sessWrapper.ConnID(), rc.ConnID())

		sc, err := sessWrapper.NewSender(context.Background(), "target", "1", nil)
		require.NoError(t, err)
		require.Equal(t, sessWrapper.ConnID(), sc.ConnID())
	})

	t.Run("errors are wrapped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sess := NewMockgoamqpSession(ctrl)

		sess.EXPECT().NewReceiver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("new receiver failed"))
		sess.EXPECT().NewSender(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("new sender failed"))
		sess.EXPECT().Close(test.CancelledAndHasTimeout).Return(context.Canceled)

		sw := &AMQPSessionWrapper{
			connID:               uint64(101),
			Inner:                sess,
			ContextWithTimeoutFn: test.NewContextWithTimeoutForTests}

		assertErr := func(expectedPartitionID string, err error, msg string) {
			t.Helper()
			var wrapErr Error

			require.ErrorAs(t, err, &wrapErr)

			require.EqualError(t, wrapErr, msg)
			require.Equal(t, uint64(101), wrapErr.ConnID)
			require.Empty(t, wrapErr.LinkName)
			require.Equal(t, expectedPartitionID, wrapErr.PartitionID)
		}

		partitionID := "1"

		_, err := sw.NewReceiver(context.Background(), "source", partitionID, nil)
		assertErr(partitionID, err, "new receiver failed")

		_, err = sw.NewSender(context.Background(), "target", partitionID, nil)
		assertErr(partitionID, err, "new sender failed")

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err = sw.Close(ctx)
		assertErr("", err, "context canceled")
		require.ErrorIs(t, err, context.Canceled)
	})
}

func TestAMQPConnWrapper(t *testing.T) {
	t.Run("ConnID is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		innerConn := NewMockgoamqpConn(ctrl)

		innerConn.EXPECT().NewSession(gomock.Any(), gomock.Any()).Return(&amqp.Session{}, nil)

		cw := AMQPClientWrapper{
			ConnID: uint64(101),
			Inner:  innerConn,
		}

		sess, err := cw.NewSession(context.Background(), nil)
		require.NoError(t, err)

		require.Equal(t, uint64(101), sess.ConnID())
	})

	t.Run("errors are wrapped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		innerConn := NewMockgoamqpConn(ctrl)

		assertErr := func(err error, msg string) {
			t.Helper()
			var wrapErr Error
			require.ErrorAs(t, err, &wrapErr)
			require.EqualError(t, wrapErr, msg)
			require.Equal(t, uint64(101), wrapErr.ConnID)
			require.Empty(t, wrapErr.LinkName)
		}

		innerConn.EXPECT().NewSession(gomock.Any(), gomock.Any()).Return(nil, errors.New("new session failed"))
		innerConn.EXPECT().Close().Return(errors.New("close failed"))

		cw := AMQPClientWrapper{
			ConnID: uint64(101),
			Inner:  innerConn,
		}

		require.Equal(t, uint64(101), cw.ID())

		_, err := cw.NewSession(context.Background(), nil)
		assertErr(err, "new session failed")

		err = cw.Close()
		assertErr(err, "close failed")
	})
}
