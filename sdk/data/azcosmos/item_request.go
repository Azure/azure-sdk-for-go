// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// operationKind identifies which operation the driver should build. The values mirror the driver's
// operation kinds so the binding can pass them through directly.
type operationKind int32

const (
	operationKindCreateItem operationKind = 19
	operationKindReadItem   operationKind = 20
)

// itemRequest describes one item operation, in Go types.
//
// It exists so that the operation methods stay free of build tags: they populate this, and whether
// it reaches the driver or a not-implemented stub is decided by which build is selected.
type itemRequest struct {
	kind         operationKind
	databaseID   string
	containerID  string
	itemID       string
	partitionKey PartitionKey
	body         []byte
	sessionToken SessionToken
	options      OperationOptions

	// ifNoneMatchETag is the conditional-read precondition. Empty means unconditional.
	ifNoneMatchETag string
}
