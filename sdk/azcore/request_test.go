// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"net/http"
	"net/url"
	"testing"
)

type testXML struct {
	SomeInt    int
	SomeString string
}

func TestRequestMarshalXML(t *testing.T) {
	u, err := url.Parse("https://contoso.com")
	if err != nil {
		panic(err)
	}
	pl := NewPipeline(nil)
	req := pl.NewRequest(http.MethodPost, *u)
	err = req.MarshalAsXML(testXML{SomeInt: 1, SomeString: "s"})
	if err != nil {
		t.Fatalf("marshal failure: %v", err)
	}
	if ct := req.Header.Get(HeaderContentType); ct != contentTypeAppXML {
		t.Fatalf("unexpected content type, got %s wanted %s", ct, contentTypeAppXML)
	}
	if req.Body == nil {
		t.Fatal("unexpected nil request body")
	}
	if req.ContentLength == 0 {
		t.Fatal("unexpected zero content length")
	}
}
