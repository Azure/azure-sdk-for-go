// Package reservations implements the Azure ARM Reservations service API version 2019-04-01.
//
// This API describe Azure Reservation
package reservations

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

const (
	// DefaultBaseURI is the default URI used for the service Reservations
	DefaultBaseURI = "https://management.azure.com"
)

// BaseClient is the base client for Reservations.
type BaseClient struct {
	autorest.Client
	BaseURI string
}

// New creates an instance of the BaseClient client.
func New() BaseClient {
	return NewWithBaseURI(DefaultBaseURI)
}

// NewWithBaseURI creates an instance of the BaseClient client.
func NewWithBaseURI(baseURI string) BaseClient {
	return BaseClient{
		Client:  autorest.NewClientWithUserAgent(UserAgent()),
		BaseURI: baseURI,
	}
}

// GetAppliedReservationList get applicable `Reservation`s that are applied to this subscription.
// Parameters:
// subscriptionID - id of the subscription
func (client BaseClient) GetAppliedReservationList(ctx context.Context, subscriptionID string) (result AppliedReservations, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BaseClient.GetAppliedReservationList")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetAppliedReservationListPreparer(ctx, subscriptionID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "reservations.BaseClient", "GetAppliedReservationList", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetAppliedReservationListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "reservations.BaseClient", "GetAppliedReservationList", resp, "Failure sending request")
		return
	}

	result, err = client.GetAppliedReservationListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "reservations.BaseClient", "GetAppliedReservationList", resp, "Failure responding to request")
	}

	return
}

// GetAppliedReservationListPreparer prepares the GetAppliedReservationList request.
func (client BaseClient) GetAppliedReservationListPreparer(ctx context.Context, subscriptionID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", subscriptionID),
	}

	const APIVersion = "2019-04-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Capacity/appliedReservations", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetAppliedReservationListSender sends the GetAppliedReservationList request. The method will close the
// http.Response Body if it receives an error.
func (client BaseClient) GetAppliedReservationListSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// GetAppliedReservationListResponder handles the response to the GetAppliedReservationList request. The method always
// closes the http.Response Body.
func (client BaseClient) GetAppliedReservationListResponder(resp *http.Response) (result AppliedReservations, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetCatalog sends the get catalog request.
// Parameters:
// subscriptionID - id of the subscription
// reservedResourceType - the type of the resource for which the skus should be provided.
// location - filters the skus based on the location specified in this parameter. This can be an azure region
// or global
func (client BaseClient) GetCatalog(ctx context.Context, subscriptionID string, reservedResourceType string, location string) (result ListCatalog, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BaseClient.GetCatalog")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetCatalogPreparer(ctx, subscriptionID, reservedResourceType, location)
	if err != nil {
		err = autorest.NewErrorWithError(err, "reservations.BaseClient", "GetCatalog", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetCatalogSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "reservations.BaseClient", "GetCatalog", resp, "Failure sending request")
		return
	}

	result, err = client.GetCatalogResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "reservations.BaseClient", "GetCatalog", resp, "Failure responding to request")
	}

	return
}

// GetCatalogPreparer prepares the GetCatalog request.
func (client BaseClient) GetCatalogPreparer(ctx context.Context, subscriptionID string, reservedResourceType string, location string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", subscriptionID),
	}

	const APIVersion = "2019-04-01"
	queryParameters := map[string]interface{}{
		"api-version":          APIVersion,
		"reservedResourceType": autorest.Encode("query", reservedResourceType),
	}
	if len(location) > 0 {
		queryParameters["location"] = autorest.Encode("query", location)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Capacity/catalogs", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetCatalogSender sends the GetCatalog request. The method will close the
// http.Response Body if it receives an error.
func (client BaseClient) GetCatalogSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// GetCatalogResponder handles the response to the GetCatalog request. The method always
// closes the http.Response Body.
func (client BaseClient) GetCatalogResponder(resp *http.Response) (result ListCatalog, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
