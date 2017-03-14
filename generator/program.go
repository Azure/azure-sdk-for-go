package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

const (
	defaultAzureRestAPISpecsPath = "https://github.com/Azure/azure-rest-api-specs.git"
	defaultAzureRESTAPIBranch    = "master"
)

var (
	azureRestAPISpecsPath string
	azureRestAPIBranch    string
	outputLocation        string
)

func init() {
	flag.StringVar(&azureRestAPIBranch, "branch", defaultAzureRESTAPIBranch, "The branch, tag, or SHA1 identifier in the Azure Rest API Specs repository to use during API generation.")
	flag.StringVar(&azureRestAPISpecsPath, "repo", defaultAzureRestAPISpecsPath, "The path to the location of the Azure REST API Specs repository that should be used for generation.")
	flag.StringVar(&outputLocation, "output", getDefaultOutputLocation(), "a directory in which all output should be placed.")
	flag.Parse()

	optionalTools := []string{"gofmt", "golint"}
	requiredTools := []string{"autorest", "git", "gulp"}

	for _, tool := range optionalTools {
		if _, err := exec.LookPath(tool); err != nil {
			log.Printf("WARNING: Could not find \"%s\" usage of this tool will be skipped.", tool)
		}
	}

	anyMissing := false
	for _, tool := range requiredTools {
		if _, err := exec.LookPath(tool); err != nil {
			log.Printf("ERROR: Could not find \"%s\" this tool will exit without executing.", tool)
			anyMissing = true
		}
	}

	if anyMissing {
		os.Exit(1)
	}
}

func main() {

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

// getDefaultOutputLocation returns the path to the local enlistment of the Azure SDK for Go.
// If there is no local enlistment, it creates a temporary directory for the output.
func getDefaultOutputLocation() string {
	goPath := os.Getenv("GOPATH")

	if goPath != "" {
		sdkLocation := path.Join(goPath, "src", "github.com", "Azure", "azure-sdk-for-go")
		if isGitDir(sdkLocation) {
			return path.Join(sdkLocation, "arm")
		}
	}

	if tempDir, err := ioutil.TempDir("", "azure-sdk-for-go-arm"); err == nil {
		return tempDir
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
