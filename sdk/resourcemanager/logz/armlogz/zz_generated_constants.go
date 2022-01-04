//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armlogz

const (
	module  = "armlogz"
	version = "v0.1.0"
)

// CreatedByType - The type of identity that created the resource.
type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

// PossibleCreatedByTypeValues returns the possible values for the CreatedByType const type.
func PossibleCreatedByTypeValues() []CreatedByType {
	return []CreatedByType{
		CreatedByTypeApplication,
		CreatedByTypeKey,
		CreatedByTypeManagedIdentity,
		CreatedByTypeUser,
	}
}

// ToPtr returns a *CreatedByType pointing to the current value.
func (c CreatedByType) ToPtr() *CreatedByType {
	return &c
}

type LiftrResourceCategories string

const (
	LiftrResourceCategoriesMonitorLogs LiftrResourceCategories = "MonitorLogs"
	LiftrResourceCategoriesUnknown     LiftrResourceCategories = "Unknown"
)

// PossibleLiftrResourceCategoriesValues returns the possible values for the LiftrResourceCategories const type.
func PossibleLiftrResourceCategoriesValues() []LiftrResourceCategories {
	return []LiftrResourceCategories{
		LiftrResourceCategoriesMonitorLogs,
		LiftrResourceCategoriesUnknown,
	}
}

// ToPtr returns a *LiftrResourceCategories pointing to the current value.
func (c LiftrResourceCategories) ToPtr() *LiftrResourceCategories {
	return &c
}

type ManagedIdentityTypes string

const (
	ManagedIdentityTypesSystemAssigned ManagedIdentityTypes = "SystemAssigned"
	ManagedIdentityTypesUserAssigned   ManagedIdentityTypes = "UserAssigned"
)

// PossibleManagedIdentityTypesValues returns the possible values for the ManagedIdentityTypes const type.
func PossibleManagedIdentityTypesValues() []ManagedIdentityTypes {
	return []ManagedIdentityTypes{
		ManagedIdentityTypesSystemAssigned,
		ManagedIdentityTypesUserAssigned,
	}
}

// ToPtr returns a *ManagedIdentityTypes pointing to the current value.
func (c ManagedIdentityTypes) ToPtr() *ManagedIdentityTypes {
	return &c
}

// MarketplaceSubscriptionStatus - Flag specifying the Marketplace Subscription Status of the resource. If payment is not made in time, the resource will
// go in Suspended state.
type MarketplaceSubscriptionStatus string

const (
	MarketplaceSubscriptionStatusActive    MarketplaceSubscriptionStatus = "Active"
	MarketplaceSubscriptionStatusSuspended MarketplaceSubscriptionStatus = "Suspended"
)

// PossibleMarketplaceSubscriptionStatusValues returns the possible values for the MarketplaceSubscriptionStatus const type.
func PossibleMarketplaceSubscriptionStatusValues() []MarketplaceSubscriptionStatus {
	return []MarketplaceSubscriptionStatus{
		MarketplaceSubscriptionStatusActive,
		MarketplaceSubscriptionStatusSuspended,
	}
}

// ToPtr returns a *MarketplaceSubscriptionStatus pointing to the current value.
func (c MarketplaceSubscriptionStatus) ToPtr() *MarketplaceSubscriptionStatus {
	return &c
}

// MonitoringStatus - Flag specifying if the resource monitoring is enabled or disabled.
type MonitoringStatus string

const (
	MonitoringStatusDisabled MonitoringStatus = "Disabled"
	MonitoringStatusEnabled  MonitoringStatus = "Enabled"
)

// PossibleMonitoringStatusValues returns the possible values for the MonitoringStatus const type.
func PossibleMonitoringStatusValues() []MonitoringStatus {
	return []MonitoringStatus{
		MonitoringStatusDisabled,
		MonitoringStatusEnabled,
	}
}

// ToPtr returns a *MonitoringStatus pointing to the current value.
func (c MonitoringStatus) ToPtr() *MonitoringStatus {
	return &c
}

// ProvisioningState - Flag specifying if the resource provisioning state as tracked by ARM.
type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleted      ProvisioningState = "Deleted"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateNotSpecified ProvisioningState = "NotSpecified"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

// PossibleProvisioningStateValues returns the possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{
		ProvisioningStateAccepted,
		ProvisioningStateCanceled,
		ProvisioningStateCreating,
		ProvisioningStateDeleted,
		ProvisioningStateDeleting,
		ProvisioningStateFailed,
		ProvisioningStateNotSpecified,
		ProvisioningStateSucceeded,
		ProvisioningStateUpdating,
	}
}

// ToPtr returns a *ProvisioningState pointing to the current value.
func (c ProvisioningState) ToPtr() *ProvisioningState {
	return &c
}

// SingleSignOnStates - Various states of the SSO resource
type SingleSignOnStates string

const (
	SingleSignOnStatesDisable  SingleSignOnStates = "Disable"
	SingleSignOnStatesEnable   SingleSignOnStates = "Enable"
	SingleSignOnStatesExisting SingleSignOnStates = "Existing"
	SingleSignOnStatesInitial  SingleSignOnStates = "Initial"
)

// PossibleSingleSignOnStatesValues returns the possible values for the SingleSignOnStates const type.
func PossibleSingleSignOnStatesValues() []SingleSignOnStates {
	return []SingleSignOnStates{
		SingleSignOnStatesDisable,
		SingleSignOnStatesEnable,
		SingleSignOnStatesExisting,
		SingleSignOnStatesInitial,
	}
}

// ToPtr returns a *SingleSignOnStates pointing to the current value.
func (c SingleSignOnStates) ToPtr() *SingleSignOnStates {
	return &c
}

// TagAction - Valid actions for a filtering tag. Exclusion takes priority over inclusion.
type TagAction string

const (
	TagActionExclude TagAction = "Exclude"
	TagActionInclude TagAction = "Include"
)

// PossibleTagActionValues returns the possible values for the TagAction const type.
func PossibleTagActionValues() []TagAction {
	return []TagAction{
		TagActionExclude,
		TagActionInclude,
	}
}

// ToPtr returns a *TagAction pointing to the current value.
func (c TagAction) ToPtr() *TagAction {
	return &c
}

// UserRole - User roles on configured in Logz.io account.
type UserRole string

const (
	UserRoleAdmin UserRole = "Admin"
	UserRoleNone  UserRole = "None"
	UserRoleUser  UserRole = "User"
)

// PossibleUserRoleValues returns the possible values for the UserRole const type.
func PossibleUserRoleValues() []UserRole {
	return []UserRole{
		UserRoleAdmin,
		UserRoleNone,
		UserRoleUser,
	}
}

// ToPtr returns a *UserRole pointing to the current value.
func (c UserRole) ToPtr() *UserRole {
	return &c
}

// VMHostUpdateStates - Various states of the updating vm extension on resource
type VMHostUpdateStates string

const (
	VMHostUpdateStatesDelete  VMHostUpdateStates = "Delete"
	VMHostUpdateStatesInstall VMHostUpdateStates = "Install"
)

// PossibleVMHostUpdateStatesValues returns the possible values for the VMHostUpdateStates const type.
func PossibleVMHostUpdateStatesValues() []VMHostUpdateStates {
	return []VMHostUpdateStates{
		VMHostUpdateStatesDelete,
		VMHostUpdateStatesInstall,
	}
}

// ToPtr returns a *VMHostUpdateStates pointing to the current value.
func (c VMHostUpdateStates) ToPtr() *VMHostUpdateStates {
	return &c
}
