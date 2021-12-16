//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armagrifood

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
	"time"
)

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

// MarshalJSON implements the json.Marshaller interface for type DetailedInformation.
func (d DetailedInformation) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "apiInputParameters", d.APIInputParameters)
	populate(objectMap, "apiName", d.APIName)
	populate(objectMap, "customParameters", d.CustomParameters)
	populate(objectMap, "platformParameters", d.PlatformParameters)
	populate(objectMap, "unitsSupported", d.UnitsSupported)
	return json.Marshal(objectMap)
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info map[string]interface{} `json:"info,omitempty" azure:"ro"`

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

// MarshalJSON implements the json.Marshaller interface for type ErrorDetail.
func (e ErrorDetail) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "additionalInfo", e.AdditionalInfo)
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "message", e.Message)
	populate(objectMap, "target", e.Target)
	return json.Marshal(objectMap)
}

// ErrorResponse - Common error response for all Azure Resource Manager APIs to return error details for failed operations. (This also follows the OData
// error response format.).
// Implements the error and azcore.HTTPResponse interfaces.
type ErrorResponse struct {
	raw string
	// The error object.
	InnerError *ErrorDetail `json:"error,omitempty"`
}

// Error implements the error interface for type ErrorResponse.
// The contents of the error text are not contractual and subject to change.
func (e ErrorResponse) Error() string {
	return e.raw
}

// Extension resource.
type Extension struct {
	ProxyResource
	// Extension resource properties.
	Properties *ExtensionProperties `json:"properties,omitempty"`

	// READ-ONLY; The ETag value to implement optimistic concurrency.
	ETag *string `json:"eTag,omitempty" azure:"ro"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type Extension.
func (e Extension) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	e.ProxyResource.marshalInternal(objectMap)
	populate(objectMap, "eTag", e.ETag)
	populate(objectMap, "properties", e.Properties)
	populate(objectMap, "systemData", e.SystemData)
	return json.Marshal(objectMap)
}

// ExtensionListResponse - Paged response contains list of requested objects and a URL link to get the next set of results.
type ExtensionListResponse struct {
	// List of requested objects.
	Value []*Extension `json:"value,omitempty"`

	// READ-ONLY; Continuation link (absolute URI) to the next page of results in the list.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type ExtensionListResponse.
func (e ExtensionListResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", e.NextLink)
	populate(objectMap, "value", e.Value)
	return json.Marshal(objectMap)
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

// ExtensionsCreateOptions contains the optional parameters for the Extensions.Create method.
type ExtensionsCreateOptions struct {
	// placeholder for future optional parameters
}

// ExtensionsDeleteOptions contains the optional parameters for the Extensions.Delete method.
type ExtensionsDeleteOptions struct {
	// placeholder for future optional parameters
}

// ExtensionsGetOptions contains the optional parameters for the Extensions.Get method.
type ExtensionsGetOptions struct {
	// placeholder for future optional parameters
}

// ExtensionsListByFarmBeatsOptions contains the optional parameters for the Extensions.ListByFarmBeats method.
type ExtensionsListByFarmBeatsOptions struct {
	// Installed extension categories.
	ExtensionCategories []string
	// Installed extension ids.
	ExtensionIDs []string
	// Maximum number of items needed (inclusive).
	// Minimum = 10, Maximum = 1000, Default value = 50.
	MaxPageSize *int32
	// Skip token for getting next set of results.
	SkipToken *string
}

// ExtensionsUpdateOptions contains the optional parameters for the Extensions.Update method.
type ExtensionsUpdateOptions struct {
	// placeholder for future optional parameters
}

// FarmBeats ARM Resource.
type FarmBeats struct {
	TrackedResource
	// FarmBeats ARM Resource properties.
	Properties *FarmBeatsProperties `json:"properties,omitempty"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type FarmBeats.
func (f FarmBeats) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	f.TrackedResource.marshalInternal(objectMap)
	populate(objectMap, "properties", f.Properties)
	populate(objectMap, "systemData", f.SystemData)
	return json.Marshal(objectMap)
}

