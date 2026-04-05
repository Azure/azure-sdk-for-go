// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	HTTPHeaderAuthorization                             = "authorization"
	HTTPHeaderContentType                               = "Content-Type"
	HTTPHeaderDate                                      = "date"
	HTTPHeaderXDate                                     = "x-ms-date"
	HTTPHeaderIfMatch                                   = "If-Match"
	HTTPHeaderIfNoneMatch                               = "If-None-Match"
	HTTPHeaderIfModifiedSince                           = "If-Modified-Since"
	HTTPHeaderActivityID                                = "x-ms-activity-id"
	HTTPHeaderConsistencyLevel                          = "x-ms-consistency-level"
	HTTPHeaderContinuation                              = "x-ms-continuation"
	HTTPHeaderPageSize                                  = "x-ms-max-item-count"
	HTTPHeaderPartitionKey                              = "x-ms-documentdb-partitionkey"
	HTTPHeaderPartitionKeyRangeID                       = "x-ms-documentdb-partitionkeyrangeid"
	HTTPHeaderIndexingDirective                         = "x-ms-indexing-directive"
	HTTPHeaderPreTriggerInclude                         = "x-ms-documentdb-pre-trigger-include"
	HTTPHeaderPostTriggerInclude                        = "x-ms-documentdb-post-trigger-include"
	HTTPHeaderPreTriggerExclude                         = "x-ms-documentdb-pre-trigger-exclude"
	HTTPHeaderPostTriggerExclude                        = "x-ms-documentdb-post-trigger-exclude"
	HTTPHeaderEnableScanInQuery                         = "x-ms-documentdb-query-enable-scan"
	HTTPHeaderEmitVerboseTracesInQuery                  = "x-ms-documentdb-query-emit-traces"
	HTTPHeaderEnableLowPrecisionOrderBy                 = "x-ms-documentdb-query-enable-low-precision-order-by"
	HTTPHeaderEnableLogging                             = "x-ms-documentdb-script-enable-logging"
	HTTPHeaderA_IM                                      = "A-IM"
	HTTPHeaderPopulateQuotaInfo                         = "x-ms-documentdb-populatequotainfo"
	HTTPHeaderDisableRUPerMinuteUsage                   = "x-ms-documentdb-disable-ru-per-minute-usage"
	HTTPHeaderPopulateQueryMetrics                      = "x-ms-documentdb-populatequerymetrics"
	HTTPHeaderResponseContinuationTokenLimitInKb        = "x-ms-documentdb-responsecontinuationtokenlimitinkb"
	HTTPHeaderPopulatePartitionStatistics               = "x-ms-documentdb-populatepartitionstatistics"
	HTTPHeaderPopulateCollectionThroughputInfo          = "x-ms-documentdb-populatecollectionthroughputinfo"
	HTTPHeaderRemainingTimeInMs                         = "x-ms-remaining-time-in-ms-on-client"
	HTTPHeaderClientRetryAttemptCount                   = "x-ms-client-retry-attempt-count"
	HTTPHeaderTargetLSN                                 = "x-ms-target-lsn"
	HTTPHeaderTargetGlobalCommittedLSN                  = "x-ms-target-global-committed-lsn"
	HTTPHeaderTransportRequestID                        = "x-ms-transport-request-id"
	HTTPHeaderResourceTokenExpiry                       = "x-ms-documentdb-expiry-seconds"
	HTTPHeaderFilterBySchemaRid                         = "x-ms-documentdb-filterby-schema-rid"
	HTTPHeaderGatewaySignature                          = "x-ms-gateway-signature"
	HTTPHeaderCollectionRemoteStorageSecurityIdentifier = "x-ms-collection-security-identifier"
	HTTPHeaderVersion                                   = "x-ms-version"
	HTTPHeaderEnumerationDirection                      = "x-ms-enumeration-direction"
	HTTPHeaderContentSerializationFormat                = "x-ms-content-serialization-format"
	HTTPHeaderCanCharge                                 = "x-ms-can-charge"
	HTTPHeaderCanThrottle                               = "x-ms-can-throttle"
	HTTPHeaderProfileRequest                            = "x-ms-profile-request"
	HTTPHeaderForceQueryScan                            = "x-ms-documentdb-force-query-scan"
	HTTPHeaderSupportSpatialLegacyCoordinates           = "x-ms-documentdb-supportspatiallegacycoordinates"
	HTTPHeaderUsePolygonsSmallerThanAHemisphere         = "x-ms-documentdb-usepolygonssmallerthanahemisphere"
	HTTPHeaderCanOfferReplaceComplete                   = "x-ms-can-offer-replace-complete"
	HTTPHeaderIsReadOnlyScript                          = "x-ms-is-readonly-script"
	HTTPHeaderIsAutoScaleRequest                        = "x-ms-is-auto-scale"
	HTTPHeaderMigrateCollectionDirective                = "x-ms-migratecollectiondirective"
	HTTPHeaderSharedOfferThroughput                     = "x-ms-cosmos-shared-offer-throughput"
	HTTPHeaderReadFeedKeyType                           = "x-ms-read-feed-key-type"
	HTTPHeaderStartID                                   = "x-ms-start-id"
	HTTPHeaderEndID                                     = "x-ms-end-id"
	HTTPHeaderStartEPK                                  = "x-ms-start-epk"
	HTTPHeaderEndEPK                                    = "x-ms-end-epk"
	HTTPHeaderRestoreMetadataFilter                     = "x-ms-restore-metadata-filter"
)

