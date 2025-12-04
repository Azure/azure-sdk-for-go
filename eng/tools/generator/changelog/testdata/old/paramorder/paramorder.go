// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package paramorder

import "context"

type AllPoliciesClient struct {
}

type AllPoliciesClientListByServiceOptions struct {
}

type AllPoliciesClientListByServiceResponse struct {
}

// Example from the issue - parameters in original order
func (client *AllPoliciesClient) NewListByServicePager(resourceGroupName string, serviceName string, options *AllPoliciesClientListByServiceOptions) *AllPoliciesClientListByServiceResponse {
	return nil
}

// Another test case - no change (should not be detected)
func (client *AllPoliciesClient) NoChange(ctx context.Context, name string, value int) error {
	return nil
}

// Test case with same types but different names (should not be detected as order change)
func (client *AllPoliciesClient) DifferentNames(oldName string, newName string) error {
	return nil
}

// Test case with order change and same parameter names
func (client *AllPoliciesClient) OrderChanged(resourceGroupName string, serviceName string, subscriptionID string) error {
	return nil
}
