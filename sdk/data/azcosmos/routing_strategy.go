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
//
// It is not supported yet: a client configured with it fails to construct rather than silently
// leaving the order to the account. Use [PreferredRegions] until it is.
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

// errProximityRoutingUnsupported is returned when a client is created with [ProximityTo].
//
// Resolving a region to an order is a client-side table lookup, not something the driver does: the
// C ABI takes a region list, so the SDK has to supply one. The Rust SDK's table is ~10k generated
// lines and has not been ported yet, so asking for proximity routing fails rather than quietly
// leaving the order to the account.
//
// It is reported when the client is constructed, matching errTokenCredentialUnsupported: resolving
// the order needs no driver, so there is no reason to defer it to the first operation.
var errProximityRoutingUnsupported = &Error{
	Code: CodeClientError,
	Message: "ProximityTo is not supported yet; list the regions explicitly with PreferredRegions, " +
		"most preferred first",
}

// preferredRegionOrder resolves the strategy to the region order the driver takes. The zero value
// resolves to none, which leaves the order to the account.
func (r RoutingStrategy) preferredRegionOrder() ([]Region, error) {
	if r.proximityTo != "" {
		return nil, errProximityRoutingUnsupported
	}
	return r.preferredRegions, nil
}
