package monitor

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

// MetricsClient is the monitor Management Client
type MetricsClient struct {
	BaseClient
}

// NewMetricsClient creates an instance of the MetricsClient client.
func NewMetricsClient() MetricsClient {
	return NewMetricsClientWithBaseURI(DefaultBaseURI)
}

// NewMetricsClientWithBaseURI creates an instance of the MetricsClient client using a custom endpoint.  Use this when
// interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewMetricsClientWithBaseURI(baseURI string) MetricsClient {
	return MetricsClient{NewWithBaseURI(baseURI)}
}

// Create **Post the metric values for a resource**.
// Parameters:
// contentType - supports application/json and application/x-ndjson
// contentLength - content length of the payload
// subscriptionID - the azure subscription id
// resourceGroupName - the ARM resource group name
// resourceProvider - the ARM resource provider name
// resourceTypeName - the ARM resource type name
// resourceName - the ARM resource name
// body - the Azure metrics document json payload
func (client MetricsClient) Create(ctx context.Context, contentType string, contentLength int32, subscriptionID string, resourceGroupName string, resourceProvider string, resourceTypeName string, resourceName string, body AzureMetricsDocument) (result AzureMetricsResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/MetricsClient.Create")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: body,
			Constraints: []validation.Constraint{{Target: "body.Time", Name: validation.Null, Rule: true, Chain: nil},
				{Target: "body.Data", Name: validation.Null, Rule: true,
					Chain: []validation.Constraint{{Target: "body.Data.BaseData", Name: validation.Null, Rule: true,
						Chain: []validation.Constraint{{Target: "body.Data.BaseData.Metric", Name: validation.Null, Rule: true, Chain: nil},
							{Target: "body.Data.BaseData.Namespace", Name: validation.Null, Rule: true, Chain: nil},
							{Target: "body.Data.BaseData.Series", Name: validation.Null, Rule: true, Chain: nil},
						}},
					}}}}}); err != nil {
		return result, validation.NewError("monitor.MetricsClient", "Create", err.Error())
	}

	req, err := client.CreatePreparer(ctx, contentType, contentLength, subscriptionID, resourceGroupName, resourceProvider, resourceTypeName, resourceName, body)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitor.MetricsClient", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "monitor.MetricsClient", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitor.MetricsClient", "Create", resp, "Failure responding to request")
		return
	}

	return
}

// CreatePreparer prepares the Create request.
func (client MetricsClient) CreatePreparer(ctx context.Context, contentType string, contentLength int32, subscriptionID string, resourceGroupName string, resourceProvider string, resourceTypeName string, resourceName string, body AzureMetricsDocument) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"resourceName":      autorest.Encode("path", resourceName),
		"resourceProvider":  autorest.Encode("path", resourceProvider),
		"resourceTypeName":  autorest.Encode("path", resourceTypeName),
		"subscriptionId":    autorest.Encode("path", subscriptionID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/{resourceProvider}/{resourceTypeName}/{resourceName}/metrics", pathParameters),
		autorest.WithJSON(body),
		autorest.WithHeader("Content-Type", autorest.String(contentType)),
		autorest.WithHeader("Content-Length", autorest.String(contentLength)))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client MetricsClient) CreateSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client MetricsClient) CreateResponder(resp *http.Response) (result AzureMetricsResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
