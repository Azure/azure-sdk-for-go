// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSessionContainer_Basic(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")

	numCollections := 2
	numPartitionKeyRangeIDs := 5

	for i := 0; i < numCollections; i++ {
		collectionResourceID := fmt.Sprintf("collection_rid_%d", i)
		collectionFullName := fmt.Sprintf("dbs/db1/colls/collName_%d", i)

		for j := 0; j < numPartitionKeyRangeIDs; j++ {
			partitionKeyRangeID := fmt.Sprintf("range_%d", j)
			lsn := fmt.Sprintf("1#%d#4=90#5=2", j)

			sessionContainer.SetSessionTokenFromRID(
				collectionResourceID,
				collectionFullName,
				map[string]string{HTTPHeaderSessionToken: partitionKeyRangeID + ":" + lsn})
		}
	}

	request := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: "dbs/db1/colls/collName_1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationReadFeed,
	}

	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(request, "range_1")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(1), sessionToken.GetLSN())

	request.RequestContext = &SessionRequestContext{
		ResolvedPartitionKeyRange: &PartitionKeyRangeInfo{
			ID:      fmt.Sprintf("range_%d", numPartitionKeyRangeIDs+10),
			Parents: []string{"range_2", "range_x"},
		},
	}

	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(request, request.RequestContext.ResolvedPartitionKeyRange.ID)
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(2), sessionToken.GetLSN())
}

func TestSessionContainer_SetSessionToken_NoSessionTokenForPartitionKeyRangeID(t *testing.T) {
	collectionRID := "uf4PAK6T-Cw="
	partitionKeyRangeID := "test_range_id"
	sessionToken := "1#100#1=20#2=5#3=30"
	collectionName := "dbs/db1/colls/collName_1"

	sessionContainer := NewSessionContainer("127.0.0.1")

	request1 := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionName + "/docs",
		ResourceType:    ResourceDocument,
		OperationType:   OperationCreate,
	}

	respHeaders := map[string]string{
		HTTPHeaderSessionToken:  partitionKeyRangeID + ":" + sessionToken,
		HTTPHeaderOwnerFullName: collectionName,
		HTTPHeaderOwnerID:       collectionRID,
	}
	sessionContainer.SetSessionToken(request1, respHeaders)

	request2 := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionName + "/docs",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}

	resolvedSessionToken := sessionContainer.ResolvePartitionLocalSessionToken(request2, partitionKeyRangeID)
	require.NotNil(t, resolvedSessionToken)
	require.Equal(t, sessionToken, resolvedSessionToken.ConvertToString())
}

func TestSessionContainer_SetSessionToken_MergeOldWithNew(t *testing.T) {
	collectionRID := "uf4PAK6T-Cw="
	collectionName := "dbs/db1/colls/collName_1"
	initialSessionToken := "1#100#1=20#2=5#3=30"
	newSessionTokenInServerResponse := "1#100#1=31#2=5#3=21"
	partitionKeyRangeID := "test_range_id"
	expectedMergedSessionToken := "1#100#1=31#2=5#3=30"

	sessionContainer := NewSessionContainer("127.0.0.1")

	request1 := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionName + "/docs",
		ResourceType:    ResourceDocument,
		OperationType:   OperationCreate,
	}

	respHeaders := map[string]string{
		HTTPHeaderSessionToken:  partitionKeyRangeID + ":" + initialSessionToken,
		HTTPHeaderOwnerFullName: collectionName,
		HTTPHeaderOwnerID:       collectionRID,
	}
	sessionContainer.SetSessionToken(request1, respHeaders)

	respHeaders[HTTPHeaderSessionToken] = partitionKeyRangeID + ":" + newSessionTokenInServerResponse
	sessionContainer.SetSessionToken(request1, respHeaders)

	request2 := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionName + "/docs",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}

	resolvedSessionToken := sessionContainer.ResolvePartitionLocalSessionToken(request2, partitionKeyRangeID)
	require.NotNil(t, resolvedSessionToken)
	require.Equal(t, expectedMergedSessionToken, resolvedSessionToken.ConvertToString())
}

