//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atomic

import (
	"sync/atomic"
	"time"
)

// Int64 is an atomic wrapper around an int64.
type Int64 int64

// NewInt64 creates a new Int64.
func NewInt64(i int64) Int64 {
	return Int64(i)
}

// CAS is an atomic compare-and-swap.
func (i *Int64) CAS(old, new int64) bool {
	return atomic.CompareAndSwapInt64((*int64)(i), old, new)
}

// Load atomically loads the value.
func (i *Int64) Load() int64 {
	return atomic.LoadInt64((*int64)(i))
}

// Store atomically stores the value.
func (i *Int64) Store(v int64) {
	atomic.StoreInt64((*int64)(i), v)
}

// String is an atomic wrapper around a string.
type String struct {
	v atomic.Value
}

// NewString creats a new String.
func NewString(s string) *String {
	ss := String{}
	ss.v.Store(s)
	return &ss
}

// Load atomically loads the string.
func (s *String) Load() string {
	return s.v.Load().(string)
}

// Store atomically stores the string.
func (s *String) Store(v string) {
	s.v.Store(v)
}

// Time is an atomic wrapper around a time.Time.
type Time struct {
	v atomic.Value
}

// NewTime creates a new Time.
func NewTime(t time.Time) *Time {
	tt := Time{}
	tt.v.Store(t)
	return &tt
}

// Load atomically loads the time.Time.
func (t *Time) Load() time.Time {
	return t.v.Load().(time.Time)
}

// Store atomically stores the time.Time.
func (t *Time) Store(v time.Time) {
	t.v.Store(v)
}
