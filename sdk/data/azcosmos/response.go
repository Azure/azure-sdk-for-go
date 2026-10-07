// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Response holds the values every Cosmos DB operation reports, whatever it operated on. It is
// embedded in the per-operation response types, which add the values specific to them.
type Response struct {
	// Diagnostics is a copy of the native operation diagnostics, if available.
	Diagnostics *Diagnostics

	// StatusCode is the final HTTP status reported by the native driver, or zero when unavailable.
	StatusCode int

	// SubStatus is the final Cosmos sub-status, or zero when none was reported.
	SubStatus int

	// AttemptCount is the number of native requests made for this operation, including retries.
	// It is zero when the driver did not attach diagnostics.
	AttemptCount uint32

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

// DatabaseResponse is the response from an operation on a database.
type DatabaseResponse struct {
	Response

	// DatabaseProperties is the database the operation addressed. It is the zero value for a
	// delete, which returns no body.
	DatabaseProperties DatabaseProperties
}

// ContainerResponse is the response from an operation on a container.
type ContainerResponse struct {
	Response

	// ContainerProperties is the container the operation addressed. It is the zero value for a
	// delete, which returns no body.
	ContainerProperties ContainerProperties
}

// QueryItemsResponse contains one query page. An empty page does not imply exhaustion.
type QueryItemsResponse struct {
	// Response includes the metadata-validation charge on the first fetch.
	Response

	// Items contains independently owned JSON values, including scalar SELECT VALUE results.
	Items [][]byte

	// SessionToken is the session token returned for this page.
	SessionToken SessionToken

	// IndexMetrics is the decoded index-utilization JSON, when requested and available.
	IndexMetrics string

	// QueryMetrics contains the service's query execution metrics, when requested and available.
	QueryMetrics string
}
