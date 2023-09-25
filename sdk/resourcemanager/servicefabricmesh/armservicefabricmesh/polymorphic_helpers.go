//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armservicefabricmesh

import "encoding/json"

func unmarshalApplicationScopedVolumeCreationParametersClassification(rawMsg json.RawMessage) (ApplicationScopedVolumeCreationParametersClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ApplicationScopedVolumeCreationParametersClassification
	switch m["kind"] {
	case string(ApplicationScopedVolumeKindServiceFabricVolumeDisk):
		b = &ApplicationScopedVolumeCreationParametersServiceFabricVolumeDisk{}
	default:
		b = &ApplicationScopedVolumeCreationParameters{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalAutoScalingMechanismClassification(rawMsg json.RawMessage) (AutoScalingMechanismClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b AutoScalingMechanismClassification
	switch m["kind"] {
	case string(AutoScalingMechanismKindAddRemoveReplica):
		b = &AddRemoveReplicaScalingMechanism{}
	default:
		b = &AutoScalingMechanism{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalAutoScalingMetricClassification(rawMsg json.RawMessage) (AutoScalingMetricClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b AutoScalingMetricClassification
	switch m["kind"] {
	case string(AutoScalingMetricKindResource):
		b = &AutoScalingResourceMetric{}
	default:
		b = &AutoScalingMetric{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalAutoScalingTriggerClassification(rawMsg json.RawMessage) (AutoScalingTriggerClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b AutoScalingTriggerClassification
	switch m["kind"] {
	case string(AutoScalingTriggerKindAverageLoad):
		b = &AverageLoadScalingTrigger{}
	default:
		b = &AutoScalingTrigger{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalDiagnosticsSinkPropertiesClassification(rawMsg json.RawMessage) (DiagnosticsSinkPropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b DiagnosticsSinkPropertiesClassification
	switch m["kind"] {
	case string(DiagnosticsSinkKindAzureInternalMonitoringPipeline):
		b = &AzureInternalMonitoringPipelineSinkDescription{}
	default:
		b = &DiagnosticsSinkProperties{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalDiagnosticsSinkPropertiesClassificationArray(rawMsg json.RawMessage) ([]DiagnosticsSinkPropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]DiagnosticsSinkPropertiesClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalDiagnosticsSinkPropertiesClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalNetworkResourcePropertiesClassification(rawMsg json.RawMessage) (NetworkResourcePropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b NetworkResourcePropertiesClassification
	switch m["kind"] {
	case string(NetworkKindLocal):
		b = &LocalNetworkResourceProperties{}
	default:
		b = &NetworkResourceProperties{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalSecretResourcePropertiesClassification(rawMsg json.RawMessage) (SecretResourcePropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b SecretResourcePropertiesClassification
	switch m["kind"] {
	case string(SecretKindInlinedValue):
		b = &InlinedValueSecretResourceProperties{}
	default:
		b = &SecretResourceProperties{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

