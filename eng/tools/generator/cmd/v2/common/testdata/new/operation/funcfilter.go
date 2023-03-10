// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package operation

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
