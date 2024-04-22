//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/stretchr/testify/require"
)

func TestNewResponseErrorNoBodyNoErrorCode(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	if err != nil {
		t.Fatal(err)
	}
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Body:       http.NoBody,
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	if re.ErrorCode != "" {
		t.Fatal("expected empty error code")
	}
	if c := re.StatusCode; c != http.StatusInternalServerError {
		t.Fatalf("unexpected status code %d", c)
	}
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE UNAVAILABLE
--------------------------------------------------------------------------------
Response contained no body
--------------------------------------------------------------------------------
`
	if got := re.Error(); got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestNewResponseErrorNoBody(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	if err != nil {
		t.Fatal(err)
	}
	respHeader := http.Header{}
	const errorCode = "ErrorTooManyCheats"
	respHeader.Set(shared.HeaderXMSErrorCode, errorCode)
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Body:       http.NoBody,
		Header:     respHeader,
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	if ec := re.ErrorCode; ec != errorCode {
		t.Fatalf("unexpected error code %s", ec)
	}
	if c := re.StatusCode; c != http.StatusInternalServerError {
		t.Fatalf("unexpected status code %d", c)
	}
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE: ErrorTooManyCheats
--------------------------------------------------------------------------------
Response contained no body
--------------------------------------------------------------------------------
`
	if got := re.Error(); got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestNewResponseErrorNoErrorCode(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	if err != nil {
		t.Fatal(err)
	}
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(`{ "code": "ErrorItsBroken", "message": "it's not working" }`)),
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	if c := re.StatusCode; c != http.StatusInternalServerError {
		t.Fatalf("unexpected status code %d", c)
	}
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE: ErrorItsBroken
--------------------------------------------------------------------------------
{
  "code": "ErrorItsBroken",
  "message": "it's not working"
}
--------------------------------------------------------------------------------
`
	if got := re.Error(); got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestNewResponseErrorPreferErrorCodeHeader(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	if err != nil {
		t.Fatal(err)
	}
	respHeader := http.Header{}
	respHeader.Set(shared.HeaderXMSErrorCode, "ErrorTooManyCheats")
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(`{ "code": "ErrorItsBroken", "message": "it's not working" }`)),
		Header:     respHeader,
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	if c := re.StatusCode; c != http.StatusInternalServerError {
		t.Fatalf("unexpected status code %d", c)
	}
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE: ErrorTooManyCheats
--------------------------------------------------------------------------------
{
  "code": "ErrorItsBroken",
  "message": "it's not working"
}
--------------------------------------------------------------------------------
`
	if got := re.Error(); got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestNewResponseErrorNoErrorCodeWrappedError(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	if err != nil {
		t.Fatal(err)
	}
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(`{ "error": { "code": "ErrorItsBroken", "message": "it's not working" } }`)),
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	if c := re.StatusCode; c != http.StatusInternalServerError {
		t.Fatalf("unexpected status code %d", c)
	}
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE: ErrorItsBroken
--------------------------------------------------------------------------------
{
  "error": {
    "code": "ErrorItsBroken",
    "message": "it's not working"
  }
}
--------------------------------------------------------------------------------
`
	if got := re.Error(); got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestNewResponseErrorNoErrorCodeInvalidBody(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	if err != nil {
		t.Fatal(err)
	}
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("JSON error string")),
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	if c := re.StatusCode; c != http.StatusInternalServerError {
		t.Fatalf("unexpected status code %d", c)
	}
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE UNAVAILABLE
--------------------------------------------------------------------------------
JSON error string
--------------------------------------------------------------------------------
`
	if got := re.Error(); got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

type readFailer struct{}

func (r *readFailer) Close() error {
	return nil
}

func (r *readFailer) Read(p []byte) (int, error) {
	return 0, errors.New("mock read failure")
}

func TestNewResponseErrorNoErrorCodeCantReadBody(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	if err != nil {
		t.Fatal(err)
	}
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Body:       &readFailer{},
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	_, ok := err.(*ResponseError)
	if ok {
		t.Fatalf("unexpected error type %T", err)
	}
	const want = `mock read failure`
	if got := err.Error(); got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestNewResponseErrorNoErrorCodeXML(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	if err != nil {
		t.Fatal(err)
	}
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Error><Code>ContainerAlreadyExists</Code><Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</Message></Error>`)),
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	if c := re.StatusCode; c != http.StatusInternalServerError {
		t.Fatalf("unexpected status code %d", c)
	}
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE: ContainerAlreadyExists
--------------------------------------------------------------------------------
<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Error><Code>ContainerAlreadyExists</Code><Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</Message></Error>
--------------------------------------------------------------------------------
`
	if got := re.Error(); got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestNewResponseErrorErrorCodeHeaderXML(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	if err != nil {
		t.Fatal(err)
	}
	respHeader := http.Header{}
	respHeader.Set(shared.HeaderXMSErrorCode, "ContainerAlreadyExists")
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Header:     respHeader,
		Body:       io.NopCloser(strings.NewReader(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Error><Code>ContainerAlreadyExists</Code><Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</Message></Error>`)),
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	if c := re.StatusCode; c != http.StatusInternalServerError {
		t.Fatalf("unexpected status code %d", c)
	}
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE: ContainerAlreadyExists
--------------------------------------------------------------------------------
<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Error><Code>ContainerAlreadyExists</Code><Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</Message></Error>
--------------------------------------------------------------------------------
`
	if got := re.Error(); got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestNewResponseErrorErrorCodeHeaderXMLWithNamespace(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	if err != nil {
		t.Fatal(err)
	}
	respHeader := http.Header{}
	respHeader.Set(shared.HeaderXMSErrorCode, "ContainerAlreadyExists")
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Header:     respHeader,
		Body:       io.NopCloser(strings.NewReader(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><m:Error xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata"><m:Code>ContainerAlreadyExists</m:Code><m:Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</m:Message></m:Error>`)),
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	if c := re.StatusCode; c != http.StatusInternalServerError {
		t.Fatalf("unexpected status code %d", c)
	}
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE: ContainerAlreadyExists
--------------------------------------------------------------------------------
<?xml version="1.0" encoding="UTF-8" standalone="yes"?><m:Error xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata"><m:Code>ContainerAlreadyExists</m:Code><m:Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</m:Message></m:Error>
--------------------------------------------------------------------------------
`
	if got := re.Error(); got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestNewResponseErrorAllMissingXML(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	if err != nil {
		t.Fatal(err)
	}
	respHeader := http.Header{}
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Header:     respHeader,
		Body:       io.NopCloser(strings.NewReader(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Error><Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</Message></Error>`)),
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("unexpected error type %T", err)
	}
	if c := re.StatusCode; c != http.StatusInternalServerError {
		t.Fatalf("unexpected status code %d", c)
	}
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE UNAVAILABLE
--------------------------------------------------------------------------------
<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Error><Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</Message></Error>
--------------------------------------------------------------------------------
`
	if got := re.Error(); got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestExtractErrorCodeFromJSON(t *testing.T) {
	errorBody := []byte(`{"odata.error": {
		"code": "ResourceNotFound",
		"message": {
		  "lang": "en-us",
		  "value": "The specified resource does not exist.\nRequestID:b2437f3b-ca2d-47a1-95a7-92f73a768a1c\n"
		}
	  }
	}`)
	code := extractErrorCodeJSON(errorBody)
	if code != "ResourceNotFound" {
		t.Fatalf("expected %s got %s", "ResourceNotFound", code)
	}

	errorBody = []byte(`{"error": {
		"code": "ResourceNotFound",
		"message": {
		  "lang": "en-us",
		  "value": "The specified resource does not exist.\nRequestID:b2437f3b-ca2d-47a1-95a7-92f73a768a1c\n"
		}
	  }
	}`)
	code = extractErrorCodeJSON(errorBody)
	if code != "ResourceNotFound" {
		t.Fatalf("expected %s got %s", "ResourceNotFound", code)
	}
}

func TestNilRawResponse(t *testing.T) {
	const expected = "Missing RawResponse\n--------------------------------------------------------------------------------\nERROR CODE UNAVAILABLE\n--------------------------------------------------------------------------------\n"
	require.EqualValues(t, expected, (&ResponseError{}).Error())
}

func TestNilRequestInRawResponse(t *testing.T) {
	const expected = "Request information not available\n--------------------------------------------------------------------------------\nRESPONSE 400: status\nERROR CODE UNAVAILABLE\n--------------------------------------------------------------------------------\nResponse contained no body\n--------------------------------------------------------------------------------\n"
	respErr := &ResponseError{
		RawResponse: &http.Response{
			Body:       http.NoBody,
			Status:     "status",
			StatusCode: http.StatusBadRequest,
		},
	}
	require.EqualValues(t, expected, respErr.Error())
}

func TestNilResponseBody(t *testing.T) {
	const expected = "Request information not available\n--------------------------------------------------------------------------------\nRESPONSE 0: \nERROR CODE UNAVAILABLE\n--------------------------------------------------------------------------------\nResponse contained no body\n--------------------------------------------------------------------------------\n"
	require.EqualValues(t, expected, (&ResponseError{RawResponse: &http.Response{}}).Error())
}

func TestLogResponseErrorCodeInBody(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	require.NoError(t, err)
	rawlog := map[log.Event][]string{}
	log.SetListener(func(cls log.Event, s string) {
		rawlog[cls] = append(rawlog[cls], s)
	})
	defer log.SetListener(nil)
	_ = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(`{ "error": { "code": "ErrorItsBroken", "message": "it's not working" } }`)),
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE: ErrorItsBroken
--------------------------------------------------------------------------------
{
  "error": {
    "code": "ErrorItsBroken",
    "message": "it's not working"
  }
}
--------------------------------------------------------------------------------
`
	msg, ok := rawlog[azlog.EventResponseError]
	require.True(t, ok)
	require.Len(t, msg, 1)
	require.EqualValues(t, want, msg[0])
}

func TestLogResponseErrorCodeInHeader(t *testing.T) {
	fakeURL, err := url.Parse("https://fakeurl.com/the/path?qp=removed")
	require.NoError(t, err)
	rawlog := map[log.Event][]string{}
	log.SetListener(func(cls log.Event, s string) {
		rawlog[cls] = append(rawlog[cls], s)
	})
	defer log.SetListener(nil)
	respHeader := http.Header{}
	respHeader.Set(shared.HeaderXMSErrorCode, "ErrorTooManyCheats")
	_ = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(`{ "code": "ErrorItsBroken", "message": "it's not working" }`)),
		Header:     respHeader,
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	const want = `GET https://fakeurl.com/the/path
--------------------------------------------------------------------------------
RESPONSE 500: the system is down
ERROR CODE: ErrorTooManyCheats
--------------------------------------------------------------------------------
{
  "code": "ErrorItsBroken",
  "message": "it's not working"
}
--------------------------------------------------------------------------------
`
	msg, ok := rawlog[azlog.EventResponseError]
	require.True(t, ok)
	require.Len(t, msg, 1)
	require.EqualValues(t, want, msg[0])
}
