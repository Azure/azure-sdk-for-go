// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
)

type FakeNSForPartClient struct {
	NamespaceForAMQPLinks
	Receiver *FakeAMQPReceiver
}

type FakeAMQPSession struct {
	amqpwrap.AMQPSession
	NS *FakeNSForPartClient
}

type FakeAMQPReceiver struct {
	amqpwrap.AMQPReceiverCloser

	// ActiveCredits are incremented and decremented by IssueCredit and Receive.
	ActiveCredits uint32

	// IssuedCredit just accumulates, so we can get an idea of how many credits we issued overall.
	IssuedCredit []uint32

	// CreditsSetFromOptions is similar to issuedCredit, but only tracks credits added in via the LinkOptions.Credit
	// field (ie, enabling prefetch).
	CreditsSetFromOptions uint32

	// ManualCreditsSetFromOptions is the value of the LinkOptions.ManualCredits value.
	ManualCreditsSetFromOptions bool

	Messages []*amqp.Message
}

func (ns *FakeNSForPartClient) NegotiateClaim(ctx context.Context, entityPath string) (context.CancelFunc, <-chan struct{}, error) {
	ctx, cancel := context.WithCancel(ctx)
	return cancel, ctx.Done(), nil
}

func (ns *FakeNSForPartClient) NewAMQPSession(ctx context.Context) (amqpwrap.AMQPSession, uint64, error) {
	return &FakeAMQPSession{
		NS: ns,
	}, 1, nil
}

func (sess *FakeAMQPSession) NewReceiver(ctx context.Context, source string, opts *amqp.ReceiverOptions) (amqpwrap.AMQPReceiverCloser, error) {
	sess.NS.Receiver.ManualCreditsSetFromOptions = opts.ManualCredits
	sess.NS.Receiver.CreditsSetFromOptions = opts.Credit

	if !opts.ManualCredits {
		sess.NS.Receiver.ActiveCredits = opts.Credit
	}

	return sess.NS.Receiver, nil
}

func (r *FakeAMQPReceiver) Credits() uint32 {
	return r.ActiveCredits
}

func (r *FakeAMQPReceiver) IssueCredit(credit uint32) error {
	r.ActiveCredits += credit
	r.IssuedCredit = append(r.IssuedCredit, credit)
	return nil
}

func (r *FakeAMQPReceiver) Receive(ctx context.Context) (*amqp.Message, error) {
	if len(r.Messages) > 0 {
		r.ActiveCredits--
		m := r.Messages[0]
		r.Messages = r.Messages[1:]
		return m, nil
	}

	return nil, nil
}
