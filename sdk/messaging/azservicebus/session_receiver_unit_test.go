// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/exported"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestSessionReceiverUserFacingErrors(t *testing.T) {
	fakeAMQPLinks := &internal.FakeAMQPLinks{
		Err: amqp.ErrConnClosed,
	}

	receiver, err := newSessionReceiver(context.Background(), newSessionReceiverArgs{
		ns: &internal.FakeNS{
			AMQPLinks: fakeAMQPLinks,
		},
		retryOptions:   noRetriesNeeded,
		sessionID:      to.Ptr("session ID"),
		entity:         entity{Queue: "queue"},
		cleanupOnClose: func() {},
	}, nil)

	require.Nil(t, receiver)
	var asSBError *Error
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	// now let's let the receiver work for the first shot and then we'll invalidate it.
	fakeAMQPLinks.Err = nil
	fakeAMQPLinks.Receiver = &internal.FakeAMQPReceiver{}
	fakeRPCLink := &internal.FakeRPCLink{}
	fakeAMQPLinks.RPC = fakeRPCLink

	fakeRPCLink.Resp = &internal.RPCResponse{
		Message: &amqp.Message{
			Value: map[string]interface{}{
				"expiration": time.Now(),
			},
		},
	}

	receiver, err = newSessionReceiver(context.Background(), newSessionReceiverArgs{
		ns: &internal.FakeNS{
			AMQPLinks: fakeAMQPLinks,
		},
		retryOptions:   noRetriesNeeded,
		sessionID:      to.Ptr("session ID"),
		entity:         entity{Queue: "queue"},
		cleanupOnClose: func() {},
	}, nil)

	require.NoError(t, err)
	require.NotNil(t, receiver)

	fakeRPCLink.Resp = nil
	fakeRPCLink.Error = internal.RPCError{Resp: &internal.RPCResponse{Code: internal.RPCResponseCodeLockLost}}

	state, err := receiver.GetSessionState(context.Background(), nil)
	require.Nil(t, state)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.SetSessionState(context.Background(), []byte{}, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.RenewSessionLock(context.Background(), nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	// there's an init() method that's a little harder to trigger, so we'll do that here.
	// Unlike the others above it doesn't rely on the management$ link.
	fakeAMQPLinks.Err = amqp.ErrConnClosed
	err = receiver.init(context.Background())
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)
}

var noRetriesNeeded = exported.RetryOptions{
	MaxRetries:    0,
	RetryDelay:    0,
	MaxRetryDelay: 0,
}
