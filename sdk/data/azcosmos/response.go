// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Response holds the values every Cosmos DB operation reports, regardless of what it operated on.
// It is embedded in the per-operation response types.
//
// Unlike v1, it carries no *http.Response: v2 operations are executed by the Cosmos driver rather
// than by an azcore HTTP pipeline, so there is no HTTP response for the caller to inspect. The
// fields below carry the values that were previously read off response headers.
type Response struct {
	// RequestCharge is the number of request units the operation consumed.
	RequestCharge float32

	// ActivityID correlates the operation with server-side telemetry.
	ActivityID string

	// ETag is the entity tag of the resource the operation returned. It is empty for operations
	// that do not address a single resource.
	ETag azcore.ETag

	// SessionToken is the session token the operation produced. Pass it to a later operation to
	// read your own writes under session consistency.
	SessionToken SessionToken
}

// ItemResponse is the response from an operation on a single item.
type ItemResponse struct {
	Response

	// Value is the raw item content the service returned. It is nil when the operation did not
	// request a content response, and for operations that do not return an item.
	Value []byte
}
