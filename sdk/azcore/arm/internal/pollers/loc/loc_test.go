//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package loc

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
)

const (
	fakePollingURL1 = "https://foo.bar.baz/status"
	fakePollingURL2 = "https://foo.bar.baz/updated"
)

func initialResponse(method string) *http.Response {
	return &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusAccepted,
	}
}

func pollingResponse(statusCode int, body io.Reader) *http.Response {
	return &http.Response{
		Body:       ioutil.NopCloser(body),
		Header:     http.Header{},
		StatusCode: statusCode,
	}
}

func TestApplicable(t *testing.T) {
	resp := &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusAccepted,
	}
	if Applicable(resp) {
		t.Fatal("missing Location should not be applicable")
	}
	resp.Header.Set(shared.HeaderLocation, fakePollingURL1)
	if !Applicable(resp) {
		t.Fatal("having Location should be applicable")
	}
}

func TestNew(t *testing.T) {
	resp := initialResponse(http.MethodPut)
	resp.Header.Set(shared.HeaderLocation, fakePollingURL1)
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
	if u := poller.URL(); u != fakePollingURL1 {
		t.Fatalf("unexpected polling URL %s", u)
	}
	pr := pollingResponse(http.StatusAccepted, http.NoBody)
	pr.Header.Set(shared.HeaderLocation, fakePollingURL2)
	if err := poller.Update(pr); err != nil {
		t.Fatal(err)
	}
	if u := poller.URL(); u != fakePollingURL2 {
		t.Fatalf("unexpected polling URL %s", u)
	}
	if err := poller.Update(pollingResponse(http.StatusNoContent, http.NoBody)); err != nil {
		t.Fatal(err)
	}
	if s := poller.Status(); s != "Succeeded" {
		t.Fatalf("unexpected status %s", s)
	}
	if err := poller.Update(pollingResponse(http.StatusConflict, http.NoBody)); err != nil {
		t.Fatal(err)
	}
	if s := poller.Status(); s != "Failed" {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestUpdateWithProvState(t *testing.T) {
	resp := initialResponse(http.MethodPut)
	resp.Header.Set(shared.HeaderLocation, fakePollingURL1)
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
	if u := poller.URL(); u != fakePollingURL1 {
		t.Fatalf("unexpected polling URL %s", u)
	}
	pr := pollingResponse(http.StatusAccepted, http.NoBody)
	pr.Header.Set(shared.HeaderLocation, fakePollingURL2)
	if err := poller.Update(pr); err != nil {
		t.Fatal(err)
	}
	if u := poller.URL(); u != fakePollingURL2 {
		t.Fatalf("unexpected polling URL %s", u)
	}
	if err := poller.Update(pollingResponse(http.StatusOK, strings.NewReader(`{ "properties": { "provisioningState": "Updating" } }`))); err != nil {
		t.Fatal(err)
	}
	if s := poller.Status(); s != "Updating" {
		t.Fatalf("unexpected status %s", s)
	}
	if err := poller.Update(pollingResponse(http.StatusOK, http.NoBody)); err != nil {
		t.Fatal(err)
	}
	if s := poller.Status(); s != "Succeeded" {
		t.Fatalf("unexpected status %s", s)
	}
}
