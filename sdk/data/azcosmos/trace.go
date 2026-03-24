// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"sync"
	"time"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	aztracing "github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

type trace struct {
	mu            sync.Mutex
	name          string
	startTime     time.Time
	endTime       *time.Time
	parent        *trace
	children      []*trace
	data          map[string]any
	dataOrder     []string
	isBeingWalked bool
	summary       *traceSummary
}

type traceSnapshot struct {
	name      string
	startTime time.Time
	endTime   *time.Time
	children  []*trace
	data      map[string]any
	dataOrder []string
}

type traceSummary struct {
	mu                 sync.Mutex
	failedRequestCount int
}

func newRootTrace(name string) *trace {
	return &trace{
		name:      name,
		startTime: time.Now().UTC(),
		children:  []*trace{},
		summary:   &traceSummary{},
	}
}

func (t *trace) root() *trace {
	current := t
	for current != nil && current.parent != nil {
		current = current.parent
	}

	return current
}

func (t *trace) StartChild(name string) *trace {
	child := &trace{
		name:      name,
		startTime: time.Now().UTC(),
		parent:    t,
		children:  []*trace{},
		summary:   t.summary,
	}

	t.AddChild(child)
	return child
}

func (t *trace) End() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.endTime != nil {
		return
	}

	endTime := time.Now().UTC()
	t.endTime = &endTime
}

func (t *trace) duration() time.Duration {
	t.mu.Lock()
	startTime := t.startTime
	endTime := t.endTime
	t.mu.Unlock()

	if endTime != nil {
		return endTime.Sub(startTime)
	}

	return time.Since(startTime)
}

func (t *trace) AddChild(child *trace) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.isBeingWalked {
		t.children = append(t.children, child)
		return
	}

	child.setWalkingStateRecursively()

	nextChildren := make([]*trace, 0, len(t.children)+1)
	nextChildren = append(nextChildren, t.children...)
	nextChildren = append(nextChildren, child)
	t.children = nextChildren
}

func (t *trace) AddDatum(key string, value any) {
	t.addOrUpdateDatum(key, value, false)
}

func (t *trace) AddOrUpdateDatum(key string, value any) {
	t.addOrUpdateDatum(key, value, true)
}

func (t *trace) addOrUpdateDatum(key string, value any, overwrite bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.data == nil {
		t.data = map[string]any{}
	}

	_, exists := t.data[key]
	if exists && !overwrite {
		return
	}

	if !t.isBeingWalked {
		t.data[key] = value
		if !exists {
			t.dataOrder = append(t.dataOrder, key)
		}
		return
	}

	nextData := make(map[string]any, len(t.data)+1)
	for k, v := range t.data {
		nextData[k] = v
	}
	nextData[key] = value
	t.data = nextData

	if !exists {
		nextOrder := make([]string, 0, len(t.dataOrder)+1)
		nextOrder = append(nextOrder, t.dataOrder...)
		nextOrder = append(nextOrder, key)
		t.dataOrder = nextOrder
	}
}

func (t *trace) setWalkingStateRecursively() {
	t.mu.Lock()
	if t.isBeingWalked {
		t.mu.Unlock()
		return
	}

	children := append([]*trace(nil), t.children...)
	t.mu.Unlock()

	for _, child := range children {
		child.setWalkingStateRecursively()
	}

	t.mu.Lock()
	t.isBeingWalked = true
	t.mu.Unlock()
}

func (t *trace) snapshot() traceSnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()

	return traceSnapshot{
		name:      t.name,
		startTime: t.startTime,
		endTime:   t.endTime,
		children:  t.children,
		data:      t.data,
		dataOrder: t.dataOrder,
	}
}

func (s *traceSummary) incrementFailedCount() {
	s.mu.Lock()
	s.failedRequestCount++
	s.mu.Unlock()
}

func (s *traceSummary) failedCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.failedRequestCount
}

type traceContextKey struct{}

func traceFromContext(ctx context.Context) *trace {
	current, _ := ctx.Value(traceContextKey{}).(*trace)
	return current
}

func withTrace(ctx context.Context, t *trace) context.Context {
	return context.WithValue(ctx, traceContextKey{}, t)
}

func diagnosticsFromContext(ctx context.Context) Diagnostics {
	return newDiagnostics(traceFromContext(ctx))
}

func ensureOperationTrace(ctx context.Context, name string) context.Context {
	if traceFromContext(ctx) != nil {
		return ctx
	}

	return withTrace(ctx, newRootTrace(name))
}

func startSpan(ctx context.Context, name string, tracer aztracing.Tracer, options *azruntime.StartSpanOptions) (context.Context, func(error)) {
	ctx, endSpan := azruntime.StartSpan(ctx, name, tracer, options)

	currentTrace := traceFromContext(ctx)
	if currentTrace == nil {
		root := newRootTrace(name)
		ctx = withTrace(ctx, root)

		return ctx, func(err error) {
			root.End()
			endSpan(err)
		}
	}

	child := currentTrace.StartChild(name)
	ctx = withTrace(ctx, child)

	return ctx, func(err error) {
		child.End()
		endSpan(err)
	}
}
