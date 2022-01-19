//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcustomproviders

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// Association - The resource definition of this association.
type Association struct {
	// The properties of the association.
	Properties *AssociationProperties `json:"properties,omitempty"`

	// READ-ONLY; The association id.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The association name.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The association type.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// AssociationProperties - The properties of the association.
type AssociationProperties struct {
	// The REST resource instance of the target resource for this association.
	TargetResourceID *string `json:"targetResourceId,omitempty"`

	// READ-ONLY; The provisioning state of the association.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// AssociationsClientBeginCreateOrUpdateOptions contains the optional parameters for the AssociationsClient.BeginCreateOrUpdate
// method.
type AssociationsClientBeginCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// AssociationsClientBeginDeleteOptions contains the optional parameters for the AssociationsClient.BeginDelete method.
type AssociationsClientBeginDeleteOptions struct {
	// placeholder for future optional parameters
}

// AssociationsClientGetOptions contains the optional parameters for the AssociationsClient.Get method.
type AssociationsClientGetOptions struct {
	// placeholder for future optional parameters
}

// AssociationsClientListAllOptions contains the optional parameters for the AssociationsClient.ListAll method.
type AssociationsClientListAllOptions struct {
	// placeholder for future optional parameters
}

// AssociationsList - List of associations.
type AssociationsList struct {
	// The URL to use for getting the next set of results.
	NextLink *string `json:"nextLink,omitempty"`

	// The array of associations.
	Value []*Association `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type AssociationsList.
func (a AssociationsList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", a.NextLink)
	populate(objectMap, "value", a.Value)
	return json.Marshal(objectMap)
}

// CustomRPActionRouteDefinition - The route definition for an action implemented by the custom resource provider.
type CustomRPActionRouteDefinition struct {
	// REQUIRED; The route definition endpoint URI that the custom resource provider will proxy requests to. This can be in the
	// form of a flat URI (e.g. 'https://testendpoint/') or can specify to route via a path
	// (e.g. 'https://testendpoint/{requestPath}')
	Endpoint *string `json:"endpoint,omitempty"`

	// REQUIRED; The name of the route definition. This becomes the name for the ARM extension (e.g.
	// '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CustomProviders/resourceProviders/{resourceProviderName}/{name}')
	Name *string `json:"name,omitempty"`

	// The routing types that are supported for action requests.
	RoutingType *ActionRouting `json:"routingType,omitempty"`
}

// CustomRPManifest - A manifest file that defines the custom resource provider resources.
type CustomRPManifest struct {
	// REQUIRED; Resource location
	Location *string `json:"location,omitempty"`

	// The manifest for the custom resource provider
	Properties *CustomRPManifestProperties `json:"properties,omitempty"`

	// Resource tags
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Resource Id
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; Resource name
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Resource type
	Type *string `json:"type,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type CustomRPManifest.
func (c CustomRPManifest) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", c.ID)
	populate(objectMap, "location", c.Location)
	populate(objectMap, "name", c.Name)
	populate(objectMap, "properties", c.Properties)
	populate(objectMap, "tags", c.Tags)
	populate(objectMap, "type", c.Type)
	return json.Marshal(objectMap)
}

// CustomRPManifestProperties - The manifest for the custom resource provider
type CustomRPManifestProperties struct {
	// A list of actions that the custom resource provider implements.
	Actions []*CustomRPActionRouteDefinition `json:"actions,omitempty"`

	// A list of resource types that the custom resource provider implements.
	ResourceTypes []*CustomRPResourceTypeRouteDefinition `json:"resourceTypes,omitempty"`

	// A list of validations to run on the custom resource provider's requests.
	Validations []*CustomRPValidations `json:"validations,omitempty"`

	// READ-ONLY; The provisioning state of the resource provider.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type CustomRPManifestProperties.
func (c CustomRPManifestProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "actions", c.Actions)
	populate(objectMap, "provisioningState", c.ProvisioningState)
	populate(objectMap, "resourceTypes", c.ResourceTypes)
	populate(objectMap, "validations", c.Validations)
	return json.Marshal(objectMap)
}

// CustomRPResourceTypeRouteDefinition - The route definition for a resource implemented by the custom resource provider.
type CustomRPResourceTypeRouteDefinition struct {
	// REQUIRED; The route definition endpoint URI that the custom resource provider will proxy requests to. This can be in the
	// form of a flat URI (e.g. 'https://testendpoint/') or can specify to route via a path
	// (e.g. 'https://testendpoint/{requestPath}')
	Endpoint *string `json:"endpoint,omitempty"`

	// REQUIRED; The name of the route definition. This becomes the name for the ARM extension (e.g.
	// '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CustomProviders/resourceProviders/{resourceProviderName}/{name}')
	Name *string `json:"name,omitempty"`

	// The routing types that are supported for resource requests.
	RoutingType *ResourceTypeRouting `json:"routingType,omitempty"`
}

// CustomRPRouteDefinition - A route definition that defines an action or resource that can be interacted with through the
// custom resource provider.
type CustomRPRouteDefinition struct {
	// REQUIRED; The route definition endpoint URI that the custom resource provider will proxy requests to. This can be in the
	// form of a flat URI (e.g. 'https://testendpoint/') or can specify to route via a path
	// (e.g. 'https://testendpoint/{requestPath}')
	Endpoint *string `json:"endpoint,omitempty"`

	// REQUIRED; The name of the route definition. This becomes the name for the ARM extension (e.g.
	// '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CustomProviders/resourceProviders/{resourceProviderName}/{name}')
	Name *string `json:"name,omitempty"`
}

// CustomRPValidations - A validation to apply on custom resource provider requests.
type CustomRPValidations struct {
	// REQUIRED; A link to the validation specification. The specification must be hosted on raw.githubusercontent.com.
	Specification *string `json:"specification,omitempty"`

	// The type of validation to run against a matching request.
	ValidationType *ValidationType `json:"validationType,omitempty"`
}

// CustomResourceProviderClientBeginCreateOrUpdateOptions contains the optional parameters for the CustomResourceProviderClient.BeginCreateOrUpdate
// method.
type CustomResourceProviderClientBeginCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// CustomResourceProviderClientBeginDeleteOptions contains the optional parameters for the CustomResourceProviderClient.BeginDelete
// method.
type CustomResourceProviderClientBeginDeleteOptions struct {
	// placeholder for future optional parameters
}

// CustomResourceProviderClientGetOptions contains the optional parameters for the CustomResourceProviderClient.Get method.
type CustomResourceProviderClientGetOptions struct {
	// placeholder for future optional parameters
}

// CustomResourceProviderClientListByResourceGroupOptions contains the optional parameters for the CustomResourceProviderClient.ListByResourceGroup
// method.
type CustomResourceProviderClientListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// CustomResourceProviderClientListBySubscriptionOptions contains the optional parameters for the CustomResourceProviderClient.ListBySubscription
// method.
type CustomResourceProviderClientListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// CustomResourceProviderClientUpdateOptions contains the optional parameters for the CustomResourceProviderClient.Update
// method.
type CustomResourceProviderClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// ErrorDefinition - Error definition.
type ErrorDefinition struct {
	// READ-ONLY; Service specific error code which serves as the substatus for the HTTP error code.
	Code *string `json:"code,omitempty" azure:"ro"`

	// READ-ONLY; Internal error details.
	Details []*ErrorDefinition `json:"details,omitempty" azure:"ro"`

	// READ-ONLY; Description of the error.
	Message *string `json:"message,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type ErrorDefinition.
func (e ErrorDefinition) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "message", e.Message)
	return json.Marshal(objectMap)
}

// ErrorResponse - Error response.
type ErrorResponse struct {
	// The error details.
	Error *ErrorDefinition `json:"error,omitempty"`
}

// ListByCustomRPManifest - List of custom resource providers.
type ListByCustomRPManifest struct {
	// The URL to use for getting the next set of results.
	NextLink *string `json:"nextLink,omitempty"`

	// The array of custom resource provider manifests.
	Value []*CustomRPManifest `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type ListByCustomRPManifest.
func (l ListByCustomRPManifest) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", l.NextLink)
	populate(objectMap, "value", l.Value)
	return json.Marshal(objectMap)
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.List method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// Resource - The resource definition.
type Resource struct {
	// REQUIRED; Resource location
	Location *string `json:"location,omitempty"`

	// Resource tags
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Resource Id
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; Resource name
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Resource type
	Type *string `json:"type,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type Resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", r.ID)
	populate(objectMap, "location", r.Location)
	populate(objectMap, "name", r.Name)
	populate(objectMap, "tags", r.Tags)
	populate(objectMap, "type", r.Type)
	return json.Marshal(objectMap)
}

// ResourceProviderOperation - Supported operations of this resource provider.
type ResourceProviderOperation struct {
	// Display metadata associated with the operation.
	Display *ResourceProviderOperationDisplay `json:"display,omitempty"`

	// Operation name, in format of {provider}/{resource}/{operation}
	Name *string `json:"name,omitempty"`
}

// ResourceProviderOperationDisplay - Display metadata associated with the operation.
type ResourceProviderOperationDisplay struct {
	// Description of this operation.
	Description *string `json:"description,omitempty"`

	// Type of operation: get, read, delete, etc.
	Operation *string `json:"operation,omitempty"`

	// Resource provider: Microsoft Custom Providers.
	Provider *string `json:"provider,omitempty"`

	// Resource on which the operation is performed.
	Resource *string `json:"resource,omitempty"`
}

// ResourceProviderOperationList - Results of the request to list operations.
type ResourceProviderOperationList struct {
	// The URL to use for getting the next set of results.
	NextLink *string `json:"nextLink,omitempty"`

	// List of operations supported by this resource provider.
	Value []*ResourceProviderOperation `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type ResourceProviderOperationList.
func (r ResourceProviderOperationList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", r.NextLink)
	populate(objectMap, "value", r.Value)
	return json.Marshal(objectMap)
}

// ResourceProvidersUpdate - custom resource provider update information.
type ResourceProvidersUpdate struct {
	// Resource tags
	Tags map[string]*string `json:"tags,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type ResourceProvidersUpdate.
func (r ResourceProvidersUpdate) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "tags", r.Tags)
	return json.Marshal(objectMap)
}

func populate(m map[string]interface{}, k string, v interface{}) {
	if v == nil {
		return
	} else if azcore.IsNullValue(v) {
		m[k] = nil
	} else if !reflect.ValueOf(v).IsNil() {
		m[k] = v
	}
}
