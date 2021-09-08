// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFormatTimes(t *testing.T) {
	start := time.Date(2021, time.September, 8, 13, 6, 0, 0, time.UTC)
	expiry := start.AddDate(1, 0, 0)
	startString, expiryString := FormatTimesForSASSigning(start, expiry)
	require.Equal(t, "2021-09-08T13:06:00Z", startString)
	require.Equal(t, "2022-09-08T13:06:00Z", expiryString)
}

func TestFormatIPRange(t *testing.T) {
	i := IPRange{
		Start: net.IPv4(224, 0, 0, 250),
	}
	require.Equal(t, i.String(), "224.0.0.250")

	i2 := IPRange{
		End: net.IPv4(192, 0, 0, 168),
	}
	require.Equal(t, i2.String(), "")

	i3 := IPRange{
		Start: net.IPv4(192, 0, 0, 168),
		End: net.IPv4(224, 0, 0, 250),
	}
	require.Equal(t, i3.String(), "192.0.0.168-224.0.0.250")
}
