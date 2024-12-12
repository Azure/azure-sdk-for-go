// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

var retryOptionsOnlyOnce = exported.RetryOptions{
	MaxRetries: 0,
}

type fakeNetError struct {
	temp    bool
	timeout bool
}

func (pe fakeNetError) Timeout() bool   { return pe.timeout }
func (pe fakeNetError) Temporary() bool { return pe.temp }
func (pe fakeNetError) Error() string   { return "Fake but very permanent error" }

func assertFailedLinks[T error, T2 error](t *testing.T, lwid *LinksWithID, expectedErr T, expectedRPCError T2) {
	err := lwid.Sender.Send(context.TODO(), &amqp.Message{
		Data: [][]byte{
			{0},
		},
	}, nil)

	require.True(t, errors.Is(err, expectedErr) || errors.As(err, &expectedErr))
	require.ErrorIs(t, err, expectedErr)

	_, err = PeekMessages(context.TODO(), lwid.RPC, lwid.Receiver.LinkName(), 0, 1)
	require.True(t, errors.Is(err, expectedRPCError) || errors.As(err, &expectedRPCError))

	msg, err := lwid.Receiver.Receive(context.TODO(), nil)
	require.True(t, errors.Is(err, expectedErr) || errors.As(err, &expectedErr))
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, msg)

}

func assertLinks(t *testing.T, lwid *LinksWithID) {
	err := lwid.Sender.Send(context.TODO(), &amqp.Message{
		Data: [][]byte{
			{0},
		},
	}, nil)
	require.NoError(t, err)

	_, err = PeekMessages(context.TODO(), lwid.RPC, lwid.Receiver.LinkName(), 0, 1)
	require.NoError(t, err)

	require.NoError(t, lwid.Receiver.IssueCredit(1))
	msg, err := lwid.Receiver.Receive(context.TODO(), nil)
	require.NoError(t, err)
	require.NotNil(t, msg)
}

func newNamespaceForTest(t *testing.T) *Namespace {
	iv := test.GetIdentityVars(t)
	ns, err := NewNamespace(NamespaceWithTokenCredential(iv.Endpoint, iv.Cred))
	require.NoError(t, err)

	return ns
}