// FarmBeatsExtension - FarmBeats extension resource.
type FarmBeatsExtension struct {
	ProxyResource
	// FarmBeatsExtension properties.
	Properties *FarmBeatsExtensionProperties `json:"properties,omitempty"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type FarmBeatsExtension.
func (f FarmBeatsExtension) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	f.ProxyResource.marshalInternal(objectMap)
	populate(objectMap, "properties", f.Properties)
	populate(objectMap, "systemData", f.SystemData)
	return json.Marshal(objectMap)
}

// FarmBeatsExtensionListResponse - Paged response contains list of requested objects and a URL link to get the next set of results.
type FarmBeatsExtensionListResponse struct {
	// List of requested objects.
	Value []*FarmBeatsExtension `json:"value,omitempty"`

	// READ-ONLY; Continuation link (absolute URI) to the next page of results in the list.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type FarmBeatsExtensionListResponse.
func (f FarmBeatsExtensionListResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", f.NextLink)
	populate(objectMap, "value", f.Value)
	return json.Marshal(objectMap)
}

// FarmBeatsExtensionProperties - FarmBeatsExtension properties.
type FarmBeatsExtensionProperties struct {
	// READ-ONLY; Textual description.
	Description *string `json:"description,omitempty" azure:"ro"`

	// READ-ONLY; Detailed information which shows summary of requested data. Used in descriptive get extension metadata call. Information for weather category
	// per api included are apisSupported, customParameters,
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

// MarshalJSON implements the json.Marshaller interface for type FarmBeatsExtensionProperties.
func (f FarmBeatsExtensionProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "description", f.Description)
	populate(objectMap, "detailedInformation", f.DetailedInformation)
	populate(objectMap, "extensionApiDocsLink", f.ExtensionAPIDocsLink)
	populate(objectMap, "extensionAuthLink", f.ExtensionAuthLink)
	populate(objectMap, "extensionCategory", f.ExtensionCategory)
	populate(objectMap, "farmBeatsExtensionId", f.FarmBeatsExtensionID)
	populate(objectMap, "farmBeatsExtensionName", f.FarmBeatsExtensionName)
	populate(objectMap, "farmBeatsExtensionVersion", f.FarmBeatsExtensionVersion)
	populate(objectMap, "publisherId", f.PublisherID)
	populate(objectMap, "targetResourceType", f.TargetResourceType)
	return json.Marshal(objectMap)
}

// FarmBeatsExtensionsGetOptions contains the optional parameters for the FarmBeatsExtensions.Get method.
type FarmBeatsExtensionsGetOptions struct {
	// placeholder for future optional parameters
}

// FarmBeatsExtensionsListOptions contains the optional parameters for the FarmBeatsExtensions.List method.
type FarmBeatsExtensionsListOptions struct {
	// Extension categories.
	ExtensionCategories []string
	// FarmBeatsExtension ids.
	FarmBeatsExtensionIDs []string
	// FarmBeats extension names.
	FarmBeatsExtensionNames []string
	// Maximum number of items needed (inclusive).
	// Minimum = 10, Maximum = 1000, Default value = 50.
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

// MarshalJSON implements the json.Marshaller interface for type FarmBeatsListResponse.
func (f FarmBeatsListResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", f.NextLink)
	populate(objectMap, "value", f.Value)
	return json.Marshal(objectMap)
}

// FarmBeatsModelsCreateOrUpdateOptions contains the optional parameters for the FarmBeatsModels.CreateOrUpdate method.
type FarmBeatsModelsCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// FarmBeatsModelsDeleteOptions contains the optional parameters for the FarmBeatsModels.Delete method.
type FarmBeatsModelsDeleteOptions struct {
	// placeholder for future optional parameters
}

// FarmBeatsModelsGetOptions contains the optional parameters for the FarmBeatsModels.Get method.
type FarmBeatsModelsGetOptions struct {
	// placeholder for future optional parameters
}

// FarmBeatsModelsListByResourceGroupOptions contains the optional parameters for the FarmBeatsModels.ListByResourceGroup method.
type FarmBeatsModelsListByResourceGroupOptions struct {
	// Maximum number of items needed (inclusive).
	// Minimum = 10, Maximum = 1000, Default value = 50.
	MaxPageSize *int32
	// Continuation token for getting next set of results.
	SkipToken *string
}

// FarmBeatsModelsListBySubscriptionOptions contains the optional parameters for the FarmBeatsModels.ListBySubscription method.
type FarmBeatsModelsListBySubscriptionOptions struct {
	// Maximum number of items needed (inclusive).
	// Minimum = 10, Maximum = 1000, Default value = 50.
	MaxPageSize *int32
	// Skip token for getting next set of results.
	SkipToken *string
}

// FarmBeatsModelsUpdateOptions contains the optional parameters for the FarmBeatsModels.Update method.
type FarmBeatsModelsUpdateOptions struct {
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

// MarshalJSON implements the json.Marshaller interface for type FarmBeatsUpdateRequestModel.
func (f FarmBeatsUpdateRequestModel) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "location", f.Location)
	populate(objectMap, "tags", f.Tags)
	return json.Marshal(objectMap)
}

// LocationsCheckNameAvailabilityOptions contains the optional parameters for the Locations.CheckNameAvailability method.
type LocationsCheckNameAvailabilityOptions struct {
	// placeholder for future optional parameters
}

// Operation - Details of a REST API operation, returned from the Resource Provider Operations API
type Operation struct {
	// Localized display information for this particular operation.
	Display *OperationDisplay `json:"display,omitempty"`

	// READ-ONLY; Enum. Indicates the action type. "Internal" refers to actions that are for internal only APIs.
	ActionType *ActionType `json:"actionType,omitempty" azure:"ro"`

	// READ-ONLY; Whether the operation applies to data-plane. This is "true" for data-plane operations and "false" for ARM/control-plane operations.
	IsDataAction *bool `json:"isDataAction,omitempty" azure:"ro"`

	// READ-ONLY; The name of the operation, as per Resource-Based Access Control (RBAC). Examples: "Microsoft.Compute/virtualMachines/write", "Microsoft.Compute/virtualMachines/capture/action"
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The intended executor of the operation; as in Resource Based Access Control (RBAC) and audit logs UX. Default value is "user,system"
	Origin *Origin `json:"origin,omitempty" azure:"ro"`
}

// OperationDisplay - Localized display information for this particular operation.
type OperationDisplay struct {
	// READ-ONLY; The short, localized friendly description of the operation; suitable for tool tips and detailed views.
	Description *string `json:"description,omitempty" azure:"ro"`

	// READ-ONLY; The concise, localized friendly name for the operation; suitable for dropdowns. E.g. "Create or Update Virtual Machine", "Restart Virtual
	// Machine".
	Operation *string `json:"operation,omitempty" azure:"ro"`

	// READ-ONLY; The localized friendly form of the resource provider name, e.g. "Microsoft Monitoring Insights" or "Microsoft Compute".
	Provider *string `json:"provider,omitempty" azure:"ro"`

	// READ-ONLY; The localized friendly name of the resource type related to this operation. E.g. "Virtual Machines" or "Job Schedule Collections".
	Resource *string `json:"resource,omitempty" azure:"ro"`
}

// OperationListResult - A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to get the next set of results.
type OperationListResult struct {
	// READ-ONLY; URL to get the next set of operation list results (if there are any).
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; List of operations supported by the resource provider
	Value []*Operation `json:"value,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type OperationListResult.
func (o OperationListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", o.NextLink)
	populate(objectMap, "value", o.Value)
	return json.Marshal(objectMap)
}

// OperationsListOptions contains the optional parameters for the Operations.List method.
type OperationsListOptions struct {
	// placeholder for future optional parameters
}

// ProxyResource - The resource model definition for a Azure Resource Manager proxy resource. It will not have tags and a location
type ProxyResource struct {
	Resource
}

func (p ProxyResource) marshalInternal(objectMap map[string]interface{}) {
	p.Resource.marshalInternal(objectMap)
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

// MarshalJSON implements the json.Marshaller interface for type Resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	r.marshalInternal(objectMap)
	return json.Marshal(objectMap)
}

func (r Resource) marshalInternal(objectMap map[string]interface{}) {
	populate(objectMap, "id", r.ID)
	populate(objectMap, "name", r.Name)
	populate(objectMap, "type", r.Type)
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

// MarshalJSON implements the json.Marshaller interface for type SystemData.
func (s SystemData) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populateTimeRFC3339(objectMap, "createdAt", s.CreatedAt)
	populate(objectMap, "createdBy", s.CreatedBy)
	populate(objectMap, "createdByType", s.CreatedByType)
	populateTimeRFC3339(objectMap, "lastModifiedAt", s.LastModifiedAt)
	populate(objectMap, "lastModifiedBy", s.LastModifiedBy)
	populate(objectMap, "lastModifiedByType", s.LastModifiedByType)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type SystemData.
func (s *SystemData) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "createdAt":
			err = unpopulateTimeRFC3339(val, &s.CreatedAt)
			delete(rawMsg, key)
		case "createdBy":
			err = unpopulate(val, &s.CreatedBy)
			delete(rawMsg, key)
		case "createdByType":
			err = unpopulate(val, &s.CreatedByType)
			delete(rawMsg, key)
		case "lastModifiedAt":
			err = unpopulateTimeRFC3339(val, &s.LastModifiedAt)
			delete(rawMsg, key)
		case "lastModifiedBy":
			err = unpopulate(val, &s.LastModifiedBy)
			delete(rawMsg, key)
		case "lastModifiedByType":
			err = unpopulate(val, &s.LastModifiedByType)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags' and a 'location'
type TrackedResource struct {
	Resource
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type TrackedResource.
func (t TrackedResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	t.marshalInternal(objectMap)
	return json.Marshal(objectMap)
}

func (t TrackedResource) marshalInternal(objectMap map[string]interface{}) {
	t.Resource.marshalInternal(objectMap)
	populate(objectMap, "location", t.Location)
	populate(objectMap, "tags", t.Tags)
}

// UnitSystemsInfo - Unit systems info for the data provider.
type UnitSystemsInfo struct {
	// REQUIRED; UnitSystem key sent as part of ProviderInput.
	Key *string `json:"key,omitempty"`

	// REQUIRED; List of unit systems supported by this data provider.
	Values []*string `json:"values,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type UnitSystemsInfo.
func (u UnitSystemsInfo) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "key", u.Key)
	populate(objectMap, "values", u.Values)
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

func unpopulate(data json.RawMessage, v interface{}) error {
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, v)
}
