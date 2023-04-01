// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/golang/mock/gomock"
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
	err = links.RecoverIfNeeded(context.Background(), "0", staleLWID, &amqp.LinkError{})
	require.NoError(t, err)
	require.Nil(t, links.links["0"], "closed link is removed from the cache")
	require.Equal(t, 1, receivers[0].CloseCalled, "original receiver is closed, and replaced")

	// trying to recover again is a no-op (if nothing is in the cache)
	err = links.RecoverIfNeeded(context.Background(), "0", staleLWID, &amqp.LinkError{})
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

	err = links.RecoverIfNeeded(context.Background(), "0", staleLWID, &amqp.LinkError{})
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

	err = links.RecoverIfNeeded(context.Background(), "0", lwid, &amqp.LinkError{})
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
	ctrl := gomock.NewController(t)
	ns := mock.NewMockNamespaceForAMQPLinks(ctrl)
	receiver := mock.NewMockAMQPReceiverCloser(ctrl)
	session := mock.NewMockAMQPSession(ctrl)

	negotiateClaimCtx, cancelNegotiateClaim := context.WithCancel(context.Background())

	ns.EXPECT().NegotiateClaim(mock.NotCancelled, gomock.Any()).Return(cancelNegotiateClaim, negotiateClaimCtx.Done(), nil)
	ns.EXPECT().NewAMQPSession(mock.NotCancelled).Return(session, uint64(1), nil)

	receiver.EXPECT().LinkName().Return("link1").AnyTimes()

	links := NewLinks(ns, "managementPath", func(partitionID string) string {
		return fmt.Sprintf("part:%s", partitionID)
	}, func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (amqpwrap.AMQPReceiverCloser, error) {
		return receiver, nil
	})

	require.NotNil(t, links.contextWithTimeoutFn, "sanity check, we are setting the context.WithTimeout func")
	links.contextWithTimeoutFn = mock.NewContextWithTimeoutForTests

	lwid, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)
	require.NotNil(t, links.links["0"])
	require.Equal(t, 1, len(links.links))

	// if the connection has closed in response to an error then it'll propagate it's error to
	// the children, including receivers. Which means closing the receiver here will _also_ return
	// a connection error.
	receiver.EXPECT().Close(mock.NotCancelledAndHasTimeout).Return(&amqp.ConnError{})

	ns.EXPECT().Recover(mock.NotCancelledAndHasTimeout, gomock.Any()).Return(nil)

	// initiate a connection level recovery
	err = links.RecoverIfNeeded(context.Background(), "0", lwid, &amqp.ConnError{})
	require.NoError(t, err)

	// we still cleanup what we can (including cancelling our background negotiate claim loop)
	require.ErrorIs(t, context.Canceled, negotiateClaimCtx.Err())
	require.Empty(t, links.links, "link is removed")
}

func TestLinks_LinkRecoveryUpgradedToConnectionRecovery(t *testing.T) {
	connectionRecoverCalled := 0

	fakeNS := &FakeNSForPartClient{
		RecoverFn: func(ctx context.Context, clientRevision uint64) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				connectionRecoverCalled++
				return nil
			}
		},
	}

	getLogsFn := test.CaptureLogsForTest()

	var nextID int
	var receivers []*FakeAMQPReceiver

	links := NewLinks(fakeNS, "managementPath", func(partitionID string) string {
		return fmt.Sprintf("part:%s", partitionID)
	},
		func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (*FakeAMQPReceiver, error) {
			nextID++
			receivers = append(receivers, &FakeAMQPReceiver{
				NameForLink: fmt.Sprintf("Link%d", nextID),
				CloseError:  amqpwrap.ErrConnResetNeeded, // simulate a connection level failure
			})
			return receivers[len(receivers)-1], nil
		})

	lwid, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)
	require.NotNil(t, lwid)
	require.NotNil(t, links.links["0"], "cache contains the newly created link for partition 0")

	err = links.RecoverIfNeeded(context.Background(), "0", lwid, &amqp.LinkError{})
	require.NoError(t, err)
	require.Nil(t, links.links["0"], "cache will no longer a link for partition 0")
	require.Equal(t, 1, connectionRecoverCalled, "Connection was recovered since link.Close() returned a connection level error")

	logs := getLogsFn()

	require.Equal(t, logs, []string{
		"[azeh.Conn] Creating link for partition ID '0'",
		"[azeh.Conn] (c:1,l:Link1,p:0): Succesfully created link for partition ID '0'",
		"[azeh.Conn] (c:1,l:Link1,p:0) Error when cleaning up old link for link recovery: connection must be reset, link/connection state may be inconsistent",
		"[azeh.Conn] Upgrading to connection reset for recovery instead of link"})
}

