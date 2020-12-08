package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/tools/generator/pipeline"
	"github.com/Azure/azure-sdk-for-go/tools/internal/ioext"
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
	log.Printf("Backuping azure-sdk-for-go to temp directory...")
	backupRoot, err := backupSDKRepository(cwd)
	if err != nil {
		return err
	}
	defer eraseBackup(backupRoot)
	log.Printf("Finished backuping to '%s'", backupRoot)
	ctx := generateContext{
		sdkRoot:    normalizePath(cwd),
		backupRoot: backupRoot,
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

func backupSDKRepository(sdk string) (string, error) {
	tempRepoDir := filepath.Join(tempDir(), fmt.Sprintf("generator-%v", time.Now().Unix()))
	if err := ioext.CopyDir(sdk, tempRepoDir); err != nil {
		return "", fmt.Errorf("failed to backup azure-sdk-for-go to '%s': %+v", tempRepoDir, err)
	}
	return tempRepoDir, nil
}

func eraseBackup(tempDir string) error {
	return os.RemoveAll(tempDir)
}

func tempDir() string {
	if dir := os.Getenv("TMP_DIR"); dir != "" {
		return dir
	}
	return os.TempDir()
}

type generateContext struct {
	sdkRoot    string
	backupRoot string
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
		for tag, metadata := range metadataMap {
			// first validate the output folder is valid
			outputFolder := normalizePath(filepath.Clean(metadata.PackagePath()))
			if !strings.HasPrefix(outputFolder, ctx.sdkRoot) {
				// TODO -- we might need to record this result, and throw an error when the script ends, because this usually means the output-folder is not configured correctly
				log.Printf("[WARNING] Output folder '%s' of tag '%s' is not under root of azure-sdk-for-go, skipping", outputFolder, tag)
				continue
			}
			// first format the package
			if err := autorest.FormatPackage(outputFolder); err != nil {
				return nil, err
			}
			// get the package path - which is a relative path to the sdk root
			packagePath, err := filepath.Rel(ctx.sdkRoot, outputFolder)
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
			exporter := changelog.Exporter{
				SDKRoot:    ctx.sdkRoot,
				BackupRoot: ctx.backupRoot,
			}
			c, err := exporter.ExportForPackage(p)
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
