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

func Test_recoveryKind(t *testing.T) {
	t.Run("link", func(t *testing.T) {
		linkErrorCodes := []string{
			string(amqp.ErrorDetachForced),
		}

		for _, code := range linkErrorCodes {
			t.Run(code, func(t *testing.T) {
				sbe := GetSBErrInfo(&amqp.Error{Condition: amqp.ErrorCondition(code)})
				require.EqualValues(t, RecoveryKindLink, sbe.RecoveryKind, fmt.Sprintf("requires link recovery: %s", code))
			})
		}

		t.Run("sentintel errors", func(t *testing.T) {
			sbe := GetSBErrInfo(amqp.ErrLinkClosed)
			require.EqualValues(t, RecoveryKindLink, sbe.RecoveryKind)

			sbe = GetSBErrInfo(amqp.ErrSessionClosed)
			require.EqualValues(t, RecoveryKindLink, sbe.RecoveryKind)
		})
	})

	t.Run("connection", func(t *testing.T) {
		codes := []string{
			string(amqp.ErrorConnectionForced),
		}

		for _, code := range codes {
			t.Run(code, func(t *testing.T) {
				sbe := GetSBErrInfo(&amqp.Error{Condition: amqp.ErrorCondition(code)})
				require.EqualValues(t, RecoveryKindConn, sbe.RecoveryKind, fmt.Sprintf("requires connection recovery: %s", code))
			})
		}

		t.Run("sentinel errors", func(t *testing.T) {
			sbe := GetSBErrInfo(amqp.ErrConnClosed)
			require.EqualValues(t, RecoveryKindConn, sbe.RecoveryKind)
		})
	})

	t.Run("nonretriable", func(t *testing.T) {
		codes := []string{
			string(amqp.ErrorTransferLimitExceeded),
			string(amqp.ErrorInternalError),
			string(amqp.ErrorUnauthorizedAccess),
			string(amqp.ErrorNotFound),
			string(amqp.ErrorMessageSizeExceeded),
		}

		for _, code := range codes {
			t.Run(code, func(t *testing.T) {
				sbe := GetSBErrInfo(&amqp.Error{Condition: amqp.ErrorCondition(code)})
				require.EqualValues(t, RecoveryKindFatal, sbe.RecoveryKind, fmt.Sprintf("cannot be recovered: %s", code))
			})
		}
	})

	t.Run("none", func(t *testing.T) {
		codes := []string{
			string("com.microsoft:operation-cancelled"),
			string("com.microsoft:server-busy"),
			string("com.microsoft:timeout"),
		}

		for _, code := range codes {
			t.Run(code, func(t *testing.T) {
				sbe := GetSBErrInfo(&amqp.Error{Condition: amqp.ErrorCondition(code)})
				require.EqualValues(t, RecoveryKindNone, sbe.RecoveryKind, fmt.Sprintf("no recovery needed: %s", code))
			})
		}
	})
}

func Test_IsNonRetriable(t *testing.T) {
	errs := []error{
		context.Canceled,
		context.DeadlineExceeded,
		ErrNonRetriable{Message: "any message"},
		fmt.Errorf("wrapped: %w", context.Canceled),
		fmt.Errorf("wrapped: %w", context.DeadlineExceeded),
		fmt.Errorf("wrapped: %w", ErrNonRetriable{Message: "any message"}),
	}

	for _, err := range errs {
		require.EqualValues(t, RecoveryKindFatal, GetSBErrInfo(err).RecoveryKind)
	}
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
		// simple timeouts from the mgmt link
		mgmtError{Resp: &RPCResponse{Code: 408}},
		mgmtError{Resp: &RPCResponse{Code: 503}},
		mgmtError{Resp: &RPCResponse{Code: 500}},
	}

	for i, err := range tempErrors {
		rk := GetSBErrInfo(err).RecoveryKind
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
		rk := GetSBErrInfo(err).RecoveryKind
		require.EqualValues(t, RecoveryKindConn, rk, fmt.Sprintf("[%d] %v", i, err))
	}
}

func Test_ServiceBusError_LinkRecoveryNeeded(t *testing.T) {
	var linkErrors = []error{
		amqp.ErrSessionClosed,
		amqp.ErrLinkClosed,
		&amqp.DetachError{},
		&amqp.Error{Condition: amqp.ErrorDetachForced},
		// we lost the session lock, attempt link recovery
		mgmtError{Resp: &RPCResponse{Code: 410}},
		// this can happen when we're recovering the link - the client gets closed and the old link is still being
		// used by this instance of the client. It needs to recover and attempt it again.
		mgmtError{Resp: &RPCResponse{Code: 401}},
	}

	for i, err := range linkErrors {
		rk := GetSBErrInfo(err).RecoveryKind
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
		rk := GetSBErrInfo(&amqp.Error{Condition: cond}).RecoveryKind
		require.EqualValues(t, RecoveryKindFatal, rk, fmt.Sprintf("[%d] %s", i, cond))
	}

	// unknown errors are also considered fatal
	rk := GetSBErrInfo(errors.New("Some unknown error")).RecoveryKind
	require.EqualValues(t, RecoveryKindFatal, rk, "some unknown error")
}
