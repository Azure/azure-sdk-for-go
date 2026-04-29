// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func (client *ServiceClient) Endpoint() string {
	return client.url
}

func (client *ServiceClient) InternalClient() *azcore.Client {
	return client.internal
}

// NewServiceClient creates a new instance of ServiceClient with the specified values.
//   - endpoint - The URL of the service account, container, or blob that is the target of the desired operation.
//   - azClient - azcore.Client is a basic HTTP client. It consists of a pipeline and tracing provider.
func NewServiceClient(endpoint string, azClient *azcore.Client) *ServiceClient {
	client := &ServiceClient{
		internal: azClient,
		url:      endpoint,
	}
	return client
}

// ListContainersSegmentCreateRequest creates the ListContainersSegment request.
func (client *ServiceClient) ListContainersSegmentCreateRequest(ctx context.Context, options *ServiceClientListContainersSegmentOptions) (*policy.Request, error) {
	return client.listContainersSegmentCreateRequest(ctx, options)
}

// ListContainersSegmentHandleResponse handles the ListContainersSegment response.
func (client *ServiceClient) ListContainersSegmentHandleResponse(resp *http.Response) (ServiceClientListContainersSegmentResponse, error) {
	return client.listContainersSegmentHandleResponse(resp)
}

// SubmitBatch - The Batch operation allows multiple API calls to be embedded into a single HTTP request.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2026-04-06
//   - contentLength - The length of the request.
//   - multipartContentType - Required. The value of this header must be multipart/mixed with a batch boundary. Example header
//     value: multipart/mixed; boundary=batch_
//   - body - Initial data
//   - options - ServiceClientSubmitBatchOptions contains the optional parameters for the ServiceClient.SubmitBatch method.
func (client *ServiceClient) SubmitBatch(ctx context.Context, contentLength int64, multipartContentType string, body io.ReadSeekCloser, options *ServiceClientSubmitBatchOptions) (ServiceClientSubmitBatchResponse, error) {
	var err error
	req, err := client.submitBatchCreateRequest(ctx, contentLength, multipartContentType, body, options)
	if err != nil {
		return ServiceClientSubmitBatchResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ServiceClientSubmitBatchResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return ServiceClientSubmitBatchResponse{}, err
	}
	resp, err := client.submitBatchHandleResponse(httpResp)
	return resp, err
}

// submitBatchCreateRequest creates the SubmitBatch request.
func (client *ServiceClient) submitBatchCreateRequest(ctx context.Context, contentLength int64, multipartContentType string, body io.ReadSeekCloser, options *ServiceClientSubmitBatchOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPost, client.url)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("comp", "batch")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = strings.Replace(reqQP.Encode(), "+", "%20", -1)
	runtime.SkipBodyDownload(req)
	req.Raw().Header["Accept"] = []string{"application/xml"}
	req.Raw().Header["Content-Length"] = []string{strconv.FormatInt(contentLength, 10)}
	req.Raw().Header["Content-Type"] = []string{multipartContentType}
	if options != nil && options.RequestID != nil {
		req.Raw().Header["x-ms-client-request-id"] = []string{*options.RequestID}
	}
	req.Raw().Header["x-ms-version"] = []string{defaultServiceClientVersion}
	if err := req.SetBody(body, multipartContentType); err != nil {
		return nil, err
	}
	return req, nil
}

// submitBatchHandleResponse handles the SubmitBatch response.
func (client *ServiceClient) submitBatchHandleResponse(resp *http.Response) (ServiceClientSubmitBatchResponse, error) {
	result := ServiceClientSubmitBatchResponse{Body: resp.Body}
	if val := resp.Header.Get("Content-Type"); val != "" {
		result.ContentType = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	return result, nil
}

// ServiceClientSubmitBatchResponse contains the response from method ServiceClient.SubmitBatch.
type ServiceClientSubmitBatchResponse struct {
	// Body contains the streaming response.
	Body io.ReadCloser

	// ContentType contains the information returned from the Content-Type header response.
	ContentType *string

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

// ServiceClientSubmitBatchOptions contains the optional parameters for the ServiceClient.SubmitBatch method.
type ServiceClientSubmitBatchOptions struct {
	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string

	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://learn.microsoft.com/rest/api/storageservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}
