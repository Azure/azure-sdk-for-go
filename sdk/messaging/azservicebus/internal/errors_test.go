// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/exported"
	"github.com/Azure/go-amqp"
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
	var b map[string]any
	var c *float64

	types := map[reflect.Type]any{
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

func Test_recoveryKind(t *testing.T) {
	t.Run("link", func(t *testing.T) {
		linkErrorCodes := []amqp.ErrCond{
			amqp.ErrCondDetachForced,
		}

		for _, code := range linkErrorCodes {
			amqpConditionTakesPrecedence(t, code, RecoveryKindLink)
		}

		t.Run("sentinel errors", func(t *testing.T) {
			rk := GetRecoveryKind(&amqp.LinkError{})
			require.EqualValues(t, RecoveryKindLink, rk)

			rk = GetRecoveryKind(&amqp.SessionError{})
			require.EqualValues(t, RecoveryKindConn, rk)

			rk = GetRecoveryKind(&amqp.ConnError{})
			require.EqualValues(t, RecoveryKindConn, rk)

		})
	})

	t.Run("connection", func(t *testing.T) {
		codes := []amqp.ErrCond{
			amqp.ErrCondConnectionForced,
			amqp.ErrCondInternalError,
		}

		for _, code := range codes {
			amqpConditionTakesPrecedence(t, code, RecoveryKindConn)
		}

		t.Run("sentinel errors", func(t *testing.T) {
			rk := GetRecoveryKind(&amqp.ConnError{})
			require.EqualValues(t, RecoveryKindConn, rk)
		})
	})

	t.Run("nonretriable", func(t *testing.T) {
		codes := []amqp.ErrCond{
			amqp.ErrCondUnauthorizedAccess,
			amqp.ErrCondNotFound,
			amqp.ErrCondMessageSizeExceeded,
		}

		for _, code := range codes {
			amqpConditionTakesPrecedence(t, code, RecoveryKindFatal)
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
				rk := GetRecoveryKind(&amqp.Error{Condition: amqp.ErrCond(code)})
				require.EqualValues(t, RecoveryKindNone, rk, fmt.Sprintf("no recovery needed: %s", code))
			})
		}
	})
}

func Test_IsNonRetriable(t *testing.T) {
	errs := []error{
		context.Canceled,
		context.DeadlineExceeded,
		NewErrNonRetriable("any message"),
		fmt.Errorf("wrapped: %w", context.Canceled),
		fmt.Errorf("wrapped: %w", context.DeadlineExceeded),
		fmt.Errorf("wrapped: %w", NewErrNonRetriable("any message")),
	}

	for _, err := range errs {
		require.EqualValues(t, RecoveryKindFatal, GetRecoveryKind(err))
	}
}

func Test_ServiceBusError_NoRecoveryNeeded(t *testing.T) {
	var tempErrors = []error{
		&amqp.Error{Condition: amqp.ErrCond("com.microsoft:server-busy")},
		&amqp.Error{Condition: amqp.ErrCond("com.microsoft:timeout")},
		&amqp.Error{Condition: amqp.ErrCond("com.microsoft:operation-cancelled")},
		// simple timeouts from the mgmt link
		RPCError{Resp: &amqpwrap.RPCResponse{Code: 408}},
		RPCError{Resp: &amqpwrap.RPCResponse{Code: 503}},
		RPCError{Resp: &amqpwrap.RPCResponse{Code: 500}},
	}

	for i, err := range tempErrors {
		rk := GetRecoveryKind(err)
		require.EqualValues(t, RecoveryKindNone, rk, fmt.Sprintf("[%d] %v", i, err))
	}
}

func Test_ServiceBusError_ConnectionRecoveryNeeded(t *testing.T) {
	var connErrors = []error{
		&amqp.Error{Condition: amqp.ErrCondConnectionForced},
		&amqp.Error{Condition: amqp.ErrCondInternalError},
		&amqp.ConnError{},
		&amqp.SessionError{},
		io.EOF,
		fakeNetError{temp: true},
		fakeNetError{timeout: true},
		fakeNetError{temp: false, timeout: false},
		errors.New("*frames.PerformTransfer: didn't find channel 10 in sessionsByRemoteChannel"),
	}

	for i, err := range connErrors {
		rk := GetRecoveryKind(err)
		require.EqualValues(t, RecoveryKindConn, rk, fmt.Sprintf("[%d] %v", i, err))
	}

	// unknown errors will just result in a connection recovery
	rk := GetRecoveryKind(errors.New("Some unknown error"))
	require.EqualValues(t, RecoveryKindConn, rk, "some unknown error")

	require.Equal(t, RecoveryKindConn, GetRecoveryKind(amqpwrap.ErrConnResetNeeded))
}

func Test_ServiceBusError_LinkRecoveryNeeded(t *testing.T) {
	var linkErrors = []error{
		&amqp.LinkError{},
		&amqp.Error{Condition: amqp.ErrCondDetachForced},
		&amqp.Error{Condition: amqp.ErrCondTransferLimitExceeded},
	}

	for i, err := range linkErrors {
		rk := GetRecoveryKind(err)
		require.EqualValues(t, RecoveryKindLink, rk, fmt.Sprintf("[%d] %v", i, err))
	}
}

