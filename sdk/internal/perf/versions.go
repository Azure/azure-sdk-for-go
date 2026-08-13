// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"fmt"
	"runtime"
	debugpkg "runtime/debug"
	"sort"
	"strings"
)

// printVersions writes a "=== Versions ===" block to stdout that lists the Go
// runtime version and the versions of every Azure SDK module compiled into the
// running perf binary. The format mirrors what perf-automation's other
// language adapters (e.g. Net.cs) parse from stdout to populate the
// IterationResult.PackageVersions dictionary:
//
//	=== Versions ===
//	go:                          1.22.3
//	github.com/Azure/azure-sdk-for-go/sdk/azcore:        Informational: v1.14.0
//	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob: Informational: v1.4.0
//
// The "Informational:" prefix matches the .NET adapter's
// `Informational: (\S*)` capture group so a future Go.cs adapter that
// reuses the .NET parsing pattern works without modification.
func printVersions() {
	fmt.Println("\n=== Versions ===")
	fmt.Printf("go:                          Informational: %s\n", runtime.Version())

	info, ok := debugpkg.ReadBuildInfo()
	if !ok {
		return
	}

	type modVersion struct {
		path    string
		version string
	}
	var mods []modVersion
	seen := map[string]struct{}{}
	add := func(path, version string) {
		if path == "" || version == "" {
			return
		}
		if _, dup := seen[path]; dup {
			return
		}
		seen[path] = struct{}{}
		mods = append(mods, modVersion{path: path, version: version})
	}

	if info.Main.Path != "" {
		add(info.Main.Path, info.Main.Version)
	}
	for _, dep := range info.Deps {
		if dep == nil {
			continue
		}
		// Prefer the resolved version when a replace directive is in effect.
		path := dep.Path
		version := dep.Version
		if dep.Replace != nil && dep.Replace.Version != "" {
			version = dep.Replace.Version
		}
		// Only emit Azure SDK and golang.org/x modules to keep the block
		// focused on shipping packages relevant to perf comparison.
		if !strings.HasPrefix(path, "github.com/Azure/") &&
			!strings.HasPrefix(path, "golang.org/x/") {
			continue
		}
		add(path, version)
	}

	sort.Slice(mods, func(i, j int) bool { return mods[i].path < mods[j].path })
	for _, m := range mods {
		fmt.Printf("%s: Informational: %s\n", m.path, m.version)
	}
}
