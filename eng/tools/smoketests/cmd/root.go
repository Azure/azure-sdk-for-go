package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "smoketests [Options]",
	Short: "Run smoketests for azure SDK for go",
	RunE: func(c *cobra.Command, args []string) error {
		return smokeTestCmd()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

var serviceDirectory string
var daily bool

func init() {
	rootCmd.PersistentFlags().StringVarP(&serviceDirectory, "serviceDirectory", "s", "", "serviceDirectory to test against")
	rootCmd.PersistentFlags().BoolVarP(&daily, "daily", "d", false, "For running daily tests set to true")
}

var smoketestModFile string
var smoketestDir string
var exampleFuncs []string
var envVars []string

func handle(e error) {
	if e != nil {
		panic(e)
	}
}

func getVersion() string {
	ver := regexp.MustCompile(`\d\.\d{2}`)
	if m := ver.FindString(runtime.Version()); m != "" {
		return fmt.Sprintf("go %s", m)
	}

	// Default, go is not from a tag
	return "go 1.19"
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

// Walks the sdk directory to find all modules based on a go.mod file
func findModuleDirectories(root string) []string {
	var ret []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		handle(err)
		if d.Name() == "go.mod" && !inIgnoredDirectories(path) {
			path = strings.ReplaceAll(path, "\\", "/")
			path = strings.ReplaceAll(path, "/go.mod", "")
			parts := strings.Split(path, "/sdk/")
			formatted := fmt.Sprintf("github.com/Azure/azure-sdk-for-go/sdk/%s", parts[1])
			ret = append(ret, formatted)
		}
		return nil
	})
	handle(err)

	return ret
}

// Reads all tags using the 'git tag -l' command and returns them as a string slice
func getAllTags() []string {
	result, err := exec.Command("git", "tag", "-l").Output()
	handle(err)
	res := bytes.NewBuffer(result).String()
	return strings.Split(res, "\n")
}

// Create a new SemVer type
func NewSemVerFromTag(s string) (*semver.Version, error) {
	path := strings.Split(s, "/")
	versionStr := path[len(path)-1]
	versionStr = strings.TrimLeft(versionStr, "v")
	return semver.NewVersion(versionStr)
}

// Find the most recent SemVer tag for a given package.
func findLatestTag(p string, tags []string) (*semver.Version, error) {
	var v *semver.Version
	var err error
	for i, tag := range tags {
		if strings.Contains(tag, p) {
			v, err = NewSemVerFromTag(tag)
			if err != nil {
				return nil, fmt.Errorf("could not parse version for tag %s", tag)
			}
			for strings.Contains(tags[i+1], p) {
				newV, err := NewSemVerFromTag(tags[i+1])
				if err != nil {
					return nil, fmt.Errorf("could not parse version for tag %s", tags[i+1])
				}
				if newV.GreaterThan(v) {
					v = newV
				}
				i += 1
			}
			return v, nil
		}
	}
	return nil, fmt.Errorf("could not find a version for module %s", p)
}

