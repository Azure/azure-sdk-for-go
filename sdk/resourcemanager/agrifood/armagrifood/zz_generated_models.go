//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armagrifood

import "time"

// CheckNameAvailabilityRequest - The check availability request body.
type CheckNameAvailabilityRequest struct {
	// The name of the resource for which availability needs to be checked.
	Name *string `json:"name,omitempty"`

	// The resource type.
	Type *string `json:"type,omitempty"`
}

// CheckNameAvailabilityResponse - The check availability result.
type CheckNameAvailabilityResponse struct {
	// Detailed reason why the given name is available.
	Message *string `json:"message,omitempty"`

	// Indicates if the resource name is available.
	NameAvailable *bool `json:"nameAvailable,omitempty"`

	// The reason why the given name is not available.
	Reason *CheckNameAvailabilityReason `json:"reason,omitempty"`
}

// DetailedInformation - Model to capture detailed information for farmBeatsExtensions.
type DetailedInformation struct {
	// List of apiInputParameters.
	APIInputParameters []*string `json:"apiInputParameters,omitempty"`

	// ApiName available for the farmBeatsExtension.
	APIName *string `json:"apiName,omitempty"`

	// List of customParameters.
	CustomParameters []*string `json:"customParameters,omitempty"`

	// List of platformParameters.
	PlatformParameters []*string `json:"platformParameters,omitempty"`

	// Unit systems info for the data provider.
	UnitsSupported *UnitSystemsInfo `json:"unitsSupported,omitempty"`
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info interface{} `json:"info,omitempty" azure:"ro"`

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

// Extension resource.
type Extension struct {
	// Extension resource properties.
	Properties *ExtensionProperties `json:"properties,omitempty"`

	// READ-ONLY; The ETag value to implement optimistic concurrency.
	ETag *string `json:"eTag,omitempty" azure:"ro"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ExtensionListResponse - Paged response contains list of requested objects and a URL link to get the next set of results.
type ExtensionListResponse struct {
	// List of requested objects.
	Value []*Extension `json:"value,omitempty"`

	// READ-ONLY; Continuation link (absolute URI) to the next page of results in the list.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`
}

// ExtensionProperties - Extension resource properties.
type ExtensionProperties struct {
	// READ-ONLY; Extension api docs link.
	ExtensionAPIDocsLink *string `json:"extensionApiDocsLink,omitempty" azure:"ro"`

	// READ-ONLY; Extension auth link.
	ExtensionAuthLink *string `json:"extensionAuthLink,omitempty" azure:"ro"`

	// READ-ONLY; Extension category. e.g. weather/sensor/satellite.
	ExtensionCategory *string `json:"extensionCategory,omitempty" azure:"ro"`

	// READ-ONLY; Extension Id.
	ExtensionID *string `json:"extensionId,omitempty" azure:"ro"`

	// READ-ONLY; Installed extension version.
	InstalledExtensionVersion *string `json:"installedExtensionVersion,omitempty" azure:"ro"`
}

// ExtensionsClientCreateOptions contains the optional parameters for the ExtensionsClient.Create method.
type ExtensionsClientCreateOptions struct {
	// placeholder for future optional parameters
}

// ExtensionsClientDeleteOptions contains the optional parameters for the ExtensionsClient.Delete method.
type ExtensionsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// ExtensionsClientGetOptions contains the optional parameters for the ExtensionsClient.Get method.
type ExtensionsClientGetOptions struct {
	// placeholder for future optional parameters
}

// ExtensionsClientListByFarmBeatsOptions contains the optional parameters for the ExtensionsClient.ListByFarmBeats method.
type ExtensionsClientListByFarmBeatsOptions struct {
	// Installed extension categories.
	ExtensionCategories []string
	// Installed extension ids.
	ExtensionIDs []string
	// Maximum number of items needed (inclusive). Minimum = 10, Maximum = 1000, Default value = 50.
	MaxPageSize *int32
	// Skip token for getting next set of results.
	SkipToken *string
}

// ExtensionsClientUpdateOptions contains the optional parameters for the ExtensionsClient.Update method.
type ExtensionsClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// FarmBeats ARM Resource.
type FarmBeats struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// FarmBeats ARM Resource properties.
	Properties *FarmBeatsProperties `json:"properties,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// FarmBeatsExtension - FarmBeats extension resource.
type FarmBeatsExtension struct {
	// FarmBeatsExtension properties.
	Properties *FarmBeatsExtensionProperties `json:"properties,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// FarmBeatsExtensionListResponse - Paged response contains list of requested objects and a URL link to get the next set of
// results.
type FarmBeatsExtensionListResponse struct {
	// List of requested objects.
	Value []*FarmBeatsExtension `json:"value,omitempty"`

	// READ-ONLY; Continuation link (absolute URI) to the next page of results in the list.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`
}

// FarmBeatsExtensionProperties - FarmBeatsExtension properties.
type FarmBeatsExtensionProperties struct {
	// READ-ONLY; Textual description.
	Description *string `json:"description,omitempty" azure:"ro"`

	// READ-ONLY; Detailed information which shows summary of requested data. Used in descriptive get extension metadata call.
	// Information for weather category per api included are apisSupported, customParameters,
	// PlatformParameters and Units supported.
	DetailedInformation []*DetailedInformation `json:"detailedInformation,omitempty" azure:"ro"`

	// READ-ONLY; FarmBeatsExtension api docs link.
	ExtensionAPIDocsLink *string `json:"extensionApiDocsLink,omitempty" azure:"ro"`

	// READ-ONLY; FarmBeatsExtension auth link.
	ExtensionAuthLink *string `json:"extensionAuthLink,omitempty" azure:"ro"`

	// READ-ONLY; Category of the extension. e.g. weather/sensor/satellite.
	ExtensionCategory *string `json:"extensionCategory,omitempty" azure:"ro"`

	// READ-ONLY; FarmBeatsExtension ID.
	FarmBeatsExtensionID *string `json:"farmBeatsExtensionId,omitempty" azure:"ro"`

	// READ-ONLY; FarmBeatsExtension name.
	FarmBeatsExtensionName *string `json:"farmBeatsExtensionName,omitempty" azure:"ro"`

	// READ-ONLY; FarmBeatsExtension version.
	FarmBeatsExtensionVersion *string `json:"farmBeatsExtensionVersion,omitempty" azure:"ro"`

	// READ-ONLY; Publisher ID.
	PublisherID *string `json:"publisherId,omitempty" azure:"ro"`

	// READ-ONLY; Target ResourceType of the farmBeatsExtension.
	TargetResourceType *string `json:"targetResourceType,omitempty" azure:"ro"`
}

// FarmBeatsExtensionsClientGetOptions contains the optional parameters for the FarmBeatsExtensionsClient.Get method.
type FarmBeatsExtensionsClientGetOptions struct {
	// placeholder for future optional parameters
}

// FarmBeatsExtensionsClientListOptions contains the optional parameters for the FarmBeatsExtensionsClient.List method.
type FarmBeatsExtensionsClientListOptions struct {
	// Extension categories.
	ExtensionCategories []string
	// FarmBeatsExtension ids.
	FarmBeatsExtensionIDs []string
	// FarmBeats extension names.
	FarmBeatsExtensionNames []string
	// Maximum number of items needed (inclusive). Minimum = 10, Maximum = 1000, Default value = 50.
	MaxPageSize *int32
	// Publisher ids.
	PublisherIDs []string
}

// FarmBeatsListResponse - Paged response contains list of requested objects and a URL link to get the next set of results.
type FarmBeatsListResponse struct {
	// List of requested objects.
	Value []*FarmBeats `json:"value,omitempty"`

	// READ-ONLY; Continuation link (absolute URI) to the next page of results in the list.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`
}

// FarmBeatsModelsClientCreateOrUpdateOptions contains the optional parameters for the FarmBeatsModelsClient.CreateOrUpdate
// method.
type FarmBeatsModelsClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// FarmBeatsModelsClientDeleteOptions contains the optional parameters for the FarmBeatsModelsClient.Delete method.
type FarmBeatsModelsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// FarmBeatsModelsClientGetOptions contains the optional parameters for the FarmBeatsModelsClient.Get method.
type FarmBeatsModelsClientGetOptions struct {
	// placeholder for future optional parameters
}

// FarmBeatsModelsClientListByResourceGroupOptions contains the optional parameters for the FarmBeatsModelsClient.ListByResourceGroup
// method.
type FarmBeatsModelsClientListByResourceGroupOptions struct {
	// Maximum number of items needed (inclusive). Minimum = 10, Maximum = 1000, Default value = 50.
	MaxPageSize *int32
	// Continuation token for getting next set of results.
	SkipToken *string
}

// FarmBeatsModelsClientListBySubscriptionOptions contains the optional parameters for the FarmBeatsModelsClient.ListBySubscription
// method.
type FarmBeatsModelsClientListBySubscriptionOptions struct {
	// Maximum number of items needed (inclusive). Minimum = 10, Maximum = 1000, Default value = 50.
	MaxPageSize *int32
	// Skip token for getting next set of results.
	SkipToken *string
}

// FarmBeatsModelsClientUpdateOptions contains the optional parameters for the FarmBeatsModelsClient.Update method.
type FarmBeatsModelsClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// FarmBeatsProperties - FarmBeats ARM Resource properties.
type FarmBeatsProperties struct {
	// READ-ONLY; Uri of the FarmBeats instance.
	InstanceURI *string `json:"instanceUri,omitempty" azure:"ro"`

	// READ-ONLY; FarmBeats instance provisioning state.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// FarmBeatsUpdateRequestModel - FarmBeats update request.
type FarmBeatsUpdateRequestModel struct {
	// Geo-location where the resource lives.
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`
}

// LocationsClientCheckNameAvailabilityOptions contains the optional parameters for the LocationsClient.CheckNameAvailability
// method.
type LocationsClientCheckNameAvailabilityOptions struct {
	// placeholder for future optional parameters
}

// Operation - Details of a REST API operation, returned from the Resource Provider Operations API
type Operation struct {
	// Localized display information for this particular operation.
	Display *OperationDisplay `json:"display,omitempty"`

	// READ-ONLY; Enum. Indicates the action type. "Internal" refers to actions that are for internal only APIs.
	ActionType *ActionType `json:"actionType,omitempty" azure:"ro"`

	// READ-ONLY; Whether the operation applies to data-plane. This is "true" for data-plane operations and "false" for ARM/control-plane
	// operations.
	IsDataAction *bool `json:"isDataAction,omitempty" azure:"ro"`

	// READ-ONLY; The name of the operation, as per Resource-Based Access Control (RBAC). Examples: "Microsoft.Compute/virtualMachines/write",
	// "Microsoft.Compute/virtualMachines/capture/action"
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The intended executor of the operation; as in Resource Based Access Control (RBAC) and audit logs UX. Default
	// value is "user,system"
	Origin *Origin `json:"origin,omitempty" azure:"ro"`
}

// OperationDisplay - Localized display information for this particular operation.
type OperationDisplay struct {
	// READ-ONLY; The short, localized friendly description of the operation; suitable for tool tips and detailed views.
	Description *string `json:"description,omitempty" azure:"ro"`

	// READ-ONLY; The concise, localized friendly name for the operation; suitable for dropdowns. E.g. "Create or Update Virtual
	// Machine", "Restart Virtual Machine".
	Operation *string `json:"operation,omitempty" azure:"ro"`

	// READ-ONLY; The localized friendly form of the resource provider name, e.g. "Microsoft Monitoring Insights" or "Microsoft
	// Compute".
	Provider *string `json:"provider,omitempty" azure:"ro"`

	// READ-ONLY; The localized friendly name of the resource type related to this operation. E.g. "Virtual Machines" or "Job
	// Schedule Collections".
	Resource *string `json:"resource,omitempty" azure:"ro"`
}

// OperationListResult - A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to
// get the next set of results.
type OperationListResult struct {
	// READ-ONLY; URL to get the next set of operation list results (if there are any).
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; List of operations supported by the resource provider
	Value []*Operation `json:"value,omitempty" azure:"ro"`
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.List method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
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

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// The identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`

	// The type of identity that created the resource.
	CreatedByType *CreatedByType `json:"createdByType,omitempty"`

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`

	// The identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType `json:"lastModifiedByType,omitempty"`
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

// UnitSystemsInfo - Unit systems info for the data provider.
type UnitSystemsInfo struct {
	// REQUIRED; UnitSystem key sent as part of ProviderInput.
	Key *string `json:"key,omitempty"`

	// REQUIRED; List of unit systems supported by this data provider.
	Values []*string `json:"values,omitempty"`
}
