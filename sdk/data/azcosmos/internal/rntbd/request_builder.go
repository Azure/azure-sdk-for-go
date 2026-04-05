// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"encoding/hex"
	"strconv"
	"strings"
)

var httpHeaderToRntbdHeader = map[string]RequestHeader{
	HTTPHeaderAuthorization:                             RequestHeaderAuthorizationToken,
	HTTPHeaderXDate:                                     RequestHeaderDate,
	HTTPHeaderConsistencyLevel:                          RequestHeaderConsistencyLevel,
	HTTPHeaderContinuation:                              RequestHeaderContinuationToken,
	HTTPHeaderSessionToken:                              RequestHeaderSessionToken,
	HTTPHeaderPageSize:                                  RequestHeaderPageSize,
	HTTPHeaderPartitionKey:                              RequestHeaderPartitionKey,
	HTTPHeaderPartitionKeyRangeID:                       RequestHeaderPartitionKeyRangeId,
	"x-ms-cosmos-collection-rid":                        RequestHeaderCollectionRid,
	"x-ms-documentdb-collection-rid":                    RequestHeaderCollectionRid,
	HTTPHeaderIndexingDirective:                         RequestHeaderIndexingDirective,
	HTTPHeaderPreTriggerInclude:                         RequestHeaderPreTriggerInclude,
	HTTPHeaderPostTriggerInclude:                        RequestHeaderPostTriggerInclude,
	HTTPHeaderPreTriggerExclude:                         RequestHeaderPreTriggerExclude,
	HTTPHeaderPostTriggerExclude:                        RequestHeaderPostTriggerExclude,
	HTTPHeaderEnableScanInQuery:                         RequestHeaderEnableScanInQuery,
	HTTPHeaderEmitVerboseTracesInQuery:                  RequestHeaderEmitVerboseTracesInQuery,
	HTTPHeaderEnableLowPrecisionOrderBy:                 RequestHeaderEnableLowPrecisionOrderBy,
	HTTPHeaderEnableLogging:                             RequestHeaderEnableLogging,
	HTTPHeaderA_IM:                                      RequestHeaderA_IM,
	HTTPHeaderPopulateQuotaInfo:                         RequestHeaderPopulateQuotaInfo,
	HTTPHeaderDisableRUPerMinuteUsage:                   RequestHeaderDisableRUPerMinuteUsage,
	HTTPHeaderPopulateQueryMetrics:                      RequestHeaderPopulateQueryMetrics,
	HTTPHeaderResponseContinuationTokenLimitInKb:        RequestHeaderResponseContinuationTokenLimitInKb,
	HTTPHeaderPopulatePartitionStatistics:               RequestHeaderPopulatePartitionStatistics,
	HTTPHeaderPopulateCollectionThroughputInfo:          RequestHeaderPopulateCollectionThroughputInfo,
	HTTPHeaderRemainingTimeInMs:                         RequestHeaderRemainingTimeInMsOnClientRequest,
	HTTPHeaderClientRetryAttemptCount:                   RequestHeaderClientRetryAttemptCount,
	HTTPHeaderTargetLSN:                                 RequestHeaderTargetLsn,
	HTTPHeaderTargetGlobalCommittedLSN:                  RequestHeaderTargetGlobalCommittedLsn,
	HTTPHeaderTransportRequestID:                        RequestHeaderTransportRequestID,
	HTTPHeaderResourceTokenExpiry:                       RequestHeaderResourceTokenExpiry,
	HTTPHeaderFilterBySchemaRid:                         RequestHeaderFilterBySchemaRid,
	HTTPHeaderGatewaySignature:                          RequestHeaderGatewaySignature,
	HTTPHeaderCollectionRemoteStorageSecurityIdentifier: RequestHeaderCollectionRemoteStorageSecurityIdentifier,
	HTTPHeaderEnumerationDirection:                      RequestHeaderEnumerationDirection,
	HTTPHeaderContentSerializationFormat:                RequestHeaderContentSerializationFormat,
	HTTPHeaderCanCharge:                                 RequestHeaderCanCharge,
	HTTPHeaderCanThrottle:                               RequestHeaderCanThrottle,
	HTTPHeaderProfileRequest:                            RequestHeaderProfileRequest,
	HTTPHeaderForceQueryScan:                            RequestHeaderForceQueryScan,
	HTTPHeaderSupportSpatialLegacyCoordinates:           RequestHeaderSupportSpatialLegacyCoordinates,
	HTTPHeaderUsePolygonsSmallerThanAHemisphere:         RequestHeaderUsePolygonsSmallerThanAHemisphere,
	HTTPHeaderCanOfferReplaceComplete:                   RequestHeaderCanOfferReplaceComplete,
	HTTPHeaderIsReadOnlyScript:                          RequestHeaderIsReadOnlyScript,
	HTTPHeaderIsAutoScaleRequest:                        RequestHeaderIsAutoScaleRequest,
	HTTPHeaderMigrateCollectionDirective:                RequestHeaderMigrateCollectionDirective,
	HTTPHeaderSharedOfferThroughput:                     RequestHeaderSharedOfferThroughput,
	HTTPHeaderReadFeedKeyType:                           RequestHeaderReadFeedKeyType,
	HTTPHeaderStartID:                                   RequestHeaderStartId,
	HTTPHeaderEndID:                                     RequestHeaderEndId,
	HTTPHeaderStartEPK:                                  RequestHeaderStartEpk,
	HTTPHeaderEndEPK:                                    RequestHeaderEndEpk,
	HTTPHeaderRestoreMetadataFilter:                     RequestHeaderRestoreMetadaFilter,
	HTTPHeaderIfMatch:                                   RequestHeaderMatch,
	HTTPHeaderIfModifiedSince:                           RequestHeaderIfModifiedSince,
	BackendHeaderBinaryID:                               RequestHeaderBinaryId,
	// EffectivePartitionKey (0x005A) - Java SDK's fillTokenFromHeader doesn't handle Bytes type, so EPK is NOT sent
	BackendHeaderBindReplicaDirective:            RequestHeaderBindReplicaDirective,
	BackendHeaderPrimaryMasterKey:                RequestHeaderPrimaryMasterKey,
	BackendHeaderSecondaryMasterKey:              RequestHeaderSecondaryMasterKey,
	BackendHeaderPrimaryReadonlyKey:              RequestHeaderPrimaryReadonlyKey,
	BackendHeaderSecondaryReadonlyKey:            RequestHeaderSecondaryReadonlyKey,
	BackendHeaderEntityID:                        RequestHeaderEntityId,
	BackendHeaderResourceSchemaName:              RequestHeaderResourceSchemaName,
	BackendHeaderIsFanoutRequest:                 RequestHeaderIsFanout,
	BackendHeaderCollectionPartitionIndex:        RequestHeaderCollectionPartitionIndex,
	BackendHeaderCollectionServiceIndex:          RequestHeaderCollectionServiceIndex,
	BackendHeaderCollectionRid:                   RequestHeaderCollectionRid,
	BackendHeaderPartitionCount:                  RequestHeaderPartitionCount,
	BackendHeaderPartitionResourceFilter:         RequestHeaderPartitionResourceFilter,
	BackendHeaderEnableDynamicRidRangeAllocation: RequestHeaderEnableDynamicRidRangeAllocation,
	BackendHeaderExcludeSystemProperties:         RequestHeaderExcludeSystemProperties,
	BackendHeaderBinaryPassthroughRequest:        RequestHeaderBinaryPassthroughRequest,
	BackendHeaderTimeToLiveInSeconds:             RequestHeaderTimeToLiveInSeconds,
	BackendHeaderRemoteStorageType:               RequestHeaderRemoteStorageType,
	BackendHeaderShareThroughput:                 RequestHeaderShareThroughput,
	BackendHeaderFanoutOperationState:            RequestHeaderFanoutOperationState,
	BackendHeaderRestoreParams:                   RequestHeaderRestoreParams,
	BackendHeaderIsUserRequest:                   RequestHeaderIsUserRequest,
	BackendHeaderAllowTentativeWrites:            RequestHeaderAllowTentativeWrites,
	// Also map the "cosmos" variant of the header used by global endpoint manager
	"x-ms-cosmos-allow-tentative-writes": RequestHeaderAllowTentativeWrites,
	HTTPHeaderPrefer:                     RequestHeaderReturnPreference,
	HTTPHeaderSDKSupportedCapabilities:   RequestHeaderSDKSupportedCapabilities,
	HTTPHeaderVersion:                    RequestHeaderClientVersion,
}

