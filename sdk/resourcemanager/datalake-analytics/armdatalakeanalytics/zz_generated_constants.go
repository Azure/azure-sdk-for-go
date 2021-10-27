//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatalakeanalytics

const (
	module  = "armdatalakeanalytics"
	version = "v0.1.1"
)

// AADObjectType - The type of AAD object the object identifier refers to.
type AADObjectType string

const (
	AADObjectTypeGroup            AADObjectType = "Group"
	AADObjectTypeServicePrincipal AADObjectType = "ServicePrincipal"
	AADObjectTypeUser             AADObjectType = "User"
)

// PossibleAADObjectTypeValues returns the possible values for the AADObjectType const type.
func PossibleAADObjectTypeValues() []AADObjectType {
	return []AADObjectType{
		AADObjectTypeGroup,
		AADObjectTypeServicePrincipal,
		AADObjectTypeUser,
	}
}

// ToPtr returns a *AADObjectType pointing to the current value.
func (c AADObjectType) ToPtr() *AADObjectType {
	return &c
}

// CheckNameAvailabilityParametersType - The resource type. Note: This should not be set by the user, as the constant value is Microsoft.DataLakeAnalytics/accounts
type CheckNameAvailabilityParametersType string

const (
	CheckNameAvailabilityParametersTypeMicrosoftDataLakeAnalyticsAccounts CheckNameAvailabilityParametersType = "Microsoft.DataLakeAnalytics/accounts"
)

// PossibleCheckNameAvailabilityParametersTypeValues returns the possible values for the CheckNameAvailabilityParametersType const type.
func PossibleCheckNameAvailabilityParametersTypeValues() []CheckNameAvailabilityParametersType {
	return []CheckNameAvailabilityParametersType{
		CheckNameAvailabilityParametersTypeMicrosoftDataLakeAnalyticsAccounts,
	}
}

// ToPtr returns a *CheckNameAvailabilityParametersType pointing to the current value.
func (c CheckNameAvailabilityParametersType) ToPtr() *CheckNameAvailabilityParametersType {
	return &c
}

// DataLakeAnalyticsAccountState - The state of the Data Lake Analytics account.
type DataLakeAnalyticsAccountState string

const (
	DataLakeAnalyticsAccountStateActive    DataLakeAnalyticsAccountState = "Active"
	DataLakeAnalyticsAccountStateSuspended DataLakeAnalyticsAccountState = "Suspended"
)

// PossibleDataLakeAnalyticsAccountStateValues returns the possible values for the DataLakeAnalyticsAccountState const type.
func PossibleDataLakeAnalyticsAccountStateValues() []DataLakeAnalyticsAccountState {
	return []DataLakeAnalyticsAccountState{
		DataLakeAnalyticsAccountStateActive,
		DataLakeAnalyticsAccountStateSuspended,
	}
}

// ToPtr returns a *DataLakeAnalyticsAccountState pointing to the current value.
func (c DataLakeAnalyticsAccountState) ToPtr() *DataLakeAnalyticsAccountState {
	return &c
}

// DataLakeAnalyticsAccountStatus - The provisioning status of the Data Lake Analytics account.
type DataLakeAnalyticsAccountStatus string

const (
	DataLakeAnalyticsAccountStatusFailed     DataLakeAnalyticsAccountStatus = "Failed"
	DataLakeAnalyticsAccountStatusCreating   DataLakeAnalyticsAccountStatus = "Creating"
	DataLakeAnalyticsAccountStatusRunning    DataLakeAnalyticsAccountStatus = "Running"
	DataLakeAnalyticsAccountStatusSucceeded  DataLakeAnalyticsAccountStatus = "Succeeded"
	DataLakeAnalyticsAccountStatusPatching   DataLakeAnalyticsAccountStatus = "Patching"
	DataLakeAnalyticsAccountStatusSuspending DataLakeAnalyticsAccountStatus = "Suspending"
	DataLakeAnalyticsAccountStatusResuming   DataLakeAnalyticsAccountStatus = "Resuming"
	DataLakeAnalyticsAccountStatusDeleting   DataLakeAnalyticsAccountStatus = "Deleting"
	DataLakeAnalyticsAccountStatusDeleted    DataLakeAnalyticsAccountStatus = "Deleted"
	DataLakeAnalyticsAccountStatusUndeleting DataLakeAnalyticsAccountStatus = "Undeleting"
	DataLakeAnalyticsAccountStatusCanceled   DataLakeAnalyticsAccountStatus = "Canceled"
)

