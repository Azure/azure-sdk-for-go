// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBuildRequestMessage_BasicRequest(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-5678-1234-567812345678")
	req := &ServiceRequest{
		OperationType:   OperationRead,
		ResourceType:    ResourceDocument,
		ResourceAddress: "/dbs/testdb/colls/testcoll/docs/testdoc",
		IsNameBased:     true,
		ActivityID:      activityID,
		Headers:         make(map[string]string),
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)
	require.NotNil(t, msg)

	require.Equal(t, OperationRead, msg.Frame.OperationType)
	require.Equal(t, ResourceDocument, msg.Frame.ResourceType)
	require.Equal(t, activityID, msg.Frame.ActivityID)

	payloadPresent := msg.Headers.GetByte(uint16(RequestHeaderPayloadPresent))
	require.Equal(t, byte(0), payloadPresent)
}

func TestBuildRequestMessage_WithPayload(t *testing.T) {
	activityID := uuid.New()
	content := []byte(`{"id":"test","pk":"value"}`)
	req := &ServiceRequest{
		OperationType:   OperationCreate,
		ResourceType:    ResourceDocument,
		ResourceAddress: "/dbs/testdb/colls/testcoll/docs",
		IsNameBased:     true,
		ActivityID:      activityID,
		Headers:         make(map[string]string),
		Content:         content,
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)
	require.NotNil(t, msg)

	payloadPresent := msg.Headers.GetByte(uint16(RequestHeaderPayloadPresent))
	require.Equal(t, byte(1), payloadPresent)
	require.Equal(t, content, msg.Payload)
}

func TestBuildRequestMessage_NameBasedHeaders(t *testing.T) {
	activityID := uuid.New()
	req := &ServiceRequest{
		OperationType:   OperationRead,
		ResourceType:    ResourceDocument,
		ResourceAddress: "/dbs/mydb/colls/mycoll/docs/mydoc",
		IsNameBased:     true,
		ActivityID:      activityID,
		Headers:         make(map[string]string),
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)

	dbName := msg.Headers.GetString(uint16(RequestHeaderDatabaseName))
	require.Equal(t, "mydb", dbName)

	collName := msg.Headers.GetString(uint16(RequestHeaderCollectionName))
	require.Equal(t, "mycoll", collName)

	docName := msg.Headers.GetString(uint16(RequestHeaderDocumentName))
	require.Equal(t, "mydoc", docName)
}

func TestBuildRequestMessage_ConsistencyLevel(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected ConsistencyLevel
	}{
		{"Strong", "Strong", ConsistencyStrong},
		{"BoundedStaleness", "BoundedStaleness", ConsistencyBoundedStaleness},
		{"Bounded", "Bounded", ConsistencyBoundedStaleness},
		{"Session", "Session", ConsistencySession},
		{"Eventual", "Eventual", ConsistencyEventual},
		{"ConsistentPrefix", "ConsistentPrefix", ConsistencyConsistentPrefix},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &ServiceRequest{
				OperationType:   OperationRead,
				ResourceType:    ResourceDocument,
				ResourceAddress: "dbs/db1/colls/coll1/docs/doc1",
				IsNameBased:     true,
				ActivityID:      uuid.New(),
				Headers: map[string]string{
					HTTPHeaderConsistencyLevel: tc.header,
				},
			}

			msg, err := BuildRequestMessage(req)
			require.NoError(t, err)

			level := msg.Headers.GetByte(uint16(RequestHeaderConsistencyLevel))
			require.Equal(t, byte(tc.expected), level)
		})
	}
}

func TestBuildRequestMessage_SessionToken(t *testing.T) {
	req := &ServiceRequest{
		OperationType:   OperationRead,
		ResourceType:    ResourceDocument,
		ResourceAddress: "dbs/db1/colls/coll1/docs/doc1",
		IsNameBased:     true,
		ActivityID:      uuid.New(),
		Headers: map[string]string{
			HTTPHeaderSessionToken: "0:1#1234#56=789",
		},
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)

	sessionToken := msg.Headers.GetString(uint16(RequestHeaderSessionToken))
	require.Equal(t, "0:1#1234#56=789", sessionToken)
}

func TestBuildRequestMessage_PartitionKey(t *testing.T) {
	req := &ServiceRequest{
		OperationType:   OperationRead,
		ResourceType:    ResourceDocument,
		ResourceAddress: "dbs/db1/colls/coll1/docs/doc1",
		IsNameBased:     true,
		ActivityID:      uuid.New(),
		Headers: map[string]string{
			HTTPHeaderPartitionKey: "[\"myPartitionKey\"]",
		},
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)

	partitionKey := msg.Headers.GetString(uint16(RequestHeaderPartitionKey))
	require.Equal(t, "[\"myPartitionKey\"]", partitionKey)
}

