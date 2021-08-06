// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"
	"time"
)

func TestContainerPropertiesSerialization(t *testing.T) {
	nowAsUnix := time.Now().Unix()

	now := UnixTime{
		Time: time.Unix(nowAsUnix, 0),
	}

	properties := &CosmosContainerProperties{
		Id:           "someId",
		ETag:         "someEtag",
		SelfLink:     "someSelfLink",
		ResourceId:   "someResourceId",
		LastModified: &now,
		PartitionKeyDefinition: &PartitionKeyDefinition{
			Paths:   []string{"somePath"},
			Version: PartitionKeyDefinitionVersion2,
		},
	}

	jsonString, err := json.Marshal(properties)
	if err != nil {
		t.Fatal(err)
	}

	otherProperties := &CosmosContainerProperties{}
	err = json.Unmarshal(jsonString, otherProperties)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	if properties.Id != otherProperties.Id {
		t.Errorf("Expected otherProperties.Id to be %s, but got %s", properties.Id, otherProperties.Id)
	}

	if properties.ETag != otherProperties.ETag {
		t.Errorf("Expected otherProperties.ETag to be %s, but got %s", properties.ETag, otherProperties.ETag)
	}

	if properties.SelfLink != otherProperties.SelfLink {
		t.Errorf("Expected otherProperties.SelfLink to be %s, but got %s", properties.SelfLink, otherProperties.SelfLink)
	}

	if properties.ResourceId != otherProperties.ResourceId {
		t.Errorf("Expected otherProperties.ResourceId to be %s, but got %s", properties.ResourceId, otherProperties.ResourceId)
	}

	if properties.LastModified.Time != otherProperties.LastModified.Time {
		t.Errorf("Expected otherProperties.LastModified.Time to be %s, but got %s", properties.LastModified.Time.UTC(), otherProperties.LastModified.Time.UTC())
	}

	if otherProperties.AnalyticalStoreTimeToLiveInSeconds != nil {
		t.Errorf("Expected otherProperties.AnalyticalStoreTimeToLiveInSeconds to be nil, but got %d", *otherProperties.AnalyticalStoreTimeToLiveInSeconds)
	}

	if otherProperties.DefaultTimeToLive != nil {
		t.Errorf("Expected otherProperties.DefaultTimeToLive to be nil, but got %d", *otherProperties.DefaultTimeToLive)
	}

	if otherProperties.PartitionKeyDefinition == nil {
		t.Errorf("Expected otherProperties.PartitionKeyDefinition to be not nil, but got nil")
	}

	if properties.PartitionKeyDefinition.Paths[0] != otherProperties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Expected otherProperties.PartitionKeyDefinition.Paths[0] to be %s, but got %s", properties.PartitionKeyDefinition.Paths[0], otherProperties.PartitionKeyDefinition.Paths[0])
	}

	if properties.PartitionKeyDefinition.Version != otherProperties.PartitionKeyDefinition.Version {
		t.Errorf("Expected otherProperties.PartitionKeyDefinition.Version to be %d, but got %d", properties.PartitionKeyDefinition.Version, otherProperties.PartitionKeyDefinition.Version)
	}
}

func TestContainerPropertiesSerializationWithTTL(t *testing.T) {
	jsonString := []byte(`{"defaultTtl": 10, "analyticalStorageTtl": 20}`)

	properties := &CosmosContainerProperties{}
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
