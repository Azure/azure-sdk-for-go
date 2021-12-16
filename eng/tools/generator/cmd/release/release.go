// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package release

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config/validate"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/flags"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Command() *cobra.Command {
	releaseCmd := &cobra.Command{
		Use:   "release <azure-sdk-for-go directory> <azure-rest-api-specs directory> [config json file path]",
		Short: "Generate a release of azure-sdk-for-go",
		Long: `This command will generate a new release for azure-sdk-for-go using the given configuration as JSON.
if the [config json file path] is set, the configs are read from the file specified, otherwise this command will 
read the config from stdin.

azure-sdk-for-go directory: the directory path of the azure-sdk-for-go with git control
azure-rest-api-specs directory: the directory path of the azure-rest-api-specs with git control
`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			sdkPath, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get the directory of azure-sdk-for-go: %v", err)
			}
			specPath, err := filepath.Abs(args[1])
			if err != nil {
				return fmt.Errorf("failed to get the directory of azure-rest-api-specs: %v", err)
			}
			baseContext, err := repo.NewCommandContext(sdkPath, specPath, true)
			if err != nil {
				return err
			}
			configPath := ""
			if len(args) > 2 {
				configPath = args[2]
			}
			ctx := commandContext{
				CommandContext: baseContext,
				configPath:     configPath,
				flags:          ParseFlags(cmd.Flags()),
			}
			return ctx.execute()
		},
	}

	BindFlags(releaseCmd.Flags())

	return releaseCmd
}

type Flags struct {
	common.GlobalFlags
	SkipDep        bool
	SkipValidate   bool
	SkipProfile    bool
	Major          bool
	KeepTempBranch bool
	RefreshAll     bool
	CompareCommit  string
	VersionNumber  string
}

func BindFlags(flagSet *pflag.FlagSet) {
	flagSet.Bool("skip-dep-ensure", false, "Skip the \"dep ensure\" in the afterscript for more robustness.")
	flagSet.BoolP("skip-validate", "l", false, "Skip the validate for readme files and tags.")
	flagSet.BoolP("skip-profile", "p", false, "Skip the profile regeneration.")
	flagSet.Bool("major", false, "Allow the breaking changes in stable packages in this release")
	flagSet.Bool("keep-temp-branch", false, "Keep the temp branch after the release is done")
	flagSet.Bool("refresh-all", false, "Refresh all the packages even if the `refresh` block is not specified in the configuration")
	flagSet.String("compare-commit", "", "Specify the commit all changelogs need to generate against")
	flagSet.String("version-number", "", "Specify the version number of this release")
}

func ParseFlags(flagSet *pflag.FlagSet) Flags {
	return Flags{
		GlobalFlags:    common.ParseGlobalFlags(flagSet),
		SkipDep:        flags.GetBool(flagSet, "skip-dep-ensure"),
		SkipValidate:   flags.GetBool(flagSet, "skip-validate"),
		SkipProfile:    flags.GetBool(flagSet, "skip-profile"),
		Major:          flags.GetBool(flagSet, "major"),
		KeepTempBranch: flags.GetBool(flagSet, "keep-temp-branch"),
		RefreshAll:     flags.GetBool(flagSet, "refresh-all"),
		CompareCommit:  flags.GetString(flagSet, "compare-commit"),
		VersionNumber:  flags.GetString(flagSet, "version-number"),
	}
}

type commandContext struct {
	repo.CommandContext

	configPath        string
	flags             Flags
	options           model.Options
	additionalOptions []model.Option
	sdkRef            *plumbing.Reference
	specRef           *plumbing.Reference

	repoContent repo.RepoContent
}

