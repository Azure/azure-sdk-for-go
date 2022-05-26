// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package stress

import (
	"context"
	"log"

	"github.com/devigned/tab"
)

// StderrTracer is a wrapper around a NoOpTracer so we can add in a
// really simple stderr logger. Useful for seeing some of the internal
// state changes (like retries) that aren't normally customer visible.
type StderrTracer struct {
	NoOpTracer *tab.NoOpTracer
}

// StartSpan forwards to NoOpTracer.StartSpan.
func (ft *StderrTracer) StartSpan(ctx context.Context, operationName string, opts ...interface{}) (context.Context, tab.Spanner) {
	return ft.NoOpTracer.StartSpan(ctx, operationName, opts...)
}

// StartSpanWithRemoteParent forwards to NoOpTracer.StartSpanWithRemoteParent.
func (ft *StderrTracer) StartSpanWithRemoteParent(ctx context.Context, operationName string, carrier tab.Carrier, opts ...interface{}) (context.Context, tab.Spanner) {
	return ft.NoOpTracer.StartSpanWithRemoteParent(ctx, operationName, carrier, opts...)
}

// FromContext creates a stderrSpanner to allow for our stderrLogger to be created.
func (ft *StderrTracer) FromContext(ctx context.Context) tab.Spanner {
	return &stderrSpanner{
		spanner: ft.NoOpTracer.FromContext(ctx),
	}
}

// NewContext forwards to NoOpTracer.NewContext
func (ft *StderrTracer) NewContext(parent context.Context, span tab.Spanner) context.Context {
	return ft.NoOpTracer.NewContext(parent, span)
}

type stderrSpanner struct {
	spanner tab.Spanner
}

func (s *stderrSpanner) AddAttributes(attributes ...tab.Attribute) {}

func (s *stderrSpanner) End() {}

func (s *stderrSpanner) Logger() tab.Logger {
	return &stderrLogger{}
}

func (s *stderrSpanner) Inject(carrier tab.Carrier) error {
	return nil
}

func (s *stderrSpanner) InternalSpan() interface{} {
	return s.spanner.InternalSpan()
}

type stderrLogger struct{}

func (l *stderrLogger) Info(msg string, attributes ...tab.Attribute) {
	log.Printf("INFO: %s", msg)
}
func (l *stderrLogger) Error(err error, attributes ...tab.Attribute) {
	log.Printf("ERROR: %s", err.Error())
}

func (l *stderrLogger) Fatal(msg string, attributes ...tab.Attribute) {
	log.Printf("FATAL: %s", msg)
}

func (l *stderrLogger) Debug(msg string, attributes ...tab.Attribute) {
	log.Printf("DEBUG: %s", msg)
}
