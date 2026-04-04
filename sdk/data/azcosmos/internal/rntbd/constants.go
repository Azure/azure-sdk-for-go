// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package rntbd implements the RNTBD (Remote Native Transport Binary Direct) protocol
// for direct connectivity to Azure Cosmos DB.
package rntbd

// CurrentProtocolVersion is the RNTBD protocol version supported by this implementation.
const CurrentProtocolVersion uint32 = 0x00000001

// ConsistencyLevel represents the RNTBD consistency level encoding.
type ConsistencyLevel byte

const (
	ConsistencyStrong           ConsistencyLevel = 0x00
	ConsistencyBoundedStaleness ConsistencyLevel = 0x01
	ConsistencySession          ConsistencyLevel = 0x02
	ConsistencyEventual         ConsistencyLevel = 0x03
	ConsistencyConsistentPrefix ConsistencyLevel = 0x04
	ConsistencyInvalid          ConsistencyLevel = 0xFF
)

// String returns the string representation of the consistency level.
func (c ConsistencyLevel) String() string {
	switch c {
	case ConsistencyStrong:
		return "Strong"
	case ConsistencyBoundedStaleness:
		return "BoundedStaleness"
	case ConsistencySession:
		return "Session"
	case ConsistencyEventual:
		return "Eventual"
	case ConsistencyConsistentPrefix:
		return "ConsistentPrefix"
	case ConsistencyInvalid:
		return "Invalid"
	default:
		return "Unknown"
	}
}

// ContentSerializationFormat represents the content serialization format.
type ContentSerializationFormat byte

const (
	ContentSerializationJsonText     ContentSerializationFormat = 0x00
	ContentSerializationCosmosBinary ContentSerializationFormat = 0x01
	ContentSerializationInvalid      ContentSerializationFormat = 0xFF
)

// String returns the string representation of the content serialization format.
func (f ContentSerializationFormat) String() string {
	switch f {
	case ContentSerializationJsonText:
		return "JsonText"
	case ContentSerializationCosmosBinary:
		return "CosmosBinary"
	case ContentSerializationInvalid:
		return "Invalid"
	default:
		return "Unknown"
	}
}

// EnumerationDirection represents the direction of enumeration.
type EnumerationDirection byte

const (
	EnumerationDirectionInvalid EnumerationDirection = 0x00
	EnumerationDirectionForward EnumerationDirection = 0x01
	EnumerationDirectionReverse EnumerationDirection = 0x02
)

// String returns the string representation of the enumeration direction.
func (d EnumerationDirection) String() string {
	switch d {
	case EnumerationDirectionInvalid:
		return "Invalid"
	case EnumerationDirectionForward:
		return "Forward"
	case EnumerationDirectionReverse:
		return "Reverse"
	default:
		return "Unknown"
	}
}

// FanoutOperationState represents the state of a fanout operation.
type FanoutOperationState byte

const (
	FanoutOperationStarted   FanoutOperationState = 0x01
	FanoutOperationCompleted FanoutOperationState = 0x02
)

// String returns the string representation of the fanout operation state.
func (s FanoutOperationState) String() string {
	switch s {
	case FanoutOperationStarted:
		return "Started"
	case FanoutOperationCompleted:
		return "Completed"
	default:
		return "Unknown"
	}
}

// IndexingDirective represents the indexing directive.
type IndexingDirective byte

const (
	IndexingDirectiveDefault IndexingDirective = 0x00
	IndexingDirectiveInclude IndexingDirective = 0x01
	IndexingDirectiveExclude IndexingDirective = 0x02
	IndexingDirectiveInvalid IndexingDirective = 0xFF
)

// String returns the string representation of the indexing directive.
func (d IndexingDirective) String() string {
	switch d {
	case IndexingDirectiveDefault:
		return "Default"
	case IndexingDirectiveInclude:
		return "Include"
	case IndexingDirectiveExclude:
		return "Exclude"
	case IndexingDirectiveInvalid:
		return "Invalid"
	default:
		return "Unknown"
	}
}

// MigrateCollectionDirective represents the migrate collection directive.
type MigrateCollectionDirective byte

const (
	MigrateCollectionThaw    MigrateCollectionDirective = 0x00
	MigrateCollectionFreeze  MigrateCollectionDirective = 0x01
	MigrateCollectionInvalid MigrateCollectionDirective = 0xFF
)

// String returns the string representation of the migrate collection directive.
func (d MigrateCollectionDirective) String() string {
	switch d {
	case MigrateCollectionThaw:
		return "Thaw"
	case MigrateCollectionFreeze:
		return "Freeze"
	case MigrateCollectionInvalid:
		return "Invalid"
	default:
		return "Unknown"
	}
}

// ReadFeedKeyType represents the read feed key type.
type ReadFeedKeyType byte

const (
	ReadFeedKeyTypeInvalid               ReadFeedKeyType = 0x00
	ReadFeedKeyTypeResourceId            ReadFeedKeyType = 0x01
	ReadFeedKeyTypeEffectivePartitionKey ReadFeedKeyType = 0x02
)

// String returns the string representation of the read feed key type.
func (t ReadFeedKeyType) String() string {
	switch t {
	case ReadFeedKeyTypeInvalid:
		return "Invalid"
	case ReadFeedKeyTypeResourceId:
		return "ResourceId"
	case ReadFeedKeyTypeEffectivePartitionKey:
		return "EffectivePartitionKey"
	default:
		return "Unknown"
	}
}

// RemoteStorageType represents the remote storage type.
type RemoteStorageType byte

const (
	RemoteStorageTypeInvalid      RemoteStorageType = 0x00
	RemoteStorageTypeNotSpecified RemoteStorageType = 0x01
	RemoteStorageTypeStandard     RemoteStorageType = 0x02
	RemoteStorageTypePremium      RemoteStorageType = 0x03
)

// String returns the string representation of the remote storage type.
func (t RemoteStorageType) String() string {
	switch t {
	case RemoteStorageTypeInvalid:
		return "Invalid"
	case RemoteStorageTypeNotSpecified:
		return "NotSpecified"
	case RemoteStorageTypeStandard:
		return "Standard"
	case RemoteStorageTypePremium:
		return "Premium"
	default:
		return "Unknown"
	}
}

