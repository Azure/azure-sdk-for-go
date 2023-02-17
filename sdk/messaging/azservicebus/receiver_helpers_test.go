// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
)

type StubAMQPReceiver struct {
	stubClose                       func(inner amqpwrap.AMQPReceiverCloser, ctx context.Context) error
	stubCloseCalled                 int
	stubIssueCredit                 func(inner amqpwrap.AMQPReceiverCloser, credit uint32) error
	stubIssueCreditCalled           int
	stubCredits                     func(inner amqpwrap.AMQPReceiverCloser) uint32
	stubCreditsCalled               int
	stubReceive                     func(inner amqpwrap.AMQPReceiverCloser, ctx context.Context) (*amqp.Message, error)
	stubReceiveCalled               int
	stubPrefetched                  func(inner amqpwrap.AMQPReceiverCloser) *amqp.Message
	stubPrefetchedCalled            int
	stubAcceptMessage               func(inner amqpwrap.AMQPReceiverCloser, ctx context.Context, msg *amqp.Message) error
	stubAcceptMessageCalled         int
	stubRejectMessage               func(inner amqpwrap.AMQPReceiverCloser, ctx context.Context, msg *amqp.Message, e *amqp.Error) error
	stubRejectMessageCalled         int
	stubReleaseMessage              func(inner amqpwrap.AMQPReceiverCloser, ctx context.Context, msg *amqp.Message) error
	stubReleaseMessageCalled        int
	stubModifyMessage               func(inner amqpwrap.AMQPReceiverCloser, ctx context.Context, msg *amqp.Message, options *amqp.ModifyMessageOptions) error
	stubModifyMessageCalled         int
	stubLinkName                    func(inner amqpwrap.AMQPReceiverCloser) string
	stubLinkNameCalled              int
	stubLinkSourceFilterValue       func(inner amqpwrap.AMQPReceiverCloser, name string) interface{}
	stubLinkSourceFilterValueCalled int
	inner                           amqpwrap.AMQPReceiverCloser
}

func (r *StubAMQPReceiver) Close(ctx context.Context) error {
	r.stubCloseCalled++
	if r.stubClose != nil {
		return r.stubClose(r.inner, ctx)
	}
	return r.inner.Close(ctx)
}

func (r *StubAMQPReceiver) IssueCredit(credit uint32) error {
	r.stubIssueCreditCalled++
	if r.stubIssueCredit != nil {
		return r.stubIssueCredit(r.inner, credit)
	}
	return r.inner.IssueCredit(credit)
}

func (r *StubAMQPReceiver) Receive(ctx context.Context) (*amqp.Message, error) {
	r.stubReceiveCalled++
	if r.stubReceive != nil {
		return r.stubReceive(r.inner, ctx)
	}
	return r.inner.Receive(ctx)
}

func (r *StubAMQPReceiver) Prefetched() *amqp.Message {
	r.stubPrefetchedCalled++
	if r.stubPrefetched != nil {
		return r.stubPrefetched(r.inner)
	}

	return r.inner.Prefetched()
}

func (r *StubAMQPReceiver) AcceptMessage(ctx context.Context, msg *amqp.Message) error {
	r.stubAcceptMessageCalled++
	if r.stubAcceptMessage != nil {
		return r.stubAcceptMessage(r.inner, ctx, msg)
	}
	return r.inner.AcceptMessage(ctx, msg)
}

func (r *StubAMQPReceiver) RejectMessage(ctx context.Context, msg *amqp.Message, e *amqp.Error) error {
	r.stubRejectMessageCalled++
	if r.stubRejectMessage != nil {
		return r.stubRejectMessage(r.inner, ctx, msg, e)
	}
	return r.inner.RejectMessage(ctx, msg, e)
}

func (r *StubAMQPReceiver) ReleaseMessage(ctx context.Context, msg *amqp.Message) error {
	r.stubReleaseMessageCalled++
	if r.stubReleaseMessage != nil {
		return r.stubReleaseMessage(r.inner, ctx, msg)
	}
	return r.inner.ReleaseMessage(ctx, msg)
}

func (r *StubAMQPReceiver) ModifyMessage(ctx context.Context, msg *amqp.Message, options *amqp.ModifyMessageOptions) error {
	r.stubModifyMessageCalled++
	if r.stubModifyMessage != nil {
		return r.stubModifyMessage(r.inner, ctx, msg, options)
	}
	return r.inner.ModifyMessage(ctx, msg, options)
}

func (r *StubAMQPReceiver) LinkName() string {
	r.stubLinkNameCalled++
	if r.stubLinkName != nil {
		return r.stubLinkName(r.inner)
	}
	return r.inner.LinkName()
}

func (r *StubAMQPReceiver) LinkSourceFilterValue(name string) interface{} {
	r.stubLinkSourceFilterValueCalled++
	if r.stubLinkSourceFilterValue != nil {
		return r.stubLinkSourceFilterValue(r.inner, name)
	}
	return r.inner.LinkSourceFilterValue(name)
}

func (r *StubAMQPReceiver) Credits() uint32 {
	r.stubCreditsCalled++
	if r.stubCredits != nil {
		return r.stubCredits(r.inner)
	}
	return r.inner.Credits()
}