func TestBuildRequestMessage_WithTransportRequestID(t *testing.T) {
	req := &ServiceRequest{
		OperationType:      OperationRead,
		ResourceType:       ResourceDocument,
		ResourceAddress:    "dbs/db1/colls/coll1/docs/doc1",
		IsNameBased:        true,
		ActivityID:         uuid.New(),
		Headers:            make(map[string]string),
		TransportRequestID: 12345,
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)

	transportID := msg.Headers.GetULong(uint16(RequestHeaderTransportRequestID))
	require.Equal(t, uint32(12345), transportID)
}

func TestBuildRequestMessage_WithReplicaPath(t *testing.T) {
	req := &ServiceRequest{
		OperationType:   OperationRead,
		ResourceType:    ResourceDocument,
		ResourceAddress: "dbs/db1/colls/coll1/docs/doc1",
		IsNameBased:     true,
		ActivityID:      uuid.New(),
		Headers:         make(map[string]string),
		ReplicaPath:     "rntbd://host:443/replica1",
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)

	replicaPath := msg.Headers.GetString(uint16(RequestHeaderReplicaPath))
	require.Equal(t, "rntbd://host:443/replica1", replicaPath)
}

func TestBuildRequestMessage_ResourceIDPath(t *testing.T) {
	req := &ServiceRequest{
		OperationType:   OperationRead,
		ResourceType:    ResourceDocument,
		ResourceID:      "SomeResourceId",
		ResourceAddress: "SomeResourceId",
		IsNameBased:     false,
		ActivityID:      uuid.New(),
		Headers:         make(map[string]string),
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)

	resourceID := msg.Headers.GetBytes(uint16(RequestHeaderResourceId))
	require.NotNil(t, resourceID)
	require.Equal(t, "SomeResourceId", string(resourceID))
}

func TestBuildRequestMessage_PageSize(t *testing.T) {
	req := &ServiceRequest{
		OperationType:   OperationReadFeed,
		ResourceType:    ResourceDocument,
		ResourceAddress: "dbs/db1/colls/coll1",
		IsNameBased:     true,
		ActivityID:      uuid.New(),
		Headers: map[string]string{
			HTTPHeaderPageSize: "100",
		},
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)

	pageSize := msg.Headers.GetULong(uint16(RequestHeaderPageSize))
	require.Equal(t, uint32(100), pageSize)
}

func TestBuildRequestMessage_IndexingDirective(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected IndexingDirective
	}{
		{"Default", "Default", IndexingDirectiveDefault},
		{"Include", "Include", IndexingDirectiveInclude},
		{"Exclude", "Exclude", IndexingDirectiveExclude},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &ServiceRequest{
				OperationType:   OperationCreate,
				ResourceType:    ResourceDocument,
				ResourceAddress: "dbs/db1/colls/coll1/docs",
				IsNameBased:     true,
				ActivityID:      uuid.New(),
				Headers: map[string]string{
					HTTPHeaderIndexingDirective: tc.header,
				},
			}

			msg, err := BuildRequestMessage(req)
			require.NoError(t, err)

			directive := msg.Headers.GetByte(uint16(RequestHeaderIndexingDirective))
			require.Equal(t, byte(tc.expected), directive)
		})
	}
}

func TestBuildRequestMessage_BooleanHeaders(t *testing.T) {
	req := &ServiceRequest{
		OperationType:   OperationSQLQuery,
		ResourceType:    ResourceDocument,
		ResourceAddress: "dbs/db1/colls/coll1",
		IsNameBased:     true,
		ActivityID:      uuid.New(),
		Headers: map[string]string{
			HTTPHeaderEnableScanInQuery:       "true",
			HTTPHeaderPopulateQueryMetrics:    "true",
			HTTPHeaderDisableRUPerMinuteUsage: "false",
		},
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)

	enableScan := msg.Headers.GetByte(uint16(RequestHeaderEnableScanInQuery))
	require.Equal(t, byte(1), enableScan)

	populateMetrics := msg.Headers.GetByte(uint16(RequestHeaderPopulateQueryMetrics))
	require.Equal(t, byte(1), populateMetrics)

	disableRUPerMinute := msg.Headers.GetByte(uint16(RequestHeaderDisableRUPerMinuteUsage))
	require.Equal(t, byte(0), disableRUPerMinute)
}

func TestBuildRequestMessage_Continuation(t *testing.T) {
	continuationToken := `{"token":"abc123","range":{"min":"","max":"FF"}}`
	req := &ServiceRequest{
		OperationType:   OperationReadFeed,
		ResourceType:    ResourceDocument,
		ResourceAddress: "dbs/db1/colls/coll1",
		IsNameBased:     true,
		ActivityID:      uuid.New(),
		Headers: map[string]string{
			HTTPHeaderContinuation: continuationToken,
		},
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)

	continuation := msg.Headers.GetString(uint16(RequestHeaderContinuationToken))
	require.Equal(t, continuationToken, continuation)
}

