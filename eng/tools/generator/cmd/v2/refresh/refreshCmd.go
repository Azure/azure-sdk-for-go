// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package refresh

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
		Use:   "refresh-v2 <azure-sdk-for-go directory/commitid> <azure-rest-api-specs directory/commitid>",
		Short: "Regenerate all v2 release of azure-sdk-for-go",
		Long: `This command will regenerate all v2 release for azure-sdk-for-go.

azure-sdk-for-go directory/commitid: the directory path of the azure-sdk-for-go with git control or just a commitid for remote repo
azure-rest-api-specs directory: the directory path of the azure-rest-api-specs with git control or just a commitid for remote repo
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := commandContext{
				flags: ParseFlags(cmd.Flags()),
			}
			return ctx.execute(args[0], args[1])
		},
	}

	BindFlags(releaseCmd.Flags())

	return releaseCmd
}

type Flags struct {
	SwaggerRepo         string
	SDKRepo             string
	ReleaseDate         string
	SkipCreateBranch    bool
	SkipGenerateExample bool
	RPs                 string
	UpdateSpecVersion   bool
}

func BindFlags(flagSet *pflag.FlagSet) {
	flagSet.String("sdk-repo", "https://github.com/Azure/azure-sdk-for-go", "Specifies the sdk repo URL for generation")
	flagSet.String("spec-repo", "https://github.com/Azure/azure-rest-api-specs", "Specifies the swagger repo URL for generation")
	flagSet.String("release-date", "", "Specifies the release date in changelog")
	flagSet.Bool("skip-create-branch", false, "Skip create release branch after generation")
	flagSet.Bool("skip-generate-example", false, "Skip generate example for SDK in the same time")
	flagSet.String("rps", "", "Specify RP list to refresh, seperated by ','")
	flagSet.Bool("update-spec-version", true, "Whether to update the commit id, the default is true")
}

func ParseFlags(flagSet *pflag.FlagSet) Flags {
	return Flags{
		SDKRepo:             flags.GetString(flagSet, "sdk-repo"),
		SwaggerRepo:         flags.GetString(flagSet, "spec-repo"),
		ReleaseDate:         flags.GetString(flagSet, "release-date"),
		SkipCreateBranch:    flags.GetBool(flagSet, "skip-create-branch"),
		SkipGenerateExample: flags.GetBool(flagSet, "skip-generate-example"),
		RPs:                 flags.GetString(flagSet, "rps"),
		UpdateSpecVersion:   flags.GetBool(flagSet, "update-spec-version"),
	}
}

type commandContext struct {
	flags Flags
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

	generateCtx := common.GenerateContext{
		SDKPath:           sdkRepo.Root(),
		SDKRepo:           &sdkRepo,
		SpecCommitHash:    specCommitHash,
		SpecRepoURL:       c.flags.SwaggerRepo,
		UpdateSpecVersion: c.flags.UpdateSpecVersion,
	}

	if !c.flags.SkipCreateBranch {
		log.Printf("Create new branch for release")
		releaseBranchName := fmt.Sprintf(releaseBranchNamePattern, "refresh", "all", "package", time.Now().Unix())
		if err := sdkRepo.CreateReleaseBranch(releaseBranchName); err != nil {
			return fmt.Errorf("failed to create release branch: %+v", err)
		}
	}

	var rpNames []string
	if c.flags.RPs == "" {
		rps, err := os.ReadDir(filepath.Join(generateCtx.SDKPath, "sdk", "resourcemanager"))
		if err != nil {
			return fmt.Errorf("failed to get all rps: %+v", err)
		}
		for _, rp := range rps {
			rpNames = append(rpNames, rp.Name())
		}
	} else {
		rpNames = strings.Split(c.flags.RPs, ",")
	}

	for _, rpName := range rpNames {
		namespaces, err := os.ReadDir(filepath.Join(generateCtx.SDKPath, "sdk", "resourcemanager", rpName))
		if err != nil {
			continue
		}

		for _, namespace := range namespaces {
			log.Printf("Release generation for rp: %s, namespace: %s", rpName, namespace.Name())
			specRpName, err := common.GetSpecRpName(filepath.Join(generateCtx.SDKPath, "sdk", "resourcemanager", rpName, namespace.Name()))
			if err != nil {
				continue
			}
			result, err := generateCtx.GenerateForSingleRPNamespace(&common.GenerateParam{
				RPName:               rpName,
				NamespaceName:        namespace.Name(),
				SpecificPackageTitle: "",
				SpecRPName:           specRpName,
				ReleaseDate:          c.flags.ReleaseDate,
				SkipGenerateExample:  c.flags.SkipGenerateExample,
			})
			if err != nil {
				fmt.Printf("failed to finish release generation process: %+v", err)
				continue
			}
			// print generation result
			log.Printf("Generation result: %s", result)

			if !c.flags.SkipCreateBranch {
				log.Printf("Include the packages that is about to release in this release and do release commit...")
				// append a time in long to avoid collision of branch names
				if err := sdkRepo.AddReleaseCommit(rpName, namespace.Name(), generateCtx.SpecCommitHash, result.Version); err != nil {
					return fmt.Errorf("failed to add release package or do release commit: %+v", err)
				}
			}
		}
	}

	return nil
}
