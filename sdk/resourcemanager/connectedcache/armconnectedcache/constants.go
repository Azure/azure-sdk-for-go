// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armconnectedcache

const (
	moduleName    = "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/connectedcache/armconnectedcache"
	moduleVersion = "v0.1.0"
)

// ActionType - Extensible enum. Indicates the action type. "Internal" refers to actions that are for internal only APIs.
type ActionType string

const (
	// ActionTypeInternal - Actions are for internal-only APIs.
	ActionTypeInternal ActionType = "Internal"
)

// PossibleActionTypeValues returns the possible values for the ActionType const type.
func PossibleActionTypeValues() []ActionType {
	return []ActionType{
		ActionTypeInternal,
	}
}

// AutoUpdateRingType - Auto update Ring type
type AutoUpdateRingType string

const (
	// AutoUpdateRingTypeFast - customer selection of fast / auto update to install mcc on their physical vm
	AutoUpdateRingTypeFast AutoUpdateRingType = "Fast"
	// AutoUpdateRingTypePreview - customer selection of preview update install mcc on their physical vm
	AutoUpdateRingTypePreview AutoUpdateRingType = "Preview"
	// AutoUpdateRingTypeSlow - customer selection of slow update to install mcc on their physical vm
	AutoUpdateRingTypeSlow AutoUpdateRingType = "Slow"
)

// PossibleAutoUpdateRingTypeValues returns the possible values for the AutoUpdateRingType const type.
func PossibleAutoUpdateRingTypeValues() []AutoUpdateRingType {
	return []AutoUpdateRingType{
		AutoUpdateRingTypeFast,
		AutoUpdateRingTypePreview,
		AutoUpdateRingTypeSlow,
	}
}

// BgpReviewStateEnum - Cache node resource Bgp review state as integer
type BgpReviewStateEnum string

const (
	// BgpReviewStateEnumApproved - bgp is in Approved state
	BgpReviewStateEnumApproved BgpReviewStateEnum = "Approved"
	// BgpReviewStateEnumAttentionRequired - bgp is setup need an attention for more troubleshoot
	BgpReviewStateEnumAttentionRequired BgpReviewStateEnum = "AttentionRequired"
	// BgpReviewStateEnumInReview - bgp is in review state
	BgpReviewStateEnumInReview BgpReviewStateEnum = "InReview"
	// BgpReviewStateEnumNotConfigured - bgp not configured
	BgpReviewStateEnumNotConfigured BgpReviewStateEnum = "NotConfigured"
)

// PossibleBgpReviewStateEnumValues returns the possible values for the BgpReviewStateEnum const type.
func PossibleBgpReviewStateEnumValues() []BgpReviewStateEnum {
	return []BgpReviewStateEnum{
		BgpReviewStateEnumApproved,
		BgpReviewStateEnumAttentionRequired,
		BgpReviewStateEnumInReview,
		BgpReviewStateEnumNotConfigured,
	}
}

// ConfigurationState - Cache node configuration setup state
type ConfigurationState string

const (
	// ConfigurationStateConfigured - connected cache setup configured
	ConfigurationStateConfigured ConfigurationState = "Configured"
	// ConfigurationStateNotConfiguredIP - connected cache setup not configured
	ConfigurationStateNotConfiguredIP ConfigurationState = "NotConfigured_Ip"
)

// PossibleConfigurationStateValues returns the possible values for the ConfigurationState const type.
func PossibleConfigurationStateValues() []ConfigurationState {
	return []ConfigurationState{
		ConfigurationStateConfigured,
		ConfigurationStateNotConfiguredIP,
	}
}

// CreatedByType - The kind of entity that created the resource.
type CreatedByType string

