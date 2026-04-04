// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	MinimumInclusiveEffectivePartitionKey = ""
	MaximumExclusiveEffectivePartitionKey = "FF"
)

type mockCollectionCache struct {
	currentCollection          *DocumentCollection
	collectionAfterRefresh     *DocumentCollection
	refreshCount               int
	returnNotFoundAfterRefresh bool
}

func (m *mockCollectionCache) ResolveCollectionAsync(ctx context.Context, request *AddressResolverRequest) (*DocumentCollection, error) {
	if request.ForceNameCacheRefresh {
		m.refreshCount++
		request.ForceNameCacheRefresh = false

		if m.collectionAfterRefresh == nil {
			m.currentCollection = nil
			m.returnNotFoundAfterRefresh = true
			return nil, NewNotFoundException("Collection not found")
		}
		m.currentCollection = m.collectionAfterRefresh
		return m.currentCollection, nil
	}

	if m.returnNotFoundAfterRefresh || m.currentCollection == nil {
		return nil, NewNotFoundException("Collection not found")
	}

	return m.currentCollection, nil
}

type mockRoutingMapCache struct {
	currentRoutingMaps      map[string]CollectionRoutingMap
	routingMapsAfterRefresh map[string]CollectionRoutingMap
	refreshCount            map[string]int
}

func newMockRoutingMapCache() *mockRoutingMapCache {
	return &mockRoutingMapCache{
		currentRoutingMaps:      make(map[string]CollectionRoutingMap),
		routingMapsAfterRefresh: make(map[string]CollectionRoutingMap),
		refreshCount:            make(map[string]int),
	}
}

func (m *mockRoutingMapCache) TryLookupAsync(ctx context.Context, collectionRID string, previousValue CollectionRoutingMap, forceRefresh bool) (CollectionRoutingMap, error) {
	if previousValue == nil {
		return m.currentRoutingMaps[collectionRID], nil
	}

	if current, ok := m.currentRoutingMaps[previousValue.GetCollectionUniqueID()]; ok && current == previousValue {
		if previousValue.GetCollectionUniqueID() != collectionRID {
			return nil, errors.New("invalid operation: collection ID mismatch")
		}

		if afterRefresh, ok := m.routingMapsAfterRefresh[collectionRID]; ok {
			m.currentRoutingMaps[collectionRID] = afterRefresh
		} else {
			delete(m.currentRoutingMaps, collectionRID)
		}

		m.refreshCount[collectionRID]++
		return m.currentRoutingMaps[collectionRID], nil
	}

	return nil, errors.New("not implemented in mock")
}

type mockAddressCache struct {
	currentAddresses      map[string][]AddressInformation
	addressesAfterRefresh map[string][]AddressInformation
	refreshCount          map[string]int
}

func newMockAddressCache() *mockAddressCache {
	return &mockAddressCache{
		currentAddresses:      make(map[string][]AddressInformation),
		addressesAfterRefresh: make(map[string][]AddressInformation),
		refreshCount:          make(map[string]int),
	}
}

func (m *mockAddressCache) findKeyForIdentity(addresses map[string][]AddressInformation, pkri *PartitionKeyRangeIdentity) string {
	key := pkri.CollectionRID + ":" + pkri.PartitionKeyRangeID
	if _, ok := addresses[key]; ok {
		return key
	}
	return ""
}

func (m *mockAddressCache) TryGetAddresses(ctx context.Context, request *AddressResolverRequest, partitionKeyRangeIdentity *PartitionKeyRangeIdentity, forceRefresh bool) ([]AddressInformation, error) {
	key := partitionKeyRangeIdentity.CollectionRID + ":" + partitionKeyRangeIdentity.PartitionKeyRangeID

	if !forceRefresh {
		return m.currentAddresses[key], nil
	}

	if afterRefresh, ok := m.addressesAfterRefresh[key]; ok {
		m.currentAddresses[key] = afterRefresh
	} else {
		delete(m.currentAddresses, key)
	}

	m.refreshCount[key]++
	return m.currentAddresses[key], nil
}

