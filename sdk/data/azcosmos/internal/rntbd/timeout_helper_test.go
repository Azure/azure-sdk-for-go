// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestTimeoutHelper_IsElapsed verifies that IsElapsed returns false immediately
// after creation and true after the timeout duration passes.
// Port of Java: TimeoutHelperTest.isElapsed()
func TestTimeoutHelper_IsElapsed(t *testing.T) {
	helper1 := NewTimeoutHelper(100 * time.Millisecond)
	require.False(t, helper1.IsElapsed(), "should not be elapsed immediately after creation")

	helper2 := NewTimeoutHelper(100 * time.Millisecond)
	time.Sleep(100 * time.Millisecond)
	require.True(t, helper2.IsElapsed(), "should be elapsed after timeout duration")
}

// TestTimeoutHelper_GetRemainingTime verifies that GetRemainingTime decreases
// as time passes and is approximately timeout - elapsed.
// Port of Java: TimeoutHelperTest.getRemainingTime()
func TestTimeoutHelper_GetRemainingTime(t *testing.T) {
	const bufferMs = 15

	for i := 1; i <= 5; i++ {
		helper := NewTimeoutHelper(100 * time.Millisecond)
		sleepMs := 10 * i
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)

		remaining := helper.GetRemainingTime()
		expectedMaxMs := int64(100 - sleepMs + bufferMs)

		require.LessOrEqual(t, remaining.Milliseconds(), expectedMaxMs,
			"remaining time should be <= %dms (100ms - %dms sleep + %dms buffer), got %dms",
			expectedMaxMs, sleepMs, bufferMs, remaining.Milliseconds())
	}
}

// TestTimeoutHelper_GetElapsedTime verifies elapsed time increases over time.
func TestTimeoutHelper_GetElapsedTime(t *testing.T) {
	helper := NewTimeoutHelper(1 * time.Second)

	elapsed1 := helper.GetElapsedTime()
	require.Less(t, elapsed1, 50*time.Millisecond, "elapsed should be minimal initially")

	time.Sleep(50 * time.Millisecond)

	elapsed2 := helper.GetElapsedTime()
	require.GreaterOrEqual(t, elapsed2, 50*time.Millisecond, "elapsed should increase after sleep")
}

// TestTimeoutHelper_ThrowTimeoutIfElapsed verifies the method returns nil when
// not elapsed and returns a RequestTimeoutError when elapsed.
func TestTimeoutHelper_ThrowTimeoutIfElapsed(t *testing.T) {
	helper1 := NewTimeoutHelper(100 * time.Millisecond)
	err := helper1.ThrowTimeoutIfElapsed()
	require.NoError(t, err, "should not return error when not elapsed")

	helper2 := NewTimeoutHelper(50 * time.Millisecond)
	time.Sleep(60 * time.Millisecond)
	err = helper2.ThrowTimeoutIfElapsed()
	require.Error(t, err, "should return error when elapsed")

	var timeoutErr *RequestTimeoutError
	require.ErrorAs(t, err, &timeoutErr, "error should be RequestTimeoutError")
	require.Equal(t, int32(StatusRequestTimeout), timeoutErr.StatusCode)
}

// TestTimeoutHelper_ThrowGoneIfElapsed verifies the method returns nil when
// not elapsed and returns a GoneError when elapsed.
func TestTimeoutHelper_ThrowGoneIfElapsed(t *testing.T) {
	helper1 := NewTimeoutHelper(100 * time.Millisecond)
	err := helper1.ThrowGoneIfElapsed()
	require.NoError(t, err, "should not return error when not elapsed")

	helper2 := NewTimeoutHelper(50 * time.Millisecond)
	time.Sleep(60 * time.Millisecond)
	err = helper2.ThrowGoneIfElapsed()
	require.Error(t, err, "should return error when elapsed")

	var goneErr *GoneError
	require.ErrorAs(t, err, &goneErr, "error should be GoneError")
	require.Equal(t, int32(StatusGone), goneErr.StatusCode)
}

// TestTimeoutHelper_NegativeRemainingTime verifies that GetRemainingTime can
// return negative durations when the timeout has been exceeded.
func TestTimeoutHelper_NegativeRemainingTime(t *testing.T) {
	helper := NewTimeoutHelper(20 * time.Millisecond)
	time.Sleep(50 * time.Millisecond)

	remaining := helper.GetRemainingTime()
	require.Less(t, remaining, time.Duration(0), "remaining should be negative when elapsed")
}

// TestTimeoutHelper_ZeroTimeout verifies behavior with zero timeout.
func TestTimeoutHelper_ZeroTimeout(t *testing.T) {
	helper := NewTimeoutHelper(0)

	require.True(t, helper.IsElapsed(), "zero timeout should be immediately elapsed")
	require.LessOrEqual(t, helper.GetRemainingTime(), time.Duration(0))
}
