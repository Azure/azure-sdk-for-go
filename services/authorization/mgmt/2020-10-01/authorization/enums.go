package authorization

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// ApprovalMode enumerates the values for approval mode.
type ApprovalMode string

const (
	// ApprovalModeNoApproval ...
	ApprovalModeNoApproval ApprovalMode = "NoApproval"
	// ApprovalModeParallel ...
	ApprovalModeParallel ApprovalMode = "Parallel"
	// ApprovalModeSerial ...
	ApprovalModeSerial ApprovalMode = "Serial"
	// ApprovalModeSingleStage ...
	ApprovalModeSingleStage ApprovalMode = "SingleStage"
)

// PossibleApprovalModeValues returns an array of possible values for the ApprovalMode const type.
func PossibleApprovalModeValues() []ApprovalMode {
	return []ApprovalMode{ApprovalModeNoApproval, ApprovalModeParallel, ApprovalModeSerial, ApprovalModeSingleStage}
}

// AssignmentType enumerates the values for assignment type.
type AssignmentType string

const (
	// AssignmentTypeActivated ...
	AssignmentTypeActivated AssignmentType = "Activated"
	// AssignmentTypeAssigned ...
	AssignmentTypeAssigned AssignmentType = "Assigned"
)

// PossibleAssignmentTypeValues returns an array of possible values for the AssignmentType const type.
func PossibleAssignmentTypeValues() []AssignmentType {
	return []AssignmentType{AssignmentTypeActivated, AssignmentTypeAssigned}
}

// EnablementRules enumerates the values for enablement rules.
type EnablementRules string

const (
	// EnablementRulesJustification ...
	EnablementRulesJustification EnablementRules = "Justification"
	// EnablementRulesMultiFactorAuthentication ...
	EnablementRulesMultiFactorAuthentication EnablementRules = "MultiFactorAuthentication"
	// EnablementRulesTicketing ...
	EnablementRulesTicketing EnablementRules = "Ticketing"
)

// PossibleEnablementRulesValues returns an array of possible values for the EnablementRules const type.
func PossibleEnablementRulesValues() []EnablementRules {
	return []EnablementRules{EnablementRulesJustification, EnablementRulesMultiFactorAuthentication, EnablementRulesTicketing}
}

// MemberType enumerates the values for member type.
type MemberType string

const (
	// MemberTypeDirect ...
	MemberTypeDirect MemberType = "Direct"
	// MemberTypeGroup ...
	MemberTypeGroup MemberType = "Group"
	// MemberTypeInherited ...
	MemberTypeInherited MemberType = "Inherited"
)

// PossibleMemberTypeValues returns an array of possible values for the MemberType const type.
func PossibleMemberTypeValues() []MemberType {
	return []MemberType{MemberTypeDirect, MemberTypeGroup, MemberTypeInherited}
}

// NotificationDeliveryMechanism enumerates the values for notification delivery mechanism.
type NotificationDeliveryMechanism string

const (
	// NotificationDeliveryMechanismEmail ...
	NotificationDeliveryMechanismEmail NotificationDeliveryMechanism = "Email"
)

// PossibleNotificationDeliveryMechanismValues returns an array of possible values for the NotificationDeliveryMechanism const type.
func PossibleNotificationDeliveryMechanismValues() []NotificationDeliveryMechanism {
	return []NotificationDeliveryMechanism{NotificationDeliveryMechanismEmail}
}

// NotificationLevel enumerates the values for notification level.
type NotificationLevel string

const (
	// NotificationLevelAll ...
	NotificationLevelAll NotificationLevel = "All"
	// NotificationLevelCritical ...
	NotificationLevelCritical NotificationLevel = "Critical"
	// NotificationLevelNone ...
	NotificationLevelNone NotificationLevel = "None"
)

// PossibleNotificationLevelValues returns an array of possible values for the NotificationLevel const type.
func PossibleNotificationLevelValues() []NotificationLevel {
	return []NotificationLevel{NotificationLevelAll, NotificationLevelCritical, NotificationLevelNone}
}

// PrincipalType enumerates the values for principal type.
type PrincipalType string

const (
	// PrincipalTypeDevice ...
	PrincipalTypeDevice PrincipalType = "Device"
	// PrincipalTypeForeignGroup ...
	PrincipalTypeForeignGroup PrincipalType = "ForeignGroup"
	// PrincipalTypeGroup ...
	PrincipalTypeGroup PrincipalType = "Group"
	// PrincipalTypeServicePrincipal ...
	PrincipalTypeServicePrincipal PrincipalType = "ServicePrincipal"
	// PrincipalTypeUser ...
	PrincipalTypeUser PrincipalType = "User"
)

// PossiblePrincipalTypeValues returns an array of possible values for the PrincipalType const type.
func PossiblePrincipalTypeValues() []PrincipalType {
	return []PrincipalType{PrincipalTypeDevice, PrincipalTypeForeignGroup, PrincipalTypeGroup, PrincipalTypeServicePrincipal, PrincipalTypeUser}
}

