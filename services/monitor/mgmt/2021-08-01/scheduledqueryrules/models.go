package scheduledqueryrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"encoding/json"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// The package's fully qualified name.
const fqdn = "github.com/Azure/azure-sdk-for-go/services/monitor/mgmt/2021-08-01/scheduledqueryrules"

// Actions actions to invoke when the alert fires.
type Actions struct {
	// ActionGroups - Action Group resource Ids to invoke when the alert fires.
	ActionGroups *[]string `json:"actionGroups,omitempty"`
	// CustomProperties - The properties of an alert payload.
	CustomProperties map[string]*string `json:"customProperties"`
}

// MarshalJSON is the custom marshaler for Actions.
func (a Actions) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if a.ActionGroups != nil {
		objectMap["actionGroups"] = a.ActionGroups
	}
	if a.CustomProperties != nil {
		objectMap["customProperties"] = a.CustomProperties
	}
	return json.Marshal(objectMap)
}

// AzureEntityResource the resource model definition for an Azure Resource Manager resource with an etag.
type AzureEntityResource struct {
	// Etag - READ-ONLY; Resource Etag.
	Etag *string `json:"etag,omitempty"`
	// ID - READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty"`
}

// MarshalJSON is the custom marshaler for AzureEntityResource.
func (aer AzureEntityResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	return json.Marshal(objectMap)
}

// Condition a condition of the scheduled query rule.
type Condition struct {
	// Query - Log query alert
	Query *string `json:"query,omitempty"`
	// TimeAggregation - Aggregation type. Relevant and required only for rules of the kind LogAlert. Possible values include: 'TimeAggregationCount', 'TimeAggregationAverage', 'TimeAggregationMinimum', 'TimeAggregationMaximum', 'TimeAggregationTotal'
	TimeAggregation TimeAggregation `json:"timeAggregation,omitempty"`
	// MetricMeasureColumn - The column containing the metric measure number. Relevant only for rules of the kind LogAlert.
	MetricMeasureColumn *string `json:"metricMeasureColumn,omitempty"`
	// ResourceIDColumn - The column containing the resource id. The content of the column must be a uri formatted as resource id. Relevant only for rules of the kind LogAlert.
	ResourceIDColumn *string `json:"resourceIdColumn,omitempty"`
	// Dimensions - List of Dimensions conditions
	Dimensions *[]Dimension `json:"dimensions,omitempty"`
	// Operator - The criteria operator. Relevant and required only for rules of the kind LogAlert. Possible values include: 'ConditionOperatorEquals', 'ConditionOperatorGreaterThan', 'ConditionOperatorGreaterThanOrEqual', 'ConditionOperatorLessThan', 'ConditionOperatorLessThanOrEqual'
	Operator ConditionOperator `json:"operator,omitempty"`
	// Threshold - the criteria threshold value that activates the alert. Relevant and required only for rules of the kind LogAlert.
	Threshold *float64 `json:"threshold,omitempty"`
	// FailingPeriods - The minimum number of violations required within the selected lookback time window required to raise an alert. Relevant only for rules of the kind LogAlert.
	FailingPeriods *ConditionFailingPeriods `json:"failingPeriods,omitempty"`
	// MetricName - The name of the metric to be sent. Relevant and required only for rules of the kind LogToMetric.
	MetricName *string `json:"metricName,omitempty"`
}

// ConditionFailingPeriods the minimum number of violations required within the selected lookback time
// window required to raise an alert. Relevant only for rules of the kind LogAlert.
type ConditionFailingPeriods struct {
	// NumberOfEvaluationPeriods - The number of aggregated lookback points. The lookback time window is calculated based on the aggregation granularity (windowSize) and the selected number of aggregated points. Default value is 1
	NumberOfEvaluationPeriods *int64 `json:"numberOfEvaluationPeriods,omitempty"`
	// MinFailingPeriodsToAlert - The number of violations to trigger an alert. Should be smaller or equal to numberOfEvaluationPeriods. Default value is 1
	MinFailingPeriodsToAlert *int64 `json:"minFailingPeriodsToAlert,omitempty"`
}

// Criteria the rule criteria that defines the conditions of the scheduled query rule.
type Criteria struct {
	// AllOf - A list of conditions to evaluate against the specified scopes
	AllOf *[]Condition `json:"allOf,omitempty"`
}