var consistencyLevelMap = map[string]ConsistencyLevel{
	"strong":           ConsistencyStrong,
	"bounded":          ConsistencyBoundedStaleness,
	"boundedstaleness": ConsistencyBoundedStaleness,
	"session":          ConsistencySession,
	"eventual":         ConsistencyEventual,
	"consistentprefix": ConsistencyConsistentPrefix,
}

var indexingDirectiveMap = map[string]IndexingDirective{
	"default": IndexingDirectiveDefault,
	"include": IndexingDirectiveInclude,
	"exclude": IndexingDirectiveExclude,
}

var contentSerializationFormatMap = map[string]ContentSerializationFormat{
	"jsontext":     ContentSerializationJsonText,
	"cosmosbinary": ContentSerializationCosmosBinary,
}

var enumerationDirectionMap = map[string]EnumerationDirection{
	"forward": EnumerationDirectionForward,
	"reverse": EnumerationDirectionReverse,
}

var readFeedKeyTypeMap = map[string]ReadFeedKeyType{
	"resourceid":            ReadFeedKeyTypeResourceId,
	"effectivepartitionkey": ReadFeedKeyTypeEffectivePartitionKey,
}

var migrateCollectionDirectiveMap = map[string]MigrateCollectionDirective{
	"thaw":   MigrateCollectionThaw,
	"freeze": MigrateCollectionFreeze,
}

