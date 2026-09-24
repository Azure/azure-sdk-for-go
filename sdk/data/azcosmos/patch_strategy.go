// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "fmt"

// PatchStrategy selects how the driver executes [ContainerClient.PatchItem].
//
// PatchStrategy is provisional. It may change or be removed before azcosmos/v2 reaches a stable
// release.
type PatchStrategy string

const (
	// PatchStrategyUnset inherits the strategy configured for the driver, runtime, or environment.
	// It is the zero value and is not the same as [PatchStrategyAuto], which explicitly selects
	// automatic strategy selection.
	PatchStrategyUnset PatchStrategy = ""

	// PatchStrategyAuto lets the driver choose server-side PATCH or client-side read-modify-write
	// execution based on instruction safety and service limits.
	PatchStrategyAuto PatchStrategy = "Auto"

	// PatchStrategyClientSide makes the driver use client-side read-modify-write execution.
	PatchStrategyClientSide PatchStrategy = "ClientSide"

	// PatchStrategyServerSide makes the driver send PATCH directly to the service. The request can
	// fail when it exceeds a service limit that automatic execution would handle client-side.
	PatchStrategyServerSide PatchStrategy = "ServerSide"
)

func (s PatchStrategy) validate() error {
	switch s {
	case PatchStrategyUnset,
		PatchStrategyAuto,
		PatchStrategyClientSide,
		PatchStrategyServerSide:
		return nil
	default:
		return fmt.Errorf("azcosmos: unknown patch strategy %q", s)
	}
}
