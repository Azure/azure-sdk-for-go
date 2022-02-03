// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

var retryOptionsOnlyOnce = utils.RetryOptions{
	MaxRetries: 0,
}

type fakeNetError struct {
	temp    bool
	timeout bool
}

func (pe fakeNetError) Timeout() bool   { return pe.timeout }
func (pe fakeNetError) Temporary() bool { return pe.temp }
func (pe fakeNetError) Error() string   { return "Fake but very permanent error" }

func assertFailedLinks(t *testing.T, lwid *LinksWithID, expectedErr error) {
	err := lwid.Sender.Send(context.TODO(), &amqp.Message{
		Data: [][]byte{
			{0},
		},
	})
	require.ErrorIs(t, err, expectedErr)

	_, err = PeekMessages(context.TODO(), lwid.RPC, 0, 1)
	require.ErrorIs(t, err, expectedErr)

	msg, err := lwid.Receiver.Receive(context.TODO())
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, msg)

}

func assertLinks(t *testing.T, lwid *LinksWithID) {
	err := lwid.Sender.Send(context.TODO(), &amqp.Message{
		Data: [][]byte{
			{0},
		},
	})
	require.NoError(t, err)

	_, err = PeekMessages(context.TODO(), lwid.RPC, 0, 1)
	require.NoError(t, err)

	require.NoError(t, lwid.Receiver.IssueCredit(1))
	msg, err := lwid.Receiver.Receive(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, msg)
}

func TestAMQPLinksBasic(t *testing.T) {
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	cs := test.GetConnectionString(t)
	ns, err := NewNamespace(NamespaceWithConnectionString(cs))
	require.NoError(t, err)

	links := NewAMQPLinks(ns, entityPath, func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
		return newLinksForAMQPLinksTest(entityPath, session)
	})

	lwr, err := links.Get(context.Background())
	require.NoError(t, err)

	assertLinks(t, lwr)

	require.EqualValues(t, entityPath, links.EntityPath())
}

func TestAMQPLinksLive(t *testing.T) {
	// we're not going to use this client for tehse tests.
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	cs := test.GetConnectionString(t)
	ns, err := NewNamespace(NamespaceWithConnectionString(cs))
	require.NoError(t, err)

	defer func() { _ = ns.Close(context.Background()) }()

	createLinksCalled := 0

	links := NewAMQPLinks(ns, entityPath, func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
		createLinksCalled++
		return newLinksForAMQPLinksTest(entityPath, session)
	})

	require.EqualValues(t, 0, createLinksCalled)
	require.NoError(t, links.RecoverIfNeeded(context.Background(), LinkID{}, amqp.ErrConnClosed))
	require.EqualValues(t, 1, createLinksCalled)

	lwr, err := links.Get(context.Background())
	require.NoError(t, err)

	amqpClient, clientRev, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 1, clientRev)
	require.NoError(t, amqpClient.Close())

	// all the links are dead because the connection is dead.
	assertFailedLinks(t, lwr, amqp.ErrConnClosed)

	// now we'll recover, which should recreate everything
	require.NoError(t, links.RecoverIfNeeded(context.Background(), lwr.ID, amqp.ErrConnClosed))
	require.EqualValues(t, 2, createLinksCalled)

	lwr, err = links.Get(context.Background())
	require.NoError(t, err)

	// should work now, connection should be reopened
	assertLinks(t, lwr)

	// cheat a bit and close the links out from under us (but leave them in place)
	actualLinks := links.(*AMQPLinksImpl)
	_ = actualLinks.Sender.Close(context.Background())
	_ = actualLinks.Receiver.Close(context.Background())
	_ = actualLinks.RPCLink.Close(context.Background())

	assertFailedLinks(t, lwr, amqp.ErrLinkClosed)

	lwr, err = links.Get(context.Background())
	require.NoError(t, err)

	require.NoError(t, links.RecoverIfNeeded(context.Background(), lwr.ID, amqp.ErrLinkClosed))
	require.EqualValues(t, 3, createLinksCalled)

	lwr, err = links.Get(context.Background())
	require.NoError(t, err)

	assertLinks(t, lwr)
}

func TestAMQPLinksLiveRecoverLink(t *testing.T) {
	// we're not going to use this client for tehse tests.
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	cs := test.GetConnectionString(t)
	ns, err := NewNamespace(NamespaceWithConnectionString(cs))
	require.NoError(t, err)

	defer func() { _ = ns.Close(context.Background()) }()

	createLinksCalled := 0

	links := NewAMQPLinks(ns, entityPath, func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
		createLinksCalled++
		return newLinksForAMQPLinksTest(entityPath, session)
	})

	require.EqualValues(t, 0, createLinksCalled)
	require.NoError(t, links.RecoverIfNeeded(context.Background(), LinkID{}, amqp.ErrConnClosed))
	require.EqualValues(t, 1, createLinksCalled)

	lwr, err := links.Get(context.Background())
	require.NoError(t, err)

	require.NoError(t, links.RecoverIfNeeded(context.Background(), lwr.ID, amqp.ErrLinkClosed))
	require.EqualValues(t, 2, createLinksCalled)
}

