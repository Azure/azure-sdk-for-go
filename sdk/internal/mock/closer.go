//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package mock

import "io"

type trackedCloser struct {
	io.Reader
	closed bool
}

// Close records that Close was called.
func (t *trackedCloser) Close() error {
	t.closed = true
	return nil
}

// Closed returns true if Close was called.
func (t *trackedCloser) Closed() bool {
	return t.closed
}

// NewTrackedCloser is similar to io.NopCloser but tracks that Close
// was called.  Call TrackedClose to check if Close has been called.
func NewTrackedCloser(r io.Reader) (io.ReadCloser, TrackedClose) {
	tc := &trackedCloser{Reader: r}
	return tc, tc.Closed
}

// TrackedClose, when invoked, returns true if Close was called.
type TrackedClose func() bool
