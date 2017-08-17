package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/marstr/collection"
)

type LatestStrategy struct {
	Root      string
	Predicate func(packageName string) bool
}

// AcceptAll is a predefined value for `LatestStrategy.Predicate` which always returns true.
func AcceptAll(name string) (result bool) {
	result = true
	return
}

// IgnorePreview searches a packages "API Version" to see if it contains the word "preview". It only returns true when a package is not a preview.
func IgnorePreview(name string) (result bool) {
	matches := packageName.FindStringSubmatch(name)
	if len(matches) >= 4 {
		result = !strings.Contains(matches[3], "preview")
	}
	return
}

var packageName = regexp.MustCompile(`service[/\\](?P<provider>[\w\-\.\d]+)[/\\](?P<type>[\w\-\.\d]+)[/\\](?:management[/\\])?(?P<version>[\d\-\w\.]+)[/\\](?P<group>[\w\d\-\.]+)`)

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

			version := pathMatches[3]
			currentGroup := operationGroup{
				provider:     pathMatches[1],
				resourceType: pathMatches[2],
				group:        pathMatches[4],
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
			files := token.NewFileSet()
			parsed, err := parser.ParseDir(files, entry.rawpath, nil, parser.ParseComments)
			if err != nil {
				continue
			}

			for _, pkg := range parsed {
				select {
				case results <- pkg:
					// Intentionally Left Blank
				case <-cancel:
					return
				}
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
