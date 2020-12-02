package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/tools/generator/pipeline"
	"github.com/spf13/cobra"
)

const (
	defaultOptionPath = "generate_options.json"
)

// Command returns the command for the generator. Note that this command is designed to run in the root directory of
// azure-sdk-for-go. It might not work if you are running this tool in somewhere else
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
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	ctx := generateContext{
		cwd:        cwd,
		optionPath: flags.OptionPath,
	}
	output, err := ctx.generate(input)
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

func readInputFrom(inputPath string) (*pipeline.GenerateInput, error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	return pipeline.NewGenerateInputFrom(inputFile)
}

func writeOutputTo(outputPath string, output *pipeline.GenerateOutput) error {
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

type generateContext struct {
	cwd        string
	optionPath string
}

// TODO -- support dry run
func (ctx generateContext) generate(input *pipeline.GenerateInput) (*pipeline.GenerateOutput, error) {
	if input.DryRun {
		return nil, fmt.Errorf("dry run not supported yet")
	}
	log.Printf("Reading options from file '%s'...", ctx.optionPath)

	optionFile, err := os.Open(ctx.optionPath)
	if err != nil {
		return nil, err
	}

	options, err := model.NewOptionsFrom(optionFile)
	if err != nil {
		return nil, err
	}
	log.Printf("Autorest options: \n%+v", options)

	// iterate over all the readme
	results := make([]pipeline.PackageResult, 0)
	for _, readme := range input.RelatedReadmeMdFiles {
		// TODO -- maintain a map from readme files to corresponding output folders, so that we could detect the situation that a package was deleted
		log.Printf("Processing readme '%s'...", readme)
		absReadme := filepath.Join(input.SpecFolder, readme)
		// generate code
		g := autorestContext{
			absReadme:      absReadme,
			metadataOutput: filepath.Dir(absReadme),
			options:        options,
		}
		if err := g.generate(); err != nil {
			return nil, err
		}
		// get the metadata map
		m := autorest.NewMetadataProcessorFromLocation(g.metadataOutput)
		metadataMap, err := m.Process()
		if err != nil {
			return nil, err
		}
		var packages []string
		for _, metadata := range metadataMap {
			// TODO -- first validate the output folder is valid
			outputFolder := filepath.Clean(metadata.PackagePath())
			// first format the package
			if err := autorest.FormatPackage(outputFolder); err != nil {
				return nil, err
			}
			// get the package path - which is a relative path to the sdk root
			packagePath, err := filepath.Rel(ctx.cwd, outputFolder)
			if err != nil {
				return nil, err
			}
			packages = append(packages, packagePath)
		}
		log.Printf("Packages changed: %+v", packages)
		// iterate over the changed packages
		for _, p := range packages {
			p = normalizePath(p)
			log.Printf("Getting package result for package '%s'", p)
			c, err := changelog.NewChangelogForPackage(p)
			if err != nil {
				return nil, err
			}
			content := c.ToMarkdown()
			breaking := c.HasBreakingChanges()
			results = append(results, pipeline.PackageResult{
				PackageName: getPackageIdentifier(p),
				Path:        []string{p},
				ReadmeMd:    []string{readme},
				Changelog: &pipeline.Changelog{
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

	return &pipeline.GenerateOutput{
		Packages: results,
	}, nil
}

func normalizePath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func getPackageIdentifier(pkg string) string {
	return strings.TrimPrefix(pkg, "services/")
}

func getPackageAPIVersionSegment(pkg string) string {
	segments := strings.Split(pkg, "/")
	return segments[len(segments)-2]
}
