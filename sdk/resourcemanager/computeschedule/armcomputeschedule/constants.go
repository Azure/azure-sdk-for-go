//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcomputeschedule

const (
	moduleName    = "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/computeschedule/armcomputeschedule"
	moduleVersion = "v1.0.0"
)

// ActionType - Enum. Indicates the action type. "Internal" refers to actions that are for internal only APIs.
type ActionType string

const (
	ActionTypeInternal ActionType = "Internal"
)

// PossibleActionTypeValues returns the possible values for the ActionType const type.
func PossibleActionTypeValues() []ActionType {
	return []ActionType{
		ActionTypeInternal,
	}
}

// DeadlineType - The types of deadlines supported by ScheduledActions
type DeadlineType string

const (
	// DeadlineTypeCompleteBy - Complete the operation by the given deadline.
	DeadlineTypeCompleteBy DeadlineType = "CompleteBy"
	// DeadlineTypeInitiateAt - Initiate the operation at the given deadline.
	DeadlineTypeInitiateAt DeadlineType = "InitiateAt"
	// DeadlineTypeUnknown - Default value of Unknown.
	DeadlineTypeUnknown DeadlineType = "Unknown"
)

// PossibleDeadlineTypeValues returns the possible values for the DeadlineType const type.
func PossibleDeadlineTypeValues() []DeadlineType {
	return []DeadlineType{
		DeadlineTypeCompleteBy,
		DeadlineTypeInitiateAt,
		DeadlineTypeUnknown,
	}
}

// OperationState - Values that define the states of operations in Scheduled Actions
type OperationState string

const (
	// OperationStateBlocked - Operations that are blocked
	OperationStateBlocked OperationState = "Blocked"
	// OperationStateCancelled - Operations that have been Cancelled by the user
	OperationStateCancelled OperationState = "Cancelled"
	// OperationStateExecuting - Operations that are in the process of being executed
	OperationStateExecuting OperationState = "Executing"
	// OperationStateFailed - Operations that have failed
	OperationStateFailed OperationState = "Failed"
	// OperationStatePendingExecution - Operations that are waiting to be executed
	OperationStatePendingExecution OperationState = "PendingExecution"
	// OperationStatePendingScheduling - Operations that are pending scheduling
	OperationStatePendingScheduling OperationState = "PendingScheduling"
	// OperationStateScheduled - Operations that have been scheduled
	OperationStateScheduled OperationState = "Scheduled"
	// OperationStateSucceeded - Operations that suceeded
	OperationStateSucceeded OperationState = "Succeeded"
	// OperationStateUnknown - The default value for the operation state enum
	OperationStateUnknown OperationState = "Unknown"
)

// PossibleOperationStateValues returns the possible values for the OperationState const type.
func PossibleOperationStateValues() []OperationState {
	return []OperationState{
		OperationStateBlocked,
		OperationStateCancelled,
		OperationStateExecuting,
		OperationStateFailed,
		OperationStatePendingExecution,
		OperationStatePendingScheduling,
		OperationStateScheduled,
		OperationStateSucceeded,
		OperationStateUnknown,
	}
}

// OptimizationPreference - The preferences customers can select to optimize their requests to ScheduledActions
type OptimizationPreference string

const (
	// OptimizationPreferenceAvailability - Optimize while considering availability of resources
	OptimizationPreferenceAvailability OptimizationPreference = "Availability"
	// OptimizationPreferenceCost - Optimize while considering cost savings
	OptimizationPreferenceCost OptimizationPreference = "Cost"
	// OptimizationPreferenceCostAvailabilityBalanced - Optimize while considering a balance of cost and availability
	OptimizationPreferenceCostAvailabilityBalanced OptimizationPreference = "CostAvailabilityBalanced"
)

// PossibleOptimizationPreferenceValues returns the possible values for the OptimizationPreference const type.
func PossibleOptimizationPreferenceValues() []OptimizationPreference {
	return []OptimizationPreference{
		OptimizationPreferenceAvailability,
		OptimizationPreferenceCost,
		OptimizationPreferenceCostAvailabilityBalanced,
	}
}

// Origin - The intended executor of the operation; as in Resource Based Access Control (RBAC) and audit logs UX. Default
// value is "user,system"
type Origin string

const (
	OriginSystem     Origin = "system"
	OriginUser       Origin = "user"
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

// ResourceOperationType - Type of operation performed on the resources
type ResourceOperationType string

const (
	// ResourceOperationTypeDeallocate - Deallocate operations on the resources
	ResourceOperationTypeDeallocate ResourceOperationType = "Deallocate"
	// ResourceOperationTypeHibernate - Hibernate operations on the resources
	ResourceOperationTypeHibernate ResourceOperationType = "Hibernate"
	// ResourceOperationTypeStart - Start operations on the resources
	ResourceOperationTypeStart ResourceOperationType = "Start"
	// ResourceOperationTypeUnknown - The default value for this enum type
	ResourceOperationTypeUnknown ResourceOperationType = "Unknown"
)

// PossibleResourceOperationTypeValues returns the possible values for the ResourceOperationType const type.
func PossibleResourceOperationTypeValues() []ResourceOperationType {
	return []ResourceOperationType{
		ResourceOperationTypeDeallocate,
		ResourceOperationTypeHibernate,
		ResourceOperationTypeStart,
		ResourceOperationTypeUnknown,
	}
}
