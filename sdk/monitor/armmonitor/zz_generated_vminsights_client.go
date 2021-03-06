// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmonitor

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"strings"
)

// VMInsightsClient contains the methods for the VMInsights group.
// Don't use this type directly, use NewVMInsightsClient() instead.
type VMInsightsClient struct {
	con *armcore.Connection
}

// NewVMInsightsClient creates a new instance of VMInsightsClient with the specified values.
func NewVMInsightsClient(con *armcore.Connection) *VMInsightsClient {
	return &VMInsightsClient{con: con}
}

// GetOnboardingStatus - Retrieves the VM Insights onboarding status for the specified resource or resource scope.
// If the operation fails it returns the *ResponseWithError error type.
func (client *VMInsightsClient) GetOnboardingStatus(ctx context.Context, resourceURI string, options *VMInsightsGetOnboardingStatusOptions) (VMInsightsOnboardingStatusResponse, error) {
	req, err := client.getOnboardingStatusCreateRequest(ctx, resourceURI, options)
	if err != nil {
		return VMInsightsOnboardingStatusResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return VMInsightsOnboardingStatusResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return VMInsightsOnboardingStatusResponse{}, client.getOnboardingStatusHandleError(resp)
	}
	return client.getOnboardingStatusHandleResponse(resp)
}

// getOnboardingStatusCreateRequest creates the GetOnboardingStatus request.
func (client *VMInsightsClient) getOnboardingStatusCreateRequest(ctx context.Context, resourceURI string, options *VMInsightsGetOnboardingStatusOptions) (*azcore.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.Insights/vmInsightsOnboardingStatuses/default"
	if resourceURI == "" {
		return nil, errors.New("parameter resourceURI cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2018-11-27-preview")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getOnboardingStatusHandleResponse handles the GetOnboardingStatus response.
func (client *VMInsightsClient) getOnboardingStatusHandleResponse(resp *azcore.Response) (VMInsightsOnboardingStatusResponse, error) {
	var val *VMInsightsOnboardingStatus
	if err := resp.UnmarshalAsJSON(&val); err != nil {
		return VMInsightsOnboardingStatusResponse{}, err
	}
	return VMInsightsOnboardingStatusResponse{RawResponse: resp.Response, VMInsightsOnboardingStatus: val}, nil
}

// getOnboardingStatusHandleError handles the GetOnboardingStatus error response.
func (client *VMInsightsClient) getOnboardingStatusHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	errType := ResponseWithError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}
