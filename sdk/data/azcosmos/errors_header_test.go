// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore IMDS

package azcosmos

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Each header symbol needs an explicit classification decision, including reserved values.
// CodeClientError deliberately preserves the numeric substatus without promising a finer category.
var nativeSubStatusMappings = []struct {
	symbol string
	value  int
	code   Code
}{
	{"TRANSPORT_GENERATED_503", 20003, CodeTransportFailure},
	{"CLIENT_CPU_OVERLOAD", 20004, CodeClientError},
	{"CLIENT_THREAD_STARVATION", 20005, CodeClientError},
	{"CHANNEL_CLOSED", 20006, CodeClientError},
	{"MALFORMED_CONTINUATION_TOKEN", 20007, CodeClientError},
	{"CLIENT_OPERATION_TIMEOUT", 20008, CodeClientOperationTimeout},
	{"TRANSPORT_CONNECTION_FAILED", 20010, CodeTransportFailure},
	{"TRANSPORT_IO_FAILED", 20011, CodeTransportFailure},
	{"TRANSPORT_DNS_FAILED", 20012, CodeTransportFailure},
	{"TRANSPORT_BODY_READ_FAILED", 20014, CodeTransportFailure},
	{"TRANSPORT_HTTP2_INCOMPATIBLE", 20015, CodeTransportFailure},
	{"SERIALIZATION_RESPONSE_BODY_INVALID", 20020, CodeSerializationFailed},
	{"CLIENT_PARTITION_KEY_EMPTY", 20100, CodeClientError},
	{"CLIENT_PARTITION_KEY_TOO_MANY_COMPONENTS", 20101, CodeClientError},
	{"CLIENT_PREFIX_PARTITION_KEY_REQUIRES_MULTIHASH", 20102, CodeClientError},
	{"CLIENT_CONNECTION_STRING_EMPTY", 20104, CodeClientError},
	{"CLIENT_CONNECTION_STRING_MALFORMED_PART", 20105, CodeClientError},
	{"CLIENT_CONNECTION_STRING_MISSING_ACCOUNT_KEY", 20107, CodeClientError},
	{"CLIENT_INVALID_ACCOUNT_ENDPOINT_URL", 20108, CodeClientError},
	{"CLIENT_INVALID_URL", 20109, CodeClientError},
	{"CLIENT_UNKNOWN_CONSISTENCY_LEVEL", 20110, CodeClientError},
	{"CLIENT_UNKNOWN_PRIORITY_LEVEL", 20111, CodeClientError},
	{"CLIENT_FEED_RANGE_REQUIRES_FANOUT_PIPELINE", 20112, CodeClientError},
	{"CLIENT_UNSUPPORTED_QUERY_FEATURE", 20113, CodeClientError},
	{"CLIENT_QUERY_PLAN_INVALID_TOP_OFFSET_LIMIT", 20114, CodeClientError},
	{"CLIENT_CONTINUATION_TOKEN_NON_QUERY_OPERATION", 20117, CodeClientError},
	{"CLIENT_DUPLICATE_FAULT_INJECTION_RULE_ID", 20150, CodeClientError},
	{"CLIENT_HTTP_CLIENT_CONSTRUCTION_FAILED", 20153, CodeClientError},
	{"CLIENT_REQWEST_FEATURE_REQUIRED", 20154, CodeClientError},
	{"CLIENT_REQUEST_URL_MISSING_HOST", 20155, CodeClientError},
	{"CLIENT_REQUEST_URL_MISSING_KNOWN_PORT", 20156, CodeClientError},
	{"CLIENT_IMDS_HTTP_CLIENT_CONSTRUCTION_FAILED", 20157, CodeClientError},
	{"CLIENT_IMDS_REQWEST_FEATURE_REQUIRED", 20158, CodeClientError},
	{"CLIENT_CONTINUATION_TOKEN_FETCH_IN_FLIGHT", 20200, CodeClientError},
	{"CLIENT_TOPOLOGY_PROVIDER_MISSING", 20201, CodeClientError},
	{"CLIENT_DRIVER_NOT_INITIALIZED", 20202, CodeClientError},
	{"CLIENT_CONTINUATION_TOKEN_SHAPE_MISMATCH", 20203, CodeClientError},
	{"CLIENT_CONTINUATION_TOKEN_INVALID_EPK_RANGE", 20205, CodeClientError},
	{"CLIENT_SPLIT_RETRIES_EXHAUSTED", 20206, CodeClientError},
	{"CLIENT_BUILD_RESPONSE_INVOKED_ON_FAILURE", 20207, CodeClientError},
	{"CLIENT_ROOT_NODE_CANNOT_REQUEST_SPLIT", 20208, CodeClientError},
	{"CLIENT_SINGLETON_OPERATION_RETURNED_EMPTY_PAGE", 20210, CodeClientError},
	{"CLIENT_CONTINUATION_TOKEN_SAVED_RANGE_UNHONORED", 20213, CodeClientError},
	{"CLIENT_NO_THROUGHPUT_OFFER_FOR_RESOURCE", 20301, CodeClientError},
	{"CLIENT_QUERY_PLAN_PRODUCED_EMPTY_RANGES", 20302, CodeClientError},
	{"SERVICE_RETURNED_OFFER_WITHOUT_ID", 20303, CodeClientError},
	{"CLIENT_THROUGHPUT_POLLER_INCOMPLETE", 20304, CodeClientError},
	{"CLIENT_TOPOLOGY_RESOLUTION_FAILED", 20305, CodeClientError},
	{"SERVICE_RETURNED_OBJECT_WITHOUT_RID", 20306, CodeClientError},
	{"CLIENT_READ_SESSION_NOT_AVAILABLE", 20310, CodeServiceUnavailable},
	{"CLIENT_WRITE_FORBIDDEN", 20311, CodeServiceUnavailable},
	{"CLIENT_DATABASE_ACCOUNT_NOT_FOUND", 20312, CodeServiceUnavailable},
	{"CLIENT_FFI_NULL_ARGUMENT", 20350, CodeClientError},
	{"CLIENT_FFI_INVALID_UTF8", 20351, CodeClientError},
	{"CLIENT_FFI_INVALID_HEADER", 20352, CodeClientError},
	{"CLIENT_FFI_INVALID_OPTION_VALUE", 20353, CodeClientError},
	{"CLIENT_FFI_OPERATION_CONSUMED", 20354, CodeClientError},
	{"CLIENT_FFI_PRECONDITION_ALREADY_SET", 20355, CodeClientError},
	{"CLIENT_FFI_UNSUPPORTED_OPERATION_FOR_MUTATOR", 20356, CodeClientError},
	{"CLIENT_FFI_FEED_EXHAUSTED", 20357, CodeClientError},
	{"CLIENT_FFI_QUEUE_SHUTDOWN", 20358, CodeClientError},
	{"CLIENT_FFI_QUEUE_FULL", 20359, CodeClientError},
	{"CLIENT_FFI_OPERATION_CANCELLED", 20360, CodeClientError},
	{"CLIENT_FFI_RUNTIME_BUILD_FAILED", 20361, CodeClientError},
	{"CLIENT_FFI_PANIC", 20362, CodeClientError},
	{"CLIENT_FFI_CURSOR_BUSY", 20363, CodeClientError},
	{"CLIENT_FFI_QUEUE_FORMAT", 20364, CodeClientError},
	{"CLIENT_FFI_CURSOR_CLOSED", 20365, CodeClientError},
	{"CLIENT_FFI_REPRESENTATION_UNSUPPORTED", 20366, CodeClientError},
	{"CLIENT_FFI_DELIVERY_LOST", 20367, CodeClientError},
	{"CLIENT_GENERATED_401", 20401, CodeAuthenticationFailed},
	{"AUTHENTICATION_TOKEN_ACQUISITION_FAILED", 20402, CodeAuthenticationFailed},
	{"TRANSIT_TIMEOUT", 20911, CodeClientError},
	{"SERVER_BARRIER_THROTTLED", 21011, CodeClientError},
}

