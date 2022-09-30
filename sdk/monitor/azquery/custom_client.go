//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery

// this file contains handwritten additions to the generated code

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// MetricsClientOptions contains optional settings for MetricsClient.
type MetricsClientOptions struct {
	azcore.ClientOptions
}

// LogsClientOptions contains optional settings for LogsClient.
type LogsClientOptions struct {
	azcore.ClientOptions
}

// NewLogsClient creates a client that accesses Azure Monitor logs data.
func NewLogsClient(credential azcore.TokenCredential, options *LogsClientOptions) *LogsClient {
	if options == nil {
		options = &LogsClientOptions{}
	}
	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{"https://api.loganalytics.io/.default"}, nil)
	pl := runtime.NewPipeline(moduleName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &LogsClient{pl: pl}
}

// NewMetricsClient creates a client that accesses Azure Monitor metrics data.
func NewMetricsClient(credential azcore.TokenCredential, options *MetricsClientOptions) *MetricsClient {
	if options == nil {
		options = &MetricsClientOptions{}
	}
	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{"https://management.azure.com/.default"}, nil)
	pl := runtime.NewPipeline(moduleName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &MetricsClient{pl: pl}
}

const metricsHost string = "https://management.azure.com"

// ErrorInfo - The code and message for an error.
type ErrorInfo struct {
	// REQUIRED; A machine readable error code.
	Code *string `json:"code,omitempty"`

	// REQUIRED; A human readable error message.
	Message *string `json:"message,omitempty"`

	// A description in detail why the operation failed.
	RawMessage json.RawMessage
}

func (e *ErrorInfo) UnmarshalJSON(data []byte) error {
	e.RawMessage = data
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "code":
			err = unpopulate(val, "Code", &e.Code)
			delete(rawMsg, key)
		case "message":
			err = unpopulate(val, "Message", &e.Message)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", e, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ErrorInfo.
func (e ErrorInfo) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "code", e.Code)
	populate(objectMap, "message", e.Message)
	return json.Marshal(objectMap)
}
