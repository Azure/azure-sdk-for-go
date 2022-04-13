//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmarketplace

const (
	moduleName    = "armmarketplace"
	moduleVersion = "v0.3.0"
)

// Accessibility - Plan accessibility
type Accessibility string

const (
	AccessibilityPrivateSubscriptionOnLevel Accessibility = "PrivateSubscriptionOnLevel"
	AccessibilityPrivateTenantOnLevel       Accessibility = "PrivateTenantOnLevel"
	AccessibilityPublic                     Accessibility = "Public"
	AccessibilityUnknown                    Accessibility = "Unknown"
)

// PossibleAccessibilityValues returns the possible values for the Accessibility const type.
func PossibleAccessibilityValues() []Accessibility {
	return []Accessibility{
		AccessibilityPrivateSubscriptionOnLevel,
		AccessibilityPrivateTenantOnLevel,
		AccessibilityPublic,
		AccessibilityUnknown,
	}
}

// AdminAction - Gets or sets admin action
type AdminAction string

const (
	AdminActionApproved AdminAction = "Approved"
	AdminActionRejected AdminAction = "Rejected"
)

// PossibleAdminActionValues returns the possible values for the AdminAction const type.
func PossibleAdminActionValues() []AdminAction {
	return []AdminAction{
		AdminActionApproved,
		AdminActionRejected,
	}
}

// Availability - Indicates private store availability
type Availability string

const (
	AvailabilityDisabled Availability = "disabled"
	AvailabilityEnabled  Availability = "enabled"
)

// PossibleAvailabilityValues returns the possible values for the Availability const type.
func PossibleAvailabilityValues() []Availability {
	return []Availability{
		AvailabilityDisabled,
		AvailabilityEnabled,
	}
}

// IdentityType - The type of identity that creates/modifies resources
type IdentityType string

const (
	IdentityTypeApplication     IdentityType = "Application"
	IdentityTypeKey             IdentityType = "Key"
	IdentityTypeManagedIdentity IdentityType = "ManagedIdentity"
	IdentityTypeUser            IdentityType = "User"
)

// PossibleIdentityTypeValues returns the possible values for the IdentityType const type.
func PossibleIdentityTypeValues() []IdentityType {
	return []IdentityType{
		IdentityTypeApplication,
		IdentityTypeKey,
		IdentityTypeManagedIdentity,
		IdentityTypeUser,
	}
}

// Operation - Set the Operation for the POST method. Ping or Delete
type Operation string

const (
	OperationDeletePrivateStoreCollection      Operation = "DeletePrivateStoreCollection"
	OperationDeletePrivateStoreCollectionOffer Operation = "DeletePrivateStoreCollectionOffer"
	OperationDeletePrivateStoreOffer           Operation = "DeletePrivateStoreOffer"
	OperationPing                              Operation = "Ping"
)

// PossibleOperationValues returns the possible values for the Operation const type.
func PossibleOperationValues() []Operation {
	return []Operation{
		OperationDeletePrivateStoreCollection,
		OperationDeletePrivateStoreCollectionOffer,
		OperationDeletePrivateStoreOffer,
		OperationPing,
	}
}

// Status - Gets the plan status
type Status string

const (
	StatusApproved Status = "Approved"
	StatusNone     Status = "None"
	StatusPending  Status = "Pending"
	StatusRejected Status = "Rejected"
)

// PossibleStatusValues returns the possible values for the Status const type.
func PossibleStatusValues() []Status {
	return []Status{
		StatusApproved,
		StatusNone,
		StatusPending,
		StatusRejected,
	}
}

// SubscriptionState - The subscription state. Possible values are Enabled, Warned, PastDue, Disabled, and Deleted.
type SubscriptionState string

const (
	SubscriptionStateDeleted  SubscriptionState = "Deleted"
	SubscriptionStateDisabled SubscriptionState = "Disabled"
	SubscriptionStateEnabled  SubscriptionState = "Enabled"
	SubscriptionStatePastDue  SubscriptionState = "PastDue"
	SubscriptionStateWarned   SubscriptionState = "Warned"
)

// PossibleSubscriptionStateValues returns the possible values for the SubscriptionState const type.
func PossibleSubscriptionStateValues() []SubscriptionState {
	return []SubscriptionState{
		SubscriptionStateDeleted,
		SubscriptionStateDisabled,
		SubscriptionStateEnabled,
		SubscriptionStatePastDue,
		SubscriptionStateWarned,
	}
}
