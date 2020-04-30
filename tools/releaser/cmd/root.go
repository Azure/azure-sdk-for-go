// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/Azure/azure-sdk-for-go/tools/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	tagPrefix   = "// tag: "
	versionFile = "version.go"
)

// Command returns the root command for this tool
func Command() *cobra.Command {
	// the root command
	root := &cobra.Command{
		Use:   "releaser <command>",
		Short: "Use this tool to make the module release and the legacy release of go SDK",
		Args:  cobra.ExactArgs(1),

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetLevel("warn")
			if verbose := viper.GetBool("verbose"); verbose {
				log.SetLevel("debug")
			}
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	// register flags
	pFlags := root.PersistentFlags()
	pFlags.Bool("verbose", false, "verbose output")
	if err := viper.BindPFlag("verbose", pFlags.Lookup("verbose")); err != nil {
		log.Fatalf("failed to bind flag: %+v", err)
	}

	root.AddCommand(taggerCommand())

	root.AddCommand(profilesCommand())

	return root
}