// OperationType represents the RNTBD operation type.
type OperationType uint16

const (
	OperationConnection OperationType = 0x0000
	OperationCreate     OperationType = 0x0001
	OperationUpdate     OperationType = 0x0002
	OperationRead       OperationType = 0x0003
	OperationReadFeed   OperationType = 0x0004
	OperationDelete     OperationType = 0x0005
	OperationReplace    OperationType = 0x0006
	// 0x0007 is obsolete (JPathQuery)
	OperationExecuteJavaScript                OperationType = 0x0008
	OperationSQLQuery                         OperationType = 0x0009
	OperationPause                            OperationType = 0x000A
	OperationResume                           OperationType = 0x000B
	OperationStop                             OperationType = 0x000C
	OperationRecycle                          OperationType = 0x000D
	OperationCrash                            OperationType = 0x000E
	OperationQuery                            OperationType = 0x000F
	OperationForceConfigRefresh               OperationType = 0x0010
	OperationHead                             OperationType = 0x0011
	OperationHeadFeed                         OperationType = 0x0012
	OperationUpsert                           OperationType = 0x0013
	OperationRecreate                         OperationType = 0x0014
	OperationThrottle                         OperationType = 0x0015
	OperationGetSplitPoint                    OperationType = 0x0016
	OperationPreCreateValidation              OperationType = 0x0017
	OperationBatchApply                       OperationType = 0x0018
	OperationAbortSplit                       OperationType = 0x0019
	OperationCompleteSplit                    OperationType = 0x001A
	OperationOfferUpdateOperation             OperationType = 0x001B
	OperationOfferPreGrowValidation           OperationType = 0x001C
	OperationBatchReportThroughputUtilization OperationType = 0x001D
	OperationCompletePartitionMigration       OperationType = 0x001E
	OperationAbortPartitionMigration          OperationType = 0x001F
	OperationPreReplaceValidation             OperationType = 0x0020
	OperationAddComputeGatewayRequestCharges  OperationType = 0x0021
	OperationMigratePartition                 OperationType = 0x0022
)

// String returns the string representation of the operation type.
func (o OperationType) String() string {
	switch o {
	case OperationConnection:
		return "Connection"
	case OperationCreate:
		return "Create"
	case OperationUpdate:
		return "Update"
	case OperationRead:
		return "Read"
	case OperationReadFeed:
		return "ReadFeed"
	case OperationDelete:
		return "Delete"
	case OperationReplace:
		return "Replace"
	case OperationExecuteJavaScript:
		return "ExecuteJavaScript"
	case OperationSQLQuery:
		return "SQLQuery"
	case OperationPause:
		return "Pause"
	case OperationResume:
		return "Resume"
	case OperationStop:
		return "Stop"
	case OperationRecycle:
		return "Recycle"
	case OperationCrash:
		return "Crash"
	case OperationQuery:
		return "Query"
	case OperationForceConfigRefresh:
		return "ForceConfigRefresh"
	case OperationHead:
		return "Head"
	case OperationHeadFeed:
		return "HeadFeed"
	case OperationUpsert:
		return "Upsert"
	case OperationRecreate:
		return "Recreate"
	case OperationThrottle:
		return "Throttle"
	case OperationGetSplitPoint:
		return "GetSplitPoint"
	case OperationPreCreateValidation:
		return "PreCreateValidation"
	case OperationBatchApply:
		return "BatchApply"
	case OperationAbortSplit:
		return "AbortSplit"
	case OperationCompleteSplit:
		return "CompleteSplit"
	case OperationOfferUpdateOperation:
		return "OfferUpdateOperation"
	case OperationOfferPreGrowValidation:
		return "OfferPreGrowValidation"
	case OperationBatchReportThroughputUtilization:
		return "BatchReportThroughputUtilization"
	case OperationCompletePartitionMigration:
		return "CompletePartitionMigration"
	case OperationAbortPartitionMigration:
		return "AbortPartitionMigration"
	case OperationPreReplaceValidation:
		return "PreReplaceValidation"
	case OperationAddComputeGatewayRequestCharges:
		return "AddComputeGatewayRequestCharges"
	case OperationMigratePartition:
		return "MigratePartition"
	default:
		return "Unknown"
	}
}

// ResourceType represents the RNTBD resource type.
type ResourceType uint16

const (
	ResourceConnection              ResourceType = 0x0000
	ResourceDatabase                ResourceType = 0x0001
	ResourceCollection              ResourceType = 0x0002
	ResourceDocument                ResourceType = 0x0003
	ResourceAttachment              ResourceType = 0x0004
	ResourceUser                    ResourceType = 0x0005
	ResourcePermission              ResourceType = 0x0006
	ResourceStoredProcedure         ResourceType = 0x0007
	ResourceConflict                ResourceType = 0x0008
	ResourceTrigger                 ResourceType = 0x0009
	ResourceUserDefinedFunction     ResourceType = 0x000A
	ResourceModule                  ResourceType = 0x000B
	ResourceReplica                 ResourceType = 0x000C
	ResourceModuleCommand           ResourceType = 0x000D
	ResourceRecord                  ResourceType = 0x000E
	ResourceOffer                   ResourceType = 0x000F
	ResourcePartitionSetInformation ResourceType = 0x0010
	ResourceXPReplicatorAddress     ResourceType = 0x0011
	ResourceMasterPartition         ResourceType = 0x0012
	ResourceServerPartition         ResourceType = 0x0013
	ResourceDatabaseAccount         ResourceType = 0x0014
	ResourceTopology                ResourceType = 0x0015
	ResourcePartitionKeyRange       ResourceType = 0x0016
	// 0x0017 is obsolete (Timestamp)
	ResourceSchema                ResourceType = 0x0018
	ResourceBatchApply            ResourceType = 0x0019
	ResourceRestoreMetadata       ResourceType = 0x001A
	ResourceComputeGatewayCharges ResourceType = 0x001B
	ResourceRidRange              ResourceType = 0x001C
	ResourceUserDefinedType       ResourceType = 0x001D
)

