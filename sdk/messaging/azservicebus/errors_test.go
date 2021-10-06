// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/stretchr/testify/require"
)

func TestErrClosed(t *testing.T) {
	var err error = errClosed{link: "hello"}

	_, ok := err.(errorinfo.NonRetriable)
	require.True(t, ok, "ErrClosed is a errorinfo.NonRetriable")
	require.EqualValues(t, "hello is closed and can no longer be used", err.Error())
}
