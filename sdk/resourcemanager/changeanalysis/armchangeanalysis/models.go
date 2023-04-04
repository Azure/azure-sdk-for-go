//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armchangeanalysis

import "time"

// Change - The detected change.
type Change struct {
	// The properties of a change.
	Properties *ChangeProperties `json:"properties,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ChangeList - The list of detected changes.
type ChangeList struct {
	// The URI that can be used to request the next page of changes.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of changes.
	Value []*Change `json:"value,omitempty"`
}

// ChangeProperties - The properties of a change.
type ChangeProperties struct {
	// The type of the change.
	ChangeType *ChangeType `json:"changeType,omitempty"`

	// The list of identities who might initiated the change. The identity could be user name (email address) or the object ID
	// of the Service Principal.
	InitiatedByList []*string `json:"initiatedByList,omitempty"`

	// The list of detailed changes at json property level.
	PropertyChanges []*PropertyChange `json:"propertyChanges,omitempty"`

	// The resource id that the change is attached to.
	ResourceID *string `json:"resourceId,omitempty"`

	// The time when the change is detected.
	TimeStamp *time.Time `json:"timeStamp,omitempty"`
}

// ChangesClientListChangesByResourceGroupOptions contains the optional parameters for the ChangesClient.NewListChangesByResourceGroupPager
// method.
type ChangesClientListChangesByResourceGroupOptions struct {
	// A skip token is used to continue retrieving items after an operation returns a partial result. If a previous response contains
	// a nextLink element, the value of the nextLink element will include a
	// skipToken parameter that specifies a starting point to use for subsequent calls.
	SkipToken *string
}

// ChangesClientListChangesBySubscriptionOptions contains the optional parameters for the ChangesClient.NewListChangesBySubscriptionPager
// method.
type ChangesClientListChangesBySubscriptionOptions struct {
	// A skip token is used to continue retrieving items after an operation returns a partial result. If a previous response contains
	// a nextLink element, the value of the nextLink element will include a
	// skipToken parameter that specifies a starting point to use for subsequent calls.
	SkipToken *string
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info any `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ErrorDetail - The error detail.
type ErrorDetail struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty" azure:"ro"`

	// READ-ONLY; The error code.
	Code *string `json:"code,omitempty" azure:"ro"`

	// READ-ONLY; The error details.
	Details []*ErrorDetail `json:"details,omitempty" azure:"ro"`

	// READ-ONLY; The error message.
	Message *string `json:"message,omitempty" azure:"ro"`

	// READ-ONLY; The error target.
	Target *string `json:"target,omitempty" azure:"ro"`
}

// ErrorResponse - Common error response for all Azure Resource Manager APIs to return error details for failed operations.
// (This also follows the OData error response format.).
type ErrorResponse struct {
	// The error object.
	Error *ErrorDetail `json:"error,omitempty"`
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.NewListPager method.
type OperationsClientListOptions struct {
	// A skip token is used to continue retrieving items after an operation returns a partial result. If a previous response contains
	// a nextLink element, the value of the nextLink element will include a
	// skipToken parameter that specifies a starting point to use for subsequent calls.
	SkipToken *string
}

// PropertyChange - Data of a property change.
type PropertyChange struct {
	// The change category.
	ChangeCategory *ChangeCategory `json:"changeCategory,omitempty"`

	// The type of the change.
	ChangeType *ChangeType `json:"changeType,omitempty"`

	// The description of the changed property.
	Description *string `json:"description,omitempty"`

	// The enhanced display name of the json path. E.g., the json path value[0].properties will be translated to something meaningful
	// like slots["Staging"].properties.
	DisplayName *string `json:"displayName,omitempty"`

	// The boolean indicating whether the oldValue and newValue are masked. The values are masked if it contains sensitive information
	// that the user doesn't have access to.
	IsDataMasked *bool `json:"isDataMasked,omitempty"`

	// The json path of the changed property.
	JSONPath *string `json:"jsonPath,omitempty"`
	Level    *Level  `json:"level,omitempty"`

	// The value of the property after the change.
	NewValue *string `json:"newValue,omitempty"`

	// The value of the property before the change.
	OldValue *string `json:"oldValue,omitempty"`
}

// ProxyResource - The resource model definition for a Azure Resource Manager proxy resource. It will not have tags and a
// location
type ProxyResource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ResourceChangesClientListOptions contains the optional parameters for the ResourceChangesClient.NewListPager method.
type ResourceChangesClientListOptions struct {
	// A skip token is used to continue retrieving items after an operation returns a partial result. If a previous response contains
	// a nextLink element, the value of the nextLink element will include a
	// skipToken parameter that specifies a starting point to use for subsequent calls.
	SkipToken *string
}

// ResourceProviderOperationDefinition - The resource provider operation definition.
type ResourceProviderOperationDefinition struct {
	// The resource provider operation details.
	Display *ResourceProviderOperationDisplay `json:"display,omitempty"`

	// The resource provider operation name.
	Name *string `json:"name,omitempty"`
}

// ResourceProviderOperationDisplay - The resource provider operation details.
type ResourceProviderOperationDisplay struct {
	// Description of the resource provider operation.
	Description *string `json:"description,omitempty"`

	// Name of the resource provider operation.
	Operation *string `json:"operation,omitempty"`

	// Name of the resource provider.
	Provider *string `json:"provider,omitempty"`

	// Name of the resource type.
	Resource *string `json:"resource,omitempty"`
}

// ResourceProviderOperationList - The resource provider operation list.
type ResourceProviderOperationList struct {
	// The URI that can be used to request the next page for list of Azure operations.
	NextLink *string `json:"nextLink,omitempty"`

	// Resource provider operations list.
	Value []*ResourceProviderOperationDefinition `json:"value,omitempty"`
}
