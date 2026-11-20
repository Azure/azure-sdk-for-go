// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package loader loads Go modules for API report generation.
package loader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

// Result contains a module's metadata and loaded Go packages.
type Result struct {
	// Dir is the absolute path to the module directory.
	Dir string

	// ModFile is the module's parsed go.mod file.
	ModFile *modfile.File

	// Packages contains the module's non-test packages in import-path order,
	// including internal packages. Dependency packages are reachable through Imports.
	Packages []*packages.Package
}

// Load loads the module rooted at dir using the current Go environment and build
// constraints. Module packages and their dependencies retain syntax and type
// information. Tests and nested modules are excluded.
// Load returns an error if module metadata or any loaded package is invalid.
func Load(ctx context.Context, dir string) (*Result, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolve module directory: %w", err)
	}

	modPath := filepath.Join(dir, "go.mod")
	content, err := os.ReadFile(modPath)
	if err != nil {
		return nil, fmt.Errorf("read module file: %w", err)
	}
	modFile, err := modfile.Parse(modPath, content, nil)
	if err != nil {
		return nil, fmt.Errorf("parse module file: %w", err)
	}
	if modFile.Module == nil {
		return nil, fmt.Errorf("module directive is missing from %q", modPath)
	}

	loaded, err := packages.Load(&packages.Config{
		Context: ctx,
		Dir:     dir,
		Mode:    packages.LoadAllSyntax | packages.NeedModule,
		Tests:   false,
	}, "./...")
	if err != nil {
		return nil, fmt.Errorf("load packages for module %q: %w", modFile.Module.Mod.Path, err)
	}

	var loadErrors []error
	packages.Visit(loaded, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			loadErrors = append(loadErrors, err)
		}
	})
	if err := errors.Join(loadErrors...); err != nil {
		return nil, fmt.Errorf("load packages for module %q: %w", modFile.Module.Mod.Path, err)
	}

	var modulePackages []*packages.Package
	for _, pkg := range loaded {
		if pkg.Module != nil && pkg.Module.Path == modFile.Module.Mod.Path {
			modulePackages = append(modulePackages, pkg)
		}
	}
	if len(modulePackages) == 0 {
		return nil, fmt.Errorf("module %q contains no non-test Go packages", modFile.Module.Mod.Path)
	}
	slices.SortFunc(modulePackages, func(a, b *packages.Package) int {
		return strings.Compare(a.PkgPath, b.PkgPath)
	})

	return &Result{
		Dir:      dir,
		ModFile:  modFile,
		Packages: modulePackages,
	}, nil
}
