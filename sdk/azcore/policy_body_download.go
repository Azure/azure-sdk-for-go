// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
)

// newBodyDownloadPolicy creates a policy object that downloads the response's body to a []byte.
func newBodyDownloadPolicy() Policy {
	return PolicyFunc(func(ctx context.Context, req *Request) (*Response, error) {
		resp, err := req.Next(ctx)
		if err != nil {
			return resp, err
		}
		var opValues bodyDownloadPolicyOpValues
		if req.OperationValue(&opValues); !opValues.skip && resp.Body != nil {
			// Either bodyDownloadPolicyOpValues was not specified (so skip is false)
			// or it was specified and skip is false: don't skip downloading the body
			b, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				err = fmt.Errorf("body download policy: %w", err)
			}
			resp.Body = &nopClosingBytesReader{s: b}
		}
		return resp, err
	})
}

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
