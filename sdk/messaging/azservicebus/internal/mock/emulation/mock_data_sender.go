// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock"
	"github.com/golang/mock/gomock"
)

type MockSender struct {
	*mock.MockAMQPSenderCloser
	Opts    *amqp.SenderOptions
	Session *MockSession
	Status  *Status
	Target  string
}

func (m *MockSender) LinkEvent() LinkEvent {
	return LinkEvent{
		ConnID: m.Session.Conn.Name(),
		SessID: m.Session.ID,
		Entity: m.Target,
		Name:   m.LinkName(),
		Role:   LinkRoleSender,
	}
}

func (md *MockData) NewSender(ctx context.Context, target string, opts *amqp.SenderOptions, sess *MockSession) (amqpwrap.AMQPSenderCloser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-sess.Status.Done():
		return nil, sess.Status.Err()
	default:
	}

	sender := &MockSender{
		MockAMQPSenderCloser: mock.NewMockAMQPSenderCloser(md.Ctrl),
		Session:              sess,
		Status:               NewStatus(sess.Status),
		Opts:                 opts,
		Target:               target,
	}

	id := fmt.Sprintf("%s|%s|%s|e:%s", sess.Conn.Name(), sess.ID, md.nextUniqueName("s"), target)
	sender.EXPECT().LinkName().Return(id).AnyTimes()

	md.Events.OpenLink(sender.LinkEvent())

	md.mocksMu.Lock()
	md.senders[target] = append(md.senders[target], sender)
	md.mocksMu.Unlock()

	if err := md.options.PreSenderMock(sender, ctx); err != nil {
		return nil, err
	}

	sender.EXPECT().MaxMessageSize().Return(uint64(1024 * 1024 * 100)).AnyTimes()

	// this should work fine even for RPC links like $cbs or $management
	q := md.upsertQueue(target)
	sender.EXPECT().Send(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, msg *amqp.Message) error {
		return q.Send(ctx, msg, sender.LinkEvent(), sender.Status)
	}).AnyTimes()

	sender.EXPECT().Close(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
		md.Events.CloseLink(sender.LinkEvent())

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sess.Status.Done():
			return sess.Status.Err()
		default:
			sender.Status.CloseWithError(amqp.ErrLinkClosed)
		}

		return nil
	}).AnyTimes()

	return sender, nil
}
