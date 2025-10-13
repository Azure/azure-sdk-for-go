package azcosmos

import (
	"fmt"
	"sync"
	"testing"
)

// mockEffectivePartitionKeyComputer is a mock implementation of the CosmosEffectivePartitionKey interface
// similar in spirit to the mock query engine used in query tests. It lets us verify that higher level
// code (once added) can depend on the interface without needing a real hashing engine.
type mockEffectivePartitionKeyComputer struct {
	mu     sync.Mutex
	Calls  []mockEPKCall
	Result string
	Err    error
	// Optional map of JSON input -> result to simulate different outputs
	PerInput map[string]string
}

type mockEPKCall struct {
	JSON    string
	Version int
	Kind    PartitionKeyKind
}

func (m *mockEffectivePartitionKeyComputer) ComputeEffectivePartitionKey(partitionKeyJSON string, version int, kind PartitionKeyKind) (string, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, mockEPKCall{JSON: partitionKeyJSON, Version: version, Kind: kind})
	m.mu.Unlock()

	if m.Err != nil {
		return "", m.Err
	}
	if m.PerInput != nil {
		if r, ok := m.PerInput[partitionKeyJSON]; ok {
			return r, nil
		}
	}
	if m.Result != "" {
		return m.Result, nil
	}
	return fmt.Sprintf("stub_%s_v%d", kind, version), nil
}

var _ CosmosEffectivePartitionKey = (*mockEffectivePartitionKeyComputer)(nil)

// TestMockEffectivePartitionKeyComputer_Basic ensures that the mock records calls and returns the expected result.
func TestMockEffectivePartitionKeyComputer_Basic(t *testing.T) {
	mock := &mockEffectivePartitionKeyComputer{Result: "epk_stub_1234"}
	res, err := mock.ComputeEffectivePartitionKey("\"abc\"", 2, PartitionKeyKindHash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "epk_stub_1234" {
		t.Fatalf("expected result epk_stub_1234, got %s", res)
	}
	if len(mock.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(mock.Calls))
	}
	call := mock.Calls[0]
	if call.JSON != "\"abc\"" || call.Version != 2 || call.Kind != PartitionKeyKindHash {
		t.Fatalf("unexpected call record: %+v", call)
	}
}

// TestMockEffectivePartitionKeyComputer_PerInput ensures per-input overrides work.
func TestMockEffectivePartitionKeyComputer_PerInput(t *testing.T) {
	mock := &mockEffectivePartitionKeyComputer{PerInput: map[string]string{
		"[\"a\",1]":    "epk_a_1",
		"\"Infinity\"": "epk_infinity",
	}}

	cases := []struct {
		json   string
		expect string
	}{
		{"[\"a\",1]", "epk_a_1"},
		{"\"Infinity\"", "epk_infinity"},
		{"\"other\"", "stub_Hash_v1"}, // falls back to default stub format
	}

	for _, c := range cases {
		got, err := mock.ComputeEffectivePartitionKey(c.json, 1, PartitionKeyKindHash)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", c.json, err)
		}
		if got != c.expect {
			t.Fatalf("for %s expected %s, got %s", c.json, c.expect, got)
		}
	}

	if len(mock.Calls) != len(cases) {
		t.Fatalf("expected %d calls, got %d", len(cases), len(mock.Calls))
	}
}

// TestMockEffectivePartitionKeyComputer_Error ensures configured error is surfaced.
func TestMockEffectivePartitionKeyComputer_Error(t *testing.T) {
	injectedErr := fmt.Errorf("boom")
	mock := &mockEffectivePartitionKeyComputer{Err: injectedErr}
	_, err := mock.ComputeEffectivePartitionKey("\"abc\"", 1, PartitionKeyKindHash)
	if err == nil || err.Error() != injectedErr.Error() {
		t.Fatalf("expected error %v, got %v", injectedErr, err)
	}
}

// TestMockEffectivePartitionKeyComputer_Concurrency ensures the mock is safe for concurrent use.
func TestMockEffectivePartitionKeyComputer_Concurrency(t *testing.T) {
	mock := &mockEffectivePartitionKeyComputer{}
	const workers = 16
	const callsPer = 50
	done := make(chan struct{})
	for w := 0; w < workers; w++ {
		go func(id int) {
			for i := 0; i < callsPer; i++ {
				_, _ = mock.ComputeEffectivePartitionKey(fmt.Sprintf("\"val-%d-%d\"", id, i), 2, PartitionKeyKindHash)
			}
			done <- struct{}{}
		}(w)
	}
	for w := 0; w < workers; w++ {
		<-done
	}
	expected := workers * callsPer
	if len(mock.Calls) != expected {
		t.Fatalf("expected %d calls recorded, got %d", expected, len(mock.Calls))
	}
}
