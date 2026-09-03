// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// The converters write into C structs through unsafe pointers, so a wrong union offset or kind is
// silent memory corruption rather than a compile error. These read back what was actually written
// to the C memory, through the inspectors in convert_native.go, which is the only way to catch that
// without a live service.

func TestPartitionKeyToNativeWritesEveryKind(t *testing.T) {
	pk := NewPartitionKey().
		AppendString("hello").
		AppendNumber(42.5).
		AppendBool(true).
		AppendNull().
		AppendUndefined()

	components, release := inspectNativePartitionKey(pk)
	t.Cleanup(release)
	require.Len(t, components, 5)

	require.Equal(t, nativeKindString, components[0].kind)
	require.Equal(t, "hello", components[0].stringValue)

	require.Equal(t, nativeKindNumber, components[1].kind)
	require.InDelta(t, 42.5, components[1].numberValue, 0)

	require.Equal(t, nativeKindBool, components[2].kind)
	require.True(t, components[2].boolValue)

	// Null and undefined are distinct kinds carrying no value, which is the distinction v1 could
	// not express.
	require.Equal(t, nativeKindNull, components[3].kind)
	require.Equal(t, nativeKindUndefined, components[4].kind)
}

// Component order is the partition key's order: a hierarchical key routes by it, so a converter
// that reordered them would route to the wrong partition.
func TestPartitionKeyToNativePreservesOrder(t *testing.T) {
	pk := NewPartitionKey().AppendString("tenant").AppendString("user").AppendString("session")

	components, release := inspectNativePartitionKey(pk)
	t.Cleanup(release)

	require.Len(t, components, 3)
	require.Equal(t, "tenant", components[0].stringValue)
	require.Equal(t, "user", components[1].stringValue)
	require.Equal(t, "session", components[2].stringValue)
}

// An empty partition key has to produce a null array rather than a zero-length allocation, because
// the request struct treats a null components pointer as "use the handle instead".
func TestPartitionKeyToNativeIsNilWhenEmpty(t *testing.T) {
	components, release := inspectNativePartitionKey(PartitionKey{})
	t.Cleanup(release)
	require.Nil(t, components)
}

// Releasing has to be safe for every shape, including one with no strings to free.
func TestPartitionKeyReleaseIsSafe(t *testing.T) {
	for _, pk := range []PartitionKey{
		{},
		NewPartitionKeyString("a"),
		NewPartitionKey().AppendNull(),
		NewPartitionKey().AppendString("a").AppendString("b").AppendString("c"),
	} {
		_, release := pk.toNative()
		require.NotPanics(t, release)
	}
}

// Unset options must leave the driver's defaults in place rather than zeroing every field, which
// is why the struct is built from cosmos_operation_options_default.
func TestOperationOptionsToNativeKeepsDefaults(t *testing.T) {
	options, release := inspectNativeOperationOptions(OperationOptions{})
	t.Cleanup(release)

	require.Equal(t, defaultNativeOperationOptions(), options)
}

// The content-response setting is tri-state at the ABI, so false has to be distinguishable from
// unset: 1 means false and 2 means true, and 0 leaves the client-level setting to apply.
func TestOperationOptionsContentResponseIsTriState(t *testing.T) {
	for _, tt := range []struct {
		name    string
		enabled *bool
		want    int32
	}{
		{"unset", nil, 0},
		{"disabled", to(false), 1},
		{"enabled", to(true), 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			options, release := inspectNativeOperationOptions(
				OperationOptions{EnableContentResponseOnWrite: tt.enabled})
			t.Cleanup(release)

			require.Equal(t, tt.want, options.contentResponseOnWrite)
		})
	}
}

func TestOperationOptionsToNativeCarriesEveryField(t *testing.T) {
	strategy, ok := nativeReadConsistencyStrategy(ReadConsistencyStrategySession)
	require.True(t, ok)

	options, release := inspectNativeOperationOptions(OperationOptions{
		ConsistencyStrategy: ReadConsistencyStrategySession,
		EndToEndTimeout:     3 * time.Second,
		ExcludedRegions:     []Region{RegionEastUS, RegionWestEurope},
	})
	t.Cleanup(release)

	require.Equal(t, strategy, options.readConsistencyStrategy)
	// The ABI takes milliseconds, so a Duration has to be converted rather than passed through.
	require.Equal(t, int64(3000), options.endToEndTimeoutMillis)
	require.Equal(t, []string{string(RegionEastUS), string(RegionWestEurope)}, options.excludedRegions)
}

