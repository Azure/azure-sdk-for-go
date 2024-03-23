//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents

import (
	"encoding/json"
	"fmt"
)

//type IotHubDeviceConnectedEventData DeviceConnectionStateEventProperties

func fixNAValue(s **string) {
	if *s != nil && **s == "n/a" {
		*s = nil
	}
}

func unmarshalInternalACSAdvancedMessageChannelEventError(val json.RawMessage, fn string, re **Error) error {
	var realErr *internalACSAdvancedMessageChannelEventError

	if err := json.Unmarshal(val, &realErr); err != nil {
		return fmt.Errorf("struct field %s: %v", fn, err)
	}

	*re = &Error{
		Code:    *realErr.ChannelCode,
		message: *realErr.ChannelMessage,
	}

	return nil
}

func unmarshalInternalACSRouterCommunicationError(val json.RawMessage, fn string, re *[]*Error) error {
	type superError struct {
		Code       string
		InnerError *superError
		Details    []*superError
		Message    string
	}

	var unpack func(err *superError) *Error

	unpack = func(err *superError) *Error {
		if err == nil {
			return nil
		}

		e := &Error{
			Code:       err.Code,
			message:    err.Message,
			InnerError: unpack(err.InnerError),
		}

		for _, detailErr := range err.Details {
			e.Details = append(e.Details, unpack(detailErr))
		}

		return e
	}

	var tmp []*superError

	if err := json.Unmarshal(val, &tmp); err != nil {
		return fmt.Errorf("struct field %s: %v", fn, err)
	}

	for _, se := range tmp {
		*re = append(*re, unpack(se))
	}

	return nil
}
