//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

// AdditionalDataClassification provides polymorphic access to related types.
// Call the interface's GetAdditionalData() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AdditionalData, *ContainerRegistryVulnerabilityProperties, *SQLServerVulnerabilityProperties, *ServerVulnerabilityProperties
type AdditionalDataClassification interface {
	// GetAdditionalData returns the AdditionalData content of the underlying type.
	GetAdditionalData() *AdditionalData
}

// AlertSimulatorRequestPropertiesClassification provides polymorphic access to related types.
// Call the interface's GetAlertSimulatorRequestProperties() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AlertSimulatorBundlesRequestProperties, *AlertSimulatorRequestProperties
type AlertSimulatorRequestPropertiesClassification interface {
	// GetAlertSimulatorRequestProperties returns the AlertSimulatorRequestProperties content of the underlying type.
	GetAlertSimulatorRequestProperties() *AlertSimulatorRequestProperties
}

// AllowlistCustomAlertRuleClassification provides polymorphic access to related types.
// Call the interface's GetAllowlistCustomAlertRule() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AllowlistCustomAlertRule, *ConnectionFromIPNotAllowed, *ConnectionToIPNotAllowed, *LocalUserNotAllowed, *ProcessNotAllowed
type AllowlistCustomAlertRuleClassification interface {
	ListCustomAlertRuleClassification
	// GetAllowlistCustomAlertRule returns the AllowlistCustomAlertRule content of the underlying type.
	GetAllowlistCustomAlertRule() *AllowlistCustomAlertRule
}

// AuthenticationDetailsPropertiesClassification provides polymorphic access to related types.
// Call the interface's GetAuthenticationDetailsProperties() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AuthenticationDetailsProperties, *AwAssumeRoleAuthenticationDetailsProperties, *AwsCredsAuthenticationDetailsProperties,
// - *GcpCredentialsDetailsProperties
type AuthenticationDetailsPropertiesClassification interface {
	// GetAuthenticationDetailsProperties returns the AuthenticationDetailsProperties content of the underlying type.
	GetAuthenticationDetailsProperties() *AuthenticationDetailsProperties
}

// AutomationActionClassification provides polymorphic access to related types.
// Call the interface's GetAutomationAction() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AutomationAction, *AutomationActionEventHub, *AutomationActionLogicApp, *AutomationActionWorkspace
type AutomationActionClassification interface {
	// GetAutomationAction returns the AutomationAction content of the underlying type.
	GetAutomationAction() *AutomationAction
}

// AwsOrganizationalDataClassification provides polymorphic access to related types.
// Call the interface's GetAwsOrganizationalData() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AwsOrganizationalData, *AwsOrganizationalDataMaster, *AwsOrganizationalDataMember
type AwsOrganizationalDataClassification interface {
	// GetAwsOrganizationalData returns the AwsOrganizationalData content of the underlying type.
	GetAwsOrganizationalData() *AwsOrganizationalData
}

// CloudOfferingClassification provides polymorphic access to related types.
// Call the interface's GetCloudOffering() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *CloudOffering, *CspmMonitorAwsOffering, *CspmMonitorAzureDevOpsOffering, *CspmMonitorGcpOffering, *CspmMonitorGitLabOffering,
// - *CspmMonitorGithubOffering, *DefenderCspmAwsOffering, *DefenderCspmGcpOffering, *DefenderFoDatabasesAwsOffering, *DefenderForContainersAwsOffering,
// - *DefenderForContainersGcpOffering, *DefenderForDatabasesGcpOffering, *DefenderForDevOpsAzureDevOpsOffering, *DefenderForDevOpsGitLabOffering,
// - *DefenderForDevOpsGithubOffering, *DefenderForServersAwsOffering, *DefenderForServersGcpOffering, *InformationProtectionAwsOffering
type CloudOfferingClassification interface {
	// GetCloudOffering returns the CloudOffering content of the underlying type.
	GetCloudOffering() *CloudOffering
}

