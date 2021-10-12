// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func init() {
	initExamples()

	if client != nil {
		sender, err := client.NewSender(exampleSessionQueue)
		exitOnError("Failed to create sender for session examples", err)

		for i := 0; i < 2; i++ {
			err = sender.SendMessage(context.Background(), &azservicebus.Message{
				SessionID: to.StringPtr("Example Session ID"),
			})
			exitOnError("Failed to send example message to session", err)
		}
	}
}

func ExampleClient_AcceptSessionForQueue() {
	sessionReceiver, err := client.AcceptSessionForQueue(context.TODO(), "exampleSessionQueue", "Example Session ID", nil)
	exitOnError("Failed to create session receiver", err)

	// session receivers function the same as any other receiver
	message, err := sessionReceiver.ReceiveMessage(context.TODO(), nil)
	exitOnError("Failed to receive a message", err)

	err = sessionReceiver.CompleteMessage(context.TODO(), message)
	exitOnError("Failed to complete message", err)

	fmt.Printf("Received message from session ID \"%s\" and completed it", *message.SessionID)
}

func ExampleClient_AcceptNextSessionForQueue() {
	sessionReceiver, err := client.AcceptNextSessionForQueue(context.TODO(), "exampleSessionQueue", nil)
	exitOnError("Failed to create session receiver", err)

	fmt.Printf("Session receiver was assigned session ID \"%s\"", sessionReceiver.SessionID())
}
