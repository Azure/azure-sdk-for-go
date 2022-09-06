// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package profile

import (
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GenerateSDKCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generate-profile-sdk <azure-sdk-for-go directory> <azure-rest-api-specs directory> <profile name>",
		Args: cobra.ExactArgs(3),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFlags(0) // remove the time stamp prefix
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := generateSDK(args[0], args[1], args[2], ParseSDKFlags(cmd.Flags())); err != nil {
				logError(err)
				return err
			}
			return nil
		},
		SilenceUsage: true, // this command is used for a pipeline, the usage should never show
	}

	BindSDKFlags(cmd.Flags())

	return cmd
}

func logError(err error) {
	for _, line := range strings.Split(err.Error(), "\n") {
		if l := strings.TrimSpace(line); l != "" {
			log.Printf("[ERROR] %s", l)
		}
	}
}

type SDKFlags struct {
	VersionNumber string
	GoVersion     string
}

func BindSDKFlags(flagSet *pflag.FlagSet) {
	flagSet.String("version-number", "", "Specify the version number of this release")
	flagSet.String("go-version", common.DefaultGoVersion, "Go version")
}

func ParseSDKFlags(flagSet *pflag.FlagSet) SDKFlags {
	return SDKFlags{
		VersionNumber: flags.GetString(flagSet, "version-number"),
		GoVersion:     flags.GetString(flagSet, "go-version"),
	}
}