// Dimension dimension splitting and filtering definition
type Dimension struct {
	// Name - Name of the dimension
	Name *string `json:"name,omitempty"`
	// Operator - Operator for dimension values. Possible values include: 'DimensionOperatorInclude', 'DimensionOperatorExclude'
	Operator DimensionOperator `json:"operator,omitempty"`
	// Values - List of dimension values
	Values *[]string `json:"values,omitempty"`
}

// ErrorAdditionalInfo the resource management error additional info.
type ErrorAdditionalInfo struct {
	// Type - READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty"`
	// Info - READ-ONLY; The additional info.
	Info interface{} `json:"info,omitempty"`
}

// MarshalJSON is the custom marshaler for ErrorAdditionalInfo.
func (eai ErrorAdditionalInfo) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	return json.Marshal(objectMap)
}

// ErrorContract describes the format of Error response.
type ErrorContract struct {
	// Error - The error details.
	Error *ErrorResponse `json:"error,omitempty"`
}

// ErrorResponse common error response for all Azure Resource Manager APIs to return error details for
// failed operations. (This also follows the OData error response format.)
type ErrorResponse struct {
	// Code - READ-ONLY; The error code.
	Code *string `json:"code,omitempty"`
	// Message - READ-ONLY; The error message.
	Message *string `json:"message,omitempty"`
	// Target - READ-ONLY; The error target.
	Target *string `json:"target,omitempty"`
	// Details - READ-ONLY; The error details.
	Details *[]ErrorResponse `json:"details,omitempty"`
	// AdditionalInfo - READ-ONLY; The error additional info.
	AdditionalInfo *[]ErrorAdditionalInfo `json:"additionalInfo,omitempty"`
}

// MarshalJSON is the custom marshaler for ErrorResponse.
func (er ErrorResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	return json.Marshal(objectMap)
}

// Properties scheduled query rule Definition
type Properties struct {
	// CreatedWithAPIVersion - READ-ONLY; The api-version used when creating this alert rule
	CreatedWithAPIVersion *string `json:"createdWithApiVersion,omitempty"`
	// IsLegacyLogAnalyticsRule - READ-ONLY; True if alert rule is legacy Log Analytic rule
	IsLegacyLogAnalyticsRule *bool `json:"isLegacyLogAnalyticsRule,omitempty"`
	// Description - The description of the scheduled query rule.
	Description *string `json:"description,omitempty"`
	// DisplayName - The display name of the alert rule
	DisplayName *string `json:"displayName,omitempty"`
	// Severity - Severity of the alert. Should be an integer between [0-4]. Value of 0 is severest. Relevant and required only for rules of the kind LogAlert.
	Severity *int64 `json:"severity,omitempty"`
	// Enabled - The flag which indicates whether this scheduled query rule is enabled. Value should be true or false
	Enabled *bool `json:"enabled,omitempty"`
	// Scopes - The list of resource id's that this scheduled query rule is scoped to.
	Scopes *[]string `json:"scopes,omitempty"`
	// EvaluationFrequency - How often the scheduled query rule is evaluated represented in ISO 8601 duration format. Relevant and required only for rules of the kind LogAlert.
	EvaluationFrequency *string `json:"evaluationFrequency,omitempty"`
	// WindowSize - The period of time (in ISO 8601 duration format) on which the Alert query will be executed (bin size). Relevant and required only for rules of the kind LogAlert.
	WindowSize *string `json:"windowSize,omitempty"`
	// OverrideQueryTimeRange - If specified then overrides the query time range (default is WindowSize*NumberOfEvaluationPeriods). Relevant only for rules of the kind LogAlert.
	OverrideQueryTimeRange *string `json:"overrideQueryTimeRange,omitempty"`
	// TargetResourceTypes - List of resource type of the target resource(s) on which the alert is created/updated. For example if the scope is a resource group and targetResourceTypes is Microsoft.Compute/virtualMachines, then a different alert will be fired for each virtual machine in the resource group which meet the alert criteria. Relevant only for rules of the kind LogAlert
	TargetResourceTypes *[]string `json:"targetResourceTypes,omitempty"`
	// Criteria - The rule criteria that defines the conditions of the scheduled query rule.
	Criteria *Criteria `json:"criteria,omitempty"`
	// MuteActionsDuration - Mute actions for the chosen period of time (in ISO 8601 duration format) after the alert is fired. Relevant only for rules of the kind LogAlert.
	MuteActionsDuration *string `json:"muteActionsDuration,omitempty"`
	// Actions - Actions to invoke when the alert fires.
	Actions *Actions `json:"actions,omitempty"`
	// IsWorkspaceAlertsStorageConfigured - READ-ONLY; The flag which indicates whether this scheduled query rule has been configured to be stored in the customer's storage. The default is false.
	IsWorkspaceAlertsStorageConfigured *bool `json:"isWorkspaceAlertsStorageConfigured,omitempty"`
	// CheckWorkspaceAlertsStorageConfigured - The flag which indicates whether this scheduled query rule should be stored in the customer's storage. The default is false. Relevant only for rules of the kind LogAlert.
	CheckWorkspaceAlertsStorageConfigured *bool `json:"checkWorkspaceAlertsStorageConfigured,omitempty"`
	// SkipQueryValidation - The flag which indicates whether the provided query should be validated or not. The default is false. Relevant only for rules of the kind LogAlert.
	SkipQueryValidation *bool `json:"skipQueryValidation,omitempty"`
	// AutoMitigate - The flag that indicates whether the alert should be automatically resolved or not. The default is true. Relevant only for rules of the kind LogAlert.
	AutoMitigate *bool `json:"autoMitigate,omitempty"`
}