const (
	BackendHeaderBinaryID                        = "x-ms-binary-id"
	BackendHeaderEffectivePartitionKey           = "x-ms-effective-partition-key"
	BackendHeaderBindReplicaDirective            = "x-ms-bind-replica-directive"
	BackendHeaderPrimaryMasterKey                = "x-ms-primary-master-key"
	BackendHeaderSecondaryMasterKey              = "x-ms-secondary-master-key"
	BackendHeaderPrimaryReadonlyKey              = "x-ms-primary-readonly-key"
	BackendHeaderSecondaryReadonlyKey            = "x-ms-secondary-readonly-key"
	BackendHeaderEntityID                        = "x-ms-entity-id"
	BackendHeaderResourceSchemaName              = "x-ms-resource-schema-name"
	BackendHeaderIsFanoutRequest                 = "x-ms-is-fanout-request"
	BackendHeaderCollectionPartitionIndex        = "x-ms-collection-partition-index"
	BackendHeaderCollectionServiceIndex          = "x-ms-collection-service-index"
	BackendHeaderCollectionRid                   = "x-ms-collection-rid"
	BackendHeaderPartitionCount                  = "x-ms-partition-count"
	BackendHeaderPartitionResourceFilter         = "x-ms-partition-resource-filter"
	BackendHeaderEnableDynamicRidRangeAllocation = "x-ms-enable-dynamic-rid-range-allocation"
	BackendHeaderExcludeSystemProperties         = "x-ms-exclude-system-properties"
	BackendHeaderBinaryPassthroughRequest        = "x-ms-binary-passthrough-request"
	BackendHeaderTimeToLiveInSeconds             = "x-ms-time-to-live-in-seconds"
	BackendHeaderRemoteStorageType               = "x-ms-remote-storage-type"
	BackendHeaderShareThroughput                 = "x-ms-share-throughput"
	BackendHeaderFanoutOperationState            = "x-ms-fanout-operation-state"
	BackendHeaderRestoreParams                   = "x-ms-restore-params"
	BackendHeaderIsUserRequest                   = "x-ms-is-user-request"
	BackendHeaderAllowTentativeWrites            = "x-ms-allow-tentative-writes"
)

const (
	RespHeaderRequestCharge             = "x-ms-request-charge"
	RespHeaderSessionToken              = "x-ms-session-token"
	RespHeaderContinuation              = "x-ms-continuation"
	RespHeaderETag                      = "etag"
	RespHeaderActivityID                = "x-ms-activity-id"
	RespHeaderRetryAfterMs              = "x-ms-retry-after-ms"
	RespHeaderSubStatus                 = "x-ms-substatus"
	RespHeaderItemCount                 = "x-ms-item-count"
	RespHeaderLSN                       = "lsn"
	RespHeaderGlobalCommittedLSN        = "x-ms-global-committed-lsn"
	RespHeaderQuorumAckedLSN            = "x-ms-quorum-acked-lsn"
	RespHeaderCurrentWriteQuorum        = "x-ms-current-write-quorum"
	RespHeaderCurrentReplicaSetSize     = "x-ms-current-replica-set-size"
	RespHeaderOwnerID                   = "x-ms-owner-id"
	RespHeaderOwnerFullName             = "x-ms-owner-full-name"
	RespHeaderQueryMetrics              = "x-ms-documentdb-query-metrics"
	RespHeaderPartitionKeyRangeID       = "x-ms-documentdb-partitionkeyrangeid"
	RespHeaderIsRUPerMinuteUsed         = "x-ms-documentdb-is-ru-per-minute-used"
	RespHeaderTransportRequestID        = "x-ms-transport-request-id"
	RespHeaderServerDateTimeUtc         = "x-ms-server-date-time-utc"
	RespHeaderLastStateChangeDateTime   = "x-ms-last-state-change-utc"
	RespHeaderSchemaVersion             = "x-ms-schemaversion"
	RespHeaderIndexingDirective         = "x-ms-indexing-directive"
	RespHeaderStorageMaxResoureQuota    = "x-ms-resource-quota"
	RespHeaderStorageResourceQuotaUsage = "x-ms-resource-usage"
	RespHeaderItemLSN                   = "x-ms-item-lsn"
	RespHeaderLocalLSN                  = "x-ms-local-lsn"
	RespHeaderQuorumAckedLocalLSN       = "x-ms-quorum-acked-local-lsn"
	RespHeaderItemLocalLSN              = "x-ms-item-local-lsn"
	RespHeaderHasTentativeWrites        = "x-ms-cosmosdb-has-tentative-writes"
)

