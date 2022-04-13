//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armwindowsesu

import "time"

// ErrorDefinition - Error definition.
type ErrorDefinition struct {
	// READ-ONLY; Service specific error code which serves as the substatus for the HTTP error code.
	Code *string `json:"code,omitempty" azure:"ro"`

	// READ-ONLY; Internal error details.
	Details []*ErrorDefinition `json:"details,omitempty" azure:"ro"`

	// READ-ONLY; Description of the error.
	Message *string `json:"message,omitempty" azure:"ro"`
}

// ErrorResponse - Error response.
type ErrorResponse struct {
	// The error details.
	Error *ErrorDefinition `json:"error,omitempty"`
}

// MultipleActivationKey - MAK key details.
type MultipleActivationKey struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// MAK key specific properties.
	Properties *MultipleActivationKeyProperties `json:"properties,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// MultipleActivationKeyList - List of MAK keys.
type MultipleActivationKeyList struct {
	// List of MAK keys.
	Value []*MultipleActivationKey `json:"value,omitempty"`

	// READ-ONLY; Link to the next page of resources.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`
}

// MultipleActivationKeyProperties - MAK key specific properties.
type MultipleActivationKeyProperties struct {
	// Agreement number under which the key is requested.
	AgreementNumber *string `json:"agreementNumber,omitempty"`

	// Number of activations/servers using the MAK key.
	InstalledServerNumber *int32 `json:"installedServerNumber,omitempty"`

	// true if user has eligible on-premises Windows physical or virtual machines, and that the requested key will only be used
	// in their organization; false otherwise.
	IsEligible *bool `json:"isEligible,omitempty"`

	// Type of OS for which the key is requested.
	OSType *OsType `json:"osType,omitempty"`

	// Type of support
	SupportType *SupportType `json:"supportType,omitempty"`

	// READ-ONLY; End of support of security updates activated by the MAK key.
	ExpirationDate *time.Time `json:"expirationDate,omitempty" azure:"ro"`

	// READ-ONLY; MAK 5x5 key.
	MultipleActivationKey *string `json:"multipleActivationKey,omitempty" azure:"ro"`

	// READ-ONLY
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// MultipleActivationKeyUpdate - MAK key details.
type MultipleActivationKeyUpdate struct {
	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`
}

// MultipleActivationKeysClientBeginCreateOptions contains the optional parameters for the MultipleActivationKeysClient.BeginCreate
// method.
type MultipleActivationKeysClientBeginCreateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// MultipleActivationKeysClientDeleteOptions contains the optional parameters for the MultipleActivationKeysClient.Delete
// method.
type MultipleActivationKeysClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// MultipleActivationKeysClientGetOptions contains the optional parameters for the MultipleActivationKeysClient.Get method.
type MultipleActivationKeysClientGetOptions struct {
	// placeholder for future optional parameters
}

// MultipleActivationKeysClientListByResourceGroupOptions contains the optional parameters for the MultipleActivationKeysClient.ListByResourceGroup
// method.
type MultipleActivationKeysClientListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// MultipleActivationKeysClientListOptions contains the optional parameters for the MultipleActivationKeysClient.List method.
type MultipleActivationKeysClientListOptions struct {
	// placeholder for future optional parameters
}

// MultipleActivationKeysClientUpdateOptions contains the optional parameters for the MultipleActivationKeysClient.Update
// method.
type MultipleActivationKeysClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// Operation - REST API operation details.
type Operation struct {
	// Meta data about operation used for display in portal.
	Display *OperationDisplay `json:"display,omitempty"`

	// READ-ONLY; Name of the operation.
	Name *string `json:"name,omitempty" azure:"ro"`
}

// OperationDisplay - Meta data about operation used for display in portal.
type OperationDisplay struct {
	Description *string `json:"description,omitempty"`
	Operation   *string `json:"operation,omitempty"`
	Provider    *string `json:"provider,omitempty"`
	Resource    *string `json:"resource,omitempty"`
}

// OperationList - List of available REST API operations.
type OperationList struct {
	// List of operations.
	Value []*Operation `json:"value,omitempty"`

	// READ-ONLY; Link to the next page of resources.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.List method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
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

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags'
// and a 'location'
type TrackedResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}
