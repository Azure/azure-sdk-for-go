//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/stretchr/testify/require"
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
	err = MarshalAsXML(req, testXML{SomeInt: 1, SomeString: "s"})
	if err != nil {
		t.Fatalf("marshal failure: %v", err)
	}
	if ct := req.Raw().Header.Get(shared.HeaderContentType); ct != shared.ContentTypeAppXML {
		t.Fatalf("unexpected content type, got %s wanted %s", ct, shared.ContentTypeAppXML)
	}
	if req.Raw().Body == nil {
		t.Fatal("unexpected nil request body")
	}
	if req.Raw().ContentLength == 0 {
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
	if err == nil {
		t.Fatal("unexpected nil error")
	}
}

func TestRequestMarshalJSON(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, "https://contoso.com")
	if err != nil {
		t.Fatal(err)
	}
	err = MarshalAsJSON(req, testJSON{SomeInt: 1, SomeString: "s"})
	if err != nil {
		t.Fatalf("marshal failure: %v", err)
	}
	if ct := req.Raw().Header.Get(shared.HeaderContentType); ct != shared.ContentTypeAppJSON {
		t.Fatalf("unexpected content type, got %s wanted %s", ct, shared.ContentTypeAppJSON)
	}
	if req.Raw().Body == nil {
		t.Fatal("unexpected nil request body")
	}
	if req.Raw().ContentLength == 0 {
		t.Fatal("unexpected zero content length")
	}
}

func TestRequestMarshalAsByteArrayURLFormat(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, "https://contoso.com")
	if err != nil {
		t.Fatal(err)
	}
	const payload = "a string that gets encoded with base64url"
	err = MarshalAsByteArray(req, []byte(payload), Base64URLFormat)
	if err != nil {
		t.Fatalf("marshal failure: %v", err)
	}
	if ct := req.Raw().Header.Get(shared.HeaderContentType); ct != shared.ContentTypeAppJSON {
		t.Fatalf("unexpected content type, got %s wanted %s", ct, shared.ContentTypeAppJSON)
	}
	if req.Raw().Body == nil {
		t.Fatal("unexpected nil request body")
	}
	if req.Raw().ContentLength == 0 {
		t.Fatal("unexpected zero content length")
	}
	b, err := io.ReadAll(req.Raw().Body)
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
	err = MarshalAsByteArray(req, []byte(payload), Base64StdFormat)
	if err != nil {
		t.Fatalf("marshal failure: %v", err)
	}
	if ct := req.Raw().Header.Get(shared.HeaderContentType); ct != shared.ContentTypeAppJSON {
		t.Fatalf("unexpected content type, got %s wanted %s", ct, shared.ContentTypeAppJSON)
	}
	if req.Raw().Body == nil {
		t.Fatal("unexpected nil request body")
	}
	if req.Raw().ContentLength == 0 {
		t.Fatal("unexpected zero content length")
	}
	b, err := io.ReadAll(req.Raw().Body)
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

	t.Setenv("AZURE_SDK_GO_OMIT_READONLY", "true")

	err = MarshalAsJSON(req, nro)
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(req.Raw().Body)
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

func TestCloneWithReadOnlyFieldsEndToEnd(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, "https://contoso.com")
	require.NoError(t, err)

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

	err = MarshalAsJSON(req, nro)
	require.NoError(t, err)

	b, err := io.ReadAll(req.Raw().Body)
	require.NoError(t, err)

	um := withReadOnly{}
	err = json.Unmarshal(b, &um)
	require.NoError(t, err)

	require.NotNil(t, um.ID)
	require.Equal(t, int32(123), *um.ID)
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

func TestCloneWithoutReadOnlyFieldsTime(t *testing.T) {
	id := int32(123)
	expires := time.Date(2021, 10, 13, 8, 48, 31, 0, time.UTC)
	type withTime struct {
		ID      *int32     `json:"id" azure:"ro"`
		Expires *time.Time `json:"expires"`
	}
	nro := withTime{
		ID:      &id,
		Expires: &expires,
	}
	v := cloneWithoutReadOnlyFields(nro)
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	um := withTime{}
	err = json.Unmarshal(b, &um)
	if err != nil {
		t.Fatal(err)
	}
	if um.ID != nil {
		t.Fatalf("expected nil ID, got %d", *um.ID)
	}
	if um.Expires == nil {
		t.Fatal("unexpected nil Expires")
	} else if *um.Expires != expires {
		t.Fatalf("unexpected Expires %v", *um.Expires)
	}
}

