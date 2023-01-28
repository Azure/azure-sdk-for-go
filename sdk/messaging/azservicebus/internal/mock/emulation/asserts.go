// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func RequireNoLeaks(t *testing.T, events *Events) {
	links := events.GetOpenLinks()
	require.Empty(t, links, "No leaked links")

	conns := events.GetOpenConns()
	require.Empty(t, conns, "No leaked connections")
}
