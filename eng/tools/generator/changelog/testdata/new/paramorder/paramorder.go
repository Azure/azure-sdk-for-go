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

// Example from the issue - parameters swapped order
func (client *AllPoliciesClient) NewListByServicePager(serviceName string, resourceGroupName string, options *AllPoliciesClientListByServiceOptions) *AllPoliciesClientListByServiceResponse {
	return nil
}

// Another test case - no change (should not be detected)
func (client *AllPoliciesClient) NoChange(ctx context.Context, name string, value int) error {
	return nil
}

// Test case with same types but different names (parameter renamed, should not be order change)
func (client *AllPoliciesClient) DifferentNames(firstName string, lastName string) error {
	return nil
}

// Test case with order change and same parameter names
func (client *AllPoliciesClient) OrderChanged(serviceName string, subscriptionID string, resourceGroupName string) error {
	return nil
}
