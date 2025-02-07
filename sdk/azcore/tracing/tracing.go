//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package tracing contains the definitions needed to support distributed tracing.
package tracing

import (
	"context"
)

// ProviderOptions contains the optional values when creating a Provider.
type ProviderOptions struct {
	// NewPropagatorFn contains the underlying implementation for creating Propagator instances.
	// If not set, a no-op propagator will be used.
	NewPropagatorFn func() Propagator
}

// NewProvider creates a new Provider with the specified values.
//   - newTracerFn is the underlying implementation for creating Tracer instances
//   - options contains optional values; pass nil to accept the default value
func NewProvider(newTracerFn func(name, version string) Tracer, options *ProviderOptions) Provider {
	if options == nil {
		options = &ProviderOptions{}
	}
	if options.NewPropagatorFn == nil {
		options.NewPropagatorFn = func() Propagator { return Propagator{} }
	}
	return Provider{
		newTracerFn:     newTracerFn,
		newPropagatorFn: options.NewPropagatorFn,
	}
}

// Provider is the factory that creates Tracer and Propagator instances.
// It defaults to a no-op provider.
type Provider struct {
	newTracerFn     func(name, version string) Tracer
	newPropagatorFn func() Propagator
}

// NewTracer creates a new Tracer for the specified module name and version.
//   - module - the fully qualified name of the module
//   - version - the version of the module
func (p Provider) NewTracer(module, version string) (tracer Tracer) {
	if p.newTracerFn != nil {
		tracer = p.newTracerFn(module, version)
	}
	return
}

