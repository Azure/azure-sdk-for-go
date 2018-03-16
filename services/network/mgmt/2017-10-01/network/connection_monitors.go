package network

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Azure/azure-pipeline-go/pipeline"
	"io/ioutil"
	"net/http"
)

// ConnectionMonitorsClient is the network Client
type ConnectionMonitorsClient struct {
	ManagementClient
}

// NewConnectionMonitorsClient creates an instance of the ConnectionMonitorsClient client.
func NewConnectionMonitorsClient(p pipeline.Pipeline) ConnectionMonitorsClient {
	return ConnectionMonitorsClient{NewManagementClient(p)}
}

// CreateOrUpdate create or update a connection monitor. This method may poll for completion. Polling can be canceled
// by passing the cancel channel argument. The channel will be used to cancel polling and any outstanding HTTP
// requests.
//
// resourceGroupName is the name of the resource group containing Network Watcher. networkWatcherName is the name of
// the Network Watcher resource. connectionMonitorName is the name of the connection monitor. parameters is parameters
// that define the operation to create a connection monitor.
func (client ConnectionMonitorsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, networkWatcherName string, connectionMonitorName string, parameters ConnectionMonitor) (*ConnectionMonitorResult, error) {
	if err := validate([]validation{
		{targetValue: parameters,
			constraints: []constraint{{target: "parameters.ConnectionMonitorParameters", name: null, rule: true,
				chain: []constraint{{target: "parameters.ConnectionMonitorParameters.Source", name: null, rule: true,
					chain: []constraint{{target: "parameters.ConnectionMonitorParameters.Source.ResourceID", name: null, rule: true, chain: nil}}},
					{target: "parameters.ConnectionMonitorParameters.Destination", name: null, rule: true, chain: nil},
				}}}}}); err != nil {
		return nil, err
	}
	req, err := client.createOrUpdatePreparer(resourceGroupName, networkWatcherName, connectionMonitorName, parameters)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(ctx, responderPolicyFactory{responder: client.createOrUpdateResponder}, req)
	if err != nil {
		return nil, err
	}
	return resp.(*ConnectionMonitorResult), err
}

// createOrUpdatePreparer prepares the CreateOrUpdate request.
func (client ConnectionMonitorsClient) createOrUpdatePreparer(resourceGroupName string, networkWatcherName string, connectionMonitorName string, parameters ConnectionMonitor) (pipeline.Request, error) {
	u := client.url
	u.Path = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/connectionMonitors/{connectionMonitorName}"
	req, err := pipeline.NewRequest("PUT", u, nil)
	if err != nil {
		return req, pipeline.NewError(err, "failed to create request")
	}
	params := req.URL.Query()
	params.Set("api-version", "2017-10-01")
	req.URL.RawQuery = params.Encode()
	b, err := json.Marshal(parameters)
	if err != nil {
		return req, pipeline.NewError(err, "failed to marshal request body")
	}
	req.Header.Set("Content-Type", "application/json")
	err = req.SetBody(bytes.NewReader(b))
	if err != nil {
		return req, pipeline.NewError(err, "failed to set request body")
	}
	return req, nil
}

// createOrUpdateResponder handles the response to the CreateOrUpdate request.
func (client ConnectionMonitorsClient) createOrUpdateResponder(resp pipeline.Response) (pipeline.Response, error) {
	err := validateResponse(resp, http.StatusOK, http.StatusCreated)
	if resp == nil {
		return nil, err
	}
	result := &ConnectionMonitorResult{rawResponse: resp.Response()}
	if err != nil {
		return result, err
	}
	defer resp.Response().Body.Close()
	b, err := ioutil.ReadAll(resp.Response().Body)
	if err != nil {
		return result, NewResponseError(err, resp.Response(), "failed to read response body")
	}
	if len(b) > 0 {
		err = json.Unmarshal(b, result)
		if err != nil {
			return result, NewResponseError(err, resp.Response(), "failed to unmarshal response body")
		}
	}
	return result, nil
}

// Delete deletes the specified connection monitor. This method may poll for completion. Polling can be canceled by
// passing the cancel channel argument. The channel will be used to cancel polling and any outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group containing Network Watcher. networkWatcherName is the name of
// the Network Watcher resource. connectionMonitorName is the name of the connection monitor.
func (client ConnectionMonitorsClient) Delete(ctx context.Context, resourceGroupName string, networkWatcherName string, connectionMonitorName string) (*http.Response, error) {
	req, err := client.deletePreparer(resourceGroupName, networkWatcherName, connectionMonitorName)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(ctx, responderPolicyFactory{responder: client.deleteResponder}, req)
	if err != nil {
		return nil, err
	}
	return resp.Response(), err
}

// deletePreparer prepares the Delete request.
func (client ConnectionMonitorsClient) deletePreparer(resourceGroupName string, networkWatcherName string, connectionMonitorName string) (pipeline.Request, error) {
	u := client.url
	u.Path = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/connectionMonitors/{connectionMonitorName}"
	req, err := pipeline.NewRequest("DELETE", u, nil)
	if err != nil {
		return req, pipeline.NewError(err, "failed to create request")
	}
	params := req.URL.Query()
	params.Set("api-version", "2017-10-01")
	req.URL.RawQuery = params.Encode()
	return req, nil
}

