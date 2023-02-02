//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// httpGetterInfo is passed to an HTTPGetter function passing it parameters
// that should be used to make an HTTP GET request.
type httpGetterInfo struct {
	Range HTTPRange

	// ETag specifies the resource's etag that should be used when creating
	// the HTTP GET request's If-Match header
	ETag *azcore.ETag
}

// RetryReaderOptions configures the retry reader's behavior.
// Zero-value fields will have their specified default values applied during use.
// This allows for modification of a subset of fields.
type RetryReaderOptions struct {
	// MaxRetries specifies the maximum number of attempts a failed read will be retried
	// before producing an error.
	// The default value is three.
	MaxRetries int32

	// OnFailedRead, when non-nil, is called after any failure to read. Expected usage is diagnostic logging.
	OnFailedRead func(failureCount int32, lastError error, rnge HTTPRange, willRetry bool)

	// EarlyCloseAsError can be set to true to prevent retries after "read on closed response body". By default,
	// retryReader has the following special behaviour: closing the response body before it is all read is treated as a
	// retryable error. This is to allow callers to force a retry by closing the body from another goroutine (e.g. if the =
	// read is too slow, caller may want to force a retry in the hope that the retry will be quicker).  If
	// TreatEarlyCloseAsError is true, then retryReader's special behaviour is suppressed, and "read on closed body" is instead
	// treated as a fatal (non-retryable) error.
	// Note that setting TreatEarlyCloseAsError only guarantees that Closing will produce a fatal error if the Close happens
	// from the same "thread" (goroutine) as Read.  Concurrent Close calls from other goroutines may instead produce network errors
	// which will be retried.
	// The default value is false.
	EarlyCloseAsError bool

	doInjectError      bool
	doInjectErrorRound int32
	injectedError      error
}
