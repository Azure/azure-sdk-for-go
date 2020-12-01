package subscription

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
const fqdn = "github.com/Azure/azure-sdk-for-go/services/subscription/mgmt/2020-09-01/subscription"

// AliasCreateFuture an abstraction for monitoring and retrieving the results of a long-running operation.
type AliasCreateFuture struct {
	azure.Future
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future *AliasCreateFuture) Result(client AliasClient) (par PutAliasResponse, err error) {
	var done bool
	done, err = future.DoneWithContext(context.Background(), client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscription.AliasCreateFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("subscription.AliasCreateFuture")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if par.Response.Response, err = future.GetResult(sender); err == nil && par.Response.Response.StatusCode != http.StatusNoContent {
		par, err = client.CreateResponder(par.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "subscription.AliasCreateFuture", "Result", par.Response.Response, "Failure responding to request")
		}
	}
	return
}

// CanceledSubscriptionID the ID of the canceled subscription
type CanceledSubscriptionID struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; The ID of the canceled subscription
	Value *string `json:"value,omitempty"`
}

// EnabledSubscriptionID the ID of the subscriptions that is being enabled
type EnabledSubscriptionID struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; The ID of the subscriptions that is being enabled
	Value *string `json:"value,omitempty"`
}

// ErrorResponse describes the format of Error response.
type ErrorResponse struct {
	// Code - Error code
	Code *string `json:"code,omitempty"`
	// Message - Error message indicating why the operation failed.
	Message *string `json:"message,omitempty"`
}

// ErrorResponseBody error response indicates that the service is not able to process the incoming request.
// The reason is provided in the error message.
type ErrorResponseBody struct {
	// Error - The details of the error.
	Error *ErrorResponse `json:"error,omitempty"`
}

// ListResult subscription list operation response.
type ListResult struct {
	autorest.Response `json:"-"`
	// Value - An array of subscriptions.
	Value *[]Model `json:"value,omitempty"`
	// NextLink - The URL to get the next set of results.
	NextLink *string `json:"nextLink,omitempty"`
}

// ListResultIterator provides access to a complete listing of Model values.
type ListResultIterator struct {
	i    int
	page ListResultPage
}

// NextWithContext advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
func (iter *ListResultIterator) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ListResultIterator.NextWithContext")
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
func (iter *ListResultIterator) Next() error {
	return iter.NextWithContext(context.Background())
}

// NotDone returns true if the enumeration should be started or is not yet complete.
func (iter ListResultIterator) NotDone() bool {
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}

// Response returns the raw server response from the last page request.
func (iter ListResultIterator) Response() ListResult {
	return iter.page.Response()
}

// Value returns the current value or a zero-initialized value if the
// iterator has advanced beyond the end of the collection.
func (iter ListResultIterator) Value() Model {
	if !iter.page.NotDone() {
		return Model{}
	}
	return iter.page.Values()[iter.i]
}

// Creates a new instance of the ListResultIterator type.
func NewListResultIterator(page ListResultPage) ListResultIterator {
	return ListResultIterator{page: page}
}

// IsEmpty returns true if the ListResult contains no values.
func (lr ListResult) IsEmpty() bool {
	return lr.Value == nil || len(*lr.Value) == 0
}

// hasNextLink returns true if the NextLink is not empty.
func (lr ListResult) hasNextLink() bool {
	return lr.NextLink != nil && len(*lr.NextLink) != 0
}

// listResultPreparer prepares a request to retrieve the next set of results.
// It returns nil if no more results exist.
func (lr ListResult) listResultPreparer(ctx context.Context) (*http.Request, error) {
	if !lr.hasNextLink() {
		return nil, nil
	}
	return autorest.Prepare((&http.Request{}).WithContext(ctx),
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(lr.NextLink)))
}

