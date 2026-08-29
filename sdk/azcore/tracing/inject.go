// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"context"
	"net/http"
)

// HeaderInjector propagates the active span context from ctx into header.
type HeaderInjector func(ctx context.Context, header http.Header)
