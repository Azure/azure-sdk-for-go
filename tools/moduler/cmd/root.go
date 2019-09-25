package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/repo"
	"github.com/spf13/cobra"
)

const (
	tagPrefix = "// tag: "
)

var rootCmd = &cobra.Command{
	Use:   "module-releaser <sdk root path>",
	Short: "Release a new module by pushing new tags to github",
	Long: `This tool search the whole SDK folder for tags in version.go files which are produced by the versioner tool, 
and push the new tags to github.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return theCommand(args)
	},
}

var (
	getTagsHook func(string) ([]string, error)
	dryRunFlag  bool
	verboseFlag bool
)

func init() {
	getTagsHook = getTags
	rootCmd.PersistentFlags().BoolVarP(&dryRunFlag, "dry-run", "d", false, "dry run, only list the detected tags, do not add or push them")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "verbose output")
}

// Execute the tool
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func theCommand(args []string) error {
	path, err := filepath.Abs(args[0])
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}
	tags, err := getTags(path)
	if err != nil {
		return fmt.Errorf("failed to list tags: %v", err)
	}
	vprintf("Get %d tags\n", len(tags))
	// find all version.go files
	files, err := findAllVersionFiles(path)
	if err != nil {
		return fmt.Errorf("failed to find all version.go files: %v", err)
	}
	// read new tags in version.go
	for _, file := range files {
		newTag, err := findNewTagInFile(file)
		if err != nil {
			return fmt.Errorf("failed to read tag in file '%s': %v", file, err)
		}
		vprintf("Found tag '%s' in file '%s'\n", newTag, file)
		if newTag == "" {
			printf("Found empty tag in file '%s'\n", file)
			continue
		}
		if tagExists(tags, newTag) {
			vprintf("Tag '%s' already exists, skipping", newTag)
			continue
		}
		// push new tag
		if err := pushNewTag(newTag); err != nil {
			return fmt.Errorf("failed to push tag '%s': %v", newTag, err)
		}
	}
	return nil
}

func tagExists(tags []string, tag string) bool {
	for _, item := range tags {
		if item == tag {
			return true
		}
	}
	return false
}

func findNewTagInFile(path string) (string, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	tag := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, tagPrefix) {
			tag = strings.ReplaceAll(line, tagPrefix, "")
			break
		}
	}
	return tag, nil
}

func findAllVersionFiles(root string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && info.Name() == "version.go" {
			files = append(files, path)
			return nil
		}
		return nil
	})
	return files, err
}

func pushNewTag(newTag string) error {
	if dryRunFlag {
		vprintf("Found new tag '%s'\n", newTag)
	} else {
		vprintf("Adding new tag '%s'\n", newTag)
		cmd := exec.Command("git", "tag", "-a", newTag, "-m", newTag)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return errors.New(string(output))
		}
		vprintf("Pushing new tag '%s'\n", newTag)
		cmd = exec.Command("git", "push", "origin", "--tags")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return errors.New(string(output))
		}
	}
	return nil
}

func getTags(repoPath string) ([]string, error) {
	wt, err := repo.Get(repoPath)
	if err != nil {
		return nil, err
	}
	return wt.ListTags("*")
}