// RecipientType enumerates the values for recipient type.
type RecipientType string

const (
	// RecipientTypeAdmin ...
	RecipientTypeAdmin RecipientType = "Admin"
	// RecipientTypeApprover ...
	RecipientTypeApprover RecipientType = "Approver"
	// RecipientTypeRequestor ...
	RecipientTypeRequestor RecipientType = "Requestor"
)

// PossibleRecipientTypeValues returns an array of possible values for the RecipientType const type.
func PossibleRecipientTypeValues() []RecipientType {
	return []RecipientType{RecipientTypeAdmin, RecipientTypeApprover, RecipientTypeRequestor}
}

// RequestType enumerates the values for request type.
type RequestType string

const (
	// RequestTypeAdminAssign ...
	RequestTypeAdminAssign RequestType = "AdminAssign"
	// RequestTypeAdminExtend ...
	RequestTypeAdminExtend RequestType = "AdminExtend"
	// RequestTypeAdminRemove ...
	RequestTypeAdminRemove RequestType = "AdminRemove"
	// RequestTypeAdminRenew ...
	RequestTypeAdminRenew RequestType = "AdminRenew"
	// RequestTypeAdminUpdate ...
	RequestTypeAdminUpdate RequestType = "AdminUpdate"
	// RequestTypeSelfActivate ...
	RequestTypeSelfActivate RequestType = "SelfActivate"
	// RequestTypeSelfDeactivate ...
	RequestTypeSelfDeactivate RequestType = "SelfDeactivate"
	// RequestTypeSelfExtend ...
	RequestTypeSelfExtend RequestType = "SelfExtend"
	// RequestTypeSelfRenew ...
	RequestTypeSelfRenew RequestType = "SelfRenew"
)

// PossibleRequestTypeValues returns an array of possible values for the RequestType const type.
func PossibleRequestTypeValues() []RequestType {
	return []RequestType{RequestTypeAdminAssign, RequestTypeAdminExtend, RequestTypeAdminRemove, RequestTypeAdminRenew, RequestTypeAdminUpdate, RequestTypeSelfActivate, RequestTypeSelfDeactivate, RequestTypeSelfExtend, RequestTypeSelfRenew}
}

// RoleManagementPolicyRuleType enumerates the values for role management policy rule type.
type RoleManagementPolicyRuleType string

const (
	// RoleManagementPolicyRuleTypeRoleManagementPolicyApprovalRule ...
	RoleManagementPolicyRuleTypeRoleManagementPolicyApprovalRule RoleManagementPolicyRuleType = "RoleManagementPolicyApprovalRule"
	// RoleManagementPolicyRuleTypeRoleManagementPolicyAuthenticationContextRule ...
	RoleManagementPolicyRuleTypeRoleManagementPolicyAuthenticationContextRule RoleManagementPolicyRuleType = "RoleManagementPolicyAuthenticationContextRule"
	// RoleManagementPolicyRuleTypeRoleManagementPolicyEnablementRule ...
	RoleManagementPolicyRuleTypeRoleManagementPolicyEnablementRule RoleManagementPolicyRuleType = "RoleManagementPolicyEnablementRule"
	// RoleManagementPolicyRuleTypeRoleManagementPolicyExpirationRule ...
	RoleManagementPolicyRuleTypeRoleManagementPolicyExpirationRule RoleManagementPolicyRuleType = "RoleManagementPolicyExpirationRule"
	// RoleManagementPolicyRuleTypeRoleManagementPolicyNotificationRule ...
	RoleManagementPolicyRuleTypeRoleManagementPolicyNotificationRule RoleManagementPolicyRuleType = "RoleManagementPolicyNotificationRule"
)

// PossibleRoleManagementPolicyRuleTypeValues returns an array of possible values for the RoleManagementPolicyRuleType const type.
func PossibleRoleManagementPolicyRuleTypeValues() []RoleManagementPolicyRuleType {
	return []RoleManagementPolicyRuleType{RoleManagementPolicyRuleTypeRoleManagementPolicyApprovalRule, RoleManagementPolicyRuleTypeRoleManagementPolicyAuthenticationContextRule, RoleManagementPolicyRuleTypeRoleManagementPolicyEnablementRule, RoleManagementPolicyRuleTypeRoleManagementPolicyExpirationRule, RoleManagementPolicyRuleTypeRoleManagementPolicyNotificationRule}
}

// RuleType enumerates the values for rule type.
type RuleType string

