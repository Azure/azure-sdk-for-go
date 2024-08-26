//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package generated

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ServiceClient contains the methods for the Service group.
// Don't use this type directly, use NewServiceClient() instead.
type ServiceClient struct {
	endpoint string
	pl       runtime.Pipeline
}

// NewServiceClient creates a new instance of ServiceClient with the specified values.
//   - endpoint - The URL of the service account, queue or message that is the target of the desired operation.
//   - pl - the pipeline used for sending requests and handling responses.
func NewServiceClient(endpoint string, pl runtime.Pipeline) *ServiceClient {
	client := &ServiceClient{
		endpoint: endpoint,
		pl:       pl,
	}
	return client
}

// GetProperties - gets the properties of a storage account's Queue service, including properties for Storage Analytics and
// CORS (Cross-Origin Resource Sharing) rules.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-03-28
//   - options - ServiceClientGetPropertiesOptions contains the optional parameters for the ServiceClient.GetProperties method.
func (client *ServiceClient) GetProperties(ctx context.Context, options *ServiceClientGetPropertiesOptions) (ServiceClientGetPropertiesResponse, error) {
	req, err := client.getPropertiesCreateRequest(ctx, options)
	if err != nil {
		return ServiceClientGetPropertiesResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ServiceClientGetPropertiesResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ServiceClientGetPropertiesResponse{}, runtime.NewResponseError(resp)
	}
	return client.getPropertiesHandleResponse(resp)
}

// getPropertiesCreateRequest creates the GetProperties request.
func (client *ServiceClient) getPropertiesCreateRequest(ctx context.Context, options *ServiceClientGetPropertiesOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "service")
	reqQP.Set("comp", "properties")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["x-ms-version"] = []string{"2024-08-04"}
	if options != nil && options.RequestID != nil {
		req.Raw().Header["x-ms-client-request-id"] = []string{*options.RequestID}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// getPropertiesHandleResponse handles the GetProperties response.
func (client *ServiceClient) getPropertiesHandleResponse(resp *http.Response) (ServiceClientGetPropertiesResponse, error) {
	result := ServiceClientGetPropertiesResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if err := runtime.UnmarshalAsXML(resp, &result.StorageServiceProperties); err != nil {
		return ServiceClientGetPropertiesResponse{}, err
	}
	return result, nil
}

// GetStatistics - Retrieves statistics related to replication for the Queue service. It is only available on the secondary
// location endpoint when read-access geo-redundant replication is enabled for the storage
// account.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-03-28
//   - options - ServiceClientGetStatisticsOptions contains the optional parameters for the ServiceClient.GetStatistics method.
func (client *ServiceClient) GetStatistics(ctx context.Context, options *ServiceClientGetStatisticsOptions) (ServiceClientGetStatisticsResponse, error) {
	req, err := client.getStatisticsCreateRequest(ctx, options)
	if err != nil {
		return ServiceClientGetStatisticsResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ServiceClientGetStatisticsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ServiceClientGetStatisticsResponse{}, runtime.NewResponseError(resp)
	}
	return client.getStatisticsHandleResponse(resp)
}

// getStatisticsCreateRequest creates the GetStatistics request.
func (client *ServiceClient) getStatisticsCreateRequest(ctx context.Context, options *ServiceClientGetStatisticsOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "service")
	reqQP.Set("comp", "stats")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["x-ms-version"] = []string{"2024-08-04"}
	if options != nil && options.RequestID != nil {
		req.Raw().Header["x-ms-client-request-id"] = []string{*options.RequestID}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// getStatisticsHandleResponse handles the GetStatistics response.
func (client *ServiceClient) getStatisticsHandleResponse(resp *http.Response) (ServiceClientGetStatisticsResponse, error) {
	result := ServiceClientGetStatisticsResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return ServiceClientGetStatisticsResponse{}, err
		}
		result.Date = &date
	}
	if err := runtime.UnmarshalAsXML(resp, &result.StorageServiceStats); err != nil {
		return ServiceClientGetStatisticsResponse{}, err
	}
	return result, nil
}

// NewListQueuesSegmentPager - The List Queues Segment operation returns a list of the queues under the specified account
//
// Generated from API version 2018-03-28
//   - options - ServiceClientListQueuesSegmentOptions contains the optional parameters for the ServiceClient.NewListQueuesSegmentPager
//     method.
//
// ListQueuesSegmentCreateRequest creates the ListQueuesFlatSegment ListQueuesSegment.
func (client *ServiceClient) ListQueuesSegmentCreateRequest(ctx context.Context, options *ServiceClientListQueuesSegmentOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("comp", "list")
	if options != nil && options.Prefix != nil {
		reqQP.Set("prefix", *options.Prefix)
	}
	if options != nil && options.Marker != nil {
		reqQP.Set("marker", *options.Marker)
	}
	if options != nil && options.Maxresults != nil {
		reqQP.Set("maxresults", strconv.FormatInt(int64(*options.Maxresults), 10))
	}
	if options != nil && options.Include != nil {
		reqQP.Set("include", strings.Join(strings.Fields(strings.Trim(fmt.Sprint(options.Include), "[]")), ","))
	}
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["x-ms-version"] = []string{"2024-08-04"}
	if options != nil && options.RequestID != nil {
		req.Raw().Header["x-ms-client-request-id"] = []string{*options.RequestID}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, nil
}

// listQueuesSegmentHandleResponse handles the ListQueuesSegment response.
func (client *ServiceClient) ListQueuesSegmentHandleResponse(resp *http.Response) (ServiceClientListQueuesSegmentResponse, error) {
	result := ServiceClientListQueuesSegmentResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return ServiceClientListQueuesSegmentResponse{}, err
		}
		result.Date = &date
	}
	if err := runtime.UnmarshalAsXML(resp, &result.ListQueuesSegmentResponse); err != nil {
		return ServiceClientListQueuesSegmentResponse{}, err
	}
	return result, nil
}

// SetProperties - Sets properties for a storage account's Queue service endpoint, including properties for Storage Analytics
// and CORS (Cross-Origin Resource Sharing) rules
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-03-28
//   - storageServiceProperties - The StorageService properties.
//   - options - ServiceClientSetPropertiesOptions contains the optional parameters for the ServiceClient.SetProperties method.
func (client *ServiceClient) SetProperties(ctx context.Context, storageServiceProperties StorageServiceProperties, options *ServiceClientSetPropertiesOptions) (ServiceClientSetPropertiesResponse, error) {
	req, err := client.setPropertiesCreateRequest(ctx, storageServiceProperties, options)
	if err != nil {
		return ServiceClientSetPropertiesResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ServiceClientSetPropertiesResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted) {
		return ServiceClientSetPropertiesResponse{}, runtime.NewResponseError(resp)
	}
	return client.setPropertiesHandleResponse(resp)
}

// setPropertiesCreateRequest creates the SetProperties request.
func (client *ServiceClient) setPropertiesCreateRequest(ctx context.Context, storageServiceProperties StorageServiceProperties, options *ServiceClientSetPropertiesOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPut, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "service")
	reqQP.Set("comp", "properties")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["x-ms-version"] = []string{"2024-08-04"}
	if options != nil && options.RequestID != nil {
		req.Raw().Header["x-ms-client-request-id"] = []string{*options.RequestID}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, runtime.MarshalAsXML(req, storageServiceProperties)
}

// setPropertiesHandleResponse handles the SetProperties response.
func (client *ServiceClient) setPropertiesHandleResponse(resp *http.Response) (ServiceClientSetPropertiesResponse, error) {
	result := ServiceClientSetPropertiesResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	return result, nil
}
