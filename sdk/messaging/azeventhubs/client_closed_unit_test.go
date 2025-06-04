// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/exported"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestClientClosedBehavior(t *testing.T) {
	t.Run("Links return ErrorCodeClientClosed after close", func(t *testing.T) {
		fakeNS := &fakeNamespaceForLinks{}
		links := internal.NewLinks(fakeNS, "managementPath", func(partitionID string) string {
			return "entityPath"
		}, func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string, partitionID string) (*fakeAMQPReceiver, error) {
			return &fakeAMQPReceiver{}, nil
		})

		// Before closing, operations should work
		_, err := links.GetManagementLink(context.Background())
		require.NoError(t, err)

		// Close the links
		err = links.Close(context.Background())
		require.NoError(t, err)

		// After closing, operations should return ErrorCodeClientClosed
		_, err = links.GetManagementLink(context.Background())
		require.Error(t, err)

		var ehErr *exported.Error
		require.True(t, errors.As(err, &ehErr))
		require.Equal(t, exported.ErrorCodeClientClosed, ehErr.Code)
	})

	t.Run("namespace returns ErrorCodeClientClosed after permanent close", func(t *testing.T) {
		ns := &internal.Namespace{}

		// Close permanently
		err := ns.Close(context.Background(), true)
		require.NoError(t, err)

		// Check should return ErrorCodeClientClosed
		err = ns.Check()
		require.Error(t, err)

		var ehErr *exported.Error
		require.True(t, errors.As(err, &ehErr))
		require.Equal(t, exported.ErrorCodeClientClosed, ehErr.Code)
	})
}

// Fake implementations for testing
type fakeNamespaceForLinks struct{}

func (f *fakeNamespaceForLinks) Check() error { return nil }
func (f *fakeNamespaceForLinks) NegotiateClaim(ctx context.Context, entityPath string) (context.CancelFunc, <-chan struct{}, error) {
	return func() {}, make(<-chan struct{}), nil
}
func (f *fakeNamespaceForLinks) NewAMQPSession(ctx context.Context) (amqpwrap.AMQPSession, uint64, error) {
	return &fakeAMQPSession{}, 1, nil
}
func (f *fakeNamespaceForLinks) NewRPCLink(ctx context.Context, managementPath string) (amqpwrap.RPCLink, uint64, error) {
	return &fakeRPCLink{}, 1, nil
}
func (f *fakeNamespaceForLinks) GetEntityAudience(entityPath string) string { return "audience" }
func (f *fakeNamespaceForLinks) Recover(ctx context.Context, clientRevision uint64) error { return nil }
func (f *fakeNamespaceForLinks) Close(ctx context.Context, permanently bool) error { return nil }

type fakeAMQPReceiver struct {
	closed bool
}

func (f *fakeAMQPReceiver) Close(ctx context.Context) error {
	f.closed = true
	return nil
}

func (f *fakeAMQPReceiver) LinkName() string {
	return "fake-link"
}

type fakeAMQPSession struct{}

func (f *fakeAMQPSession) Close(ctx context.Context) error { return nil }
func (f *fakeAMQPSession) ConnID() uint64 { return 1 }
func (f *fakeAMQPSession) NewReceiver(ctx context.Context, source string, partitionID string, opts *amqp.ReceiverOptions) (amqpwrap.AMQPReceiverCloser, error) {
	return nil, nil
}
func (f *fakeAMQPSession) NewSender(ctx context.Context, target string, partitionID string, opts *amqp.SenderOptions) (amqpwrap.AMQPSenderCloser, error) {
	return nil, nil
}

type fakeRPCLink struct{}

func (f *fakeRPCLink) Close(ctx context.Context) error { return nil }
func (f *fakeRPCLink) ConnID() uint64 { return 1 }
func (f *fakeRPCLink) RPC(ctx context.Context, msg *amqp.Message) (*amqpwrap.RPCResponse, error) {
	return nil, nil
}
func (f *fakeRPCLink) LinkName() string { return "fake-rpc-link" }