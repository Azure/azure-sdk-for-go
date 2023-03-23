// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package parameter

type Client struct {
}

func (c *Client) NotChange(ctx context.Context, resourceGroupName string, serviceName string, value string, option ClientOption) {

}

func (c *Client) OnlyToAny(ctx context.Context, resourceGroupName string, serviceName string, value interface{}, option ClientOption) {

}

func (c *Client) BeforeAny(ctx context.Context, resourceGroupName string, serviceName string, value interface{}, option ClientOption) {

}

func (c *Client) AfterAny(ctx context.Context, resourceGroupName string, serviceName string, value interface{}, option ClientOption) {

}