// String returns the string representation of the resource type.
func (r ResourceType) String() string {
	switch r {
	case ResourceConnection:
		return "Connection"
	case ResourceDatabase:
		return "Database"
	case ResourceCollection:
		return "Collection"
	case ResourceDocument:
		return "Document"
	case ResourceAttachment:
		return "Attachment"
	case ResourceUser:
		return "User"
	case ResourcePermission:
		return "Permission"
	case ResourceStoredProcedure:
		return "StoredProcedure"
	case ResourceConflict:
		return "Conflict"
	case ResourceTrigger:
		return "Trigger"
	case ResourceUserDefinedFunction:
		return "UserDefinedFunction"
	case ResourceModule:
		return "Module"
	case ResourceReplica:
		return "Replica"
	case ResourceModuleCommand:
		return "ModuleCommand"
	case ResourceRecord:
		return "Record"
	case ResourceOffer:
		return "Offer"
	case ResourcePartitionSetInformation:
		return "PartitionSetInformation"
	case ResourceXPReplicatorAddress:
		return "XPReplicatorAddress"
	case ResourceMasterPartition:
		return "MasterPartition"
	case ResourceServerPartition:
		return "ServerPartition"
	case ResourceDatabaseAccount:
		return "DatabaseAccount"
	case ResourceTopology:
		return "Topology"
	case ResourcePartitionKeyRange:
		return "PartitionKeyRange"
	case ResourceSchema:
		return "Schema"
	case ResourceBatchApply:
		return "BatchApply"
	case ResourceRestoreMetadata:
		return "RestoreMetadata"
	case ResourceComputeGatewayCharges:
		return "ComputeGatewayCharges"
	case ResourceRidRange:
		return "RidRange"
	case ResourceUserDefinedType:
		return "UserDefinedType"
	default:
		return "Unknown"
	}
}

// TokenType represents the RNTBD token type for binary protocol encoding.
type TokenType byte

const (
	TokenByte        TokenType = 0x00
	TokenUShort      TokenType = 0x01
	TokenULong       TokenType = 0x02
	TokenLong        TokenType = 0x03
	TokenULongLong   TokenType = 0x04
	TokenLongLong    TokenType = 0x05
	TokenGuid        TokenType = 0x06
	TokenSmallString TokenType = 0x07
	TokenString      TokenType = 0x08
	TokenULongString TokenType = 0x09
	TokenSmallBytes  TokenType = 0x0A
	TokenBytes       TokenType = 0x0B
	TokenULongBytes  TokenType = 0x0C
	TokenFloat       TokenType = 0x0D
	TokenDouble      TokenType = 0x0E
	TokenInvalid     TokenType = 0xFF
)

// String returns the string representation of the token type.
func (t TokenType) String() string {
	switch t {
	case TokenByte:
		return "Byte"
	case TokenUShort:
		return "UShort"
	case TokenULong:
		return "ULong"
	case TokenLong:
		return "Long"
	case TokenULongLong:
		return "ULongLong"
	case TokenLongLong:
		return "LongLong"
	case TokenGuid:
		return "Guid"
	case TokenSmallString:
		return "SmallString"
	case TokenString:
		return "String"
	case TokenULongString:
		return "ULongString"
	case TokenSmallBytes:
		return "SmallBytes"
	case TokenBytes:
		return "Bytes"
	case TokenULongBytes:
		return "ULongBytes"
	case TokenFloat:
		return "Float"
	case TokenDouble:
		return "Double"
	case TokenInvalid:
		return "Invalid"
	default:
		return "Unknown"
	}
}

// ContextHeader represents headers sent in context negotiation.
type ContextHeader uint16

const (
	ContextHeaderProtocolVersion                 ContextHeader = 0x0000
	ContextHeaderClientVersion                   ContextHeader = 0x0001
	ContextHeaderServerAgent                     ContextHeader = 0x0002
	ContextHeaderServerVersion                   ContextHeader = 0x0003
	ContextHeaderIdleTimeoutInSeconds            ContextHeader = 0x0004
	ContextHeaderUnauthenticatedTimeoutInSeconds ContextHeader = 0x0005
)

// ContextHeaderInfo provides metadata about context headers.
type ContextHeaderInfo struct {
	ID         ContextHeader
	Type       TokenType
	IsRequired bool
}

// ContextHeaders maps header IDs to their metadata.
var ContextHeaders = map[ContextHeader]ContextHeaderInfo{
	ContextHeaderProtocolVersion:                 {ContextHeaderProtocolVersion, TokenULong, false},
	ContextHeaderClientVersion:                   {ContextHeaderClientVersion, TokenSmallString, false},
	ContextHeaderServerAgent:                     {ContextHeaderServerAgent, TokenSmallString, true},
	ContextHeaderServerVersion:                   {ContextHeaderServerVersion, TokenSmallString, true},
	ContextHeaderIdleTimeoutInSeconds:            {ContextHeaderIdleTimeoutInSeconds, TokenULong, false},
	ContextHeaderUnauthenticatedTimeoutInSeconds: {ContextHeaderUnauthenticatedTimeoutInSeconds, TokenULong, false},
}

// ContextRequestHeader represents headers sent in context request.
type ContextRequestHeader uint16

const (
	ContextRequestHeaderProtocolVersion ContextRequestHeader = 0x0000
	ContextRequestHeaderClientVersion   ContextRequestHeader = 0x0001
	ContextRequestHeaderUserAgent       ContextRequestHeader = 0x0002
)

// ContextRequestHeaderInfo provides metadata about context request headers.
type ContextRequestHeaderInfo struct {
	ID         ContextRequestHeader
	Type       TokenType
	IsRequired bool
}

// ContextRequestHeaders maps header IDs to their metadata.
var ContextRequestHeaders = map[ContextRequestHeader]ContextRequestHeaderInfo{
	ContextRequestHeaderProtocolVersion: {ContextRequestHeaderProtocolVersion, TokenULong, true},
	ContextRequestHeaderClientVersion:   {ContextRequestHeaderClientVersion, TokenSmallString, true},
	ContextRequestHeaderUserAgent:       {ContextRequestHeaderUserAgent, TokenSmallString, true},
}

