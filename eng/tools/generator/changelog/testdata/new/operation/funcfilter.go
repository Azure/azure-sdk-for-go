// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package operation

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

type Client struct{}

func (client *Client) BeginCreateOrUpdate(resourceGroupName string, options *ClientBeginCreateOrUpdateOptions) (ClientBeginCreateOrUpdateResponse, error) {

	return ClientBeginCreateOrUpdateResponse{}, nil
}

type ClientBeginCreateOrUpdateOptions struct{}

type ClientBeginCreateOrUpdateResponse struct{}

func (client *Client) NewListBySubscriptionPager(options *ClientListBySubscriptionOptions) *runtime.Pager[ClientListBySubscriptionResponse] {
	return &runtime.Pager[ClientListBySubscriptionResponse]
}

type ClientListBySubscriptionOptions struct{}

type ClientListBySubscriptionResponse struct{}

// This function has same params but different return type - creates breaking change with nil params
func (client *Client) Get(resourceGroupName string, options *ClientGetOptions) (ClientGetResponse, error) {
	return ClientGetResponse{Data: "data"}, nil
}

type ClientGetOptions struct{}

type ClientGetResponse struct {
	Data string
}
