// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

// autoRenewMessageLock runs a loop calling [azservicebus.Receiver.RenewLock] until the message has been completed, the lock
// has expired or the passed in context has been cancelled.
//
// This function returns nil if the message's lock has expired.
//
// - receiver - the receiver that received 'msg'
// - msg - the received message
// - configuredLockDuration - the lock duration, configured at the queue or subscription level.
func autoRenewMessageLock(ctx context.Context, receiver *azservicebus.Receiver, msg *azservicebus.ReceivedMessage, configuredLockDuration time.Duration) error {
	for {
		log.Printf("Renewing lock for %s", msg.MessageID)
		if err := receiver.RenewMessageLock(ctx, msg, nil); err != nil {
			// You get this error if your lock has expired _or_ if you've settled the message.
			// We can ignore it here because the error will also be returned by the settlement
			// method (CompleteMessage, AbandonMessage, etc..)
			if sbErr := (*azservicebus.Error)(nil); errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeLockLost {
				log.Printf("Lock was lost for message with ID %s", msg.MessageID)
				return nil
			}

			return err
		}
		log.Printf("Renewed lock for %s", msg.MessageID)

		select {
		case <-ctx.Done():
			log.Printf("Cancelling lock renewal for %s", msg.MessageID)
			return ctx.Err()
		case <-time.After(configuredLockDuration / 2):
		}
	}
}

func Example_autoRenewLocks() {
	// This is configurable, on the service. Change this value to match since it's used
	// to figure out how often to renew a lock.
	const queueOrSubscriptionLockDuration = 30 * time.Second

	messages, err := receiver.ReceiveMessages(context.TODO(), 1, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	for _, msg := range messages {
		go func(msg *azservicebus.ReceivedMessage) {
			//
			// Each goroutine will _eventually_ stop once its message has been settled.
			//
			// If you want to be more proactive you can create a separate context for each message
			// and cancel that when you settle.
			//
			err := autoRenewMessageLock(context.TODO(), receiver, msg, queueOrSubscriptionLockDuration)

			if err != nil {
				//  TODO: Update the following line with your application specific error handling logic
				log.Printf("ERROR: %s", err)
				return
			}
		}(msg)
	}

	// TODO: Process the messages.
	//
	// Our lock renewal goroutines will stop once RenewLock returns that the lock has been lost. This happens
	// if the lock has expired or if we've settled the message.
}
