package iotcentral

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"encoding/json"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"net/http"
)

// AppSku enumerates the values for app sku.
type AppSku string

const (
	// F1 ...
	F1 AppSku = "F1"
	// S1 ...
	S1 AppSku = "S1"
)

// PossibleAppSkuValues returns an array of possible values for the AppSku const type.
func PossibleAppSkuValues() []AppSku {
	return []AppSku{F1, S1}
}

// App the IoT Central application.
type App struct {
	autorest.Response `json:"-"`
	// AppProperties - The common properties of an IoT Central application.
	*AppProperties `json:"properties,omitempty"`
	// Sku - A valid instance SKU.
	Sku *AppSkuInfo `json:"sku,omitempty"`
	// ID - The ARM resource identifier.
	ID *string `json:"id,omitempty"`
	// Name - The ARM resource name.
	Name *string `json:"name,omitempty"`
	// Type - The resource type.
	Type *string `json:"type,omitempty"`
	// Location - The resource location.
	Location *string `json:"location,omitempty"`
	// Tags - The resource tags.
	Tags map[string]*string `json:"tags"`
}

// MarshalJSON is the custom marshaler for App.
func (a App) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if a.AppProperties != nil {
		objectMap["properties"] = a.AppProperties
	}
	if a.Sku != nil {
		objectMap["sku"] = a.Sku
	}
	if a.ID != nil {
		objectMap["id"] = a.ID
	}
	if a.Name != nil {
		objectMap["name"] = a.Name
	}
	if a.Type != nil {
		objectMap["type"] = a.Type
	}
	if a.Location != nil {
		objectMap["location"] = a.Location
	}
	if a.Tags != nil {
		objectMap["tags"] = a.Tags
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for App struct.
func (a *App) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "properties":
			if v != nil {
				var appProperties AppProperties
				err = json.Unmarshal(*v, &appProperties)
				if err != nil {
					return err
				}
				a.AppProperties = &appProperties
			}
		case "sku":
			if v != nil {
				var sku AppSkuInfo
				err = json.Unmarshal(*v, &sku)
				if err != nil {
					return err
				}
				a.Sku = &sku
			}
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				a.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				a.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				a.Type = &typeVar
			}
		case "location":
			if v != nil {
				var location string
				err = json.Unmarshal(*v, &location)
				if err != nil {
					return err
				}
				a.Location = &location
			}
		case "tags":
			if v != nil {
				var tags map[string]*string
				err = json.Unmarshal(*v, &tags)
				if err != nil {
					return err
				}
				a.Tags = tags
			}
		}
	}

	return nil
}

// AppAvailabilityInfo the properties indicating whether a given IoT Central application information is available.
type AppAvailabilityInfo struct {
	autorest.Response `json:"-"`
	// NameAvailable - The value which indicates whether the provided name is available.
	NameAvailable *bool `json:"nameAvailable,omitempty"`
	// Subdomain - The value which indicates whether the provided subdomain is available.
	Subdomain *string `json:"subdomain,omitempty"`
	// Reason - The reason for unavailability.
	Reason *string `json:"reason,omitempty"`
	// Message - The detailed reason message.
	Message *string `json:"message,omitempty"`
}

// AppListResult a list of IoT Central Applications with a next link.
type AppListResult struct {
	autorest.Response `json:"-"`
	// NextLink - The link used to get the next page of IoT Central Applications.
	NextLink *string `json:"nextLink,omitempty"`
	// Value - A list of IoT Central Applications.
	Value *[]App `json:"value,omitempty"`
}

// AppListResultIterator provides access to a complete listing of App values.
type AppListResultIterator struct {
	i    int
	page AppListResultPage
}

// Next advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
func (iter *AppListResultIterator) Next() error {
	iter.i++
	if iter.i < len(iter.page.Values()) {
		return nil
	}
	err := iter.page.Next()
	if err != nil {
		iter.i--
		return err
	}
	iter.i = 0
	return nil
}

// NotDone returns true if the enumeration should be started or is not yet complete.
func (iter AppListResultIterator) NotDone() bool {
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}

// Response returns the raw server response from the last page request.
func (iter AppListResultIterator) Response() AppListResult {
	return iter.page.Response()
}

// Value returns the current value or a zero-initialized value if the
// iterator has advanced beyond the end of the collection.
func (iter AppListResultIterator) Value() App {
	if !iter.page.NotDone() {
		return App{}
	}
	return iter.page.Values()[iter.i]
}

// IsEmpty returns true if the ListResult contains no values.
func (alr AppListResult) IsEmpty() bool {
	return alr.Value == nil || len(*alr.Value) == 0
}

// appListResultPreparer prepares a request to retrieve the next set of results.
// It returns nil if no more results exist.
func (alr AppListResult) appListResultPreparer() (*http.Request, error) {
	if alr.NextLink == nil || len(to.String(alr.NextLink)) < 1 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(alr.NextLink)))
}

