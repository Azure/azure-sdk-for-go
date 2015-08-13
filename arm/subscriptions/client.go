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
// Code generated by Microsoft (R) AutoRest Code Generator 0.11.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

package subscriptions

import (
	"fmt"
	"github.com/azure/go-autorest/autorest"
	"net/http"
	"net/url"
	"time"
)

const (
	ApiVersion             = "2014-04-01-preview"
	DefaultBaseUri         = "https://management.azure.com"
	DefaultPollingDuration = 10 * time.Minute
)

type SubscriptionClient struct {
	autorest.Client
	BaseUri        string
	SubscriptionId string
}

func New(subscriptionId string) *SubscriptionClient {
	return NewWithBaseUri(DefaultBaseUri, subscriptionId)
}

func NewWithBaseUri(baseUri string, subscriptionId string) *SubscriptionClient {
	client := &SubscriptionClient{BaseUri: baseUri, SubscriptionId: subscriptionId}
	client.PollingMode = autorest.PollUntilDuration
	client.PollingDuration = DefaultPollingDuration
	return client
}

////////////////////////////////////////////////////////////////////////////////
//
// Subscriptions Client
//
////////////////////////////////////////////////////////////////////////////////
type SubscriptionsClient struct {
	SubscriptionClient
}

func NewSubscriptionsClient(subscriptionId string) *SubscriptionsClient {
	return NewSubscriptionsClientWithBaseUri(DefaultBaseUri, subscriptionId)
}

func NewSubscriptionsClientWithBaseUri(baseUri string, subscriptionId string) *SubscriptionsClient {
	return &SubscriptionsClient{*NewWithBaseUri(baseUri, subscriptionId)}
}

// Get gets details about particular subscription.;
//
// subscriptionId is id of the subscription.
func (client *SubscriptionsClient) Get(subscriptionId string) (result Subscription, err error) {

	req, err := client.NewGetRequest(subscriptionId)
	if err != nil {
		return result, fmt.Errorf("subscriptions: Failure preparing SubscriptionsClient.Get request (%v)", err)
	}

	req, err = autorest.Prepare(
		req,
		client.WithAuthorization(),
		client.WithInspection())
	if err != nil {
		return result, fmt.Errorf("subscriptions: Failure sending SubscriptionsClient.Get request (%v)", err)
	}

	resp, err := autorest.SendWithSender(client, req)
	if err != nil {
		return result, fmt.Errorf("subscriptions: Failure sending SubscriptionsClient.Get request (%v)", err)
	}

	result = Subscription{}
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorest.WithErrorUnlessOK(),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())

	return result, err
}

func (client *SubscriptionsClient) NewGetRequest(subscriptionId string) (*http.Request, error) {

	pathParameters := map[string]interface{}{
		"subscriptionId": url.QueryEscape(subscriptionId),
	}

	queryParameters := map[string]interface{}{
		"api-version": ApiVersion,
	}

	return autorest.DecoratePreparer(
		client.GetRequestPreparer(),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters)).Prepare(&http.Request{})
}

func (client *SubscriptionsClient) GetRequestPreparer() autorest.Preparer {
	return autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseUri),
		autorest.WithPath("/subscriptions/{subscriptionId}"))
}

// List gets a list of the subscriptionIds.;
func (client *SubscriptionsClient) List() (result SubscriptionListResult, err error) {

	req, err := client.NewListRequest()
	if err != nil {
		return result, fmt.Errorf("subscriptions: Failure preparing SubscriptionsClient.List request (%v)", err)
	}

	req, err = autorest.Prepare(
		req,
		client.WithAuthorization(),
		client.WithInspection())
	if err != nil {
		return result, fmt.Errorf("subscriptions: Failure sending SubscriptionsClient.List request (%v)", err)
	}

	resp, err := autorest.SendWithSender(client, req)
	if err != nil {
		return result, fmt.Errorf("subscriptions: Failure sending SubscriptionsClient.List request (%v)", err)
	}

	result = SubscriptionListResult{}
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorest.WithErrorUnlessOK(),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())

	return result, err
}

func (client *SubscriptionsClient) NewListRequest() (*http.Request, error) {

	pathParameters := map[string]interface{}{
		"subscriptionId": url.QueryEscape(client.SubscriptionId),
	}

	queryParameters := map[string]interface{}{
		"api-version": ApiVersion,
	}

	return autorest.DecoratePreparer(
		client.ListRequestPreparer(),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters)).Prepare(&http.Request{})
}

func (client *SubscriptionsClient) ListRequestPreparer() autorest.Preparer {
	return autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseUri),
		autorest.WithPath("/subscriptions"))
}

////////////////////////////////////////////////////////////////////////////////
//
// Tenants Client
//
////////////////////////////////////////////////////////////////////////////////
type TenantsClient struct {
	SubscriptionClient
}

func NewTenantsClient(subscriptionId string) *TenantsClient {
	return NewTenantsClientWithBaseUri(DefaultBaseUri, subscriptionId)
}

func NewTenantsClientWithBaseUri(baseUri string, subscriptionId string) *TenantsClient {
	return &TenantsClient{*NewWithBaseUri(baseUri, subscriptionId)}
}

// List gets a list of the tenantIds.;
func (client *TenantsClient) List() (result TenantListResult, err error) {

	req, err := client.NewListRequest()
	if err != nil {
		return result, fmt.Errorf("subscriptions: Failure preparing TenantsClient.List request (%v)", err)
	}

	req, err = autorest.Prepare(
		req,
		client.WithAuthorization(),
		client.WithInspection())
	if err != nil {
		return result, fmt.Errorf("subscriptions: Failure sending TenantsClient.List request (%v)", err)
	}

	resp, err := autorest.SendWithSender(client, req)
	if err != nil {
		return result, fmt.Errorf("subscriptions: Failure sending TenantsClient.List request (%v)", err)
	}

	result = TenantListResult{}
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorest.WithErrorUnlessOK(),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())

	return result, err
}

func (client *TenantsClient) NewListRequest() (*http.Request, error) {

	pathParameters := map[string]interface{}{
		"subscriptionId": url.QueryEscape(client.SubscriptionId),
	}

	queryParameters := map[string]interface{}{
		"api-version": ApiVersion,
	}

	return autorest.DecoratePreparer(
		client.ListRequestPreparer(),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters)).Prepare(&http.Request{})
}

func (client *TenantsClient) ListRequestPreparer() autorest.Preparer {
	return autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseUri),
		autorest.WithPath("/tenants"))
}
