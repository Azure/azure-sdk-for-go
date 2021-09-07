package contentmoderator

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
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

// ListManagementTermClient is the you use the API to scan your content as it is generated. Content Moderator then
// processes your content and sends the results along with relevant information either back to your systems or to the
// built-in review tool. You can use this information to take decisions e.g. take it down, send to human judge, etc.
//
// When using the API, images need to have a minimum of 128 pixels and a maximum file size of 4MB.
// Text can be at most 1024 characters long.
// If the content passed to the text API or the image API exceeds the size limits, the API will return an error code
// that informs about the issue.
type ListManagementTermClient struct {
	BaseClient
}

// NewListManagementTermClient creates an instance of the ListManagementTermClient client.
func NewListManagementTermClient(endpoint string) ListManagementTermClient {
	return ListManagementTermClient{New(endpoint)}
}

// AddTerm add a term to the term list with list Id equal to list Id passed.
// Parameters:
// listID - list Id of the image list.
// term - term to be deleted
// language - language of the terms.
func (client ListManagementTermClient) AddTerm(ctx context.Context, listID string, term string, language string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ListManagementTermClient.AddTerm")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.AddTermPreparer(ctx, listID, term, language)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "AddTerm", nil, "Failure preparing request")
		return
	}

	resp, err := client.AddTermSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "AddTerm", resp, "Failure sending request")
		return
	}

	result, err = client.AddTermResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "AddTerm", resp, "Failure responding to request")
		return
	}

	return
}

// AddTermPreparer prepares the AddTerm request.
func (client ListManagementTermClient) AddTermPreparer(ctx context.Context, listID string, term string, language string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"listId": autorest.Encode("path", listID),
		"term":   autorest.Encode("path", term),
	}

	queryParameters := map[string]interface{}{
		"language": autorest.Encode("query", language),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{Endpoint}", urlParameters),
		autorest.WithPathParameters("/contentmoderator/lists/v1.0/termlists/{listId}/terms/{term}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// AddTermSender sends the AddTerm request. The method will close the
// http.Response Body if it receives an error.
func (client ListManagementTermClient) AddTermSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// AddTermResponder handles the response to the AddTerm request. The method always
// closes the http.Response Body.
func (client ListManagementTermClient) AddTermResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByClosing())
	result.Response = resp
	return
}

// DeleteAllTerms deletes all terms from the list with list Id equal to the list Id passed.
// Parameters:
// listID - list Id of the image list.
// language - language of the terms.
func (client ListManagementTermClient) DeleteAllTerms(ctx context.Context, listID string, language string) (result String, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ListManagementTermClient.DeleteAllTerms")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeleteAllTermsPreparer(ctx, listID, language)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "DeleteAllTerms", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteAllTermsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "DeleteAllTerms", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteAllTermsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "DeleteAllTerms", resp, "Failure responding to request")
		return
	}

	return
}

// DeleteAllTermsPreparer prepares the DeleteAllTerms request.
func (client ListManagementTermClient) DeleteAllTermsPreparer(ctx context.Context, listID string, language string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"listId": autorest.Encode("path", listID),
	}

	queryParameters := map[string]interface{}{
		"language": autorest.Encode("query", language),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithCustomBaseURL("{Endpoint}", urlParameters),
		autorest.WithPathParameters("/contentmoderator/lists/v1.0/termlists/{listId}/terms", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteAllTermsSender sends the DeleteAllTerms request. The method will close the
// http.Response Body if it receives an error.
func (client ListManagementTermClient) DeleteAllTermsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// DeleteAllTermsResponder handles the response to the DeleteAllTerms request. The method always
// closes the http.Response Body.
func (client ListManagementTermClient) DeleteAllTermsResponder(resp *http.Response) (result String, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// DeleteTerm deletes a term from the list with list Id equal to the list Id passed.
// Parameters:
// listID - list Id of the image list.
// term - term to be deleted
// language - language of the terms.
func (client ListManagementTermClient) DeleteTerm(ctx context.Context, listID string, term string, language string) (result String, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ListManagementTermClient.DeleteTerm")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeleteTermPreparer(ctx, listID, term, language)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "DeleteTerm", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteTermSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "DeleteTerm", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteTermResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "DeleteTerm", resp, "Failure responding to request")
		return
	}

	return
}

// DeleteTermPreparer prepares the DeleteTerm request.
func (client ListManagementTermClient) DeleteTermPreparer(ctx context.Context, listID string, term string, language string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"listId": autorest.Encode("path", listID),
		"term":   autorest.Encode("path", term),
	}

	queryParameters := map[string]interface{}{
		"language": autorest.Encode("query", language),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithCustomBaseURL("{Endpoint}", urlParameters),
		autorest.WithPathParameters("/contentmoderator/lists/v1.0/termlists/{listId}/terms/{term}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteTermSender sends the DeleteTerm request. The method will close the
// http.Response Body if it receives an error.
func (client ListManagementTermClient) DeleteTermSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// DeleteTermResponder handles the response to the DeleteTerm request. The method always
// closes the http.Response Body.
func (client ListManagementTermClient) DeleteTermResponder(resp *http.Response) (result String, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetAllTerms gets all terms from the list with list Id equal to the list Id passed.
// Parameters:
// listID - list Id of the image list.
// language - language of the terms.
// offset - the pagination start index.
// limit - the max limit.
func (client ListManagementTermClient) GetAllTerms(ctx context.Context, listID string, language string, offset *int32, limit *int32) (result Terms, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ListManagementTermClient.GetAllTerms")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetAllTermsPreparer(ctx, listID, language, offset, limit)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "GetAllTerms", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetAllTermsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "GetAllTerms", resp, "Failure sending request")
		return
	}

	result, err = client.GetAllTermsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentmoderator.ListManagementTermClient", "GetAllTerms", resp, "Failure responding to request")
		return
	}

	return
}

// GetAllTermsPreparer prepares the GetAllTerms request.
func (client ListManagementTermClient) GetAllTermsPreparer(ctx context.Context, listID string, language string, offset *int32, limit *int32) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"listId": autorest.Encode("path", listID),
	}

	queryParameters := map[string]interface{}{
		"language": autorest.Encode("query", language),
	}
	if offset != nil {
		queryParameters["offset"] = autorest.Encode("query", *offset)
	}
	if limit != nil {
		queryParameters["limit"] = autorest.Encode("query", *limit)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{Endpoint}", urlParameters),
		autorest.WithPathParameters("/contentmoderator/lists/v1.0/termlists/{listId}/terms", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetAllTermsSender sends the GetAllTerms request. The method will close the
// http.Response Body if it receives an error.
func (client ListManagementTermClient) GetAllTermsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetAllTermsResponder handles the response to the GetAllTerms request. The method always
// closes the http.Response Body.
func (client ListManagementTermClient) GetAllTermsResponder(resp *http.Response) (result Terms, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
