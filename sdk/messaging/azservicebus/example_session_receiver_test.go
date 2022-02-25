// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"fmt"
)

func ExampleClient_AcceptSessionForQueue() {
	sessionReceiver, err := client.AcceptSessionForQueue(context.TODO(), "exampleSessionQueue", "Example Session ID", nil)
	exitOnError("Failed to create session receiver", err)

	// session receivers function the same as any other receiver
	messages, err := sessionReceiver.ReceiveMessages(context.TODO(), 5, nil)
	exitOnError("Failed to receive a message", err)

	for _, message := range messages {
		err = sessionReceiver.CompleteMessage(context.TODO(), message)
		exitOnError("Failed to complete message", err)

		fmt.Printf("Received message from session ID \"%s\" and completed it", *message.SessionID)
	}
}

func ExampleClient_AcceptNextSessionForQueue() {
	sessionReceiver, err := client.AcceptNextSessionForQueue(context.TODO(), "exampleSessionQueue", nil)
	exitOnError("Failed to create session receiver", err)

	fmt.Printf("Session receiver was assigned session ID \"%s\"", sessionReceiver.SessionID())
}
