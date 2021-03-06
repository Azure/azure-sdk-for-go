package personalizer

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

// EvaluationsClient is the personalizer Service is an Azure Cognitive Service that makes it easy to target content and
// experiences without complex pre-analysis or cleanup of past data. Given a context and featurized content, the
// Personalizer Service returns which content item to show to users in rewardActionId. As rewards are sent in response
// to the use of rewardActionId, the reinforcement learning algorithm will improve the model and improve performance of
// future rank calls.
type EvaluationsClient struct {
	BaseClient
}

// NewEvaluationsClient creates an instance of the EvaluationsClient client.
func NewEvaluationsClient(endpoint string) EvaluationsClient {
	return EvaluationsClient{New(endpoint)}
}

// Create submit a new evaluation job.
// Parameters:
// evaluation - the evaluation job definition.
func (client EvaluationsClient) Create(ctx context.Context, evaluation EvaluationContract) (result Evaluation, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/EvaluationsClient.Create")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: evaluation,
			Constraints: []validation.Constraint{{Target: "evaluation.Name", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "evaluation.Name", Name: validation.MaxLength, Rule: 256, Chain: nil}}},
				{Target: "evaluation.StartTime", Name: validation.Null, Rule: true, Chain: nil},
				{Target: "evaluation.EndTime", Name: validation.Null, Rule: true, Chain: nil},
				{Target: "evaluation.Policies", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("personalizer.EvaluationsClient", "Create", err.Error())
	}

	req, err := client.CreatePreparer(ctx, evaluation)
	if err != nil {
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "Create", resp, "Failure responding to request")
		return
	}

	return
}

// CreatePreparer prepares the Create request.
func (client EvaluationsClient) CreatePreparer(ctx context.Context, evaluation EvaluationContract) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{Endpoint}/personalizer/v1.0", urlParameters),
		autorest.WithPath("/evaluations"),
		autorest.WithJSON(evaluation))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client EvaluationsClient) CreateSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client EvaluationsClient) CreateResponder(resp *http.Response) (result Evaluation, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete delete the evaluation associated with the Id.
// Parameters:
// evaluationID - id of the evaluation to delete.
func (client EvaluationsClient) Delete(ctx context.Context, evaluationID string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/EvaluationsClient.Delete")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: evaluationID,
			Constraints: []validation.Constraint{{Target: "evaluationID", Name: validation.MaxLength, Rule: 256, Chain: nil}}}}); err != nil {
		return result, validation.NewError("personalizer.EvaluationsClient", "Delete", err.Error())
	}

	req, err := client.DeletePreparer(ctx, evaluationID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "Delete", resp, "Failure responding to request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client EvaluationsClient) DeletePreparer(ctx context.Context, evaluationID string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"evaluationId": autorest.Encode("path", evaluationID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithCustomBaseURL("{Endpoint}/personalizer/v1.0", urlParameters),
		autorest.WithPathParameters("/evaluations/{evaluationId}", pathParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client EvaluationsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client EvaluationsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get get the evaluation associated with the Id.
// Parameters:
// evaluationID - id of the evaluation.
func (client EvaluationsClient) Get(ctx context.Context, evaluationID string) (result Evaluation, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/EvaluationsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: evaluationID,
			Constraints: []validation.Constraint{{Target: "evaluationID", Name: validation.MaxLength, Rule: 256, Chain: nil}}}}); err != nil {
		return result, validation.NewError("personalizer.EvaluationsClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, evaluationID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client EvaluationsClient) GetPreparer(ctx context.Context, evaluationID string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"evaluationId": autorest.Encode("path", evaluationID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{Endpoint}/personalizer/v1.0", urlParameters),
		autorest.WithPathParameters("/evaluations/{evaluationId}", pathParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client EvaluationsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client EvaluationsClient) GetResponder(resp *http.Response) (result Evaluation, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List list all the submitted evaluations.
func (client EvaluationsClient) List(ctx context.Context) (result ListEvaluation, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/EvaluationsClient.List")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "personalizer.EvaluationsClient", "List", resp, "Failure responding to request")
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client EvaluationsClient) ListPreparer(ctx context.Context) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{Endpoint}/personalizer/v1.0", urlParameters),
		autorest.WithPath("/evaluations"))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client EvaluationsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client EvaluationsClient) ListResponder(resp *http.Response) (result ListEvaluation, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
