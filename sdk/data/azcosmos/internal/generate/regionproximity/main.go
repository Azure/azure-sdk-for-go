// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Command regionproximity generates the Go region-proximity table from the pinned Rust SDK source.
package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"go/format"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

const (
	rustRevision   = "294719d5da11318f5d49d090302cc61c461bf83a"
	rustSourceHash = "1f067621a59ea59fdc8c60d48ba6673998204fca941820f74fd1876bb9947bdb"
	rustSourceURL  = "https://raw.githubusercontent.com/Azure/azure-sdk-for-rust/" + rustRevision +
		"/sdk/cosmos/azure_data_cosmos/src/region_proximity.rs"
)

var (
	matchArmPattern = regexp.MustCompile(`"([^"]+)" => Some\(&([A-Z0-9_]+)_REGIONS\)`)
	arrayPattern    = regexp.MustCompile(`(?s)static ([A-Z0-9_]+)_REGIONS: \[Region; 96\] = \[(.*?)\n\];`)
	regionPattern   = regexp.MustCompile(`Region::([A-Z0-9_]+),`)
)

func main() {
	source := downloadSource()
	matchArms := parseMatchArms(source)
	arrays := parseArrays(source)
	validate(matchArms, arrays)

	var output bytes.Buffer
	fmt.Fprintln(&output, "// Copyright (c) Microsoft Corporation. All rights reserved.")
	fmt.Fprintln(&output, "// Licensed under the MIT License.")
	fmt.Fprintln(&output)
	fmt.Fprintf(&output, "// Code generated from Azure/azure-sdk-for-rust region_proximity.rs at %s; DO NOT EDIT.\n", rustRevision)
	fmt.Fprintln(&output)
	fmt.Fprintln(&output, "package azcosmos")
	fmt.Fprintln(&output)
	fmt.Fprintln(&output, "var proximityRegionOrderBySource = map[Region][]Region{")

	sources := make([]string, 0, len(matchArms))
	for source := range matchArms {
		sources = append(sources, source)
	}
	sort.Strings(sources)
	for _, source := range sources {
		fmt.Fprintf(&output, "\t%q: {\n", source)
		for _, region := range arrays[matchArms[source]] {
			fmt.Fprintf(&output, "\t\t%q,\n", region)
		}
		fmt.Fprintln(&output, "\t},")
	}
	fmt.Fprintln(&output, "}")

	formatted, err := format.Source(output.Bytes())
	if err != nil {
		panic(fmt.Errorf("formatting generated source: %w", err))
	}
	if err := os.WriteFile("region_proximity.go", formatted, 0o644); err != nil {
		panic(fmt.Errorf("writing generated source: %w", err))
	}
}

func downloadSource() []byte {
	response, err := http.Get(rustSourceURL)
	if err != nil {
		panic(fmt.Errorf("downloading Rust proximity source: %w", err))
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		panic(fmt.Errorf("downloading Rust proximity source: %s", response.Status))
	}
	source, err := io.ReadAll(response.Body)
	if err != nil {
		panic(fmt.Errorf("reading Rust proximity source: %w", err))
	}
	if actual := fmt.Sprintf("%x", sha256.Sum256(source)); actual != rustSourceHash {
		panic(fmt.Errorf("Rust proximity source hash is %s, want %s", actual, rustSourceHash))
	}
	return source
}

func parseMatchArms(source []byte) map[string]string {
	arms := make(map[string]string)
	for _, match := range matchArmPattern.FindAllSubmatch(source, -1) {
		arms[string(match[1])] = string(match[2])
	}
	return arms
}

func parseArrays(source []byte) map[string][]string {
	arrays := make(map[string][]string)
	for _, match := range arrayPattern.FindAllSubmatch(source, -1) {
		regions := make([]string, 0, 96)
		for _, region := range regionPattern.FindAllSubmatch(match[2], -1) {
			regions = append(regions, strings.ToLower(strings.ReplaceAll(string(region[1]), "_", "")))
		}
		arrays[string(match[1])] = regions
	}
	return arrays
}

func validate(matchArms map[string]string, arrays map[string][]string) {
	if len(matchArms) != 96 || len(arrays) != 96 {
		panic(fmt.Errorf("parsed %d match arms and %d arrays, want 96 of each", len(matchArms), len(arrays)))
	}
	for source, arrayName := range matchArms {
		regions, ok := arrays[arrayName]
		if !ok {
			panic(fmt.Errorf("source %q references missing array %q", source, arrayName))
		}
		if len(regions) != 96 || regions[0] != source {
			panic(fmt.Errorf("source %q has an invalid proximity list", source))
		}
		seen := make(map[string]struct{}, len(regions))
		for _, region := range regions {
			if _, ok := seen[region]; ok {
				panic(fmt.Errorf("source %q repeats region %q", source, region))
			}
			seen[region] = struct{}{}
		}
	}
}
