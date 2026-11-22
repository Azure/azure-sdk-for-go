//go:build js && wasm

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"fmt"
	"math"
)

// mmb is a plain heap-allocated buffer used on platforms without mmap support (e.g. WASM).
type mmb []byte

// newMMB allocates a plain byte slice. No mmap is used.
func newMMB(size int64) (mmb, error) {
	if size < 0 || size > math.MaxInt {
		return nil, fmt.Errorf("size %d out of range", size)
	}

	return make([]byte, int(size)), nil
}

// delete is a no-op; GC handles reclamation.
func (m *mmb) delete() {
	*m = nil
}
