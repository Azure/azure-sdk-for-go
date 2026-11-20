// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"io"
	"math"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManagementNativeTypedThroughputPolicies(t *testing.T) {
	for _, tt := range []struct {
		name       string
		throughput ThroughputProperties
		autoscale  string
	}{
		{name: "full width manual", throughput: ThroughputProperties{ManualThroughput: to(uint64(math.MaxUint64))}},
		{name: "no policy", throughput: ThroughputProperties{Autoscale: &AutoscaleSettings{MaxThroughput: 4000}}, autoscale: `{"maxThroughput":4000}`},
		{name: "empty policy", throughput: ThroughputProperties{Autoscale: &AutoscaleSettings{MaxThroughput: 4000, AutoUpgradePolicy: &AutoscaleAutoUpgradePolicy{}}}, autoscale: `{"maxThroughput":4000,"autoUpgradePolicy":{}}`},
		{name: "explicit zero percent", throughput: ThroughputProperties{Autoscale: &AutoscaleSettings{MaxThroughput: 4000, AutoUpgradePolicy: &AutoscaleAutoUpgradePolicy{ThroughputPolicy: &AutoscaleThroughputPolicy{}}}}, autoscale: `{"maxThroughput":4000,"autoUpgradePolicy":{"throughputPolicy":{"incrementPercent":0}}}`},
		{name: "nonzero percent", throughput: ThroughputProperties{Autoscale: &AutoscaleSettings{MaxThroughput: 4000, AutoUpgradePolicy: &AutoscaleAutoUpgradePolicy{ThroughputPolicy: &AutoscaleThroughputPolicy{IncrementPercent: 10}}}}, autoscale: `{"maxThroughput":4000,"autoUpgradePolicy":{"throughputPolicy":{"incrementPercent":10}}}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			options := &CreateDatabaseOptions{Throughput: &tt.throughput}
			wantHeader := "inherited"
			if !nativeManagementThroughputAvailable {
				options.Operation.CustomHeaders = map[string]string{}
				wantHeader = ""
			}

			server := newMetadataServer(t, func(w http.ResponseWriter, r *http.Request) bool {
				if r.URL.Path != "/dbs" || r.Method != http.MethodPost {
					return false
				}
				require.Equal(t, wantHeader, r.Header.Get("x-test-inherited"))
				if tt.autoscale == "" {
					require.Equal(t, "18446744073709551615", r.Header.Get("x-ms-offer-throughput"))
				} else {
					require.JSONEq(t, tt.autoscale, r.Header.Get("x-ms-cosmos-offer-autopilot-settings"))
				}
				w.WriteHeader(http.StatusCreated)
				_, _ = io.WriteString(w, `{"id":"db"}`)
				return true
			})
			client, err := NewClientWithKey(t.Context(), server.URL, mustKeyCredential(t), &ClientOptions{
				Operation: OperationOptions{CustomHeaders: map[string]string{"x-test-inherited": "inherited"}},
			})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			_, err = client.CreateDatabase(t.Context(), DatabaseProperties{ID: "db"}, options)
			require.NoError(t, err)
		})
	}
}

func TestManagementNativeThroughputRejectsInheritedOppositeMode(t *testing.T) {
	for _, tt := range []struct {
		name       string
		headers    map[string]string
		throughput ThroughputProperties
	}{
		{name: "manual with inherited autoscale", headers: map[string]string{"x-ms-cosmos-offer-autopilot-settings": `{"maxThroughput":4000}`}, throughput: ThroughputProperties{ManualThroughput: to(uint64(400))}},
		{name: "autoscale with inherited manual", headers: map[string]string{"x-ms-offer-throughput": "400"}, throughput: ThroughputProperties{Autoscale: &AutoscaleSettings{MaxThroughput: 4000}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var creates atomic.Int32
			server := newMetadataServer(t, func(w http.ResponseWriter, r *http.Request) bool {
				if r.URL.Path != "/dbs" || r.Method != http.MethodPost {
					return false
				}
				creates.Add(1)
				w.WriteHeader(http.StatusCreated)
				_, _ = io.WriteString(w, `{"id":"db"}`)
				return true
			})
			client, err := NewClientWithKey(t.Context(), server.URL, mustKeyCredential(t), &ClientOptions{
				Operation: OperationOptions{CustomHeaders: tt.headers},
			})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			response, err := client.CreateDatabase(t.Context(), DatabaseProperties{ID: "db"}, &CreateDatabaseOptions{Throughput: &tt.throughput})
			require.Equal(t, DatabaseResponse{}, response)
			var cosmosErr *Error
			require.ErrorAs(t, err, &cosmosErr)
			require.False(t, cosmosErr.FromWire)
			require.Zero(t, creates.Load(), "conflicting throughput modes must not reach the service")
		})
	}
}
