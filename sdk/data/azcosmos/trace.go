// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	aztracing "github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

type trace struct {
	mu        sync.Mutex
	name      string
	startTime time.Time
	endTime   *time.Time
	parent    *trace
	children  []*trace
	data      map[string]any
	dataOrder []string
	summary   *traceSummary
}

type finalizedTrace struct {
	name      string
	startTime time.Time
	endTime   *time.Time
	children  []*finalizedTrace
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

func (t *trace) AddChild(child *trace) {
	t.mu.Lock()
	t.children = append(t.children, child)
	t.mu.Unlock()
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

	t.data[key] = value
	if !exists {
		t.dataOrder = append(t.dataOrder, key)
	}
}

func (t *trace) finalize(freezeTime time.Time) *finalizedTrace {
	if t == nil {
		return nil
	}

	t.mu.Lock()
	name := t.name
	startTime := t.startTime
	endTime := cloneTraceTime(t.endTime, freezeTime)
	children := append([]*trace(nil), t.children...)
	dataOrder := append([]string(nil), t.dataOrder...)
	data := make(map[string]any, len(dataOrder))
	for _, key := range dataOrder {
		data[key] = finalizeTraceDatum(t.data[key])
	}
	t.mu.Unlock()

	finalizedChildren := make([]*finalizedTrace, len(children))
	for index, child := range children {
		finalizedChildren[index] = child.finalize(freezeTime)
	}

	return &finalizedTrace{
		name:      name,
		startTime: startTime,
		endTime:   endTime,
		children:  finalizedChildren,
		data:      data,
		dataOrder: dataOrder,
	}
}

func cloneTraceTime(value *time.Time, freezeTime time.Time) *time.Time {
	frozen := freezeTime.UTC()
	if value != nil {
		frozen = value.UTC()
	}
	return &frozen
}

func finalizeTraceDatum(value any) any {
	switch typed := value.(type) {
	case *clientSideRequestStatisticsTraceDatum:
		return typed.snapshot()
	case clientSideRequestStatisticsSnapshot:
		return typed
	case pointOperationStatisticsTraceDatum:
		return typed
	case *pointOperationStatisticsTraceDatum:
		return *typed
	case json.RawMessage:
		return append(json.RawMessage(nil), typed...)
	case string, float64, float32, int, int32, int64, bool, nil:
		return typed
	case []string:
		return append([]string(nil), typed...)
	case []any:
		next := make([]any, len(typed))
		for index, item := range typed {
			next[index] = finalizeTraceDatum(item)
		}
		return next
	case map[string]any:
		next := make(map[string]any, len(typed))
		for key, item := range typed {
			next[key] = finalizeTraceDatum(item)
		}
		return next
	default:
		return fmt.Sprint(value)
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

func ensureOperationTrace(ctx context.Context, name string) (context.Context, func()) {
	if traceFromContext(ctx) != nil {
		return ctx, func() {}
	}

	root := newRootTrace(name)
	return withTrace(ctx, root), func() {
		root.End()
	}
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