type ServiceRequest struct {
	OperationType             OperationType
	ResourceType              ResourceType
	ResourceID                string
	ResourceAddress           string
	IsNameBased               bool
	ActivityID                uuid.UUID
	Headers                   map[string]string
	Content                   []byte
	Continuation              string
	ReplicaPath               string
	TransportRequestID        uint32
	PartitionKeyRangeIdentity *PartitionKeyRangeIdentity
	DefaultReplicaIndex       *int
}

func NewServiceRequest(
	opType OperationType,
	resType ResourceType,
	resourceAddress string,
	headers map[string]string,
) *ServiceRequest {
	if headers == nil {
		headers = make(map[string]string)
	}
	return &ServiceRequest{
		OperationType:   opType,
		ResourceType:    resType,
		ResourceAddress: resourceAddress,
		Headers:         headers,
		ActivityID:      uuid.New(),
		IsNameBased:     strings.HasPrefix(resourceAddress, "/dbs/") || strings.HasPrefix(resourceAddress, "dbs/"),
	}
}

func (r *ServiceRequest) GetHeader(name string) string {
	if r.Headers == nil {
		return ""
	}
	return r.Headers[name]
}

func (r *ServiceRequest) SetHeader(name, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[name] = value
}

func (r *ServiceRequest) GetContinuation() string {
	if r.Continuation != "" {
		return r.Continuation
	}
	return r.GetHeader(HTTPHeaderContinuation)
}

func (r *ServiceRequest) HasContent() bool {
	return len(r.Content) > 0
}

func (r *ServiceRequest) GetDefaultReplicaIndex() *int {
	return r.DefaultReplicaIndex
}

type StoreResponse struct {
	StatusCode int
	ActivityID uuid.UUID
	Headers    map[string]string
	Content    []byte
	Endpoint   string
}

func NewStoreResponse(statusCode int, activityID uuid.UUID, endpoint string) *StoreResponse {
	return &StoreResponse{
		StatusCode: statusCode,
		ActivityID: activityID,
		Headers:    make(map[string]string),
		Endpoint:   endpoint,
	}
}

func (r *StoreResponse) GetHeader(name string) string {
	if r.Headers == nil {
		return ""
	}
	return r.Headers[strings.ToLower(name)]
}

func (r *StoreResponse) SetHeader(name, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[strings.ToLower(name)] = value
}

func (r *StoreResponse) GetRequestCharge() float64 {
	val := r.GetHeader(RespHeaderRequestCharge)
	if val == "" {
		return 0
	}
	charge, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0
	}
	return charge
}

func (r *StoreResponse) GetSessionTokenString() string {
	return r.GetHeader(RespHeaderSessionToken)
}

func (r *StoreResponse) GetContinuation() string {
	return r.GetHeader(RespHeaderContinuation)
}

func (r *StoreResponse) GetETag() string {
	return r.GetHeader(RespHeaderETag)
}

func (r *StoreResponse) GetLSN() int64 {
	val := r.GetHeader(RespHeaderLSN)
	if val == "" {
		return -1
	}
	lsn, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return -1
	}
	return lsn
}

func (r *StoreResponse) GetItemCount() int {
	val := r.GetHeader(RespHeaderItemCount)
	if val == "" {
		return 0
	}
	count, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return count
}

func (r *StoreResponse) GetSubStatusCode() int {
	val := r.GetHeader(RespHeaderSubStatus)
	if val == "" {
		return 0
	}
	code, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return code
}

func (r *StoreResponse) GetRetryAfterMs() int64 {
	val := r.GetHeader(RespHeaderRetryAfterMs)
	if val == "" {
		return 0
	}
	ms, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0
	}
	return ms
}

func (r *StoreResponse) IsSuccessStatusCode() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

func ParseResourcePath(path string) map[string]string {
	result := make(map[string]string)
	path = strings.Trim(path, "/")
	if path == "" {
		return result
	}

	fragments := strings.Split(path, "/")
	count := len(fragments)

	for i := 0; i < count-1; i += 2 {
		segmentType := strings.ToLower(fragments[i])
		segmentName := fragments[i+1]

		switch segmentType {
		case "dbs":
			result["database"] = segmentName
		case "colls":
			result["collection"] = segmentName
		case "docs":
			result["document"] = segmentName
		case "sprocs":
			result["storedProcedure"] = segmentName
		case "triggers":
			result["trigger"] = segmentName
		case "udfs":
			result["userDefinedFunction"] = segmentName
		case "users":
			result["user"] = segmentName
		case "permissions":
			result["permission"] = segmentName
		case "conflicts":
			result["conflict"] = segmentName
		case "attachments":
			result["attachment"] = segmentName
		case "pkranges":
			result["partitionKeyRange"] = segmentName
		case "schemas":
			result["schema"] = segmentName
		case "udts":
			result["userDefinedType"] = segmentName
		}
	}

	return result
}

func DecodeBase64(s string) []byte {
	if s == "" {
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		data, err = base64.URLEncoding.DecodeString(s)
		if err != nil {
			return nil
		}
	}
	return data
}

func ParseBool(s string) bool {
	if s == "" {
		return false
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return b
}

func ParseInt64(s string) (int64, bool) {
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func ParseUint32(s string) (uint32, bool) {
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint32(v), true
}
