// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestPriorityLevelValues(t *testing.T) {
	values := PriorityLevelValues()
	if len(values) != 2 {
		t.Fatalf("expected 2 priority levels, got %d", len(values))
	}
	if values[0] != PriorityLevelHigh {
		t.Errorf("expected first value to be High, got %v", values[0])
	}
	if values[1] != PriorityLevelLow {
		t.Errorf("expected second value to be Low, got %v", values[1])
	}
}

func TestPriorityLevelToPtr(t *testing.T) {
	ptr := PriorityLevelHigh.ToPtr()
	if ptr == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *ptr != PriorityLevelHigh {
		t.Errorf("expected High, got %v", *ptr)
	}
}