// CustomAlertRuleClassification provides polymorphic access to related types.
// Call the interface's GetCustomAlertRule() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ActiveConnectionsNotInAllowedRange, *AllowlistCustomAlertRule, *AmqpC2DMessagesNotInAllowedRange, *AmqpC2DRejectedMessagesNotInAllowedRange,
// - *AmqpD2CMessagesNotInAllowedRange, *ConnectionFromIPNotAllowed, *ConnectionToIPNotAllowed, *CustomAlertRule, *DenylistCustomAlertRule,
// - *DirectMethodInvokesNotInAllowedRange, *FailedLocalLoginsNotInAllowedRange, *FileUploadsNotInAllowedRange, *HTTPC2DMessagesNotInAllowedRange,
// - *HTTPC2DRejectedMessagesNotInAllowedRange, *HTTPD2CMessagesNotInAllowedRange, *ListCustomAlertRule, *LocalUserNotAllowed,
// - *MqttC2DMessagesNotInAllowedRange, *MqttC2DRejectedMessagesNotInAllowedRange, *MqttD2CMessagesNotInAllowedRange, *ProcessNotAllowed,
// - *QueuePurgesNotInAllowedRange, *ThresholdCustomAlertRule, *TimeWindowCustomAlertRule, *TwinUpdatesNotInAllowedRange,
// - *UnauthorizedOperationsNotInAllowedRange
type CustomAlertRuleClassification interface {
	// GetCustomAlertRule returns the CustomAlertRule content of the underlying type.
	GetCustomAlertRule() *CustomAlertRule
}

// EnvironmentDataClassification provides polymorphic access to related types.
// Call the interface's GetEnvironmentData() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AwsEnvironmentData, *AzureDevOpsScopeEnvironmentData, *EnvironmentData, *GcpProjectEnvironmentData, *GithubScopeEnvironmentData,
// - *GitlabScopeEnvironmentData
type EnvironmentDataClassification interface {
	// GetEnvironmentData returns the EnvironmentData content of the underlying type.
	GetEnvironmentData() *EnvironmentData
}

// GcpOrganizationalDataClassification provides polymorphic access to related types.
// Call the interface's GetGcpOrganizationalData() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *GcpOrganizationalData, *GcpOrganizationalDataMember, *GcpOrganizationalDataOrganization
type GcpOrganizationalDataClassification interface {
	// GetGcpOrganizationalData returns the GcpOrganizationalData content of the underlying type.
	GetGcpOrganizationalData() *GcpOrganizationalData
}

// ListCustomAlertRuleClassification provides polymorphic access to related types.
// Call the interface's GetListCustomAlertRule() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AllowlistCustomAlertRule, *ConnectionFromIPNotAllowed, *ConnectionToIPNotAllowed, *DenylistCustomAlertRule, *ListCustomAlertRule,
// - *LocalUserNotAllowed, *ProcessNotAllowed
type ListCustomAlertRuleClassification interface {
	CustomAlertRuleClassification
	// GetListCustomAlertRule returns the ListCustomAlertRule content of the underlying type.
	GetListCustomAlertRule() *ListCustomAlertRule
}

// ListCustomAlertRuleClassification provides polymorphic access to related types.
// Call the interface's GetListCustomAlertRule() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AllowlistCustomAlertRule, *ConnectionFromIPNotAllowed, *ConnectionToIPNotAllowed, *DenylistCustomAlertRule, *ListCustomAlertRule,
// - *LocalUserNotAllowed, *ProcessNotAllowed
type ListCustomAlertRuleClassification interface {
	CustomAlertRuleClassification
	// GetListCustomAlertRule returns the ListCustomAlertRule content of the underlying type.
	GetListCustomAlertRule() *ListCustomAlertRule
}

// OnPremiseResourceDetailsClassification provides polymorphic access to related types.
// Call the interface's GetOnPremiseResourceDetails() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *OnPremiseResourceDetails, *OnPremiseSQLResourceDetails
type OnPremiseResourceDetailsClassification interface {
	ResourceDetailsClassification
	// GetOnPremiseResourceDetails returns the OnPremiseResourceDetails content of the underlying type.
	GetOnPremiseResourceDetails() *OnPremiseResourceDetails
}

// ResourceDetailsClassification provides polymorphic access to related types.
// Call the interface's GetResourceDetails() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AzureResourceDetails, *OnPremiseResourceDetails, *OnPremiseSQLResourceDetails, *ResourceDetails
type ResourceDetailsClassification interface {
	// GetResourceDetails returns the ResourceDetails content of the underlying type.
	GetResourceDetails() *ResourceDetails
}

// ResourceIdentifierClassification provides polymorphic access to related types.
// Call the interface's GetResourceIdentifier() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AzureResourceIdentifier, *LogAnalyticsIdentifier, *ResourceIdentifier
type ResourceIdentifierClassification interface {
	// GetResourceIdentifier returns the ResourceIdentifier content of the underlying type.
	GetResourceIdentifier() *ResourceIdentifier
}

