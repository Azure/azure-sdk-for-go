package compute

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
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// CloudServiceOperatingSystemsClient is the compute Client
type CloudServiceOperatingSystemsClient struct {
	BaseClient
}

// NewCloudServiceOperatingSystemsClient creates an instance of the CloudServiceOperatingSystemsClient client.
func NewCloudServiceOperatingSystemsClient(subscriptionID string) CloudServiceOperatingSystemsClient {
	return NewCloudServiceOperatingSystemsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewCloudServiceOperatingSystemsClientWithBaseURI creates an instance of the CloudServiceOperatingSystemsClient
// client using a custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI
// (sovereign clouds, Azure stack).
func NewCloudServiceOperatingSystemsClientWithBaseURI(baseURI string, subscriptionID string) CloudServiceOperatingSystemsClient {
	return CloudServiceOperatingSystemsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// GetOSFamily gets properties of a guest operating system family that can be specified in the XML service
// configuration (.cscfg) for a cloud service.
// Parameters:
// location - name of the location that the OS family pertains to.
// osFamilyName - name of the OS family.
func (client CloudServiceOperatingSystemsClient) GetOSFamily(ctx context.Context, location string, osFamilyName string) (result OSFamily, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CloudServiceOperatingSystemsClient.GetOSFamily")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetOSFamilyPreparer(ctx, location, osFamilyName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "GetOSFamily", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetOSFamilySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "GetOSFamily", resp, "Failure sending request")
		return
	}

	result, err = client.GetOSFamilyResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "GetOSFamily", resp, "Failure responding to request")
		return
	}

	return
}

// GetOSFamilyPreparer prepares the GetOSFamily request.
func (client CloudServiceOperatingSystemsClient) GetOSFamilyPreparer(ctx context.Context, location string, osFamilyName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       autorest.Encode("path", location),
		"osFamilyName":   autorest.Encode("path", osFamilyName),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/cloudServiceOsFamilies/{osFamilyName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetOSFamilySender sends the GetOSFamily request. The method will close the
// http.Response Body if it receives an error.
func (client CloudServiceOperatingSystemsClient) GetOSFamilySender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetOSFamilyResponder handles the response to the GetOSFamily request. The method always
// closes the http.Response Body.
func (client CloudServiceOperatingSystemsClient) GetOSFamilyResponder(resp *http.Response) (result OSFamily, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetOSVersion gets properties of a guest operating system version that can be specified in the XML service
// configuration (.cscfg) for a cloud service.
// Parameters:
// location - name of the location that the OS version pertains to.
// osVersionName - name of the OS version.
func (client CloudServiceOperatingSystemsClient) GetOSVersion(ctx context.Context, location string, osVersionName string) (result OSVersion, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CloudServiceOperatingSystemsClient.GetOSVersion")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetOSVersionPreparer(ctx, location, osVersionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "GetOSVersion", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetOSVersionSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "GetOSVersion", resp, "Failure sending request")
		return
	}

	result, err = client.GetOSVersionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "GetOSVersion", resp, "Failure responding to request")
		return
	}

	return
}

