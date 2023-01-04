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

type ISO8601Duration string

// Common ISO8601 durations
// NEED TO RENAME
const (
	SevenDays        ISO8601Duration = "P7D"
	ThreeDays        ISO8601Duration = "P3D"
	TwoDays          ISO8601Duration = "P2D"
	OneDay           ISO8601Duration = "P1D"
	OneHour          ISO8601Duration = "PT1H"
	FourHours        ISO8601Duration = "PT4H"
	TwentyFourHours  ISO8601Duration = "PT24H"
	FourtyEightHours ISO8601Duration = "PT48H"
	ThirtyMinutes    ISO8601Duration = "PT30M"
	FiveMinutes      ISO8601Duration = "PT5M"
)

// Don't use this type directly, use NewISO8601TimeIntervalFromStartEnd() instead.
type ISO8601TimeInterval struct {
	start    time.Time
	end      time.Time
	duration ISO8601Duration
}

func (i ISO8601TimeInterval) String() string {
	var timespan string
	if !i.start.IsZero() && !i.end.IsZero() {
		timespan = i.start.Format(time.RFC3339) + "/" + i.end.Format(time.RFC3339)
	} else if !i.start.IsZero() && i.duration != "" {
		timespan = i.start.Format(time.RFC3339) + "/" + string(i.duration)
	} else if i.duration != "" && !i.end.IsZero() {
		timespan = string(i.duration) + "/" + i.end.Format(time.RFC3339)
	} else if i.duration != "" {
		timespan = string(i.duration)
	}

	return timespan
}

func NewISO8601TimeIntervalFromStartEnd(start time.Time, end time.Time) (*ISO8601TimeInterval, error) {
	if start.IsZero() {
		return nil, errors.New("start time is zero")
	}
	if end.IsZero() {
		return nil, errors.New("end time is zero")
	}
	if end.Before(start) {
		return nil, errors.New("end time occurs before start time")
	}

	return &ISO8601TimeInterval{start: start, end: end}, nil
}

// Timespan in the format start_time/duration
func NewISO8601TimeIntervalFromStartDuration(start time.Time, duration ISO8601Duration) (*ISO8601TimeInterval, error) {
	if start.IsZero() {
		return nil, errors.New("start time is zero")
	}
	if duration == "" {
		return nil, errors.New("duration string is empty")
	}
	return &ISO8601TimeInterval{start: start, duration: duration}, nil
}

func NewISO8601TimeIntervalFromDurationEnd(duration ISO8601Duration, end time.Time) (*ISO8601TimeInterval, error) {
	if duration == "" {
		return nil, errors.New("duration string is empty")
	}
	if end.IsZero() {
		return nil, errors.New("end time is zero")
	}
	return &ISO8601TimeInterval{duration: duration, end: end}, nil
}

// NewISO8601TimeIntervalFromDuration creates a ISO8601TimeInterval used to specify the timespan for a logs query
// Uses
func NewISO8601TimeIntervalFromDuration(duration ISO8601Duration) (*ISO8601TimeInterval, error) {
	if duration == "" {
		return nil, errors.New("duration string is empty")
	}
	return &ISO8601TimeInterval{duration: duration}, nil
}

func (i ISO8601TimeInterval) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}