func TestSessionContainer_ResolveGlobalSessionTokenReturnsEmptyStringOnEmptyCache(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	request := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: "dbs/db1/colls/collName/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}
	require.Equal(t, "", sessionContainer.ResolveGlobalSessionToken(request))
}

func TestSessionContainer_ResolveGlobalSessionTokenReturnsEmptyStringOnCacheMiss(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	partitionKeyRangeID := "range_0"
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	initialSessionToken := "1#100#1=20#2=5#3=30"

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, "dbs/db1/colls1/collName",
		map[string]string{HTTPHeaderSessionToken: partitionKeyRangeID + ":" + initialSessionToken})

	request := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: "dbs/db1/colls1/collName2/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}
	require.Equal(t, "", sessionContainer.ResolveGlobalSessionToken(request))
}

func TestSessionContainer_ResolveGlobalSessionTokenReturnsTokenMapUsingName(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#1=20#2=5#3=30"})
	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_1:1#101#1=20#2=5#3=30"})

	request := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}
	sessionToken := sessionContainer.ResolveGlobalSessionToken(request)
	tokens := strings.Split(sessionToken, ",")

	require.Equal(t, 2, len(tokens))
	tokenSet := make(map[string]bool)
	for _, tok := range tokens {
		tokenSet[tok] = true
	}
	require.True(t, tokenSet["range_0:1#100#1=20#2=5#3=30"])
	require.True(t, tokenSet["range_1:1#101#1=20#2=5#3=30"])
}

func TestSessionContainer_ResolveGlobalSessionTokenReturnsTokenMapUsingResourceID(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	request := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#1=20#2=5#3=30"})
	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_1:1#101#1=20#2=5#3=30"})

	sessionToken := sessionContainer.ResolveGlobalSessionToken(request)
	tokens := strings.Split(sessionToken, ",")

	require.Equal(t, 2, len(tokens))
	tokenSet := make(map[string]bool)
	for _, tok := range tokens {
		tokenSet[tok] = true
	}
	require.True(t, tokenSet["range_0:1#100#1=20#2=5#3=30"])
	require.True(t, tokenSet["range_1:1#101#1=20#2=5#3=30"])
}

func TestSessionContainer_ResolveLocalSessionTokenReturnsTokenMapUsingName(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#1=20#2=5#3=30"})
	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_1:1#101#1=20#2=5#3=30"})

	request := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}

	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(request, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())

	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(request, "range_1")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(101), sessionToken.GetLSN())
}

func TestSessionContainer_ResolveLocalSessionTokenReturnsTokenMapUsingResourceID(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	request := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#1=20#2=5#3=30"})
	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_1:1#101#1=20#2=5#3=30"})

	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(request, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())

	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(request, "range_1")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(101), sessionToken.GetLSN())
}

func TestSessionContainer_ResolveLocalSessionTokenReturnsNullOnPartitionMiss(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	request := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
		RequestContext: &SessionRequestContext{
			ResolvedPartitionKeyRange: &PartitionKeyRangeInfo{},
		},
	}

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#1=20#2=5#3=30"})
	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_1:1#101#1=20#2=5#3=30"})

	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(request, "range_2")
	require.Nil(t, sessionToken)
}

func TestSessionContainer_ResolveLocalSessionTokenReturnsNullOnCollectionMiss(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	randomCollectionID := rand.Int()
	documentCollectionID := fmt.Sprintf("collection_rid_%d", randomCollectionID)
	collectionFullName := "dbs/db1/colls1/collName"

	request := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    fmt.Sprintf("collection_rid_%d", randomCollectionID-1),
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
		RequestContext: &SessionRequestContext{
			ResolvedPartitionKeyRange: &PartitionKeyRangeInfo{},
		},
	}

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#1=20#2=5#3=30"})
	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_1:1#101#1=20#2=5#3=30"})

	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(request, "range_1")
	require.Nil(t, sessionToken)
}

