//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcustomerlockbox

import "time"

// Approval - Request content object, in the use of Approve or Deny a Lockbox request.
type Approval struct {
	// Reason of the decision
	Reason *string `json:"reason,omitempty"`

	// Approval decision to the Lockbox request.
	Status *Status `json:"status,omitempty"`
}

// ErrorAdditionalInfo - An error additional info for the Lockbox service.
type ErrorAdditionalInfo struct {
	// Additional information about the request that is in error state.
	Info *ErrorAdditionalInfoInfo `json:"info,omitempty"`

	// The type of error info.
	Type *string `json:"type,omitempty"`
}

// ErrorAdditionalInfoInfo - Additional information about the request that is in error state.
type ErrorAdditionalInfoInfo struct {
	// The current status/state of the request quired.
	CurrentStatus *Status `json:"currentStatus,omitempty"`
}

// ErrorBody - An error response body from the Lockbox service.
type ErrorBody struct {
	// A list of error details about the error.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty"`

	// An identifier for the error. Codes are invariant and are intended to be consumed programmatically.
	Code *string `json:"code,omitempty"`

	// A message describing the error, intended to be suitable for display in a user interface.
	Message *string `json:"message,omitempty"`

	// The target of the particular error. For example, the name of the property in error.
	Target *string `json:"target,omitempty"`
}

// ErrorResponse - An error response from the Lockbox service.
type ErrorResponse struct {
	// Detailed information about the error encountered.
	Error *ErrorBody `json:"error,omitempty"`
}

// GetClientTenantOptedInOptions contains the optional parameters for the GetClient.TenantOptedIn method.
type GetClientTenantOptedInOptions struct {
	// placeholder for future optional parameters
}

// LockboxRequestResponse - A Lockbox request response object, containing all information associated with the request.
type LockboxRequestResponse struct {
	// The properties that are associated with a lockbox request.
	Properties *LockboxRequestResponseProperties `json:"properties,omitempty"`

	// READ-ONLY; The Arm resource id of the Lockbox request.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the Lockbox request.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the Lockbox request.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// LockboxRequestResponseProperties - The properties that are associated with a lockbox request.
type LockboxRequestResponseProperties struct {
	// The status of the request.
	Status *Status `json:"status,omitempty"`

	// READ-ONLY; Access level for requested resource
	AccessLevel *string `json:"accessLevel,omitempty" azure:"ro"`

	// READ-ONLY; The creation time of the request.
	CreatedDateTime *time.Time `json:"createdDateTime,omitempty" azure:"ro"`

	// READ-ONLY; The duration of the request in hours.
	Duration *string `json:"duration,omitempty" azure:"ro"`

	// READ-ONLY; The expiration time of the request.
	ExpirationDateTime *time.Time `json:"expirationDateTime,omitempty" azure:"ro"`

	// READ-ONLY; The justification of the requestor.
	Justification *string `json:"justification,omitempty" azure:"ro"`

	// READ-ONLY; The Lockbox request ID.
	RequestID *string `json:"requestId,omitempty" azure:"ro"`

	// READ-ONLY; A list of resource IDs associated with the Lockbox request separated by ','.
	ResourceIDs *string `json:"resourceIds,omitempty" azure:"ro"`

	// READ-ONLY; The resource type of the requested resources.
	ResourceType *string `json:"resourceType,omitempty" azure:"ro"`

	// READ-ONLY; The subscription ID.
	SubscriptionID *string `json:"subscriptionId,omitempty" azure:"ro"`

	// READ-ONLY; The url of the support case.
	SupportCaseURL *string `json:"supportCaseUrl,omitempty" azure:"ro"`

	// READ-ONLY; The id of the support request associated.
	SupportRequest *string `json:"supportRequest,omitempty" azure:"ro"`

	// READ-ONLY; The support case system that was used to initiate the request.
	Workitemsource *string `json:"workitemsource,omitempty" azure:"ro"`
}

// Operation result model for ARM RP
type Operation struct {
	// READ-ONLY; Contains the localized display information for this particular operation / action.
	Display *OperationDisplay `json:"display,omitempty" azure:"ro"`

	// READ-ONLY; Gets or sets a value indicating whether it is a data plane action
	IsDataAction *string `json:"isDataAction,omitempty" azure:"ro"`

	// READ-ONLY; Gets or sets action name
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Gets or sets origin
	Origin *string `json:"origin,omitempty" azure:"ro"`

	// READ-ONLY; Gets or sets properties
	Properties *string `json:"properties,omitempty" azure:"ro"`
}

// OperationDisplay - Contains the localized display information for this particular operation / action.
type OperationDisplay struct {
	// READ-ONLY; The localized friendly description for the operation.
	Description *string `json:"description,omitempty" azure:"ro"`

	// READ-ONLY; The localized friendly name for the operation.
	Operation *string `json:"operation,omitempty" azure:"ro"`

	// READ-ONLY; The localized friendly form of the resource provider name.
	Provider *string `json:"provider,omitempty" azure:"ro"`

	// READ-ONLY; The localized friendly form of the resource type related to this action/operation.
	Resource *string `json:"resource,omitempty" azure:"ro"`
}

// OperationListResult - Result of the request to list Customer Lockbox operations. It contains a list of operations.
type OperationListResult struct {
	// READ-ONLY; URL to get the next set of operation list results if there are any.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; List of Customer Lockbox operations supported by the Microsoft.StreamAnalytics resource provider.
	Value []*Operation `json:"value,omitempty" azure:"ro"`
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.List method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// PostClientDisableLockboxOptions contains the optional parameters for the PostClient.DisableLockbox method.
type PostClientDisableLockboxOptions struct {
	// placeholder for future optional parameters
}

// PostClientEnableLockboxOptions contains the optional parameters for the PostClient.EnableLockbox method.
type PostClientEnableLockboxOptions struct {
	// placeholder for future optional parameters
}

// RequestListResult - Object containing a list of streaming jobs.
type RequestListResult struct {
	// READ-ONLY; URL to get the next set of operation list results if there are any.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A list of Lockbox requests. Populated by a 'List' operation.
	Value []*LockboxRequestResponse `json:"value,omitempty" azure:"ro"`
}

// RequestsClientGetOptions contains the optional parameters for the RequestsClient.Get method.
type RequestsClientGetOptions struct {
	// placeholder for future optional parameters
}

// RequestsClientListOptions contains the optional parameters for the RequestsClient.List method.
type RequestsClientListOptions struct {
	// The $filter OData query parameter. Only filter by request status is supported, e.g $filter=properties/status eq 'Pending'
	Filter *string
}

// RequestsClientUpdateStatusOptions contains the optional parameters for the RequestsClient.UpdateStatus method.
type RequestsClientUpdateStatusOptions struct {
	// placeholder for future optional parameters
}

// TenantOptInResponse - TenantOptIn Response object
type TenantOptInResponse struct {
	// READ-ONLY; True if tenant is opted in, false otherwise
	IsOptedIn *bool `json:"isOptedIn,omitempty" azure:"ro"`
}
