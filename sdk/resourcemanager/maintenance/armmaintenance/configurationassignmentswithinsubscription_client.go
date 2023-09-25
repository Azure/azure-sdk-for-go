//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmaintenance

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ConfigurationAssignmentsWithinSubscriptionClient contains the methods for the ConfigurationAssignmentsWithinSubscription group.
// Don't use this type directly, use NewConfigurationAssignmentsWithinSubscriptionClient() instead.
type ConfigurationAssignmentsWithinSubscriptionClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewConfigurationAssignmentsWithinSubscriptionClient creates a new instance of ConfigurationAssignmentsWithinSubscriptionClient with the specified values.
//   - subscriptionID - Subscription credentials that uniquely identify a Microsoft Azure subscription. The subscription ID forms
//     part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewConfigurationAssignmentsWithinSubscriptionClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConfigurationAssignmentsWithinSubscriptionClient, error) {
	cl, err := arm.NewClient(moduleName+".ConfigurationAssignmentsWithinSubscriptionClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ConfigurationAssignmentsWithinSubscriptionClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// NewListPager - Get configuration assignment within a subscription
//
// Generated from API version 2023-04-01
//   - options - ConfigurationAssignmentsWithinSubscriptionClientListOptions contains the optional parameters for the ConfigurationAssignmentsWithinSubscriptionClient.NewListPager
//     method.
func (client *ConfigurationAssignmentsWithinSubscriptionClient) NewListPager(options *ConfigurationAssignmentsWithinSubscriptionClientListOptions) (*runtime.Pager[ConfigurationAssignmentsWithinSubscriptionClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ConfigurationAssignmentsWithinSubscriptionClientListResponse]{
		More: func(page ConfigurationAssignmentsWithinSubscriptionClientListResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ConfigurationAssignmentsWithinSubscriptionClientListResponse) (ConfigurationAssignmentsWithinSubscriptionClientListResponse, error) {
			req, err := client.listCreateRequest(ctx, options)
			if err != nil {
				return ConfigurationAssignmentsWithinSubscriptionClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ConfigurationAssignmentsWithinSubscriptionClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ConfigurationAssignmentsWithinSubscriptionClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ConfigurationAssignmentsWithinSubscriptionClient) listCreateRequest(ctx context.Context, options *ConfigurationAssignmentsWithinSubscriptionClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Maintenance/configurationAssignments"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ConfigurationAssignmentsWithinSubscriptionClient) listHandleResponse(resp *http.Response) (ConfigurationAssignmentsWithinSubscriptionClientListResponse, error) {
	result := ConfigurationAssignmentsWithinSubscriptionClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ListConfigurationAssignmentsResult); err != nil {
		return ConfigurationAssignmentsWithinSubscriptionClientListResponse{}, err
	}
	return result, nil
}