func TestNativeSubStatusHeaderDrift(t *testing.T) {
	header, err := os.ReadFile("azurecosmosdriver.h")
	require.NoError(t, err)
	// Strip comments before parsing so commented-out symbols cannot satisfy the inventory.
	text := regexp.MustCompile(`(?s)/\*.*?\*/|//[^\n]*`).ReplaceAllString(string(header), "")
	enum := regexp.MustCompile(`(?s)enum cosmos_sub_status_t\b[^{]*\{([^}]+)\}`).FindStringSubmatch(text)
	require.Len(t, enum, 2, "cannot find complete cosmos_sub_status_t; review the native ABI")
	entry := regexp.MustCompile(`^\s*(COSMOS_SUB_STATUS_[A-Z0-9_]+)\s*=\s*([0-9]+)\s*$`)
	actual := map[string]int{}
	for _, declaration := range strings.Split(enum[1], ",") {
		if strings.TrimSpace(declaration) == "" {
			continue
		}
		match := entry.FindStringSubmatch(declaration)
		require.Len(t, match, 3, "unrecognized native declaration %q; update the drift parser and make a mapping/fallback decision", declaration)
		value, err := strconv.Atoi(match[2])
		require.NoError(t, err)
		require.NotContains(t, actual, match[1], "duplicate native substatus")
		actual[match[1]] = value
	}
	expected := map[string]int{}
	for _, mapping := range nativeSubStatusMappings {
		symbol := "COSMOS_SUB_STATUS_" + mapping.symbol
		require.NotContains(t, expected, symbol, "duplicate mapping decision")
		expected[symbol] = mapping.value
	}
	require.Equal(t, expected, actual,
		"native substatus symbols/values changed: review errors.go and update nativeSubStatusMappings with an explicit mapping or CodeClientError fallback; do not edit the generated header")
}
