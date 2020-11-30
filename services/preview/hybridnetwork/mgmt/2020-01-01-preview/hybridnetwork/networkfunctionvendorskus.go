package hybridnetwork

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

// NetworkFunctionVendorSkusClient is the client for the NetworkFunctionVendorSkus methods of the Hybridnetwork
// service.
type NetworkFunctionVendorSkusClient struct {
	BaseClient
}

// NewNetworkFunctionVendorSkusClient creates an instance of the NetworkFunctionVendorSkusClient client.
func NewNetworkFunctionVendorSkusClient(subscriptionID string) NetworkFunctionVendorSkusClient {
	return NewNetworkFunctionVendorSkusClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewNetworkFunctionVendorSkusClientWithBaseURI creates an instance of the NetworkFunctionVendorSkusClient client
// using a custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign
// clouds, Azure stack).
func NewNetworkFunctionVendorSkusClientWithBaseURI(baseURI string, subscriptionID string) NetworkFunctionVendorSkusClient {
	return NetworkFunctionVendorSkusClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// ListBySku lists information about network function vendor sku details.
// Parameters:
// vendorName - the name of the network function vendor.
// vendorSkuName - the name of the network function sku.
func (client NetworkFunctionVendorSkusClient) ListBySku(ctx context.Context, vendorName string, vendorSkuName string) (result NetworkFunctionSkuDetailsPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/NetworkFunctionVendorSkusClient.ListBySku")
		defer func() {
			sc := -1
			if result.nfsd.Response.Response != nil {
				sc = result.nfsd.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("hybridnetwork.NetworkFunctionVendorSkusClient", "ListBySku", err.Error())
	}

	result.fn = client.listBySkuNextResults
	req, err := client.ListBySkuPreparer(ctx, vendorName, vendorSkuName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "ListBySku", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListBySkuSender(req)
	if err != nil {
		result.nfsd.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "ListBySku", resp, "Failure sending request")
		return
	}

	result.nfsd, err = client.ListBySkuResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "ListBySku", resp, "Failure responding to request")
	}
	if result.nfsd.hasNextLink() && result.nfsd.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListBySkuPreparer prepares the ListBySku request.
func (client NetworkFunctionVendorSkusClient) ListBySkuPreparer(ctx context.Context, vendorName string, vendorSkuName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
		"vendorName":     autorest.Encode("path", vendorName),
		"vendorSkuName":  autorest.Encode("path", vendorSkuName),
	}

	const APIVersion = "2020-01-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.HybridNetwork/networkFunctionVendors/{vendorName}/vendorSkus/{vendorSkuName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListBySkuSender sends the ListBySku request. The method will close the
// http.Response Body if it receives an error.
func (client NetworkFunctionVendorSkusClient) ListBySkuSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListBySkuResponder handles the response to the ListBySku request. The method always
// closes the http.Response Body.
func (client NetworkFunctionVendorSkusClient) ListBySkuResponder(resp *http.Response) (result NetworkFunctionSkuDetails, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listBySkuNextResults retrieves the next set of results, if any.
func (client NetworkFunctionVendorSkusClient) listBySkuNextResults(ctx context.Context, lastResults NetworkFunctionSkuDetails) (result NetworkFunctionSkuDetails, err error) {
	req, err := lastResults.networkFunctionSkuDetailsPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "listBySkuNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListBySkuSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "listBySkuNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListBySkuResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "listBySkuNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListBySkuComplete enumerates all values, automatically crossing page boundaries as required.
func (client NetworkFunctionVendorSkusClient) ListBySkuComplete(ctx context.Context, vendorName string, vendorSkuName string) (result NetworkFunctionSkuDetailsIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/NetworkFunctionVendorSkusClient.ListBySku")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListBySku(ctx, vendorName, vendorSkuName)
	return
}

// ListByVendor lists all network function vendor sku details in a vendor.
// Parameters:
// vendorName - the name of the network function vendor.
func (client NetworkFunctionVendorSkusClient) ListByVendor(ctx context.Context, vendorName string) (result NetworkFunctionSkuListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/NetworkFunctionVendorSkusClient.ListByVendor")
		defer func() {
			sc := -1
			if result.nfslr.Response.Response != nil {
				sc = result.nfslr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("hybridnetwork.NetworkFunctionVendorSkusClient", "ListByVendor", err.Error())
	}

	result.fn = client.listByVendorNextResults
	req, err := client.ListByVendorPreparer(ctx, vendorName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "ListByVendor", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByVendorSender(req)
	if err != nil {
		result.nfslr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "ListByVendor", resp, "Failure sending request")
		return
	}

	result.nfslr, err = client.ListByVendorResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "ListByVendor", resp, "Failure responding to request")
	}
	if result.nfslr.hasNextLink() && result.nfslr.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListByVendorPreparer prepares the ListByVendor request.
func (client NetworkFunctionVendorSkusClient) ListByVendorPreparer(ctx context.Context, vendorName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
		"vendorName":     autorest.Encode("path", vendorName),
	}

	const APIVersion = "2020-01-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.HybridNetwork/networkFunctionVendors/{vendorName}/vendorSkus", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByVendorSender sends the ListByVendor request. The method will close the
// http.Response Body if it receives an error.
func (client NetworkFunctionVendorSkusClient) ListByVendorSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByVendorResponder handles the response to the ListByVendor request. The method always
// closes the http.Response Body.
func (client NetworkFunctionVendorSkusClient) ListByVendorResponder(resp *http.Response) (result NetworkFunctionSkuListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByVendorNextResults retrieves the next set of results, if any.
func (client NetworkFunctionVendorSkusClient) listByVendorNextResults(ctx context.Context, lastResults NetworkFunctionSkuListResult) (result NetworkFunctionSkuListResult, err error) {
	req, err := lastResults.networkFunctionSkuListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "listByVendorNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByVendorSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "listByVendorNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByVendorResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridnetwork.NetworkFunctionVendorSkusClient", "listByVendorNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByVendorComplete enumerates all values, automatically crossing page boundaries as required.
func (client NetworkFunctionVendorSkusClient) ListByVendorComplete(ctx context.Context, vendorName string) (result NetworkFunctionSkuListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/NetworkFunctionVendorSkusClient.ListByVendor")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByVendor(ctx, vendorName)
	return
}
