package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/pipeline"
	"github.com/Azure/azure-sdk-for-go/tools/generator/utils"
	"github.com/Azure/azure-sdk-for-go/tools/internal/ioext"
	"github.com/spf13/cobra"
)

const (
	defaultOptionPath = "generate_options.json"
)

// Command returns the command for the generator. Note that this command is designed to run in the root directory of
// azure-sdk-for-go. It does not work if you are running this tool in somewhere else
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
	input, err := pipeline.ReadInput(inputPath)
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
		sdkRoot:    utils.NormalizePath(cwd),
		clnRoot:    backupRoot,
		specRoot:   input.SpecFolder,
		commitHash: input.HeadSha,
		optionPath: flags.OptionPath,
	}
	output, err := ctx.generate(input)
	if err != nil {
		return err
	}
	log.Printf("Output generated: \n%s", output.String())
	log.Printf("Writing output to file '%s'...", outputPath)
	if err := pipeline.WriteOutput(outputPath, output); err != nil {
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
	clnRoot    string
	specRoot   string
	commitHash string
	optionPath string
}

// TODO -- support dry run
func (ctx generateContext) generate(input *pipeline.GenerateInput) (*pipeline.GenerateOutput, error) {
	if input.DryRun {
		return nil, fmt.Errorf("dry run not supported yet")
	}
	log.Printf("Reading options from file '%s'...", ctx.optionPath)

	// now we summary all the metadata in sdk
	log.Printf("Cleaning up all the packages related with the following readme files: [%s]", strings.Join(input.RelatedReadmeMdFiles, ", "))
	cleanUpCtx := cleanUpContext{
		root:        filepath.Join(ctx.sdkRoot, "services"),
		readmeFiles: input.RelatedReadmeMdFiles,
	}
	removedPackages, err := cleanUpCtx.clean()
	if err != nil {
		return nil, err
	}
	var removedPackagePaths []string
	for _, p := range removedPackages {
		removedPackagePaths = append(removedPackagePaths, p.outputFolder)
	}
	log.Printf("The following %d package(s) have been cleaned up: [%s]", len(removedPackagePaths), strings.Join(removedPackagePaths, ", "))

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
		m := changelogContext{
			sdkRoot:         ctx.sdkRoot,
			clnRoot:         ctx.clnRoot,
			specRoot:        ctx.specRoot,
			commitHash:      ctx.commitHash,
			codeGenVer:      options.CodeGeneratorVersion(),
			readme:          readme,
			removedPackages: removedPackages,
		}
		log.Printf("Processing metadata generated in readme '%s'...", readme)
		packages, err := m.process(g.metadataOutput)
		if err != nil {
			return nil, err
		}

		// iterate over the changed packages
		for _, p := range packages {
			log.Printf("Getting package result for package '%s'", p.PackageName)
			content := p.Changelog.ToCompactMarkdown()
			breaking := p.Changelog.HasBreakingChanges()
			results = append(results, pipeline.PackageResult{
				PackageName: getPackageIdentifier(p.PackageName),
				Path:        []string{p.PackageName},
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

func getPackageIdentifier(pkg string) string {
	return strings.TrimPrefix(utils.NormalizePath(pkg), "services/")
}

func getPackageAPIVersionSegment(pkg string) string {
	segments := strings.Split(pkg, "/")
	return segments[len(segments)-2]
}
