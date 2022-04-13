//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armauthorization

const (
	moduleName    = "armauthorization"
	moduleVersion = "v0.4.0"
)

// AccessRecommendationType - The feature- generated recommendation shown to the reviewer.
type AccessRecommendationType string

const (
	AccessRecommendationTypeApprove         AccessRecommendationType = "Approve"
	AccessRecommendationTypeDeny            AccessRecommendationType = "Deny"
	AccessRecommendationTypeNoInfoAvailable AccessRecommendationType = "NoInfoAvailable"
)

// PossibleAccessRecommendationTypeValues returns the possible values for the AccessRecommendationType const type.
func PossibleAccessRecommendationTypeValues() []AccessRecommendationType {
	return []AccessRecommendationType{
		AccessRecommendationTypeApprove,
		AccessRecommendationTypeDeny,
		AccessRecommendationTypeNoInfoAvailable,
	}
}

// AccessReviewActorIdentityType - The identity type : user/servicePrincipal
type AccessReviewActorIdentityType string

const (
	AccessReviewActorIdentityTypeServicePrincipal AccessReviewActorIdentityType = "servicePrincipal"
	AccessReviewActorIdentityTypeUser             AccessReviewActorIdentityType = "user"
)

// PossibleAccessReviewActorIdentityTypeValues returns the possible values for the AccessReviewActorIdentityType const type.
func PossibleAccessReviewActorIdentityTypeValues() []AccessReviewActorIdentityType {
	return []AccessReviewActorIdentityType{
		AccessReviewActorIdentityTypeServicePrincipal,
		AccessReviewActorIdentityTypeUser,
	}
}

// AccessReviewApplyResult - The outcome of applying the decision.
type AccessReviewApplyResult string

const (
	AccessReviewApplyResultAppliedSuccessfully                  AccessReviewApplyResult = "AppliedSuccessfully"
	AccessReviewApplyResultAppliedSuccessfullyButObjectNotFound AccessReviewApplyResult = "AppliedSuccessfullyButObjectNotFound"
	AccessReviewApplyResultAppliedWithUnknownFailure            AccessReviewApplyResult = "AppliedWithUnknownFailure"
	AccessReviewApplyResultApplyNotSupported                    AccessReviewApplyResult = "ApplyNotSupported"
	AccessReviewApplyResultApplying                             AccessReviewApplyResult = "Applying"
	AccessReviewApplyResultNew                                  AccessReviewApplyResult = "New"
)

// PossibleAccessReviewApplyResultValues returns the possible values for the AccessReviewApplyResult const type.
func PossibleAccessReviewApplyResultValues() []AccessReviewApplyResult {
	return []AccessReviewApplyResult{
		AccessReviewApplyResultAppliedSuccessfully,
		AccessReviewApplyResultAppliedSuccessfullyButObjectNotFound,
		AccessReviewApplyResultAppliedWithUnknownFailure,
		AccessReviewApplyResultApplyNotSupported,
		AccessReviewApplyResultApplying,
		AccessReviewApplyResultNew,
	}
}

// AccessReviewHistoryDefinitionStatus - This read-only field specifies the of the requested review history data. This is
// either requested, in-progress, done or error.
type AccessReviewHistoryDefinitionStatus string

const (
	AccessReviewHistoryDefinitionStatusDone       AccessReviewHistoryDefinitionStatus = "Done"
	AccessReviewHistoryDefinitionStatusError      AccessReviewHistoryDefinitionStatus = "Error"
	AccessReviewHistoryDefinitionStatusInProgress AccessReviewHistoryDefinitionStatus = "InProgress"
	AccessReviewHistoryDefinitionStatusRequested  AccessReviewHistoryDefinitionStatus = "Requested"
)

// PossibleAccessReviewHistoryDefinitionStatusValues returns the possible values for the AccessReviewHistoryDefinitionStatus const type.
func PossibleAccessReviewHistoryDefinitionStatusValues() []AccessReviewHistoryDefinitionStatus {
	return []AccessReviewHistoryDefinitionStatus{
		AccessReviewHistoryDefinitionStatusDone,
		AccessReviewHistoryDefinitionStatusError,
		AccessReviewHistoryDefinitionStatusInProgress,
		AccessReviewHistoryDefinitionStatusRequested,
	}
}

// AccessReviewInstanceReviewersType - This field specifies the type of reviewers for a review. Usually for a review, reviewers
// are explicitly assigned. However, in some cases, the reviewers may not be assigned and instead be chosen
// dynamically. For example managers review or self review.
type AccessReviewInstanceReviewersType string