func TestSessionContainer_ResolvePartitionLocalSessionTokenReturnsTokenOnParentMatch(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	request := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
		RequestContext: &SessionRequestContext{
			ResolvedPartitionKeyRange: &PartitionKeyRangeInfo{
				Parents: []string{"range_1"},
			},
		},
	}

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#1=20#2=5#3=30"})
	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_1:1#101#1=20#2=5#3=30"})

	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(request, "range_2")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(101), sessionToken.GetLSN())
}

func TestSessionContainer_ResolvePartitionLocalSessionTokenMergesMultipleParents(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	// Two parents with different region LSNs:
	// parent_a: globalLSN=100, region1=20, region2=5
	// parent_b: globalLSN=90,  region1=10, region2=15
	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "parent_a:1#100#1=20#2=5"})
	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "parent_b:1#90#1=10#2=15"})

	// Request for a child range whose parents are both parent_a and parent_b
	request := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
		RequestContext: &SessionRequestContext{
			ResolvedPartitionKeyRange: &PartitionKeyRangeInfo{
				Parents: []string{"parent_a", "parent_b"},
			},
		},
	}

	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(request, "child_range")
	require.NotNil(t, sessionToken)
	// Merged token should have max of each region:
	// globalLSN = max(100, 90) = 100
	// region1 = max(20, 10) = 20
	// region2 = max(5, 15) = 15
	require.Equal(t, int64(100), sessionToken.GetLSN())
	require.Equal(t, "1#100#1=20#2=15", sessionToken.ConvertToString())
}

func TestSessionContainer_ClearTokenByCollectionFullNameRemovesToken(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#1=20#2=5#3=30"})

	requestByRID := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}
	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(requestByRID, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())

	requestByName := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}
	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(requestByName, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())

	sessionContainer.ClearTokenByCollectionFullName(collectionFullName)

	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(requestByRID, "range_0")
	require.Nil(t, sessionToken)

	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(requestByName, "range_0")
	require.Nil(t, sessionToken)
}

func TestSessionContainer_ClearTokenByResourceIDRemovesToken(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName,
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#1=20#2=5#3=30"})

	requestByRID := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}
	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(requestByRID, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())

	requestByName := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}
	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(requestByName, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())

	sessionContainer.ClearTokenByResourceID(documentCollectionID)

	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(requestByRID, "range_0")
	require.Nil(t, sessionToken)

	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(requestByName, "range_0")
	require.Nil(t, sessionToken)
}

func TestSessionContainer_ClearTokenKeepsUnmatchedCollection(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	randomCollectionID := rand.Int()
	documentCollectionID1 := fmt.Sprintf("collection_rid_%d", randomCollectionID)
	collectionFullName1 := "dbs/db1/colls1/collName1"

	sessionContainer.SetSessionTokenFromRID(documentCollectionID1, collectionFullName1,
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#1=20#2=5#3=30"})

	request1 := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID1,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}

	documentCollectionID2 := fmt.Sprintf("collection_rid_%d", randomCollectionID-1)
	collectionFullName2 := "dbs/db1/colls1/collName2"

	request2 := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID2,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}

	sessionContainer.SetSessionTokenFromRID(documentCollectionID2, collectionFullName2,
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#1=20#2=5#3=30"})

	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(request1, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())

	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(request2, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())

	sessionContainer.ClearTokenByResourceID(documentCollectionID2)

	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(request1, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())

	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(request2, "range_0")
	require.Nil(t, sessionToken)
}

func TestSessionContainer_SetSessionTokenDoesntFailOnEmptySessionTokenHeader(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	sessionContainer.SetSessionToken(nil, map[string]string{})
}

