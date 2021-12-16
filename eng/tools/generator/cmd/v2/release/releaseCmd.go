// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package release

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	releaseBranchNamePattern = "release-%s-%s-%s-%v"
)

// Release command
func Command() *cobra.Command {
	releaseCmd := &cobra.Command{
		Use:   "release-v2 <azure-sdk-for-go directory/commitid> <azure-rest-api-specs directory/commitid> <rp-name> [namespaceName]",
		Short: "Generate a v2 release of azure-sdk-for-go",
		Long: `This command will generate a new v2 release for azure-sdk-for-go with given rp name and namespace name.

azure-sdk-for-go directory/commitid: the directory path of the azure-sdk-for-go with git control or just a commitid for remote repo
azure-rest-api-specs directory: the directory path of the azure-rest-api-specs with git control or just a commitid for remote repo
rp-name: name of resource provider to be released, same for the swagger folder name
namespaceName: name of namespace to be released, default value is arm+rp-name

`,
		Args: cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			rpName := args[2]
			namespaceName := "arm" + rpName
			if len(args) == 4 {
				namespaceName = args[3]
			}

			ctx := commandContext{
				rpName:        rpName,
				namespaceName: namespaceName,
				flags:         ParseFlags(cmd.Flags()),
			}
			return ctx.execute(args[0], args[1])
		},
	}

	BindFlags(releaseCmd.Flags())

	return releaseCmd
}

type Flags struct {
	VersionNumber       string
	SwaggerRepo         string
	PackageTitle        string
	SDKRepo             string
	SpecRPName          string
	ReleaseDate         string
	SkipCreateBranch    bool
	SkipGenerateExample bool
	PackageConfig       string
	GoVersion           string
}

func BindFlags(flagSet *pflag.FlagSet) {
	flagSet.String("version-number", "", "Specify the version number of this release")
	flagSet.String("package-title", "", "Specifies the title of this package")
	flagSet.String("sdk-repo", "https://github.com/Azure/azure-sdk-for-go", "Specifies the sdk repo URL for generation")
	flagSet.String("spec-repo", "https://github.com/Azure/azure-rest-api-specs", "Specifies the swagger repo URL for generation")
	flagSet.String("spec-rp-name", "", "Specifies the swagger spec RP name, default is RP name")
	flagSet.String("release-date", "", "Specifies the release date in changelog")
	flagSet.Bool("skip-create-branch", false, "Skip create release branch after generation")
	flagSet.Bool("skip-generate-example", false, "Skip generate example for SDK in the same time")
	flagSet.String("package-config", "", "Additional config for package")
	flagSet.String("go-version", "1.16", "Go version")
}

func ParseFlags(flagSet *pflag.FlagSet) Flags {
	return Flags{
		VersionNumber:       flags.GetString(flagSet, "version-number"),
		PackageTitle:        flags.GetString(flagSet, "package-title"),
		SDKRepo:             flags.GetString(flagSet, "sdk-repo"),
		SwaggerRepo:         flags.GetString(flagSet, "spec-repo"),
		SpecRPName:          flags.GetString(flagSet, "spec-rp-name"),
		ReleaseDate:         flags.GetString(flagSet, "release-date"),
		SkipCreateBranch:    flags.GetBool(flagSet, "skip-create-branch"),
		SkipGenerateExample: flags.GetBool(flagSet, "skip-generate-example"),
		PackageConfig:       flags.GetString(flagSet, "package-config"),
		GoVersion:           flags.GetString(flagSet, "go-version"),
	}
}

type commandContext struct {
	rpName        string
	namespaceName string
	flags         Flags
}

func (c *commandContext) execute(sdkRepoParam, specRepoParam string) error {
	sdkRepo, err := common.GetSDKRepo(sdkRepoParam, c.flags.SDKRepo)
	if err != nil {
		return err
	}

	specCommitHash, err := common.GetSpecCommit(specRepoParam)
	if err != nil {
		return err
	}

	log.Printf("Release generation for rp: %s, namespace: %s", c.rpName, c.namespaceName)
	generateCtx := common.GenerateContext{
		SDKPath:        sdkRepo.Root(),
		SDKRepo:        &sdkRepo,
		SpecCommitHash: specCommitHash,
		SpecRepoURL:    c.flags.SwaggerRepo,
	}

	if c.flags.SpecRPName == "" {
		c.flags.SpecRPName = c.rpName
	}
	result, err := generateCtx.GenerateForSingleRPNamespace(&common.GenerateParam{
		RPName:              c.rpName,
		NamespaceName:       c.namespaceName,
		NamespaceConfig:     c.flags.PackageConfig,
		SpecficPackageTitle: c.flags.PackageTitle,
		SpecficVersion:      c.flags.VersionNumber,
		SpecRPName:          c.flags.SpecRPName,
		ReleaseDate:         c.flags.ReleaseDate,
		SkipGenerateExample: c.flags.SkipGenerateExample,
		GoVersion:           c.flags.GoVersion,
	})
	if err != nil {
		return fmt.Errorf("failed to finish release generation process: %+v", err)
	}
	// print generation result
	log.Printf("Generation result: %s", result)

	if !c.flags.SkipCreateBranch {
		log.Printf("Create new branch for release")
		releaseBranchName := fmt.Sprintf(releaseBranchNamePattern, c.rpName, c.namespaceName, result.Version, time.Now().Unix())
		if err := sdkRepo.CreateReleaseBranch(releaseBranchName); err != nil {
			return fmt.Errorf("failed to create release branch: %+v", err)
		}

		log.Printf("Include the packages that is about to release in this release and do release commit...")
		// append a time in long to avoid collision of branch names
		if err := sdkRepo.AddReleaseCommit(c.rpName, c.namespaceName, generateCtx.SpecCommitHash, result.Version); err != nil {
			return fmt.Errorf("failed to add release package or do release commit: %+v", err)
		}
	}

	return nil
}
