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

// PipelineRunClient is the client for the PipelineRun methods of the Artifacts service.
type PipelineRunClient struct {
	BaseClient
}

// NewPipelineRunClient creates an instance of the PipelineRunClient client.
func NewPipelineRunClient(endpoint string) PipelineRunClient {
	return PipelineRunClient{New(endpoint)}
}

// CancelPipelineRun cancel a pipeline run by its run ID.
// Parameters:
// runID - the pipeline run identifier.
// isRecursive - if true, cancel all the Child pipelines that are triggered by the current pipeline.
func (client PipelineRunClient) CancelPipelineRun(ctx context.Context, runID string, isRecursive *bool) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PipelineRunClient.CancelPipelineRun")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CancelPipelineRunPreparer(ctx, runID, isRecursive)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "CancelPipelineRun", nil, "Failure preparing request")
		return
	}

	resp, err := client.CancelPipelineRunSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "CancelPipelineRun", resp, "Failure sending request")
		return
	}

	result, err = client.CancelPipelineRunResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "CancelPipelineRun", resp, "Failure responding to request")
		return
	}

	return
}

// CancelPipelineRunPreparer prepares the CancelPipelineRun request.
func (client PipelineRunClient) CancelPipelineRunPreparer(ctx context.Context, runID string, isRecursive *bool) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"runId": autorest.Encode("path", runID),
	}

	const APIVersion = "2019-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if isRecursive != nil {
		queryParameters["isRecursive"] = autorest.Encode("query", *isRecursive)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{endpoint}", urlParameters),
		autorest.WithPathParameters("/pipelineruns/{runId}/cancel", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CancelPipelineRunSender sends the CancelPipelineRun request. The method will close the
// http.Response Body if it receives an error.
func (client PipelineRunClient) CancelPipelineRunSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// CancelPipelineRunResponder handles the response to the CancelPipelineRun request. The method always
// closes the http.Response Body.
func (client PipelineRunClient) CancelPipelineRunResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// GetPipelineRun get a pipeline run by its run ID.
// Parameters:
// runID - the pipeline run identifier.
func (client PipelineRunClient) GetPipelineRun(ctx context.Context, runID string) (result PipelineRun, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PipelineRunClient.GetPipelineRun")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPipelineRunPreparer(ctx, runID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "GetPipelineRun", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetPipelineRunSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "GetPipelineRun", resp, "Failure sending request")
		return
	}

	result, err = client.GetPipelineRunResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "GetPipelineRun", resp, "Failure responding to request")
		return
	}

	return
}

