package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type CodeCoverage struct {
	Packages []Package `json:"Packages"`
}

type Package struct {
	Name         string  `json:"name"`
	CoverageGoal float64 `json:"CoverageGoal"`
}

const (
	coverageXmlFile = "coverage.xml"
)

var configFile, _ = filepath.Abs(filepath.Join("..", "eng", "config.json"))

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var coverageFiles []string

func filterFiles(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if strings.Contains(path, "coverage.xml") {
		coverageFiles = append(coverageFiles, path)
	}
	return nil
}

func FindCoverageFiles(p string) {
	err := filepath.Walk(".", filterFiles)
	check(err)
}

func ReadConfigData() *CodeCoverage {
	jsonFile, err := os.Open(configFile)
	check(err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	check(err)

	var cov CodeCoverage
	err = json.Unmarshal([]byte(byteValue), &cov)
	check(err)
	return &cov
}

// This supports doing a single package at a time. If this needs to be expanded in the future
// this method will have to return a []*float64 for each packages goal
func findCoverageGoal(covFiles []string, configData *CodeCoverage) float64 {
	for _, covFile := range covFiles {
		for _, p := range configData.Packages {
			if strings.Contains(covFile, p.Name) {
				return p.CoverageGoal
			}
		}
	}
	fmt.Println("WARNING: Could not find a coverage goal, defaulting to 95%.")
	return 0.95
}

func main() {

	serviceDir := flag.String("serviceDirectory", "", "Service Directory")
	flag.Parse()

	coverageFiles = make([]string, 0)
	rootPath, err := filepath.Abs(".")
	check(err)

	FindCoverageFiles(rootPath)

	configData := ReadConfigData()
	coverageGoal := findCoverageGoal([]string{*serviceDir}, configData)

	fmt.Printf("Failing if the coverage is below %.2f\n", coverageGoal)

	coverageValues := make([]float64, 0)
	for _, coverageFile := range coverageFiles {
		xmlFile, err := os.Open(coverageFile)
		check(err)
		defer xmlFile.Close()

		byteValue, err := ioutil.ReadAll(xmlFile)
		check(err)

		re := regexp.MustCompile(`<coverage line-rate=\"\d\.\d+\"`)
		coverageValue := re.Find(byteValue) //, -1)
		if coverageValue == nil {
			log.Fatalf("Could not match regexp to coverage.xml file.")
		}

		parts := strings.Split(string(coverageValue), "=")
		coverageNumber := parts[1]

		coverageNumber = coverageNumber[1 : len(coverageNumber)-1]
		coverageFloat, err := strconv.ParseFloat(coverageNumber, 32)
		check(err)
		coverageValues = append(coverageValues, coverageFloat)
	}

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
		fmt.Println("Coverage step failed")
		os.Exit(1)
	}
}