func TestOperationOptionsClampsSubMillisecondTimeout(t *testing.T) {
	options, release := inspectNativeOperationOptions(OperationOptions{
		EndToEndTimeout: time.Nanosecond,
	})
	t.Cleanup(release)

	require.Equal(t, int64(1), options.endToEndTimeoutMillis)
}

// Every strategy has to map to a distinct discriminant, which catches both a mis-mapping and two
// cases having been given the same value.
func TestReadConsistencyStrategyToNativeIsInjective(t *testing.T) {
	unset, ok := nativeReadConsistencyStrategy(ReadConsistencyStrategyUnset)
	require.False(t, ok, "unset leaves the driver's default in place")
	require.Zero(t, unset)

	seen := make(map[int32]ReadConsistencyStrategy)
	for _, strategy := range []ReadConsistencyStrategy{
		ReadConsistencyStrategyDefault,
		ReadConsistencyStrategyEventual,
		ReadConsistencyStrategySession,
		ReadConsistencyStrategyLatestCommitted,
		ReadConsistencyStrategyGlobalStrong,
	} {
		value, ok := nativeReadConsistencyStrategy(strategy)
		require.True(t, ok, "%q should map", strategy)

		previous, duplicated := seen[value]
		require.False(t, duplicated, "%q and %q both map to %d", previous, strategy, value)
		seen[value] = strategy
	}
}

// The request struct is built as a Go composite literal, so every field it does not name takes
// Go's zero value. That is correct for the fields whose unset value is zero and wrong for the ones
// whose is not: the driver rejects a zero max-item-count outright, with an invalid-option-value
// status raised before any network I/O, so leaving it zero fails every operation rather than
// degrading quietly.
//
// This pins the sentinel for both item kinds. A new field with a non-zero unset value has to be
// added here as well as to newOperationRequest.
func TestOperationRequestUsesTheDriversUnsetSentinels(t *testing.T) {
	for _, kind := range []operationKind{operationKindReadItem, operationKindCreateItem} {
		maxItemCount := inspectRequestSentinels(kind)
		require.Negative(t, maxItemCount,
			"the driver reads < 0 as unset and rejects 0, so a zero here fails the operation")
	}
}

// Client options reach the driver through a flat config, so the conversion is what decides whether
// a setting has any effect at all.
func TestClientOptionsConvertToTheDriversConfig(t *testing.T) {
	t.Run("preferred regions are passed in order", func(t *testing.T) {
		options, release, err := inspectNativeClientOptions(ClientOptions{
			Routing: PreferredRegions(RegionWestUS, RegionEastUS),
		})
		require.NoError(t, err)
		defer release()

		require.Equal(t, []string{string(RegionWestUS), string(RegionEastUS)}, options.preferredRegions,
			"the order is the preference, so it has to survive the conversion")
	})

	t.Run("no routing leaves the order to the account", func(t *testing.T) {
		options, release, err := inspectNativeClientOptions(ClientOptions{})
		require.NoError(t, err)
		defer release()

		require.Empty(t, options.preferredRegions)
	})

	t.Run("proximity routing is reported rather than ignored", func(t *testing.T) {
		_, release, err := inspectNativeClientOptions(ClientOptions{Routing: ProximityTo(RegionEastUS)})
		defer release()

		require.ErrorIs(t, err, errProximityRoutingUnsupported,
			"silently leaving the order to the account would hide that the request had no effect")
	})

	// The client-level value is sent explicitly rather than left unset, so the documented Go
	// default holds even if the driver's default changes.
	t.Run("the content response setting is sent explicitly", func(t *testing.T) {
		for _, tt := range []struct {
			name    string
			enabled bool
			want    int32
		}{
			{"disabled", false, 1},
			{"enabled", true, 2},
		} {
			t.Run(tt.name, func(t *testing.T) {
				options, release, err := inspectNativeClientOptions(ClientOptions{
					EnableContentResponseOnWrite: tt.enabled,
				})
				require.NoError(t, err)
				defer release()

				require.Equal(t, tt.want, options.operationOptions.contentResponseOnWrite)
			})
		}
	})
}
