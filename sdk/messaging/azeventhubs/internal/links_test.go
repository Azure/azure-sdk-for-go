// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestLinksCBSLinkStillOpen(t *testing.T) {
	// we're not going to use this client for these tests.
	testParams := test.GetConnectionParamsForTest(t)
	ns, err := NewNamespace(NamespaceWithTokenCredential(testParams.EventHubNamespace, testParams.Cred))
	require.NoError(t, err)

	defer func() { _ = ns.Close(context.Background(), true) }()

	session, oldConnID, err := ns.NewAMQPSession(context.Background())
	require.NoError(t, err)

	// opening a Sender to the $cbs endpoint. This endpoint can only be opened by a single
	// sender/receiver pair in a connection.
	_, err = session.NewSender(context.Background(), "$cbs", "", nil)
	require.NoError(t, err)

	newLinkFn := func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string, partitionID string) (AMQPSenderCloser, error) {
		return session.NewSender(ctx, entityPath, "", &amqp.SenderOptions{
			SettlementMode:              to.Ptr(amqp.SenderSettleModeMixed),
			RequestedReceiverSettleMode: to.Ptr(amqp.ReceiverSettleModeFirst),
		})
	}

	formatEntityPath := func(partitionID string) string {
		return fmt.Sprintf("%s/Partitions/%s", testParams.EventHubName, partitionID)
	}

	links := NewLinks(ns, fmt.Sprintf("%s/$management", testParams.EventHubName), formatEntityPath, newLinkFn)

	var lwid LinkWithID[AMQPSenderCloser]

	err = links.Retry(context.Background(), exported.EventConn, "test", "0", exported.RetryOptions{
		RetryDelay:    -1,
		MaxRetryDelay: time.Millisecond,
	}, func(ctx context.Context, innerLWID LinkWithID[AMQPSenderCloser]) error {
		lwid = innerLWID
		return nil
	})
	require.NoError(t, err)

	defer func() {
		err := links.Close(context.Background())
		require.NoError(t, err)
	}()

	require.NoError(t, err)
	require.Equal(t, oldConnID+1, lwid.ConnID(), "Connection gets incremented since it had to be reset")
}

func TestLinksRecoverLinkWithConnectionFailure(t *testing.T) {
	ns, links := newLinksForTest(t)
	defer test.RequireClose(t, links)
	defer test.RequireNSClose(t, ns)

	oldLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)

	// cause a connection level failure by closing the connection out from underneath
	// this.
	origConn, _, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	err = origConn.Close()
	require.NoError(t, err)

	err = oldLWID.Link().Send(context.Background(), &amqp.Message{}, nil)
	require.Error(t, err)
	require.Equal(t, RecoveryKindConn, GetRecoveryKind(err))

	// now recover like normal

	err = links.lr.RecoverIfNeeded(context.Background(), lwidToError(err, oldLWID))
	require.NoError(t, err)

	newLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)

	requireNewLinkNewConn(t, oldLWID, newLWID)

	err = newLWID.Link().Send(context.Background(), &amqp.Message{
		Data: [][]byte{[]byte("TestLinksRecoverLinkWithConnectionFailure")},
	}, nil)
	require.NoError(t, err)
}

// TestLinksRecoverLinkWithConnectionFailureAndExpiredContext checks that we're able to recover
// after a "partial" recovery, where the user or the passed in context was already cancelled. The
// recovery, in those cases, should leave us in a state that the next call to GetLinks()
// will reinstantiate everything.
func TestLinksRecoverLinkWithConnectionFailureAndExpiredContext(t *testing.T) {
	ns, links := newLinksForTest(t)
	defer test.RequireClose(t, links)
	defer test.RequireNSClose(t, ns)

	oldLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)

	// cause a connection level failure by closing the connection out from underneath
	// this.
	origConn, _, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	err = origConn.Close()
	require.NoError(t, err)

	err = oldLWID.Link().Send(context.Background(), &amqp.Message{}, nil)
	require.Error(t, err)
	require.Equal(t, RecoveryKindConn, GetRecoveryKind(err))

	// Try to recover, but using an expired context. We'll get a network error (not enough time to resolve or
	// create a connection), which would normally be a connection level recovery event.
	cancelledCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Hour))
	defer cancel()

	err = links.lr.RecoverIfNeeded(cancelledCtx, lwidToError(err, oldLWID))
	var netErr net.Error
	require.ErrorAs(t, err, &netErr)

	// now recover like normal
	err = links.lr.RecoverIfNeeded(context.Background(), lwidToError(err, oldLWID))
	require.NoError(t, err)

	newLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)

	requireNewLinkNewConn(t, oldLWID, newLWID)

	err = newLWID.Link().Send(context.Background(), &amqp.Message{
		Data: [][]byte{[]byte("hello world")},
	}, nil)
	require.NoError(t, err)
}