var fanoutOperationStateMap = map[string]FanoutOperationState{
	"started":   FanoutOperationStarted,
	"completed": FanoutOperationCompleted,
}

var remoteStorageTypeMap = map[string]RemoteStorageType{
	"notspecified": RemoteStorageTypeNotSpecified,
	"standard":     RemoteStorageTypeStandard,
	"premium":      RemoteStorageTypePremium,
}

func BuildRequestMessage(req *ServiceRequest) (*RequestMessage, error) {
	msg := NewRequestMessage(req.ResourceType, req.OperationType, req.ActivityID)

	hasPayload := req.HasContent()
	msg.Headers.SetByte(uint16(RequestHeaderPayloadPresent), boolToByte(hasPayload))

	// ContentSerializationFormat: Java SDK only sets this token when the HTTP header
	// "x-ms-documentdb-content-serialization-format" is explicitly present (see RntbdRequestHeaders.addContentSerializationFormat).
	// Do NOT set it unconditionally - setting it when not expected can cause server-side deserialization issues.

	if req.ReplicaPath != "" {
		msg.Headers.SetString(uint16(RequestHeaderReplicaPath), req.ReplicaPath)
	}

	if req.TransportRequestID != 0 {
		msg.Headers.SetULong(uint16(RequestHeaderTransportRequestID), req.TransportRequestID)
	}

	// Java SDK sets ResourceId BEFORE checking if name-based (see addResourceIdOrPathHeaders lines 1169-1177).
	// "Name-based can also have ResourceId because gateway might have generated it"
	// This means ResourceId is sent regardless of whether request is name-based or not.
	if req.ResourceID != "" {
		resourceIDBytes := DecodeBase64(req.ResourceID)
		if resourceIDBytes == nil {
			resourceIDBytes = []byte(req.ResourceID)
		}
		msg.Headers.SetBytes(uint16(RequestHeaderResourceId), resourceIDBytes)
	}

	// THEN add name-based path headers if this is a name-based request
	if req.IsNameBased {
		addNameBasedHeaders(msg.Headers, req.ResourceAddress, req.ResourceType)
	}

	for httpHeader, value := range req.Headers {
		if value == "" {
			continue
		}
		addHeader(msg.Headers, httpHeader, value)
	}

	if req.HasContent() {
		msg.Payload = req.Content
	}

	return msg, nil
}

func addNameBasedHeaders(headers *TokenStream, resourceAddress string, resourceType ResourceType) {
	components := ParseResourcePath(resourceAddress)

	if db, ok := components["database"]; ok {
		headers.SetString(uint16(RequestHeaderDatabaseName), db)
	}
	if coll, ok := components["collection"]; ok {
		headers.SetString(uint16(RequestHeaderCollectionName), coll)
	}
	if doc, ok := components["document"]; ok {
		headers.SetString(uint16(RequestHeaderDocumentName), doc)
	}
	if sproc, ok := components["storedProcedure"]; ok {
		headers.SetString(uint16(RequestHeaderStoredProcedureName), sproc)
	}
	if trigger, ok := components["trigger"]; ok {
		headers.SetString(uint16(RequestHeaderTriggerName), trigger)
	}
	if udf, ok := components["userDefinedFunction"]; ok {
		headers.SetString(uint16(RequestHeaderUserDefinedFunctionName), udf)
	}
	if user, ok := components["user"]; ok {
		headers.SetString(uint16(RequestHeaderUserName), user)
	}
	if perm, ok := components["permission"]; ok {
		headers.SetString(uint16(RequestHeaderPermissionName), perm)
	}
	if conflict, ok := components["conflict"]; ok {
		headers.SetString(uint16(RequestHeaderConflictName), conflict)
	}
	if attachment, ok := components["attachment"]; ok {
		headers.SetString(uint16(RequestHeaderAttachmentName), attachment)
	}
	if pkRange, ok := components["partitionKeyRange"]; ok {
		headers.SetString(uint16(RequestHeaderPartitionKeyRangeName), pkRange)
	}
	if schema, ok := components["schema"]; ok {
		headers.SetString(uint16(RequestHeaderSchemaName), schema)
	}
	if udt, ok := components["userDefinedType"]; ok {
		headers.SetString(uint16(RequestHeaderUserDefinedTypeName), udt)
	}
}

