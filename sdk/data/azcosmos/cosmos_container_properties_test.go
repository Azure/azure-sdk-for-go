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
			VectorIndexes: []VectorIndex{
				{
					Path: "/vector1",
					Type: VectorIndexTypeFlat,
				},
				{
					Path: "/embeddings/textVector",
					Type: VectorIndexTypeDiskANN,
				},
			},
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
		VectorEmbeddingPolicy: &VectorEmbeddingPolicy{
			VectorEmbeddings: []VectorEmbedding{
				{
					Path:             "/vector1",
					DataType:         VectorDataTypeFloat32,
					DistanceFunction: VectorDistanceFunctionCosine,
					Dimensions:       1536,
				},
				{
					Path:             "/embeddings/textVector",
					DataType:         VectorDataTypeUint8,
					DistanceFunction: VectorDistanceFunctionEuclidean,
					Dimensions:       768,
				},
			},
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

	if properties.ID != otherProperties.ID {
		t.Errorf("Expected Id to be %s, but got %s", properties.ID, otherProperties.ID)
	}

	if *properties.ETag != *otherProperties.ETag {
		t.Errorf("Expected ETag to be %s, but got %s", *properties.ETag, *otherProperties.ETag)
	}

	if properties.SelfLink != otherProperties.SelfLink {
		t.Errorf("Expected SelfLink to be %s, but got %s", properties.SelfLink, otherProperties.SelfLink)
	}

	if properties.ResourceID != otherProperties.ResourceID {
		t.Errorf("Expected ResourceId to be %s, but got %s", properties.ResourceID, otherProperties.ResourceID)
	}

	if properties.LastModified != otherProperties.LastModified {
		t.Errorf("Expected LastModified.Time to be %v, but got %v", properties.LastModified, otherProperties.LastModified)
	}

	if otherProperties.AnalyticalStoreTimeToLiveInSeconds != nil {
		t.Errorf("Expected AnalyticalStoreTimeToLiveInSeconds to be nil, but got %d", *otherProperties.AnalyticalStoreTimeToLiveInSeconds)
	}

	if otherProperties.DefaultTimeToLive != nil {
		t.Errorf("Expected DefaultTimeToLive to be nil, but got %d", *otherProperties.DefaultTimeToLive)
	}

	if properties.PartitionKeyDefinition.Paths[0] != otherProperties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Expected PartitionKeyDefinition.Paths[0] to be %s, but got %s", properties.PartitionKeyDefinition.Paths[0], otherProperties.PartitionKeyDefinition.Paths[0])
	}

	if properties.PartitionKeyDefinition.Version != otherProperties.PartitionKeyDefinition.Version {
		t.Errorf("Expected PartitionKeyDefinition.Version to be %d, but got %d", properties.PartitionKeyDefinition.Version, otherProperties.PartitionKeyDefinition.Version)
	}

	if otherProperties.IndexingPolicy == nil {
		t.Errorf("Expected IndexingPolicy to be not nil, but got nil")
	}

	if otherProperties.IndexingPolicy.Automatic != properties.IndexingPolicy.Automatic {
		t.Errorf("Expected IndexingPolicy.Automatic to be %t, but got %t", properties.IndexingPolicy.Automatic, otherProperties.IndexingPolicy.Automatic)
	}

	if otherProperties.IndexingPolicy.IndexingMode != properties.IndexingPolicy.IndexingMode {
		t.Errorf("Expected IndexingPolicy.IndexingMode to be %v, but got %v", properties.IndexingPolicy.IndexingMode, otherProperties.IndexingPolicy.IndexingMode)
	}

	if otherProperties.IndexingPolicy.IncludedPaths[0].Path != properties.IndexingPolicy.IncludedPaths[0].Path {
		t.Errorf("Expected IndexingPolicy.IncludedPaths[0].Path to be %s, but got %s", properties.IndexingPolicy.IncludedPaths[0].Path, otherProperties.IndexingPolicy.IncludedPaths[0].Path)
	}

	if otherProperties.IndexingPolicy.ExcludedPaths[0].Path != properties.IndexingPolicy.ExcludedPaths[0].Path {
		t.Errorf("Expected IndexingPolicy.ExcludedPaths[0].Path to be %s, but got %s", properties.IndexingPolicy.ExcludedPaths[0].Path, otherProperties.IndexingPolicy.ExcludedPaths[0].Path)
	}

	if otherProperties.IndexingPolicy.SpatialIndexes[0].Path != properties.IndexingPolicy.SpatialIndexes[0].Path {
		t.Errorf("Expected IndexingPolicy.SpatialIndexes[0].Path to be %s, but got %s", properties.IndexingPolicy.SpatialIndexes[0].Path, otherProperties.IndexingPolicy.SpatialIndexes[0].Path)
	}

	if otherProperties.IndexingPolicy.SpatialIndexes[0].SpatialTypes[0] != properties.IndexingPolicy.SpatialIndexes[0].SpatialTypes[0] {
		t.Errorf("Expected IndexingPolicy.SpatialIndexes[0].SpatialTypes[0] to be %v, but got %v", properties.IndexingPolicy.SpatialIndexes[0].SpatialTypes[0], otherProperties.IndexingPolicy.SpatialIndexes[0].SpatialTypes[0])
	}

	if otherProperties.IndexingPolicy.CompositeIndexes[0][0].Path != properties.IndexingPolicy.CompositeIndexes[0][0].Path {
		t.Errorf("Expected IndexingPolicy.CompositeIndexes[0][0].Path to be %s, but got %s", properties.IndexingPolicy.CompositeIndexes[0][0].Path, otherProperties.IndexingPolicy.CompositeIndexes[0][0].Path)
	}

	if len(otherProperties.IndexingPolicy.VectorIndexes) != len(properties.IndexingPolicy.VectorIndexes) {
		t.Errorf("Expected VectorIndexes length to be %d, but got %d", len(properties.IndexingPolicy.VectorIndexes), len(otherProperties.IndexingPolicy.VectorIndexes))
	}

	if otherProperties.IndexingPolicy.VectorIndexes[0].Path != properties.IndexingPolicy.VectorIndexes[0].Path {
		t.Errorf("Expected VectorIndexes[0].Path to be %s, but got %s", properties.IndexingPolicy.VectorIndexes[0].Path, otherProperties.IndexingPolicy.VectorIndexes[0].Path)
	}

	if otherProperties.IndexingPolicy.VectorIndexes[0].Type != properties.IndexingPolicy.VectorIndexes[0].Type {
		t.Errorf("Expected VectorIndexes[0].Type to be %s, but got %s", properties.IndexingPolicy.VectorIndexes[0].Type, otherProperties.IndexingPolicy.VectorIndexes[0].Type)
	}

	if otherProperties.IndexingPolicy.VectorIndexes[1].Path != properties.IndexingPolicy.VectorIndexes[1].Path {
		t.Errorf("Expected VectorIndexes[1].Path to be %s, but got %s", properties.IndexingPolicy.VectorIndexes[1].Path, otherProperties.IndexingPolicy.VectorIndexes[1].Path)
	}

	if otherProperties.IndexingPolicy.VectorIndexes[1].Type != properties.IndexingPolicy.VectorIndexes[1].Type {
		t.Errorf("Expected VectorIndexes[1].Type to be %s, but got %s", properties.IndexingPolicy.VectorIndexes[1].Type, otherProperties.IndexingPolicy.VectorIndexes[1].Type)
	}

	if otherProperties.UniqueKeyPolicy == nil {
		t.Errorf("Expected UniqueKeyPolicy to be not nil, but got nil")
	}

	if otherProperties.UniqueKeyPolicy.UniqueKeys[0].Paths[0] != properties.UniqueKeyPolicy.UniqueKeys[0].Paths[0] {
		t.Errorf("Expected UniqueKeyPolicy.UniqueKeys[0].Paths[0] to be %s, but got %s", properties.UniqueKeyPolicy.UniqueKeys[0].Paths[0], otherProperties.UniqueKeyPolicy.UniqueKeys[0].Paths[0])
	}

	if otherProperties.ConflictResolutionPolicy == nil {
		t.Errorf("Expected ConflictResolutionPolicy to be not nil, but got nil")
	}

	if otherProperties.ConflictResolutionPolicy.Mode != properties.ConflictResolutionPolicy.Mode {
		t.Errorf("Expected ConflictResolutionPolicy.Mode to be %v, but got %v", properties.ConflictResolutionPolicy.Mode, otherProperties.ConflictResolutionPolicy.Mode)
	}

	if otherProperties.ConflictResolutionPolicy.ResolutionPath != properties.ConflictResolutionPolicy.ResolutionPath {
		t.Errorf("Expected ConflictResolutionPolicy.ResolutionPath to be %s, but got %s", properties.ConflictResolutionPolicy.ResolutionPath, otherProperties.ConflictResolutionPolicy.ResolutionPath)
	}

	if otherProperties.VectorEmbeddingPolicy == nil {
		t.Errorf("Expected VectorEmbeddingPolicy to be not nil, but got nil")
	}

	if len(otherProperties.VectorEmbeddingPolicy.VectorEmbeddings) != len(properties.VectorEmbeddingPolicy.VectorEmbeddings) {
		t.Errorf("Expected VectorEmbeddings length to be %d, but got %d", len(properties.VectorEmbeddingPolicy.VectorEmbeddings), len(otherProperties.VectorEmbeddingPolicy.VectorEmbeddings))
	}

	// Test first vector embedding
	if otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[0].Path != properties.VectorEmbeddingPolicy.VectorEmbeddings[0].Path {
		t.Errorf("Expected VectorEmbeddings[0].Path to be %s, but got %s", properties.VectorEmbeddingPolicy.VectorEmbeddings[0].Path, otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[0].Path)
	}

	if otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[0].DataType != properties.VectorEmbeddingPolicy.VectorEmbeddings[0].DataType {
		t.Errorf("Expected VectorEmbeddings[0].DataType to be %s, but got %s", properties.VectorEmbeddingPolicy.VectorEmbeddings[0].DataType, otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[0].DataType)
	}

	if otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[0].DistanceFunction != properties.VectorEmbeddingPolicy.VectorEmbeddings[0].DistanceFunction {
		t.Errorf("Expected VectorEmbeddings[0].DistanceFunction to be %s, but got %s", properties.VectorEmbeddingPolicy.VectorEmbeddings[0].DistanceFunction, otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[0].DistanceFunction)
	}

	if otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[0].Dimensions != properties.VectorEmbeddingPolicy.VectorEmbeddings[0].Dimensions {
		t.Errorf("Expected VectorEmbeddings[0].Dimensions to be %d, but got %d", properties.VectorEmbeddingPolicy.VectorEmbeddings[0].Dimensions, otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[0].Dimensions)
	}

	// Test second vector embedding
	if otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[1].Path != properties.VectorEmbeddingPolicy.VectorEmbeddings[1].Path {
		t.Errorf("Expected VectorEmbeddings[1].Path to be %s, but got %s", properties.VectorEmbeddingPolicy.VectorEmbeddings[1].Path, otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[1].Path)
	}

	if otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[1].DataType != properties.VectorEmbeddingPolicy.VectorEmbeddings[1].DataType {
		t.Errorf("Expected VectorEmbeddings[1].DataType to be %s, but got %s", properties.VectorEmbeddingPolicy.VectorEmbeddings[1].DataType, otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[1].DataType)
	}

	if otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[1].DistanceFunction != properties.VectorEmbeddingPolicy.VectorEmbeddings[1].DistanceFunction {
		t.Errorf("Expected VectorEmbeddings[1].DistanceFunction to be %s, but got %s", properties.VectorEmbeddingPolicy.VectorEmbeddings[1].DistanceFunction, otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[1].DistanceFunction)
	}

	if otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[1].Dimensions != properties.VectorEmbeddingPolicy.VectorEmbeddings[1].Dimensions {
		t.Errorf("Expected VectorEmbeddings[1].Dimensions to be %d, but got %d", properties.VectorEmbeddingPolicy.VectorEmbeddings[1].Dimensions, otherProperties.VectorEmbeddingPolicy.VectorEmbeddings[1].Dimensions)
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
