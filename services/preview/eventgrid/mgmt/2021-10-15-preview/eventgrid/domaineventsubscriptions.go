package eventgrid

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// DomainEventSubscriptionsClient is the azure EventGrid Management Client
type DomainEventSubscriptionsClient struct {
	BaseClient
}

// NewDomainEventSubscriptionsClient creates an instance of the DomainEventSubscriptionsClient client.
func NewDomainEventSubscriptionsClient(subscriptionID string) DomainEventSubscriptionsClient {
	return NewDomainEventSubscriptionsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewDomainEventSubscriptionsClientWithBaseURI creates an instance of the DomainEventSubscriptionsClient client using
// a custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign
// clouds, Azure stack).
func NewDomainEventSubscriptionsClientWithBaseURI(baseURI string, subscriptionID string) DomainEventSubscriptionsClient {
	return DomainEventSubscriptionsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate asynchronously creates a new event subscription or updates an existing event subscription.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription.
// domainName - name of the domain topic.
// eventSubscriptionName - name of the event subscription to be created. Event subscription names must be
// between 3 and 100 characters in length and use alphanumeric letters only.
// eventSubscriptionInfo - event subscription properties containing the destination and filter information.
func (client DomainEventSubscriptionsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string, eventSubscriptionInfo EventSubscription) (result DomainEventSubscriptionsCreateOrUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DomainEventSubscriptionsClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, domainName, eventSubscriptionName, eventSubscriptionInfo)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client DomainEventSubscriptionsClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string, eventSubscriptionInfo EventSubscription) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"domainName":            autorest.Encode("path", domainName),
		"eventSubscriptionName": autorest.Encode("path", eventSubscriptionName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-10-15-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	eventSubscriptionInfo.SystemData = nil
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/domains/{domainName}/eventSubscriptions/{eventSubscriptionName}", pathParameters),
		autorest.WithJSON(eventSubscriptionInfo),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client DomainEventSubscriptionsClient) CreateOrUpdateSender(req *http.Request) (future DomainEventSubscriptionsCreateOrUpdateFuture, err error) {
	var resp *http.Response
	future.FutureAPI = &azure.Future{}
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client DomainEventSubscriptionsClient) CreateOrUpdateResponder(resp *http.Response) (result EventSubscription, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete delete an existing event subscription for a domain.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription.
// domainName - name of the domain.
// eventSubscriptionName - name of the event subscription to be deleted. Event subscription names must be
// between 3 and 100 characters in length and use alphanumeric letters only.
func (client DomainEventSubscriptionsClient) Delete(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string) (result DomainEventSubscriptionsDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DomainEventSubscriptionsClient.Delete")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, resourceGroupName, domainName, eventSubscriptionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "Delete", result.Response(), "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client DomainEventSubscriptionsClient) DeletePreparer(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"domainName":            autorest.Encode("path", domainName),
		"eventSubscriptionName": autorest.Encode("path", eventSubscriptionName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-10-15-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/domains/{domainName}/eventSubscriptions/{eventSubscriptionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client DomainEventSubscriptionsClient) DeleteSender(req *http.Request) (future DomainEventSubscriptionsDeleteFuture, err error) {
	var resp *http.Response
	future.FutureAPI = &azure.Future{}
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client DomainEventSubscriptionsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get get properties of an event subscription of a domain.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription.
// domainName - name of the partner topic.
// eventSubscriptionName - name of the event subscription to be found. Event subscription names must be between
// 3 and 100 characters in length and use alphanumeric letters only.
func (client DomainEventSubscriptionsClient) Get(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string) (result EventSubscription, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DomainEventSubscriptionsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, resourceGroupName, domainName, eventSubscriptionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client DomainEventSubscriptionsClient) GetPreparer(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"domainName":            autorest.Encode("path", domainName),
		"eventSubscriptionName": autorest.Encode("path", eventSubscriptionName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-10-15-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/domains/{domainName}/eventSubscriptions/{eventSubscriptionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client DomainEventSubscriptionsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client DomainEventSubscriptionsClient) GetResponder(resp *http.Response) (result EventSubscription, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetDeliveryAttributes get all delivery attributes for an event subscription for domain.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription.
// domainName - name of the domain topic.
// eventSubscriptionName - name of the event subscription.
func (client DomainEventSubscriptionsClient) GetDeliveryAttributes(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string) (result DeliveryAttributeListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DomainEventSubscriptionsClient.GetDeliveryAttributes")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetDeliveryAttributesPreparer(ctx, resourceGroupName, domainName, eventSubscriptionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "GetDeliveryAttributes", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetDeliveryAttributesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "GetDeliveryAttributes", resp, "Failure sending request")
		return
	}

	result, err = client.GetDeliveryAttributesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "GetDeliveryAttributes", resp, "Failure responding to request")
		return
	}

	return
}

// GetDeliveryAttributesPreparer prepares the GetDeliveryAttributes request.
func (client DomainEventSubscriptionsClient) GetDeliveryAttributesPreparer(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"domainName":            autorest.Encode("path", domainName),
		"eventSubscriptionName": autorest.Encode("path", eventSubscriptionName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-10-15-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/domains/{domainName}/eventSubscriptions/{eventSubscriptionName}/getDeliveryAttributes", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetDeliveryAttributesSender sends the GetDeliveryAttributes request. The method will close the
// http.Response Body if it receives an error.
func (client DomainEventSubscriptionsClient) GetDeliveryAttributesSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetDeliveryAttributesResponder handles the response to the GetDeliveryAttributes request. The method always
// closes the http.Response Body.
func (client DomainEventSubscriptionsClient) GetDeliveryAttributesResponder(resp *http.Response) (result DeliveryAttributeListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetFullURL get the full endpoint URL for an event subscription for domain.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription.
// domainName - name of the domain topic.
// eventSubscriptionName - name of the event subscription.
func (client DomainEventSubscriptionsClient) GetFullURL(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string) (result EventSubscriptionFullURL, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DomainEventSubscriptionsClient.GetFullURL")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetFullURLPreparer(ctx, resourceGroupName, domainName, eventSubscriptionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "GetFullURL", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetFullURLSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "GetFullURL", resp, "Failure sending request")
		return
	}

	result, err = client.GetFullURLResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "GetFullURL", resp, "Failure responding to request")
		return
	}

	return
}

// GetFullURLPreparer prepares the GetFullURL request.
func (client DomainEventSubscriptionsClient) GetFullURLPreparer(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"domainName":            autorest.Encode("path", domainName),
		"eventSubscriptionName": autorest.Encode("path", eventSubscriptionName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-10-15-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/domains/{domainName}/eventSubscriptions/{eventSubscriptionName}/getFullUrl", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetFullURLSender sends the GetFullURL request. The method will close the
// http.Response Body if it receives an error.
func (client DomainEventSubscriptionsClient) GetFullURLSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetFullURLResponder handles the response to the GetFullURL request. The method always
// closes the http.Response Body.
func (client DomainEventSubscriptionsClient) GetFullURLResponder(resp *http.Response) (result EventSubscriptionFullURL, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List list all event subscriptions that have been created for a specific topic.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription.
// domainName - name of the domain.
// filter - the query used to filter the search results using OData syntax. Filtering is permitted on the
// 'name' property only and with limited number of OData operations. These operations are: the 'contains'
// function as well as the following logical operations: not, and, or, eq (for equal), and ne (for not equal).
// No arithmetic operations are supported. The following is a valid filter example: $filter=contains(namE,
// 'PATTERN') and name ne 'PATTERN-1'. The following is not a valid filter example: $filter=location eq
// 'westus'.
// top - the number of results to return per page for the list operation. Valid range for top parameter is 1 to
// 100. If not specified, the default number of results to be returned is 20 items per page.
func (client DomainEventSubscriptionsClient) List(ctx context.Context, resourceGroupName string, domainName string, filter string, top *int32) (result EventSubscriptionsListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DomainEventSubscriptionsClient.List")
		defer func() {
			sc := -1
			if result.eslr.Response.Response != nil {
				sc = result.eslr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName, domainName, filter, top)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.eslr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "List", resp, "Failure sending request")
		return
	}

	result.eslr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "List", resp, "Failure responding to request")
		return
	}
	if result.eslr.hasNextLink() && result.eslr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client DomainEventSubscriptionsClient) ListPreparer(ctx context.Context, resourceGroupName string, domainName string, filter string, top *int32) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"domainName":        autorest.Encode("path", domainName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-10-15-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/domains/{domainName}/eventSubscriptions", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client DomainEventSubscriptionsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client DomainEventSubscriptionsClient) ListResponder(resp *http.Response) (result EventSubscriptionsListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client DomainEventSubscriptionsClient) listNextResults(ctx context.Context, lastResults EventSubscriptionsListResult) (result EventSubscriptionsListResult, err error) {
	req, err := lastResults.eventSubscriptionsListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client DomainEventSubscriptionsClient) ListComplete(ctx context.Context, resourceGroupName string, domainName string, filter string, top *int32) (result EventSubscriptionsListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DomainEventSubscriptionsClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, resourceGroupName, domainName, filter, top)
	return
}

// Update update an existing event subscription for a topic.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription.
// domainName - name of the domain.
// eventSubscriptionName - name of the event subscription to be updated.
// eventSubscriptionUpdateParameters - updated event subscription information.
func (client DomainEventSubscriptionsClient) Update(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string, eventSubscriptionUpdateParameters EventSubscriptionUpdateParameters) (result DomainEventSubscriptionsUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DomainEventSubscriptionsClient.Update")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UpdatePreparer(ctx, resourceGroupName, domainName, eventSubscriptionName, eventSubscriptionUpdateParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "Update", nil, "Failure preparing request")
		return
	}

	result, err = client.UpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.DomainEventSubscriptionsClient", "Update", result.Response(), "Failure sending request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client DomainEventSubscriptionsClient) UpdatePreparer(ctx context.Context, resourceGroupName string, domainName string, eventSubscriptionName string, eventSubscriptionUpdateParameters EventSubscriptionUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"domainName":            autorest.Encode("path", domainName),
		"eventSubscriptionName": autorest.Encode("path", eventSubscriptionName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-10-15-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/domains/{domainName}/eventSubscriptions/{eventSubscriptionName}", pathParameters),
		autorest.WithJSON(eventSubscriptionUpdateParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client DomainEventSubscriptionsClient) UpdateSender(req *http.Request) (future DomainEventSubscriptionsUpdateFuture, err error) {
	var resp *http.Response
	future.FutureAPI = &azure.Future{}
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client DomainEventSubscriptionsClient) UpdateResponder(resp *http.Response) (result EventSubscription, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
