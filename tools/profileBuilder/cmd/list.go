// +build go1.9

// Copyright © 2018 Microsoft Corporation
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
	"io/ioutil"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/tools/profileBuilder/model"
	"github.com/marstr/randname"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Creates a profile from a set of packages.",
	Long: `Reads a list of packages from stdin, where each line is treated as a Go package
identifier. These packages are then used to create a profile.

Often, the easiest way of invoking this command will be using a pipe operator
to specify the packages to include.

Example:
$> ../model/testdata/smallProfile.txt > profileBuilder list --name small_profile
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logWriter := ioutil.Discard
		if viper.GetBool("verbose") {
			logWriter = os.Stdout
		}

		outputLog := log.New(logWriter, "[STATUS] ", 0)
		errLog := log.New(os.Stderr, "[ERROR] ", 0)

		outputLog.Printf("Output-Location set to: %s", viper.GetString(outputLocationLongName))

		model.BuildProfile(
			&model.ListStrategy{Reader: os.Stdin},
			*nameToUse,
			*outputLocationName,
			outputLog,
			errLog)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	outputLocationName = listCmd.Flags().StringP(outputLocationLongName, outputLocationShortName, outputLocationDefault, outputLocationDescription)
	nameToUse = listCmd.Flags().StringP(nameLongName, nameShortName, nameDefault, nameDescription)

	viper.BindPFlags(listCmd.Flags())

	viper.SetDefault(nameLongName, randname.Generate())
}
