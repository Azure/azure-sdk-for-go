// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAtomicFloat(t *testing.T) {
	a := newAtomicFloat64(0)

	wg := &sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100*i; j++ {
				a.SetFloat(float64(j))
			}
		}(i)
		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()

	require.Equal(t, float64(99), a.GetFloat())
}
