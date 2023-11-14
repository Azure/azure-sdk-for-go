// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/stretchr/testify/require"
)

// NewSpanValidator creates a tracing.Provider that verifies a span was created that matches the specified SpanMatcher.
func NewSpanValidator(t *testing.T, matcher SpanMatcher) tracing.Provider {
	return tracing.NewProvider(func(name, version string) tracing.Tracer {
		tt := matchingTracer{
			matcher: matcher,
		}

		t.Cleanup(func() {
			require.NotNil(t, tt.match, "didn't find a span with name %s", tt.matcher.Name)
			require.True(t, tt.match.ended, "span wasn't ended")
			require.EqualValues(t, matcher.Status, tt.match.status, "span status values don't match")
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
type SpanMatcher struct {
	Name   string
	Status tracing.SpanStatus
}

type matchingTracer struct {
	matcher SpanMatcher
	match   *span
}

func (mt *matchingTracer) Start(ctx context.Context, spanName string, kind tracing.SpanKind) (context.Context, tracing.Span) {
	if spanName != mt.matcher.Name {
		return ctx, tracing.Span{}
	}
	// span name matches our matcher, track it
	mt.match = &span{
		name: spanName,
	}
	return ctx, tracing.NewSpan(tracing.SpanImpl{
		End:       mt.match.End,
		SetStatus: mt.match.SetStatus,
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
