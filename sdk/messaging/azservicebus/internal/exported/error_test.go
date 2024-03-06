// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"testing"

	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestBasicError(t *testing.T) {
	err := NewError(CodeConnectionLost, nil)
	require.Equal(t, err.Error(), "(connlost): unknown error")

	err = NewError(CodeLockLost, &amqp.Error{Condition: amqp.ErrCond("com.microsoft:message-lock-lost")})
	require.Equal(t, err.Error(), "(locklost): *Error{Condition: com.microsoft:message-lock-lost, Description: , Info: map[]}")
}