// ListResultPage contains a page of Model values.
type ListResultPage struct {
	fn func(context.Context, ListResult) (ListResult, error)
	lr ListResult
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *ListResultPage) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ListResultPage.NextWithContext")
		defer func() {
			sc := -1
			if page.Response().Response.Response != nil {
				sc = page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	for {
		next, err := page.fn(ctx, page.lr)
		if err != nil {
			return err
		}
		page.lr = next
		if !next.hasNextLink() || !next.IsEmpty() {
			break
		}
	}
	return nil
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (page *ListResultPage) Next() error {
	return page.NextWithContext(context.Background())
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page ListResultPage) NotDone() bool {
	return !page.lr.IsEmpty()
}

// Response returns the raw server response from the last page request.
func (page ListResultPage) Response() ListResult {
	return page.lr
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page ListResultPage) Values() []Model {
	if page.lr.IsEmpty() {
		return nil
	}
	return *page.lr.Value
}

// Creates a new instance of the ListResultPage type.
func NewListResultPage(cur ListResult, getNextPage func(context.Context, ListResult) (ListResult, error)) ListResultPage {
	return ListResultPage{
		fn: getNextPage,
		lr: cur,
	}
}

// Location location information.
type Location struct {
	// ID - READ-ONLY; The fully qualified ID of the location. For example, /subscriptions/00000000-0000-0000-0000-000000000000/locations/westus.
	ID *string `json:"id,omitempty"`
	// SubscriptionID - READ-ONLY; The subscription ID.
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	// Name - READ-ONLY; The location name.
	Name *string `json:"name,omitempty"`
	// DisplayName - READ-ONLY; The display name of the location.
	DisplayName *string `json:"displayName,omitempty"`
	// Latitude - READ-ONLY; The latitude of the location.
	Latitude *string `json:"latitude,omitempty"`
	// Longitude - READ-ONLY; The longitude of the location.
	Longitude *string `json:"longitude,omitempty"`
}

// LocationListResult location list operation response.
type LocationListResult struct {
	autorest.Response `json:"-"`
	// Value - An array of locations.
	Value *[]Location `json:"value,omitempty"`
}

// Model subscription information.
type Model struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The fully qualified ID for the subscription. For example, /subscriptions/00000000-0000-0000-0000-000000000000.
	ID *string `json:"id,omitempty"`
	// SubscriptionID - READ-ONLY; The subscription ID.
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	// DisplayName - READ-ONLY; The subscription display name.
	DisplayName *string `json:"displayName,omitempty"`
	// State - READ-ONLY; The subscription state. Possible values are Enabled, Warned, PastDue, Disabled, and Deleted. Possible values include: 'Enabled', 'Warned', 'PastDue', 'Disabled', 'Deleted'
	State State `json:"state,omitempty"`
	// SubscriptionPolicies - The subscription policies.
	SubscriptionPolicies *Policies `json:"subscriptionPolicies,omitempty"`
	// AuthorizationSource - The authorization source of the request. Valid values are one or more combinations of Legacy, RoleBased, Bypassed, Direct and Management. For example, 'Legacy, RoleBased'.
	AuthorizationSource *string `json:"authorizationSource,omitempty"`
}

// MarshalJSON is the custom marshaler for Model.
func (mVar Model) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if mVar.SubscriptionPolicies != nil {
		objectMap["subscriptionPolicies"] = mVar.SubscriptionPolicies
	}
	if mVar.AuthorizationSource != nil {
		objectMap["authorizationSource"] = mVar.AuthorizationSource
	}
	return json.Marshal(objectMap)
}

// Name the new name of the subscription.
type Name struct {
	// SubscriptionName - New subscription name
	SubscriptionName *string `json:"subscriptionName,omitempty"`
}

// Operation REST API operation
type Operation struct {
	// Name - Operation name: {provider}/{resource}/{operation}
	Name *string `json:"name,omitempty"`
	// Display - The object that represents the operation.
	Display *OperationDisplay `json:"display,omitempty"`
}

