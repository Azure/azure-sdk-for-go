// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"sync"
)

type atomicFloat64 struct {
	f  float64
	mu *sync.RWMutex
}

func newAtomicFloat64(f float64) *atomicFloat64 {
	return &atomicFloat64{
		f:  f,
		mu: &sync.RWMutex{},
	}
}

func (a *atomicFloat64) GetFloat() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.f
}

func (a *atomicFloat64) SetFloat(f float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.f = f
}
