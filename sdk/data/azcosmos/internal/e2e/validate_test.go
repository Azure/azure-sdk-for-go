// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func compiledTestNames(implementations []Implementation) map[string]bool {
	names := make(map[string]bool)
	for _, impl := range implementations {
		if impl.Test != "" {
			names[impl.Test] = true
		}
	}
	return names
}

func TestValidatePassesForTheRealImplementationMap(t *testing.T) {
	catalog, err := LoadCatalog()
	require.NoError(t, err)

	err = Validate(catalog, Implementations, compiledTestNames(Implementations))
	require.NoError(t, err, "every scenario in the pinned catalog must be active, blocked, or notApplicable with no drift")
}

func TestValidateFailsOnMissingScenario(t *testing.T) {
	catalog, err := LoadCatalog()
	require.NoError(t, err)

	incomplete := make([]Implementation, 0, len(Implementations)-1)
	for _, impl := range Implementations {
		if impl.ScenarioID == "item.lifecycle" {
			continue
		}
		incomplete = append(incomplete, impl)
	}

	err = Validate(catalog, incomplete, compiledTestNames(incomplete))
	require.Error(t, err)
	require.Contains(t, err.Error(), `"item.lifecycle" has no implementation-map entry`)
}

func TestValidateFailsOnUnknownScenarioID(t *testing.T) {
	catalog, err := LoadCatalog()
	require.NoError(t, err)

	impls := append([]Implementation{}, Implementations...)
	impls = append(impls, Implementation{ScenarioID: "item.does-not-exist", Status: StatusBlocked, Reason: "test"})

	err = Validate(catalog, impls, compiledTestNames(impls))
	require.Error(t, err)
	require.Contains(t, err.Error(), `unknown scenario id "item.does-not-exist"`)
}

func TestValidateFailsOnDuplicateScenarioID(t *testing.T) {
	catalog, err := LoadCatalog()
	require.NoError(t, err)

	impls := append([]Implementation{}, Implementations...)
	impls = append(impls, Implementation{ScenarioID: "item.lifecycle", Status: StatusBlocked, Reason: "test"})

	err = Validate(catalog, impls, compiledTestNames(impls))
	require.Error(t, err)
	require.Contains(t, err.Error(), `duplicate implementation-map entry for scenario "item.lifecycle"`)
}

func TestValidateFailsOnActiveEntryWithUncompiledTest(t *testing.T) {
	catalog, err := LoadCatalog()
	require.NoError(t, err)

	compiled := compiledTestNames(Implementations)
	delete(compiled, "TestItemLifecycleSmoke")

	err = Validate(catalog, Implementations, compiled)
	require.Error(t, err)
	require.Contains(t, err.Error(), `names Go test "TestItemLifecycleSmoke", which was not found among compiled tests`)
}

func TestValidateFailsOnActiveEntryMissingTestName(t *testing.T) {
	catalog, err := LoadCatalog()
	require.NoError(t, err)

	impls := append([]Implementation{}, Implementations...)
	for i, impl := range impls {
		if impl.ScenarioID == "item.lifecycle" {
			impls[i].Test = ""
		}
	}

	err = Validate(catalog, impls, compiledTestNames(impls))
	require.Error(t, err)
	require.Contains(t, err.Error(), `"item.lifecycle" is active but names no Go test function`)
}

func TestValidateFailsOnBlockedEntryMissingReason(t *testing.T) {
	catalog, err := LoadCatalog()
	require.NoError(t, err)

	impls := append([]Implementation{}, Implementations...)
	for i, impl := range impls {
		if impl.ScenarioID == "management.capabilities" {
			impls[i].Reason = ""
		}
	}

	err = Validate(catalog, impls, compiledTestNames(impls))
	require.Error(t, err)
	require.Contains(t, err.Error(), `"management.capabilities" has status "blocked" but gives no Reason`)
}

func TestCoverageAccountsForEveryScenarioExactlyOnce(t *testing.T) {
	catalog, err := LoadCatalog()
	require.NoError(t, err)
	require.NoError(t, Validate(catalog, Implementations, compiledTestNames(Implementations)))

	report := Coverage(catalog, Implementations)
	total := len(report.Active) + len(report.Blocked) + len(report.NotApplicable)
	require.Equal(t, len(catalog.Scenarios), total)
	require.NotEmpty(t, report.Active, "at least one scenario must be active after this implementation pass")
}
