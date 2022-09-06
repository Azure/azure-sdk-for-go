//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armautomation

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RunbookDraftClient contains the methods for the RunbookDraft group.
// Don't use this type directly, use NewRunbookDraftClient() instead.
type RunbookDraftClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewRunbookDraftClient creates a new instance of RunbookDraftClient with the specified values.
// subscriptionID - Gets subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID
// forms part of the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewRunbookDraftClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RunbookDraftClient, error) {
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
	client := &RunbookDraftClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Get - Retrieve the runbook draft identified by runbook name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-06-30
// resourceGroupName - Name of an Azure Resource group.
// automationAccountName - The name of the automation account.
// runbookName - The runbook name.
// options - RunbookDraftClientGetOptions contains the optional parameters for the RunbookDraftClient.Get method.
func (client *RunbookDraftClient) Get(ctx context.Context, resourceGroupName string, automationAccountName string, runbookName string, options *RunbookDraftClientGetOptions) (RunbookDraftClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, automationAccountName, runbookName, options)
	if err != nil {
		return RunbookDraftClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return RunbookDraftClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return RunbookDraftClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *RunbookDraftClient) getCreateRequest(ctx context.Context, resourceGroupName string, automationAccountName string, runbookName string, options *RunbookDraftClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/runbooks/{runbookName}/draft"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if automationAccountName == "" {
		return nil, errors.New("parameter automationAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{automationAccountName}", url.PathEscape(automationAccountName))
	if runbookName == "" {
		return nil, errors.New("parameter runbookName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{runbookName}", url.PathEscape(runbookName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-30")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *RunbookDraftClient) getHandleResponse(resp *http.Response) (RunbookDraftClientGetResponse, error) {
	result := RunbookDraftClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RunbookDraft); err != nil {
		return RunbookDraftClientGetResponse{}, err
	}
	return result, nil
}

// GetContent - Retrieve the content of runbook draft identified by runbook name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-06-30
// resourceGroupName - Name of an Azure Resource group.
// automationAccountName - The name of the automation account.
// runbookName - The runbook name.
// options - RunbookDraftClientGetContentOptions contains the optional parameters for the RunbookDraftClient.GetContent method.
func (client *RunbookDraftClient) GetContent(ctx context.Context, resourceGroupName string, automationAccountName string, runbookName string, options *RunbookDraftClientGetContentOptions) (RunbookDraftClientGetContentResponse, error) {
	req, err := client.getContentCreateRequest(ctx, resourceGroupName, automationAccountName, runbookName, options)
	if err != nil {
		return RunbookDraftClientGetContentResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return RunbookDraftClientGetContentResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent) {
		return RunbookDraftClientGetContentResponse{}, runtime.NewResponseError(resp)
	}
	return RunbookDraftClientGetContentResponse{}, nil
}

// getContentCreateRequest creates the GetContent request.
func (client *RunbookDraftClient) getContentCreateRequest(ctx context.Context, resourceGroupName string, automationAccountName string, runbookName string, options *RunbookDraftClientGetContentOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/runbooks/{runbookName}/draft/content"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if automationAccountName == "" {
		return nil, errors.New("parameter automationAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{automationAccountName}", url.PathEscape(automationAccountName))
	if runbookName == "" {
		return nil, errors.New("parameter runbookName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{runbookName}", url.PathEscape(runbookName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-30")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"text/powershell"}
	return req, nil
}

// BeginReplaceContent - Replaces the runbook draft content.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-06-30
// resourceGroupName - Name of an Azure Resource group.
// automationAccountName - The name of the automation account.
// runbookName - The runbook name.
// runbookContent - The runbook draft content.
// options - RunbookDraftClientBeginReplaceContentOptions contains the optional parameters for the RunbookDraftClient.BeginReplaceContent
// method.
func (client *RunbookDraftClient) BeginReplaceContent(ctx context.Context, resourceGroupName string, automationAccountName string, runbookName string, runbookContent io.ReadSeekCloser, options *RunbookDraftClientBeginReplaceContentOptions) (*runtime.Poller[RunbookDraftClientReplaceContentResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.replaceContent(ctx, resourceGroupName, automationAccountName, runbookName, runbookContent, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[RunbookDraftClientReplaceContentResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[RunbookDraftClientReplaceContentResponse](options.ResumeToken, client.pl, nil)
	}
}

// ReplaceContent - Replaces the runbook draft content.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-06-30
func (client *RunbookDraftClient) replaceContent(ctx context.Context, resourceGroupName string, automationAccountName string, runbookName string, runbookContent io.ReadSeekCloser, options *RunbookDraftClientBeginReplaceContentOptions) (*http.Response, error) {
	req, err := client.replaceContentCreateRequest(ctx, resourceGroupName, automationAccountName, runbookName, runbookContent, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// replaceContentCreateRequest creates the ReplaceContent request.
func (client *RunbookDraftClient) replaceContentCreateRequest(ctx context.Context, resourceGroupName string, automationAccountName string, runbookName string, runbookContent io.ReadSeekCloser, options *RunbookDraftClientBeginReplaceContentOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/runbooks/{runbookName}/draft/content"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if automationAccountName == "" {
		return nil, errors.New("parameter automationAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{automationAccountName}", url.PathEscape(automationAccountName))
	if runbookName == "" {
		return nil, errors.New("parameter runbookName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{runbookName}", url.PathEscape(runbookName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-30")
	req.Raw().URL.RawQuery = reqQP.Encode()
	runtime.SkipBodyDownload(req)
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, req.SetBody(runbookContent, "text/powershell")
}

// UndoEdit - Undo draft edit to last known published state identified by runbook name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-06-30
// resourceGroupName - Name of an Azure Resource group.
// automationAccountName - The name of the automation account.
// runbookName - The runbook name.
// options - RunbookDraftClientUndoEditOptions contains the optional parameters for the RunbookDraftClient.UndoEdit method.
func (client *RunbookDraftClient) UndoEdit(ctx context.Context, resourceGroupName string, automationAccountName string, runbookName string, options *RunbookDraftClientUndoEditOptions) (RunbookDraftClientUndoEditResponse, error) {
	req, err := client.undoEditCreateRequest(ctx, resourceGroupName, automationAccountName, runbookName, options)
	if err != nil {
		return RunbookDraftClientUndoEditResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return RunbookDraftClientUndoEditResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return RunbookDraftClientUndoEditResponse{}, runtime.NewResponseError(resp)
	}
	return client.undoEditHandleResponse(resp)
}

// undoEditCreateRequest creates the UndoEdit request.
func (client *RunbookDraftClient) undoEditCreateRequest(ctx context.Context, resourceGroupName string, automationAccountName string, runbookName string, options *RunbookDraftClientUndoEditOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/runbooks/{runbookName}/draft/undoEdit"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if automationAccountName == "" {
		return nil, errors.New("parameter automationAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{automationAccountName}", url.PathEscape(automationAccountName))
	if runbookName == "" {
		return nil, errors.New("parameter runbookName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{runbookName}", url.PathEscape(runbookName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-30")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// undoEditHandleResponse handles the UndoEdit response.
func (client *RunbookDraftClient) undoEditHandleResponse(resp *http.Response) (RunbookDraftClientUndoEditResponse, error) {
	result := RunbookDraftClientUndoEditResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RunbookDraftUndoEditResult); err != nil {
		return RunbookDraftClientUndoEditResponse{}, err
	}
	return result, nil
}