func TestAMQPLinksBasic(t *testing.T) {
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	ns := newNamespaceForTest(t)

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: entityPath,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			return newLinksForAMQPLinksTest(entityPath, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer func() {
		err := links.Close(context.Background(), true)
		require.NoError(t, err)
	}()

	lwr, err := links.Get(context.Background())
	require.NoError(t, err)

	assertLinks(t, lwr)
}

func TestAMQPLinksLiveCloseConnectionUnexpectedly(t *testing.T) {
	// we're not going to use this client for these tests.
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	ns := newNamespaceForTest(t)

	defer func() { _ = ns.Close(false) }()

	createLinksCalled := 0

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: entityPath,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			createLinksCalled++
			return newLinksForAMQPLinksTest(entityPath, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer func() {
		err := links.Close(context.Background(), true)
		require.NoError(t, err)
	}()

	require.EqualValues(t, 0, createLinksCalled)
	require.NoError(t, links.RecoverIfNeeded(context.Background(), LinkID{}, &amqp.ConnError{}))
	require.EqualValues(t, 1, createLinksCalled)

	lwr, err := links.Get(context.Background())
	require.NoError(t, err)

	amqpClient, clientRev, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 1, clientRev)
	require.NoError(t, amqpClient.Close())

	// all the links are dead because the connection is dead.
	assertFailedLinks(t, lwr, &amqp.ConnError{}, &amqp.ConnError{})

	// now we'll recover, which should recreate everything
	require.NoError(t, links.RecoverIfNeeded(context.Background(), lwr.ID, &amqp.ConnError{}))
	require.EqualValues(t, 2, createLinksCalled)

	lwr, err = links.Get(context.Background())
	require.NoError(t, err)

	// should work now, connection should be reopened
	assertLinks(t, lwr)

	lwr, err = links.Get(context.Background())
	require.NoError(t, err)

	require.NoError(t, links.RecoverIfNeeded(context.Background(), lwr.ID, &amqp.LinkError{}))
	require.EqualValues(t, 3, createLinksCalled)

	lwr, err = links.Get(context.Background())
	require.NoError(t, err)

	assertLinks(t, lwr)
}

func TestAMQPLinksLiveCloseLinksUnexpectedly(t *testing.T) {
	ns := newNamespaceForTest(t)

	defer func() { _ = ns.Close(false) }()

	createLinksCalled := 0

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: test.BuiltInTestQueue,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			createLinksCalled++
			return newLinksForAMQPLinksTest(test.BuiltInTestQueue, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer func() {
		err := links.Close(context.Background(), true)
		require.NoError(t, err)
	}()

	oldLWR, err := links.Get(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 1, createLinksCalled)

	// shut down all the links as if they were killed on their own.
	actualLinks := links.(*AMQPLinksImpl)
	_ = actualLinks.Sender.Close(context.Background())
	_ = actualLinks.Receiver.Close(context.Background())
	_ = actualLinks.RPCLink.Close(context.Background())

	_, clientRev, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 1, clientRev)

	actualLinkErr := actualLinks.Sender.Send(context.Background(), &amqp.Message{}, nil)
	require.Error(t, actualLinkErr)

	// now we'll recover, which should recreate everything
	err = links.RecoverIfNeeded(context.Background(), oldLWR.ID, actualLinkErr)
	require.NoError(t, err)
	require.EqualValues(t, 2, createLinksCalled)

	newLWR, err := links.Get(context.Background())
	require.NoError(t, err)

	// should work now, connection should be reopened
	assertLinks(t, newLWR)

	requireNewLinkSameConn(t, oldLWR, newLWR)
}

func TestAMQPLinks_LinkWithConnectionFailure(t *testing.T) {
	ns := newNamespaceForTest(t)

	defer func() { _ = ns.Close(false) }()

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: test.BuiltInTestQueue,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			return newLinksForAMQPLinksTest(test.BuiltInTestQueue, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer func() {
		err := links.Close(context.Background(), true)
		require.NoError(t, err)
	}()

	oldLWR, err := links.Get(context.Background())
	require.NoError(t, err)

	// shut down the connection
	conn, clientRev, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 1, clientRev)
	err = conn.Close()
	require.NoError(t, err)

	// verify the errors we expect
	actualLinkErr := oldLWR.Sender.Send(context.Background(), &amqp.Message{}, nil)
	require.Error(t, actualLinkErr)
	require.Equal(t, RecoveryKindConn, GetRecoveryKind(actualLinkErr))

	// now we'll recover, which should recreate everything
	require.NoError(t, links.RecoverIfNeeded(context.Background(), oldLWR.ID, actualLinkErr))

	newLWR, err := links.Get(context.Background())
	require.NoError(t, err)

	requireNewLinkNewConn(t, oldLWR, newLWR)

	err = newLWR.Sender.Send(context.Background(), &amqp.Message{Value: "hello world"}, nil)
	require.NoError(t, err)
}

func TestAMQPLinks_LinkWithConnectionFailureAndExpiredContext(t *testing.T) {
	ns := newNamespaceForTest(t)

	defer test.RequireNSClose(t, ns)

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: test.BuiltInTestQueue,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			return newLinksForAMQPLinksTest(test.BuiltInTestQueue, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer test.RequireLinksClose(t, links)

	oldLWR, err := links.Get(context.Background())
	require.NoError(t, err)

	// shut down the connection
	conn, clientRev, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, 1, clientRev)
	err = conn.Close()
	require.NoError(t, err)

	// verify the errors we expect
	connErr := oldLWR.Sender.Send(context.Background(), &amqp.Message{}, nil)
	require.Error(t, connErr)
	require.Equal(t, RecoveryKindConn, GetRecoveryKind(connErr))

	// when we try to recover and the user cancels we still need to leave it in a state
	// where the next call in will recover it.
	cancelledCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Hour))
	defer cancel()
	err = links.RecoverIfNeeded(cancelledCtx, oldLWR.ID, connErr)
	var netErr net.Error
	require.ErrorAs(t, err, &netErr)

	newLWR, err := links.Get(context.Background())
	require.NoError(t, err)

	requireNewLinkNewConn(t, oldLWR, newLWR)

	err = newLWR.Sender.Send(context.Background(), &amqp.Message{Value: "hello world"}, nil)
	require.NoError(t, err)
}

