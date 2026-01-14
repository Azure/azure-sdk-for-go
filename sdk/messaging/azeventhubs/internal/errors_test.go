// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/exported"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestOwnershipLost(t *testing.T) {
	detachErr := &amqp.LinkError{
		RemoteErr: &amqp.Error{
			Condition: amqp.ErrCond("amqp:link:stolen"),
		},
	}

	require.Equal(t, RecoveryKindFatal, GetRecoveryKind(detachErr))
	require.False(t, IsQuickRecoveryError(detachErr))

	transformedErr := TransformError(detachErr)

	var err *exported.Error
	require.ErrorAs(t, transformedErr, &err)
	require.Equal(t, exported.ErrorCodeOwnershipLost, err.Code)

	require.False(t, IsOwnershipLostError(&amqp.LinkError{}))
	require.False(t, IsOwnershipLostError(&amqp.ConnError{}))
	require.False(t, IsOwnershipLostError(errors.New("definitely not an ownership lost error")))
}

func TestIsQuickRecoveryError(t *testing.T) {
	t.Run("LinkError returns true", func(t *testing.T) {
		require.True(t, IsQuickRecoveryError(&amqp.LinkError{}))
	})

	t.Run("LinkError with remote error returns true", func(t *testing.T) {
		require.True(t, IsQuickRecoveryError(&amqp.LinkError{
			RemoteErr: &amqp.Error{
				Condition: amqp.ErrCondDetachForced,
			},
		}))
	})

	t.Run("ownership lost error returns false", func(t *testing.T) {
		require.False(t, IsQuickRecoveryError(&amqp.LinkError{
			RemoteErr: &amqp.Error{
				Condition: amqp.ErrCond("amqp:link:stolen"),
			},
		}))
	})

	t.Run("ConnError returns false", func(t *testing.T) {
		require.False(t, IsQuickRecoveryError(&amqp.ConnError{}))
	})

	t.Run("generic error returns false", func(t *testing.T) {
		require.False(t, IsQuickRecoveryError(errors.New("some error")))
	})

	t.Run("nil error returns false", func(t *testing.T) {
		require.False(t, IsQuickRecoveryError(nil))
	})

	t.Run("wrapped LinkError returns true", func(t *testing.T) {
		wrappedErr := fmt.Errorf("wrapped: %w", &amqp.LinkError{})
		require.True(t, IsQuickRecoveryError(wrappedErr))
	})
}

func TestGetRecoveryKind(t *testing.T) {
	require.Equal(t, GetRecoveryKind(nil), RecoveryKindNone)
	require.Equal(t, GetRecoveryKind(amqpwrap.ErrConnResetNeeded), RecoveryKindConn)
	require.Equal(t, GetRecoveryKind(&amqp.LinkError{}), RecoveryKindLink)
	require.Equal(t, GetRecoveryKind(RPCLinkClosedErr), RecoveryKindFatal)
	require.Equal(t, GetRecoveryKind(context.Canceled), RecoveryKindFatal)
	require.Equal(t, GetRecoveryKind(&amqp.Error{Condition: amqp.ErrCondResourceLimitExceeded}), RecoveryKindFatal)

	// fatal RPC errors
	for _, code := range []int{http.StatusUnauthorized, http.StatusNotFound} {
		t.Run(fmt.Sprintf("RPCError.Code==%d is fatal", code), func(t *testing.T) {
			actual := GetRecoveryKind(RPCError{Resp: &amqpwrap.RPCResponse{Code: code}})
			require.Equal(t, RecoveryKindFatal, actual)
		})
	}

	// recoverable RPC errors
	for _, code := range []int{http.StatusRequestTimeout, http.StatusServiceUnavailable, http.StatusInternalServerError} {
		t.Run(fmt.Sprintf("RPCError.Code==%d is retriable", code), func(t *testing.T) {
			actual := GetRecoveryKind(RPCError{Resp: &amqpwrap.RPCResponse{Code: code}})
			require.Equal(t, RecoveryKindNone, actual)
		})
	}
}

func TestIsNotAllowedError(t *testing.T) {
	require.True(t, IsNotAllowedError(&amqp.Error{
		Condition: amqp.ErrCondNotAllowed,
	}))

	require.False(t, IsNotAllowedError(&amqp.Error{
		Condition: amqp.ErrCondConnectionForced,
	}))

	require.False(t, IsNotAllowedError(errors.New("hello")))
}

func Test_TransformError(t *testing.T) {
	var asExportedErr *exported.Error

	err := TransformError(&amqp.LinkError{
		RemoteErr: &amqp.Error{
			Condition: amqp.ErrCond("amqp:link:stolen"),
		},
	})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.ErrorCodeOwnershipLost, asExportedErr.Code)

	err = TransformError(RPCError{Resp: &amqpwrap.RPCResponse{Code: http.StatusUnauthorized}})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.ErrorCodeUnauthorizedAccess, asExportedErr.Code)

	err = TransformError(&amqp.Error{Condition: amqp.ErrCondUnauthorizedAccess})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.ErrorCodeUnauthorizedAccess, asExportedErr.Code)

	// make sure we don't translate errors that are already usable, like Azure Identity failures.
	err = TransformError(&azidentity.AuthenticationFailedError{})
	afe := &azidentity.AuthenticationFailedError{}
	require.ErrorAs(t, err, &afe)

	// sanity check, an RPCError but it's not a azservicebus.Code type error.
	err = TransformError(RPCError{Resp: &amqpwrap.RPCResponse{Code: http.StatusNotFound}})
	require.False(t, errors.As(err, &asExportedErr))

	// sanity check, an RPCError but it's not a azservicebus.Code type error.
	err = TransformError(&amqp.Error{Condition: amqp.ErrCondNotFound})
	require.False(t, errors.As(err, &asExportedErr))

	err = TransformError(&amqp.LinkError{})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.ErrorCodeConnectionLost, asExportedErr.Code)

	err = TransformError(&amqp.ConnError{})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.ErrorCodeConnectionLost, asExportedErr.Code)

	// don't double wrap an already wrapped error
	alreadyWrappedErr := &exported.Error{Code: exported.ErrorCodeConnectionLost}
	err = TransformError(alreadyWrappedErr)
	require.Equal(t, alreadyWrappedErr, err)

	// and it's okay, for convenience, to pass a nil.
	require.Nil(t, TransformError(nil))
}
