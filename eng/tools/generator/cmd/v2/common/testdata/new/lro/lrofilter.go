// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lro

type Client struct{}

func (client *Client) BeginCreateOrUpdate(resourceGroupName string, options *ClientBeginCreateOrUpdateOptions) (ClientBeginCreateOrUpdateResponse, error) {

	return ClientBeginCreateOrUpdateResponse{}, nil
}

type ClientBeginCreateOrUpdateOptions struct{}

type ClientBeginCreateOrUpdateResponse struct{}

func (client *Client) Delete(resourceGroupName string, options *ClientDeleteOptions) (ClientDeleteResponse, error) {

	return ClientDeleteResponse{}, nil
}

type ClientDeleteOptions struct{}

type ClientDeleteResponse struct{}
