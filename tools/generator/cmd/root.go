package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/tools/generator/model"
	"github.com/spf13/cobra"
)

const (
	defaultOptionPath = "generate_options.json"
)

// Command ...
func Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "generator <generate input filepath> <generate output filepath>",
		Args: cobra.ExactArgs(2),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFlags(0) // remove the time stamp prefix
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			optionPath, err := cmd.Flags().GetString("options")
			if err != nil {
				return err
			}
			return execute(args[0], args[1], Flags{
				OptionPath: optionPath,
			})
		},
		SilenceUsage: true, // this command is used for a pipeline, the usage should never show
	}

	flags := rootCmd.Flags()
	flags.String("options", defaultOptionPath, "Specify a file with the autorest options")

	return rootCmd
}

// Flags ...
type Flags struct {
	OptionPath string
}

func execute(inputPath, outputPath string, flags Flags) error {
	log.Printf("Reading generate input file from '%s'...", inputPath)
	input, err := readInputFrom(inputPath)
	if err != nil {
		return fmt.Errorf("cannot read generate input: %+v", err)
	}
	log.Printf("Generating using the following GenerateInput...\n%s", input.String())
	output, err := generate(input, flags.OptionPath)
	if err != nil {
		return fmt.Errorf("cannot generate: %+v", err)
	}
	log.Printf("Output generated: \n%s", output.String())
	log.Printf("Writing output to file '%s'...", outputPath)
	if err := writeOutputTo(outputPath, output); err != nil {
		return fmt.Errorf("cannot write generate output: %+v", err)
	}
	return nil
}

func readInputFrom(inputPath string) (*model.GenerateInput, error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	return model.NewGenerateInputFrom(inputFile)
}

func writeOutputTo(outputPath string, output *model.GenerateOutput) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := output.WriteTo(file); err != nil {
		return err
	}
	return nil
}

// TODO -- support dry run
func generate(input *model.GenerateInput, optionPath string) (*model.GenerateOutput, error) {
	if input.DryRun {
		return nil, fmt.Errorf("dry run not supported yet")
	}
	log.Printf("Reading options from file '%s'...", optionPath)

	optionFile, err := os.Open(optionPath)
	if err != nil {
		return nil, err
	}

	options, err := autorest.NewOptionsFrom(optionFile)
	if err != nil {
		return nil, err
	}
	log.Printf("Autorest options: \n%s", options.String())

	// iterate over all the readme
	results := make([]model.PackageResult, 0)
	for _, readme := range input.RelatedReadmeMdFiles {
		log.Printf("Processing readme '%s'...", readme)
		task := autorest.Task{
			AbsReadmeMd: filepath.Join(input.SpecFolder, readme),
		}
		if err := task.Execute(*options); err != nil {
			return nil, err
		}
		// get changed file list
		changedFiles, err := getChangedFiles()
		if err != nil {
			return nil, err
		}
		log.Printf("Files changed in the SDK: %+v", changedFiles)
		// get packages using the changed file list
		// returns a map, key is package path, value is files that have changed
		packages, err := autorest.GetChangedPackages(changedFiles)
		if err != nil {
			return nil, err
		}
		log.Printf("Packages changed: %+v", packages)
		// iterate over the changed packages
		for p, files := range packages {
			log.Printf("Getting package result for package '%s', changed files are: [%s]", p, strings.Join(files, ", "))
			c, err := changelog.NewChangelogForPackage(p)
			if err != nil {
				return nil, err
			}
			content := c.ToMarkdown()
			breaking := c.HasBreakingChanges()
			results = append(results, model.PackageResult{
				PackageName: getPackageIdentifier(p),
				Path:        []string{p},
				ReadmeMd:    []string{readme},
				Changelog: &model.Changelog{
					Content:           &content,
					HasBreakingChange: &breaking,
				},
			})
		}
	}

	// sort results
	sort.SliceStable(results, func(i, j int) bool {
		apiI := getPackageAPIVersionSegment(results[i].PackageName)
		apiJ := getPackageAPIVersionSegment(results[j].PackageName)
		return apiI > apiJ
	})

	return &model.GenerateOutput{
		Packages: results,
	}, nil
}

func getChangedFiles() ([]string, error) {
	var files []string
	// get the file changed
	changed, err := getDiffFiles()
	if err != nil {
		return nil, err
	}
	files = append(files, changed...)
	// get the untracked files
	untracked, err := getUntrackedFiles()
	if err != nil {
		return nil, err
	}
	files = append(files, untracked...)
	return files, nil
}

func getDiffFiles() ([]string, error) {
	c := exec.Command("git", "diff", "--name-only")
	output, err := c.Output()
	if err != nil {
		return nil, err
	}
	var files []string
	for _, f := range strings.Split(string(output), "\n") {
		f = strings.TrimSpace(f)
		if f != "" {
			files = append(files, f)
		}
	}
	return files, nil
}

func getUntrackedFiles() ([]string, error) {
	c := exec.Command("git", "ls-files", "--other", "--exclude-standard")
	output, err := c.Output()
	if err != nil {
		return nil, err
	}
	var files []string
	for _, f := range strings.Split(string(output), "\n") {
		f = strings.TrimSpace(f)
		if f != "" {
			files = append(files, f)
		}
	}
	return files, nil
}

func getPackageIdentifier(pkg string) string {
	return strings.TrimPrefix(pkg, "services/")
}

func getPackageAPIVersionSegment(pkg string) string {
	segments := strings.Split(pkg, "/")
	return segments[len(segments)-2]
}