func TestLinkFailureWhenConnectionIsDead(t *testing.T) {
	ns, links := newLinksForTest(t)
	defer test.RequireClose(t, links)
	defer test.RequireNSClose(t, ns)

	oldLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)

	// cause a connection level failure by closing the connection out from underneath
	// this.
	origConn, _, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	err = origConn.Close()
	require.NoError(t, err)

	err = oldLWID.Link().Send(context.Background(), &amqp.Message{}, nil)
	require.Error(t, err)
	require.Equal(t, RecoveryKindConn, GetRecoveryKind(err))

	err = links.lr.RecoverIfNeeded(context.Background(), lwidToError(&amqp.LinkError{}, oldLWID))
	var connErr *amqp.ConnError
	require.ErrorAs(t, err, &connErr)
	require.Nil(t, connErr.RemoteErr, "is the forwarded error from the closed connection")
	require.Equal(t, RecoveryKindConn, GetRecoveryKind(connErr), "next recovery would force a connection level recovery")

	err = links.lr.RecoverIfNeeded(context.Background(), lwidToError(connErr, oldLWID))
	require.NoError(t, err)

	newLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)

	requireNewLinkNewConn(t, oldLWID, newLWID)

	err = newLWID.Link().Send(context.Background(), &amqp.Message{
		Data: [][]byte{[]byte("TestLinkFailureWhenConnectionIsDead")},
	}, nil)
	require.NoError(t, err)
}

func TestLinkFailure(t *testing.T) {
	ns, links := newLinksForTest(t)
	defer test.RequireClose(t, links)
	defer test.RequireNSClose(t, ns)

	oldLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)

	// close the Receiver out from under the Links
	err = oldLWID.Link().Close(context.Background())
	require.NoError(t, err)

	err = oldLWID.Link().Send(context.Background(), &amqp.Message{Value: "hello"}, nil)
	require.Error(t, err)
	require.Equal(t, RecoveryKindLink, GetRecoveryKind(err))

	// we only close the link here, it actually opens up on the next time we call links.Get()
	cancelledCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Hour))
	defer cancel()

	err = links.lr.RecoverIfNeeded(cancelledCtx, lwidToError(err, oldLWID))
	require.NoError(t, err)

	newLWID, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)

	requireNewLinkSameConn(t, oldLWID, newLWID)
}

