// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestCreateClientAssertionJWT(t *testing.T) {
	_, err := url.Parse(defaultAuthorityHost)
	if err != nil {
		t.Fatalf("Failed to parse default authority host: %v", err)
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

// func (f RoundTripFunc) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
// 	return f(ctx, req), nil
// }

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestSendAuthRequest(t *testing.T) {
	// server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	// 	// Test request parameters
	// 	equals(t, req.URL.String(), "https://login.microsoftonline.com/")
	// 	// Send response to be tested
	// 	rw.Write([]byte(`OK`))
	// }))
	// // Close the server when test finishes
	// defer server.Close()

	// // Use Client & URL from our local test server
	// urlStr, err := url.Parse("https://login.microsoftonline.com/")
	// if err != nil {
	// 	t.Fatalf("Error: %v", err)
	// }
	// pOpts := &IdentityClientOptions{AuthorityHost: urlStr, PipelineOptions: azcore.PipelineOptions{
	// 	HTTPClient: azcore.TransportFunc(func(ctx context.Context, req *http.Request) (*http.Response, error) {
	// 		return server.Client().Do(req.WithContext(ctx))
	// 	})}}
	// // body := NewClientSecretCredential("expected_tenant", "expected_client", "secret", pOpts)
	// body := NewClientSecretCredential("72f988bf-86f1-41af-91ab-2d7cd011db47", "31334978-f7d6-49a6-bf4f-8cebe115f455", "QBiIj7j54L3EvtM[@AZ?/CG/k3iJZuS8", pOpts)
	// _, err = body.GetToken(context.Background(), &azcore.TokenRequestOptions{Scopes: []string{"www.storage.azure.com/.default"}})
	// ok(t, err)
	// equals(t, []byte("OK"), body)
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
