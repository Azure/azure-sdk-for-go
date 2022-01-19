// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/rpc"
	"github.com/Azure/go-amqp"
	"github.com/devigned/tab"
)

type NonRetriable interface {
	error
	NonRetriable()
}

// IsNonRetriable indicates an error is fatal. Typically, this means
// the connection or link has been closed.
func IsNonRetriable(err error) bool {
	var nonRetriable NonRetriable
	return errors.As(err, &nonRetriable)
}

type ErrNonRetriable struct {
	Message string
}

func (e ErrNonRetriable) Error() string { return e.Message }

// Error Conditions
const (
	// Service Bus Errors
	errorServerBusy         amqp.ErrorCondition = "com.microsoft:server-busy"
	errorTimeout            amqp.ErrorCondition = "com.microsoft:timeout"
	errorOperationCancelled amqp.ErrorCondition = "com.microsoft:operation-cancelled"
	errorContainerClose     amqp.ErrorCondition = "com.microsoft:container-close"
)

const (
	amqpRetryDefaultTimes int           = 3
	amqpRetryDefaultDelay time.Duration = time.Second
)

type (
	// ErrMissingField indicates that an expected property was missing from an AMQP message. This should only be
	// encountered when there is an error with this library, or the server has altered its behavior unexpectedly.
	ErrMissingField string

	// ErrMalformedMessage indicates that a message was expected in the form of []byte was not a []byte. This is likely
	// a bug and should be reported.
	ErrMalformedMessage string

	// ErrIncorrectType indicates that type assertion failed. This should only be encountered when there is an error
	// with this library, or the server has altered its behavior unexpectedly.
	ErrIncorrectType struct {
		Key          string
		ExpectedType reflect.Type
		ActualValue  interface{}
	}

	// ErrAMQP indicates that the server communicated an AMQP error with a particular
	ErrAMQP rpc.Response

	// ErrNoMessages is returned when an operation returned no messages. It is not indicative that there will not be
	// more messages in the future.
	ErrNoMessages struct{}

	// ErrNotFound is returned when an entity is not found (404)
	ErrNotFound struct {
		EntityPath string
	}

	// ErrConnectionClosed indicates that the connection has been closed.
	ErrConnectionClosed string
)

func (e ErrMissingField) Error() string {
	return fmt.Sprintf("missing value %q", string(e))
}

func (e ErrMalformedMessage) Error() string {
	return "message was expected in the form of []byte was not a []byte"
}

// NewErrIncorrectType lets you skip using the `reflect` package. Just provide a variable of the desired type as
// 'expected'.
func NewErrIncorrectType(key string, expected, actual interface{}) ErrIncorrectType {
	return ErrIncorrectType{
		Key:          key,
		ExpectedType: reflect.TypeOf(expected),
		ActualValue:  actual,
	}
}

func (e ErrIncorrectType) Error() string {
	return fmt.Sprintf(
		"value at %q was expected to be of type %q but was actually of type %q",
		e.Key,
		e.ExpectedType,
		reflect.TypeOf(e.ActualValue))
}

func (e ErrAMQP) Error() string {
	return fmt.Sprintf("server says (%d) %s", e.Code, e.Description)
}

func (e ErrNoMessages) Error() string {
	return "no messages available"
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("entity at %s not found", e.EntityPath)
}

// IsErrNotFound returns true if the error argument is an ErrNotFound type
func IsErrNotFound(err error) bool {
	_, ok := err.(ErrNotFound)
	return ok
}

func (e ErrConnectionClosed) Error() string {
	return fmt.Sprintf("the connection has been closed: %s", string(e))
}

// Leveraging @serbrech's fine work from go-shuttle:
// https://github.com/Azure/go-shuttle/blob/ea882947109ade9b34d4d69642fdf7aec4570fee/common/errorhandling/recovery.go

var retryableAMQPConditions = map[string]bool{
	string(amqp.ErrorInternalError):         true,
	string(errorServerBusy):                 true, // "com.microsoft:server-busy"
	string(errorTimeout):                    true, // "com.microsoft:timeout"
	string(errorOperationCancelled):         true, // "com.microsoft:operation-cancelled"
	"client.sender:not-enough-link-credit":  true,
	string(amqp.ErrorUnauthorizedAccess):    true,
	string(amqp.ErrorDetachForced):          true,
	string(amqp.ErrorConnectionForced):      true,
	string(amqp.ErrorTransferLimitExceeded): true,
	"amqp: connection closed":               true,
	"unexpected frame":                      true,
	string(amqp.ErrorNotFound):              true,
}

func isRetryableAMQPError(ctxForLogging context.Context, err error) bool {
	var amqpErr *amqp.Error
	var isAMQPError = errors.As(err, &amqpErr)

	if isAMQPError {
		_, ok := retryableAMQPConditions[string(amqpErr.Condition)]
		return ok
	}

	// TODO: there is a bug somewhere that seems to be errorString'ing errors. Need to track that down.
	// In the meantime, try string matching instead
	for condition := range retryableAMQPConditions {
		if strings.Contains(err.Error(), condition) {
			tab.For(ctxForLogging).Error(fmt.Errorf("error needed to be matched by a string matcher, rather than by type: %w", err))
			return true
		}
	}

	return false
}

