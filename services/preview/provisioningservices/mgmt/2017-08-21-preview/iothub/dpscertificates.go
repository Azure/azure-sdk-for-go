package iothub

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

// DpsCertificatesClient is the API for using the Azure IoT Hub Device Provisioning Service features.
type DpsCertificatesClient struct {
	BaseClient
}

// NewDpsCertificatesClient creates an instance of the DpsCertificatesClient client.
func NewDpsCertificatesClient(subscriptionID string) DpsCertificatesClient {
	return NewDpsCertificatesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewDpsCertificatesClientWithBaseURI creates an instance of the DpsCertificatesClient client.
func NewDpsCertificatesClientWithBaseURI(baseURI string, subscriptionID string) DpsCertificatesClient {
	return DpsCertificatesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// List get all the certificates tied to the provisioning service.
// Parameters:
// resourceGroupName - name of resource group.
// provisioningServiceName - name of provisioning service to retrieve certificates for.
func (client DpsCertificatesClient) List(ctx context.Context, resourceGroupName string, provisioningServiceName string) (result CertificateListDescription, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DpsCertificatesClient.List")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListPreparer(ctx, resourceGroupName, provisioningServiceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iothub.DpsCertificatesClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "iothub.DpsCertificatesClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iothub.DpsCertificatesClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client DpsCertificatesClient) ListPreparer(ctx context.Context, resourceGroupName string, provisioningServiceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"provisioningServiceName": autorest.Encode("path", provisioningServiceName),
		"resourceGroupName":       autorest.Encode("path", resourceGroupName),
		"subscriptionId":          autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-08-21-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Devices/provisioningServices/{provisioningServiceName}/certificates", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client DpsCertificatesClient) ListSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client DpsCertificatesClient) ListResponder(resp *http.Response) (result CertificateListDescription, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
