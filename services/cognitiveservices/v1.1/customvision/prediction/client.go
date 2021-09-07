// Package prediction implements the Azure ARM Prediction service API version 2.0.
//
//
package prediction

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
	"github.com/gofrs/uuid"
	"io"
	"net/http"
)

// BaseClient is the base client for Prediction.
type BaseClient struct {
	autorest.Client
	APIKey   string
	Endpoint string
}

// New creates an instance of the BaseClient client.
func New(aPIKey string, endpoint string) BaseClient {
	return NewWithoutDefaults(aPIKey, endpoint)
}

// NewWithoutDefaults creates an instance of the BaseClient client.
func NewWithoutDefaults(aPIKey string, endpoint string) BaseClient {
	return BaseClient{
		Client:   autorest.NewClientWithUserAgent(UserAgent()),
		APIKey:   aPIKey,
		Endpoint: endpoint,
	}
}

// PredictImage sends the predict image request.
// Parameters:
// projectID - the project id
// iterationID - optional. Specifies the id of a particular iteration to evaluate against.
// The default iteration for the project will be used when not specified
// application - optional. Specifies the name of application using the endpoint
func (client BaseClient) PredictImage(ctx context.Context, projectID uuid.UUID, imageData io.ReadCloser, iterationID *uuid.UUID, application string) (result ImagePrediction, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BaseClient.PredictImage")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.PredictImagePreparer(ctx, projectID, imageData, iterationID, application)
	if err != nil {
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImage", nil, "Failure preparing request")
		return
	}

	resp, err := client.PredictImageSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImage", resp, "Failure sending request")
		return
	}

	result, err = client.PredictImageResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImage", resp, "Failure responding to request")
		return
	}

	return
}

// PredictImagePreparer prepares the PredictImage request.
func (client BaseClient) PredictImagePreparer(ctx context.Context, projectID uuid.UUID, imageData io.ReadCloser, iterationID *uuid.UUID, application string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"projectId": autorest.Encode("path", projectID),
	}

	queryParameters := map[string]interface{}{}
	if iterationID != nil {
		queryParameters["iterationId"] = autorest.Encode("query", *iterationID)
	}
	if len(application) > 0 {
		queryParameters["application"] = autorest.Encode("query", application)
	}

	formDataParameters := map[string]interface{}{
		"imageData": imageData,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{Endpoint}/customvision/v2.0/Prediction", urlParameters),
		autorest.WithPathParameters("/{projectId}/image", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithMultiPartFormData(formDataParameters),
		autorest.WithHeader("Prediction-Key", client.APIKey))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PredictImageSender sends the PredictImage request. The method will close the
// http.Response Body if it receives an error.
func (client BaseClient) PredictImageSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// PredictImageResponder handles the response to the PredictImage request. The method always
// closes the http.Response Body.
func (client BaseClient) PredictImageResponder(resp *http.Response) (result ImagePrediction, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// PredictImageURL sends the predict image url request.
// Parameters:
// projectID - the project id
// imageURL - an {Iris.Web.Api.Models.ImageUrl} that contains the url of the image to be evaluated
// iterationID - optional. Specifies the id of a particular iteration to evaluate against.
// The default iteration for the project will be used when not specified
// application - optional. Specifies the name of application using the endpoint
func (client BaseClient) PredictImageURL(ctx context.Context, projectID uuid.UUID, imageURL ImageURL, iterationID *uuid.UUID, application string) (result ImagePrediction, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BaseClient.PredictImageURL")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.PredictImageURLPreparer(ctx, projectID, imageURL, iterationID, application)
	if err != nil {
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImageURL", nil, "Failure preparing request")
		return
	}

	resp, err := client.PredictImageURLSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImageURL", resp, "Failure sending request")
		return
	}

	result, err = client.PredictImageURLResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImageURL", resp, "Failure responding to request")
		return
	}

	return
}

// PredictImageURLPreparer prepares the PredictImageURL request.
func (client BaseClient) PredictImageURLPreparer(ctx context.Context, projectID uuid.UUID, imageURL ImageURL, iterationID *uuid.UUID, application string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"projectId": autorest.Encode("path", projectID),
	}

	queryParameters := map[string]interface{}{}
	if iterationID != nil {
		queryParameters["iterationId"] = autorest.Encode("query", *iterationID)
	}
	if len(application) > 0 {
		queryParameters["application"] = autorest.Encode("query", application)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{Endpoint}/customvision/v2.0/Prediction", urlParameters),
		autorest.WithPathParameters("/{projectId}/url", pathParameters),
		autorest.WithJSON(imageURL),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeader("Prediction-Key", client.APIKey))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PredictImageURLSender sends the PredictImageURL request. The method will close the
