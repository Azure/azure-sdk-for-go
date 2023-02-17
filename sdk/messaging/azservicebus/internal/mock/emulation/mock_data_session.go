// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock"
	"github.com/golang/mock/gomock"
)

type MockSession struct {
	Status *Status
	*mock.MockAMQPSession

	ID   string
	Conn *MockConnection
}

func (sess *MockSession) Done() <-chan error {
	errCh := make(chan error)

	go func() {
		select {
		case <-sess.Conn.Status.Done():
			errCh <- sess.Conn.Status.Err()
		case <-sess.Status.Done():
			errCh <- sess.Status.Err()
		}
	}()

	return errCh
}

func (md *MockData) newSession(ctx context.Context, opts *amqp.SessionOptions, conn *MockConnection) (amqpwrap.AMQPSession, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-conn.Status.Done():
		return nil, conn.Status.Err()
	default:
	}

	sess := &MockSession{
		Status:          NewStatus(conn.Status),
		ID:              md.nextUniqueName("sess"),
		Conn:            conn,
		MockAMQPSession: mock.NewMockAMQPSession(md.Ctrl),
	}

	if err := md.options.PreSessionMock(sess, ctx, opts); err != nil {
		return nil, err
	}

	sess.EXPECT().NewReceiver(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, source string, opts *amqp.ReceiverOptions) (amqpwrap.AMQPReceiverCloser, error) {
		return md.NewReceiver(ctx, source, opts, sess)
	}).AnyTimes()

	sess.EXPECT().NewSender(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, target string, opts *amqp.SenderOptions) (amqpwrap.AMQPSenderCloser, error) {
		return md.NewSender(ctx, target, opts, sess)
	}).AnyTimes()

	sess.EXPECT().Close(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
		select {
		case <-conn.Status.Done():
			return conn.Status.Err()
		default:
			sess.Status.CloseWithError(amqp.ErrSessionClosed)
			return nil
		}
	}).AnyTimes()

	return sess, nil
}