const (
	// RuleTypeRoleManagementPolicyApprovalRule ...
	RuleTypeRoleManagementPolicyApprovalRule RuleType = "RoleManagementPolicyApprovalRule"
	// RuleTypeRoleManagementPolicyAuthenticationContextRule ...
	RuleTypeRoleManagementPolicyAuthenticationContextRule RuleType = "RoleManagementPolicyAuthenticationContextRule"
	// RuleTypeRoleManagementPolicyEnablementRule ...
	RuleTypeRoleManagementPolicyEnablementRule RuleType = "RoleManagementPolicyEnablementRule"
	// RuleTypeRoleManagementPolicyExpirationRule ...
	RuleTypeRoleManagementPolicyExpirationRule RuleType = "RoleManagementPolicyExpirationRule"
	// RuleTypeRoleManagementPolicyNotificationRule ...
	RuleTypeRoleManagementPolicyNotificationRule RuleType = "RoleManagementPolicyNotificationRule"
	// RuleTypeRoleManagementPolicyRule ...
	RuleTypeRoleManagementPolicyRule RuleType = "RoleManagementPolicyRule"
)

// PossibleRuleTypeValues returns an array of possible values for the RuleType const type.
func PossibleRuleTypeValues() []RuleType {
	return []RuleType{RuleTypeRoleManagementPolicyApprovalRule, RuleTypeRoleManagementPolicyAuthenticationContextRule, RuleTypeRoleManagementPolicyEnablementRule, RuleTypeRoleManagementPolicyExpirationRule, RuleTypeRoleManagementPolicyNotificationRule, RuleTypeRoleManagementPolicyRule}
}

// Status enumerates the values for status.
type Status string

const (
	// StatusAccepted ...
	StatusAccepted Status = "Accepted"
	// StatusAdminApproved ...
	StatusAdminApproved Status = "AdminApproved"
	// StatusAdminDenied ...
	StatusAdminDenied Status = "AdminDenied"
	// StatusCanceled ...
	StatusCanceled Status = "Canceled"
	// StatusDenied ...
	StatusDenied Status = "Denied"
	// StatusFailed ...
	StatusFailed Status = "Failed"
	// StatusFailedAsResourceIsLocked ...
	StatusFailedAsResourceIsLocked Status = "FailedAsResourceIsLocked"
	// StatusGranted ...
	StatusGranted Status = "Granted"
	// StatusInvalid ...
	StatusInvalid Status = "Invalid"
	// StatusPendingAdminDecision ...
	StatusPendingAdminDecision Status = "PendingAdminDecision"
	// StatusPendingApproval ...
	StatusPendingApproval Status = "PendingApproval"
	// StatusPendingApprovalProvisioning ...
	StatusPendingApprovalProvisioning Status = "PendingApprovalProvisioning"
	// StatusPendingEvaluation ...
	StatusPendingEvaluation Status = "PendingEvaluation"
	// StatusPendingExternalProvisioning ...
	StatusPendingExternalProvisioning Status = "PendingExternalProvisioning"
	// StatusPendingProvisioning ...
	StatusPendingProvisioning Status = "PendingProvisioning"
	// StatusPendingRevocation ...
	StatusPendingRevocation Status = "PendingRevocation"
	// StatusPendingScheduleCreation ...
	StatusPendingScheduleCreation Status = "PendingScheduleCreation"
	// StatusProvisioned ...
	StatusProvisioned Status = "Provisioned"
	// StatusProvisioningStarted ...
	StatusProvisioningStarted Status = "ProvisioningStarted"
	// StatusRevoked ...
	StatusRevoked Status = "Revoked"
	// StatusScheduleCreated ...
	StatusScheduleCreated Status = "ScheduleCreated"
	// StatusTimedOut ...
	StatusTimedOut Status = "TimedOut"
)

// PossibleStatusValues returns an array of possible values for the Status const type.
func PossibleStatusValues() []Status {
	return []Status{StatusAccepted, StatusAdminApproved, StatusAdminDenied, StatusCanceled, StatusDenied, StatusFailed, StatusFailedAsResourceIsLocked, StatusGranted, StatusInvalid, StatusPendingAdminDecision, StatusPendingApproval, StatusPendingApprovalProvisioning, StatusPendingEvaluation, StatusPendingExternalProvisioning, StatusPendingProvisioning, StatusPendingRevocation, StatusPendingScheduleCreation, StatusProvisioned, StatusProvisioningStarted, StatusRevoked, StatusScheduleCreated, StatusTimedOut}
}

// Type enumerates the values for type.
type Type string

const (
	// TypeAfterDateTime ...
	TypeAfterDateTime Type = "AfterDateTime"
	// TypeAfterDuration ...
	TypeAfterDuration Type = "AfterDuration"
	// TypeNoExpiration ...
	TypeNoExpiration Type = "NoExpiration"
)

// PossibleTypeValues returns an array of possible values for the Type const type.
func PossibleTypeValues() []Type {
	return []Type{TypeAfterDateTime, TypeAfterDuration, TypeNoExpiration}
}

// UserType enumerates the values for user type.
type UserType string

const (
	// UserTypeGroup ...
	UserTypeGroup UserType = "Group"
	// UserTypeUser ...
	UserTypeUser UserType = "User"
)

// PossibleUserTypeValues returns an array of possible values for the UserType const type.
func PossibleUserTypeValues() []UserType {
	return []UserType{UserTypeGroup, UserTypeUser}
}
