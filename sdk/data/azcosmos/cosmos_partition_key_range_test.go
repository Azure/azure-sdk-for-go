// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestPartitionKeyRangeSerialization(t *testing.T) {
	etag := azcore.ETag("\"00000000-0000-0000-0000-000000000000\"")
	pkr := partitionKeyRange{
		ID:                      "0",
		ResourceID:              "rid1",
		ETag:                    &etag,
		MinInclusive:            "00000000",
		MaxExclusive:            "FFFFFFFF",
		ResourceIDPrefix:        123,
		SelfLink:                "self",
		ThroughputFraction:      0.5,
		Status:                  "online",
		Parents:                 []string{"parent1", "parent2"},
		OwnedArchivalPKRangeIds: []string{"pkr1", "pkr2"},
		LastModified:            time.Unix(1610000000, 0),
		LSN:                     789,
	}

	jsonBytes, err := json.Marshal(pkr)
	if err != nil {
		t.Fatalf("Failed to marshal PartitionKeyRange: %v", err)
	}

	var newPkr partitionKeyRange
	err = json.Unmarshal(jsonBytes, &newPkr)
	if err != nil {
		t.Fatalf("Failed to unmarshal PartitionKeyRange: %v", err)
	}

	if pkr.ID != newPkr.ID {
		t.Errorf("ID mismatch: expected %s, got %s", pkr.ID, newPkr.ID)
	}
	if pkr.ResourceID != newPkr.ResourceID {
		t.Errorf("ResourceID mismatch: expected %s, got %s", pkr.ResourceID, newPkr.ResourceID)
	}
	if pkr.MinInclusive != newPkr.MinInclusive {
		t.Errorf("MinInclusive mismatch: expected %s, got %s", pkr.MinInclusive, newPkr.MinInclusive)
	}
	if pkr.MaxExclusive != newPkr.MaxExclusive {
		t.Errorf("MaxExclusive mismatch: expected %s, got %s", pkr.MaxExclusive, newPkr.MaxExclusive)
	}
	if pkr.ResourceIDPrefix != newPkr.ResourceIDPrefix {
		t.Errorf("ResourceIDPrefix mismatch: expected %d, got %d", pkr.ResourceIDPrefix, newPkr.ResourceIDPrefix)
	}
	if pkr.SelfLink != newPkr.SelfLink {
		t.Errorf("SelfLink mismatch: expected %s, got %s", pkr.SelfLink, newPkr.SelfLink)
	}
	if pkr.ThroughputFraction != newPkr.ThroughputFraction {
		t.Errorf("ThroughputFraction mismatch: expected %f, got %f", pkr.ThroughputFraction, newPkr.ThroughputFraction)
	}
	if pkr.Status != newPkr.Status {
		t.Errorf("Status mismatch: expected %s, got %s", pkr.Status, newPkr.Status)
	}
	if pkr.LastModified.Unix() != newPkr.LastModified.Unix() {
		t.Errorf("LastModified mismatch: expected %v, got %v", pkr.LastModified.Unix(), newPkr.LastModified.Unix())
	}
	if pkr.LSN != newPkr.LSN {
		t.Errorf("LSN mismatch: expected %d, got %d", pkr.LSN, newPkr.LSN)
	}
}
