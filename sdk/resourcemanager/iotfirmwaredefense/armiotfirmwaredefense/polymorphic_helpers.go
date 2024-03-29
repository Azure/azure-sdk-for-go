//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armiotfirmwaredefense

import "encoding/json"

func unmarshalSummaryResourcePropertiesClassification(rawMsg json.RawMessage) (SummaryResourcePropertiesClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b SummaryResourcePropertiesClassification
	switch m["summaryType"] {
	case string(SummaryTypeBinaryHardening):
		b = &BinaryHardeningSummaryResource{}
	case string(SummaryTypeCVE):
		b = &CveSummary{}
	case string(SummaryTypeCryptoCertificate):
		b = &CryptoCertificateSummaryResource{}
	case string(SummaryTypeCryptoKey):
		b = &CryptoKeySummaryResource{}
	case string(SummaryTypeFirmware):
		b = &FirmwareSummary{}
	default:
		b = &SummaryResourceProperties{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}
