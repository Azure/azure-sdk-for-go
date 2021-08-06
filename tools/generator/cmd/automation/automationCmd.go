// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package automation

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/automation/pipeline"
	"github.com/Azure/azure-sdk-for-go/tools/generator/common"
	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/tools/internal/packages/track1"
	"github.com/Azure/azure-sdk-for-go/tools/internal/utils"
	"github.com/spf13/cobra"
)

// Command returns the automation command. Note that this command is designed to run in the root directory of
// azure-sdk-for-go. It does not work if you are running this tool in somewhere else
func Command() *cobra.Command {
	automationCmd := &cobra.Command{
		Use:  "automation <generate input filepath> <generate output filepath>",
		Args: cobra.ExactArgs(2),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFlags(0) // remove the time stamp prefix
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			optionPath, err := cmd.Flags().GetString("options")
			if err != nil {
				logError(err)
				return err
			}
			if err := execute(args[0], args[1], Flags{
				OptionPath: optionPath,
			}); err != nil {
				logError(err)
				return err
			}
			return nil
		},
		SilenceUsage: true, // this command is used for a pipeline, the usage should never show
	}

	flags := automationCmd.Flags()
	flags.String("options", common.DefaultOptionPath, "Specify a file with the autorest options")

	return automationCmd
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
	log.Printf("Using current directory as SDK root: %s", cwd)

	ctx := automationContext{
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

type automationContext struct {
	sdkRoot    string
	specRoot   string
	commitHash string
	optionPath string

	repoContent map[string]exports.Content

	sdkVersion string

	existingPackages existingPackageMap

	defaultOptions    model.Options
	additionalOptions []model.Option
}

func (ctx *automationContext) categorizePackages() error {
	ctx.existingPackages = existingPackageMap{}

	serviceRoot := filepath.Join(ctx.sdkRoot, "services")
	m, err := autorest.CollectGenerationMetadata(serviceRoot)
	if err != nil {
		return err
	}

	for path, metadata := range m {
		// the path in the metadata map is the absolute path
		relPath, err := filepath.Rel(ctx.sdkRoot, path)
		if err != nil {
			return err
		}
		ctx.existingPackages.add(utils.NormalizePath(relPath), metadata)
	}

	return nil
}

func (ctx *automationContext) readDefaultOptions() error {
	log.Printf("Reading defaultOptions from file '%s'...", ctx.optionPath)
	optionFile, err := os.Open(ctx.optionPath)
	if err != nil {
		return err
	}

	generateOptions, err := model.NewGenerateOptionsFrom(optionFile)
	if err != nil {
		return err
	}

	// parsing the default options
	defaultOptions, err := model.ParseOptions(generateOptions.AutorestArguments)
	if err != nil {
		return fmt.Errorf("cannot parse default options from %v: %+v", generateOptions.AutorestArguments, err)
	}

	// remove the `--multiapi` in default options
	var options []model.Option
	for _, o := range defaultOptions.Arguments() {
		if v, ok := o.(model.FlagOption); ok && v.Flag() == "multiapi" {
			continue
		}
		options = append(options, o)
	}

	ctx.defaultOptions = model.NewOptions(options...)
	log.Printf("Autorest defaultOptions: \n%+v", ctx.defaultOptions.Arguments())

	// parsing the additional options
	additionalOptions, err := model.ParseOptions(generateOptions.AdditionalOptions)
	if err != nil {
		return fmt.Errorf("cannot parse additional options from %v: %+v", generateOptions.AdditionalOptions, err)
	}
	ctx.additionalOptions = additionalOptions.Arguments()

	return nil
}

// TODO -- support dry run
func (ctx *automationContext) generate(input *pipeline.GenerateInput) (*pipeline.GenerateOutput, error) {
	if input.DryRun {
		return nil, fmt.Errorf("dry run not supported yet")
	}

	log.Printf("Reading packages in azure-sdk-for-go...")
	if err := ctx.readRepoContent(); err != nil {
		return nil, err
	}

	log.Printf("Reading metadata information in azure-sdk-for-go...")
	if err := ctx.categorizePackages(); err != nil {
		return nil, err
	}

	log.Printf("Reading default options...")
	if err := ctx.readDefaultOptions(); err != nil {
		return nil, err
	}

	log.Printf("Reading version number...")
	if err := ctx.readVersion(); err != nil {
		return nil, err
	}

	// iterate over all the readme
	results := make([]pipeline.PackageResult, 0)
	errorBuilder := generateErrorBuilder{}
	for _, readme := range input.RelatedReadmeMdFiles {
		generateCtx := generateContext{
			sdkRoot:          ctx.sdkRoot,
			specRoot:         ctx.specRoot,
			commitHash:       ctx.commitHash,
			repoContent:      ctx.repoContent,
			existingPackages: ctx.existingPackages[readme],
			defaultOptions:   ctx.defaultOptions,
		}

		packageResults, errors := generateCtx.generate(readme)
		if len(errors) != 0 {
			errorBuilder.add(errors...)
			continue
		}

		// iterate over the changed packages
		set := packageResultSet{}
		for _, p := range packageResults {
			log.Printf("Getting package result for package '%s'", p.Package.PackageName)
			content := p.Package.Changelog.ToCompactMarkdown()
			breaking := p.Package.Changelog.HasBreakingChanges()
			breakingChangeItems := p.Package.Changelog.GetBreakingChangeItems()
			set.add(pipeline.PackageResult{
				Version:     ctx.sdkVersion,
				PackageName: getPackageIdentifier(p.Package.PackageName),
				Path:        []string{p.Package.PackageName},
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
	if err != nil {
		return nil, err
	}
	if err := track1.VerifyWithDefaultVerifiers(filepath.Join(ctx.sdkRoot, "services"), exceptions); err != nil {
		return nil, err
	}

	return &pipeline.GenerateOutput{
		Packages: squashResults(results),
	}, errorBuilder.build()
}

// squashResults squashes the package results by appending all of the `path`s in the following items to the first item
// By doing this, the SDK automation pipeline will only create one PR that contains all of the generation results
// instead of creating one PR for each generation result.
// This is to reduce the resource cost on GitHub
func squashResults(packages []pipeline.PackageResult) []pipeline.PackageResult {
	if len(packages) == 0 {
		return packages
	}
	for i := 1; i < len(packages); i++ {
		// append the path of the i-th item to the first
		packages[0].Path = append(packages[0].Path, packages[i].Path...)
		// erase the path on the i-th item
		packages[i].Path = make([]string, 0)
	}

	return packages
}

func (ctx *automationContext) readRepoContent() error {
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

func (ctx *automationContext) readVersion() error {
	v, err := ReadVersion(filepath.Join(ctx.sdkRoot, "version"))
	if err != nil {
		return err
	}
	ctx.sdkVersion = v
	return nil
}

func contains(array []autorest.GenerateResult, item string) bool {
	for _, r := range array {
		if utils.NormalizePath(r.Package.PackageName) == utils.NormalizePath(item) {
			return true
		}
	}
	return false
}

type generateErrorBuilder struct {
	errors []error
}

func (b *generateErrorBuilder) add(err ...error) {
	b.errors = append(b.errors, err...)
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

func (s packageResultSet) contains(r pipeline.PackageResult) bool {
	_, ok := s[r.PackageName]
	return ok
}

func (s packageResultSet) add(r pipeline.PackageResult) {
	if s.contains(r) {
		log.Printf("[WARNING] The result set already contains key %s with value %+v, but we are still trying to insert a new value %+v on the same key", r.PackageName, s[r.PackageName], r)
	}
	s[r.PackageName] = r
}

func (s packageResultSet) toSlice() []pipeline.PackageResult {
	results := make([]pipeline.PackageResult, 0)
	for _, r := range s {
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

func logError(err error) {
	for _, line := range strings.Split(err.Error(), "\n") {
		if l := strings.TrimSpace(line); l != "" {
			log.Printf("[ERROR] %s", l)
		}
	}
}