// Creates a slice of modules matched with the most recent version
func matchModulesAndTags(goModFiles []string, tags []string) []Module {
	var m []Module

	for _, goModFile := range goModFiles {
		packagePath := strings.Split(goModFile, "github.com/Azure/azure-sdk-for-go/sdk/")
		relativePackagePath := packagePath[1]
		version, err := findLatestTag(relativePackagePath, tags)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if version.Major() > 1 {
				m = append(m, Module{
					Name:    goModFile,
					Package: fmt.Sprintf("%s/v%d", goModFile, version.Major()),
					Replace: fmt.Sprintf("../%s", relativePackagePath),
					Version: "v" + version.String(),
				})
			} else {
				m = append(m, Module{
					Name:    goModFile,
					Package: goModFile,
					Replace: fmt.Sprintf("../%s", relativePackagePath),
					Version: "v" + version.String(),
				})
			}
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
	f, err := os.OpenFile(smoketestModFile, os.O_RDWR, 0666)
	handle(err)
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("module github.com/Azure/azure-sdk-for-go/sdk/smoketests\n\n%s\n\n", getVersion()))
	if err != nil {
		return err
	}

	replaceString := "replace %s => %s\n"
	if serviceDirectory == "notset" {
		for _, module := range modules {
			s := fmt.Sprintf(replaceString, module.Package, module.Replace)
			_, err = f.Write([]byte(s))
			handle(err)
		}
	} else {
		fmt.Printf("Replace directive for %s\n", serviceDirectory)
		for _, module := range modules {
			if strings.Contains(module.Name, serviceDirectory) {
				s := fmt.Sprintf(replaceString, module.Package, module.Replace)
				_, err = f.Write([]byte(s))
				handle(err)
			}
		}
	}

	_, err = f.WriteString("\n\nrequire (\n")

	if err != nil {
		return err
	}

	requireString := "\t%s %s\n"
	for _, module := range modules {
		s := fmt.Sprintf(requireString, module.Package, module.Version)
		_, err = f.Write([]byte(s))
		handle(err)
	}

	_, err = f.WriteString(")")
	handle(err)
	return nil
}

// FindExampleFiles finds all files that are named "example_*.go".
// If serviceDirectory
func FindExampleFiles(root, serviceDirectory string) ([]string, error) {
	var ret []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		handle(err)
		if strings.HasPrefix(d.Name(), "example_") && !inIgnoredDirectories(path) && filepath.Ext(d.Name()) == ".go" {
			path = strings.ReplaceAll(path, "\\", "/")
			if serviceDirectory == "" || strings.Contains(path, serviceDirectory) {
				ret = append(ret, path)
			}

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

	stringData := bytes.NewBuffer(data).String()

	functionRegex := regexp.MustCompile(`(?m)^func\sExampleNew(.*)\(\)\s{`)
	functions := functionRegex.FindAllString(stringData, -1)
	for i, l := range functions {
		l = strings.ReplaceAll(l, "()", "")
		parts := strings.Split(l, " ")
		functions[i] = parts[1]
	}

	// do a name change for clashing
	var newNames []string
	for _, f := range functions {
		newNames = append(newNames, fmt.Sprintf("%s%d", f, rand.Intn(100000)))
	}

	for i := 0; i < len(functions); i++ {
		stringData = strings.Replace(stringData, functions[i], newNames[i], 1)
		exampleFuncs = append(exampleFuncs, newNames[i])
	}

	err = ioutil.WriteFile(dest, []byte(stringData), 0644)
	return err
}

// CopyExampleFiles copies all the example files to the destination directory.
// This creates a hash of the fileName for the destination path
func CopyExampleFiles(exFiles []string, dest string) {
	fmt.Printf("Copying %d example files to %s\n", len(exFiles), dest)

	for _, exFile := range exFiles {
		newFileName := strings.ReplaceAll(exFile[10:], "/", "_")
		newFileName = strings.ReplaceAll(newFileName, " ", "")
		newFileName = strings.ReplaceAll(newFileName, "_test", "")
		destinationPath := filepath.Join(dest, newFileName)

		err := copyFile(exFile, destinationPath)
		handle(err)
	}
}

// ReplacePackageStatement replaces all "package ***" with a common "package main" statement
func ReplacePackageStatement(root string) error {
	packageName := "package main"
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(d.Name()) == ".go" {
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

// Create the main.go file with environment variables, example functions, and imports.
func BuildMainFile(root string, c ConfigFile) error {
	mainFile := filepath.Join(root, "main.go")
	f, err := os.Create(mainFile)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}

	// Write the main.go file

	src := "package main\nimport \"os\"\nfunc main() {"
	if len(envVars) == 0 {
		src = "package main\nfunc main() {\n"
	}

	for _, envVar := range envVars {
		src += fmt.Sprintf(`os.Setenv("%s", "%s")`, envVar, FindEnvVarFromConfig(c, envVar))
		src += "\n"
	}

	src += "\n"

	for _, exampleFunc := range exampleFuncs {
		src += fmt.Sprintf("%s()\n", exampleFunc)
	}

	src += "}"

	err = ioutil.WriteFile(mainFile, []byte(src), 0666)
	return err
}

// Find all the environment variables looked up within a .go file
func FindEnvVars(root string) error {
	fmt.Println("Find all environment variables using `os.Getenv` or `os.LookupEnv`")

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".go" {
			// Find Env Vars
			searchFile(path)
		}
		return nil
	})
	return err
}