func TestAMQPLinksLiveRace(t *testing.T) {
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	cs := test.GetConnectionString(t)
	ns, err := NewNamespace(NamespaceWithConnectionString(cs))
	require.NoError(t, err)

	defer func() { _ = ns.Close(context.Background()) }()

	createLinksCalled := 0

	links := NewAMQPLinks(ns, entityPath, func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
		createLinksCalled++
		return newLinksForAMQPLinksTest(entityPath, session)
	})

	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := links.RecoverIfNeeded(context.Background(), LinkID{}, amqp.ErrConnClosed)
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	// TODO: also check that the connection hasn't recycled multiple times.
	require.EqualValues(t, 1, createLinksCalled)
}

func TestAMQPLinksLiveRaceLink(t *testing.T) {
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	cs := test.GetConnectionString(t)
	ns, err := NewNamespace(NamespaceWithConnectionString(cs))
	require.NoError(t, err)

	defer func() { _ = ns.Close(context.Background()) }()

	createLinksCalled := 0

	enableLogging()

	links := NewAMQPLinks(ns, entityPath, func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
		createLinksCalled++
		return newLinksForAMQPLinksTest(entityPath, session)
	})

	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := links.RecoverIfNeeded(context.Background(), LinkID{}, &amqp.DetachError{})
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	// TODO: also check that the connection hasn't recycled multiple times.
	require.EqualValues(t, 1, createLinksCalled)
}

func TestAMQPLinksRetry(t *testing.T) {
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	cs := test.GetConnectionString(t)
	ns, err := NewNamespace(NamespaceWithConnectionString(cs))
	require.NoError(t, err)

	defer func() { _ = ns.Close(context.Background()) }()

	createLinksCalled := 0

	links := NewAMQPLinks(ns, entityPath, func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
		createLinksCalled++
		return newLinksForAMQPLinksTest(entityPath, session)
	})

	err = links.Retry(context.Background(), "retryOp", func(ctx context.Context, lwid *LinksWithID, args *utils.RetryFnArgs) error {
		// force recoveries
		return &amqp.DetachError{}
	}, utils.RetryOptions{
		MaxRetries: 2,
		// note: omitting MaxRetries just to give a sanity check that
		// we do setDefaults() before we run.
		RetryDelay:    time.Millisecond,
		MaxRetryDelay: time.Millisecond,
	})

	var detachErr *amqp.DetachError
	require.ErrorAs(t, err, &detachErr)
	require.EqualValues(t, 3, createLinksCalled)
}

func TestAMQPLinksMultipleWithSameConnection(t *testing.T) {
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	cs := test.GetConnectionString(t)
	ns, err := NewNamespace(NamespaceWithConnectionString(cs))
	require.NoError(t, err)

	defer func() { _ = ns.Close(context.Background()) }()

	createLinksCalled := 0

	links := NewAMQPLinks(ns, entityPath, func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
		createLinksCalled++
		return newLinksForAMQPLinksTest(entityPath, session)
	})

	createLinksCalled2 := 0

	links2 := NewAMQPLinks(ns, entityPath, func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
		createLinksCalled2++
		return newLinksForAMQPLinksTest(entityPath, session)
	})

	wg := sync.WaitGroup{}

	lwr, err := links.Get(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 1, createLinksCalled)
	require.EqualValues(t, 1, lwr.ID.Link)

	lwr2, err := links2.Get(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 1, createLinksCalled2)
	require.EqualValues(t, 1, lwr2.ID.Link)

	wg.Add(1)

	go func() {
		defer wg.Done()
		err = links.RecoverIfNeeded(context.Background(), lwr.ID, &amqp.DetachError{})
		require.NoError(t, err)
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		err = links2.RecoverIfNeeded(context.Background(), lwr2.ID, &amqp.DetachError{})
		require.NoError(t, err)
	}()

	wg.Wait()

	// TODO: also check that the connection hasn't recycled multiple times.
	require.EqualValues(t, 2, createLinksCalled)
	require.EqualValues(t, 2, createLinksCalled2)

	_, clientRev, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 1, clientRev)

	recovered, err := ns.Recover(context.Background(), clientRev)
	require.NoError(t, err)
	require.True(t, recovered)

	_, clientRev, err = ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 2, clientRev)

	// now attempt a recover but with an older revision (won't do anything since we've
	// already recovered past that older rev. They should just call Get())
	recovered, err = ns.Recover(context.Background(), clientRev-1)
	require.NoError(t, err)
	require.False(t, recovered)

	_, clientRev, err = ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 2, clientRev)
}

func newLinksForAMQPLinksTest(entityPath string, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
	receiveMode := amqp.ModeSecond

	opts := []amqp.LinkOption{
		amqp.LinkSourceAddress(entityPath),
		amqp.LinkReceiverSettle(receiveMode),
		amqp.LinkWithManualCredits(),
		amqp.LinkCredit(1000),
	}

	receiver, err := session.NewReceiver(opts...)

	if err != nil {
		return nil, nil, err
	}

	sender, err := session.NewSender(
		amqp.LinkSenderSettle(amqp.ModeMixed),
		amqp.LinkReceiverSettle(amqp.ModeFirst),
		amqp.LinkTargetAddress(entityPath))

	if err != nil {
		_ = receiver.Close(context.Background())
		return nil, nil, err
	}

	return sender, receiver, nil
}

func enableLogging() {
	azlog.SetListener(func(e azlog.Event, s string) {
		log.Printf("%s %s", e, s)
	})
}
