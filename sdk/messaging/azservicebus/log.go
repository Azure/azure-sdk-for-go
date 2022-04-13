// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"

const (
	// EventConn is used whenever we create a connection or any links (ie: receivers, senders).
	EventConn log.Event = "azsb.Conn"

	// EventAuth is used when we're doing authentication/claims negotiation.
	EventAuth log.Event = "azsb.Auth"

	// EventReceiver represents operations that happen on Receivers.
	EventReceiver log.Event = "azsb.Receiver"

	// EventSender represents operations that happen on Senders.
	EventSender log.Event = "azsb.Sender"

	// EventAdmin is used for operations in the azservicebus/admin.Client
	EventAdmin log.Event = "azsb.Admin"
)
