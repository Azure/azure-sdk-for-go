package operationalinsights

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
    "github.com/Azure/go-autorest/autorest"
    "github.com/Azure/go-autorest/autorest/azure"
    "net/http"
    "context"
    "github.com/Azure/go-autorest/tracing"
    "github.com/Azure/go-autorest/autorest/validation"
)

// DataSourcesClient is the operational Insights Client
type DataSourcesClient struct {
    BaseClient
}
// NewDataSourcesClient creates an instance of the DataSourcesClient client.
func NewDataSourcesClient(subscriptionID string) DataSourcesClient {
    return NewDataSourcesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewDataSourcesClientWithBaseURI creates an instance of the DataSourcesClient client.
    func NewDataSourcesClientWithBaseURI(baseURI string, subscriptionID string) DataSourcesClient {
        return DataSourcesClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CreateOrUpdate create or update a data source.
    // Parameters:
        // resourceGroupName - the name of the resource group to get. The name is case insensitive.
        // workspaceName - name of the Log Analytics Workspace that will contain the datasource
        // dataSourceName - the name of the datasource resource.
        // parameters - the parameters required to create or update a datasource.
func (client DataSourcesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, dataSourceName string, parameters DataSource) (result DataSource, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DataSourcesClient.CreateOrUpdate")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}},
            { TargetValue: parameters,
             Constraints: []validation.Constraint{	{Target: "parameters.Properties", Name: validation.Null, Rule: true, Chain: nil }}}}); err != nil {
            return result, validation.NewError("operationalinsights.DataSourcesClient", "CreateOrUpdate", err.Error())
            }

                req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, workspaceName, dataSourceName, parameters)
    if err != nil {
    err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            resp, err := client.CreateOrUpdateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "CreateOrUpdate", resp, "Failure sending request")
            return
            }

            result, err = client.CreateOrUpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "CreateOrUpdate", resp, "Failure responding to request")
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client DataSourcesClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, workspaceName string, dataSourceName string, parameters DataSource) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "dataSourceName": autorest.Encode("path",dataSourceName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "workspaceName": autorest.Encode("path",workspaceName),
            }

                        const APIVersion = "2015-11-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/dataSources/{dataSourceName}",pathParameters),
    autorest.WithJSON(parameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client DataSourcesClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client DataSourcesClient) CreateOrUpdateResponder(resp *http.Response) (result DataSource, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusCreated),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete deletes a data source instance.
    // Parameters:
        // resourceGroupName - the name of the resource group to get. The name is case insensitive.
        // workspaceName - name of the Log Analytics Workspace that contains the datasource.
        // dataSourceName - name of the datasource.
