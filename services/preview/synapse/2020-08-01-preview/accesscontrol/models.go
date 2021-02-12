package accesscontrol

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/satori/go.uuid"
)

// The package's fully qualified name.
const fqdn = "github.com/Azure/azure-sdk-for-go/services/preview/synapse/2020-08-01-preview/accesscontrol"

// CheckAccessDecision check access response details
type CheckAccessDecision struct {
	// AccessDecision - Access Decision.
	AccessDecision *string `json:"accessDecision,omitempty"`
	// ActionID - Action Id.
	ActionID       *string                `json:"actionId,omitempty"`
	RoleAssignment *RoleAssignmentDetails `json:"roleAssignment,omitempty"`
}

// CheckPrincipalAccessRequest check access request details
type CheckPrincipalAccessRequest struct {
	// Subject - Subject details
	Subject *SubjectInfo `json:"subject,omitempty"`
	// Actions - List of actions.
	Actions *[]RequiredAction `json:"actions,omitempty"`
	// Scope - Scope at which the check access is done.
	Scope *string `json:"scope,omitempty"`
}

// CheckPrincipalAccessResponse check access response details
type CheckPrincipalAccessResponse struct {
	autorest.Response `json:"-"`
	// AccessDecisions - To check if the current user, group, or service principal has permission to read artifacts in the specified workspace.
	AccessDecisions *[]CheckAccessDecision `json:"AccessDecisions,omitempty"`
}

// ErrorContract contains details when the response code indicates an error.
type ErrorContract struct {
	// Error - The error details.
	Error *ErrorResponse `json:"error,omitempty"`
}

// ErrorDetail ...
type ErrorDetail struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
	Target  *string `json:"target,omitempty"`
}

// ErrorResponse ...
type ErrorResponse struct {
	Code    *string        `json:"code,omitempty"`
	Message *string        `json:"message,omitempty"`
	Target  *string        `json:"target,omitempty"`
	Details *[]ErrorDetail `json:"details,omitempty"`
}

// ListString ...
type ListString struct {
	autorest.Response `json:"-"`
	Value             *[]string `json:"value,omitempty"`
}

// ListSynapseRoleDefinition ...
type ListSynapseRoleDefinition struct {
	autorest.Response `json:"-"`
	Value             *[]SynapseRoleDefinition `json:"value,omitempty"`
}

// RequiredAction action Info
type RequiredAction struct {
	// ID - Action Id.
	ID *string `json:"id,omitempty"`
	// IsDataAction - Is a data action or not.
	IsDataAction *bool `json:"isDataAction,omitempty"`
}

// RoleAssignmentDetails role Assignment response details
type RoleAssignmentDetails struct {
	autorest.Response `json:"-"`
	// ID - Role Assignment ID
	ID *string `json:"id,omitempty"`
	// RoleDefinitionID - Role ID of the Synapse Built-In Role
	RoleDefinitionID *uuid.UUID `json:"roleDefinitionId,omitempty"`
	// PrincipalID - Object ID of the AAD principal or security-group
	PrincipalID *uuid.UUID `json:"principalId,omitempty"`
	// Scope - Scope at the role assignment is created
	Scope *string `json:"scope,omitempty"`
	// PrincipalType - Type of the principal Id: User, Group or ServicePrincipal
	PrincipalType *string `json:"principalType,omitempty"`
}

// RoleAssignmentDetailsList role Assignment response details
type RoleAssignmentDetailsList struct {
	autorest.Response `json:"-"`
	// Count - Number of role assignments
	Count *int32 `json:"count,omitempty"`
	// Value - A list of role assignments
	Value *[]RoleAssignmentDetails `json:"value,omitempty"`
}

// RoleAssignmentRequest role Assignment request details
type RoleAssignmentRequest struct {
	// RoleID - Role ID of the Synapse Built-In Role
	RoleID *uuid.UUID `json:"roleId,omitempty"`
	// PrincipalID - Object ID of the AAD principal or security-group
	PrincipalID *uuid.UUID `json:"principalId,omitempty"`
	// Scope - Scope at which the role assignment is created
	Scope *string `json:"scope,omitempty"`
	// PrincipalType - Type of the principal Id: User, Group or ServicePrincipal
	PrincipalType *string `json:"principalType,omitempty"`
}

// SubjectInfo subject details
type SubjectInfo struct {
	// PrincipalID - Principal Id
	PrincipalID *uuid.UUID `json:"principalId,omitempty"`
	// GroupIds - List of group Ids that the principalId is part of.
	GroupIds *[]uuid.UUID `json:"groupIds,omitempty"`
}

// SynapseRbacPermission synapse role definition details
type SynapseRbacPermission struct {
	// Actions - List of actions
	Actions *[]string `json:"actions,omitempty"`
	// NotActions - List of Not actions
	NotActions *[]string `json:"notActions,omitempty"`
	// DataActions - List of data actions
	DataActions *[]string `json:"dataActions,omitempty"`
	// NotDataActions - List of Not data actions
	NotDataActions *[]string `json:"notDataActions,omitempty"`
}

// SynapseRoleDefinition synapse role definition details
type SynapseRoleDefinition struct {
	autorest.Response `json:"-"`
	// ID - Role Definition ID
	ID *uuid.UUID `json:"id,omitempty"`
	// Name - Name of the Synapse role
	Name *string `json:"name,omitempty"`
	// IsBuiltIn - Is a built-in role or not
	IsBuiltIn *bool `json:"isBuiltIn,omitempty"`
	// Description - Description for the Synapse role
	Description *string `json:"description,omitempty"`
	// Permissions - Permissions for the Synapse role
	Permissions *[]SynapseRbacPermission `json:"permissions,omitempty"`
	// Scopes - Allowed scopes for the Synapse role
	Scopes *[]string `json:"scopes,omitempty"`
	// AvailabilityStatus - Availability of the Synapse role
	AvailabilityStatus *string `json:"availabilityStatus,omitempty"`
}
