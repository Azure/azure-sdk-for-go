// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armkeyvault

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// VaultsOperations contains the methods for the Vaults group.
type VaultsOperations interface {
	// CheckNameAvailability - Checks that the vault name is valid and is not already in use.
	CheckNameAvailability(ctx context.Context, vaultName VaultCheckNameAvailabilityParameters, options *VaultsCheckNameAvailabilityOptions) (*CheckNameAvailabilityResultResponse, error)
	// BeginCreateOrUpdate - Create or update a key vault in the specified subscription.
	BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, vaultName string, parameters VaultCreateOrUpdateParameters, options *VaultsCreateOrUpdateOptions) (*VaultPollerResponse, error)
	// ResumeCreateOrUpdate - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumeCreateOrUpdate(token string) (VaultPoller, error)
	// Delete - Deletes the specified Azure key vault.
	Delete(ctx context.Context, resourceGroupName string, vaultName string, options *VaultsDeleteOptions) (*http.Response, error)
	// Get - Gets the specified Azure key vault.
	Get(ctx context.Context, resourceGroupName string, vaultName string, options *VaultsGetOptions) (*VaultResponse, error)
	// GetDeleted - Gets the deleted Azure key vault.
	GetDeleted(ctx context.Context, vaultName string, location string, options *VaultsGetDeletedOptions) (*DeletedVaultResponse, error)
	// List - The List operation gets information about the vaults associated with the subscription.
	List(options *VaultsListOptions) ResourceListResultPager
	// ListByResourceGroup - The List operation gets information about the vaults associated with the subscription and within the specified resource group.
	ListByResourceGroup(resourceGroupName string, options *VaultsListByResourceGroupOptions) VaultListResultPager
	// ListBySubscription - The List operation gets information about the vaults associated with the subscription.
	ListBySubscription(options *VaultsListBySubscriptionOptions) VaultListResultPager
	// ListDeleted - Gets information about the deleted vaults in a subscription.
	ListDeleted(options *VaultsListDeletedOptions) DeletedVaultListResultPager
	// BeginPurgeDeleted - Permanently deletes the specified vault. aka Purges the deleted Azure key vault.
	BeginPurgeDeleted(ctx context.Context, vaultName string, location string, options *VaultsPurgeDeletedOptions) (*HTTPPollerResponse, error)
	// ResumePurgeDeleted - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumePurgeDeleted(token string) (HTTPPoller, error)
	// Update - Update a key vault in the specified subscription.
	Update(ctx context.Context, resourceGroupName string, vaultName string, parameters VaultPatchParameters, options *VaultsUpdateOptions) (*VaultResponse, error)
	// UpdateAccessPolicy - Update access policies in a key vault in the specified subscription.
	UpdateAccessPolicy(ctx context.Context, resourceGroupName string, vaultName string, operationKind AccessPolicyUpdateKind, parameters VaultAccessPolicyParameters, options *VaultsUpdateAccessPolicyOptions) (*VaultAccessPolicyParametersResponse, error)
}

// VaultsClient implements the VaultsOperations interface.
// Don't use this type directly, use NewVaultsClient() instead.
type VaultsClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewVaultsClient creates a new instance of VaultsClient with the specified values.
func NewVaultsClient(con *armcore.Connection, subscriptionID string) VaultsOperations {
	return &VaultsClient{con: con, subscriptionID: subscriptionID}
}

// Pipeline returns the pipeline associated with this client.
func (client *VaultsClient) Pipeline() azcore.Pipeline {
	return client.con.Pipeline()
}

// CheckNameAvailability - Checks that the vault name is valid and is not already in use.
func (client *VaultsClient) CheckNameAvailability(ctx context.Context, vaultName VaultCheckNameAvailabilityParameters, options *VaultsCheckNameAvailabilityOptions) (*CheckNameAvailabilityResultResponse, error) {
	req, err := client.CheckNameAvailabilityCreateRequest(ctx, vaultName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.CheckNameAvailabilityHandleError(resp)
	}
	result, err := client.CheckNameAvailabilityHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CheckNameAvailabilityCreateRequest creates the CheckNameAvailability request.
func (client *VaultsClient) CheckNameAvailabilityCreateRequest(ctx context.Context, vaultName VaultCheckNameAvailabilityParameters, options *VaultsCheckNameAvailabilityOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.KeyVault/checkNameAvailability"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2019-09-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(vaultName)
}

// CheckNameAvailabilityHandleResponse handles the CheckNameAvailability response.
func (client *VaultsClient) CheckNameAvailabilityHandleResponse(resp *azcore.Response) (*CheckNameAvailabilityResultResponse, error) {
	result := CheckNameAvailabilityResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.CheckNameAvailabilityResult)
}