// SettingClassification provides polymorphic access to related types.
// Call the interface's GetSetting() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AlertSyncSettings, *DataExportSettings, *Setting
type SettingClassification interface {
	// GetSetting returns the Setting content of the underlying type.
	GetSetting() *Setting
}

// ThresholdCustomAlertRuleClassification provides polymorphic access to related types.
// Call the interface's GetThresholdCustomAlertRule() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ActiveConnectionsNotInAllowedRange, *AmqpC2DMessagesNotInAllowedRange, *AmqpC2DRejectedMessagesNotInAllowedRange, *AmqpD2CMessagesNotInAllowedRange,
// - *DirectMethodInvokesNotInAllowedRange, *FailedLocalLoginsNotInAllowedRange, *FileUploadsNotInAllowedRange, *HTTPC2DMessagesNotInAllowedRange,
// - *HTTPC2DRejectedMessagesNotInAllowedRange, *HTTPD2CMessagesNotInAllowedRange, *MqttC2DMessagesNotInAllowedRange, *MqttC2DRejectedMessagesNotInAllowedRange,
// - *MqttD2CMessagesNotInAllowedRange, *QueuePurgesNotInAllowedRange, *ThresholdCustomAlertRule, *TimeWindowCustomAlertRule,
// - *TwinUpdatesNotInAllowedRange, *UnauthorizedOperationsNotInAllowedRange
type ThresholdCustomAlertRuleClassification interface {
	CustomAlertRuleClassification
	// GetThresholdCustomAlertRule returns the ThresholdCustomAlertRule content of the underlying type.
	GetThresholdCustomAlertRule() *ThresholdCustomAlertRule
}

// ThresholdCustomAlertRuleClassification provides polymorphic access to related types.
// Call the interface's GetThresholdCustomAlertRule() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ActiveConnectionsNotInAllowedRange, *AmqpC2DMessagesNotInAllowedRange, *AmqpC2DRejectedMessagesNotInAllowedRange, *AmqpD2CMessagesNotInAllowedRange,
// - *DirectMethodInvokesNotInAllowedRange, *FailedLocalLoginsNotInAllowedRange, *FileUploadsNotInAllowedRange, *HTTPC2DMessagesNotInAllowedRange,
// - *HTTPC2DRejectedMessagesNotInAllowedRange, *HTTPD2CMessagesNotInAllowedRange, *MqttC2DMessagesNotInAllowedRange, *MqttC2DRejectedMessagesNotInAllowedRange,
// - *MqttD2CMessagesNotInAllowedRange, *QueuePurgesNotInAllowedRange, *ThresholdCustomAlertRule, *TimeWindowCustomAlertRule,
// - *TwinUpdatesNotInAllowedRange, *UnauthorizedOperationsNotInAllowedRange
type ThresholdCustomAlertRuleClassification interface {
	CustomAlertRuleClassification
	// GetThresholdCustomAlertRule returns the ThresholdCustomAlertRule content of the underlying type.
	GetThresholdCustomAlertRule() *ThresholdCustomAlertRule
}

// TimeWindowCustomAlertRuleClassification provides polymorphic access to related types.
// Call the interface's GetTimeWindowCustomAlertRule() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ActiveConnectionsNotInAllowedRange, *AmqpC2DMessagesNotInAllowedRange, *AmqpC2DRejectedMessagesNotInAllowedRange, *AmqpD2CMessagesNotInAllowedRange,
// - *DirectMethodInvokesNotInAllowedRange, *FailedLocalLoginsNotInAllowedRange, *FileUploadsNotInAllowedRange, *HTTPC2DMessagesNotInAllowedRange,
// - *HTTPC2DRejectedMessagesNotInAllowedRange, *HTTPD2CMessagesNotInAllowedRange, *MqttC2DMessagesNotInAllowedRange, *MqttC2DRejectedMessagesNotInAllowedRange,
// - *MqttD2CMessagesNotInAllowedRange, *QueuePurgesNotInAllowedRange, *TimeWindowCustomAlertRule, *TwinUpdatesNotInAllowedRange,
// - *UnauthorizedOperationsNotInAllowedRange
type TimeWindowCustomAlertRuleClassification interface {
	ThresholdCustomAlertRuleClassification
	// GetTimeWindowCustomAlertRule returns the TimeWindowCustomAlertRule content of the underlying type.
	GetTimeWindowCustomAlertRule() *TimeWindowCustomAlertRule
}