func (c *commandContext) execute() error {
	log.Printf("Parsing the config...")
	cfg, err := config.ParseConfig(c.configPath)
	if err != nil {
		return err
	}
	if !c.flags.SkipValidate {
		validator := validate.NewLocalValidator(c.Spec().Root())
		if err := validator.Validate(*cfg); err != nil {
			return err
		}
	}
	log.Printf("Configuration: %s", cfg.String())

	c.sdkRef, err = c.SDK().Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD ref of azure-sdk-for-go: %+v", err)
	}
	log.Printf("The release branch is based on HEAD ref '%s' (commit %s) of azure-sdk-for-go", c.sdkRef.Name(), c.sdkRef.Hash())

	// get the repo content to be compared with
	log.Printf("Reading packages in azure-sdk-for-go...")
	c.repoContent, err = c.SDK().ReportForCommit(c.flags.CompareCommit)
	if err != nil {
		return fmt.Errorf("failed to get the initial status of the SDK repository: %+v", err)
	}

	c.specRef, err = c.Spec().Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD ref of azure-rest-api-specs: %+v", err)
	}
	log.Printf("The new version is generated from HEAD ref '%s' (commit %s) of azure-rest-api-specs", c.specRef.Name(), c.specRef.Hash())

	log.Printf("Reading autorest options in SDK repository...")
	if err := c.readOptions(); err != nil {
		return err
	}

	log.Printf("Refreshing packages...")
	refreshRef, err := c.refresh(&cfg.RefreshInfo)
	if err != nil {
		return err
	}

	// create a temporary branch to hold the generation result
	log.Printf("Creating temporary branch...")
	tempBranchName, err := c.CreateReleaseBranch("temp")
	if err != nil {
		return err
	}
	log.Printf("Temporary branch '%s' created", tempBranchName)
	defer func() {
		if c.flags.KeepTempBranch {
			return
		}
		log.Printf("Deleting temp branch '%s'...", tempBranchName)
		if err := c.SDK().DeleteBranch(tempBranchName); err != nil {
			log.Printf("Error deleting temp branch '%s': %+v", tempBranchName, err)
			return
		}
		log.Printf("Temp branch '%s' deleted", tempBranchName)
	}()

	log.Printf("Reading additionalOptions...")
	c.additionalOptions, err = cfg.AdditionalOptions()
	if err != nil {
		return fmt.Errorf("failed to parse additional options from the configuration: %+v", err)
	}

	if c.flags.Major {
		log.Printf("Applying additionalOptions since we are making a major release")
		c.options = c.options.(model.Options).MergeOptions(c.additionalOptions...)
		log.Printf("The generate options: %+v", c.options)
	}

	log.Printf("Executing autorest generation...")
	results, errors := c.generate(cfg)
	// handle the errors during generation
	if len(errors) > 0 {
		var errorMessages []string
		for key, err := range errors {
			errorMessages = append(errorMessages, fmt.Sprintf("failed to execute on %s: %+v", key, err))
		}
		log.Printf("Generation has %d errors: \n%s", len(errors), strings.Join(errorMessages, "\n"))
	}
	// calculate the commits to be released
	includedPackages, excludedPackages := c.getPackagesToRelease(results)
	if c.flags.Major {
		// if we have the major flag, we add those excluded packages back
		includedPackages = append(includedPackages, excludedPackages...)
	}
	log.Printf("The following %d packages will be released: ", len(includedPackages))
	for _, r := range includedPackages {
		log.Println(r.PackageInfo())
	}

	log.Printf("Switch back to branch '%s'", refreshRef.Name())
	if err := c.SDK().Checkout(&repo.CheckoutOptions{
		Branch: refreshRef.Name(),
		Force:  true,
	}); err != nil {
		return err
	}

	log.Printf("Reading published versions...")
	version, err := c.getVersion(len(excludedPackages) != 0 && c.flags.Major)
	if err != nil {
		return err
	}
	log.Printf("Creating actual release branch...")
	if _, err := c.CreateReleaseBranch(version.NewVersion); err != nil {
		return err
	}
	log.Printf("Include the packages that is about to release in this release...")
	for _, p := range includedPackages {
		log.Printf("Including package '%s' in this release (commit hash %s)...", p.Package.PackageName, p.CommitHash)
		if err := c.SDK().CherryPick(p.CommitHash); err != nil {
			return err
		}
	}

	// generate the changelog
	log.Printf("Generating CHANGELOG...")
	r, err := repo.GetPackagesReportFromContent(c.repoContent, c.SDK().Root())
	if err != nil {
		return err
	}

	// write changelog
	if err := autorest.NewWriterFromFile(common.ChangelogPath(c.SDK().Root())).WithVersion(version.NewVersion).Write(r); err != nil {
		return err
	}
	// write version
	if err := repo.ModifyVersionFile(c.SDK().Root(), version.LatestVersion, version.NewVersion); err != nil {
		return err
	}
	// add commit
	if err := repo.AddCommit(c.SDK(), version.NewVersion); err != nil {
		return err
	}

	if len(errors) > 0 {
		return fmt.Errorf("release completes, but generation has %d error(s)", len(errors))
	}

	// print the release results
	for _, line := range GetReleaseResult(includedPackages) {
		fmt.Println(line)
	}
	fmt.Println("")

	return nil
}

