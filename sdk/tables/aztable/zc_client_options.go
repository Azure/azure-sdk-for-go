// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
)

type ClientOptions struct {
	// HTTPClient sets the transport for making HTTP requests.
	HTTPClient azcore.Transporter
	// Retry configures the built-in retry policy behavior.
	Retry azcore.RetryOptions
	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry azcore.TelemetryOptions
	// PerCallOptions are options to run on every request
	PerCallOptions []azcore.Policy
	// Scopes are the authentication scopes for AAD Authentication
	Scopes []string
}

func (o *ClientOptions) getConnectionOptions() *generated.ConnectionOptions {
	if o == nil {
		return nil
	}

	return &generated.ConnectionOptions{
		HTTPClient: o.HTTPClient,
		Retry:      o.Retry,
		Telemetry:  o.Telemetry,
	}
}
