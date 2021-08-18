// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pollers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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

func TestGetStatusSuccess(t *testing.T) {
	const jsonBody = `{ "status": "InProgress" }`
	resp := azcore.Response{
		Response: &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(jsonBody)),
		},
	}
	status, err := GetStatus(&resp)
	if err != nil {
		t.Fatal(err)
	}
	if status != "InProgress" {
		t.Fatalf("unexpected status %s", status)
	}
}

func TestGetNoBody(t *testing.T) {
	resp := azcore.Response{
		Response: &http.Response{
			Body: http.NoBody,
		},
	}
	status, err := GetStatus(&resp)
	if !errors.Is(err, ErrNoBody) {
		t.Fatalf("unexpected error %T", err)
	}
	if status != "" {
		t.Fatal("expected empty status")
	}
	status, err = GetProvisioningState(&resp)
	if !errors.Is(err, ErrNoBody) {
		t.Fatalf("unexpected error %T", err)
	}
	if status != "" {
		t.Fatal("expected empty status")
	}
}

func TestGetStatusError(t *testing.T) {
	resp := azcore.Response{
		Response: &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{}")),
		},
	}
	status, err := GetStatus(&resp)
	if err != nil {
		t.Fatal(err)
	}
	if status != "" {
		t.Fatalf("expected empty status, got %s", status)
	}
}

func TestGetProvisioningState(t *testing.T) {
	const jsonBody = `{ "properties": { "provisioningState": "Canceled" } }`
	resp := azcore.Response{
		Response: &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(jsonBody)),
		},
	}
	state, err := GetProvisioningState(&resp)
	if err != nil {
		t.Fatal(err)
	}
	if state != "Canceled" {
		t.Fatalf("unexpected status %s", state)
	}
}

func TestGetProvisioningStateError(t *testing.T) {
	resp := azcore.Response{
		Response: &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{}")),
		},
	}
	state, err := GetProvisioningState(&resp)
	if err != nil {
		t.Fatal(err)
	}
	if state != "" {
		t.Fatalf("expected empty provisioning state, got %s", state)
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
