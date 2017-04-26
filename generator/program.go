package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	defaultRemoteAzureRestAPISpecsPath = "https://github.com/Azure/azure-rest-api-specs.git"
	defaultAzureRESTAPIBranch          = "master"
)

// ExitCode gives a hint to the end user about why the program exited without relying on seeing stderr.
const (
	ExitCodeOkay int = iota
	ExitCodeMissingRequirements
	ExitCodeFailedToClone
)

var (
	remoteAzureRestAPISpecsPath string
	localAzureRestAPISpecsPath  string
	azureRestAPIBranch          string
	outputLocation              string
	dryRun                      bool
	help                        bool
	anyMissing                  bool
	noClone                     bool
	verbose                     bool
	targetFile                  string
)

func init() {
	flag.StringVar(&azureRestAPIBranch, "branch", defaultAzureRESTAPIBranch, "The branch, tag, or SHA1 identifier in the Azure Rest API Specs repository to use during API generation.")
	flag.StringVar(&remoteAzureRestAPISpecsPath, "repo", defaultRemoteAzureRestAPISpecsPath, "The path to the location of the Azure REST API Specs repository that should be used for generation.")
	flag.StringVar(&outputLocation, "output", getDefaultOutputLocation(), "a directory in which all output should be placed.")
	flag.StringVar(&targetFile, "target", "", "If a target file is provided, generator will only run on this file instead of all swaggers it encounters in the repository.")
	flag.BoolVar(&dryRun, "preview", false, "Use this flag to print a list of packages that would be generated instead of actually generating the new sdk.")
	flag.BoolVar(&help, "help", false, "Provide this flag to print usage information instead of running the program.")
	flag.BoolVar(&noClone, "noClone", false, "Use this flag to prevent cloning a new copy of Azure REST API Specs. The existing enlistment should be used instead.")
	flag.BoolVar(&verbose, "verbose", false, "Print status messages as processing is done.")

	flag.Parse()

	if help {
		return
	}

	optionalTools := []string{"gofmt", "golint"}
	requiredTools := []string{"autorest", "git"}

	for _, tool := range optionalTools {
		if _, err := exec.LookPath(tool); err != nil {
			log.Printf("WARNING: Could not find \"%s\" usage of this tool will be skipped.", tool)
		}
	}

	anyMissing = false
	for _, tool := range requiredTools {
		if _, err := exec.LookPath(tool); err != nil {
			log.Printf("ERROR: Could not find \"%s\" this tool will exit without executing.", tool)
			anyMissing = true
		}
	}

	if noClone {
		localAzureRestAPISpecsPath = remoteAzureRestAPISpecsPath
	} else {
		var err error
		localAzureRestAPISpecsPath, err = ioutil.TempDir("./", "")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(ExitCodeFailedToClone)
		}
	}
}

func main() {
	var repoLoc string
	exitStatus := ExitCodeOkay
	defer os.Exit(exitStatus)

	errLog := log.New(os.Stderr, "[ERROR] ", 0)

	var statusWriter io.Writer
	if verbose {
		statusWriter = os.Stdout
	} else {
		statusWriter = ioutil.Discard
	}
	statusLog := log.New(statusWriter, "[STATUS] ", 0)

	if help {
		flag.Usage()
		return
	}

	if anyMissing {
		exitStatus = ExitCodeMissingRequirements
		return
	}

	if noClone == false {
		if temp, err := cloneAPISpecs(remoteAzureRestAPISpecsPath, localAzureRestAPISpecsPath); err == nil {
			repoLoc = temp
			defer func() {
				if err := os.RemoveAll(repoLoc); err != nil {
					errLog.Print(err)
				}
			}()
		} else {
			errLog.Print(err)
			exitStatus = ExitCodeFailedToClone
			return
		}
	}

	if err := checkoutAPISpecsVer(azureRestAPIBranch, repoLoc); err != nil {
		errLog.Print(err)
		exitStatus = ExitCodeFailedToClone
		return
	}

	swaggersToProcess := getSwaggers(repoLoc, statusLog, errLog)
	if dryRun {
		for manifest := range swaggersToProcess {
			namespace, err := getNamespace(manifest)
			if err != nil {
				errLog.Printf("Error with '%s': %v", manifest.Path, err)
			}
			fmt.Printf("%s -> %s\n", manifest.Path, namespace)
		}
		return
	}
}

func vetAll(packages <-chan string) (<-chan string, *log.Logger) {
	vetPackages := make(chan string)
	violationLog := log.New(os.Stdout, "vet", log.LstdFlags)

	go func() {
		defer close(vetPackages)
		for pkg := range packages {
			cmd := exec.Command("go", "vet", pkg)

			if err := cmd.Run(); err != nil {
				violationLog.Printf("error while vetting \"%s\": %v", pkg, err)
			}
		}
	}()

	return vetPackages, violationLog
}

