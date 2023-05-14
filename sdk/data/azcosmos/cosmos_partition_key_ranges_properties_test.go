// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestPartitionKeyRangesPropertiesSerialization(t *testing.T) {
	nowAsUnix := time.Unix(time.Now().Unix(), 0)

	etag := azcore.ETag("etag")

	properties := &PartitionKeyRangesProperties{
		ResourceID: "someResourceId",
		Count:      1,
		PartitionKeyRanges: []PartitionKeyRange{
			{
				ID:           "someId",
				ETag:         &etag,
				SelfLink:     "someSelfLink",
				LastModified: nowAsUnix,
				MinInclusive: "someMinInclusive",
				MaxExclusive: "someMaxExclusive",
			},
		},
	}

	jsonString, err := json.Marshal(properties)
	if err != nil {
		t.Fatal(err)
	}

	otherProperties := &PartitionKeyRangesProperties{}
	err = json.Unmarshal(jsonString, otherProperties)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	if properties.ResourceID != otherProperties.ResourceID {
		t.Errorf("Expected ResourceID to be %s, but got %s", properties.ResourceID, otherProperties.ResourceID)
	}

	if properties.Count != otherProperties.Count {
		t.Errorf("Expected Count to be %d, but got %d", properties.Count, otherProperties.Count)
	}

	if len(properties.PartitionKeyRanges) != len(otherProperties.PartitionKeyRanges) {
		t.Errorf("Expected PartitionKeyRanges length to be %d, but got %d", len(properties.PartitionKeyRanges), len(otherProperties.PartitionKeyRanges))
	}

	if properties.PartitionKeyRanges[0].ID != otherProperties.PartitionKeyRanges[0].ID {
		t.Errorf("Expected ID to be %s, but got %s", properties.PartitionKeyRanges[0].ID, otherProperties.PartitionKeyRanges[0].ID)
	}

	if *properties.PartitionKeyRanges[0].ETag != *otherProperties.PartitionKeyRanges[0].ETag {
		t.Errorf("Expected ETag to be %v, but got %v", properties.PartitionKeyRanges[0].ETag, otherProperties.PartitionKeyRanges[0].ETag)
	}

	if properties.PartitionKeyRanges[0].SelfLink != otherProperties.PartitionKeyRanges[0].SelfLink {
		t.Errorf("Expected SelfLink to be %s, but got %s", properties.PartitionKeyRanges[0].SelfLink, otherProperties.PartitionKeyRanges[0].SelfLink)
	}

	if properties.PartitionKeyRanges[0].LastModified != otherProperties.PartitionKeyRanges[0].LastModified {
		t.Errorf("Expected LastModified to be %s, but got %s", properties.PartitionKeyRanges[0].LastModified, otherProperties.PartitionKeyRanges[0].LastModified)
	}

	if properties.PartitionKeyRanges[0].MinInclusive != otherProperties.PartitionKeyRanges[0].MinInclusive {
		t.Errorf("Expected MinInclusive to be %s, but got %s", properties.PartitionKeyRanges[0].MinInclusive, otherProperties.PartitionKeyRanges[0].MinInclusive)
	}

	if properties.PartitionKeyRanges[0].MaxExclusive != otherProperties.PartitionKeyRanges[0].MaxExclusive {
		t.Errorf("Expected MaxExclusive to be %s, but got %s", properties.PartitionKeyRanges[0].MaxExclusive, otherProperties.PartitionKeyRanges[0].MaxExclusive)
	}
}
