// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/indexer/util"
)

// adds any missing SDK packages to godoc.org
func main() {
	// by default assume we're running from the source dir
	// and calculate the relative path to the services directory.
	dir := "../../services"
	if len(os.Args) > 1 {
		// assume second arg is source dir
		dir = os.Args[1]
	}

	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	pkgs, err := util.GetPackagesForIndexing(dir)
	if err != nil {
		panic(err)
	}

	// this URL will return the set of packages that have been indexed
	resp, err := http.DefaultClient.Get("https://godoc.org/github.com/Azure/azure-sdk-for-go/services")
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(err)
	}
	defer resp.Body.Close()

	indexedPkgs, err := util.GetIndexedPackages(resp.Body)
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

		// performing a GET on the package URL will cause the service to index it
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
		time.Sleep(10 * time.Second)
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
