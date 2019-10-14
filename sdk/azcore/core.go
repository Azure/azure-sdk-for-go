// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"io"
)

// The Policy interface represents a mutable Policy object created by a Factory. The object can mutate/process
// the HTTP request and then forward it on to the next Policy object in the linked-list. The returned
// Response goes backward through the linked-list for additional processing.
// NOTE: Request is passed by value so changes do not change the caller's version of
// the request. However, Request has some fields that reference mutable objects (not strings).
// These references are copied; a deep copy is not performed. Specifically, this means that
// you should avoid modifying the objects referred to by these fields: URL, Header, Body,
// GetBody, TransferEncoding, Form, MultipartForm, Trailer, TLS, Cancel, and Response.
type Policy interface {
	Do(ctx context.Context, req *Request) (*Response, error)
}

// NewPipeline creates a new goroutine-safe Pipeline object from the slice of Factory objects and the specified options.
func NewPipeline(policies ...Policy) Pipeline {
	return Pipeline{
		policies: append(policies[0:0:0], policies...), // append to an empty slice all the policies => makes a copy of the policy slice so it's immutable
	}
}

// Pipeline represents a linked-list of policies and pipeline options that apply to all policies in the pipeline.
// You construct a Pipeline by calling the pipeline.NewPipeline function. To send an HTTP request, call pipeline.NewRequest
// and then call Pipeline's Do method passing a context, the request, and a method-specific Factory (or nil). Passing a
// method-specific Factory allows this one call to Do to inject a Policy into the linked-list. The policy is injected where
// the MethodFactoryMarker (see the pipeline.MethodFactoryMarker function) is in the slice of Factory objects.
//
// When Do is called, the Pipeline object asks each Factory object to construct its Policy object and adds each Policy to a linked-list.
// THen, Do sends the Context and Request through all the Policy objects. The final Policy object sends the request over the network
// (via the HTTPSender object passed to NewPipeline) and the response is returned backwards through all the Policy objects.
// Since Pipeline and Factory objects are goroutine-safe, you typically create 1 Pipeline object and reuse it to make many HTTP requests.
type Pipeline struct {
	policies []Policy
}

// PipelineOptions is used to configure a request policy pipeline's retry policy and logging.
type PipelineOptions struct {
	// Retry configures the built-in retry policy behavior.
	Retry RetryOptions

	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry TelemetryOptions

	HTTPClient Policy

	LogOptions RequestLogOptions
}

type ReadSeekCloser interface {
	io.ReadCloser
	io.Seeker
}

type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}

func NopCloser(rs io.ReadSeeker) ReadSeekCloser {
	return nopCloser{rs}
}
