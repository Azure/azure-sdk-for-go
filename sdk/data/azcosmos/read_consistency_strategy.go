// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// ReadConsistencyStrategy selects how fresh a read must be.
//
// This is deliberately not the account's consistency level. The driver exposes a smaller set of
// read strategies, and a read can only relax what the account guarantees, never strengthen it. See
// https://learn.microsoft.com/azure/cosmos-db/consistency-levels for the account-level concept.
type ReadConsistencyStrategy string

const (
	// ReadConsistencyStrategyUnset inherits the strategy configured for the account, the runtime
	// and the environment. It is the zero value, and is not the same as
	// [ReadConsistencyStrategyDefault], which overrides that inheritance.
	ReadConsistencyStrategyUnset ReadConsistencyStrategy = ""

	// ReadConsistencyStrategyDefault reads with the default behavior for the account's
	// consistency level, rather than inheriting a strategy configured elsewhere.
	ReadConsistencyStrategyDefault ReadConsistencyStrategy = "Default"

	// ReadConsistencyStrategyEventual reads without any ordering guarantee. It is the cheapest
	// strategy and may serve stale data.
	ReadConsistencyStrategyEventual ReadConsistencyStrategy = "Eventual"

	// ReadConsistencyStrategySession reads under a session token, so a client observes its own
	// writes.
	ReadConsistencyStrategySession ReadConsistencyStrategy = "Session"

	// ReadConsistencyStrategyGlobalStrong reads the latest committed version across regions. It
	// is only available on accounts whose consistency level permits it.
	ReadConsistencyStrategyGlobalStrong ReadConsistencyStrategy = "GlobalStrong"
)
