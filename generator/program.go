package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

const (
	defaultRemoteAzureRestAPISpecsPath = "https://github.com/Azure/azure-rest-api-specs.git"
	defaultAzureRESTAPIBranch          = "master"
)

// ExitCode
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
)

func init() {
	flag.StringVar(&azureRestAPIBranch, "branch", defaultAzureRESTAPIBranch, "The branch, tag, or SHA1 identifier in the Azure Rest API Specs repository to use during API generation.")
	flag.StringVar(&remoteAzureRestAPISpecsPath, "repo", defaultRemoteAzureRestAPISpecsPath, "The path to the location of the Azure REST API Specs repository that should be used for generation.")
	flag.StringVar(&outputLocation, "output", getDefaultOutputLocation(), "a directory in which all output should be placed.")
	flag.BoolVar(&dryRun, "preview", false, "Use this flag to print a list of packages that would be generated instead of actually generating the new sdk.")
	flag.BoolVar(&help, "help", false, "Provide this flag to print usage information instead of running the program.")
	flag.BoolVar(&noClone, "noClone", false, "Use this flag to prevent cloning a new copy of Azure REST API Specs. The existing enlistment should be used instead.")

	flag.Parse()

	if help {
		return
	}

	optionalTools := []string{"gofmt", "golint"}
	requiredTools := []string{"autorest", "git", "gulp"}

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
	exitStatus := ExitCodeOkay
	defer os.Exit(exitStatus)

	if help {
		flag.Usage()
		return
	}

	if anyMissing {
		exitStatus = ExitCodeMissingRequirements
		return
	}

	if noClone == false {
		repoLoc, err := cloneAPISpecs(remoteAzureRestAPISpecsPath, localAzureRestAPISpecsPath)
		if err == nil {
			defer func() {
				if err := os.RemoveAll(repoLoc); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}()
		} else {
			fmt.Fprintln(os.Stderr, err)
			exitStatus = ExitCodeFailedToClone
			return
		}
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

func getSwaggers(dir string) <-chan string {
	retval := make(chan string)

	go func() {
		defer close(retval)
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) (result error) {
			if err != nil {
				return
			}

			if filepath.Ext(path) != "json" {
				return
			}

			retval <- path
			return
		})
	}()

	return retval
}

func generate(swagger string) {
	autorest := exec.Command(
		"autorest",
		"-Input", swagger,
		"-CodeGenerator", "Go",
		"-Header", "MICROSOFT_APACHE",
		// "-Namespace", foo,
		// "-OutputDirectory", bar,
		"-Modeler", "Swagger",
		"-pv",
		"-SkipValidation")
	autorest.Run()
}

// generateAll takes a channel of swaggers, generates them, and returns a channel of
// generated paths to the package created.
func generateAll(swaggers <-chan string) <-chan string {
	retval := make(chan string)

	go func() {
		defer close(retval)
		for swagger := range swaggers {
			generate(swagger)
		}
	}()

	return retval
}
