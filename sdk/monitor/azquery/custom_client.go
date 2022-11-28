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

// Table - Contains the columns and rows for one table in a query response.
type Table struct {
	// REQUIRED; The list of columns in this table.
	Columns []*Column `json:"columns,omitempty"`

	// REQUIRED; The name of the table.
	Name *string `json:"name,omitempty"`

	// REQUIRED; The resulting rows from this query.
	Rows []Row `json:"rows,omitempty"`

	// maps column name to index for easy lookup, helper for accessing Row data
	ColumnIndexLookup map[string]int `json:"-"`
}

// UnmarshalJSON implements the json.Unmarshaller interface for type Table.
func (t *Table) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", t, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "columns":
			err = unpopulate(val, "Columns", &t.Columns)
			delete(rawMsg, key)
			t.ColumnIndexLookup = map[string]int{}
			for i, v := range t.Columns {
				t.ColumnIndexLookup[*v.Name] = i
			}
		case "name":
			err = unpopulate(val, "Name", &t.Name)
			delete(rawMsg, key)
		case "rows":
			err = unpopulate(val, "Rows", &t.Rows)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", t, err)
		}
	}
	return nil
}
