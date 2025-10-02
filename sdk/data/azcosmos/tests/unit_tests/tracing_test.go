// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/stretchr/testify/require"
)

type spanContextKey struct{}

// newSpanValidator creates a tracing.Provider that verifies a span was created that matches the specified SpanMatcher.
func newSpanValidator(t *testing.T, matcher *spanMatcher) tracing.Provider {
	return tracing.NewProvider(func(name, version string) tracing.Tracer {
		tt := matchingTracer{
			matcher: matcher,
		}

		t.Cleanup(func() {
			for _, expectedSpan := range matcher.ExpectedSpans {
				found := false
				for _, match := range matcher.MatchedSpans {
					if match.name == expectedSpan {
						found = true
						require.True(t, match.ended, "span %s wasn't ended", match.name)
						break
					}
				}
				require.True(t, found, "span %s wasn't found", expectedSpan)
			}
		})

		return tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
			kind := tracing.SpanKindInternal
			if options != nil {
				kind = options.Kind
			}
			return tt.Start(ctx, spanName, kind)
		}, &tracing.TracerOptions{
			SpanFromContext: func(ctx context.Context) tracing.Span {
				if span, ok := ctx.Value(spanContextKey{}).(tracing.Span); ok {
					return span
				}
				return tracing.Span{}
			},
		})
	}, nil)
}

// SpanMatcher contains the values to match when a span is created.
type spanMatcher struct {
	ExpectedSpans []string
	MatchedSpans  []*matchingSpan
}

type matchingTracer struct {
	matcher *spanMatcher
}

func (mt *matchingTracer) Start(ctx context.Context, spanName string, kind tracing.SpanKind) (context.Context, tracing.Span) {

	if slices.IndexFunc(mt.matcher.ExpectedSpans, func(i string) bool { return i == spanName }) < 0 && !strings.Contains(spanName, "NextPage") {
		return ctx, tracing.Span{}
	}
	// span name matches our matcher, track it
	newSpan := &matchingSpan{
		name: spanName,
	}
	mt.matcher.MatchedSpans = append(mt.matcher.MatchedSpans, newSpan)
	tracingSpan := tracing.NewSpan(tracing.SpanImpl{
		End:           newSpan.End,
		SetStatus:     newSpan.SetStatus,
		SetAttributes: newSpan.SetAttributes,
	})
	ctx = context.WithValue(ctx, spanContextKey{}, tracingSpan)
	return ctx, tracingSpan
}

type matchingSpan struct {
	name       string
	status     tracing.SpanStatus
	desc       string
	attributes []tracing.Attribute
	ended      bool
}

func (s *matchingSpan) End() {
	s.ended = true
}

func (s *matchingSpan) SetStatus(code tracing.SpanStatus, desc string) {
	s.status = code
	s.desc = desc
	s.ended = true
}

func (s *matchingSpan) SetAttributes(attrs ...tracing.Attribute) {
	s.attributes = append(s.attributes, attrs...)
}

func attributeValueForKey(attributes []tracing.Attribute, key string) any {
	i := slices.IndexFunc(attributes, func(attr tracing.Attribute) bool {
		return attr.Key == key
	})

	if i < 0 {
		return nil
	}

	return attributes[i].Value
}
