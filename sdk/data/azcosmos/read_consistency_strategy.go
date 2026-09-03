// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "fmt"

// ReadConsistencyStrategy selects how fresh a read must be.
//
// This is deliberately not the account's consistency level. A strategy can relax a read or request
// a quorum/stronger read where the account permits it. See
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

	// ReadConsistencyStrategyLatestCommitted reads the latest committed version using a quorum
	// read, independently of the account's default consistency level. It does not use the session
	// token lane.
	ReadConsistencyStrategyLatestCommitted ReadConsistencyStrategy = "LatestCommitted"

	// ReadConsistencyStrategyGlobalStrong reads the latest committed version across regions. It
	// is only available on accounts whose consistency level permits it.
	ReadConsistencyStrategyGlobalStrong ReadConsistencyStrategy = "GlobalStrong"
)

// validate rejects values the numeric driver ABI cannot represent. The empty value is the only
// unset sentinel; unlike service string enums, an unknown value cannot be forwarded.
func (s ReadConsistencyStrategy) validate() error {
	switch s {
	case ReadConsistencyStrategyUnset,
		ReadConsistencyStrategyDefault,
		ReadConsistencyStrategyEventual,
		ReadConsistencyStrategySession,
		ReadConsistencyStrategyLatestCommitted,
		ReadConsistencyStrategyGlobalStrong:
		return nil
	default:
		return fmt.Errorf("azcosmos: unknown read consistency strategy %q", s)
	}
}
