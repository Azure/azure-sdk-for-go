package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/marstr/collection"
)

type LatestStrategy struct {
	Root      string
	predicate func(packageName string) bool
}

var packageName = regexp.MustCompile(`service/(?P<provider>[\w\-\.\d]+)/(?P<type>[\w\-\.\d]+)/(?P<version>[\d\-\w\.]+)/(?P<group>[\w\d\-\.]+)`)

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

		maxFound := make(map[operationGroup]string)

		filepath.Walk(latest.Root, func(currentPath string, info os.FileInfo, openErr error) (err error) {
			pathMatches := packageName.FindStringSubmatch(currentPath)
			if len(pathMatches) == 0 {
				return
			}

			version := pathMatches[3]
			currentGroup := operationGroup{
				provider:     pathMatches[1],
				resourceType: pathMatches[2],
				group:        pathMatches[4],
			}

			prev, ok := maxFound[currentGroup]
			if !ok {
				maxFound[currentGroup] = version
				return
			}

			if le, _ := versionle(prev, version); le {
				maxFound[currentGroup] = version
			}

			return
		})

		// TODO: regurgitate the packages that were found to be the most recent versions.
	}()

	return results
}

// ErrNotVersionString is instantiated when a string not conforming to known Azure API Version patterns is parsed is if it did.
type ErrNotVersionString string

func (err ErrNotVersionString) Error() string {
	return fmt.Sprintf("`%s` is not a recognized Azure version string", string(err))
}

var versionle = func() func(string, string) (bool, error) {
	versionPattern := regexp.MustCompile(`^(?P<year>[\d]{4})-(?P<month>[\d]{2})-(?P<day>[\d]{2})(?:[\.\-](?P<tag>.+))?$`)

	return func(left, right string) (result bool, err error) {
		leftMatch := versionPattern.FindStringSubmatch(left)
		rightMatch := versionPattern.FindStringSubmatch(right)

		if len(leftMatch) < 3 {
			err = ErrNotVersionString(left)
			return
		}

		if len(rightMatch) < 3 {
			err = ErrNotVersionString(right)
			return
		}

		var leftLERight = func(left, right string) (result bool, err error) {
			var leftNum, rightNum int
			leftNum, err = strconv.Atoi(left)
			if err != nil {
				return
			}
			rightNum, err = strconv.Atoi(right)
			if err != nil {
				return
			}

			result = leftNum <= rightNum
			return
		}

		for i := 1; i <= 3; i++ {
			var canShortCircuit bool
			canShortCircuit, err = leftLERight(leftMatch[i], rightMatch[i])
			if err != nil {
				return
			}

			if canShortCircuit {
				result = true
				return
			}
		}

		result = leftMatch[4] <= rightMatch[4]

		return
	}
}()
