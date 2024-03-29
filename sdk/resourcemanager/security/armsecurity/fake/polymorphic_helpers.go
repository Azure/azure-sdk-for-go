//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
)

func unmarshalServerVulnerabilityAssessmentsSettingClassification(rawMsg json.RawMessage) (armsecurity.ServerVulnerabilityAssessmentsSettingClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b armsecurity.ServerVulnerabilityAssessmentsSettingClassification
	switch m["kind"] {
	case string(armsecurity.ServerVulnerabilityAssessmentsSettingKindAzureServersSetting):
		b = &armsecurity.AzureServersSetting{}
	default:
		b = &armsecurity.ServerVulnerabilityAssessmentsSetting{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalSettingClassification(rawMsg json.RawMessage) (armsecurity.SettingClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b armsecurity.SettingClassification
	switch m["kind"] {
	case string(armsecurity.SettingKindAlertSyncSettings):
		b = &armsecurity.AlertSyncSettings{}
	case string(armsecurity.SettingKindDataExportSettings):
		b = &armsecurity.DataExportSettings{}
	default:
		b = &armsecurity.Setting{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}
