// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmigrationassessment

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

// AksOptionsOperationsClient contains the methods for the AksOptionsOperations group.
// Don't use this type directly, use NewAksOptionsOperationsClient() instead.
type AksOptionsOperationsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewAksOptionsOperationsClient creates a new instance of AksOptionsOperationsClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAksOptionsOperationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AksOptionsOperationsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AksOptionsOperationsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Get a AKSAssessmentOptions
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - projectName - Assessment Project Name
//   - assessmentOptionsName - AKS Assessment Options Name.
//   - options - AksOptionsOperationsClientGetOptions contains the optional parameters for the AksOptionsOperationsClient.Get
//     method.
func (client *AksOptionsOperationsClient) Get(ctx context.Context, resourceGroupName string, projectName string, assessmentOptionsName string, options *AksOptionsOperationsClientGetOptions) (AksOptionsOperationsClientGetResponse, error) {
	var err error
	const operationName = "AksOptionsOperationsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, projectName, assessmentOptionsName, options)
	if err != nil {
		return AksOptionsOperationsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AksOptionsOperationsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AksOptionsOperationsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *AksOptionsOperationsClient) getCreateRequest(ctx context.Context, resourceGroupName string, projectName string, assessmentOptionsName string, _ *AksOptionsOperationsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/aksAssessmentOptions/{assessmentOptionsName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if projectName == "" {
		return nil, errors.New("parameter projectName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{projectName}", url.PathEscape(projectName))
	if assessmentOptionsName == "" {
		return nil, errors.New("parameter assessmentOptionsName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{assessmentOptionsName}", url.PathEscape(assessmentOptionsName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AksOptionsOperationsClient) getHandleResponse(resp *http.Response) (AksOptionsOperationsClientGetResponse, error) {
	result := AksOptionsOperationsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AKSAssessmentOptions); err != nil {
		return AksOptionsOperationsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByAssessmentProjectPager - List AKSAssessmentOptions resources by AssessmentProject
//
// Generated from API version 2024-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - projectName - Assessment Project Name
//   - options - AksOptionsOperationsClientListByAssessmentProjectOptions contains the optional parameters for the AksOptionsOperationsClient.NewListByAssessmentProjectPager
//     method.
func (client *AksOptionsOperationsClient) NewListByAssessmentProjectPager(resourceGroupName string, projectName string, options *AksOptionsOperationsClientListByAssessmentProjectOptions) *runtime.Pager[AksOptionsOperationsClientListByAssessmentProjectResponse] {
	return runtime.NewPager(runtime.PagingHandler[AksOptionsOperationsClientListByAssessmentProjectResponse]{
		More: func(page AksOptionsOperationsClientListByAssessmentProjectResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AksOptionsOperationsClientListByAssessmentProjectResponse) (AksOptionsOperationsClientListByAssessmentProjectResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "AksOptionsOperationsClient.NewListByAssessmentProjectPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByAssessmentProjectCreateRequest(ctx, resourceGroupName, projectName, options)
			}, nil)
			if err != nil {
				return AksOptionsOperationsClientListByAssessmentProjectResponse{}, err
			}
			return client.listByAssessmentProjectHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByAssessmentProjectCreateRequest creates the ListByAssessmentProject request.
func (client *AksOptionsOperationsClient) listByAssessmentProjectCreateRequest(ctx context.Context, resourceGroupName string, projectName string, _ *AksOptionsOperationsClientListByAssessmentProjectOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/aksAssessmentOptions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if projectName == "" {
		return nil, errors.New("parameter projectName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{projectName}", url.PathEscape(projectName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByAssessmentProjectHandleResponse handles the ListByAssessmentProject response.
func (client *AksOptionsOperationsClient) listByAssessmentProjectHandleResponse(resp *http.Response) (AksOptionsOperationsClientListByAssessmentProjectResponse, error) {
	result := AksOptionsOperationsClientListByAssessmentProjectResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AKSAssessmentOptionsListResult); err != nil {
		return AksOptionsOperationsClientListByAssessmentProjectResponse{}, err
	}
	return result, nil
}
