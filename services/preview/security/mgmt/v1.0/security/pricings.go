package security

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
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// PricingsClient is the API spec for Microsoft.Security (Azure Security Center) resource provider
type PricingsClient struct {
	BaseClient
}

// NewPricingsClient creates an instance of the PricingsClient client.
func NewPricingsClient(subscriptionID string, ascLocation string) PricingsClient {
	return NewPricingsClientWithBaseURI(DefaultBaseURI, subscriptionID, ascLocation)
}

// NewPricingsClientWithBaseURI creates an instance of the PricingsClient client using a custom endpoint.  Use this
// when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewPricingsClientWithBaseURI(baseURI string, subscriptionID string, ascLocation string) PricingsClient {
	return PricingsClient{NewWithBaseURI(baseURI, subscriptionID, ascLocation)}
}

// CreateOrUpdateResourceGroupPricing security pricing configuration in the resource group
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription. The name is case
// insensitive.
// pricingName - name of the pricing configuration
// pricing - pricing object
func (client PricingsClient) CreateOrUpdateResourceGroupPricing(ctx context.Context, resourceGroupName string, pricingName string, pricing Pricing) (result Pricing, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PricingsClient.CreateOrUpdateResourceGroupPricing")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.PricingsClient", "CreateOrUpdateResourceGroupPricing", err.Error())
	}

	req, err := client.CreateOrUpdateResourceGroupPricingPreparer(ctx, resourceGroupName, pricingName, pricing)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "CreateOrUpdateResourceGroupPricing", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateOrUpdateResourceGroupPricingSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "CreateOrUpdateResourceGroupPricing", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResourceGroupPricingResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "CreateOrUpdateResourceGroupPricing", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdateResourceGroupPricingPreparer prepares the CreateOrUpdateResourceGroupPricing request.
func (client PricingsClient) CreateOrUpdateResourceGroupPricingPreparer(ctx context.Context, resourceGroupName string, pricingName string, pricing Pricing) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"pricingName":       autorest.Encode("path", pricingName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-08-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/pricings/{pricingName}", pathParameters),
		autorest.WithJSON(pricing),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateResourceGroupPricingSender sends the CreateOrUpdateResourceGroupPricing request. The method will close the
// http.Response Body if it receives an error.
func (client PricingsClient) CreateOrUpdateResourceGroupPricingSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// CreateOrUpdateResourceGroupPricingResponder handles the response to the CreateOrUpdateResourceGroupPricing request. The method always
// closes the http.Response Body.
func (client PricingsClient) CreateOrUpdateResourceGroupPricingResponder(resp *http.Response) (result Pricing, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetResourceGroupPricing security pricing configuration in the resource group
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription. The name is case
// insensitive.
// pricingName - name of the pricing configuration
func (client PricingsClient) GetResourceGroupPricing(ctx context.Context, resourceGroupName string, pricingName string) (result Pricing, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PricingsClient.GetResourceGroupPricing")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.PricingsClient", "GetResourceGroupPricing", err.Error())
	}

	req, err := client.GetResourceGroupPricingPreparer(ctx, resourceGroupName, pricingName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "GetResourceGroupPricing", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetResourceGroupPricingSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "GetResourceGroupPricing", resp, "Failure sending request")
		return
	}

	result, err = client.GetResourceGroupPricingResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "GetResourceGroupPricing", resp, "Failure responding to request")
	}

	return
}

// GetResourceGroupPricingPreparer prepares the GetResourceGroupPricing request.
func (client PricingsClient) GetResourceGroupPricingPreparer(ctx context.Context, resourceGroupName string, pricingName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"pricingName":       autorest.Encode("path", pricingName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-08-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/pricings/{pricingName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetResourceGroupPricingSender sends the GetResourceGroupPricing request. The method will close the
// http.Response Body if it receives an error.
func (client PricingsClient) GetResourceGroupPricingSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResourceGroupPricingResponder handles the response to the GetResourceGroupPricing request. The method always
// closes the http.Response Body.
func (client PricingsClient) GetResourceGroupPricingResponder(resp *http.Response) (result Pricing, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetSubscriptionPricing security pricing configuration in the subscriptionSecurity pricing configuration in the
// subscription
// Parameters:
// pricingName - name of the pricing configuration
func (client PricingsClient) GetSubscriptionPricing(ctx context.Context, pricingName string) (result Pricing, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PricingsClient.GetSubscriptionPricing")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.PricingsClient", "GetSubscriptionPricing", err.Error())
	}

	req, err := client.GetSubscriptionPricingPreparer(ctx, pricingName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "GetSubscriptionPricing", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSubscriptionPricingSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "GetSubscriptionPricing", resp, "Failure sending request")
		return
	}

	result, err = client.GetSubscriptionPricingResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "GetSubscriptionPricing", resp, "Failure responding to request")
	}

	return
}

