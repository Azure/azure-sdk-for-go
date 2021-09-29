// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"
)

type fakeNS struct {
	claimNegotiated int
	MgmtClient      MgmtClient
	Session         AMQPSessionCloser
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

func (ns *fakeNS) NegotiateClaim(ctx context.Context, entityPath string) (func() <-chan struct{}, error) {
	ch := make(chan struct{})
	close(ch)

	ns.claimNegotiated++

	return func() <-chan struct{} {
		return ch
	}, nil
}

func (ns *fakeNS) GetEntityAudience(entityPath string) string {
	return fmt.Sprintf("audience: %s", entityPath)
}

func (ns *fakeNS) NewAMQPSession(ctx context.Context) (AMQPSessionCloser, error) {
	return ns.Session, nil
}

func (ns *fakeNS) NewMgmtClient(ctx context.Context, managementPath string) (MgmtClient, error) {
	return ns.MgmtClient, nil
}

type createLinkResponse struct {
	sender   AMQPSenderCloser
	receiver AMQPReceiverCloser
	err      error
}
