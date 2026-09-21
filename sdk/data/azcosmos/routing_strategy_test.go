// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProximityTo(t *testing.T) {
	strategy := ProximityTo(RegionEastUS)

	require.Equal(t, RegionEastUS, strategy.proximityTo)
	require.Empty(t, strategy.preferredRegions)
}

func TestPreferredRegions(t *testing.T) {
	regions := []Region{RegionWestUS, RegionEastUS}
	strategy := PreferredRegions(regions...)

	require.Empty(t, strategy.proximityTo)
	require.Equal(t, regions, strategy.preferredRegions)

	regions[0] = RegionNorthEurope
	require.Equal(t, []Region{RegionWestUS, RegionEastUS}, strategy.preferredRegions)
}

func TestRoutingStrategyZeroValue(t *testing.T) {
	var strategy RoutingStrategy

	require.Empty(t, strategy.proximityTo)
	require.Empty(t, strategy.preferredRegions)
}

func TestProximityRegionTableMatchesRust(t *testing.T) {
	require.Len(t, proximityRegionOrderBySource, 96)

	sources := make([]string, 0, len(proximityRegionOrderBySource))
	for source := range proximityRegionOrderBySource {
		sources = append(sources, string(source))
	}
	sort.Strings(sources)

	hash := sha256.New()
	for _, sourceName := range sources {
		source := Region(sourceName)
		regions := proximityRegionOrderBySource[source]

		require.Len(t, regions, 96, "source %q", source)
		require.Equal(t, source, regions[0], "source %q", source)

		seen := make(map[Region]struct{}, len(regions))
		for _, region := range regions {
			_, duplicate := seen[region]
			require.False(t, duplicate, "source %q repeats region %q", source, region)
			seen[region] = struct{}{}
		}

		hash.Write([]byte(source))
		hash.Write([]byte{0})
		for _, region := range regions {
			hash.Write([]byte(region))
			hash.Write([]byte{0})
		}
	}

	require.Equal(t,
		"0694969409a58ea7881a88f42ac8626227452f331a48da5976aff07159b06bc6",
		fmt.Sprintf("%x", hash.Sum(nil)))
}

func TestProximityToKnownRegionUsesEstimatedOrder(t *testing.T) {
	regions, err := ProximityTo(RegionEastUS).preferredRegionOrder()
	require.NoError(t, err)
	require.Equal(t, []Region{
		RegionEastUS,
		RegionEastUS2,
		RegionEastUS3,
		RegionNorthCentralUS,
		RegionNortheastUS5,
		RegionCentralUS,
		RegionSoutheastUS5,
		RegionCanadaCentral,
	}, regions[:8])
}

func TestProximityRegionOrderCannotBeMutated(t *testing.T) {
	regions, err := ProximityTo(RegionEastUS).preferredRegionOrder()
	require.NoError(t, err)
	regions[0] = RegionWestUS

	again, err := ProximityTo(RegionEastUS).preferredRegionOrder()
	require.NoError(t, err)
	require.Equal(t, RegionEastUS, again[0])
}

func TestProximityToUnknownRegionLeavesAccountOrder(t *testing.T) {
	regions, err := ProximityTo("not-a-real-region").preferredRegionOrder()
	require.NoError(t, err)
	require.Empty(t, regions)
}

func TestOtherRoutingStrategiesKeepTheirOrder(t *testing.T) {
	regions, err := (RoutingStrategy{}).preferredRegionOrder()
	require.NoError(t, err)
	require.Empty(t, regions)

	regions, err = PreferredRegions(RegionWestUS, RegionEastUS).preferredRegionOrder()
	require.NoError(t, err)
	require.Equal(t, []Region{RegionWestUS, RegionEastUS}, regions)
}
