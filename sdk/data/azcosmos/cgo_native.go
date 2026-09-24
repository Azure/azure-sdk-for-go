// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include <stddef.h>
#include "azurecosmosdriver.h"

#define COSMOS_ASSERT_OFFSET(type, field, expected) \
	_Static_assert(offsetof(type, field) == expected, #type "." #field " ABI offset changed")

_Static_assert(sizeof(cosmos_status_code_t) == 4, "cosmos_status_code_t must remain 32-bit");
_Static_assert((cosmos_status_code_t)-1 < 0, "cosmos_status_code_t must remain signed");
_Static_assert(COSMOS_OPERATION_KIND_CREATE_ITEM == 19, "create-item operation discriminant changed");
_Static_assert(COSMOS_OPERATION_KIND_READ_ITEM == 20, "read-item operation discriminant changed");
_Static_assert(COSMOS_OPERATION_KIND_UPSERT_ITEM == 21, "upsert-item operation discriminant changed");
_Static_assert(COSMOS_OPERATION_KIND_REPLACE_ITEM == 22, "replace-item operation discriminant changed");
_Static_assert(COSMOS_OPERATION_KIND_DELETE_ITEM == 23, "delete-item operation discriminant changed");
_Static_assert(COSMOS_OPERATION_KIND_PATCH_ITEM == 24, "patch-item operation discriminant changed");
_Static_assert(COSMOS_PRECONDITION_KIND_NONE == 0, "no-precondition discriminant changed");
_Static_assert(COSMOS_PRECONDITION_KIND_IF_MATCH == 1, "if-match precondition discriminant changed");
_Static_assert(COSMOS_PRECONDITION_KIND_IF_NONE_MATCH == 2, "if-none-match precondition discriminant changed");
_Static_assert(COSMOS_PATCH_STRATEGY_UNSET == 0, "unset PATCH strategy discriminant changed");
_Static_assert(COSMOS_PATCH_STRATEGY_AUTO == 1, "automatic PATCH strategy discriminant changed");
_Static_assert(COSMOS_PATCH_STRATEGY_CLIENT_SIDE == 2, "client-side PATCH strategy discriminant changed");
_Static_assert(COSMOS_PATCH_STRATEGY_SERVER_SIDE == 3, "server-side PATCH strategy discriminant changed");

_Static_assert(sizeof(cosmos_string_view_t) == 16, "cosmos_string_view_t ABI size changed");
_Static_assert(_Alignof(cosmos_string_view_t) == 8, "cosmos_string_view_t ABI alignment changed");
COSMOS_ASSERT_OFFSET(cosmos_string_view_t, data, 0);
COSMOS_ASSERT_OFFSET(cosmos_string_view_t, len, 8);

_Static_assert(sizeof(cosmos_value_payload_t) == 8, "cosmos_value_payload_t ABI size changed");
_Static_assert(_Alignof(cosmos_value_payload_t) == 8, "cosmos_value_payload_t ABI alignment changed");
_Static_assert(sizeof(cosmos_value_t) == 16, "cosmos_value_t ABI size changed");
_Static_assert(_Alignof(cosmos_value_t) == 8, "cosmos_value_t ABI alignment changed");
COSMOS_ASSERT_OFFSET(cosmos_value_t, payload, 8);

_Static_assert(sizeof(cosmos_partition_key_component_value_t) == 16, "partition-key value ABI size changed");
_Static_assert(_Alignof(cosmos_partition_key_component_value_t) == 8, "partition-key value ABI alignment changed");
_Static_assert(sizeof(cosmos_partition_key_component_t) == 24, "partition-key component ABI size changed");
_Static_assert(_Alignof(cosmos_partition_key_component_t) == 8, "partition-key component ABI alignment changed");
COSMOS_ASSERT_OFFSET(cosmos_partition_key_component_t, value, 8);

_Static_assert(sizeof(cosmos_completion_queue_options_t) == 12, "completion-queue options ABI size changed");
_Static_assert(_Alignof(cosmos_completion_queue_options_t) == 4, "completion-queue options ABI alignment changed");
COSMOS_ASSERT_OFFSET(cosmos_completion_queue_options_t, include_error_details, 8);

_Static_assert(sizeof(cosmos_runtime_options_t) == 64, "runtime options ABI size changed");
_Static_assert(_Alignof(cosmos_runtime_options_t) == 8, "runtime options ABI alignment changed");
COSMOS_ASSERT_OFFSET(cosmos_runtime_options_t, correlation_id, 8);
COSMOS_ASSERT_OFFSET(cosmos_runtime_options_t, user_agent_suffix, 24);
COSMOS_ASSERT_OFFSET(cosmos_runtime_options_t, wrapping_sdk_identifier, 40);
COSMOS_ASSERT_OFFSET(cosmos_runtime_options_t, cpu_refresh_interval_ms, 56);

_Static_assert(sizeof(cosmos_operation_options_t) == 96, "operation options ABI size changed");
_Static_assert(_Alignof(cosmos_operation_options_t) == 8, "operation options ABI alignment changed");
COSMOS_ASSERT_OFFSET(cosmos_operation_options_t, end_to_end_timeout_ms, 24);
COSMOS_ASSERT_OFFSET(cosmos_operation_options_t, excluded_regions, 56);
COSMOS_ASSERT_OFFSET(cosmos_operation_options_t, excluded_regions_len, 64);
COSMOS_ASSERT_OFFSET(cosmos_operation_options_t, binary_encoding_request_text_response, 89);

_Static_assert(sizeof(cosmos_driver_options_config_t) == 24, "driver options ABI size changed");
_Static_assert(_Alignof(cosmos_driver_options_config_t) == 8, "driver options ABI alignment changed");
COSMOS_ASSERT_OFFSET(cosmos_driver_options_config_t, preferred_regions_len, 8);
COSMOS_ASSERT_OFFSET(cosmos_driver_options_config_t, operation_options, 16);

_Static_assert(sizeof(cosmos_completion_t) == 112, "cosmos_completion_t ABI size changed");
_Static_assert(_Alignof(cosmos_completion_t) == 8, "cosmos_completion_t ABI alignment changed");
COSMOS_ASSERT_OFFSET(cosmos_completion_t, outcome, 0);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, status, 4);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, user_data, 8);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, was_cancel_requested, 16);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, http_status_code, 18);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, is_from_wire, 20);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, message, 24);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, headers, 48);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, headers_len, 56);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, body, 64);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, body_len, 72);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, diagnostics, 80);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, driver, 88);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, container, 96);
COSMOS_ASSERT_OFFSET(cosmos_completion_t, backing, 104);

_Static_assert(sizeof(cosmos_operation_request_t) == 224, "cosmos_operation_request_t ABI size changed");
_Static_assert(_Alignof(cosmos_operation_request_t) == 8, "cosmos_operation_request_t ABI alignment changed");
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, kind, 0);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, container, 24);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, item_id, 32);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, partition_key_components, 72);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, partition_key_len, 80);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, body, 96);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, body_len, 104);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, session_token, 112);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, max_item_count, 160);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, precondition_kind, 172);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, precondition_etag, 176);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, options, 192);
*/
import "C"

// This file carries the cgo directives for the whole package. They are declared once here rather
// than repeated in each file that imports "C", because cgo unions the directives across the
// package: repeating them means every file has to be kept in step, and a file that drifts is a
// link error rather than a compile error.
//
// Target-specific external modules carry the archives and system-linker flags. The conditionally
// imported internal packages compile the Go token callback shim for the same target. Keeping both
// behind root-package constraints matters because Go treats ios as darwin and android as linux.
