//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"

// ClientOptions contains optional settings for a client's pipeline.
// All zero-value fields will be initialized with default values.
type ClientOptions struct {
	// Logging configures the built-in logging policy.
	Logging policy.LogOptions

	// Retry configures the built-in retry policy.
	Retry policy.RetryOptions

	// Telemetry configures the built-in telemetry policy.
	Telemetry policy.TelemetryOptions

	// Transport sets the transport for HTTP requests.
	Transport policy.Transporter

	// PerCallPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request.
	PerCallPolicies []policy.Policy

	// PerTryPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request, and for each retry of that request.
	PerTryPolicies []policy.Policy
}
