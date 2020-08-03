package customerlockbox

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
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// The package's fully qualified name.
const fqdn = "github.com/Azure/azure-sdk-for-go/services/preview/customerlockbox/mgmt/2018-02-28-preview/customerlockbox"

// Approval request content object, in the use of Approve or Deny a Lockbox request.
type Approval struct {
	autorest.Response `json:"-"`
	// Decision - Approval decision to the Lockbox request. Possible values include: 'Approve', 'Deny'
	Decision Decision `json:"decision,omitempty"`
	// Reason - Reason of the decision
	Reason *string `json:"reason,omitempty"`
}

// ErrorAdditionalInfo an error additional info for the Lockbox service.
type ErrorAdditionalInfo struct {
	// Type - The type of error info.
	Type *string                  `json:"type,omitempty"`
	Info *ErrorAdditionalInfoInfo `json:"info,omitempty"`
}

// ErrorAdditionalInfoInfo ...
type ErrorAdditionalInfoInfo struct {
	// CurrentStatus - The current status/state of the request quired. Possible values include: 'Initializing', 'Pending', 'Approving', 'Denying', 'Approved', 'Denied', 'Expired', 'Revoking', 'Revoked', 'Error', 'Unknown', 'Completed', 'Completing'
	CurrentStatus Status `json:"currentStatus,omitempty"`
}

// ErrorBody an error response body from the Lockbox service.
type ErrorBody struct {
	// Code - An identifier for the error. Codes are invariant and are intended to be consumed programmatically.
	Code *string `json:"code,omitempty"`
	// Message - A message describing the error, intended to be suitable for display in a user interface.
	Message *string `json:"message,omitempty"`
	// Target - The target of the particular error. For example, the name of the property in error.
	Target *string `json:"target,omitempty"`
	// AdditionalInfo - A list of error details about the error.
	AdditionalInfo *[]ErrorAdditionalInfo `json:"additionalInfo,omitempty"`
}

// ErrorResponse an error response from the Lockbox service.
type ErrorResponse struct {
	Error *ErrorBody `json:"error,omitempty"`
}

// LockboxRequestResponse a Lockbox request response object, containing all information associated with the
// request.
type LockboxRequestResponse struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The Arm resource id of the Lockbox request.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the Lockbox request.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The type of the Lockbox request.
	Type *string `json:"type,omitempty"`
	// Properties - The properties that are associated with a lockbox request.
	Properties *LockboxRequestResponseProperties `json:"properties,omitempty"`
}

// MarshalJSON is the custom marshaler for LockboxRequestResponse.
func (lrr LockboxRequestResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if lrr.Properties != nil {
		objectMap["properties"] = lrr.Properties
	}
	return json.Marshal(objectMap)
}

// LockboxRequestResponseProperties the properties that are associated with a lockbox request.
type LockboxRequestResponseProperties struct {
	// RequestID - READ-ONLY; The Lockbox request ID.
	RequestID *string `json:"requestId,omitempty"`
	// Justification - READ-ONLY; The justification of the requestor.
	Justification *string `json:"justification,omitempty"`
	// Status - The status of the request. Possible values include: 'Initializing', 'Pending', 'Approving', 'Denying', 'Approved', 'Denied', 'Expired', 'Revoking', 'Revoked', 'Error', 'Unknown', 'Completed', 'Completing'
	Status Status `json:"status,omitempty"`
	// CreatedDateTime - READ-ONLY; The creation time of the request.
	CreatedDateTime *date.Time `json:"createdDateTime,omitempty"`
	// ExpirationDateTime - READ-ONLY; The expiration time of the request.
	ExpirationDateTime *date.Time `json:"expirationDateTime,omitempty"`
	// Duration - READ-ONLY; The duration of the request in hours.
	Duration *int32 `json:"duration,omitempty"`
	// RequestedResourceIds - READ-ONLY; A list of resource IDs associated with the Lockbox request separated by ','.
	RequestedResourceIds *[]string `json:"requestedResourceIds,omitempty"`
	// ResourceType - READ-ONLY; The resource type of the requested resources.
	ResourceType *string `json:"resourceType,omitempty"`
	// SupportRequest - READ-ONLY; The id of the support request associated.
	SupportRequest *string `json:"supportRequest,omitempty"`
	// SupportCaseURL - READ-ONLY; The url of the support case.
	SupportCaseURL *string `json:"supportCaseUrl,omitempty"`
	// SubscriptionID - READ-ONLY; The subscription ID.
	SubscriptionID *string `json:"subscriptionId,omitempty"`
}

// MarshalJSON is the custom marshaler for LockboxRequestResponseProperties.
func (lrrp LockboxRequestResponseProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if lrrp.Status != "" {
		objectMap["status"] = lrrp.Status
	}
	return json.Marshal(objectMap)
}

