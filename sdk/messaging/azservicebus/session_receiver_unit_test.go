// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalculateAcceptNextSessionTimeout(t *testing.T) {
	const maxTimeoutMilliseconds = uint32(1<<32 - 1)

	tests := []struct {
		name           string
		timeout        time.Duration
		jitterFraction float64
		expected       uint32
		ok             bool
	}{
		{
			name:    "expired",
			timeout: -time.Second,
		},
		{
			name:    "less than one millisecond",
			timeout: time.Microsecond,
		},
		{
			name:           "jitter reduces timeout below one millisecond",
			timeout:        time.Millisecond,
			jitterFraction: 1,
		},
		{
			name:     "no jitter",
			timeout:  90 * time.Second,
			expected: 90000,
			ok:       true,
		},
		{
			name:           "jitter capped at one hundred milliseconds",
			timeout:        90 * time.Second,
			jitterFraction: 1,
			expected:       89900,
			ok:             true,
		},
		{
			name:     "milliseconds capped at uint32",
			timeout:  (time.Duration(maxTimeoutMilliseconds) * time.Millisecond) + time.Hour,
			expected: maxTimeoutMilliseconds,
			ok:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, ok := calculateAcceptNextSessionTimeout(test.timeout, test.jitterFraction)
			require.Equal(t, test.ok, ok)
			require.Equal(t, test.expected, actual)
		})
	}
}
