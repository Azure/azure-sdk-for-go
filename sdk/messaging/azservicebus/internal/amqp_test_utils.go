// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/Azure/go-amqp"
)

type FakeNS struct {
	claimNegotiated int
	recovered       uint64
	clientRevisions []uint64
	RPCLink         RPCLink
	Session         AMQPSessionCloser
	AMQPLinks       *FakeAMQPLinks
}

type FakeAMQPSender struct {
	Closed int
	AMQPSender
}

type FakeAMQPSession struct {
	AMQPSessionCloser
	closed int
}

type FakeAMQPLinks struct {
	AMQPLinks

	Closed int

	// values to be returned for each `Get` call
	Revision LinkID
	Receiver AMQPReceiver
	Sender   AMQPSender
	RPC      RPCLink
	Err      error

	permanently bool
}

type FakeAMQPReceiver struct {
	AMQPReceiver
	Closed           int
	Drain            int
	RequestedCredits uint32

	ReceiveResults chan struct {
		M *amqp.Message
		E error
	}
}

func (r *FakeAMQPReceiver) IssueCredit(credit uint32) error {
	r.RequestedCredits += credit
	return nil
}

func (r *FakeAMQPReceiver) DrainCredit(ctx context.Context) error {
	r.Drain++

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (r *FakeAMQPReceiver) Receive(ctx context.Context) (*amqp.Message, error) {
	select {
	case res := <-r.ReceiveResults:
		return res.M, res.E
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *FakeAMQPReceiver) Prefetched(ctx context.Context) (*amqp.Message, error) {
	select {
	case res := <-r.ReceiveResults:
		return res.M, res.E
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, nil
	}
}

func (r *FakeAMQPReceiver) Close(ctx context.Context) error {
	r.Closed++
	return nil
}

func (l *FakeAMQPLinks) Get(ctx context.Context) (*LinksWithID, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &LinksWithID{
			Sender:   l.Sender,
			Receiver: l.Receiver,
			RPC:      l.RPC,
			ID:       l.Revision,
		}, l.Err
	}
}

func (l *FakeAMQPLinks) Retry(ctx context.Context, name string, fn RetryWithLinksFn, o utils.RetryOptions) error {
	lwr, err := l.Get(ctx)

	if err != nil {
		return err
	}

	return fn(ctx, lwr, &utils.RetryFnArgs{})
}

func (l *FakeAMQPLinks) Close(ctx context.Context, permanently bool) error {
	if permanently {
		l.permanently = true
	}

	l.Closed++
	return nil
}

func (l *FakeAMQPLinks) ClosedPermanently() bool {
	return l.permanently
}

func (s *FakeAMQPSender) Close(ctx context.Context) error {
	s.Closed++
	return nil
}

func (s *FakeAMQPSession) Close(ctx context.Context) error {
	s.closed++
	return nil
}

func (ns *FakeNS) NegotiateClaim(ctx context.Context, entityPath string) (func() <-chan struct{}, error) {
	ch := make(chan struct{})
	close(ch)

	ns.claimNegotiated++

	return func() <-chan struct{} {
		return ch
	}, nil
}

func (ns *FakeNS) GetEntityAudience(entityPath string) string {
	return fmt.Sprintf("audience: %s", entityPath)
}

func (ns *FakeNS) NewAMQPSession(ctx context.Context) (AMQPSessionCloser, uint64, error) {
	return ns.Session, ns.recovered + 100, nil
}

func (ns *FakeNS) NewRPCLink(ctx context.Context, managementPath string) (RPCLink, error) {
	return ns.RPCLink, nil
}

func (ns *FakeNS) Recover(ctx context.Context, clientRevision uint64) (bool, error) {
	ns.clientRevisions = append(ns.clientRevisions, clientRevision)
	ns.recovered++
	return true, nil
}

func (ns *FakeNS) CloseIfNeeded(ctx context.Context, clientRevision uint64) error {
	return nil
}

func (ns *FakeNS) NewAMQPLinks(entityPath string, createLinkFunc CreateLinkFunc) AMQPLinks {
	return ns.AMQPLinks
}
