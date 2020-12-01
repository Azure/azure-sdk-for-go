// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
)

// AvailableDelegationsClient contains the methods for the AvailableDelegations group.
// Don't use this type directly, use NewAvailableDelegationsClient() instead.
type AvailableDelegationsClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewAvailableDelegationsClient creates a new instance of AvailableDelegationsClient with the specified values.
func NewAvailableDelegationsClient(con *armcore.Connection, subscriptionID string) AvailableDelegationsClient {
	return AvailableDelegationsClient{con: con, subscriptionID: subscriptionID}
}

// Pipeline returns the pipeline associated with this client.
func (client AvailableDelegationsClient) Pipeline() azcore.Pipeline {
	return client.con.Pipeline()
}

// List - Gets all of the available subnet delegations for this subscription in this region.
func (client AvailableDelegationsClient) List(location string, options *AvailableDelegationsListOptions) AvailableDelegationsResultPager {
	return &availableDelegationsResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listCreateRequest(ctx, location, options)
		},
		responder: client.listHandleResponse,
		errorer:   client.listHandleError,
		advancer: func(ctx context.Context, resp AvailableDelegationsResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.AvailableDelegationsResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listCreateRequest creates the List request.
func (client AvailableDelegationsClient) listCreateRequest(ctx context.Context, location string, options *AvailableDelegationsListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/locations/{location}/availableDelegations"
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-07-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client AvailableDelegationsClient) listHandleResponse(resp *azcore.Response) (AvailableDelegationsResultResponse, error) {
	result := AvailableDelegationsResultResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.AvailableDelegationsResult)
	return result, err
}

// listHandleError handles the List error response.
func (client AvailableDelegationsClient) listHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}
