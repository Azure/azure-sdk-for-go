//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pollers

import (
	"net/http"
	"strings"
	"testing"
)

func TestIsTerminalState(t *testing.T) {
	if IsTerminalState("Updating") {
		t.Fatal("Updating is not a terminal state")
	}
	if !IsTerminalState("Succeeded") {
		t.Fatal("Succeeded is a terminal state")
	}
	if !IsTerminalState("failed") {
		t.Fatal("failed is a terminal state")
	}
	if !IsTerminalState("canceled") {
		t.Fatal("canceled is a terminal state")
	}
}

func TestStatusCodeValid(t *testing.T) {
	if !StatusCodeValid(&http.Response{StatusCode: http.StatusOK}) {
		t.Fatal("unexpected valid code")
	}
	if !StatusCodeValid(&http.Response{StatusCode: http.StatusAccepted}) {
		t.Fatal("unexpected valid code")
	}
	if !StatusCodeValid(&http.Response{StatusCode: http.StatusCreated}) {
		t.Fatal("unexpected valid code")
	}
	if !StatusCodeValid(&http.Response{StatusCode: http.StatusNoContent}) {
		t.Fatal("unexpected valid code")
	}
	if StatusCodeValid(&http.Response{StatusCode: http.StatusPartialContent}) {
		t.Fatal("unexpected valid code")
	}
	if StatusCodeValid(&http.Response{StatusCode: http.StatusBadRequest}) {
		t.Fatal("unexpected valid code")
	}
	if StatusCodeValid(&http.Response{StatusCode: http.StatusInternalServerError}) {
		t.Fatal("unexpected valid code")
	}
}

func TestMakeID(t *testing.T) {
	const (
		pollerID = "pollerID"
		kind     = "kind"
	)
	id := MakeID(pollerID, kind)
	parts := strings.Split(id, idSeparator)
	if l := len(parts); l != 2 {
		t.Fatalf("unexpected length %d", l)
	}
	if p := parts[0]; p != pollerID {
		t.Fatalf("unexpected poller ID %s", p)
	}
	if p := parts[1]; p != kind {
		t.Fatalf("unexpected poller kind %s", p)
	}
}

func TestDecodeID(t *testing.T) {
	_, _, err := DecodeID("")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	_, _, err = DecodeID("invalid_token")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	_, _, err = DecodeID("invalid_token;")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	_, _, err = DecodeID("  ;invalid_token")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	_, _, err = DecodeID("invalid;token;too")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	id, kind, err := DecodeID("pollerID;kind")
	if err != nil {
		t.Fatal(err)
	}
	if id != "pollerID" {
		t.Fatalf("unexpected ID %s", id)
	}
	if kind != "kind" {
		t.Fatalf("unexpected kin %s", kind)
	}
}

func TestIsValidURL(t *testing.T) {
	if IsValidURL("/foo") {
		t.Fatal("unexpected valid URL")
	}
	if !IsValidURL("https://foo.bar/baz") {
		t.Fatal("expected valid URL")
	}
}

func TestFailed(t *testing.T) {
	if Failed("Succeeded") || Failed("Updating") {
		t.Fatal("unexpected failure")
	}
	if !Failed("failed") {
		t.Fatal("expected failure")
	}
}

func TestNopPoller(t *testing.T) {
	np := NopPoller{}
	if !np.Done() {
		t.Fatal("expected done")
	}
	if np.FinalGetURL() != "" {
		t.Fatal("expected empty final get URL")
	}
	if np.Status() != StatusSucceeded {
		t.Fatal("expected Succeeded")
	}
	if np.URL() != "" {
		t.Fatal("expected empty URL")
	}
	if err := np.Update(nil); err != nil {
		t.Fatal(err)
	}
}

/*func TestNewPollerNop(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	resp := initialResponse(http.MethodPost, srv.URL(), strings.NewReader(successResp))
	resp.StatusCode = http.StatusOK
	poller, err := NewPoller("pollerID", "", resp, getPipeline(srv), handleError)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := poller.lro.(*nopPoller); !ok {
		t.Fatalf("unexpected poller type %T", poller.lro)
	}
	tk, err := poller.ResumeToken()
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if tk != "" {
		t.Fatal("expected empty token")
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), 10*time.Millisecond, &result)
	if err != nil {
		t.Fatal(err)
	}
	if v := *result.Field; v != "value" {
		t.Fatalf("unexpected value %s", v)
	}
}*/
