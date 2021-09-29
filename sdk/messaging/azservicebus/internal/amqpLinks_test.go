// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAMQPLinks(t *testing.T) {
	fakeSender := &fakeAMQPSender{}
	fakeSession := &fakeAMQPSession{}
	fakeMgmtClient := &fakeMgmtClient{}

	createLinkFunc, createLinkCallCount := setupCreateLinkResponses(t, []createLinkResponse{
		{sender: fakeSender},
	})

	links := newAMQPLinks(&fakeNS{
		Session:    fakeSession,
		MgmtClient: fakeMgmtClient,
	}, "entityPath", createLinkFunc)

	require.EqualValues(t, "entityPath", links.EntityPath())
	require.EqualValues(t, "audience: entityPath", links.Audience())

	// successful Get() where a Sender was initialized
	sender, receiver, mgmt, linkRevision, err := links.Get(context.Background())
	require.NotNil(t, sender)
	require.NotNil(t, mgmt) // you always get a free mgmt link
	require.Nil(t, receiver)
	require.Nil(t, err)
	require.EqualValues(t, 0, linkRevision)
	require.EqualValues(t, 1, *createLinkCallCount)

	// further calls should just be cached instances
	sender2, receiver2, mgmt2, linkRevision2, err2 := links.Get(context.Background())
	require.EqualValues(t, sender, sender2)
	require.EqualValues(t, mgmt, mgmt2)
	require.Nil(t, receiver2)
	require.Nil(t, err2)
	require.EqualValues(t, 0, linkRevision2, "No recover calls, so link revision remains the same")
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
	require.EqualValues(t, 1, fakeSender.closed)
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
