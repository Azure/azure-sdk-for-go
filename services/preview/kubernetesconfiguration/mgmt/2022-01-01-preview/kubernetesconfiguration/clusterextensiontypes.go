package kubernetesconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
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

// ClusterExtensionTypesClient is the kubernetesConfiguration Client
type ClusterExtensionTypesClient struct {
	BaseClient
}

// NewClusterExtensionTypesClient creates an instance of the ClusterExtensionTypesClient client.
func NewClusterExtensionTypesClient(subscriptionID string) ClusterExtensionTypesClient {
	return NewClusterExtensionTypesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewClusterExtensionTypesClientWithBaseURI creates an instance of the ClusterExtensionTypesClient client using a
// custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds,
// Azure stack).
func NewClusterExtensionTypesClientWithBaseURI(baseURI string, subscriptionID string) ClusterExtensionTypesClient {
	return ClusterExtensionTypesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// List get Extension Types
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// clusterRp - the Kubernetes cluster RP - either Microsoft.ContainerService (for AKS clusters) or
// Microsoft.Kubernetes (for OnPrem K8S clusters).
// clusterResourceName - the Kubernetes cluster resource name - either managedClusters (for AKS clusters) or
// connectedClusters (for OnPrem K8S clusters).
// clusterName - the name of the kubernetes cluster.
func (client ClusterExtensionTypesClient) List(ctx context.Context, resourceGroupName string, clusterRp ExtensionsClusterRp, clusterResourceName ExtensionsClusterResourceName, clusterName string) (result ExtensionTypeListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ClusterExtensionTypesClient.List")
		defer func() {
			sc := -1
			if result.etl.Response.Response != nil {
				sc = result.etl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("kubernetesconfiguration.ClusterExtensionTypesClient", "List", err.Error())
	}

	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName, clusterRp, clusterResourceName, clusterName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kubernetesconfiguration.ClusterExtensionTypesClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.etl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "kubernetesconfiguration.ClusterExtensionTypesClient", "List", resp, "Failure sending request")
		return
	}

	result.etl, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kubernetesconfiguration.ClusterExtensionTypesClient", "List", resp, "Failure responding to request")
		return
	}
	if result.etl.hasNextLink() && result.etl.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client ClusterExtensionTypesClient) ListPreparer(ctx context.Context, resourceGroupName string, clusterRp ExtensionsClusterRp, clusterResourceName ExtensionsClusterResourceName, clusterName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"clusterName":         autorest.Encode("path", clusterName),
		"clusterResourceName": autorest.Encode("path", clusterResourceName),
		"clusterRp":           autorest.Encode("path", clusterRp),
		"resourceGroupName":   autorest.Encode("path", resourceGroupName),
		"subscriptionId":      autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2022-01-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{clusterRp}/{clusterResourceName}/{clusterName}/providers/Microsoft.KubernetesConfiguration/extensionTypes", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client ClusterExtensionTypesClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ClusterExtensionTypesClient) ListResponder(resp *http.Response) (result ExtensionTypeList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client ClusterExtensionTypesClient) listNextResults(ctx context.Context, lastResults ExtensionTypeList) (result ExtensionTypeList, err error) {
	req, err := lastResults.extensionTypeListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "kubernetesconfiguration.ClusterExtensionTypesClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "kubernetesconfiguration.ClusterExtensionTypesClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kubernetesconfiguration.ClusterExtensionTypesClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client ClusterExtensionTypesClient) ListComplete(ctx context.Context, resourceGroupName string, clusterRp ExtensionsClusterRp, clusterResourceName ExtensionsClusterResourceName, clusterName string) (result ExtensionTypeListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ClusterExtensionTypesClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, resourceGroupName, clusterRp, clusterResourceName, clusterName)
	return
}
