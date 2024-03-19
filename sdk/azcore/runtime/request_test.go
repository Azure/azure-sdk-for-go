//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"bytes"
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"testing"

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
	err = SetMultipartFormData(req, map[string]any{
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

func TestEncodeQueryParams(t *testing.T) {
	const testURL = "https://contoso.com/"
	nextLink, err := EncodeQueryParams(testURL + "query?$skip=5&$filter='foo eq bar'")
	require.NoError(t, err)
	require.EqualValues(t, testURL+"query?%24filter=%27foo+eq+bar%27&%24skip=5", nextLink)
	nextLink, err = EncodeQueryParams(testURL + "query?%24filter=%27foo+eq+bar%27&%24skip=5")
	require.NoError(t, err)
	require.EqualValues(t, testURL+"query?%24filter=%27foo+eq+bar%27&%24skip=5", nextLink)
	nextLink, err = EncodeQueryParams(testURL + "query?foo=bar&one=two")
	require.NoError(t, err)
	require.EqualValues(t, testURL+"query?foo=bar&one=two", nextLink)
	nextLink, err = EncodeQueryParams(testURL)
	require.NoError(t, err)
	require.EqualValues(t, testURL, nextLink)
	nextLink, err = EncodeQueryParams(testURL + "query?compound=thing1;thing2;thing3")
	require.NoError(t, err)
	require.EqualValues(t, testURL+"query?compound=thing1%3Bthing2%3Bthing3", nextLink)
}
