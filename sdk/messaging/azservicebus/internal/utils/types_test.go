// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestISO8601StringToDuration(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		iso      string
		duration time.Duration
	}{
		{iso: "PT1S", duration: time.Second},
		{iso: "PT45S", duration: 45 * time.Second},
		{iso: "PT6H", duration: 6 * time.Hour},
		{iso: "PT21600S", duration: 21600 * time.Second},
		{iso: "PT4H", duration: 4 * time.Hour},
		{iso: "PT10M", duration: 10 * time.Minute},
	}

	for _, tc := range testCases {
		t.Run(tc.iso, func(t *testing.T) {
			actualDuration, err := ISO8601StringToDuration(&tc.iso)
			require.NoError(t, err)
			require.EqualValues(t, *actualDuration, tc.duration)
		})
	}
}
