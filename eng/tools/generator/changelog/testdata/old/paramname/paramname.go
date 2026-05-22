// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package paramname

import "context"

type Client struct {
}

// Function where parameter name changes to underscore (should be filtered)
func (c *Client) NameToUnderscore(ctx context.Context, resourceGroupName string, serviceName string) error {
	return nil
}

// Function where multiple parameter names change to underscore (should be filtered)
func (c *Client) MultipleNamesToUnderscore(ctx context.Context, resourceGroupName string, serviceName string, value int) error {
	return nil
}

// Function where parameter type also changes (should NOT be filtered)
func (c *Client) TypeAndNameChange(ctx context.Context, resourceGroupName string) error {
	return nil
}

// Function where parameter name changes but not to underscore (should NOT be filtered)
func (c *Client) NameChangeNotToUnderscore(ctx context.Context, resourceGroupName string) error {
	return nil
}

// Function with no changes (should not appear in breaking changes)
func (c *Client) NoChange(ctx context.Context, resourceGroupName string) error {
	return nil
}
