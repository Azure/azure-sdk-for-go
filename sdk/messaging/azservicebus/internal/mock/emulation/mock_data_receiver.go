// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"context"
	"fmt"
	"math"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock"
	"github.com/golang/mock/gomock"
)

type MockReceiver struct {
	*mock.MockAMQPReceiverCloser
	Opts          *amqp.ReceiverOptions
	Session       *MockSession
	Source        string
	Status        *Status
	TargetAddress string

	// InternalReceive will receive from our default mock. Useful if you want to
	// change the default EXPECT() for AMQPReceiver.Receive().
	InternalReceive func(ctx context.Context) (*amqp.Message, error)

	// InternalIssueCredit will issue credit for our default mock. Useful if you
	// want to change the default EXPECT() for AMQPReceiver.IssueCredit().
	InternalIssueCredit func(credit uint32) error
}

func (rcvr *MockReceiver) Done() <-chan error {
	errCh := make(chan error)

	go func() {
		select {
		case <-rcvr.Session.Done():
			errCh <- rcvr.Session.Status.Err()
		case <-rcvr.Status.Done():
			errCh <- rcvr.Status.Err()
		}
	}()

	return errCh
}

func (rcvr *MockReceiver) LinkEvent() LinkEvent {
	return LinkEvent{
		ConnID:        rcvr.Session.Conn.Name(),
		SessID:        rcvr.Session.ID,
		Entity:        rcvr.Source,
		Name:          rcvr.LinkName(),
		Role:          LinkRoleReceiver,
		TargetAddress: rcvr.TargetAddress,
	}
}

func (md *MockData) NewReceiver(ctx context.Context, source string, opts *amqp.ReceiverOptions, sess *MockSession) (amqpwrap.AMQPReceiverCloser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-sess.Status.Done():
		return nil, sess.Status.Err()
	default:
	}

	if opts == nil {
		opts = &amqp.ReceiverOptions{}
	}

	rcvr := &MockReceiver{
		MockAMQPReceiverCloser: mock.NewMockAMQPReceiverCloser(md.Ctrl),
		Session:                sess,
		Source:                 source,
		Status:                 NewStatus(sess.Status),
		TargetAddress:          opts.TargetAddress,
		Opts:                   opts,
	}

	id := fmt.Sprintf("%s|%s|%s|e:%s", sess.Conn.Name(), sess.ID, md.nextUniqueName("r"), source)
	rcvr.EXPECT().LinkName().Return(id).AnyTimes()

	md.Events.OpenLink(rcvr.LinkEvent())

	md.mocksMu.Lock()
	md.receivers[source] = append(md.receivers[source], rcvr)
	md.mocksMu.Unlock()

	var credits uint32
	var q *Queue

	rcvr.InternalReceive = func(ctx context.Context) (*amqp.Message, error) {
		m, err := q.Receive(ctx, rcvr.LinkEvent(), rcvr.Status)

		if err != nil {
			return nil, err
		}

		credits--
		return m, nil
	}

	rcvr.InternalIssueCredit = func(credit uint32) error {
		credits += credit
		return q.IssueCredit(credit, rcvr.LinkEvent(), rcvr.Status)
	}

	if err := md.options.PreReceiverMock(rcvr, ctx); err != nil {
		return nil, err
	}

	if source == "$cbs" {
		q = md.upsertQueue(opts.TargetAddress)

		md.cbsRouterOnce.Do(func() {
			cbs := md.upsertQueue(source)
			go func() { md.cbsRouter(md.cbsContext, cbs, md.getQueue) }()
		})
	} else {
		q = md.upsertQueue(source)
	}

	rcvr.EXPECT().Receive(gomock.Any()).DoAndReturn(rcvr.InternalReceive).AnyTimes()

	rcvr.EXPECT().Close(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
		md.Events.CloseLink(rcvr.LinkEvent())

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sess.Status.Done():
			return sess.Status.Err()
		default:
			rcvr.Status.CloseWithError(amqp.ErrLinkClosed)
		}

		return nil
	}).AnyTimes()

	rcvr.EXPECT().Credits().DoAndReturn(func() uint32 {
		return credits
	}).AnyTimes()

	rcvr.EXPECT().AcceptMessage(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, msg *amqp.Message) error {
		return q.AcceptMessage(ctx, msg, rcvr.LinkEvent(), rcvr.Status)
	}).AnyTimes()

	rcvr.EXPECT().RejectMessage(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, msg *amqp.Message, e *amqp.Error) error {
		return q.RejectMessage(ctx, msg, e, rcvr.LinkEvent(), rcvr.Status)
	}).AnyTimes()

	rcvr.EXPECT().ReleaseMessage(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, msg *amqp.Message) error {
		return q.ReleaseMessage(ctx, msg, rcvr.LinkEvent())
	}).AnyTimes()

	rcvr.EXPECT().ModifyMessage(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, msg *amqp.Message, options *amqp.ModifyMessageOptions) error {
		return q.ModifyMessage(ctx, msg, options, rcvr.LinkEvent())
	}).AnyTimes()

	rcvr.EXPECT().Prefetched().Return((*amqp.Message)(nil)).AnyTimes()

	if opts.ManualCredits {
		rcvr.EXPECT().IssueCredit(gomock.Any()).DoAndReturn(rcvr.InternalIssueCredit).AnyTimes()
	} else {
		// assume unlimited credits for this receiver - the AMQP stack is going to take care of replenishing credits.
		_ = q.IssueCredit(math.MaxUint32, rcvr.LinkEvent(), rcvr.Status)
	}

	return rcvr, nil
}
