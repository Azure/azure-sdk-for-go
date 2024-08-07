//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armservicefabricmanagedclusters

import "encoding/json"

func unmarshalPartitionClassification(rawMsg json.RawMessage) (PartitionClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b PartitionClassification
	switch m["partitionScheme"] {
	case string(PartitionSchemeNamed):
		b = &NamedPartitionScheme{}
	case string(PartitionSchemeSingleton):
		b = &SingletonPartitionScheme{}
	case string(PartitionSchemeUniformInt64Range):
		b = &UniformInt64RangePartitionScheme{}
	default:
		b = &Partition{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalScalingMechanismClassification(rawMsg json.RawMessage) (ScalingMechanismClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ScalingMechanismClassification
	switch m["kind"] {
	case string(ServiceScalingMechanismKindAddRemoveIncrementalNamedPartition):
		b = &AddRemoveIncrementalNamedPartitionScalingMechanism{}
	case string(ServiceScalingMechanismKindScalePartitionInstanceCount):
		b = &PartitionInstanceCountScaleMechanism{}
	default:
		b = &ScalingMechanism{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalScalingTriggerClassification(rawMsg json.RawMessage) (ScalingTriggerClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ScalingTriggerClassification
	switch m["kind"] {
	case string(ServiceScalingTriggerKindAveragePartitionLoadTrigger):
		b = &AveragePartitionLoadScalingTrigger{}
	case string(ServiceScalingTriggerKindAverageServiceLoadTrigger):
		b = &AverageServiceLoadScalingTrigger{}
	default:
		b = &ScalingTrigger{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalServicePlacementPolicyClassification(rawMsg json.RawMessage) (ServicePlacementPolicyClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ServicePlacementPolicyClassification
	switch m["type"] {
	case string(ServicePlacementPolicyTypeInvalidDomain):
		b = &ServicePlacementInvalidDomainPolicy{}
	case string(ServicePlacementPolicyTypeNonPartiallyPlaceService):
		b = &ServicePlacementNonPartiallyPlaceServicePolicy{}
	case string(ServicePlacementPolicyTypePreferredPrimaryDomain):
		b = &ServicePlacementPreferPrimaryDomainPolicy{}
	case string(ServicePlacementPolicyTypeRequiredDomain):
		b = &ServicePlacementRequiredDomainPolicy{}
	case string(ServicePlacementPolicyTypeRequiredDomainDistribution):
		b = &ServicePlacementRequireDomainDistributionPolicy{}
	default:
		b = &ServicePlacementPolicy{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalServicePlacementPolicyClassificationArray(rawMsg json.RawMessage) ([]ServicePlacementPolicyClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]ServicePlacementPolicyClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalServicePlacementPolicyClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalServiceResourcePropertiesClassification(rawMsg json.RawMessage) (ServiceResourcePropertiesClassification, error) {
	if rawMsg == nil || string(rawMsg) == "null" {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ServiceResourcePropertiesClassification
	switch m["serviceKind"] {
	case string(ServiceKindStateful):
		b = &StatefulServiceProperties{}
	case string(ServiceKindStateless):
		b = &StatelessServiceProperties{}
	default:
		b = &ServiceResourceProperties{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}
