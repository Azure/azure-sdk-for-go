// +build !emulator
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestPathGeneration(t *testing.T) {

	expected := "dbs/testdb/colls/testcoll"
	actual := getPath("dbs/testdb", pathSegmentCollection, "testcoll")
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	expected = "dbs/testdb"
	actual = getPath("", pathSegmentDatabase, "testdb")
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	expected = "dbs/esc%40ped"
	actual = getPath("", pathSegmentDatabase, "esc@ped")
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestPathToResourceTypeMapping(t *testing.T) {

	expected := pathSegmentDatabase
	actual, err := getResourcePath(resourceTypeDatabase)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	expected = pathSegmentCollection
	actual, err = getResourcePath(resourceTypeCollection)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
