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
	"context"
	"encoding/json"
	"github.com/Azure/azure-pipeline-go/pipeline"
	"io/ioutil"
	"net/http"
)

// LoadBalancerBackendAddressPoolsClient is the network Client
type LoadBalancerBackendAddressPoolsClient struct {
	ManagementClient
}

// NewLoadBalancerBackendAddressPoolsClient creates an instance of the LoadBalancerBackendAddressPoolsClient client.
func NewLoadBalancerBackendAddressPoolsClient(p pipeline.Pipeline) LoadBalancerBackendAddressPoolsClient {
	return LoadBalancerBackendAddressPoolsClient{NewManagementClient(p)}
}

// Get gets load balancer backend address pool.
//
// resourceGroupName is the name of the resource group. loadBalancerName is the name of the load balancer.
// backendAddressPoolName is the name of the backend address pool.
func (client LoadBalancerBackendAddressPoolsClient) Get(ctx context.Context, resourceGroupName string, loadBalancerName string, backendAddressPoolName string) (*BackendAddressPool, error) {
	req, err := client.getPreparer(resourceGroupName, loadBalancerName, backendAddressPoolName)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(ctx, responderPolicyFactory{responder: client.getResponder}, req)
	if err != nil {
		return nil, err
	}
	return resp.(*BackendAddressPool), err
}

// getPreparer prepares the Get request.
func (client LoadBalancerBackendAddressPoolsClient) getPreparer(resourceGroupName string, loadBalancerName string, backendAddressPoolName string) (pipeline.Request, error) {
	u := client.url
	u.Path = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/loadBalancers/{loadBalancerName}/backendAddressPools/{backendAddressPoolName}"
	req, err := pipeline.NewRequest("GET", u, nil)
	if err != nil {
		return req, pipeline.NewError(err, "failed to create request")
	}
	params := req.URL.Query()
	params.Set("api-version", "2018-01-01")
	req.URL.RawQuery = params.Encode()
	return req, nil
}

// getResponder handles the response to the Get request.
func (client LoadBalancerBackendAddressPoolsClient) getResponder(resp pipeline.Response) (pipeline.Response, error) {
	err := validateResponse(resp, http.StatusOK)
	if resp == nil {
		return nil, err
	}
	result := &BackendAddressPool{rawResponse: resp.Response()}
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

// List gets all the load balancer backed address pools.
//
// resourceGroupName is the name of the resource group. loadBalancerName is the name of the load balancer.
func (client LoadBalancerBackendAddressPoolsClient) List(ctx context.Context, resourceGroupName string, loadBalancerName string) (*LoadBalancerBackendAddressPoolListResult, error) {
	req, err := client.listPreparer(resourceGroupName, loadBalancerName)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(ctx, responderPolicyFactory{responder: client.listResponder}, req)
	if err != nil {
		return nil, err
	}
	return resp.(*LoadBalancerBackendAddressPoolListResult), err
}

// listPreparer prepares the List request.
func (client LoadBalancerBackendAddressPoolsClient) listPreparer(resourceGroupName string, loadBalancerName string) (pipeline.Request, error) {
	u := client.url
	u.Path = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/loadBalancers/{loadBalancerName}/backendAddressPools"
	req, err := pipeline.NewRequest("GET", u, nil)
	if err != nil {
		return req, pipeline.NewError(err, "failed to create request")
	}
	params := req.URL.Query()
	params.Set("api-version", "2018-01-01")
	req.URL.RawQuery = params.Encode()
	return req, nil
}

// listResponder handles the response to the List request.
func (client LoadBalancerBackendAddressPoolsClient) listResponder(resp pipeline.Response) (pipeline.Response, error) {
	err := validateResponse(resp, http.StatusOK)
	if resp == nil {
		return nil, err
	}
	result := &LoadBalancerBackendAddressPoolListResult{rawResponse: resp.Response()}
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
