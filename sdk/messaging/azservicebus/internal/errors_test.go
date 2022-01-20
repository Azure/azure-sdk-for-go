// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrMissingField_Error(t *testing.T) {
	const fieldName = "fieldName"
	var subject ErrMissingField = fieldName
	var cast error = subject

	got := cast.Error()
	const want = `missing value "` + fieldName + `"`

	if got != want {
		t.Logf("\n\tgot: \t%q\n\twant:\t%q", got, want)
		t.Fail()
	}
}

func TestErrIncorrectType_Error(t *testing.T) {
	var a int
	var b map[string]interface{}
	var c *float64

	types := map[reflect.Type]interface{}{
		reflect.TypeOf(a): 7.0,
		reflect.TypeOf(b): map[string]string{},
		reflect.TypeOf(c): int(2),
	}

	const key = "myFieldName"
	for expected, actual := range types {
		actualType := reflect.TypeOf(actual)
		t.Run(fmt.Sprintf("%s-%s", expected, actualType), func(t *testing.T) {
			expectedMessage := fmt.Sprintf(
				"value at %q was expected to be of type %q but was actually of type %q",
				key,
				expected.String(),
				actualType.String())

			subject := ErrIncorrectType{
				Key:          key,
				ActualValue:  actual,
				ExpectedType: expected,
			}

			var cast error = subject

			got := cast.Error()
			if got != expectedMessage {
				t.Logf("\n\tgot: \t%q\n\twant:\t%q", got, expectedMessage)
				t.Fail()
			}
		})
	}
}

func TestErrNotFound_Error(t *testing.T) {
	err := ErrNotFound{EntityPath: "/foo/bar"}
	assert.Equal(t, "entity at /foo/bar not found", err.Error())
	assert.True(t, IsErrNotFound(err))

	otherErr := errors.New("foo")
	assert.False(t, IsErrNotFound(otherErr))
}

func Test_isPermanentNetError(t *testing.T) {
	require.False(t, isPermanentNetError(&fakeNetError{
		temp: true,
	}))

	require.False(t, isPermanentNetError(&fakeNetError{
		timeout: true,
	}))

	require.False(t, isPermanentNetError(errors.New("not a net error")))

	require.True(t, isPermanentNetError(&fakeNetError{}))
}

func Test_isRetryableAMQPError(t *testing.T) {
	ctx := context.Background()

	retryableCodes := []string{
		string(amqp.ErrorInternalError),
		string(errorServerBusy),
		string(errorTimeout),
		string(errorOperationCancelled),
		"client.sender:not-enough-link-credit",
		string(amqp.ErrorUnauthorizedAccess),
		string(amqp.ErrorDetachForced),
		string(amqp.ErrorConnectionForced),
		string(amqp.ErrorTransferLimitExceeded),
		"amqp: connection closed",
		"unexpected frame",
		string(amqp.ErrorNotFound),
	}

	for _, code := range retryableCodes {
		require.True(t, isRetryableAMQPError(ctx, &amqp.Error{
			Condition: amqp.ErrorCondition(code),
		}))

		// it works equally well if the error is just in the String().
		// Need to narrow this down some more to see where the errors
		// might not be getting converted properly.
		require.True(t, isRetryableAMQPError(ctx, errors.New(code)))
	}

	require.False(t, isRetryableAMQPError(ctx, errors.New("some non-amqp related error")))
}

func Test_shouldRecreateLink(t *testing.T) {
	require.False(t, shouldRecreateLink(nil))

	require.True(t, shouldRecreateLink(&amqp.DetachError{}))

	// going to treat these as "connection troubles" and throw them into the
	// connection recovery scenario instead.
	require.False(t, shouldRecreateLink(amqp.ErrLinkClosed))
	require.False(t, shouldRecreateLink(amqp.ErrSessionClosed))
}

func Test_shouldRecreateConnection(t *testing.T) {
	ctx := context.Background()

	require.False(t, shouldRecreateConnection(ctx, nil))
	require.True(t, shouldRecreateConnection(ctx, &fakeNetError{}))
	require.True(t, shouldRecreateConnection(ctx, fmt.Errorf("%w", &fakeNetError{})))

	require.False(t, shouldRecreateLink(amqp.ErrLinkClosed))
	require.False(t, shouldRecreateLink(fmt.Errorf("wrapped: %w", amqp.ErrLinkClosed)))

	require.False(t, shouldRecreateLink(amqp.ErrSessionClosed))
	require.False(t, shouldRecreateLink(fmt.Errorf("wrapped: %w", amqp.ErrSessionClosed)))
}

