// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package loc

import (
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	fakePollingURL  = "https://foo.bar.baz/status"
	fakeResourceURL = "https://foo.bar.baz/resource"
)

func initialResponse(method string) *azcore.Response {
	return &azcore.Response{
		Response: &http.Response{
			Header:     http.Header{},
			StatusCode: http.StatusAccepted,
		},
	}
}

func pollingResponse(statusCode int) *azcore.Response {
	return &azcore.Response{
		Response: &http.Response{
			Header:     http.Header{},
			StatusCode: statusCode,
		},
	}
}

func TestApplicable(t *testing.T) {
	resp := azcore.Response{
		Response: &http.Response{
			Header:     http.Header{},
			StatusCode: http.StatusAccepted,
		},
	}
	if Applicable(&resp) {
		t.Fatal("missing Location should not be applicable")
	}
	resp.Response.Header.Set(pollers.HeaderLocation, fakePollingURL)
	if !Applicable(&resp) {
		t.Fatal("having Location should be applicable")
	}
}

func TestNew(t *testing.T) {
	resp := initialResponse(http.MethodPut)
	resp.Header.Set(pollers.HeaderLocation, fakePollingURL)
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("poller should not be done")
	}
	if u := poller.FinalGetURL(); u != "" {
		t.Fatal("expected empty final GET URL")
	}
	if s := poller.Status(); s != "InProgress" {
		t.Fatalf("unexpected status %s", s)
	}
	if u := poller.URL(); u != fakePollingURL {
		t.Fatalf("unexpected polling URL %s", u)
	}
	if err := poller.Update(pollingResponse(http.StatusNoContent)); err != nil {
		t.Fatal(err)
	}
	if s := poller.Status(); s != "Succeeded" {
		t.Fatalf("unexpected status %s", s)
	}
	if err := poller.Update(pollingResponse(http.StatusConflict)); err != nil {
		t.Fatal(err)
	}
	if s := poller.Status(); s != "Failed" {
		t.Fatalf("unexpected status %s", s)
	}
}
