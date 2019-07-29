package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/repo"
	"github.com/spf13/cobra"
)

const (
	latest              = "latest"
	master              = "master"
	specUpstream        = "origin"
	branchPattern       = "major-version-release-v%d.0.0"
	autorestArgsPattern = "--use=@microsoft.azure/autorest.go@~2.1.99 %s --go --multiapi --go-sdk-folder=%s --use-onever"
)

// flags
var upstream string
var quietFlag bool
var debugFlag bool
var verboseFlag bool
var skip []string
var batch int

// global variables
var initialBranch string
var initialDir string
var pattern *regexp.Regexp
var majorVersion int
var absolutePathOfSDK string
var absolutePathOfSpecs string

var rootCmd = &cobra.Command{
	Use:   "major-updater <SDK dir> <specs dir>",
	Short: "Do a whole procedure of monthly regular major update",
	Long:  `This tool will execute a procedure of releasing a new major update of the azure-sdk-for-go`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := theCommand(args)
		restoreDirAndBranch()
		return err
	},
}

func init() {
	pattern = regexp.MustCompile(`^v([0-9]+)\..*$`)
	rootCmd.PersistentFlags().StringVar(&upstream, "upstream", "origin", "specify the upstream of the SDK repo")
	rootCmd.PersistentFlags().IntVarP(&batch, "batch", "b", 4, "batch count")
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "quiet output")
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "debug output")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringArrayVarP(&skip, "skip", "s", []string{}, "skiped procedures")
}

// Execute executes the specified command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func theCommand(args []string) error {
	sdkDir := args[0]
	specsDir := args[1]
	var err error
	initialDir, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get the current working directory: %v", err)
	}
	err = changeDir(sdkDir)
	absolutePathOfSDK, _ = os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot change dir to SDK folder: %v", err)
	}
	if !contains(skip, "dep") {
		println("Executing dep")
		err = executeDep()
		if err != nil {
			return err
		}
	}
	if !contains(skip, "sdk") {
		println("Update SDK repo...")
		err = updateSDKRepo()
		if err != nil {
			return err
		}
	}
	if !contains(skip, "spec") && !contains(skip, "specs") {
		println("Update specs repo...")
		err = updateSpecsRepo(specsDir)
		if err != nil {
			return err
		}
	}
	err = runAutorest()
	return nil
}

func executeDep() error {
	args := "ensure -update"
	if verboseFlag {
		args = args + " -v"
	}
	c := exec.Command("dep", strings.Split(args, " ")...)
	err := startCmd(c)
	if err != nil {
		return err
	}
	return c.Wait()
}

func updateSDKRepo() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get the current working directory: %v", err)
	}
	wt, err := repo.Get(cwd)
	if err != nil {
		return fmt.Errorf("failed to get the working tree: %v", err)
	}
	initialBranch, err = wt.Branch()
	if err != nil {
		return fmt.Errorf("failed to get the current branch: %v", err)
	}
	currentMajorVersion, err := findNextMajorVersionNumber(wt)
	majorVersion = currentMajorVersion + 1
	printf("Next major version: %d\n", majorVersion)
	vprintf("Checking out to latest branch in %s\n", cwd)
	err = wt.Checkout(latest)
	if err != nil {
		return fmt.Errorf("checkout failed: %v", err)
	}
	err = wt.Pull(upstream, latest)
	if err != nil {
		return fmt.Errorf("pull failed: %v", err)
	}
	printf("Checking out to new branch based on %s", latest)
	err = createNewBranch(wt)
	return err
}

func updateSpecsRepo(cwd string) error {
	err := changeDir(cwd)
	absolutePathOfSpecs, _ = os.Getwd()
	wt, err := repo.Get(cwd)
	if err != nil {
		return fmt.Errorf("failed to get the working tree: %v", err)
	}
	vprintf("Checking out to master branch in %s\n", cwd)
	err = wt.Checkout(master)
	if err != nil {
		return fmt.Errorf("checkout failed: %v", err)
	}
	err = wt.Pull(specUpstream, master)
	if err != nil {
		return fmt.Errorf("pull failed: %v", err)
	}
	return nil
}

func runAutorest() error {
	err := os.Setenv("NODE_OPTIONS", "--max-old-space-size=8192")
	if err != nil {
		return fmt.Errorf("failed to set environment variable: %v", err)
	}
	// call runAutorestSingle here in a loop
	return nil
}

func runAutorestSingle(filename string) error {
	autorestArgs := fmt.Sprintf(autorestArgsPattern, filename, absolutePathOfSDK)
	c := exec.Command("autorest", strings.Split(autorestArgs, " ")...)
	err := c.Start()
	if err != nil {
		return fmt.Errorf("failed to start autorest on %s: %v", filename, err)
	}
	return nil
}

func createNewBranch(wt repo.WorkingTree) error {
	branchName := fmt.Sprintf(branchPattern, majorVersion)
	vprintf("creating branch %s\n", branchName)
	err := wt.CreateAndCheckout(branchName)
	return err
}

func restoreDirAndBranch() {
	if err := changeDir(initialDir); err != nil {
		return
	}
	if wt, err := repo.Get(initialDir); err == nil {
		wt.Checkout(initialBranch)
	}
}

func changeDir(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("failed to change directory to %s: %v", path, err)
	}
	return nil
}

func findNextMajorVersionNumber(wt repo.WorkingTree) (int, error) {
	tags, err := wt.ListTags("v*")
	if err != nil {
		return 0, fmt.Errorf("failed to list tags: %v", err)
	}
	number := 0
	for _, tag := range tags {
		matches := pattern.FindStringSubmatch(tag)
		cc, err := strconv.ParseInt(matches[1], 10, 16)
		c := int(cc)
		if err == nil && c > number {
			number = c
		}
	}
	return number, nil
}
