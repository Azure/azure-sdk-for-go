//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armtestbase

import "encoding/json"

func unmarshalAnalysisResultSingletonResourcePropertiesClassification(rawMsg json.RawMessage) (AnalysisResultSingletonResourcePropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b AnalysisResultSingletonResourcePropertiesClassification
	switch m["analysisResultType"] {
	case string(AnalysisResultTypeCPURegression):
		b = &CPURegressionResultSingletonResourceProperties{}
	case string(AnalysisResultTypeCPUUtilization):
		b = &CPUUtilizationResultSingletonResourceProperties{}
	case string(AnalysisResultTypeMemoryRegression):
		b = &MemoryRegressionResultSingletonResourceProperties{}
	case string(AnalysisResultTypeMemoryUtilization):
		b = &MemoryUtilizationResultSingletonResourceProperties{}
	case string(AnalysisResultTypeReliability):
		b = &ReliabilityResultSingletonResourceProperties{}
	case string(AnalysisResultTypeScriptExecution):
		b = &ScriptExecutionResultSingletonResourceProperties{}
	case string(AnalysisResultTypeTestAnalysis):
		b = &TestAnalysisResultSingletonResourceProperties{}
	default:
		b = &AnalysisResultSingletonResourceProperties{}
	}
	return b, json.Unmarshal(rawMsg, b)
}