func (c *commandContext) getPackagesToRelease(results []GenerateResult) (include []GenerateResult, exclude []GenerateResult) {
	for _, r := range results {
		if autorest.CanIncludeInMinor(r.Package) {
			include = append(include, r)
		} else {
			log.Printf("Package '%s' contains breaking changes and is not a preview package", r.Package.PackageName)
			exclude = append(exclude, r)
		}
	}
	return
}

func (c *commandContext) readOptions() error {
	optionFile, err := os.Open(filepath.Join(c.SDK().Root(), c.flags.OptionsPath))
	if err != nil {
		return fmt.Errorf("failed to open option files: %+v", err)
	}
	rawOptions, err := model.NewRawOptionsFrom(optionFile)
	if err != nil {
		return fmt.Errorf("failed to parse options: %+v", err)
	}
	c.options, err = rawOptions.Parse(c.SDK().Root())
	if err != nil {
		return err
	}
	return nil
}

func (c *commandContext) generate(requests *config.Config) ([]GenerateResult, map[string]error) {
	errorMap := make(map[string]error)
	var results []GenerateResult
	for readme, infoMap := range requests.Track1Requests {
		c.CheckExternalChanges()

		log.Printf("Executing autorest task for readme '%s'", readme)
		for tag, infoList := range infoMap {
			generateCtx := generateContext{
				sdkRepo:            c.SDK(),
				specRepo:           c.Spec(),
				skipProfiles:       c.flags.SkipProfile,
				readme:             readme,
				specLastCommitHash: c.Spec().LastHead().Hash().String(),
				defaultOptions:     c.options,
				repoContent:        c.repoContent,
				additionalOptions:  c.additionalOptions,
			}
			resultsForTag, err := generateCtx.generate(tag, infoList)
			if err != nil {
				log.Printf("Generation for readme '%s' failed: %+v", readme, err)
				key := fmt.Sprintf("%s/%s", readme, tag)
				errorMap[key] = err
				continue
			}
			results = append(results, *resultsForTag)
		}
	}
	return results, errorMap
}

func (c *commandContext) getVersion(major bool) (*common.VersionInfo, error) {
	latestVersion, err := repo.GetLatestVersion(c.SDK())
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %+v", err)
	}
	if len(c.flags.VersionNumber) > 0 {
		return &common.VersionInfo{
			LatestVersion: "v" + latestVersion.String(),
			NewVersion:    c.flags.VersionNumber,
		}, nil
	}
	newVersion := latestVersion.IncMinor()
	if major {
		newVersion = latestVersion.IncMajor()
	}
	return &common.VersionInfo{
		LatestVersion: "v" + latestVersion.String(),
		NewVersion:    "v" + newVersion.String(),
	}, nil
}
