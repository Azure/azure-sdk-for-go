package main

import (
	"bytes"
	"flag"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
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
		handle(err)

		m = append(m, Module{
			Name:    goModFile,
			Replace: fmt.Sprintf("../%s", relativePackagePath),
			Version: version,
		})
	}

	return m
}

// GetTopLevel runs "git rev-parse --show-toplevel" to get the an absolute path to the current repo
func GetTopLevel() string {
	topLevel, err := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	handle(err)
	return strings.ReplaceAll(bytes.NewBuffer(topLevel).String(), "\n", "")
}

// Creates a smoketests directory and creates a go.mod file
func buildSmokeTestDirectory() {
	topLevel := GetTopLevel()
	root := strings.ReplaceAll(topLevel, "\n", "")
	smoketestDir = filepath.Join(root, "sdk", "smoketests")

	_ = os.MkdirAll(smoketestDir, 0777)

	err := os.Chdir(smoketestDir)
	handle(err)

	// Create go.mod file
	f, err := os.Create(filepath.Join(smoketestDir, "go.mod"))
	handle(err)
	smoketestModFile = f.Name()

	err = f.Close()
	handle(err)
}

func BuildModFile(modules []Module) error {
	fmt.Println("Creating mod file manully...")

	f, err := os.OpenFile(smoketestModFile, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("module github.com/Azure/azure-sdk-for-go/sdk/smoketests\n\n%s\n\n", getVersion()))

	if err != nil {
		return err
	}

	fmt.Println("Starting with replace")
	replaceString := "replace %s => %s\n"
	for _, module := range modules {
		s := fmt.Sprintf(replaceString, module.Name, module.Replace)
		_, err = f.Write([]byte(s))

		if err != nil {
			return err
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

		if err != nil {
			return err
		}
	}

	_, err = f.WriteString(")")
	return err
}

func VerifyGoMod() {
	// Make sure in sdk/smoketests
	dir, err := os.Getwd()
	handle(err)
	if !strings.Contains(dir, "sdk/smoketests") {
		// Navigate to sdk/smoketests
		os.Chdir(smoketestDir)
	}

	output, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		fmt.Printf("Error running go mod tidy: %s\n", bytes.NewBuffer(output).String())
		panic(err)
	}
}

// CleanUp removes the sdk/moketests directory. We are okay with this failing (however it shouldn't)
func CleanUp() {
	fmt.Println("Cleaning up...")
	err := os.RemoveAll(smoketestDir)
	if err != nil {
		log.Printf("Could not remove smoketest directory\n\t%s\n", err.Error())
	}
}

// Find the example func which will be in the format func ExampleNewClient or ExampleNew***Client
func FindClientExample(packageDir string) string {
	// Open the example_test.go file which will have the examples
	topLevel := GetTopLevel()
	packageDir = strings.TrimPrefix(packageDir, "github.com/Azure/azure-sdk-for-go")
	exampleFile := filepath.Join(topLevel, packageDir, "example_test.go")

	f, err := os.OpenFile(exampleFile, os.O_RDONLY, 0666)
	defer func() {
		err = f.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	srcBytes, err := ioutil.ReadFile(exampleFile)
	if err != nil {
		log.Printf("Could not read example file at %s. Failed with error: \n%s", exampleFile, err.Error())
		panic(err)
	}
	src := bytes.NewBuffer(srcBytes).String()

	// Rename the package portion to `package smoketests`
	m := regexp.MustCompile("package [a-zA-Z_]*")
	replaceString := "package smoketests"
	res := m.ReplaceAllString(src, replaceString)
	fmt.Println(res[:300])

	return ""
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
		destinationPath := filepath.Join(dest, fmt.Sprintf("%d.go", h.Sum32()))

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
	packageName := "smoketests"
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".go") {
			handle(err)
			data, err := ioutil.ReadFile(path)
			handle(err)

			datastring := bytes.NewBuffer(data).String()

			m := regexp.MustCompile("package (.*)&")
			fmt.Println(m.FindAllStringIndex(datastring, -1))
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
	fmt.Println("Running smoketest")

	rootDirectory := flag.String("rootDirectory", "", "root directory to find packages in")

	flag.Parse()

	if *rootDirectory == "" {
		fmt.Println("-rootDirectory argument must be provided")
		os.Exit(1)
	}

	absSDKPath, err := filepath.Abs(fmt.Sprintf("%s/sdk", *rootDirectory))
	handle(err)
	fmt.Println("Root SDK directory: ", absSDKPath)

	smoketestDir = filepath.Join(absSDKPath, "smoketests")
	fmt.Println("Smoke test directory: ", smoketestDir)
	err = os.Mkdir(smoketestDir, 0666)
	if err != os.ErrExist {
		handle(err)
	}

	smoketestModFile = filepath.Join(smoketestDir, "go.mod")
	f, err := os.Create(smoketestModFile)
	handle(err)
	err = f.Close()
	handle(err)

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

	err = BuildModFile(modules)
	handle(err)

	err = ReplacePackageStatement(smoketestDir)
	handle(err)
}
