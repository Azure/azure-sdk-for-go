//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery

// this file contains handwritten additions to the generated code

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ClientOptions contains optional settings for Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewLogsClient creates a client that accesses a monitor.
func NewLogsClient(credential azcore.TokenCredential, options *ClientOptions) *LogsClient {
	if options == nil {
		options = &ClientOptions{}
	}
	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{"https://api.loganalytics.io/.default"}, nil)
	pl := runtime.NewPipeline(moduleName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &LogsClient{pl: pl}
}

// NewMetricsClient creates a client that accesses a monitor.
func NewMetricsClient(credential azcore.TokenCredential, options *ClientOptions) *MetricsClient {
	if options == nil {
		options = &ClientOptions{}
	}
	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{"https://management.azure.com"}, nil)
	pl := runtime.NewPipeline(moduleName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &MetricsClient{pl: pl}
}

func QueryTimeInterval(startTime time.Time, endTime time.Time) string {
	return startTime.Format(time.RFC3339) + "/" + endTime.Format(time.RFC3339)
}

const metricsHost string = "https://management.azure.com"