// OperationDisplay the object that represents the operation.
type OperationDisplay struct {
	// Provider - Service provider: Microsoft.Subscription
	Provider *string `json:"provider,omitempty"`
	// Resource - Resource on which the operation is performed: Profile, endpoint, etc.
	Resource *string `json:"resource,omitempty"`
	// Operation - Operation type: Read, write, delete, etc.
	Operation *string `json:"operation,omitempty"`
}

// OperationListResult result of the request to list operations. It contains a list of operations and a URL
// link to get the next set of results.
type OperationListResult struct {
	autorest.Response `json:"-"`
	// Value - List of operations.
	Value *[]Operation `json:"value,omitempty"`
	// NextLink - URL to get the next set of operation list results if there are any.
	NextLink *string `json:"nextLink,omitempty"`
}

// Policies subscription policies.
type Policies struct {
	// LocationPlacementID - READ-ONLY; The subscription location placement ID. The ID indicates which regions are visible for a subscription. For example, a subscription with a location placement Id of Public_2014-09-01 has access to Azure public regions.
	LocationPlacementID *string `json:"locationPlacementId,omitempty"`
	// QuotaID - READ-ONLY; The subscription quota ID.
	QuotaID *string `json:"quotaId,omitempty"`
	// SpendingLimit - READ-ONLY; The subscription spending limit. Possible values include: 'On', 'Off', 'CurrentPeriodOff'
	SpendingLimit SpendingLimit `json:"spendingLimit,omitempty"`
}

// PutAliasListResult the list of aliases.
type PutAliasListResult struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; The list of alias.
	Value *[]PutAliasResponse `json:"value,omitempty"`
	// NextLink - READ-ONLY; The link (url) to the next page of results.
	NextLink *string `json:"nextLink,omitempty"`
}

// PutAliasRequest the parameters required to create a new subscription.
type PutAliasRequest struct {
	// Properties - Put alias request properties.
	Properties *PutAliasRequestProperties `json:"properties,omitempty"`
}

// PutAliasRequestProperties put subscription properties.
type PutAliasRequestProperties struct {
	// DisplayName - The friendly name of the subscription.
	DisplayName *string `json:"displayName,omitempty"`
	// Workload - The workload type of the subscription. It can be either Production or DevTest. Possible values include: 'Production', 'DevTest'
	Workload Workload `json:"workload,omitempty"`
	// BillingScope - Determines whether subscription is fieldLed, partnerLed or LegacyEA
	BillingScope *string `json:"billingScope,omitempty"`
	// SubscriptionID - This parameter can be used to create alias for existing subscription Id
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	// ResellerID - Reseller ID, basically MPN Id
	ResellerID *string `json:"resellerId,omitempty"`
}

// PutAliasResponse subscription Information with the alias.
type PutAliasResponse struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; Fully qualified ID for the alias resource.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Alias ID.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type, Microsoft.Subscription/aliases.
	Type *string `json:"type,omitempty"`
	// Properties - Put Alias response properties.
	Properties *PutAliasResponseProperties `json:"properties,omitempty"`
}

// MarshalJSON is the custom marshaler for PutAliasResponse.
func (par PutAliasResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if par.Properties != nil {
		objectMap["properties"] = par.Properties
	}
	return json.Marshal(objectMap)
}

// PutAliasResponseProperties put subscription creation result properties.
type PutAliasResponseProperties struct {
	// SubscriptionID - READ-ONLY; Newly created subscription Id.
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	// ProvisioningState - The provisioning state of the resource. Possible values include: 'Accepted', 'Succeeded', 'Failed'
	ProvisioningState ProvisioningState `json:"provisioningState,omitempty"`
}

// MarshalJSON is the custom marshaler for PutAliasResponseProperties.
func (parp PutAliasResponseProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if parp.ProvisioningState != "" {
		objectMap["provisioningState"] = parp.ProvisioningState
	}
	return json.Marshal(objectMap)
}

// RenamedSubscriptionID the ID of the subscriptions that is being renamed
type RenamedSubscriptionID struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; The ID of the subscriptions that is being renamed
	Value *string `json:"value,omitempty"`
}