// RequestHeader represents RNTBD request header IDs.
type RequestHeader uint16

const (
	RequestHeaderResourceId                RequestHeader = 0x0000
	RequestHeaderAuthorizationToken        RequestHeader = 0x0001
	RequestHeaderPayloadPresent            RequestHeader = 0x0002
	RequestHeaderDate                      RequestHeader = 0x0003
	RequestHeaderPageSize                  RequestHeader = 0x0004
	RequestHeaderSessionToken              RequestHeader = 0x0005
	RequestHeaderContinuationToken         RequestHeader = 0x0006
	RequestHeaderIndexingDirective         RequestHeader = 0x0007
	RequestHeaderMatch                     RequestHeader = 0x0008
	RequestHeaderPreTriggerInclude         RequestHeader = 0x0009
	RequestHeaderPostTriggerInclude        RequestHeader = 0x000A
	RequestHeaderIsFanout                  RequestHeader = 0x000B
	RequestHeaderCollectionPartitionIndex  RequestHeader = 0x000C
	RequestHeaderCollectionServiceIndex    RequestHeader = 0x000D
	RequestHeaderPreTriggerExclude         RequestHeader = 0x000E
	RequestHeaderPostTriggerExclude        RequestHeader = 0x000F
	RequestHeaderConsistencyLevel          RequestHeader = 0x0010
	RequestHeaderEntityId                  RequestHeader = 0x0011
	RequestHeaderResourceSchemaName        RequestHeader = 0x0012
	RequestHeaderReplicaPath               RequestHeader = 0x0013
	RequestHeaderResourceTokenExpiry       RequestHeader = 0x0014
	RequestHeaderDatabaseName              RequestHeader = 0x0015
	RequestHeaderCollectionName            RequestHeader = 0x0016
	RequestHeaderDocumentName              RequestHeader = 0x0017
	RequestHeaderAttachmentName            RequestHeader = 0x0018
	RequestHeaderUserName                  RequestHeader = 0x0019
	RequestHeaderPermissionName            RequestHeader = 0x001A
	RequestHeaderStoredProcedureName       RequestHeader = 0x001B
	RequestHeaderUserDefinedFunctionName   RequestHeader = 0x001C
	RequestHeaderTriggerName               RequestHeader = 0x001D
	RequestHeaderEnableScanInQuery         RequestHeader = 0x001E
	RequestHeaderEmitVerboseTracesInQuery  RequestHeader = 0x001F
	RequestHeaderConflictName              RequestHeader = 0x0020
	RequestHeaderBindReplicaDirective      RequestHeader = 0x0021
	RequestHeaderPrimaryMasterKey          RequestHeader = 0x0022
	RequestHeaderSecondaryMasterKey        RequestHeader = 0x0023
	RequestHeaderPrimaryReadonlyKey        RequestHeader = 0x0024
	RequestHeaderSecondaryReadonlyKey      RequestHeader = 0x0025
	RequestHeaderProfileRequest            RequestHeader = 0x0026
	RequestHeaderEnableLowPrecisionOrderBy RequestHeader = 0x0027
	RequestHeaderClientVersion             RequestHeader = 0x0028
	RequestHeaderCanCharge                 RequestHeader = 0x0029
	RequestHeaderCanThrottle               RequestHeader = 0x002A
	RequestHeaderPartitionKey              RequestHeader = 0x002B
	RequestHeaderPartitionKeyRangeId       RequestHeader = 0x002C
	RequestHeaderNotUsed2D                 RequestHeader = 0x002D // Not used
	RequestHeaderNotUsed2E                 RequestHeader = 0x002E // Not used
	RequestHeaderNotUsed2F                 RequestHeader = 0x002F // Not used
	// 0x0030 not used
	RequestHeaderMigrateCollectionDirective      RequestHeader = 0x0031
	RequestHeaderNotUsed32                       RequestHeader = 0x0032 // Not used
	RequestHeaderSupportSpatialLegacyCoordinates RequestHeader = 0x0033
	RequestHeaderPartitionCount                  RequestHeader = 0x0034
	RequestHeaderCollectionRid                   RequestHeader = 0x0035
	RequestHeaderPartitionKeyRangeName           RequestHeader = 0x0036
	// 0x0037-0x0039 not used (RoundTripTimeInMsec, RequestMessageSentTime, RequestMessageTimeOffset)
	RequestHeaderSchemaName                                RequestHeader = 0x003A
	RequestHeaderFilterBySchemaRid                         RequestHeader = 0x003B
	RequestHeaderUsePolygonsSmallerThanAHemisphere         RequestHeader = 0x003C
	RequestHeaderGatewaySignature                          RequestHeader = 0x003D
	RequestHeaderEnableLogging                             RequestHeader = 0x003E
	RequestHeaderA_IM                                      RequestHeader = 0x003F
	RequestHeaderPopulateQuotaInfo                         RequestHeader = 0x0040
	RequestHeaderDisableRUPerMinuteUsage                   RequestHeader = 0x0041
	RequestHeaderPopulateQueryMetrics                      RequestHeader = 0x0042
	RequestHeaderResponseContinuationTokenLimitInKb        RequestHeader = 0x0043
	RequestHeaderPopulatePartitionStatistics               RequestHeader = 0x0044
	RequestHeaderRemoteStorageType                         RequestHeader = 0x0045
	RequestHeaderCollectionRemoteStorageSecurityIdentifier RequestHeader = 0x0046
	RequestHeaderIfModifiedSince                           RequestHeader = 0x0047
	RequestHeaderPopulateCollectionThroughputInfo          RequestHeader = 0x0048
	RequestHeaderRemainingTimeInMsOnClientRequest          RequestHeader = 0x0049
	RequestHeaderClientRetryAttemptCount                   RequestHeader = 0x004A
	RequestHeaderTargetLsn                                 RequestHeader = 0x004B
	RequestHeaderTargetGlobalCommittedLsn                  RequestHeader = 0x004C
	RequestHeaderTransportRequestID                        RequestHeader = 0x004D
	RequestHeaderRestoreMetadaFilter                       RequestHeader = 0x004E
	RequestHeaderRestoreParams                             RequestHeader = 0x004F
	RequestHeaderShareThroughput                           RequestHeader = 0x0050
	RequestHeaderPartitionResourceFilter                   RequestHeader = 0x0051
	RequestHeaderIsReadOnlyScript                          RequestHeader = 0x0052
	RequestHeaderIsAutoScaleRequest                        RequestHeader = 0x0053
	RequestHeaderForceQueryScan                            RequestHeader = 0x0054
	// 0x0055 not used (LeaseSeqNumber)
	RequestHeaderCanOfferReplaceComplete         RequestHeader = 0x0056
	RequestHeaderExcludeSystemProperties         RequestHeader = 0x0057
	RequestHeaderBinaryId                        RequestHeader = 0x0058
	RequestHeaderTimeToLiveInSeconds             RequestHeader = 0x0059
	RequestHeaderEffectivePartitionKey           RequestHeader = 0x005A
	RequestHeaderBinaryPassthroughRequest        RequestHeader = 0x005B
	RequestHeaderUserDefinedTypeName             RequestHeader = 0x005C
	RequestHeaderEnableDynamicRidRangeAllocation RequestHeader = 0x005D
	RequestHeaderEnumerationDirection            RequestHeader = 0x005E
	RequestHeaderStartId                         RequestHeader = 0x005F
	RequestHeaderEndId                           RequestHeader = 0x0060
	RequestHeaderFanoutOperationState            RequestHeader = 0x0061
	RequestHeaderStartEpk                        RequestHeader = 0x0062
	RequestHeaderEndEpk                          RequestHeader = 0x0063
	RequestHeaderReadFeedKeyType                 RequestHeader = 0x0064
	RequestHeaderContentSerializationFormat      RequestHeader = 0x0065
	RequestHeaderAllowTentativeWrites            RequestHeader = 0x0066
	RequestHeaderIsUserRequest                   RequestHeader = 0x0067
	RequestHeaderSharedOfferThroughput           RequestHeader = 0x0068
)

