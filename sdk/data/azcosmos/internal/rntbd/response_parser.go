// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"fmt"
	"strconv"
)

var rntbdHeaderToHTTPHeader = map[ResponseHeader]string{
	ResponseHeaderPayloadPresent:            "",
	ResponseHeaderLastStateChangeDateTime:   RespHeaderLastStateChangeDateTime,
	ResponseHeaderContinuationToken:         RespHeaderContinuation,
	ResponseHeaderETag:                      RespHeaderETag,
	ResponseHeaderRetryAfterMilliseconds:    RespHeaderRetryAfterMs,
	ResponseHeaderStorageMaxResoureQuota:    RespHeaderStorageMaxResoureQuota,
	ResponseHeaderStorageResourceQuotaUsage: RespHeaderStorageResourceQuotaUsage,
	ResponseHeaderSchemaVersion:             RespHeaderSchemaVersion,
	ResponseHeaderLSN:                       RespHeaderLSN,
	ResponseHeaderItemCount:                 RespHeaderItemCount,
	ResponseHeaderRequestCharge:             RespHeaderRequestCharge,
	ResponseHeaderOwnerFullName:             RespHeaderOwnerFullName,
	ResponseHeaderOwnerId:                   RespHeaderOwnerID,
	ResponseHeaderQuorumAckedLSN:            RespHeaderQuorumAckedLSN,
	ResponseHeaderSubStatus:                 RespHeaderSubStatus,
	ResponseHeaderCurrentWriteQuorum:        RespHeaderCurrentWriteQuorum,
	ResponseHeaderCurrentReplicaSetSize:     RespHeaderCurrentReplicaSetSize,
	ResponseHeaderPartitionKeyRangeId:       RespHeaderPartitionKeyRangeID,
	ResponseHeaderIsRUPerMinuteUsed:         RespHeaderIsRUPerMinuteUsed,
	ResponseHeaderQueryMetrics:              RespHeaderQueryMetrics,
	ResponseHeaderGlobalCommittedLSN:        RespHeaderGlobalCommittedLSN,
	ResponseHeaderItemLSN:                   RespHeaderItemLSN,
	ResponseHeaderTransportRequestID:        RespHeaderTransportRequestID,
	ResponseHeaderServerDateTimeUtc:         RespHeaderServerDateTimeUtc,
	ResponseHeaderLocalLSN:                  RespHeaderLocalLSN,
	ResponseHeaderQuorumAckedLocalLSN:       RespHeaderQuorumAckedLocalLSN,
	ResponseHeaderItemLocalLSN:              RespHeaderItemLocalLSN,
	ResponseHeaderHasTentativeWrites:        RespHeaderHasTentativeWrites,
	ResponseHeaderSessionToken:              RespHeaderSessionToken,
	ResponseHeaderIndexingDirective:         RespHeaderIndexingDirective,
}

func ParseResponseMessage(msg *ResponseMessage, endpoint string) (*StoreResponse, error) {
	resp := NewStoreResponse(int(msg.Frame.Status), msg.Frame.ActivityID, endpoint)
	resp.Content = msg.Payload

	resp.SetHeader(RespHeaderActivityID, msg.Frame.ActivityID.String())

	for id, token := range msg.Headers.tokens {
		if !token.IsPresent() {
			continue
		}

		headerID := ResponseHeader(id)
		httpHeader, ok := rntbdHeaderToHTTPHeader[headerID]
		if !ok || httpHeader == "" {
			continue
		}

		value, err := token.GetValue()
		if err != nil {
			continue
		}

		strValue := convertTokenValueToString(headerID, value)
		if strValue != "" {
			resp.SetHeader(httpHeader, strValue)
		}
	}

	return resp, nil
}

func convertTokenValueToString(headerID ResponseHeader, value interface{}) string {
	if value == nil {
		return ""
	}

	headerInfo, ok := ResponseHeaders[headerID]
	if !ok {
		return fmt.Sprintf("%v", value)
	}

	switch headerInfo.Type {
	case TokenByte:
		if headerID == ResponseHeaderIsRUPerMinuteUsed ||
			headerID == ResponseHeaderHasTentativeWrites {
			if b, ok := value.(byte); ok {
				if b != 0 {
					return "true"
				}
				return "false"
			}
		}
		if b, ok := value.(byte); ok {
			return strconv.Itoa(int(b))
		}

	case TokenUShort:
		if v, ok := value.(uint16); ok {
			return strconv.FormatUint(uint64(v), 10)
		}

	case TokenULong:
		if v, ok := value.(int64); ok {
			return strconv.FormatInt(v, 10)
		}
		if v, ok := value.(uint32); ok {
			return strconv.FormatUint(uint64(v), 10)
		}

	case TokenLong:
		if v, ok := value.(int32); ok {
			return strconv.FormatInt(int64(v), 10)
		}

	case TokenLongLong, TokenULongLong:
		if v, ok := value.(int64); ok {
			return strconv.FormatInt(v, 10)
		}
		if v, ok := value.(uint64); ok {
			return strconv.FormatUint(v, 10)
		}

	case TokenDouble:
		if f, ok := value.(float64); ok {
			return strconv.FormatFloat(f, 'f', -1, 64)
		}

	case TokenFloat:
		if f, ok := value.(float32); ok {
			return strconv.FormatFloat(float64(f), 'f', -1, 32)
		}

	case TokenString, TokenSmallString, TokenULongString:
		if s, ok := value.(string); ok {
			return s
		}

	case TokenBytes, TokenSmallBytes, TokenULongBytes:
		if b, ok := value.([]byte); ok {
			return string(b)
		}
	}

	return fmt.Sprintf("%v", value)
}

func GetResponseHeaderString(msg *ResponseMessage, header ResponseHeader) string {
	token := msg.Headers.Get(uint16(header))
	if token == nil || !token.IsPresent() {
		return ""
	}

	value, err := token.GetValue()
	if err != nil {
		return ""
	}

	return convertTokenValueToString(header, value)
}

func GetResponseHeaderByte(msg *ResponseMessage, header ResponseHeader) byte {
	return msg.Headers.GetByte(uint16(header))
}

func GetResponseHeaderULong(msg *ResponseMessage, header ResponseHeader) uint32 {
	return msg.Headers.GetULong(uint16(header))
}

func GetResponseHeaderLongLong(msg *ResponseMessage, header ResponseHeader) int64 {
	return msg.Headers.GetLongLong(uint16(header))
}

func GetResponseHeaderDouble(msg *ResponseMessage, header ResponseHeader) float64 {
	return msg.Headers.GetDouble(uint16(header))
}

func IsSuccessStatus(statusCode int32) bool {
	return statusCode >= 200 && statusCode < 300
}

func IsRetriableStatus(statusCode int32, subStatus int) bool {
	switch statusCode {
	case 408, 503, 410:
		return true
	case 449:
		return true
	case 429:
		return true
	case 500:
		return subStatus == 0
	}
	return false
}
