//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package async

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
)

const (
	fakePollingURL  = "https://foo.bar.baz/status"
	fakeResourceURL = "https://foo.bar.baz/resource"
)

func initialResponse(method string, resp io.Reader) *http.Response {
	req, err := http.NewRequest(method, fakeResourceURL, nil)
	if err != nil {
		panic(err)
	}
	return &http.Response{
		Body:    ioutil.NopCloser(resp),
		Header:  http.Header{},
		Request: req,
	}
}

func pollingResponse(resp io.Reader) *http.Response {
	return &http.Response{
		Body:   ioutil.NopCloser(resp),
		Header: http.Header{},
	}
}

func TestApplicable(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}
	if Applicable(resp) {
		t.Fatal("missing Azure-AsyncOperation should not be applicable")
	}
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	if !Applicable(resp) {
		t.Fatal("having Azure-AsyncOperation should be applicable")
	}
}

func TestNew(t *testing.T) {
	const jsonBody = `{ "properties": { "provisioningState": "Started" } }`
	resp := initialResponse(http.MethodPut, strings.NewReader(jsonBody))
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New(resp, "", "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("poller should not be done")
	}
	if u := poller.FinalGetURL(); u != fakeResourceURL {
		t.Fatalf("unexpected final get URL %s", u)
	}
	if s := poller.Status(); s != "Started" {
		t.Fatalf("unexpected status %s", s)
	}
	if u := poller.URL(); u != fakePollingURL {
		t.Fatalf("unexpected polling URL %s", u)
	}
	if err := poller.Update(pollingResponse(strings.NewReader(`{ "status": "InProgress" }`))); err != nil {
		t.Fatal(err)
	}
	if s := poller.Status(); s != "InProgress" {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestNewDeleteNoProvState(t *testing.T) {
	resp := initialResponse(http.MethodDelete, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New(resp, "", "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("poller should not be done")
	}
	if s := poller.Status(); s != "InProgress" {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestNewPutNoProvState(t *testing.T) {
	// missing provisioning state on initial response
	// NOTE: ARM RPC forbids this but we allow it for back-compat
	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New(resp, "", "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("poller should not be done")
	}
	if s := poller.Status(); s != "InProgress" {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestNewFinalGetLocation(t *testing.T) {
	const (
		jsonBody = `{ "properties": { "provisioningState": "Started" } }`
		locURL   = "https://foo.bar.baz/location"
	)
	resp := initialResponse(http.MethodPost, strings.NewReader(jsonBody))
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, locURL)
	poller, err := New(resp, "location", "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("poller should not be done")
	}
	if u := poller.FinalGetURL(); u != locURL {
		t.Fatalf("unexpected final get URL %s", u)
	}
	if u := poller.URL(); u != fakePollingURL {
		t.Fatalf("unexpected polling URL %s", u)
	}
}

func TestNewFinalGetOrigin(t *testing.T) {
	const (
		jsonBody = `{ "properties": { "provisioningState": "Started" } }`
		locURL   = "https://foo.bar.baz/location"
	)
	resp := initialResponse(http.MethodPost, strings.NewReader(jsonBody))
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, locURL)
	poller, err := New(resp, "original-uri", "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("poller should not be done")
	}
	if u := poller.FinalGetURL(); u != fakeResourceURL {
		t.Fatalf("unexpected final get URL %s", u)
	}
	if u := poller.URL(); u != fakePollingURL {
		t.Fatalf("unexpected polling URL %s", u)
	}
}

func TestNewPutNoProvStateOnUpdate(t *testing.T) {
	// missing provisioning state on initial response
	// NOTE: ARM RPC forbids this but we allow it for back-compat
	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderAzureAsync, fakePollingURL)
	poller, err := New(resp, "", "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("poller should not be done")
	}
	if s := poller.Status(); s != "InProgress" {
		t.Fatalf("unexpected status %s", s)
	}
	if err := poller.Update(pollingResponse(strings.NewReader("{}"))); err == nil {
		t.Fatal("unexpected nil error")
	}
}
