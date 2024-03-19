//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestNullValue(t *testing.T) {
	v := NullValue[*string]()
	vv := NullValue[*string]()
	if v != vv {
		t.Fatal("null values should match for the same types")
	}
}

func TestIsNullValue(t *testing.T) {
	if IsNullValue("") {
		t.Fatal("string literal can't be a null value")
	}
	s := ""
	if IsNullValue(&s) {
		t.Fatal("&s isn't a null value")
	}
	var i *int
	if IsNullValue(i) {
		t.Fatal("i isn't a null value")
	}
	i = NullValue[*int]()
	if !IsNullValue(i) {
		t.Fatal("expected null value for i")
	}
	i2 := 0
	i = &i2
	if IsNullValue(i) {
		t.Fatal("i should no longer be null value")
	}
}

func TestNullValueMapSlice(t *testing.T) {
	v := NullValue[[]string]()
	vv := NullValue[[]string]()
	if reflect.TypeOf(v) != reflect.TypeOf(vv) {
		t.Fatal("null values should match for the same types")
	}
	m := NullValue[map[string]int]()
	if reflect.TypeOf(v) == reflect.TypeOf(m) {
		t.Fatal("null values for string and int should not match")
	}
}

func TestIsNullValueMapSlice(t *testing.T) {
	if IsNullValue([]string{}) {
		t.Fatal("slice literal can't be a null value")
	}
	if IsNullValue(map[int]string{}) {
		t.Fatal("map literal can't be a null value")
	}
	s := NullValue[[]int]()
	if !IsNullValue(s) {
		t.Fatal("expected null value for s")
	}
	m := NullValue[map[string]any]()
	if !IsNullValue(m) {
		t.Fatal("expected null value for s")
	}

	type nullFields struct {
		Map   map[string]int
		Slice []string
	}

	nf := nullFields{}
	if IsNullValue(nf.Map) {
		t.Fatal("unexpected null map")
	}
	if IsNullValue(nf.Slice) {
		t.Fatal("unexpected null slice")
	}

	nf.Map = map[string]int{}
	nf.Slice = []string{}
	if IsNullValue(nf.Map) {
		t.Fatal("unexpected null map")
	}
	if IsNullValue(nf.Slice) {
		t.Fatal("unexpected null slice")
	}

	nf.Map = NullValue[map[string]int]()
	nf.Slice = NullValue[[]string]()
	if !IsNullValue(nf.Map) {
		t.Fatal("expected null map")
	}
	if !IsNullValue(nf.Slice) {
		t.Fatal("expected null slice")
	}
}

func TestNewClient(t *testing.T) {
	client, err := NewClient("package.Client", "v1.0.0", runtime.PipelineOptions{}, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotZero(t, client.Pipeline())
	require.Zero(t, client.Tracer())

	client, err = NewClient("package.Client", "", runtime.PipelineOptions{}, &ClientOptions{
		Telemetry: policy.TelemetryOptions{
			Disabled: true,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNewClientTracingEnabled(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()

	var attrString string
	client, err := NewClient("package.Client", "v1.0.0", runtime.PipelineOptions{
		Tracing: runtime.TracingOptions{
			Namespace: "Widget.Factory",
		},
	}, &policy.ClientOptions{
		TracingProvider: tracing.NewProvider(func(name, version string) tracing.Tracer {
			return tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
				require.NotNil(t, options)
				for _, attr := range options.Attributes {
					if attr.Key == shared.TracingNamespaceAttrName {
						v, ok := attr.Value.(string)
						require.True(t, ok)
						attrString = attr.Key + ":" + v
					}
				}
				return ctx, tracing.Span{}
			}, nil)
		}, nil),
		Transport: srv,
	})
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotZero(t, client.Pipeline())
	require.NotZero(t, client.Tracer())

	const requestEndpoint = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/fakeResourceGroupo/providers/Microsoft.Storage/storageAccounts/fakeAccountName"
	req, err := exported.NewRequest(context.WithValue(context.Background(), shared.CtxWithTracingTracer{}, client.Tracer()), http.MethodGet, srv.URL()+requestEndpoint)
	require.NoError(t, err)
	srv.AppendResponse()
	_, err = client.Pipeline().Do(req)
	require.NoError(t, err)
	require.EqualValues(t, "az.namespace:Widget.Factory", attrString)
}

func TestClientWithClientName(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()

	var clientName string
	var modVersion string
	var attrString string
	client, err := NewClient("module", "v1.0.0", runtime.PipelineOptions{
		Tracing: runtime.TracingOptions{
			Namespace: "Widget.Factory",
		},
	}, &policy.ClientOptions{
		TracingProvider: tracing.NewProvider(func(name, version string) tracing.Tracer {
			clientName = name
			modVersion = version
			return tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
				require.NotNil(t, options)
				for _, attr := range options.Attributes {
					if attr.Key == shared.TracingNamespaceAttrName {
						v, ok := attr.Value.(string)
						require.True(t, ok)
						attrString = attr.Key + ":" + v
					}
				}
				return ctx, tracing.Span{}
			}, nil)
		}, nil),
		Transport: srv,
	})
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotZero(t, client.Pipeline())
	require.NotZero(t, client.Tracer())
	require.EqualValues(t, "module", clientName)
	require.EqualValues(t, "v1.0.0", modVersion)

	const requestEndpoint = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/fakeResourceGroupo/providers/Microsoft.Storage/storageAccounts/fakeAccountName"
	req, err := exported.NewRequest(context.WithValue(context.Background(), shared.CtxWithTracingTracer{}, client.Tracer()), http.MethodGet, srv.URL()+requestEndpoint)
	require.NoError(t, err)
	srv.SetResponse()
	_, err = client.Pipeline().Do(req)
	require.NoError(t, err)
	require.EqualValues(t, "az.namespace:Widget.Factory", attrString)

	newClient := client.WithClientName("other")
	require.EqualValues(t, "other", clientName)
	require.EqualValues(t, "v1.0.0", modVersion)
	require.EqualValues(t, client.Pipeline(), newClient.Pipeline())
	_, err = newClient.Pipeline().Do(req)
	require.NoError(t, err)
	require.EqualValues(t, "az.namespace:Widget.Factory", attrString)
}

func TestNewKeyCredential(t *testing.T) {
	require.NotNil(t, NewKeyCredential("foo"))
}

func TestNewSASCredential(t *testing.T) {
	require.NotNil(t, NewSASCredential("foo"))
}
