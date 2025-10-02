// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestDatabasePropertiesSerialization(t *testing.T) {
	nowAsUnix := time.Unix(time.Now().Unix(), 0)

	etag := azcore.ETag("someETag")
	properties := DatabaseProperties{
		ID:           "someId",
		ETag:         &etag,
		SelfLink:     "someSelfLink",
		ResourceID:   "someResourceId",
		LastModified: nowAsUnix,
	}

	jsonString, err := json.Marshal(properties)
	if err != nil {
		t.Fatal(err)
	}

	otherProperties := &DatabaseProperties{}
	err = json.Unmarshal(jsonString, otherProperties)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	if properties.ID != otherProperties.ID {
		t.Errorf("Expected otherProperties.Id to be %s, but got %s", properties.ID, otherProperties.ID)
	}

	if *properties.ETag != *otherProperties.ETag {
		t.Errorf("Expected otherProperties.ETag to be %s, but got %s", *properties.ETag, *otherProperties.ETag)
	}

	if properties.SelfLink != otherProperties.SelfLink {
		t.Errorf("Expected otherProperties.SelfLink to be %s, but got %s", properties.SelfLink, otherProperties.SelfLink)
	}

	if properties.ResourceID != otherProperties.ResourceID {
		t.Errorf("Expected otherProperties.ResourceId to be %s, but got %s", properties.ResourceID, otherProperties.ResourceID)
	}

	if properties.LastModified != otherProperties.LastModified {
		t.Errorf("Expected otherProperties.LastModified.Time to be %v, but got %v", properties.LastModified, otherProperties.LastModified)
	}

}
