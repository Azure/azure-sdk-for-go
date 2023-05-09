// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock"
	"github.com/Azure/go-amqp"
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

	var createdEntities []interface {
		Close(ctx context.Context) error
	}

	sess.EXPECT().NewReceiver(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, source string, opts *amqp.ReceiverOptions) (amqpwrap.AMQPReceiverCloser, error) {
		receiver, err := md.NewReceiver(ctx, source, opts, sess)

		if err == nil {
			createdEntities = append(createdEntities, receiver)
		}

		return receiver, err
	}).AnyTimes()

	sess.EXPECT().NewSender(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, target string, opts *amqp.SenderOptions) (amqpwrap.AMQPSenderCloser, error) {
		sender, err := md.NewSender(ctx, target, opts, sess)

		if err == nil {
			createdEntities = append(createdEntities, sender)
		}

		return sender, err
	}).AnyTimes()

	sess.EXPECT().Close(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
		// if the receiver or sender are still open close them as well - this mimics the expected state
		// after ending a session.
		for _, e := range createdEntities {
			_ = e.Close(context.Background())
		}

		select {
		case <-conn.Status.Done():
			return conn.Status.Err()
		default:
			sess.Status.CloseWithError(&amqp.SessionError{})
			return nil
		}
	}).AnyTimes()

	return sess, nil
}
