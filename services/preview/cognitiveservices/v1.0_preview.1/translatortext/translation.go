package translatortext

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

// TranslationClient is the client for the Translation methods of the Translatortext service.
type TranslationClient struct {
	BaseClient
}

// NewTranslationClient creates an instance of the TranslationClient client.
func NewTranslationClient() TranslationClient {
	return TranslationClient{New()}
}

// CancelOperation cancel a currently processing or queued operation.
// An operation will not be cancelled if it is already completed or failed or cancelling.  A bad request will be
// returned.
// All documents that have completed translation will not be cancelled and will be charged.
// All pending documents will be cancelled if possible.
// Parameters:
// endpoint - supported Cognitive Services endpoints (protocol and hostname, for example:
// https://westus.api.cognitive.microsoft.com).
// ID - format - uuid.  The operation-id
func (client TranslationClient) CancelOperation(ctx context.Context, endpoint string, ID uuid.UUID) (result BatchStatusDetail, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TranslationClient.CancelOperation")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CancelOperationPreparer(ctx, endpoint, ID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "CancelOperation", nil, "Failure preparing request")
		return
	}

	resp, err := client.CancelOperationSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "CancelOperation", resp, "Failure sending request")
		return
	}

	result, err = client.CancelOperationResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "CancelOperation", resp, "Failure responding to request")
		return
	}

	return
}

// CancelOperationPreparer prepares the CancelOperation request.
func (client TranslationClient) CancelOperationPreparer(ctx context.Context, endpoint string, ID uuid.UUID) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": endpoint,
	}

	pathParameters := map[string]interface{}{
		"id": autorest.Encode("path", ID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithCustomBaseURL("{endpoint}/translator/text/batch/v1.0-preview.1", urlParameters),
		autorest.WithPathParameters("/batches/{id}", pathParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CancelOperationSender sends the CancelOperation request. The method will close the
// http.Response Body if it receives an error.
func (client TranslationClient) CancelOperationSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// CancelOperationResponder handles the response to the CancelOperation request. The method always
// closes the http.Response Body.
func (client TranslationClient) CancelOperationResponder(resp *http.Response) (result BatchStatusDetail, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetDocumentFormats the list of supported document formats supported by our service.
// The list will include the common file extension used and supported as well as the content-type if using the upload
// API.
// Parameters:
// endpoint - supported Cognitive Services endpoints (protocol and hostname, for example:
// https://westus.api.cognitive.microsoft.com).
func (client TranslationClient) GetDocumentFormats(ctx context.Context, endpoint string) (result FileFormatListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TranslationClient.GetDocumentFormats")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetDocumentFormatsPreparer(ctx, endpoint)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetDocumentFormats", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetDocumentFormatsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetDocumentFormats", resp, "Failure sending request")
		return
	}

	result, err = client.GetDocumentFormatsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetDocumentFormats", resp, "Failure responding to request")
		return
	}

	return
}

// GetDocumentFormatsPreparer prepares the GetDocumentFormats request.
func (client TranslationClient) GetDocumentFormatsPreparer(ctx context.Context, endpoint string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": endpoint,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{endpoint}/translator/text/batch/v1.0-preview.1", urlParameters),
		autorest.WithPath("/documents/formats"))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetDocumentFormatsSender sends the GetDocumentFormats request. The method will close the
// http.Response Body if it receives an error.
func (client TranslationClient) GetDocumentFormatsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetDocumentFormatsResponder handles the response to the GetDocumentFormats request. The method always
// closes the http.Response Body.
func (client TranslationClient) GetDocumentFormatsResponder(resp *http.Response) (result FileFormatListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetDocumentStatus returns the status of the translation of the document.
// Parameters:
// endpoint - supported Cognitive Services endpoints (protocol and hostname, for example:
// https://westus.api.cognitive.microsoft.com).
// ID - format - uuid.  The batch id
// documentID - format - uuid.  The document id
func (client TranslationClient) GetDocumentStatus(ctx context.Context, endpoint string, ID uuid.UUID, documentID uuid.UUID) (result DocumentStatusDetail, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TranslationClient.GetDocumentStatus")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetDocumentStatusPreparer(ctx, endpoint, ID, documentID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetDocumentStatus", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetDocumentStatusSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetDocumentStatus", resp, "Failure sending request")
		return
	}

	result, err = client.GetDocumentStatusResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetDocumentStatus", resp, "Failure responding to request")
		return
	}

	return
}

