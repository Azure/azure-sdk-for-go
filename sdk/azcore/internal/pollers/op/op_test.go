//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package op

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
)

const (
	fakePollingURL  = "https://foo.bar.baz/status"
	fakePollingURL2 = "https://foo.bar.baz/status/updated"
	fakeLocationURL = "https://foo.bar.baz/location"
	fakeResourceURL = "https://foo.bar.baz/resource"
)

func initialResponse(method string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, fakeResourceURL, nil)
	if err != nil {
		panic(err)
	}
	return &http.Response{
		Body:    ioutil.NopCloser(body),
		Header:  http.Header{},
		Request: req,
	}
}

func createResponse(body io.Reader) *http.Response {
	return &http.Response{
		Body:   ioutil.NopCloser(body),
		Header: http.Header{},
	}
}

func TestApplicable(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}
	if Applicable(resp) {
		t.Fatal("missing Operation-Location should not be applicable")
	}
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	if !Applicable(resp) {
		t.Fatal("having Operation-Location should be applicable")
	}
}

func TestNew(t *testing.T) {
	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("unexpected done")
	}
	if u := poller.FinalGetURL(); u != fakeResourceURL {
		t.Fatalf("unexpected final get URL %s", u)
	}
	if s := poller.Status(); s != pollers.StatusInProgress {
		t.Fatalf("unexpected status %s", s)
	}
	if u := poller.URL(); u != fakePollingURL {
		t.Fatalf("unexpected URL %s", u)
	}
}

func TestNewWithInitialStatus(t *testing.T) {
	resp := initialResponse(http.MethodPut, strings.NewReader(`{ "status": "Updating" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("unexpected done")
	}
	if s := poller.Status(); s != "Updating" {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestNewWithPost(t *testing.T) {
	resp := initialResponse(http.MethodPost, http.NoBody)
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("unexpected done")
	}
	if u := poller.FinalGetURL(); u != fakeLocationURL {
		t.Fatalf("unexpected final get URL %s", u)
	}
}

func TestNewWithDelete(t *testing.T) {
	resp := initialResponse(http.MethodDelete, http.NoBody)
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("unexpected done")
	}
	if u := poller.FinalGetURL(); u != "" {
		t.Fatalf("unexpected final get URL %s", u)
	}
}

func TestNewFail(t *testing.T) {
	resp := initialResponse(http.MethodPut, http.NoBody)
	poller, err := New(resp, "pollerID")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if poller != nil {
		t.Fatal("expected nil poller")
	}
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, "/must/be/absolute")
	poller, err = New(resp, "pollerID")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if poller != nil {
		t.Fatal("expected nil poller")
	}
	resp.Header.Set(shared.HeaderOperationLocation, "/must/be/absolute")
	poller, err = New(resp, "pollerID")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if poller != nil {
		t.Fatal("expected nil poller")
	}
}

func TestUpdateSucceeded(t *testing.T) {
	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	resp = createResponse(strings.NewReader(`{ "status": "Running" }`))
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL2)
	if err := poller.Update(resp); err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("unexpected done")
	}
	if s := poller.Status(); s != "Running" {
		t.Fatalf("unexpected status %s", s)
	}
	if u := poller.URL(); u != fakePollingURL2 {
		t.Fatalf("unexpected URL %s", u)
	}
	resp = createResponse(strings.NewReader(`{ "status": "Succeeded" }`))
	if err := poller.Update(resp); err != nil {
		t.Fatal(err)
	}
	if !poller.Done() {
		t.Fatal("expected done")
	}
	if s := poller.Status(); s != pollers.StatusSucceeded {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestUpdateResourceLocation(t *testing.T) {
	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	resp = createResponse(strings.NewReader(`{ "status": "Succeeded", "resourceLocation": "https://foo.bar.baz/resource2" }`))
	if err := poller.Update(resp); err != nil {
		t.Fatal(err)
	}
	if !poller.Done() {
		t.Fatal("expected done")
	}
	if s := poller.Status(); s != pollers.StatusSucceeded {
		t.Fatalf("unexpected status %s", s)
	}
	if u := poller.FinalGetURL(); u != "https://foo.bar.baz/resource2" {
		t.Fatalf("unexpected final get url %s", u)
	}
}

func TestUpdateFailed(t *testing.T) {
	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	resp = createResponse(strings.NewReader(`{ "status": "Failed" }`))
	if err := poller.Update(resp); err != nil {
		t.Fatal(err)
	}
	if !poller.Done() {
		t.Fatal("expected done")
	}
	if s := poller.Status(); s != pollers.StatusFailed {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestUpdateMissingStatus(t *testing.T) {
	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.Header.Set(shared.HeaderOperationLocation, fakePollingURL)
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	resp = createResponse(http.NoBody)
	if err := poller.Update(resp); err == nil {
		t.Fatal("unexpected nil error")
	}
	if poller.Done() {
		t.Fatal("unexpected done")
	}
}
