// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestOwnershipLost(t *testing.T) {
	detachErr := &amqp.DetachError{
		RemoteError: &amqp.Error{
			Condition: amqp.ErrorCondition("amqp:link:stolen"),
		},
	}

	require.Equal(t, RecoveryKindFatal, GetRecoveryKind(detachErr))
	require.False(t, IsQuickRecoveryError(detachErr))

	transformedErr := TransformError(detachErr)

	var err *exported.Error
	require.ErrorAs(t, transformedErr, &err)
	require.Equal(t, exported.ErrorCodeOwnershipLost, err.Code)

	require.False(t, IsOwnershipLostError(&amqp.DetachError{}))
	require.False(t, IsOwnershipLostError(&amqp.ConnectionError{}))
	require.False(t, IsOwnershipLostError(errors.New("definitely not an ownership lost error")))
}

func TestGetRecoveryKind(t *testing.T) {
	require.Equal(t, GetRecoveryKind(nil), RecoveryKindNone)
	require.Equal(t, GetRecoveryKind(errConnResetNeeded), RecoveryKindConn)
	require.Equal(t, GetRecoveryKind(&amqp.DetachError{}), RecoveryKindLink)
	require.Equal(t, GetRecoveryKind(context.Canceled), RecoveryKindFatal)
	require.Equal(t, GetRecoveryKind(RPCError{Resp: &amqpwrap.RPCResponse{Code: http.StatusUnauthorized}}), RecoveryKindFatal)
	require.Equal(t, GetRecoveryKind(RPCError{Resp: &amqpwrap.RPCResponse{Code: http.StatusNotFound}}), RecoveryKindFatal)
}

func Test_TransformError(t *testing.T) {
	var asExportedErr *exported.Error

	err := TransformError(&amqp.DetachError{
		RemoteError: &amqp.Error{
			Condition: amqp.ErrorCondition("amqp:link:stolen"),
		},
	})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.ErrorCodeOwnershipLost, asExportedErr.Code)

	err = TransformError(RPCError{Resp: &amqpwrap.RPCResponse{Code: http.StatusUnauthorized}})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.ErrorCodeUnauthorizedAccess, asExportedErr.Code)

	err = TransformError(&amqp.Error{Condition: amqp.ErrorUnauthorizedAccess})
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
	err = TransformError(&amqp.Error{Condition: amqp.ErrorNotFound})
	require.False(t, errors.As(err, &asExportedErr))

	err = TransformError(amqp.ErrLinkClosed)
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.ErrorCodeConnectionLost, asExportedErr.Code)

	err = TransformError(&amqp.ConnectionError{})
	require.ErrorAs(t, err, &asExportedErr)
	require.Equal(t, exported.ErrorCodeConnectionLost, asExportedErr.Code)

	// don't double wrap an already wrapped error
	alreadyWrappedErr := &exported.Error{Code: exported.ErrorCodeConnectionLost}
	err = TransformError(alreadyWrappedErr)
	require.Equal(t, alreadyWrappedErr, err)

	// and it's okay, for convenience, to pass a nil.
	require.Nil(t, TransformError(nil))
}
