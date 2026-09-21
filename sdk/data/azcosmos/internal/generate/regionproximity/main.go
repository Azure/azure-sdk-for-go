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
	"time"
)

const (
	maxRustSourceSize = 1 << 20
	rustRevision      = "294719d5da11318f5d49d090302cc61c461bf83a"
	rustSourceHash    = "1f067621a59ea59fdc8c60d48ba6673998204fca941820f74fd1876bb9947bdb"
	rustSourceURL     = "https://raw.githubusercontent.com/Azure/azure-sdk-for-rust/" + rustRevision +
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
	fmt.Fprintf(&output, "// Code generated from Azure/azure-sdk-for-rust region_proximity.rs at %s (SHA-256 %s); DO NOT EDIT.\n",
		rustRevision, rustSourceHash)
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
	client := http.Client{Timeout: 30 * time.Second}
	response, err := client.Get(rustSourceURL)
	if err != nil {
		panic(fmt.Errorf("downloading Rust proximity source: %w", err))
	}
	source, err := io.ReadAll(io.LimitReader(response.Body, maxRustSourceSize+1))
	closeErr := response.Body.Close()
	if err != nil {
		panic(fmt.Errorf("reading Rust proximity source: %w", err))
	}
	if closeErr != nil {
		panic(fmt.Errorf("closing Rust proximity source: %w", closeErr))
	}
	if response.StatusCode != http.StatusOK {
		panic(fmt.Errorf("downloading Rust proximity source: %s", response.Status))
	}
	if len(source) > maxRustSourceSize {
		panic(fmt.Errorf("rust proximity source exceeds %d bytes", maxRustSourceSize))
	}
	if actual := fmt.Sprintf("%x", sha256.Sum256(source)); actual != rustSourceHash {
		panic(fmt.Errorf("rust proximity source hash is %s, want %s", actual, rustSourceHash))
	}
	return source
}

func parseMatchArms(source []byte) map[string]string {
	arms := make(map[string]string)
	for _, match := range matchArmPattern.FindAllSubmatch(source, -1) {
		name := string(match[1])
		if _, exists := arms[name]; exists {
			panic(fmt.Errorf("source %q appears more than once", name))
		}
		arms[name] = string(match[2])
	}
	return arms
}

func parseArrays(source []byte) map[string][]string {
	arrays := make(map[string][]string)
	for _, match := range arrayPattern.FindAllSubmatch(source, -1) {
		name := string(match[1])
		if _, exists := arrays[name]; exists {
			panic(fmt.Errorf("array %q appears more than once", name))
		}
		regions := make([]string, 0, 96)
		for _, region := range regionPattern.FindAllSubmatch(match[2], -1) {
			regions = append(regions, strings.ToLower(strings.ReplaceAll(string(region[1]), "_", "")))
		}
		arrays[name] = regions
	}
	return arrays
}

func validate(matchArms map[string]string, arrays map[string][]string) {
	if len(matchArms) != 96 || len(arrays) != 96 {
		panic(fmt.Errorf("parsed %d match arms and %d arrays, want 96 of each", len(matchArms), len(arrays)))
	}
	usedArrays := make(map[string]struct{}, len(arrays))
	for source, arrayName := range matchArms {
		regions, ok := arrays[arrayName]
		if !ok {
			panic(fmt.Errorf("source %q references missing array %q", source, arrayName))
		}
		if _, used := usedArrays[arrayName]; used {
			panic(fmt.Errorf("array %q is referenced more than once", arrayName))
		}
		usedArrays[arrayName] = struct{}{}
		if len(regions) != 96 || regions[0] != source {
			panic(fmt.Errorf("source %q has an invalid proximity list", source))
		}
		seen := make(map[string]struct{}, len(regions))
		for _, region := range regions {
			if _, known := matchArms[region]; !known {
				panic(fmt.Errorf("source %q contains unknown region %q", source, region))
			}
			if _, ok := seen[region]; ok {
				panic(fmt.Errorf("source %q repeats region %q", source, region))
			}
			seen[region] = struct{}{}
		}
	}
	if len(usedArrays) != len(arrays) {
		panic(fmt.Errorf("%d proximity arrays are not referenced", len(arrays)-len(usedArrays)))
	}
}
