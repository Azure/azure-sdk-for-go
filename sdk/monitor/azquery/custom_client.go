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
// SORT IN AMOUNT OF TIME ORDER
const (
	SevenDays        ISO8601Duration = "P7D"
	ThreeDays        ISO8601Duration = "P3D"
	TwoDays          Duration        = "P2D"
	OneDay           Duration        = "P1D"
	OneHour          Duration        = "PT1H"
	FourHours        Duration        = "PT4H"
	TwentyFourHours  Duration        = "PT24H"
	FourtyEightHours Duration        = "PT48H"
	ThirtyMinutes    Duration        = "PT30M"
	FiveMinutes      Duration        = "PT5M"
)

// ADD CONSTRUCTOR??

// TIME>TIME instead of string
type ISO8601TimeInterval struct {
	start    time.Time
	end      time.Time
	duration ISO8601Duration

	//
}

// not general purpose enough for round trip (if this was general purpose), how can the customer know what they can look at
// DURATION happened before- keyvault cert experation, go check
// once you create, can't look at what's in it, (add string??)
// long naming- methods on the type instead?? SetStartEnd
// TODO- show to team

// show both, round trip vs not round trip

func NewISO8601TimeIntervalFromStartEnd(start time.Time, end time.Time) *ISO8601TimeInterval {
	return &ISO8601TimeInterval{start: start, end: end}
}

// Timespan in the format start_time/duration
func NewISO8601TimeIntervalFromStartDuration(start time.Time, duration ISO8601Duration) *ISO8601TimeInterval {
	return &ISO8601TimeInterval{start: start, duration: duration}
}

func NewISO8601TimeIntervalFromDurationEnd(duration ISO8601Duration, end time.Time) *ISO8601TimeInterval {
	return &ISO8601TimeInterval{duration: duration, end: end}
}

// timespan in the format duration
func NewISO8601TimeIntervalFromDuration(duration ISO8601Duration) *ISO8601TimeInterval {
	return &ISO8601TimeInterval{duration: duration}
}

// ADD SOME ERROR CHECKING
func (t Timespan) MarshalJSON() ([]byte, error) {
	var timespan string
	if t.StartTime != "" && t.EndTime != "" {
		timespan = t.StartTime + "/" + t.EndTime
	} else if t.StartTime != "" && t.Duration != "" {
		timespan = t.StartTime + "/" + string(t.Duration)
	} else if t.Duration != "" && t.EndTime != "" {
		timespan = string(t.Duration) + "/" + t.EndTime
	}
	if t.Duration != "" {
		timespan = string(t.Duration)
	}
	return json.Marshal(timespan)
}