func TestLinksManagementRetry(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)
	ns, links := newLinksForTest(t)
	defer ns.Close(context.Background(), true)
	defer test.RequireClose(t, links)

	var prevLWID LinkWithID[amqpwrap.RPCLink]
	called := 0

	getEventHubProps := func(ctx context.Context, lwid LinkWithID[amqpwrap.RPCLink]) error {
		called++
		// mostly lifted from mgmt.go/getEventHubProperties
		token, err := ns.GetTokenForEntity(testParams.EventHubName)

		if err != nil {
			return err
		}

		amqpMsg := &amqp.Message{
			ApplicationProperties: map[string]any{
				"operation":      "READ",
				"name":           testParams.EventHubName,
				"type":           "com.microsoft:eventhub",
				"security_token": token.Token,
			},
		}

		resp, err := lwid.Link().RPC(context.Background(), amqpMsg)

		if err != nil {
			return err
		}

		if resp.Code >= 300 {
			return fmt.Errorf("failed getting partition properties: %v", resp.Description)
		}

		partitionIDs := resp.Message.Value.(map[string]any)["partition_ids"]
		require.NotEmpty(t, partitionIDs)

		prevLWID = lwid
		return nil
	}

	err := links.RetryManagement(context.Background(), "test", "op", exported.RetryOptions{}, getEventHubProps)
	require.NoError(t, err)
	require.Equal(t, 1, called, "nothing broken, should work on the first time")

	// we can do a quick check of another bit - that we don't just arbitrarily reset a management link
	// if the link _name_ doesn't match.
	err = links.closeManagementLinkIfMatch(context.Background(), "not the management link name")
	require.NoError(t, err)
	require.NotNil(t, links.managementLink)
	origMgmtLink := links.managementLink

	// let's trigger connection recovery by closing the amqp.Conn
	// behind `Links`'s back.
	client, connID, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	require.Equal(t, prevLWID.ConnID(), connID, "connection is stable")

	err = client.Close()
	require.NoError(t, err)

	called = 0

	err = links.RetryManagement(context.Background(), "test", "op", exported.RetryOptions{
		MaxRetries:    1,
		RetryDelay:    time.Nanosecond,
		MaxRetryDelay: time.Nanosecond,
	}, getEventHubProps)
	require.NoError(t, err)

	require.Equal(t, connID+1, prevLWID.ConnID(), "new connection was created")
	require.Equal(t, 2, called, "first usage failed due to dead connection, second call worked after recovery")
	require.NotEqual(t, origMgmtLink.Link().LinkName(), links.managementLink.Link().LinkName(), "management link also recreated")

	// and now let's try it with the mgmt link dead.
	origMgmtLWID, err := links.GetManagementLink(context.Background())
	require.NoError(t, err)

	err = origMgmtLWID.Link().(*rpcLink).receiver.Close(context.Background())
	require.NoError(t, err)

	err = links.RetryManagement(context.Background(), "test", "op", exported.RetryOptions{
		MaxRetries:    1,
		RetryDelay:    time.Nanosecond,
		MaxRetryDelay: time.Nanosecond,
	}, getEventHubProps)
	require.NoError(t, err)

	require.Equal(t, origMgmtLWID.ConnID(), prevLWID.ConnID(), "connection wasn't touched")
	require.NotEqual(t, origMgmtLWID.Link().LinkName(), prevLWID.Link().LinkName(), "management link recreated")

	test.RequireClose(t, links)

	require.Nil(t, links.managementLink)
}

func requireNewLinkSameConn(t *testing.T, oldLWID LinkWithID[AMQPSenderCloser], newLWID LinkWithID[AMQPSenderCloser]) {
	t.Helper()
	require.NotEqual(t, oldLWID.Link().LinkName(), newLWID.Link().LinkName(), "Link should have a new ID because it was recreated")
	require.Equal(t, oldLWID.ConnID(), newLWID.ConnID(), "Connection ID should be the same since recreation wasn't needed")
}

func requireNewLinkNewConn(t *testing.T, oldLWID LinkWithID[AMQPSenderCloser], newLWID LinkWithID[AMQPSenderCloser]) {
	t.Helper()
	require.NotEqual(t, oldLWID.Link().LinkName(), newLWID.Link().LinkName(), "Link should have a new ID because it was recreated")
	require.Equal(t, oldLWID.ConnID()+1, newLWID.ConnID(), "Connection ID should be recreated")
}

func newLinksForTest(t *testing.T) (*Namespace, *Links[amqpwrap.AMQPSenderCloser]) {
	testParams := test.GetConnectionParamsForTest(t)
	cred, err := credential.New(nil)
	require.NoError(t, err)

	ns, err := NewNamespace(NamespaceWithTokenCredential(testParams.EventHubNamespace, cred))
	require.NoError(t, err)

	links := NewLinks(ns, fmt.Sprintf("%s/$management", testParams.EventHubLinksOnlyName), func(partitionID string) string {
		return fmt.Sprintf("%s/Partitions/%s", testParams.EventHubLinksOnlyName, partitionID)
	}, func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string, partitionID string) (AMQPSenderCloser, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return session.NewSender(ctx, entityPath, "0", &amqp.SenderOptions{
				SettlementMode:              to.Ptr(amqp.SenderSettleModeMixed),
				RequestedReceiverSettleMode: to.Ptr(amqp.ReceiverSettleModeFirst),
			})
		}
	})

	err = links.Retry(context.Background(), exported.EventConn, "test", "0", exported.RetryOptions{
		RetryDelay:    -1,
		MaxRetryDelay: time.Millisecond,
	}, func(ctx context.Context, innerLWID LinkWithID[AMQPSenderCloser]) error {
		return nil
	})
	require.NoError(t, err)

	return ns, links
}

func lwidToError[LinkT AMQPLink](err error, lwid LinkWithID[LinkT]) error {
	return amqpwrap.WrapError(err, lwid.ConnID(), lwid.Link().LinkName(), lwid.PartitionID())
}
