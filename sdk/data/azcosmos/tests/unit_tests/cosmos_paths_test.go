// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestPathCreateLink(t *testing.T) {

	expected := "dbs/testdb/colls/testcoll"
	actual := createLink("dbs/testdb", pathSegmentCollection, "testcoll")
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	expected = "dbs/testdb"
	actual = createLink("", pathSegmentDatabase, "testdb")
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	expected = "dbs/with%20space"
	actual = createLink("", pathSegmentDatabase, "with space")
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestPathToResourceTypeMapping(t *testing.T) {
	verifyPathResultAndExpectation(t, resourceTypeDatabase, pathSegmentDatabase)
	verifyPathResultAndExpectation(t, resourceTypeCollection, pathSegmentCollection)
	verifyPathResultAndExpectation(t, resourceTypeDocument, pathSegmentDocument)
	verifyPathResultAndExpectation(t, resourceTypeDatabaseAccount, pathSegmentDatabaseAccount)
	verifyPathResultAndExpectation(t, resourceTypeOffer, pathSegmentOffer)
	verifyPathResultAndExpectation(t, resourceTypeUser, pathSegmentUser)
	verifyPathResultAndExpectation(t, resourceTypeStoredProcedure, pathSegmentStoredProcedure)
	verifyPathResultAndExpectation(t, resourceTypeUserDefinedFunction, pathSegmentUserDefinedFunction)
	verifyPathResultAndExpectation(t, resourceTypeTrigger, pathSegmentTrigger)
	verifyPathResultAndExpectation(t, resourceTypePermission, pathSegmentPermission)
	verifyPathResultAndExpectation(t, resourceTypePartitionKeyRange, pathSegmentPartitionKeyRange)
	verifyPathResultAndExpectation(t, resourceTypeClientEncryptionKey, pathSegmentClientEncryptionKey)
	verifyPathResultAndExpectation(t, resourceTypeUser, pathSegmentUser)
	verifyPathResultAndExpectation(t, resourceTypeConflict, pathSegmentConflict)
}

func verifyPathResultAndExpectation(t *testing.T, resourceType resourceType, expected string) {
	actual, err := getResourcePath(resourceType)
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
