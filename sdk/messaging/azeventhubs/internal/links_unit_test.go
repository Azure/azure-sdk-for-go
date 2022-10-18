// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestLinks_NoOp(t *testing.T) {
	fakeNS := &FakeNSForPartClient{}
	links := NewLinks(fakeNS, "managementPath", func(partitionID string) string {
		return fmt.Sprintf("part:%s", partitionID)
	},
		func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (*FakeAMQPReceiver, error) {
			panic("Nothing should be created for a nil error")
		})

	// no error just no-ops
	err := links.RecoverIfNeeded(context.Background(), "0", nil, nil)
	require.NoError(t, err)
}

func TestLinks_LinkStale(t *testing.T) {
	fakeNS := &FakeNSForPartClient{}

	var nextID int
	var receivers []*FakeAMQPReceiver

	links := NewLinks(fakeNS, "managementPath", func(partitionID string) string {
		return fmt.Sprintf("part:%s", partitionID)
	},
		func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (*FakeAMQPReceiver, error) {
			nextID++
			receivers = append(receivers, &FakeAMQPReceiver{
				NameForLink: fmt.Sprintf("Link%d", nextID),
			})
			return receivers[len(receivers)-1], nil
		})

	staleLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)
	require.NotNil(t, staleLWID)
	require.NotNil(t, links.links["0"], "cache contains the newly created link for partition 0")

	// we'll recover first, but our lwid (after this recovery) is stale since
	// the link cache will be updated after this is done.
	err = links.RecoverIfNeeded(context.Background(), "0", staleLWID, &amqp.DetachError{})
	require.NoError(t, err)
	require.Nil(t, links.links["0"], "closed link is removed from the cache")
	require.Equal(t, 1, receivers[0].CloseCalled, "original receiver is closed, and replaced")

	// trying to recover again is a no-op (if nothing is in the cache)
	err = links.RecoverIfNeeded(context.Background(), "0", staleLWID, &amqp.DetachError{})
	require.NoError(t, err)
	require.Nil(t, links.links["0"], "closed link is removed from the cache")
	require.Equal(t, 1, receivers[0].CloseCalled, "original receiver is closed, and replaced")

	receivers = nil

	// now let's create a new link, and attempt using the old stale lwid
	// it'll no-op then too - we don't need to do anything, they should just call GetLink() again.
	newLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)
	require.NotNil(t, newLWID)
	require.Equal(t, (*links.links["0"].Link).LinkName(), newLWID.Link.LinkName(), "cache contains the newly created link for partition 0")

	err = links.RecoverIfNeeded(context.Background(), "0", staleLWID, &amqp.DetachError{})
	require.NoError(t, err)
	require.Equal(t, 0, receivers[0].CloseCalled, "receiver is NOT closed - we didn't need to replace it since the lwid with the error was stale")
}

func TestLinks_LinkRecoveryOnly(t *testing.T) {
	fakeNS := &FakeNSForPartClient{}

	var nextID int
	var receivers []*FakeAMQPReceiver

	links := NewLinks(fakeNS, "managementPath", func(partitionID string) string {
		return fmt.Sprintf("part:%s", partitionID)
	},
		func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (*FakeAMQPReceiver, error) {
			nextID++
			receivers = append(receivers, &FakeAMQPReceiver{
				NameForLink: fmt.Sprintf("Link%d", nextID),
			})
			return receivers[len(receivers)-1], nil
		})

	lwid, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)
	require.NotNil(t, lwid)
	require.NotNil(t, links.links["0"], "cache contains the newly created link for partition 0")

	err = links.RecoverIfNeeded(context.Background(), "0", lwid, &amqp.DetachError{})
	require.NoError(t, err)
	require.Nil(t, links.links["0"], "cache will no longer a link for partition 0")

	// no new links are create - we'll need to do something that requires a link
	// to cause it to come back.
	require.Equal(t, 1, len(receivers))
	require.Equal(t, 1, receivers[0].CloseCalled)

	receivers = nil

	// cause a new link to get created to replace the old one.
	newLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)
	require.NotEqual(t, lwid, newLWID, "new link gets a new ID")
	require.NotNil(t, links.links["0"], "cache contains the newly created link for partition 0")

	require.Equal(t, 1, len(receivers))
	require.Equal(t, 0, receivers[0].CloseCalled)
}

func TestLinks_ConnectionRecovery(t *testing.T) {
	recoverClientCalled := 0

	fakeNS := &FakeNSForPartClient{
		RecoverFn: func(ctx context.Context, clientRevision uint64) error {
			// we'll just always recover for our test.
			recoverClientCalled++
			return nil
		},
	}

	var nextID int
	var receivers []*FakeAMQPReceiver

	links := NewLinks(fakeNS, "managementPath", func(partitionID string) string {
		return fmt.Sprintf("part:%s", partitionID)
	},
		func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (*FakeAMQPReceiver, error) {
			nextID++
			receivers = append(receivers, &FakeAMQPReceiver{
				NameForLink: fmt.Sprintf("Link%d", nextID),
			})
			return receivers[len(receivers)-1], nil
		})

	lwid, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)
	require.NotNil(t, lwid)
	require.NotNil(t, links.links["0"], "cache contains the newly created link for partition 0")

	require.Equal(t, recoverClientCalled, 0)

	err = links.RecoverIfNeeded(context.Background(), "0", lwid, &amqp.ConnectionError{})
	require.NoError(t, err)
	require.Nil(t, links.links["0"], "cache will no longer a link for partition 0")

	require.Equal(t, recoverClientCalled, 1, "client was recovered")

	// no new links are create - we'll need to do something that requires a link
	// to cause it to come back.
	require.Equal(t, 1, len(receivers))
	require.Equal(t, 1, receivers[0].CloseCalled)

	// cause a new link to get created to replace the old one.
	receivers = nil

	newLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)
	require.NotEqual(t, lwid, newLWID, "new link gets a new ID")
	require.NotNil(t, links.links["0"], "cache contains the newly created link for partition 0")

	require.Equal(t, 1, len(receivers))
	require.Equal(t, 0, receivers[0].CloseCalled)
}
