package vmwarecloudsimple

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

// VirtualMachineTemplateByPCClient is the description of the new service
type VirtualMachineTemplateByPCClient struct {
	BaseClient
}

// NewVirtualMachineTemplateByPCClient creates an instance of the VirtualMachineTemplateByPCClient client.
func NewVirtualMachineTemplateByPCClient(referer string, regionID string, subscriptionID string) VirtualMachineTemplateByPCClient {
	return NewVirtualMachineTemplateByPCClientWithBaseURI(DefaultBaseURI, referer, regionID, subscriptionID)
}

// NewVirtualMachineTemplateByPCClientWithBaseURI creates an instance of the VirtualMachineTemplateByPCClient client.
func NewVirtualMachineTemplateByPCClientWithBaseURI(baseURI string, referer string, regionID string, subscriptionID string) VirtualMachineTemplateByPCClient {
	return VirtualMachineTemplateByPCClient{NewWithBaseURI(baseURI, referer, regionID, subscriptionID)}
}

// Get returns virtual machine templates by its name
// Parameters:
// pcName - the private cloud name
// virtualMachineTemplateName - virtual machine template id (vsphereId)
func (client VirtualMachineTemplateByPCClient) Get(ctx context.Context, pcName string, virtualMachineTemplateName string) (result VirtualMachineTemplate, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualMachineTemplateByPCClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, pcName, virtualMachineTemplateName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vmwarecloudsimple.VirtualMachineTemplateByPCClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "vmwarecloudsimple.VirtualMachineTemplateByPCClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vmwarecloudsimple.VirtualMachineTemplateByPCClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client VirtualMachineTemplateByPCClient) GetPreparer(ctx context.Context, pcName string, virtualMachineTemplateName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"pcName":                     autorest.Encode("path", pcName),
		"regionId":                   autorest.Encode("path", client.RegionID),
		"subscriptionId":             autorest.Encode("path", client.SubscriptionID),
		"virtualMachineTemplateName": autorest.Encode("path", virtualMachineTemplateName),
	}

	const APIVersion = "2019-04-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.VMwareCloudSimple/locations/{regionId}/privateClouds/{pcName}/virtualMachineTemplates/{virtualMachineTemplateName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineTemplateByPCClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client VirtualMachineTemplateByPCClient) GetResponder(resp *http.Response) (result VirtualMachineTemplate, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
