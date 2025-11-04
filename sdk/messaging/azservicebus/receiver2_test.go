// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/stretchr/testify/require"
)

// NOTE: tests should migrate into here since this is the actual _test package.

func TestRenewLocks(t *testing.T) {
	client := test.NewClient(t, test.NewClientArgs[azservicebus.ClientOptions, azservicebus.Client]{
		NewClientFromConnectionString: azservicebus.NewClientFromConnectionString,
		NewClient:                     azservicebus.NewClient,
	}, nil)

	queueName, cleanupQueue := test.CreateExpiringQueue(t, nil)

	defer cleanupQueue()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	defer test.RequireClose(t, sender)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	defer test.RequireClose(t, receiver)

	err = sender.SendMessage(context.Background(), &azservicebus.Message{
		Body: []byte("hello world"),
	}, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	origTime := *messages[0].LockedUntil

	errCh := make(chan error)

	go func() {
		err = autoRenewMessageLock(context.Background(), receiver, messages[0], 5*time.Second)
		errCh <- err
	}()

	time.Sleep(10 * time.Second) // make sure we renew a few times.

	log.Printf("Completing message")
	err = receiver.CompleteMessage(context.Background(), messages[0], nil)
	require.NoError(t, err)
	log.Printf("Completed message")

	require.Greater(t, *messages[0].LockedUntil, origTime)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		require.Fail(t, "goroutine took longer than 5 seconds to complete")
	}
}

func TestRenewLocksEarlyCancel(t *testing.T) {
	client := test.NewClient(t, test.NewClientArgs[azservicebus.ClientOptions, azservicebus.Client]{
		NewClientFromConnectionString: azservicebus.NewClientFromConnectionString,
		NewClient:                     azservicebus.NewClient,
	}, nil)

	queueName, cleanupQueue := test.CreateExpiringQueue(t, &atom.QueueDescription{
		// this time we're bumping up the time. We want to make sure we're cancelling.
		LockDuration: to.Ptr("PT1M"),
	})

	defer cleanupQueue()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	defer test.RequireClose(t, sender)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	defer test.RequireClose(t, receiver)

	err = sender.SendMessage(context.Background(), &azservicebus.Message{
		Body: []byte("hello world"),
	}, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	t.Run("ContextAlreadyCancelled", func(t *testing.T) {
		// pre-cancelled.
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		log.Printf("Auto renewing lock, but the context is already cancelled")
		err = autoRenewMessageLock(ctx, receiver, messages[0], time.Minute)
		log.Printf("Auto renewer returned %s", err)
		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("ContextCancelledAfterTenSeconds", func(t *testing.T) {
		// now cancel _after_ the first renewal has taken place (by sleeping a bit...)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		errCh := make(chan error)

		log.Printf("Auto renewing lock, context will cancel after 10 seconds")
		go func() {
			err = autoRenewMessageLock(ctx, receiver, messages[0], time.Minute)
			errCh <- err
		}()

		select {
		case err := <-errCh:
			log.Printf("Auto renewer returned %s", err)
			require.ErrorIs(t, err, context.DeadlineExceeded)
		case <-time.After(15 * time.Second):
			require.Fail(t, "goroutine took longer than 15 seconds to cancel")
		}
	})
}