// http.Response Body if it receives an error.
func (client BaseClient) PredictImageURLSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// PredictImageURLResponder handles the response to the PredictImageURL request. The method always
// closes the http.Response Body.
func (client BaseClient) PredictImageURLResponder(resp *http.Response) (result ImagePrediction, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// PredictImageURLWithNoStore sends the predict image url with no store request.
// Parameters:
// projectID - the project id
// imageURL - an {Iris.Web.Api.Models.ImageUrl} that contains the url of the image to be evaluated
// iterationID - optional. Specifies the id of a particular iteration to evaluate against.
// The default iteration for the project will be used when not specified
// application - optional. Specifies the name of application using the endpoint
func (client BaseClient) PredictImageURLWithNoStore(ctx context.Context, projectID uuid.UUID, imageURL ImageURL, iterationID *uuid.UUID, application string) (result ImagePrediction, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BaseClient.PredictImageURLWithNoStore")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.PredictImageURLWithNoStorePreparer(ctx, projectID, imageURL, iterationID, application)
	if err != nil {
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImageURLWithNoStore", nil, "Failure preparing request")
		return
	}

	resp, err := client.PredictImageURLWithNoStoreSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImageURLWithNoStore", resp, "Failure sending request")
		return
	}

	result, err = client.PredictImageURLWithNoStoreResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImageURLWithNoStore", resp, "Failure responding to request")
		return
	}

	return
}

// PredictImageURLWithNoStorePreparer prepares the PredictImageURLWithNoStore request.
func (client BaseClient) PredictImageURLWithNoStorePreparer(ctx context.Context, projectID uuid.UUID, imageURL ImageURL, iterationID *uuid.UUID, application string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"projectId": autorest.Encode("path", projectID),
	}

	queryParameters := map[string]interface{}{}
	if iterationID != nil {
		queryParameters["iterationId"] = autorest.Encode("query", *iterationID)
	}
	if len(application) > 0 {
		queryParameters["application"] = autorest.Encode("query", application)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{Endpoint}/customvision/v2.0/Prediction", urlParameters),
		autorest.WithPathParameters("/{projectId}/url/nostore", pathParameters),
		autorest.WithJSON(imageURL),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeader("Prediction-Key", client.APIKey))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PredictImageURLWithNoStoreSender sends the PredictImageURLWithNoStore request. The method will close the
// http.Response Body if it receives an error.
func (client BaseClient) PredictImageURLWithNoStoreSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// PredictImageURLWithNoStoreResponder handles the response to the PredictImageURLWithNoStore request. The method always
// closes the http.Response Body.
func (client BaseClient) PredictImageURLWithNoStoreResponder(resp *http.Response) (result ImagePrediction, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// PredictImageWithNoStore sends the predict image with no store request.
// Parameters:
// projectID - the project id
// iterationID - optional. Specifies the id of a particular iteration to evaluate against.
// The default iteration for the project will be used when not specified
// application - optional. Specifies the name of application using the endpoint
func (client BaseClient) PredictImageWithNoStore(ctx context.Context, projectID uuid.UUID, imageData io.ReadCloser, iterationID *uuid.UUID, application string) (result ImagePrediction, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BaseClient.PredictImageWithNoStore")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.PredictImageWithNoStorePreparer(ctx, projectID, imageData, iterationID, application)
	if err != nil {
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImageWithNoStore", nil, "Failure preparing request")
		return
	}

	resp, err := client.PredictImageWithNoStoreSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImageWithNoStore", resp, "Failure sending request")
		return
	}

	result, err = client.PredictImageWithNoStoreResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "prediction.BaseClient", "PredictImageWithNoStore", resp, "Failure responding to request")
		return
	}

	return
}

// PredictImageWithNoStorePreparer prepares the PredictImageWithNoStore request.
func (client BaseClient) PredictImageWithNoStorePreparer(ctx context.Context, projectID uuid.UUID, imageData io.ReadCloser, iterationID *uuid.UUID, application string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":         autorest.Encode("path"),
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"projectId": autorest.Encode("path", projectID),
	}

	queryParameters := map[string]interface{}{}
	if iterationID != nil {
		queryParameters["iterationId"] = autorest.Encode("query", *iterationID)
	}
	if len(application) > 0 {
		queryParameters["application"] = autorest.Encode("query", application)
	}

	formDataParameters := map[string]interface{}{
		"imageData": imageData,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{Endpoint}/customvision/v2.0/Prediction", urlParameters),
		autorest.WithPathParameters("/{projectId}/image/nostore", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithMultiPartFormData(formDataParameters),
		autorest.WithHeader("Prediction-Key", client.APIKey))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PredictImageWithNoStoreSender sends the PredictImageWithNoStore request. The method will close the
// http.Response Body if it receives an error.
func (client BaseClient) PredictImageWithNoStoreSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// PredictImageWithNoStoreResponder handles the response to the PredictImageWithNoStore request. The method always
// closes the http.Response Body.
func (client BaseClient) PredictImageWithNoStoreResponder(resp *http.Response) (result ImagePrediction, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
