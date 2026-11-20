// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package report builds normalized models for API review reports.
package report

import (
	"slices"
	"strings"

	"apigen/internal/loader"
)

// Module contains module metadata and direct, non-test dependencies.
type Module struct {
	// Path is the module's import path.
	Path string

	// GoVersion is the version from the go directive, or empty if it is absent.
	GoVersion string

	// Dependencies contains external modules imported directly by non-test code,
	// including internal packages, sorted by module path.
	Dependencies []Dependency
}

// Dependency identifies an external module imported directly by module code.
type Dependency struct {
	// Path is the dependency's module path.
	Path string

	// Version is the selected module version, or empty for a workspace module.
	Version string
}

// ExtractModule extracts module metadata and direct, non-test dependencies.
// Standard library imports, same-module imports, and transitive-only dependencies
// are excluded. loaded must be a successful result from loader.Load.
func ExtractModule(loaded *loader.Result) Module {
	module := Module{Path: loaded.ModFile.Module.Mod.Path}
	if loaded.ModFile.Go != nil {
		module.GoVersion = loaded.ModFile.Go.Version
	}

	dependencies := make(map[string]Dependency)
	for _, pkg := range loaded.Packages {
		for _, imported := range pkg.Imports {
			dependency := imported.Module
			if dependency == nil || dependency.Path == module.Path {
				continue
			}
			dependencies[dependency.Path] = Dependency{
				Path:    dependency.Path,
				Version: dependency.Version,
			}
		}
	}
	for _, dependency := range dependencies {
		module.Dependencies = append(module.Dependencies, dependency)
	}
	slices.SortFunc(module.Dependencies, func(a, b Dependency) int {
		return strings.Compare(a.Path, b.Path)
	})
	return module
}