func createAddresses(uri string) []AddressInformation {
	return []AddressInformation{
		NewAddressInformation(true, true, uri, ProtocolHTTPS),
	}
}

func setupTestCollections() (*DocumentCollection, *DocumentCollection) {
	collection1 := &DocumentCollection{
		ID:         "coll",
		ResourceID: "rid1",
		PartitionKey: &PartitionKeyDefinition{
			Paths: []string{"/field1"},
		},
	}

	collection2 := &DocumentCollection{
		ID:         "coll",
		ResourceID: "rid2",
		PartitionKey: &PartitionKeyDefinition{
			Paths: []string{"/field1"},
		},
	}

	return collection1, collection2
}

func setupRoutingMaps(collection1, collection2 *DocumentCollection) (
	routingMapCollection1BeforeSplit CollectionRoutingMap,
	routingMapCollection1AfterSplit CollectionRoutingMap,
	routingMapCollection2BeforeSplit CollectionRoutingMap,
	routingMapCollection2AfterSplit CollectionRoutingMap,
	serviceIdentities map[string]*ServiceIdentity,
) {
	serviceIdentities = make(map[string]*ServiceIdentity)

	si1 := NewServiceIdentity("federation1", "fabric://serverservice1", false)
	range1 := &PartitionKeyRange{
		ID:           "0",
		MinInclusive: MinimumInclusiveEffectivePartitionKey,
		MaxExclusive: MaximumExclusiveEffectivePartitionKey,
	}
	si1.AddPartitionKeyRangeIdentity(NewPartitionKeyRangeIdentity(collection1.ResourceID, range1.ID))
	serviceIdentities["c1_before_0"] = si1

	routingMapCollection1BeforeSplit = NewInMemoryCollectionRoutingMap(
		[]*PartitionKeyRange{range1},
		[]IServerIdentity{si1},
		collection1.ResourceID,
	)

	si2 := NewServiceIdentity("federation1", "fabric://serverservice2", false)
	si3 := NewServiceIdentity("federation1", "fabric://serverservice3", false)
	range2 := &PartitionKeyRange{ID: "1", MinInclusive: MinimumInclusiveEffectivePartitionKey, MaxExclusive: "5E", Parents: []string{"0"}}
	range3 := &PartitionKeyRange{ID: "2", MinInclusive: "5E", MaxExclusive: MaximumExclusiveEffectivePartitionKey, Parents: []string{"0"}}
	si2.AddPartitionKeyRangeIdentity(NewPartitionKeyRangeIdentity(collection1.ResourceID, range2.ID))
	si3.AddPartitionKeyRangeIdentity(NewPartitionKeyRangeIdentity(collection1.ResourceID, range3.ID))
	serviceIdentities["c1_after_0"] = si2
	serviceIdentities["c1_after_1"] = si3

	routingMapCollection1AfterSplit = NewInMemoryCollectionRoutingMap(
		[]*PartitionKeyRange{range2, range3},
		[]IServerIdentity{si2, si3},
		collection1.ResourceID,
	)

	si4 := NewServiceIdentity("federation1", "fabric://serverservice4", false)
	range4 := &PartitionKeyRange{
		ID:           "0",
		MinInclusive: MinimumInclusiveEffectivePartitionKey,
		MaxExclusive: MaximumExclusiveEffectivePartitionKey,
	}
	si4.AddPartitionKeyRangeIdentity(NewPartitionKeyRangeIdentity(collection2.ResourceID, range4.ID))
	serviceIdentities["c2_before_0"] = si4

	routingMapCollection2BeforeSplit = NewInMemoryCollectionRoutingMap(
		[]*PartitionKeyRange{range4},
		[]IServerIdentity{si4},
		collection2.ResourceID,
	)

	si5 := NewServiceIdentity("federation1", "fabric://serverservice5", false)
	si6 := NewServiceIdentity("federation1", "fabric://serverservice6", false)
	range5 := &PartitionKeyRange{ID: "1", MinInclusive: MinimumInclusiveEffectivePartitionKey, MaxExclusive: "5E", Parents: []string{"0"}}
	range6 := &PartitionKeyRange{ID: "2", MinInclusive: "5E", MaxExclusive: MaximumExclusiveEffectivePartitionKey, Parents: []string{"0"}}
	si5.AddPartitionKeyRangeIdentity(NewPartitionKeyRangeIdentity(collection2.ResourceID, range5.ID))
	si6.AddPartitionKeyRangeIdentity(NewPartitionKeyRangeIdentity(collection2.ResourceID, range6.ID))
	serviceIdentities["c2_after_0"] = si5
	serviceIdentities["c2_after_1"] = si6

	routingMapCollection2AfterSplit = NewInMemoryCollectionRoutingMap(
		[]*PartitionKeyRange{range5, range6},
		[]IServerIdentity{si5, si6},
		collection2.ResourceID,
	)

	return
}

