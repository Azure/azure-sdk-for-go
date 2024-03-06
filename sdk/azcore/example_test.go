//go:build go1.18
// +build go1.18

// Copyright 2017 Microsoft Corporation. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package azcore_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// false positive by linter
func ExampleSetEvents() {
	// only log HTTP requests and responses
	log.SetEvents(log.EventRequest, log.EventResponse)
}

// false positive by linter
func ExampleSetListener() {
	// a simple logger that writes to stdout
	log.SetListener(func(cls log.Event, msg string) {
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
		Count: azcore.NullValue[*int](),
	}
	b, _ := json.Marshal(w)
	fmt.Println(string(b))
	// Output:
	// {"count":null}
}

func ExampleResponseError() {
	pipeline := runtime.NewPipeline("module", "version", runtime.PipelineOptions{}, nil)
	req, err := runtime.NewRequest(context.Background(), "POST", "https://fakecontainerregisty.azurecr.io/acr/v1/nonexisteng/_tags")
	if err != nil {
		panic(err)
	}
	resp, err := pipeline.Do(req)
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		// Handle Error
		if respErr.StatusCode == http.StatusNotFound {
			fmt.Printf("Repository could not be found: %v", respErr)
		} else if respErr.StatusCode == http.StatusForbidden {
			fmt.Printf("You do not have permission to access this repository: %v", respErr)
		} else {
			// ...
		}
	}
	// Do something with response
	fmt.Println(io.ReadAll(resp.Body))
}
