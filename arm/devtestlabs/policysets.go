package devtestlabs

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
// Code generated by Microsoft (R) AutoRest Code Generator 2.2.21.0
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// PolicySetsClient is the the DevTest Labs Client.
type PolicySetsClient struct {
	ManagementClient
}

// NewPolicySetsClient creates an instance of the PolicySetsClient client.
func NewPolicySetsClient(subscriptionID string) PolicySetsClient {
	return NewPolicySetsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewPolicySetsClientWithBaseURI creates an instance of the PolicySetsClient client.
func NewPolicySetsClientWithBaseURI(baseURI string, subscriptionID string) PolicySetsClient {
	return PolicySetsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// EvaluatePolicies evaluates lab policy.
//
// resourceGroupName is the name of the resource group. labName is the name of the lab. name is the name of the policy
// set. evaluatePoliciesRequest is request body for evaluating a policy set.
func (client PolicySetsClient) EvaluatePolicies(resourceGroupName string, labName string, name string, evaluatePoliciesRequest EvaluatePoliciesRequest) (result EvaluatePoliciesResponse, err error) {
	req, err := client.EvaluatePoliciesPreparer(resourceGroupName, labName, name, evaluatePoliciesRequest)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devtestlabs.PolicySetsClient", "EvaluatePolicies", nil, "Failure preparing request")
		return
	}

	resp, err := client.EvaluatePoliciesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "devtestlabs.PolicySetsClient", "EvaluatePolicies", resp, "Failure sending request")
		return
	}

	result, err = client.EvaluatePoliciesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devtestlabs.PolicySetsClient", "EvaluatePolicies", resp, "Failure responding to request")
	}

	return
}

// EvaluatePoliciesPreparer prepares the EvaluatePolicies request.
func (client PolicySetsClient) EvaluatePoliciesPreparer(resourceGroupName string, labName string, name string, evaluatePoliciesRequest EvaluatePoliciesRequest) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"labName":           autorest.Encode("path", labName),
		"name":              autorest.Encode("path", name),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-05-15"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/policysets/{name}/evaluatePolicies", pathParameters),
		autorest.WithJSON(evaluatePoliciesRequest),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// EvaluatePoliciesSender sends the EvaluatePolicies request. The method will close the
// http.Response Body if it receives an error.
func (client PolicySetsClient) EvaluatePoliciesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// EvaluatePoliciesResponder handles the response to the EvaluatePolicies request. The method always
// closes the http.Response Body.
func (client PolicySetsClient) EvaluatePoliciesResponder(resp *http.Response) (result EvaluatePoliciesResponse, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
