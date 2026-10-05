// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "time"

// OperationOptions holds the driver-level settings every Cosmos DB operation accepts. It is
// carried by each operation's own options type, which adds the settings specific to it.
//
// It exposes a subset of the driver's per-operation settings, shared by the operation APIs
// instead of being restated for each. A setting can be present but inert for a given operation:
// [OperationOptions.EnableContentResponseOnWrite] means nothing to a read, for instance.
type OperationOptions struct {
	// ConsistencyStrategy selects how fresh a read must be. The zero value reads with whatever the
	// account's consistency level implies. It has no effect on writes.
	ConsistencyStrategy ReadConsistencyStrategy

	// EnableContentResponseOnWrite requests that a write return the resulting item. Nil uses the
	// client-level [ClientOptions.EnableContentResponseOnWrite], and a non-nil value overrides it
	// in either direction. Explicitly disabling it reduces network and CPU cost. It has no effect
	// on reads.
	EnableContentResponseOnWrite *bool

	// ExcludedRegions removes regions from consideration for this operation, in addition to any
	// the client is already avoiding. Unlike [RoutingStrategy], which only orders regions, this
	// keeps the operation away from them entirely.
	ExcludedRegions []Region

	// EndToEndTimeout bounds the whole operation, including the retries the driver performs on
	// the caller's behalf. Zero means no bound beyond the context's deadline.
	EndToEndTimeout time.Duration
}
