package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func inIgnoredDirectories(path string) bool {
	if strings.Contains(path, "internal") {
		return true
	}
	if strings.Contains(path, "samples") {
		return true
	}
	if strings.Contains(path, "smoketests") {
		return true
	}
	if strings.Contains(path, "template") {
		return true
	}

	return false
}

func findModuleDirectories(root string) []string {
	var ret []string

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if strings.Contains(info.Name(), "go.mod") && !inIgnoredDirectories(path) {
			path = strings.ReplaceAll(path, "\\", "/")
			path = strings.ReplaceAll(path, "/go.mod", "")
			parts := strings.Split(path, "github.com/")
			formatted := fmt.Sprintf("github.com/%s", parts[1])
			ret = append(ret, formatted)
		}
		return nil
	})

	return ret
}

func getAllTags() []string {
	result, err := exec.Command("git", "tag", "-l").Output()
	if err != nil {
		panic(err)
	}
	res := bytes.NewBuffer(result).String()
	return strings.Split(res, "\n")
}

type Module struct {
	Name    string
	Version string
	Replace string
}

type SemVer struct {
	Major, Minor, Patch int
}

func (s SemVer) Newer(s2 SemVer) bool {
	if s.Major > s2.Major {
		return true
	} else if s.Major == s2.Major && s.Minor > s2.Minor {
		return true
	} else if s.Major == s2.Major && s.Minor == s2.Minor && s.Patch > s2.Patch {
		return true
	}
	return false
}

func (s SemVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.Major, s.Minor, s.Patch)
}

func toInt(a string) int {
	r, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return r
}

func NewSemVerFromTag(s string) SemVer {
	path := strings.Split(s, "/")
	versionStr := path[len(path)-1]
	versionStr = strings.TrimLeft(versionStr, "v")
	parts := strings.Split(versionStr, ".")
	return SemVer{
		Major: toInt(parts[0]),
		Minor: toInt(parts[1]),
		Patch: toInt(parts[2]),
	}
}

func findLatestTag(p string, tags []string) (string, error) {
	fmt.Println("Searching for latest tag for ", p)
	var v SemVer
	for i, tag := range tags {
		if strings.Contains(tag, p) {
			v = NewSemVerFromTag(tag)
			for strings.Contains(tags[i+1], p) {
				newV := NewSemVerFromTag(tags[i+1])
				if newV.Newer(v) {
					v = newV
				}
				i += 1
			}
			return v.String(), nil
		}
	}
	return "", fmt.Errorf("could not find a version for module %s", p)
}

func matchModulesAndTags(goModFiles []string, tags []string) []Module {
	var m []Module

	for _, goModFile := range goModFiles {
		packagePath := strings.Split(goModFile, "github.com/Azure/azure-sdk-for-go/")
		relativePackagePath := packagePath[1]
		version, err := findLatestTag(relativePackagePath, tags)
		if err != nil {
			panic(err)
		}

		m = append(m, Module{
			Name:    goModFile,
			Replace: fmt.Sprintf("../%s", relativePackagePath),
			Version: version,
		})
	}

	return m
}

func main() {
	fmt.Println("Running smoketest")

	rootDirectory := flag.String("rootDirectory", "", "root directory to find packages in")

	flag.Parse()

	if *rootDirectory == "" {
		fmt.Println("-rootDirectory command must be provided")
		os.Exit(1)
	}

	absPath, err := filepath.Abs(fmt.Sprintf("%s/sdk", *rootDirectory))
	if err != nil {
		panic(err)
	}
	fmt.Println("Root directory: ", absPath)
	moduleDirectories := findModuleDirectories(absPath)
	fmt.Printf("Found %d modules\n", len(moduleDirectories))
	allTags := getAllTags()
	fmt.Printf("Found %d tags\n", len(allTags))

	modules := matchModulesAndTags(moduleDirectories, allTags)
	fmt.Println(modules)
}
