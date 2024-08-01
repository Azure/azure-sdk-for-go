// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"slices"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/stretchr/testify/require"
)

// newSpanValidator creates a tracing.Provider that verifies a span was created that matches the specified SpanMatcher.
func newSpanValidator(t *testing.T, matcher spanMatcher) tracing.Provider {
	return tracing.NewProvider(func(name, version string) tracing.Tracer {
		tt := matchingTracer{
			matcher: matcher,
		}

		t.Cleanup(func() {
			for _, expectedSpan := range matcher.ExpectedSpans {
				found := false
				for _, match := range tt.matches {
					if match.name == expectedSpan {
						found = true
						require.True(t, match.ended, "span %s wasn't ended", match.name)
						break;
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
		}, nil)
	}, nil)
}

// SpanMatcher contains the values to match when a span is created.
type spanMatcher struct {
	ExpectedSpans []string
}

type matchingTracer struct {
	matcher spanMatcher
	matches []*span
}

func (mt *matchingTracer) Start(ctx context.Context, spanName string, kind tracing.SpanKind) (context.Context, tracing.Span) {

	if slices.IndexFunc(mt.matcher.ExpectedSpans, func(i string) bool { return i == spanName }) < 0 {
		return ctx, tracing.Span{}
	}
	// span name matches our matcher, track it
	newSpan := &span{
		name: spanName,
	}
	mt.matches = append(mt.matches, newSpan)
	return ctx, tracing.NewSpan(tracing.SpanImpl{
		End:       newSpan.End,
		SetStatus: newSpan.SetStatus,
	})
}

type span struct {
	name   string
	status tracing.SpanStatus
	desc   string
	ended  bool
}

func (s *span) End() {
	s.ended = true
}

func (s *span) SetStatus(code tracing.SpanStatus, desc string) {
	s.status = code
	s.desc = desc
	s.ended = true
}