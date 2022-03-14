package perf

import (
	"sync"
)

type atomicFloat64 struct {
	f  float64
	mu sync.RWMutex
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
