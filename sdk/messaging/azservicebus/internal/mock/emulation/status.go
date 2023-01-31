// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"sync/atomic"
)

type Status struct {
	done   chan struct{}
	errVal atomic.Value
}

func NewStatus(parent *Status) *Status {
	s := &Status{
		done: make(chan struct{}, 1),
	}

	if parent != nil {
		select {
		case <-parent.Done():
			s.CloseWithError(parent.Err())
		default:
			go func() {
				<-parent.Done()
				s.CloseWithError(parent.Err())
			}()
		}
	}

	return s
}

func (s *Status) Done() <-chan struct{} {
	if s == nil {
		return nil
	}

	return s.done
}

func (s *Status) CloseWithError(err error) {
	if err == nil {
		panic(err)
	}

	if s.errVal.CompareAndSwap(nil, &err) {
		close(s.done)
	}
}

func (s *Status) Err() error {
	return *s.errVal.Load().(*error)
}
