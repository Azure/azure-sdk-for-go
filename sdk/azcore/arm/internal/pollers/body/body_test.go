//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package body

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
)

const (
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

func pollingResponse(status int, resp io.Reader) *http.Response {
	return &http.Response{
		Body:       ioutil.NopCloser(resp),
		Header:     http.Header{},
		StatusCode: status,
	}
}

func TestApplicable(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
		Request: &http.Request{
			Method: http.MethodDelete,
		},
	}
	if Applicable(resp) {
		t.Fatal("method DELETE should not be applicable")
	}
	resp.Request.Method = http.MethodPatch
	if !Applicable(resp) {
		t.Fatal("method PATCH should be applicable")
	}
	resp.Request.Method = http.MethodPut
	if !Applicable(resp) {
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
	if err := poller.Update(pollingResponse(http.StatusOK, strings.NewReader(`{ "properties": { "provisioningState": "InProgress" } }`))); err != nil {
		t.Fatal(err)
	}
	if s := poller.Status(); s != "InProgress" {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestUpdateNoProvStateFail(t *testing.T) {
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
	err = poller.Update(pollingResponse(http.StatusOK, http.NoBody))
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if !errors.Is(err, shared.ErrNoBody) {
		t.Fatalf("unexpected error type %T", err)
	}
}

func TestUpdateNoProvStateSuccess(t *testing.T) {
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
	err = poller.Update(pollingResponse(http.StatusOK, strings.NewReader(`{}`)))
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateNoProvState204(t *testing.T) {
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
	err = poller.Update(pollingResponse(http.StatusNoContent, http.NoBody))
	if err != nil {
		t.Fatal(err)
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
