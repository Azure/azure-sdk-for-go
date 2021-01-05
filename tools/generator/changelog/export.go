package changelog

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/exports"
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/report"
)

// Exporter ...
type Exporter struct {
	// SDKRoot ...
	SDKRoot string
	// BackupRoot ...
	BackupRoot string
}

// ExportForPackage generates the Changelog for the given relPkgDir
// relPkgDir is the package path relative to the root of SDK, e.g. services/compute/mgmt/2020-06-01/compute
func (p Exporter) ExportForPackage(relPkgDir string) (*Changelog, error) {
	lhs, err := getExportForPackage(filepath.Join(p.BackupRoot, relPkgDir))
	if err != nil {
		return nil, err
	}
	rhs, err := getExportForPackage(filepath.Join(p.SDKRoot, relPkgDir))
	if err != nil {
		return nil, err
	}
	return getChangelogForPackage(relPkgDir, lhs, rhs)
}

func getChangelogForPackage(pkgDir string, lhs, rhs *exports.Content) (*Changelog, error) {
	if lhs == nil && rhs == nil {
		return nil, fmt.Errorf("this package does not exist even after the generation, this should never happen")
	}
	if lhs == nil {
		// the package does not exist before the generation: this is a new package
		return &Changelog{
			PackageName: pkgDir,
			NewPackage:  true,
		}, nil
	}
	if rhs == nil {
		// the package no longer exists after the generation: this package was removed
		return &Changelog{
			PackageName:    pkgDir,
			RemovedPackage: true,
		}, nil
	}
	// lhs and rhs are both non-nil
	p := report.Generate(*lhs, *rhs, nil)
	return &Changelog{
		PackageName: pkgDir,
		Modified:    &p,
	}, nil
}

func getExportForPackage(pkgDir string) (*exports.Content, error) {
	// The function exports.Get does not handle the circumstance that the package does not exist
	// therefore we have to check if it exists and if not exit early to ensure we do not return an error
	if _, err := os.Stat(pkgDir); os.IsNotExist(err) {
		return nil, nil
	}
	exp, err := exports.Get(pkgDir)
	if err != nil {
		return nil, err
	}
	return &exp, nil
}
