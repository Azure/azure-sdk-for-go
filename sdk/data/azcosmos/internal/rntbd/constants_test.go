// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"testing"
)

func TestCurrentProtocolVersion(t *testing.T) {
	// Protocol version must match Java SDK: 0x00000001
	if CurrentProtocolVersion != 0x00000001 {
		t.Errorf("CurrentProtocolVersion = %#x, want %#x", CurrentProtocolVersion, 0x00000001)
	}
}

func TestConsistencyLevel(t *testing.T) {
	tests := []struct {
		level ConsistencyLevel
		id    byte
		name  string
	}{
		{ConsistencyStrong, 0x00, "Strong"},
		{ConsistencyBoundedStaleness, 0x01, "BoundedStaleness"},
		{ConsistencySession, 0x02, "Session"},
		{ConsistencyEventual, 0x03, "Eventual"},
		{ConsistencyConsistentPrefix, 0x04, "ConsistentPrefix"},
		{ConsistencyInvalid, 0xFF, "Invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if byte(tt.level) != tt.id {
				t.Errorf("ConsistencyLevel %s = %#x, want %#x", tt.name, byte(tt.level), tt.id)
			}
			if tt.level.String() != tt.name {
				t.Errorf("ConsistencyLevel.String() = %q, want %q", tt.level.String(), tt.name)
			}
		})
	}
}

func TestTokenType(t *testing.T) {
	// Token types must match Java SDK exactly (from RntbdTokenType.java)
	tests := []struct {
		tokenType TokenType
		id        byte
		name      string
	}{
		{TokenByte, 0x00, "Byte"},
		{TokenUShort, 0x01, "UShort"},
		{TokenULong, 0x02, "ULong"},
		{TokenLong, 0x03, "Long"},
		{TokenULongLong, 0x04, "ULongLong"},
		{TokenLongLong, 0x05, "LongLong"},
		{TokenGuid, 0x06, "Guid"},
		{TokenSmallString, 0x07, "SmallString"},
		{TokenString, 0x08, "String"},
		{TokenULongString, 0x09, "ULongString"},
		{TokenSmallBytes, 0x0A, "SmallBytes"},
		{TokenBytes, 0x0B, "Bytes"},
		{TokenULongBytes, 0x0C, "ULongBytes"},
		{TokenFloat, 0x0D, "Float"},
		{TokenDouble, 0x0E, "Double"},
		{TokenInvalid, 0xFF, "Invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if byte(tt.tokenType) != tt.id {
				t.Errorf("TokenType %s = %#x, want %#x", tt.name, byte(tt.tokenType), tt.id)
			}
			if tt.tokenType.String() != tt.name {
				t.Errorf("TokenType.String() = %q, want %q", tt.tokenType.String(), tt.name)
			}
		})
	}
}