const (
	AccessReviewInstanceReviewersTypeAssigned AccessReviewInstanceReviewersType = "Assigned"
	AccessReviewInstanceReviewersTypeManagers AccessReviewInstanceReviewersType = "Managers"
	AccessReviewInstanceReviewersTypeSelf     AccessReviewInstanceReviewersType = "Self"
)

// PossibleAccessReviewInstanceReviewersTypeValues returns the possible values for the AccessReviewInstanceReviewersType const type.
func PossibleAccessReviewInstanceReviewersTypeValues() []AccessReviewInstanceReviewersType {
	return []AccessReviewInstanceReviewersType{
		AccessReviewInstanceReviewersTypeAssigned,
		AccessReviewInstanceReviewersTypeManagers,
		AccessReviewInstanceReviewersTypeSelf,
	}
}

// AccessReviewInstanceStatus - This read-only field specifies the status of an access review instance.
type AccessReviewInstanceStatus string

const (
	AccessReviewInstanceStatusApplied       AccessReviewInstanceStatus = "Applied"
	AccessReviewInstanceStatusApplying      AccessReviewInstanceStatus = "Applying"
	AccessReviewInstanceStatusAutoReviewed  AccessReviewInstanceStatus = "AutoReviewed"
	AccessReviewInstanceStatusAutoReviewing AccessReviewInstanceStatus = "AutoReviewing"
	AccessReviewInstanceStatusCompleted     AccessReviewInstanceStatus = "Completed"
	AccessReviewInstanceStatusCompleting    AccessReviewInstanceStatus = "Completing"
	AccessReviewInstanceStatusInProgress    AccessReviewInstanceStatus = "InProgress"
	AccessReviewInstanceStatusInitializing  AccessReviewInstanceStatus = "Initializing"
	AccessReviewInstanceStatusNotStarted    AccessReviewInstanceStatus = "NotStarted"
	AccessReviewInstanceStatusScheduled     AccessReviewInstanceStatus = "Scheduled"
	AccessReviewInstanceStatusStarting      AccessReviewInstanceStatus = "Starting"
)

// PossibleAccessReviewInstanceStatusValues returns the possible values for the AccessReviewInstanceStatus const type.
func PossibleAccessReviewInstanceStatusValues() []AccessReviewInstanceStatus {
	return []AccessReviewInstanceStatus{
		AccessReviewInstanceStatusApplied,
		AccessReviewInstanceStatusApplying,
		AccessReviewInstanceStatusAutoReviewed,
		AccessReviewInstanceStatusAutoReviewing,
		AccessReviewInstanceStatusCompleted,
		AccessReviewInstanceStatusCompleting,
		AccessReviewInstanceStatusInProgress,
		AccessReviewInstanceStatusInitializing,
		AccessReviewInstanceStatusNotStarted,
		AccessReviewInstanceStatusScheduled,
		AccessReviewInstanceStatusStarting,
	}
}

// AccessReviewRecurrencePatternType - The recurrence type : weekly, monthly, etc.
type AccessReviewRecurrencePatternType string

const (
	AccessReviewRecurrencePatternTypeAbsoluteMonthly AccessReviewRecurrencePatternType = "absoluteMonthly"
	AccessReviewRecurrencePatternTypeWeekly          AccessReviewRecurrencePatternType = "weekly"
)

// PossibleAccessReviewRecurrencePatternTypeValues returns the possible values for the AccessReviewRecurrencePatternType const type.
func PossibleAccessReviewRecurrencePatternTypeValues() []AccessReviewRecurrencePatternType {
	return []AccessReviewRecurrencePatternType{
		AccessReviewRecurrencePatternTypeAbsoluteMonthly,
		AccessReviewRecurrencePatternTypeWeekly,
	}
}

// AccessReviewRecurrenceRangeType - The recurrence range type. The possible values are: endDate, noEnd, numbered.
type AccessReviewRecurrenceRangeType string

const (
	AccessReviewRecurrenceRangeTypeEndDate  AccessReviewRecurrenceRangeType = "endDate"
	AccessReviewRecurrenceRangeTypeNoEnd    AccessReviewRecurrenceRangeType = "noEnd"
	AccessReviewRecurrenceRangeTypeNumbered AccessReviewRecurrenceRangeType = "numbered"
)

// PossibleAccessReviewRecurrenceRangeTypeValues returns the possible values for the AccessReviewRecurrenceRangeType const type.
func PossibleAccessReviewRecurrenceRangeTypeValues() []AccessReviewRecurrenceRangeType {
	return []AccessReviewRecurrenceRangeType{
		AccessReviewRecurrenceRangeTypeEndDate,
		AccessReviewRecurrenceRangeTypeNoEnd,
		AccessReviewRecurrenceRangeTypeNumbered,
	}
}

