// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"testing"

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
}