func TestSessionContainer_SetSessionTokenSetsTokenWhenRequestIsntNameBased(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	request := &SessionContainerRequest{
		IsNameBased:     false,
		ResourceID:      documentCollectionID,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}

	sessionContainer.SetSessionToken(request, map[string]string{HTTPHeaderSessionToken: "range_0:1#100#4=90#5=1"})

	requestByRID := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}
	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(requestByRID, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())

	requestByName := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}
	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(requestByName, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())
}

func TestSessionContainer_SetSessionTokenGivesPriorityToOwnerFullNameOverResourceAddress(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName1 := "dbs/db1/colls1/collName1"
	collectionFullName2 := "dbs/db1/colls1/collName2"

	request := &SessionContainerRequest{
		IsNameBased:     false,
		ResourceID:      documentCollectionID,
		ResourceAddress: collectionFullName1 + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}
	sessionContainer.SetSessionToken(request, map[string]string{
		HTTPHeaderSessionToken:  "range_0:1#100#4=90#5=1",
		HTTPHeaderOwnerFullName: collectionFullName2,
	})

	requestByName1 := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionFullName1 + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}
	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(requestByName1, "range_0")
	require.Nil(t, sessionToken)

	requestByName2 := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionFullName2 + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}
	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(requestByName2, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())
}

func TestSessionContainer_SetSessionTokenIgnoresOwnerIDWhenRequestIsntNameBased(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	randomCollectionID := rand.Int()
	documentCollectionID1 := fmt.Sprintf("collection_rid_%d", randomCollectionID)
	documentCollectionID2 := fmt.Sprintf("collection_rid_%d", randomCollectionID-1)
	collectionFullName := "dbs/db1/colls1/collName1"

	request := &SessionContainerRequest{
		IsNameBased:     false,
		ResourceID:      documentCollectionID1,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}

	sessionContainer.SetSessionToken(request, map[string]string{
		HTTPHeaderSessionToken: "range_0:1#100#4=90#5=1",
		HTTPHeaderOwnerID:      documentCollectionID2,
	})

	requestByRID1 := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID1,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}
	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(requestByRID1, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())

	requestByRID2 := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID2,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}
	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(requestByRID2, "range_0")
	require.Nil(t, sessionToken)
}

func TestSessionContainer_SetSessionTokenGivesPriorityToOwnerIDOverResourceIDWhenRequestIsNameBased(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	randomCollectionID := rand.Int()
	documentCollectionID1 := fmt.Sprintf("collection_rid_%d", randomCollectionID)
	documentCollectionID2 := fmt.Sprintf("collection_rid_%d", randomCollectionID-1)
	collectionFullName := "dbs/db1/colls1/collName1"

	request := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceID:      documentCollectionID1,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}

	sessionContainer.SetSessionToken(request, map[string]string{
		HTTPHeaderSessionToken: "range_0:1#100#4=90#5=1",
		HTTPHeaderOwnerID:      documentCollectionID2,
	})

	requestByRID1 := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID1,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}
	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(requestByRID1, "range_0")
	require.Nil(t, sessionToken)

	requestByRID2 := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID2,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}
	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(requestByRID2, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(100), sessionToken.GetLSN())
}

func TestSessionContainer_SetSessionTokenDoesntWorkForMasterQueries(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	request := &SessionContainerRequest{
		IsNameBased:     false,
		ResourceID:      documentCollectionID,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceCollection,
		OperationType:   OperationReadFeed,
	}

	sessionContainer.SetSessionToken(request, map[string]string{HTTPHeaderSessionToken: "range_0:1"})

	requestByRID := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}
	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(requestByRID, "range_0")
	require.Nil(t, sessionToken)

	requestByName := &SessionContainerRequest{
		IsNameBased:     true,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}
	sessionToken = sessionContainer.ResolvePartitionLocalSessionToken(requestByName, "range_0")
	require.Nil(t, sessionToken)
}

