// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func ExampleClient_AcceptSessionForQueue() {
	sessionReceiver, err := client.AcceptSessionForQueue(context.TODO(), "exampleSessionQueue", "exampleSessionId", nil)

	if err != nil {
		panic(err)
	}

	defer sessionReceiver.Close(context.TODO())

	// session receivers work just like non-session receivers
	// with one difference - instead of a lock per message there is a lock
	// for the session itself.

	// Like a message lock, you'll want to periodically renew your session lock.
	err = sessionReceiver.RenewSessionLock(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	messages, err := sessionReceiver.ReceiveMessages(context.TODO(), 5, nil)

	if err != nil {
		panic(err)
	}

	for _, message := range messages {
		err = sessionReceiver.CompleteMessage(context.TODO(), message, nil)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Received message from session ID \"%s\" and completed it", *message.SessionID)
	}
}

func ExampleClient_AcceptSessionForSubscription() {
	sessionReceiver, err := client.AcceptSessionForSubscription(context.TODO(), "exampleTopic", "exampleSubscription", "exampleSessionId", nil)

	if err != nil {
		panic(err)
	}

	defer sessionReceiver.Close(context.TODO())

	// see ExampleClient_AcceptSessionForQueue() for usage of the SessionReceiver.
}

func ExampleClient_AcceptNextSessionForQueue() {
	for {
		func() {
			sessionReceiver, err := client.AcceptNextSessionForQueue(context.TODO(), "exampleSessionQueue", nil)

			if err != nil {
				var sbErr *azservicebus.Error

				if errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeTimeout {
					// there are no sessions available. This isn't fatal - we can use the client and
					// try to AcceptNextSessionForQueue() again.
					fmt.Printf("No session available\n")
					return
				} else {
					panic(err)
				}
			}

			defer sessionReceiver.Close(context.TODO())

			fmt.Printf("Session receiver was assigned session ID \"%s\"", sessionReceiver.SessionID())

			// see ExampleClient_AcceptSessionForQueue() for usage of the SessionReceiver.
		}()
	}
}

func ExampleClient_AcceptNextSessionForSubscription() {
	for {
		func() {
			sessionReceiver, err := client.AcceptNextSessionForSubscription(context.TODO(), "exampleTopicName", "exampleSubscriptionName", nil)

			if err != nil {
				var sbErr *azservicebus.Error

				if errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeTimeout {
					// there are no sessions available. This isn't fatal - we can use the client and
					// try to AcceptNextSessionForSubscription() again.
					fmt.Printf("No session available\n")
					return
				} else {
					panic(err)
				}
			}

			defer sessionReceiver.Close(context.TODO())

			fmt.Printf("Session receiver was assigned session ID \"%s\"", sessionReceiver.SessionID())

			// see AcceptSessionForSubscription() for some usage of the SessionReceiver itself.
		}()
	}
}