// PossibleDataLakeAnalyticsAccountStatusValues returns the possible values for the DataLakeAnalyticsAccountStatus const type.
func PossibleDataLakeAnalyticsAccountStatusValues() []DataLakeAnalyticsAccountStatus {
	return []DataLakeAnalyticsAccountStatus{
		DataLakeAnalyticsAccountStatusFailed,
		DataLakeAnalyticsAccountStatusCreating,
		DataLakeAnalyticsAccountStatusRunning,
		DataLakeAnalyticsAccountStatusSucceeded,
		DataLakeAnalyticsAccountStatusPatching,
		DataLakeAnalyticsAccountStatusSuspending,
		DataLakeAnalyticsAccountStatusResuming,
		DataLakeAnalyticsAccountStatusDeleting,
		DataLakeAnalyticsAccountStatusDeleted,
		DataLakeAnalyticsAccountStatusUndeleting,
		DataLakeAnalyticsAccountStatusCanceled,
	}
}

// ToPtr returns a *DataLakeAnalyticsAccountStatus pointing to the current value.
func (c DataLakeAnalyticsAccountStatus) ToPtr() *DataLakeAnalyticsAccountStatus {
	return &c
}

// DebugDataAccessLevel - The current state of the DebugDataAccessLevel for this account.
type DebugDataAccessLevel string

const (
	DebugDataAccessLevelAll      DebugDataAccessLevel = "All"
	DebugDataAccessLevelCustomer DebugDataAccessLevel = "Customer"
	DebugDataAccessLevelNone     DebugDataAccessLevel = "None"
)

// PossibleDebugDataAccessLevelValues returns the possible values for the DebugDataAccessLevel const type.
func PossibleDebugDataAccessLevelValues() []DebugDataAccessLevel {
	return []DebugDataAccessLevel{
		DebugDataAccessLevelAll,
		DebugDataAccessLevelCustomer,
		DebugDataAccessLevelNone,
	}
}

// ToPtr returns a *DebugDataAccessLevel pointing to the current value.
func (c DebugDataAccessLevel) ToPtr() *DebugDataAccessLevel {
	return &c
}

// FirewallAllowAzureIPsState - The current state of allowing or disallowing IPs originating within Azure through the firewall. If the firewall is disabled,
// this is not enforced.
type FirewallAllowAzureIPsState string

const (
	FirewallAllowAzureIPsStateEnabled  FirewallAllowAzureIPsState = "Enabled"
	FirewallAllowAzureIPsStateDisabled FirewallAllowAzureIPsState = "Disabled"
)

// PossibleFirewallAllowAzureIPsStateValues returns the possible values for the FirewallAllowAzureIPsState const type.
func PossibleFirewallAllowAzureIPsStateValues() []FirewallAllowAzureIPsState {
	return []FirewallAllowAzureIPsState{
		FirewallAllowAzureIPsStateEnabled,
		FirewallAllowAzureIPsStateDisabled,
	}
}

// ToPtr returns a *FirewallAllowAzureIPsState pointing to the current value.
func (c FirewallAllowAzureIPsState) ToPtr() *FirewallAllowAzureIPsState {
	return &c
}

// FirewallState - The current state of the IP address firewall for this account.
type FirewallState string

const (
	FirewallStateEnabled  FirewallState = "Enabled"
	FirewallStateDisabled FirewallState = "Disabled"
)

// PossibleFirewallStateValues returns the possible values for the FirewallState const type.
func PossibleFirewallStateValues() []FirewallState {
	return []FirewallState{
		FirewallStateEnabled,
		FirewallStateDisabled,
	}
}

// ToPtr returns a *FirewallState pointing to the current value.
func (c FirewallState) ToPtr() *FirewallState {
	return &c
}

// NestedResourceProvisioningState - The current state of the NestedResourceProvisioning for this account.
type NestedResourceProvisioningState string

const (
	NestedResourceProvisioningStateSucceeded NestedResourceProvisioningState = "Succeeded"
	NestedResourceProvisioningStateCanceled  NestedResourceProvisioningState = "Canceled"
	NestedResourceProvisioningStateFailed    NestedResourceProvisioningState = "Failed"
)

// PossibleNestedResourceProvisioningStateValues returns the possible values for the NestedResourceProvisioningState const type.
func PossibleNestedResourceProvisioningStateValues() []NestedResourceProvisioningState {
	return []NestedResourceProvisioningState{
		NestedResourceProvisioningStateSucceeded,
		NestedResourceProvisioningStateCanceled,
		NestedResourceProvisioningStateFailed,
	}
}

