// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package e2e

import (
	"fmt"
	"sort"
	"strings"
)

// CoverageReport is an exact covered/blocked/not-applicable matrix with reasons, suitable for
// printing in CI output or a PR description. It never omits a scenario: every ID present in the
// catalog appears in exactly one of Active, Blocked, or NotApplicable.
type CoverageReport struct {
	Active        []Implementation
	Blocked       []Implementation
	NotApplicable []Implementation
}

// Coverage builds a CoverageReport from a catalog and an already-Validate'd implementation map.
// Call Validate first; Coverage does not re-check consistency.
func Coverage(catalog Catalog, implementations []Implementation) CoverageReport {
	var report CoverageReport
	for _, impl := range implementations {
		switch impl.Status {
		case StatusActive:
			report.Active = append(report.Active, impl)
		case StatusBlocked:
			report.Blocked = append(report.Blocked, impl)
		case StatusNotApplicable:
			report.NotApplicable = append(report.NotApplicable, impl)
		}
	}
	byID := func(list []Implementation) func(i, j int) bool {
		return func(i, j int) bool { return list[i].ScenarioID < list[j].ScenarioID }
	}
	sort.Slice(report.Active, byID(report.Active))
	sort.Slice(report.Blocked, byID(report.Blocked))
	sort.Slice(report.NotApplicable, byID(report.NotApplicable))
	return report
}

// String renders the report as a human-readable matrix.
func (r CoverageReport) String() string {
	var b strings.Builder
	section := func(title string, entries []Implementation, detail func(Implementation) string) {
		fmt.Fprintf(&b, "%s (%d):\n", title, len(entries))
		for _, impl := range entries {
			fmt.Fprintf(&b, "  - %s: %s\n", impl.ScenarioID, detail(impl))
		}
	}
	section("Active", r.Active, func(i Implementation) string { return i.Test })
	section("Blocked", r.Blocked, func(i Implementation) string { return i.Reason })
	section("Not applicable", r.NotApplicable, func(i Implementation) string { return i.Reason })
	return b.String()
}
