//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armappcomplianceautomation

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

// ReportClient contains the methods for the Report group.
// Don't use this type directly, use NewReportClient() instead.
type ReportClient struct {
	internal *arm.Client
}

// NewReportClient creates a new instance of ReportClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewReportClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ReportClient, error) {
	cl, err := arm.NewClient(moduleName+".ReportClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ReportClient{
		internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a new AppComplianceAutomation report or update an exiting AppComplianceAutomation report.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-16-preview
//   - reportName - Report Name.
//   - parameters - Parameters for the create or update operation
//   - options - ReportClientBeginCreateOrUpdateOptions contains the optional parameters for the ReportClient.BeginCreateOrUpdate
//     method.
func (client *ReportClient) BeginCreateOrUpdate(ctx context.Context, reportName string, parameters ReportResource, options *ReportClientBeginCreateOrUpdateOptions) (*runtime.Poller[ReportClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, reportName, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ReportClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[ReportClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Create a new AppComplianceAutomation report or update an exiting AppComplianceAutomation report.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-16-preview
func (client *ReportClient) createOrUpdate(ctx context.Context, reportName string, parameters ReportResource, options *ReportClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, reportName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ReportClient) createOrUpdateCreateRequest(ctx context.Context, reportName string, parameters ReportResource, options *ReportClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.AppComplianceAutomation/reports/{reportName}"
	if reportName == "" {
		return nil, errors.New("parameter reportName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reportName}", url.PathEscape(reportName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-16-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// BeginDelete - Delete an AppComplianceAutomation report.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-16-preview
//   - reportName - Report Name.
//   - options - ReportClientBeginDeleteOptions contains the optional parameters for the ReportClient.BeginDelete method.
func (client *ReportClient) BeginDelete(ctx context.Context, reportName string, options *ReportClientBeginDeleteOptions) (*runtime.Poller[ReportClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, reportName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ReportClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[ReportClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Delete an AppComplianceAutomation report.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-16-preview
func (client *ReportClient) deleteOperation(ctx context.Context, reportName string, options *ReportClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, reportName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ReportClient) deleteCreateRequest(ctx context.Context, reportName string, options *ReportClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.AppComplianceAutomation/reports/{reportName}"
	if reportName == "" {
		return nil, errors.New("parameter reportName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reportName}", url.PathEscape(reportName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-16-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get the AppComplianceAutomation report and its properties.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-16-preview
//   - reportName - Report Name.
//   - options - ReportClientGetOptions contains the optional parameters for the ReportClient.Get method.
func (client *ReportClient) Get(ctx context.Context, reportName string, options *ReportClientGetOptions) (ReportClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, reportName, options)
	if err != nil {
		return ReportClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ReportClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ReportClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ReportClient) getCreateRequest(ctx context.Context, reportName string, options *ReportClientGetOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.AppComplianceAutomation/reports/{reportName}"
	if reportName == "" {
		return nil, errors.New("parameter reportName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reportName}", url.PathEscape(reportName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-16-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ReportClient) getHandleResponse(resp *http.Response) (ReportClientGetResponse, error) {
	result := ReportClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReportResource); err != nil {
		return ReportClientGetResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update an exiting AppComplianceAutomation report.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-16-preview
//   - reportName - Report Name.
//   - parameters - Parameters for the create or update operation
//   - options - ReportClientBeginUpdateOptions contains the optional parameters for the ReportClient.BeginUpdate method.
func (client *ReportClient) BeginUpdate(ctx context.Context, reportName string, parameters ReportResourcePatch, options *ReportClientBeginUpdateOptions) (*runtime.Poller[ReportClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, reportName, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ReportClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[ReportClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - Update an exiting AppComplianceAutomation report.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-16-preview
func (client *ReportClient) update(ctx context.Context, reportName string, parameters ReportResourcePatch, options *ReportClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, reportName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *ReportClient) updateCreateRequest(ctx context.Context, reportName string, parameters ReportResourcePatch, options *ReportClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.AppComplianceAutomation/reports/{reportName}"
	if reportName == "" {
		return nil, errors.New("parameter reportName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reportName}", url.PathEscape(reportName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-16-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}
