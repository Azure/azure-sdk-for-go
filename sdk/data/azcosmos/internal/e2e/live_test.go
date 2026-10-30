// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLiveSelection(t *testing.T) {
	catalog, err := LoadCatalog()
	require.NoError(t, err)
	require.NoError(t, ValidateLiveSelection(catalog, LiveScenarios))
	for _, tc := range []struct {
		name string
		edit func(*Catalog, *[]LiveScenario)
	}{
		{"empty", func(_ *Catalog, s *[]LiveScenario) { *s = nil }},
		{"duplicate", func(_ *Catalog, s *[]LiveScenario) { *s = append(*s, (*s)[0]) }},
		{"unknown", func(_ *Catalog, s *[]LiveScenario) { (*s)[0].ID = "unknown" }},
		{"wrong profile", func(_ *Catalog, s *[]LiveScenario) { (*s)[0].Profile = "coreOperations" }},
		{"not applicable", func(c *Catalog, _ *[]LiveScenario) {
			s := c.Scenarios["item.lifecycle"]
			s.Backends.AzureLive.Applicability = ApplicabilityNotApplicable
			c.Scenarios[s.ID] = s
		}},
		{"unsupported capability", func(c *Catalog, _ *[]LiveScenario) {
			s := c.Scenarios["item.lifecycle"]
			s.Backends.AzureLive.Requires = []string{"partitionSplit"}
			c.Scenarios[s.ID] = s
		}},
		{"simulated fidelity", func(c *Catalog, _ *[]LiveScenario) {
			s := c.Scenarios["item.lifecycle"]
			s.Backends.AzureLive.Fidelity = "simulated"
			c.Scenarios[s.ID] = s
		}},
		{"unknown profile", func(c *Catalog, _ *[]LiveScenario) {
			delete(c.Profiles, "smokeTests")
		}},
		{"missing required scenario", func(c *Catalog, _ *[]LiveScenario) {
			s := c.Scenarios["item.create-conflict"]
			s.Backends.AzureLive.Applicability = ApplicabilityRequired
			c.Scenarios[s.ID] = s
		}},
		{"profile drift", func(c *Catalog, _ *[]LiveScenario) {
			p := c.Profiles["smokeTests"]
			p.Accounts[0]["consistency"] = "strong"
			c.Profiles[p.ID] = p
		}},
		{"non-default runtime", func(c *Catalog, _ *[]LiveScenario) {
			p := c.Profiles["smokeTests"]
			p.Runtimes[0]["gatewayV2"] = "enabled"
			c.Profiles[p.ID] = p
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			catalog, err := LoadCatalog()
			require.NoError(t, err)
			selection := append([]LiveScenario(nil), LiveScenarios...)
			tc.edit(&catalog, &selection)
			require.Error(t, ValidateLiveSelection(catalog, selection))
		})
	}
}
