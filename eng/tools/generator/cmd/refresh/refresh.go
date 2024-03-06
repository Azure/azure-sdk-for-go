// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package refresh

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/flags"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	sdkutils "github.com/Azure/azure-sdk-for-go/eng/tools/internal/utils"
	"github.com/ahmetb/go-linq/v3"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Command() *cobra.Command {
	refreshCmd := &cobra.Command{
		Use:   "refresh <azure-sdk-for-go directory> <azure-rest-api-specs directory> [config json file path]",
		Short: "Regenerate all the packages in azure-sdk-for-go",
		Long: `This command will regenerate the specified packages in azure-sdk-for-go using the autorest.go version
specified in the option, but using the same swagger as it is using now.
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
			// this command by design will be checking out from commit to commit in azure-rest-api-specs,
			// therefore we explicitly turn out the panic
			baseContext, err := repo.NewCommandContext(sdkPath, specPath, false)
			if err != nil {
				return err
			}
			configPath := ""
			if len(args) > 2 {
				configPath = args[2]
			}
			ctx := CommandContext{
				CommandContext: baseContext,
				configPath:     configPath,
				Flags:          ParseFlags(cmd.Flags()),
			}
			return ctx.execute()
		},
	}

	BindFlags(refreshCmd.Flags())

	return refreshCmd
}

type Flags struct {
	common.GlobalFlags
	SkipProfile bool
	All         bool
}

func BindFlags(flagSet *pflag.FlagSet) {
	flagSet.Bool("skip-profile", false, "Skip the profile regeneration.")
	flagSet.Bool("all", false, "Refresh all packages without a configuration.")
}

func ParseFlags(flagSet *pflag.FlagSet) Flags {
	return Flags{
		GlobalFlags: common.ParseGlobalFlags(flagSet),
		SkipProfile: flags.GetBool(flagSet, "skip-profile"),
		All:         flags.GetBool(flagSet, "all"),
	}
}

type CommandContext struct {
	repo.CommandContext
	Flags Flags

	RepoContent map[string]exports.Content

	configPath string
}

func (c *CommandContext) parseConfig() (*config.Config, error) {
	if c.Flags.All {
		return &config.Config{}, nil
	}
	return config.ParseConfig(c.configPath)
}

func (c *CommandContext) execute() error {
	log.Printf("Parsing the config...")
	cfg, err := c.parseConfig()
	if err != nil {
		return err
	}
	log.Printf("Configuration: %s", cfg.String())

	// get the repo content to be compared with
	log.Printf("Reading packages in azure-sdk-for-go...")
	c.RepoContent, err = c.SDK().ReportForCommit("")
	if err != nil {
		return fmt.Errorf("failed to get the initial status of the SDK repository: %+v", err)
	}

	ref, err := c.Spec().Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD ref of azure-rest-api-specs: %+v", err)
	}
	defer func() {
		if err := c.Spec().Checkout(&repo.CheckoutOptions{
			Branch: ref.Name(),
			Force:  true,
		}); err != nil {
			log.Printf("Error checking out azure-rest-api-specs to %v", ref)
		}
	}()

	// create a temporary branch to hold the generation result
	log.Printf("Creating temporary branch...")
	tempBranchName, err := c.CreateReleaseBranch("temp")
	if err != nil {
		return err
	}
	log.Printf("Temporary branch '%s' created", tempBranchName)

	if _, err := c.Refresh(&cfg.RefreshInfo); err != nil {
		return err
	}

	return nil
}

func (c *CommandContext) Refresh(refreshConfig *config.RefreshInfo) (*plumbing.Reference, error) {
	if refreshConfig == nil {
		return nil, nil
	}

	// append the additional options
	additionalOptions, err := refreshConfig.AdditionalOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to parse additional options: %+v", err)
	}

	log.Printf("Getting the packages to be refreshed...")
	infoMap, err := c.getPackagesToRefresh(refreshConfig.RelativePackages())
	if err != nil {
		return nil, fmt.Errorf("failed to get the packages to refresh: %+v", err)
	}

	log.Printf("Total %d packages pending refresh after categorization", infoMap.Count())

	log.Printf("Regenerating all the packages...")
	if err := c.generate(infoMap, additionalOptions); err != nil {
		return nil, fmt.Errorf("failed to generate: %+v", err)
	}

	// commit changes
	log.Printf("Committing generated files...")
	if err := c.commitGeneratedContent(); err != nil {
		return nil, fmt.Errorf("failed to commit generated content: %+v", err)
	}

	// regenerate profiles
	log.Printf("Regenerating profiles...")
	if err := c.regenerateProfiles(); err != nil {
		return nil, fmt.Errorf("failed to regenerate profiles: %+v", err)
	}

	// commit profiles
	log.Printf("Commiting profiles...")
	if err := c.commitProfiles(); err != nil {
		return nil, fmt.Errorf("failed to commit profiles: %+v", err)
	}

	return c.SDK().Head()
}

func (c *CommandContext) getPackagesToRefresh(packages []string) (GenerationMap, error) {
	m, err := autorest.CollectGenerationMetadata(common.ServicesPath(c.SDK().Root()))
	if err != nil {
		return nil, fmt.Errorf("cannot read the metadata map: %+v", err)
	}

	newInfoMap := make(map[string]autorest.GenerationMetadata)
	if len(packages) > 0 {
		log.Printf("Picking the following packages to refresh: \n%s", strings.Join(packages, "\n"))
		// only take the following packages. Note that the package path in packages should be relative to the SDK root
		for _, relativePath := range packages {
			fullPath := sdkutils.NormalizePath(filepath.Join(c.SDK().Root(), relativePath))
			if v, ok := m[fullPath]; ok {
				log.Printf("picking package '%s'", relativePath)
				newInfoMap[fullPath] = v
			} else {
				log.Printf("do not find package '%s', ignoring", relativePath)
			}
		}
	}

	if len(newInfoMap) > 0 {
		return NewGenerationMap(newInfoMap), nil
	}

	return NewGenerationMap(m), nil
}

func (c *CommandContext) generate(infoMap GenerationMap, additionalOptions []model.Option) error {
	var errResult error
	for commit, infoList := range infoMap {
		errorsOnCommit := c.generateOnCommit(commit, infoList, additionalOptions)
		if errorsOnCommit != nil {
			errResult = errors.Join(errResult, errorsOnCommit)
		}
	}

	return errResult
}

func (c *CommandContext) generateOnCommit(commit string, infoList []GenerationInfo, additionalOptions []model.Option) error {
	log.Printf("Regenerate on commit %s starts", commit)
	// first we checkout the spec repo to that commit
	log.Printf("Checking out to commit %s...", commit)
	if err := c.Spec().Checkout(&repo.CheckoutOptions{
		Hash: plumbing.NewHash(commit),
	}); err != nil {
		var messages []string
		linq.From(infoList).SelectT(func(item GenerationInfo) string {
			return item.String()
		}).ToSlice(&messages)
		log.Printf("Error in checking out to '%s' which contains the following packages: \n%s", commit, strings.Join(messages, "\n"))
		return error(fmt.Errorf("cannot checkout to commit '%s': %+v", commit, err))
	}
	var errResult error
	var errCount int
	for _, info := range infoList {
		log.Printf("start generation task (readme '%s' / tag '%s')", info.Readme, info.Tag)
		// build the options from the metadata
		options, err := c.buildOptions(info.GenerationMetadata)
		if err != nil {
			errResult = errors.Join(errResult, err)
			errCount++
			continue
		}
		options = options.MergeOptions(additionalOptions...)
		start := time.Now()
		generateCtx := generateContext{
			sdkRoot:        c.SDK().Root(),
			specRoot:       c.Spec().Root(),
			specCommitHash: commit,
			options:        options,
			repoContent:    c.RepoContent,
		}
		_, err = generateCtx.generate(info)
		if err != nil {
			log.Printf("fails in generation task (readme %s / tag %s): %+v", info.Readme, info.Tag, err)
			errResult = errors.Join(errResult, fmt.Errorf("generate on commit %s failed: %+v", commit, err))
			errCount++
			continue
		}
		log.Printf("done generation of generation task %v (%v)", info, time.Since(start))
	}
	log.Printf("Regenerate on commit %s finished with %d errors", commit, errCount)
	return errResult
}

func (c *CommandContext) buildOptions(metadata autorest.GenerationMetadata) (model.Options, error) {
	rawOptions := strings.Split(metadata.AdditionalProperties.AdditionalOptions, " ")
	additionalOptions, err := model.ParseOptions(rawOptions)
	if err != nil {
		return nil, err
	}
	// the raw options do not contain `go-sdk-folder` or `use`, add them
	options := additionalOptions.MergeOptions(
		model.NewKeyValueOption("go-sdk-folder", c.SDK().Root()),
		model.NewKeyValueOption("use", metadata.CodeGenVersion),
	)
	return options, nil
}

func (c *CommandContext) regenerateProfiles() error {
	return autorest.RegenerateProfiles(c.SDK().Root())
}

func (c *CommandContext) commitGeneratedContent() error {
	if err := c.SDK().Add("services"); err != nil {
		return fmt.Errorf("failed to add `services`: %+v", err)
	}

	message := "Regenerated packages from their original commit hash"
	if err := c.SDK().Commit(message); err != nil {
		if repo.IsNothingToCommit(err) {
			log.Printf("There is nothing to commit. Message: %s", message)
			return nil
		}
		return fmt.Errorf("failed to commit changes: %+v", err)
	}

	return nil
}

func (c *CommandContext) commitProfiles() error {
	if err := c.SDK().Add("profiles"); err != nil {
		return fmt.Errorf("failed to add `profiles`: %+v", err)
	}

	if err := c.SDK().Commit("Refresh profiles"); err != nil {
		if repo.IsNothingToCommit(err) {
			return nil
		}
		return fmt.Errorf("failed to commit changes: %+v", err)
	}

	return nil
}
