// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && azcosmos_driver

package azcosmos

/*
#cgo CFLAGS: -I${SRCDIR}/internal/native
#cgo LDFLAGS: -lazurecosmosdriver

#include <stdlib.h>
#include "azurecosmosdriver.h"
*/
import "C"

import (
	"runtime"
	"runtime/cgo"
	"sync"
)

// The driver is asynchronous: submitting an operation returns immediately and its result is posted
// to a completion queue, which the host drains. This is the reactor that does the draining.
//
// One goroutine per client blocks in cosmos_completion_queue_wait and hands each completion to the
// operation that submitted it, matched by the pointer-sized cookie the driver round-trips
// verbatim. Operations therefore block on a channel rather than on the queue, so many can be in
// flight against a single queue and none of them holds an OS thread.

// completionBatch is how many completions one wait call may drain. The queue hands back everything
// already available once it has woken, so a larger batch means fewer wait calls under load.
const completionBatch = 64

// completionWaitMillis bounds how long a wait call blocks before returning empty. The reactor only
// needs to wake to notice that it has been asked to stop, so this is a shutdown-latency knob
// rather than a throughput one.
const completionWaitMillis = 250

// reactor drains one completion queue and delivers each completion to the operation that submitted
// it.
type reactor struct {
	queue *C.cosmos_completion_queue_t

	// stop is closed to ask the loop to finish. done is closed by the loop when it has.
	stop chan struct{}
	done chan struct{}

	stopOnce sync.Once
}

// pendingOperation receives the single completion for one submitted operation.
//
// The channel is buffered so the reactor never blocks delivering a result, even if the submitting
// goroutine has already abandoned the wait because its context was cancelled.
type pendingOperation struct {
	result chan completionResult
}

// newReactor creates a completion queue on the runtime and starts draining it.
func newReactor(rt *C.cosmos_runtime_t) (*reactor, error) {
	options := C.cosmos_completion_queue_options_t{
		capacity_hint: C.uint32_t(completionBatch),
		max_capacity:  0, // unbounded; the driver's own backpressure applies
		// Rich error details are what let a failed operation report a message and headers rather
		// than a bare status.
		include_error_details: C.bool(true),
	}

	queue := C.cosmos_completion_queue_create(rt, &options)
	if queue == nil {
		return nil, &Error{
			Code:    CodeClientError,
			Message: "azcosmos: creating the completion queue",
		}
	}

	r := &reactor{
		queue: queue,
		stop:  make(chan struct{}),
		done:  make(chan struct{}),
	}
	go r.run()
	return r, nil
}

// run drains the queue until stopped. It owns the queue handle for its lifetime.
func (r *reactor) run() {
	defer close(r.done)

	// The wait call blocks in C, which parks the calling thread. Pinning the goroutine to its own
	// OS thread keeps the Go scheduler from having to grow the thread pool around it, and keeps
	// the blocking call off a thread that other goroutines are sharing.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	batch := make([]C.cosmos_completion_t, completionBatch)
	for {
		select {
		case <-r.stop:
			// Drain what is already queued before leaving, so operations that completed while
			// shutdown was starting still get their result rather than hanging.
			r.drainOnce(batch, 0)
			return
		default:
		}
		r.drainOnce(batch, completionWaitMillis)
	}
}

// drainOnce waits up to timeoutMillis for completions and delivers everything it receives.
func (r *reactor) drainOnce(batch []C.cosmos_completion_t, timeoutMillis uint32) {
	n := C.cosmos_completion_queue_wait(
		r.queue,
		&batch[0],
		C.uintptr_t(len(batch)),
		C.uint32_t(timeoutMillis),
	)
	if n == 0 {
		return
	}
	// Every completion the wait wrote has to be freed, whether or not it was delivered, and the
	// whole run can be freed at once. Translation copies out of it first.
	defer C.cosmos_completion_queue_free_completions(&batch[0], n)

	for i := 0; i < int(n); i++ {
		r.deliver(&batch[i])
	}
}

// deliver hands one completion to the operation that submitted it.
//
// The cookie is a cgo.Handle for the pendingOperation. A completion whose cookie is zero or no
// longer valid is dropped: that means the submitting goroutine has already gone away, which is not
// an error, and dereferencing it would be a use-after-free.
func (r *reactor) deliver(completion *C.cosmos_completion_t) {
	cookie := uintptr(completion.user_data)
	if cookie == 0 {
		return
	}

	handle := cgo.Handle(cookie)
	pending, ok := handle.Value().(*pendingOperation)
	if !ok {
		return
	}

	// Translation copies every field it needs, because the completion's memory is reclaimed as
	// soon as this drain returns.
	select {
	case pending.result <- translateCompletion(completion):
	default:
		// The buffer is sized for exactly one result, so a full channel means the operation was
		// already answered. Nothing to do.
	}
}

// close stops the reactor and releases the queue. It is idempotent.
//
// Shutting the queue down first unblocks the wait call rather than waiting out its timeout, and
// causes the driver to post cancelled completions for anything still in flight, so operations
// blocked on a result are answered rather than left hanging.
func (r *reactor) close() {
	if r == nil {
		return
	}
	r.stopOnce.Do(func() {
		close(r.stop)
		C.cosmos_completion_queue_shutdown(r.queue)
		<-r.done
		// Freed only after the loop has returned, so nothing can be waiting on it.
		C.cosmos_completion_queue_free(r.queue)
		r.queue = nil
	})
}
