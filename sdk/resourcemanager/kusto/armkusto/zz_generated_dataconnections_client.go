//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armkusto

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

// DataConnectionsClient contains the methods for the DataConnections group.
// Don't use this type directly, use NewDataConnectionsClient() instead.
type DataConnectionsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewDataConnectionsClient creates a new instance of DataConnectionsClient with the specified values.
// subscriptionID - Gets subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID
// forms part of the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewDataConnectionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DataConnectionsClient, error) {
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
	client := &DataConnectionsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CheckNameAvailability - Checks that the data connection name is valid and is not already in use.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// databaseName - The name of the database in the Kusto cluster.
// dataConnectionName - The name of the data connection.
// options - DataConnectionsClientCheckNameAvailabilityOptions contains the optional parameters for the DataConnectionsClient.CheckNameAvailability
// method.
func (client *DataConnectionsClient) CheckNameAvailability(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName DataConnectionCheckNameRequest, options *DataConnectionsClientCheckNameAvailabilityOptions) (DataConnectionsClientCheckNameAvailabilityResponse, error) {
	req, err := client.checkNameAvailabilityCreateRequest(ctx, resourceGroupName, clusterName, databaseName, dataConnectionName, options)
	if err != nil {
		return DataConnectionsClientCheckNameAvailabilityResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DataConnectionsClientCheckNameAvailabilityResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DataConnectionsClientCheckNameAvailabilityResponse{}, runtime.NewResponseError(resp)
	}
	return client.checkNameAvailabilityHandleResponse(resp)
}

// checkNameAvailabilityCreateRequest creates the CheckNameAvailability request.
func (client *DataConnectionsClient) checkNameAvailabilityCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName DataConnectionCheckNameRequest, options *DataConnectionsClientCheckNameAvailabilityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/checkNameAvailability"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, dataConnectionName)
}

// checkNameAvailabilityHandleResponse handles the CheckNameAvailability response.
func (client *DataConnectionsClient) checkNameAvailabilityHandleResponse(resp *http.Response) (DataConnectionsClientCheckNameAvailabilityResponse, error) {
	result := DataConnectionsClientCheckNameAvailabilityResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CheckNameResult); err != nil {
		return DataConnectionsClientCheckNameAvailabilityResponse{}, err
	}
	return result, nil
}