func TestAMQPLinks_LinkFailure(t *testing.T) {
	ns := newNamespaceForTest(t)

	defer test.RequireNSClose(t, ns)

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: test.BuiltInTestQueue,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			return newLinksForAMQPLinksTest(test.BuiltInTestQueue, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer test.RequireLinksClose(t, links)

	oldLWR, err := links.Get(context.Background())
	require.NoError(t, err)

	// shut down the links as if they were closed out from underneath us
	err = oldLWR.Receiver.(amqpwrap.AMQPReceiverCloser).Close(context.Background())
	require.NoError(t, err)
	err = oldLWR.Sender.(amqpwrap.AMQPSenderCloser).Close(context.Background())
	require.NoError(t, err)

	// verify the errors we expect
	linkErr := oldLWR.Sender.Send(context.Background(), &amqp.Message{}, nil)
	require.Error(t, linkErr)
	require.Equal(t, RecoveryKindLink, GetRecoveryKind(linkErr))

	err = links.RecoverIfNeeded(context.Background(), oldLWR.ID, linkErr)
	require.NoError(t, err)

	newLWR, err := links.Get(context.Background())
	require.NoError(t, err)

	requireNewLinkSameConn(t, oldLWR, newLWR)

	err = newLWR.Sender.Send(context.Background(), &amqp.Message{Value: "hello world"}, nil)
	require.NoError(t, err)
}

func TestAMQPLinks_LinkFailureUpgradedToConnectionError(t *testing.T) {
	ns := newNamespaceForTest(t)

	defer test.RequireNSClose(t, ns)

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: test.BuiltInTestQueue,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			return newLinksForAMQPLinksTest(test.BuiltInTestQueue, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer test.RequireLinksClose(t, links)

	oldLWR, err := links.Get(context.Background())
	require.NoError(t, err)

	// shut down the links as if they were closed out from underneath us
	test.RequireClose(t, oldLWR.Receiver.(amqpwrap.AMQPReceiverCloser))
	test.RequireClose(t, oldLWR.Sender.(amqpwrap.AMQPSenderCloser))

	// verify the errors we expect
	linkErr := oldLWR.Sender.Send(context.Background(), &amqp.Message{}, nil)
	require.Error(t, linkErr)
	require.Equal(t, RecoveryKindLink, GetRecoveryKind(linkErr))

	test.EnableStdoutLogging(t)

	expiredCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Hour))
	defer cancel()

	t.Logf("Sender.Send() returned error %v", linkErr)

	// this is a bit tricky - what we're testing here is that if we can still properly recover, even after the user has cancelled
	// the passed in context. We, for a short time, had a version of go-amqp that would destabilize the entire AMQP
	// connection if a NewSession(ctx) call was cancelled.
	//
	// What will ultimately end up happening is that we'll close all of our old state (cancellation doesn't stop that)
	// and no new stuff will end up getting created since the context is expired.
	err = links.RecoverIfNeeded(expiredCtx, oldLWR.ID, linkErr)
	var netErr net.Error
	require.ErrorAs(t, err, &netErr)

	// all the previous state was closed out and removed.
	require.Nil(t, links.(*AMQPLinksImpl).Sender)
	require.Nil(t, links.(*AMQPLinksImpl).Receiver)
	require.Nil(t, links.(*AMQPLinksImpl).session)
	require.Nil(t, links.(*AMQPLinksImpl).RPCLink)

	// TODO: temporarily commented out just to see if we've broken any other tests.
	// newLWR, err := links.Get(context.Background())
	// require.NoError(t, err)

	// requireNewLinkNewConn(t, oldLWR, newLWR)

	// err = newLWR.Sender.Send(context.Background(), &amqp.Message{Value: "hello world"}, nil)
	// require.NoError(t, err)
}