const (
	// CreatedByTypeApplication - The entity was created by an application.
	CreatedByTypeApplication CreatedByType = "Application"
	// CreatedByTypeKey - The entity was created by a key.
	CreatedByTypeKey CreatedByType = "Key"
	// CreatedByTypeManagedIdentity - The entity was created by a managed identity.
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	// CreatedByTypeUser - The entity was created by a user.
	CreatedByTypeUser CreatedByType = "User"
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

// CustomerTransitState - Customer resource transit states
type CustomerTransitState string

const (
	// CustomerTransitStateCombinedTransit - transit provider and have own subscribers
	CustomerTransitStateCombinedTransit CustomerTransitState = "CombinedTransit"
	// CustomerTransitStateNoTransit - do not have transit
	CustomerTransitStateNoTransit CustomerTransitState = "NoTransit"
	// CustomerTransitStateTransitOnly - pure transit provider or network service provider
	CustomerTransitStateTransitOnly CustomerTransitState = "TransitOnly"
)

// PossibleCustomerTransitStateValues returns the possible values for the CustomerTransitState const type.
func PossibleCustomerTransitStateValues() []CustomerTransitState {
	return []CustomerTransitState{
		CustomerTransitStateCombinedTransit,
		CustomerTransitStateNoTransit,
		CustomerTransitStateTransitOnly,
	}
}

// CycleType - Update Cycle type
type CycleType string

const (
	// CycleTypeFast - customer selection of fast / auto update to install mcc on their physical vm
	CycleTypeFast CycleType = "Fast"
	// CycleTypePreview - customer selection of preview update install mcc on their physical vm
	CycleTypePreview CycleType = "Preview"
	// CycleTypeSlow - customer selection of slow update to install mcc on their physical vm
	CycleTypeSlow CycleType = "Slow"
)

// PossibleCycleTypeValues returns the possible values for the CycleType const type.
func PossibleCycleTypeValues() []CycleType {
	return []CycleType{
		CycleTypeFast,
		CycleTypePreview,
		CycleTypeSlow,
	}
}

// Origin - The intended executor of the operation; as in Resource Based Access Control (RBAC) and audit logs UX. Default
// value is "user,system"
type Origin string

const (
	// OriginSystem - Indicates the operation is initiated by a system.
	OriginSystem Origin = "system"
	// OriginUser - Indicates the operation is initiated by a user.
	OriginUser Origin = "user"
	// OriginUserSystem - Indicates the operation is initiated by a user or system.
	OriginUserSystem Origin = "user,system"
)

// PossibleOriginValues returns the possible values for the Origin const type.
func PossibleOriginValues() []Origin {
	return []Origin{
		OriginSystem,
		OriginUser,
		OriginUserSystem,
	}
}

// OsType - Operating System of the cache node
type OsType string

const (
	// OsTypeEflow - cache node installs on Azure Eflow
	OsTypeEflow OsType = "Eflow"
	// OsTypeLinux - cache node installs on Linux Operating system
	OsTypeLinux OsType = "Linux"
	// OsTypeWindows - cache node installs on windows operating system
	OsTypeWindows OsType = "Windows"
)

// PossibleOsTypeValues returns the possible values for the OsType const type.
func PossibleOsTypeValues() []OsType {
	return []OsType{
		OsTypeEflow,
		OsTypeLinux,
		OsTypeWindows,
	}
}

// ProvisioningState - provisioning state of the resource
type ProvisioningState string

const (
	// ProvisioningStateAccepted - Accepted state of the provisioning state during the Async Operations
	ProvisioningStateAccepted ProvisioningState = "Accepted"
	// ProvisioningStateCanceled - Resource creation was canceled.
	ProvisioningStateCanceled ProvisioningState = "Canceled"
	// ProvisioningStateDeleting - Deleting state of the provisioning state
	ProvisioningStateDeleting ProvisioningState = "Deleting"
	// ProvisioningStateFailed - Resource creation failed.
	ProvisioningStateFailed ProvisioningState = "Failed"
	// ProvisioningStateSucceeded - Resource has been created.
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	// ProvisioningStateUnknown - unknown state of the provisioning state
	ProvisioningStateUnknown ProvisioningState = "Unknown"
	// ProvisioningStateUpgrading - Upgrading state of the provisioning state
	ProvisioningStateUpgrading ProvisioningState = "Upgrading"
)

// PossibleProvisioningStateValues returns the possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{
		ProvisioningStateAccepted,
		ProvisioningStateCanceled,
		ProvisioningStateDeleting,
		ProvisioningStateFailed,
		ProvisioningStateSucceeded,
		ProvisioningStateUnknown,
		ProvisioningStateUpgrading,
	}
}

// ProxyRequired - Proxy details enum
type ProxyRequired string

const (
	// ProxyRequiredNone - Proxy is not required in setup
	ProxyRequiredNone ProxyRequired = "None"
	// ProxyRequiredRequired - proxy is required in setup
	ProxyRequiredRequired ProxyRequired = "Required"
)

// PossibleProxyRequiredValues returns the possible values for the ProxyRequired const type.
func PossibleProxyRequiredValues() []ProxyRequired {
	return []ProxyRequired{
		ProxyRequiredNone,
		ProxyRequiredRequired,
	}
}