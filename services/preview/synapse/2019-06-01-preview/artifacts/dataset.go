package artifacts

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

// DatasetClient is the client for the Dataset methods of the Artifacts service.
type DatasetClient struct {
	BaseClient
}

// NewDatasetClient creates an instance of the DatasetClient client.
func NewDatasetClient(endpoint string) DatasetClient {
	return DatasetClient{New(endpoint)}
}

// CreateOrUpdateDataset creates or updates a dataset.
// Parameters:
// datasetName - the dataset name.
// dataset - dataset resource definition.
// ifMatch - eTag of the dataset entity.  Should only be specified for update, for which it should match
// existing entity or can be * for unconditional update.
func (client DatasetClient) CreateOrUpdateDataset(ctx context.Context, datasetName string, dataset DatasetResource, ifMatch string) (result DatasetCreateOrUpdateDatasetFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DatasetClient.CreateOrUpdateDataset")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: datasetName,
			Constraints: []validation.Constraint{{Target: "datasetName", Name: validation.MaxLength, Rule: 260, Chain: nil},
				{Target: "datasetName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "datasetName", Name: validation.Pattern, Rule: `^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("artifacts.DatasetClient", "CreateOrUpdateDataset", err.Error())
	}

	req, err := client.CreateOrUpdateDatasetPreparer(ctx, datasetName, dataset, ifMatch)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "CreateOrUpdateDataset", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateDatasetSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "CreateOrUpdateDataset", nil, "Failure sending request")
		return
	}

	return
}

// CreateOrUpdateDatasetPreparer prepares the CreateOrUpdateDataset request.
func (client DatasetClient) CreateOrUpdateDatasetPreparer(ctx context.Context, datasetName string, dataset DatasetResource, ifMatch string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"datasetName": autorest.Encode("path", datasetName),
	}

	const APIVersion = "2019-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithCustomBaseURL("{endpoint}", urlParameters),
		autorest.WithPathParameters("/datasets/{datasetName}", pathParameters),
		autorest.WithJSON(dataset),
		autorest.WithQueryParameters(queryParameters))
	if len(ifMatch) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("If-Match", autorest.String(ifMatch)))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateDatasetSender sends the CreateOrUpdateDataset request. The method will close the
// http.Response Body if it receives an error.
func (client DatasetClient) CreateOrUpdateDatasetSender(req *http.Request) (future DatasetCreateOrUpdateDatasetFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// CreateOrUpdateDatasetResponder handles the response to the CreateOrUpdateDataset request. The method always
// closes the http.Response Body.
func (client DatasetClient) CreateOrUpdateDatasetResponder(resp *http.Response) (result DatasetResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// DeleteDataset deletes a dataset.
// Parameters:
// datasetName - the dataset name.
func (client DatasetClient) DeleteDataset(ctx context.Context, datasetName string) (result DatasetDeleteDatasetFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DatasetClient.DeleteDataset")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: datasetName,
			Constraints: []validation.Constraint{{Target: "datasetName", Name: validation.MaxLength, Rule: 260, Chain: nil},
				{Target: "datasetName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "datasetName", Name: validation.Pattern, Rule: `^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("artifacts.DatasetClient", "DeleteDataset", err.Error())
	}

	req, err := client.DeleteDatasetPreparer(ctx, datasetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "DeleteDataset", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteDatasetSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "DeleteDataset", nil, "Failure sending request")
		return
	}

	return
}

// DeleteDatasetPreparer prepares the DeleteDataset request.
func (client DatasetClient) DeleteDatasetPreparer(ctx context.Context, datasetName string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"datasetName": autorest.Encode("path", datasetName),
	}

	const APIVersion = "2019-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithCustomBaseURL("{endpoint}", urlParameters),
		autorest.WithPathParameters("/datasets/{datasetName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteDatasetSender sends the DeleteDataset request. The method will close the
// http.Response Body if it receives an error.
func (client DatasetClient) DeleteDatasetSender(req *http.Request) (future DatasetDeleteDatasetFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// DeleteDatasetResponder handles the response to the DeleteDataset request. The method always
// closes the http.Response Body.
func (client DatasetClient) DeleteDatasetResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// GetDataset gets a dataset.
// Parameters:
// datasetName - the dataset name.
// ifNoneMatch - eTag of the dataset entity. Should only be specified for get. If the ETag matches the existing
// entity tag, or if * was provided, then no content will be returned.
func (client DatasetClient) GetDataset(ctx context.Context, datasetName string, ifNoneMatch string) (result DatasetResource, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DatasetClient.GetDataset")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: datasetName,
			Constraints: []validation.Constraint{{Target: "datasetName", Name: validation.MaxLength, Rule: 260, Chain: nil},
				{Target: "datasetName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "datasetName", Name: validation.Pattern, Rule: `^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("artifacts.DatasetClient", "GetDataset", err.Error())
	}

	req, err := client.GetDatasetPreparer(ctx, datasetName, ifNoneMatch)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "GetDataset", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetDatasetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "GetDataset", resp, "Failure sending request")
		return
	}

	result, err = client.GetDatasetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "GetDataset", resp, "Failure responding to request")
		return
	}

	return
}

// GetDatasetPreparer prepares the GetDataset request.
func (client DatasetClient) GetDatasetPreparer(ctx context.Context, datasetName string, ifNoneMatch string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"datasetName": autorest.Encode("path", datasetName),
	}

	const APIVersion = "2019-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{endpoint}", urlParameters),
		autorest.WithPathParameters("/datasets/{datasetName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if len(ifNoneMatch) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("If-None-Match", autorest.String(ifNoneMatch)))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetDatasetSender sends the GetDataset request. The method will close the
// http.Response Body if it receives an error.
func (client DatasetClient) GetDatasetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetDatasetResponder handles the response to the GetDataset request. The method always
// closes the http.Response Body.
func (client DatasetClient) GetDatasetResponder(resp *http.Response) (result DatasetResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNotModified),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetDatasetsByWorkspace lists datasets.
func (client DatasetClient) GetDatasetsByWorkspace(ctx context.Context) (result DatasetListResponsePage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DatasetClient.GetDatasetsByWorkspace")
		defer func() {
			sc := -1
			if result.dlr.Response.Response != nil {
				sc = result.dlr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.getDatasetsByWorkspaceNextResults
	req, err := client.GetDatasetsByWorkspacePreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "GetDatasetsByWorkspace", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetDatasetsByWorkspaceSender(req)
	if err != nil {
		result.dlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "GetDatasetsByWorkspace", resp, "Failure sending request")
		return
	}

	result.dlr, err = client.GetDatasetsByWorkspaceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "GetDatasetsByWorkspace", resp, "Failure responding to request")
		return
	}
	if result.dlr.hasNextLink() && result.dlr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// GetDatasetsByWorkspacePreparer prepares the GetDatasetsByWorkspace request.
func (client DatasetClient) GetDatasetsByWorkspacePreparer(ctx context.Context) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"endpoint": client.Endpoint,
	}

	const APIVersion = "2019-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{endpoint}", urlParameters),
		autorest.WithPath("/datasets"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetDatasetsByWorkspaceSender sends the GetDatasetsByWorkspace request. The method will close the
// http.Response Body if it receives an error.
func (client DatasetClient) GetDatasetsByWorkspaceSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetDatasetsByWorkspaceResponder handles the response to the GetDatasetsByWorkspace request. The method always
// closes the http.Response Body.
func (client DatasetClient) GetDatasetsByWorkspaceResponder(resp *http.Response) (result DatasetListResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// getDatasetsByWorkspaceNextResults retrieves the next set of results, if any.
func (client DatasetClient) getDatasetsByWorkspaceNextResults(ctx context.Context, lastResults DatasetListResponse) (result DatasetListResponse, err error) {
	req, err := lastResults.datasetListResponsePreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "artifacts.DatasetClient", "getDatasetsByWorkspaceNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.GetDatasetsByWorkspaceSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "artifacts.DatasetClient", "getDatasetsByWorkspaceNextResults", resp, "Failure sending next results request")
	}
	result, err = client.GetDatasetsByWorkspaceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "getDatasetsByWorkspaceNextResults", resp, "Failure responding to next results request")
	}
	return
}