// AppListResultPage contains a page of App values.
type AppListResultPage struct {
	fn  func(AppListResult) (AppListResult, error)
	alr AppListResult
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *AppListResultPage) Next() error {
	next, err := page.fn(page.alr)
	if err != nil {
		return err
	}
	page.alr = next
	return nil
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page AppListResultPage) NotDone() bool {
	return !page.alr.IsEmpty()
}

// Response returns the raw server response from the last page request.
func (page AppListResultPage) Response() AppListResult {
	return page.alr
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page AppListResultPage) Values() []App {
	if page.alr.IsEmpty() {
		return nil
	}
	return *page.alr.Value
}

// AppPatch the description of the IoT Central application.
type AppPatch struct {
	// Tags - Instance tags
	Tags map[string]*string `json:"tags"`
	// AppProperties - The common properties of an IoT Central application.
	*AppProperties `json:"properties,omitempty"`
}

// MarshalJSON is the custom marshaler for AppPatch.
func (ap AppPatch) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if ap.Tags != nil {
		objectMap["tags"] = ap.Tags
	}
	if ap.AppProperties != nil {
		objectMap["properties"] = ap.AppProperties
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for AppPatch struct.
func (ap *AppPatch) UnmarshalJSON(body []byte) error {
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
				ap.Tags = tags
			}
		case "properties":
			if v != nil {
				var appProperties AppProperties
				err = json.Unmarshal(*v, &appProperties)
				if err != nil {
					return err
				}
				ap.AppProperties = &appProperties
			}
		}
	}

	return nil
}

// AppProperties the properties of an IoT Central application.
type AppProperties struct {
	// ApplicationID - The ID of the application.
	ApplicationID *string `json:"applicationId,omitempty"`
	// DisplayName - The display name of the application.
	DisplayName *string `json:"displayName,omitempty"`
	// Subdomain - The subdomain of the application.
	Subdomain *string `json:"subdomain,omitempty"`
	// Template - The ID of the application template, which is a blueprint that defines the characteristics and behaviors of an application. Optional; if not specified, defaults to a blank blueprint and allows the application to be defined from scratch.
	Template *string `json:"template,omitempty"`
}

// AppsCreateOrUpdateFuture an abstraction for monitoring and retrieving the results of a long-running operation.
type AppsCreateOrUpdateFuture struct {
	azure.Future
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future *AppsCreateOrUpdateFuture) Result(client AppsClient) (a App, err error) {
	var done bool
	done, err = future.Done(client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotcentral.AppsCreateOrUpdateFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("iotcentral.AppsCreateOrUpdateFuture")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if a.Response.Response, err = future.GetResult(sender); err == nil && a.Response.Response.StatusCode != http.StatusNoContent {
		a, err = client.CreateOrUpdateResponder(a.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "iotcentral.AppsCreateOrUpdateFuture", "Result", a.Response.Response, "Failure responding to request")
		}
	}
	return
}

// AppsDeleteFuture an abstraction for monitoring and retrieving the results of a long-running operation.
type AppsDeleteFuture struct {
	azure.Future
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future *AppsDeleteFuture) Result(client AppsClient) (ar autorest.Response, err error) {
	var done bool
	done, err = future.Done(client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotcentral.AppsDeleteFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("iotcentral.AppsDeleteFuture")
		return
	}
	ar.Response = future.Response()
	return
}

// AppSkuInfo information about the SKU of the IoT Central application.
type AppSkuInfo struct {
	// Name - The name of the SKU. Possible values include: 'F1', 'S1'
	Name AppSku `json:"name,omitempty"`
}

// AppsUpdateFuture an abstraction for monitoring and retrieving the results of a long-running operation.
type AppsUpdateFuture struct {
	azure.Future
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future *AppsUpdateFuture) Result(client AppsClient) (a App, err error) {
	var done bool
	done, err = future.Done(client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotcentral.AppsUpdateFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("iotcentral.AppsUpdateFuture")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if a.Response.Response, err = future.GetResult(sender); err == nil && a.Response.Response.StatusCode != http.StatusNoContent {
		a, err = client.UpdateResponder(a.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "iotcentral.AppsUpdateFuture", "Result", a.Response.Response, "Failure responding to request")
		}
	}
	return
}

// ErrorDetails error details.
type ErrorDetails struct {
	Error *ErrorResponseBody `json:"error,omitempty"`
}

// ErrorResponseBody details of error response.
type ErrorResponseBody struct {
	// Code - Error code, intended to be consumed programmatically.
	Code *string `json:"code,omitempty"`
	// Message - Description of the error, intended for display in user interface.
	Message *string `json:"message,omitempty"`
	// Target - Target of the particular error, for example name of the property.
	Target *string `json:"target,omitempty"`
	// Details - A list of additional details about the error.
	Details *[]ErrorResponseBody `json:"details,omitempty"`
}

