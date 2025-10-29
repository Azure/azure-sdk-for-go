//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// HasStatusCode returns true if the Response's status code is one of the specified values.
// Exported as runtime.HasStatusCode().
func HasStatusCode(resp *http.Response, statusCodes ...int) bool {
	if resp == nil {
		return false
	}
	for _, sc := range statusCodes {
		if resp.StatusCode == sc {
			return true
		}
	}
	return false
}

// PayloadOptions contains the optional values for the Payload func.
// NOT exported but used by azcore.
type PayloadOptions struct {
	// BytesModifier receives the downloaded byte slice and returns an updated byte slice.
	// Use this to modify the downloaded bytes in a payload (e.g. removing a BOM).
	BytesModifier func([]byte) []byte
}

// Payload reads and returns the response body or an error.
// On a successful read, the response body is cached.
// Subsequent reads will access the cached value.
// Exported as runtime.Payload() WITHOUT the opts parameter.
func Payload(resp *http.Response, opts *PayloadOptions) ([]byte, error) {
	if resp.Body == nil {
		// this shouldn't happen in real-world scenarios as a
		// response with no body should set it to http.NoBody
		return nil, nil
	}
	modifyBytes := func(b []byte) []byte { return b }
	if opts != nil && opts.BytesModifier != nil {
		modifyBytes = opts.BytesModifier
	}

	// r.Body won't be a nopClosingBytesReader if downloading was skipped
	if buf, ok := resp.Body.(*nopClosingBytesReader); ok {
		bytesBody := modifyBytes(buf.Bytes())
		buf.Set(bytesBody)
		return bytesBody, nil
	}

	bytesBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	bytesBody = modifyBytes(bytesBody)
	resp.Body = &nopClosingBytesReader{s: bytesBody}
	return bytesBody, nil
}

// PayloadDownloaded returns true if the response body has already been downloaded.
// This implies that the Payload() func above has been previously called.
// NOT exported but used by azcore.
func PayloadDownloaded(resp *http.Response) bool {
	_, ok := resp.Body.(*nopClosingBytesReader)
	return ok
}

// nopClosingBytesReader is an io.ReadSeekCloser around a byte slice.
// It also provides direct access to the byte slice to avoid rereading.
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

// Seek implements the io.Seeker interface.
func (r *nopClosingBytesReader) Seek(offset int64, whence int) (int64, error) {
	var i int64
	switch whence {
	case io.SeekStart:
		i = offset
	case io.SeekCurrent:
		i = r.i + offset
	case io.SeekEnd:
		i = int64(len(r.s)) + offset
	default:
		return 0, errors.New("nopClosingBytesReader: invalid whence")
	}
	if i < 0 {
		return 0, errors.New("nopClosingBytesReader: negative position")
	}
	r.i = i
	return i, nil
}

// Unmarshaler is an interface for custom JSON unmarshaling implementations.
// This is NOT exported directly but used internally by azcore.
// The public API is exposed via runtime.SetUnmarshaler().
type Unmarshaler interface {
	// Unmarshal parses the JSON-encoded data and stores the result
	// in the value pointed to by v.
	Unmarshal(data []byte, v any) error
}

// DefaultUnmarshaler is the global JSON unmarshaler used by all Azure SDK operations.
// By default, it uses the standard library's encoding/json.
// This can be customized via runtime.SetUnmarshaler().
var DefaultUnmarshaler Unmarshaler = &stdlibUnmarshaler{}

// stdlibUnmarshaler is the default unmarshaler that uses encoding/json from the standard library.
type stdlibUnmarshaler struct{}

func (s *stdlibUnmarshaler) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// SetUnmarshaler replaces the default JSON unmarshaler with a custom implementation.
// This is NOT exported directly - use runtime.SetUnmarshaler() instead.
// WARNING: This is NOT thread-safe. Only call this in init() functions or at application startup.
// For concurrent/per-request usage, use WithUnmarshaler() with context instead.
func SetUnmarshaler(u Unmarshaler) {
	if u == nil {
		panic("exported: SetUnmarshaler called with nil")
	}
	DefaultUnmarshaler = u
}

// unmarshalerKey is the context key type for storing custom Unmarshaler instances.
type unmarshalerKey struct{}

// WithUnmarshaler returns a new context with the specified Unmarshaler attached.
// This is thread-safe and allows per-request unmarshaler customization.
// NOT exported directly - use runtime.WithUnmarshaler() instead.
func WithUnmarshaler(ctx context.Context, u Unmarshaler) context.Context {
	if u == nil {
		panic("exported: WithUnmarshaler called with nil")
	}
	return context.WithValue(ctx, unmarshalerKey{}, u)
}

// GetUnmarshaler returns the Unmarshaler from the context, or the global DefaultUnmarshaler if not set.
// This provides the lookup mechanism for context-based unmarshaler selection.
// NOT exported - used internally by runtime and pollers.
func GetUnmarshaler(ctx context.Context) Unmarshaler {
	if ctx != nil {
		if u, ok := ctx.Value(unmarshalerKey{}).(Unmarshaler); ok {
			return u
		}
	}
	return DefaultUnmarshaler
}