func Test_ServiceBusError_Fatal(t *testing.T) {
	var fatalConditions = []amqp.ErrCond{
		amqp.ErrCondMessageSizeExceeded,
		amqp.ErrCondResourceLimitExceeded,
		amqp.ErrCondUnauthorizedAccess,
		amqp.ErrCondNotFound,
		amqp.ErrCondNotAllowed,
		amqp.ErrCond("com.microsoft:entity-disabled"),
		amqp.ErrCond("com.microsoft:session-cannot-be-locked"),
		amqp.ErrCond("com.microsoft:message-lock-lost"),
	}

	for i, cond := range fatalConditions {
		rk := GetRecoveryKind(&amqp.Error{Condition: cond})
		require.EqualValues(t, RecoveryKindFatal, rk, fmt.Sprintf("[%d] %s", i, cond))
	}

	require.Equal(t, RecoveryKindFatal, GetRecoveryKind(RPCError{Resp: &amqpwrap.RPCResponse{Code: http.StatusNotFound}}))
	require.Equal(t, RecoveryKindFatal, GetRecoveryKind(RPCError{Resp: &amqpwrap.RPCResponse{Code: RPCResponseCodeLockLost}}))
	require.Equal(t, RecoveryKindFatal, GetRecoveryKind(RPCError{Resp: &amqpwrap.RPCResponse{Code: http.StatusUnauthorized}}))
	require.Equal(t, RecoveryKindFatal, GetRecoveryKind(errClosed))

}

func Test_IsLockLostError(t *testing.T) {
	require.True(t, isLockLostError(RPCError{Resp: &amqpwrap.RPCResponse{Code: RPCResponseCodeLockLost}}))
	require.True(t, isLockLostError(&amqp.Error{Condition: amqp.ErrCond("com.microsoft:message-lock-lost")}))

	require.True(t, isLockLostError(RPCError{Resp: &amqpwrap.RPCResponse{Code: RPCResponseCodeLockLost}}))
}

func Test_TransformError(t *testing.T) {
	var asExportedErr *exported.Error

	err := TransformError(RPCError{Resp: &amqpwrap.RPCResponse{Code: RPCResponseCodeLockLost}})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.CodeLockLost, asExportedErr.Code)

	err = TransformError(RPCError{Resp: &amqpwrap.RPCResponse{Code: http.StatusUnauthorized}})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.CodeUnauthorizedAccess, asExportedErr.Code)

	err = TransformError(&amqp.Error{Condition: amqp.ErrCondUnauthorizedAccess})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.CodeUnauthorizedAccess, asExportedErr.Code)

	// make sure we don't translate errors that are already usable, like Azure Identity failures.
	err = TransformError(&azidentity.AuthenticationFailedError{})
	afe := &azidentity.AuthenticationFailedError{}
	require.ErrorAs(t, err, &afe)

	// sanity check, an RPCError but it's not a azservicebus.Code type error.
	err = TransformError(RPCError{Resp: &amqpwrap.RPCResponse{Code: http.StatusNotFound}})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.CodeNotFound, asExportedErr.Code)

	err = TransformError(&amqp.Error{Condition: amqp.ErrCond("com.microsoft:message-lock-lost")})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.CodeLockLost, asExportedErr.Code)

	// sanity check, an RPCError but it's not a azservicebus.Code type error.
	err = TransformError(&amqp.Error{Condition: amqp.ErrCondNotFound})
	require.False(t, errors.As(err, &asExportedErr))

	err = TransformError(&amqp.LinkError{})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.CodeConnectionLost, asExportedErr.Code)

	err = TransformError(&amqp.ConnError{})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.CodeConnectionLost, asExportedErr.Code)

	// don't double wrap an already wrapped error
	alreadyWrappedErr := &exported.Error{Code: exported.CodeConnectionLost}
	err = TransformError(alreadyWrappedErr)
	require.Equal(t, alreadyWrappedErr, err)

	// and it's okay, for convenience, to pass a nil.
	require.Nil(t, TransformError(nil))
}

func amqpConditionTakesPrecedence(t *testing.T, cond amqp.ErrCond, kind RecoveryKind) {
	t.Run(fmt.Sprintf("code:%s", string(cond)), func(t *testing.T) {
		// NOTE: some combinations of *Error and condition are non-sensical. The important part
		// is the inner amqp.Error's Condition should take precedence over where it was returned.
		t.Run("amqp.Error", func(t *testing.T) {
			rk := GetRecoveryKind(&amqp.Error{Condition: cond})
			require.Equal(t, kind, rk)
		})

		t.Run("amqp.LinkError", func(t *testing.T) {
			rk := GetRecoveryKind(&amqp.LinkError{RemoteErr: &amqp.Error{Condition: cond}})
			require.Equal(t, kind, rk)
		})

		t.Run("amqp.SessionError", func(t *testing.T) {
			rk := GetRecoveryKind(&amqp.SessionError{RemoteErr: &amqp.Error{Condition: cond}})
			require.Equal(t, kind, rk)
		})

		t.Run("amqp.ConnError", func(t *testing.T) {
			rk := GetRecoveryKind(&amqp.ConnError{RemoteErr: &amqp.Error{Condition: cond}})
			require.Equal(t, kind, rk)
		})
	})
}
