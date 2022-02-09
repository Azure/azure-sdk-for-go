// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"

const (
	// Link/connection creation
	EventConn = "azsb.Conn"

	// authentication/claims negotiation
	EventAuth = "azsb.Auth"

	// receiver operations
	EventReceiver = "azsb.Receiver"

	// mgmt link
	EventMgmtLink = "azsb.Mgmt"

	// internal operations
	EventRetry = utils.EventRetry
)