// TODO: while testing it appeared there were some errors that were getting string-ized
// We want to eliminate these. 'stress.go' reproduces most of these as you disconnect
// and reconnect.
func Test_stringErrorsToEliminate(t *testing.T) {
	require.True(t, shouldRecreateLink(errors.New("detach frame link detached")))
	require.True(t, isRetryableAMQPError(context.Background(), errors.New("amqp: connection closed")))
	require.True(t, IsCancelError(errors.New("context canceled")))
}

func Test_IsCancelError(t *testing.T) {
	require.False(t, IsCancelError(nil))
	require.False(t, IsCancelError(errors.New("not a cancel error")))

	require.True(t, IsCancelError(errors.New("context canceled")))

	require.True(t, IsCancelError(context.Canceled))
	require.True(t, IsCancelError(context.DeadlineExceeded))
	require.True(t, IsCancelError(fmt.Errorf("wrapped: %w", context.Canceled)))
	require.True(t, IsCancelError(fmt.Errorf("wrapped: %w", context.DeadlineExceeded)))
}

func Test_ServiceBusError_NoRecoveryNeeded(t *testing.T) {
	var tempErrors = []error{
		&amqp.Error{Condition: amqp.ErrorCondition("com.microsoft:server-busy")},
		&amqp.Error{Condition: amqp.ErrorCondition("com.microsoft:timeout")},
		&amqp.Error{Condition: amqp.ErrorCondition("com.microsoft:operation-cancelled")},
		errors.New("link is currently draining"), // not yet exposed from go-amqp
		fakeNetError{temp: true},
		fakeNetError{timeout: true},
		fakeNetError{temp: false, timeout: false},
	}

	for i, err := range tempErrors {
		rk := ToSBE(context.Background(), err).RecoveryKind
		require.EqualValues(t, RecoveryKindNone, rk, fmt.Sprintf("[%d] %v", i, err))
	}
}

func Test_ServiceBusError_ConnectionRecoveryNeeded(t *testing.T) {
	var connErrors = []error{
		&amqp.Error{Condition: amqp.ErrorConnectionForced},
		amqp.ErrConnClosed,
		io.EOF,
	}

	for i, err := range connErrors {
		rk := ToSBE(context.Background(), err).RecoveryKind
		require.EqualValues(t, RecoveryKindConn, rk, fmt.Sprintf("[%d] %v", i, err))
	}
}

func Test_ServiceBusError_LinkRecoveryNeeded(t *testing.T) {
	var linkErrors = []error{
		amqp.ErrSessionClosed,
		amqp.ErrLinkClosed,
		&amqp.DetachError{},
		&amqp.Error{Condition: amqp.ErrorDetachForced},
	}

	for i, err := range linkErrors {
		rk := ToSBE(context.Background(), err).RecoveryKind
		require.EqualValues(t, RecoveryKindLink, rk, fmt.Sprintf("[%d] %v", i, err))
	}
}

func Test_ServiceBusError_Fatal(t *testing.T) {
	var fatalConditions = []amqp.ErrorCondition{
		amqp.ErrorMessageSizeExceeded,
		amqp.ErrorUnauthorizedAccess,
		amqp.ErrorNotFound,
		amqp.ErrorNotAllowed,
		amqp.ErrorInternalError,
		amqp.ErrorCondition("com.microsoft:entity-disabled"),
		amqp.ErrorCondition("com.microsoft:session-cannot-be-locked"),
		amqp.ErrorCondition("com.microsoft:message-lock-lost"),
	}

	for i, cond := range fatalConditions {
		rk := ToSBE(context.Background(), &amqp.Error{Condition: cond}).RecoveryKind
		require.EqualValues(t, RecoveryKindFatal, rk, fmt.Sprintf("[%d] %s", i, cond))
	}

	// unknown errors are also considered fatal
	rk := ToSBE(context.Background(), errors.New("Some unknown error")).RecoveryKind
	require.EqualValues(t, RecoveryKindFatal, rk, "some unknown error")
}
