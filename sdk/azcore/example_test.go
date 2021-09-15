//go:build go1.16
// +build go1.16

// Copyright 2017 Microsoft Corporation. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package azcore_test

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
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
