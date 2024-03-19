//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
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
	require.NoError(t, err)
	err = SetMultipartFormData(req, map[string]any{
		"json":   []byte(`{"id":123}`),
		"string": "value",
		"int":    1,
		"data":   exported.NopCloser(strings.NewReader("some data")),
		"datum": []io.ReadSeekCloser{
			exported.NopCloser(strings.NewReader("first part")),
			exported.NopCloser(strings.NewReader("second part")),
			exported.NopCloser(strings.NewReader("third part")),
		},
	})
	require.NoError(t, err)
	mt, params, err := mime.ParseMediaType(req.Raw().Header.Get(shared.HeaderContentType))
	require.NoError(t, err)
	require.EqualValues(t, "multipart/form-data", mt)
	reader := multipart.NewReader(req.Raw().Body, params["boundary"])
	var datum []io.ReadSeekCloser
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		switch fn := part.FormName(); fn {
		case "json":
			data, err := io.ReadAll(part)
			require.NoError(t, err)
			require.EqualValues(t, '{', data[0])
			type thing struct {
				ID int `json:"id"`
			}
			thing1 := thing{}
			require.NoError(t, json.Unmarshal(data, &thing1))
			require.EqualValues(t, 123, thing1.ID)
			require.EqualValues(t, "application/json", part.Header.Get(shared.HeaderContentType))
		case "string":
			strPart, err := io.ReadAll(part)
			require.NoError(t, err)
			require.EqualValues(t, "value", strPart)
			require.EqualValues(t, "text/plain", part.Header.Get(shared.HeaderContentType))
		case "int":
			intPart, err := io.ReadAll(part)
			require.NoError(t, err)
			require.EqualValues(t, "1", intPart)
			require.EqualValues(t, "text/plain", part.Header.Get(shared.HeaderContentType))
		case "data":
			dataPart, err := io.ReadAll(part)
			require.NoError(t, err)
			require.EqualValues(t, "some data", dataPart)
			require.EqualValues(t, "application/octet-stream", part.Header.Get(shared.HeaderContentType))
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

func TestSetMultipartContent(t *testing.T) {
	req, err := NewRequest(context.Background(), http.MethodPost, "https://contoso.com")
	require.NoError(t, err)
	err = SetMultipartFormData(req, map[string]any{
		"default": streaming.MultipartContent{
			Body: exported.NopCloser(strings.NewReader("default body")),
		},
		"withContentType": streaming.MultipartContent{
			Body:        exported.NopCloser(strings.NewReader("body with content type")),
			ContentType: "text/plain",
		},
		"withFilename": streaming.MultipartContent{
			Body:     exported.NopCloser(strings.NewReader("body with filename")),
			Filename: "content.txt",
		},
		"allSet": streaming.MultipartContent{
			Body:        exported.NopCloser(strings.NewReader("body with everything set")),
			ContentType: "text/plain",
			Filename:    "content.txt",
		},
		"multiple": []streaming.MultipartContent{
			{
				Body:     exported.NopCloser(bytes.NewReader([]byte{1, 2, 3, 4, 5})),
				Filename: "data.bin",
			},
			{
				Body:        exported.NopCloser(strings.NewReader("some text")),
				ContentType: "text/plain",
			},
		},
	})
	require.NoError(t, err)
	mt, params, err := mime.ParseMediaType(req.Raw().Header.Get(shared.HeaderContentType))
	require.NoError(t, err)
	require.EqualValues(t, "multipart/form-data", mt)
	reader := multipart.NewReader(req.Raw().Body, params["boundary"])
	countMultiple := 0
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		switch fn := part.FormName(); fn {
		case "default":
			require.EqualValues(t, "default", part.FileName())
			require.EqualValues(t, "application/octet-stream", part.Header.Get(shared.HeaderContentType))
			body, err := io.ReadAll(part)
			require.NoError(t, err)
			require.EqualValues(t, "default body", body)
		case "withContentType":
			require.EqualValues(t, "withContentType", part.FileName())
			require.EqualValues(t, "text/plain", part.Header.Get(shared.HeaderContentType))
			body, err := io.ReadAll(part)
			require.NoError(t, err)
			require.EqualValues(t, "body with content type", body)
		case "withFilename":
			require.EqualValues(t, "content.txt", part.FileName())
			require.EqualValues(t, "application/octet-stream", part.Header.Get(shared.HeaderContentType))
			body, err := io.ReadAll(part)
			require.NoError(t, err)
			require.EqualValues(t, "body with filename", body)
		case "allSet":
			require.EqualValues(t, "content.txt", part.FileName())
			require.EqualValues(t, "text/plain", part.Header.Get(shared.HeaderContentType))
			body, err := io.ReadAll(part)
			require.NoError(t, err)
			require.EqualValues(t, "body with everything set", body)
		case "multiple":
			body, err := io.ReadAll(part)
			require.NoError(t, err)
			if fn := part.FileName(); fn == "data.bin" {
				require.EqualValues(t, "application/octet-stream", part.Header.Get(shared.HeaderContentType))
				require.EqualValues(t, []byte{1, 2, 3, 4, 5}, body)
			} else if fn == "multiple" {
				require.EqualValues(t, "text/plain", part.Header.Get(shared.HeaderContentType))
				require.EqualValues(t, "some text", body)
			} else {
				t.Fatalf("unexpected file %s", fn)
			}
			countMultiple++
		default:
			t.Fatalf("unexpected part %s", fn)
		}
	}
	require.EqualValues(t, 2, countMultiple)
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