func TestOperationType(t *testing.T) {
	// Operation types must match Java SDK exactly
	tests := []struct {
		opType OperationType
		id     uint16
		name   string
	}{
		{OperationConnection, 0x0000, "Connection"},
		{OperationCreate, 0x0001, "Create"},
		{OperationPatch, 0x0002, "Patch"},
		{OperationRead, 0x0003, "Read"},
		{OperationReadFeed, 0x0004, "ReadFeed"},
		{OperationDelete, 0x0005, "Delete"},
		{OperationReplace, 0x0006, "Replace"},
		// 0x0007 is obsolete (JPathQuery)
		{OperationExecuteJavaScript, 0x0008, "ExecuteJavaScript"},
		{OperationSQLQuery, 0x0009, "SQLQuery"},
		{OperationPause, 0x000A, "Pause"},
		{OperationResume, 0x000B, "Resume"},
		{OperationStop, 0x000C, "Stop"},
		{OperationRecycle, 0x000D, "Recycle"},
		{OperationCrash, 0x000E, "Crash"},
		{OperationQuery, 0x000F, "Query"},
		{OperationForceConfigRefresh, 0x0010, "ForceConfigRefresh"},
		{OperationHead, 0x0011, "Head"},
		{OperationHeadFeed, 0x0012, "HeadFeed"},
		{OperationUpsert, 0x0013, "Upsert"},
		{OperationRecreate, 0x0014, "Recreate"},
		{OperationThrottle, 0x0015, "Throttle"},
		{OperationGetSplitPoint, 0x0016, "GetSplitPoint"},
		{OperationPreCreateValidation, 0x0017, "PreCreateValidation"},
		{OperationBatchApply, 0x0018, "BatchApply"},
		{OperationAbortSplit, 0x0019, "AbortSplit"},
		{OperationCompleteSplit, 0x001A, "CompleteSplit"},
		{OperationOfferUpdateOperation, 0x001B, "OfferUpdateOperation"},
		{OperationOfferPreGrowValidation, 0x001C, "OfferPreGrowValidation"},
		{OperationBatchReportThroughputUtilization, 0x001D, "BatchReportThroughputUtilization"},
		{OperationCompletePartitionMigration, 0x001E, "CompletePartitionMigration"},
		{OperationAbortPartitionMigration, 0x001F, "AbortPartitionMigration"},
		{OperationPreReplaceValidation, 0x0020, "PreReplaceValidation"},
		{OperationAddComputeGatewayRequestCharges, 0x0021, "AddComputeGatewayRequestCharges"},
		{OperationMigratePartition, 0x0022, "MigratePartition"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if uint16(tt.opType) != tt.id {
				t.Errorf("OperationType %s = %#x, want %#x", tt.name, uint16(tt.opType), tt.id)
			}
			if tt.opType.String() != tt.name {
				t.Errorf("OperationType.String() = %q, want %q", tt.opType.String(), tt.name)
			}
		})
	}
}

func TestResourceType(t *testing.T) {
	// Resource types must match Java SDK exactly
	tests := []struct {
		resType ResourceType
		id      uint16
		name    string
	}{
		{ResourceConnection, 0x0000, "Connection"},
		{ResourceDatabase, 0x0001, "Database"},
		{ResourceCollection, 0x0002, "Collection"},
		{ResourceDocument, 0x0003, "Document"},
		{ResourceAttachment, 0x0004, "Attachment"},
		{ResourceUser, 0x0005, "User"},
		{ResourcePermission, 0x0006, "Permission"},
		{ResourceStoredProcedure, 0x0007, "StoredProcedure"},
		{ResourceConflict, 0x0008, "Conflict"},
		{ResourceTrigger, 0x0009, "Trigger"},
		{ResourceUserDefinedFunction, 0x000A, "UserDefinedFunction"},
		{ResourceModule, 0x000B, "Module"},
		{ResourceReplica, 0x000C, "Replica"},
		{ResourceModuleCommand, 0x000D, "ModuleCommand"},
		{ResourceRecord, 0x000E, "Record"},
		{ResourceOffer, 0x000F, "Offer"},
		{ResourcePartitionSetInformation, 0x0010, "PartitionSetInformation"},
		{ResourceXPReplicatorAddress, 0x0011, "XPReplicatorAddress"},
		{ResourceMasterPartition, 0x0012, "MasterPartition"},
		{ResourceServerPartition, 0x0013, "ServerPartition"},
		{ResourceDatabaseAccount, 0x0014, "DatabaseAccount"},
		{ResourceTopology, 0x0015, "Topology"},
		{ResourcePartitionKeyRange, 0x0016, "PartitionKeyRange"},
		// 0x0017 is obsolete (Timestamp)
		{ResourceSchema, 0x0018, "Schema"},
		{ResourceBatchApply, 0x0019, "BatchApply"},
		{ResourceRestoreMetadata, 0x001A, "RestoreMetadata"},
		{ResourceComputeGatewayCharges, 0x001B, "ComputeGatewayCharges"},
		{ResourceRidRange, 0x001C, "RidRange"},
		{ResourceUserDefinedType, 0x001D, "UserDefinedType"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if uint16(tt.resType) != tt.id {
				t.Errorf("ResourceType %s = %#x, want %#x", tt.name, uint16(tt.resType), tt.id)
			}
			if tt.resType.String() != tt.name {
				t.Errorf("ResourceType.String() = %q, want %q", tt.resType.String(), tt.name)
			}
		})
	}
}

