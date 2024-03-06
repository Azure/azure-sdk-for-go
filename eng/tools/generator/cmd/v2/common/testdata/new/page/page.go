// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package page

type Client struct{}

func (client *Client) NewGetLogPager(resourceGroupName string, options *ClientNewGetLogPagerOptions) (ClientNewGetLogPagerResponse, error) {

	return ClientNewGetLogPagerResponse{}, nil
}

type ClientNewGetLogPagerOptions struct{}

type ClientNewGetLogPagerResponse struct{}

func (client *Client) List(resourceGroupName string, options *ClientListOptions) (ClientListResponse, error) {

	return ClientListResponse{}, nil
}

type ClientListOptions struct{}

type ClientListResponse struct{}