// NewPropagator creates a new Propagator.
func (p Provider) NewPropagator() Propagator {
	if p.newPropagatorFn != nil {
		return p.newPropagatorFn()
	}
	return Propagator{}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// TracerOptions contains the optional values when creating a Tracer.
type TracerOptions struct {
	// SpanFromContext contains the implementation for the Tracer.SpanFromContext method.
	SpanFromContext func(context.Context) Span
	// LinkFromContext contains the implementation for the Tracer.LinkFromContext method.
	LinkFromContext func(context.Context, ...Attribute) Link
}

// NewTracer creates a Tracer with the specified values.
//   - newSpanFn is the underlying implementation for creating Span instances
//   - options contains optional values; pass nil to accept the default value
func NewTracer(newSpanFn func(ctx context.Context, spanName string, options *SpanOptions) (context.Context, Span), options *TracerOptions) Tracer {
	if options == nil {
		options = &TracerOptions{}
	}
	return Tracer{
		newSpanFn:         newSpanFn,
		spanFromContextFn: options.SpanFromContext,
		linkFromContextFn: options.LinkFromContext,
	}
}

// Tracer is the factory that creates Span instances.
type Tracer struct {
	attrs             []Attribute
	links             []Link
	newSpanFn         func(ctx context.Context, spanName string, options *SpanOptions) (context.Context, Span)
	spanFromContextFn func(ctx context.Context) Span
	linkFromContextFn func(ctx context.Context, attrs ...Attribute) Link
}

// Start creates a new span and a context.Context that contains it.
//   - ctx is the parent context for this span. If it contains a Span, the newly created span will be a child of that span, else it will be a root span
//   - spanName identifies the span within a trace, it's typically the fully qualified API name
//   - options contains optional values for the span, pass nil to accept any defaults
func (t Tracer) Start(ctx context.Context, spanName string, options *SpanOptions) (context.Context, Span) {
	if t.newSpanFn != nil {
		opts := SpanOptions{}
		if options != nil {
			opts = *options
		}
		opts.Attributes = append(opts.Attributes, t.attrs...)
		return t.newSpanFn(ctx, spanName, &opts)
	}
	return ctx, Span{}
}

// SetAttributes sets attrs to be applied to each Span. If a key from attrs
// already exists for an attribute of the Span it will be overwritten with
// the value contained in attrs.
func (t *Tracer) SetAttributes(attrs ...Attribute) {
	t.attrs = append(t.attrs, attrs...)
}

// Enabled returns true if this Tracer is capable of creating Spans.
func (t Tracer) Enabled() bool {
	return t.newSpanFn != nil
}

// SpanFromContext returns the Span associated with the current context.
// If the provided context has no Span, false is returned.
func (t Tracer) SpanFromContext(ctx context.Context) Span {
	if t.spanFromContextFn != nil {
		return t.spanFromContextFn(ctx)
	}
	return Span{}
}

// LinkFromContext returns a link encapsulating the SpanContext of the current context.
// If the provided context has no Span, an empty Link is returned.
func (t Tracer) LinkFromContext(ctx context.Context, attrs ...Attribute) Link {
	if t.linkFromContextFn != nil {
		return t.linkFromContextFn(ctx, attrs...)
	}
	return Link{}
}

// SpanOptions contains optional settings for creating a span.
type SpanOptions struct {
	// Kind indicates the kind of Span.
	Kind SpanKind

	// Attributes contains key-value pairs of attributes for the span.
	Attributes []Attribute

	// Links contains the links to other spans.
	Links []Link
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// SpanImpl abstracts the underlying implementation for Span,
// allowing it to work with various tracing implementations.
// Any zero-values will have their default, no-op behavior.
type SpanImpl struct {
	// End contains the implementation for the Span.End method.
	End func()

	// SetAttributes contains the implementation for the Span.SetAttributes method.
	SetAttributes func(...Attribute)

	// AddEvent contains the implementation for the Span.AddEvent method.
	AddEvent func(string, ...Attribute)

	// AddLink contains the implementation for the Span.AddLink method.
	AddLink func(Link)

	// SpanContext returns the SpanContext of the Span.
	SpanContext func() SpanContext

	// SetStatus contains the implementation for the Span.SetStatus method.
	SetStatus func(SpanStatus, string)
}

// NewSpan creates a Span with the specified implementation.
func NewSpan(impl SpanImpl) Span {
	return Span{
		impl: impl,
	}
}

// Span is a single unit of a trace.  A trace can contain multiple spans.
// A zero-value Span provides a no-op implementation.
type Span struct {
	impl SpanImpl
}

// End terminates the span and MUST be called before the span leaves scope.
// Any further updates to the span will be ignored after End is called.
func (s Span) End() {
	if s.impl.End != nil {
		s.impl.End()
	}
}

// SetAttributes sets the specified attributes on the Span.
// Any existing attributes with the same keys will have their values overwritten.
func (s Span) SetAttributes(attrs ...Attribute) {
	if s.impl.SetAttributes != nil {
		s.impl.SetAttributes(attrs...)
	}
}

// AddEvent adds a named event with an optional set of attributes to the span.
func (s Span) AddEvent(name string, attrs ...Attribute) {
	if s.impl.AddEvent != nil {
		s.impl.AddEvent(name, attrs...)
	}
}

// AddLink adds a link to the span.
func (s Span) AddLink(link Link) {
	if s.impl.AddLink != nil {
		s.impl.AddLink(link)
	}
}

// SpanContext returns the SpanContext of the Span.
func (s Span) SpanContext() SpanContext {
	if s.impl.SpanContext != nil {
		return s.impl.SpanContext()
	}
	return SpanContext{}
}

// SetStatus sets the status on the span along with a description.
func (s Span) SetStatus(code SpanStatus, desc string) {
	if s.impl.SetStatus != nil {
		s.impl.SetStatus(code, desc)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Attribute is a key-value pair.
type Attribute struct {
	// Key is the name of the attribute.
	Key string

	// Value is the attribute's value.
	// Types that are natively supported include int64, float64, int, bool, string.
	// Any other type will be formatted per rules of fmt.Sprintf("%v").
	Value any
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Link is the relationship between two Spans.
type Link struct {
	// SpanContext of the linked Span.
	SpanContext SpanContext

	// Attributes describe the aspects of the link.
	Attributes []Attribute
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// SpanID is a unique identity of a span in a trace.
type SpanID [8]byte

// TraceID is a unique identity of a trace.
type TraceID [16]byte

// TraceFlags contains flags that can be set on a SpanContext.
type TraceFlags byte

// TraceStateImpl contains the implementation for TraceState.
type TraceStateImpl struct {
	// String contains the implementation for the TraceState.String method.
	String func() string
}

// NewTraceState creates a TraceState with the specified implementation.
func NewTraceState(impl TraceStateImpl) TraceState {
	return TraceState{
		impl: impl,
	}
}

// TraceState provides additional vendor-specific trace identification information across different distributed tracing systems.
type TraceState struct {
	impl TraceStateImpl
}

// String encodes the TraceState into a string compliant with the W3C Trace Context specification.
func (ts TraceState) String() string {
	if ts.impl.String != nil {
		return ts.impl.String()
	}
	return ""
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// SpanContext contains identifying trace information about a Span.
type SpanContext struct {
	spanID     SpanID
	traceID    TraceID
	traceFlags TraceFlags
	traceState TraceState
	remote     bool
}

// SpanID returns the SpanID from the SpanContext.
func (sc SpanContext) SpanID() SpanID {
	return sc.spanID
}

// TraceID returns the TraceID from the SpanContext.
func (sc SpanContext) TraceID() TraceID {
	return sc.traceID
}

// TraceFlags returns the flags from the SpanContext.
func (sc SpanContext) TraceFlags() TraceFlags {
	return sc.traceFlags
}

// TraceState returns the TraceState from the SpanContext.
func (sc SpanContext) TraceState() TraceState {
	return sc.traceState
}

// IsRemote indicates whether the SpanContext represents a remotely-created Span.
func (sc SpanContext) IsRemote() bool {
	return sc.remote
}

// SpanContextConfig contains mutable fields usable for constructing an immutable SpanContext.
type SpanContextConfig struct {
	TraceID    TraceID
	SpanID     SpanID
	TraceFlags TraceFlags
	TraceState TraceState
	Remote     bool
}

// NewSpanContext constructs a SpanContext using values from the provided SpanContextConfig.
func NewSpanContext(config SpanContextConfig) SpanContext {
	return SpanContext{
		traceID:    config.TraceID,
		spanID:     config.SpanID,
		traceFlags: config.TraceFlags,
		traceState: config.TraceState,
		remote:     config.Remote,
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// PropagatorImpl contains the implementation for the Propagator.
type PropagatorImpl struct {
	// Inject contains the implementation for the Propagator.Inject method.
	Inject func(ctx context.Context, carrier Carrier)

	// Extract contains the implementation for the Propagator.Extract method.
	Extract func(ctx context.Context, carrier Carrier) context.Context

	// Fields contains the implementation for the Propagator.Fields method.
	Fields func() []string
}

// NewPropagator creates a Propagator with the specified implementation.
func NewPropagator(impl PropagatorImpl) Propagator {
	return Propagator{
		impl: impl,
	}
}

// Propagator is used to extract and inject context data from and into messages exchanged by applications.
type Propagator struct {
	impl PropagatorImpl
}

// Inject injects the span context into the carrier.
func (p Propagator) Inject(ctx context.Context, carrier Carrier) {
	if p.impl.Inject != nil {
		p.impl.Inject(ctx, carrier)
	}
}

// Extract extracts the span context from the carrier.
func (p Propagator) Extract(ctx context.Context, carrier Carrier) context.Context {
	if p.impl.Extract != nil {
		return p.impl.Extract(ctx, carrier)
	}
	return ctx
}

// Fields returns the fields that the propagator can inject and extract.
func (p Propagator) Fields() []string {
	if p.impl.Fields != nil {
		return p.impl.Fields()
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

type CarrierImpl struct {
	// Get contains the implementation for the Carrier.Get method.
	Get func(key string) string

	// Set contains the implementation for the Carrier.Set method.
	Set func(key string, value string)

	// Keys contains the implementation for the Carrier.Keys method.
	Keys func() []string
}

// NewCarrier creates a Carrier with the specified implementation.
func NewCarrier(impl CarrierImpl) Carrier {
	return Carrier{
		impl: impl,
	}
}

// Carrier is the storage medium used by the Propagator.
type Carrier struct {
	impl CarrierImpl
}

// Get returns the value associated with the passed key.
func (c Carrier) Get(key string) string {
	if c.impl.Get != nil {
		return c.impl.Get(key)
	}
	return ""
}

// Set stores the key-value pair.
func (c Carrier) Set(key string, value string) {
	if c.impl.Set != nil {
		c.impl.Set(key, value)
	}
}

// Keys lists the keys stored in this carrier.
func (c Carrier) Keys() []string {
	if c.impl.Keys != nil {
		return c.impl.Keys()
	}
	return nil
}
