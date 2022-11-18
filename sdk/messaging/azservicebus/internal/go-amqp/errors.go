// Copyright (C) 2017 Kale Blankenship
// Portions Copyright (c) Microsoft Corporation

package amqp

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp/internal/encoding"
)

// ErrCond is an AMQP defined error condition.
// See http://docs.oasis-open.org/amqp/core/v1.0/os/amqp-core-transport-v1.0-os.html#type-amqp-error for info on their meaning.
type ErrCond = encoding.ErrCond

// Error Conditions
const (
	// AMQP Errors
	ErrCondDecodeError           ErrCond = "amqp:decode-error"
	ErrCondFrameSizeTooSmall     ErrCond = "amqp:frame-size-too-small"
	ErrCondIllegalState          ErrCond = "amqp:illegal-state"
	ErrCondInternalError         ErrCond = "amqp:internal-error"
	ErrCondInvalidField          ErrCond = "amqp:invalid-field"
	ErrCondNotAllowed            ErrCond = "amqp:not-allowed"
	ErrCondNotFound              ErrCond = "amqp:not-found"
	ErrCondNotImplemented        ErrCond = "amqp:not-implemented"
	ErrCondPreconditionFailed    ErrCond = "amqp:precondition-failed"
	ErrCondResourceDeleted       ErrCond = "amqp:resource-deleted"
	ErrCondResourceLimitExceeded ErrCond = "amqp:resource-limit-exceeded"
	ErrCondResourceLocked        ErrCond = "amqp:resource-locked"
	ErrCondUnauthorizedAccess    ErrCond = "amqp:unauthorized-access"

	// Connection Errors
	ErrCondConnectionForced   ErrCond = "amqp:connection:forced"
	ErrCondConnectionRedirect ErrCond = "amqp:connection:redirect"
	ErrCondFramingError       ErrCond = "amqp:connection:framing-error"

	// Session Errors
	ErrCondErrantLink       ErrCond = "amqp:session:errant-link"
	ErrCondHandleInUse      ErrCond = "amqp:session:handle-in-use"
	ErrCondUnattachedHandle ErrCond = "amqp:session:unattached-handle"
	ErrCondWindowViolation  ErrCond = "amqp:session:window-violation"

	// Link Errors
	ErrCondDetachForced          ErrCond = "amqp:link:detach-forced"
	ErrCondLinkRedirect          ErrCond = "amqp:link:redirect"
	ErrCondMessageSizeExceeded   ErrCond = "amqp:link:message-size-exceeded"
	ErrCondStolen                ErrCond = "amqp:link:stolen"
	ErrCondTransferLimitExceeded ErrCond = "amqp:link:transfer-limit-exceeded"
)

// Error is an AMQP error.
// DetachError and SessionError will contain instances of this type
// when detached/closed by the peer with an AMQP error. In this case,
// use errors.As() to unwrap the inner *Error.
type Error = encoding.Error

// DetachError is returned by methods on Sender/Receiver when the link has become detached/closed.
type DetachError struct {
	inner error
}

// Error implements the error interface for DetachError.
func (e *DetachError) Error() string {
	if e.inner == nil {
		return "amqp: link closed"
	}
	return e.inner.Error()
}

// Unwrap returns the inner *Error or nil.
func (e *DetachError) Unwrap() error {
	var err *Error
	if errors.As(e.inner, &err) {
		return err
	}
	return nil
}

// ConnectionError is propagated to Session and Senders/Receivers
// when the connection has been closed.
type ConnectionError struct {
	inner error
}

// Error implements the error interface for ConnectionError.
func (c *ConnectionError) Error() string {
	if c.inner == nil {
		return "amqp: connection closed"
	}
	return c.inner.Error()
}

// Unwrap returns the inner *Error or nil.
func (e *ConnectionError) Unwrap() error {
	var err *Error
	if errors.As(e.inner, &err) {
		return err
	}
	return nil
}

// SessionError is propagated to Senders/Receivers when the session
// has been closed.
type SessionError struct {
	inner error
}

// Error implements the error interface for SessionError.
func (s *SessionError) Error() string {
	if s.inner == nil {
		return "amqp: session closed"
	}
	return s.inner.Error()
}

// Unwrap returns the inner *Error or nil.
func (e *SessionError) Unwrap() error {
	var err *Error
	if errors.As(e.inner, &err) {
		return err
	}
	return nil
}
