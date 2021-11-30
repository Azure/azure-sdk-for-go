// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package utils

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/devigned/tab"
)

// SimpleTracer is a wrapper around a NoOpTracer so we can add in a
// really simple stderr logger. Useful for seeing some of the internal
// state changes (like retries) that aren't normally customer visible.
// Ex:
// tab.Register(&utils.SimpleTracer{
// 		Include: map[string]bool{
// 			tracing.SpanProcessorClose: true,
// 			tracing.SpanProcessorLoop:  true,
// 			tracing.SpanNegotiateClaim: true,
// 			tracing.SpanRecover:        true,
// 			tracing.SpanRecoverLink:    true,
// 			tracing.SpanRecoverClient:  true,
// 		},
// 	})
type SimpleTracer struct {
	noOpTracer  *tab.NoOpTracer
	spanCounter int64
	include     map[string]bool
	printf      func(format string, v ...interface{})
}

func NewSimpleTracer(include map[string]bool, printf func(format string, v ...interface{})) *SimpleTracer {
	ch := make(chan string, 20000)

	if printf == nil {
		printf = func(format string, v ...interface{}) {
			// try to avoid blocking the main thread
			text := fmt.Sprintf(format, v...)
			ch <- text
		}
	}

	go func() {
		for msg := range ch {
			log.Println(msg)
		}
	}()

	return &SimpleTracer{
		include: include,
		printf:  printf,
	}
}

// StartSpan forwards to NoOpTracer.StartSpan.
func (t *SimpleTracer) StartSpan(ctx context.Context, operationName string, opts ...interface{}) (context.Context, tab.Spanner) {
	id := t.getSpanID(operationName)

	if id == -1 {
		return t.noOpTracer.StartSpan(ctx, operationName, opts...)
	}

	t.printf("[%d] START(%s)", id, operationName)

	ctx = context.WithValue(ctx, spanOpKey("operationName"), operationName)
	return ctx, &stderrSpanner{name: operationName, id: id, printf: t.printf}
}

// StartSpanWithRemoteParent forwards to NoOpTracer.StartSpanWithRemoteParent.
func (t *SimpleTracer) StartSpanWithRemoteParent(ctx context.Context, operationName string, carrier tab.Carrier, opts ...interface{}) (context.Context, tab.Spanner) {
	id := t.getSpanID(operationName)

	if id == -1 {
		return t.noOpTracer.StartSpan(ctx, operationName, opts...)
	}

	t.printf("[%d] START(%s), remote parent", id, operationName)
	ctx = context.WithValue(ctx, spanOpKey("operationName"), operationName)

	return ctx, &stderrSpanner{name: operationName, id: id, printf: t.printf}
}

// FromContext creates a stderrSpanner to allow for our stderrLogger to be created.
func (t *SimpleTracer) FromContext(ctx context.Context) tab.Spanner {
	operationName := ctx.Value(spanOpKey("operationName"))

	val, ok := operationName.(string)

	if !ok {
		val = "<unknown>"
	}

	id := t.getSpanID(val)

	if id == -1 {
		_, span := t.noOpTracer.StartSpan(ctx, val)
		return span
	}

	t.printf("[%d] START(%s), from context", id, val)
	return &stderrSpanner{name: val, id: id, printf: t.printf}
}

// NewContext forwards to NoOpTracer.NewContext
func (t *SimpleTracer) NewContext(parent context.Context, span tab.Spanner) context.Context {
	return t.noOpTracer.NewContext(parent, span)
}

func (t *SimpleTracer) getSpanID(operationName string) int64 {
	if len(t.include) > 0 {
		_, ok := t.include[operationName]

		if !ok {
			return -1
		}
	}

	return atomic.AddInt64(&t.spanCounter, 1)
}

type stderrSpanner struct {
	id      int64
	name    string
	spanner tab.Spanner
	attrs   []tab.Attribute
	printf  func(format string, v ...interface{})
}

func (s *stderrSpanner) AddAttributes(attributes ...tab.Attribute) {
	s.attrs = append(s.attrs, attributes...)
}

func (s *stderrSpanner) End() {
	s.printf("[%d] END(%s)\n%s", s.id, s.name, sprintfAttributes(s.attrs))
}

func (s *stderrSpanner) Logger() tab.Logger {
	return &stderrLogger{id: s.id, printf: s.printf}
}

func (s *stderrSpanner) Inject(carrier tab.Carrier) error {
	return nil
}

func (s *stderrSpanner) InternalSpan() interface{} {
	return s.spanner.InternalSpan()
}

type stderrLogger struct {
	id     int64
	printf func(format string, v ...interface{})
}

func sprintfAttributes(attributes []tab.Attribute) string {
	s := ""

	for _, attr := range attributes {
		s += fmt.Sprintf("  %s=%v\n", attr.Key, attr.Value)
	}

	return s
}

func (l *stderrLogger) Info(msg string, attributes ...tab.Attribute) {
	l.printf("[%d] INFO: %s\n%s", l.id, msg, sprintfAttributes(attributes))
}

func (l *stderrLogger) Error(err error, attributes ...tab.Attribute) {
	l.printf("[%d] ERROR: error: %v\n%s", l.id, err, sprintfAttributes(attributes))
}

func (l *stderrLogger) Fatal(msg string, attributes ...tab.Attribute) {
	l.printf("[%d] FATAL: %s\n%s", l.id, msg, sprintfAttributes(attributes))
}

func (l *stderrLogger) Debug(msg string, attributes ...tab.Attribute) {
	l.printf("[%d] DEBUG: %s\n%s", l.id, msg, sprintfAttributes(attributes))
}

type spanOpKey string