func TestRequestHeaderMapping(t *testing.T) {
	// Verify all request headers are in the map with correct types
	// Selected critical headers to verify
	tests := []struct {
		header     RequestHeader
		id         uint16
		tokenType  TokenType
		isRequired bool
	}{
		{RequestHeaderResourceId, 0x0000, TokenBytes, false},
		{RequestHeaderAuthorizationToken, 0x0001, TokenString, false},
		{RequestHeaderPayloadPresent, 0x0002, TokenByte, true},
		{RequestHeaderSessionToken, 0x0005, TokenString, false},
		{RequestHeaderConsistencyLevel, 0x0010, TokenByte, false},
		{RequestHeaderPartitionKey, 0x002B, TokenString, false},
		{RequestHeaderPartitionKeyRangeId, 0x002C, TokenString, false},
		{RequestHeaderEffectivePartitionKey, 0x005A, TokenBytes, false},
	}

	for _, tt := range tests {
		t.Run(RequestHeaders[tt.header].Type.String(), func(t *testing.T) {
			if uint16(tt.header) != tt.id {
				t.Errorf("RequestHeader ID = %#x, want %#x", uint16(tt.header), tt.id)
			}

			info, ok := RequestHeaders[tt.header]
			if !ok {
				t.Fatalf("RequestHeader %#x not found in map", tt.header)
			}

			if info.Type != tt.tokenType {
				t.Errorf("RequestHeader[%#x].Type = %s, want %s", tt.header, info.Type, tt.tokenType)
			}

			if info.IsRequired != tt.isRequired {
				t.Errorf("RequestHeader[%#x].IsRequired = %v, want %v", tt.header, info.IsRequired, tt.isRequired)
			}
		})
	}
}

func TestResponseHeaderMapping(t *testing.T) {
	// Verify all response headers are in the map with correct types
	// Selected critical headers to verify
	tests := []struct {
		header     ResponseHeader
		id         uint16
		tokenType  TokenType
		isRequired bool
	}{
		{ResponseHeaderPayloadPresent, 0x0000, TokenByte, true},
		{ResponseHeaderContinuationToken, 0x0003, TokenString, false},
		{ResponseHeaderETag, 0x0004, TokenString, false},
		{ResponseHeaderLSN, 0x0013, TokenLongLong, false},
		{ResponseHeaderRequestCharge, 0x0015, TokenDouble, false},
		{ResponseHeaderSubStatus, 0x001C, TokenULong, false},
		{ResponseHeaderSessionToken, 0x003E, TokenString, false},
	}

	for _, tt := range tests {
		t.Run(ResponseHeaders[tt.header].Type.String(), func(t *testing.T) {
			if uint16(tt.header) != tt.id {
				t.Errorf("ResponseHeader ID = %#x, want %#x", uint16(tt.header), tt.id)
			}

			info, ok := ResponseHeaders[tt.header]
			if !ok {
				t.Fatalf("ResponseHeader %#x not found in map", tt.header)
			}

			if info.Type != tt.tokenType {
				t.Errorf("ResponseHeader[%#x].Type = %s, want %s", tt.header, info.Type, tt.tokenType)
			}

			if info.IsRequired != tt.isRequired {
				t.Errorf("ResponseHeader[%#x].IsRequired = %v, want %v", tt.header, info.IsRequired, tt.isRequired)
			}
		})
	}
}