// GetDocumentStatusPreparer prepares the GetDocumentStatus request.
func (client TranslationClient) GetDocumentStatusPreparer(ctx context.Context, endpoint string, ID uuid.UUID, documentID uuid.UUID) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": endpoint,
	}

	pathParameters := map[string]interface{}{
		"documentId": autorest.Encode("path", documentID),
		"id":         autorest.Encode("path", ID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{endpoint}/translator/text/batch/v1.0-preview.1", urlParameters),
		autorest.WithPathParameters("/batches/{id}/documents/{documentId}", pathParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetDocumentStatusSender sends the GetDocumentStatus request. The method will close the
// http.Response Body if it receives an error.
func (client TranslationClient) GetDocumentStatusSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetDocumentStatusResponder handles the response to the GetDocumentStatus request. The method always
// closes the http.Response Body.
func (client TranslationClient) GetDocumentStatusResponder(resp *http.Response) (result DocumentStatusDetail, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetDocumentStorageSource the list of storage sources supported by our service.
// Parameters:
// endpoint - supported Cognitive Services endpoints (protocol and hostname, for example:
// https://westus.api.cognitive.microsoft.com).
func (client TranslationClient) GetDocumentStorageSource(ctx context.Context, endpoint string) (result StorageSourceListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TranslationClient.GetDocumentStorageSource")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetDocumentStorageSourcePreparer(ctx, endpoint)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetDocumentStorageSource", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetDocumentStorageSourceSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetDocumentStorageSource", resp, "Failure sending request")
		return
	}

	result, err = client.GetDocumentStorageSourceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetDocumentStorageSource", resp, "Failure responding to request")
		return
	}

	return
}

// GetDocumentStorageSourcePreparer prepares the GetDocumentStorageSource request.
func (client TranslationClient) GetDocumentStorageSourcePreparer(ctx context.Context, endpoint string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": endpoint,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{endpoint}/translator/text/batch/v1.0-preview.1", urlParameters),
		autorest.WithPath("/storagesources"))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetDocumentStorageSourceSender sends the GetDocumentStorageSource request. The method will close the
// http.Response Body if it receives an error.
func (client TranslationClient) GetDocumentStorageSourceSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetDocumentStorageSourceResponder handles the response to the GetDocumentStorageSource request. The method always
// closes the http.Response Body.
func (client TranslationClient) GetDocumentStorageSourceResponder(resp *http.Response) (result StorageSourceListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetGlossaryFormats the list of supported glossary formats supported by our service.
// The list will include the common file extension used.
// Parameters:
// endpoint - supported Cognitive Services endpoints (protocol and hostname, for example:
// https://westus.api.cognitive.microsoft.com).
func (client TranslationClient) GetGlossaryFormats(ctx context.Context, endpoint string) (result FileFormatListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TranslationClient.GetGlossaryFormats")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetGlossaryFormatsPreparer(ctx, endpoint)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetGlossaryFormats", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetGlossaryFormatsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetGlossaryFormats", resp, "Failure sending request")
		return
	}

	result, err = client.GetGlossaryFormatsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetGlossaryFormats", resp, "Failure responding to request")
		return
	}

	return
}

// GetGlossaryFormatsPreparer prepares the GetGlossaryFormats request.
func (client TranslationClient) GetGlossaryFormatsPreparer(ctx context.Context, endpoint string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": endpoint,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{endpoint}/translator/text/batch/v1.0-preview.1", urlParameters),
		autorest.WithPath("/glossaries/formats"))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetGlossaryFormatsSender sends the GetGlossaryFormats request. The method will close the
// http.Response Body if it receives an error.
func (client TranslationClient) GetGlossaryFormatsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetGlossaryFormatsResponder handles the response to the GetGlossaryFormats request. The method always
// closes the http.Response Body.
func (client TranslationClient) GetGlossaryFormatsResponder(resp *http.Response) (result FileFormatListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetOperationDocumentsStatus returns the status of the list of documents translation operation by a given operation
// id.
//
// The documents are sorted by the document Id
//
// If the number of documents exceed our paging limit, server side paging will be used.
// Paginated responses will indicate a partial result by including a continuation token in the response. The absence of
// a continuation token means that no additional pages are available.
//
// Clients MAY use $top and $skip query parameters to specify a number of results to return and an offset into the
// collection.
// The server will honor the values specified by the client; however, clients MUST be prepared to handle responses that
// contain a different page size or contain a continuation token.
// When both $top and $skip are given by a client, the server SHOULD first apply $skip and then $top on the collection.
// Note: If the server can't honor $top and/or $skip, the server MUST return an error to the client informing about it
// instead of just ignoring the query options. This will avoid the risk of the client making assumptions about the data
// returned.
// Parameters:
// endpoint - supported Cognitive Services endpoints (protocol and hostname, for example:
// https://westus.api.cognitive.microsoft.com).
// ID - format - uuid.  The operation id
// top - take the $top entries in the collection
// When both $top and $skip are supplied, $skip is applied first
// skip - skip the $skip entries in the collection
// When both $top and $skip are supplied, $skip is applied first
func (client TranslationClient) GetOperationDocumentsStatus(ctx context.Context, endpoint string, ID uuid.UUID, top *int32, skip *int32) (result DocumentStatusResponse, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TranslationClient.GetOperationDocumentsStatus")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMaximum, Rule: int64(100), Chain: nil},
					{Target: "top", Name: validation.InclusiveMinimum, Rule: int64(1), Chain: nil},
				}}}},
		{TargetValue: skip,
			Constraints: []validation.Constraint{{Target: "skip", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "skip", Name: validation.InclusiveMaximum, Rule: int64(2147483647), Chain: nil},
					{Target: "skip", Name: validation.InclusiveMinimum, Rule: int64(0), Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewError("translatortext.TranslationClient", "GetOperationDocumentsStatus", err.Error())
	}

	req, err := client.GetOperationDocumentsStatusPreparer(ctx, endpoint, ID, top, skip)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetOperationDocumentsStatus", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetOperationDocumentsStatusSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetOperationDocumentsStatus", resp, "Failure sending request")
		return
	}

	result, err = client.GetOperationDocumentsStatusResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetOperationDocumentsStatus", resp, "Failure responding to request")
		return
	}

	return
}

