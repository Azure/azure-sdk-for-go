// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestContainerPropertiesSerialization(t *testing.T) {
	nowAsUnix := time.Unix(time.Now().Unix(), 0)

	etag := azcore.ETag("etag")

	properties := ContainerProperties{
		ID:           "someId",
		ETag:         &etag,
		SelfLink:     "someSelfLink",
		ResourceID:   "someResourceId",
		LastModified: nowAsUnix,
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"somePath"},
			Version: 2,
		},
		IndexingPolicy: &IndexingPolicy{
			IncludedPaths: []IncludedPath{
				{Path: "/someIncludedPath"},
			},
			ExcludedPaths: []ExcludedPath{
				{Path: "/someExcludedPath"},
			},
			Automatic:    true,
			IndexingMode: IndexingModeNone,
			SpatialIndexes: []SpatialIndex{
				{Path: "/someSpatialIndex",
					SpatialTypes: []SpatialType{SpatialTypePoint}}},
			CompositeIndexes: [][]CompositeIndex{
				{
					{Path: "/someCompositeIndex",
						Order: CompositeIndexAscending},
				}},
		},
		UniqueKeyPolicy: &UniqueKeyPolicy{
			UniqueKeys: []UniqueKey{
				{Paths: []string{"/someUniqueKey"}},
			},
		},
		ConflictResolutionPolicy: &ConflictResolutionPolicy{
			Mode:           ConflictResolutionModeLastWriteWins,
			ResolutionPath: "/someResolutionPath",
		},
	}

	jsonString, err := json.Marshal(properties)
	if err != nil {
		t.Fatal(err)
	}

	otherProperties := &ContainerProperties{}
	err = json.Unmarshal(jsonString, otherProperties)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"ID", otherProperties.ID, properties.ID},
		{"ETag", *otherProperties.ETag, *properties.ETag},
		{"SelfLink", otherProperties.SelfLink, properties.SelfLink},
		{"ResourceID", otherProperties.ResourceID, properties.ResourceID},
		{"LastModified", otherProperties.LastModified, properties.LastModified},
		{"AnalyticalStoreTimtoLiveInSeconds", otherProperties.AnalyticalStoreTimeToLiveInSeconds, properties.AnalyticalStoreTimeToLiveInSeconds},
		{"DefaultTimeToLive", otherProperties.DefaultTimeToLive, properties.DefaultTimeToLive},
		{"PartitionKeyDefinitionPaths", otherProperties.PartitionKeyDefinition.Paths[0], properties.PartitionKeyDefinition.Paths[0]},
		{"PartitionKeyDefinitionVersion", otherProperties.PartitionKeyDefinition.Version, properties.PartitionKeyDefinition.Version},
		{"IndexingPolicy", otherProperties.IndexingPolicy != nil, properties.IndexingPolicy != nil},
		{"IndexingPolicyAutomatic", otherProperties.IndexingPolicy.Automatic, properties.IndexingPolicy.Automatic},
		{"IndexingPolicyIndexingMode", otherProperties.IndexingPolicy.IndexingMode, properties.IndexingPolicy.IndexingMode},
		{"IndexingPolicyIncludedPaths", otherProperties.IndexingPolicy.IncludedPaths[0], properties.IndexingPolicy.IncludedPaths[0]},
		{"IndexingPolicyExcludedPaths", otherProperties.IndexingPolicy.ExcludedPaths[0], properties.IndexingPolicy.ExcludedPaths[0]},
		{"IndexingPolicySpatialIndexesPath", otherProperties.IndexingPolicy.SpatialIndexes[0].Path, properties.IndexingPolicy.SpatialIndexes[0].Path},
		{"IndexingPolicySpatialIndexesSpatialTypes", otherProperties.IndexingPolicy.SpatialIndexes[0].SpatialTypes[0], properties.IndexingPolicy.SpatialIndexes[0].SpatialTypes[0]},
		{"IndexingPolicyCompositeIndexesPath", otherProperties.IndexingPolicy.CompositeIndexes[0][0].Path, properties.IndexingPolicy.CompositeIndexes[0][0].Path},
		{"UniqueKeyPolicy", otherProperties.UniqueKeyPolicy != nil, properties.UniqueKeyPolicy != nil},
		{"UniqueKeyPolicyUniqueKeys", otherProperties.UniqueKeyPolicy.UniqueKeys[0].Paths[0], properties.UniqueKeyPolicy.UniqueKeys[0].Paths[0]},
		{"ConflictResolutionPolicy", otherProperties.ConflictResolutionPolicy != nil, properties.ConflictResolutionPolicy !=nil},
		{"ConflictResolutionPolicyMode", otherProperties.ConflictResolutionPolicy.Mode, properties.ConflictResolutionPolicy.Mode},
		{"ConflictResolutionPolicyResolutionPath", otherProperties.ConflictResolutionPolicy.ResolutionPath, properties.ConflictResolutionPolicy.ResolutionPath},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Expected %s to be %v, but got %v", tt.name, tt.expected, tt.got)
			}
		})
	}
}

func TestContainerPropertiesSerializationWithTTL(t *testing.T) {
	jsonString := []byte(`{"defaultTtl": 10, "analyticalStorageTtl": 20}`)

	properties := &ContainerProperties{}
	err := json.Unmarshal(jsonString, properties)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	if *properties.DefaultTimeToLive != 10 {
		t.Errorf("Expected properties.DefaultTimeToLive to be %d, but got %d", 10, properties.DefaultTimeToLive)
	}

	if *properties.AnalyticalStoreTimeToLiveInSeconds != 20 {
		t.Errorf("Expected properties.AnalyticalStoreTimeToLiveInSeconds to be %d, but got %d", 20, properties.AnalyticalStoreTimeToLiveInSeconds)
	}
}