// deleteResponder handles the response to the Delete request.
func (client ConnectionMonitorsClient) deleteResponder(resp pipeline.Response) (pipeline.Response, error) {
	err := validateResponse(resp, http.StatusOK, http.StatusNoContent, http.StatusAccepted)
	if resp == nil {
		return nil, err
	}
	resp.Response().Body.Close()
	return resp, err
}

// Get gets a connection monitor by name.
//
// resourceGroupName is the name of the resource group containing Network Watcher. networkWatcherName is the name of
// the Network Watcher resource. connectionMonitorName is the name of the connection monitor.
func (client ConnectionMonitorsClient) Get(ctx context.Context, resourceGroupName string, networkWatcherName string, connectionMonitorName string) (*ConnectionMonitorResult, error) {
	req, err := client.getPreparer(resourceGroupName, networkWatcherName, connectionMonitorName)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(ctx, responderPolicyFactory{responder: client.getResponder}, req)
	if err != nil {
		return nil, err
	}
	return resp.(*ConnectionMonitorResult), err
}

// getPreparer prepares the Get request.
func (client ConnectionMonitorsClient) getPreparer(resourceGroupName string, networkWatcherName string, connectionMonitorName string) (pipeline.Request, error) {
	u := client.url
	u.Path = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/connectionMonitors/{connectionMonitorName}"
	req, err := pipeline.NewRequest("GET", u, nil)
	if err != nil {
		return req, pipeline.NewError(err, "failed to create request")
	}
	params := req.URL.Query()
	params.Set("api-version", "2017-10-01")
	req.URL.RawQuery = params.Encode()
	return req, nil
}

// getResponder handles the response to the Get request.
func (client ConnectionMonitorsClient) getResponder(resp pipeline.Response) (pipeline.Response, error) {
	err := validateResponse(resp, http.StatusOK)
	if resp == nil {
		return nil, err
	}
	result := &ConnectionMonitorResult{rawResponse: resp.Response()}
	if err != nil {
		return result, err
	}
	defer resp.Response().Body.Close()
	b, err := ioutil.ReadAll(resp.Response().Body)
	if err != nil {
		return result, NewResponseError(err, resp.Response(), "failed to read response body")
	}
	if len(b) > 0 {
		err = json.Unmarshal(b, result)
		if err != nil {
			return result, NewResponseError(err, resp.Response(), "failed to unmarshal response body")
		}
	}
	return result, nil
}

// List lists all connection monitors for the specified Network Watcher.
//
// resourceGroupName is the name of the resource group containing Network Watcher. networkWatcherName is the name of
// the Network Watcher resource.
func (client ConnectionMonitorsClient) List(ctx context.Context, resourceGroupName string, networkWatcherName string) (*ConnectionMonitorListResult, error) {
	req, err := client.listPreparer(resourceGroupName, networkWatcherName)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(ctx, responderPolicyFactory{responder: client.listResponder}, req)
	if err != nil {
		return nil, err
	}
	return resp.(*ConnectionMonitorListResult), err
}

// listPreparer prepares the List request.
func (client ConnectionMonitorsClient) listPreparer(resourceGroupName string, networkWatcherName string) (pipeline.Request, error) {
	u := client.url
	u.Path = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/connectionMonitors"
	req, err := pipeline.NewRequest("GET", u, nil)
	if err != nil {
		return req, pipeline.NewError(err, "failed to create request")
	}
	params := req.URL.Query()
	params.Set("api-version", "2017-10-01")
	req.URL.RawQuery = params.Encode()
	return req, nil
}

// listResponder handles the response to the List request.
func (client ConnectionMonitorsClient) listResponder(resp pipeline.Response) (pipeline.Response, error) {
	err := validateResponse(resp, http.StatusOK)
	if resp == nil {
		return nil, err
	}
	result := &ConnectionMonitorListResult{rawResponse: resp.Response()}
	if err != nil {
		return result, err
	}
	defer resp.Response().Body.Close()
	b, err := ioutil.ReadAll(resp.Response().Body)
	if err != nil {
		return result, NewResponseError(err, resp.Response(), "failed to read response body")
	}
	if len(b) > 0 {
		err = json.Unmarshal(b, result)
		if err != nil {
			return result, NewResponseError(err, resp.Response(), "failed to unmarshal response body")
		}
	}
	return result, nil
}

// Query query a snapshot of the most recent connection states. This method may poll for completion. Polling can be
// canceled by passing the cancel channel argument. The channel will be used to cancel polling and any outstanding HTTP
// requests.
//
// resourceGroupName is the name of the resource group containing Network Watcher. networkWatcherName is the name of
// the Network Watcher resource. connectionMonitorName is the name given to the connection monitor.
func (client ConnectionMonitorsClient) Query(ctx context.Context, resourceGroupName string, networkWatcherName string, connectionMonitorName string) (*ConnectionMonitorQueryResult, error) {
	req, err := client.queryPreparer(resourceGroupName, networkWatcherName, connectionMonitorName)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(ctx, responderPolicyFactory{responder: client.queryResponder}, req)
	if err != nil {
		return nil, err
	}
	return resp.(*ConnectionMonitorQueryResult), err
}

