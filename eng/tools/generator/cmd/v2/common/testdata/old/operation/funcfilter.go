// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package operation

type Client struct{}

func (client *Client) Update(resourceGroupName string, options *ClientUpdateOptions) (ClientUpdateResponse, error) {

	return ClientUpdateResponse{}, nil
}

type ClientUpdateOptions struct{}

type ClientUpdateResponse struct{}

func (client *Client) BeingDelete(resourceGroupName string, options *ClientBeginDeleteOptions) (*runtime.Poller[ClientBeginDeleteResponse], error) {

	return &runtime.Poller[ClientBeginDeleteResponse]{}, nil
}

type ClientBeginDeleteOptions struct{}

type ClientBeginDeleteResponse struct{}

func (client *Client) NewListPager(resourceGroupName string, options *ClientListOptions) *runtime.Pager[ClientListResponse] {
	return &runtime.Pager[ClientListResponse]
}

type ClientListOptions struct{}

type ClientListResponse struct{}
