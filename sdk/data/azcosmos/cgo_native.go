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

_Static_assert(sizeof(cosmos_operation_request_t) == 168, "cosmos_operation_request_t ABI size changed");
_Static_assert(_Alignof(cosmos_operation_request_t) == 8, "cosmos_operation_request_t ABI alignment changed");
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, kind, 0);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, container, 24);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, item_id, 32);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, partition_key_components, 56);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, partition_key_len, 64);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, body, 80);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, body_len, 88);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, session_token, 96);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, max_item_count, 120);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, precondition_kind, 132);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, precondition_etag, 136);
COSMOS_ASSERT_OFFSET(cosmos_operation_request_t, options, 144);
*/
import "C"

// This file carries the cgo directives for the whole package. They are declared once here rather
// than repeated in each file that imports "C", because cgo unions the directives across the
// package: repeating them means every file has to be kept in step, and a file that drifts is a
// link error rather than a compile error.
//
// Target-specific internal packages carry the .syso archives and system-linker flags. Keeping each
// archive behind a conditionally imported package matters because Go treats ios as darwin and
// android as linux for build tags and .syso filename selection.
