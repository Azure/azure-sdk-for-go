//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcdn

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

// SecretsClient contains the methods for the Secrets group.
// Don't use this type directly, use NewSecretsClient() instead.
type SecretsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewSecretsClient creates a new instance of SecretsClient with the specified values.
// subscriptionID - Azure Subscription ID.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewSecretsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SecretsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublicCloud.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &SecretsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreate - Creates a new Secret within the specified profile.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// secretName - Name of the Secret under the profile.
// secret - The Secret properties.
// options - SecretsClientBeginCreateOptions contains the optional parameters for the SecretsClient.BeginCreate method.
func (client *SecretsClient) BeginCreate(ctx context.Context, resourceGroupName string, profileName string, secretName string, secret Secret, options *SecretsClientBeginCreateOptions) (*armruntime.Poller[SecretsClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, profileName, secretName, secret, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller(resp, client.pl, &armruntime.NewPollerOptions[SecretsClientCreateResponse]{
			FinalStateVia: armruntime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return armruntime.NewPollerFromResumeToken[SecretsClientCreateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Create - Creates a new Secret within the specified profile.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *SecretsClient) create(ctx context.Context, resourceGroupName string, profileName string, secretName string, secret Secret, options *SecretsClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, profileName, secretName, secret, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *SecretsClient) createCreateRequest(ctx context.Context, resourceGroupName string, profileName string, secretName string, secret Secret, options *SecretsClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/secrets/{secretName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if secretName == "" {
		return nil, errors.New("parameter secretName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secretName}", url.PathEscape(secretName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, secret)
}

// BeginDelete - Deletes an existing Secret within profile.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// secretName - Name of the Secret under the profile.
// options - SecretsClientBeginDeleteOptions contains the optional parameters for the SecretsClient.BeginDelete method.
func (client *SecretsClient) BeginDelete(ctx context.Context, resourceGroupName string, profileName string, secretName string, options *SecretsClientBeginDeleteOptions) (*armruntime.Poller[SecretsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, profileName, secretName, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller(resp, client.pl, &armruntime.NewPollerOptions[SecretsClientDeleteResponse]{
			FinalStateVia: armruntime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return armruntime.NewPollerFromResumeToken[SecretsClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Deletes an existing Secret within profile.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *SecretsClient) deleteOperation(ctx context.Context, resourceGroupName string, profileName string, secretName string, options *SecretsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, profileName, secretName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *SecretsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, profileName string, secretName string, options *SecretsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/secrets/{secretName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if secretName == "" {
		return nil, errors.New("parameter secretName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secretName}", url.PathEscape(secretName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Gets an existing Secret within a profile.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// secretName - Name of the Secret under the profile.
// options - SecretsClientGetOptions contains the optional parameters for the SecretsClient.Get method.
func (client *SecretsClient) Get(ctx context.Context, resourceGroupName string, profileName string, secretName string, options *SecretsClientGetOptions) (SecretsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, profileName, secretName, options)
	if err != nil {
		return SecretsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SecretsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SecretsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *SecretsClient) getCreateRequest(ctx context.Context, resourceGroupName string, profileName string, secretName string, options *SecretsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/secrets/{secretName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if secretName == "" {
		return nil, errors.New("parameter secretName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secretName}", url.PathEscape(secretName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SecretsClient) getHandleResponse(resp *http.Response) (SecretsClientGetResponse, error) {
	result := SecretsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Secret); err != nil {
		return SecretsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByProfilePager - Lists existing AzureFrontDoor secrets.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// options - SecretsClientListByProfileOptions contains the optional parameters for the SecretsClient.ListByProfile method.
func (client *SecretsClient) NewListByProfilePager(resourceGroupName string, profileName string, options *SecretsClientListByProfileOptions) *runtime.Pager[SecretsClientListByProfileResponse] {
	return runtime.NewPager(runtime.PageProcessor[SecretsClientListByProfileResponse]{
		More: func(page SecretsClientListByProfileResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SecretsClientListByProfileResponse) (SecretsClientListByProfileResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByProfileCreateRequest(ctx, resourceGroupName, profileName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SecretsClientListByProfileResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return SecretsClientListByProfileResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SecretsClientListByProfileResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByProfileHandleResponse(resp)
		},
	})
}

// listByProfileCreateRequest creates the ListByProfile request.
func (client *SecretsClient) listByProfileCreateRequest(ctx context.Context, resourceGroupName string, profileName string, options *SecretsClientListByProfileOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/secrets"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByProfileHandleResponse handles the ListByProfile response.
func (client *SecretsClient) listByProfileHandleResponse(resp *http.Response) (SecretsClientListByProfileResponse, error) {
	result := SecretsClientListByProfileResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SecretListResult); err != nil {
		return SecretsClientListByProfileResponse{}, err
	}
	return result, nil
}
