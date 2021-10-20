//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmysql

import "encoding/json"

func unmarshalServerPropertiesForCreateClassification(rawMsg json.RawMessage) (ServerPropertiesForCreateClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ServerPropertiesForCreateClassification
	switch m["createMode"] {
	case string(CreateModeDefault):
		b = &ServerPropertiesForDefaultCreate{}
	case string(CreateModeGeoRestore):
		b = &ServerPropertiesForGeoRestore{}
	case string(CreateModePointInTimeRestore):
		b = &ServerPropertiesForRestore{}
	case string(CreateModeReplica):
		b = &ServerPropertiesForReplica{}
	default:
		b = &ServerPropertiesForCreate{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalServerPropertiesForCreateClassificationArray(rawMsg json.RawMessage) ([]ServerPropertiesForCreateClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]ServerPropertiesForCreateClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalServerPropertiesForCreateClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}