func TestCacheRefreshesWhileRoutingByPartitionKey(t *testing.T) {
	collection1, collection2 := setupTestCollections()
	routingMapC1Before, routingMapC1After, routingMapC2Before, routingMapC2After, serviceIdentities := setupRoutingMaps(collection1, collection2)

	addresses1 := createAddresses("tcp://host/partition1")
	addresses2 := createAddresses("tcp://host/partition2")
	addresses3 := createAddresses("tcp://host/partition3")

	testCases := []struct {
		name                      string
		collectionBeforeRefresh   *DocumentCollection
		collectionAfterRefresh    *DocumentCollection
		routingMapBeforeRefresh   map[string]CollectionRoutingMap
		routingMapAfterRefresh    map[string]CollectionRoutingMap
		addressesBeforeRefresh    map[string][]AddressInformation
		addressesAfterRefresh     map[string][]AddressInformation
		targetAddresses           []AddressInformation
		targetPartitionKeyRange   *PartitionKeyRange
		forceNameCacheRefresh     bool
		forceRoutingMapRefresh    bool
		forceAddressRefresh       bool
		expectedCollectionRefresh int
		expectedRoutingMapRefresh int
		expectedAddressRefresh    int
		nameBased                 bool
		expectError               error
	}{
		{
			name:                    "All caches are up to date. Name Based",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  nil,
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":0": addresses1,
			},
			addressesAfterRefresh:     nil,
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1Before.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    0,
			nameBased:                 true,
		},
		{
			name:                    "All caches are up to date. Rid Based",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  nil,
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":0": addresses1,
			},
			addressesAfterRefresh:     nil,
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1Before.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    0,
			nameBased:                 false,
		},
		{
			name:                    "Address cache is stale. Force Refresh. Name Based",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  nil,
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":0": addresses1,
			},
			addressesAfterRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":0": addresses2,
			},
			targetAddresses:           addresses2,
			targetPartitionKeyRange:   routingMapC1Before.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       true,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    1,
			nameBased:                 true,
		},
		{
			name:                    "Address cache is stale. Force Refresh. Rid Based",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  nil,
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":0": addresses1,
			},
			addressesAfterRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":0": addresses2,
			},
			targetAddresses:           addresses2,
			targetPartitionKeyRange:   routingMapC1Before.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       true,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    1,
			nameBased:                 false,
		},
		{
			name:                    "Routing map cache is stale. Force Refresh. Name based",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1After},
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":0": addresses1,
				collection1.ResourceID + ":1": addresses2,
				collection1.ResourceID + ":2": addresses3,
			},
			addressesAfterRefresh:     nil,
			targetAddresses:           addresses2,
			targetPartitionKeyRange:   routingMapC1After.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    true,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 1,
			expectedAddressRefresh:    0,
			nameBased:                 true,
		},
		{
			name:                    "Routing map cache is stale. Force Refresh. Rid based",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1After},
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":0": addresses1,
				collection1.ResourceID + ":1": addresses2,
				collection1.ResourceID + ":2": addresses3,
			},
			addressesAfterRefresh:     nil,
			targetAddresses:           addresses2,
			targetPartitionKeyRange:   routingMapC1After.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    true,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 1,
			expectedAddressRefresh:    0,
			nameBased:                 false,
		},
		{
			name:                    "Name cache is stale. Force Refresh. Name based",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection2,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection2.ResourceID: routingMapC2Before},
			routingMapAfterRefresh:  nil,
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection2.ResourceID + ":0": addresses1,
			},
			addressesAfterRefresh:     nil,
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC2Before.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     true,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 1,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    0,
			nameBased:                 true,
		},
		{
			name:                    "Name cache is stale (collection deleted new one created same name). Routing Map Cache returns null",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection2,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection2.ResourceID: routingMapC2Before},
			routingMapAfterRefresh:  nil,
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection2.ResourceID + ":0": addresses1,
			},
			addressesAfterRefresh:     nil,
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC2Before.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 1,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    0,
			nameBased:                 true,
		},
		{
			name:                    "Routing map cache is stale (split happened). Address Cache returns null",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1After},
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":1": addresses1,
			},
			addressesAfterRefresh:     nil,
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1After.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 1,
			expectedRoutingMapRefresh: 1,
			expectedAddressRefresh:    0,
			nameBased:                 true,
		},
		{
			name:                    "Collection cache is stale (deleted created same name). Routing map cache is stale for new collection (split happened). Address Cache returns null",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection2,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{
				collection1.ResourceID: routingMapC1Before,
				collection2.ResourceID: routingMapC2Before,
			},
			routingMapAfterRefresh: map[string]CollectionRoutingMap{collection2.ResourceID: routingMapC2After},
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection2.ResourceID + ":1": addresses1,
			},
			addressesAfterRefresh:     nil,
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC2After.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 1,
			expectedRoutingMapRefresh: 1,
			expectedAddressRefresh:    0,
			nameBased:                 true,
		},
		{
			name:                      "Collection cache is stale (collection deleted). Routing map cache is stale (collection deleted). Address Cache returns null",
			collectionBeforeRefresh:   collection1,
			collectionAfterRefresh:    nil,
			routingMapBeforeRefresh:   map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:    map[string]CollectionRoutingMap{},
			addressesBeforeRefresh:    map[string][]AddressInformation{},
			addressesAfterRefresh:     nil,
			targetAddresses:           nil,
			targetPartitionKeyRange:   nil,
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 1,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    0,
			nameBased:                 true,
			expectError:               &NotFoundException{},
		},
		{
			name:                      "Collection cache is stale (collection deleted). Routing map cache returns null.",
			collectionBeforeRefresh:   collection1,
			collectionAfterRefresh:    nil,
			routingMapBeforeRefresh:   map[string]CollectionRoutingMap{},
			routingMapAfterRefresh:    nil,
			addressesBeforeRefresh:    map[string][]AddressInformation{},
			addressesAfterRefresh:     nil,
			targetAddresses:           nil,
			targetPartitionKeyRange:   nil,
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 1,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    0,
			nameBased:                 true,
			expectError:               &NotFoundException{},
		},
	}

	_ = serviceIdentities

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			collectionCache := &mockCollectionCache{
				currentCollection:      tc.collectionBeforeRefresh,
				collectionAfterRefresh: tc.collectionAfterRefresh,
			}

			routingMapCache := newMockRoutingMapCache()
			for k, v := range tc.routingMapBeforeRefresh {
				routingMapCache.currentRoutingMaps[k] = v
			}
			if tc.routingMapAfterRefresh != nil {
				for k, v := range tc.routingMapAfterRefresh {
					routingMapCache.routingMapsAfterRefresh[k] = v
				}
			} else {
				routingMapCache.routingMapsAfterRefresh = routingMapCache.currentRoutingMaps
			}

			addressCache := newMockAddressCache()
			for k, v := range tc.addressesBeforeRefresh {
				addressCache.currentAddresses[k] = v
			}
			if tc.addressesAfterRefresh != nil {
				for k, v := range tc.addressesAfterRefresh {
					addressCache.addressesAfterRefresh[k] = v
				}
			} else {
				addressCache.addressesAfterRefresh = addressCache.currentAddresses
			}

			resolver := NewAddressResolver()
			resolver.InitializeCaches(collectionCache, routingMapCache, addressCache)

			request := NewAddressResolverRequest(tc.nameBased, "dbs/db/colls/coll/docs/doc1")
			request.ForceNameCacheRefresh = tc.forceNameCacheRefresh
			request.ForcePartitionKeyRangeRefresh = tc.forceRoutingMapRefresh
			request.PartitionKey = `["foo"]`

			resolvedAddresses, err := resolver.ResolveAddressAsync(context.Background(), request, tc.forceAddressRefresh)

			require.Equal(t, tc.expectedCollectionRefresh, collectionCache.refreshCount, "collection cache refresh count mismatch")

			totalRoutingMapRefresh := 0
			for _, count := range routingMapCache.refreshCount {
				totalRoutingMapRefresh += count
				require.LessOrEqual(t, count, 1, "routing map should not be refreshed more than once")
			}
			require.Equal(t, tc.expectedRoutingMapRefresh, totalRoutingMapRefresh, "routing map cache refresh count mismatch")

			totalAddressRefresh := 0
			for _, count := range addressCache.refreshCount {
				totalAddressRefresh += count
				require.LessOrEqual(t, count, 1, "address should not be refreshed more than once")
			}
			require.Equal(t, tc.expectedAddressRefresh, totalAddressRefresh, "address cache refresh count mismatch")

			if tc.expectError != nil {
				require.Error(t, err)
				switch tc.expectError.(type) {
				case *NotFoundException:
					var notFoundErr *NotFoundException
					require.True(t, errors.As(err, &notFoundErr), "expected NotFoundException")
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resolvedAddresses)
				require.Equal(t, tc.targetAddresses[0].GetPhysicalUri().GetURIAsString(), resolvedAddresses[0].GetPhysicalUri().GetURIAsString())
				require.Equal(t, tc.targetPartitionKeyRange.ID, request.RequestContext.ResolvedPartitionKeyRange.ID)
			}
		})
	}
}

