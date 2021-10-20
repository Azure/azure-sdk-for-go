//go:build go1.16
// +build go1.16

// Copyright 2017 Microsoft Corporation. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package azcore_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// false positive by linter
func ExampleSetClassifications() { //nolint:govet
	// only log HTTP requests and responses
	log.SetClassifications(log.Request, log.Response)
}

// false positive by linter
func ExampleSetListener() { //nolint:govet
	// a simple logger that writes to stdout
	log.SetListener(func(cls log.Classification, msg string) {
		fmt.Printf("%s: %s\n", cls, msg)
	})
}

type Widget struct {
	Name  *string `json:",omitempty"`
	Count *int    `json:",omitempty"`
}

func (w Widget) MarshalJSON() ([]byte, error) {
	msg := map[string]interface{}{}
	if azcore.IsNullValue(w.Name) {
		msg["name"] = nil
	} else if w.Name != nil {
		msg["name"] = w.Name
	}
	if azcore.IsNullValue(w.Count) {
		msg["count"] = nil
	} else if w.Count != nil {
		msg["count"] = w.Count
	}
	return json.Marshal(msg)
}

func ExampleNullValue() {
	w := Widget{
		Count: azcore.NullValue(0).(*int),
	}
	b, _ := json.Marshal(w)
	fmt.Println(string(b))
	// Output:
	// {"count":null}
}

func ExampleHTTPResponse() {
	pipeline := runtime.NewPipeline("module", "version", nil, nil, nil)
	req, err := runtime.NewRequest(context.Background(), "POST", "https://fakecontainerregisty.azurecr.io/acr/v1/nonexisteng/_tags")
	if err != nil {
		panic(err)
	}
	resp, err := pipeline.Do(req)
	var httpErr azcore.HTTPResponse
	if errors.As(err, &httpErr) {
		// Handle Error
		if httpErr.RawResponse().StatusCode == http.StatusNotFound {
			fmt.Printf("Repository could not be found: %v", httpErr.RawResponse())
		} else if httpErr.RawResponse().StatusCode == http.StatusForbidden {
			fmt.Printf("You do not have permission to access this repository: %v", httpErr.RawResponse())
		} else {
			// ...
		}
	}
	// Do something with response
	fmt.Println(ioutil.ReadAll(resp.Body))
}