// TestAMQPLinksCBSLinkStillOpen makes sure we can recover from an incompletely
// closed $cbs link, which can happen if a user cancels and we can't properly close
// the link as a result.
func TestAMQPLinksCBSLinkStillOpen(t *testing.T) {
	// we're not going to use this client for these tests.
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	ns := newNamespaceForTest(t)

	defer func() { _ = ns.Close(false) }()

	session, oldConnID, err := ns.NewAMQPSession(context.Background())
	require.NoError(t, err)

	// opening a Sender to the $cbs endpoint. This endpoint can only be opened by a single
	// sender/receiver pair in a connection.
	_, err = session.NewSender(context.Background(), "$cbs", nil)
	require.NoError(t, err)

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: entityPath,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			return newLinksForAMQPLinksTest(entityPath, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	var lwid *LinksWithID

	err = links.Retry(context.Background(), exported.EventConn, "test", func(ctx context.Context, innerLwid *LinksWithID, args *utils.RetryFnArgs) error {
		lwid = innerLwid
		return nil
	}, exported.RetryOptions{
		RetryDelay:    -1,
		MaxRetryDelay: time.Millisecond,
	}, nil)

	defer func() {
		err := links.Close(context.Background(), true)
		require.NoError(t, err)
	}()

	require.NoError(t, err)
	require.Equal(t, oldConnID+1, lwid.ID.Conn, "Connection gets incremented since it had to be reset")
}

func TestAMQPLinksLiveRecoverLink(t *testing.T) {
	// we're not going to use this client for these tests.
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	ns := newNamespaceForTest(t)

	defer func() { _ = ns.Close(false) }()

	createLinksCalled := 0

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: entityPath,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			createLinksCalled++
			return newLinksForAMQPLinksTest(entityPath, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer func() {
		err := links.Close(context.Background(), true)
		require.NoError(t, err)
	}()

	require.EqualValues(t, 0, createLinksCalled)
	require.NoError(t, links.RecoverIfNeeded(context.Background(), LinkID{}, &amqp.ConnError{}))
	require.EqualValues(t, 1, createLinksCalled)

	lwr, err := links.Get(context.Background())
	require.NoError(t, err)

	require.NoError(t, links.RecoverIfNeeded(context.Background(), lwr.ID, &amqp.LinkError{}))
	require.EqualValues(t, 2, createLinksCalled)
}

func TestAMQPLinksLiveRace(t *testing.T) {
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	ns := newNamespaceForTest(t)

	defer func() { _ = ns.Close(false) }()

	createLinksCalled := 0

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: entityPath,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			createLinksCalled++
			return newLinksForAMQPLinksTest(entityPath, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer func() {
		err := links.Close(context.Background(), true)
		require.NoError(t, err)
	}()

	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := links.RecoverIfNeeded(context.Background(), LinkID{}, &amqp.ConnError{})
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	// TODO: also check that the connection hasn't recycled multiple times.
	require.EqualValues(t, 1, createLinksCalled)
}

func TestAMQPLinksLiveRaceLink(t *testing.T) {
	endCapture := test.CaptureLogsForTest(false)
	defer func() {
		messages := endCapture()
		for _, msg := range messages {
			fmt.Printf("%s\n", msg)
		}
	}()

	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	ns := newNamespaceForTest(t)

	defer func() { _ = ns.Close(false) }()

	createLinksCalled := 0

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: entityPath,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			createLinksCalled++
			return newLinksForAMQPLinksTest(entityPath, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer func() {
		err := links.Close(context.Background(), true)
		require.NoError(t, err)
	}()

	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := links.RecoverIfNeeded(context.Background(), LinkID{}, &amqp.LinkError{})
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

	ns := newNamespaceForTest(t)

	defer func() { _ = ns.Close(false) }()

	createLinksCalled := 0

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: entityPath,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			createLinksCalled++
			return newLinksForAMQPLinksTest(entityPath, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer func() {
		err := links.Close(context.Background(), true)
		require.NoError(t, err)
	}()

	err := links.Retry(context.Background(), log.Event("NotUsed"), "NotUsed", func(ctx context.Context, lwid *LinksWithID, args *utils.RetryFnArgs) error {
		// force recoveries
		return &amqp.ConnError{}
	}, exported.RetryOptions{
		MaxRetries: 2,
		// note: omitting MaxRetries just to give a sanity check that
		// we do setDefaults() before we run.
		RetryDelay:    time.Millisecond,
		MaxRetryDelay: time.Millisecond,
	}, nil)

	var connErr *amqp.ConnError
	require.ErrorAs(t, err, &connErr)
	require.EqualValues(t, 3, createLinksCalled)
}

func TestAMQPLinksMultipleWithSameConnection(t *testing.T) {
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	ns := newNamespaceForTest(t)

	defer func() { _ = ns.Close(false) }()

	createLinksCalled := 0

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: entityPath,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			createLinksCalled++
			return newLinksForAMQPLinksTest(entityPath, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer func() {
		_ = links.Close(context.Background(), true)
	}()

	createLinksCalled2 := 0

	links2 := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: entityPath,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			createLinksCalled2++
			return newLinksForAMQPLinksTest(entityPath, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer func() {
		_ = links2.Close(context.Background(), true)
	}()

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
		err = links.RecoverIfNeeded(context.Background(), lwr.ID, &amqp.LinkError{})
		require.NoError(t, err)
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		err := links2.RecoverIfNeeded(context.Background(), lwr2.ID, &amqp.LinkError{})
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

func TestAMQPLinksCloseIfNeeded(t *testing.T) {
	t.Run("fatal", func(t *testing.T) {
		for _, fatalErr := range []error{NewErrNonRetriable("")} {
			receiver := &FakeAMQPReceiver{}
			sender := &FakeAMQPSender{}
			ns := &FakeNS{}

			links := NewAMQPLinks(NewAMQPLinksArgs{
				NS:         ns,
				EntityPath: "entityPath",
				CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
					return sender, receiver, nil
				},
				GetRecoveryKindFunc: GetRecoveryKind,
			})

			defer func() {
				err := links.Close(context.Background(), true)
				require.NoError(t, err)
			}()

			_, err := links.Get(context.Background())
			require.NoError(t, err)

			rk := links.CloseIfNeeded(context.Background(), fatalErr)
			require.Equal(t, RecoveryKindFatal, rk)
			require.Equal(t, 1, receiver.Closed)
			require.Equal(t, 1, sender.Closed)
			require.Equal(t, 1, ns.CloseCalled)
		}
	})

	t.Run("link", func(t *testing.T) {
		receiver := &FakeAMQPReceiver{}
		sender := &FakeAMQPSender{}
		ns := &FakeNS{}

		links := NewAMQPLinks(NewAMQPLinksArgs{
			NS:         ns,
			EntityPath: "entityPath",
			CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
				return sender, receiver, nil
			},
			GetRecoveryKindFunc: GetRecoveryKind,
		})

		defer func() {
			err := links.Close(context.Background(), true)
			require.NoError(t, err)
		}()

		_, err := links.Get(context.Background())
		require.NoError(t, err)

		rk := links.CloseIfNeeded(context.Background(), &amqp.LinkError{})
		require.Equal(t, RecoveryKindLink, rk)
		require.Equal(t, 1, receiver.Closed)
		require.Equal(t, 1, sender.Closed)
		require.Equal(t, 0, ns.CloseCalled)
	})

	t.Run("conn", func(t *testing.T) {
		receiver := &FakeAMQPReceiver{}
		sender := &FakeAMQPSender{}
		ns := &FakeNS{}

		links := NewAMQPLinks(NewAMQPLinksArgs{
			NS:         ns,
			EntityPath: "entityPath",
			CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
				return sender, receiver, nil
			},
			GetRecoveryKindFunc: GetRecoveryKind,
		})

		defer func() {
			err := links.Close(context.Background(), true)
			require.NoError(t, err)
		}()

		_, err := links.Get(context.Background())
		require.NoError(t, err)

		rk := links.CloseIfNeeded(context.Background(), &amqp.ConnError{})
		require.Equal(t, RecoveryKindConn, rk)
		require.Equal(t, 1, receiver.Closed)
		require.Equal(t, 1, sender.Closed)
		require.Equal(t, 1, ns.CloseCalled)
	})

	t.Run("none", func(t *testing.T) {
		receiver := &FakeAMQPReceiver{}
		sender := &FakeAMQPSender{}
		ns := &FakeNS{}

		links := NewAMQPLinks(NewAMQPLinksArgs{
			NS:         ns,
			EntityPath: "entityPath",
			CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
				return sender, receiver, nil
			},
			GetRecoveryKindFunc: GetRecoveryKind,
		})

		defer func() {
			err := links.Close(context.Background(), true)
			require.NoError(t, err)
		}()

		_, err := links.Get(context.Background())
		require.NoError(t, err)

		rk := links.CloseIfNeeded(context.Background(), nil)
		require.Equal(t, RecoveryKindNone, rk)
		require.Equal(t, 0, receiver.Closed)
		require.Equal(t, 0, sender.Closed)
		require.Equal(t, 0, ns.CloseCalled)

		rk = links.CloseIfNeeded(context.Background(), context.Canceled)
		require.Equal(t, RecoveryKindNone, rk)
		require.Equal(t, 0, receiver.Closed)
		require.Equal(t, 0, sender.Closed)
		require.Equal(t, 0, ns.CloseCalled)

		rk = links.CloseIfNeeded(context.Background(), context.DeadlineExceeded)
		require.Equal(t, RecoveryKindNone, rk)
		require.Equal(t, 0, receiver.Closed)
		require.Equal(t, 0, sender.Closed)
		require.Equal(t, 0, ns.CloseCalled)
	})
}

func TestAMQPLinksCreditTracking(t *testing.T) {
	entityPath, cleanup := test.CreateExpiringQueue(t, nil)
	defer cleanup()

	ns := newNamespaceForTest(t)

	links := NewAMQPLinks(NewAMQPLinksArgs{
		NS:         ns,
		EntityPath: entityPath,
		CreateLinkFunc: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			return newLinksForAMQPLinksTest(entityPath, session)
		},
		GetRecoveryKindFunc: GetRecoveryKind,
	})

	defer func() {
		err := links.Close(context.Background(), true)
		require.NoError(t, err)
	}()

	lwr, err := links.Get(context.Background())
	require.NoError(t, err)

	t.Run("credits are decremented when messages are amqpReceiver.Receive()'d", func(t *testing.T) {
		err = lwr.Sender.Send(context.Background(), &amqp.Message{
			Data: [][]byte{[]byte("Received")},
		}, nil)
		require.NoError(t, err)

		err = lwr.Receiver.IssueCredit(1)
		require.NoError(t, err)
		require.Equal(t, uint32(1), lwr.Receiver.Credits())

		message, err := lwr.Receiver.Receive(context.Background(), nil)
		require.NoError(t, err)
		require.Equal(t, [][]byte{[]byte("Received")}, message.Data)
		require.Equal(t, uint32(0), lwr.Receiver.Credits())

		err = lwr.Receiver.AcceptMessage(context.Background(), message)
		require.NoError(t, err)
	})

	t.Run("credits are decremented when messages are amqpReceiver.Prefetched()", func(t *testing.T) {
		err = lwr.Sender.Send(context.Background(), &amqp.Message{
			Data: [][]byte{[]byte("Received")},
		}, nil)
		require.NoError(t, err)

		err = lwr.Receiver.IssueCredit(1)
		require.NoError(t, err)
		require.Equal(t, uint32(1), lwr.Receiver.Credits())

		// prefetched messages arrive, but we don't block in Prefetched() so
		// we'll have to poll our receiver for this part.
		deadline := time.Now().Add(time.Minute)

		for time.Until(deadline) > 0 {
			message := lwr.Receiver.Prefetched()

			if message != nil {
				require.Equal(t, [][]byte{[]byte("Received")}, message.Data)
				require.Equal(t, uint32(0), lwr.Receiver.Credits())

				err = lwr.Receiver.AcceptMessage(context.Background(), message)
				require.NoError(t, err)
				break
			}

			time.Sleep(time.Second)
		}
	})

	t.Run("credits are not altered if an error comes back from Prefetched() or Receive()", func(t *testing.T) {
		// now that the link is empty, let's test:

		// A receive where an error happens (cancellation, in this case)
		// this won't touch the credit since nothing is actually received.
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err = lwr.Receiver.Receive(ctx, nil)
		require.ErrorIs(t, err, context.Canceled)
		require.Equal(t, uint32(0), lwr.Receiver.Credits())

		// a prefetch where there isn't anything.
		message := lwr.Receiver.Prefetched()
		require.Nil(t, message)
		require.Equal(t, uint32(0), lwr.Receiver.Credits())
	})
}

func requireNewLinkSameConn(t *testing.T, oldLWID *LinksWithID, newLWID *LinksWithID) {
	t.Helper()
	require.NotEqual(t, oldLWID.Sender.LinkName(), newLWID.Sender.LinkName(), "Link should have a new ID because it was recreated")
	require.Equal(t, oldLWID.ID.Conn, newLWID.ID.Conn, "Connection ID should be the same since recreation wasn't needed")
}

func requireNewLinkNewConn(t *testing.T, oldLWID *LinksWithID, newLWID *LinksWithID) {
	t.Helper()
	require.NotEqual(t, oldLWID.Sender.LinkName(), newLWID.Sender.LinkName(), "Link should have a new ID because it was recreated")
	require.Equal(t, oldLWID.ID.Conn+1, newLWID.ID.Conn, "Connection ID should be recreated")
}