// RequestHeaderInfo provides metadata about request headers.
type RequestHeaderInfo struct {
	ID         RequestHeader
	Type       TokenType
	IsRequired bool
}

// RequestHeaders maps header IDs to their metadata.
var RequestHeaders = map[RequestHeader]RequestHeaderInfo{
	RequestHeaderResourceId:                                {RequestHeaderResourceId, TokenBytes, false},
	RequestHeaderAuthorizationToken:                        {RequestHeaderAuthorizationToken, TokenString, false},
	RequestHeaderPayloadPresent:                            {RequestHeaderPayloadPresent, TokenByte, true},
	RequestHeaderDate:                                      {RequestHeaderDate, TokenSmallString, false},
	RequestHeaderPageSize:                                  {RequestHeaderPageSize, TokenULong, false},
	RequestHeaderSessionToken:                              {RequestHeaderSessionToken, TokenString, false},
	RequestHeaderContinuationToken:                         {RequestHeaderContinuationToken, TokenString, false},
	RequestHeaderIndexingDirective:                         {RequestHeaderIndexingDirective, TokenByte, false},
	RequestHeaderMatch:                                     {RequestHeaderMatch, TokenString, false},
	RequestHeaderPreTriggerInclude:                         {RequestHeaderPreTriggerInclude, TokenString, false},
	RequestHeaderPostTriggerInclude:                        {RequestHeaderPostTriggerInclude, TokenString, false},
	RequestHeaderIsFanout:                                  {RequestHeaderIsFanout, TokenByte, false},
	RequestHeaderCollectionPartitionIndex:                  {RequestHeaderCollectionPartitionIndex, TokenULong, false},
	RequestHeaderCollectionServiceIndex:                    {RequestHeaderCollectionServiceIndex, TokenULong, false},
	RequestHeaderPreTriggerExclude:                         {RequestHeaderPreTriggerExclude, TokenString, false},
	RequestHeaderPostTriggerExclude:                        {RequestHeaderPostTriggerExclude, TokenString, false},
	RequestHeaderConsistencyLevel:                          {RequestHeaderConsistencyLevel, TokenByte, false},
	RequestHeaderEntityId:                                  {RequestHeaderEntityId, TokenString, false},
	RequestHeaderResourceSchemaName:                        {RequestHeaderResourceSchemaName, TokenSmallString, false},
	RequestHeaderReplicaPath:                               {RequestHeaderReplicaPath, TokenString, true},
	RequestHeaderResourceTokenExpiry:                       {RequestHeaderResourceTokenExpiry, TokenULong, false},
	RequestHeaderDatabaseName:                              {RequestHeaderDatabaseName, TokenString, false},
	RequestHeaderCollectionName:                            {RequestHeaderCollectionName, TokenString, false},
	RequestHeaderDocumentName:                              {RequestHeaderDocumentName, TokenString, false},
	RequestHeaderAttachmentName:                            {RequestHeaderAttachmentName, TokenString, false},
	RequestHeaderUserName:                                  {RequestHeaderUserName, TokenString, false},
	RequestHeaderPermissionName:                            {RequestHeaderPermissionName, TokenString, false},
	RequestHeaderStoredProcedureName:                       {RequestHeaderStoredProcedureName, TokenString, false},
	RequestHeaderUserDefinedFunctionName:                   {RequestHeaderUserDefinedFunctionName, TokenString, false},
	RequestHeaderTriggerName:                               {RequestHeaderTriggerName, TokenString, false},
	RequestHeaderEnableScanInQuery:                         {RequestHeaderEnableScanInQuery, TokenByte, false},
	RequestHeaderEmitVerboseTracesInQuery:                  {RequestHeaderEmitVerboseTracesInQuery, TokenByte, false},
	RequestHeaderConflictName:                              {RequestHeaderConflictName, TokenString, false},
	RequestHeaderBindReplicaDirective:                      {RequestHeaderBindReplicaDirective, TokenString, false},
	RequestHeaderPrimaryMasterKey:                          {RequestHeaderPrimaryMasterKey, TokenString, false},
	RequestHeaderSecondaryMasterKey:                        {RequestHeaderSecondaryMasterKey, TokenString, false},
	RequestHeaderPrimaryReadonlyKey:                        {RequestHeaderPrimaryReadonlyKey, TokenString, false},
	RequestHeaderSecondaryReadonlyKey:                      {RequestHeaderSecondaryReadonlyKey, TokenString, false},
	RequestHeaderProfileRequest:                            {RequestHeaderProfileRequest, TokenByte, false},
	RequestHeaderEnableLowPrecisionOrderBy:                 {RequestHeaderEnableLowPrecisionOrderBy, TokenByte, false},
	RequestHeaderClientVersion:                             {RequestHeaderClientVersion, TokenSmallString, false},
	RequestHeaderCanCharge:                                 {RequestHeaderCanCharge, TokenByte, false},
	RequestHeaderCanThrottle:                               {RequestHeaderCanThrottle, TokenByte, false},
	RequestHeaderPartitionKey:                              {RequestHeaderPartitionKey, TokenString, false},
	RequestHeaderPartitionKeyRangeId:                       {RequestHeaderPartitionKeyRangeId, TokenString, false},
	RequestHeaderNotUsed2D:                                 {RequestHeaderNotUsed2D, TokenInvalid, false},
	RequestHeaderNotUsed2E:                                 {RequestHeaderNotUsed2E, TokenInvalid, false},
	RequestHeaderNotUsed2F:                                 {RequestHeaderNotUsed2F, TokenInvalid, false},
	RequestHeaderMigrateCollectionDirective:                {RequestHeaderMigrateCollectionDirective, TokenByte, false},
	RequestHeaderNotUsed32:                                 {RequestHeaderNotUsed32, TokenInvalid, false},
	RequestHeaderSupportSpatialLegacyCoordinates:           {RequestHeaderSupportSpatialLegacyCoordinates, TokenByte, false},
	RequestHeaderPartitionCount:                            {RequestHeaderPartitionCount, TokenULong, false},
	RequestHeaderCollectionRid:                             {RequestHeaderCollectionRid, TokenString, false},
	RequestHeaderPartitionKeyRangeName:                     {RequestHeaderPartitionKeyRangeName, TokenString, false},
	RequestHeaderSchemaName:                                {RequestHeaderSchemaName, TokenString, false},
	RequestHeaderFilterBySchemaRid:                         {RequestHeaderFilterBySchemaRid, TokenString, false},
	RequestHeaderUsePolygonsSmallerThanAHemisphere:         {RequestHeaderUsePolygonsSmallerThanAHemisphere, TokenByte, false},
	RequestHeaderGatewaySignature:                          {RequestHeaderGatewaySignature, TokenString, false},
	RequestHeaderEnableLogging:                             {RequestHeaderEnableLogging, TokenByte, false},
	RequestHeaderA_IM:                                      {RequestHeaderA_IM, TokenString, false},
	RequestHeaderPopulateQuotaInfo:                         {RequestHeaderPopulateQuotaInfo, TokenByte, false},
	RequestHeaderDisableRUPerMinuteUsage:                   {RequestHeaderDisableRUPerMinuteUsage, TokenByte, false},
	RequestHeaderPopulateQueryMetrics:                      {RequestHeaderPopulateQueryMetrics, TokenByte, false},
	RequestHeaderResponseContinuationTokenLimitInKb:        {RequestHeaderResponseContinuationTokenLimitInKb, TokenULong, false},
	RequestHeaderPopulatePartitionStatistics:               {RequestHeaderPopulatePartitionStatistics, TokenByte, false},
	RequestHeaderRemoteStorageType:                         {RequestHeaderRemoteStorageType, TokenByte, false},
	RequestHeaderCollectionRemoteStorageSecurityIdentifier: {RequestHeaderCollectionRemoteStorageSecurityIdentifier, TokenString, false},
	RequestHeaderIfModifiedSince:                           {RequestHeaderIfModifiedSince, TokenString, false},
	RequestHeaderPopulateCollectionThroughputInfo:          {RequestHeaderPopulateCollectionThroughputInfo, TokenByte, false},
	RequestHeaderRemainingTimeInMsOnClientRequest:          {RequestHeaderRemainingTimeInMsOnClientRequest, TokenULong, false},
	RequestHeaderClientRetryAttemptCount:                   {RequestHeaderClientRetryAttemptCount, TokenULong, false},
	RequestHeaderTargetLsn:                                 {RequestHeaderTargetLsn, TokenLongLong, false},
	RequestHeaderTargetGlobalCommittedLsn:                  {RequestHeaderTargetGlobalCommittedLsn, TokenLongLong, false},
	RequestHeaderTransportRequestID:                        {RequestHeaderTransportRequestID, TokenULong, false},
	RequestHeaderRestoreMetadaFilter:                       {RequestHeaderRestoreMetadaFilter, TokenString, false},
	RequestHeaderRestoreParams:                             {RequestHeaderRestoreParams, TokenString, false},
	RequestHeaderShareThroughput:                           {RequestHeaderShareThroughput, TokenByte, false},
	RequestHeaderPartitionResourceFilter:                   {RequestHeaderPartitionResourceFilter, TokenString, false},
	RequestHeaderIsReadOnlyScript:                          {RequestHeaderIsReadOnlyScript, TokenByte, false},
	RequestHeaderIsAutoScaleRequest:                        {RequestHeaderIsAutoScaleRequest, TokenByte, false},
	RequestHeaderForceQueryScan:                            {RequestHeaderForceQueryScan, TokenByte, false},
	RequestHeaderCanOfferReplaceComplete:                   {RequestHeaderCanOfferReplaceComplete, TokenByte, false},
	RequestHeaderExcludeSystemProperties:                   {RequestHeaderExcludeSystemProperties, TokenByte, false},
	RequestHeaderBinaryId:                                  {RequestHeaderBinaryId, TokenBytes, false},
	RequestHeaderTimeToLiveInSeconds:                       {RequestHeaderTimeToLiveInSeconds, TokenLong, false},
	RequestHeaderEffectivePartitionKey:                     {RequestHeaderEffectivePartitionKey, TokenBytes, false},
	RequestHeaderBinaryPassthroughRequest:                  {RequestHeaderBinaryPassthroughRequest, TokenByte, false},
	RequestHeaderUserDefinedTypeName:                       {RequestHeaderUserDefinedTypeName, TokenString, false},
	RequestHeaderEnableDynamicRidRangeAllocation:           {RequestHeaderEnableDynamicRidRangeAllocation, TokenByte, false},
	RequestHeaderEnumerationDirection:                      {RequestHeaderEnumerationDirection, TokenByte, false},
	RequestHeaderStartId:                                   {RequestHeaderStartId, TokenBytes, false},
	RequestHeaderEndId:                                     {RequestHeaderEndId, TokenBytes, false},
	RequestHeaderFanoutOperationState:                      {RequestHeaderFanoutOperationState, TokenByte, false},
	RequestHeaderStartEpk:                                  {RequestHeaderStartEpk, TokenBytes, false},
	RequestHeaderEndEpk:                                    {RequestHeaderEndEpk, TokenBytes, false},
	RequestHeaderReadFeedKeyType:                           {RequestHeaderReadFeedKeyType, TokenByte, false},
	RequestHeaderContentSerializationFormat:                {RequestHeaderContentSerializationFormat, TokenByte, false},
	RequestHeaderAllowTentativeWrites:                      {RequestHeaderAllowTentativeWrites, TokenByte, false},
	RequestHeaderIsUserRequest:                             {RequestHeaderIsUserRequest, TokenByte, false},
	RequestHeaderSharedOfferThroughput:                     {RequestHeaderSharedOfferThroughput, TokenULong, false},
}