// Operation ioT Central REST API operation
type Operation struct {
	// Name - Operation name: {provider}/{resource}/{read | write | action | delete}
	Name *string `json:"name,omitempty"`
	// Display - The object that represents the operation.
	Display *OperationDisplay `json:"display,omitempty"`
}

// OperationDisplay the object that represents the operation.
type OperationDisplay struct {
	// Provider - Service provider: Microsoft IoT Central
	Provider *string `json:"provider,omitempty"`
	// Resource - Resource Type: IoT Central
	Resource *string `json:"resource,omitempty"`
	// Operation - Name of the operation
	Operation *string `json:"operation,omitempty"`
	// Description - Friendly description for the operation,
	Description *string `json:"description,omitempty"`
}

// OperationInputs input values.
type OperationInputs struct {
	// Name - The name of the IoT Central application instance to check.
	Name *string `json:"name,omitempty"`
	// Subdomain - The subdomain of the IoT Central application instance to check.
	Subdomain *string `json:"subdomain,omitempty"`
	// Type - The name of the IoT Central resource name to query.
	Type *string `json:"type,omitempty"`
}

// OperationListResult a list of IoT Central operations. It contains a list of operations and a URL link to get the
// next set of results.
type OperationListResult struct {
	autorest.Response `json:"-"`
	// NextLink - The link used to get the next page of IoT Central description objects.
	NextLink *string `json:"nextLink,omitempty"`
	// Value - A list of operations supported by the Microsoft.IoTCentral resource provider.
	Value *[]Operation `json:"value,omitempty"`
}

// OperationListResultIterator provides access to a complete listing of Operation values.
type OperationListResultIterator struct {
	i    int
	page OperationListResultPage
}

// Next advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
func (iter *OperationListResultIterator) Next() error {
	iter.i++
	if iter.i < len(iter.page.Values()) {
		return nil
	}
	err := iter.page.Next()
	if err != nil {
		iter.i--
		return err
	}
	iter.i = 0
	return nil
}

// NotDone returns true if the enumeration should be started or is not yet complete.
func (iter OperationListResultIterator) NotDone() bool {
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}

// Response returns the raw server response from the last page request.
func (iter OperationListResultIterator) Response() OperationListResult {
	return iter.page.Response()
}

// Value returns the current value or a zero-initialized value if the
// iterator has advanced beyond the end of the collection.
func (iter OperationListResultIterator) Value() Operation {
	if !iter.page.NotDone() {
		return Operation{}
	}
	return iter.page.Values()[iter.i]
}

// IsEmpty returns true if the ListResult contains no values.
func (olr OperationListResult) IsEmpty() bool {
	return olr.Value == nil || len(*olr.Value) == 0
}

// operationListResultPreparer prepares a request to retrieve the next set of results.
// It returns nil if no more results exist.
func (olr OperationListResult) operationListResultPreparer() (*http.Request, error) {
	if olr.NextLink == nil || len(to.String(olr.NextLink)) < 1 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(olr.NextLink)))
}

// OperationListResultPage contains a page of Operation values.
type OperationListResultPage struct {
	fn  func(OperationListResult) (OperationListResult, error)
	olr OperationListResult
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *OperationListResultPage) Next() error {
	next, err := page.fn(page.olr)
	if err != nil {
		return err
	}
	page.olr = next
	return nil
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page OperationListResultPage) NotDone() bool {
	return !page.olr.IsEmpty()
}

// Response returns the raw server response from the last page request.
func (page OperationListResultPage) Response() OperationListResult {
	return page.olr
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page OperationListResultPage) Values() []Operation {
	if page.olr.IsEmpty() {
		return nil
	}
	return *page.olr.Value
}

// Resource the common properties of an ARM resource.
type Resource struct {
	// ID - The ARM resource identifier.
	ID *string `json:"id,omitempty"`
	// Name - The ARM resource name.
	Name *string `json:"name,omitempty"`
	// Type - The resource type.
	Type *string `json:"type,omitempty"`
	// Location - The resource location.
	Location *string `json:"location,omitempty"`
	// Tags - The resource tags.
	Tags map[string]*string `json:"tags"`
}

// MarshalJSON is the custom marshaler for Resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if r.ID != nil {
		objectMap["id"] = r.ID
	}
	if r.Name != nil {
		objectMap["name"] = r.Name
	}
	if r.Type != nil {
		objectMap["type"] = r.Type
	}
	if r.Location != nil {
		objectMap["location"] = r.Location
	}
	if r.Tags != nil {
		objectMap["tags"] = r.Tags
	}
	return json.Marshal(objectMap)
}