// GetSubscriptionPricingPreparer prepares the GetSubscriptionPricing request.
func (client PricingsClient) GetSubscriptionPricingPreparer(ctx context.Context, pricingName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"pricingName":    autorest.Encode("path", pricingName),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-08-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Security/pricings/{pricingName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSubscriptionPricingSender sends the GetSubscriptionPricing request. The method will close the
// http.Response Body if it receives an error.
func (client PricingsClient) GetSubscriptionPricingSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetSubscriptionPricingResponder handles the response to the GetSubscriptionPricing request. The method always
// closes the http.Response Body.
func (client PricingsClient) GetSubscriptionPricingResponder(resp *http.Response) (result Pricing, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List security pricing configurations in the subscription
func (client PricingsClient) List(ctx context.Context) (result PricingListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PricingsClient.List")
		defer func() {
			sc := -1
			if result.pl.Response.Response != nil {
				sc = result.pl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.PricingsClient", "List", err.Error())
	}

	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.pl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "List", resp, "Failure sending request")
		return
	}

	result.pl, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "List", resp, "Failure responding to request")
	}
	if result.pl.hasNextLink() && result.pl.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListPreparer prepares the List request.
func (client PricingsClient) ListPreparer(ctx context.Context) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-08-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Security/pricings", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client PricingsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client PricingsClient) ListResponder(resp *http.Response) (result PricingList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client PricingsClient) listNextResults(ctx context.Context, lastResults PricingList) (result PricingList, err error) {
	req, err := lastResults.pricingListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "security.PricingsClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "security.PricingsClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client PricingsClient) ListComplete(ctx context.Context) (result PricingListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PricingsClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx)
	return
}

// ListByResourceGroup security pricing configurations in the resource group
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription. The name is case
// insensitive.
func (client PricingsClient) ListByResourceGroup(ctx context.Context, resourceGroupName string) (result PricingListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PricingsClient.ListByResourceGroup")
		defer func() {
			sc := -1
			if result.pl.Response.Response != nil {
				sc = result.pl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.PricingsClient", "ListByResourceGroup", err.Error())
	}

	result.fn = client.listByResourceGroupNextResults
	req, err := client.ListByResourceGroupPreparer(ctx, resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "ListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.pl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "ListByResourceGroup", resp, "Failure sending request")
		return
	}

	result.pl, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "ListByResourceGroup", resp, "Failure responding to request")
	}
	if result.pl.hasNextLink() && result.pl.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListByResourceGroupPreparer prepares the ListByResourceGroup request.
func (client PricingsClient) ListByResourceGroupPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-08-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/pricings", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByResourceGroupSender sends the ListByResourceGroup request. The method will close the
// http.Response Body if it receives an error.
func (client PricingsClient) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByResourceGroupResponder handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (client PricingsClient) ListByResourceGroupResponder(resp *http.Response) (result PricingList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByResourceGroupNextResults retrieves the next set of results, if any.
func (client PricingsClient) listByResourceGroupNextResults(ctx context.Context, lastResults PricingList) (result PricingList, err error) {
	req, err := lastResults.pricingListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "security.PricingsClient", "listByResourceGroupNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "security.PricingsClient", "listByResourceGroupNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "listByResourceGroupNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByResourceGroupComplete enumerates all values, automatically crossing page boundaries as required.
func (client PricingsClient) ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result PricingListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PricingsClient.ListByResourceGroup")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByResourceGroup(ctx, resourceGroupName)
	return
}

// UpdateSubscriptionPricing security pricing configuration in the subscription
// Parameters:
// pricingName - name of the pricing configuration
// pricing - pricing object
func (client PricingsClient) UpdateSubscriptionPricing(ctx context.Context, pricingName string, pricing Pricing) (result Pricing, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PricingsClient.UpdateSubscriptionPricing")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.PricingsClient", "UpdateSubscriptionPricing", err.Error())
	}

	req, err := client.UpdateSubscriptionPricingPreparer(ctx, pricingName, pricing)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "UpdateSubscriptionPricing", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateSubscriptionPricingSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "UpdateSubscriptionPricing", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateSubscriptionPricingResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.PricingsClient", "UpdateSubscriptionPricing", resp, "Failure responding to request")
	}

	return
}

// UpdateSubscriptionPricingPreparer prepares the UpdateSubscriptionPricing request.
func (client PricingsClient) UpdateSubscriptionPricingPreparer(ctx context.Context, pricingName string, pricing Pricing) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"pricingName":    autorest.Encode("path", pricingName),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-08-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Security/pricings/{pricingName}", pathParameters),
		autorest.WithJSON(pricing),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSubscriptionPricingSender sends the UpdateSubscriptionPricing request. The method will close the
// http.Response Body if it receives an error.
func (client PricingsClient) UpdateSubscriptionPricingSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// UpdateSubscriptionPricingResponder handles the response to the UpdateSubscriptionPricing request. The method always
// closes the http.Response Body.
func (client PricingsClient) UpdateSubscriptionPricingResponder(resp *http.Response) (result Pricing, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
