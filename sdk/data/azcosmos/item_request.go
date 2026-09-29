// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"time"
)

// operationKind identifies which operation the driver should build. The values mirror the driver's
// operation kinds so the binding can pass them through directly.
type operationKind int32

const (
	operationKindCreateItem  operationKind = 19
	operationKindReadItem    operationKind = 20
	operationKindUpsertItem  operationKind = 21
	operationKindReplaceItem operationKind = 22
	operationKindDeleteItem  operationKind = 23
	operationKindPatchItem   operationKind = 24
)

// preconditionKind identifies the conditional request header the driver should send. The values
// mirror the driver's precondition kinds so the binding can pass them through directly.
type preconditionKind int32

const (
	preconditionKindNone        preconditionKind = 0
	preconditionKindIfMatch     preconditionKind = 1
	preconditionKindIfNoneMatch preconditionKind = 2
)

// itemRequest describes one item operation, in Go types.
//
// It exists so that the operation methods stay free of build tags: they populate this, and whether
// it reaches the driver or a not-implemented stub is decided by which build is selected.
type itemRequest struct {
	kind                   operationKind
	databaseID             string
	containerID            string
	itemID                 string
	partitionKey           PartitionKey
	body                   []byte
	sessionToken           SessionToken
	options                OperationOptions
	patchStrategy          PatchStrategy
	patchMaxAttempts       *uint8
	patchTrackingID        PatchTrackingID
	patchTrackingCapacity  *uint16
	patchTrackingRetention *time.Duration

	preconditionKind preconditionKind
	preconditionETag string
}

// newDriverUnavailableError says what this build is missing, rather than reporting the operation as
// merely unimplemented: the operation is implemented, but not in a build that cannot reach the
// driver.
func newDriverUnavailableError() *Error {
	return &Error{
		Code: CodeClientError,
		Message: "azcosmos: this build cannot reach the Cosmos driver. " +
			"v2 requires CGO_ENABLED=1, a supported target, and a target-compatible C toolchain",
	}
}

// contextWithEndToEndTimeout starts an explicit operation budget at public operation entry, so
// lazy driver creation and container resolution consume the same budget as the item request.
func contextWithEndToEndTimeout(
	ctx context.Context,
	explicit time.Duration,
) (context.Context, context.CancelFunc) {
	if explicit <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, explicit)
}
