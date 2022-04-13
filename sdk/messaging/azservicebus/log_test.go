// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/stretchr/testify/require"
)

func TestEventConstantsDontDriftApart(t *testing.T) {
	require.Equal(t, EventConn, internal.EventConn)
	require.Equal(t, EventAuth, internal.EventAuth)
	require.Equal(t, EventReceiver, internal.EventReceiver)
	require.Equal(t, EventSender, internal.EventSender)
	require.Equal(t, EventAdmin, internal.EventAdmin)
}
