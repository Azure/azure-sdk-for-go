package delegatednetwork

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
	"context"
	"encoding/json"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// The package's fully qualified name.
const fqdn = "github.com/Azure/azure-sdk-for-go/services/preview/delegatednetwork/mgmt/2020-08-08-preview/delegatednetwork"

// ControllerCreateFuture an abstraction for monitoring and retrieving the results of a long-running operation.
type ControllerCreateFuture struct {
	azure.Future
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future *ControllerCreateFuture) Result(client ControllerClient) (dc DelegatedController, err error) {
	var done bool
	done, err = future.DoneWithContext(context.Background(), client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "delegatednetwork.ControllerCreateFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("delegatednetwork.ControllerCreateFuture")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if dc.Response.Response, err = future.GetResult(sender); err == nil && dc.Response.Response.StatusCode != http.StatusNoContent {
		dc, err = client.CreateResponder(dc.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "delegatednetwork.ControllerCreateFuture", "Result", dc.Response.Response, "Failure responding to request")
		}
	}
	return
}

// ControllerDeleteFuture an abstraction for monitoring and retrieving the results of a long-running operation.
type ControllerDeleteFuture struct {
	azure.Future
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future *ControllerDeleteFuture) Result(client ControllerClient) (ar autorest.Response, err error) {
	var done bool
	done, err = future.DoneWithContext(context.Background(), client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "delegatednetwork.ControllerDeleteFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("delegatednetwork.ControllerDeleteFuture")
		return
	}
	ar.Response = future.Response()
	return
}

// ControllerInstanceProperties properties of orchestrator
type ControllerInstanceProperties struct {
	// ServerAppID - AAD ID used with apiserver
	ServerAppID *string `json:"serverAppID,omitempty"`
	// ServerTenantID - TenantID of server App ID
	ServerTenantID *string `json:"serverTenantID,omitempty"`
	// ClusterRootCA - RootCA certificate of kubernetes cluster
	ClusterRootCA *string `json:"clusterRootCA,omitempty"`
	// APIServerEndpoint - APIServer url
	APIServerEndpoint *string `json:"apiServerEndpoint,omitempty"`
}

// ControllerResource represents an instance of an DNC controller resource.
type ControllerResource struct {
	// ID - READ-ONLY; An identifier that represents the DNC controller resource.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the DNC controller resource.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The type of the DNC controller  resource.(Microsoft.DelegatedNetwork/controller)
	Type *string `json:"type,omitempty"`
	// Location - Location of the DNC controller resource.
	Location *string `json:"location,omitempty"`
}

// MarshalJSON is the custom marshaler for ControllerResource.
func (cr ControllerResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if cr.Location != nil {
		objectMap["location"] = cr.Location
	}
	return json.Marshal(objectMap)
}

// ControllerResponseProperties properties of Delegated controller resource.
type ControllerResponseProperties struct {
	// State - READ-ONLY; The current state of dnc controller resource. Possible values include: 'Deleting', 'Succeeded', 'Failed', 'Provisioning'
	State ControllerState `json:"state,omitempty"`
	// DncAppID - The current state of dnc controller resource.
	DncAppID *string `json:"dncAppId,omitempty"`
	// DncEndpoint - dnc endpoint url that customers can use to connect to
	DncEndpoint *string `json:"dncEndpoint,omitempty"`
}

// MarshalJSON is the custom marshaler for ControllerResponseProperties.
func (crp ControllerResponseProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if crp.DncAppID != nil {
		objectMap["dncAppId"] = crp.DncAppID
	}
	if crp.DncEndpoint != nil {
		objectMap["dncEndpoint"] = crp.DncEndpoint
	}
	return json.Marshal(objectMap)
}

// ControllerTypeParameters details of controller type.
type ControllerTypeParameters struct {
	// ControllerType - Type of controller. Possible values include: 'Kubernetes'
	ControllerType ControllerType `json:"controllerType,omitempty"`
	// ControllerInstanceProperties - Controller properties
	*ControllerInstanceProperties `json:"properties,omitempty"`
}