// AccessReviewResult - Represents a reviewer's decision for a given review
type AccessReviewResult string

const (
	AccessReviewResultApprove     AccessReviewResult = "Approve"
	AccessReviewResultDeny        AccessReviewResult = "Deny"
	AccessReviewResultDontKnow    AccessReviewResult = "DontKnow"
	AccessReviewResultNotNotified AccessReviewResult = "NotNotified"
	AccessReviewResultNotReviewed AccessReviewResult = "NotReviewed"
)

// PossibleAccessReviewResultValues returns the possible values for the AccessReviewResult const type.
func PossibleAccessReviewResultValues() []AccessReviewResult {
	return []AccessReviewResult{
		AccessReviewResultApprove,
		AccessReviewResultDeny,
		AccessReviewResultDontKnow,
		AccessReviewResultNotNotified,
		AccessReviewResultNotReviewed,
	}
}

// AccessReviewReviewerType - The identity type : user/servicePrincipal
type AccessReviewReviewerType string

const (
	AccessReviewReviewerTypeServicePrincipal AccessReviewReviewerType = "servicePrincipal"
	AccessReviewReviewerTypeUser             AccessReviewReviewerType = "user"
)

// PossibleAccessReviewReviewerTypeValues returns the possible values for the AccessReviewReviewerType const type.
func PossibleAccessReviewReviewerTypeValues() []AccessReviewReviewerType {
	return []AccessReviewReviewerType{
		AccessReviewReviewerTypeServicePrincipal,
		AccessReviewReviewerTypeUser,
	}
}

// AccessReviewScheduleDefinitionReviewersType - This field specifies the type of reviewers for a review. Usually for a review,
// reviewers are explicitly assigned. However, in some cases, the reviewers may not be assigned and instead be chosen
// dynamically. For example managers review or self review.
type AccessReviewScheduleDefinitionReviewersType string

const (
	AccessReviewScheduleDefinitionReviewersTypeAssigned AccessReviewScheduleDefinitionReviewersType = "Assigned"
	AccessReviewScheduleDefinitionReviewersTypeManagers AccessReviewScheduleDefinitionReviewersType = "Managers"
	AccessReviewScheduleDefinitionReviewersTypeSelf     AccessReviewScheduleDefinitionReviewersType = "Self"
)

// PossibleAccessReviewScheduleDefinitionReviewersTypeValues returns the possible values for the AccessReviewScheduleDefinitionReviewersType const type.
func PossibleAccessReviewScheduleDefinitionReviewersTypeValues() []AccessReviewScheduleDefinitionReviewersType {
	return []AccessReviewScheduleDefinitionReviewersType{
		AccessReviewScheduleDefinitionReviewersTypeAssigned,
		AccessReviewScheduleDefinitionReviewersTypeManagers,
		AccessReviewScheduleDefinitionReviewersTypeSelf,
	}
}

// AccessReviewScheduleDefinitionStatus - This read-only field specifies the status of an accessReview.
type AccessReviewScheduleDefinitionStatus string

const (
	AccessReviewScheduleDefinitionStatusApplied       AccessReviewScheduleDefinitionStatus = "Applied"
	AccessReviewScheduleDefinitionStatusApplying      AccessReviewScheduleDefinitionStatus = "Applying"
	AccessReviewScheduleDefinitionStatusAutoReviewed  AccessReviewScheduleDefinitionStatus = "AutoReviewed"
	AccessReviewScheduleDefinitionStatusAutoReviewing AccessReviewScheduleDefinitionStatus = "AutoReviewing"
	AccessReviewScheduleDefinitionStatusCompleted     AccessReviewScheduleDefinitionStatus = "Completed"
	AccessReviewScheduleDefinitionStatusCompleting    AccessReviewScheduleDefinitionStatus = "Completing"
	AccessReviewScheduleDefinitionStatusInProgress    AccessReviewScheduleDefinitionStatus = "InProgress"
	AccessReviewScheduleDefinitionStatusInitializing  AccessReviewScheduleDefinitionStatus = "Initializing"
	AccessReviewScheduleDefinitionStatusNotStarted    AccessReviewScheduleDefinitionStatus = "NotStarted"
	AccessReviewScheduleDefinitionStatusScheduled     AccessReviewScheduleDefinitionStatus = "Scheduled"
	AccessReviewScheduleDefinitionStatusStarting      AccessReviewScheduleDefinitionStatus = "Starting"
)

