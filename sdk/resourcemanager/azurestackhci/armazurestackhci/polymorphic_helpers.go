//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armazurestackhci

import "encoding/json"

func unmarshalEdgeDeviceClassification(rawMsg json.RawMessage) (EdgeDeviceClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b EdgeDeviceClassification
	switch m["kind"] {
	case string(DeviceKindHCI):
		b = &HciEdgeDevice{}
	default:
		b = &EdgeDevice{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalEdgeDeviceClassificationArray(rawMsg json.RawMessage) ([]EdgeDeviceClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]EdgeDeviceClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalEdgeDeviceClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}
