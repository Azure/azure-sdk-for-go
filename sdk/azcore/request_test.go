// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"unsafe"
)

type testJSON struct {
	SomeInt    int
	SomeString string
}

type testXML struct {
	SomeInt    int
	SomeString string
}

func TestRequestMarshalXML(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, "https://contoso.com")
	if err != nil {
		t.Fatal(err)
	}
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

func TestRequestEmptyPipeline(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, "https://contoso.com")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := req.Next()
	if resp != nil {
		t.Fatal("expected nil response")
	}
	if !errors.Is(err, ErrNoMorePolicies) {
		t.Fatalf("expected ErrNoMorePolicies, got %v", err)
	}
}

func TestRequestMarshalJSON(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, "https://contoso.com")
	if err != nil {
		t.Fatal(err)
	}
	err = req.MarshalAsJSON(testJSON{SomeInt: 1, SomeString: "s"})
	if err != nil {
		t.Fatalf("marshal failure: %v", err)
	}
	if ct := req.Header.Get(HeaderContentType); ct != contentTypeAppJSON {
		t.Fatalf("unexpected content type, got %s wanted %s", ct, contentTypeAppJSON)
	}
	if req.Body == nil {
		t.Fatal("unexpected nil request body")
	}
	if req.ContentLength == 0 {
		t.Fatal("unexpected zero content length")
	}
}

func TestRequestMarshalAsByteArrayURLFormat(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, "https://contoso.com")
	if err != nil {
		t.Fatal(err)
	}
	const payload = "a string that gets encoded with base64url"
	err = req.MarshalAsByteArray([]byte(payload), Base64URLFormat)
	if err != nil {
		t.Fatalf("marshal failure: %v", err)
	}
	if ct := req.Header.Get(HeaderContentType); ct != contentTypeAppJSON {
		t.Fatalf("unexpected content type, got %s wanted %s", ct, contentTypeAppJSON)
	}
	if req.Body == nil {
		t.Fatal("unexpected nil request body")
	}
	if req.ContentLength == 0 {
		t.Fatal("unexpected zero content length")
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `"YSBzdHJpbmcgdGhhdCBnZXRzIGVuY29kZWQgd2l0aCBiYXNlNjR1cmw"` {
		t.Fatalf("bad body, got %s", string(b))
	}
}

func TestRequestMarshalAsByteArrayStdFormat(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, "https://contoso.com")
	if err != nil {
		t.Fatal(err)
	}
	const payload = "a string that gets encoded with base64url"
	err = req.MarshalAsByteArray([]byte(payload), Base64StdFormat)
	if err != nil {
		t.Fatalf("marshal failure: %v", err)
	}
	if ct := req.Header.Get(HeaderContentType); ct != contentTypeAppJSON {
		t.Fatalf("unexpected content type, got %s wanted %s", ct, contentTypeAppJSON)
	}
	if req.Body == nil {
		t.Fatal("unexpected nil request body")
	}
	if req.ContentLength == 0 {
		t.Fatal("unexpected zero content length")
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `"YSBzdHJpbmcgdGhhdCBnZXRzIGVuY29kZWQgd2l0aCBiYXNlNjR1cmw="` {
		t.Fatalf("bad body, got %s", string(b))
	}
}

func TestCloneWithoutReadOnlyFieldsNoClone(t *testing.T) {
	nonStruct := "don't clone this"
	v := cloneWithoutReadOnlyFields(&nonStruct)
	if reflect.ValueOf(v).Pointer() != uintptr(unsafe.Pointer(&nonStruct)) {
		t.Fatal()
	}
	type noReadOnly struct {
		ID   int32
		Name *string
	}
	nro := noReadOnly{
		ID:   123,
		Name: &nonStruct,
	}
	v = cloneWithoutReadOnlyFields(&nro)
	if reflect.ValueOf(v).Pointer() != uintptr(unsafe.Pointer(&nro)) {
		t.Fatal("pointers don't match, clone was made")
	}
}

