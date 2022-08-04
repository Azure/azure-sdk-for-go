//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// testing our custom Error interface for the authentication failed error type
func TestAuthenticationFailedErrorInterface(t *testing.T) {
	urlString := "http://www.test.com"
	resBodyString := "Test Error Response Body"
	url, err := url.Parse(urlString)
	if err != nil {
		t.Fatal("URL parse failed")
	}
	req := &http.Request{Method: "POST", URL: url}
	res := &http.Response{
		Status:     "400 Bad Request",
		StatusCode: 400,
		Body:       ioutil.NopCloser(bytes.NewBufferString(resBodyString)),
		Request:    req,
	}
	err = newAuthenticationFailedError(credNameAzureCLI, "error message", res)
	if e, ok := err.(*AuthenticationFailedError); ok {
		if e.RawResponse == nil {
			t.Fatal("expected a non-nil RawResponse")
		}
	} else {
		t.Fatalf("expected AuthenticationFailedError, received %T", err)
	}
	errMsg := err.Error()
	if !strings.HasPrefix(errMsg, credNameAzureCLI) {
		t.Fatal("missing credential type prefix")
	}
	if !strings.Contains(errMsg, fmt.Sprint(res.StatusCode)) {
		t.Fatalf("unexpected error code: expected error to contain status code %s", fmt.Sprint(res.StatusCode))
	}
	if !strings.Contains(errMsg, res.Request.Method) {
		t.Fatalf("unexpected error method: expected error to contain method %s", res.Request.Method)
	}
	if !strings.Contains(errMsg, urlString) {
		t.Fatalf("unexpected error url: expected error to contain url %s", urlString)
	}
	if !strings.Contains(errMsg, resBodyString) {
		t.Fatalf("unexpected error body: expected error to contain response body %s", resBodyString)
	}
	if !strings.Contains(errMsg, "https://aka.ms/azsdk/go/identity/troubleshoot#azure-cli") {
		t.Fatalf("unexpected error link: expected error to contain troubleshooting link %s", resBodyString)
	}
}