// MarshalJSON is the custom marshaler for Properties.
func (p Properties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if p.Description != nil {
		objectMap["description"] = p.Description
	}
	if p.DisplayName != nil {
		objectMap["displayName"] = p.DisplayName
	}
	if p.Severity != nil {
		objectMap["severity"] = p.Severity
	}
	if p.Enabled != nil {
		objectMap["enabled"] = p.Enabled
	}
	if p.Scopes != nil {
		objectMap["scopes"] = p.Scopes
	}
	if p.EvaluationFrequency != nil {
		objectMap["evaluationFrequency"] = p.EvaluationFrequency
	}
	if p.WindowSize != nil {
		objectMap["windowSize"] = p.WindowSize
	}
	if p.OverrideQueryTimeRange != nil {
		objectMap["overrideQueryTimeRange"] = p.OverrideQueryTimeRange
	}
	if p.TargetResourceTypes != nil {
		objectMap["targetResourceTypes"] = p.TargetResourceTypes
	}
	if p.Criteria != nil {
		objectMap["criteria"] = p.Criteria
	}
	if p.MuteActionsDuration != nil {
		objectMap["muteActionsDuration"] = p.MuteActionsDuration
	}
	if p.Actions != nil {
		objectMap["actions"] = p.Actions
	}
	if p.CheckWorkspaceAlertsStorageConfigured != nil {
		objectMap["checkWorkspaceAlertsStorageConfigured"] = p.CheckWorkspaceAlertsStorageConfigured
	}
	if p.SkipQueryValidation != nil {
		objectMap["skipQueryValidation"] = p.SkipQueryValidation
	}
	if p.AutoMitigate != nil {
		objectMap["autoMitigate"] = p.AutoMitigate
	}
	return json.Marshal(objectMap)
}

// ProxyResource the resource model definition for a Azure Resource Manager proxy resource. It will not
// have tags and a location
type ProxyResource struct {
	// ID - READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty"`
}

// MarshalJSON is the custom marshaler for ProxyResource.
func (pr ProxyResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	return json.Marshal(objectMap)
}

// Resource common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// ID - READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty"`
}

// MarshalJSON is the custom marshaler for Resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	return json.Marshal(objectMap)
}

// ResourceCollection represents a collection of scheduled query rule resources.
type ResourceCollection struct {
	autorest.Response `json:"-"`
	// Value - The values for the scheduled query rule resources.
	Value *[]ResourceType `json:"value,omitempty"`
	// NextLink - READ-ONLY; Provides the link to retrieve the next set of elements.
	NextLink *string `json:"nextLink,omitempty"`
}

// MarshalJSON is the custom marshaler for ResourceCollection.
func (rc ResourceCollection) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if rc.Value != nil {
		objectMap["value"] = rc.Value
	}
	return json.Marshal(objectMap)
}

// ResourceCollectionIterator provides access to a complete listing of ResourceType values.
type ResourceCollectionIterator struct {
	i    int
	page ResourceCollectionPage
}

