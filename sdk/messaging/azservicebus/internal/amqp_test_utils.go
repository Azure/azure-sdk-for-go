// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"

	"github.com/Azure/azure-amqp-common-go/v3/rpc"
)

type FakeNS struct {
	claimNegotiated int
	recovered       uint64
	clientRevisions []uint64
	MgmtClient      MgmtClient
	RPCLink         *rpc.Link
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

type fakeMgmtClient struct {
	MgmtClient
	closed int
}

type FakeAMQPLinks struct {
	AMQPLinks

	Closed int

	// values to be returned for each `Get` call
	Revision uint64
	Receiver AMQPReceiver
	Sender   AMQPSender
	Mgmt     MgmtClient
	Err      error

	permanently bool
}

type FakeAMQPReceiver struct {
	AMQPReceiver
	Closed int
	Drain  int
}

func (r *FakeAMQPReceiver) DrainCredit(ctx context.Context) error {
	r.Drain++
	return nil
}

func (r *FakeAMQPReceiver) Close(ctx context.Context) error {
	r.Closed++
	return nil
}

func (l *FakeAMQPLinks) Get(ctx context.Context) (AMQPSender, AMQPReceiver, MgmtClient, uint64, error) {
	return l.Sender, l.Receiver, l.Mgmt, l.Revision, l.Err
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

func (m *fakeMgmtClient) Close(ctx context.Context) error {
	m.closed++
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

func (ns *FakeNS) NewMgmtClient(ctx context.Context, links AMQPLinks) (MgmtClient, error) {
	return ns.MgmtClient, nil
}

func (ns *FakeNS) NewRPCLink(ctx context.Context, managementPath string) (*rpc.Link, error) {
	return ns.RPCLink, nil
}

func (ns *FakeNS) Recover(ctx context.Context, clientRevision uint64) error {
	ns.clientRevisions = append(ns.clientRevisions, clientRevision)
	ns.recovered++
	return nil
}

func (ns *FakeNS) NewAMQPLinks(entityPath string, createLinkFunc CreateLinkFunc) AMQPLinks {
	return ns.AMQPLinks
}

type createLinkResponse struct {
	sender   AMQPSenderCloser
	receiver AMQPReceiverCloser
	err      error
}