func (client DataSourcesClient) Delete(ctx context.Context, resourceGroupName string, workspaceName string, dataSourceName string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DataSourcesClient.Delete")
        defer func() {
            sc := -1
            if result.Response != nil {
                sc = result.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("operationalinsights.DataSourcesClient", "Delete", err.Error())
            }

                req, err := client.DeletePreparer(ctx, resourceGroupName, workspaceName, dataSourceName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "Delete", nil , "Failure preparing request")
    return
    }

            resp, err := client.DeleteSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "Delete", resp, "Failure sending request")
            return
            }

            result, err = client.DeleteResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "Delete", resp, "Failure responding to request")
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client DataSourcesClient) DeletePreparer(ctx context.Context, resourceGroupName string, workspaceName string, dataSourceName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "dataSourceName": autorest.Encode("path",dataSourceName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "workspaceName": autorest.Encode("path",workspaceName),
            }

                        const APIVersion = "2015-11-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/dataSources/{dataSourceName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client DataSourcesClient) DeleteSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client DataSourcesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Get gets a datasource instance.
    // Parameters:
        // resourceGroupName - the name of the resource group to get. The name is case insensitive.
        // workspaceName - name of the Log Analytics Workspace that contains the datasource.
        // dataSourceName - name of the datasource
func (client DataSourcesClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, dataSourceName string) (result DataSource, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DataSourcesClient.Get")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("operationalinsights.DataSourcesClient", "Get", err.Error())
            }

                req, err := client.GetPreparer(ctx, resourceGroupName, workspaceName, dataSourceName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client DataSourcesClient) GetPreparer(ctx context.Context, resourceGroupName string, workspaceName string, dataSourceName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "dataSourceName": autorest.Encode("path",dataSourceName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "workspaceName": autorest.Encode("path",workspaceName),
            }

                        const APIVersion = "2015-11-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/dataSources/{dataSourceName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client DataSourcesClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client DataSourcesClient) GetResponder(resp *http.Response) (result DataSource, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// ListByWorkspace gets the first page of data source instances in a workspace with the link to the next page.
    // Parameters:
        // resourceGroupName - the name of the resource group to get. The name is case insensitive.
        // workspaceName - the workspace that contains the data sources.
        // filter - the filter to apply on the operation.
        // skiptoken - starting point of the collection of data source instances.
func (client DataSourcesClient) ListByWorkspace(ctx context.Context, resourceGroupName string, workspaceName string, filter string, skiptoken string) (result DataSourceListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DataSourcesClient.ListByWorkspace")
        defer func() {
            sc := -1
            if result.dslr.Response.Response != nil {
                sc = result.dslr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("operationalinsights.DataSourcesClient", "ListByWorkspace", err.Error())
            }

                        result.fn = client.listByWorkspaceNextResults
    req, err := client.ListByWorkspacePreparer(ctx, resourceGroupName, workspaceName, filter, skiptoken)
    if err != nil {
    err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "ListByWorkspace", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByWorkspaceSender(req)
            if err != nil {
            result.dslr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "ListByWorkspace", resp, "Failure sending request")
            return
            }

            result.dslr, err = client.ListByWorkspaceResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "ListByWorkspace", resp, "Failure responding to request")
            }

    return
    }

    // ListByWorkspacePreparer prepares the ListByWorkspace request.
    func (client DataSourcesClient) ListByWorkspacePreparer(ctx context.Context, resourceGroupName string, workspaceName string, filter string, skiptoken string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "workspaceName": autorest.Encode("path",workspaceName),
            }

                        const APIVersion = "2015-11-01-preview"
        queryParameters := map[string]interface{} {
        "$filter": autorest.Encode("query",filter),
        "api-version": APIVersion,
        }
            if len(skiptoken) > 0 {
            queryParameters["$skiptoken"] = autorest.Encode("query",skiptoken)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/dataSources",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByWorkspaceSender sends the ListByWorkspace request. The method will close the
    // http.Response Body if it receives an error.
    func (client DataSourcesClient) ListByWorkspaceSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListByWorkspaceResponder handles the response to the ListByWorkspace request. The method always
// closes the http.Response Body.
func (client DataSourcesClient) ListByWorkspaceResponder(resp *http.Response) (result DataSourceListResult, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

            // listByWorkspaceNextResults retrieves the next set of results, if any.
            func (client DataSourcesClient) listByWorkspaceNextResults(ctx context.Context, lastResults DataSourceListResult) (result DataSourceListResult, err error) {
            req, err := lastResults.dataSourceListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "listByWorkspaceNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListByWorkspaceSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "listByWorkspaceNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListByWorkspaceResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "operationalinsights.DataSourcesClient", "listByWorkspaceNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListByWorkspaceComplete enumerates all values, automatically crossing page boundaries as required.
    func (client DataSourcesClient) ListByWorkspaceComplete(ctx context.Context, resourceGroupName string, workspaceName string, filter string, skiptoken string) (result DataSourceListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/DataSourcesClient.ListByWorkspace")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.ListByWorkspace(ctx, resourceGroupName, workspaceName, filter, skiptoken)
                return
        }