// NextWithContext advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
func (iter *ResourceCollectionIterator) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ResourceCollectionIterator.NextWithContext")
		defer func() {
			sc := -1
			if iter.Response().Response.Response != nil {
				sc = iter.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	iter.i++
	if iter.i < len(iter.page.Values()) {
		return nil
	}
	err = iter.page.NextWithContext(ctx)
	if err != nil {
		iter.i--
		return err
	}
	iter.i = 0
	return nil
}

// Next advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (iter *ResourceCollectionIterator) Next() error {
	return iter.NextWithContext(context.Background())
}

// NotDone returns true if the enumeration should be started or is not yet complete.
func (iter ResourceCollectionIterator) NotDone() bool {
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}

// Response returns the raw server response from the last page request.
func (iter ResourceCollectionIterator) Response() ResourceCollection {
	return iter.page.Response()
}

// Value returns the current value or a zero-initialized value if the
// iterator has advanced beyond the end of the collection.
func (iter ResourceCollectionIterator) Value() ResourceType {
	if !iter.page.NotDone() {
		return ResourceType{}
	}
	return iter.page.Values()[iter.i]
}

// Creates a new instance of the ResourceCollectionIterator type.
func NewResourceCollectionIterator(page ResourceCollectionPage) ResourceCollectionIterator {
	return ResourceCollectionIterator{page: page}
}

// IsEmpty returns true if the ListResult contains no values.
func (rc ResourceCollection) IsEmpty() bool {
	return rc.Value == nil || len(*rc.Value) == 0
}

// hasNextLink returns true if the NextLink is not empty.
func (rc ResourceCollection) hasNextLink() bool {
	return rc.NextLink != nil && len(*rc.NextLink) != 0
}

// resourceCollectionPreparer prepares a request to retrieve the next set of results.
// It returns nil if no more results exist.
func (rc ResourceCollection) resourceCollectionPreparer(ctx context.Context) (*http.Request, error) {
	if !rc.hasNextLink() {
		return nil, nil
	}
	return autorest.Prepare((&http.Request{}).WithContext(ctx),
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(rc.NextLink)))
}

// ResourceCollectionPage contains a page of ResourceType values.
type ResourceCollectionPage struct {
	fn func(context.Context, ResourceCollection) (ResourceCollection, error)
	rc ResourceCollection
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *ResourceCollectionPage) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ResourceCollectionPage.NextWithContext")
		defer func() {
			sc := -1
			if page.Response().Response.Response != nil {
				sc = page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	for {
		next, err := page.fn(ctx, page.rc)
		if err != nil {
			return err
		}
		page.rc = next
		if !next.hasNextLink() || !next.IsEmpty() {
			break
		}
	}
	return nil
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (page *ResourceCollectionPage) Next() error {
	return page.NextWithContext(context.Background())
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page ResourceCollectionPage) NotDone() bool {
	return !page.rc.IsEmpty()
}

// Response returns the raw server response from the last page request.
func (page ResourceCollectionPage) Response() ResourceCollection {
	return page.rc
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page ResourceCollectionPage) Values() []ResourceType {
	if page.rc.IsEmpty() {
		return nil
	}
	return *page.rc.Value
}

// Creates a new instance of the ResourceCollectionPage type.
func NewResourceCollectionPage(cur ResourceCollection, getNextPage func(context.Context, ResourceCollection) (ResourceCollection, error)) ResourceCollectionPage {
	return ResourceCollectionPage{
		fn: getNextPage,
		rc: cur,
	}
}

// ResourcePatch the scheduled query rule resource for patch operations.
type ResourcePatch struct {
	// Tags - Resource tags
	Tags map[string]*string `json:"tags"`
	// Properties - The scheduled query rule properties of the resource.
	*Properties `json:"properties,omitempty"`
}

// MarshalJSON is the custom marshaler for ResourcePatch.
func (rp ResourcePatch) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if rp.Tags != nil {
		objectMap["tags"] = rp.Tags
	}
	if rp.Properties != nil {
		objectMap["properties"] = rp.Properties
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for ResourcePatch struct.
func (rp *ResourcePatch) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "tags":
			if v != nil {
				var tags map[string]*string
				err = json.Unmarshal(*v, &tags)
				if err != nil {
					return err
				}
				rp.Tags = tags
			}
		case "properties":
			if v != nil {
				var properties Properties
				err = json.Unmarshal(*v, &properties)
				if err != nil {
					return err
				}
				rp.Properties = &properties
			}
		}
	}

	return nil
}

