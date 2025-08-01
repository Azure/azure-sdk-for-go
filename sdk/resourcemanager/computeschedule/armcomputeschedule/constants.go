// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armcomputeschedule

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

// Language - The notification languages currently supported
type Language string

const (
	// LanguageEnUs - American english language
	LanguageEnUs Language = "en-us"
)

// PossibleLanguageValues returns the possible values for the Language const type.
func PossibleLanguageValues() []Language {
	return []Language{
		LanguageEnUs,
	}
}

// Month - Representation of the months available selection in a gregorian calendar
type Month string

const (
	// MonthAll - All months
	MonthAll Month = "All"
	// MonthApril - The April month.
	MonthApril Month = "April"
	// MonthAugust - The August month.
	MonthAugust Month = "August"
	// MonthDecember - The December month.
	MonthDecember Month = "December"
	// MonthFebruary - The February month.
	MonthFebruary Month = "February"
	// MonthJanuary - The January month.
	MonthJanuary Month = "January"
	// MonthJuly - The July month.
	MonthJuly Month = "July"
	// MonthJune - The June month.
	MonthJune Month = "June"
	// MonthMarch - The March month.
	MonthMarch Month = "March"
	// MonthMay - The May month.
	MonthMay Month = "May"
	// MonthNovember - The November month.
	MonthNovember Month = "November"
	// MonthOctober - The October month.
	MonthOctober Month = "October"
	// MonthSeptember - The September month.
	MonthSeptember Month = "September"
)

// PossibleMonthValues returns the possible values for the Month const type.
func PossibleMonthValues() []Month {
	return []Month{
		MonthAll,
		MonthApril,
		MonthAugust,
		MonthDecember,
		MonthFebruary,
		MonthJanuary,
		MonthJuly,
		MonthJune,
		MonthMarch,
		MonthMay,
		MonthNovember,
		MonthOctober,
		MonthSeptember,
	}
}

// NotificationType - The type of notification supported
type NotificationType string

const (
	// NotificationTypeEmail - Notify through e-mail
	NotificationTypeEmail NotificationType = "Email"
)

// PossibleNotificationTypeValues returns the possible values for the NotificationType const type.
func PossibleNotificationTypeValues() []NotificationType {
	return []NotificationType{
		NotificationTypeEmail,
	}
}

// OccurrenceState - The state the occurrence is at a given time
type OccurrenceState string

const (
	// OccurrenceStateCanceled - The occurrence has been canceled
	OccurrenceStateCanceled OccurrenceState = "Canceled"
	// OccurrenceStateCancelling - The occurrence is going through cancellation
	OccurrenceStateCancelling OccurrenceState = "Cancelling"
	// OccurrenceStateCreated - The occurrence was created
	OccurrenceStateCreated OccurrenceState = "Created"
	// OccurrenceStateFailed - The occurrence has failed during its scheduling
	OccurrenceStateFailed OccurrenceState = "Failed"
	// OccurrenceStateRescheduling - The occurrence is being rescheduled
	OccurrenceStateRescheduling OccurrenceState = "Rescheduling"
	// OccurrenceStateScheduled - The occurrence has been scheduled
	OccurrenceStateScheduled OccurrenceState = "Scheduled"
	// OccurrenceStateSucceeded - The occurrence has successfully ran
	OccurrenceStateSucceeded OccurrenceState = "Succeeded"
)