// Operation operation result model for ARM RP
type Operation struct {
	// Name - READ-ONLY; Gets or sets action name
	Name *string `json:"name,omitempty"`
	// IsDataAction - READ-ONLY; Gets or sets a value indicating whether it is a data plane action
	IsDataAction *string `json:"isDataAction,omitempty"`
	// Display - READ-ONLY; Contains the localized display information for this particular operation / action.
	Display *OperationDisplay `json:"display,omitempty"`
	// Properties - READ-ONLY; Gets or sets properties
	Properties *string `json:"properties,omitempty"`
	// Origin - READ-ONLY; Gets or sets origin
	Origin *string `json:"origin,omitempty"`
}

// OperationDisplay contains the localized display information for this particular operation / action.
type OperationDisplay struct {
	// Provider - READ-ONLY; The localized friendly form of the resource provider name.
	Provider *string `json:"provider,omitempty"`
	// Resource - READ-ONLY; The localized friendly form of the resource type related to this action/operation.
	Resource *string `json:"resource,omitempty"`
	// Operation - READ-ONLY; The localized friendly name for the operation.
	Operation *string `json:"operation,omitempty"`
	// Description - READ-ONLY; The localized friendly description for the operation.
	Description *string `json:"description,omitempty"`
}

// OperationListResult result of the request to list Customer Lockbox operations. It contains a list of
// operations.
type OperationListResult struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; List of Customer Lockbox operations supported by the Microsoft.StreamAnalytics resource provider.
	Value *[]Operation `json:"value,omitempty"`
	// NextLink - READ-ONLY; URL to get the next set of operation list results if there are any.
	NextLink *string `json:"nextLink,omitempty"`
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

// RequestListResult object containing a list of streaming jobs.
type RequestListResult struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; A list of Lockbox requests. Populated by a 'List' operation.
	Value *[]LockboxRequestResponse `json:"value,omitempty"`
	// NextLink - READ-ONLY; URL to get the next set of operation list results if there are any.
	NextLink *string `json:"nextLink,omitempty"`
}

// RequestListResultIterator provides access to a complete listing of LockboxRequestResponse values.
type RequestListResultIterator struct {
	i    int
	page RequestListResultPage
}

// NextWithContext advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
func (iter *RequestListResultIterator) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/RequestListResultIterator.NextWithContext")
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
func (iter *RequestListResultIterator) Next() error {
	return iter.NextWithContext(context.Background())
}

// NotDone returns true if the enumeration should be started or is not yet complete.
func (iter RequestListResultIterator) NotDone() bool {
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}

// Response returns the raw server response from the last page request.
func (iter RequestListResultIterator) Response() RequestListResult {
	return iter.page.Response()
}

// Value returns the current value or a zero-initialized value if the
// iterator has advanced beyond the end of the collection.
func (iter RequestListResultIterator) Value() LockboxRequestResponse {
	if !iter.page.NotDone() {
		return LockboxRequestResponse{}
	}
	return iter.page.Values()[iter.i]
}

// Creates a new instance of the RequestListResultIterator type.
func NewRequestListResultIterator(page RequestListResultPage) RequestListResultIterator {
	return RequestListResultIterator{page: page}
}

// IsEmpty returns true if the ListResult contains no values.
func (rlr RequestListResult) IsEmpty() bool {
	return rlr.Value == nil || len(*rlr.Value) == 0
}

// hasNextLink returns true if the NextLink is not empty.
func (rlr RequestListResult) hasNextLink() bool {
	return rlr.NextLink != nil && len(*rlr.NextLink) != 0
}

// requestListResultPreparer prepares a request to retrieve the next set of results.
// It returns nil if no more results exist.
func (rlr RequestListResult) requestListResultPreparer(ctx context.Context) (*http.Request, error) {
	if !rlr.hasNextLink() {
		return nil, nil
	}
	return autorest.Prepare((&http.Request{}).WithContext(ctx),
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(rlr.NextLink)))
}

// RequestListResultPage contains a page of LockboxRequestResponse values.
type RequestListResultPage struct {
	fn  func(context.Context, RequestListResult) (RequestListResult, error)
	rlr RequestListResult
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *RequestListResultPage) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/RequestListResultPage.NextWithContext")
		defer func() {
			sc := -1
			if page.Response().Response.Response != nil {
				sc = page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	for {
		next, err := page.fn(ctx, page.rlr)
		if err != nil {
			return err
		}
		page.rlr = next
		if !next.hasNextLink() || !next.IsEmpty() {
			break
		}
	}
	return nil
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (page *RequestListResultPage) Next() error {
	return page.NextWithContext(context.Background())
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page RequestListResultPage) NotDone() bool {
	return !page.rlr.IsEmpty()
}

// Response returns the raw server response from the last page request.
func (page RequestListResultPage) Response() RequestListResult {
	return page.rlr
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page RequestListResultPage) Values() []LockboxRequestResponse {
	if page.rlr.IsEmpty() {
		return nil
	}
	return *page.rlr.Value
}

// Creates a new instance of the RequestListResultPage type.
func NewRequestListResultPage(getNextPage func(context.Context, RequestListResult) (RequestListResult, error)) RequestListResultPage {
	return RequestListResultPage{fn: getNextPage}
}
