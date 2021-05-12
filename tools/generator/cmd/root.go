// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/pipeline"
	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/tools/internal/ioext"
	"github.com/Azure/azure-sdk-for-go/tools/internal/packages/track1"
	"github.com/Azure/azure-sdk-for-go/tools/internal/utils"
	"github.com/Azure/azure-sdk-for-go/version"
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

	// we no longer need to back up the repo
	//log.Printf("Backuping azure-sdk-for-go to temp directory...")
	//backupRoot, err := backupSDKRepository(cwd)
	//if err != nil {
	//	return err
	//}
	//defer eraseBackup(backupRoot)
	//log.Printf("Finished backuping to '%s'", backupRoot)

	ctx := generateContext{
		sdkRoot:    utils.NormalizePath(cwd),
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

	repoContent map[string]exports.Content
}

// TODO -- support dry run
func (ctx generateContext) generate(input *pipeline.GenerateInput) (*pipeline.GenerateOutput, error) {
	if input.DryRun {
		return nil, fmt.Errorf("dry run not supported yet")
	}

	log.Printf("Reading packages in azure-sdk-for-go...")
	if err := ctx.readRepoContent(); err != nil {
		return nil, err
	}

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
	for _, p := range removedPackages.packages() {
		removedPackagePaths = append(removedPackagePaths, p.outputFolder)
	}
	log.Printf("The following %d package(s) have been cleaned up: [%s]", len(removedPackagePaths), strings.Join(removedPackagePaths, ", "))

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
	errorBuilder := generateErrorBuilder{}
	for _, readme := range input.RelatedReadmeMdFiles {
		log.Printf("Processing readme '%s'...", readme)
		absReadme := filepath.Join(input.SpecFolder, readme)
		metadataOutput := filepath.Dir(absReadme)
		// generate code
		g := autorestContext{
			generator: autorest.NewGeneratorFromOptions(options).WithReadme(absReadme).WithMetadataOutput(metadataOutput),
		}
		if err := g.generate(); err != nil {
			errorBuilder.add(fmt.Errorf("cannot generate readme '%s': %+v", readme, err))
			continue
		}
		m := changelogContext{
			sdkRoot:         ctx.sdkRoot,
			readme:          readme,
			removedPackages: removedPackages[readme],
			commonMetadata: autorest.GenerationMetadata{
				CommitHash:     ctx.commitHash,
				Readme:         autorest.NormalizedSpecRoot + utils.NormalizePath(readme),
				CodeGenVersion: options.CodeGeneratorVersion(),
				RepositoryURL:  "https://github.com/Azure/azure-rest-api-specs.git",
			},
			repoContent:       ctx.repoContent,
			autorestArguments: g.autorestArguments(),
		}
		log.Printf("Processing metadata generated in readme '%s'...", readme)
		packages, err := m.process(metadataOutput)
		if err != nil {
			errorBuilder.add(fmt.Errorf("cannot process metadata for readme '%s': %+v", readme, err))
			continue
		}

		// iterate over the changed packages
		set := packageResultSet{}
		for _, p := range packages {
			log.Printf("Getting package result for package '%s'", p.PackageName)
			content := p.Changelog.ToCompactMarkdown()
			breaking := p.Changelog.HasBreakingChanges()
			breakingChangeItems := p.Changelog.GetBreakingChangeItems()
			set.add(pipeline.PackageResult{
				Version:     version.Number,
				PackageName: getPackageIdentifier(p.PackageName),
				Path:        []string{p.PackageName},
				ReadmeMd:    []string{readme},
				Changelog: &pipeline.Changelog{
					Content:             &content,
					HasBreakingChange:   &breaking,
					BreakingChangeItems: &breakingChangeItems,
				},
			})
		}
		results = append(results, set.toSlice()...)
	}

	// validate the sdk structure
	log.Printf("Validating services directory structure...")
	exceptions, err := loadExceptions(filepath.Join(ctx.sdkRoot, "tools/pkgchk/exceptions.txt"))
	if err := track1.VerifyWithDefaultVerifiers(filepath.Join(ctx.sdkRoot, "services"), exceptions); err != nil {
		return nil, err
	}

	return &pipeline.GenerateOutput{
		Packages: results,
	}, errorBuilder.build()
}

func (ctx *generateContext) readRepoContent() error {
	ctx.repoContent = make(map[string]exports.Content)
	pkgs, err := track1.List(filepath.Join(ctx.sdkRoot, "services"))
	if err != nil {
		return fmt.Errorf("failed to list track 1 packages: %+v", err)
	}

	for _, pkg := range pkgs {
		relativePath, err := filepath.Rel(ctx.sdkRoot, pkg.FullPath())
		if err != nil {
			return err
		}
		relativePath = utils.NormalizePath(relativePath)
		if _, ok := ctx.repoContent[relativePath]; ok {
			return fmt.Errorf("duplicate package: %s", pkg.Path())
		}
		exp, err := exports.Get(pkg.FullPath())
		if err != nil {
			return err
		}
		ctx.repoContent[relativePath] = exp
	}

	return nil
}

type generateErrorBuilder struct {
	errors []error
}

func (b *generateErrorBuilder) add(err error) {
	b.errors = append(b.errors, err)
}

func (b *generateErrorBuilder) build() error {
	if len(b.errors) == 0 {
		return nil
	}
	var messages []string
	for _, err := range b.errors {
		messages = append(messages, err.Error())
	}
	return fmt.Errorf("total %d error(s): \n%s", len(b.errors), strings.Join(messages, "\n"))
}

type packageResultSet map[string]pipeline.PackageResult

func (s *packageResultSet) contains(r pipeline.PackageResult) bool {
	_, ok := (*s)[r.PackageName]
	return ok
}

func (s *packageResultSet) add(r pipeline.PackageResult) {
	if s.contains(r) {
		log.Printf("[WARNING] The result set already contains key %s with value %+v, but we are still trying to insert a new value %+v on the same key", r.PackageName, (*s)[r.PackageName], r)
	}
	(*s)[r.PackageName] = r
}

func (s *packageResultSet) toSlice() []pipeline.PackageResult {
	results := make([]pipeline.PackageResult, 0)
	for _, r := range *s {
		results = append(results, r)
	}
	// sort the results
	sort.SliceStable(results, func(i, j int) bool {
		// we first clip the preview segment and then sort by string literal
		pI := strings.Replace(results[i].PackageName, "preview/", "/", 1)
		pJ := strings.Replace(results[j].PackageName, "preview/", "/", 1)
		return pI > pJ
	})
	return results
}

func getPackageIdentifier(pkg string) string {
	return strings.TrimPrefix(utils.NormalizePath(pkg), "services/")
}

func loadExceptions(exceptFile string) (map[string]bool, error) {
	if exceptFile == "" {
		return nil, nil
	}
	f, err := os.Open(exceptFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	exceptions := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		exceptions[scanner.Text()] = true
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return exceptions, nil
}