// GetPipelineRunPreparer prepares the GetPipelineRun request.
func (client PipelineRunClient) GetPipelineRunPreparer(ctx context.Context, runID string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"runId": autorest.Encode("path", runID),
	}

	const APIVersion = "2019-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{endpoint}", urlParameters),
		autorest.WithPathParameters("/pipelineruns/{runId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetPipelineRunSender sends the GetPipelineRun request. The method will close the
// http.Response Body if it receives an error.
func (client PipelineRunClient) GetPipelineRunSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetPipelineRunResponder handles the response to the GetPipelineRun request. The method always
// closes the http.Response Body.
func (client PipelineRunClient) GetPipelineRunResponder(resp *http.Response) (result PipelineRun, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// QueryActivityRuns query activity runs based on input filter conditions.
// Parameters:
// pipelineName - the pipeline name.
// runID - the pipeline run identifier.
// filterParameters - parameters to filter the activity runs.
func (client PipelineRunClient) QueryActivityRuns(ctx context.Context, pipelineName string, runID string, filterParameters RunFilterParameters) (result ActivityRunsQueryResponse, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PipelineRunClient.QueryActivityRuns")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: pipelineName,
			Constraints: []validation.Constraint{{Target: "pipelineName", Name: validation.MaxLength, Rule: 260, Chain: nil},
				{Target: "pipelineName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "pipelineName", Name: validation.Pattern, Rule: `^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`, Chain: nil}}},
		{TargetValue: filterParameters,
			Constraints: []validation.Constraint{{Target: "filterParameters.LastUpdatedAfter", Name: validation.Null, Rule: true, Chain: nil},
				{Target: "filterParameters.LastUpdatedBefore", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("artifacts.PipelineRunClient", "QueryActivityRuns", err.Error())
	}

	req, err := client.QueryActivityRunsPreparer(ctx, pipelineName, runID, filterParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "QueryActivityRuns", nil, "Failure preparing request")
		return
	}

	resp, err := client.QueryActivityRunsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "QueryActivityRuns", resp, "Failure sending request")
		return
	}

	result, err = client.QueryActivityRunsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "QueryActivityRuns", resp, "Failure responding to request")
		return
	}

	return
}

// QueryActivityRunsPreparer prepares the QueryActivityRuns request.
func (client PipelineRunClient) QueryActivityRunsPreparer(ctx context.Context, pipelineName string, runID string, filterParameters RunFilterParameters) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"pipelineName": autorest.Encode("path", pipelineName),
		"runId":        autorest.Encode("path", runID),
	}

	const APIVersion = "2019-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{endpoint}", urlParameters),
		autorest.WithPathParameters("/pipelines/{pipelineName}/pipelineruns/{runId}/queryActivityruns", pathParameters),
		autorest.WithJSON(filterParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// QueryActivityRunsSender sends the QueryActivityRuns request. The method will close the
// http.Response Body if it receives an error.
func (client PipelineRunClient) QueryActivityRunsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// QueryActivityRunsResponder handles the response to the QueryActivityRuns request. The method always
// closes the http.Response Body.
func (client PipelineRunClient) QueryActivityRunsResponder(resp *http.Response) (result ActivityRunsQueryResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// QueryPipelineRunsByWorkspace query pipeline runs in the workspace based on input filter conditions.
// Parameters:
// filterParameters - parameters to filter the pipeline run.
func (client PipelineRunClient) QueryPipelineRunsByWorkspace(ctx context.Context, filterParameters RunFilterParameters) (result PipelineRunsQueryResponse, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PipelineRunClient.QueryPipelineRunsByWorkspace")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: filterParameters,
			Constraints: []validation.Constraint{{Target: "filterParameters.LastUpdatedAfter", Name: validation.Null, Rule: true, Chain: nil},
				{Target: "filterParameters.LastUpdatedBefore", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("artifacts.PipelineRunClient", "QueryPipelineRunsByWorkspace", err.Error())
	}

	req, err := client.QueryPipelineRunsByWorkspacePreparer(ctx, filterParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "QueryPipelineRunsByWorkspace", nil, "Failure preparing request")
		return
	}

	resp, err := client.QueryPipelineRunsByWorkspaceSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "QueryPipelineRunsByWorkspace", resp, "Failure sending request")
		return
	}

	result, err = client.QueryPipelineRunsByWorkspaceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "artifacts.PipelineRunClient", "QueryPipelineRunsByWorkspace", resp, "Failure responding to request")
		return
	}

	return
}

// QueryPipelineRunsByWorkspacePreparer prepares the QueryPipelineRunsByWorkspace request.
func (client PipelineRunClient) QueryPipelineRunsByWorkspacePreparer(ctx context.Context, filterParameters RunFilterParameters) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": client.Endpoint,
	}

	const APIVersion = "2019-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{endpoint}", urlParameters),
		autorest.WithPath("/queryPipelineRuns"),
		autorest.WithJSON(filterParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// QueryPipelineRunsByWorkspaceSender sends the QueryPipelineRunsByWorkspace request. The method will close the
// http.Response Body if it receives an error.
func (client PipelineRunClient) QueryPipelineRunsByWorkspaceSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// QueryPipelineRunsByWorkspaceResponder handles the response to the QueryPipelineRunsByWorkspace request. The method always
// closes the http.Response Body.
func (client PipelineRunClient) QueryPipelineRunsByWorkspaceResponder(resp *http.Response) (result PipelineRunsQueryResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
