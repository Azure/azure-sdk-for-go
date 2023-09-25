//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpurview

const (
	moduleName = "armpurview"
	moduleVersion = "v1.1.0"
)

// CreatedByType - The type of identity that created the resource.
type CreatedByType string

const (
	CreatedByTypeApplication CreatedByType = "Application"
	CreatedByTypeKey CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
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

// LastModifiedByType - The type of identity that last modified the resource.
type LastModifiedByType string

const (
	LastModifiedByTypeApplication LastModifiedByType = "Application"
	LastModifiedByTypeKey LastModifiedByType = "Key"
	LastModifiedByTypeManagedIdentity LastModifiedByType = "ManagedIdentity"
	LastModifiedByTypeUser LastModifiedByType = "User"
)

// PossibleLastModifiedByTypeValues returns the possible values for the LastModifiedByType const type.
func PossibleLastModifiedByTypeValues() []LastModifiedByType {
	return []LastModifiedByType{	
		LastModifiedByTypeApplication,
		LastModifiedByTypeKey,
		LastModifiedByTypeManagedIdentity,
		LastModifiedByTypeUser,
	}
}

// Name - Gets or sets the sku name.
type Name string

const (
	NameStandard Name = "Standard"
)

// PossibleNameValues returns the possible values for the Name const type.
func PossibleNameValues() []Name {
	return []Name{	
		NameStandard,
	}
}

// ProvisioningState - Gets or sets the state of the provisioning.
type ProvisioningState string

const (
	ProvisioningStateCanceled ProvisioningState = "Canceled"
	ProvisioningStateCreating ProvisioningState = "Creating"
	ProvisioningStateDeleting ProvisioningState = "Deleting"
	ProvisioningStateFailed ProvisioningState = "Failed"
	ProvisioningStateMoving ProvisioningState = "Moving"
	ProvisioningStateSoftDeleted ProvisioningState = "SoftDeleted"
	ProvisioningStateSoftDeleting ProvisioningState = "SoftDeleting"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUnknown ProvisioningState = "Unknown"
)

// PossibleProvisioningStateValues returns the possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{	
		ProvisioningStateCanceled,
		ProvisioningStateCreating,
		ProvisioningStateDeleting,
		ProvisioningStateFailed,
		ProvisioningStateMoving,
		ProvisioningStateSoftDeleted,
		ProvisioningStateSoftDeleting,
		ProvisioningStateSucceeded,
		ProvisioningStateUnknown,
	}
}

// PublicNetworkAccess - Gets or sets the public network access.
type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled PublicNetworkAccess = "Enabled"
	PublicNetworkAccessNotSpecified PublicNetworkAccess = "NotSpecified"
)

// PossiblePublicNetworkAccessValues returns the possible values for the PublicNetworkAccess const type.
func PossiblePublicNetworkAccessValues() []PublicNetworkAccess {
	return []PublicNetworkAccess{	
		PublicNetworkAccessDisabled,
		PublicNetworkAccessEnabled,
		PublicNetworkAccessNotSpecified,
	}
}

// Reason - The reason the name is not available.
type Reason string

const (
	ReasonAlreadyExists Reason = "AlreadyExists"
	ReasonInvalid Reason = "Invalid"
)

// PossibleReasonValues returns the possible values for the Reason const type.
func PossibleReasonValues() []Reason {
	return []Reason{	
		ReasonAlreadyExists,
		ReasonInvalid,
	}
}

// ScopeType - The scope where the default account is set.
type ScopeType string

const (
	ScopeTypeSubscription ScopeType = "Subscription"
	ScopeTypeTenant ScopeType = "Tenant"
)

// PossibleScopeTypeValues returns the possible values for the ScopeType const type.
func PossibleScopeTypeValues() []ScopeType {
	return []ScopeType{	
		ScopeTypeSubscription,
		ScopeTypeTenant,
	}
}

// Status - The status.
type Status string

const (
	StatusApproved Status = "Approved"
	StatusDisconnected Status = "Disconnected"
	StatusPending Status = "Pending"
	StatusRejected Status = "Rejected"
	StatusUnknown Status = "Unknown"
)

// PossibleStatusValues returns the possible values for the Status const type.
func PossibleStatusValues() []Status {
	return []Status{	
		StatusApproved,
		StatusDisconnected,
		StatusPending,
		StatusRejected,
		StatusUnknown,
	}
}

// Type - Identity Type
type Type string

const (
	TypeNone Type = "None"
	TypeSystemAssigned Type = "SystemAssigned"
	TypeUserAssigned Type = "UserAssigned"
)

// PossibleTypeValues returns the possible values for the Type const type.
func PossibleTypeValues() []Type {
	return []Type{	
		TypeNone,
		TypeSystemAssigned,
		TypeUserAssigned,
	}
}

