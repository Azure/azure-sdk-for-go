// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include "azurecosmosdriver.h"
*/
import "C"

import "sync"

type cursorDelivery struct {
	completion *C.cosmos_cursor_completion_t
	err        error
}

// One client-owned drainer serves all retained cursors, separately from legacy completions.
type cursorReactor struct {
	queue     *C.cosmos_completion_queue_t
	mu        sync.Mutex
	next      uintptr
	pending   map[uintptr]chan cursorDelivery
	err       error
	stop      chan struct{}
	done      chan struct{}
	closeOnce sync.Once
}

func (d *nativeDriver) queryReactor() (*cursorReactor, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return nil, errClientClosed()
	}
	if d.cursorReactor != nil {
		return d.cursorReactor, nil
	}
	queue := C.cosmos_cursor_queue_create(d.runtime, 0)
	if queue == nil {
		return nil, invalidCursorCompletion("creating shared cursor queue")
	}
	r := &cursorReactor{
		queue: queue, pending: make(map[uintptr]chan cursorDelivery),
		stop: make(chan struct{}), done: make(chan struct{}),
	}
	d.cursorReactor = r
	go r.run()
	return r, nil
}

func (r *cursorReactor) register() (uintptr, chan cursorDelivery, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return 0, nil, cloneError(r.err)
	}
	// Refuse wraparound rather than reusing an outstanding correlation cookie.
	if r.next == ^uintptr(0) {
		return 0, nil, invalidCursorCompletion("cursor correlation space exhausted")
	}
	r.next++
	result := make(chan cursorDelivery, 1)
	r.pending[r.next] = result
	return r.next, result, nil
}

func (r *cursorReactor) forget(cookie uintptr) {
	r.mu.Lock()
	delete(r.pending, cookie)
	r.mu.Unlock()
}

func (r *cursorReactor) deliver(delivery cursorDelivery, cookie uintptr) {
	r.mu.Lock()
	waiter := r.pending[cookie]
	delete(r.pending, cookie)
	r.mu.Unlock()
	if waiter != nil {
		waiter <- delivery
	} else if delivery.completion != nil {
		C.cosmos_cursor_completion_free(delivery.completion)
	}
}

func (r *cursorReactor) fail(err error) {
	r.mu.Lock()
	r.err = err
	waiters := r.pending
	r.pending = make(map[uintptr]chan cursorDelivery)
	r.mu.Unlock()
	for _, waiter := range waiters {
		waiter <- cursorDelivery{err: cloneError(err)}
	}
}

func (r *cursorReactor) run() {
	defer close(r.done)
	batch := make([]*C.cosmos_cursor_completion_t, completionBatch)
	for {
		select {
		case <-r.stop:
			return
		default:
		}
		var count C.uintptr_t
		status := C.cosmos_cursor_queue_wait(r.queue, &batch[0], C.uintptr_t(len(batch)), completionWaitMillis, &count)
		if status != 0 {
			r.fail(statusError(status, nil, "waiting for shared cursor completions"))
			return
		}
		for i := 0; i < int(count); i++ {
			completion := batch[i]
			batch[i] = nil
			r.deliver(cursorDelivery{completion: completion}, uintptr(completion.common.user_data))
		}
	}
}

// The client first waits for active callers and canceled-operation cleanup.
func (r *cursorReactor) close() {
	if r == nil {
		return
	}
	r.closeOnce.Do(func() {
		r.fail(errClientClosed())
		close(r.stop)
		C.cosmos_completion_queue_shutdown(r.queue)
		<-r.done
		C.cosmos_completion_queue_free(r.queue)
		r.queue = nil
	})
}
