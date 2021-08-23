// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package release

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/tools/generator/flags"
	"github.com/Azure/azure-sdk-for-go/tools/generator/repo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Release command
func Command() *cobra.Command {
	releaseCmd := &cobra.Command{
		Use:   "release-v2 <azure-sdk-for-go directory> <azure-rest-api-specs directory> <rp-name> [namespaceName]",
		Short: "Generate a v2 release of azure-sdk-for-go",
		Long: `This command will generate a new v2 release for azure-sdk-for-go with given rp name and namespace name.

azure-sdk-for-go directory: the directory path of the azure-sdk-for-go with git control
azure-rest-api-specs directory: the directory path of the azure-rest-api-specs with git control
rp-name: name of resource provider to be released, same for the swagger folder name
namespaceName: name of namespace to be released, default value is arm+rp-name

`,
		Args: cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			sdkPath, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get the directory of azure-sdk-for-go: %v", err)
			}
			specPath, err := filepath.Abs(args[1])
			if err != nil {
				return fmt.Errorf("failed to get the directory of azure-rest-api-specs: %v", err)
			}
			rpName := args[2]
			namespaceName := "arm" + rpName
			if len(args) == 4 {
				namespaceName = args[3]
			}

			ctx := commandContext{
				rpName:        rpName,
				namespaceName: namespaceName,
				sdkPath:       sdkPath,
				specPath:      specPath,
				flags:         ParseFlags(cmd.Flags()),
			}
			return ctx.execute()
		},
	}

	BindFlags(releaseCmd.Flags())

	return releaseCmd
}

type Flags struct {
	VersionNumber string
	RepoURL       string
	PackageTitle  string
}

func BindFlags(flagSet *pflag.FlagSet) {
	flagSet.String("version-number", "", "Specify the version number of this release")
	flagSet.String("repo-url", "", "Specifies the swagger repo url for generation")
	flagSet.String("package-title", "", "Specifies the title of this package")
}

func ParseFlags(flagSet *pflag.FlagSet) Flags {
	return Flags{
		VersionNumber: flags.GetString(flagSet, "version-number"),
		RepoURL:       flags.GetString(flagSet, "repo-url"),
		PackageTitle:  flags.GetString(flagSet, "package-title"),
	}
}

type commandContext struct {
	rpName        string
	namespaceName string
	sdkPath       string
	sdkRepo       repo.SDKRepository
	specPath      string
	specRepo      repo.SpecRepository
	flags         Flags
}

func (c *commandContext) execute() error {
	var err error
	// create sdk and spec git repo ref
	c.sdkRepo, err = repo.OpenSDKRepository(c.sdkPath)
	if err != nil {
		return fmt.Errorf("failed to get sdk repo: %+v", err)
	}
	c.specRepo, err = repo.OpenSpecRepository(c.specPath)
	if err != nil {
		return fmt.Errorf("failed to get spec repo: %+v", err)
	}

	// get sdk and spec repo head
	sdkRef, err := c.sdkRepo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD ref of azure-sdk-for-go: %+v", err)
	}
	log.Printf("The release branch is based on HEAD ref '%s' (commit %s) of azure-sdk-for-go", sdkRef.Name(), sdkRef.Hash())

	specRef, err := c.specRepo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD ref of azure-rest-api-specs: %+v", err)
	}
	log.Printf("The new version is generated from HEAD ref '%s' (commit %s) of azure-rest-api-specs", specRef.Name(), specRef.Hash())

	log.Printf("Release generation for rp: %s, namespace: %s", c.rpName, c.namespaceName)
	generateCtx := common.GenerateContext{
		SdkPath:    c.sdkPath,
		SpecPath:   c.specPath,
		CommitHash: specRef.Hash().String(),
	}

	result, err := generateCtx.GenerateForSingleRpNamespace(c.rpName, c.namespaceName, c.flags.PackageTitle, c.flags.VersionNumber, c.flags.RepoURL)
	if err != nil {
		return fmt.Errorf("failed to finish release generation process: %+v", err)
	}
	// print generation result
	log.Printf("Generation result: %s", result)

	log.Printf("Create new branch for release")
	releaseBranchName := fmt.Sprintf(releaseBranchNamePattern, c.rpName, c.namespaceName, result.Version, time.Now().Unix())
	if err := c.sdkRepo.CreateReleaseBranch(releaseBranchName); err != nil {
		return fmt.Errorf("failed to create release branch: %+v", err)
	}

	log.Printf("Include the packages that is about to release in this release and do release commit...")
	// append a time in long to avoid collision of branch names
	if err := c.sdkRepo.AddReleaseCommit(c.rpName, c.namespaceName, generateCtx.CommitHash, result.Version); err != nil {
		return fmt.Errorf("failed to add release package or do release commit: %+v", err)
	}

	return nil
}

const (
	releaseBranchNamePattern = "release-%s-%s-%s-%v"
)