func shouldRecreateLink(err error) bool {
	if err == nil {
		return false
	}

	var detachError *amqp.DetachError

	return errors.As(err, &detachError) ||
		// TODO: proper error types needs to happen
		strings.Contains(err.Error(), "detach frame link detached")
}

func shouldRecreateConnection(ctxForLogging context.Context, err error) bool {
	if err == nil {
		return false
	}

	shouldRecreate := isPermanentNetError(err) ||
		isRetryableAMQPError(ctxForLogging, err) ||
		errors.Is(err, io.EOF) ||
		// these are distinct from a detach and probably indicate something
		// wrong with the connection itself, rather than just the link
		errors.Is(err, amqp.ErrSessionClosed) ||
		errors.Is(err, amqp.ErrLinkClosed)

	return shouldRecreate
}

type recoveryKind string

const RecoveryKindNone recoveryKind = ""
const RecoveryKindFatal recoveryKind = "fatal"
const RecoveryKindLink recoveryKind = "link"
const RecoveryKindConn recoveryKind = "connection"

type ServiceBusError struct {
	inner        error
	RecoveryKind recoveryKind
}

func (sbe *ServiceBusError) String() string {
	return sbe.inner.Error()
}

func (sbe *ServiceBusError) AsError() error {
	return sbe.inner
}

// ToSBE wraps the passed in 'err' with a proper error with one of either:
// - `fatalServiceBusError` if no recovery is possible.
// - `serviceBusError` if the error is recoverable. The `recoveryKind` field contains the
//   type of recovery needed.
func ToSBE(loggingCtx context.Context, err error) *ServiceBusError {
	if err == nil {
		return nil
	}

	sbe := &ServiceBusError{
		inner:        err,
		RecoveryKind: GetRecoveryKind(loggingCtx, err),
	}

	return sbe
}

func isPermanentNetError(err error) bool {
	var netErr net.Error

	if errors.As(err, &netErr) {
		temp := netErr.Temporary()
		timeout := netErr.Timeout()
		return !temp && !timeout
	}

	return false
}

func IsCancelError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	if err.Error() == "context canceled" { // go-amqp is returning this when I cancel
		return true
	}

	return false
}

func IsDrainingError(err error) bool {
	// TODO: we should be able to identify these errors programatically
	return strings.Contains(err.Error(), "link is currently draining")
}

var amqpConditionsToRecoveryKind = map[amqp.ErrorCondition]recoveryKind{
	// no recovery needed, these are temporary errors.
	amqp.ErrorCondition("com.microsoft:server-busy"):         RecoveryKindNone,
	amqp.ErrorCondition("com.microsoft:timeout"):             RecoveryKindNone,
	amqp.ErrorCondition("com.microsoft:operation-cancelled"): RecoveryKindNone,

	// Link recovery needed
	amqp.ErrorDetachForced: RecoveryKindLink, // "amqp:link:detach-forced"

	// Connection recovery needed
	amqp.ErrorConnectionForced: RecoveryKindConn, // "amqp:connection:forced"

	// No recovery possible - this operation is non retriable.
	amqp.ErrorMessageSizeExceeded:                                 RecoveryKindFatal,
	amqp.ErrorUnauthorizedAccess:                                  RecoveryKindFatal, // creds are bad
	amqp.ErrorNotFound:                                            RecoveryKindFatal,
	amqp.ErrorNotAllowed:                                          RecoveryKindFatal,
	amqp.ErrorInternalError:                                       RecoveryKindFatal, // "amqp:internal-error"
	amqp.ErrorCondition("com.microsoft:entity-disabled"):          RecoveryKindFatal, // entity is disabled in the portal
	amqp.ErrorCondition("com.microsoft:session-cannot-be-locked"): RecoveryKindFatal,
	amqp.ErrorCondition("com.microsoft:message-lock-lost"):        RecoveryKindFatal,
}

func GetRecoveryKind(ctxForLogging context.Context, err error) recoveryKind {
	if IsCancelError(err) {
		return RecoveryKindFatal
	}

	var netErr net.Error

	if errors.As(err, &netErr) {
		// ie, just retry
		return RecoveryKindNone
	}

	// this is a carryover from another library. I haven't seen this in the wild.
	if errors.Is(err, io.EOF) {
		return RecoveryKindConn
	}

	var errNonRetriable *ErrNonRetriable

	if errors.As(err, &errNonRetriable) {
		return RecoveryKindFatal
	}

	var de *amqp.DetachError

	// check the "special" AMQP errors that aren't condition-based.
	if errors.Is(err, amqp.ErrSessionClosed) ||
		errors.Is(err, amqp.ErrLinkClosed) ||
		errors.As(err, &de) {
		return RecoveryKindLink
	}

	if errors.Is(err, amqp.ErrConnClosed) {
		return RecoveryKindConn
	}

	if IsDrainingError(err) {
		// temporary, operation should just be retryable since drain will
		// eventually complete.
		return RecoveryKindNone
	}

	// then it's _probably_ an actual *amqp.Error, in which case we bucket it by
	// the 'condition'.
	var amqpError *amqp.Error

	if errors.As(err, &amqpError) {
		recoveryKind, ok := amqpConditionsToRecoveryKind[amqpError.Condition]

		if ok {
			return recoveryKind
		}
	}

	// this is some error type we've never seen.
	tab.For(ctxForLogging).Fatal(fmt.Sprintf("No recovery possible with error: %#v", err))
	return RecoveryKindFatal
}
