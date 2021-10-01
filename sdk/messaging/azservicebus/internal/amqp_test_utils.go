// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"
)

type FakeNS struct {
	claimNegotiated int
	recovered       uint64
	clientRevisions []uint64
	MgmtClient      MgmtClient
	Session         AMQPSessionCloser
	AMQPLinks       *FakeAMQPLinks
}

type fakeAMQPSender struct {
	closed int
	AMQPSender
}

type fakeAMQPSession struct {
	AMQPSessionCloser
	closed int
}

type fakeMgmtClient struct {
	MgmtClient
	closed int
}

type FakeAMQPLinks struct {
	AMQPLinks

	// values to be returned for each `Get` call
	Revision uint64
	Receiver AMQPReceiver
	Sender   AMQPSender
	Mgmt     MgmtClient
	Err      error
}

func (l FakeAMQPLinks) Get(ctx context.Context) (AMQPSender, AMQPReceiver, MgmtClient, uint64, error) {
	return l.Sender, l.Receiver, l.Mgmt, l.Revision, l.Err
}

func (s *fakeAMQPSender) Close(ctx context.Context) error {
	s.closed++
	return nil
}

func (s *fakeAMQPSession) Close(ctx context.Context) error {
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
