// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package body

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	fakeResourceURL = "https://foo.bar.baz/resource"
)

func initialResponse(method string, resp io.Reader) *azcore.Response {
	req, err := http.NewRequest(method, fakeResourceURL, nil)
	if err != nil {
		panic(err)
	}
	return &azcore.Response{
		Response: &http.Response{
			Body:    ioutil.NopCloser(resp),
			Header:  http.Header{},
			Request: req,
		},
	}
}

func pollingResponse(resp io.Reader) *azcore.Response {
	return &azcore.Response{
		Response: &http.Response{
			Body:   ioutil.NopCloser(resp),
			Header: http.Header{},
		},
	}
}

func TestApplicable(t *testing.T) {
	resp := azcore.Response{
		Response: &http.Response{
			Header: http.Header{},
			Request: &http.Request{
				Method: http.MethodDelete,
			},
		},
	}
	if Applicable(&resp) {
		t.Fatal("method DELETE should not be applicable")
	}
	resp.Request.Method = http.MethodPatch
	if !Applicable(&resp) {
		t.Fatal("method PATCH should be applicable")
	}
	resp.Request.Method = http.MethodPut
	if !Applicable(&resp) {
		t.Fatal("method PUT should be applicable")
	}
}

func TestNew(t *testing.T) {
	const jsonBody = `{ "properties": { "provisioningState": "Started" } }`
	resp := initialResponse(http.MethodPut, strings.NewReader(jsonBody))
	resp.StatusCode = http.StatusCreated
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
	if s := poller.Status(); s != "Started" {
		t.Fatalf("unexpected status %s", s)
	}
	if u := poller.URL(); u != fakeResourceURL {
		t.Fatalf("unexpected polling URL %s", u)
	}
	if err := poller.Update(pollingResponse(strings.NewReader(`{ "properties": { "provisioningState": "InProgress" } }`))); err != nil {
		t.Fatal(err)
	}
	if s := poller.Status(); s != "InProgress" {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestUpdateNoProvState(t *testing.T) {
	const jsonBody = `{ "properties": { "provisioningState": "Started" } }`
	resp := initialResponse(http.MethodPut, strings.NewReader(jsonBody))
	resp.StatusCode = http.StatusOK
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
	if s := poller.Status(); s != "Started" {
		t.Fatalf("unexpected status %s", s)
	}
	if u := poller.URL(); u != fakeResourceURL {
		t.Fatalf("unexpected polling URL %s", u)
	}
	if err := poller.Update(pollingResponse(http.NoBody)); err != nil {
		t.Fatal(err)
	}
	if s := poller.Status(); s != "Succeeded" {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestNewNoInitialProvStateOK(t *testing.T) {
	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.StatusCode = http.StatusOK
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if !poller.Done() {
		t.Fatal("poller not be done")
	}
	if u := poller.FinalGetURL(); u != "" {
		t.Fatal("expected empty final GET URL")
	}
	if s := poller.Status(); s != "Succeeded" {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestNewNoInitialProvStateNC(t *testing.T) {
	resp := initialResponse(http.MethodPut, http.NoBody)
	resp.StatusCode = http.StatusNoContent
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	if !poller.Done() {
		t.Fatal("poller not be done")
	}
	if u := poller.FinalGetURL(); u != "" {
		t.Fatal("expected empty final GET URL")
	}
	if s := poller.Status(); s != "Succeeded" {
		t.Fatalf("unexpected status %s", s)
	}
}
