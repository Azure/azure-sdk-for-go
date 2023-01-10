//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery

// this file contains handwritten additions to the generated code

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
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

// LogsClient contains the methods for the LogsClient group.
// Don't use this type directly, use NewLogsClient() instead.
type LogsClient struct {
	host string
	pl   runtime.Pipeline
}

// MetricsClient contains the methods for the Metrics group.
// Don't use this type directly, use NewMetricsClient() instead.
type MetricsClient struct {
	host string
	pl   runtime.Pipeline
}

// NewLogsClient creates a client that accesses Azure Monitor logs data.
func NewLogsClient(credential azcore.TokenCredential, options *LogsClientOptions) (*LogsClient, error) {
	if options == nil {
		options = &LogsClientOptions{}
	}
	if reflect.ValueOf(options.Cloud).IsZero() {
		options.Cloud = cloud.AzurePublic
	}
	c, ok := options.Cloud.Services[ServiceNameLogs]
	if !ok || c.Audience == "" || c.Endpoint == "" {
		return nil, errors.New("provided Cloud field is missing Azure Monitor Logs configuration")
	}

	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{c.Audience + "/.default"}, nil)
	pl := runtime.NewPipeline(moduleName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &LogsClient{host: c.Endpoint, pl: pl}, nil
}

// NewMetricsClient creates a client that accesses Azure Monitor metrics data.
func NewMetricsClient(credential azcore.TokenCredential, options *MetricsClientOptions) (*MetricsClient, error) {
	if options == nil {
		options = &MetricsClientOptions{}
	}
	if reflect.ValueOf(options.Cloud).IsZero() {
		options.Cloud = cloud.AzurePublic
	}
	c, ok := options.Cloud.Services[ServiceNameMetrics]
	if !ok || c.Audience == "" || c.Endpoint == "" {
		return nil, errors.New("provided Cloud field is missing Azure Monitor Metrics configuration")
	}

	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{c.Audience + "/.default"}, nil)
	pl := runtime.NewPipeline(moduleName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &MetricsClient{host: c.Endpoint, pl: pl}, nil
}

// ErrorInfo - The code and message for an error.
type ErrorInfo struct {
	// REQUIRED; A machine readable error code.
	Code string

	// full error message detailing why the operation failed.
	data []byte
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ErrorInfo.
func (e *ErrorInfo) UnmarshalJSON(data []byte) error {
	e.data = data
	ei := struct{ Code string }{}
	if err := json.Unmarshal(data, &ei); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	e.Code = ei.Code

	return nil
}

// Error implements a custom error for type ErrorInfo.
func (e *ErrorInfo) Error() string {
	return string(e.data)
}

// Row of data in a table, types of data used by service specified in LogsColumnType
type Row []any

// ISO8601TimeInterval specifies the time range over which to query.
// Use NewISO8601TimeInterval() for help formatting.
// Follows the ISO8601 time interval standard with most common format being startISOTime/endISOTime.
// Use UTC for all times.
type ISO8601TimeInterval string

// NewISO8601TimeInterval creates a ISO8601TimeInterval for use in a query.
// Use UTC for start and end times.
// Start time must be before end time.
func NewISO8601TimeInterval(start time.Time, end time.Time) ISO8601TimeInterval {
	return ISO8601TimeInterval(start.Format(time.RFC3339) + "/" + end.Format(time.RFC3339))
}

// Times returns the interval's start and end times if it's in the format startISOTime/endISOTime, else it will return an error.
func (i ISO8601TimeInterval) Times() (time.Time, time.Time, error) {
	// split into different start and end times
	times := strings.Split(string(i), "/")
	if len(times) != 2 {
		return time.Time{}, time.Time{}, errors.New("time interval should be in format startISOTime/endISOTime")
	}
	start, err := time.Parse(time.RFC3339, times[0])
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("error parsing start time")
	}
	end, err := time.Parse(time.RFC3339, times[1])
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("error parsing end time")
	}
	// return times
	return start, end, nil
}