func TestBuildRequestMessage_RoundTrip(t *testing.T) {
	activityID := uuid.MustParse("12345678-abcd-5678-abcd-567812345678")
	content := []byte(`{"id":"roundtrip","data":"test"}`)
	req := &ServiceRequest{
		OperationType:      OperationCreate,
		ResourceType:       ResourceDocument,
		ResourceAddress:    "/dbs/testdb/colls/testcoll/docs",
		IsNameBased:        true,
		ActivityID:         activityID,
		TransportRequestID: 9999,
		Content:            content,
		Headers: map[string]string{
			HTTPHeaderConsistencyLevel: "Session",
			HTTPHeaderPartitionKey:     "[\"pk\"]",
			HTTPHeaderSessionToken:     "0:1#100#5=200",
		},
	}

	msg, err := BuildRequestMessage(req)
	require.NoError(t, err)

	encoded, err := EncodeRequestToBytes(msg)
	require.NoError(t, err)

	decoded, err := DecodeRequestFromBytes(encoded)
	require.NoError(t, err)

	require.Equal(t, msg.Frame.OperationType, decoded.Frame.OperationType)
	require.Equal(t, msg.Frame.ResourceType, decoded.Frame.ResourceType)
	require.Equal(t, msg.Frame.ActivityID, decoded.Frame.ActivityID)
	require.Equal(t, msg.Payload, decoded.Payload)

	require.Equal(t, msg.Headers.GetByte(uint16(RequestHeaderPayloadPresent)),
		decoded.Headers.GetByte(uint16(RequestHeaderPayloadPresent)))
	require.Equal(t, msg.Headers.GetByte(uint16(RequestHeaderConsistencyLevel)),
		decoded.Headers.GetByte(uint16(RequestHeaderConsistencyLevel)))
	require.Equal(t, msg.Headers.GetString(uint16(RequestHeaderPartitionKey)),
		decoded.Headers.GetString(uint16(RequestHeaderPartitionKey)))
}

func TestParseResourcePath(t *testing.T) {
	tests := []struct {
		path     string
		expected map[string]string
	}{
		{
			path: "/dbs/mydb/colls/mycoll/docs/mydoc",
			expected: map[string]string{
				"database":   "mydb",
				"collection": "mycoll",
				"document":   "mydoc",
			},
		},
		{
			path: "dbs/db1/colls/coll1",
			expected: map[string]string{
				"database":   "db1",
				"collection": "coll1",
			},
		},
		{
			path: "/dbs/db/colls/coll/sprocs/sp1",
			expected: map[string]string{
				"database":        "db",
				"collection":      "coll",
				"storedProcedure": "sp1",
			},
		},
		{
			path: "/dbs/db/colls/coll/triggers/tr1",
			expected: map[string]string{
				"database":   "db",
				"collection": "coll",
				"trigger":    "tr1",
			},
		},
		{
			path: "/dbs/db/colls/coll/udfs/udf1",
			expected: map[string]string{
				"database":            "db",
				"collection":          "coll",
				"userDefinedFunction": "udf1",
			},
		},
		{
			path: "/dbs/db/users/user1/permissions/perm1",
			expected: map[string]string{
				"database":   "db",
				"user":       "user1",
				"permission": "perm1",
			},
		},
		{
			path:     "",
			expected: map[string]string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			result := ParseResourcePath(tc.path)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestConvertToByte(t *testing.T) {
	tests := []struct {
		header   string
		value    string
		expected byte
	}{
		{HTTPHeaderConsistencyLevel, "Strong", byte(ConsistencyStrong)},
		{HTTPHeaderConsistencyLevel, "strong", byte(ConsistencyStrong)},
		{HTTPHeaderConsistencyLevel, "Session", byte(ConsistencySession)},
		{HTTPHeaderIndexingDirective, "Include", byte(IndexingDirectiveInclude)},
		{HTTPHeaderIndexingDirective, "exclude", byte(IndexingDirectiveExclude)},
		{"x-some-bool-header", "true", byte(1)},
		{"x-some-bool-header", "True", byte(1)},
		{"x-some-bool-header", "1", byte(1)},
		{"x-some-bool-header", "false", byte(0)},
		{"x-some-bool-header", "0", byte(0)},
	}

	for _, tc := range tests {
		t.Run(tc.header+"_"+tc.value, func(t *testing.T) {
			result := convertToByte(tc.header, tc.value)
			require.Equal(t, tc.expected, result)
		})
	}
}
