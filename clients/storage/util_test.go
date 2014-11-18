package storage

import (
	"testing"
	"time"
)

func Test_timeRfc1123Formatted(t *testing.T) {
	now := time.Now().UTC()

	expectedLayout := "Mon, 02 Jan 2006 15:04:05 GMT"
	expected := now.Format(expectedLayout)

	if output := timeRfc1123Formatted(now); output != expected {
		t.Errorf("Expected: %s, got: %s", expected, output)
	}
}
