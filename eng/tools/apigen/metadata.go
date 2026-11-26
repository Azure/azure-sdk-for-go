// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go/constant"
	"go/types"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"apigen/internal/loader"

	"golang.org/x/mod/semver"
)

const apigenVersion = "0.1.0"

func reviewMetadata(loaded *loader.Result, repositoryRoot string, document []byte) ([]byte, error) {
	version, err := moduleVersion(loaded)
	if err != nil {
		return nil, err
	}
	relativePath, err := packageRelativePath(loaded.Dir, repositoryRoot)
	if err != nil {
		return nil, err
	}
	// JSON string quoting is valid YAML and preserves unusual directory names.
	quotedPath, err := json.Marshal(relativePath)
	if err != nil {
		return nil, fmt.Errorf("encode package relative path: %w", err)
	}
	return fmt.Appendf(nil, "apiHash: %s\npackageVersion: %s\npackageRelativePath: %s\napigenVersion: %s\ngoVersion: %s\n",
		apiHash(document), version, quotedPath, apigenVersion, strings.TrimPrefix(runtime.Version(), "go")), nil
}

func apiHash(document []byte) string {
	sum := sha256.Sum256(document)
	return hex.EncodeToString(sum[:])
}

func moduleVersion(loaded *loader.Result) (string, error) {
	var version, source string
	for _, pkg := range loaded.Packages {
		for _, name := range pkg.Types.Scope().Names() {
			switch strings.ToLower(name) {
			case "moduleversion", "servicelibversion", "version":
			default:
				continue
			}
			declaration, ok := pkg.Types.Scope().Lookup(name).(*types.Const)
			if !ok {
				continue
			}
			filename := pkg.Fset.PositionFor(declaration.Pos(), false).Filename
			base := filepath.Base(filename)
			if !strings.Contains(base, "constant") && !strings.Contains(base, "version") {
				continue
			}
			if declaration.Val().Kind() != constant.String {
				return "", fmt.Errorf("module version constant %s in %q must be a string", name, filename)
			}
			value := constant.StringVal(declaration.Val())
			if !strings.HasPrefix(value, "v") {
				return "", fmt.Errorf("module version in %q should be %q, not %q", filename, "v"+value, value)
			}
			if !semver.IsValid(value) || !strings.HasPrefix(value, semver.MajorMinor(value)+".") {
				return "", fmt.Errorf("module version %q in %q must have the form vMAJOR.MINOR.PATCH with an optional prerelease or build suffix", value, filename)
			}
			if version != "" && version != value {
				return "", fmt.Errorf("conflicting module versions %q in %q and %q in %q", version, source, value, filename)
			}
			version, source = value, filename
		}
	}
	if version == "" {
		return "", fmt.Errorf("no module version constant found in %q: expected moduleVersion, serviceLibVersion, or version in a constant or version Go file", loaded.Dir)
	}
	return strings.TrimPrefix(version, "v"), nil
}

func packageRelativePath(moduleDirectory, repositoryRoot string) (string, error) {
	moduleDirectory, err := filepath.EvalSymlinks(moduleDirectory)
	if err != nil {
		return "", fmt.Errorf("resolve module directory: %w", err)
	}
	moduleDirectory, err = filepath.Abs(moduleDirectory)
	if err != nil {
		return "", fmt.Errorf("resolve module directory: %w", err)
	}
	if repositoryRoot == "" {
		for directory := moduleDirectory; ; directory = filepath.Dir(directory) {
			_, err := os.Stat(filepath.Join(directory, ".git"))
			if err == nil {
				repositoryRoot = directory
				break
			}
			if !errors.Is(err, os.ErrNotExist) {
				return "", fmt.Errorf("inspect Git root at %q: %w", directory, err)
			}
			if filepath.Dir(directory) == directory {
				return "", fmt.Errorf("cannot find Git root for %q; specify --repo-root", moduleDirectory)
			}
		}
	}
	repositoryRoot, err = filepath.Abs(repositoryRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}
	repositoryRoot, err = filepath.EvalSymlinks(repositoryRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}
	info, err := os.Stat(repositoryRoot)
	if err != nil {
		return "", fmt.Errorf("inspect repository root: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("repository root %q is not a directory", repositoryRoot)
	}
	relative, err := filepath.Rel(repositoryRoot, moduleDirectory)
	if err != nil {
		return "", fmt.Errorf("make module path relative to repository root: %w", err)
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("module directory %q is outside repository root %q", moduleDirectory, repositoryRoot)
	}
	return filepath.ToSlash(relative), nil
}
