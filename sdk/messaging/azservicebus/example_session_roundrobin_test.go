// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func ExampleClient_AcceptNextSessionForQueue_roundrobin() {
	var sessionEnabledQueueName = "exampleSessionQueue"

	for {
		// You can have multiple active session receivers, provided they're each receiving
		// from different sessions.
		//
		// AcceptNextSessionForQueue (or AcceptNextSessionForSubscription) makes it simple to implement
		// this pattern, consuming multiple session receivers in parallel.
		sessionReceiver, err := client.AcceptNextSessionForQueue(context.TODO(), sessionEnabledQueueName, nil)

		if err != nil {
			var sbErr *azservicebus.Error

			if errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeTimeout {
				fmt.Printf("No sessions available\n")

				// NOTE: you could also continue here, which will block and wait again for a
				// session to become available.
				break
			}

			panic(err)
		}

		fmt.Printf("Got receiving for session '%s'\n", sessionReceiver.SessionID())

		// consume the session
		go func() {
			defer func() {
				fmt.Printf("Closing receiver for session '%s'\n", sessionReceiver.SessionID())
				err := sessionReceiver.Close(context.TODO())

				if err != nil {
					panic(err)
				}
			}()

			// we're only reading a few messages here, but you can also receive in a loop
			messages, err := sessionReceiver.ReceiveMessages(context.TODO(), 10, nil)

			if err != nil {
				panic(err)
			}

			fmt.Printf("Received %d messages from session '%s'\n", len(messages), sessionReceiver.SessionID())

			for _, m := range messages {
				if err := processMessageFromSession(context.TODO(), sessionReceiver, m); err != nil {
					panic(err)
				}
			}
		}()
	}
}

func processMessageFromSession(ctx context.Context, receiver *azservicebus.SessionReceiver, receivedMsg *azservicebus.ReceivedMessage) error {
	fmt.Printf("Received message for session %s, message ID %s\n", *receivedMsg.SessionID, receivedMsg.MessageID)
	return receiver.CompleteMessage(ctx, receivedMsg, nil)
}
