package changelog

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/exports"
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/report"
)

// NewChangelogForPackage returns the changelog for the given pkgDir.
// This function must be used when the new changes are generated, but not committed to git
// This function will first fetch all the exported content of the given package
// Then it add everything to git and then do a `git stash`, temporary revert the package to the previous state
// and fetch all the exported content of the given package again
// compare and generate a changelog report for the package
// This function will undo all the git changes before return
func NewChangelogForPackage(pkgDir string) (c *Changelog, err error) {
	// first we need to get the current status of the package
	rhs, err := getExportForPackage(pkgDir)
	if err != nil {
		return nil, err
	}
	// stash everything and get the previous status of the package
	if err := stashEverything(); err != nil {
		return nil, err
	}
	// reset everything when we are done
	defer func() {
		err = resetEverything()
	}()
	// get the original state of the package
	lhs, err := getExportForPackage(pkgDir)
	if err != nil {
		return nil, err
	}
	return getChangelogForPackage(pkgDir, lhs, rhs)
}

func stashEverything() error {
	if err := gitAddAll(); err != nil {
		return err
	}
	if err := gitStash(); err != nil {
		return err
	}
	return nil
}

func resetEverything() error {
	if err := gitStashPop(); err != nil {
		return err
	}
	if err := gitResetHead(); err != nil {
		return err
	}
	return nil
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
	p := report.Generate(*lhs, *rhs, false, false)
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
