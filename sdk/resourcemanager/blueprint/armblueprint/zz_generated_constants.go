//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armblueprint

const (
	moduleName    = "armblueprint"
	moduleVersion = "v0.2.0"
)

// ArtifactKind - Specifies the kind of blueprint artifact.
type ArtifactKind string

const (
	ArtifactKindPolicyAssignment ArtifactKind = "policyAssignment"
	ArtifactKindRoleAssignment   ArtifactKind = "roleAssignment"
	ArtifactKindTemplate         ArtifactKind = "template"
)

// PossibleArtifactKindValues returns the possible values for the ArtifactKind const type.
func PossibleArtifactKindValues() []ArtifactKind {
	return []ArtifactKind{
		ArtifactKindPolicyAssignment,
		ArtifactKindRoleAssignment,
		ArtifactKindTemplate,
	}
}

// ToPtr returns a *ArtifactKind pointing to the current value.
func (c ArtifactKind) ToPtr() *ArtifactKind {
	return &c
}

type AssignmentDeleteBehavior string

const (
	AssignmentDeleteBehaviorAll  AssignmentDeleteBehavior = "all"
	AssignmentDeleteBehaviorNone AssignmentDeleteBehavior = "none"
)

// PossibleAssignmentDeleteBehaviorValues returns the possible values for the AssignmentDeleteBehavior const type.
func PossibleAssignmentDeleteBehaviorValues() []AssignmentDeleteBehavior {
	return []AssignmentDeleteBehavior{
		AssignmentDeleteBehaviorAll,
		AssignmentDeleteBehaviorNone,
	}
}

// ToPtr returns a *AssignmentDeleteBehavior pointing to the current value.
func (c AssignmentDeleteBehavior) ToPtr() *AssignmentDeleteBehavior {
	return &c
}

// AssignmentLockMode - Lock mode.
type AssignmentLockMode string

const (
	AssignmentLockModeAllResourcesDoNotDelete AssignmentLockMode = "AllResourcesDoNotDelete"
	AssignmentLockModeAllResourcesReadOnly    AssignmentLockMode = "AllResourcesReadOnly"
	AssignmentLockModeNone                    AssignmentLockMode = "None"
)

// PossibleAssignmentLockModeValues returns the possible values for the AssignmentLockMode const type.
func PossibleAssignmentLockModeValues() []AssignmentLockMode {
	return []AssignmentLockMode{
		AssignmentLockModeAllResourcesDoNotDelete,
		AssignmentLockModeAllResourcesReadOnly,
		AssignmentLockModeNone,
	}
}

// ToPtr returns a *AssignmentLockMode pointing to the current value.
func (c AssignmentLockMode) ToPtr() *AssignmentLockMode {
	return &c
}

// AssignmentProvisioningState - State of the blueprint assignment.
type AssignmentProvisioningState string

const (
	AssignmentProvisioningStateCanceled   AssignmentProvisioningState = "canceled"
	AssignmentProvisioningStateCancelling AssignmentProvisioningState = "cancelling"
	AssignmentProvisioningStateCreating   AssignmentProvisioningState = "creating"
	AssignmentProvisioningStateDeleting   AssignmentProvisioningState = "deleting"
	AssignmentProvisioningStateDeploying  AssignmentProvisioningState = "deploying"
	AssignmentProvisioningStateFailed     AssignmentProvisioningState = "failed"
	AssignmentProvisioningStateLocking    AssignmentProvisioningState = "locking"
	AssignmentProvisioningStateSucceeded  AssignmentProvisioningState = "succeeded"
	AssignmentProvisioningStateValidating AssignmentProvisioningState = "validating"
	AssignmentProvisioningStateWaiting    AssignmentProvisioningState = "waiting"
)

