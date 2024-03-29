//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

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

// AssessmentsMetadataClient contains the methods for the AssessmentsMetadata group.
// Don't use this type directly, use NewAssessmentsMetadataClient() instead.
type AssessmentsMetadataClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewAssessmentsMetadataClient creates a new instance of AssessmentsMetadataClient with the specified values.
//   - subscriptionID - Azure subscription ID
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAssessmentsMetadataClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AssessmentsMetadataClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AssessmentsMetadataClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateInSubscription - Create metadata information on an assessment type in a specific subscription
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - assessmentMetadataName - The Assessment Key - Unique key for the assessment type
//   - assessmentMetadata - AssessmentMetadata object
//   - options - AssessmentsMetadataClientCreateInSubscriptionOptions contains the optional parameters for the AssessmentsMetadataClient.CreateInSubscription
//     method.
func (client *AssessmentsMetadataClient) CreateInSubscription(ctx context.Context, assessmentMetadataName string, assessmentMetadata AssessmentMetadataResponse, options *AssessmentsMetadataClientCreateInSubscriptionOptions) (AssessmentsMetadataClientCreateInSubscriptionResponse, error) {
	var err error
	const operationName = "AssessmentsMetadataClient.CreateInSubscription"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createInSubscriptionCreateRequest(ctx, assessmentMetadataName, assessmentMetadata, options)
	if err != nil {
		return AssessmentsMetadataClientCreateInSubscriptionResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AssessmentsMetadataClientCreateInSubscriptionResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AssessmentsMetadataClientCreateInSubscriptionResponse{}, err
	}
	resp, err := client.createInSubscriptionHandleResponse(httpResp)
	return resp, err
}

