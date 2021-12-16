//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

import "encoding/json"

func unmarshalAdditionalDataClassification(rawMsg json.RawMessage) (AdditionalDataClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b AdditionalDataClassification
	switch m["assessedResourceType"] {
	case string(AssessedResourceTypeContainerRegistryVulnerability):
		b = &ContainerRegistryVulnerabilityProperties{}
	case "ServerVulnerabilityAssessment":
		b = &ServerVulnerabilityProperties{}
	case string(AssessedResourceTypeSQLServerVulnerability):
		b = &SQLServerVulnerabilityProperties{}
	default:
		b = &AdditionalData{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalAdditionalDataClassificationArray(rawMsg json.RawMessage) ([]AdditionalDataClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]AdditionalDataClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalAdditionalDataClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalAdditionalDataClassificationMap(rawMsg json.RawMessage) (map[string]AdditionalDataClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]AdditionalDataClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalAdditionalDataClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalAlertSimulatorRequestPropertiesClassification(rawMsg json.RawMessage) (AlertSimulatorRequestPropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b AlertSimulatorRequestPropertiesClassification
	switch m["kind"] {
	case string(KindBundles):
		b = &AlertSimulatorBundlesRequestProperties{}
	default:
		b = &AlertSimulatorRequestProperties{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalAlertSimulatorRequestPropertiesClassificationArray(rawMsg json.RawMessage) ([]AlertSimulatorRequestPropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]AlertSimulatorRequestPropertiesClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalAlertSimulatorRequestPropertiesClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalAlertSimulatorRequestPropertiesClassificationMap(rawMsg json.RawMessage) (map[string]AlertSimulatorRequestPropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]AlertSimulatorRequestPropertiesClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalAlertSimulatorRequestPropertiesClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalAllowlistCustomAlertRuleClassification(rawMsg json.RawMessage) (AllowlistCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b AllowlistCustomAlertRuleClassification
	switch m["ruleType"] {
	case "ConnectionFromIpNotAllowed":
		b = &ConnectionFromIPNotAllowed{}
	case "ConnectionToIpNotAllowed":
		b = &ConnectionToIPNotAllowed{}
	case "LocalUserNotAllowed":
		b = &LocalUserNotAllowed{}
	case "ProcessNotAllowed":
		b = &ProcessNotAllowed{}
	default:
		b = &AllowlistCustomAlertRule{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalAllowlistCustomAlertRuleClassificationArray(rawMsg json.RawMessage) ([]AllowlistCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]AllowlistCustomAlertRuleClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalAllowlistCustomAlertRuleClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalAllowlistCustomAlertRuleClassificationMap(rawMsg json.RawMessage) (map[string]AllowlistCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]AllowlistCustomAlertRuleClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalAllowlistCustomAlertRuleClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalAuthenticationDetailsPropertiesClassification(rawMsg json.RawMessage) (AuthenticationDetailsPropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b AuthenticationDetailsPropertiesClassification
	switch m["authenticationType"] {
	case string(AuthenticationTypeAwsAssumeRole):
		b = &AwAssumeRoleAuthenticationDetailsProperties{}
	case string(AuthenticationTypeAwsCreds):
		b = &AwsCredsAuthenticationDetailsProperties{}
	case string(AuthenticationTypeGcpCredentials):
		b = &GcpCredentialsDetailsProperties{}
	default:
		b = &AuthenticationDetailsProperties{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalAuthenticationDetailsPropertiesClassificationArray(rawMsg json.RawMessage) ([]AuthenticationDetailsPropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]AuthenticationDetailsPropertiesClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalAuthenticationDetailsPropertiesClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalAuthenticationDetailsPropertiesClassificationMap(rawMsg json.RawMessage) (map[string]AuthenticationDetailsPropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]AuthenticationDetailsPropertiesClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalAuthenticationDetailsPropertiesClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalAutomationActionClassification(rawMsg json.RawMessage) (AutomationActionClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b AutomationActionClassification
	switch m["actionType"] {
	case string(ActionTypeEventHub):
		b = &AutomationActionEventHub{}
	case string(ActionTypeLogicApp):
		b = &AutomationActionLogicApp{}
	case string(ActionTypeWorkspace):
		b = &AutomationActionWorkspace{}
	default:
		b = &AutomationAction{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalAutomationActionClassificationArray(rawMsg json.RawMessage) ([]AutomationActionClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]AutomationActionClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalAutomationActionClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalAutomationActionClassificationMap(rawMsg json.RawMessage) (map[string]AutomationActionClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]AutomationActionClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalAutomationActionClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalCloudOfferingClassification(rawMsg json.RawMessage) (CloudOfferingClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b CloudOfferingClassification
	switch m["offeringType"] {
	case string(OfferingTypeCspmMonitorAws):
		b = &CspmMonitorAwsOffering{}
	case string(OfferingTypeDefenderForContainersAws):
		b = &DefenderForContainersAwsOffering{}
	case "DefenderForServersAWS":
		b = &DefenderForServersAwsOffering{}
	default:
		b = &CloudOffering{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalCloudOfferingClassificationArray(rawMsg json.RawMessage) ([]CloudOfferingClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]CloudOfferingClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalCloudOfferingClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalCloudOfferingClassificationMap(rawMsg json.RawMessage) (map[string]CloudOfferingClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]CloudOfferingClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalCloudOfferingClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalCustomAlertRuleClassification(rawMsg json.RawMessage) (CustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b CustomAlertRuleClassification
	switch m["ruleType"] {
	case "ActiveConnectionsNotInAllowedRange":
		b = &ActiveConnectionsNotInAllowedRange{}
	case "AllowlistCustomAlertRule":
		b = &AllowlistCustomAlertRule{}
	case "AmqpC2DMessagesNotInAllowedRange":
		b = &AmqpC2DMessagesNotInAllowedRange{}
	case "AmqpC2DRejectedMessagesNotInAllowedRange":
		b = &AmqpC2DRejectedMessagesNotInAllowedRange{}
	case "AmqpD2CMessagesNotInAllowedRange":
		b = &AmqpD2CMessagesNotInAllowedRange{}
	case "ConnectionFromIpNotAllowed":
		b = &ConnectionFromIPNotAllowed{}
	case "ConnectionToIpNotAllowed":
		b = &ConnectionToIPNotAllowed{}
	case "DenylistCustomAlertRule":
		b = &DenylistCustomAlertRule{}
	case "DirectMethodInvokesNotInAllowedRange":
		b = &DirectMethodInvokesNotInAllowedRange{}
	case "FailedLocalLoginsNotInAllowedRange":
		b = &FailedLocalLoginsNotInAllowedRange{}
	case "FileUploadsNotInAllowedRange":
		b = &FileUploadsNotInAllowedRange{}
	case "HttpC2DMessagesNotInAllowedRange":
		b = &HTTPC2DMessagesNotInAllowedRange{}
	case "HttpC2DRejectedMessagesNotInAllowedRange":
		b = &HTTPC2DRejectedMessagesNotInAllowedRange{}
	case "HttpD2CMessagesNotInAllowedRange":
		b = &HTTPD2CMessagesNotInAllowedRange{}
	case "ListCustomAlertRule":
		b = &ListCustomAlertRule{}
	case "LocalUserNotAllowed":
		b = &LocalUserNotAllowed{}
	case "MqttC2DMessagesNotInAllowedRange":
		b = &MqttC2DMessagesNotInAllowedRange{}
	case "MqttC2DRejectedMessagesNotInAllowedRange":
		b = &MqttC2DRejectedMessagesNotInAllowedRange{}
	case "MqttD2CMessagesNotInAllowedRange":
		b = &MqttD2CMessagesNotInAllowedRange{}
	case "ProcessNotAllowed":
		b = &ProcessNotAllowed{}
	case "QueuePurgesNotInAllowedRange":
		b = &QueuePurgesNotInAllowedRange{}
	case "ThresholdCustomAlertRule":
		b = &ThresholdCustomAlertRule{}
	case "TimeWindowCustomAlertRule":
		b = &TimeWindowCustomAlertRule{}
	case "TwinUpdatesNotInAllowedRange":
		b = &TwinUpdatesNotInAllowedRange{}
	case "UnauthorizedOperationsNotInAllowedRange":
		b = &UnauthorizedOperationsNotInAllowedRange{}
	default:
		b = &CustomAlertRule{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalCustomAlertRuleClassificationArray(rawMsg json.RawMessage) ([]CustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]CustomAlertRuleClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalCustomAlertRuleClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalCustomAlertRuleClassificationMap(rawMsg json.RawMessage) (map[string]CustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]CustomAlertRuleClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalCustomAlertRuleClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalListCustomAlertRuleClassification(rawMsg json.RawMessage) (ListCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ListCustomAlertRuleClassification
	switch m["ruleType"] {
	case "AllowlistCustomAlertRule":
		b = &AllowlistCustomAlertRule{}
	case "ConnectionFromIpNotAllowed":
		b = &ConnectionFromIPNotAllowed{}
	case "ConnectionToIpNotAllowed":
		b = &ConnectionToIPNotAllowed{}
	case "DenylistCustomAlertRule":
		b = &DenylistCustomAlertRule{}
	case "LocalUserNotAllowed":
		b = &LocalUserNotAllowed{}
	case "ProcessNotAllowed":
		b = &ProcessNotAllowed{}
	default:
		b = &ListCustomAlertRule{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalListCustomAlertRuleClassificationArray(rawMsg json.RawMessage) ([]ListCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]ListCustomAlertRuleClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalListCustomAlertRuleClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalListCustomAlertRuleClassificationMap(rawMsg json.RawMessage) (map[string]ListCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]ListCustomAlertRuleClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalListCustomAlertRuleClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalResourceDetailsClassification(rawMsg json.RawMessage) (ResourceDetailsClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ResourceDetailsClassification
	switch m["source"] {
	default:
		b = &ResourceDetails{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalResourceDetailsClassificationArray(rawMsg json.RawMessage) ([]ResourceDetailsClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]ResourceDetailsClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalResourceDetailsClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalResourceDetailsClassificationMap(rawMsg json.RawMessage) (map[string]ResourceDetailsClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]ResourceDetailsClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalResourceDetailsClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalResourceIdentifierClassification(rawMsg json.RawMessage) (ResourceIdentifierClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ResourceIdentifierClassification
	switch m["type"] {
	case string(ResourceIdentifierTypeAzureResource):
		b = &AzureResourceIdentifier{}
	case string(ResourceIdentifierTypeLogAnalytics):
		b = &LogAnalyticsIdentifier{}
	default:
		b = &ResourceIdentifier{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalResourceIdentifierClassificationArray(rawMsg json.RawMessage) ([]ResourceIdentifierClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]ResourceIdentifierClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalResourceIdentifierClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalResourceIdentifierClassificationMap(rawMsg json.RawMessage) (map[string]ResourceIdentifierClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]ResourceIdentifierClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalResourceIdentifierClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalSettingClassification(rawMsg json.RawMessage) (SettingClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b SettingClassification
	switch m["kind"] {
	case string(SettingKindAlertSyncSettings):
		b = &AlertSyncSettings{}
	case string(SettingKindDataExportSettings):
		b = &DataExportSettings{}
	default:
		b = &Setting{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalSettingClassificationArray(rawMsg json.RawMessage) ([]SettingClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]SettingClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalSettingClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalSettingClassificationMap(rawMsg json.RawMessage) (map[string]SettingClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]SettingClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalSettingClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalThresholdCustomAlertRuleClassification(rawMsg json.RawMessage) (ThresholdCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ThresholdCustomAlertRuleClassification
	switch m["ruleType"] {
	case "ActiveConnectionsNotInAllowedRange":
		b = &ActiveConnectionsNotInAllowedRange{}
	case "AmqpC2DMessagesNotInAllowedRange":
		b = &AmqpC2DMessagesNotInAllowedRange{}
	case "AmqpC2DRejectedMessagesNotInAllowedRange":
		b = &AmqpC2DRejectedMessagesNotInAllowedRange{}
	case "AmqpD2CMessagesNotInAllowedRange":
		b = &AmqpD2CMessagesNotInAllowedRange{}
	case "DirectMethodInvokesNotInAllowedRange":
		b = &DirectMethodInvokesNotInAllowedRange{}
	case "FailedLocalLoginsNotInAllowedRange":
		b = &FailedLocalLoginsNotInAllowedRange{}
	case "FileUploadsNotInAllowedRange":
		b = &FileUploadsNotInAllowedRange{}
	case "HttpC2DMessagesNotInAllowedRange":
		b = &HTTPC2DMessagesNotInAllowedRange{}
	case "HttpC2DRejectedMessagesNotInAllowedRange":
		b = &HTTPC2DRejectedMessagesNotInAllowedRange{}
	case "HttpD2CMessagesNotInAllowedRange":
		b = &HTTPD2CMessagesNotInAllowedRange{}
	case "MqttC2DMessagesNotInAllowedRange":
		b = &MqttC2DMessagesNotInAllowedRange{}
	case "MqttC2DRejectedMessagesNotInAllowedRange":
		b = &MqttC2DRejectedMessagesNotInAllowedRange{}
	case "MqttD2CMessagesNotInAllowedRange":
		b = &MqttD2CMessagesNotInAllowedRange{}
	case "QueuePurgesNotInAllowedRange":
		b = &QueuePurgesNotInAllowedRange{}
	case "TimeWindowCustomAlertRule":
		b = &TimeWindowCustomAlertRule{}
	case "TwinUpdatesNotInAllowedRange":
		b = &TwinUpdatesNotInAllowedRange{}
	case "UnauthorizedOperationsNotInAllowedRange":
		b = &UnauthorizedOperationsNotInAllowedRange{}
	default:
		b = &ThresholdCustomAlertRule{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalThresholdCustomAlertRuleClassificationArray(rawMsg json.RawMessage) ([]ThresholdCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]ThresholdCustomAlertRuleClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalThresholdCustomAlertRuleClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalThresholdCustomAlertRuleClassificationMap(rawMsg json.RawMessage) (map[string]ThresholdCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]ThresholdCustomAlertRuleClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalThresholdCustomAlertRuleClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}

func unmarshalTimeWindowCustomAlertRuleClassification(rawMsg json.RawMessage) (TimeWindowCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b TimeWindowCustomAlertRuleClassification
	switch m["ruleType"] {
	case "ActiveConnectionsNotInAllowedRange":
		b = &ActiveConnectionsNotInAllowedRange{}
	case "AmqpC2DMessagesNotInAllowedRange":
		b = &AmqpC2DMessagesNotInAllowedRange{}
	case "AmqpC2DRejectedMessagesNotInAllowedRange":
		b = &AmqpC2DRejectedMessagesNotInAllowedRange{}
	case "AmqpD2CMessagesNotInAllowedRange":
		b = &AmqpD2CMessagesNotInAllowedRange{}
	case "DirectMethodInvokesNotInAllowedRange":
		b = &DirectMethodInvokesNotInAllowedRange{}
	case "FailedLocalLoginsNotInAllowedRange":
		b = &FailedLocalLoginsNotInAllowedRange{}
	case "FileUploadsNotInAllowedRange":
		b = &FileUploadsNotInAllowedRange{}
	case "HttpC2DMessagesNotInAllowedRange":
		b = &HTTPC2DMessagesNotInAllowedRange{}
	case "HttpC2DRejectedMessagesNotInAllowedRange":
		b = &HTTPC2DRejectedMessagesNotInAllowedRange{}
	case "HttpD2CMessagesNotInAllowedRange":
		b = &HTTPD2CMessagesNotInAllowedRange{}
	case "MqttC2DMessagesNotInAllowedRange":
		b = &MqttC2DMessagesNotInAllowedRange{}
	case "MqttC2DRejectedMessagesNotInAllowedRange":
		b = &MqttC2DRejectedMessagesNotInAllowedRange{}
	case "MqttD2CMessagesNotInAllowedRange":
		b = &MqttD2CMessagesNotInAllowedRange{}
	case "QueuePurgesNotInAllowedRange":
		b = &QueuePurgesNotInAllowedRange{}
	case "TwinUpdatesNotInAllowedRange":
		b = &TwinUpdatesNotInAllowedRange{}
	case "UnauthorizedOperationsNotInAllowedRange":
		b = &UnauthorizedOperationsNotInAllowedRange{}
	default:
		b = &TimeWindowCustomAlertRule{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalTimeWindowCustomAlertRuleClassificationArray(rawMsg json.RawMessage) ([]TimeWindowCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]TimeWindowCustomAlertRuleClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalTimeWindowCustomAlertRuleClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalTimeWindowCustomAlertRuleClassificationMap(rawMsg json.RawMessage) (map[string]TimeWindowCustomAlertRuleClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages map[string]json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fMap := make(map[string]TimeWindowCustomAlertRuleClassification, len(rawMessages))
	for key, rawMessage := range rawMessages {
		f, err := unmarshalTimeWindowCustomAlertRuleClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fMap[key] = f
	}
	return fMap, nil
}
