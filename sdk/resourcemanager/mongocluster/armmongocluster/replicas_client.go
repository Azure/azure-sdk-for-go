// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armmongocluster

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

// ReplicasClient contains the methods for the Replicas group.
// Don't use this type directly, use NewReplicasClient() instead.
type ReplicasClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewReplicasClient creates a new instance of ReplicasClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewReplicasClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ReplicasClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ReplicasClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// NewListByParentPager - List all the replicas for the mongo cluster.
//
// Generated from API version 2024-07-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - mongoClusterName - The name of the mongo cluster.
//   - options - ReplicasClientListByParentOptions contains the optional parameters for the ReplicasClient.NewListByParentPager
//     method.
func (client *ReplicasClient) NewListByParentPager(resourceGroupName string, mongoClusterName string, options *ReplicasClientListByParentOptions) *runtime.Pager[ReplicasClientListByParentResponse] {
	return runtime.NewPager(runtime.PagingHandler[ReplicasClientListByParentResponse]{
		More: func(page ReplicasClientListByParentResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ReplicasClientListByParentResponse) (ReplicasClientListByParentResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ReplicasClient.NewListByParentPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByParentCreateRequest(ctx, resourceGroupName, mongoClusterName, options)
			}, nil)
			if err != nil {
				return ReplicasClientListByParentResponse{}, err
			}
			return client.listByParentHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByParentCreateRequest creates the ListByParent request.
func (client *ReplicasClient) listByParentCreateRequest(ctx context.Context, resourceGroupName string, mongoClusterName string, _ *ReplicasClientListByParentOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/mongoClusters/{mongoClusterName}/replicas"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if mongoClusterName == "" {
		return nil, errors.New("parameter mongoClusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{mongoClusterName}", url.PathEscape(mongoClusterName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByParentHandleResponse handles the ListByParent response.
func (client *ReplicasClient) listByParentHandleResponse(resp *http.Response) (ReplicasClientListByParentResponse, error) {
	result := ReplicasClientListByParentResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReplicaListResult); err != nil {
		return ReplicasClientListByParentResponse{}, err
	}
	return result, nil
}
