// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lro

type Client struct{}

func (client *Client) CreateOrUpdate(resourceGroupName string, options *ClientCreateOrUpdateOptions) (ClientCreateOrUpdateResponse, error) {

	return ClientCreateOrUpdateResponse{}, nil
}

type ClientCreateOrUpdateOptions struct{}

type ClientCreateOrUpdateResponse struct{}

func (client *Client) BeginDelete(resourceGroupName string, options *ClientBeginDeleteOptions) (ClientBeginDeleteResponse, error) {

	return ClientBeginDeleteResponse{}, nil
}

type ClientBeginDeleteOptions struct{}

type ClientBeginDeleteResponse struct{}
