// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// RoutingStrategy decides the order in which a client considers the account's regions.
//
// It orders regions, it does not restrict them: once failover exhausts the order, the client may
// still use any remaining region. Use [ReadItemOptions.ExcludedRegions] and its equivalents to
// keep an operation away from a region entirely.
//
// Build one with [ProximityTo] or [PreferredRegions]. The zero value leaves the order to the
// account, which is rarely what a latency-sensitive application wants.
type RoutingStrategy struct {
	proximityTo      Region
	preferredRegions []Region
}

// ProximityTo orders regions by estimated geographic proximity to the given region, which should
// be where the application runs, or the nearest Azure region to it.
//
// The estimates are built into the SDK and may not match the round-trip times actually observed.
func ProximityTo(region Region) RoutingStrategy {
	return RoutingStrategy{proximityTo: region}
}

// PreferredRegions orders regions explicitly, most preferred first.
func PreferredRegions(regions ...Region) RoutingStrategy {
	return RoutingStrategy{preferredRegions: append([]Region(nil), regions...)}
}

// clone copies the region slice so a caller mutating theirs after constructing a client cannot
// silently change how that client routes.
func (r RoutingStrategy) clone() RoutingStrategy {
	r.preferredRegions = append([]Region(nil), r.preferredRegions...)
	return r
}
