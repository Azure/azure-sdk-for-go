// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package e2e

import (
	"fmt"
	"sort"
	"strings"
)

// Validate checks the Go implementation map against the shared catalog and the Go test
// functions actually compiled for the selected build. It fails configuration (returns a
// non-nil, multi-issue error) rather than silently reducing coverage when it finds:
//
//   - a scenario ID with no implementation-map entry at all;
//   - an implementation-map entry naming an unknown scenario ID;
//   - a duplicate implementation-map entry for the same scenario ID;
//   - an active entry whose Test name is empty, or whose Test does not appear in compiledTests;
//   - an active entry that also sets Reason, or a non-active entry with an empty Reason;
//   - an active entry naming a Go test for a scenario whose declared profiles do not include any
//     profile actually present in the catalog (a mapping drift signal).
//
// compiledTests should be every Go test function name actually discovered by `go test -list` for
// the build under validation, not merely text that resembles a function declaration; see
// ../../e2e's TestImplementationMapMatchesCompiledTests for how that set is obtained.
func Validate(catalog Catalog, implementations []Implementation, compiledTests map[string]bool) error {
	var issues []string

	seen := make(map[string]bool, len(implementations))
	for _, impl := range implementations {
		if impl.ScenarioID == "" {
			issues = append(issues, "implementation entry has an empty ScenarioID")
			continue
		}
		if seen[impl.ScenarioID] {
			issues = append(issues, fmt.Sprintf("duplicate implementation-map entry for scenario %q", impl.ScenarioID))
			continue
		}
		seen[impl.ScenarioID] = true

		scenario, known := catalog.Scenarios[impl.ScenarioID]
		if !known {
			issues = append(issues, fmt.Sprintf("implementation-map entry references unknown scenario id %q", impl.ScenarioID))
			continue
		}

		switch impl.Status {
		case StatusActive:
			if impl.Reason != "" {
				issues = append(issues, fmt.Sprintf("scenario %q is active but also sets Reason %q; active entries are self-evidencing via Test", impl.ScenarioID, impl.Reason))
			}
			if impl.Test == "" {
				issues = append(issues, fmt.Sprintf("scenario %q is active but names no Go test function", impl.ScenarioID))
			} else if !compiledTests[impl.Test] {
				issues = append(issues, fmt.Sprintf("scenario %q names Go test %q, which was not found among compiled tests", impl.ScenarioID, impl.Test))
			}
			if !scenarioHasKnownProfile(scenario, catalog) {
				issues = append(issues, fmt.Sprintf("scenario %q declares no profile present in the pinned catalog", impl.ScenarioID))
			}
		case StatusBlocked, StatusNotApplicable:
			if impl.Test != "" {
				issues = append(issues, fmt.Sprintf("scenario %q has status %q but also names Go test %q; non-active entries must not claim a test", impl.ScenarioID, impl.Status, impl.Test))
			}
			if impl.Reason == "" {
				issues = append(issues, fmt.Sprintf("scenario %q has status %q but gives no Reason", impl.ScenarioID, impl.Status))
			}
		default:
			issues = append(issues, fmt.Sprintf("scenario %q has unknown status %q", impl.ScenarioID, impl.Status))
		}
	}

	for id := range catalog.Scenarios {
		if !seen[id] {
			issues = append(issues, fmt.Sprintf("scenario %q has no implementation-map entry (every scenario must be active, blocked, or notApplicable)", id))
		}
	}

	if len(issues) == 0 {
		return nil
	}
	sort.Strings(issues)
	return fmt.Errorf("e2e: implementation map validation failed:\n  - %s", strings.Join(issues, "\n  - "))
}

func scenarioHasKnownProfile(scenario Scenario, catalog Catalog) bool {
	for _, profileID := range scenario.Profiles {
		if _, ok := catalog.Profiles[profileID]; ok {
			return true
		}
	}
	return false
}
