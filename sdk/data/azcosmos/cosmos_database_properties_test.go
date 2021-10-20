// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDatabasePropertiesSerialization(t *testing.T) {
	nowAsUnix := time.Now().Unix()

	now := UnixTime{
		Time: time.Unix(nowAsUnix, 0),
	}

	properties := &DatabaseProperties{
		Id:           "someId",
		ETag:         "someEtag",
		SelfLink:     "someSelfLink",
		ResourceId:   "someResourceId",
		LastModified: &now,
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

}
