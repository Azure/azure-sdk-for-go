// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmonitor

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
)

// MetricBaselineOperations contains the methods for the MetricBaseline group.
type MetricBaselineOperations interface {
	// CalculateBaseline - **Lists the baseline values for a resource**.
	CalculateBaseline(ctx context.Context, resourceUri string, timeSeriesInformation TimeSeriesInformation, options *MetricBaselineCalculateBaselineOptions) (*CalculateBaselineResponseResponse, error)
	// Get - **Gets the baseline values for a specific metric**.
	Get(ctx context.Context, resourceUri string, metricName string, options *MetricBaselineGetOptions) (*BaselineResponseResponse, error)
}

// MetricBaselineClient implements the MetricBaselineOperations interface.
// Don't use this type directly, use NewMetricBaselineClient() instead.
type MetricBaselineClient struct {
	*Client
}

// NewMetricBaselineClient creates a new instance of MetricBaselineClient with the specified values.
func NewMetricBaselineClient(c *Client) MetricBaselineOperations {
	return &MetricBaselineClient{Client: c}
}

// Do invokes the Do() method on the pipeline associated with this client.
func (client *MetricBaselineClient) Do(req *azcore.Request) (*azcore.Response, error) {
	return client.p.Do(req)
}

// CalculateBaseline - **Lists the baseline values for a resource**.
func (client *MetricBaselineClient) CalculateBaseline(ctx context.Context, resourceUri string, timeSeriesInformation TimeSeriesInformation, options *MetricBaselineCalculateBaselineOptions) (*CalculateBaselineResponseResponse, error) {
	req, err := client.CalculateBaselineCreateRequest(ctx, resourceUri, timeSeriesInformation, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.CalculateBaselineHandleError(resp)
	}
	result, err := client.CalculateBaselineHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CalculateBaselineCreateRequest creates the CalculateBaseline request.
func (client *MetricBaselineClient) CalculateBaselineCreateRequest(ctx context.Context, resourceUri string, timeSeriesInformation TimeSeriesInformation, options *MetricBaselineCalculateBaselineOptions) (*azcore.Request, error) {
	urlPath := "/{resourceUri}/providers/microsoft.insights/calculatebaseline"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceUri)
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2017-11-01-preview")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(timeSeriesInformation)
}

// CalculateBaselineHandleResponse handles the CalculateBaseline response.
func (client *MetricBaselineClient) CalculateBaselineHandleResponse(resp *azcore.Response) (*CalculateBaselineResponseResponse, error) {
	result := CalculateBaselineResponseResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.CalculateBaselineResponse)
}

// CalculateBaselineHandleError handles the CalculateBaseline error response.
func (client *MetricBaselineClient) CalculateBaselineHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// Get - **Gets the baseline values for a specific metric**.
func (client *MetricBaselineClient) Get(ctx context.Context, resourceUri string, metricName string, options *MetricBaselineGetOptions) (*BaselineResponseResponse, error) {
	req, err := client.GetCreateRequest(ctx, resourceUri, metricName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.GetHandleError(resp)
	}
	result, err := client.GetHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetCreateRequest creates the Get request.
func (client *MetricBaselineClient) GetCreateRequest(ctx context.Context, resourceUri string, metricName string, options *MetricBaselineGetOptions) (*azcore.Request, error) {
	urlPath := "/{resourceUri}/providers/microsoft.insights/baseline/{metricName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceUri)
	urlPath = strings.ReplaceAll(urlPath, "{metricName}", url.PathEscape(metricName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if options != nil && options.Timespan != nil {
		query.Set("timespan", *options.Timespan)
	}
	if options != nil && options.Interval != nil {
		query.Set("interval", *options.Interval)
	}
	if options != nil && options.Aggregation != nil {
		query.Set("aggregation", *options.Aggregation)
	}
	if options != nil && options.Sensitivities != nil {
		query.Set("sensitivities", *options.Sensitivities)
	}
	if options != nil && options.ResultType != nil {
		query.Set("resultType", string(*options.ResultType))
	}
	query.Set("api-version", "2017-11-01-preview")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// GetHandleResponse handles the Get response.
func (client *MetricBaselineClient) GetHandleResponse(resp *azcore.Response) (*BaselineResponseResponse, error) {
	result := BaselineResponseResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.BaselineResponse)
}

// GetHandleError handles the Get error response.
func (client *MetricBaselineClient) GetHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}
