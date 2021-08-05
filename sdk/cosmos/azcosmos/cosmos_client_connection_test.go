// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestPathGeneration(t *testing.T) {
	connection := &cosmosClientConnection{}

	expected := "dbs/testdb/colls/testcoll"
	actual := connection.getPath("dbs/testdb", pathSegmentCollection, "testcoll")
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	expected = "dbs/testdb"
	actual = connection.getPath("", pathSegmentDatabase, "testdb")
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
