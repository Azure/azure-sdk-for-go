// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/stretchr/testify/require"
)

type StubAMQPReceiver struct {
	stubClose                       func(inner internal.AMQPReceiverCloser, ctx context.Context) error
	stubCloseCalled                 int
	stubIssueCredit                 func(inner internal.AMQPReceiverCloser, credit uint32) error
	stubIssueCreditCalled           int
	stubDrainCredit                 func(inner internal.AMQPReceiverCloser, ctx context.Context) error
	stubDrainCreditCalled           int
	stubReceive                     func(inner internal.AMQPReceiverCloser, ctx context.Context) (*amqp.Message, error)
	stubReceiveCalled               int
	stubPrefetched                  func(inner internal.AMQPReceiverCloser) *amqp.Message
	stubPrefetchedCalled            int
	stubAcceptMessage               func(inner internal.AMQPReceiverCloser, ctx context.Context, msg *amqp.Message) error
	stubAcceptMessageCalled         int
	stubRejectMessage               func(inner internal.AMQPReceiverCloser, ctx context.Context, msg *amqp.Message, e *amqp.Error) error
	stubRejectMessageCalled         int
	stubReleaseMessage              func(inner internal.AMQPReceiverCloser, ctx context.Context, msg *amqp.Message) error
	stubReleaseMessageCalled        int
	stubModifyMessage               func(inner internal.AMQPReceiverCloser, ctx context.Context, msg *amqp.Message, options *amqp.ModifyMessageOptions) error
	stubModifyMessageCalled         int
	stubLinkName                    func(inner internal.AMQPReceiverCloser) string
	stubLinkNameCalled              int
	stubLinkSourceFilterValue       func(inner internal.AMQPReceiverCloser, name string) interface{}
	stubLinkSourceFilterValueCalled int
	inner                           internal.AMQPReceiverCloser
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

func (r *StubAMQPReceiver) DrainCredit(ctx context.Context) error {
	r.stubDrainCreditCalled++
	if r.stubDrainCredit != nil {
		return r.stubDrainCredit(r.inner, ctx)
	}
	return r.inner.DrainCredit(ctx)
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

// addStub initializes the links and wraps the internal AMQPReceiver with a stub.
//  The stub allows you to either forward calls to the actual underlying AMQPReceiver instance
// or take it over.
func addStub(t *testing.T, receiver *Receiver, stub *StubAMQPReceiver) *StubAMQPReceiver {
	actualLinks := receiver.amqpLinks.(*internal.AMQPLinksImpl)

	// make sure the links are live
	_, err := actualLinks.Get(context.Background())
	require.NoError(t, err)

	stub.inner = actualLinks.Receiver
	actualLinks.Receiver = stub

	return stub
}
