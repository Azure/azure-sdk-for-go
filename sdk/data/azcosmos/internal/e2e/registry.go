// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package e2e

// ImplementationStatus records whether a shared scenario has an active Go test, is blocked by a
// missing capability, or is genuinely not applicable to this SDK/backend combination.
type ImplementationStatus string

const (
	// StatusActive means Test names a compiled Go test function that exercises this scenario.
	StatusActive ImplementationStatus = "active"
	// StatusBlocked means the scenario is in scope but cannot be implemented yet; Reason must
	// explain the missing Go API or emulator capability. Blocked scenarios must never be
	// silently dropped from the registry - their absence from Go coverage must stay visible.
	StatusBlocked ImplementationStatus = "blocked"
	// StatusNotApplicable means the scenario is permanently out of scope for Go (for example, it
	// is a Rust-internals precedent with no Go-observable counterpart). Reason must explain why.
	StatusNotApplicable ImplementationStatus = "notApplicable"
)

// Implementation maps one shared scenario ID to its Go implementation state. This is the
// versioned Go equivalent of Rust's implementations/rust.json: a single, compiled, type-checked
// source of truth instead of a second JSON file to keep in sync by hand.
type Implementation struct {
	ScenarioID string
	Status     ImplementationStatus
	// Test is the exact Go test function name (for example "TestItemLifecycleSmoke") that
	// implements this scenario. Required when Status is StatusActive; must be empty otherwise.
	Test string
	// Reason explains why a scenario is blocked or not applicable. Required when Status is not
	// StatusActive; must be empty for StatusActive (the Test name is the evidence).
	Reason string
}

// Implementations is the Go implementation map. Every shared scenario ID must appear here
// exactly once, whether or not Go has implemented it, so coverage gaps stay visible instead of
// silently shrinking the registry to only what happens to be implemented.
var Implementations = buildImplementations()
