// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// ClientOptions are the optional parameters for the NewClient method
type ClientOptions struct {
	azcore.ClientOptions
}

func (c *ClientOptions) toPolicyOptions() *policy.ClientOptions {
	if c == nil {
		return &policy.ClientOptions{}
	}

	return &policy.ClientOptions{
		Logging:          c.Logging,
		Retry:            c.Retry,
		Telemetry:        c.Telemetry,
		Transport:        c.Transport,
		PerCallPolicies:  c.PerCallPolicies,
		PerRetryPolicies: c.PerRetryPolicies,
	}
}

func policyOptionsToClientOptions(c *policy.ClientOptions) *ClientOptions {
	if c == nil {
		return nil
	}

	return &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Logging:          c.Logging,
			Retry:            c.Retry,
			Telemetry:        c.Telemetry,
			Transport:        c.Transport,
			PerCallPolicies:  c.PerCallPolicies,
			PerRetryPolicies: c.PerRetryPolicies,
		},
	}
}
