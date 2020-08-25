// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// newBodyDownloadPolicy creates a policy object that downloads the response's body to a []byte.
func newBodyDownloadPolicy() Policy {
	return PolicyFunc(func(ctx context.Context, req *Request) (*Response, error) {
		resp, err := req.Next(ctx)
		if err != nil {
			return resp, err
		}
		var opValues bodyDownloadPolicyOpValues
		// don't skip downloading error response bodies
		if req.OperationValue(&opValues); opValues.skip && resp.StatusCode < 400 {
			return resp, err
		}
		// Either bodyDownloadPolicyOpValues was not specified (so skip is false)
		// or it was specified and skip is false: don't skip downloading the body
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return resp, newBodyDownloadError(err, req)
		}
		resp.Body = &nopClosingBytesReader{s: b}
		return resp, err
	})
}

type bodyDownloadError struct {
	err     error
	noRetry bool
}

func newBodyDownloadError(err error, req *Request) error {
	// on failure, only retry the request for idempotent operations.
	// we currently identify them as DELETE, GET, and PUT requests.
	noRetry := true
	if m := strings.ToUpper(req.Method); m == http.MethodDelete || m == http.MethodGet || m == http.MethodPut {
		noRetry = false
	}
	return &bodyDownloadError{
		err:     err,
		noRetry: noRetry,
	}
}

func (b *bodyDownloadError) Error() string {
	return fmt.Sprintf("body download policy: %s", b.err.Error())
}

func (b *bodyDownloadError) IsNotRetriable() bool {
	return b.noRetry
}

func (b *bodyDownloadError) Unwrap() error {
	return b.err
}

var _ Retrier = (*bodyDownloadError)(nil)

// bodyDownloadPolicyOpValues is the struct containing the per-operation values
type bodyDownloadPolicyOpValues struct {
	skip bool
}

// nopClosingBytesReader is an io.ReadCloser around a byte slice.
// It also provides direct access to the byte slice.
type nopClosingBytesReader struct {
	s []byte
	i int64
}

// Bytes returns the underlying byte slice.
func (r *nopClosingBytesReader) Bytes() []byte {
	return r.s
}

// Close implements the io.Closer interface.
func (*nopClosingBytesReader) Close() error {
	return nil
}

// Read implements the io.Reader interface.
func (r *nopClosingBytesReader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

// Set replaces the existing byte slice with the specified byte slice and resets the reader.
func (r *nopClosingBytesReader) Set(b []byte) {
	r.s = b
	r.i = 0
}