// ResponseHeader represents RNTBD response header IDs.
type ResponseHeader uint16

const (
	ResponseHeaderPayloadPresent ResponseHeader = 0x0000
	// 0x0001 not used
	ResponseHeaderLastStateChangeDateTime ResponseHeader = 0x0002
	ResponseHeaderContinuationToken       ResponseHeader = 0x0003
	ResponseHeaderETag                    ResponseHeader = 0x0004
	// 0x0005-0x0006 not used
	ResponseHeaderReadsPerformed            ResponseHeader = 0x0007
	ResponseHeaderWritesPerformed           ResponseHeader = 0x0008
	ResponseHeaderQueriesPerformed          ResponseHeader = 0x0009
	ResponseHeaderIndexTermsGenerated       ResponseHeader = 0x000A
	ResponseHeaderScriptsExecuted           ResponseHeader = 0x000B
	ResponseHeaderRetryAfterMilliseconds    ResponseHeader = 0x000C
	ResponseHeaderIndexingDirective         ResponseHeader = 0x000D
	ResponseHeaderStorageMaxResoureQuota    ResponseHeader = 0x000E
	ResponseHeaderStorageResourceQuotaUsage ResponseHeader = 0x000F
	ResponseHeaderSchemaVersion             ResponseHeader = 0x0010
	ResponseHeaderCollectionPartitionIndex  ResponseHeader = 0x0011
	ResponseHeaderCollectionServiceIndex    ResponseHeader = 0x0012
	ResponseHeaderLSN                       ResponseHeader = 0x0013
	ResponseHeaderItemCount                 ResponseHeader = 0x0014
	ResponseHeaderRequestCharge             ResponseHeader = 0x0015
	// 0x0016 not used
	ResponseHeaderOwnerFullName               ResponseHeader = 0x0017
	ResponseHeaderOwnerId                     ResponseHeader = 0x0018
	ResponseHeaderDatabaseAccountId           ResponseHeader = 0x0019
	ResponseHeaderQuorumAckedLSN              ResponseHeader = 0x001A
	ResponseHeaderRequestValidationFailure    ResponseHeader = 0x001B
	ResponseHeaderSubStatus                   ResponseHeader = 0x001C
	ResponseHeaderCollectionUpdateProgress    ResponseHeader = 0x001D
	ResponseHeaderCurrentWriteQuorum          ResponseHeader = 0x001E
	ResponseHeaderCurrentReplicaSetSize       ResponseHeader = 0x001F
	ResponseHeaderCollectionLazyIndexProgress ResponseHeader = 0x0020
	ResponseHeaderPartitionKeyRangeId         ResponseHeader = 0x0021
	// 0x0022-0x0024 not used (RequestMessageReceivedTime, ResponseMessageSentTime, ResponseMessageTimeOffset)
	ResponseHeaderLogResults                   ResponseHeader = 0x0025
	ResponseHeaderXPRole                       ResponseHeader = 0x0026
	ResponseHeaderIsRUPerMinuteUsed            ResponseHeader = 0x0027
	ResponseHeaderQueryMetrics                 ResponseHeader = 0x0028
	ResponseHeaderGlobalCommittedLSN           ResponseHeader = 0x0029
	ResponseHeaderNumberOfReadRegions          ResponseHeader = 0x0030
	ResponseHeaderOfferReplacePending          ResponseHeader = 0x0031
	ResponseHeaderItemLSN                      ResponseHeader = 0x0032
	ResponseHeaderRestoreState                 ResponseHeader = 0x0033
	ResponseHeaderCollectionSecurityIdentifier ResponseHeader = 0x0034
	ResponseHeaderTransportRequestID           ResponseHeader = 0x0035
	ResponseHeaderShareThroughput              ResponseHeader = 0x0036
	// 0x0037 not used (LeaseSeqNumber)
	ResponseHeaderDisableRntbdChannel ResponseHeader = 0x0038
	ResponseHeaderServerDateTimeUtc   ResponseHeader = 0x0039
	ResponseHeaderLocalLSN            ResponseHeader = 0x003A
	ResponseHeaderQuorumAckedLocalLSN ResponseHeader = 0x003B
	ResponseHeaderItemLocalLSN        ResponseHeader = 0x003C
	ResponseHeaderHasTentativeWrites  ResponseHeader = 0x003D
	ResponseHeaderSessionToken        ResponseHeader = 0x003E
)

