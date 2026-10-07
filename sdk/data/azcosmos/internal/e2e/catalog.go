// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package e2e ports the shared Cosmos SDK functional E2E certification catalog (scenario
// metadata, setup profiles, and an implementation map) into Go. The catalog data under
// catalog/ is a pinned, vendored snapshot of Azure/azure-sdk-for-rust's sdk/cosmos/e2e_tests;
// see catalog/README.md for the exact pinned revision and refresh procedure.
//
// This package owns only the language-neutral contract: loading and validating the catalog,
// and matching it against Go's own implementation map. It intentionally knows nothing about
// how to run a scenario; that lives in the sibling ../../e2e test package, which imports only
// the public azcosmos/v2 API.
package e2e

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"sort"
)

//go:embed catalog/schema catalog/profiles catalog/scenarios
var catalogFS embed.FS

// BackendApplicability mirrors the shared scenario schema's backends.*.applicability enum.
type BackendApplicability string

const (
	ApplicabilityRequired      BackendApplicability = "required"
	ApplicabilitySupported     BackendApplicability = "supported"
	ApplicabilitySimulated     BackendApplicability = "simulated"
	ApplicabilityNotApplicable BackendApplicability = "notApplicable"
)

// Backend records one backend's applicability, fidelity, and required capabilities for a
// scenario, mirroring the shared scenario.v1.json schema's "backend" definition.
type Backend struct {
	Applicability BackendApplicability `json:"applicability"`
	Fidelity      string               `json:"fidelity"`
	Reason        string               `json:"reason,omitempty"`
	Requires      []string             `json:"requires,omitempty"`
}

// Scenario is one permanent, language-neutral shared scenario definition.
type Scenario struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Maturity string   `json:"maturity"`
	Profiles []string `json:"profiles"`
	Tags     []string `json:"tags"`
	Backends struct {
		HostedEmulatorGatewayV1 Backend `json:"hostedEmulatorGatewayV1"`
		HostedEmulatorGatewayV2 Backend `json:"hostedEmulatorGatewayV2"`
		AzureLive               Backend `json:"azureLive"`
	} `json:"backends"`
}

// Profile is one reusable account/runtime/client setup definition.
type Profile struct {
	ID       string           `json:"id"`
	Accounts []map[string]any `json:"accounts"`
	Runtimes []map[string]any `json:"runtimes"`
	Clients  []map[string]any `json:"clients"`
}

// Catalog is the loaded, parsed set of shared scenario and profile documents.
type Catalog struct {
	Scenarios map[string]Scenario // keyed by Scenario.ID
	Profiles  map[string]Profile  // keyed by Profile.ID
}

// SortedScenarioIDs returns every scenario ID in the catalog, sorted for stable output.
func (c Catalog) SortedScenarioIDs() []string {
	ids := make([]string, 0, len(c.Scenarios))
	for id := range c.Scenarios {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// LoadCatalog parses every vendored scenario and profile document. It fails on malformed JSON,
// a missing required field from the shared schema's "required" list, or a duplicate ID; it does
// not perform full JSON Schema validation (no schema-validation dependency is justified for a
// vendored, already-schema-checked-upstream snapshot), but it does check every structural
// invariant this package's validation relies on.
func LoadCatalog() (Catalog, error) {
	catalog := Catalog{
		Scenarios: map[string]Scenario{},
		Profiles:  map[string]Profile{},
	}

	profileFiles, err := catalogFS.ReadDir("catalog/profiles")
	if err != nil {
		return Catalog{}, fmt.Errorf("e2e: reading catalog/profiles: %w", err)
	}
	for _, entry := range profileFiles {
		if entry.IsDir() {
			continue
		}
		p, err := loadProfile(path.Join("catalog/profiles", entry.Name()))
		if err != nil {
			return Catalog{}, err
		}
		if _, exists := catalog.Profiles[p.ID]; exists {
			return Catalog{}, fmt.Errorf("e2e: duplicate profile id %q", p.ID)
		}
		catalog.Profiles[p.ID] = p
	}

	var scenarioFiles []string
	if err := walkJSON(catalogFS, "catalog/scenarios", &scenarioFiles); err != nil {
		return Catalog{}, err
	}
	for _, file := range scenarioFiles {
		s, err := loadScenario(file)
		if err != nil {
			return Catalog{}, err
		}
		if _, exists := catalog.Scenarios[s.ID]; exists {
			return Catalog{}, fmt.Errorf("e2e: duplicate scenario id %q (file %s)", s.ID, file)
		}
		catalog.Scenarios[s.ID] = s
	}

	return catalog, nil
}

func loadProfile(file string) (Profile, error) {
	data, err := catalogFS.ReadFile(file)
	if err != nil {
		return Profile{}, fmt.Errorf("e2e: reading %s: %w", file, err)
	}
	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return Profile{}, fmt.Errorf("e2e: parsing %s: %w", file, err)
	}
	if p.ID == "" {
		return Profile{}, fmt.Errorf("e2e: %s is missing required field \"id\"", file)
	}
	if len(p.Accounts) == 0 || len(p.Runtimes) == 0 || len(p.Clients) == 0 {
		return Profile{}, fmt.Errorf("e2e: profile %q (%s) must declare at least one account, runtime, and client", p.ID, file)
	}
	return p, nil
}

func loadScenario(file string) (Scenario, error) {
	data, err := catalogFS.ReadFile(file)
	if err != nil {
		return Scenario{}, fmt.Errorf("e2e: reading %s: %w", file, err)
	}
	var s Scenario
	if err := json.Unmarshal(data, &s); err != nil {
		return Scenario{}, fmt.Errorf("e2e: parsing %s: %w", file, err)
	}
	switch {
	case s.ID == "":
		return Scenario{}, fmt.Errorf("e2e: %s is missing required field \"id\"", file)
	case s.Title == "":
		return Scenario{}, fmt.Errorf("e2e: scenario %q (%s) is missing required field \"title\"", s.ID, file)
	case s.Maturity == "":
		return Scenario{}, fmt.Errorf("e2e: scenario %q (%s) is missing required field \"maturity\"", s.ID, file)
	case len(s.Profiles) == 0:
		return Scenario{}, fmt.Errorf("e2e: scenario %q (%s) must declare at least one profile", s.ID, file)
	case len(s.Tags) == 0:
		return Scenario{}, fmt.Errorf("e2e: scenario %q (%s) must declare at least one tag", s.ID, file)
	}
	for _, backend := range []Backend{s.Backends.HostedEmulatorGatewayV1, s.Backends.HostedEmulatorGatewayV2, s.Backends.AzureLive} {
		if backend.Applicability == "" || backend.Fidelity == "" {
			return Scenario{}, fmt.Errorf("e2e: scenario %q (%s) has a backend missing applicability/fidelity", s.ID, file)
		}
	}
	return s, nil
}

func walkJSON(fsys embed.FS, dir string, out *[]string) error {
	entries, err := fsys.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("e2e: reading %s: %w", dir, err)
	}
	for _, entry := range entries {
		full := path.Join(dir, entry.Name())
		if entry.IsDir() {
			if err := walkJSON(fsys, full, out); err != nil {
				return err
			}
			continue
		}
		if path.Ext(entry.Name()) == ".json" {
			*out = append(*out, full)
		}
	}
	return nil
}
