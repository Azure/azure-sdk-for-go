// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
)

const (
	// EventEndpointManager logs related to endpoint management initialization,
	// failover, and endpoint priority recomputation
	EventEndpointManager azlog.Event = "azcosmos.EndpointManager"
)
