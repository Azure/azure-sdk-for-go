package authoring

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
	"github.com/gofrs/uuid"
	"net/http"
)

// PatternClient is the client for the Pattern methods of the Authoring service.
type PatternClient struct {
	BaseClient
}

// NewPatternClient creates an instance of the PatternClient client.
func NewPatternClient(endpoint string) PatternClient {
	return PatternClient{New(endpoint)}
}

// AddPattern sends the add pattern request.
// Parameters:
// appID - the application ID.
// versionID - the version ID.
// pattern - the input pattern.
func (client PatternClient) AddPattern(ctx context.Context, appID uuid.UUID, versionID string, pattern PatternRuleCreateObject) (result PatternRuleInfo, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PatternClient.AddPattern")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.AddPatternPreparer(ctx, appID, versionID, pattern)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "AddPattern", nil, "Failure preparing request")
		return
	}

	resp, err := client.AddPatternSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "AddPattern", resp, "Failure sending request")
		return
	}

	result, err = client.AddPatternResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "AddPattern", resp, "Failure responding to request")
		return
	}

	return
}

// AddPatternPreparer prepares the AddPattern request.
func (client PatternClient) AddPatternPreparer(ctx context.Context, appID uuid.UUID, versionID string, pattern PatternRuleCreateObject) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"appId":     autorest.Encode("path", appID),
		"versionId": autorest.Encode("path", versionID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{Endpoint}/luis/api/v2.0", urlParameters),
		autorest.WithPathParameters("/apps/{appId}/versions/{versionId}/patternrule", pathParameters),
		autorest.WithJSON(pattern))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// AddPatternSender sends the AddPattern request. The method will close the
// http.Response Body if it receives an error.
func (client PatternClient) AddPatternSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// AddPatternResponder handles the response to the AddPattern request. The method always
// closes the http.Response Body.
func (client PatternClient) AddPatternResponder(resp *http.Response) (result PatternRuleInfo, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// BatchAddPatterns sends the batch add patterns request.
// Parameters:
// appID - the application ID.
// versionID - the version ID.
// patterns - a JSON array containing patterns.
func (client PatternClient) BatchAddPatterns(ctx context.Context, appID uuid.UUID, versionID string, patterns []PatternRuleCreateObject) (result ListPatternRuleInfo, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PatternClient.BatchAddPatterns")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: patterns,
			Constraints: []validation.Constraint{{Target: "patterns", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("authoring.PatternClient", "BatchAddPatterns", err.Error())
	}

	req, err := client.BatchAddPatternsPreparer(ctx, appID, versionID, patterns)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "BatchAddPatterns", nil, "Failure preparing request")
		return
	}

	resp, err := client.BatchAddPatternsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "BatchAddPatterns", resp, "Failure sending request")
		return
	}

	result, err = client.BatchAddPatternsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "BatchAddPatterns", resp, "Failure responding to request")
		return
	}

	return
}

// BatchAddPatternsPreparer prepares the BatchAddPatterns request.
func (client PatternClient) BatchAddPatternsPreparer(ctx context.Context, appID uuid.UUID, versionID string, patterns []PatternRuleCreateObject) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"appId":     autorest.Encode("path", appID),
		"versionId": autorest.Encode("path", versionID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{Endpoint}/luis/api/v2.0", urlParameters),
		autorest.WithPathParameters("/apps/{appId}/versions/{versionId}/patternrules", pathParameters),
		autorest.WithJSON(patterns))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// BatchAddPatternsSender sends the BatchAddPatterns request. The method will close the
