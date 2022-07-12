//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapimanagement

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

// DeletedServicesClient contains the methods for the DeletedServices group.
// Don't use this type directly, use NewDeletedServicesClient() instead.
type DeletedServicesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewDeletedServicesClient creates a new instance of DeletedServicesClient with the specified values.
// subscriptionID - Subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID forms
// part of the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewDeletedServicesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DeletedServicesClient, error) {
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
	client := &DeletedServicesClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// GetByName - Get soft-deleted Api Management Service by name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-08-01
// serviceName - The name of the API Management service.
// location - The location of the deleted API Management service.
// options - DeletedServicesClientGetByNameOptions contains the optional parameters for the DeletedServicesClient.GetByName
// method.
func (client *DeletedServicesClient) GetByName(ctx context.Context, serviceName string, location string, options *DeletedServicesClientGetByNameOptions) (DeletedServicesClientGetByNameResponse, error) {
	req, err := client.getByNameCreateRequest(ctx, serviceName, location, options)
	if err != nil {
		return DeletedServicesClientGetByNameResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DeletedServicesClientGetByNameResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DeletedServicesClientGetByNameResponse{}, runtime.NewResponseError(resp)
	}
	return client.getByNameHandleResponse(resp)
}

// getByNameCreateRequest creates the GetByName request.
func (client *DeletedServicesClient) getByNameCreateRequest(ctx context.Context, serviceName string, location string, options *DeletedServicesClientGetByNameOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ApiManagement/locations/{location}/deletedservices/{serviceName}"
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getByNameHandleResponse handles the GetByName response.
func (client *DeletedServicesClient) getByNameHandleResponse(resp *http.Response) (DeletedServicesClientGetByNameResponse, error) {
	result := DeletedServicesClientGetByNameResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeletedServiceContract); err != nil {
		return DeletedServicesClientGetByNameResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Lists all soft-deleted services available for undelete for the given subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-08-01
// options - DeletedServicesClientListBySubscriptionOptions contains the optional parameters for the DeletedServicesClient.ListBySubscription
// method.
func (client *DeletedServicesClient) NewListBySubscriptionPager(options *DeletedServicesClientListBySubscriptionOptions) *runtime.Pager[DeletedServicesClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[DeletedServicesClientListBySubscriptionResponse]{
		More: func(page DeletedServicesClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *DeletedServicesClientListBySubscriptionResponse) (DeletedServicesClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return DeletedServicesClientListBySubscriptionResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return DeletedServicesClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return DeletedServicesClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *DeletedServicesClient) listBySubscriptionCreateRequest(ctx context.Context, options *DeletedServicesClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ApiManagement/deletedservices"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *DeletedServicesClient) listBySubscriptionHandleResponse(resp *http.Response) (DeletedServicesClientListBySubscriptionResponse, error) {
	result := DeletedServicesClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeletedServicesCollection); err != nil {
		return DeletedServicesClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// BeginPurge - Purges Api Management Service (deletes it with no option to undelete).
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-08-01
// serviceName - The name of the API Management service.
// location - The location of the deleted API Management service.
// options - DeletedServicesClientBeginPurgeOptions contains the optional parameters for the DeletedServicesClient.BeginPurge
// method.
func (client *DeletedServicesClient) BeginPurge(ctx context.Context, serviceName string, location string, options *DeletedServicesClientBeginPurgeOptions) (*runtime.Poller[DeletedServicesClientPurgeResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.purge(ctx, serviceName, location, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[DeletedServicesClientPurgeResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[DeletedServicesClientPurgeResponse](options.ResumeToken, client.pl, nil)
	}
}

// Purge - Purges Api Management Service (deletes it with no option to undelete).
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-08-01
func (client *DeletedServicesClient) purge(ctx context.Context, serviceName string, location string, options *DeletedServicesClientBeginPurgeOptions) (*http.Response, error) {
	req, err := client.purgeCreateRequest(ctx, serviceName, location, options)
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

// purgeCreateRequest creates the Purge request.
func (client *DeletedServicesClient) purgeCreateRequest(ctx context.Context, serviceName string, location string, options *DeletedServicesClientBeginPurgeOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ApiManagement/locations/{location}/deletedservices/{serviceName}"
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}
