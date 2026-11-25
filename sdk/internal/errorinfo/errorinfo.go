// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package errorinfo

import "errors"

// AuthenticationFailureInfo describes a credential error for classification by Azure SDK modules.
// Fields may be added. Define each so that its zero value means "no" or "unknown", because
// a credential module built against an earlier version won't set it.
type AuthenticationFailureInfo struct {
	// Unavailable is true when the credential couldn't attempt authentication.
	Unavailable bool
	// Chained is true when the error aggregates a credential chain's failures. Classify
	// each of Causes individually rather than this error's message.
	Chained bool
	// Causes are the errors that caused this one, if any.
	Causes []error
}

// AuthenticationFailure is implemented by credential errors so that Azure SDK modules can classify
// them without depending on the credential module. Unlike Unwrap, it doesn't expose causes to
// errors.Is and errors.As. Use errors.As() to access this interface in the error chain.
type AuthenticationFailure interface {
	error
	AuthenticationFailure() AuthenticationFailureInfo
}

// NonRetriable represents a non-transient error.  This works in
// conjunction with the retry policy, indicating that the error condition
// is idempotent, so no retries will be attempted.
// Use errors.As() to access this interface in the error chain.
type NonRetriable interface {
	error
	NonRetriable()
}

// NonRetriableError ensures the specified error is [NonRetriable]. If
// the error is already [NonRetriable], it returns that error unchanged.
// Otherwise, it returns a new, [NonRetriable] error.
func NonRetriableError(err error) error {
	if !errors.As(err, new(NonRetriable)) {
		err = &nonRetriableError{err}
	}
	return err
}

// nonRetriableError is a struct that embeds the error interface.
// It is used to represent errors that should not be retried.
type nonRetriableError struct {
	error
}

// Error method for nonRetriableError struct.
// It returns the error message of the embedded error.
func (p *nonRetriableError) Error() string {
	return p.error.Error()
}

// NonRetriable is a marker method for nonRetriableError struct.
// Non-functional and indicates that the error is non-retriable.
func (*nonRetriableError) NonRetriable() {
	// marker method
}

// Unwrap method for nonRetriableError struct.
// It returns the original error that was marked as non-retriable.
func (p *nonRetriableError) Unwrap() error {
	return p.error
}
