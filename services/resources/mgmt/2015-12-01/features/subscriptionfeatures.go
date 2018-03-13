package features

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
	"net/http"
)

// SubscriptionFeaturesClient is the azure Feature Exposure Control (AFEC) provides a mechanism for the resource
// providers to control feature exposure to users. Resource providers typically use this mechanism to provide
// public/private preview for new features prior to making them generally available. Users need to explicitly register
// for AFEC features to get access to such functionality.
type SubscriptionFeaturesClient struct {
	BaseClient
}

// NewSubscriptionFeaturesClient creates an instance of the SubscriptionFeaturesClient client.
func NewSubscriptionFeaturesClient(subscriptionID string) SubscriptionFeaturesClient {
	return NewSubscriptionFeaturesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewSubscriptionFeaturesClientWithBaseURI creates an instance of the SubscriptionFeaturesClient client.
func NewSubscriptionFeaturesClientWithBaseURI(baseURI string, subscriptionID string) SubscriptionFeaturesClient {
	return SubscriptionFeaturesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Register registers the preview feature for the subscription.
//
// resourceProviderNamespace is the namespace of the resource provider. featureName is the name of the feature to
// register.
func (client SubscriptionFeaturesClient) Register(ctx context.Context, resourceProviderNamespace string, featureName string) (result Result, err error) {
	req, err := client.RegisterPreparer(ctx, resourceProviderNamespace, featureName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.SubscriptionFeaturesClient", "Register", nil, "Failure preparing request")
		return
	}

	resp, err := client.RegisterSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "features.SubscriptionFeaturesClient", "Register", resp, "Failure sending request")
		return
	}

	result, err = client.RegisterResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.SubscriptionFeaturesClient", "Register", resp, "Failure responding to request")
	}

	return
}

// RegisterPreparer prepares the Register request.
func (client SubscriptionFeaturesClient) RegisterPreparer(ctx context.Context, resourceProviderNamespace string, featureName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"featureName":               autorest.Encode("path", featureName),
		"resourceProviderNamespace": autorest.Encode("path", resourceProviderNamespace),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Features/providers/{resourceProviderNamespace}/features/{featureName}/register", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// RegisterSender sends the Register request. The method will close the
// http.Response Body if it receives an error.
func (client SubscriptionFeaturesClient) RegisterSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// RegisterResponder handles the response to the Register request. The method always
// closes the http.Response Body.
func (client SubscriptionFeaturesClient) RegisterResponder(resp *http.Response) (result Result, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
