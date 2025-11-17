// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package operation

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

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

// This function has same params but different return type - creates breaking change with nil params
func (client *Client) Get(resourceGroupName string, options *ClientGetOptions) error {
	return nil
}

type ClientGetOptions struct{}