// CheckNameAvailabilityHandleError handles the CheckNameAvailability error response.
func (client *VaultsClient) CheckNameAvailabilityHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

func (client *VaultsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, vaultName string, parameters VaultCreateOrUpdateParameters, options *VaultsCreateOrUpdateOptions) (*VaultPollerResponse, error) {
	resp, err := client.CreateOrUpdate(ctx, resourceGroupName, vaultName, parameters, options)
	if err != nil {
		return nil, err
	}
	result := &VaultPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("VaultsClient.CreateOrUpdate", "", resp, client.CreateOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	poller := &vaultPoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*VaultResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

func (client *VaultsClient) ResumeCreateOrUpdate(token string) (VaultPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("VaultsClient.CreateOrUpdate", token, client.CreateOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	return &vaultPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// CreateOrUpdate - Create or update a key vault in the specified subscription.
func (client *VaultsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, vaultName string, parameters VaultCreateOrUpdateParameters, options *VaultsCreateOrUpdateOptions) (*azcore.Response, error) {
	req, err := client.CreateOrUpdateCreateRequest(ctx, resourceGroupName, vaultName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusCreated) {
		return nil, client.CreateOrUpdateHandleError(resp)
	}
	return resp, nil
}

// CreateOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *VaultsClient) CreateOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, parameters VaultCreateOrUpdateParameters, options *VaultsCreateOrUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{vaultName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2019-09-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(parameters)
}

// CreateOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *VaultsClient) CreateOrUpdateHandleResponse(resp *azcore.Response) (*VaultResponse, error) {
	result := VaultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.Vault)
}

// CreateOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *VaultsClient) CreateOrUpdateHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// Delete - Deletes the specified Azure key vault.
func (client *VaultsClient) Delete(ctx context.Context, resourceGroupName string, vaultName string, options *VaultsDeleteOptions) (*http.Response, error) {
	req, err := client.DeleteCreateRequest(ctx, resourceGroupName, vaultName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusNoContent) {
		return nil, client.DeleteHandleError(resp)
	}
	return resp.Response, nil
}

// DeleteCreateRequest creates the Delete request.
func (client *VaultsClient) DeleteCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, options *VaultsDeleteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{vaultName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2019-09-01")
	req.URL.RawQuery = query.Encode()
	return req, nil
}

// DeleteHandleError handles the Delete error response.
func (client *VaultsClient) DeleteHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// Get - Gets the specified Azure key vault.
func (client *VaultsClient) Get(ctx context.Context, resourceGroupName string, vaultName string, options *VaultsGetOptions) (*VaultResponse, error) {
	req, err := client.GetCreateRequest(ctx, resourceGroupName, vaultName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.GetHandleError(resp)
	}
	result, err := client.GetHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetCreateRequest creates the Get request.
func (client *VaultsClient) GetCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, options *VaultsGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{vaultName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2019-09-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// GetHandleResponse handles the Get response.
func (client *VaultsClient) GetHandleResponse(resp *azcore.Response) (*VaultResponse, error) {
	result := VaultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.Vault)
}

// GetHandleError handles the Get error response.
func (client *VaultsClient) GetHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// GetDeleted - Gets the deleted Azure key vault.
func (client *VaultsClient) GetDeleted(ctx context.Context, vaultName string, location string, options *VaultsGetDeletedOptions) (*DeletedVaultResponse, error) {
	req, err := client.GetDeletedCreateRequest(ctx, vaultName, location, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.GetDeletedHandleError(resp)
	}
	result, err := client.GetDeletedHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetDeletedCreateRequest creates the GetDeleted request.
func (client *VaultsClient) GetDeletedCreateRequest(ctx context.Context, vaultName string, location string, options *VaultsGetDeletedOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.KeyVault/locations/{location}/deletedVaults/{vaultName}"
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2019-09-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// GetDeletedHandleResponse handles the GetDeleted response.
func (client *VaultsClient) GetDeletedHandleResponse(resp *azcore.Response) (*DeletedVaultResponse, error) {
	result := DeletedVaultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.DeletedVault)
}

// GetDeletedHandleError handles the GetDeleted error response.
func (client *VaultsClient) GetDeletedHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// List - The List operation gets information about the vaults associated with the subscription.
func (client *VaultsClient) List(options *VaultsListOptions) ResourceListResultPager {
	return &resourceListResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.ListCreateRequest(ctx, options)
		},
		responder: client.ListHandleResponse,
		errorer:   client.ListHandleError,
		advancer: func(ctx context.Context, resp *ResourceListResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.ResourceListResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// ListCreateRequest creates the List request.
func (client *VaultsClient) ListCreateRequest(ctx context.Context, options *VaultsListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resources"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("$filter", "resourceType eq 'Microsoft.KeyVault/vaults'")
	if options != nil && options.Top != nil {
		query.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	query.Set("api-version", "2015-11-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListHandleResponse handles the List response.
func (client *VaultsClient) ListHandleResponse(resp *azcore.Response) (*ResourceListResultResponse, error) {
	result := ResourceListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.ResourceListResult)
}

// ListHandleError handles the List error response.
func (client *VaultsClient) ListHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// ListByResourceGroup - The List operation gets information about the vaults associated with the subscription and within the specified resource group.
func (client *VaultsClient) ListByResourceGroup(resourceGroupName string, options *VaultsListByResourceGroupOptions) VaultListResultPager {
	return &vaultListResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.ListByResourceGroupCreateRequest(ctx, resourceGroupName, options)
		},
		responder: client.ListByResourceGroupHandleResponse,
		errorer:   client.ListByResourceGroupHandleError,
		advancer: func(ctx context.Context, resp *VaultListResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.VaultListResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// ListByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *VaultsClient) ListByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *VaultsListByResourceGroupOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	if options != nil && options.Top != nil {
		query.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	query.Set("api-version", "2019-09-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *VaultsClient) ListByResourceGroupHandleResponse(resp *azcore.Response) (*VaultListResultResponse, error) {
	result := VaultListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.VaultListResult)
}

// ListByResourceGroupHandleError handles the ListByResourceGroup error response.
func (client *VaultsClient) ListByResourceGroupHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// ListBySubscription - The List operation gets information about the vaults associated with the subscription.
func (client *VaultsClient) ListBySubscription(options *VaultsListBySubscriptionOptions) VaultListResultPager {
	return &vaultListResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.ListBySubscriptionCreateRequest(ctx, options)
		},
		responder: client.ListBySubscriptionHandleResponse,
		errorer:   client.ListBySubscriptionHandleError,
		advancer: func(ctx context.Context, resp *VaultListResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.VaultListResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// ListBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *VaultsClient) ListBySubscriptionCreateRequest(ctx context.Context, options *VaultsListBySubscriptionOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.KeyVault/vaults"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	if options != nil && options.Top != nil {
		query.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	query.Set("api-version", "2019-09-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *VaultsClient) ListBySubscriptionHandleResponse(resp *azcore.Response) (*VaultListResultResponse, error) {
	result := VaultListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.VaultListResult)
}

// ListBySubscriptionHandleError handles the ListBySubscription error response.
func (client *VaultsClient) ListBySubscriptionHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// ListDeleted - Gets information about the deleted vaults in a subscription.
func (client *VaultsClient) ListDeleted(options *VaultsListDeletedOptions) DeletedVaultListResultPager {
	return &deletedVaultListResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.ListDeletedCreateRequest(ctx, options)
		},
		responder: client.ListDeletedHandleResponse,
		errorer:   client.ListDeletedHandleError,
		advancer: func(ctx context.Context, resp *DeletedVaultListResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.DeletedVaultListResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// ListDeletedCreateRequest creates the ListDeleted request.
func (client *VaultsClient) ListDeletedCreateRequest(ctx context.Context, options *VaultsListDeletedOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.KeyVault/deletedVaults"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2019-09-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListDeletedHandleResponse handles the ListDeleted response.
func (client *VaultsClient) ListDeletedHandleResponse(resp *azcore.Response) (*DeletedVaultListResultResponse, error) {
	result := DeletedVaultListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.DeletedVaultListResult)
}

// ListDeletedHandleError handles the ListDeleted error response.
func (client *VaultsClient) ListDeletedHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

func (client *VaultsClient) BeginPurgeDeleted(ctx context.Context, vaultName string, location string, options *VaultsPurgeDeletedOptions) (*HTTPPollerResponse, error) {
	resp, err := client.PurgeDeleted(ctx, vaultName, location, options)
	if err != nil {
		return nil, err
	}
	result := &HTTPPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("VaultsClient.PurgeDeleted", "", resp, client.PurgeDeletedHandleError)
	if err != nil {
		return nil, err
	}
	poller := &httpPoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*http.Response, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

func (client *VaultsClient) ResumePurgeDeleted(token string) (HTTPPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("VaultsClient.PurgeDeleted", token, client.PurgeDeletedHandleError)
	if err != nil {
		return nil, err
	}
	return &httpPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// PurgeDeleted - Permanently deletes the specified vault. aka Purges the deleted Azure key vault.
func (client *VaultsClient) PurgeDeleted(ctx context.Context, vaultName string, location string, options *VaultsPurgeDeletedOptions) (*azcore.Response, error) {
	req, err := client.PurgeDeletedCreateRequest(ctx, vaultName, location, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted) {
		return nil, client.PurgeDeletedHandleError(resp)
	}
	return resp, nil
}

// PurgeDeletedCreateRequest creates the PurgeDeleted request.
func (client *VaultsClient) PurgeDeletedCreateRequest(ctx context.Context, vaultName string, location string, options *VaultsPurgeDeletedOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.KeyVault/locations/{location}/deletedVaults/{vaultName}/purge"
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2019-09-01")
	req.URL.RawQuery = query.Encode()
	return req, nil
}

// PurgeDeletedHandleError handles the PurgeDeleted error response.
func (client *VaultsClient) PurgeDeletedHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// Update - Update a key vault in the specified subscription.
func (client *VaultsClient) Update(ctx context.Context, resourceGroupName string, vaultName string, parameters VaultPatchParameters, options *VaultsUpdateOptions) (*VaultResponse, error) {
	req, err := client.UpdateCreateRequest(ctx, resourceGroupName, vaultName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusCreated) {
		return nil, client.UpdateHandleError(resp)
	}
	result, err := client.UpdateHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateCreateRequest creates the Update request.
func (client *VaultsClient) UpdateCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, parameters VaultPatchParameters, options *VaultsUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{vaultName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPatch, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2019-09-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(parameters)
}

// UpdateHandleResponse handles the Update response.
func (client *VaultsClient) UpdateHandleResponse(resp *azcore.Response) (*VaultResponse, error) {
	result := VaultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.Vault)
}

// UpdateHandleError handles the Update error response.
func (client *VaultsClient) UpdateHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// UpdateAccessPolicy - Update access policies in a key vault in the specified subscription.
func (client *VaultsClient) UpdateAccessPolicy(ctx context.Context, resourceGroupName string, vaultName string, operationKind AccessPolicyUpdateKind, parameters VaultAccessPolicyParameters, options *VaultsUpdateAccessPolicyOptions) (*VaultAccessPolicyParametersResponse, error) {
	req, err := client.UpdateAccessPolicyCreateRequest(ctx, resourceGroupName, vaultName, operationKind, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusCreated) {
		return nil, client.UpdateAccessPolicyHandleError(resp)
	}
	result, err := client.UpdateAccessPolicyHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateAccessPolicyCreateRequest creates the UpdateAccessPolicy request.
func (client *VaultsClient) UpdateAccessPolicyCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, operationKind AccessPolicyUpdateKind, parameters VaultAccessPolicyParameters, options *VaultsUpdateAccessPolicyOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{vaultName}/accessPolicies/{operationKind}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	urlPath = strings.ReplaceAll(urlPath, "{operationKind}", url.PathEscape(string(operationKind)))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2019-09-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(parameters)
}

// UpdateAccessPolicyHandleResponse handles the UpdateAccessPolicy response.
func (client *VaultsClient) UpdateAccessPolicyHandleResponse(resp *azcore.Response) (*VaultAccessPolicyParametersResponse, error) {
	result := VaultAccessPolicyParametersResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.VaultAccessPolicyParameters)
}

// UpdateAccessPolicyHandleError handles the UpdateAccessPolicy error response.
func (client *VaultsClient) UpdateAccessPolicyHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}