// http.Response Body if it receives an error.
func (client PatternClient) BatchAddPatternsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// BatchAddPatternsResponder handles the response to the BatchAddPatterns request. The method always
// closes the http.Response Body.
func (client PatternClient) BatchAddPatternsResponder(resp *http.Response) (result ListPatternRuleInfo, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// DeletePattern sends the delete pattern request.
// Parameters:
// appID - the application ID.
// versionID - the version ID.
// patternID - the pattern ID.
func (client PatternClient) DeletePattern(ctx context.Context, appID uuid.UUID, versionID string, patternID uuid.UUID) (result OperationStatus, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PatternClient.DeletePattern")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePatternPreparer(ctx, appID, versionID, patternID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "DeletePattern", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeletePatternSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "DeletePattern", resp, "Failure sending request")
		return
	}

	result, err = client.DeletePatternResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "DeletePattern", resp, "Failure responding to request")
		return
	}

	return
}

// DeletePatternPreparer prepares the DeletePattern request.
func (client PatternClient) DeletePatternPreparer(ctx context.Context, appID uuid.UUID, versionID string, patternID uuid.UUID) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"appId":     autorest.Encode("path", appID),
		"patternId": autorest.Encode("path", patternID),
		"versionId": autorest.Encode("path", versionID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithCustomBaseURL("{Endpoint}/luis/api/v2.0", urlParameters),
		autorest.WithPathParameters("/apps/{appId}/versions/{versionId}/patternrules/{patternId}", pathParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeletePatternSender sends the DeletePattern request. The method will close the
// http.Response Body if it receives an error.
func (client PatternClient) DeletePatternSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// DeletePatternResponder handles the response to the DeletePattern request. The method always
// closes the http.Response Body.
func (client PatternClient) DeletePatternResponder(resp *http.Response) (result OperationStatus, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// DeletePatterns sends the delete patterns request.
// Parameters:
// appID - the application ID.
// versionID - the version ID.
// patternIds - the patterns IDs.
func (client PatternClient) DeletePatterns(ctx context.Context, appID uuid.UUID, versionID string, patternIds []uuid.UUID) (result OperationStatus, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PatternClient.DeletePatterns")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: patternIds,
			Constraints: []validation.Constraint{{Target: "patternIds", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("authoring.PatternClient", "DeletePatterns", err.Error())
	}

	req, err := client.DeletePatternsPreparer(ctx, appID, versionID, patternIds)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "DeletePatterns", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeletePatternsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "DeletePatterns", resp, "Failure sending request")
		return
	}

	result, err = client.DeletePatternsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "DeletePatterns", resp, "Failure responding to request")
		return
	}

	return
}

// DeletePatternsPreparer prepares the DeletePatterns request.
func (client PatternClient) DeletePatternsPreparer(ctx context.Context, appID uuid.UUID, versionID string, patternIds []uuid.UUID) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"appId":     autorest.Encode("path", appID),
		"versionId": autorest.Encode("path", versionID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithCustomBaseURL("{Endpoint}/luis/api/v2.0", urlParameters),
		autorest.WithPathParameters("/apps/{appId}/versions/{versionId}/patternrules", pathParameters),
		autorest.WithJSON(patternIds))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeletePatternsSender sends the DeletePatterns request. The method will close the
// http.Response Body if it receives an error.
func (client PatternClient) DeletePatternsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// DeletePatternsResponder handles the response to the DeletePatterns request. The method always
// closes the http.Response Body.
func (client PatternClient) DeletePatternsResponder(resp *http.Response) (result OperationStatus, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListIntentPatterns sends the list intent patterns request.
// Parameters:
// appID - the application ID.
// versionID - the version ID.
// intentID - the intent classifier ID.
// skip - the number of entries to skip. Default value is 0.
// take - the number of entries to return. Maximum page size is 500. Default is 100.
func (client PatternClient) ListIntentPatterns(ctx context.Context, appID uuid.UUID, versionID string, intentID uuid.UUID, skip *int32, take *int32) (result ListPatternRuleInfo, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PatternClient.ListIntentPatterns")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: skip,
			Constraints: []validation.Constraint{{Target: "skip", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "skip", Name: validation.InclusiveMinimum, Rule: int64(0), Chain: nil}}}}},
		{TargetValue: take,
			Constraints: []validation.Constraint{{Target: "take", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "take", Name: validation.InclusiveMaximum, Rule: int64(500), Chain: nil},
					{Target: "take", Name: validation.InclusiveMinimum, Rule: int64(0), Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewError("authoring.PatternClient", "ListIntentPatterns", err.Error())
	}

	req, err := client.ListIntentPatternsPreparer(ctx, appID, versionID, intentID, skip, take)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "ListIntentPatterns", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListIntentPatternsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "ListIntentPatterns", resp, "Failure sending request")
		return
	}

	result, err = client.ListIntentPatternsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "ListIntentPatterns", resp, "Failure responding to request")
		return
	}

	return
}

