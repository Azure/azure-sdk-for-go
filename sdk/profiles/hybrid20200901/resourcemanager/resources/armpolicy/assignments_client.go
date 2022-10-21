//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armpolicy

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/profiles/hybrid20200901"
	"net/http"
	"net/url"
	"strings"
)

// AssignmentsClient contains the methods for the PolicyAssignments group.
// Don't use this type directly, use NewAssignmentsClient() instead.
type AssignmentsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewAssignmentsClient creates a new instance of AssignmentsClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewAssignmentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AssignmentsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(hybrid20200901.ModuleName, hybrid20200901.ModuleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &AssignmentsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Create - Policy assignments are inherited by child resources. For example, when you apply a policy to a resource group
// that policy is assigned to all resources in the group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2016-12-01
// scope - The scope of the policy assignment.
// policyAssignmentName - The name of the policy assignment.
// parameters - Parameters for the policy assignment.
// options - AssignmentsClientCreateOptions contains the optional parameters for the AssignmentsClient.Create method.
func (client *AssignmentsClient) Create(ctx context.Context, scope string, policyAssignmentName string, parameters Assignment, options *AssignmentsClientCreateOptions) (AssignmentsClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, scope, policyAssignmentName, parameters, options)
	if err != nil {
		return AssignmentsClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AssignmentsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusCreated) {
		return AssignmentsClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *AssignmentsClient) createCreateRequest(ctx context.Context, scope string, policyAssignmentName string, parameters Assignment, options *AssignmentsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/policyAssignments/{policyAssignmentName}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if policyAssignmentName == "" {
		return nil, errors.New("parameter policyAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyAssignmentName}", url.PathEscape(policyAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json, text/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createHandleResponse handles the Create response.
func (client *AssignmentsClient) createHandleResponse(resp *http.Response) (AssignmentsClientCreateResponse, error) {
	result := AssignmentsClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Assignment); err != nil {
		return AssignmentsClientCreateResponse{}, err
	}
	return result, nil
}

// CreateByID - Policy assignments are inherited by child resources. For example, when you apply a policy to a resource group
// that policy is assigned to all resources in the group. When providing a scope for the
// assignment, use '/subscriptions/{subscription-id}/' for subscriptions, '/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}'
// for resource groups, and
// '/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/{resource-provider-namespace}/{resource-type}/{resource-name}'
// for resources.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2016-12-01
// policyAssignmentID - The ID of the policy assignment to create. Use the format '/{scope}/providers/Microsoft.Authorization/policyAssignments/{policy-assignment-name}'.
// parameters - Parameters for policy assignment.
// options - AssignmentsClientCreateByIDOptions contains the optional parameters for the AssignmentsClient.CreateByID method.
func (client *AssignmentsClient) CreateByID(ctx context.Context, policyAssignmentID string, parameters Assignment, options *AssignmentsClientCreateByIDOptions) (AssignmentsClientCreateByIDResponse, error) {
	req, err := client.createByIDCreateRequest(ctx, policyAssignmentID, parameters, options)
	if err != nil {
		return AssignmentsClientCreateByIDResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AssignmentsClientCreateByIDResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusCreated) {
		return AssignmentsClientCreateByIDResponse{}, runtime.NewResponseError(resp)
	}
	return client.createByIDHandleResponse(resp)
}

// createByIDCreateRequest creates the CreateByID request.
func (client *AssignmentsClient) createByIDCreateRequest(ctx context.Context, policyAssignmentID string, parameters Assignment, options *AssignmentsClientCreateByIDOptions) (*policy.Request, error) {
	urlPath := "/{policyAssignmentId}"
	urlPath = strings.ReplaceAll(urlPath, "{policyAssignmentId}", policyAssignmentID)
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json, text/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createByIDHandleResponse handles the CreateByID response.
func (client *AssignmentsClient) createByIDHandleResponse(resp *http.Response) (AssignmentsClientCreateByIDResponse, error) {
	result := AssignmentsClientCreateByIDResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Assignment); err != nil {
		return AssignmentsClientCreateByIDResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a policy assignment.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2016-12-01
// scope - The scope of the policy assignment.
// policyAssignmentName - The name of the policy assignment to delete.
// options - AssignmentsClientDeleteOptions contains the optional parameters for the AssignmentsClient.Delete method.
func (client *AssignmentsClient) Delete(ctx context.Context, scope string, policyAssignmentName string, options *AssignmentsClientDeleteOptions) (AssignmentsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, scope, policyAssignmentName, options)
	if err != nil {
		return AssignmentsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AssignmentsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return AssignmentsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return client.deleteHandleResponse(resp)
}

// deleteCreateRequest creates the Delete request.
func (client *AssignmentsClient) deleteCreateRequest(ctx context.Context, scope string, policyAssignmentName string, options *AssignmentsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/policyAssignments/{policyAssignmentName}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if policyAssignmentName == "" {
		return nil, errors.New("parameter policyAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyAssignmentName}", url.PathEscape(policyAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json, text/json"}
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *AssignmentsClient) deleteHandleResponse(resp *http.Response) (AssignmentsClientDeleteResponse, error) {
	result := AssignmentsClientDeleteResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Assignment); err != nil {
		return AssignmentsClientDeleteResponse{}, err
	}
	return result, nil
}

// DeleteByID - When providing a scope for the assignment, use '/subscriptions/{subscription-id}/' for subscriptions, '/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}'
// for resource groups, and
// '/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/{resource-provider-namespace}/{resource-type}/{resource-name}'
// for resources.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2016-12-01
// policyAssignmentID - The ID of the policy assignment to delete. Use the format '/{scope}/providers/Microsoft.Authorization/policyAssignments/{policy-assignment-name}'.
// options - AssignmentsClientDeleteByIDOptions contains the optional parameters for the AssignmentsClient.DeleteByID method.
func (client *AssignmentsClient) DeleteByID(ctx context.Context, policyAssignmentID string, options *AssignmentsClientDeleteByIDOptions) (AssignmentsClientDeleteByIDResponse, error) {
	req, err := client.deleteByIDCreateRequest(ctx, policyAssignmentID, options)
	if err != nil {
		return AssignmentsClientDeleteByIDResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AssignmentsClientDeleteByIDResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AssignmentsClientDeleteByIDResponse{}, runtime.NewResponseError(resp)
	}
	return client.deleteByIDHandleResponse(resp)
}

// deleteByIDCreateRequest creates the DeleteByID request.
func (client *AssignmentsClient) deleteByIDCreateRequest(ctx context.Context, policyAssignmentID string, options *AssignmentsClientDeleteByIDOptions) (*policy.Request, error) {
	urlPath := "/{policyAssignmentId}"
	urlPath = strings.ReplaceAll(urlPath, "{policyAssignmentId}", policyAssignmentID)
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json, text/json"}
	return req, nil
}

// deleteByIDHandleResponse handles the DeleteByID response.
func (client *AssignmentsClient) deleteByIDHandleResponse(resp *http.Response) (AssignmentsClientDeleteByIDResponse, error) {
	result := AssignmentsClientDeleteByIDResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Assignment); err != nil {
		return AssignmentsClientDeleteByIDResponse{}, err
	}
	return result, nil
}

// Get - Gets a policy assignment.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2016-12-01
// scope - The scope of the policy assignment.
// policyAssignmentName - The name of the policy assignment to get.
// options - AssignmentsClientGetOptions contains the optional parameters for the AssignmentsClient.Get method.
func (client *AssignmentsClient) Get(ctx context.Context, scope string, policyAssignmentName string, options *AssignmentsClientGetOptions) (AssignmentsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, scope, policyAssignmentName, options)
	if err != nil {
		return AssignmentsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AssignmentsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AssignmentsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *AssignmentsClient) getCreateRequest(ctx context.Context, scope string, policyAssignmentName string, options *AssignmentsClientGetOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/policyAssignments/{policyAssignmentName}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if policyAssignmentName == "" {
		return nil, errors.New("parameter policyAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyAssignmentName}", url.PathEscape(policyAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json, text/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AssignmentsClient) getHandleResponse(resp *http.Response) (AssignmentsClientGetResponse, error) {
	result := AssignmentsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Assignment); err != nil {
		return AssignmentsClientGetResponse{}, err
	}
	return result, nil
}

// GetByID - When providing a scope for the assignment, use '/subscriptions/{subscription-id}/' for subscriptions, '/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}'
// for resource groups, and
// '/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/{resource-provider-namespace}/{resource-type}/{resource-name}'
// for resources.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2016-12-01
// policyAssignmentID - The ID of the policy assignment to get. Use the format '/{scope}/providers/Microsoft.Authorization/policyAssignments/{policy-assignment-name}'.
// options - AssignmentsClientGetByIDOptions contains the optional parameters for the AssignmentsClient.GetByID method.
func (client *AssignmentsClient) GetByID(ctx context.Context, policyAssignmentID string, options *AssignmentsClientGetByIDOptions) (AssignmentsClientGetByIDResponse, error) {
	req, err := client.getByIDCreateRequest(ctx, policyAssignmentID, options)
	if err != nil {
		return AssignmentsClientGetByIDResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AssignmentsClientGetByIDResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AssignmentsClientGetByIDResponse{}, runtime.NewResponseError(resp)
	}
	return client.getByIDHandleResponse(resp)
}

// getByIDCreateRequest creates the GetByID request.
func (client *AssignmentsClient) getByIDCreateRequest(ctx context.Context, policyAssignmentID string, options *AssignmentsClientGetByIDOptions) (*policy.Request, error) {
	urlPath := "/{policyAssignmentId}"
	urlPath = strings.ReplaceAll(urlPath, "{policyAssignmentId}", policyAssignmentID)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json, text/json"}
	return req, nil
}

// getByIDHandleResponse handles the GetByID response.
func (client *AssignmentsClient) getByIDHandleResponse(resp *http.Response) (AssignmentsClientGetByIDResponse, error) {
	result := AssignmentsClientGetByIDResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Assignment); err != nil {
		return AssignmentsClientGetByIDResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets all the policy assignments for a subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2016-12-01
// options - AssignmentsClientListOptions contains the optional parameters for the AssignmentsClient.List method.
func (client *AssignmentsClient) NewListPager(options *AssignmentsClientListOptions) *runtime.Pager[AssignmentsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[AssignmentsClientListResponse]{
		More: func(page AssignmentsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AssignmentsClientListResponse) (AssignmentsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AssignmentsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return AssignmentsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AssignmentsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *AssignmentsClient) listCreateRequest(ctx context.Context, options *AssignmentsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Authorization/policyAssignments"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	reqQP.Set("api-version", "2016-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json, text/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *AssignmentsClient) listHandleResponse(resp *http.Response) (AssignmentsClientListResponse, error) {
	result := AssignmentsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AssignmentListResult); err != nil {
		return AssignmentsClientListResponse{}, err
	}
	return result, nil
}

// NewListForResourcePager - Gets policy assignments for a resource.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2016-12-01
// resourceGroupName - The name of the resource group containing the resource. The name is case insensitive.
// resourceProviderNamespace - The namespace of the resource provider.
// parentResourcePath - The parent resource path.
// resourceType - The resource type.
// resourceName - The name of the resource with policy assignments.
// options - AssignmentsClientListForResourceOptions contains the optional parameters for the AssignmentsClient.ListForResource
// method.
func (client *AssignmentsClient) NewListForResourcePager(resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, options *AssignmentsClientListForResourceOptions) *runtime.Pager[AssignmentsClientListForResourceResponse] {
	return runtime.NewPager(runtime.PagingHandler[AssignmentsClientListForResourceResponse]{
		More: func(page AssignmentsClientListForResourceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AssignmentsClientListForResourceResponse) (AssignmentsClientListForResourceResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listForResourceCreateRequest(ctx, resourceGroupName, resourceProviderNamespace, parentResourcePath, resourceType, resourceName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AssignmentsClientListForResourceResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return AssignmentsClientListForResourceResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AssignmentsClientListForResourceResponse{}, runtime.NewResponseError(resp)
			}
			return client.listForResourceHandleResponse(resp)
		},
	})
}

// listForResourceCreateRequest creates the ListForResource request.
func (client *AssignmentsClient) listForResourceCreateRequest(ctx context.Context, resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, options *AssignmentsClientListForResourceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{parentResourcePath}/{resourceType}/{resourceName}/providers/Microsoft.Authorization/policyAssignments"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceProviderNamespace == "" {
		return nil, errors.New("parameter resourceProviderNamespace cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceProviderNamespace}", url.PathEscape(resourceProviderNamespace))
	urlPath = strings.ReplaceAll(urlPath, "{parentResourcePath}", parentResourcePath)
	urlPath = strings.ReplaceAll(urlPath, "{resourceType}", resourceType)
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	reqQP.Set("api-version", "2016-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json, text/json"}
	return req, nil
}

// listForResourceHandleResponse handles the ListForResource response.
func (client *AssignmentsClient) listForResourceHandleResponse(resp *http.Response) (AssignmentsClientListForResourceResponse, error) {
	result := AssignmentsClientListForResourceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AssignmentListResult); err != nil {
		return AssignmentsClientListForResourceResponse{}, err
	}
	return result, nil
}

// NewListForResourceGroupPager - Gets policy assignments for the resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2016-12-01
// resourceGroupName - The name of the resource group that contains policy assignments.
// options - AssignmentsClientListForResourceGroupOptions contains the optional parameters for the AssignmentsClient.ListForResourceGroup
// method.
func (client *AssignmentsClient) NewListForResourceGroupPager(resourceGroupName string, options *AssignmentsClientListForResourceGroupOptions) *runtime.Pager[AssignmentsClientListForResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[AssignmentsClientListForResourceGroupResponse]{
		More: func(page AssignmentsClientListForResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AssignmentsClientListForResourceGroupResponse) (AssignmentsClientListForResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listForResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AssignmentsClientListForResourceGroupResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return AssignmentsClientListForResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AssignmentsClientListForResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listForResourceGroupHandleResponse(resp)
		},
	})
}

// listForResourceGroupCreateRequest creates the ListForResourceGroup request.
func (client *AssignmentsClient) listForResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *AssignmentsClientListForResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Authorization/policyAssignments"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	unencodedParams := []string{req.Raw().URL.RawQuery}
	if options != nil && options.Filter != nil {
		unencodedParams = append(unencodedParams, "$filter="+*options.Filter)
	}
	req.Raw().URL.RawQuery = strings.Join(unencodedParams, "&")
	req.Raw().Header["Accept"] = []string{"application/json, text/json"}
	return req, nil
}

// listForResourceGroupHandleResponse handles the ListForResourceGroup response.
func (client *AssignmentsClient) listForResourceGroupHandleResponse(resp *http.Response) (AssignmentsClientListForResourceGroupResponse, error) {
	result := AssignmentsClientListForResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AssignmentListResult); err != nil {
		return AssignmentsClientListForResourceGroupResponse{}, err
	}
	return result, nil
}
