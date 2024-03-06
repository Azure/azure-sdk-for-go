// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

const pkgPathRegex = `services[/\\](?P<prv>[\w\-\._/\\]+)[/\\](?P<ver>v?\d{4}-\d{2}-\d{2}[\w\.\-]*|v?\d+[\.\w\-]*)[/\\](?P<grp>[/\\\w\-\._]+)`

var (
	validPackagePath   = regexp.MustCompile(pkgPathRegex)
	validPkgPathModVer = regexp.MustCompile(pkgPathRegex + `[/\\](?P<mod>v\d+)[/\\]?(?P<api>\w+api)?`)
	apiPkgRegex        = regexp.MustCompile(`[/\\]\w+api$`)
)

// PathInfo provides information about a package's path.
type PathInfo struct {
	IsArm    bool
	Provider string
	Version  string
	Group    string
	ModVer   string
	APIPkg   string
}

// DeconstructPath takes a full package path and deconstructs it into its constituent parts.
func DeconstructPath(path string) (PathInfo, error) {
	path = filepath.Clean(path)
	// must check the module regex first, this is due to greedy regex and lack of negative look-aheads
	matches := validPkgPathModVer.FindAllStringSubmatch(path, -1)
	if matches == nil {
		matches = validPackagePath.FindAllStringSubmatch(path, -1)
	}
	if matches == nil {
		return PathInfo{}, fmt.Errorf("path '%s' does not resemble a known package path", path)
	}

	// matches[0][0] is the full match, we don't care about that
	// matches[0][1] is <prv>, it might end with /mgmt due to greedy regex
	// matches[0][2] is <ver>
	// matches[0][3] is <grp>, it might end with /*api for non-module case due to greedy regex
	prv := matches[0][1]
	ver := matches[0][2]
	grp := matches[0][3]
	mod := ""
	api := ""

	if len(matches[0]) == 6 {
		// this package contains a major version
		// matches[0][4] is <mod>
		// matches[0][5] is <api>
		mod = matches[0][4]
		api = matches[0][5]
	} else if loc := apiPkgRegex.FindStringIndex(grp); loc != nil {
		// for non-module case, if this is the *api package strip it off grp and place in api
		api = grp[loc[0]+1:]
		grp = grp[:loc[0]]
	}

	arm := false
	if index := strings.LastIndex(prv, string(filepath.Separator)+armPathModifier); index > 0 {
		prv = prv[:index]
		arm = true
	}
	return PathInfo{
		Provider: prv,
		IsArm:    arm,
		Version:  ver,
		Group:    grp,
		ModVer:   mod,
		APIPkg:   api,
	}, nil
}