// ListIntentPatternsPreparer prepares the ListIntentPatterns request.
func (client PatternClient) ListIntentPatternsPreparer(ctx context.Context, appID uuid.UUID, versionID string, intentID uuid.UUID, skip *int32, take *int32) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"appId":     autorest.Encode("path", appID),
		"intentId":  autorest.Encode("path", intentID),
		"versionId": autorest.Encode("path", versionID),
	}

	queryParameters := map[string]interface{}{}
	if skip != nil {
		queryParameters["skip"] = autorest.Encode("query", *skip)
	} else {
		queryParameters["skip"] = autorest.Encode("query", 0)
	}
	if take != nil {
		queryParameters["take"] = autorest.Encode("query", *take)
	} else {
		queryParameters["take"] = autorest.Encode("query", 100)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{Endpoint}/luis/api/v2.0", urlParameters),
		autorest.WithPathParameters("/apps/{appId}/versions/{versionId}/intents/{intentId}/patternrules", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListIntentPatternsSender sends the ListIntentPatterns request. The method will close the
// http.Response Body if it receives an error.
func (client PatternClient) ListIntentPatternsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListIntentPatternsResponder handles the response to the ListIntentPatterns request. The method always
// closes the http.Response Body.
func (client PatternClient) ListIntentPatternsResponder(resp *http.Response) (result ListPatternRuleInfo, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListPatterns sends the list patterns request.
// Parameters:
// appID - the application ID.
// versionID - the version ID.
// skip - the number of entries to skip. Default value is 0.
// take - the number of entries to return. Maximum page size is 500. Default is 100.
func (client PatternClient) ListPatterns(ctx context.Context, appID uuid.UUID, versionID string, skip *int32, take *int32) (result ListPatternRuleInfo, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PatternClient.ListPatterns")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: skip,
			Constraints: []validation.Constraint{{Target: "skip", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "skip", Name: validation.InclusiveMinimum, Rule: int64(0), Chain: nil}}}}},
		{TargetValue: take,
			Constraints: []validation.Constraint{{Target: "take", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "take", Name: validation.InclusiveMaximum, Rule: int64(500), Chain: nil},
					{Target: "take", Name: validation.InclusiveMinimum, Rule: int64(0), Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewError("authoring.PatternClient", "ListPatterns", err.Error())
	}

	req, err := client.ListPatternsPreparer(ctx, appID, versionID, skip, take)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "ListPatterns", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListPatternsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "ListPatterns", resp, "Failure sending request")
		return
	}

	result, err = client.ListPatternsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "ListPatterns", resp, "Failure responding to request")
		return
	}

	return
}

