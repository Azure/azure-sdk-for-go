package databoxedge

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

// NodesClient is the client for the Nodes methods of the Databoxedge service.
type NodesClient struct {
	BaseClient
}

// NewNodesClient creates an instance of the NodesClient client.
func NewNodesClient(subscriptionID string) NodesClient {
	return NewNodesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewNodesClientWithBaseURI creates an instance of the NodesClient client using a custom endpoint.  Use this when
// interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewNodesClientWithBaseURI(baseURI string, subscriptionID string) NodesClient {
	return NodesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// ListByDataBoxEdgeDevice gets all the nodes currently configured under this Data Box Edge device
// Parameters:
// deviceName - the device name.
// resourceGroupName - the resource group name.
func (client NodesClient) ListByDataBoxEdgeDevice(ctx context.Context, deviceName string, resourceGroupName string) (result NodeList, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/NodesClient.ListByDataBoxEdgeDevice")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListByDataBoxEdgeDevicePreparer(ctx, deviceName, resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databoxedge.NodesClient", "ListByDataBoxEdgeDevice", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByDataBoxEdgeDeviceSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "databoxedge.NodesClient", "ListByDataBoxEdgeDevice", resp, "Failure sending request")
		return
	}

	result, err = client.ListByDataBoxEdgeDeviceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databoxedge.NodesClient", "ListByDataBoxEdgeDevice", resp, "Failure responding to request")
	}

	return
}

// ListByDataBoxEdgeDevicePreparer prepares the ListByDataBoxEdgeDevice request.
func (client NodesClient) ListByDataBoxEdgeDevicePreparer(ctx context.Context, deviceName string, resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"deviceName":        autorest.Encode("path", deviceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-05-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/{deviceName}/nodes", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByDataBoxEdgeDeviceSender sends the ListByDataBoxEdgeDevice request. The method will close the
// http.Response Body if it receives an error.
func (client NodesClient) ListByDataBoxEdgeDeviceSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByDataBoxEdgeDeviceResponder handles the response to the ListByDataBoxEdgeDevice request. The method always
// closes the http.Response Body.
func (client NodesClient) ListByDataBoxEdgeDeviceResponder(resp *http.Response) (result NodeList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
