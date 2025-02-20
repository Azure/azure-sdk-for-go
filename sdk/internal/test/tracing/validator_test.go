// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/stretchr/testify/require"
)

func ExampleNewSpanValidator() {
	t := &testing.T{}

	// attributes and links used when starting the span
	initialAttr := tracing.Attribute{Key: "initialAttrKey", Value: "initialAttrValue"}
	initialLink := tracing.Link{Attributes: []tracing.Attribute{{Key: "initialLinkKey", Value: "initialLinkValue"}}}
	// attributes and links added after starting the span
	testAttr := tracing.Attribute{Key: "testSetAttrKey", Value: "testSetAttrValue"}
	testLink := tracing.Link{Attributes: []tracing.Attribute{{Key: "testAddLinkKey", Value: "testAddLinkValue"}}}

	// create a span validator that will verify spans with the given name, kind, status, attributes, and links
	// the span matcher will verify the span at the end of the tests
	provider := NewSpanValidator(t, SpanMatcher{
		Name:       "TestSpan",
		Kind:       tracing.SpanKindClient,
		Status:     tracing.SpanStatusUnset,
		Attributes: []tracing.Attribute{initialAttr, testAttr},
		Links:      []tracing.Link{initialLink, testLink},
	}, nil)
	tracer := provider.NewTracer("module", "version")

	// start a test span with initial attributes and links
	ctx, endSpan := runtime.StartSpan(context.Background(), "TestSpan", tracer, &runtime.StartSpanOptions{
		Kind:       tracing.SpanKindClient,
		Attributes: []tracing.Attribute{initialAttr},
		Links:      []tracing.Link{initialLink},
	})
	defer func() { endSpan(nil) }()

	// get the created span from context and add attributes and links
	// they will get verified with SpanValidator provider
	spn := tracer.SpanFromContext(ctx)
	spn.SetAttributes(testAttr)
	spn.AddLink(testLink)
}

func TestNewSpanValidator(t *testing.T) {
	// attributes and links used when starting the span
	initialAttr := tracing.Attribute{Key: "initialAttrKey", Value: "initialAttrValue"}
	initialLink := tracing.Link{Attributes: []tracing.Attribute{{Key: "initialLinkKey", Value: "initialLinkValue"}}}
	// attributes and links added after starting the span
	testAttr := tracing.Attribute{Key: "testSetAttrKey", Value: "testSetAttrValue"}
	testLink := tracing.Link{Attributes: []tracing.Attribute{{Key: "testAddLinkKey", Value: "testAddLinkValue"}}}

	provider := NewSpanValidator(t, SpanMatcher{
		Name:   "TestSpan",
		Kind:   tracing.SpanKindClient,
		Status: tracing.SpanStatusError,
		// error.type is also expected because the span ended with an error
		Attributes: []tracing.Attribute{initialAttr, testAttr, {Key: "error.type", Value: "*errors.errorString"}},
		Links:      []tracing.Link{initialLink, testLink},
	}, nil)
	tracer := provider.NewTracer("module", "version")
	require.NotNil(t, tracer)
	require.True(t, tracer.Enabled())

	ctx, endSpan := runtime.StartSpan(context.Background(), "TestSpan", tracer, &runtime.StartSpanOptions{
		Kind:       tracing.SpanKindClient,
		Attributes: []tracing.Attribute{initialAttr},
		Links:      []tracing.Link{initialLink},
	})
	spn := tracer.SpanFromContext(ctx)
	spn.SetAttributes(testAttr)
	spn.AddLink(testLink)
	defer func() { endSpan(errors.New("test error")) }()

	require.NotNil(t, tracer.SpanFromContext(ctx))
	require.NotNil(t, tracer.LinkFromContext(ctx))
	require.Zero(t, provider.NewPropagator())
}

func TestMatchingTracerStart(t *testing.T) {
	matcher := SpanMatcher{
		Name:   "TestSpan",
		Kind:   tracing.SpanKindProducer,
		Status: tracing.SpanStatusUnset,
		Attributes: []tracing.Attribute{
			{Key: "testKey1", Value: "testValue1"},
			{Key: "testKey2", Value: "testValue2"},
		},
	}
	tracer := matchingTracer{
		matcher: matcher,
	}
	ctx := context.Background()
	// no-op when SpanName doesn't match
	_, spn := tracer.Start(ctx, "BadSpanName", tracing.SpanKindProducer, nil, nil)
	spn.End()
	require.EqualValues(t, spn, tracing.Span{})
	// tracks span when SpanName matches
	_, spn = tracer.Start(ctx, "TestSpan", tracing.SpanKindProducer, []tracing.Attribute{
		{Key: "testKey1", Value: "testValue1"},
		{Key: "testKey2", Value: "testValue2"},
	}, nil)
	require.NotNil(t, spn)
	spn.SetAttributes(tracing.Attribute{
		Key:   "setAttributeKey",
		Value: "setAttributeValue",
	})
	spn.AddLink(tracing.Link{})
	spn.SetStatus(tracing.SpanStatusOK, "ok")
	spn.End()
}