func addHeader(headers *TokenStream, httpHeader string, value string) {
	rntbdHeader, ok := httpHeaderToRntbdHeader[httpHeader]
	if !ok {
		return
	}

	headerInfo, ok := RequestHeaders[rntbdHeader]
	if !ok {
		return
	}

	switch headerInfo.Type {
	case TokenByte:
		headers.SetByte(uint16(rntbdHeader), convertToByte(httpHeader, value))
	case TokenUShort:
		if v, err := strconv.ParseUint(value, 10, 16); err == nil {
			headers.SetUShort(uint16(rntbdHeader), uint16(v))
		}
	case TokenULong:
		if v, err := strconv.ParseUint(value, 10, 32); err == nil {
			headers.SetULong(uint16(rntbdHeader), uint32(v))
		}
	case TokenLong:
		if v, err := strconv.ParseInt(value, 10, 32); err == nil {
			headers.SetLong(uint16(rntbdHeader), int32(v))
		}
	case TokenULongLong:
		if v, err := strconv.ParseUint(value, 10, 64); err == nil {
			headers.SetULongLong(uint16(rntbdHeader), v)
		}
	case TokenLongLong:
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			headers.SetLongLong(uint16(rntbdHeader), v)
		}
	case TokenString, TokenSmallString, TokenULongString:
		headers.SetValue(uint16(rntbdHeader), headerInfo.Type, value) //nolint:errcheck
	case TokenBytes, TokenSmallBytes, TokenULongBytes:
		// For EffectivePartitionKey, the value is a hex string that needs to be decoded
		if rntbdHeader == RequestHeaderEffectivePartitionKey {
			if decoded, err := hex.DecodeString(value); err == nil {
				headers.SetValue(uint16(rntbdHeader), headerInfo.Type, decoded) //nolint:errcheck
			}
		} else if decoded := DecodeBase64(value); decoded != nil {
			headers.SetValue(uint16(rntbdHeader), headerInfo.Type, decoded) //nolint:errcheck
		} else {
			headers.SetValue(uint16(rntbdHeader), headerInfo.Type, []byte(value)) //nolint:errcheck
		}
	case TokenDouble:
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			headers.SetDouble(uint16(rntbdHeader), v)
		}
	case TokenFloat:
		if v, err := strconv.ParseFloat(value, 32); err == nil {
			headers.SetFloat(uint16(rntbdHeader), float32(v))
		}
	}
}

func convertToByte(httpHeader string, value string) byte {
	valueLower := strings.ToLower(value)

	switch httpHeader {
	case HTTPHeaderConsistencyLevel:
		if level, ok := consistencyLevelMap[valueLower]; ok {
			return byte(level)
		}
	case HTTPHeaderIndexingDirective:
		if directive, ok := indexingDirectiveMap[valueLower]; ok {
			return byte(directive)
		}
	case HTTPHeaderContentSerializationFormat:
		if format, ok := contentSerializationFormatMap[valueLower]; ok {
			return byte(format)
		}
	case HTTPHeaderEnumerationDirection:
		if dir, ok := enumerationDirectionMap[valueLower]; ok {
			return byte(dir)
		}
	case HTTPHeaderReadFeedKeyType:
		if keyType, ok := readFeedKeyTypeMap[valueLower]; ok {
			return byte(keyType)
		}
	case HTTPHeaderMigrateCollectionDirective:
		if directive, ok := migrateCollectionDirectiveMap[valueLower]; ok {
			return byte(directive)
		}
	case BackendHeaderFanoutOperationState:
		if state, ok := fanoutOperationStateMap[valueLower]; ok {
			return byte(state)
		}
	case BackendHeaderRemoteStorageType:
		if storageType, ok := remoteStorageTypeMap[valueLower]; ok {
			return byte(storageType)
		}
	case HTTPHeaderPrefer:
		if strings.Contains(valueLower, "return=minimal") {
			return 1
		}
		return 0
	}

	if valueLower == "true" || value == "1" {
		return 1
	}
	return 0
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