// GetOperationDocumentsStatusPreparer prepares the GetOperationDocumentsStatus request.
func (client TranslationClient) GetOperationDocumentsStatusPreparer(ctx context.Context, endpoint string, ID uuid.UUID, top *int32, skip *int32) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": endpoint,
	}

	pathParameters := map[string]interface{}{
		"id": autorest.Encode("path", ID),
	}

	queryParameters := map[string]interface{}{}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	} else {
		queryParameters["$top"] = autorest.Encode("query", 50)
	}
	if skip != nil {
		queryParameters["$skip"] = autorest.Encode("query", *skip)
	} else {
		queryParameters["$skip"] = autorest.Encode("query", 0)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{endpoint}/translator/text/batch/v1.0-preview.1", urlParameters),
		autorest.WithPathParameters("/batches/{id}/documents", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetOperationDocumentsStatusSender sends the GetOperationDocumentsStatus request. The method will close the
// http.Response Body if it receives an error.
func (client TranslationClient) GetOperationDocumentsStatusSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetOperationDocumentsStatusResponder handles the response to the GetOperationDocumentsStatus request. The method always
// closes the http.Response Body.
func (client TranslationClient) GetOperationDocumentsStatusResponder(resp *http.Response) (result DocumentStatusResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetOperations returns the list of status of the translation batch operation.
// The list will consist only of the batch request submitted by the user (based on their subscription)
//
// The operation status are sorted by the operation Id
//
// If the number of operations exceed our paging limit, server side paging will be used.
// Paginated responses will indicate a partial result by including a continuation token in the response. The absence of
// a continuation token means that no additional pages are available.
//
// Clients MAY use $top and $skip query parameters to specify a number of results to return and an offset into the
// collection.
// The server will honor the values specified by the client; however, clients MUST be prepared to handle responses that
// contain a different page size or contain a continuation token.
// When both $top and $skip are given by a client, the server SHOULD first apply $skip and then $top on the collection.
// Note: If the server can't honor $top and/or $skip, the server MUST return an error to the client informing about it
// instead of just ignoring the query options. This will avoid the risk of the client making assumptions about the data
// returned.
// Parameters:
// endpoint - supported Cognitive Services endpoints (protocol and hostname, for example:
// https://westus.api.cognitive.microsoft.com).
// top - take the $top entries in the collection
// When both $top and $skip are supplied, $skip is applied first
// skip - skip the $skip entries in the collection
// When both $top and $skip are supplied, $skip is applied first
func (client TranslationClient) GetOperations(ctx context.Context, endpoint string, top *int32, skip *int32) (result BatchStatusResponse, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TranslationClient.GetOperations")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMaximum, Rule: int64(100), Chain: nil},
					{Target: "top", Name: validation.InclusiveMinimum, Rule: int64(1), Chain: nil},
				}}}},
		{TargetValue: skip,
			Constraints: []validation.Constraint{{Target: "skip", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "skip", Name: validation.InclusiveMaximum, Rule: int64(2147483647), Chain: nil},
					{Target: "skip", Name: validation.InclusiveMinimum, Rule: int64(0), Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewError("translatortext.TranslationClient", "GetOperations", err.Error())
	}

	req, err := client.GetOperationsPreparer(ctx, endpoint, top, skip)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetOperations", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetOperationsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetOperations", resp, "Failure sending request")
		return
	}

	result, err = client.GetOperationsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetOperations", resp, "Failure responding to request")
		return
	}

	return
}