// BeginCreateOrUpdate - Creates or updates a data connection.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// databaseName - The name of the database in the Kusto cluster.
// dataConnectionName - The name of the data connection.
// parameters - The data connection parameters supplied to the CreateOrUpdate operation.
// options - DataConnectionsClientBeginCreateOrUpdateOptions contains the optional parameters for the DataConnectionsClient.BeginCreateOrUpdate
// method.
func (client *DataConnectionsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string, parameters DataConnectionClassification, options *DataConnectionsClientBeginCreateOrUpdateOptions) (*armruntime.Poller[DataConnectionsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, clusterName, databaseName, dataConnectionName, parameters, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller[DataConnectionsClientCreateOrUpdateResponse](resp, client.pl, nil)
	} else {
		return armruntime.NewPollerFromResumeToken[DataConnectionsClientCreateOrUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// CreateOrUpdate - Creates or updates a data connection.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *DataConnectionsClient) createOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string, parameters DataConnectionClassification, options *DataConnectionsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, clusterName, databaseName, dataConnectionName, parameters, options)
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

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *DataConnectionsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string, parameters DataConnectionClassification, options *DataConnectionsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/dataConnections/{dataConnectionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if dataConnectionName == "" {
		return nil, errors.New("parameter dataConnectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataConnectionName}", url.PathEscape(dataConnectionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// BeginDataConnectionValidation - Checks that the data connection parameters are valid.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// databaseName - The name of the database in the Kusto cluster.
// parameters - The data connection parameters supplied to the CreateOrUpdate operation.
// options - DataConnectionsClientBeginDataConnectionValidationOptions contains the optional parameters for the DataConnectionsClient.BeginDataConnectionValidation
// method.
func (client *DataConnectionsClient) BeginDataConnectionValidation(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, parameters DataConnectionValidation, options *DataConnectionsClientBeginDataConnectionValidationOptions) (*armruntime.Poller[DataConnectionsClientDataConnectionValidationResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.dataConnectionValidation(ctx, resourceGroupName, clusterName, databaseName, parameters, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller(resp, client.pl, &armruntime.NewPollerOptions[DataConnectionsClientDataConnectionValidationResponse]{
			FinalStateVia: armruntime.FinalStateViaLocation,
		})
	} else {
		return armruntime.NewPollerFromResumeToken[DataConnectionsClientDataConnectionValidationResponse](options.ResumeToken, client.pl, nil)
	}
}

// DataConnectionValidation - Checks that the data connection parameters are valid.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *DataConnectionsClient) dataConnectionValidation(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, parameters DataConnectionValidation, options *DataConnectionsClientBeginDataConnectionValidationOptions) (*http.Response, error) {
	req, err := client.dataConnectionValidationCreateRequest(ctx, resourceGroupName, clusterName, databaseName, parameters, options)
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

// dataConnectionValidationCreateRequest creates the DataConnectionValidation request.
func (client *DataConnectionsClient) dataConnectionValidationCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, parameters DataConnectionValidation, options *DataConnectionsClientBeginDataConnectionValidationOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/dataConnectionValidation"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// BeginDelete - Deletes the data connection with the given name.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// databaseName - The name of the database in the Kusto cluster.
// dataConnectionName - The name of the data connection.
// options - DataConnectionsClientBeginDeleteOptions contains the optional parameters for the DataConnectionsClient.BeginDelete
// method.
func (client *DataConnectionsClient) BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string, options *DataConnectionsClientBeginDeleteOptions) (*armruntime.Poller[DataConnectionsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, clusterName, databaseName, dataConnectionName, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller[DataConnectionsClientDeleteResponse](resp, client.pl, nil)
	} else {
		return armruntime.NewPollerFromResumeToken[DataConnectionsClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Deletes the data connection with the given name.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *DataConnectionsClient) deleteOperation(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string, options *DataConnectionsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, clusterName, databaseName, dataConnectionName, options)
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
func (client *DataConnectionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string, options *DataConnectionsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/dataConnections/{dataConnectionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if dataConnectionName == "" {
		return nil, errors.New("parameter dataConnectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataConnectionName}", url.PathEscape(dataConnectionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Returns a data connection.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// databaseName - The name of the database in the Kusto cluster.
// dataConnectionName - The name of the data connection.
// options - DataConnectionsClientGetOptions contains the optional parameters for the DataConnectionsClient.Get method.
func (client *DataConnectionsClient) Get(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string, options *DataConnectionsClientGetOptions) (DataConnectionsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, clusterName, databaseName, dataConnectionName, options)
	if err != nil {
		return DataConnectionsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DataConnectionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DataConnectionsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *DataConnectionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string, options *DataConnectionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/dataConnections/{dataConnectionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if dataConnectionName == "" {
		return nil, errors.New("parameter dataConnectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataConnectionName}", url.PathEscape(dataConnectionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DataConnectionsClient) getHandleResponse(resp *http.Response) (DataConnectionsClientGetResponse, error) {
	result := DataConnectionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result); err != nil {
		return DataConnectionsClientGetResponse{}, err
	}
	return result, nil
}

// ListByDatabase - Returns the list of data connections of the given Kusto database.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// databaseName - The name of the database in the Kusto cluster.
// options - DataConnectionsClientListByDatabaseOptions contains the optional parameters for the DataConnectionsClient.ListByDatabase
// method.
func (client *DataConnectionsClient) ListByDatabase(resourceGroupName string, clusterName string, databaseName string, options *DataConnectionsClientListByDatabaseOptions) *runtime.Pager[DataConnectionsClientListByDatabaseResponse] {
	return runtime.NewPager(runtime.PageProcessor[DataConnectionsClientListByDatabaseResponse]{
		More: func(page DataConnectionsClientListByDatabaseResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *DataConnectionsClientListByDatabaseResponse) (DataConnectionsClientListByDatabaseResponse, error) {
			req, err := client.listByDatabaseCreateRequest(ctx, resourceGroupName, clusterName, databaseName, options)
			if err != nil {
				return DataConnectionsClientListByDatabaseResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return DataConnectionsClientListByDatabaseResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return DataConnectionsClientListByDatabaseResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByDatabaseHandleResponse(resp)
		},
	})
}

// listByDatabaseCreateRequest creates the ListByDatabase request.
func (client *DataConnectionsClient) listByDatabaseCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, options *DataConnectionsClientListByDatabaseOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/dataConnections"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByDatabaseHandleResponse handles the ListByDatabase response.
func (client *DataConnectionsClient) listByDatabaseHandleResponse(resp *http.Response) (DataConnectionsClientListByDatabaseResponse, error) {
	result := DataConnectionsClientListByDatabaseResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataConnectionListResult); err != nil {
		return DataConnectionsClientListByDatabaseResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Updates a data connection.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// databaseName - The name of the database in the Kusto cluster.
// dataConnectionName - The name of the data connection.
// parameters - The data connection parameters supplied to the Update operation.
// options - DataConnectionsClientBeginUpdateOptions contains the optional parameters for the DataConnectionsClient.BeginUpdate
// method.
func (client *DataConnectionsClient) BeginUpdate(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string, parameters DataConnectionClassification, options *DataConnectionsClientBeginUpdateOptions) (*armruntime.Poller[DataConnectionsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, clusterName, databaseName, dataConnectionName, parameters, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller[DataConnectionsClientUpdateResponse](resp, client.pl, nil)
	} else {
		return armruntime.NewPollerFromResumeToken[DataConnectionsClientUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Update - Updates a data connection.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *DataConnectionsClient) update(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string, parameters DataConnectionClassification, options *DataConnectionsClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, clusterName, databaseName, dataConnectionName, parameters, options)
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

// updateCreateRequest creates the Update request.
func (client *DataConnectionsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string, parameters DataConnectionClassification, options *DataConnectionsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/dataConnections/{dataConnectionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if dataConnectionName == "" {
		return nil, errors.New("parameter dataConnectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataConnectionName}", url.PathEscape(dataConnectionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}