func cloneAPISpecs(origin, local string) (string, error) {
	_, cloneLoc := filepath.Split(local)
	clone := exec.Command("git", "clone", origin, cloneLoc)
	clone.Stderr = os.Stderr
	clone.Stdout = os.Stdout
	return cloneLoc, clone.Run()
}

func checkoutAPISpecsVer(target, repoLocation string) error {
	checkout := exec.Command("git", "checkout", target)
	checkout.Stdout = os.Stdout
	checkout.Stderr = os.Stderr
	checkout.Dir = repoLocation
	return checkout.Run()
}

// getDefaultOutputLocation returns the path to the local enlistment of the Azure SDK for Go.
// If there is no local enlistment, it creates a temporary directory for the output.
func getDefaultOutputLocation() string {
	goPath := os.Getenv("GOPATH")

	if goPath != "" {
		sdkLocation := path.Join(goPath, "src", "github.com", "Azure", "azure-sdk-for-go")
		if isGitDir(sdkLocation) {
			return filepath.Clean(path.Join(sdkLocation, "arm"))
		}
	}

	if tempDir, err := ioutil.TempDir("", "azure-sdk-for-go-arm"); err == nil {
		return filepath.Clean(tempDir)
	}
	return "./"
}

func isGitDir(dir string) bool {
	retval := false
	if children, err := ioutil.ReadDir(dir); err == nil {
		for _, child := range children {
			if child.IsDir() && child.Name() == ".git" {
				retval = true
				break
			}
		}
	}
	return retval
}

// getSwaggers dives through the entire Azure/azure-rest-api-specs repository and
// publishes a list of swaggers that should be generated by autorest.
func getSwaggers(dir string, statusLog *log.Logger, errLog *log.Logger) <-chan Swagger {
	retval := make(chan Swagger)

	go func() {
		defer close(retval)

		seen := map[string][]string{}

		seenContains := func(needle Swagger) bool {
			if previouslySeen, ok := seen[needle.Info.Title]; ok {
				for _, version := range previouslySeen {
					if version == needle.Info.Version {
						return true
					}
				}
			}
			return false
		}

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) (result error) {
			if err != nil {
				return
			}

			if strings.ToLower(filepath.Ext(path)) == ".json" {
				var contents []byte
				if temp, err := ioutil.ReadFile(path); err == nil {
					contents = temp
				} else {
					return
				}

				var manifest Swagger
				if err := json.Unmarshal(contents, &manifest); err != nil {
					return
				}
				manifest.Path = path

				title := manifest.Info.Title

				if title == "" {
					return
				}

				if seenContains(manifest) {
					return
				} else if versions, ok := seen[manifest.Info.Title]; ok {
					seen[manifest.Info.Title] = append(versions, manifest.Info.Version)
				} else {
					seen[manifest.Info.Title] = []string{manifest.Info.Version}
				}

				retval <- manifest
			}
			return
		})
	}()

	return retval
}

func generate(swag Swagger, outputRootPath string) error {

	namespace, err := getNamespace(swag)
	if err != nil {
		return err
	}
	finalOutputDir := path.Clean(filepath.Join(outputRootPath, namespace))

	autorest := exec.Command(
		"autorest",
		"-Input", swag.Path,
		"-CodeGenerator", "Go",
		"-Header", "MICROSOFT_APACHE",
		"-Namespace", namespace,
		"-OutputDirectory", finalOutputDir,
		"-Modeler", "Swagger",
		"-pv",
		"-SkipValidation")
	return autorest.Run()
}

// generateAll takes a channel of swaggers, generates them, and returns a channel of
// generated paths to the package created.
func generateAll(swaggers <-chan Swagger, repoPath, outputRootPath string) <-chan string {
	retval := make(chan string)

	go func() {
		defer close(retval)
		for swagger := range swaggers {
			generate(swagger, outputRootPath)
		}
	}()

	return retval
}

// getNamespace takes a swagger
var getNamespace = func() func(Swagger) (string, error) {
	baseNamespace := []string{"github.com", "Azure", "azure-sdk-for-go"}
	namespacePattern := regexp.MustCompile(`(?P<plane>\w+)-(?P<package>.*)[/\\](?P<version>\d{4}-\d{2}-\d{2}[\w\d\-\.]*)[/\\]swagger[/\\].*\.json`)

	return func(swag Swagger) (string, error) {
		results := namespacePattern.FindAllStringSubmatch(swag.Path, -1)
		if len(results) == 0 {
			return "", fmt.Errorf("%s is not in a recognized namespace format", swag.Path)
		}
		plane := results[0][1]
		pkg := results[0][2]
		version := results[0][3]
		namespace := append(baseNamespace, plane, pkg, version)

		return strings.Replace(strings.Join(namespace, "/"), `\`, "/", -1), nil
	}
}()