// GetOperationsPreparer prepares the GetOperations request.
func (client TranslationClient) GetOperationsPreparer(ctx context.Context, endpoint string, top *int32, skip *int32) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": endpoint,
	}

	queryParameters := map[string]interface{}{}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	} else {
		queryParameters["$top"] = autorest.Encode("query", 50)
	}
	if skip != nil {
		queryParameters["$skip"] = autorest.Encode("query", *skip)
	} else {
		queryParameters["$skip"] = autorest.Encode("query", 0)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{endpoint}/translator/text/batch/v1.0-preview.1", urlParameters),
		autorest.WithPath("/batches"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetOperationsSender sends the GetOperations request. The method will close the
// http.Response Body if it receives an error.
func (client TranslationClient) GetOperationsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetOperationsResponder handles the response to the GetOperations request. The method always
// closes the http.Response Body.
func (client TranslationClient) GetOperationsResponder(resp *http.Response) (result BatchStatusResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetOperationStatus returns the status of the translation batch operation.
// The status will include the overall job status as well as a summary of the current progress of all the documents
// being translated.
// Parameters:
// endpoint - supported Cognitive Services endpoints (protocol and hostname, for example:
// https://westus.api.cognitive.microsoft.com).
// ID - format - uuid.  The operation id
func (client TranslationClient) GetOperationStatus(ctx context.Context, endpoint string, ID uuid.UUID) (result BatchStatusDetail, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TranslationClient.GetOperationStatus")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetOperationStatusPreparer(ctx, endpoint, ID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetOperationStatus", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetOperationStatusSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetOperationStatus", resp, "Failure sending request")
		return
	}

	result, err = client.GetOperationStatusResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "GetOperationStatus", resp, "Failure responding to request")
		return
	}

	return
}

// GetOperationStatusPreparer prepares the GetOperationStatus request.
func (client TranslationClient) GetOperationStatusPreparer(ctx context.Context, endpoint string, ID uuid.UUID) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": endpoint,
	}

	pathParameters := map[string]interface{}{
		"id": autorest.Encode("path", ID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{endpoint}/translator/text/batch/v1.0-preview.1", urlParameters),
		autorest.WithPathParameters("/batches/{id}", pathParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetOperationStatusSender sends the GetOperationStatus request. The method will close the
// http.Response Body if it receives an error.
func (client TranslationClient) GetOperationStatusSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetOperationStatusResponder handles the response to the GetOperationStatus request. The method always
// closes the http.Response Body.
func (client TranslationClient) GetOperationStatusResponder(resp *http.Response) (result BatchStatusDetail, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// SubmitBatchRequest submit a batch request to the document translation service.
//
// Each request can consists of multiple inputs.
// Each input will contains both a source and destination container for source and target language pair.
//
// The prefix and suffix filter (if supplied) will be used to filter the folders.
// The prefix will be applied to the subpath after the container name
//
// Glossaries / Translation memory can be supplied and will be applied when the document is being translated.
// If the glossary is invalid or unreachable during translation time.  An error will be indicated in the document
// status.
//
// If the file with the same name already exists in the destination, it will be overwritten.
// TargetUrl for each target language needs to be unique.
// Parameters:
// endpoint - supported Cognitive Services endpoints (protocol and hostname, for example:
// https://westus.api.cognitive.microsoft.com).
// body - request details
func (client TranslationClient) SubmitBatchRequest(ctx context.Context, endpoint string, body *BatchSubmissionRequest) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TranslationClient.SubmitBatchRequest")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: body,
			Constraints: []validation.Constraint{{Target: "body", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "body.Inputs", Name: validation.Null, Rule: true, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("translatortext.TranslationClient", "SubmitBatchRequest", err.Error())
	}

	req, err := client.SubmitBatchRequestPreparer(ctx, endpoint, body)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "SubmitBatchRequest", nil, "Failure preparing request")
		return
	}

	resp, err := client.SubmitBatchRequestSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "SubmitBatchRequest", resp, "Failure sending request")
		return
	}

	result, err = client.SubmitBatchRequestResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "translatortext.TranslationClient", "SubmitBatchRequest", resp, "Failure responding to request")
		return
	}

	return
}

// SubmitBatchRequestPreparer prepares the SubmitBatchRequest request.
func (client TranslationClient) SubmitBatchRequestPreparer(ctx context.Context, endpoint string, body *BatchSubmissionRequest) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"endpoint": endpoint,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{endpoint}/translator/text/batch/v1.0-preview.1", urlParameters),
		autorest.WithPath("/batches"))
	if body != nil {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithJSON(body))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SubmitBatchRequestSender sends the SubmitBatchRequest request. The method will close the
// http.Response Body if it receives an error.
func (client TranslationClient) SubmitBatchRequestSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// SubmitBatchRequestResponder handles the response to the SubmitBatchRequest request. The method always
// closes the http.Response Body.
func (client TranslationClient) SubmitBatchRequestResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}
