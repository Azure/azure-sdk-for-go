// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestCosmosErrorOnEmptyResponse(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(404))

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)

	var azErr *azcore.ResponseError
	if err := azruntime.NewResponseErrorWithErrorCode(resp, resp.Status); !errors.As(err, &azErr) {
		t.Fatalf("unexpected error type %T", err)
	}
	if azErr.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code %d", azErr.StatusCode)
	}
	if azErr.ErrorCode != "404 Not Found" {
		t.Errorf("unexpected error code %s", azErr.ErrorCode)
	}
	if azErr.RawResponse == nil {
		t.Error("unexpected nil RawResponse")
	}
}

func TestCosmosErrorOnNonJsonBody(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody([]byte("This is not JSON")),
		mock.WithStatusCode(404))

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)

	var azErr *azcore.ResponseError
	if err := azruntime.NewResponseErrorWithErrorCode(resp, resp.Status); !errors.As(err, &azErr) {
		t.Fatalf("unexpected error type %T", err)
	}
	if azErr.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code %d", azErr.StatusCode)
	}
	if azErr.ErrorCode != "404 Not Found" {
		t.Errorf("unexpected error code %s", azErr.ErrorCode)
	}
	if azErr.RawResponse == nil {
		t.Error("unexpected nil RawResponse")
	}
	if !strings.Contains(azErr.Error(), "This is not JSON") {
		t.Error("missing error message")
	}
}

func TestCosmosErrorOnJsonBody(t *testing.T) {
	someError := map[string]string{"Code": "SomeCode"}

	jsonString, err := json.Marshal(someError)
	if err != nil {
		t.Fatal(err)
	}

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithStatusCode(404))

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)

	var azErr *azcore.ResponseError
	err2 := azruntime.NewResponseErrorWithErrorCode(resp, resp.Status)
	assert.Error(t, err2)
	if err := azruntime.NewResponseErrorWithErrorCode(resp, resp.Status); !errors.As(err, &azErr) {
		t.Fatalf("unexpected error type %T", err)
	}
	if azErr.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code %d", azErr.StatusCode)
	}
	if azErr.ErrorCode != "404 Not Found" {
		t.Errorf("unexpected error code %s", azErr.ErrorCode)
	}
	if azErr.RawResponse == nil {
		t.Error("unexpected nil RawResponse")
	}
	if !strings.Contains(azErr.Error(), `"Code": "SomeCode"`) {
		t.Error("missing error JSON")
	}
}