func TestCacheRefreshesWhileRoutingByPartitionKeyRangeId(t *testing.T) {
	collection1, collection2 := setupTestCollections()
	routingMapC1Before, routingMapC1After, routingMapC2Before, routingMapC2After, _ := setupRoutingMaps(collection1, collection2)

	addresses1 := createAddresses("tcp://host/partition1")
	addresses2 := createAddresses("tcp://host/partition2")

	testCases := []struct {
		name                      string
		collectionBeforeRefresh   *DocumentCollection
		collectionAfterRefresh    *DocumentCollection
		routingMapBeforeRefresh   map[string]CollectionRoutingMap
		routingMapAfterRefresh    map[string]CollectionRoutingMap
		addressesBeforeRefresh    map[string][]AddressInformation
		addressesAfterRefresh     map[string][]AddressInformation
		rangeIdentity             *PartitionKeyRangeIdentity
		targetAddresses           []AddressInformation
		targetPartitionKeyRange   *PartitionKeyRange
		forceNameCacheRefresh     bool
		forceRoutingMapRefresh    bool
		forceAddressRefresh       bool
		expectedCollectionRefresh int
		expectedRoutingMapRefresh int
		expectedAddressRefresh    int
		nameBased                 bool
		expectError               error
	}{
		{
			name:                    "All caches are up to date. Name Based",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  nil,
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":0": addresses1,
			},
			addressesAfterRefresh:     nil,
			rangeIdentity:             NewPartitionKeyRangeIdentity(collection1.ResourceID, "0"),
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1Before.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    0,
			nameBased:                 true,
		},
		{
			name:                    "All caches are up to date. Name Based. Non existent range with collection rid",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  nil,
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":0": addresses1,
			},
			addressesAfterRefresh:     nil,
			rangeIdentity:             NewPartitionKeyRangeIdentity(collection1.ResourceID, "1"),
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1Before.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 1,
			expectedAddressRefresh:    0,
			nameBased:                 true,
			expectError:               &PartitionKeyRangeGoneException{},
		},
		{
			name:                    "All caches are up to date. Name Based. Non existent range without collection rid",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  nil,
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":0": addresses1,
			},
			addressesAfterRefresh:     nil,
			rangeIdentity:             NewPartitionKeyRangeIdentityWithRangeOnly("1"),
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1Before.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 1,
			expectedRoutingMapRefresh: 1,
			expectedAddressRefresh:    0,
			nameBased:                 true,
			expectError:               &PartitionKeyRangeGoneException{},
		},
		{
			name:                    "All caches are up to date. Name Based.Range doesn't exist in routing map because split happened.",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1After},
			routingMapAfterRefresh:  nil,
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":1": addresses1,
			},
			addressesAfterRefresh:     nil,
			rangeIdentity:             NewPartitionKeyRangeIdentity(collection1.ResourceID, "0"),
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1Before.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    0,
			nameBased:                 true,
			expectError:               &PartitionKeyRangeGoneException{},
		},
		{
			name:                    "Name Based.Routing map cache is outdated because split happened. Address cache returns null.",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1After},
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":1": addresses1,
			},
			addressesAfterRefresh:     nil,
			rangeIdentity:             NewPartitionKeyRangeIdentity(collection1.ResourceID, "0"),
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1After.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 1,
			expectedAddressRefresh:    0,
			nameBased:                 true,
			expectError:               &PartitionKeyRangeGoneException{},
		},
		{
			name:                    "Name Based.Routing map cache is outdated because split happened.",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection1,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1After},
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":1": addresses1,
			},
			addressesAfterRefresh:     nil,
			rangeIdentity:             NewPartitionKeyRangeIdentity(collection1.ResourceID, "1"),
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1After.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 1,
			expectedAddressRefresh:    0,
			nameBased:                 true,
		},
		{
			name:                    "Collection cache is outdated. Routing map cache returns null. Collection is deleted. Range with collection rid.",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  nil,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{},
			routingMapAfterRefresh:  map[string]CollectionRoutingMap{},
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":1": addresses1,
			},
			addressesAfterRefresh:     nil,
			rangeIdentity:             NewPartitionKeyRangeIdentity(collection1.ResourceID, "0"),
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1After.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    0,
			nameBased:                 true,
			expectError:               &InvalidPartitionException{},
		},
		{
			name:                    "Collection cache is outdated. Routing map cache returns null. Collection is deleted. Range without collection rid",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  nil,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{},
			routingMapAfterRefresh:  map[string]CollectionRoutingMap{},
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":1": addresses1,
			},
			addressesAfterRefresh:     nil,
			rangeIdentity:             NewPartitionKeyRangeIdentityWithRangeOnly("0"),
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1After.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 1,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    0,
			nameBased:                 true,
			expectError:               &NotFoundException{},
		},
		{
			name:                    "Collection cache is outdated. Routing map cache returns null. Collection is deleted. Range with collection rid. Rid based.",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  nil,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{},
			routingMapAfterRefresh:  map[string]CollectionRoutingMap{},
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":1": addresses1,
			},
			addressesAfterRefresh:     nil,
			rangeIdentity:             NewPartitionKeyRangeIdentity(collection1.ResourceID, "0"),
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1After.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       false,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 0,
			expectedAddressRefresh:    0,
			nameBased:                 false,
			expectError:               &NotFoundException{},
		},
		{
			name:                    "Collection cache is outdated. Routing map cache is outdated. Address cache is outdated. ForceAddressRefresh. Range with collection rid. Name based.",
			collectionBeforeRefresh: collection1,
			collectionAfterRefresh:  collection2,
			routingMapBeforeRefresh: map[string]CollectionRoutingMap{collection1.ResourceID: routingMapC1Before},
			routingMapAfterRefresh:  map[string]CollectionRoutingMap{collection2.ResourceID: routingMapC2Before},
			addressesBeforeRefresh: map[string][]AddressInformation{
				collection1.ResourceID + ":1": addresses1,
			},
			addressesAfterRefresh: map[string][]AddressInformation{
				collection2.ResourceID + ":1": addresses2,
			},
			rangeIdentity:             NewPartitionKeyRangeIdentity(collection1.ResourceID, "0"),
			targetAddresses:           addresses1,
			targetPartitionKeyRange:   routingMapC1After.GetOrderedPartitionKeyRanges()[0],
			forceNameCacheRefresh:     false,
			forceRoutingMapRefresh:    false,
			forceAddressRefresh:       true,
			expectedCollectionRefresh: 0,
			expectedRoutingMapRefresh: 1,
			expectedAddressRefresh:    1,
			nameBased:                 true,
			expectError:               &InvalidPartitionException{},
		},
	}

	_ = routingMapC2Before
	_ = routingMapC2After

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			collectionCache := &mockCollectionCache{
				currentCollection:      tc.collectionBeforeRefresh,
				collectionAfterRefresh: tc.collectionAfterRefresh,
			}

			routingMapCache := newMockRoutingMapCache()
			for k, v := range tc.routingMapBeforeRefresh {
				routingMapCache.currentRoutingMaps[k] = v
			}
			if tc.routingMapAfterRefresh != nil {
				for k, v := range tc.routingMapAfterRefresh {
					routingMapCache.routingMapsAfterRefresh[k] = v
				}
			} else {
				routingMapCache.routingMapsAfterRefresh = routingMapCache.currentRoutingMaps
			}

			addressCache := newMockAddressCache()
			for k, v := range tc.addressesBeforeRefresh {
				addressCache.currentAddresses[k] = v
			}
			if tc.addressesAfterRefresh != nil {
				for k, v := range tc.addressesAfterRefresh {
					addressCache.addressesAfterRefresh[k] = v
				}
			} else {
				addressCache.addressesAfterRefresh = addressCache.currentAddresses
			}

			resolver := NewAddressResolver()
			resolver.InitializeCaches(collectionCache, routingMapCache, addressCache)

			request := NewAddressResolverRequest(tc.nameBased, "dbs/db/colls/coll/docs/doc1")
			request.ForceNameCacheRefresh = tc.forceNameCacheRefresh
			request.ForcePartitionKeyRangeRefresh = tc.forceRoutingMapRefresh
			request.RouteTo(tc.rangeIdentity)

			resolvedAddresses, err := resolver.ResolveAddressAsync(context.Background(), request, tc.forceAddressRefresh)

			require.Equal(t, tc.expectedCollectionRefresh, collectionCache.refreshCount, "collection cache refresh count mismatch")

			totalRoutingMapRefresh := 0
			for _, count := range routingMapCache.refreshCount {
				totalRoutingMapRefresh += count
				require.LessOrEqual(t, count, 1, "routing map should not be refreshed more than once")
			}
			require.Equal(t, tc.expectedRoutingMapRefresh, totalRoutingMapRefresh, "routing map cache refresh count mismatch")

			totalAddressRefresh := 0
			for _, count := range addressCache.refreshCount {
				totalAddressRefresh += count
				require.LessOrEqual(t, count, 1, "address should not be refreshed more than once")
			}
			require.Equal(t, tc.expectedAddressRefresh, totalAddressRefresh, "address cache refresh count mismatch")

			if tc.expectError != nil {
				require.Error(t, err)
				switch tc.expectError.(type) {
				case *NotFoundException:
					var notFoundErr *NotFoundException
					require.True(t, errors.As(err, &notFoundErr), "expected NotFoundException, got: %v", err)
				case *InvalidPartitionException:
					var invalidPartitionErr *InvalidPartitionException
					require.True(t, errors.As(err, &invalidPartitionErr), "expected InvalidPartitionException, got: %v", err)
				case *PartitionKeyRangeGoneException:
					var partitionKeyRangeGoneErr *PartitionKeyRangeGoneException
					require.True(t, errors.As(err, &partitionKeyRangeGoneErr), "expected PartitionKeyRangeGoneException, got: %v", err)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resolvedAddresses)
				require.Equal(t, tc.targetAddresses[0].GetPhysicalUri().GetURIAsString(), resolvedAddresses[0].GetPhysicalUri().GetURIAsString())
				require.Equal(t, tc.targetPartitionKeyRange.ID, request.RequestContext.ResolvedPartitionKeyRange.ID)
			}
		})
	}
}
