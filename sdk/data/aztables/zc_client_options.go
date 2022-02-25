// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type ClientOptions struct {
	azcore.ClientOptions
}

func (c *ClientOptions) toPolicyOptions() *azcore.ClientOptions {
	return &azcore.ClientOptions{
		Logging:          c.Logging,
		Retry:            c.Retry,
		Telemetry:        c.Telemetry,
		Transport:        c.Transport,
		PerCallPolicies:  c.PerCallPolicies,
		PerRetryPolicies: c.PerRetryPolicies,
	}
}