// queryPreparer prepares the Query request.
func (client ConnectionMonitorsClient) queryPreparer(resourceGroupName string, networkWatcherName string, connectionMonitorName string) (pipeline.Request, error) {
	u := client.url
	u.Path = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/connectionMonitors/{connectionMonitorName}/query"
	req, err := pipeline.NewRequest("POST", u, nil)
	if err != nil {
		return req, pipeline.NewError(err, "failed to create request")
	}
	params := req.URL.Query()
	params.Set("api-version", "2017-10-01")
	req.URL.RawQuery = params.Encode()
	return req, nil
}

// queryResponder handles the response to the Query request.
func (client ConnectionMonitorsClient) queryResponder(resp pipeline.Response) (pipeline.Response, error) {
	err := validateResponse(resp, http.StatusOK, http.StatusAccepted)
	if resp == nil {
		return nil, err
	}
	result := &ConnectionMonitorQueryResult{rawResponse: resp.Response()}
	if err != nil {
		return result, err
	}
	defer resp.Response().Body.Close()
	b, err := ioutil.ReadAll(resp.Response().Body)
	if err != nil {
		return result, NewResponseError(err, resp.Response(), "failed to read response body")
	}
	if len(b) > 0 {
		err = json.Unmarshal(b, result)
		if err != nil {
			return result, NewResponseError(err, resp.Response(), "failed to unmarshal response body")
		}
	}
	return result, nil
}

// Start starts the specified connection monitor. This method may poll for completion. Polling can be canceled by
// passing the cancel channel argument. The channel will be used to cancel polling and any outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group containing Network Watcher. networkWatcherName is the name of
// the Network Watcher resource. connectionMonitorName is the name of the connection monitor.
func (client ConnectionMonitorsClient) Start(ctx context.Context, resourceGroupName string, networkWatcherName string, connectionMonitorName string) (*http.Response, error) {
	req, err := client.startPreparer(resourceGroupName, networkWatcherName, connectionMonitorName)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(ctx, responderPolicyFactory{responder: client.startResponder}, req)
	if err != nil {
		return nil, err
	}
	return resp.Response(), err
}

// startPreparer prepares the Start request.
func (client ConnectionMonitorsClient) startPreparer(resourceGroupName string, networkWatcherName string, connectionMonitorName string) (pipeline.Request, error) {
	u := client.url
	u.Path = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/connectionMonitors/{connectionMonitorName}/start"
	req, err := pipeline.NewRequest("POST", u, nil)
	if err != nil {
		return req, pipeline.NewError(err, "failed to create request")
	}
	params := req.URL.Query()
	params.Set("api-version", "2017-10-01")
	req.URL.RawQuery = params.Encode()
	return req, nil
}

// startResponder handles the response to the Start request.
func (client ConnectionMonitorsClient) startResponder(resp pipeline.Response) (pipeline.Response, error) {
	err := validateResponse(resp, http.StatusOK, http.StatusAccepted)
	if resp == nil {
		return nil, err
	}
	resp.Response().Body.Close()
	return resp, err
}

// Stop stops the specified connection monitor. This method may poll for completion. Polling can be canceled by passing
// the cancel channel argument. The channel will be used to cancel polling and any outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group containing Network Watcher. networkWatcherName is the name of
// the Network Watcher resource. connectionMonitorName is the name of the connection monitor.
func (client ConnectionMonitorsClient) Stop(ctx context.Context, resourceGroupName string, networkWatcherName string, connectionMonitorName string) (*http.Response, error) {
	req, err := client.stopPreparer(resourceGroupName, networkWatcherName, connectionMonitorName)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(ctx, responderPolicyFactory{responder: client.stopResponder}, req)
	if err != nil {
		return nil, err
	}
	return resp.Response(), err
}

// stopPreparer prepares the Stop request.
func (client ConnectionMonitorsClient) stopPreparer(resourceGroupName string, networkWatcherName string, connectionMonitorName string) (pipeline.Request, error) {
	u := client.url
	u.Path = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/connectionMonitors/{connectionMonitorName}/stop"
	req, err := pipeline.NewRequest("POST", u, nil)
	if err != nil {
		return req, pipeline.NewError(err, "failed to create request")
	}
	params := req.URL.Query()
	params.Set("api-version", "2017-10-01")
	req.URL.RawQuery = params.Encode()
	return req, nil
}

// stopResponder handles the response to the Stop request.
func (client ConnectionMonitorsClient) stopResponder(resp pipeline.Response) (pipeline.Response, error) {
	err := validateResponse(resp, http.StatusOK, http.StatusAccepted)
	if resp == nil {
		return nil, err
	}
	resp.Response().Body.Close()
	return resp, err
}
