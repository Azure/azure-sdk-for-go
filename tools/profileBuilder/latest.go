package main

// Copyright 2017 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/marstr/collection"
)

// LatestStrategy evaluates the Azure-SDK-for-Go repository for the latest available API versions of each service.
type LatestStrategy struct {
	Root      string
	Predicate func(packageName string) bool
}

// AcceptAll is a predefined value for `LatestStrategy.Predicate` which always returns true.
func AcceptAll(name string) bool {
	return true
}

// IgnorePreview searches a packages "API Version" to see if it contains the word "preview". It only returns true when a package is not a preview.
func IgnorePreview(name string) (result bool) {
	matches := packageName.FindStringSubmatch(name)
	if len(matches) >= 4 { // Each group was captured.
		result = !strings.Contains(matches[2], "preview") // matches[2] is the `version` group
	}
	return
}

var packageName = regexp.MustCompile(`service[/\\](?P<provider>[\w\-\.\d_]+)[/\\](?:management[/\\])?(?P<version>[\d\-\w\._]+)[/\\](?P<group>[\w\d\-\._]+)`)

// Enumerate scans through the known Azure SDK for Go packages and finds each
func (latest LatestStrategy) Enumerate(cancel <-chan struct{}) collection.Enumerator {
	results := make(chan interface{})

	go func() {
		defer close(results)

		type operationGroup struct {
			provider     string
			resourceType string
			group        string
		}

		type operInfo struct {
			version string
			rawpath string
		}

		maxFound := make(map[operationGroup]operInfo)

		filepath.Walk(latest.Root, func(currentPath string, info os.FileInfo, openErr error) (err error) {
			pathMatches := packageName.FindStringSubmatch(currentPath)

			if latest.Predicate == nil {
				latest.Predicate = AcceptAll
			}

			if len(pathMatches) == 0 || !info.IsDir() || !latest.Predicate(currentPath) {
				return
			}

			version := pathMatches[2]
			currentGroup := operationGroup{
				provider:     pathMatches[1],
				resourceType: pathMatches[2],
				group:        pathMatches[3],
			}

			prev, ok := maxFound[currentGroup]
			if !ok {
				maxFound[currentGroup] = operInfo{version, currentPath}
				return
			}

			if le, _ := versionle(prev.version, version); le {
				maxFound[currentGroup] = operInfo{version, currentPath}
			}

			return
		})

		for _, entry := range maxFound {
			absolute, err := filepath.Abs(entry.rawpath)
			if err != nil {
				continue
			}

			select {
			case results <- absolute:
				// Intionally Left Blank
			case <-cancel:
				return
			}
		}
	}()

	return results
}

// ErrNotVersionString is instantiated when a string not conforming to known Azure API Version patterns is parsed is if it did.
type ErrNotVersionString string

func (err ErrNotVersionString) Error() string {
	return fmt.Sprintf("`%s` is not a recognized Azure version string", string(err))
}

var versionle = func() func(string, string) (bool, error) {
	versionPattern := regexp.MustCompile(`^(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})(?:[\.\-](?P<tag>.+))?$`)

	return func(left, right string) (result bool, err error) {
		leftMatch := versionPattern.FindStringSubmatch(left)
		rightMatch := versionPattern.FindStringSubmatch(right)

		if len(leftMatch) < 3 { // No match found
			err = ErrNotVersionString(left)
			return
		}

		if len(rightMatch) < 3 { // No match found
			err = ErrNotVersionString(right)
			return
		}

		for i := 1; i <= 3; i++ { // Start with index 1 because the first element is then entire match, not just a group. End at 3 because there are three numeric groups.
			if leftMatch[i] == rightMatch[i] {
				continue
			}

			var leftNum, rightNum int
			leftNum, err = strconv.Atoi(leftMatch[i])
			if err != nil {
				return
			}

			rightNum, err = strconv.Atoi(rightMatch[i])
			if err != nil {
				return
			}

			if leftNum < rightNum {
				result = true
				return
			}

			result = false // This line only here for readability
			return
		}

		if leftTag, rightTag := leftMatch[4], rightMatch[4]; leftTag == "" && rightTag != "" { // match[4] is the tag portion of a date based API Version label
			result = false
		} else if leftTag != "" && rightTag != "" {
			result = leftTag <= rightTag
		} else {
			result = true
		}
		return
	}
}()
