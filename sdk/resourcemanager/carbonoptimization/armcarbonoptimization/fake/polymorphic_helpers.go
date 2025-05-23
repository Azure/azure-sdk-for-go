// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package fake

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/carbonoptimization/armcarbonoptimization"
)

func unmarshalQueryFilterClassification(rawMsg json.RawMessage) (armcarbonoptimization.QueryFilterClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b armcarbonoptimization.QueryFilterClassification
	switch m["reportType"] {
	case string(armcarbonoptimization.ReportTypeEnumOverallSummaryReport):
		b = &armcarbonoptimization.OverallSummaryReportQueryFilter{}
	case string(armcarbonoptimization.ReportTypeEnumMonthlySummaryReport):
		b = &armcarbonoptimization.MonthlySummaryReportQueryFilter{}
	case string(armcarbonoptimization.ReportTypeEnumTopItemsSummaryReport):
		b = &armcarbonoptimization.TopItemsSummaryReportQueryFilter{}
	case string(armcarbonoptimization.ReportTypeEnumTopItemsMonthlySummaryReport):
		b = &armcarbonoptimization.TopItemsMonthlySummaryReportQueryFilter{}
	case string(armcarbonoptimization.ReportTypeEnumItemDetailsReport):
		b = &armcarbonoptimization.ItemDetailsQueryFilter{}
	default:
		b = &armcarbonoptimization.QueryFilter{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}
