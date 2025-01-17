// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/stretchr/testify/require"
)

const (
	SpanStatusUnset = tracing.SpanStatusUnset
	SpanStatusError = tracing.SpanStatusError
	SpanStatusOK    = tracing.SpanStatusOK
)

type SpanStatus = tracing.SpanStatus

// NewSpanValidator creates a Provider that verifies a span was created that matches the specified SpanMatcher.
// The returned Provider can be used to create a client with a tracing provider that will validate spans in unit tests.
func NewSpanValidator(t *testing.T, matcher SpanMatcher) Provider {
	return tracing.NewProvider(func(name, version string) Tracer {
		tt := matchingTracer{
			matcher: matcher,
		}

		t.Cleanup(func() {
			require.NotNil(t, tt.match, "didn't find a span with name %s", tt.matcher.Name)
			require.True(t, tt.match.ended, "span wasn't ended")
			if tt.matcher.Kind != 0 {
				require.EqualValues(t, tt.matcher.Kind, tt.match.kind, "span kind values don't match")
			}
			require.EqualValues(t, matcher.Status, tt.match.status, "span status values don't match")
			require.ElementsMatch(t, matcher.Attributes, tt.match.attrs, "span attributes don't match")
		})

		return tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
			kind := tracing.SpanKindInternal
			attrs := []Attribute{}
			if options != nil {
				kind = options.Kind
				attrs = append(attrs, options.Attributes...)
			}
			return tt.Start(ctx, spanName, kind, attrs)
		}, nil)
	}, nil)
}

// SpanMatcher contains the values to match when a span is created.
type SpanMatcher struct {
	Name       string
	Kind       SpanKind
	Status     SpanStatus
	Attributes []Attribute
}

type matchingTracer struct {
	matcher SpanMatcher
	match   *span
}

func (mt *matchingTracer) Start(ctx context.Context, spanName string, kind SpanKind, attrs []Attribute) (context.Context, tracing.Span) {
	if spanName != mt.matcher.Name {
		return ctx, tracing.Span{}
	}
	// span name matches our matcher, track it
	mt.match = &span{
		name:  spanName,
		kind:  kind,
		attrs: attrs,
	}
	return ctx, tracing.NewSpan(tracing.SpanImpl{
		End:           mt.match.End,
		SetStatus:     mt.match.SetStatus,
		SetAttributes: mt.match.SetAttributes,
	})
}

type span struct {
	name   string
	kind   SpanKind
	status SpanStatus
	desc   string
	attrs  []Attribute
	ended  bool
}

func (s *span) End() {
	s.ended = true
}

func (s *span) SetAttributes(attrs ...Attribute) {
	s.attrs = append(s.attrs, attrs...)
}

func (s *span) SetStatus(code SpanStatus, desc string) {
	s.status = code
	s.desc = desc
	s.ended = true
}
