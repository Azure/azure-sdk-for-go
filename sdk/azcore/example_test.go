// +build go1.13

// Copyright 2017 Microsoft Corporation. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package azcore_test

import (
	"context"
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

func ExampleLogger_Should() {
	// you can create your own logging classification as needed
	const LogExpensiveThing azcore.LogClassification = "ExpensiveThing"
	if azcore.Log().Should(LogExpensiveThing) {
		// perform expensive calculation only when enabled
		azcore.Log().Write(LogExpensiveThing, "expensive log message")
	}
}

func ExampleLogger_SetClassifications() {
	// only log HTTP requests and responses
	azcore.Log().SetClassifications(azcore.LogRequest, azcore.LogResponse)
}

func ExampleLogger_SetListener() {
	// a simple logger that writes to stdout
	azcore.Log().SetListener(func(cls azcore.LogClassification, msg string) {
		fmt.Printf("%s: %s\n", cls, msg)
	})
}