// ListPatternsPreparer prepares the ListPatterns request.
func (client PatternClient) ListPatternsPreparer(ctx context.Context, appID uuid.UUID, versionID string, skip *int32, take *int32) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"appId":     autorest.Encode("path", appID),
		"versionId": autorest.Encode("path", versionID),
	}

	queryParameters := map[string]interface{}{}
	if skip != nil {
		queryParameters["skip"] = autorest.Encode("query", *skip)
	} else {
		queryParameters["skip"] = autorest.Encode("query", 0)
	}
	if take != nil {
		queryParameters["take"] = autorest.Encode("query", *take)
	} else {
		queryParameters["take"] = autorest.Encode("query", 100)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{Endpoint}/luis/api/v2.0", urlParameters),
		autorest.WithPathParameters("/apps/{appId}/versions/{versionId}/patternrules", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListPatternsSender sends the ListPatterns request. The method will close the
// http.Response Body if it receives an error.
func (client PatternClient) ListPatternsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListPatternsResponder handles the response to the ListPatterns request. The method always
// closes the http.Response Body.
func (client PatternClient) ListPatternsResponder(resp *http.Response) (result ListPatternRuleInfo, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// UpdatePattern sends the update pattern request.
// Parameters:
// appID - the application ID.
// versionID - the version ID.
// patternID - the pattern ID.
// pattern - an object representing a pattern.
func (client PatternClient) UpdatePattern(ctx context.Context, appID uuid.UUID, versionID string, patternID uuid.UUID, pattern PatternRuleUpdateObject) (result PatternRuleInfo, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PatternClient.UpdatePattern")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UpdatePatternPreparer(ctx, appID, versionID, patternID, pattern)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "UpdatePattern", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdatePatternSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "UpdatePattern", resp, "Failure sending request")
		return
	}

	result, err = client.UpdatePatternResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "UpdatePattern", resp, "Failure responding to request")
		return
	}

	return
}

// UpdatePatternPreparer prepares the UpdatePattern request.
func (client PatternClient) UpdatePatternPreparer(ctx context.Context, appID uuid.UUID, versionID string, patternID uuid.UUID, pattern PatternRuleUpdateObject) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"appId":     autorest.Encode("path", appID),
		"patternId": autorest.Encode("path", patternID),
		"versionId": autorest.Encode("path", versionID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithCustomBaseURL("{Endpoint}/luis/api/v2.0", urlParameters),
		autorest.WithPathParameters("/apps/{appId}/versions/{versionId}/patternrules/{patternId}", pathParameters),
		autorest.WithJSON(pattern))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdatePatternSender sends the UpdatePattern request. The method will close the
// http.Response Body if it receives an error.
func (client PatternClient) UpdatePatternSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// UpdatePatternResponder handles the response to the UpdatePattern request. The method always
// closes the http.Response Body.
func (client PatternClient) UpdatePatternResponder(resp *http.Response) (result PatternRuleInfo, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// UpdatePatterns sends the update patterns request.
// Parameters:
// appID - the application ID.
// versionID - the version ID.
// patterns - an array represents the patterns.
func (client PatternClient) UpdatePatterns(ctx context.Context, appID uuid.UUID, versionID string, patterns []PatternRuleUpdateObject) (result ListPatternRuleInfo, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PatternClient.UpdatePatterns")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: patterns,
			Constraints: []validation.Constraint{{Target: "patterns", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("authoring.PatternClient", "UpdatePatterns", err.Error())
	}

	req, err := client.UpdatePatternsPreparer(ctx, appID, versionID, patterns)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "UpdatePatterns", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdatePatternsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "UpdatePatterns", resp, "Failure sending request")
		return
	}

	result, err = client.UpdatePatternsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authoring.PatternClient", "UpdatePatterns", resp, "Failure responding to request")
		return
	}

	return
}

// UpdatePatternsPreparer prepares the UpdatePatterns request.
func (client PatternClient) UpdatePatternsPreparer(ctx context.Context, appID uuid.UUID, versionID string, patterns []PatternRuleUpdateObject) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"appId":     autorest.Encode("path", appID),
		"versionId": autorest.Encode("path", versionID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithCustomBaseURL("{Endpoint}/luis/api/v2.0", urlParameters),
		autorest.WithPathParameters("/apps/{appId}/versions/{versionId}/patternrules", pathParameters),
		autorest.WithJSON(patterns))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdatePatternsSender sends the UpdatePatterns request. The method will close the
// http.Response Body if it receives an error.
func (client PatternClient) UpdatePatternsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// UpdatePatternsResponder handles the response to the UpdatePatterns request. The method always
// closes the http.Response Body.
func (client PatternClient) UpdatePatternsResponder(resp *http.Response) (result ListPatternRuleInfo, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