// createInSubscriptionCreateRequest creates the CreateInSubscription request.
func (client *AssessmentsMetadataClient) createInSubscriptionCreateRequest(ctx context.Context, assessmentMetadataName string, assessmentMetadata AssessmentMetadataResponse, options *AssessmentsMetadataClientCreateInSubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/assessmentMetadata/{assessmentMetadataName}"
	if assessmentMetadataName == "" {
		return nil, errors.New("parameter assessmentMetadataName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{assessmentMetadataName}", url.PathEscape(assessmentMetadataName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, assessmentMetadata); err != nil {
		return nil, err
	}
	return req, nil
}

// createInSubscriptionHandleResponse handles the CreateInSubscription response.
func (client *AssessmentsMetadataClient) createInSubscriptionHandleResponse(resp *http.Response) (AssessmentsMetadataClientCreateInSubscriptionResponse, error) {
	result := AssessmentsMetadataClientCreateInSubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AssessmentMetadataResponse); err != nil {
		return AssessmentsMetadataClientCreateInSubscriptionResponse{}, err
	}
	return result, nil
}

// DeleteInSubscription - Delete metadata information on an assessment type in a specific subscription, will cause the deletion
// of all the assessments of that type in that subscription
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - assessmentMetadataName - The Assessment Key - Unique key for the assessment type
//   - options - AssessmentsMetadataClientDeleteInSubscriptionOptions contains the optional parameters for the AssessmentsMetadataClient.DeleteInSubscription
//     method.
func (client *AssessmentsMetadataClient) DeleteInSubscription(ctx context.Context, assessmentMetadataName string, options *AssessmentsMetadataClientDeleteInSubscriptionOptions) (AssessmentsMetadataClientDeleteInSubscriptionResponse, error) {
	var err error
	const operationName = "AssessmentsMetadataClient.DeleteInSubscription"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteInSubscriptionCreateRequest(ctx, assessmentMetadataName, options)
	if err != nil {
		return AssessmentsMetadataClientDeleteInSubscriptionResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AssessmentsMetadataClientDeleteInSubscriptionResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AssessmentsMetadataClientDeleteInSubscriptionResponse{}, err
	}
	return AssessmentsMetadataClientDeleteInSubscriptionResponse{}, nil
}

// deleteInSubscriptionCreateRequest creates the DeleteInSubscription request.
func (client *AssessmentsMetadataClient) deleteInSubscriptionCreateRequest(ctx context.Context, assessmentMetadataName string, options *AssessmentsMetadataClientDeleteInSubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/assessmentMetadata/{assessmentMetadataName}"
	if assessmentMetadataName == "" {
		return nil, errors.New("parameter assessmentMetadataName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{assessmentMetadataName}", url.PathEscape(assessmentMetadataName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get metadata information on an assessment type
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - assessmentMetadataName - The Assessment Key - Unique key for the assessment type
//   - options - AssessmentsMetadataClientGetOptions contains the optional parameters for the AssessmentsMetadataClient.Get method.
func (client *AssessmentsMetadataClient) Get(ctx context.Context, assessmentMetadataName string, options *AssessmentsMetadataClientGetOptions) (AssessmentsMetadataClientGetResponse, error) {
	var err error
	const operationName = "AssessmentsMetadataClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, assessmentMetadataName, options)
	if err != nil {
		return AssessmentsMetadataClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AssessmentsMetadataClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AssessmentsMetadataClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *AssessmentsMetadataClient) getCreateRequest(ctx context.Context, assessmentMetadataName string, options *AssessmentsMetadataClientGetOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Security/assessmentMetadata/{assessmentMetadataName}"
	if assessmentMetadataName == "" {
		return nil, errors.New("parameter assessmentMetadataName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{assessmentMetadataName}", url.PathEscape(assessmentMetadataName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AssessmentsMetadataClient) getHandleResponse(resp *http.Response) (AssessmentsMetadataClientGetResponse, error) {
	result := AssessmentsMetadataClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AssessmentMetadataResponse); err != nil {
		return AssessmentsMetadataClientGetResponse{}, err
	}
	return result, nil
}

// GetInSubscription - Get metadata information on an assessment type in a specific subscription
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - assessmentMetadataName - The Assessment Key - Unique key for the assessment type
//   - options - AssessmentsMetadataClientGetInSubscriptionOptions contains the optional parameters for the AssessmentsMetadataClient.GetInSubscription
//     method.
func (client *AssessmentsMetadataClient) GetInSubscription(ctx context.Context, assessmentMetadataName string, options *AssessmentsMetadataClientGetInSubscriptionOptions) (AssessmentsMetadataClientGetInSubscriptionResponse, error) {
	var err error
	const operationName = "AssessmentsMetadataClient.GetInSubscription"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getInSubscriptionCreateRequest(ctx, assessmentMetadataName, options)
	if err != nil {
		return AssessmentsMetadataClientGetInSubscriptionResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AssessmentsMetadataClientGetInSubscriptionResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AssessmentsMetadataClientGetInSubscriptionResponse{}, err
	}
	resp, err := client.getInSubscriptionHandleResponse(httpResp)
	return resp, err
}

// getInSubscriptionCreateRequest creates the GetInSubscription request.
func (client *AssessmentsMetadataClient) getInSubscriptionCreateRequest(ctx context.Context, assessmentMetadataName string, options *AssessmentsMetadataClientGetInSubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/assessmentMetadata/{assessmentMetadataName}"
	if assessmentMetadataName == "" {
		return nil, errors.New("parameter assessmentMetadataName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{assessmentMetadataName}", url.PathEscape(assessmentMetadataName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getInSubscriptionHandleResponse handles the GetInSubscription response.
func (client *AssessmentsMetadataClient) getInSubscriptionHandleResponse(resp *http.Response) (AssessmentsMetadataClientGetInSubscriptionResponse, error) {
	result := AssessmentsMetadataClientGetInSubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AssessmentMetadataResponse); err != nil {
		return AssessmentsMetadataClientGetInSubscriptionResponse{}, err
	}
	return result, nil
}

// NewListPager - Get metadata information on all assessment types
//
// Generated from API version 2021-06-01
//   - options - AssessmentsMetadataClientListOptions contains the optional parameters for the AssessmentsMetadataClient.NewListPager
//     method.
func (client *AssessmentsMetadataClient) NewListPager(options *AssessmentsMetadataClientListOptions) *runtime.Pager[AssessmentsMetadataClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[AssessmentsMetadataClientListResponse]{
		More: func(page AssessmentsMetadataClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AssessmentsMetadataClientListResponse) (AssessmentsMetadataClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "AssessmentsMetadataClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return AssessmentsMetadataClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *AssessmentsMetadataClient) listCreateRequest(ctx context.Context, options *AssessmentsMetadataClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Security/assessmentMetadata"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *AssessmentsMetadataClient) listHandleResponse(resp *http.Response) (AssessmentsMetadataClientListResponse, error) {
	result := AssessmentsMetadataClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AssessmentMetadataResponseList); err != nil {
		return AssessmentsMetadataClientListResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Get metadata information on all assessment types in a specific subscription
//
// Generated from API version 2021-06-01
//   - options - AssessmentsMetadataClientListBySubscriptionOptions contains the optional parameters for the AssessmentsMetadataClient.NewListBySubscriptionPager
//     method.
func (client *AssessmentsMetadataClient) NewListBySubscriptionPager(options *AssessmentsMetadataClientListBySubscriptionOptions) *runtime.Pager[AssessmentsMetadataClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[AssessmentsMetadataClientListBySubscriptionResponse]{
		More: func(page AssessmentsMetadataClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AssessmentsMetadataClientListBySubscriptionResponse) (AssessmentsMetadataClientListBySubscriptionResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "AssessmentsMetadataClient.NewListBySubscriptionPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listBySubscriptionCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return AssessmentsMetadataClientListBySubscriptionResponse{}, err
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *AssessmentsMetadataClient) listBySubscriptionCreateRequest(ctx context.Context, options *AssessmentsMetadataClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/assessmentMetadata"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *AssessmentsMetadataClient) listBySubscriptionHandleResponse(resp *http.Response) (AssessmentsMetadataClientListBySubscriptionResponse, error) {
	result := AssessmentsMetadataClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AssessmentMetadataResponseList); err != nil {
		return AssessmentsMetadataClientListBySubscriptionResponse{}, err
	}
	return result, nil
}