// TenantIDDescription tenant Id information.
type TenantIDDescription struct {
	// ID - READ-ONLY; The fully qualified ID of the tenant. For example, /tenants/00000000-0000-0000-0000-000000000000.
	ID *string `json:"id,omitempty"`
	// TenantID - READ-ONLY; The tenant ID. For example, 00000000-0000-0000-0000-000000000000.
	TenantID *string `json:"tenantId,omitempty"`
}

// TenantListResult tenant Ids information.
type TenantListResult struct {
	autorest.Response `json:"-"`
	// Value - An array of tenants.
	Value *[]TenantIDDescription `json:"value,omitempty"`
	// NextLink - The URL to use for getting the next set of results.
	NextLink *string `json:"nextLink,omitempty"`
}

// TenantListResultIterator provides access to a complete listing of TenantIDDescription values.
type TenantListResultIterator struct {
	i    int
	page TenantListResultPage
}

// NextWithContext advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
func (iter *TenantListResultIterator) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TenantListResultIterator.NextWithContext")
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
func (iter *TenantListResultIterator) Next() error {
	return iter.NextWithContext(context.Background())
}

// NotDone returns true if the enumeration should be started or is not yet complete.
func (iter TenantListResultIterator) NotDone() bool {
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}

// Response returns the raw server response from the last page request.
func (iter TenantListResultIterator) Response() TenantListResult {
	return iter.page.Response()
}

// Value returns the current value or a zero-initialized value if the
// iterator has advanced beyond the end of the collection.
func (iter TenantListResultIterator) Value() TenantIDDescription {
	if !iter.page.NotDone() {
		return TenantIDDescription{}
	}
	return iter.page.Values()[iter.i]
}

// Creates a new instance of the TenantListResultIterator type.
func NewTenantListResultIterator(page TenantListResultPage) TenantListResultIterator {
	return TenantListResultIterator{page: page}
}

// IsEmpty returns true if the ListResult contains no values.
func (tlr TenantListResult) IsEmpty() bool {
	return tlr.Value == nil || len(*tlr.Value) == 0
}

// hasNextLink returns true if the NextLink is not empty.
func (tlr TenantListResult) hasNextLink() bool {
	return tlr.NextLink != nil && len(*tlr.NextLink) != 0
}

// tenantListResultPreparer prepares a request to retrieve the next set of results.
// It returns nil if no more results exist.
func (tlr TenantListResult) tenantListResultPreparer(ctx context.Context) (*http.Request, error) {
	if !tlr.hasNextLink() {
		return nil, nil
	}
	return autorest.Prepare((&http.Request{}).WithContext(ctx),
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(tlr.NextLink)))
}

// TenantListResultPage contains a page of TenantIDDescription values.
type TenantListResultPage struct {
	fn  func(context.Context, TenantListResult) (TenantListResult, error)
	tlr TenantListResult
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *TenantListResultPage) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TenantListResultPage.NextWithContext")
		defer func() {
			sc := -1
			if page.Response().Response.Response != nil {
				sc = page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	for {
		next, err := page.fn(ctx, page.tlr)
		if err != nil {
			return err
		}
		page.tlr = next
		if !next.hasNextLink() || !next.IsEmpty() {
			break
		}
	}
	return nil
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (page *TenantListResultPage) Next() error {
	return page.NextWithContext(context.Background())
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page TenantListResultPage) NotDone() bool {
	return !page.tlr.IsEmpty()
}

// Response returns the raw server response from the last page request.
func (page TenantListResultPage) Response() TenantListResult {
	return page.tlr
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page TenantListResultPage) Values() []TenantIDDescription {
	if page.tlr.IsEmpty() {
		return nil
	}
	return *page.tlr.Value
}

// Creates a new instance of the TenantListResultPage type.
func NewTenantListResultPage(cur TenantListResult, getNextPage func(context.Context, TenantListResult) (TenantListResult, error)) TenantListResultPage {
	return TenantListResultPage{
		fn:  getNextPage,
		tlr: cur,
	}
}