// Search for both os.Getenv and os.LookupEnv in go file
func searchFile(path string) error {
	envVarRegex := regexp.MustCompile(`(?m)os.LookupEnv\("(.*)"\)`)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	stringData := bytes.NewBuffer(data).String()
	LookupEnvs := envVarRegex.FindAllString(stringData, -1)
	envVars = append(envVars, trimLookupEnvs(LookupEnvs)...)

	envVarRegex = regexp.MustCompile(`(?m)os.Getenv(.*)\r\n`)
	Getenvs := envVarRegex.FindAllString(stringData, -1)
	envVars = append(envVars, trimGetenvs(Getenvs)...)
	return nil
}

// Search for a default value in the eng/config.json file
func FindEnvVarFromConfig(c ConfigFile, envVar string) string {
	for _, p := range c.Packages {
		for key, value := range p.EnvironmentVariables {
			if key == envVar {
				return value
			}
		}
	}

	return ""
}

func trimLookupEnvs(values []string) []string {
	pseudoSet := make(map[string]struct{})

	for _, value := range values {
		value = strings.TrimSpace(value)
		value = strings.TrimPrefix(value, `os.LookupEnv("`)
		value = strings.TrimSuffix(value, `")`)
		pseudoSet[value] = struct{}{}
	}

	var ret []string
	for v := range pseudoSet {
		ret = append(ret, v)
	}

	return ret
}

func trimGetenvs(values []string) []string {
	var ret []string

	for _, value := range values {
		value = strings.TrimPrefix(value, `os.Getenv("`)
		value = strings.TrimSuffix(value, `")`)
		ret = append(ret, value)
	}

	return ret
}

// Read the eng/config.json file to parse environment variables
func LoadEngConfig(rootDirectory string) ConfigFile {
	fileName := filepath.Join(rootDirectory, "eng", "config.json")
	buffer, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	var p ConfigFile
	err = json.Unmarshal(buffer, &p)
	handle(err)
	return p
}

func smokeTestCmd() error {

	if serviceDirectory == "" && !daily {
		fmt.Println("Arguments could not be understood. Either specify a serviceDirectory or set --daily to true")
		os.Exit(-1)
	} else if serviceDirectory != "" && !daily {
		fmt.Printf("Running smoke tests on the %s directory\n", serviceDirectory)
	}

	rootDirectory := GetTopLevel()

	configFile := LoadEngConfig(rootDirectory)

	absSDKPath, err := filepath.Abs(fmt.Sprintf("%s/sdk", rootDirectory))
	handle(err)
	fmt.Println("Root SDK directory: ", absSDKPath)

	smoketestDir = filepath.Join(absSDKPath, "smoketests")
	fmt.Println("Smoke test directory: ", smoketestDir)

	smoketestModFile = filepath.Join(smoketestDir, "go.mod")

	exampleFiles, err := FindExampleFiles(absSDKPath, serviceDirectory)
	handle(err)
	fmt.Printf("Found %d example files for smoke tests\n", len(exampleFiles))
	for _, e := range exampleFiles {
		fmt.Printf("\t%s\n", e)
	}

	moduleDirectories := findModuleDirectories(absSDKPath)
	allTags := getAllTags()
	modules := matchModulesAndTags(moduleDirectories, allTags)

	CopyExampleFiles(exampleFiles, smoketestDir)

	err = BuildModFile(modules, serviceDirectory)
	handle(err)

	err = ReplacePackageStatement(smoketestDir)
	handle(err)

	err = FindEnvVars(smoketestDir)
	handle(err)

	return BuildMainFile(smoketestDir, configFile)
}
