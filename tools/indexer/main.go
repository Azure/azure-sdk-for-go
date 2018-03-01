// Copyright 2018 Microsoft Corporation
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

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type packageSet map[string]bool

// adds any missing SDK packages to godoc.org
func main() {
	dir := ""
	if len(os.Args) > 1 {
		// assume second arg is source dir
		dir = os.Args[1]
	} else {
		// if no args specified assume we're running from the source dir
		// and calculate the relative path to the services directory.
		var err error
		dir, err = filepath.Abs("../../services")
		if err != nil {
			panic(err)
		}
	}

	pkgs, err := getPackagesForIndexing(dir)
	if err != nil {
		panic(err)
	}

	indexedPkgs, err := getIndexedPackages()
	if err != nil {
		panic(err)
	}

	// for each package in pkgs, check if it's already been
	// indexed.  if it hasn't been indexed then do so

	for pkg := range pkgs {
		if _, already := indexedPkgs[pkg]; already {
			pkgs[pkg] = true
			continue
		}

		fmt.Printf("indexing %s...", pkg)
		resp, err := http.DefaultClient.Get(fmt.Sprintf("https://godoc.org/%s", pkg))
		if err != nil {
			panic(err)
		}

		resp.Body.Close()
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
			// a 404 means that the package exists locally but not yet in github
			fmt.Printf("completed with status '%s'\n", resp.Status)
			pkgs[pkg] = true
		} else {
			fmt.Printf("FAILED with status '%s'\n", resp.Status)
		}

		// sleep a bit between indexing
		time.Sleep(60 * time.Second)
	}

	// check if any packages failed to index
	failed := false
	for _, v := range pkgs {
		if !v {
			failed = true
			break
		}
	}

	if failed {
		fmt.Println("not all packages were indexed")
		os.Exit(1)
	}
	fmt.Println("successfully indexed all packages")
}

// returns the set of packages that have already been indexed
func getIndexedPackages() (packageSet, error) {
	// this URL will return the set of packages that have been indexed
	resp, err := http.DefaultClient.Get("https://godoc.org/github.com/Azure/azure-sdk-for-go/services")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("didn't receive 200 when looking for indexed packages, got %v: %s", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) < 1 {
		return nil, errors.New("did't receive a response body when lookinig for indexed packages")
	}

	// scrape the response body to create the package list
	pkgs := packageSet{}
	bodyStr := string(body)
	const lookFor = "github.com/Azure/azure-sdk-for-go/services/"
	for {
		// find the start of the package
		index := strings.Index(bodyStr, lookFor)
		if index < 0 {
			break
		}

		// now find the terminal " character
		terminal := strings.IndexRune(bodyStr[index+len(lookFor):], '"') + index + len(lookFor)
		pkg := bodyStr[index:terminal]
		pkgs[pkg] = true

		// advance to the next
		bodyStr = bodyStr[terminal:]
	}
	return pkgs, nil
}

// returns the set of packages, calculated from the repo, to be indexed
func getPackagesForIndexing(dir string) (packageSet, error) {
	leafDirs := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// check if leaf dir
			fi, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}
			hasSubDirs := false
			for _, f := range fi {
				if f.IsDir() {
					hasSubDirs = true
					break
				}
			}
			if !hasSubDirs {
				leafDirs = append(leafDirs, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// dirs will look something like "D:\work\src\github.com\Azure\azure-sdk-for-go\services\..."
	// strip off the stuff before the github.com and change the whacks so it looks like a package import

	pkgs := packageSet{}
	for _, dir := range leafDirs {
		i := strings.Index(dir, "github.com")
		if i < 0 {
			return nil, fmt.Errorf("didn't find github.com in directory '%s'", dir)
		}
		pkg := strings.Replace(dir[i:], "\\", "/", -1)
		pkgs[pkg] = false
	}
	return pkgs, nil
}