// MarshalJSON is the custom marshaler for ControllerTypeParameters.
func (ctp ControllerTypeParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if ctp.ControllerType != "" {
		objectMap["controllerType"] = ctp.ControllerType
	}
	if ctp.ControllerInstanceProperties != nil {
		objectMap["properties"] = ctp.ControllerInstanceProperties
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for ControllerTypeParameters struct.
func (ctp *ControllerTypeParameters) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "controllerType":
			if v != nil {
				var controllerType ControllerType
				err = json.Unmarshal(*v, &controllerType)
				if err != nil {
					return err
				}
				ctp.ControllerType = controllerType
			}
		case "properties":
			if v != nil {
				var controllerInstanceProperties ControllerInstanceProperties
				err = json.Unmarshal(*v, &controllerInstanceProperties)
				if err != nil {
					return err
				}
				ctp.ControllerInstanceProperties = &controllerInstanceProperties
			}
		}
	}

	return nil
}

// DelegatedController represents an instance of a DNC controller.
type DelegatedController struct {
	autorest.Response `json:"-"`
	// ControllerResponseProperties - Properties of the provision operation request.
	*ControllerResponseProperties `json:"properties,omitempty"`
	// ID - READ-ONLY; An identifier that represents the DNC controller resource.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the DNC controller resource.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The type of the DNC controller  resource.(Microsoft.DelegatedNetwork/controller)
	Type *string `json:"type,omitempty"`
	// Location - Location of the DNC controller resource.
	Location *string `json:"location,omitempty"`
}

// MarshalJSON is the custom marshaler for DelegatedController.
func (dc DelegatedController) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if dc.ControllerResponseProperties != nil {
		objectMap["properties"] = dc.ControllerResponseProperties
	}
	if dc.Location != nil {
		objectMap["location"] = dc.Location
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for DelegatedController struct.
func (dc *DelegatedController) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "properties":
			if v != nil {
				var controllerResponseProperties ControllerResponseProperties
				err = json.Unmarshal(*v, &controllerResponseProperties)
				if err != nil {
					return err
				}
				dc.ControllerResponseProperties = &controllerResponseProperties
			}
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				dc.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				dc.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				dc.Type = &typeVar
			}
		case "location":
			if v != nil {
				var location string
				err = json.Unmarshal(*v, &location)
				if err != nil {
					return err
				}
				dc.Location = &location
			}
		}
	}

	return nil
}

// DelegatedControllers an array of Delegated controller resources.
type DelegatedControllers struct {
	autorest.Response `json:"-"`
	// Value - An array of Delegated controller resources.
	Value *[]DelegatedController `json:"value,omitempty"`
}

// DelegatedSubnet delegated subnet details
type DelegatedSubnet struct {
	autorest.Response `json:"-"`
	// DelegatedSubnetResponseProperties - Properties of the delegated subnet request
	*DelegatedSubnetResponseProperties `json:"properties,omitempty"`
	// ID - READ-ONLY; An identifier that represents the DelegatedSubnet resource.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the DelegatedSubnet resource.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The type of the DelegatedSubnet  resource.(Microsoft.DelegatedNetwork/delegatedSubnet)
	Type *string `json:"type,omitempty"`
	// Location - Location of the DelegatedSubnet resource.
	Location *string `json:"location,omitempty"`
}

// MarshalJSON is the custom marshaler for DelegatedSubnet.
func (ds DelegatedSubnet) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if ds.DelegatedSubnetResponseProperties != nil {
		objectMap["properties"] = ds.DelegatedSubnetResponseProperties
	}
	if ds.Location != nil {
		objectMap["location"] = ds.Location
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for DelegatedSubnet struct.
func (ds *DelegatedSubnet) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "properties":
			if v != nil {
				var delegatedSubnetResponseProperties DelegatedSubnetResponseProperties
				err = json.Unmarshal(*v, &delegatedSubnetResponseProperties)
				if err != nil {
					return err
				}
				ds.DelegatedSubnetResponseProperties = &delegatedSubnetResponseProperties
			}
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				ds.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				ds.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				ds.Type = &typeVar
			}
		case "location":
			if v != nil {
				var location string
				err = json.Unmarshal(*v, &location)
				if err != nil {
					return err
				}
				ds.Location = &location
			}
		}
	}

	return nil
}

// DelegatedSubnetList an array of Delegated subnets resources.
type DelegatedSubnetList struct {
	autorest.Response `json:"-"`
	// Value - An array of Delegated subnets resources.
	Value *[]DelegatedSubnet `json:"value,omitempty"`
}

// DelegatedSubnetParameters delegatedSubnet Parameters
type DelegatedSubnetParameters struct {
	// ControllerID - Delegated Network Controller ARM resource ID
	ControllerID *string `json:"controllerID,omitempty"`
}

