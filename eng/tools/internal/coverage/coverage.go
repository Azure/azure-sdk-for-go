// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type codeCoverage struct {
	Packages []coveragePackage `json:"Packages"`
}

type coveragePackage struct {
	Name         string  `json:"name"`
	CoverageGoal float64 `json:"CoverageGoal"`
}

const (
	coverageFile = "coveragefunc.txt"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func findCoverageFiles(root string) []string {
	var coverageFiles []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == coverageFile {
			coverageFiles = append(coverageFiles, path)
		}

		return nil
	})

	check(err)

	return coverageFiles
}

func readConfigData(coverageConfig string) *codeCoverage {
	jsonFile, err := os.Open(coverageConfig)
	check(err)
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	check(err)

	var cov codeCoverage
	err = json.Unmarshal([]byte(byteValue), &cov)
	check(err)
	return &cov
}

// This supports doing a single package at a time. If this needs to be expanded in the future
// this method will have to return a []*float64 for each packages goal
//
// Packages are identified in configData (parsed from /eng/config.json) by segments of their
// module paths, for example "keyvault/azkeys", whereas covFile may be a disk path like
// "/vss/1/sdk/security/keyvault/azkeys". So, this function searches for the configData entry
// whose Name is the longest substring of covFile.
func findCoverageGoal(covFiles []string, configData *codeCoverage) float64 {
	for _, covFile := range covFiles {
		covFile = strings.ReplaceAll(covFile, `\`, "/")
		var bestMatch *coveragePackage
		for i := range configData.Packages {
			p := &configData.Packages[i]
			prx := regexp.MustCompile(`(^|/)` + regexp.QuoteMeta(p.Name) + `($|/)`)
			if prx.MatchString(covFile) {
				if bestMatch == nil || len(p.Name) > len(bestMatch.Name) {
					bestMatch = p
				}
			}
		}
		if bestMatch != nil {
			return bestMatch.CoverageGoal
		}
	}
	fmt.Println("WARNING: Could not find a coverage goal, defaulting to 95%.")
	return 0.95
}

func parseCoveragePercent(contents []byte) (float64, error) {
	re := regexp.MustCompile(`total:.*?(\d+\.\d+)%`)
	matches := re.FindStringSubmatch(string(contents))

	if len(matches) < 2 {
		return 0, errors.New("could not match regexp to coveragefunc.txt file")
	}

	coverageFloat, err := strconv.ParseFloat(matches[1], 32)
	if err != nil {
		return 0, err
	}

	return coverageFloat / 100, nil
}

func parseCoverageFiles(coverageFiles []string) []float64 {
	coverageValues := make([]float64, 0)

	for _, coverageFile := range coverageFiles {
		fmt.Println(coverageFile)
		file, err := os.Open(coverageFile)
		check(err)
		defer file.Close()

		byteValue, err := io.ReadAll(file)
		check(err)

		coveragePercent, err := parseCoveragePercent(byteValue)
		check(err)

		coverageValues = append(coverageValues, coveragePercent)
	}

	return coverageValues
}

func CheckCoverage(serviceDir string, coverageConfig string, searchDirectory string) {
	rootPath, err := filepath.Abs(searchDirectory)
	check(err)

	fmt.Printf("Searching for coverage files in %s\n", rootPath)

	coverageFiles := findCoverageFiles(rootPath)
	if len(coverageFiles) == 0 {
		fmt.Println("No coverage files found in " + rootPath)
		return
	}

	fmt.Printf("Reading config data from %s\n", coverageConfig)

	configData := readConfigData(coverageConfig)
	coverageGoal := findCoverageGoal([]string{serviceDir}, configData)

	fmt.Printf("(%s) Failing if the coverage is below %.2f\n", serviceDir, coverageGoal)

	coverageValues := parseCoverageFiles(coverageFiles)

	if len(coverageValues) != len(coverageFiles) {
		fmt.Printf("Found %d coverage values in %d coverage files\n", len(coverageValues), len(coverageFiles))
	}

	failedCoverage := false
	for i := range coverageValues {
		status := "Succeeded"
		if coverageValues[i] < coverageGoal {
			status = "Failed"
		}
		fmt.Printf("Status: %v\tCoverage file: %v\t Coverage Amount: %.4f\n", status, coverageFiles[i], coverageValues[i])
		if coverageValues[i] < coverageGoal {
			failedCoverage = true
		}
	}

	if failedCoverage {
		log.Fatal("Found insufficient coverage.")
	}
}

func main() {
	coverageConfig := flag.String("config", "", "Coverage Threshold Configuration")
	serviceDir := flag.String("serviceDirectory", "", "Service Directory")
	searchDirectory := flag.String("searchDirectory", ".", "Search Directory")
	flag.Parse()
	if *coverageConfig == "" {
		log.Fatal("Required flag: '-config'")
	}
	CheckCoverage(*serviceDir, *coverageConfig, *searchDirectory)
}
