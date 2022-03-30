//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pollers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
)

func TestGetStatusSuccess(t *testing.T) {
	const jsonBody = `{ "status": "InProgress" }`
	resp := &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(jsonBody)),
	}
	status, err := GetStatus(resp)
	if err != nil {
		t.Fatal(err)
	}
	if status != "InProgress" {
		t.Fatalf("unexpected status %s", status)
	}
}

func TestGetNoBody(t *testing.T) {
	resp := &http.Response{
		Body: http.NoBody,
	}
	status, err := GetStatus(resp)
	if !errors.Is(err, shared.ErrNoBody) {
		t.Fatalf("unexpected error %T", err)
	}
	if status != "" {
		t.Fatal("expected empty status")
	}
	status, err = GetProvisioningState(resp)
	if !errors.Is(err, shared.ErrNoBody) {
		t.Fatalf("unexpected error %T", err)
	}
	if status != "" {
		t.Fatal("expected empty status")
	}
}

func TestGetStatusError(t *testing.T) {
	resp := &http.Response{
		Body: ioutil.NopCloser(strings.NewReader("{}")),
	}
	status, err := GetStatus(resp)
	if err != nil {
		t.Fatal(err)
	}
	if status != "" {
		t.Fatalf("expected empty status, got %s", status)
	}
}

func TestGetProvisioningState(t *testing.T) {
	const jsonBody = `{ "properties": { "provisioningState": "Canceled" } }`
	resp := &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(jsonBody)),
	}
	state, err := GetProvisioningState(resp)
	if err != nil {
		t.Fatal(err)
	}
	if state != "Canceled" {
		t.Fatalf("unexpected status %s", state)
	}
}

func TestGetProvisioningStateError(t *testing.T) {
	resp := &http.Response{
		Body: ioutil.NopCloser(strings.NewReader("{}")),
	}
	state, err := GetProvisioningState(resp)
	if err != nil {
		t.Fatal(err)
	}
	if state != "" {
		t.Fatalf("expected empty provisioning state, got %s", state)
	}
}