// DelegatedSubnetResource represents an instance of a DelegatedSubnet resource.
type DelegatedSubnetResource struct {
	// ID - READ-ONLY; An identifier that represents the DelegatedSubnet resource.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the DelegatedSubnet resource.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The type of the DelegatedSubnet  resource.(Microsoft.DelegatedNetwork/delegatedSubnet)
	Type *string `json:"type,omitempty"`
	// Location - Location of the DelegatedSubnet resource.
	Location *string `json:"location,omitempty"`
}

// MarshalJSON is the custom marshaler for DelegatedSubnetResource.
func (dsr DelegatedSubnetResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if dsr.Location != nil {
		objectMap["location"] = dsr.Location
	}
	return json.Marshal(objectMap)
}

// DelegatedSubnetResponseProperties properties of delegated subnet resource.
type DelegatedSubnetResponseProperties struct {
	// State - READ-ONLY; The current state of delegated subnet resource. Possible values include: 'DelegatedSubnetStateDeleting', 'DelegatedSubnetStateSucceeded', 'DelegatedSubnetStateFailed', 'DelegatedSubnetStateProvisioning'
	State DelegatedSubnetState `json:"state,omitempty"`
	// ResourceGUID - Guid for the resource(delegatedSubnet) created
	ResourceGUID *string `json:"resourceGuid,omitempty"`
}

// MarshalJSON is the custom marshaler for DelegatedSubnetResponseProperties.
func (dsrp DelegatedSubnetResponseProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if dsrp.ResourceGUID != nil {
		objectMap["resourceGuid"] = dsrp.ResourceGUID
	}
	return json.Marshal(objectMap)
}

// DelegatedSubnetServiceDeleteDetailsFuture an abstraction for monitoring and retrieving the results of a
// long-running operation.
type DelegatedSubnetServiceDeleteDetailsFuture struct {
	azure.Future
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future *DelegatedSubnetServiceDeleteDetailsFuture) Result(client DelegatedSubnetServiceClient) (ar autorest.Response, err error) {
	var done bool
	done, err = future.DoneWithContext(context.Background(), client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "delegatednetwork.DelegatedSubnetServiceDeleteDetailsFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("delegatednetwork.DelegatedSubnetServiceDeleteDetailsFuture")
		return
	}
	ar.Response = future.Response()
	return
}

// DelegatedSubnetServicePatchDetailsFuture an abstraction for monitoring and retrieving the results of a
// long-running operation.
type DelegatedSubnetServicePatchDetailsFuture struct {
	azure.Future
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future *DelegatedSubnetServicePatchDetailsFuture) Result(client DelegatedSubnetServiceClient) (ds DelegatedSubnet, err error) {
	var done bool
	done, err = future.DoneWithContext(context.Background(), client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "delegatednetwork.DelegatedSubnetServicePatchDetailsFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("delegatednetwork.DelegatedSubnetServicePatchDetailsFuture")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if ds.Response.Response, err = future.GetResult(sender); err == nil && ds.Response.Response.StatusCode != http.StatusNoContent {
		ds, err = client.PatchDetailsResponder(ds.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "delegatednetwork.DelegatedSubnetServicePatchDetailsFuture", "Result", ds.Response.Response, "Failure responding to request")
		}
	}
	return
}

// DelegatedSubnetServicePutDetailsFuture an abstraction for monitoring and retrieving the results of a
// long-running operation.
type DelegatedSubnetServicePutDetailsFuture struct {
	azure.Future
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future *DelegatedSubnetServicePutDetailsFuture) Result(client DelegatedSubnetServiceClient) (ds DelegatedSubnet, err error) {
	var done bool
	done, err = future.DoneWithContext(context.Background(), client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "delegatednetwork.DelegatedSubnetServicePutDetailsFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("delegatednetwork.DelegatedSubnetServicePutDetailsFuture")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if ds.Response.Response, err = future.GetResult(sender); err == nil && ds.Response.Response.StatusCode != http.StatusNoContent {
		ds, err = client.PutDetailsResponder(ds.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "delegatednetwork.DelegatedSubnetServicePutDetailsFuture", "Result", ds.Response.Response, "Failure responding to request")
		}
	}
	return
}

// ErrorDefinition error definition.
type ErrorDefinition struct {
	// Code - READ-ONLY; Service specific error code which serves as the substatus for the HTTP error code.
	Code *string `json:"code,omitempty"`
	// Message - READ-ONLY; Description of the error.
	Message *string `json:"message,omitempty"`
	// Details - READ-ONLY; Internal error details.
	Details *[]ErrorDefinition `json:"details,omitempty"`
}

