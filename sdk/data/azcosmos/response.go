// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Response holds the values every Cosmos DB operation reports, whatever it operated on. It is
// embedded in the per-operation response types, which add the values specific to them.
type Response struct {
	// RequestCharge is the number of request units the operation consumed
	// (`x-ms-request-charge`). See
	// https://learn.microsoft.com/azure/cosmos-db/request-units.
	RequestCharge float64

	// ActivityID correlates the operation with server-side telemetry (`x-ms-activity-id`).
	ActivityID string
}

// ItemResponse is the response from an operation on a single item.
type ItemResponse struct {
	Response

	// ETag is the entity tag of the item the operation addressed (`etag`). Use it to make a later
	// write conditional on the item not having changed.
	ETag azcore.ETag

	// SessionToken is the session token the operation produced (`x-ms-session-token`). Pass it to
	// a later operation to read your own writes under session consistency.
	SessionToken SessionToken

	// Value is the raw item content the service returned. It is nil when the operation did not
	// request a content response, and for operations that do not return an item.
	Value []byte
}