// GetOSVersionPreparer prepares the GetOSVersion request.
func (client CloudServiceOperatingSystemsClient) GetOSVersionPreparer(ctx context.Context, location string, osVersionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       autorest.Encode("path", location),
		"osVersionName":  autorest.Encode("path", osVersionName),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/cloudServiceOsVersions/{osVersionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetOSVersionSender sends the GetOSVersion request. The method will close the
// http.Response Body if it receives an error.
func (client CloudServiceOperatingSystemsClient) GetOSVersionSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetOSVersionResponder handles the response to the GetOSVersion request. The method always
// closes the http.Response Body.
func (client CloudServiceOperatingSystemsClient) GetOSVersionResponder(resp *http.Response) (result OSVersion, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListOSFamilies gets a list of all guest operating system families available to be specified in the XML service
// configuration (.cscfg) for a cloud service. Use nextLink property in the response to get the next page of OS
// Families. Do this till nextLink is null to fetch all the OS Families.
// Parameters:
// location - name of the location that the OS families pertain to.
func (client CloudServiceOperatingSystemsClient) ListOSFamilies(ctx context.Context, location string) (result OSFamilyListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CloudServiceOperatingSystemsClient.ListOSFamilies")
		defer func() {
			sc := -1
			if result.oflr.Response.Response != nil {
				sc = result.oflr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listOSFamiliesNextResults
	req, err := client.ListOSFamiliesPreparer(ctx, location)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "ListOSFamilies", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListOSFamiliesSender(req)
	if err != nil {
		result.oflr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "ListOSFamilies", resp, "Failure sending request")
		return
	}

	result.oflr, err = client.ListOSFamiliesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "ListOSFamilies", resp, "Failure responding to request")
		return
	}
	if result.oflr.hasNextLink() && result.oflr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListOSFamiliesPreparer prepares the ListOSFamilies request.
func (client CloudServiceOperatingSystemsClient) ListOSFamiliesPreparer(ctx context.Context, location string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       autorest.Encode("path", location),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/cloudServiceOsFamilies", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListOSFamiliesSender sends the ListOSFamilies request. The method will close the
// http.Response Body if it receives an error.
func (client CloudServiceOperatingSystemsClient) ListOSFamiliesSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListOSFamiliesResponder handles the response to the ListOSFamilies request. The method always
// closes the http.Response Body.
func (client CloudServiceOperatingSystemsClient) ListOSFamiliesResponder(resp *http.Response) (result OSFamilyListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listOSFamiliesNextResults retrieves the next set of results, if any.
func (client CloudServiceOperatingSystemsClient) listOSFamiliesNextResults(ctx context.Context, lastResults OSFamilyListResult) (result OSFamilyListResult, err error) {
	req, err := lastResults.oSFamilyListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "listOSFamiliesNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListOSFamiliesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "listOSFamiliesNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListOSFamiliesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "listOSFamiliesNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListOSFamiliesComplete enumerates all values, automatically crossing page boundaries as required.
func (client CloudServiceOperatingSystemsClient) ListOSFamiliesComplete(ctx context.Context, location string) (result OSFamilyListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CloudServiceOperatingSystemsClient.ListOSFamilies")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListOSFamilies(ctx, location)
	return
}

// ListOSVersions gets a list of all guest operating system versions available to be specified in the XML service
// configuration (.cscfg) for a cloud service. Use nextLink property in the response to get the next page of OS
// versions. Do this till nextLink is null to fetch all the OS versions.
// Parameters:
// location - name of the location that the OS versions pertain to.
func (client CloudServiceOperatingSystemsClient) ListOSVersions(ctx context.Context, location string) (result OSVersionListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CloudServiceOperatingSystemsClient.ListOSVersions")
		defer func() {
			sc := -1
			if result.ovlr.Response.Response != nil {
				sc = result.ovlr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listOSVersionsNextResults
	req, err := client.ListOSVersionsPreparer(ctx, location)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "ListOSVersions", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListOSVersionsSender(req)
	if err != nil {
		result.ovlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "ListOSVersions", resp, "Failure sending request")
		return
	}

	result.ovlr, err = client.ListOSVersionsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "ListOSVersions", resp, "Failure responding to request")
		return
	}
	if result.ovlr.hasNextLink() && result.ovlr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListOSVersionsPreparer prepares the ListOSVersions request.
func (client CloudServiceOperatingSystemsClient) ListOSVersionsPreparer(ctx context.Context, location string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       autorest.Encode("path", location),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/cloudServiceOsVersions", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListOSVersionsSender sends the ListOSVersions request. The method will close the
// http.Response Body if it receives an error.
func (client CloudServiceOperatingSystemsClient) ListOSVersionsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListOSVersionsResponder handles the response to the ListOSVersions request. The method always
// closes the http.Response Body.
func (client CloudServiceOperatingSystemsClient) ListOSVersionsResponder(resp *http.Response) (result OSVersionListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listOSVersionsNextResults retrieves the next set of results, if any.
func (client CloudServiceOperatingSystemsClient) listOSVersionsNextResults(ctx context.Context, lastResults OSVersionListResult) (result OSVersionListResult, err error) {
	req, err := lastResults.oSVersionListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "listOSVersionsNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListOSVersionsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "listOSVersionsNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListOSVersionsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.CloudServiceOperatingSystemsClient", "listOSVersionsNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListOSVersionsComplete enumerates all values, automatically crossing page boundaries as required.
func (client CloudServiceOperatingSystemsClient) ListOSVersionsComplete(ctx context.Context, location string) (result OSVersionListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CloudServiceOperatingSystemsClient.ListOSVersions")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListOSVersions(ctx, location)
	return
}
