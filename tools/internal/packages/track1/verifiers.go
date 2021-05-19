// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package track1

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/Azure/azure-sdk-for-go/tools/internal/packages"
)

type verifiers []packages.VerifyFunc

func (v verifiers) Verify(pkg packages.Package) []error {
	var errors []error
	for _, verifier := range v {
		if err := verifier(pkg); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// GetDefaultVerifier gets the default track 1 verifiers which verifies the following:
// - verifies that the package name matches the leaf directory name
// - verifies that the package path only consists of lower case letters
// - verifies that the package path must follow the desired directory structure
func GetDefaultVerifier() packages.Verifier {
	return verifiers([]packages.VerifyFunc{
		verifyPkgMatchesDir,
		verifyLowerCase,
		verifyDirectoryStructure,
	})
}

// ensures that the leaf directory name matches the package name
// new to modules: if the last leaf is version suffix, find its parent as leaf folder
func verifyPkgMatchesDir(p packages.Package) error {
	leaf := findPackageFolderInPath(p.FullPath())
	if !strings.EqualFold(leaf, p.Name()) {
		return fmt.Errorf("leaf directory of '%s' doesn't match package name '%s'", p.Path(), p.Name())
	}
	return nil
}

func findPackageFolderInPath(path string) string {
	regex := regexp.MustCompile(`/v\d+$`)
	if regex.MatchString(path) {
		// folder path ends with version suffix
		path = path[:strings.LastIndex(path, "/")]
	}
	result := path[strings.LastIndex(path, "/")+1:]
	return result
}

// ensures that there are no upper-case letters in a package's directory
func verifyLowerCase(p packages.Package) error {
	// walk the package directory looking for upper-case characters
	for _, r := range p.Path() {
		if r == '/' {
			continue
		}
		if unicode.IsUpper(r) {
			return fmt.Errorf("found upper-case character in directory '%s'", p.Path())
		}
	}
	return nil
}

// ensures that the package's directory hierarchy is properly formed
func verifyDirectoryStructure(p packages.Package) error {
	// for ARM the package directory structure is highly deterministic:
	// /redis/mgmt/2015-08-01/redis
	// /resources/mgmt/2017-06-01-preview/policy
	// /preview/signalr/mgmt/2018-03-01-preview/signalr
	// /preview/security/mgmt/v2.0/security (version scheme for composite packages)
	// /network/mgmt/2019-10-01/network/v2 (new with modules)
	if !p.IsARMPackage() {
		return nil
	}
	regexStr := fmt.Sprintf(`^(preview/)?%s/mgmt/%s/%s%s$`, rpNameRegex, apiVersionRegex, nameSpaceRegex, majorSubDirRegex)
	regex := regexp.MustCompile(regexStr)
	if !regex.MatchString(p.Path()) {
		return fmt.Errorf("bad directory structure '%s'", p.Path())
	}
	return nil
}

const (
	rpNameRegex      = `[a-z0-9\-]+`
	apiVersionRegex  = `((?:\d{4}-\d{2}-\d{2}(?:-preview)?)|(?:v\d{1,2}\.\d{1,2}))`
	nameSpaceRegex   = `[a-z0-9]+`
	majorSubDirRegex = `(/v\d+)?`
)