// PossibleAssignmentProvisioningStateValues returns the possible values for the AssignmentProvisioningState const type.
func PossibleAssignmentProvisioningStateValues() []AssignmentProvisioningState {
	return []AssignmentProvisioningState{
		AssignmentProvisioningStateCanceled,
		AssignmentProvisioningStateCancelling,
		AssignmentProvisioningStateCreating,
		AssignmentProvisioningStateDeleting,
		AssignmentProvisioningStateDeploying,
		AssignmentProvisioningStateFailed,
		AssignmentProvisioningStateLocking,
		AssignmentProvisioningStateSucceeded,
		AssignmentProvisioningStateValidating,
		AssignmentProvisioningStateWaiting,
	}
}

// ToPtr returns a *AssignmentProvisioningState pointing to the current value.
func (c AssignmentProvisioningState) ToPtr() *AssignmentProvisioningState {
	return &c
}

// BlueprintTargetScope - The scope where this blueprint definition can be assigned.
type BlueprintTargetScope string

const (
	// BlueprintTargetScopeManagementGroup - The blueprint targets a management group during blueprint assignment. This is reserved
	// for future use.
	BlueprintTargetScopeManagementGroup BlueprintTargetScope = "managementGroup"
	// BlueprintTargetScopeSubscription - The blueprint targets a subscription during blueprint assignment.
	BlueprintTargetScopeSubscription BlueprintTargetScope = "subscription"
)

// PossibleBlueprintTargetScopeValues returns the possible values for the BlueprintTargetScope const type.
func PossibleBlueprintTargetScopeValues() []BlueprintTargetScope {
	return []BlueprintTargetScope{
		BlueprintTargetScopeManagementGroup,
		BlueprintTargetScopeSubscription,
	}
}

// ToPtr returns a *BlueprintTargetScope pointing to the current value.
func (c BlueprintTargetScope) ToPtr() *BlueprintTargetScope {
	return &c
}

// ManagedServiceIdentityType - Type of the managed identity.
type ManagedServiceIdentityType string

const (
	ManagedServiceIdentityTypeNone           ManagedServiceIdentityType = "None"
	ManagedServiceIdentityTypeSystemAssigned ManagedServiceIdentityType = "SystemAssigned"
	ManagedServiceIdentityTypeUserAssigned   ManagedServiceIdentityType = "UserAssigned"
)

// PossibleManagedServiceIdentityTypeValues returns the possible values for the ManagedServiceIdentityType const type.
func PossibleManagedServiceIdentityTypeValues() []ManagedServiceIdentityType {
	return []ManagedServiceIdentityType{
		ManagedServiceIdentityTypeNone,
		ManagedServiceIdentityTypeSystemAssigned,
		ManagedServiceIdentityTypeUserAssigned,
	}
}

// ToPtr returns a *ManagedServiceIdentityType pointing to the current value.
func (c ManagedServiceIdentityType) ToPtr() *ManagedServiceIdentityType {
	return &c
}

// TemplateParameterType - Allowed data types for Resource Manager template parameters.
type TemplateParameterType string

const (
	TemplateParameterTypeArray        TemplateParameterType = "array"
	TemplateParameterTypeBool         TemplateParameterType = "bool"
	TemplateParameterTypeInt          TemplateParameterType = "int"
	TemplateParameterTypeObject       TemplateParameterType = "object"
	TemplateParameterTypeSecureObject TemplateParameterType = "secureObject"
	TemplateParameterTypeSecureString TemplateParameterType = "secureString"
	TemplateParameterTypeString       TemplateParameterType = "string"
)

// PossibleTemplateParameterTypeValues returns the possible values for the TemplateParameterType const type.
func PossibleTemplateParameterTypeValues() []TemplateParameterType {
	return []TemplateParameterType{
		TemplateParameterTypeArray,
		TemplateParameterTypeBool,
		TemplateParameterTypeInt,
		TemplateParameterTypeObject,
		TemplateParameterTypeSecureObject,
		TemplateParameterTypeSecureString,
		TemplateParameterTypeString,
	}
}

// ToPtr returns a *TemplateParameterType pointing to the current value.
func (c TemplateParameterType) ToPtr() *TemplateParameterType {
	return &c
}
