// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package automation

import (
	"bufio"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/automation/validate"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/exports"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/pipeline"
	"github.com/Azure/azure-sdk-for-go/tools/generator/utils"
	"github.com/Azure/azure-sdk-for-go/tools/pkgchk/track1"
	"github.com/Azure/azure-sdk-for-go/version"
	"github.com/spf13/cobra"
)

const (
	defaultOptionPath = "generate_options.json"
)

// Command returns the command for the generator. Note that this command is designed to run in the root directory of
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
				return err
			}
			return execute(args[0], args[1], Flags{
				OptionPath: optionPath,
			})
		},
		SilenceUsage: true, // this command is used for a pipeline, the usage should never show
	}

	flags := automationCmd.Flags()
	flags.String("options", defaultOptionPath, "Specify a file with the autorest options")

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

func tempDir() string {
	if dir := os.Getenv("TMP_DIR"); dir != "" {
		return dir
	}
	return os.TempDir()
}

type generateContext struct {
	sdkRoot    string
	specRoot   string
	commitHash string
	optionPath string

	repoContent map[string]exports.Content

	existingPackages existingPackageMap

	defaultOptions model.Options
}

func (ctx generateContext) SDKRoot() string {
	return ctx.sdkRoot
}

func (ctx generateContext) SpecRoot() string {
	return ctx.specRoot
}

func (ctx generateContext) RepoContent() map[string]exports.Content {
	return ctx.repoContent
}

func (ctx *generateContext) categorizePackages() error {
	ctx.existingPackages = existingPackageMap{}

	serviceRoot := filepath.Join(ctx.sdkRoot, "services")
	m, err := autorest.CollectGenerationMetadata(serviceRoot)
	if err != nil {
		return err
	}

	for path, metadata := range m {
		// the path in the metadata map is the absolute path
		ctx.existingPackages.add(path, metadata)
	}

	return nil
}

func (ctx *generateContext) readDefaultOptions() error {
	log.Printf("Reading defaultOptions from file '%s'...", ctx.optionPath)
	optionFile, err := os.Open(ctx.optionPath)
	if err != nil {
		return err
	}

	defaultOptions, err := model.NewOptionsFrom(optionFile)
	if err != nil {
		return err
	}
	log.Printf("Autorest defaultOptions: \n%+v", defaultOptions.Arguments())

	// remove the `--multiapi` in default options
	var options []model.Option
	for _, o := range defaultOptions.Arguments() {
		if v, ok := o.(model.FlagOption); ok && v.Flag() == "multiapi" {
			continue
		}
		options = append(options, o)
	}

	ctx.defaultOptions = model.NewOptions(options...)

	return nil
}

func (ctx *generateContext) generate(input *pipeline.GenerateInput) (*pipeline.GenerateOutput, error) {
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

	// iterate over all the readme
	results := make([]pipeline.PackageResult, 0)
	errorBuilder := generateErrorBuilder{}
	for _, readme := range input.RelatedReadmeMdFiles {
		absReadme := filepath.Join(input.SpecFolder, readme)
		absReadmeGo := filepath.Join(filepath.Dir(absReadme), "readme.go.md")
		log.Printf("Reading tags from readme.go.md '%s'...", absReadmeGo)
		reader, err := os.Open(absReadmeGo)
		if err != nil {
			errorBuilder.add(fmt.Errorf("cannot read from readme.go.md: %+v", err))
			continue
		}
		log.Printf("Parsing tags from readme.go.md '%s'...", absReadmeGo)
		tags, err := autorest.ReadBatchTags(reader)
		if err != nil {
			errorBuilder.add(fmt.Errorf("cannot read batch tags in readme.go.md '%s': %+v", absReadmeGo, err))
			continue
		}

		log.Printf("Cleaning all the packages from readme '%s'...", readme)
		packages := ctx.existingPackages[readme]

		removedPackages, err := clean(packages)
		if err != nil {
			errorBuilder.add(fmt.Errorf("cannot clean packages from readme '%s': %+v", readme, err))
			continue
		}

		log.Printf("Generating the following tags: \n[%s]", strings.Join(tags, ", "))
		var options model.Options
		var packageResults []autorest.GenerateResult
		for _, tag := range tags {
			// Get the proper options to use depending on whether this tag has been already generated in the SDK or not
			if metadata, ok := packages[tag]; ok {
				// this tag has been generated, use the existing parameters in its metadata
				additionalOptions, err := model.ParseOptions(strings.Split(metadata.AdditionalProperties.AdditionalOptions, " "))
				if err != nil {
					errorBuilder.add(fmt.Errorf("cannot parse existing defaultOptions for readme '%s'/tag '%s': %+v", readme, tag, err))
					continue
				}
				options = ctx.defaultOptions.MergeOptions(additionalOptions.Arguments()...)
			} else {
				// this is a new tag
				options = ctx.defaultOptions
			}

			input := autorest.GenerateInput{
				Readme:     readme,
				Tag:        tag,
				CommitHash: ctx.commitHash,
				Options:    options,
			}
			validateCtx := validate.MetadataValidateContext{
				Readme:  readme,
				SDKRoot: ctx.sdkRoot,
			}
			result, err := autorest.GeneratePackage(ctx, input, autorest.GenerateOptions{
				Stderr:            os.Stderr,
				Stdout:            os.Stderr,
				AutoRestLogPrefix: "[AUTOREST] ",
				ChangelogTitle:    "Unreleased",
				Validators: []autorest.MetadataValidateFunc{
					validateCtx.PreviewCheck,
					validateCtx.MgmtCheck,
					validateCtx.NamespaceCheck,
				},
			})
			if err != nil {
				errorBuilder.add(err)
				continue
			}
			packageResults = append(packageResults, *result)
		}

		// also add the removed packages in the results if it is not regenerated
		for _, removedPackage := range removedPackages {
			if !contains(packageResults, removedPackage.packageFullPath) {
				// this package is not regenerated, therefore it is removed
				packageResults = append(packageResults, autorest.GenerateResult{
					Package: autorest.ChangelogResult{
						Tag:             removedPackage.Tag,
						PackageFullPath: removedPackage.packageFullPath,
						Changelog: model.Changelog{
							RemovedPackage: true,
						},
					},
				})
			}
		}

		// iterate over the changed packages
		set := packageResultSet{}
		for _, p := range packageResults {
			log.Printf("Getting package result for package '%s'", p.Package.PackageName)
			content := p.Package.Changelog.ToCompactMarkdown()
			breaking := p.Package.Changelog.HasBreakingChanges()
			breakingChangeItems := p.Package.Changelog.GetBreakingChangeItems()
			set.add(pipeline.PackageResult{
				Version:     version.Number, // TODO -- after migrate this to a module, we cannot get the version number in this way anymore
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

func contains(array []autorest.GenerateResult, item string) bool {
	for _, r := range array {
		if utils.NormalizePath(r.Package.PackageFullPath) == utils.NormalizePath(item) {
			return true
		}
	}
	return false
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