// ResourceType the scheduled query rule resource.
type ResourceType struct {
	autorest.Response `json:"-"`
	// Kind - Indicates the type of scheduled query rule. The default is LogAlert. Possible values include: 'KindLogAlert', 'KindLogToMetric'
	Kind Kind `json:"kind,omitempty"`
	// Etag - READ-ONLY; The etag field is *not* required. If it is provided in the response body, it must also be provided as a header per the normal etag convention.  Entity tags are used for comparing two or more entities from the same requested resource. HTTP/1.1 uses entity tags in the etag (section 14.19), If-Match (section 14.24), If-None-Match (section 14.26), and If-Range (section 14.27) header fields.
	Etag *string `json:"etag,omitempty"`
	// SystemData - READ-ONLY; SystemData of ScheduledQueryRule.
	SystemData *SystemData `json:"systemData,omitempty"`
	// Properties - The rule properties of the resource.
	*Properties `json:"properties,omitempty"`
	// Tags - Resource tags.
	Tags map[string]*string `json:"tags"`
	// Location - The geo-location where the resource lives
	Location *string `json:"location,omitempty"`
	// ID - READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty"`
}

// MarshalJSON is the custom marshaler for ResourceType.
func (rt ResourceType) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if rt.Kind != "" {
		objectMap["kind"] = rt.Kind
	}
	if rt.Properties != nil {
		objectMap["properties"] = rt.Properties
	}
	if rt.Tags != nil {
		objectMap["tags"] = rt.Tags
	}
	if rt.Location != nil {
		objectMap["location"] = rt.Location
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for ResourceType struct.
func (rt *ResourceType) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "kind":
			if v != nil {
				var kind Kind
				err = json.Unmarshal(*v, &kind)
				if err != nil {
					return err
				}
				rt.Kind = kind
			}
		case "etag":
			if v != nil {
				var etag string
				err = json.Unmarshal(*v, &etag)
				if err != nil {
					return err
				}
				rt.Etag = &etag
			}
		case "systemData":
			if v != nil {
				var systemData SystemData
				err = json.Unmarshal(*v, &systemData)
				if err != nil {
					return err
				}
				rt.SystemData = &systemData
			}
		case "properties":
			if v != nil {
				var properties Properties
				err = json.Unmarshal(*v, &properties)
				if err != nil {
					return err
				}
				rt.Properties = &properties
			}
		case "tags":
			if v != nil {
				var tags map[string]*string
				err = json.Unmarshal(*v, &tags)
				if err != nil {
					return err
				}
				rt.Tags = tags
			}
		case "location":
			if v != nil {
				var location string
				err = json.Unmarshal(*v, &location)
				if err != nil {
					return err
				}
				rt.Location = &location
			}
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				rt.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				rt.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				rt.Type = &typeVar
			}
		}
	}

	return nil
}

// SystemData metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// CreatedBy - The identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`
	// CreatedByType - The type of identity that created the resource. Possible values include: 'CreatedByTypeUser', 'CreatedByTypeApplication', 'CreatedByTypeManagedIdentity', 'CreatedByTypeKey'
	CreatedByType CreatedByType `json:"createdByType,omitempty"`
	// CreatedAt - The timestamp of resource creation (UTC).
	CreatedAt *date.Time `json:"createdAt,omitempty"`
	// LastModifiedBy - The identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`
	// LastModifiedByType - The type of identity that last modified the resource. Possible values include: 'CreatedByTypeUser', 'CreatedByTypeApplication', 'CreatedByTypeManagedIdentity', 'CreatedByTypeKey'
	LastModifiedByType CreatedByType `json:"lastModifiedByType,omitempty"`
	// LastModifiedAt - The timestamp of resource last modification (UTC)
	LastModifiedAt *date.Time `json:"lastModifiedAt,omitempty"`
}

// TrackedResource the resource model definition for an Azure Resource Manager tracked top level resource
// which has 'tags' and a 'location'
type TrackedResource struct {
	// Tags - Resource tags.
	Tags map[string]*string `json:"tags"`
	// Location - The geo-location where the resource lives
	Location *string `json:"location,omitempty"`
	// ID - READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty"`
}

// MarshalJSON is the custom marshaler for TrackedResource.
func (tr TrackedResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if tr.Tags != nil {
		objectMap["tags"] = tr.Tags
	}
	if tr.Location != nil {
		objectMap["location"] = tr.Location
	}
	return json.Marshal(objectMap)
}
