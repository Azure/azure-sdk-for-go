package kusto

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

// AttachedDatabaseConfigurationsClient is the the Azure Kusto management API provides a RESTful set of web services
// that interact with Azure Kusto services to manage your clusters and databases. The API enables you to create,
// update, and delete clusters and databases.
type AttachedDatabaseConfigurationsClient struct {
	BaseClient
}

// NewAttachedDatabaseConfigurationsClient creates an instance of the AttachedDatabaseConfigurationsClient client.
func NewAttachedDatabaseConfigurationsClient(subscriptionID string) AttachedDatabaseConfigurationsClient {
	return NewAttachedDatabaseConfigurationsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewAttachedDatabaseConfigurationsClientWithBaseURI creates an instance of the AttachedDatabaseConfigurationsClient
// client using a custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI
// (sovereign clouds, Azure stack).
func NewAttachedDatabaseConfigurationsClientWithBaseURI(baseURI string, subscriptionID string) AttachedDatabaseConfigurationsClient {
	return AttachedDatabaseConfigurationsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate creates or updates an attached database configuration.
// Parameters:
// resourceGroupName - the name of the resource group containing the Kusto cluster.
// clusterName - the name of the Kusto cluster.
// attachedDatabaseConfigurationName - the name of the attached database configuration.
// parameters - the database parameters supplied to the CreateOrUpdate operation.
func (client AttachedDatabaseConfigurationsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, attachedDatabaseConfigurationName string, parameters AttachedDatabaseConfiguration) (result AttachedDatabaseConfigurationsCreateOrUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/AttachedDatabaseConfigurationsClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.AttachedDatabaseConfigurationProperties", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "parameters.AttachedDatabaseConfigurationProperties.DatabaseName", Name: validation.Null, Rule: true, Chain: nil},
					{Target: "parameters.AttachedDatabaseConfigurationProperties.ClusterResourceID", Name: validation.Null, Rule: true, Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewError("kusto.AttachedDatabaseConfigurationsClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, clusterName, attachedDatabaseConfigurationName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kusto.AttachedDatabaseConfigurationsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kusto.AttachedDatabaseConfigurationsClient", "CreateOrUpdate", nil, "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client AttachedDatabaseConfigurationsClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, clusterName string, attachedDatabaseConfigurationName string, parameters AttachedDatabaseConfiguration) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"attachedDatabaseConfigurationName": autorest.Encode("path", attachedDatabaseConfigurationName),
		"clusterName":                       autorest.Encode("path", clusterName),
		"resourceGroupName":                 autorest.Encode("path", resourceGroupName),
		"subscriptionId":                    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-02-15"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/attachedDatabaseConfigurations/{attachedDatabaseConfigurationName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client AttachedDatabaseConfigurationsClient) CreateOrUpdateSender(req *http.Request) (future AttachedDatabaseConfigurationsCreateOrUpdateFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client AttachedDatabaseConfigurationsClient) CreateOrUpdateResponder(resp *http.Response) (result AttachedDatabaseConfiguration, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes the attached database configuration with the given name.
// Parameters:
// resourceGroupName - the name of the resource group containing the Kusto cluster.
// clusterName - the name of the Kusto cluster.
// attachedDatabaseConfigurationName - the name of the attached database configuration.
func (client AttachedDatabaseConfigurationsClient) Delete(ctx context.Context, resourceGroupName string, clusterName string, attachedDatabaseConfigurationName string) (result AttachedDatabaseConfigurationsDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/AttachedDatabaseConfigurationsClient.Delete")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, resourceGroupName, clusterName, attachedDatabaseConfigurationName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kusto.AttachedDatabaseConfigurationsClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kusto.AttachedDatabaseConfigurationsClient", "Delete", nil, "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client AttachedDatabaseConfigurationsClient) DeletePreparer(ctx context.Context, resourceGroupName string, clusterName string, attachedDatabaseConfigurationName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"attachedDatabaseConfigurationName": autorest.Encode("path", attachedDatabaseConfigurationName),
		"clusterName":                       autorest.Encode("path", clusterName),
		"resourceGroupName":                 autorest.Encode("path", resourceGroupName),
		"subscriptionId":                    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-02-15"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/attachedDatabaseConfigurations/{attachedDatabaseConfigurationName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client AttachedDatabaseConfigurationsClient) DeleteSender(req *http.Request) (future AttachedDatabaseConfigurationsDeleteFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client AttachedDatabaseConfigurationsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get returns an attached database configuration.
// Parameters:
// resourceGroupName - the name of the resource group containing the Kusto cluster.
// clusterName - the name of the Kusto cluster.
// attachedDatabaseConfigurationName - the name of the attached database configuration.
func (client AttachedDatabaseConfigurationsClient) Get(ctx context.Context, resourceGroupName string, clusterName string, attachedDatabaseConfigurationName string) (result AttachedDatabaseConfiguration, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/AttachedDatabaseConfigurationsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, resourceGroupName, clusterName, attachedDatabaseConfigurationName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kusto.AttachedDatabaseConfigurationsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "kusto.AttachedDatabaseConfigurationsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kusto.AttachedDatabaseConfigurationsClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client AttachedDatabaseConfigurationsClient) GetPreparer(ctx context.Context, resourceGroupName string, clusterName string, attachedDatabaseConfigurationName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"attachedDatabaseConfigurationName": autorest.Encode("path", attachedDatabaseConfigurationName),
		"clusterName":                       autorest.Encode("path", clusterName),
		"resourceGroupName":                 autorest.Encode("path", resourceGroupName),
		"subscriptionId":                    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-02-15"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/attachedDatabaseConfigurations/{attachedDatabaseConfigurationName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client AttachedDatabaseConfigurationsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client AttachedDatabaseConfigurationsClient) GetResponder(resp *http.Response) (result AttachedDatabaseConfiguration, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByCluster returns the list of attached database configurations of the given Kusto cluster.
// Parameters:
// resourceGroupName - the name of the resource group containing the Kusto cluster.
// clusterName - the name of the Kusto cluster.
func (client AttachedDatabaseConfigurationsClient) ListByCluster(ctx context.Context, resourceGroupName string, clusterName string) (result AttachedDatabaseConfigurationListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/AttachedDatabaseConfigurationsClient.ListByCluster")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListByClusterPreparer(ctx, resourceGroupName, clusterName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kusto.AttachedDatabaseConfigurationsClient", "ListByCluster", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByClusterSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "kusto.AttachedDatabaseConfigurationsClient", "ListByCluster", resp, "Failure sending request")
		return
	}

	result, err = client.ListByClusterResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kusto.AttachedDatabaseConfigurationsClient", "ListByCluster", resp, "Failure responding to request")
		return
	}

	return
}

// ListByClusterPreparer prepares the ListByCluster request.
func (client AttachedDatabaseConfigurationsClient) ListByClusterPreparer(ctx context.Context, resourceGroupName string, clusterName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"clusterName":       autorest.Encode("path", clusterName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-02-15"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/attachedDatabaseConfigurations", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByClusterSender sends the ListByCluster request. The method will close the
// http.Response Body if it receives an error.
func (client AttachedDatabaseConfigurationsClient) ListByClusterSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByClusterResponder handles the response to the ListByCluster request. The method always
// closes the http.Response Body.
func (client AttachedDatabaseConfigurationsClient) ListByClusterResponder(resp *http.Response) (result AttachedDatabaseConfigurationListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
