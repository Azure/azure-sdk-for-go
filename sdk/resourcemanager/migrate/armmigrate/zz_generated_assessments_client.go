//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmigrate

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// AssessmentsClient contains the methods for the Assessments group.
// Don't use this type directly, use NewAssessmentsClient() instead.
type AssessmentsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewAssessmentsClient creates a new instance of AssessmentsClient with the specified values.
// subscriptionID - Azure Subscription Id in which project was created.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewAssessmentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AssessmentsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &AssessmentsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Create - Create a new assessment with the given name and the specified settings. Since name of an assessment in a project
// is a unique identifier, if an assessment with the name provided already exists, then
// the existing assessment is updated.
// Any PUT operation, resulting in either create or update on an assessment, will cause the assessment to go in a "InProgress"
// state. This will be indicated by the field 'computationState' on the
// Assessment object. During this time no other PUT operation will be allowed on that assessment object, nor will a Delete
// operation. Once the computation for the assessment is complete, the field
// 'computationState' will be updated to 'Ready', and then other PUT or DELETE operations can happen on the assessment.
// When assessment is under computation, any PUT will lead to a 400 - Bad Request error.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2019-10-01
// resourceGroupName - Name of the Azure Resource Group that project is part of.
// projectName - Name of the Azure Migrate project.
// groupName - Unique name of a group within a project.
// assessmentName - Unique name of an assessment within a project.
// options - AssessmentsClientCreateOptions contains the optional parameters for the AssessmentsClient.Create method.
func (client *AssessmentsClient) Create(ctx context.Context, resourceGroupName string, projectName string, groupName string, assessmentName string, options *AssessmentsClientCreateOptions) (AssessmentsClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, projectName, groupName, assessmentName, options)
	if err != nil {
		return AssessmentsClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AssessmentsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return AssessmentsClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *AssessmentsClient) createCreateRequest(ctx context.Context, resourceGroupName string, projectName string, groupName string, assessmentName string, options *AssessmentsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/groups/{groupName}/assessments/{assessmentName}"
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
	if groupName == "" {
		return nil, errors.New("parameter groupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{groupName}", url.PathEscape(groupName))
	if assessmentName == "" {
		return nil, errors.New("parameter assessmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{assessmentName}", url.PathEscape(assessmentName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.Assessment != nil {
		return req, runtime.MarshalAsJSON(req, *options.Assessment)
	}
	return req, nil
}

// createHandleResponse handles the Create response.
func (client *AssessmentsClient) createHandleResponse(resp *http.Response) (AssessmentsClientCreateResponse, error) {
	result := AssessmentsClientCreateResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Assessment); err != nil {
		return AssessmentsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Delete an assessment from the project. The machines remain in the assessment. Deleting a non-existent assessment
// results in a no-operation.
// When an assessment is under computation, as indicated by the 'computationState' field, it cannot be deleted. Any such attempt
// will return a 400 - Bad Request.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2019-10-01
// resourceGroupName - Name of the Azure Resource Group that project is part of.
// projectName - Name of the Azure Migrate project.
// groupName - Unique name of a group within a project.
// assessmentName - Unique name of an assessment within a project.
// options - AssessmentsClientDeleteOptions contains the optional parameters for the AssessmentsClient.Delete method.
func (client *AssessmentsClient) Delete(ctx context.Context, resourceGroupName string, projectName string, groupName string, assessmentName string, options *AssessmentsClientDeleteOptions) (AssessmentsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, projectName, groupName, assessmentName, options)
	if err != nil {
		return AssessmentsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AssessmentsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return AssessmentsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return client.deleteHandleResponse(resp)
}

// deleteCreateRequest creates the Delete request.
func (client *AssessmentsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, projectName string, groupName string, assessmentName string, options *AssessmentsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/groups/{groupName}/assessments/{assessmentName}"
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
	if groupName == "" {
		return nil, errors.New("parameter groupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{groupName}", url.PathEscape(groupName))
	if assessmentName == "" {
		return nil, errors.New("parameter assessmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{assessmentName}", url.PathEscape(assessmentName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *AssessmentsClient) deleteHandleResponse(resp *http.Response) (AssessmentsClientDeleteResponse, error) {
	result := AssessmentsClientDeleteResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	return result, nil
}

// Get - Get an existing assessment with the specified name. Returns a json object of type 'assessment' as specified in Models
// section.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2019-10-01
// resourceGroupName - Name of the Azure Resource Group that project is part of.
// projectName - Name of the Azure Migrate project.
// groupName - Unique name of a group within a project.
// assessmentName - Unique name of an assessment within a project.
// options - AssessmentsClientGetOptions contains the optional parameters for the AssessmentsClient.Get method.
func (client *AssessmentsClient) Get(ctx context.Context, resourceGroupName string, projectName string, groupName string, assessmentName string, options *AssessmentsClientGetOptions) (AssessmentsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, projectName, groupName, assessmentName, options)
	if err != nil {
		return AssessmentsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AssessmentsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AssessmentsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *AssessmentsClient) getCreateRequest(ctx context.Context, resourceGroupName string, projectName string, groupName string, assessmentName string, options *AssessmentsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/groups/{groupName}/assessments/{assessmentName}"
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
	if groupName == "" {
		return nil, errors.New("parameter groupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{groupName}", url.PathEscape(groupName))
	if assessmentName == "" {
		return nil, errors.New("parameter assessmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{assessmentName}", url.PathEscape(assessmentName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AssessmentsClient) getHandleResponse(resp *http.Response) (AssessmentsClientGetResponse, error) {
	result := AssessmentsClientGetResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Assessment); err != nil {
		return AssessmentsClientGetResponse{}, err
	}
	return result, nil
}

// GetReportDownloadURL - Get the URL for downloading the assessment in a report format.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2019-10-01
// resourceGroupName - Name of the Azure Resource Group that project is part of.
// projectName - Name of the Azure Migrate project.
// groupName - Unique name of a group within a project.
// assessmentName - Unique name of an assessment within a project.
// options - AssessmentsClientGetReportDownloadURLOptions contains the optional parameters for the AssessmentsClient.GetReportDownloadURL
// method.
func (client *AssessmentsClient) GetReportDownloadURL(ctx context.Context, resourceGroupName string, projectName string, groupName string, assessmentName string, options *AssessmentsClientGetReportDownloadURLOptions) (AssessmentsClientGetReportDownloadURLResponse, error) {
	req, err := client.getReportDownloadURLCreateRequest(ctx, resourceGroupName, projectName, groupName, assessmentName, options)
	if err != nil {
		return AssessmentsClientGetReportDownloadURLResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AssessmentsClientGetReportDownloadURLResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AssessmentsClientGetReportDownloadURLResponse{}, runtime.NewResponseError(resp)
	}
	return client.getReportDownloadURLHandleResponse(resp)
}

// getReportDownloadURLCreateRequest creates the GetReportDownloadURL request.
func (client *AssessmentsClient) getReportDownloadURLCreateRequest(ctx context.Context, resourceGroupName string, projectName string, groupName string, assessmentName string, options *AssessmentsClientGetReportDownloadURLOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/groups/{groupName}/assessments/{assessmentName}/downloadUrl"
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
	if groupName == "" {
		return nil, errors.New("parameter groupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{groupName}", url.PathEscape(groupName))
	if assessmentName == "" {
		return nil, errors.New("parameter assessmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{assessmentName}", url.PathEscape(assessmentName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getReportDownloadURLHandleResponse handles the GetReportDownloadURL response.
func (client *AssessmentsClient) getReportDownloadURLHandleResponse(resp *http.Response) (AssessmentsClientGetReportDownloadURLResponse, error) {
	result := AssessmentsClientGetReportDownloadURLResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.DownloadURL); err != nil {
		return AssessmentsClientGetReportDownloadURLResponse{}, err
	}
	return result, nil
}

// NewListByGroupPager - Get all assessments created for the specified group.
// Returns a json array of objects of type 'assessment' as specified in Models section.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2019-10-01
// resourceGroupName - Name of the Azure Resource Group that project is part of.
// projectName - Name of the Azure Migrate project.
// groupName - Unique name of a group within a project.
// options - AssessmentsClientListByGroupOptions contains the optional parameters for the AssessmentsClient.ListByGroup method.
func (client *AssessmentsClient) NewListByGroupPager(resourceGroupName string, projectName string, groupName string, options *AssessmentsClientListByGroupOptions) *runtime.Pager[AssessmentsClientListByGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[AssessmentsClientListByGroupResponse]{
		More: func(page AssessmentsClientListByGroupResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *AssessmentsClientListByGroupResponse) (AssessmentsClientListByGroupResponse, error) {
			req, err := client.listByGroupCreateRequest(ctx, resourceGroupName, projectName, groupName, options)
			if err != nil {
				return AssessmentsClientListByGroupResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return AssessmentsClientListByGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AssessmentsClientListByGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByGroupHandleResponse(resp)
		},
	})
}

// listByGroupCreateRequest creates the ListByGroup request.
func (client *AssessmentsClient) listByGroupCreateRequest(ctx context.Context, resourceGroupName string, projectName string, groupName string, options *AssessmentsClientListByGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/groups/{groupName}/assessments"
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
	if groupName == "" {
		return nil, errors.New("parameter groupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{groupName}", url.PathEscape(groupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByGroupHandleResponse handles the ListByGroup response.
func (client *AssessmentsClient) listByGroupHandleResponse(resp *http.Response) (AssessmentsClientListByGroupResponse, error) {
	result := AssessmentsClientListByGroupResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.AssessmentResultList); err != nil {
		return AssessmentsClientListByGroupResponse{}, err
	}
	return result, nil
}

// NewListByProjectPager - Get all assessments created in the project.
// Returns a json array of objects of type 'assessment' as specified in Models section.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2019-10-01
// resourceGroupName - Name of the Azure Resource Group that project is part of.
// projectName - Name of the Azure Migrate project.
// options - AssessmentsClientListByProjectOptions contains the optional parameters for the AssessmentsClient.ListByProject
// method.
func (client *AssessmentsClient) NewListByProjectPager(resourceGroupName string, projectName string, options *AssessmentsClientListByProjectOptions) *runtime.Pager[AssessmentsClientListByProjectResponse] {
	return runtime.NewPager(runtime.PagingHandler[AssessmentsClientListByProjectResponse]{
		More: func(page AssessmentsClientListByProjectResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *AssessmentsClientListByProjectResponse) (AssessmentsClientListByProjectResponse, error) {
			req, err := client.listByProjectCreateRequest(ctx, resourceGroupName, projectName, options)
			if err != nil {
				return AssessmentsClientListByProjectResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return AssessmentsClientListByProjectResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AssessmentsClientListByProjectResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByProjectHandleResponse(resp)
		},
	})
}

// listByProjectCreateRequest creates the ListByProject request.
func (client *AssessmentsClient) listByProjectCreateRequest(ctx context.Context, resourceGroupName string, projectName string, options *AssessmentsClientListByProjectOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/assessments"
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByProjectHandleResponse handles the ListByProject response.
func (client *AssessmentsClient) listByProjectHandleResponse(resp *http.Response) (AssessmentsClientListByProjectResponse, error) {
	result := AssessmentsClientListByProjectResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.AssessmentResultList); err != nil {
		return AssessmentsClientListByProjectResponse{}, err
	}
	return result, nil
}
