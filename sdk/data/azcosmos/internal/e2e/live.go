// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore ppcb

package e2e

import (
	"fmt"
	"slices"
)

// LiveScenario identifies a bounded live check and its shared setup profile.
type LiveScenario struct {
	ID      string
	Profile string
}

// LiveScenarios is the initial live smoke selection, not certification of every
// scenario in either profile. Pagination belongs to coreOperations in the shared catalog.
var LiveScenarios = []LiveScenario{
	{ID: "bootstrap.primary-success", Profile: "smokeTests"},
	{ID: "item.lifecycle", Profile: "smokeTests"},
	{ID: "item.optimistic-concurrency", Profile: "smokeTests"},
	{ID: "diagnostics.success-and-error", Profile: "smokeTests"},
	{ID: "query.pagination-resume", Profile: "coreOperations"},
}

// ValidateLiveSelection rejects catalog drift and unsupported requirements before network I/O.
func ValidateLiveSelection(catalog Catalog, selection []LiveScenario) error {
	if len(selection) == 0 {
		return fmt.Errorf("e2e: live selection is empty")
	}
	selected := make(map[string]bool)
	profiles := make(map[string]bool)
	for _, entry := range selection {
		if selected[entry.ID] {
			return fmt.Errorf("e2e: duplicate live scenario %q", entry.ID)
		}
		scenario, ok := catalog.Scenarios[entry.ID]
		if !ok || !slices.Contains(scenario.Profiles, entry.Profile) {
			return fmt.Errorf("e2e: live scenario %q does not declare profile %q", entry.ID, entry.Profile)
		}
		profile, ok := catalog.Profiles[entry.Profile]
		if !ok || !liveProfileSupported(profile) {
			return fmt.Errorf("e2e: live profile %q requires unsupported setup", entry.Profile)
		}
		backend := scenario.Backends.AzureLive
		if (backend.Applicability != ApplicabilityRequired && backend.Applicability != ApplicabilitySupported) ||
			backend.Fidelity != "full" {
			return fmt.Errorf("e2e: scenario %q cannot run with full Azure Live fidelity", entry.ID)
		}
		for _, capability := range backend.Requires {
			if capability != "item" && capability != "query" {
				return fmt.Errorf("e2e: scenario %q requires unsupported live capability %q", entry.ID, capability)
			}
		}
		selected[entry.ID] = true
		profiles[entry.Profile] = true
	}
	for _, scenario := range catalog.Scenarios {
		if scenario.Backends.AzureLive.Applicability != ApplicabilityRequired {
			continue
		}
		for _, profile := range scenario.Profiles {
			if profiles[profile] && !selected[scenario.ID] {
				return fmt.Errorf("e2e: required live scenario %q is not implemented in the selected lane", scenario.ID)
			}
		}
	}
	return nil
}

func liveProfileSupported(profile Profile) bool {
	if len(profile.Accounts) != 1 || len(profile.Runtimes) != 1 || len(profile.Clients) != 1 {
		return false
	}
	account, runtime, client := profile.Accounts[0], profile.Runtimes[0], profile.Clients[0]
	regions, ok := account["regions"].([]any)
	if !ok || len(regions) != 1 {
		return false
	}
	region, ok := regions[0].(map[string]any)
	replication, replicationOK := account["replication"].(map[string]any)
	return ok && region["name"] == "East US" &&
		replicationOK && replication["minDelayMs"] == float64(0) && replication["maxDelayMs"] == float64(0) &&
		account["id"] == "sessionSingleRegion" && account["consistency"] == "session" &&
		account["writeMode"] == "single" && account["perPartitionFailover"] == false &&
		runtime["id"] == "sdkDefault" && runtime["gatewayV2"] == "backendDefault" &&
		runtime["ppcb"] == "sdkDefault" && client["id"] == "sdkDefault" &&
		client["binaryEncoding"] == "sdkDefault" && client["routing"] == "proximity"
}