// ToPtr returns a *NestedResourceProvisioningState pointing to the current value.
func (c NestedResourceProvisioningState) ToPtr() *NestedResourceProvisioningState {
	return &c
}

// OperationOrigin - The intended executor of the operation.
type OperationOrigin string

const (
	OperationOriginSystem     OperationOrigin = "system"
	OperationOriginUser       OperationOrigin = "user"
	OperationOriginUserSystem OperationOrigin = "user,system"
)

// PossibleOperationOriginValues returns the possible values for the OperationOrigin const type.
func PossibleOperationOriginValues() []OperationOrigin {
	return []OperationOrigin{
		OperationOriginSystem,
		OperationOriginUser,
		OperationOriginUserSystem,
	}
}

// ToPtr returns a *OperationOrigin pointing to the current value.
func (c OperationOrigin) ToPtr() *OperationOrigin {
	return &c
}

// SubscriptionState - The subscription state.
type SubscriptionState string

const (
	SubscriptionStateDeleted      SubscriptionState = "Deleted"
	SubscriptionStateRegistered   SubscriptionState = "Registered"
	SubscriptionStateSuspended    SubscriptionState = "Suspended"
	SubscriptionStateUnregistered SubscriptionState = "Unregistered"
	SubscriptionStateWarned       SubscriptionState = "Warned"
)

// PossibleSubscriptionStateValues returns the possible values for the SubscriptionState const type.
func PossibleSubscriptionStateValues() []SubscriptionState {
	return []SubscriptionState{
		SubscriptionStateDeleted,
		SubscriptionStateRegistered,
		SubscriptionStateSuspended,
		SubscriptionStateUnregistered,
		SubscriptionStateWarned,
	}
}

// ToPtr returns a *SubscriptionState pointing to the current value.
func (c SubscriptionState) ToPtr() *SubscriptionState {
	return &c
}

// TierType - The commitment tier for the next month.
type TierType string

const (
	TierTypeConsumption             TierType = "Consumption"
	TierTypeCommitment100AUHours    TierType = "Commitment_100AUHours"
	TierTypeCommitment500AUHours    TierType = "Commitment_500AUHours"
	TierTypeCommitment1000AUHours   TierType = "Commitment_1000AUHours"
	TierTypeCommitment5000AUHours   TierType = "Commitment_5000AUHours"
	TierTypeCommitment10000AUHours  TierType = "Commitment_10000AUHours"
	TierTypeCommitment50000AUHours  TierType = "Commitment_50000AUHours"
	TierTypeCommitment100000AUHours TierType = "Commitment_100000AUHours"
	TierTypeCommitment500000AUHours TierType = "Commitment_500000AUHours"
)

// PossibleTierTypeValues returns the possible values for the TierType const type.
func PossibleTierTypeValues() []TierType {
	return []TierType{
		TierTypeConsumption,
		TierTypeCommitment100AUHours,
		TierTypeCommitment500AUHours,
		TierTypeCommitment1000AUHours,
		TierTypeCommitment5000AUHours,
		TierTypeCommitment10000AUHours,
		TierTypeCommitment50000AUHours,
		TierTypeCommitment100000AUHours,
		TierTypeCommitment500000AUHours,
	}
}

// ToPtr returns a *TierType pointing to the current value.
func (c TierType) ToPtr() *TierType {
	return &c
}

// VirtualNetworkRuleState - The current state of the VirtualNetworkRule for this account.
type VirtualNetworkRuleState string

const (
	VirtualNetworkRuleStateActive               VirtualNetworkRuleState = "Active"
	VirtualNetworkRuleStateNetworkSourceDeleted VirtualNetworkRuleState = "NetworkSourceDeleted"
	VirtualNetworkRuleStateFailed               VirtualNetworkRuleState = "Failed"
)

// PossibleVirtualNetworkRuleStateValues returns the possible values for the VirtualNetworkRuleState const type.
func PossibleVirtualNetworkRuleStateValues() []VirtualNetworkRuleState {
	return []VirtualNetworkRuleState{
		VirtualNetworkRuleStateActive,
		VirtualNetworkRuleStateNetworkSourceDeleted,
		VirtualNetworkRuleStateFailed,
	}
}

// ToPtr returns a *VirtualNetworkRuleState pointing to the current value.
func (c VirtualNetworkRuleState) ToPtr() *VirtualNetworkRuleState {
	return &c
}