func TestCloneWithoutReadOnlyFieldsClone(t *testing.T) {
	id := int32(123)
	name := "widget"
	type withReadOnly struct {
		ID   *int32  `json:"id" azure:"ro"`
		Name *string `json:"name"`
	}
	nro := withReadOnly{
		ID:   &id,
		Name: &name,
	}
	v := cloneWithoutReadOnlyFields(&nro)
	if reflect.ValueOf(v).Pointer() == uintptr(unsafe.Pointer(&nro)) {
		t.Fatal("pointers match, clone was not made")
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	um := withReadOnly{}
	err = json.Unmarshal(b, &um)
	if err != nil {
		t.Fatal(err)
	}
	if um.ID != nil {
		t.Fatalf("expected nil ID, got %d", *um.ID)
	}
	if um.Name == nil {
		t.Fatal("unexpected nil Name")
	}
}

func TestCloneWithoutReadOnlyFieldsCloneRecursive(t *testing.T) {
	id := int32(123)
	name := "widget"
	something := "something"
	pie := float32(3.14159)
	type inner2 struct {
		Type   *string
		Unique *float32 `json:"omitempty" azure:"ro"`
	}
	type inner1 struct {
		Thing  *string `json:"omitempty" azure:"ro"`
		Color  *string `json:"color"`
		Inner2 *inner2
	}
	type withReadOnly struct {
		ID     *int32  `json:"id" azure:"ro"`
		Name   *string `json:"name"`
		Inner1 *inner1
	}
	nro := withReadOnly{
		ID:   &id,
		Name: &name,
		Inner1: &inner1{
			Thing: &something,
			Color: &something,
			Inner2: &inner2{
				Type:   &something,
				Unique: &pie,
			},
		},
	}
	v := cloneWithoutReadOnlyFields(&nro)
	if reflect.ValueOf(v).Pointer() == uintptr(unsafe.Pointer(&nro)) {
		t.Fatal("pointers match, clone was not made")
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	um := withReadOnly{}
	err = json.Unmarshal(b, &um)
	if err != nil {
		t.Fatal(err)
	}
	if um.ID != nil {
		t.Fatalf("expected nil ID, got %d", *um.ID)
	}
	if um.Name == nil {
		t.Fatal("unexpected nil Name")
	}
	if um.Inner1.Thing != nil {
		t.Fatalf("expected nil Thing, got %s", *um.Inner1.Thing)
	}
	if um.Inner1.Color == nil {
		t.Fatal("unexpected nil Color")
	}
	if um.Inner1.Inner2.Unique != nil {
		t.Fatalf("expected nil Unique, got %f", *um.Inner1.Inner2.Unique)
	}
	if um.Inner1.Inner2.Type == nil {
		t.Fatal("unexpected nil Type")
	}
}

func TestCloneWithoutReadOnlyFieldsCloneNested(t *testing.T) {
	id := int32(123)
	name := "widget"
	something := "something"
	pie := float32(3.14159)
	type inner2 struct {
		Type   *string
		Unique *float32 `json:"omitempty" azure:"ro"`
	}
	type inner1 struct {
		Thing  *string `json:"omitempty"`
		Color  *string `json:"color"`
		Inner2 *inner2
	}
	type withReadOnly struct {
		ID     *int32  `json:"id"`
		Name   *string `json:"name"`
		Inner1 *inner1
	}
	nro := withReadOnly{
		ID:   &id,
		Name: &name,
		Inner1: &inner1{
			Thing: &something,
			Color: &something,
			Inner2: &inner2{
				Type:   &something,
				Unique: &pie,
			},
		},
	}
	v := cloneWithoutReadOnlyFields(&nro)
	if reflect.ValueOf(v).Pointer() == uintptr(unsafe.Pointer(&nro)) {
		t.Fatal("pointers match, clone was not made")
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	um := withReadOnly{}
	err = json.Unmarshal(b, &um)
	if err != nil {
		t.Fatal(err)
	}
	if um.ID == nil {
		t.Fatal("unexpected nil ID")
	}
	if um.Name == nil {
		t.Fatal("unexpected nil Name")
	}
	if um.Inner1.Thing == nil {
		t.Fatal("unexpected nil Thing")
	}
	if um.Inner1.Color == nil {
		t.Fatal("unexpected nil Color")
	}
	if um.Inner1.Inner2.Unique != nil {
		t.Fatalf("expected nil Unique, got %f", *um.Inner1.Inner2.Unique)
	}
	if um.Inner1.Inner2.Type == nil {
		t.Fatal("unexpected nil Type")
	}
}

func TestCloneWithoutReadOnlyFieldsEndToEnd(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, "https://contoso.com")
	if err != nil {
		t.Fatal(err)
	}
	id := int32(123)
	name := "widget"
	type withReadOnly struct {
		ID   *int32  `json:"id" azure:"ro"`
		Name *string `json:"name"`
	}
	nro := withReadOnly{
		ID:   &id,
		Name: &name,
	}
	err = req.MarshalAsJSON(nro)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	um := withReadOnly{}
	err = json.Unmarshal(b, &um)
	if err != nil {
		t.Fatal(err)
	}
	if um.ID != nil {
		t.Fatalf("expected nil ID, got %d", *um.ID)
	}
}

func TestCloneWithoutReadOnlyFieldsCloneEmbedded(t *testing.T) {
	id := int32(123)
	name := "widget"
	something := "something"
	pie := float32(3.14159)
	type Inner2 struct {
		Type  *string
		State *float32 `json:"omitempty" azure:"ro"`
	}
	type Inner1 struct {
		ID   *int32  `json:"id" azure:"ro"`
		ETag *string `json:"omitempty" azure:"ro"`
		Inner2
	}
	type withReadOnly struct {
		Name  *string `json:"name"`
		Color *string `json:"color"`
		Inner1
	}
	nro := withReadOnly{
		Color: &something,
		Name:  &name,
		Inner1: Inner1{
			ID:   &id,
			ETag: &something,
			Inner2: Inner2{
				Type:  &something,
				State: &pie,
			},
		},
	}
	v := cloneWithoutReadOnlyFields(&nro)
	if reflect.ValueOf(v).Pointer() == uintptr(unsafe.Pointer(&nro)) {
		t.Fatal("pointers match, clone was not made")
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	um := withReadOnly{}
	err = json.Unmarshal(b, &um)
	if err != nil {
		t.Fatal(err)
	}
	if um.Name == nil {
		t.Fatal("unexpected nil Name")
	}
	if um.Color == nil {
		t.Fatal("unexpected nil Color")
	}
	if um.ID != nil {
		t.Fatalf("expected nil ID, got %d", *um.ID)
	}
	if um.ETag != nil {
		t.Fatalf("expected nil ETag, got %s", *um.ETag)
	}
	if um.Type == nil {
		t.Fatal("unexpected nil Type")
	}
	if um.State != nil {
		t.Fatalf("expected nil State, got %f", *um.State)
	}
}

func TestCloneWithoutReadOnlyFieldsCloneByVal(t *testing.T) {
	id := int32(123)
	name := "widget"
	type withReadOnly struct {
		ID   *int32  `json:"id" azure:"ro"`
		Name *string `json:"name"`
	}
	nro := withReadOnly{
		ID:   &id,
		Name: &name,
	}
	v := cloneWithoutReadOnlyFields(nro)
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	um := withReadOnly{}
	err = json.Unmarshal(b, &um)
	if err != nil {
		t.Fatal(err)
	}
	if um.ID != nil {
		t.Fatalf("expected nil ID, got %d", *um.ID)
	}
	if um.Name == nil {
		t.Fatal("unexpected nil Name")
	}
}

func TestAzureTagIsReadOnly(t *testing.T) {
	if azureTagIsReadOnly("") {
		t.Fatal("unexpected RO for empty string")
	}
	if azureTagIsReadOnly("rw") {
		t.Fatal("unexpected RO for rw")
	}
	if azureTagIsReadOnly("this,that,the,other") {
		t.Fatal("unexpected RO for this,that,the,other")
	}
	if !azureTagIsReadOnly("ro") {
		t.Fatal("expected RO")
	}
	if !azureTagIsReadOnly("copy,ro,something") {
		t.Fatal("expected RO")
	}
}

func TestRequestSetBodyContentLengthHeader(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPut, "http://test.contoso.com")
	if err != nil {
		t.Fatal(err)
	}
	buff := make([]byte, 768, 768)
	const buffLen = 768
	for i := 0; i < buffLen; i++ {
		buff[i] = 1
	}
	req.SetBody(NopCloser(bytes.NewReader(buff)), "application/octet-stream")
	if req.Header.Get(HeaderContentLength) != strconv.FormatInt(buffLen, 10) {
		t.Fatalf("expected content-length %d, got %s", buffLen, req.Header.Get(HeaderContentLength))
	}
}

func TestNewRequestFail(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodOptions, "://test.contoso.com/")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if req != nil {
		t.Fatal("unexpected request")
	}
	req, err = NewRequest(context.Background(), http.MethodPatch, "/missing/the/host")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if req != nil {
		t.Fatal("unexpected request")
	}
	req, err = NewRequest(context.Background(), http.MethodPatch, "mailto://nobody.contoso.com")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if req != nil {
		t.Fatal("unexpected request")
	}
}

func TestJoinPaths(t *testing.T) {
	if path := JoinPaths(); path != "" {
		t.Fatalf("unexpected path %s", path)
	}
	const expected = "http://test.contoso.com/path/one/path/two/path/three/path/four/"
	if path := JoinPaths("http://test.contoso.com/", "/path/one", "path/two", "/path/three/", "path/four/"); path != expected {
		t.Fatalf("got %s, expected %s", path, expected)
	}
}

func TestRequestValidFail(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodGet, "http://test.contoso.com/")
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("inval d", "header")
	p := NewPipeline(nil)
	resp, err := p.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("unexpected response")
	}
	req.Header = http.Header{}
	// the string "null\0"
	req.Header.Add("invalid", string([]byte{0x6e, 0x75, 0x6c, 0x6c, 0x0}))
	resp, err = p.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("unexpected response")
	}
}