func TestContextHeaders(t *testing.T) {
	tests := []struct {
		header     ContextHeader
		id         uint16
		tokenType  TokenType
		isRequired bool
	}{
		{ContextHeaderProtocolVersion, 0x0000, TokenULong, false},
		{ContextHeaderClientVersion, 0x0001, TokenSmallString, false},
		{ContextHeaderServerAgent, 0x0002, TokenSmallString, true},
		{ContextHeaderServerVersion, 0x0003, TokenSmallString, true},
		{ContextHeaderIdleTimeoutInSeconds, 0x0004, TokenULong, false},
		{ContextHeaderUnauthenticatedTimeoutInSeconds, 0x0005, TokenULong, false},
	}

	for _, tt := range tests {
		t.Run("ContextHeader", func(t *testing.T) {
			if uint16(tt.header) != tt.id {
				t.Errorf("ContextHeader ID = %#x, want %#x", uint16(tt.header), tt.id)
			}

			info, ok := ContextHeaders[tt.header]
			if !ok {
				t.Fatalf("ContextHeader %#x not found in map", tt.header)
			}

			if info.Type != tt.tokenType {
				t.Errorf("ContextHeader[%#x].Type = %s, want %s", tt.header, info.Type, tt.tokenType)
			}

			if info.IsRequired != tt.isRequired {
				t.Errorf("ContextHeader[%#x].IsRequired = %v, want %v", tt.header, info.IsRequired, tt.isRequired)
			}
		})
	}
}

func TestEnumerationDirection(t *testing.T) {
	tests := []struct {
		dir  EnumerationDirection
		id   byte
		name string
	}{
		{EnumerationDirectionInvalid, 0x00, "Invalid"},
		{EnumerationDirectionForward, 0x01, "Forward"},
		{EnumerationDirectionReverse, 0x02, "Reverse"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if byte(tt.dir) != tt.id {
				t.Errorf("EnumerationDirection %s = %#x, want %#x", tt.name, byte(tt.dir), tt.id)
			}
			if tt.dir.String() != tt.name {
				t.Errorf("EnumerationDirection.String() = %q, want %q", tt.dir.String(), tt.name)
			}
		})
	}
}

func TestIndexingDirective(t *testing.T) {
	tests := []struct {
		dir  IndexingDirective
		id   byte
		name string
	}{
		{IndexingDirectiveDefault, 0x00, "Default"},
		{IndexingDirectiveInclude, 0x01, "Include"},
		{IndexingDirectiveExclude, 0x02, "Exclude"},
		{IndexingDirectiveInvalid, 0xFF, "Invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if byte(tt.dir) != tt.id {
				t.Errorf("IndexingDirective %s = %#x, want %#x", tt.name, byte(tt.dir), tt.id)
			}
			if tt.dir.String() != tt.name {
				t.Errorf("IndexingDirective.String() = %q, want %q", tt.dir.String(), tt.name)
			}
		})
	}
}

func TestContentSerializationFormat(t *testing.T) {
	tests := []struct {
		format ContentSerializationFormat
		id     byte
		name   string
	}{
		{ContentSerializationJsonText, 0x00, "JsonText"},
		{ContentSerializationCosmosBinary, 0x01, "CosmosBinary"},
		{ContentSerializationInvalid, 0xFF, "Invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if byte(tt.format) != tt.id {
				t.Errorf("ContentSerializationFormat %s = %#x, want %#x", tt.name, byte(tt.format), tt.id)
			}
			if tt.format.String() != tt.name {
				t.Errorf("ContentSerializationFormat.String() = %q, want %q", tt.format.String(), tt.name)
			}
		})
	}
}

// TestRequestHeaderCount verifies we have the expected number of request headers
func TestRequestHeaderCount(t *testing.T) {
	// Java SDK has headers from 0x0000-0x0068 (105 slots)
	// We include "NotUsed" entries to maintain ID alignment for future extensions
	// Direct mode implementation added 6 additional headers (0x0064-0x0069)
	expectedCount := 106

	if len(RequestHeaders) != expectedCount {
		t.Errorf("RequestHeaders count = %d, want %d", len(RequestHeaders), expectedCount)
	}
}

// TestResponseHeaderCount verifies we have the expected number of response headers
func TestResponseHeaderCount(t *testing.T) {
	// Java SDK defines response headers from 0x0000 to 0x003E
	// Includes gaps for future expansion
	expectedCount := 49 // All defined headers in our implementation

	if len(ResponseHeaders) != expectedCount {
		t.Errorf("ResponseHeaders count = %d, want %d", len(ResponseHeaders), expectedCount)
	}
}
