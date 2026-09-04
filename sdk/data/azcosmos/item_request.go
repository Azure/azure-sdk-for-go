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

// endToEndTimeout resolves the budget the driver should bound the operation by, given the caller's
// context and whatever they set explicitly.
//
// Cancelling the context stops the operation, but only once the driver reaches a point where it can
// notice. The driver's own budget is what guarantees the operation *terminates*: without one it is
// bounded only by per-attempt transport timeouts times a retry budget, which can park a caller for
// far longer than their deadline. Passing the deadline down is what makes one number bound the
// operation at every layer.
//
// An explicit setting wins. It is a different thing from a deadline — the caller is describing how
// long the operation may spend, not when they stop waiting — so it is not second-guessed.
//
// A context with no deadline leaves the driver's default in place. Choosing a default here rather
// than reporting one is a policy decision this package does not make.
func endToEndTimeout(ctx context.Context, explicit time.Duration) time.Duration {
	if explicit > 0 {
		return explicit
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		return 0
	}
	remaining := time.Until(deadline)
	if remaining <= 0 {
		// Already expired. The operation is about to fail on the context anyway, so this only has
		// to avoid handing the driver a zero or negative budget, which it would read as unset.
		return time.Millisecond
	}
	return remaining
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