func TestSessionContainer_SetSessionTokenDoesntOverwriteHigherLSN(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	request := &SessionContainerRequest{
		IsNameBased:     false,
		ResourceID:      documentCollectionID,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}

	sessionContainer.SetSessionToken(request, map[string]string{HTTPHeaderSessionToken: "range_0:1#105#4=90#5=1"})
	sessionContainer.SetSessionToken(request, map[string]string{HTTPHeaderSessionToken: "range_0:1#100#4=90#5=1"})

	requestByRID := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}
	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(requestByRID, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(105), sessionToken.GetLSN())
}

func TestSessionContainer_SetSessionTokenOverwriteLowerLSN(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	request := &SessionContainerRequest{
		IsNameBased:     false,
		ResourceID:      documentCollectionID,
		ResourceAddress: collectionFullName + "/docs/doc1",
		ResourceType:    ResourceDocument,
		OperationType:   OperationRead,
	}

	sessionContainer.SetSessionToken(request, map[string]string{HTTPHeaderSessionToken: "range_0:1#100#4=90#5=1"})
	sessionContainer.SetSessionToken(request, map[string]string{HTTPHeaderSessionToken: "range_0:1#105#4=90#5=1"})

	requestByRID := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}
	sessionToken := sessionContainer.ResolvePartitionLocalSessionToken(requestByRID, "range_0")
	require.NotNil(t, sessionToken)
	require.Equal(t, int64(105), sessionToken.GetLSN())
}

func TestSessionContainer_SetSessionTokenDoesNothingOnEmptySessionTokenHeader(t *testing.T) {
	sessionContainer := NewSessionContainer("127.0.0.1")
	documentCollectionID := fmt.Sprintf("collection_rid_%d", rand.Int())
	collectionFullName := "dbs/db1/colls1/collName"

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName+"/docs/doc1",
		map[string]string{HTTPHeaderSessionToken: "range_0:1#100#4=90#5=1"})

	request := &SessionContainerRequest{
		IsNameBased:   false,
		ResourceID:    documentCollectionID,
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
	}
	sessionToken := sessionContainer.ResolveGlobalSessionToken(request)
	tokens := strings.Split(sessionToken, ",")
	require.Equal(t, 1, len(tokens))
	require.Equal(t, "range_0:1#100#4=90#5=1", tokens[0])

	sessionContainer.SetSessionTokenFromRID(documentCollectionID, collectionFullName, map[string]string{})
	sessionToken = sessionContainer.ResolveGlobalSessionToken(request)
	tokens = strings.Split(sessionToken, ",")
	require.Equal(t, 1, len(tokens))
	require.Equal(t, "range_0:1#100#4=90#5=1", tokens[0])
}

func TestSessionTokenHelper_Parse(t *testing.T) {
	helper := &SessionTokenHelper{}

	token, err := helper.Parse("1#100#1=20#2=5#3=30")
	require.NoError(t, err)
	require.NotNil(t, token)
	require.Equal(t, int64(100), token.GetLSN())

	token, err = helper.Parse("range_0:1#100#1=20#2=5#3=30")
	require.NoError(t, err)
	require.NotNil(t, token)
	require.Equal(t, int64(100), token.GetLSN())
}

func TestSessionTokenHelper_TryParse(t *testing.T) {
	helper := &SessionTokenHelper{}

	token, ok := helper.TryParse("")
	require.False(t, ok)
	require.Nil(t, token)

	token, ok = helper.TryParse("1#100#1=20#2=5#3=30")
	require.True(t, ok)
	require.NotNil(t, token)
	require.Equal(t, int64(100), token.GetLSN())
}

