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

// StderrTracer is a wrapper around a NoOpTracer so we can add in a
// really simple stderr logger. Useful for seeing some of the internal
// state changes (like retries) that aren't normally customer visible.
// Ex:
// tab.Register(&utils.StderrTracer{
// 		Include: map[string]bool{
// 			tracing.SpanProcessorClose: true,
// 			tracing.SpanProcessorLoop:  true,
// 			tracing.SpanNegotiateClaim: true,
// 			tracing.SpanRecover:        true,
// 			tracing.SpanRecoverLink:    true,
// 			tracing.SpanRecoverClient:  true,
// 		},
// 	})
type StderrTracer struct {
	NoOpTracer  *tab.NoOpTracer
	spanCounter int64
	Include     map[string]bool
}

type spanOpKey string

// StartSpan forwards to NoOpTracer.StartSpan.
func (t *StderrTracer) StartSpan(ctx context.Context, operationName string, opts ...interface{}) (context.Context, tab.Spanner) {
	id := t.getSpanID(operationName)

	if id == -1 {
		return t.NoOpTracer.StartSpan(ctx, operationName, opts...)
	}

	log.Printf("[%d] START(%s)", id, operationName)

	ctx = context.WithValue(ctx, spanOpKey("operationName"), operationName)
	return ctx, &stderrSpanner{name: operationName, id: id}
}

// StartSpanWithRemoteParent forwards to NoOpTracer.StartSpanWithRemoteParent.
func (t *StderrTracer) StartSpanWithRemoteParent(ctx context.Context, operationName string, carrier tab.Carrier, opts ...interface{}) (context.Context, tab.Spanner) {
	id := t.getSpanID(operationName)

	if id == -1 {
		return t.NoOpTracer.StartSpan(ctx, operationName, opts...)
	}

	log.Printf("[%d] START(%s), remote parent", id, operationName)
	ctx = context.WithValue(ctx, spanOpKey("operationName"), operationName)

	return ctx, &stderrSpanner{name: operationName, id: id}
}

// FromContext creates a stderrSpanner to allow for our stderrLogger to be created.
func (t *StderrTracer) FromContext(ctx context.Context) tab.Spanner {
	operationName := ctx.Value(spanOpKey("operationName"))

	val, ok := operationName.(string)

	if !ok {
		val = "<unknown>"
	}

	id := t.getSpanID(val)

	if id == -1 {
		_, span := t.NoOpTracer.StartSpan(ctx, val)
		return span
	}

	log.Printf("[%d] START(%s), from context", id, val)
	return &stderrSpanner{name: val, id: id}
}

// NewContext forwards to NoOpTracer.NewContext
func (t *StderrTracer) NewContext(parent context.Context, span tab.Spanner) context.Context {
	return t.NoOpTracer.NewContext(parent, span)
}

func (t *StderrTracer) getSpanID(operationName string) int64 {
	if len(t.Include) > 0 {
		_, ok := t.Include[operationName]

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
}

func (s *stderrSpanner) AddAttributes(attributes ...tab.Attribute) {
	s.attrs = append(s.attrs, attributes...)
}

func (s *stderrSpanner) End() {
	log.Printf("[%d] END(%s)\n%s", s.id, s.name, sprintfAttributes(s.attrs))
}

func (s *stderrSpanner) Logger() tab.Logger {
	return &stderrLogger{id: s.id}
}

func (s *stderrSpanner) Inject(carrier tab.Carrier) error {
	return nil
}

func (s *stderrSpanner) InternalSpan() interface{} {
	return s.spanner.InternalSpan()
}

type stderrLogger struct {
	id int64
}

func sprintfAttributes(attributes []tab.Attribute) string {
	s := ""

	for _, attr := range attributes {
		s += fmt.Sprintf("  %s=%v\n", attr.Key, attr.Value)
	}

	return s
}

func (l *stderrLogger) Info(msg string, attributes ...tab.Attribute) {
	log.Printf("[%d] INFO: %s\n%s", l.id, msg, sprintfAttributes(attributes))
}

func (l *stderrLogger) Error(err error, attributes ...tab.Attribute) {
	log.Printf("[%d] ERROR: error: %v\n%s", l.id, err, sprintfAttributes(attributes))
}

func (l *stderrLogger) Fatal(msg string, attributes ...tab.Attribute) {
	log.Printf("[%d] FATAL: %s\n%s", l.id, msg, sprintfAttributes(attributes))
}

func (l *stderrLogger) Debug(msg string, attributes ...tab.Attribute) {
	log.Printf("[%d] DEBUG: %s\n%s", l.id, msg, sprintfAttributes(attributes))
}
