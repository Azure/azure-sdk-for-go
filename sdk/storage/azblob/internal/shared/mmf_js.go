//go:build js

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

// Mmb is a plain heap-allocated buffer used on platforms without mmap support (e.g. WASM).
type Mmb []byte

// NewMMB allocates a plain byte slice. No mmap is used.
func NewMMB(size int64) (Mmb, error) {
	return make([]byte, size), nil
}

// Delete is a no-op; GC handles reclamation.
func (m *Mmb) Delete() {
	*m = nil
}
