//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents

import (
	"encoding/json"
	"fmt"
	"strings"
)

func unmarshalInternalACSMessageChannelEventError(val json.RawMessage, fn string, re **Error) error {
	var realErr *internalACSMessageChannelEventError

	if err := json.Unmarshal(val, &realErr); err != nil {
		return fmt.Errorf("struct field %s: %w", fn, err)
	}

	*re = &Error{
		Code:    *realErr.ChannelCode,
		message: *realErr.ChannelMessage,
	}

	return nil
}

func unpackMsg(err internalACSRouterCommunicationError, indent string) string {
	sb := strings.Builder{}

	code := ""
	if err.Code != nil {
		code = *err.Code
	}

	message := ""
	if err.Message != nil {
		message = *err.Message
	}

	sb.WriteString(fmt.Sprintf("%sCode: %s\n%sMessage: %s\n", indent, code, indent, message))

	if len(err.Details) > 0 {
		for i, detailErr := range err.Details {
			sb.WriteString(fmt.Sprintf("%sDetails[%d]:\n%s", indent, i, unpackMsg(detailErr, indent+"  ")))
		}
	}

	if err.Innererror != nil {
		sb.WriteString(fmt.Sprintf("%sInnerError:\n%s", indent, unpackMsg(*err.Innererror, indent+"  ")))
	}

	return sb.String()
}

func unmarshalInternalACSRouterCommunicationError(val json.RawMessage, fn string, re *[]*Error) error {
	var tmp []internalACSRouterCommunicationError

	if err := json.Unmarshal(val, &tmp); err != nil {
		return fmt.Errorf("struct field %s: %w", fn, err)
	}

	for _, se := range tmp {
		code := ""

		if se.Code != nil {
			code = *se.Code
		}

		e := &Error{
			Code: code,
			// we're going to compress the remainder of these details into a
			// string.
			message: unpackMsg(se, ""),
		}

		*re = append(*re, e)
	}

	return nil
}