// ResponseHeaderInfo provides metadata about response headers.
type ResponseHeaderInfo struct {
	ID         ResponseHeader
	Type       TokenType
	IsRequired bool
}

// ResponseHeaders maps header IDs to their metadata.
var ResponseHeaders = map[ResponseHeader]ResponseHeaderInfo{
	ResponseHeaderPayloadPresent:               {ResponseHeaderPayloadPresent, TokenByte, true},
	ResponseHeaderLastStateChangeDateTime:      {ResponseHeaderLastStateChangeDateTime, TokenSmallString, false},
	ResponseHeaderContinuationToken:            {ResponseHeaderContinuationToken, TokenString, false},
	ResponseHeaderETag:                         {ResponseHeaderETag, TokenString, false},
	ResponseHeaderReadsPerformed:               {ResponseHeaderReadsPerformed, TokenULong, false},
	ResponseHeaderWritesPerformed:              {ResponseHeaderWritesPerformed, TokenULong, false},
	ResponseHeaderQueriesPerformed:             {ResponseHeaderQueriesPerformed, TokenULong, false},
	ResponseHeaderIndexTermsGenerated:          {ResponseHeaderIndexTermsGenerated, TokenULong, false},
	ResponseHeaderScriptsExecuted:              {ResponseHeaderScriptsExecuted, TokenULong, false},
	ResponseHeaderRetryAfterMilliseconds:       {ResponseHeaderRetryAfterMilliseconds, TokenULong, false},
	ResponseHeaderIndexingDirective:            {ResponseHeaderIndexingDirective, TokenByte, false},
	ResponseHeaderStorageMaxResoureQuota:       {ResponseHeaderStorageMaxResoureQuota, TokenString, false},
	ResponseHeaderStorageResourceQuotaUsage:    {ResponseHeaderStorageResourceQuotaUsage, TokenString, false},
	ResponseHeaderSchemaVersion:                {ResponseHeaderSchemaVersion, TokenSmallString, false},
	ResponseHeaderCollectionPartitionIndex:     {ResponseHeaderCollectionPartitionIndex, TokenULong, false},
	ResponseHeaderCollectionServiceIndex:       {ResponseHeaderCollectionServiceIndex, TokenULong, false},
	ResponseHeaderLSN:                          {ResponseHeaderLSN, TokenLongLong, false},
	ResponseHeaderItemCount:                    {ResponseHeaderItemCount, TokenULong, false},
	ResponseHeaderRequestCharge:                {ResponseHeaderRequestCharge, TokenDouble, false},
	ResponseHeaderOwnerFullName:                {ResponseHeaderOwnerFullName, TokenString, false},
	ResponseHeaderOwnerId:                      {ResponseHeaderOwnerId, TokenString, false},
	ResponseHeaderDatabaseAccountId:            {ResponseHeaderDatabaseAccountId, TokenString, false},
	ResponseHeaderQuorumAckedLSN:               {ResponseHeaderQuorumAckedLSN, TokenLongLong, false},
	ResponseHeaderRequestValidationFailure:     {ResponseHeaderRequestValidationFailure, TokenByte, false},
	ResponseHeaderSubStatus:                    {ResponseHeaderSubStatus, TokenULong, false},
	ResponseHeaderCollectionUpdateProgress:     {ResponseHeaderCollectionUpdateProgress, TokenULong, false},
	ResponseHeaderCurrentWriteQuorum:           {ResponseHeaderCurrentWriteQuorum, TokenULong, false},
	ResponseHeaderCurrentReplicaSetSize:        {ResponseHeaderCurrentReplicaSetSize, TokenULong, false},
	ResponseHeaderCollectionLazyIndexProgress:  {ResponseHeaderCollectionLazyIndexProgress, TokenULong, false},
	ResponseHeaderPartitionKeyRangeId:          {ResponseHeaderPartitionKeyRangeId, TokenString, false},
	ResponseHeaderLogResults:                   {ResponseHeaderLogResults, TokenString, false},
	ResponseHeaderXPRole:                       {ResponseHeaderXPRole, TokenULong, false},
	ResponseHeaderIsRUPerMinuteUsed:            {ResponseHeaderIsRUPerMinuteUsed, TokenByte, false},
	ResponseHeaderQueryMetrics:                 {ResponseHeaderQueryMetrics, TokenString, false},
	ResponseHeaderGlobalCommittedLSN:           {ResponseHeaderGlobalCommittedLSN, TokenLongLong, false},
	ResponseHeaderNumberOfReadRegions:          {ResponseHeaderNumberOfReadRegions, TokenULong, false},
	ResponseHeaderOfferReplacePending:          {ResponseHeaderOfferReplacePending, TokenByte, false},
	ResponseHeaderItemLSN:                      {ResponseHeaderItemLSN, TokenLongLong, false},
	ResponseHeaderRestoreState:                 {ResponseHeaderRestoreState, TokenString, false},
	ResponseHeaderCollectionSecurityIdentifier: {ResponseHeaderCollectionSecurityIdentifier, TokenString, false},
	ResponseHeaderTransportRequestID:           {ResponseHeaderTransportRequestID, TokenULong, false},
	ResponseHeaderShareThroughput:              {ResponseHeaderShareThroughput, TokenByte, false},
	ResponseHeaderDisableRntbdChannel:          {ResponseHeaderDisableRntbdChannel, TokenByte, false},
	ResponseHeaderServerDateTimeUtc:            {ResponseHeaderServerDateTimeUtc, TokenSmallString, false},
	ResponseHeaderLocalLSN:                     {ResponseHeaderLocalLSN, TokenLongLong, false},
	ResponseHeaderQuorumAckedLocalLSN:          {ResponseHeaderQuorumAckedLocalLSN, TokenLongLong, false},
	ResponseHeaderItemLocalLSN:                 {ResponseHeaderItemLocalLSN, TokenLongLong, false},
	ResponseHeaderHasTentativeWrites:           {ResponseHeaderHasTentativeWrites, TokenByte, false},
	ResponseHeaderSessionToken:                 {ResponseHeaderSessionToken, TokenString, false},
}
