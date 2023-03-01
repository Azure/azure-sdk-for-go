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

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
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
				rk := GetRecoveryKind(&amqp.Error{Condition: amqp.ErrorCondition(code)})
				require.EqualValues(t, RecoveryKindLink, rk, fmt.Sprintf("requires link recovery: %s", code))
			})
		}

		t.Run("sentinel errors", func(t *testing.T) {
			rk := GetRecoveryKind(amqp.ErrLinkClosed)
			require.EqualValues(t, RecoveryKindLink, rk)

			rk = GetRecoveryKind(amqp.ErrSessionClosed)
			require.EqualValues(t, RecoveryKindConn, rk)
		})
	})

	t.Run("connection", func(t *testing.T) {
		codes := []string{
			string(amqp.ErrorConnectionForced),
			string(amqp.ErrorInternalError),
		}

		for _, code := range codes {
			t.Run(code, func(t *testing.T) {
				rk := GetRecoveryKind(&amqp.Error{Condition: amqp.ErrorCondition(code)})
				require.EqualValues(t, RecoveryKindConn, rk, fmt.Sprintf("requires connection recovery: %s", code))
			})
		}

		t.Run("sentinel errors", func(t *testing.T) {
			rk := GetRecoveryKind(&amqp.ConnectionError{})
			require.EqualValues(t, RecoveryKindConn, rk)
		})
	})

	t.Run("nonretriable", func(t *testing.T) {
		codes := []string{
			string(amqp.ErrorUnauthorizedAccess),
			string(amqp.ErrorNotFound),
			string(amqp.ErrorMessageSizeExceeded),
		}

		for _, code := range codes {
			t.Run(code, func(t *testing.T) {
				rk := GetRecoveryKind(&amqp.Error{Condition: amqp.ErrorCondition(code)})
				require.EqualValues(t, RecoveryKindFatal, rk, fmt.Sprintf("cannot be recovered: %s", code))
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
				rk := GetRecoveryKind(&amqp.Error{Condition: amqp.ErrorCondition(code)})
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
		&amqp.Error{Condition: amqp.ErrorCondition("com.microsoft:server-busy")},
		&amqp.Error{Condition: amqp.ErrorCondition("com.microsoft:timeout")},
		&amqp.Error{Condition: amqp.ErrorCondition("com.microsoft:operation-cancelled")},
		errors.New("link is currently draining"), // not yet exposed from go-amqp
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
		&amqp.Error{Condition: amqp.ErrorConnectionForced},
		&amqp.Error{Condition: amqp.ErrorInternalError},
		&amqp.ConnectionError{},
		amqp.ErrSessionClosed,
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
}

func Test_ServiceBusError_LinkRecoveryNeeded(t *testing.T) {
	var linkErrors = []error{
		amqp.ErrLinkClosed,
		&amqp.DetachError{},
		&amqp.Error{Condition: amqp.ErrorDetachForced},
		&amqp.Error{Condition: amqp.ErrorTransferLimitExceeded},
		// this can happen when we're recovering the link - the client gets closed and the old link is still being
		// used by this instance of the client. It needs to recover and attempt it again.
		RPCError{Resp: &amqpwrap.RPCResponse{Code: 401}},
	}

	for i, err := range linkErrors {
		rk := GetRecoveryKind(err)
		require.EqualValues(t, RecoveryKindLink, rk, fmt.Sprintf("[%d] %v", i, err))
	}
}

func Test_ServiceBusError_Fatal(t *testing.T) {
	var fatalConditions = []amqp.ErrorCondition{
		amqp.ErrorMessageSizeExceeded,
		amqp.ErrorUnauthorizedAccess,
		amqp.ErrorNotFound,
		amqp.ErrorNotAllowed,
		amqp.ErrorCondition("com.microsoft:entity-disabled"),
		amqp.ErrorCondition("com.microsoft:session-cannot-be-locked"),
		amqp.ErrorCondition("com.microsoft:message-lock-lost"),
	}

	for i, cond := range fatalConditions {
		rk := GetRecoveryKind(&amqp.Error{Condition: cond})
		require.EqualValues(t, RecoveryKindFatal, rk, fmt.Sprintf("[%d] %s", i, cond))
	}

	require.Equal(t, RecoveryKindFatal, GetRecoveryKind(RPCError{Resp: &amqpwrap.RPCResponse{Code: http.StatusNotFound}}))
	require.Equal(t, RecoveryKindFatal, GetRecoveryKind(RPCError{Resp: &amqpwrap.RPCResponse{Code: RPCResponseCodeLockLost}}))
}

func Test_IsLockLostError(t *testing.T) {
	require.True(t, isLockLostError(RPCError{Resp: &amqpwrap.RPCResponse{Code: RPCResponseCodeLockLost}}))
	require.True(t, isLockLostError(&amqp.Error{Condition: amqp.ErrorCondition("com.microsoft:message-lock-lost")}))

	require.True(t, isLockLostError(RPCError{Resp: &amqpwrap.RPCResponse{Code: RPCResponseCodeLockLost}}))
}

func Test_TransformError(t *testing.T) {
	var asExportedErr *exported.Error

	err := TransformError(RPCError{Resp: &amqpwrap.RPCResponse{Code: RPCResponseCodeLockLost}})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.CodeLockLost, asExportedErr.Code)

	// sanity check, an RPCError but it's not a azservicebus.Code type error.
	err = TransformError(RPCError{Resp: &amqpwrap.RPCResponse{Code: http.StatusNotFound}})
	require.False(t, errors.As(err, &asExportedErr))

	err = TransformError(&amqp.Error{Condition: amqp.ErrorCondition("com.microsoft:message-lock-lost")})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.CodeLockLost, asExportedErr.Code)

	// sanity check, an RPCError but it's not a azservicebus.Code type error.
	err = TransformError(&amqp.Error{Condition: amqp.ErrorNotFound})
	require.False(t, errors.As(err, &asExportedErr))

	err = TransformError(amqp.ErrLinkClosed)
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.CodeConnectionLost, asExportedErr.Code)

	err = TransformError(&amqp.ConnectionError{})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.CodeConnectionLost, asExportedErr.Code)

	// don't double wrap an already wrapped error
	alreadyWrappedErr := &exported.Error{Code: exported.CodeConnectionLost}
	err = TransformError(alreadyWrappedErr)
	require.Equal(t, alreadyWrappedErr, err)

	// and it's okay, for convenience, to pass a nil.
	require.Nil(t, TransformError(nil))
}
