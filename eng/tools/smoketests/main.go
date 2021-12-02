package main

import (
	"bytes"
	"flag"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var smoketestModFile string
var smoketestDir string

func handle(e error) {
	if e != nil {
		panic(e)
	}
}

func getVersion() string {
	v := runtime.Version()
	if strings.Contains(v, "go") {
		v = strings.TrimLeft(v, "go")

		return fmt.Sprintf("go %s", v)
	}

	// Default, go is not from a tag
	return "go 1.17"
}

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

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		handle(err)
		if strings.Contains(info.Name(), "go.mod") && !inIgnoredDirectories(path) {
			path = strings.ReplaceAll(path, "\\", "/")
			path = strings.ReplaceAll(path, "/go.mod", "")
			parts := strings.Split(path, "github.com/")
			formatted := fmt.Sprintf("github.com/%s", parts[1])
			ret = append(ret, formatted)
		}
		return nil
	})
	handle(err)

	return ret
}

func getAllTags() []string {
	result, err := exec.Command("git", "tag", "-l").Output()
	handle(err)
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
	handle(err)
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
		packagePath := strings.Split(goModFile, "github.com/Azure/azure-sdk-for-go/sdk/")
		relativePackagePath := packagePath[1]
		version, err := findLatestTag(relativePackagePath, tags)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			m = append(m, Module{
				Name:    goModFile,
				Replace: fmt.Sprintf("../%s", relativePackagePath),
				Version: version,
			})
		}
	}

	return m
}

// GetTopLevel runs "git rev-parse --show-toplevel" to get the an absolute path to the current repo
func GetTopLevel() string {
	topLevel, err := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	handle(err)
	return strings.ReplaceAll(bytes.NewBuffer(topLevel).String(), "\n", "")
}

// BuildModFile creates a go.mod file and adds replace directives for the appropriate modules.
// If serviceDirectory is a blank string it replaces all modules, otherwise it only replaces matching modules
func BuildModFile(modules []Module, serviceDirectory string) error {
	fmt.Println("Creating mod file manully at ", smoketestModFile)

	f, err := os.OpenFile(smoketestModFile, os.O_RDWR, 0666)
	handle(err)
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("module github.com/Azure/azure-sdk-for-go/sdk/smoketests\n\n%s\n\n", getVersion()))
	if err != nil {
		return err
	}

	replaceString := "replace %s => %s\n"
	if serviceDirectory == "notset" {
		fmt.Println("Starting with replace")
		for _, module := range modules {
			s := fmt.Sprintf(replaceString, module.Name, module.Replace)
			_, err = f.Write([]byte(s))
			handle(err)
		}
	} else {
		fmt.Printf("Replace directive for %s\n", serviceDirectory)
		for _, module := range modules {
			if strings.Contains(module.Name, serviceDirectory) {
				s := fmt.Sprintf(replaceString, module.Name, module.Replace)
				_, err = f.Write([]byte(s))
				handle(err)
			}
		}
	}

	fmt.Println("Require portion")
	_, err = f.WriteString("\n\nrequire (\n")

	if err != nil {
		return err
	}

	requireString := "\t%s %s\n"
	for _, module := range modules {
		s := fmt.Sprintf(requireString, module.Name, module.Version)
		_, err = f.Write([]byte(s))
		handle(err)
	}

	_, err = f.WriteString(")")
	handle(err)
	return nil
}

// FindExampleFiles finds all files that are named "example_*.go".
func FindExampleFiles(root string) ([]string, error) {
	var ret []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		handle(err)
		if strings.HasPrefix(info.Name(), "example_") && !inIgnoredDirectories(path) && strings.HasSuffix(info.Name(), ".go") {
			fName := path
			fName = strings.ReplaceAll(fName, "\\", "/")
			ret = append(ret, fName)
		}
		return nil
	})

	return ret, err
}

// copyFile copies the contents from src to dest. Creating the dest file first
func copyFile(src, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, data, 0644)
	return err
}

// CopyExampleFiles copies all the example files to the destination directory.
// This creates a hash of the fileName for the destination path
func CopyExampleFiles(exFiles []string, dest string) error {
	fmt.Printf("Copying %d example files to %s\n", len(exFiles), dest)

	for _, exFile := range exFiles {
		h := fnv.New32a()
		_, err := h.Write([]byte(exFile))
		if err != nil {
			return err
		}
		newFileName := strings.ReplaceAll(exFile[10:], "/", "_")
		newFileName = strings.ReplaceAll(newFileName, " ", "")
		fmt.Println(newFileName)
		destinationPath := filepath.Join(dest, fmt.Sprintf("%s.go", newFileName))

		err = copyFile(exFile, destinationPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReplacePackageStatement replaces all "package ***" with a common "package smoketests" statement
func ReplacePackageStatement(root string) error {
	fmt.Println("Fixing package names in", root)
	packageName := "package smoketests"
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".go") {
			handle(err)
			data, err := ioutil.ReadFile(path)
			handle(err)

			datastring := bytes.NewBuffer(data).String()

			m := regexp.MustCompile("(?m)^package (.*)$")
			datastring = m.ReplaceAllString(datastring, packageName)

			err = ioutil.WriteFile(path, []byte(datastring), 0666)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func main() {
	serviceDirectory := flag.String("serviceDirectory", "notset", "pass in a single service directory for nightly run")
	flag.Parse()

	fmt.Println("Running smoketest")

	rootDirectory := GetTopLevel()

	absSDKPath, err := filepath.Abs(fmt.Sprintf("%s/sdk", rootDirectory))
	handle(err)
	fmt.Println("Root SDK directory: ", absSDKPath)

	smoketestDir = filepath.Join(absSDKPath, "smoketests")
	fmt.Println("Smoke test directory: ", smoketestDir)
	// Create directory if it does not exist
	// _ = os.Mkdir(smoketestDir, 0666)

	smoketestModFile = filepath.Join(smoketestDir, "go.mod")
	// f, err := os.Create(smoketestModFile)
	// handle(err)
	// err = f.Close()
	// handle(err)

	exampleFiles, err := FindExampleFiles(absSDKPath)
	handle(err)
	fmt.Printf("Found %d example files for smoke tests\n", len(exampleFiles))
	for _, e := range exampleFiles {
		fmt.Println(e)
	}

	moduleDirectories := findModuleDirectories(absSDKPath)
	fmt.Printf("Found %d modules\n", len(moduleDirectories))

	allTags := getAllTags()
	fmt.Printf("Found %d tags\n", len(allTags))

	modules := matchModulesAndTags(moduleDirectories, allTags)
	_ = modules

	err = CopyExampleFiles(exampleFiles, smoketestDir)
	handle(err)

	err = BuildModFile(modules, *serviceDirectory)
	handle(err)

	err = ReplacePackageStatement(smoketestDir)
	handle(err)
}
