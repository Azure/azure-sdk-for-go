//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestNewAPIKeyPolicy(t *testing.T) {
	type args struct {
		header string
		cred   KeyCredential
	}
	simpleCred, err := NewKeyCredential("apiKey")
	require.NoError(t, err)

	simpleHeader := "headerName"
	tests := []struct {
		name string
		args args
		want *apiKeyPolicy
	}{
		{
			name: "simple",
			args: args{
				cred:   simpleCred,
				header: simpleHeader,
			},
			want: &apiKeyPolicy{
				header: simpleHeader,
				cred:   simpleCred,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newAPIKeyPolicy(tt.args.cred, tt.args.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIKeyPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIKeyPolicy_Success(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

	cred, err := NewKeyCredential("secret")
	require.NoError(t, err)

	authPolicy := newAPIKeyPolicy(cred, "api-key")
	pipeline := runtime.NewPipeline(
		"testmodule",
		"v0.1.0",
		runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}},
		&policy.ClientOptions{
			Transport: srv,
		})
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		t.Fatalf("Expected nil error but received one")
	}
	if hdrValue := resp.Request.Header.Get("api-key"); hdrValue != "secret" {
		t.Fatalf("expected api-key '%s', got '%s'", "secret", hdrValue)
	}
}