// PossibleOccurrenceStateValues returns the possible values for the OccurrenceState const type.
func PossibleOccurrenceStateValues() []OccurrenceState {
	return []OccurrenceState{
		OccurrenceStateCanceled,
		OccurrenceStateCancelling,
		OccurrenceStateCreated,
		OccurrenceStateFailed,
		OccurrenceStateRescheduling,
		OccurrenceStateScheduled,
		OccurrenceStateSucceeded,
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

// ProvisioningState - Provisioning state
type ProvisioningState string

const (
	// ProvisioningStateCanceled - Resource creation was canceled.
	ProvisioningStateCanceled ProvisioningState = "Canceled"
	// ProvisioningStateDeleting - Resource is being deleted.
	ProvisioningStateDeleting ProvisioningState = "Deleting"
	// ProvisioningStateFailed - Resource creation failed.
	ProvisioningStateFailed ProvisioningState = "Failed"
	// ProvisioningStateSucceeded - Resource has been created.
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

// PossibleProvisioningStateValues returns the possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{
		ProvisioningStateCanceled,
		ProvisioningStateDeleting,
		ProvisioningStateFailed,
		ProvisioningStateSucceeded,
	}
}

// ResourceOperationStatus - The state the resource is on after the resource operation is applied
type ResourceOperationStatus string

const (
	// ResourceOperationStatusFailed - The resource operation has failed.
	ResourceOperationStatusFailed ResourceOperationStatus = "Failed"
	// ResourceOperationStatusSucceeded - The resource operation was successful
	ResourceOperationStatusSucceeded ResourceOperationStatus = "Succeeded"
)

// PossibleResourceOperationStatusValues returns the possible values for the ResourceOperationStatus const type.
func PossibleResourceOperationStatusValues() []ResourceOperationStatus {
	return []ResourceOperationStatus{
		ResourceOperationStatusFailed,
		ResourceOperationStatusSucceeded,
	}
}

// ResourceOperationType - The kind of operation types that can be performed on resources using ScheduledActions
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

// ResourceProvisioningState - The provisioning state of a resource type.
type ResourceProvisioningState string

const (
	// ResourceProvisioningStateCanceled - Resource creation was canceled.
	ResourceProvisioningStateCanceled ResourceProvisioningState = "Canceled"
	// ResourceProvisioningStateFailed - Resource creation failed.
	ResourceProvisioningStateFailed ResourceProvisioningState = "Failed"
	// ResourceProvisioningStateSucceeded - Resource has been created.
	ResourceProvisioningStateSucceeded ResourceProvisioningState = "Succeeded"
)

// PossibleResourceProvisioningStateValues returns the possible values for the ResourceProvisioningState const type.
func PossibleResourceProvisioningStateValues() []ResourceProvisioningState {
	return []ResourceProvisioningState{
		ResourceProvisioningStateCanceled,
		ResourceProvisioningStateFailed,
		ResourceProvisioningStateSucceeded,
	}
}

// ResourceType - The type of resource being targeted
type ResourceType string

const (
	// ResourceTypeVirtualMachine - Resources defined are Virtual Machines
	ResourceTypeVirtualMachine ResourceType = "VirtualMachine"
	// ResourceTypeVirtualMachineScaleSet - Resources defined are Virtual Machines Scale Sets
	ResourceTypeVirtualMachineScaleSet ResourceType = "VirtualMachineScaleSet"
)

// PossibleResourceTypeValues returns the possible values for the ResourceType const type.
func PossibleResourceTypeValues() []ResourceType {
	return []ResourceType{
		ResourceTypeVirtualMachine,
		ResourceTypeVirtualMachineScaleSet,
	}
}

// ScheduledActionType - Specify which action user wants to be performed on the resources
type ScheduledActionType string

const (
	// ScheduledActionTypeDeallocate - Perform a deallocate action on the specified resources
	ScheduledActionTypeDeallocate ScheduledActionType = "Deallocate"
	// ScheduledActionTypeHibernate - Perform hibernate and deallocate on the specified resources
	ScheduledActionTypeHibernate ScheduledActionType = "Hibernate"
	// ScheduledActionTypeStart - Perform a start action on the specified resources
	ScheduledActionTypeStart ScheduledActionType = "Start"
)

// PossibleScheduledActionTypeValues returns the possible values for the ScheduledActionType const type.
func PossibleScheduledActionTypeValues() []ScheduledActionType {
	return []ScheduledActionType{
		ScheduledActionTypeDeallocate,
		ScheduledActionTypeHibernate,
		ScheduledActionTypeStart,
	}
}

// WeekDay - Representation of the possible selection of days in a week in a gregorian calendar
type WeekDay string

const (
	// WeekDayAll - All week days
	WeekDayAll WeekDay = "All"
	// WeekDayFriday - Friday weekday.
	WeekDayFriday WeekDay = "Friday"
	// WeekDayMonday - Monday weekday.
	WeekDayMonday WeekDay = "Monday"
	// WeekDaySaturday - Saturday weekday.
	WeekDaySaturday WeekDay = "Saturday"
	// WeekDaySunday - Sunday weekday.
	WeekDaySunday WeekDay = "Sunday"
	// WeekDayThursday - Thursday weekday.
	WeekDayThursday WeekDay = "Thursday"
	// WeekDayTuesday - Tuesday weekday.
	WeekDayTuesday WeekDay = "Tuesday"
	// WeekDayWednesday - Wednesday weekday.
	WeekDayWednesday WeekDay = "Wednesday"
)

// PossibleWeekDayValues returns the possible values for the WeekDay const type.
func PossibleWeekDayValues() []WeekDay {
	return []WeekDay{
		WeekDayAll,
		WeekDayFriday,
		WeekDayMonday,
		WeekDaySaturday,
		WeekDaySunday,
		WeekDayThursday,
		WeekDayTuesday,
		WeekDayWednesday,
	}
}
