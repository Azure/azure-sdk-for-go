package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
		packagePath := strings.Split(goModFile, "github.com/Azure/azure-sdk-for-go/")
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

// Creates a smoketests directory and initializes a go.mod file by running go mod init.
// It returns a function to clean up the created directory
func buildSmokeTestDirectory() {
	topLevel, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	handle(err)
	root := strings.ReplaceAll(bytes.NewBuffer(topLevel).String(), "\n", "")
	smoketestDir = filepath.Join(root, "sdk", "smoketests")

	_ = os.MkdirAll(smoketestDir, 0777)
	handle(err)

	err = os.Chdir(smoketestDir)
	handle(err)

	// Create go.mod file
	f, err := os.Create(filepath.Join(smoketestDir, "go.mod"))
	handle(err)
	smoketestModFile = f.Name()

	err = f.Close()
	handle(err)
}

func BuildModFile(modules []Module) {
	fmt.Println("Creating mod file manully...")

	f, err := os.OpenFile(smoketestModFile, os.O_RDWR, 0666)
	handle(err)
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("module github.com/Azure/azure-sdk-for-go/sdk/smoketests\n\n%s\n\n", getVersion()))
	handle(err)

	fmt.Println("Starting with replace")
	replaceString := "replace %s => %s\n"
	for _, module := range modules {
		s := fmt.Sprintf(replaceString, module.Name, module.Replace)
		_, err = f.Write([]byte(s))
		handle(err)
	}

	fmt.Println("Require portion")
	_, err = f.WriteString("\n\nrequire (\n")
	handle(err)

	requireString := "\t%s %s\n"
	for _, module := range modules {
		s := fmt.Sprintf(requireString, module.Name, module.Version)
		_, err = f.Write([]byte(s))
		handle(err)
	}

	_, err = f.WriteString(")")
	handle(err)
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

func main() {
	fmt.Println("Running smoketest")

	rootDirectory := flag.String("rootDirectory", "", "root directory to find packages in")

	flag.Parse()

	if *rootDirectory == "" {
		fmt.Println("-rootDirectory argument must be provided")
		os.Exit(1)
	}

	absPath, err := filepath.Abs(fmt.Sprintf("%s/sdk", *rootDirectory))
	handle(err)
	fmt.Println("Root directory: ", absPath)

	moduleDirectories := findModuleDirectories(absPath)
	fmt.Printf("Found %d modules\n", len(moduleDirectories))

	allTags := getAllTags()
	fmt.Printf("Found %d tags\n", len(allTags))

	modules := matchModulesAndTags(moduleDirectories, allTags)

	buildSmokeTestDirectory()

	// Build go.mod file
	BuildModFile(modules)

	// Run go.mod tidy
	VerifyGoMod()

	fmt.Println("Successfully ran smoketests")

	CleanUp()
}