// PossibleAccessReviewScheduleDefinitionStatusValues returns the possible values for the AccessReviewScheduleDefinitionStatus const type.
func PossibleAccessReviewScheduleDefinitionStatusValues() []AccessReviewScheduleDefinitionStatus {
	return []AccessReviewScheduleDefinitionStatus{
		AccessReviewScheduleDefinitionStatusApplied,
		AccessReviewScheduleDefinitionStatusApplying,
		AccessReviewScheduleDefinitionStatusAutoReviewed,
		AccessReviewScheduleDefinitionStatusAutoReviewing,
		AccessReviewScheduleDefinitionStatusCompleted,
		AccessReviewScheduleDefinitionStatusCompleting,
		AccessReviewScheduleDefinitionStatusInProgress,
		AccessReviewScheduleDefinitionStatusInitializing,
		AccessReviewScheduleDefinitionStatusNotStarted,
		AccessReviewScheduleDefinitionStatusScheduled,
		AccessReviewScheduleDefinitionStatusStarting,
	}
}

// AccessReviewScopeAssignmentState - The role assignment state eligible/active to review
type AccessReviewScopeAssignmentState string

const (
	AccessReviewScopeAssignmentStateActive   AccessReviewScopeAssignmentState = "active"
	AccessReviewScopeAssignmentStateEligible AccessReviewScopeAssignmentState = "eligible"
)

// PossibleAccessReviewScopeAssignmentStateValues returns the possible values for the AccessReviewScopeAssignmentState const type.
func PossibleAccessReviewScopeAssignmentStateValues() []AccessReviewScopeAssignmentState {
	return []AccessReviewScopeAssignmentState{
		AccessReviewScopeAssignmentStateActive,
		AccessReviewScopeAssignmentStateEligible,
	}
}

// AccessReviewScopePrincipalType - The identity type user/servicePrincipal to review
type AccessReviewScopePrincipalType string

const (
	AccessReviewScopePrincipalTypeGuestUser         AccessReviewScopePrincipalType = "guestUser"
	AccessReviewScopePrincipalTypeRedeemedGuestUser AccessReviewScopePrincipalType = "redeemedGuestUser"
	AccessReviewScopePrincipalTypeServicePrincipal  AccessReviewScopePrincipalType = "servicePrincipal"
	AccessReviewScopePrincipalTypeUser              AccessReviewScopePrincipalType = "user"
	AccessReviewScopePrincipalTypeUserGroup         AccessReviewScopePrincipalType = "user,group"
)

// PossibleAccessReviewScopePrincipalTypeValues returns the possible values for the AccessReviewScopePrincipalType const type.
func PossibleAccessReviewScopePrincipalTypeValues() []AccessReviewScopePrincipalType {
	return []AccessReviewScopePrincipalType{
		AccessReviewScopePrincipalTypeGuestUser,
		AccessReviewScopePrincipalTypeRedeemedGuestUser,
		AccessReviewScopePrincipalTypeServicePrincipal,
		AccessReviewScopePrincipalTypeUser,
		AccessReviewScopePrincipalTypeUserGroup,
	}
}

// DecisionResourceType - The type of resource: azureRole
type DecisionResourceType string

const (
	DecisionResourceTypeAzureRole DecisionResourceType = "azureRole"
)

// PossibleDecisionResourceTypeValues returns the possible values for the DecisionResourceType const type.
func PossibleDecisionResourceTypeValues() []DecisionResourceType {
	return []DecisionResourceType{
		DecisionResourceTypeAzureRole,
	}
}

// DecisionTargetType - The type of decision target : User/ServicePrincipal
type DecisionTargetType string

const (
	DecisionTargetTypeServicePrincipal DecisionTargetType = "servicePrincipal"
	DecisionTargetTypeUser             DecisionTargetType = "user"
)

// PossibleDecisionTargetTypeValues returns the possible values for the DecisionTargetType const type.
func PossibleDecisionTargetTypeValues() []DecisionTargetType {
	return []DecisionTargetType{
		DecisionTargetTypeServicePrincipal,
		DecisionTargetTypeUser,
	}
}

// DefaultDecisionType - This specifies the behavior for the autoReview feature when an access review completes.
type DefaultDecisionType string

const (
	DefaultDecisionTypeApprove        DefaultDecisionType = "Approve"
	DefaultDecisionTypeDeny           DefaultDecisionType = "Deny"
	DefaultDecisionTypeRecommendation DefaultDecisionType = "Recommendation"
)

// PossibleDefaultDecisionTypeValues returns the possible values for the DefaultDecisionType const type.
func PossibleDefaultDecisionTypeValues() []DefaultDecisionType {
	return []DefaultDecisionType{
		DefaultDecisionTypeApprove,
		DefaultDecisionTypeDeny,
		DefaultDecisionTypeRecommendation,
	}
}
