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

// DiscoveredSecuritySolutionsClient is the API spec for Microsoft.Security (Azure Security Center) resource provider
type DiscoveredSecuritySolutionsClient struct {
	BaseClient
}

// NewDiscoveredSecuritySolutionsClient creates an instance of the DiscoveredSecuritySolutionsClient client.
func NewDiscoveredSecuritySolutionsClient(subscriptionID string, ascLocation string) DiscoveredSecuritySolutionsClient {
	return NewDiscoveredSecuritySolutionsClientWithBaseURI(DefaultBaseURI, subscriptionID, ascLocation)
}

// NewDiscoveredSecuritySolutionsClientWithBaseURI creates an instance of the DiscoveredSecuritySolutionsClient client
// using a custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign
// clouds, Azure stack).
func NewDiscoveredSecuritySolutionsClientWithBaseURI(baseURI string, subscriptionID string, ascLocation string) DiscoveredSecuritySolutionsClient {
	return DiscoveredSecuritySolutionsClient{NewWithBaseURI(baseURI, subscriptionID, ascLocation)}
}

// Get gets a specific discovered Security Solution.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription. The name is case
// insensitive.
// discoveredSecuritySolutionName - name of a discovered security solution.
func (client DiscoveredSecuritySolutionsClient) Get(ctx context.Context, resourceGroupName string, discoveredSecuritySolutionName string) (result DiscoveredSecuritySolution, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DiscoveredSecuritySolutionsClient.Get")
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
		return result, validation.NewError("security.DiscoveredSecuritySolutionsClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, discoveredSecuritySolutionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client DiscoveredSecuritySolutionsClient) GetPreparer(ctx context.Context, resourceGroupName string, discoveredSecuritySolutionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"ascLocation":                    autorest.Encode("path", client.AscLocation),
		"discoveredSecuritySolutionName": autorest.Encode("path", discoveredSecuritySolutionName),
		"resourceGroupName":              autorest.Encode("path", resourceGroupName),
		"subscriptionId":                 autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-01-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/locations/{ascLocation}/discoveredSecuritySolutions/{discoveredSecuritySolutionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client DiscoveredSecuritySolutionsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client DiscoveredSecuritySolutionsClient) GetResponder(resp *http.Response) (result DiscoveredSecuritySolution, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List gets a list of discovered Security Solutions for the subscription.
func (client DiscoveredSecuritySolutionsClient) List(ctx context.Context) (result DiscoveredSecuritySolutionListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DiscoveredSecuritySolutionsClient.List")
		defer func() {
			sc := -1
			if result.dssl.Response.Response != nil {
				sc = result.dssl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.DiscoveredSecuritySolutionsClient", "List", err.Error())
	}

	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.dssl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "List", resp, "Failure sending request")
		return
	}

	result.dssl, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client DiscoveredSecuritySolutionsClient) ListPreparer(ctx context.Context) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-01-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Security/discoveredSecuritySolutions", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client DiscoveredSecuritySolutionsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client DiscoveredSecuritySolutionsClient) ListResponder(resp *http.Response) (result DiscoveredSecuritySolutionList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client DiscoveredSecuritySolutionsClient) listNextResults(ctx context.Context, lastResults DiscoveredSecuritySolutionList) (result DiscoveredSecuritySolutionList, err error) {
	req, err := lastResults.discoveredSecuritySolutionListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client DiscoveredSecuritySolutionsClient) ListComplete(ctx context.Context) (result DiscoveredSecuritySolutionListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DiscoveredSecuritySolutionsClient.List")
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

// ListByHomeRegion gets a list of discovered Security Solutions for the subscription and location.
func (client DiscoveredSecuritySolutionsClient) ListByHomeRegion(ctx context.Context) (result DiscoveredSecuritySolutionListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DiscoveredSecuritySolutionsClient.ListByHomeRegion")
		defer func() {
			sc := -1
			if result.dssl.Response.Response != nil {
				sc = result.dssl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("security.DiscoveredSecuritySolutionsClient", "ListByHomeRegion", err.Error())
	}

	result.fn = client.listByHomeRegionNextResults
	req, err := client.ListByHomeRegionPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "ListByHomeRegion", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByHomeRegionSender(req)
	if err != nil {
		result.dssl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "ListByHomeRegion", resp, "Failure sending request")
		return
	}

	result.dssl, err = client.ListByHomeRegionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "ListByHomeRegion", resp, "Failure responding to request")
	}

	return
}

// ListByHomeRegionPreparer prepares the ListByHomeRegion request.
func (client DiscoveredSecuritySolutionsClient) ListByHomeRegionPreparer(ctx context.Context) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"ascLocation":    autorest.Encode("path", client.AscLocation),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-01-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Security/locations/{ascLocation}/discoveredSecuritySolutions", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByHomeRegionSender sends the ListByHomeRegion request. The method will close the
// http.Response Body if it receives an error.
func (client DiscoveredSecuritySolutionsClient) ListByHomeRegionSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByHomeRegionResponder handles the response to the ListByHomeRegion request. The method always
// closes the http.Response Body.
func (client DiscoveredSecuritySolutionsClient) ListByHomeRegionResponder(resp *http.Response) (result DiscoveredSecuritySolutionList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByHomeRegionNextResults retrieves the next set of results, if any.
func (client DiscoveredSecuritySolutionsClient) listByHomeRegionNextResults(ctx context.Context, lastResults DiscoveredSecuritySolutionList) (result DiscoveredSecuritySolutionList, err error) {
	req, err := lastResults.discoveredSecuritySolutionListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "listByHomeRegionNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByHomeRegionSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "listByHomeRegionNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByHomeRegionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.DiscoveredSecuritySolutionsClient", "listByHomeRegionNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByHomeRegionComplete enumerates all values, automatically crossing page boundaries as required.
func (client DiscoveredSecuritySolutionsClient) ListByHomeRegionComplete(ctx context.Context) (result DiscoveredSecuritySolutionListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DiscoveredSecuritySolutionsClient.ListByHomeRegion")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByHomeRegion(ctx)
	return
}