// GetDatasetsByWorkspaceComplete enumerates all values, automatically crossing page boundaries as required.
func (client DatasetClient) GetDatasetsByWorkspaceComplete(ctx context.Context) (result DatasetListResponseIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DatasetClient.GetDatasetsByWorkspace")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.GetDatasetsByWorkspace(ctx)
	return
}

// RenameDataset renames a dataset.
// Parameters:
// datasetName - the dataset name.
// request - proposed new name.
func (client DatasetClient) RenameDataset(ctx context.Context, datasetName string, request RenameRequest) (result DatasetRenameDatasetFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DatasetClient.RenameDataset")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: datasetName,
			Constraints: []validation.Constraint{{Target: "datasetName", Name: validation.MaxLength, Rule: 260, Chain: nil},
				{Target: "datasetName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "datasetName", Name: validation.Pattern, Rule: `^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`, Chain: nil}}},
		{TargetValue: request,
			Constraints: []validation.Constraint{{Target: "request.NewName", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "request.NewName", Name: validation.MaxLength, Rule: 260, Chain: nil},
					{Target: "request.NewName", Name: validation.MinLength, Rule: 1, Chain: nil},
					{Target: "request.NewName", Name: validation.Pattern, Rule: `^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`, Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewError("artifacts.DatasetClient", "RenameDataset", err.Error())
	}

	req, err := client.RenameDatasetPreparer(ctx, datasetName, request)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "RenameDataset", nil, "Failure preparing request")
		return
	}

	result, err = client.RenameDatasetSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.DatasetClient", "RenameDataset", nil, "Failure sending request")
		return
	}

	return
}

// RenameDatasetPreparer prepares the RenameDataset request.
func (client DatasetClient) RenameDatasetPreparer(ctx context.Context, datasetName string, request RenameRequest) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"datasetName": autorest.Encode("path", datasetName),
	}

	const APIVersion = "2019-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{endpoint}", urlParameters),
		autorest.WithPathParameters("/datasets/{datasetName}/rename", pathParameters),
		autorest.WithJSON(request),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// RenameDatasetSender sends the RenameDataset request. The method will close the
// http.Response Body if it receives an error.
func (client DatasetClient) RenameDatasetSender(req *http.Request) (future DatasetRenameDatasetFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// RenameDatasetResponder handles the response to the RenameDataset request. The method always
// closes the http.Response Body.
func (client DatasetClient) RenameDatasetResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}