func TestLinks_closeWithTimeout(t *testing.T) {
	for _, errToReturn := range []error{context.DeadlineExceeded, context.Canceled} {
		t.Run(fmt.Sprintf("Close() cancels with error %v", errToReturn), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ns := mock.NewMockNamespaceForAMQPLinks(ctrl)
			receiver := mock.NewMockAMQPReceiverCloser(ctrl)
			session := mock.NewMockAMQPSession(ctrl)

			negotiateClaimCtx, cancelNegotiateClaim := context.WithCancel(context.Background())

			ns.EXPECT().NegotiateClaim(mock.NotCancelled, gomock.Any()).Return(cancelNegotiateClaim, negotiateClaimCtx.Done(), nil)
			ns.EXPECT().NewAMQPSession(mock.NotCancelled).Return(session, uint64(1), nil)
			ns.EXPECT().Recover(mock.Cancelled, gomock.Any()).Return(context.Canceled)

			receiver.EXPECT().LinkName().Return("link1").AnyTimes()

			links := NewLinks(ns, "managementPath", func(partitionID string) string {
				return fmt.Sprintf("part:%s", partitionID)
			}, func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (amqpwrap.AMQPReceiverCloser, error) {
				return receiver, nil
			})

			require.NotNil(t, links.contextWithTimeoutFn, "sanity check, we are setting the context.WithTimeout func")
			links.contextWithTimeoutFn = mock.NewContextWithTimeoutForTests

			lwid, err := links.GetLink(context.Background(), "0")
			require.NoError(t, err)

			// now set ourselves up so Close() is "slow" and we end up timing out, or
			// the user "cancels"
			receiver.EXPECT().Close(mock.NotCancelledAndHasTimeout).DoAndReturn(func(ctx context.Context) error {
				<-ctx.Done()
				return amqpwrap.HandleNewOrCloseError(ctx.Err())
			})

			// purposefully recover with what should be a link level recovery. However, the Close() failing
			// means we end up "upgrading" to a connection reset instead.
			err = links.RecoverIfNeeded(context.Background(), "0", lwid, &amqp.LinkError{})
			require.ErrorIs(t, err, amqpwrap.ErrConnResetNeeded)

			// the error that comes back when the link times out being closed can only
			// be fixed by a connection reset.
			require.Equal(t, RecoveryKindConn, GetRecoveryKind(err))

			// we still cleanup what we can (including cancelling our background negotiate claim loop)
			require.ErrorIs(t, context.Canceled, negotiateClaimCtx.Err())
		})
	}
}

func TestLinks_linkRecoveryOnly(t *testing.T) {
	ctrl := gomock.NewController(t)
	fakeNS := mock.NewMockNamespaceForAMQPLinks(ctrl)
	fakeReceiver := mock.NewMockAMQPReceiverCloser(ctrl)
	session := mock.NewMockAMQPSession(ctrl)

	negotiateClaimCtx, cancelNegotiateClaim := context.WithCancel(context.Background())

	fakeNS.EXPECT().NegotiateClaim(mock.NotCancelled, gomock.Any()).Return(
		cancelNegotiateClaim, negotiateClaimCtx.Done(), nil,
	)
	fakeNS.EXPECT().NewAMQPSession(mock.NotCancelled).Return(session, uint64(1), nil)

	fakeReceiver.EXPECT().LinkName().Return("link1").AnyTimes()

	// super important that when we close we're given a context that properly times out.
	// (in this test the Close(ctx) call doesn't time out)
	fakeReceiver.EXPECT().Close(mock.NotCancelledAndHasTimeout).Return(nil)

	links := NewLinks(fakeNS, "managementPath", func(partitionID string) string {
		return fmt.Sprintf("part:%s", partitionID)
	}, func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (amqpwrap.AMQPReceiverCloser, error) {
		return fakeReceiver, nil
	})

	links.contextWithTimeoutFn = mock.NewContextWithTimeoutForTests

	lwid, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)

	err = links.RecoverIfNeeded(context.Background(), "0", lwid, &amqp.LinkError{})
	require.NoError(t, err)

	// we still cleanup what we can (including cancelling our background negotiate claim loop)
	require.ErrorIs(t, context.Canceled, negotiateClaimCtx.Err())
}

func TestLinks_linkRecoveryFailsWithLinkFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	fakeNS := mock.NewMockNamespaceForAMQPLinks(ctrl)
	fakeReceiver := mock.NewMockAMQPReceiverCloser(ctrl)
	session := mock.NewMockAMQPSession(ctrl)

	negotiateClaimCtx, cancelNegotiateClaim := context.WithCancel(context.Background())

	fakeNS.EXPECT().NegotiateClaim(mock.NotCancelled, gomock.Any()).Return(
		cancelNegotiateClaim, negotiateClaimCtx.Done(), nil,
	)
	fakeNS.EXPECT().NewAMQPSession(mock.NotCancelled).Return(session, uint64(1), nil)

	fakeReceiver.EXPECT().LinkName().Return("link1").AnyTimes()

	// super important that when we close we're given a context that properly times out.
	// (in this test the Close(ctx) call doesn't time out)
	detachErr := &amqp.LinkError{RemoteErr: &amqp.Error{Condition: amqp.ErrCondDetachForced}}
	fakeReceiver.EXPECT().Close(mock.NotCancelledAndHasTimeout).Return(detachErr)

	links := NewLinks(fakeNS, "managementPath", func(partitionID string) string {
		return fmt.Sprintf("part:%s", partitionID)
	}, func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (amqpwrap.AMQPReceiverCloser, error) {
		return fakeReceiver, nil
	})

	links.contextWithTimeoutFn = mock.NewContextWithTimeoutForTests

	lwid, err := links.GetLink(context.Background(), "0")
	require.NoError(t, err)

	err = links.RecoverIfNeeded(context.Background(), "0", lwid, &amqp.LinkError{})
	require.Equal(t, err, detachErr)

	// we still cleanup what we can (including cancelling our background negotiate claim loop)
	require.ErrorIs(t, context.Canceled, negotiateClaimCtx.Err())
}