func TestCloneWithoutReadOnlyFieldsNilField(t *testing.T) {
	type zeroValues struct {
		A *string `json:"a"`
		B *string `json:"b" azure:"ro"`
	}
	expected := zeroValues{}
	clone := cloneWithoutReadOnlyFields(expected)
	if reflect.ValueOf(clone).Pointer() == uintptr(unsafe.Pointer(&expected)) {
		t.Fatal("pointers match, clone was not made")
	}
	b, err := json.Marshal(clone)
	if err != nil {
		t.Fatal(err)
	}
	um := zeroValues{}
	err = json.Unmarshal(b, &um)
	if err != nil {
		t.Fatal(err)
	}
	if um.A != expected.A || um.B != expected.B {
		t.Fatal("unexpected values in unmarshalled struct")
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
	buff := make([]byte, 768)
	const buffLen = 768
	for i := 0; i < buffLen; i++ {
		buff[i] = 1
	}
	err = req.SetBody(exported.NopCloser(bytes.NewReader(buff)), "application/octet-stream")
	if err != nil {
		t.Fatal(err)
	}
	if req.Raw().Header.Get(shared.HeaderContentLength) != strconv.FormatInt(buffLen, 10) {
		t.Fatalf("expected content-length %d, got %s", buffLen, req.Raw().Header.Get(shared.HeaderContentLength))
	}
}

func TestJoinPaths(t *testing.T) {
	type joinTest struct {
		root     string
		paths    []string
		expected string
	}

	tests := []joinTest{
		{
			root:     "",
			paths:    nil,
			expected: "",
		},
		{
			root:     "/",
			paths:    nil,
			expected: "/",
		},
		{
			root:     "http://test.contoso.com/",
			paths:    []string{"/path/one", "path/two", "/path/three/", "path/four/"},
			expected: "http://test.contoso.com/path/one/path/two/path/three/path/four/",
		},
		{
			root:     "http://test.contoso.com",
			paths:    []string{"path/one", "path/two", "/path/three/", "path/four/"},
			expected: "http://test.contoso.com/path/one/path/two/path/three/path/four/",
		},
		{
			root:     "http://test.contoso.com/?qp1=abc&qp2=def",
			paths:    []string{"/path/one", "path/two"},
			expected: "http://test.contoso.com/path/one/path/two?qp1=abc&qp2=def",
		},
		{
			root:     "http://test.contoso.com?qp1=abc&qp2=def",
			paths:    []string{"path/one", "path/two/"},
			expected: "http://test.contoso.com/path/one/path/two/?qp1=abc&qp2=def",
		},
		{
			root:     "http://test.contoso.com/?qp1=abc&qp2=def",
			paths:    []string{"path/one", "path/two/"},
			expected: "http://test.contoso.com/path/one/path/two/?qp1=abc&qp2=def",
		},
		{
			root:     "http://test.contoso.com/?qp1=abc&qp2=def",
			paths:    []string{"/path/one", "path/two/"},
			expected: "http://test.contoso.com/path/one/path/two/?qp1=abc&qp2=def",
		},
	}

	for _, tt := range tests {
		if path := JoinPaths(tt.root, tt.paths...); path != tt.expected {
			t.Fatalf("got %s, expected %s", path, tt.expected)
		}
	}
}

func TestRequestValidFail(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodGet, "http://test.contoso.com/")
	if err != nil {
		t.Fatal(err)
	}
	req.Raw().Header.Add("inval d", "header")
	p := exported.NewPipeline(nil)
	resp, err := p.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("unexpected response")
	}
	req.Raw().Header = http.Header{}
	// the string "null\0"
	req.Raw().Header.Add("invalid", string([]byte{0x6e, 0x75, 0x6c, 0x6c, 0x0}))
	resp, err = p.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("unexpected response")
	}
}

func TestSetMultipartFormData(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, "https://contoso.com")
	if err != nil {
		t.Fatal(err)
	}
	err = SetMultipartFormData(req, map[string]interface{}{
		"string": "value",
		"int":    1,
		"data":   exported.NopCloser(strings.NewReader("some data")),
		"datum": []io.ReadSeekCloser{
			exported.NopCloser(strings.NewReader("first part")),
			exported.NopCloser(strings.NewReader("second part")),
			exported.NopCloser(strings.NewReader("third part")),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	mt, params, err := mime.ParseMediaType(req.Raw().Header.Get(shared.HeaderContentType))
	if err != nil {
		t.Fatal(err)
	}
	if mt != "multipart/form-data" {
		t.Fatalf("unexpected media type %s", mt)
	}
	reader := multipart.NewReader(req.Raw().Body, params["boundary"])
	var datum []io.ReadSeekCloser
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		switch fn := part.FormName(); fn {
		case "string":
			strPart := make([]byte, 16)
			_, err = part.Read(strPart)
			if err != io.EOF {
				t.Fatal(err)
			}
			if tr := string(strPart[:5]); tr != "value" {
				t.Fatalf("unexpected value %s", tr)
			}
		case "int":
			intPart := make([]byte, 16)
			_, err = part.Read(intPart)
			if err != io.EOF {
				t.Fatal(err)
			}
			if tr := string(intPart[:1]); tr != "1" {
				t.Fatalf("unexpected value %s", tr)
			}
		case "data":
			dataPart := make([]byte, 16)
			_, err = part.Read(dataPart)
			if err != io.EOF {
				t.Fatal(err)
			}
			if tr := string(dataPart[:9]); tr != "some data" {
				t.Fatalf("unexpected value %s", tr)
			}
		case "datum":
			content, err := io.ReadAll(part)
			require.NoError(t, err)
			datum = append(datum, exported.NopCloser(bytes.NewReader(content)))
		default:
			t.Fatalf("unexpected part %s", fn)
		}
	}
	require.Len(t, datum, 3)
	first, err := io.ReadAll(datum[0])
	require.NoError(t, err)
	second, err := io.ReadAll(datum[1])
	require.NoError(t, err)
	third, err := io.ReadAll(datum[2])
	require.NoError(t, err)
	require.Equal(t, "first part", string(first))
	require.Equal(t, "second part", string(second))
	require.Equal(t, "third part", string(third))
}
