// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

type Attribute = tracing.Attribute
type Tracer = tracing.Tracer
type Provider = tracing.Provider

type StartSpanOptions = runtime.StartSpanOptions

func StartSpan(ctx context.Context, spanName string, tracer Tracer, options *StartSpanOptions) (context.Context, func(error)) {
	return runtime.StartSpan(ctx, spanName, tracer, options)
}
