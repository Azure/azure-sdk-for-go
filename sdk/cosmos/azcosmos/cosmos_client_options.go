// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const telemetryInfo = "azsdk-go-azcosmos/v0.1.0"

// CosmosClientOptions defines the options for the Cosmos client.
type CosmosClientOptions struct {
	// HTTPClient sets the transport for making HTTP requests.
	HTTPClient azcore.Transport
	// Retry configures the built-in retry policy behavior.
	Retry azcore.RetryOptions
	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry azcore.TelemetryOptions
	// Logging configures the built-in logging policy behavior.
	Logging azcore.LogOptions
	// PerCallPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request.
	PerCallPolicies []azcore.Policy
	// PerRetryPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request, and for each retry request.
	PerRetryPolicies []azcore.Policy
	// ApplicationPreferredRegions defines list of preferred regions for the client to connect to.
	ApplicationPreferredRegions *[]string
	// ConsistencyLevel can be used to weaken the database account consistency level for read operations. If this is not set the database account consistency level will be used for all requests.
	ConsistencyLevel ConsistencyLevel
	// When EnableContentResponseOnWrite is false will cause the response to have a null resource. This reduces networking and CPU load by not sending the resource back over the network and serializing it on the client.
	EnableContentResponseOnWrite bool
	// LimitToEndpoint limits the operations to the provided endpoint on the CosmosClient. See https://docs.microsoft.com/azure/cosmos-db/troubleshoot-sdk-availability
	LimitToEndpoint bool
	// RateLimitedRetry defines the retry configuration for rate limited requests.
	// By default, the SDK will do 9 retries.
	RateLimitedRetry *CosmosClientOptionsRateLimitedRetry
}

type CosmosClientOptionsRateLimitedRetry struct {
	// MaxRetryAttempts specifies the number of retries to perform on rate limited requests.
	MaxRetryAttempts int
	// MaxRetryWaitTime specifies the maximum time to wait for retries.
	MaxRetryWaitTime time.Duration
}

// enrichTelemetryOptions adds the current package version to the telemetry information, if any.
func (o *CosmosClientOptions) enrichTelemetryOptions() *azcore.TelemetryOptions {
	to := o.Telemetry
	if to.Value == "" {
		to.Value = telemetryInfo
	} else {
		to.Value = fmt.Sprintf("%s %s", telemetryInfo, to.Value)
	}
	return &to
}

// getSDKInternalPolicies builds a list of internal retry policies for the cosmos service.
// This includes throttling and failover policies.
func (o *CosmosClientOptions) getSDKInternalPolicies() []azcore.Policy {
	return []azcore.Policy{
		newResourceThrottleRetryPolicy(o),
		// TODO: Add more policies here.
	}
}
