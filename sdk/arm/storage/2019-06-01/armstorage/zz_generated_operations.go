// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"io/ioutil"
	"net/http"
)

// Operations contains the methods for the Operations group.
type Operations interface {
	// List - Lists all of the available Storage Rest API operations.
	List(ctx context.Context, options *OperationsListOptions) (*OperationListResultResponse, error)
}

// OperationsClient implements the Operations interface.
// Don't use this type directly, use NewOperationsClient() instead.
type OperationsClient struct {
	*Client
}

// NewOperationsClient creates a new instance of OperationsClient with the specified values.
func NewOperationsClient(c *Client) Operations {
	return &OperationsClient{Client: c}
}

// Do invokes the Do() method on the pipeline associated with this client.
func (client *OperationsClient) Do(req *azcore.Request) (*azcore.Response, error) {
	return client.p.Do(req)
}

// List - Lists all of the available Storage Rest API operations.
func (client *OperationsClient) List(ctx context.Context, options *OperationsListOptions) (*OperationListResultResponse, error) {
	req, err := client.ListCreateRequest(ctx, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.ListHandleError(resp)
	}
	result, err := client.ListHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListCreateRequest creates the List request.
func (client *OperationsClient) ListCreateRequest(ctx context.Context, options *OperationsListOptions) (*azcore.Request, error) {
	urlPath := "/providers/Microsoft.Storage/operations"
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2019-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListHandleResponse handles the List response.
func (client *OperationsClient) ListHandleResponse(resp *azcore.Response) (*OperationListResultResponse, error) {
	result := OperationListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.OperationListResult)
}

// ListHandleError handles the List error response.
func (client *OperationsClient) ListHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}