// ErrorResponse error response.
type ErrorResponse struct {
	// Error - Error description
	Error *ErrorDefinition `json:"error,omitempty"`
}

// Operation microsoft.DelegatedNetwork REST API operation definition
type Operation struct {
	// Name - READ-ONLY; Operation name: {provider}/{resource}/{operation}.
	Name *string `json:"name,omitempty"`
	// Origin - Origin of the operation
	Origin *string `json:"origin,omitempty"`
	// IsDataAction - Gets or sets a value indicating whether the operation is a data action or not.
	IsDataAction *bool `json:"isDataAction,omitempty"`
	// Display - Operation properties display
	Display *OperationDisplay `json:"display,omitempty"`
	// Properties - Properties of the operation
	Properties interface{} `json:"properties,omitempty"`
}

// MarshalJSON is the custom marshaler for Operation.
func (o Operation) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if o.Origin != nil {
		objectMap["origin"] = o.Origin
	}
	if o.IsDataAction != nil {
		objectMap["isDataAction"] = o.IsDataAction
	}
	if o.Display != nil {
		objectMap["display"] = o.Display
	}
	if o.Properties != nil {
		objectMap["properties"] = o.Properties
	}
	return json.Marshal(objectMap)
}

// OperationDisplay the object that represents the operation.
type OperationDisplay struct {
	// Provider - READ-ONLY; Service provider: Microsoft.DelegatedNetwork.
	Provider *string `json:"provider,omitempty"`
	// Resource - READ-ONLY; Resource on which the operation is performed: controller, etc.
	Resource *string `json:"resource,omitempty"`
	// Operation - READ-ONLY; Operation type: create, get, delete, etc.
	Operation *string `json:"operation,omitempty"`
	// Description - READ-ONLY; Friendly description for the operation,
	Description *string `json:"description,omitempty"`
}

// OperationListResult result of request to list controller operations.It contains a list of operations and a
// URL link to get the next set of results
type OperationListResult struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; List of operations supported by the Microsoft.DelegatedNetwork resource provider.
	Value *[]Operation `json:"value,omitempty"`
	// NextLink - URL to get the next set of operation list results if there are any.
	NextLink *string `json:"nextLink,omitempty"`
}

// MarshalJSON is the custom marshaler for OperationListResult.
func (olr OperationListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if olr.NextLink != nil {
		objectMap["nextLink"] = olr.NextLink
	}
	return json.Marshal(objectMap)
}

// OperationListResultIterator provides access to a complete listing of Operation values.
type OperationListResultIterator struct {
	i    int
	page OperationListResultPage
}

// NextWithContext advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
func (iter *OperationListResultIterator) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/OperationListResultIterator.NextWithContext")
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
func (iter *OperationListResultIterator) Next() error {
	return iter.NextWithContext(context.Background())
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

// Creates a new instance of the OperationListResultIterator type.
func NewOperationListResultIterator(page OperationListResultPage) OperationListResultIterator {
	return OperationListResultIterator{page: page}
}

// IsEmpty returns true if the ListResult contains no values.
func (olr OperationListResult) IsEmpty() bool {
	return olr.Value == nil || len(*olr.Value) == 0
}

// hasNextLink returns true if the NextLink is not empty.
func (olr OperationListResult) hasNextLink() bool {
	return olr.NextLink != nil && len(*olr.NextLink) != 0
}

// operationListResultPreparer prepares a request to retrieve the next set of results.
// It returns nil if no more results exist.
func (olr OperationListResult) operationListResultPreparer(ctx context.Context) (*http.Request, error) {
	if !olr.hasNextLink() {
		return nil, nil
	}
	return autorest.Prepare((&http.Request{}).WithContext(ctx),
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(olr.NextLink)))
}

// OperationListResultPage contains a page of Operation values.
type OperationListResultPage struct {
	fn  func(context.Context, OperationListResult) (OperationListResult, error)
	olr OperationListResult
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *OperationListResultPage) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/OperationListResultPage.NextWithContext")
		defer func() {
			sc := -1
			if page.Response().Response.Response != nil {
				sc = page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	for {
		next, err := page.fn(ctx, page.olr)
		if err != nil {
			return err
		}
		page.olr = next
		if !next.hasNextLink() || !next.IsEmpty() {
			break
		}
	}
	return nil
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (page *OperationListResultPage) Next() error {
	return page.NextWithContext(context.Background())
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

// Creates a new instance of the OperationListResultPage type.
func NewOperationListResultPage(getNextPage func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage {
	return OperationListResultPage{fn: getNextPage}
}