func TestSessionTokenHelper_ResolvePartitionLocalSessionToken(t *testing.T) {
	helper := &SessionTokenHelper{}

	// Test 1: Direct match - should return the exact token
	request := &SessionContainerRequest{
		RequestContext: &SessionRequestContext{
			ResolvedPartitionKeyRange: &PartitionKeyRangeInfo{
				Parents: []string{"range_1"},
			},
		},
	}

	globalToken := "range_0:1#100#1=20#2=5#3=30,range_1:1#101#1=20#2=5#3=30"
	token, err := helper.ResolvePartitionLocalSessionToken(request, "range_0", globalToken)
	require.NoError(t, err)
	require.NotNil(t, token)
	require.Equal(t, int64(100), token.GetLSN())

	// Test 2: No direct match, but parent exists - should return parent's token
	// Looking for range_2 (not in global token), parent range_1 exists
	token, err = helper.ResolvePartitionLocalSessionToken(request, "range_2", globalToken)
	require.NoError(t, err)
	require.NotNil(t, token)
	require.Equal(t, int64(101), token.GetLSN())

	// Test 3: No direct match, no parent match - should return nil
	// Create request with parents that don't exist in the global token
	requestNoParentMatch := &SessionContainerRequest{
		RequestContext: &SessionRequestContext{
			ResolvedPartitionKeyRange: &PartitionKeyRangeInfo{
				Parents: []string{"range_99", "range_100"},
			},
		},
	}
	token, err = helper.ResolvePartitionLocalSessionToken(requestNoParentMatch, "range_3", globalToken)
	require.NoError(t, err)
	require.Nil(t, token)

	// Test 4: No request context - should return nil for non-matching partition
	token, err = helper.ResolvePartitionLocalSessionToken(nil, "range_3", globalToken)
	require.NoError(t, err)
	require.Nil(t, token)

	// Test 5: Empty global token - should return nil
	token, err = helper.ResolvePartitionLocalSessionToken(request, "range_0", "")
	require.NoError(t, err)
	require.Nil(t, token)
}

func TestGetCollectionPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"dbs/db1/colls/coll1/docs/doc1", "dbs/db1/colls/coll1"},
		{"dbs/db1/colls/coll1", "dbs/db1/colls/coll1"},
		{"/dbs/db1/colls/coll1/docs/doc1", "dbs/db1/colls/coll1"},
		{"dbs/db1/colls/coll1/", "dbs/db1/colls/coll1"},
		{"short", "short"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := getCollectionPath(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestIsReadingFromMaster(t *testing.T) {
	// Collection: only ReadFeed, Query, SQLQuery are master reads
	require.True(t, isReadingFromMaster(ResourceCollection, OperationReadFeed))
	require.True(t, isReadingFromMaster(ResourceCollection, OperationQuery))
	require.True(t, isReadingFromMaster(ResourceCollection, OperationSQLQuery))
	require.False(t, isReadingFromMaster(ResourceCollection, OperationRead))
	require.False(t, isReadingFromMaster(ResourceCollection, OperationHead))
	require.False(t, isReadingFromMaster(ResourceCollection, OperationHeadFeed))
	require.False(t, isReadingFromMaster(ResourceCollection, OperationCreate))

	// PartitionKeyRange: master except GetSplitPoint and AbortSplit
	require.True(t, isReadingFromMaster(ResourcePartitionKeyRange, OperationRead))
	require.True(t, isReadingFromMaster(ResourcePartitionKeyRange, OperationReadFeed))
	require.False(t, isReadingFromMaster(ResourcePartitionKeyRange, OperationGetSplitPoint))
	require.False(t, isReadingFromMaster(ResourcePartitionKeyRange, OperationAbortSplit))

	// Other master resources
	require.True(t, isReadingFromMaster(ResourceDatabase, OperationCreate))
	require.True(t, isReadingFromMaster(ResourceTopology, OperationRead))
	require.True(t, isReadingFromMaster(ResourceUserDefinedType, OperationRead))

	// Non-master resources
	require.False(t, isReadingFromMaster(ResourceDocument, OperationRead))
	require.False(t, isReadingFromMaster(ResourceDocument, OperationCreate))
}
