package elastic

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// CreatedByType enumerates the values for created by type.
type CreatedByType string

const (
	// CreatedByTypeApplication ...
	CreatedByTypeApplication CreatedByType = "Application"
	// CreatedByTypeKey ...
	CreatedByTypeKey CreatedByType = "Key"
	// CreatedByTypeManagedIdentity ...
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	// CreatedByTypeUser ...
	CreatedByTypeUser CreatedByType = "User"
)

// PossibleCreatedByTypeValues returns an array of possible values for the CreatedByType const type.
func PossibleCreatedByTypeValues() []CreatedByType {
	return []CreatedByType{CreatedByTypeApplication, CreatedByTypeKey, CreatedByTypeManagedIdentity, CreatedByTypeUser}
}

// DeploymentStatus enumerates the values for deployment status.
type DeploymentStatus string

const (
	// DeploymentStatusHealthy ...
	DeploymentStatusHealthy DeploymentStatus = "Healthy"
	// DeploymentStatusUnhealthy ...
	DeploymentStatusUnhealthy DeploymentStatus = "Unhealthy"
)

// PossibleDeploymentStatusValues returns an array of possible values for the DeploymentStatus const type.
func PossibleDeploymentStatusValues() []DeploymentStatus {
	return []DeploymentStatus{DeploymentStatusHealthy, DeploymentStatusUnhealthy}
}

// LiftrResourceCategories enumerates the values for liftr resource categories.
type LiftrResourceCategories string

const (
	// LiftrResourceCategoriesMonitorLogs ...
	LiftrResourceCategoriesMonitorLogs LiftrResourceCategories = "MonitorLogs"
	// LiftrResourceCategoriesUnknown ...
	LiftrResourceCategoriesUnknown LiftrResourceCategories = "Unknown"
)

// PossibleLiftrResourceCategoriesValues returns an array of possible values for the LiftrResourceCategories const type.
func PossibleLiftrResourceCategoriesValues() []LiftrResourceCategories {
	return []LiftrResourceCategories{LiftrResourceCategoriesMonitorLogs, LiftrResourceCategoriesUnknown}
}

// ManagedIdentityTypes enumerates the values for managed identity types.
type ManagedIdentityTypes string

const (
	// ManagedIdentityTypesSystemAssigned ...
	ManagedIdentityTypesSystemAssigned ManagedIdentityTypes = "SystemAssigned"
)

// PossibleManagedIdentityTypesValues returns an array of possible values for the ManagedIdentityTypes const type.
func PossibleManagedIdentityTypesValues() []ManagedIdentityTypes {
	return []ManagedIdentityTypes{ManagedIdentityTypesSystemAssigned}
}

// MonitoringStatus enumerates the values for monitoring status.
type MonitoringStatus string

const (
	// MonitoringStatusDisabled ...
	MonitoringStatusDisabled MonitoringStatus = "Disabled"
	// MonitoringStatusEnabled ...
	MonitoringStatusEnabled MonitoringStatus = "Enabled"
)

// PossibleMonitoringStatusValues returns an array of possible values for the MonitoringStatus const type.
func PossibleMonitoringStatusValues() []MonitoringStatus {
	return []MonitoringStatus{MonitoringStatusDisabled, MonitoringStatusEnabled}
}

// OperationName enumerates the values for operation name.
type OperationName string

const (
	// OperationNameAdd ...
	OperationNameAdd OperationName = "Add"
	// OperationNameDelete ...
	OperationNameDelete OperationName = "Delete"
)

// PossibleOperationNameValues returns an array of possible values for the OperationName const type.
func PossibleOperationNameValues() []OperationName {
	return []OperationName{OperationNameAdd, OperationNameDelete}
}

// ProvisioningState enumerates the values for provisioning state.
type ProvisioningState string

const (
	// ProvisioningStateAccepted ...
	ProvisioningStateAccepted ProvisioningState = "Accepted"
	// ProvisioningStateCanceled ...
	ProvisioningStateCanceled ProvisioningState = "Canceled"
	// ProvisioningStateCreating ...
	ProvisioningStateCreating ProvisioningState = "Creating"
	// ProvisioningStateDeleted ...
	ProvisioningStateDeleted ProvisioningState = "Deleted"
	// ProvisioningStateDeleting ...
	ProvisioningStateDeleting ProvisioningState = "Deleting"
	// ProvisioningStateFailed ...
	ProvisioningStateFailed ProvisioningState = "Failed"
	// ProvisioningStateNotSpecified ...
	ProvisioningStateNotSpecified ProvisioningState = "NotSpecified"
	// ProvisioningStateSucceeded ...
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	// ProvisioningStateUpdating ...
	ProvisioningStateUpdating ProvisioningState = "Updating"
)

// PossibleProvisioningStateValues returns an array of possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{ProvisioningStateAccepted, ProvisioningStateCanceled, ProvisioningStateCreating, ProvisioningStateDeleted, ProvisioningStateDeleting, ProvisioningStateFailed, ProvisioningStateNotSpecified, ProvisioningStateSucceeded, ProvisioningStateUpdating}
}

// SendingLogs enumerates the values for sending logs.
type SendingLogs string

const (
	// SendingLogsFalse ...
	SendingLogsFalse SendingLogs = "False"
	// SendingLogsTrue ...
	SendingLogsTrue SendingLogs = "True"
)

// PossibleSendingLogsValues returns an array of possible values for the SendingLogs const type.
func PossibleSendingLogsValues() []SendingLogs {
	return []SendingLogs{SendingLogsFalse, SendingLogsTrue}
}

// TagAction enumerates the values for tag action.
type TagAction string

const (
	// TagActionExclude ...
	TagActionExclude TagAction = "Exclude"
	// TagActionInclude ...
	TagActionInclude TagAction = "Include"
)

// PossibleTagActionValues returns an array of possible values for the TagAction const type.
func PossibleTagActionValues() []TagAction {
	return []TagAction{TagActionExclude, TagActionInclude}
}
