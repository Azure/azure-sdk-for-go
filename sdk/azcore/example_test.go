// +build go1.13

// Copyright 2017 Microsoft Corporation. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package azcore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func ExamplePipeline_Do() {
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, "https://github.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	pipeline := azcore.NewPipeline(nil)
	resp, err := pipeline.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", robots)
}

func ExampleRequest_SetBody() {
	req, err := azcore.NewRequest(context.Background(), http.MethodPut, "https://contoso.com/some/endpoint")
	if err != nil {
		log.Fatal(err)
	}
	body := strings.NewReader("this is seekable content to be uploaded")
	err = req.SetBody(azcore.NopCloser(body), "text/plain")
	if err != nil {
		log.Fatal(err)
	}
}

// false positive by linter
func ExampleLogSetClassifications() { //nolint:govet
	// only log HTTP requests and responses
	azcore.LogSetClassifications(azcore.LogRequest, azcore.LogResponse)
}

// false positive by linter
func ExampleLogSetListener() { //nolint:govet
	// a simple logger that writes to stdout
	azcore.LogSetListener(func(cls azcore.LogClassification, msg string) {
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
