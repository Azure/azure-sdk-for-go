// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"crypto/sha256"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strconv"
	"strings"
	"testing"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/stretchr/testify/require"
)

var regionsWithoutProximityData = map[Region]struct{}{
	RegionGermanyCentral:    {},
	RegionGermanyNortheast:  {},
	RegionUSGOVIowa:         {},
	RegionNorthEurope2:      {},
	RegionEastEurope:        {},
	RegionAPACSoutheast2:    {},
	RegionUKSouth2:          {},
	RegionUKNorth:           {},
	RegionEastUSSTG:         {},
	RegionSouthCentralUSSTG: {},
	RegionUSGOVWyoming:      {},
	RegionUSDODSouthwest:    {},
	RegionUSDODWestCentral:  {},
	RegionUSDODSouthCentral: {},
	RegionChinaNorth10:      {},
	RegionKoreaSouth2:       {},
}

func TestProximityTo(t *testing.T) {
	strategy := ProximityTo(RegionEastUS)

	require.Equal(t, RegionEastUS, strategy.proximityTo)
	require.True(t, strategy.proximitySet)
	require.Empty(t, strategy.preferredRegions)
}

func TestPreferredRegions(t *testing.T) {
	regions := []Region{RegionWestUS, RegionEastUS}
	strategy := PreferredRegions(regions...)

	require.Empty(t, strategy.proximityTo)
	require.False(t, strategy.proximitySet)
	require.Equal(t, regions, strategy.preferredRegions)

	regions[0] = RegionNorthEurope
	require.Equal(t, []Region{RegionWestUS, RegionEastUS}, strategy.preferredRegions)
}

func TestRoutingStrategyZeroValue(t *testing.T) {
	var strategy RoutingStrategy

	require.Empty(t, strategy.proximityTo)
	require.False(t, strategy.proximitySet)
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

func TestExportedRegionConstantsHaveProximityDataOrDocumentedException(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "region.go", nil, 0)
	require.NoError(t, err)

	foundExceptions := make(map[Region]struct{}, len(regionsWithoutProximityData))
	for _, declaration := range file.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || general.Tok != token.CONST {
			continue
		}
		for _, specification := range general.Specs {
			value, ok := specification.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range value.Names {
				if !name.IsExported() || !strings.HasPrefix(name.Name, "Region") || i >= len(value.Values) {
					continue
				}
				literal, ok := value.Values[i].(*ast.BasicLit)
				require.True(t, ok, "constant %s must remain a string literal", name.Name)
				regionValue, err := strconv.Unquote(literal.Value)
				require.NoError(t, err)
				region := Region(regionValue)

				if _, excepted := regionsWithoutProximityData[region]; excepted {
					require.NotContains(t, proximityRegionOrderBySource, region,
						"remove %s from regionsWithoutProximityData now that it has proximity data", name.Name)
					foundExceptions[region] = struct{}{}
					continue
				}
				require.Contains(t, proximityRegionOrderBySource, region,
					"%s needs proximity data or an explicit documented exception", name.Name)
			}
		}
	}
	require.Equal(t, regionsWithoutProximityData, foundExceptions,
		"every documented exception must still be an exported Region constant")
}

func TestProximityToKnownRegionUsesEstimatedOrder(t *testing.T) {
	want := []Region{
		RegionEastUS,
		RegionEastUS2,
		RegionEastUS3,
		RegionNorthCentralUS,
		RegionNortheastUS5,
		RegionCentralUS,
		RegionSoutheastUS5,
		RegionCanadaCentral,
	}
	for _, region := range []Region{RegionEastUS, "EASTUS", "East US", "East\tUS"} {
		regions, err := ProximityTo(region).preferredRegionOrder()
		require.NoError(t, err)
		require.Equal(t, want, regions[:len(want)], "source %q", region)
	}
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
	var event azlog.Event
	var message string
	azlog.SetEvents(EventRouting)
	azlog.SetListener(func(receivedEvent azlog.Event, receivedMessage string) {
		event = receivedEvent
		message = receivedMessage
	})
	t.Cleanup(func() {
		azlog.SetListener(nil)
		azlog.SetEvents()
	})

	regions, err := ProximityTo("not-a-real-region").preferredRegionOrder()
	require.NoError(t, err)
	require.Empty(t, regions)
	require.Equal(t, EventRouting, event)
	require.Contains(t, message, `unrecognized application region "not-a-real-region"`)
	require.Contains(t, message, "falling back to account-defined region order")
}

func TestProximityToEmptyRegionLeavesAccountOrder(t *testing.T) {
	var event azlog.Event
	var message string
	azlog.SetEvents(EventRouting)
	azlog.SetListener(func(receivedEvent azlog.Event, receivedMessage string) {
		event = receivedEvent
		message = receivedMessage
	})
	t.Cleanup(func() {
		azlog.SetListener(nil)
		azlog.SetEvents()
	})

	regions, err := ProximityTo("").preferredRegionOrder()
	require.NoError(t, err)
	require.Empty(t, regions)
	require.Equal(t, EventRouting, event)
	require.Contains(t, message, `unrecognized application region ""`)
	require.Contains(t, message, "falling back to account-defined region order")
}

func TestOtherRoutingStrategiesKeepTheirOrder(t *testing.T) {
	regions, err := (RoutingStrategy{}).preferredRegionOrder()
	require.NoError(t, err)
	require.Empty(t, regions)

	regions, err = PreferredRegions(RegionWestUS, RegionEastUS).preferredRegionOrder()
	require.NoError(t, err)
	require.Equal(t, []Region{RegionWestUS, RegionEastUS}, regions)
}
