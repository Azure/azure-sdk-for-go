// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestAMQPLinks(t *testing.T) {
	fakeSender := &FakeAMQPSender{}
	fakeSession := &FakeAMQPSession{}
	fakeMgmtClient := &fakeMgmtClient{}

	createLinkFunc, createLinkCallCount := setupCreateLinkResponses(t, []createLinkResponse{
		{sender: fakeSender},
	})

	links := newAMQPLinks(&FakeNS{
		Session:    fakeSession,
		MgmtClient: fakeMgmtClient,
	}, "entityPath", &fakeRetrier{}, createLinkFunc)

	require.EqualValues(t, "entityPath", links.EntityPath())
	require.EqualValues(t, "audience: entityPath", links.Audience())

	// successful Get() where a Sender was initialized
	sender, receiver, mgmt, linkRevision, err := links.Get(context.Background())
	require.NotNil(t, sender)
	require.NotNil(t, mgmt) // you always get a free mgmt link
	require.Nil(t, receiver)
	require.Nil(t, err)
	require.EqualValues(t, 1, linkRevision)
	require.EqualValues(t, 1, *createLinkCallCount)

	// further calls should just be cached instances
	sender2, receiver2, mgmt2, linkRevision2, err2 := links.Get(context.Background())
	require.EqualValues(t, sender, sender2)
	require.EqualValues(t, mgmt, mgmt2)
	require.Nil(t, receiver2)
	require.Nil(t, err2)
	require.EqualValues(t, 1, linkRevision2, "No recover calls, so link revision remains the same")
	require.EqualValues(t, 1, *createLinkCallCount, "No create call needed since an instance was cached")

	// closing multiple times is fine.
	asAMQPLinks, ok := links.(*amqpLinks)
	require.True(t, ok)

	require.NoError(t, links.Close(context.Background(), false))
	require.False(t, asAMQPLinks.closedPermanently)

	require.NoError(t, links.Close(context.Background(), true))
	require.True(t, asAMQPLinks.closedPermanently)

	require.NoError(t, links.Close(context.Background(), true))
	require.True(t, asAMQPLinks.closedPermanently)

	require.NoError(t, links.Close(context.Background(), false))
	require.True(t, asAMQPLinks.closedPermanently)

	// and the individual links are closed as well
	require.EqualValues(t, 1, fakeSender.Closed)
	require.EqualValues(t, 1, fakeSession.closed)
	require.EqualValues(t, 1, fakeMgmtClient.closed)

	// and calls to Get() will indicate the amqpLinks has been closed permanently
	sender, receiver, mgmt, linkRevision, err = links.Get(context.Background())
	require.Nil(t, sender)
	require.Nil(t, receiver)
	require.Nil(t, mgmt)
	require.EqualValues(t, 0, linkRevision)

	_, ok = err.(NonRetriable)
	require.True(t, ok)
}

type permanentNetError struct {
	temp    bool
	timeout bool
}

func (pe permanentNetError) Timeout() bool   { return pe.timeout }
func (pe permanentNetError) Temporary() bool { return pe.temp }
func (pe permanentNetError) Error() string   { return "Fake but very permanent error" }

func TestAMQPLinksRecovery(t *testing.T) {
	sess := &FakeAMQPSession{}
	ns := &FakeNS{
		Session: sess,
	}
	sender := &FakeAMQPSender{}

	createLinkCalled := 0

	tmpLinks := newAMQPLinks(ns, "entity path", NewBackoffRetrier(BackoffRetrierParams{}), func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
		createLinkCalled++
		return sender, nil, nil
	})

	links, _ := tmpLinks.(*amqpLinks)

	links.clientRevision = 2001
	links.sender = sender

	ctx := context.TODO()

	require.Nil(t, links.RecoverIfNeeded(ctx, 0, nil))
	require.EqualValues(t, 0, sess.closed)
	require.EqualValues(t, 0, ns.recovered)
	require.EqualValues(t, 0, createLinkCalled, "new links aren't needed")
	require.False(t, links.closedPermanently, "link should still be usable")
	require.Empty(t, ns.clientRevisions, "no connection recoveries happened")

	require.Nil(t, links.RecoverIfNeeded(ctx, 0, errors.New("Passes through")))
	require.EqualValues(t, 0, sess.closed)
	require.EqualValues(t, 0, ns.recovered)
	require.EqualValues(t, 0, createLinkCalled, "new links aren't needed")
	require.False(t, links.closedPermanently, "link should still be usable")
	require.Empty(t, ns.clientRevisions, "no connection recoveries happened")

	// now let's initiate a recovery at the connection level
	require.NoError(t, links.RecoverIfNeeded(ctx, 0, permanentNetError{}), permanentNetError{}.Error())
	require.EqualValues(t, 1, ns.recovered, "client gets recovered")
	require.EqualValues(t, 1, sender.Closed, "link is closed")
	require.EqualValues(t, 1, createLinkCalled, "link is created")
	require.False(t, links.closedPermanently, "link should still be usable")
	require.EqualValues(t, []uint64{2001}, ns.clientRevisions, "links handed us the client revision it got last")

	// validate that our linkRevision got updated and that we're returning it.
	// (note that link revisions start at 1, so we're not at 2, even though
	// only one recover has happened)
	_, _, _, linkRevision, err := links.Get(ctx)
	require.NoError(t, err)
	require.EqualValues(t, uint64(2), linkRevision)

	ns.recovered = 0
	sender.Closed = 0
	createLinkCalled = 0

	// let's do just a link level one
	require.NoError(t, links.RecoverIfNeeded(ctx, links.revision+1, amqp.ErrLinkDetached), amqp.ErrLinkDetached.Error())
	require.EqualValues(t, 0, ns.recovered)
	require.EqualValues(t, 1, sender.Closed)
	require.EqualValues(t, 1, createLinkCalled)

	_, _, _, linkRevision, err = links.Get(ctx)
	require.NoError(t, err)
	require.EqualValues(t, uint64(3), linkRevision)

	// cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ns.recovered = 0
	sender.Closed = 0
	createLinkCalled = 0

	// cancellation overrides any other logic.
	require.Error(t, links.RecoverIfNeeded(ctx, links.revision+1, amqp.ErrLinkDetached), amqp.ErrLinkDetached.Error())
	require.EqualValues(t, 0, ns.recovered)
	require.EqualValues(t, 0, sender.Closed)
	require.EqualValues(t, 0, createLinkCalled)
}

func TestAMQPLinks_Closed(t *testing.T) {
	createLinks := func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
		return nil, nil, nil
	}

	links := newAMQPLinks(&FakeNS{}, "hello", &backoffRetrier{}, createLinks)
	links.Close(context.Background(), true)

	_, _, _, _, err := links.Get(context.Background())

	require.True(t, IsNonRetriable(err))
}

func setupCreateLinkResponses(t *testing.T, responses []createLinkResponse) (CreateLinkFunc, *int) {
	callCount := 0
	testCreateLinkFunc := func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error) {
		callCount++

		if len(responses) == 0 {
			require.Fail(t, "createLinkFunc called too many times")
		}

		r := responses[0]
		responses = responses[1:]

		return r.sender, r.receiver, r.err
	}

	return testCreateLinkFunc, &callCount
}
