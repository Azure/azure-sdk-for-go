// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Diagnostics contains request diagnostics for an Azure Cosmos DB operation.
// The diagnostics string is materialized lazily when String is called.
type Diagnostics struct {
	root *trace
}

func newDiagnostics(t *trace) Diagnostics {
	if t == nil {
		return Diagnostics{}
	}

	return Diagnostics{root: t.root()}
}

// String returns the lazily materialized diagnostics payload.
func (d Diagnostics) String() string {
	if d.root == nil {
		return ""
	}

	d.root.setWalkingStateRecursively()
	return writeTraceJSON(d.root)
}

// ClientElapsedTime returns the end-to-end duration of the request.
func (d Diagnostics) ClientElapsedTime() time.Duration {
	if d.root == nil {
		return 0
	}

	return d.root.duration()
}

// StartTimeUTC returns the UTC start time of the request.
func (d Diagnostics) StartTimeUTC() *time.Time {
	if d.root == nil {
		return nil
	}

	snapshot := d.root.snapshot()
	start := snapshot.startTime
	return &start
}

// FailedRequestCount returns the number of failed backend attempts recorded for the request.
func (d Diagnostics) FailedRequestCount() int {
	if d.root == nil || d.root.summary == nil {
		return 0
	}

	return d.root.summary.failedCount()
}

// DiagnosticsFromError extracts diagnostics from an operation error, if present.
func DiagnosticsFromError(err error) (Diagnostics, bool) {
	if err == nil {
		return Diagnostics{}, false
	}

	var reqErr *requestError
	if errors.As(err, &reqErr) {
		return reqErr.diagnostics, reqErr.diagnostics.root != nil
	}

	var responseErr *azcore.ResponseError
	if errors.As(err, &responseErr) && responseErr.RawResponse != nil {
		diagnostics := diagnosticsFromResponse(responseErr.RawResponse)
		return diagnostics, diagnostics.root != nil
	}

	return Diagnostics{}, false
}

func diagnosticsFromResponse(resp *http.Response) Diagnostics {
	if resp == nil || resp.Request == nil {
		return Diagnostics{}
	}

	state := requestDiagnosticsStateFromContext(resp.Request.Context())
	if state != nil {
		return newDiagnostics(state.requestTrace)
	}

	return diagnosticsFromContext(resp.Request.Context())
}

type requestError struct {
	cause       error
	diagnostics Diagnostics
}

func (e *requestError) Error() string {
	return e.cause.Error()
}

func (e *requestError) Unwrap() error {
	return e.cause
}

func wrapRequestError(err error, diagnostics Diagnostics) error {
	if err == nil || diagnostics.root == nil {
		return err
	}

	return &requestError{
		cause:       err,
		diagnostics: diagnostics,
	}
}
