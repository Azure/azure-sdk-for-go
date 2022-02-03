//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package loc

import (
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
)

const (
	fakeLocationURL  = "https://foo.bar.baz/status"
	fakeLocationURL2 = "https://foo.bar.baz/status/other"
)

func initialResponse() *http.Response {
	return &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusAccepted,
	}
}

func TestApplicable(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}
	if Applicable(resp) {
		t.Fatal("missing Location should not be applicable")
	}
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	if !Applicable(resp) {
		t.Fatal("having Location should be applicable")
	}
}

func TestNew(t *testing.T) {
	resp := initialResponse()
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
	if s := poller.Status(); s != pollers.StatusInProgress {
		t.Fatalf("unexpected status %s", s)
	}
	if u := poller.URL(); u != fakeLocationURL {
		t.Fatalf("unexpected polling URL %s", u)
	}
}

func TestNewFail(t *testing.T) {
	resp := initialResponse()
	poller, err := New(resp, "pollerID")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if poller != nil {
		t.Fatal("expected nil poller")
	}
	resp.Header.Set(shared.HeaderLocation, "/must/be/absolute")
	poller, err = New(resp, "pollerID")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if poller != nil {
		t.Fatal("expected nil poller")
	}
}

func TestUpdateSucceeded(t *testing.T) {
	resp := initialResponse()
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL2)
	if err := poller.Update(resp); err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("unexpected done")
	}
	if u := poller.URL(); u != fakeLocationURL2 {
		t.Fatalf("unexpected polling URL %s", u)
	}
	if err := poller.Update(&http.Response{StatusCode: http.StatusOK}); err != nil {
		t.Fatal(err)
	}
	if !poller.Done() {
		t.Fatal("expected done")
	}
	if s := poller.Status(); s != pollers.StatusSucceeded {
		t.Fatalf("unexpected status %s", s)
	}
}

func TestUpdateFailed(t *testing.T) {
	resp := initialResponse()
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL)
	poller, err := New(resp, "pollerID")
	if err != nil {
		t.Fatal(err)
	}
	resp.Header.Set(shared.HeaderLocation, fakeLocationURL2)
	if err := poller.Update(resp); err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("unexpected done")
	}
	if u := poller.URL(); u != fakeLocationURL2 {
		t.Fatalf("unexpected polling URL %s", u)
	}
	if err := poller.Update(&http.Response{StatusCode: http.StatusConflict}); err != nil {
		t.Fatal(err)
	}
	if !poller.Done() {
		t.Fatal("expected done")
	}
	if s := poller.Status(); s != pollers.StatusFailed {
		t.Fatalf("unexpected status %s", s)
	}
}
